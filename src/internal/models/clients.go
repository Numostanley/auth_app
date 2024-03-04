package models

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Client struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
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

	client.ID = uuid.New()
	result := db.Create(client)

	if result.Error != nil {
		return fmt.Errorf("error creating client: %v", result.Error)
	}

	fmt.Printf("Client created: %v \n", client.ClientID)

	return nil
}

func GetClientByClientID(clientID string, db *gorm.DB) (*Client, error) {
	client := Client{ClientID: clientID}
	fetchedClient := db.Where("client_id = ?", clientID).First(&client)

	if fetchedClient.Error != nil {
		return nil, fmt.Errorf("error returning client %s", fetchedClient.Error)
	}
	return &client, nil
}

var AppClients = struct {
	AdminAppClient  string
	MobileAppClient string
}{
	AdminAppClient:  "AdminAppClient",
	MobileAppClient: "MobileAppClient",
}
