package auth

import (
	"context"
	"encoding/json"
	"errors"
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

type contextKey string

func (c contextKey) String() string {
	return "auth." + string(c)
}

const (
	auth0Domain    = "ecovo.auth0.com"
	authContextKey = contextKey("userInfo")
)

func getUserInfo(auth string) (*UserInfo, error) {
	req, err := http.NewRequest("GET", "https://"+auth0Domain+"/userinfo", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", auth)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	var userInfo UserInfo
	err = json.NewDecoder(resp.Body).Decode(&userInfo)
	if err != nil {
		return nil, err
	}
	return &userInfo, nil
}

// TokenValidatorMiddleware looks for an access token in the request headers
// to validate it and obtain a user's basic information by calling Auth0's
// Authentication API's /userinfo endpoint.
func TokenValidatorMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userInfo, err := getUserInfo(r.Header.Get("Authorization"))
		if err != nil {
			// TODO: Find a neater way to log what went wrong
			log.Print(err)

			// TODO: Write more informative error message
			w.WriteHeader(http.StatusUnauthorized)
		} else {
			ctx := context.WithValue(r.Context(), authContextKey, userInfo)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}

// UserInfoFromContext extracts a user's basic information that was placed in a
// request context by the token validation middleware.
func UserInfoFromContext(ctx context.Context) (*UserInfo, error) {
	if ctx == nil {
		return nil, errors.New("cannot get user info from nil context")
	}

	userInfo := ctx.Value(authContextKey).(*UserInfo)
	if userInfo == nil {
		return nil, errors.New("no user info found in context")
	}

	return userInfo, nil
}
