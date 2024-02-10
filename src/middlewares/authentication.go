package middlewares

import (
	"fmt"
	"net/http"

	"github.com/Numostanley/auth_app/models"
	"github.com/Numostanley/auth_app/serializers"
	"github.com/Numostanley/auth_app/utils"
)

type authedHandler func(http.ResponseWriter, *http.Request, models.User)

func AuthenticationMiddleware(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := serializers.ResponseSerializer{
			Success: true,
			Message: "",
			Data:    struct{}{},
			Error:   "",
		}

		authenticator := utils.TokenAuthentication{}

		_, user, err := authenticator.Authenticate(r)
		if err != nil {
			data.Success = false
			data.Error = fmt.Sprintf("Authentication error: %v", err)
			utils.RespondWithError(w, http.StatusUnauthorized, data)
			return
		}

		handler(w, r, *user)
	}
}
