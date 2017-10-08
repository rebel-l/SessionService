package configuration

import (
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"runtime"
)

const Version = "0.1.0"

func Init() *Config {
	filename,
	servicePort,
	serviceLogLevel,
	version,
	redisAddr,
	redisPassword,
	redisDb := parseCliArgs()

	if *version {
		printVersion()
	}

	c := newConfig(*filename)
	if *servicePort > 0 {
		c.Service.Port = *servicePort
	}

	if *serviceLogLevel > 0 {
		c.Service.LogLevel = logrus.Level(*serviceLogLevel)
	}

	if *redisAddr != "" {
		c.Redis.Addr = *redisAddr
	}

	if *redisPassword != "" {
		c.Redis.Password = *redisPassword
	}

	if *redisDb > 0 {
		c.Redis.DB = *redisDb
	}

	return c
}

func parseCliArgs() (
	filename *string,
	servicePort *int,
	serviceLogLevel *int,
	version *bool,
	redisAddr *string,
	redisPassword *string,
	redisDb *int,
) {
	filename = flag.String(
		"config",
		"",
		"A filename from where to load a config from. The file must be in JSON format. This parameter is optional.",
	)

	servicePort = flag.Int(
		"Service.Port",
		0,
		"The port on which this service should serve. This parameter is optional.",
	)

	serviceLogLevel = flag.Int(
		"Service.LogLevel",
		0,
		"The level this service logging on. 0 = highest level, 5 = lowest level. This parameter is optional.",
	)

	version = flag.Bool("v", false, "Shows Version information.")

	redisAddr = flag.String("Redis.Addr", "", "The Hostname of the Redis server. Format: hostname:port. This parameter is optional.")

	redisPassword = flag.String("Redis.Password", "", "The password of the Redis server. This parameter is optional.")

	redisDb = flag.Int("Redis.DB", 0, "The DB of the Redis server. This parameter is optional.")

	flag.Parse()
	return
}

func printVersion() {
	fmt.Printf("SessionService Version: %s\n", Version)
	fmt.Printf("Runtime Version: %s\n\n", runtime.Version())
	os.Exit(0)
}
