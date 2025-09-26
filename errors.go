package mongodb

import "github.com/goravel/framework/errors"

var (
	FailedToGenerateURI = errors.New("failed to generate MongoDB URI, please check the database configuration")
	ConfigNotFound      = errors.New("not found database configuration")
	ConnectionFailed    = errors.New("failed to connect to MongoDB")
	DatabaseNotFound    = errors.New("database name not specified")
)
