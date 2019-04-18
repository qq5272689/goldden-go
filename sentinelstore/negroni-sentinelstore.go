package sentinelstore

import (
	nSessions "github.com/goincremental/negroni-sessions"
	gSessions "github.com/gorilla/sessions"
	"github.com/mediocregopher/radix/v3"
)

//New returns a new Sentinel store
func NewNegroniSentinelGobStore(Sentinel *radix.Sentinel, sessionExpire int, keyPairs ...[]byte) nSessions.Store {
	store := NewSentinelGobStore(Sentinel, sessionExpire, keyPairs...)
	return &NegroniSentinleStore{store}
}

func NewNegroniSentinelJsonStore(Sentinel *radix.Sentinel, sessionExpire int, keyPairs ...[]byte) nSessions.Store {
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
