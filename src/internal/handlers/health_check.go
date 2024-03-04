package handlers

import (
	"net/http"

	"github.com/Numostanley/auth_app/internal/utils"
)

func HandlerReadiness(w http.ResponseWriter, _ *http.Request) {
	response := struct {
		Message string `json:"message"`
	}{
		Message: "Silence is Golden"}
	utils.RespondWithJSON(w, 200, response)
}
