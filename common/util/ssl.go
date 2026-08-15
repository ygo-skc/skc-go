package util

import (
	"log/slog"
	"os"
	"path/filepath"
)

// combine certificate.crt and ca_bundle.crt to avoid issues with verifying the legitimacy of certificate.crt
func CombineCerts(certDir string) {
	if certData, err := os.ReadFile(filepath.Join(certDir, "certificate.crt")); err != nil {
		slog.Error("Failed to read certificate file", slog.String("cert_dir", certDir), slog.Any("err", err))
		os.Exit(1)
	} else if caBundleData, err := os.ReadFile(filepath.Join(certDir, "ca_bundle.crt")); err != nil {
		slog.Error("Failed to read CA bundle file", slog.String("cert_dir", certDir), slog.Any("err", err))
		os.Exit(1)
	} else if err = os.WriteFile(filepath.Join(certDir, "concatenated.crt"), append(certData, caBundleData...), 0600); err != nil {
		slog.Error("Failed to write combined certificate file", slog.String("cert_dir", certDir), slog.Any("err", err))
		os.Exit(1)
	}
}
