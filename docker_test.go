package mongodb

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type DockerTestSuite struct {
	suite.Suite
	database string
	docker   *Docker
}

func TestDockerTestSuite(t *testing.T) {
	suite.Run(t, new(DockerTestSuite))
}

func (s *DockerTestSuite) SetupTest() {
	s.database = "goravel"
	s.docker = NewDocker(s.database)
}

func (s *DockerTestSuite) Test_Build_Config_AddData_Fresh_Shutdown() {
	s.Nil(s.docker.Build())

	config := s.docker.Config()
	s.Equal(s.database, config.Database)
	s.Equal("mongodb", config.Driver)
	s.Equal("localhost", config.Host)
	s.Equal(27017, config.Port)

	s.Nil(s.docker.Fresh())

	databaseDriver, err := s.docker.Database("another")
	s.NoError(err)
	s.NotNil(databaseDriver)

	s.NoError(s.docker.Shutdown())
	s.NoError(databaseDriver.Shutdown())
}

func (s *DockerTestSuite) TestDatabase() {
	s.Nil(s.docker.Build())

	docker, err := s.docker.Database("another")
	s.Nil(err)
	s.NotNil(docker)

	dockerImpl := docker.(*Docker)
	s.Equal("another", dockerImpl.database)

	s.NoError(s.docker.Shutdown())
	s.NoError(docker.Shutdown())
}
