package ports

import (
	"github.com/73NN0/foe-hammer/internal/common/stdio"
	"github.com/73NN0/foe-hammer/internal/config/app"
	"github.com/73NN0/foe-hammer/internal/config/domain"
)

const eventTopic = "config.event"

// Payloads pour les commandes

type GetConfigByIDPayload struct {
	ID int `json:"id"`
}

type GetConfigByPathPayload struct {
	RootDir string `json:"root_dir"`
}

type CreateConfigPayload struct {
	Config domain.ProjectConfig `json:"config"`
}

type UpdateConfigPayload struct {
	Config domain.ProjectConfig `json:"config"`
}

type DeleteConfigPayload struct {
	ID int `json:"id"`
}

// NewConfigHandler crée un handler pour les commandes config.
func NewStdioConfigHandler(service *app.Service) stdio.MessageHandler {
	return func(msg stdio.Message, pub stdio.Publisher) error {
		var result stdio.HandlerResult

		switch msg.Type {
		case "ListConfigs":
			items, err := service.List()
			if err != nil {
				return err
			}
			result = stdio.Success("ConfigsListed", map[string]any{"configs": items})

		case "GetConfigByID":
			var p GetConfigByIDPayload
			if err := stdio.UnmarshalPayload(msg, &p); err != nil {
				return err
			}
			cfg, err := service.GetByID(p.ID)
			if err != nil {
				result = stdio.Fail("ConfigGetFailed", err, map[string]any{"id": p.ID})
			} else {
				result = stdio.Success("ConfigResolved", cfg)
			}

		case "GetConfigByPath":
			var p GetConfigByPathPayload
			if err := stdio.UnmarshalPayload(msg, &p); err != nil {
				return err
			}
			cfg, err := service.GetByPath(p.RootDir)
			if err != nil {
				result = stdio.Fail("ConfigGetFailed", err, map[string]any{"root_dir": p.RootDir})
			} else {
				result = stdio.Success("ConfigResolved", cfg)
			}

		case "CreateConfig":
			var p CreateConfigPayload
			if err := stdio.UnmarshalPayload(msg, &p); err != nil {
				return err
			}
			if err := service.Create(p.Config); err != nil {
				result = stdio.Fail("ConfigCreateFailed", err, map[string]any{"root_dir": p.Config.RootDir})
			} else {
				result = stdio.Success("ConfigCreated", p.Config)
			}

		case "UpdateConfig":
			var p UpdateConfigPayload
			if err := stdio.UnmarshalPayload(msg, &p); err != nil {
				return err
			}
			if err := service.Update(p.Config); err != nil {
				result = stdio.Fail("ConfigUpdateFailed", err, map[string]any{"id": p.Config.ID})
			} else {
				result = stdio.Success("ConfigUpdated", p.Config)
			}

		case "DeleteConfig":
			var p DeleteConfigPayload
			if err := stdio.UnmarshalPayload(msg, &p); err != nil {
				return err
			}
			if err := service.Delete(p.ID); err != nil {
				result = stdio.Fail("ConfigDeleteFailed", err, map[string]any{"id": p.ID})
			} else {
				result = stdio.Success("ConfigDeleted", map[string]any{"id": p.ID})
			}

		default:
			return nil
		}

		return result.Publish(msg, pub, eventTopic)
	}
}
