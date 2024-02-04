package utils

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type OauthConfig struct {
	Issuer                 string
	PrivateKey             string
	PublicKey              string
	SigningAlgorithm       string
	Audience               string
	TokenExpiryTime        int
	RefreshTokenExpiryTime int
}

func (core *OauthConfig) Initialize() {
	err := godotenv.Load()
	if err != nil {
		return
	}

	core.GetIssuer()
	_, err = core.GetPrivateKey()
	if err != nil {
		return
	}
	_, err = core.GetPublicKey()
	if err != nil {
		return
	}
	core.GetSigningAlgorithm()
	core.GetAudience()
	_, err = core.GetTokenExpiryTime()
	if err != nil {
		return
	}
	_, err = core.GetRefreshTokenExpiryTime()
	if err != nil {
		return
	}
}

func (core *OauthConfig) GetIssuer() string {
	issuer := os.Getenv("ISSUER")
	core.Issuer = issuer
	return issuer
}

func (core *OauthConfig) GetPrivateKey() (string, error) {
	privateKeyPath := os.Getenv("PRIVATE_KEY_FILE_PATH")
	privateKey, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return "", fmt.Errorf("error reading private key: %v", err)
	}
	core.PrivateKey = string(privateKey)
	return string(privateKey), nil
}

func (core *OauthConfig) GetPublicKey() (string, error) {
	publicKeyPath := os.Getenv("PUBLIC_KEY_FILE_PATH")
	publicKey, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return "", fmt.Errorf("error reading public key: %v", err)
	}
	core.PublicKey = string(publicKey)
	return string(publicKey), nil
}

func (core *OauthConfig) GetSigningAlgorithm() string {
	signingAlgorithm := os.Getenv("SIGNING_ALGORITHM")
	core.SigningAlgorithm = signingAlgorithm
	return signingAlgorithm
}

func (core *OauthConfig) GetAudience() string {
	audience := os.Getenv("ISSUER")
	core.Audience = audience
	return audience
}

func (core *OauthConfig) GetTokenExpiryTime() (int, error) {
	tokenExpiryTime := os.Getenv("TOKEN_EXPIRY_TIME")
	intValue, err := strconv.Atoi(tokenExpiryTime)
	if err != nil {
		return 0, fmt.Errorf("error converting string to int: %v", err)
	}
	core.TokenExpiryTime = intValue
	return intValue, nil
}

func (core *OauthConfig) GetRefreshTokenExpiryTime() (int, error) {
	refreshTokenExpiryTime := os.Getenv("REFRESH_TOKEN_EXPIRY_TIME")
	intValue, err := strconv.Atoi(refreshTokenExpiryTime)
	if err != nil {
		return 0, fmt.Errorf("error converting string to int: %v", err)
	}
	core.RefreshTokenExpiryTime = intValue
	return intValue, nil
}
