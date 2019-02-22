package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

// UserInfo contains a user's basic information extracted from an access token.
type UserInfo struct {
	ID        string `json:"sub,omitempty"`
	FirstName string `json:"given_name"`
	LastName  string `json:"family_name"`
	Picture   string `json:"picture"`
	Email     string `json:"email"`
}

// Config contains the information required to configure a validator to make
// requests to validate a request's authorization header.
type Config struct {
	// Domain represents the domain where the user info endpoint is hosted.
	Domain string
}

// Validate looks at the configuration's contents to ensure it has all the
// required fields.
func (conf *Config) validate() error {
	if conf.Domain == "" {
		return errors.New("missing domain")
	}

	return nil
}

// Validator is an interface representing the ability to validate an
// authorization header and return the authenticated user's information.
type Validator interface {
	// Validate validates an authorization and returns the authenticated user's
	// information.
	Validate(authHeader string) (*UserInfo, error)
}

// A TokenValidator is a validator that validates a bearer token in an
// authorization header by making a request to a /userinfo endpoint.
type TokenValidator struct {
	conf *Config
}

// NewTokenValidator creates a new token validator with the given
// configuration.
func NewTokenValidator(conf *Config) (Validator, error) {
	if conf == nil {
		return nil, fmt.Errorf("auth: missing configuration")
	}

	err := conf.validate()
	if err != nil {
		return nil, fmt.Errorf("auth: configuration %s", err)
	}

	return &TokenValidator{conf}, nil
}

// Validate makes a request to the /userinfo endpoint on the domain specified
// in the token validator's configuration to validate the bearer token present
// in the authorization header and returns the authenticated user's
// information.
func (validator *TokenValidator) Validate(authHeader string) (*UserInfo, error) {
	req, err := http.NewRequest("GET", "https://"+validator.conf.Domain+"/userinfo", nil)
	if err != nil {
		return nil, fmt.Errorf("auth: failed to create request (%s)", err)
	}

	req.Header.Set("Authorization", authHeader)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("auth: failed to make request (%s)", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("auth: failed to validate token")
	}

	var userInfo UserInfo
	err = json.NewDecoder(resp.Body).Decode(&userInfo)
	if err != nil {
		return nil, fmt.Errorf("auth: failed to decode user info (%s)", err)
	}
	return &userInfo, nil
}

type contextKey string

func (c contextKey) String() string {
	return "auth." + string(c)
}

const (
	userInfoContextKey = contextKey("userInfo")
)

// ValidationMiddleware validates a request's authorization header using the
// given validator to ensure that the user is authorized to access an endpoint
// and extracts the authenticated user's information.
//
// The authenticated user's information placed in the request's context and can
// be accessed by using the UserInfoFromContext utility function.
func ValidationMiddleware(validator Validator, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userInfo, err := validator.Validate(r.Header.Get("Authorization"))
		if err != nil {
			log.Println(err)

			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		} else {
			ctx := context.WithValue(r.Context(), userInfoContextKey, userInfo)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}

// UserInfoFromContext extracts an authenticated user's information from
// request's context.
func UserInfoFromContext(ctx context.Context) (*UserInfo, error) {
	if ctx == nil {
		return nil, fmt.Errorf("auth: request context is nil")
	}

	userInfo := ctx.Value(userInfoContextKey).(*UserInfo)
	if userInfo == nil {
		return nil, fmt.Errorf("auth: %s not found in context", userInfoContextKey)
	}

	return userInfo, nil
}
