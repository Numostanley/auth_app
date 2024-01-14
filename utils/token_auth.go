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

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims(payload))
	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		return "", fmt.Errorf("error generating JWT token: %v", err)
	}

	user.LastLogin = &now
	db := db.Database.DB
	db.Save(&user)

	return tokenString, nil
}

func (tAuth *TokenAuthentication) DecodeToken(tokenString string, verifySignature bool) (map[string]interface{}, error) {
	oauthConfig := OauthConfig{}
	oauthConfig.Initialize()

	var err error
	var t *jwt.Token
	var claims jwt.MapClaims
	var keyFunc jwt.Keyfunc

	if verifySignature {
		keyFunc = func(token *jwt.Token) (interface{}, error) {
			return oauthConfig.PublicKey, nil
		}
	}

	t, err = jwt.Parse(tokenString, keyFunc)

	if err != nil {
		return nil, fmt.Errorf("error decoding JWT token: %v", err)
	}

	if t != nil {
		claims, ok := t.Claims.(jwt.MapClaims)
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

	expiryTime, ok := decodedToken["expiry_time"].(float64)
	if !ok {
		return nil, nil, fmt.Errorf("invalid expiry_time in JWT payload")
	}

	expiryTimeInUnix := int64(expiryTime)
	expiryDateTime := time.Unix(expiryTimeInUnix, 0)

	if time.Now().After(expiryDateTime) {
		return nil, nil, fmt.Errorf("token is expired")
	}

	issuer, ok := decodedToken["issuer"].(string)
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

	user, err := GetUserByID(userID)
	if user == nil {
		return nil, nil, fmt.Errorf("invalid user associated with the token %s", err)
	}

	return decodedToken, user, nil
}
