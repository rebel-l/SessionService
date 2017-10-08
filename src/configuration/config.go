package configuration

import (
	"encoding/json"
	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
	"os"
)

type Config struct {
	Service *Service
	Redis *redis.Options
}

func newConfig(filename string) *Config {
	c := new(Config)
	c.Service = newService()
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
	file, err := os.Open(filename)
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
