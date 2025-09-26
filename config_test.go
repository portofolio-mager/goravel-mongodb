package mongodb

import (
	"fmt"
	"testing"

	mocksconfig "github.com/goravel/framework/mocks/config"
	"github.com/stretchr/testify/suite"

	"github.com/tonidy/goravel-mongodb/contracts"
)

type ConfigTestSuite struct {
	suite.Suite
	config     *Config
	connection string
	mockConfig *mocksconfig.Config
}

func TestConfigTestSuite(t *testing.T) {
	suite.Run(t, &ConfigTestSuite{
		connection: Name,
	})
}

func (s *ConfigTestSuite) SetupTest() {
	s.mockConfig = mocksconfig.NewConfig(s.T())
	s.config = NewConfig(s.mockConfig, s.connection)
}

func (s *ConfigTestSuite) TestReads() {
	// Test when configs is empty
	s.mockConfig.EXPECT().Get(fmt.Sprintf("database.connections.%s.read", s.connection)).Return(nil).Once()
	s.Nil(s.config.Readers())

	// Test when configs is not empty
	s.mockConfig.EXPECT().Get(fmt.Sprintf("database.connections.%s.read", s.connection)).Return([]contracts.Config{
		{
			URI:      "mongodb://localhost:27017",
			Database: "forge",
		},
	}).Once()
	s.Equal([]contracts.FullConfig{
		{
			Connection: s.connection,
			Driver:     Name,
			Config: contracts.Config{
				URI:      "mongodb://localhost:27017",
				Database: "forge",
			},
		},
	}, s.config.Readers())
}

func (s *ConfigTestSuite) TestWrites() {
	s.Run("success when configs is empty", func() {
		s.mockConfig.EXPECT().Get(fmt.Sprintf("database.connections.%s.write", s.connection)).Return(nil).Once()
		s.mockConfig.EXPECT().GetString(fmt.Sprintf("database.connections.%s.uri", s.connection)).Return("mongodb://localhost:27017").Once()
		s.mockConfig.EXPECT().GetString(fmt.Sprintf("database.connections.%s.database", s.connection)).Return("forge").Once()

		s.Equal([]contracts.FullConfig{
			{
				Connection: s.connection,
				Driver:     Name,
				Config: contracts.Config{
					URI:      "mongodb://localhost:27017",
					Database: "forge",
				},
			},
		}, s.config.Writers())
	})

	s.Run("success when configs is not empty", func() {
		s.mockConfig.EXPECT().Get(fmt.Sprintf("database.connections.%s.write", s.connection)).Return([]contracts.Config{
			{
				Database: "forge",
			},
		}).Once()
		s.mockConfig.EXPECT().GetString(fmt.Sprintf("database.connections.%s.uri", s.connection)).Return("mongodb://localhost:27017").Once()

		s.Equal([]contracts.FullConfig{
			{
				Connection: s.connection,
				Driver:     Name,
				Config: contracts.Config{
					URI:      "mongodb://localhost:27017",
					Database: "forge",
				},
			},
		}, s.config.Writers())
	})
}

func (s *ConfigTestSuite) TestFillDefault() {
	uri := "mongodb://localhost:27017"
	database := "forge"

	tests := []struct {
		name          string
		configs       []contracts.Config
		setup         func()
		expectConfigs []contracts.FullConfig
	}{
		{
			name:    "success when configs is empty",
			setup:   func() {},
			configs: []contracts.Config{},
		},
		{
			name:    "success when configs have item but key is empty",
			configs: []contracts.Config{{}},
			setup: func() {
				s.mockConfig.EXPECT().GetString(fmt.Sprintf("database.connections.%s.uri", s.connection)).Return(uri).Once()
				s.mockConfig.EXPECT().GetString(fmt.Sprintf("database.connections.%s.database", s.connection)).Return(database).Once()
			},
			expectConfigs: []contracts.FullConfig{
				{
					Connection: s.connection,
					Driver:     Name,
					Config: contracts.Config{
						URI:      uri,
						Database: database,
					},
				},
			},
		},
		{
			name: "success when configs have item",
			configs: []contracts.Config{
				{
					URI:      uri,
					Database: database,
				},
			},
			setup: func() {},
			expectConfigs: []contracts.FullConfig{
				{
					Connection: s.connection,
					Driver:     Name,
					Config: contracts.Config{
						URI:      uri,
						Database: database,
					},
				},
			},
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			test.setup()
			configs := s.config.fillDefault(test.configs)

			s.Equal(test.expectConfigs, configs)
		})
	}
}
