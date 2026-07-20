package util

import (
	"fmt"
	"log/slog"
	"os"
)

// combine certificate.crt and ca_bundle.crt to avoid issues with verifying the legitimacy of certificate.crt
func CombineCerts(certDir string) {
	if certData, err := os.ReadFile(fmt.Sprintf("%s/certificate.crt", certDir)); err != nil {
		slog.Error("Failed to read certificate file", slog.String("cert_dir", certDir), slog.Any("err", err))
		os.Exit(1)
	} else if caBundleData, err := os.ReadFile(fmt.Sprintf("%s/ca_bundle.crt", certDir)); err != nil {
		slog.Error("Failed to read CA bundle file", slog.String("cert_dir", certDir), slog.Any("err", err))
		os.Exit(1)
	} else if err = os.WriteFile(fmt.Sprintf("%s/concatenated.crt", certDir), append(certData, caBundleData...), 0600); err != nil {
		slog.Error("Failed to write combined certificate file", slog.String("cert_dir", certDir), slog.Any("err", err))
		os.Exit(1)
	}
}
