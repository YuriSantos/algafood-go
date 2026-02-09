package exception

import "net/http"

type AuthenticationException struct {
	Message string
}

func NewAuthenticationException(message string) *AuthenticationException {
	return &AuthenticationException{Message: message}
}

func (e *AuthenticationException) Error() string {
	return e.Message
}

func (e *AuthenticationException) GetStatusCode() int {
	return http.StatusUnauthorized
}
