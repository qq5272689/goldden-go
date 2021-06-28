package ldap

import (
	"errors"
	"gitee.com/goldden-go/goldden-go/pkg/models"
	"gitee.com/goldden-go/goldden-go/pkg/utils/auth"
	"gitee.com/goldden-go/goldden-go/pkg/utils/logger"
	"go.uber.org/zap"
)

// ErrNoLDAPServers is returned when there is no LDAP servers specified
var ErrNoLDAPServers = errors.New("no LDAP servers are configured")

// ErrDidNotFindUser if request for user is unsuccessful
var ErrDidNotFindUser = errors.New("did not find a user")

// ServerStatus holds the LDAP server status
type ServerStatus struct {
	Host      string
	Port      int
	Available bool
	Error     error
}

// IMultiLDAP is interface for MultiLDAP
type IMultiLDAP interface {
	Ping() ([]*ServerStatus, error)
	Login(query *auth.LoginData) (
		*models.User, error,
	)

	Users(logins []string) (
		[]*models.User, error,
	)

	User(login string) (
		*models.User, ServerConfig, error,
	)
}

// MultiLDAP is basic struct of LDAP authorization
type MultiLDAP struct {
	configs []*ServerConfig
}

// New creates the new LDAP auth
func NewMultiLDAP(configs []*ServerConfig) IMultiLDAP {
	return &MultiLDAP{
		configs: configs,
	}
}

// Ping dials each of the LDAP servers and returns their status. If the server is unavailable, it also returns the error.
func (multiples *MultiLDAP) Ping() ([]*ServerStatus, error) {
	if len(multiples.configs) == 0 {
		return nil, ErrNoLDAPServers
	}

	serverStatuses := []*ServerStatus{}
	for _, config := range multiples.configs {
		status := &ServerStatus{}

		status.Host = config.Host
		status.Port = config.Port

		server := NewLDAPServer(config)
		err := server.Dial()

		if err == nil {
			status.Available = true
			serverStatuses = append(serverStatuses, status)
			server.Close()
		} else {
			status.Available = false
			status.Error = err
			serverStatuses = append(serverStatuses, status)
		}
	}

	return serverStatuses, nil
}

// Login tries to log in the user in multiples LDAP
func (multiples *MultiLDAP) Login(query *auth.LoginData) (
	*models.User, error,
) {
	if len(multiples.configs) == 0 {
		return nil, ErrNoLDAPServers
	}

	for index, config := range multiples.configs {
		server := NewLDAPServer(config)

		if err := server.Dial(); err != nil {
			logDialFailure(err, config)

			// Only return an error if it is the last server so we can try next server
			if index == len(multiples.configs)-1 {
				return nil, err
			}
			continue
		}

		defer server.Close()

		user, err := server.Login(query)
		// FIXME
		if user != nil {
			return user, nil
		}
		if err != nil {
			if isSilentError(err) {
				logger.Debug(
					"unable to login with LDAP - skipping server",
					zap.String("host", config.Host),
					zap.Int("port", config.Port),
					zap.Error(err),
				)
				continue
			}

			return nil, err
		}
	}

	// Return invalid credentials if we couldn't find the user anywhere
	return nil, ErrInvalidCredentials
}

// User attempts to find an user by login/username by searching into all of the configured LDAP servers. Then, if the user is found it returns the user alongisde the server it was found.
func (multiples *MultiLDAP) User(login string) (
	*models.User,
	ServerConfig,
	error,
) {
	if len(multiples.configs) == 0 {
		return nil, ServerConfig{}, ErrNoLDAPServers
	}

	search := []string{login}
	for index, config := range multiples.configs {
		server := NewLDAPServer(config)

		if err := server.Dial(); err != nil {
			logDialFailure(err, config)

			// Only return an error if it is the last server so we can try next server
			if index == len(multiples.configs)-1 {
				return nil, *config, err
			}
			continue
		}

		defer server.Close()

		if err := server.Bind(); err != nil {
			return nil, *config, err
		}

		users, err := server.Users(search)
		if err != nil {
			return nil, *config, err
		}

		if len(users) != 0 {
			return users[0], *config, nil
		}
	}

	return nil, ServerConfig{}, ErrDidNotFindUser
}

// Users gets users from multiple LDAP servers
func (multiples *MultiLDAP) Users(logins []string) (
	[]*models.User,
	error,
) {
	var result []*models.User

	if len(multiples.configs) == 0 {
		return nil, ErrNoLDAPServers
	}

	for index, config := range multiples.configs {
		server := NewLDAPServer(config)

		if err := server.Dial(); err != nil {
			logDialFailure(err, config)

			// Only return an error if it is the last server so we can try next server
			if index == len(multiples.configs)-1 {
				return nil, err
			}
			continue
		}

		defer server.Close()

		if err := server.Bind(); err != nil {
			return nil, err
		}

		users, err := server.Users(logins)
		if err != nil {
			return nil, err
		}
		result = append(result, users...)
	}

	return result, nil
}

// isSilentError evaluates an error and tells whenever we should fail the LDAP request
// immediately or if we should continue into other LDAP servers
func isSilentError(err error) bool {
	continueErrs := []error{ErrInvalidCredentials, ErrCouldNotFindUser}

	for _, cerr := range continueErrs {
		if errors.Is(err, cerr) {
			return true
		}
	}

	return false
}

func logDialFailure(err error, config *ServerConfig) {
	logger.Debug(
		"unable to dial LDAP server",
		zap.String("host", config.Host),
		zap.Int("port", config.Port),
		zap.Error(err),
	)
}
