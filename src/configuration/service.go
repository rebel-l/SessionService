package configuration

import "github.com/sirupsen/logrus"

type Service struct {
	Port int
	LogLevel logrus.Level
	SessionLifetime int
}

func newService() *Service {
	s := new(Service)
	s.Port = ServiceDefaultPort
	s.LogLevel = ServiceDefaultLogLevel
	return s
}
