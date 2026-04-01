package routes

import (
	"strings"

	"balance/app/modules"
	"balance/app/utils/base"

	"github.com/gin-gonic/gin"
)

func normalizeEnvironmentLabel(value string) string {
	normalized := strings.ToLower(strings.TrimSpace(value))
	switch normalized {
	case "prod", "production":
		return "Production"
	case "stage", "staging":
		return "Staging"
	default:
		return "Development"
	}
}

type systemManifestResponse struct {
	Version         string `json:"version"`
	EncryptedStatus string `json:"encrypted_status"`
	Environment     string `json:"environment"`
}

func systemManifestHandler(mod *modules.Modules) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		environmentLabel := normalizeEnvironmentLabel(mod.Conf.Svc.Environment())
		encryptedStatus := strings.TrimSpace(mod.Conf.Svc.Config().EncryptedStatus)
		if encryptedStatus == "" {
			encryptedStatus = "AES-256 SECURE"
		}

		_ = base.Success(ctx, &systemManifestResponse{
			Version:         mod.Conf.Svc.Version(),
			EncryptedStatus: encryptedStatus,
			Environment:     environmentLabel,
		}, "system-manifest")
	}
}
