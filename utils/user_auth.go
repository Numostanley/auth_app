package utils

import (
	"fmt"

	"github.com/Numostanley/d8er_app/models"
)

func PerformAuthentication(clientID, clientSecret, grantType, email, password, scope string) (*models.Client, *models.User, error) {
	client, err := GetClientByClientID(clientID)

	if err != nil {
		return nil, nil, fmt.Errorf("invalid_client %s", err)
	}

	if !client.ValidateSecret(clientSecret) {
		return nil, nil, fmt.Errorf("invalid_client_credentials")
	}

	if !client.ValidateGrantType(grantType) {
		return nil, nil, fmt.Errorf("invalid_grant_type")
	}

	if !client.ValidateScope(scope) {
		return nil, nil, fmt.Errorf("invalid_scope")
	}

	user, err := GetUserByEmail(email)

	if err != nil {
		return client, nil, fmt.Errorf("invalid_user %s", err)
	}

	// if !user.IsEmailVerified {
	// 	return client, nil, fmt.Errorf("invalid_user")
	// }

	if !user.ValidatePassword(password) {
		return client, nil, fmt.Errorf("invalid_user_credentials")
	}

	if !user.ValidateUserAgainstClientID(clientID) {
		return client, nil, fmt.Errorf("invalid_client_and_user")
	}

	// if user.HasActiveSession() {
	// 	return client, nil, fmt.Errorf("user_has_an_active_session")
	// }

	return client, user, nil
}
