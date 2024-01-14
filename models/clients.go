package models

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Client struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	ClientID     string    `json:"client_id" gorm:"unique"`
	ClientSecret string    `json:"client_secret"`
	Name         string    `json:"name"`
	ResponseType string    `json:"response_type"`
	Scope        string    `json:"scope"`
	GrantTypes   string    `json:"grant_types"`
	RedirectURIs string    `json:"redirect_uris"`
}

func (c *Client) ValidateSecret(secret string) bool {
	encodedSecret := encodeClientSecret(secret)
	return encodedSecret == c.ClientSecret
}

func (c *Client) ValidateScope(scope string) bool {
	return validateList(scope, c.Scope)
}

func (c *Client) ValidateGrantType(grantType string) bool {
	return validateList(grantType, c.GrantTypes)
}

func (c *Client) ValidateRedirectURI(redirectURI string) bool {
	return validateList(redirectURI, c.RedirectURIs)
}

func encodeClientSecret(secret string) string {
	return base64.StdEncoding.EncodeToString([]byte(secret))
}

func validateList(input string, list string) bool {
	inputItems := strings.Split(input, " ")
	listItems := strings.Split(list, " ")

	validItems := make(map[string]bool)

	for _, item := range listItems {
		validItems[item] = true
	}

	for _, inputItem := range inputItems {
		if !validItems[inputItem] {
			return false
		}
	}

	return true
}

func CreateClient(db *gorm.DB, client *Client) error {
	client.ClientSecret = encodeClientSecret(client.ClientSecret)

	result := db.Create(client)

	if result.Error != nil {
		return fmt.Errorf("error creating client: %v", result.Error)
	}

	fmt.Printf("Client created: %v \n", client.ClientID)

	return nil
}

var D8erAppClients = struct {
	SuperAdminUserClient string
	AppUserClients       string
}{
	SuperAdminUserClient: "D8ERSuperAdminClient",
	AppUserClients:       "D8ERAppClient",
}
