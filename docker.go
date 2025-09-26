package mongodb

import (
	"github.com/goravel/framework/contracts/testing/docker"
)

var _ docker.DatabaseDriver = &Docker{}

type Docker struct {
	database string
}

func NewDocker(database string) *Docker {
	return &Docker{
		database: database,
	}
}

func (r *Docker) Build() error {
	return nil
}

func (r *Docker) Config() docker.DatabaseConfig {
	return docker.DatabaseConfig{
		Driver:   Name,
		Host:     r.Host(),
		Port:     r.Port(),
		Database: r.database,
		Username: r.Username(),
		Password: r.Password(),
	}
}

func (r *Docker) Database(name string) (docker.DatabaseDriver, error) {
	return NewDocker(name), nil
}

func (r *Docker) Driver() string {
	return Name
}

func (r *Docker) Fresh() error {
	return nil
}

func (r *Docker) Host() string {
	return "localhost"
}

func (r *Docker) Image(image docker.Image) {
}

func (r *Docker) Ready() error {
	return nil
}

func (r *Docker) Reuse(containerID string, port int) error {
	return nil
}

func (r *Docker) Shutdown() error {
	return nil
}

func (r *Docker) Password() string {
	return "123123"
}

func (r *Docker) Port() int {
	return 27017
}

func (r *Docker) Username() string {
	return "root"
}
