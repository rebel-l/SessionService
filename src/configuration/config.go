package configuration

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"github.com/rebel-l/sessionservice/src/authentication"
	log "github.com/sirupsen/logrus"
	"os"
)

type Config struct {
	Service *Service
	Redis *redis.Options
	AccountList []authentication.Account

	openFile func(string)(*os.File, error)
}

func newConfig(filename string) *Config {
	c := new(Config)
	c.openFile = os.Open
	c.Service = newService()
	c.Redis = new(redis.Options)
	err := c.loadFromFile(filename)
	if err != nil {
		log.Infof(
			"Not able to load config from %s, continue with defaults (or cli arguments). Reported error: %s",
			filename,
			err,
		)
	}
	return c
}

func (c *Config) loadFromFile(filename string) error {
	if filename == "" {
		return nil
	}

	log.Infof("Load config filename: %s", filename)
	file, err := c.openFile(filename)
	if err != nil {
		return err
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(c)
	if err != nil {
		return err
	}

	return nil
}
