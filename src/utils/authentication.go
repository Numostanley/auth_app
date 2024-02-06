package utils

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Numostanley/d8er_app/db"
	"github.com/Numostanley/d8er_app/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

type TokenAuthentication struct{}

func (tAuth *TokenAuthentication) GetOauthToken(scope string, user models.User) (string, error) {

	oauthConfig := OauthConfig{}
	oauthConfig.Initialize()

	key, err := oauthConfig.GetPrivateKey()
	if err != nil {
		return "", fmt.Errorf("error reading private key: %v", err)
	}
	aud := oauthConfig.Audience
	now := time.Now()

	tokenExpiryTime, err := oauthConfig.GetTokenExpiryTime()
	if err != nil {
		return "", fmt.Errorf("error reading tokenExpiryTime key: %v", err)
	}

	payload := map[string]interface{}{
		"iss":          oauthConfig.Issuer,
		"sub":          user.ID,
		"full_name":    user.GetFullName(),
		"email":        user.Email,
		"phone_number": user.PhoneNumber,
		"scope":        scope,
		"iat":          now,
		"aud":          aud,
		"exp":          now.Add(time.Second * time.Duration(tokenExpiryTime)),
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(key))
	if err != nil {
		return "", fmt.Errorf("error parsing RSA private key: %v", err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims(payload))
	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("error generating JWT token: %v", err)
	}

	user.LastLogin = &now
	database := db.Database.DB
	database.Save(&user)

	return tokenString, nil
}

func (tAuth *TokenAuthentication) DecodeToken(tokenString string, verifySignature bool) (map[string]interface{}, error) {
	oauthConfig := OauthConfig{}
	oauthConfig.Initialize()

	key, err := oauthConfig.GetPrivateKey()
	if err != nil {
		return map[string]interface{}{}, fmt.Errorf("error reading private key: %v", err)
	}

	var token *jwt.Token
	var claims jwt.MapClaims
	var keyFunc jwt.Keyfunc

	if verifySignature {
		privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(key))
		if err != nil {
			return nil, fmt.Errorf("error parsing RSA private key: %v", err)
		}

		keyFunc = func(token *jwt.Token) (interface{}, error) {
			return privateKey.Public(), nil
		}
	}
	token, err = jwt.Parse(tokenString, keyFunc)

	if err != nil {
		return nil, fmt.Errorf("error decoding JWT token: %v", err)
	}

	if token != nil {
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return nil, fmt.Errorf("error extracting claims from JWT token")
		}
		return claims, nil
	}
	return claims, nil
}

func (tAuth *TokenAuthentication) Authenticate(request *http.Request) (map[string]interface{}, *models.User, error) {
	authHeader := request.Header.Get("Authorization")
	auth := strings.Split(authHeader, " ")

	oauthConfig := OauthConfig{}
	oauthConfig.Initialize()

	if len(auth) != 2 || strings.ToLower(auth[0]) != "bearer" {
		return nil, nil, fmt.Errorf("invalid BEARER Authorization header")
	}

	tokenString := auth[1]

	decodedToken, err := tAuth.DecodeToken(tokenString, true)
	if err != nil {
		return nil, nil, err
	}

	expiryTime, ok := decodedToken["exp"]
	if !ok {
		return nil, nil, fmt.Errorf("invalid expiry_time in JWT payload")
	}

	expiryTimeString, ok := expiryTime.(string)
	if !ok {
		return nil, nil, fmt.Errorf("invalid expiry_time format in JWT payload")
	}

	expiryDateTime, err := time.Parse(time.RFC3339Nano, expiryTimeString)
	if err != nil {
		return nil, nil, fmt.Errorf("error parsing expiry_time: %v", err)
	}

	expiryTimeInUnix := expiryDateTime.Unix()

	if time.Now().After(time.Unix(expiryTimeInUnix, 0)) {
		return nil, nil, fmt.Errorf("token is expired")
	}

	issuer, ok := decodedToken["iss"].(string)
	if !ok {
		return nil, nil, fmt.Errorf("invalid issuer in JWT payload")
	}

	if strings.ToLower(issuer) != oauthConfig.Issuer {
		return nil, nil, fmt.Errorf("invalid JWT issuer")
	}

	sub, ok := decodedToken["sub"].(string)
	if !ok {
		return nil, nil, fmt.Errorf("invalid sub in JWT payload")
	}

	userID, err := uuid.Parse(sub)
	if err != nil {
		return nil, nil, fmt.Errorf("error parsing user ID: %v", err)
	}

	database := db.Database.DB

	user, err := models.GetUserByID(userID, database)
	if user == nil {
		return nil, nil, fmt.Errorf("invalid user associated with the token %s", err)
	}

	return decodedToken, user, nil
}

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

	database := db.Database.DB
	user, err := models.GetUserByEmail(email, database)

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
