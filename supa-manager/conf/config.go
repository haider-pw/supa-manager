package conf

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"os"
)

type PostgresSettings struct {
	DiskSize       int    `json:"disk_size" split_words:"true" default:"10"`
	DefaultVersion string `json:"default_version" split_words:"true" default:"14.2"`
	DockerImage    string `json:"docker_image" split_words:"true" default:"supabase/postgres"`
}

type DomainSettings struct {
	StudioUrl  string  `json:"studio_url" split_words:"true" required:"true"`
	Base       string  `json:"base_url" required:"true"`
	DnsHookUrl *string `json:"dns_hook_url" split_words:"true"`
	DnsHookKey *string `json:"dns_hook_key" split_words:"true"`
}

type ProvisioningSettings struct {
	Enabled          bool   `json:"enabled" default:"true"`
	DockerHost       string `json:"docker_host" split_words:"true" default:"unix:///var/run/docker.sock"`
	ProjectsDir      string `json:"projects_dir" split_words:"true" default:"./projects"`
	BasePostgresPort int    `json:"base_postgres_port" split_words:"true" default:"5433"`
	BaseKongHTTPPort int    `json:"base_kong_http_port" split_words:"true" default:"54321"`
}

type Config struct {
	DatabaseUrl       string               `json:"database_url" split_words:"true" required:"true"`
	Port              int                  `json:"port" default:"8080"`
	EncryptionSecret  string               `json:"encryption_secret" split_words:"true" required:"true"`
	JwtSecret         string               `json:"jwt_secret" split_words:"true" required:"true"`
	AllowSignup       bool                 `json:"allow_signup" split_words:"true" default:"false"`
	ServiceVersionUrl string               `json:"service_version_url" split_words:"true" required:"true" default:"https://supamanager.io/updates"`
	Domain            DomainSettings       `json:"domain" required:"true"`
	Postgres          PostgresSettings     `json:"postgres" required:"true"`
	Provisioning      ProvisioningSettings `json:"provisioning"`
}

func LoadConfig(filename string) (*Config, error) {
	if _, err := os.Stat("./.env"); !os.IsNotExist(err) {
		if err := loadEnvironment(filename); err != nil {
			return nil, err
		}
	}
	config := new(Config)
	if err := envconfig.Process("", config); err != nil {
		return nil, err
	}
	return config, nil
}

func loadEnvironment(filename string) error {
	var err error
	if filename != "" {
		err = godotenv.Load(filename)
	} else {
		err = godotenv.Load()
		// handle if .env file does not exist, this is OK
		if os.IsNotExist(err) {
			return nil
		}
	}
	return err
}
