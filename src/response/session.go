package response

import (
	"github.com/pborman/uuid"
	"math/rand"
	"time"
)

type Session struct {
	Id string `json:"id"`
	Data map[string]string `json:"data"`
	Lifetime int `json:"lifetime"`
	Expires int64 `json:"expires"`
	Domain string `json:"domain"`
}

func NewSession(id string, lifetime int) *Session {
	session := new(Session)
	if id != "" {
		session.Id = id
	} else {
		session.GenerateId()
	}

	if lifetime == 0 {
		session.Lifetime = LIFETIME
	} else {
		session.Lifetime = lifetime
	}

	t := time.Now()
	t = t.Add(time.Duration(session.Lifetime) * time.Second)
	session.Expires = t.Unix()
	return session
}

func (s *Session) GenerateId() {
	rsource := rand.New(rand.NewSource(rand.Int63()))
	uuid.SetRand(rsource)
	s.Id = uuid.NewRandom().String()
}
