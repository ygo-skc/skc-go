# ygo-service — System Design

gRPC service for the SKC (Salty Card Collection) Yu-Gi-Oh! platform. It exposes
read-only card, product, restriction, and score data backed by a MySQL database.

- **Transport:** gRPC over TLS, port `9020`
- **Proto contracts:** defined in the shared `github.com/ygo-skc/skc-go/common/v3` module
  (`ygo` and `health` packages) and *implemented* here
- **Datastore:** MySQL (`go-sql-driver/mysql`), pooled connection (`skcDBConn`)
- **Registered services:** `HealthService`, `CardService`, `ProductService`,
  `CardRestrictionService`, `ScoreService` — **17 RPCs total**

---

## 1. High-level architecture

Every request follows the same three-layer path. The gRPC handlers in `api/` are thin
adapters; the repositories in `db/` own all SQL and domain logic.

```mermaid
flowchart LR
    Client([gRPC client])

    subgraph svc["ygo-service (:9020, TLS)"]
        direction TB
        GS["grpc.Server<br/>keepalive · TLS · msg-size caps"]

        subgraph api["api/ — handlers (adapters)"]
            H["healthServiceServer"]
            C["ygoCardServiceServer"]
            P["ygoProductServiceServer"]
            R["ygoCardRestrictionServiceServer"]
            S["ygoScoreServiceServer"]
        end

        subgraph db["db/ — repositories (domain)"]
            CR["CardRepository"]
            PR["ProductRepository"]
            RR["CardRestrictionRepository"]
            SR["ScoreRepository"]
        end
    end

    DB[("MySQL<br/>skcDBConn pool")]

    Client -->|"TLS + protobuf"| GS
    GS --> H & C & P & R & S
    C --> CR
    P --> PR
    R --> RR
    S --> SR
    CR --> DB
    PR --> DB
    RR --> DB
    SR --> DB
```

### Registration & lifecycle (`main.go` → `api/server.go`)

```mermaid
sequenceDiagram
    participant M as main()
    participant DB as db.EstablishDBConn
    participant SVC as api.RunService
    participant G as grpc.Server

    M->>DB: open pool, PingContext (5s)
    Note over DB: fails fast → os.Exit(1)
    M->>SVC: go RunService()
    SVC->>SVC: CombineCerts + load TLS creds
    SVC->>G: NewServer(creds, keepalive, buffer/msg caps)
    SVC->>G: Register Health/Card/Product/Restriction/Score
    SVC->>G: Serve(listener :9020)
    M->>M: select{} (block forever)
```

---

## 2. Cross-cutting patterns

These repeat across nearly every RPC, so they are described once here and referenced
by the per-RPC diagrams below.

| Concern | Where | Behavior |
|---|---|---|
| **Named logger** | `util.NewLogger(ctx, "<label>")` | Each handler derives a request-scoped logger + context; repos retrieve it via `util.RetrieveLogger(ctx)`. |
| **Error mapping** | `db/sql_helpers.go` | `sql.ErrNoRows` → `codes.NotFound`; any other query/scan error → `codes.Internal` ("Error occurred while querying DB"). |
| **Variadic `IN (...)`** | `buildVariableQuerySubjects` + `variablePlaceholders` | Builds `?, ?, ?` placeholders and args for multi-key lookups; empty input short-circuits to an empty result (no DB hit). |
| **Missing-key reporting** | `model.FindMissingKeys` | Batch lookups return `UnknownResources` — the subset of requested IDs/names not found — instead of erroring. |
| **Card row → proto** | `model.NewYGOCardProtoBuilder(...)` | Shared builder maps 8 card columns into a `ygo.Card`. |

**Standard batch-lookup flow** (referenced as *"batch pattern"* below):

```mermaid
sequenceDiagram
    participant H as Handler
    participant Repo
    participant DB as MySQL
    H->>Repo: method(ctx, keys)
    alt keys empty
        Repo-->>H: empty result (no query)
    else
        Repo->>Repo: build placeholders + args
        Repo->>DB: SELECT ... WHERE col IN (?, ?, …)
        DB-->>Repo: rows
        Repo->>Repo: scan → proto map, FindMissingKeys
        Repo-->>H: {results, UnknownResources}
    end
```

---

## 3. HealthService

### 3.1 `APIStatus`
`api/status_handler.go` — liveness/version probe. No datastore access.

- **Request/Response:** `health.APIStatusRequest` → `health.APIStatusResponse{Version: "3.0.0"}`

```mermaid
sequenceDiagram
    participant Client
    participant H as healthServiceServer
    Client->>H: APIStatus(APIStatusRequest)
    H->>H: NewLogger("Status") · log
    H-->>Client: APIStatusResponse{Version:"3.0.0"}
```

---

## 4. CardService — `api/card_service_handler.go` → `db/card_repo.go`

Nine RPCs over the `card_info` / `card_colors` tables. All share the `cardAttributes`
column set and the `parseCardRows` scanner.

### 4.1 `GetCardColors`
List of every card color and its numeric id.

- **Types:** `GetCardColorsRequest` → `GetCardColorsResponse{Values: map[color]id}`
- **SQL:** `SELECT color_id, card_color FROM card_colors ORDER BY color_id`

```mermaid
sequenceDiagram
    participant Client
    participant H as CardService
    participant R as CardRepository
    participant DB as MySQL
    Client->>H: GetCardColors
    H->>R: GetCardColorIDs(ctx)
    R->>DB: SELECT color_id, card_color FROM card_colors
    DB-->>R: rows
    R-->>H: map[card_color]color_id
    H-->>Client: GetCardColorsResponse{Values}
```

### 4.2 `GetCardByID`
Single card by card number.

- **Types:** `GetCardByIDRequest{Subject.Id}` → `GetCardByIDResponse{Card}`
- **SQL:** `SELECT <cardAttributes> FROM card_info WHERE card_number = ?`
- **Rules:** `NotFound` logged specifically when the id doesn't exist.

```mermaid
sequenceDiagram
    participant Client
    participant H as CardService
    participant R as CardRepository
    participant DB as MySQL
    Client->>H: GetCardByID(Subject.Id)
    H->>R: GetCardByID(ctx, id)
    R->>DB: SELECT … FROM card_info WHERE card_number = ?
    alt found
        DB-->>R: row → ygo.Card
        R-->>H: Card
        H-->>Client: GetCardByIDResponse{Card}
    else no rows
        DB-->>R: ErrNoRows
        R-->>H: status NotFound
        H-->>Client: error NotFound
    end
```

### 4.3 `GetCardsByID`  *(batch pattern)*
Multiple cards keyed by id.

- **Types:** `GetCardsByIDRequest{Subjects.Ids}` → `GetCardsByIDResponse{Cards}`
- **SQL:** `SELECT <cardAttributes> FROM card_info WHERE card_number IN (?, …)`
- **Result:** `Cards.CardInfo` keyed by id + `UnknownResources` for missing ids.

```mermaid
sequenceDiagram
    participant Client
    participant H as CardService
    participant R as CardRepository
    participant DB as MySQL
    Client->>H: GetCardsByID(Subjects.Ids)
    H->>R: GetCardsByIDs(ctx, ids)
    alt ids empty
        R-->>H: empty Cards
    else
        R->>DB: SELECT … WHERE card_number IN (?, …)
        DB-->>R: rows → map[id]Card
        R->>R: FindMissingKeys(ids)
        R-->>H: Cards{CardInfo, UnknownResources}
    end
    H-->>Client: GetCardsByIDResponse{Cards}
```

### 4.4 `GetCardsByName`  *(batch pattern)*
Same as 4.3 but keyed by exact card name.

- **Types:** `GetCardsByNameRequest{Subjects.Names}` → `GetCardsByNameResponse{Cards}`
- **SQL:** `SELECT <cardAttributes> FROM card_info WHERE card_name IN (?, …)`
- **Result:** map keyed by name (`collectWithMapUsingNameKey`) + `UnknownResources`.

```mermaid
sequenceDiagram
    participant Client
    participant H as CardService
    participant R as CardRepository
    participant DB as MySQL
    Client->>H: GetCardsByName(Subjects.Names)
    H->>R: GetCardsByNames(ctx, names)
    alt names empty
        R-->>H: empty Cards
    else
        R->>DB: SELECT … WHERE card_name IN (?, …)
        DB-->>R: rows → map[name]Card
        R->>R: FindMissingKeys(names)
        R-->>H: Cards{CardInfo, UnknownResources}
    end
    H-->>Client: GetCardsByNameResponse{Cards}
```

### 4.5 `GetCardsReferencingNameInEffect`
Full-text search: which cards mention the given name(s) in their effect text.

- **Types:** `GetCardsReferencingNameInEffectRequest{Subjects.Names}` → `…Response{Cards}` (list)
- **Rules:** each name is quote-stripped and wrapped as a `"phrase"` (`convertToFullText`),
  joined with spaces, then matched against a `FULLTEXT` index.
- **SQL:** `… WHERE MATCH(card_effect) AGAINST(? IN BOOLEAN MODE) ORDER BY color_id, card_name`

```mermaid
sequenceDiagram
    participant Client
    participant H as CardService
    participant R as CardRepository
    participant DB as MySQL
    Client->>H: GetCardsReferencingNameInEffect(names)
    alt names empty
        R-->>H: [] (no query)
    else
        H->>R: GetCardsReferencingNameInEffect(ctx, names)
        R->>R: convertToFullText per name → join " "
        R->>DB: MATCH(card_effect) AGAINST(? IN BOOLEAN MODE)
        DB-->>R: rows → [ ]Card
        R-->>H: []Card
    end
    H-->>Client: Response{Cards}
```

### 4.6 `GetArchetypalCardsUsingCardName`
Loose name match — cards whose **name** contains the archetype string.

- **Types:** `…Request{Subject.Name}` → `…Response{Cards}` (list)
- **SQL:** `… WHERE card_name LIKE BINARY ? ORDER BY card_name` with `%name%`.

```mermaid
sequenceDiagram
    participant Client
    participant H as CardService
    participant R as CardRepository
    participant DB as MySQL
    Client->>H: GetArchetypalCardsUsingCardName(Subject.Name)
    H->>R: GetArchetypalCardsUsingCardName(ctx, name)
    R->>DB: SELECT … WHERE card_name LIKE BINARY '%name%'
    DB-->>R: rows → []Card
    R-->>H: []Card
    H-->>Client: Response{Cards}
```

### 4.7 `GetExplicitArchetypalInclusions`
Cards **explicitly declared** part of an archetype by their effect text.

- **Types:** `…Request{Subject.Name}` → `…Response{Cards}` (list)
- **SQL (nested):** an inner full-text query (`+"This card is always treated as" +"<name>"`)
  is wrapped by an outer `REGEXP 'always treated as a.*"<name>".* card'` filter.

```mermaid
sequenceDiagram
    participant Client
    participant H as CardService
    participant R as CardRepository
    participant DB as MySQL
    Client->>H: GetExplicitArchetypalInclusions(name)
    H->>R: GetExplicitArchetypalInclusions(ctx, name)
    R->>R: build FULLTEXT subquery + REGEXP wrapper
    R->>DB: SELECT a.* FROM (fulltext subquery) a WHERE a.card_effect REGEXP …
    DB-->>R: rows → []Card
    R-->>H: []Card
    H-->>Client: Response{Cards}
```

### 4.8 `GetExplicitArchetypalExclusions`
Mirror of 4.7 for cards explicitly **excluded** from an archetype.

- **Types:** `…Request{Subject.Name}` → `…Response{Cards}` (list)
- **SQL (nested):** inner `+"This card is not treated as" +"<name>"` full-text,
  outer `REGEXP 'not treated as.*"<name>".* card'`.

```mermaid
sequenceDiagram
    participant Client
    participant H as CardService
    participant R as CardRepository
    participant DB as MySQL
    Client->>H: GetExplicitArchetypalExclusions(name)
    H->>R: GetExplicitArchetypalExclusions(ctx, name)
    R->>R: build FULLTEXT "not treated as" subquery + REGEXP wrapper
    R->>DB: SELECT a.* FROM (fulltext subquery) a WHERE a.card_effect REGEXP …
    DB-->>R: rows → []Card
    R-->>H: []Card
    H-->>Client: Response{Cards}
```

### 4.9 `GetRandomCard`
One random non-Token card, optionally excluding a blacklist.

- **Types:** `GetRandomCardRequest{Blacklist.BlackListedRefs}` → `…Response{Card}`
- **Rules:** query variant chosen by blacklist size; always `card_color != 'Token'`,
  `ORDER BY RAND() LIMIT 1`.

```mermaid
sequenceDiagram
    participant Client
    participant H as CardService
    participant R as CardRepository
    participant DB as MySQL
    Client->>H: GetRandomCard(Blacklist.BlackListedRefs)
    H->>R: GetRandomCard(ctx, blacklist)
    alt blacklist empty
        R->>DB: … WHERE card_color!='Token' ORDER BY RAND() LIMIT 1
    else
        R->>DB: … WHERE card_number NOT IN (?, …) AND card_color!='Token' ORDER BY RAND() LIMIT 1
    end
    DB-->>R: 1 row → Card
    R-->>H: Card
    H-->>Client: GetRandomCardResponse{Card}
```

---

## 5. ProductService — `api/product_service_handler.go` → `db/product_repo.go`

Four RPCs over `products` / `product_contents` / `product_info`.

### 5.1 `GetCardsByProductID`
Full product detail: metadata + every card in the set, with rarities and a rarity histogram.

- **Types:** `…Request{Subject.Id}` → `…Response{Product}`
- **Flow:** two queries — product header, then contents — merged in `parseRowsForProductItems`,
  which folds duplicate (card, position) rows into one item with multiple rarities and
  accumulates a `RarityDistribution`.
- **SQL:**
  - `SELECT … FROM products WHERE product_id = ?`
  - `SELECT <cardAttributes>, product_position, card_rarity FROM product_contents WHERE product_id = ? ORDER BY product_position`

```mermaid
sequenceDiagram
    participant Client
    participant H as ProductService
    participant R as ProductRepository
    participant DB as MySQL
    Client->>H: GetCardsByProductID(Subject.Id)
    H->>R: GetCardsByProductID(ctx, id)
    R->>DB: SELECT … FROM products WHERE product_id = ?
    alt not found
        DB-->>R: ErrNoRows → NotFound
        R-->>H: error
    else
        DB-->>R: product header
        R->>DB: SELECT … FROM product_contents WHERE product_id = ? ORDER BY product_position
        DB-->>R: rows
        R->>R: fold (card,position)→ProductItem[], build RarityDistribution
        R-->>H: Product{Items, TotalItems, RarityDistribution}
    end
    H-->>Client: GetCardsByProductIDResponse{Product}
```

### 5.2 `GetProductSummaryByID`
Single product summary. Implemented as a thin wrapper over the batch call (5.3).

- **Types:** `…Request{Subject.Id}` → `…Response{ProductSummary}`
- **Rules:** delegates to `GetProductsSummaryByID([id])`; missing id → `codes.NotFound`.

```mermaid
sequenceDiagram
    participant Client
    participant H as ProductService
    participant R as ProductRepository
    participant DB as MySQL
    Client->>H: GetProductSummaryByID(Subject.Id)
    H->>R: GetProductSummaryByID(ctx, id)
    R->>R: GetProductsSummaryByID(ctx, [id])
    R->>DB: SELECT … FROM product_info WHERE product_id IN (?)
    DB-->>R: 0..1 row
    alt present
        R-->>H: ProductSummary
        H-->>Client: Response{ProductSummary}
    else absent
        R-->>H: status NotFound
        H-->>Client: error NotFound
    end
```

### 5.3 `GetProductsSummaryByID`  *(batch pattern)*
Multiple product summaries by id.

- **Types:** `…Request{Subjects.Ids}` → `…Response{Products}`
- **SQL:** `SELECT product_id, …, product_content_total FROM product_info WHERE product_id IN (?, …)`
- **Result:** `Products.Products` map + `UnknownResources`.

```mermaid
sequenceDiagram
    participant Client
    participant H as ProductService
    participant R as ProductRepository
    participant DB as MySQL
    Client->>H: GetProductsSummaryByID(Subjects.Ids)
    H->>R: GetProductsSummaryByID(ctx, ids)
    alt ids empty
        R-->>H: empty Products
    else
        R->>DB: SELECT … FROM product_info WHERE product_id IN (?, …)
        DB-->>R: rows → map[id]ProductSummary
        R->>R: FindMissingKeys(ids)
        R-->>H: Products{Products, UnknownResources}
    end
    H-->>Client: GetProductsSummaryByIDResponse{Products}
```

### 5.4 `GetProductsReleasedSameDay`
"On this day" — products sharing a given month/day (any year).

- **Types:** `…Request{Date string}` → `…Response{Products}` (list)
- **Rules:** handler parses `Date` as `2006-01-02`; bad format → `codes.InvalidArgument`.
  Only month/day are matched (`DATE_FORMAT(release_date,'%m-%d')`), newest year first.
- **SQL:** `… FROM product_info WHERE DATE_FORMAT(product_release_date,'%m-%d') = ? ORDER BY product_release_date DESC`

```mermaid
sequenceDiagram
    participant Client
    participant H as ProductService
    participant R as ProductRepository
    participant DB as MySQL
    Client->>H: GetProductsReleasedSameDay(Date)
    H->>H: time.Parse("2006-01-02", Date)
    alt parse fails
        H-->>Client: error InvalidArgument
    else
        H->>R: GetProductsReleasedSameDay(ctx, date)
        R->>DB: … WHERE DATE_FORMAT(release_date,'%m-%d')=? ORDER BY release_date DESC
        DB-->>R: rows → []ProductSummary
        R-->>H: []ProductSummary
        H-->>Client: Response{Products}
    end
```

---

## 6. CardRestrictionService — `api/card_restriction_service_handler.go` → `db/card_restriction_repo.go`

### 6.1 `GetEffectiveTimelineForFormat`
Given a format, return all restriction-list effective dates split into the currently
**active** date and any **future** (scheduled) dates.

- **Types:** `ygo.Format{Value}` → `ygo.EffectiveTimeline{AllDates, FutureDates, ActiveDate}`
- **Rules (in handler):**
  - Only `Genesys` is supported (case-insensitive); otherwise `codes.InvalidArgument`.
  - "Today" is computed in **America/Chicago** (`chicagoLocation`, loaded at init).
  - Dates are returned newest-first from SQL; the handler walks them, collecting dates
    `> today` as future, then takes the **first** date `<= today` as active and stops.
- **SQL:** `SELECT UNIQUE effective_date FROM card_scores WHERE format = ? ORDER BY effective_date DESC`

```mermaid
sequenceDiagram
    participant Client
    participant H as CardRestrictionService
    participant R as CardRestrictionRepository
    participant DB as MySQL
    Client->>H: GetEffectiveTimelineForFormat(Format.Value)
    alt format != "Genesys"
        H-->>Client: error InvalidArgument
    else
        H->>R: GetDatesForFormat(ctx, format)
        R->>DB: SELECT UNIQUE effective_date FROM card_scores WHERE format=? ORDER BY effective_date DESC
        DB-->>R: []date (desc)
        R-->>H: []date
        H->>H: today = now(Chicago)
        loop each date (newest→oldest)
            alt date > today
                H->>H: append to FutureDates
            else
                H->>H: ActiveDate = date; break
            end
        end
        H-->>Client: EffectiveTimeline{AllDates, FutureDates, ActiveDate}
    end
```

---

## 7. ScoreService — `api/score_service_handler.go` → `db/score_repo.go`

Three RPCs over `card_scores` (joined to `card_info` / `cards`). The two per-card
endpoints share a `parser` fold (in the handler) that reshapes flat score rows into a
per-format current/scheduled/history view.

### 7.1 `GetScoresByFormatAndDate`
Every scored card for a format at a specific effective date, with a sort option.

- **Types:** `ygo.RestrictedContentRequest{Format, EffectiveDate, SortOrder}` → `ygo.ScoresForFormatAndDate`
- **Rules:** repo chooses an `ORDER BY` from `SortOrder`
  (`score DESC, card_color, card_name` vs. default `card_color, card_name`).
  Handler returns `codes.NotFound` when the format/date pair yields zero entries.
- **SQL:** `… FROM card_scores cs FORCE INDEX(FORMAT) JOIN card_info ci … WHERE cs.format=? AND cs.effective_date=? ORDER BY <sort>`

```mermaid
sequenceDiagram
    participant Client
    participant H as ScoreService
    participant R as ScoreRepository
    participant DB as MySQL
    Client->>H: GetScoresByFormatAndDate(format, date, sortOrder)
    H->>R: GetScoresByFormatAndDate(ctx, format, date, sortOrder)
    R->>R: pick ORDER BY from sortOrder
    R->>DB: card_scores FORCE INDEX(FORMAT) JOIN card_info WHERE format=? AND effective_date=?
    DB-->>R: rows → []CardScoreEntry, count
    R-->>H: entries, count
    alt count == 0
        H-->>Client: error NotFound
    else
        H-->>Client: ScoresForFormatAndDate{Format, EffectiveDate, Entries, TotalEntries}
    end
```

### 7.2 `GetCardScoreByID`
Full score history for one card across all formats/versions.

- **Types:** `ygo.ResourceID{Id}` → `ygo.CardScore`
- **Rules:** query returns one row per (format, effective_date) version via a
  `DISTINCT` version set `LEFT JOIN`ed to the card's scores (missing → `COALESCE(...,0)`).
  The handler's `parser` fold builds, per format: `CurrentScoreByFormat` (latest date
  `< today`), `UniqueFormats`, `ScheduledChanges` (dates `> today`, as `format|date`),
  and the full `ScoreHistory`. Empty history → `codes.NotFound`. "Today" is Chicago-local.
- **SQL:** version set `SELECT DISTINCT format, effective_date FROM card_scores`
  `LEFT JOIN card_scores … AND scores.card_number = ? ORDER BY effective_date DESC`

```mermaid
sequenceDiagram
    participant Client
    participant H as ScoreService
    participant R as ScoreRepository
    participant DB as MySQL
    Client->>H: GetCardScoreByID(ResourceID.Id)
    H->>H: today = now(Chicago)
    H->>R: GetCardScoreByID(ctx, id, today, parser)
    R->>DB: version set LEFT JOIN card_scores (card_number=?) ORDER BY effective_date DESC
    DB-->>R: rows (format, date, score, id)
    loop each row
        R->>R: parser(score, entry, today)
        Note right of R: sets CurrentScoreByFormat / UniqueFormats /<br/>ScheduledChanges / ScoreHistory
    end
    R-->>H: CardScore
    alt ScoreHistory empty
        H-->>Client: error NotFound
    else
        H-->>Client: CardScore
    end
```

### 7.3 `GetCardScoresByIDs`  *(batch pattern)*
Score history for many cards at once.

- **Types:** `GetCardScoresByIDsRequest{Subjects.Ids}` → `…Response{Scores: CardScores}`
- **Rules:** `CROSS JOIN cards` builds the (card × version) matrix, `LEFT JOIN` fills scores
  (`COALESCE(...,0)`); the same `parser` fold runs per card id. Response wraps the per-id
  map plus `UnknownResources` (requested ids with no data).
- **SQL:** version set `CROSS JOIN cards LEFT JOIN card_scores … WHERE cards.card_number IN (?, …) ORDER BY effective_date DESC`

```mermaid
sequenceDiagram
    participant Client
    participant H as ScoreService
    participant R as ScoreRepository
    participant DB as MySQL
    Client->>H: GetCardScoresByIDs(Subjects.Ids)
    H->>H: today = now(Chicago)
    H->>R: GetCardScoresByIDs(ctx, ids, today, parser)
    alt ids empty
        R-->>H: empty map
    else
        R->>DB: version set CROSS JOIN cards LEFT JOIN card_scores WHERE cards.card_number IN (?, …)
        DB-->>R: rows
        loop each row
            R->>R: parser(scoresByID[cardID], entry, today)
        end
        R-->>H: map[id]CardScore
    end
    H->>H: FindMissingKeys(ids) → UnknownResources
    H-->>Client: Response{CardScores{CardInfo, UnknownResources}}
```

---

## 8. RPC index

| # | Service | RPC | Handler → Repo | Primary table(s) |
|---|---|---|---|---|
| 1 | Health | `APIStatus` | — | (none) |
| 2 | Card | `GetCardColors` | `GetCardColorIDs` | `card_colors` |
| 3 | Card | `GetCardByID` | `GetCardByID` | `card_info` |
| 4 | Card | `GetCardsByID` | `GetCardsByIDs` | `card_info` |
| 5 | Card | `GetCardsByName` | `GetCardsByNames` | `card_info` |
| 6 | Card | `GetCardsReferencingNameInEffect` | `GetCardsReferencingNameInEffect` | `card_info` (FULLTEXT) |
| 7 | Card | `GetArchetypalCardsUsingCardName` | `GetArchetypalCardsUsingCardName` | `card_info` |
| 8 | Card | `GetExplicitArchetypalInclusions` | `GetExplicitArchetypalInclusions` | `card_info` (FULLTEXT+REGEXP) |
| 9 | Card | `GetExplicitArchetypalExclusions` | `GetExplicitArchetypalExclusions` | `card_info` (FULLTEXT+REGEXP) |
| 10 | Card | `GetRandomCard` | `GetRandomCard` | `card_info` |
| 11 | Product | `GetCardsByProductID` | `GetCardsByProductID` | `products`, `product_contents` |
| 12 | Product | `GetProductSummaryByID` | `GetProductSummaryByID` | `product_info` |
| 13 | Product | `GetProductsSummaryByID` | `GetProductsSummaryByID` | `product_info` |
| 14 | Product | `GetProductsReleasedSameDay` | `GetProductsReleasedSameDay` | `product_info` |
| 15 | Restriction | `GetEffectiveTimelineForFormat` | `GetDatesForFormat` | `card_scores` |
| 16 | Score | `GetScoresByFormatAndDate` | `GetScoresByFormatAndDate` | `card_scores`, `card_info` |
| 17 | Score | `GetCardScoreByID` | `GetCardScoreByID` | `card_scores` |
| 18 | Score | `GetCardScoresByIDs` | `GetCardScoresByIDs` | `card_scores`, `cards` |

> Row count is 18 because `GetProductSummaryByID` and `GetProductsSummaryByID` are
> distinct RPCs sharing one repo query — **17 distinct RPCs** across 5 services.
