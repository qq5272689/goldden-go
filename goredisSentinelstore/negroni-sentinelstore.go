package goredisSentinelstore

import (
	"github.com/go-redis/redis/v8"
	nSessions "github.com/goincremental/negroni-sessions"
	gSessions "github.com/gorilla/sessions"
)

//New returns a new Sentinel store
func NewNegroniSentinelGobStore(Sentinel *redis.Client, sessionExpire int, keyPairs ...[]byte) nSessions.Store {
	store := NewSentinelGobStore(Sentinel, sessionExpire, keyPairs...)
	return &NegroniSentinleStore{store}
}

func NewNegroniSentinelJsonStore(Sentinel *redis.Client, sessionExpire int, keyPairs ...[]byte) nSessions.Store {
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
