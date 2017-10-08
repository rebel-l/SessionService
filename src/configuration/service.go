package configuration

import "github.com/sirupsen/logrus"

type Service struct {
	Port int
	LogLevel logrus.Level
}

func newService() *Service {
	s := new(Service)
	s.Port = ServiceDeafultPort
	s.LogLevel = ServiceDefaultLogLevel
	return s
}
