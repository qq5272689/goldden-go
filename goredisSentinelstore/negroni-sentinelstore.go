package goredisSentinelstore

import (
	nSessions "github.com/goincremental/negroni-sessions"
	"github.com/gomodule/redigo/redis"
	gSessions "github.com/gorilla/sessions"
)

//New returns a new Sentinel store
func NewNegroniSentinelGobStore(Sentinel *redis.Pool, sessionExpire int, keyPairs ...[]byte) nSessions.Store {
	store := NewSentinelGobStore(Sentinel, sessionExpire, keyPairs...)
	return &NegroniSentinleStore{store}
}

func NewNegroniSentinelJsonStore(Sentinel *redis.Pool, sessionExpire int, keyPairs ...[]byte) nSessions.Store {
	store := NewSentinelJsonStore(Sentinel, sessionExpire, keyPairs...)
	return &NegroniSentinleStore{store}
}

type NegroniSentinleStore struct {
	*SentinleStore
}

func (c *NegroniSentinleStore) Options(options nSessions.Options) {
	c.SentinleStore.Options = &gSessions.Options{
		Path:     options.Path,
		Domain:   options.Domain,
		MaxAge:   options.MaxAge,
		Secure:   options.Secure,
		HttpOnly: options.HTTPOnly,
	}
}
