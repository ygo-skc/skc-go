mkdir -p certs

aws secretsmanager get-secret-value --secret-id "/prod/skc/ygo-service/ssl" --region us-east-2 |
	jq -r '.SecretString' |
	jq -r "with_entries(select(.key | startswith(\"SSL\")))" >certs/base64-certs-json
jq -r ".SSL_PRIVATE_KEY" <certs/base64-certs-json | base64 -d >certs/private.key
jq -r ".SSL_CA_BUNDLE_CRT" <certs/base64-certs-json | base64 -d >certs/ca_bundle.crt
jq -r ".SSL_CERTIFICATE_CRT" <certs/base64-certs-json | base64 -d >certs/certificate.crt

rm certs/base64-certs-json

#############################################
createEnvFile() {
	local SKC_API_DB_INFO=$1
	local FILE_NAME=$2

	echo "$SKC_API_DB_INFO" | jq -r ". | to_entries|map(\"\(.key)=\\\"\(.value|tostring)\\\"\")|.[]" >"$FILE_NAME"
}

SKC_API_DB_INFO=$(aws secretsmanager get-secret-value --secret-id "/prod/skc/skc-api/db" --region us-east-2 |
	jq -r '.SecretString' |
	jq -c "{DB_USERNAME: .username, DB_PASSWORD: .password, DB_HOST: .host, DB_PORT: .port, DB_NAME: .dbname} | with_entries(.key |= \"SKC_\(.)\")")
createEnvFile "$SKC_API_DB_INFO" ".env_prod"

SKC_API_DB_INFO=$(aws secretsmanager get-secret-value --secret-id "/local/skc/skc-api/db" --region us-east-2 |
	jq -r '.SecretString' |
	jq -c "{DB_USERNAME: .username, DB_PASSWORD: .password, DB_HOST: .host, DB_PORT: .port, DB_NAME: .dbname} | with_entries(.key |= \"SKC_\(.)\")")
createEnvFile "$SKC_API_DB_INFO" ".env"

SKC_API_DB_INFO=$(aws secretsmanager get-secret-value --secret-id "/docker/local/skc/skc-api/db" --region us-east-2 |
	jq -r '.SecretString' |
	jq -c "{DB_USERNAME: .username, DB_PASSWORD: .password, DB_HOST: .host, DB_PORT: .port, DB_NAME: .dbname} | with_entries(.key |= \"SKC_\(.)\")")
createEnvFile "$SKC_API_DB_INFO" ".env_docker_local"
