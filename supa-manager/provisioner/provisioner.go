package provisioner

import (
	"supamanager.io/supa-manager/database"
)

type Config struct {
	DockerHost       string
	ProjectsDir      string
	BasePostgresPort int
	BaseKongHTTPPort int
	DomainBase       string
	StudioURL        string
}

type Provisioner struct {
	config  *Config
	queries *database.Queries
}

func NewProvisioner(config *Config, queries *database.Queries) *Provisioner {
	return &Provisioner{
		config:  config,
		queries: queries,
	}
}
