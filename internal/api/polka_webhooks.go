package api

import (
	"encoding/json"
	"internal/database"
	"log"
	"net/http"
	"strings"
)

type paramPolkaWebhooksData struct {
	UserId int `json:"user_id"`
}
type paramPolkaWebhooks struct {
	Event string                 `json:"event"`
	Data  paramPolkaWebhooksData `json:"data"`
}

func (cfg *ApiConfig) PostPolkaWebhooks(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := paramPolkaWebhooks{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	apiKeyString := strings.Replace(r.Header.Get("Authorization"), "ApiKey ", "", 1)
	if cfg.PolkaApiKey != apiKeyString {
		log.Printf("Polka webhooks received invalid API key: '%s'", apiKeyString)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if params.Event != "user.upgraded" {
		log.Printf("Polka webhooks received unhandled event: %s", params.Event)
		w.WriteHeader(http.StatusOK)
		return
	}

	err = cfg.DB.UpgradeUserById(params.Data.UserId, true)
	if err == database.ErrUserIdNotFound {
		log.Printf("Polka webhooks received invalid user ID: %d", params.Data.UserId)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("Polka webhooks error upgrading user:\n%v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
