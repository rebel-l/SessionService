package configuration

import (
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"runtime"
	"github.com/rebel-l/sessionservice/src/response"
)

const Version = "0.1.0"

var instance *Parser

type Parser struct {
	filename               string
	servicePort            int
	serviceLogLevel        int
	serviceSessionLifetime int
	version                bool
	redisAddr              string
	redisPassword          string
	redisDb                int
	printVersion           func()
}

func GetParser() *Parser {
	if instance == nil {
		instance = new(Parser)
		instance.printVersion = printVersion
		instance.parseCliArgs()
	}

	return instance
}

func (p *Parser) Parse() *Config {
	if p.version {
		p.printVersion()
	}

	c := newConfig(p.filename)
	if p.servicePort > 0 {
		c.Service.Port = p.servicePort
	}

	if p.serviceLogLevel > 0 {
		c.Service.LogLevel = log.Level(p.serviceLogLevel)
	}

	if p.serviceSessionLifetime > 0 {
		c.Service.SessionLifetime = p.serviceSessionLifetime
	}

	if p.redisAddr != "" {
		c.Redis.Addr = p.redisAddr
	}

	if p.redisPassword != "" {
		c.Redis.Password = p.redisPassword
	}

	if p.redisDb > 0 {
		c.Redis.DB = p.redisDb
	}

	return c
}

func (p *Parser) parseCliArgs() {
	var filename *string
	var servicePort *int
	var serviceLogLevel *int
	var serviceSessionLifetime *int
	var version *bool
	var redisAddr *string
	var redisPassword *string
	var redisDb *int

	if flag.Lookup("config") == nil {
		filename = flag.String(
			"config",
			"",
			"A filename from where to load a config from. The file must be in JSON format. This parameter is optional.",
		)
	}

	if flag.Lookup("Service.Port") == nil {
		servicePort = flag.Int(
			"Service.Port",
			0,
			"The port on which this service should serve. This parameter is optional.",
		)
	}

	if flag.Lookup("Service.LogLevel") == nil {
		serviceLogLevel = flag.Int(
			"Service.LogLevel",
			0,
			"The level this service logging on. 0 = highest level, 5 = lowest level. This parameter is optional.",
		)
	}

	if flag.Lookup("Service.SessionLifetime") == nil {
		serviceSessionLifetime = flag.Int(
			"Service.SessionLifetime",
			0,
			fmt.Sprintf("The lifetime of the session in seconds. Default is %d", response.LIFETIME),
		)
	}

	if flag.Lookup("v") == nil {
		version = flag.Bool("v", false, "Shows Version information.")
	}

	if flag.Lookup("Redis.Addr") == nil {
		redisAddr = flag.String("Redis.Addr", "", "The Hostname of the Redis server. Format: hostname:port. This parameter is optional.")
	}

	if flag.Lookup("Redis.Password") == nil {
		redisPassword = flag.String("Redis.Password", "", "The password of the Redis server. This parameter is optional.")
	}

	if flag.Lookup("Redis.DB") == nil {
		redisDb = flag.Int("Redis.DB", 0, "The DB of the Redis server. This parameter is optional.")
	}

	flag.Parse()

	p.filename = *filename
	p.servicePort = *servicePort
	p.serviceLogLevel = *serviceLogLevel
	p.serviceSessionLifetime = *serviceSessionLifetime
	p.redisAddr = *redisAddr
	p.redisDb = *redisDb
	p.redisPassword = *redisPassword
	p.version = *version

	return
}

func printVersion() {
	fmt.Printf("SessionService Version: %s\n", Version)
	fmt.Printf("Runtime Version: %s\n\n", runtime.Version())
	os.Exit(0)
}
