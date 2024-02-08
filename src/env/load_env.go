package env

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type GetEnv struct {
	PortString             string
	PostgresHost           string
	PostgresUser           string
	PostgresDB             string
	PostgresPassword       string
	PostgresPort           string
	SigningAlgorithm       string
	PrivateKeyPath         string
	Issuer                 string
	PublicKeyPath          string
	TokenExpiryTime        string
	RefreshTokenExpiryTime string
	AllowedHosts           string
	EmailUser              string
	EmailPassword          string
	EmailHost              string
}

func (env *GetEnv) LoadEnv() {
	err := godotenv.Load("env/dev/.env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}
	env.getPortSring()
	env.getPostgresHost()
	env.getPostgresUser()
	env.getPostgresDB()
	env.getPostgresPassword()
	env.getPostgresPort()
	env.getAllowedHosts()
	env.getSigningAlgorithm()
	env.getPrivateKeyPath()
	env.getIssuer()
	env.getPublicKeyPath()
	env.getTokenExpiryTime()
	env.getRefreshTokenExpiryTime()
	env.getEmailUser()
	env.getEmailPassword()
	env.getEmailHost()
}

func (env *GetEnv) getPortSring() {
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT not found in the environment")
	}
	env.PortString = string(portString)
}

func (env *GetEnv) getPostgresHost() {
	postgresHost := os.Getenv("PG_HOST")
	if postgresHost == "" {
		log.Fatal("PG_HOST not found in the environment")
	}
	env.PostgresHost = string(postgresHost)
}

func (env *GetEnv) getPostgresUser() {
	postgresUser := os.Getenv("PG_USER")
	if postgresUser == "" {
		log.Fatal("PG_USER not found in the environment")
	}
	env.PostgresUser = string(postgresUser)
}

func (env *GetEnv) getPostgresDB() {
	postgresDB := os.Getenv("PG_DATABASE")
	if postgresDB == "" {
		log.Fatal("PG_DATABASE not found in the environment")
	}
	env.PostgresDB = string(postgresDB)
}

func (env *GetEnv) getPostgresPassword() {
	postgresPassword := os.Getenv("PG_PASSWORD")
	if postgresPassword == "" {
		log.Fatal("PG_PASSWORD not found in the environment")
	}
	env.PostgresPassword = string(postgresPassword)
}

func (env *GetEnv) getPostgresPort() {
	postgresPort := os.Getenv("PG_PORT")
	if postgresPort == "" {
		log.Fatal("PG_PORT not found in the environment")
	}
	env.PostgresPort = string(postgresPort)
}

func (env *GetEnv) getAllowedHosts() {
	allowedHosts := os.Getenv("ALLOWED_HOSTS")
	if allowedHosts == "" {
		log.Fatal("ALLOWED_HOSTS not found in the environment")
	}
	env.AllowedHosts = string(allowedHosts)
}

func (env *GetEnv) getSigningAlgorithm() {
	signingAlgorithm := os.Getenv("SIGNING_ALGORITHM")
	if signingAlgorithm == "" {
		log.Fatal("SIGNING_ALGORITHM not found in the environment")
	}
	env.SigningAlgorithm = string(signingAlgorithm)
}

func (env *GetEnv) getPrivateKeyPath() {
	privateKeyPath := os.Getenv("PRIVATE_KEY_FILE_PATH")
	if privateKeyPath == "" {
		log.Fatal("PRIVATE_KEY_FILE_PATH not found in the environment")
	}
	env.PrivateKeyPath = string(privateKeyPath)
}

func (env *GetEnv) getIssuer() {
	issuer := os.Getenv("ISSUER")
	if issuer == "" {
		log.Fatal("ISSUER not found in the environment")
	}
	env.Issuer = string(issuer)
}

func (env *GetEnv) getPublicKeyPath() {
	publicKeyPath := os.Getenv("PUBLIC_KEY_FILE_PATH")
	if publicKeyPath == "" {
		log.Fatal("PUBLIC_KEY_FILE_PATH not found in the environment")
	}
	env.PublicKeyPath = string(publicKeyPath)
}

func (env *GetEnv) getTokenExpiryTime() {
	tokenExpiryTime := os.Getenv("TOKEN_EXPIRY_TIME")
	if tokenExpiryTime == "" {
		log.Fatal("TOKEN_EXPIRY_TIME not found in the environment")
	}
	env.TokenExpiryTime = string(tokenExpiryTime)
}

func (env *GetEnv) getRefreshTokenExpiryTime() {
	refreshTokenExpiryTime := os.Getenv("REFRESH_TOKEN_EXPIRY_TIME")
	if refreshTokenExpiryTime == "" {
		log.Fatal("REFRESH_TOKEN_EXPIRY_TIME not found in the environment")
	}
	env.RefreshTokenExpiryTime = string(refreshTokenExpiryTime)
}

func (env *GetEnv) getEmailUser() {
	emailUser := os.Getenv("EMAIL_USER")
	if emailUser == "" {
		log.Fatal("EMAIL_USER not found in the environment")
	}
	env.EmailUser = string(emailUser)
}

func (env *GetEnv) getEmailPassword() {
	emailPassword := os.Getenv("EMAIL_PASSWORD")
	if emailPassword == "" {
		log.Fatal("EMAIL_PASSWORD not found in the environment")
	}
	env.EmailPassword = string(emailPassword)
}

func (env *GetEnv) getEmailHost() {
	emailHost := os.Getenv("EMAIL_HOST")
	if emailHost == "" {
		log.Fatal("EMAIL_HOST not found in the environment")
	}
	env.EmailHost = string(emailHost)
}
