package server

import (
	"io/ioutil"

	"github.com/BurntSushi/toml"
)

// Config Server Configuration
type Config struct {
	Addr  string
	Debug bool
	Api   struct {
		Version   string
		JwtSecret string
	}

	Database struct {
		Driver         string
		ConnString     string
		Pool           int
		CacheEnabled   bool
		CacheNamespace string
		CacheRedisAddr string
		CacheRedisPwd  string
	}

	Worker struct {
		Enabled   bool
		Addr      string
		Namespace string
	}
}

func NewDefaultConfig() *Config {
	return &Config{}
}

// ParseToml loads a TOML configuration from a provided path and
// returns the AST produced from the TOML parser. When loading the file, it
// will find environment variables and replace them.
func (c *Config) ParseToml(fpath string) error {
	buf, err := ioutil.ReadFile(fpath)
	if err != nil {
		return err
	}

	if err := toml.Unmarshal(buf, &c); err != nil {
		return err
	}
	return nil
}
