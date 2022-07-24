package env

import (
	"os"
)

var (
	// cdn endpoint for photo
	cdnEndpoint = os.Getenv("CDN_ENDPOINT")
	// user Jwt salt 
	userJwtSalt = os.Getenv("USER_JWT_SALT")
	// user code hash salt
	userCodeSalt = os.Getenv("USER_CODE_SALT")
	// group code hash salt
	groupCodeSalt = os.Getenv("GROUP_CODE_SALT")
	// database username
	databaseUsername = os.Getenv("DATABASE_USERNAME")
	// database password
	databasePassword = os.Getenv("DATABASE_PASSWORD")
	// database host
	databaseHost = os.Getenv("DATABASE_HOST")
	// database port
	databasePort = os.Getenv("DATABASE_PORT")
	// database name
	databaseName = os.Getenv("DATABASE_NAME")
)

// GetUserJwtSalt returns jwt for user info hash salt
func GetUserJwtSalt() string {
	return userJwtSalt
}

// GetUserCodeSalt returns user code hash salt
func GetUserCodeSalt() string {
	return userCodeSalt
}

// GetGroupCodeSalt returns group code hash salt
func GetGroupCodeSalt() string {
	return groupCodeSalt
}

// GetCdnEndpoint returns cdn endpoint for photo
func GetCdnEndpoint() string {
	return cdnEndpoint
}

// GetDatabaseUsername returns database username
func GetDatabaseUsername() string {
	return databaseUsername
}

// GetDatabasePassword returns database password
func GetDatabasePassword() string {
	return databasePassword
}

// GetDatabaseHost returns database host
func GetDatabaseHost() string {
	return databaseHost
}

// GetDatabasePort returns database port
func GetDatabasePort() string {
	return databasePort
}

// GetDatabaseName returns database name
func GetDatabaseName() string {
	return databaseName
}