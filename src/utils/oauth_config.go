package utils

import (
	"fmt"
	"os"
	"strconv"

	"github.com/Numostanley/d8er_app/env"
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
	core.GetIssuer()
	_, err := core.GetPrivateKey()
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
	enV := env.GetEnv{}
	enV.LoadEnv()
	core.Issuer = enV.Issuer
	return core.Issuer
}

func (core *OauthConfig) GetPrivateKey() (string, error) {
	enV := env.GetEnv{}
	enV.LoadEnv()
	privateKey, err := os.ReadFile(enV.PrivateKeyPath)
	if err != nil {
		return "", fmt.Errorf("error reading private key: %v", err)
	}
	core.PrivateKey = string(privateKey)
	return string(privateKey), nil
}

func (core *OauthConfig) GetPublicKey() (string, error) {
	enV := env.GetEnv{}
	enV.LoadEnv()
	publicKey, err := os.ReadFile(enV.PublicKeyPath)
	if err != nil {
		return "", fmt.Errorf("error reading public key: %v", err)
	}
	core.PublicKey = string(publicKey)
	return string(publicKey), nil
}

func (core *OauthConfig) GetSigningAlgorithm() string {
	enV := env.GetEnv{}
	enV.LoadEnv()
	core.SigningAlgorithm = enV.SigningAlgorithm
	return core.SigningAlgorithm
}

func (core *OauthConfig) GetAudience() string {
	enV := env.GetEnv{}
	enV.LoadEnv()
	core.Audience = enV.Issuer
	return core.Audience
}

func (core *OauthConfig) GetTokenExpiryTime() (int, error) {
	enV := env.GetEnv{}
	enV.LoadEnv()
	intValue, err := strconv.Atoi(enV.TokenExpiryTime)
	if err != nil {
		return 0, fmt.Errorf("error converting string to int: %v", err)
	}
	core.TokenExpiryTime = intValue
	return intValue, nil
}

func (core *OauthConfig) GetRefreshTokenExpiryTime() (int, error) {
	enV := env.GetEnv{}
	enV.LoadEnv()
	intValue, err := strconv.Atoi(enV.RefreshTokenExpiryTime)
	if err != nil {
		return 0, fmt.Errorf("error converting string to int: %v", err)
	}
	core.RefreshTokenExpiryTime = intValue
	return intValue, nil
}
