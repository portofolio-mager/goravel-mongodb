package mongodb

import (
	"fmt"

	"github.com/goravel/framework/contracts/config"
	"github.com/portofolio-mager/goravel-mongodb/contracts"
)

type Config struct {
	config     config.Config
	connection string
}

func NewConfig(config config.Config, connection string) *Config {
	return &Config{
		config:     config,
		connection: connection,
	}
}

func (r *Config) Config() config.Config {
	return r.config
}

func (r *Config) Connection() string {
	return r.connection
}

func (r *Config) Readers() []contracts.FullConfig {
	configs := r.config.Get(fmt.Sprintf("database.connections.%s.read", r.connection))
	if readConfigs, ok := configs.([]contracts.Config); ok {
		return r.fillDefault(readConfigs)
	}

	return nil
}

func (r *Config) Writers() []contracts.FullConfig {
	configs := r.config.Get(fmt.Sprintf("database.connections.%s.write", r.connection))
	if writeConfigs, ok := configs.([]contracts.Config); ok {
		return r.fillDefault(writeConfigs)
	}

	// Use default db configuration when write is empty
	return r.fillDefault([]contracts.Config{{}})
}

func (r *Config) fillDefault(configs []contracts.Config) []contracts.FullConfig {
	if len(configs) == 0 {
		return nil
	}

	var fullConfigs []contracts.FullConfig
	for _, config := range configs {
		fullConfig := contracts.FullConfig{
			Config: contracts.Config{
				URI:            config.URI,
				Database:       config.Database,
				Username:       config.Username,
				Password:       config.Password,
				AuthSource:     config.AuthSource,
				ReplicaSet:     config.ReplicaSet,
				TLS:            config.TLS,
				TLSCAFile:      config.TLSCAFile,
				TLSCertFile:    config.TLSCertFile,
				TLSKeyFile:     config.TLSKeyFile,
				MaxPoolSize:    config.MaxPoolSize,
				MinPoolSize:    config.MinPoolSize,
				ConnectTimeout: config.ConnectTimeout,
				ServerTimeout:  config.ServerTimeout,
				Options:        config.Options,
			},
			Connection: r.connection,
			Driver:     Name,
		}

		// If read or write is empty, use the default config - only fill required fields
		if fullConfig.URI == "" {
			fullConfig.URI = r.config.GetString(fmt.Sprintf("database.connections.%s.uri", r.connection))
		}
		if fullConfig.Database == "" {
			fullConfig.Database = r.config.GetString(fmt.Sprintf("database.connections.%s.database", r.connection))
		}

		fullConfigs = append(fullConfigs, fullConfig)
	}

	return fullConfigs
}
