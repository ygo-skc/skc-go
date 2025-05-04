package util

import (
	"fmt"
	"log"
	"os"
)

// combine certificate.crt and ca_bundle.crt to avoid issues with verifying the legitimacy of certificate.crt
func CombineCerts(certDir string) {
	if certData, err := os.ReadFile(fmt.Sprintf("%s/certificate.crt", certDir)); err != nil {
		log.Fatalf("Failed to read certificate file: %s", err)
	} else if caBundleData, err := os.ReadFile(fmt.Sprintf("%s/ca_bundle.crt", certDir)); err != nil {
		log.Fatalf("Failed to read CA bundle file: %s", err)
	} else if err = os.WriteFile(fmt.Sprintf("%s/concatenated.crt", certDir), append(certData, caBundleData...), 0600); err != nil {
		log.Fatalf("Failed to write combined certificate file: %s", err)
	}
}
