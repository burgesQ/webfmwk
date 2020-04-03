package jwt

import (
	"net/http"
	"strings"

	"github.com/burgesQ/webfmwk/v3"
	"github.com/dgrijalva/jwt-go"
	j "github.com/dgrijalva/jwt-go"
)

var _signingKey = []byte("test_signing_key")

//SetSigningKey set the jwt signing key
func SetSigningKey(key string) {
	_signingKey = []byte(key)
}

func GenToken(name string) (string, error) {
	token := j.NewWithClaims(j.SigningMethodHS256, j.MapClaims{
		"user": name,
	})
	tokenString, err := token.SignedString(_signingKey)
	return tokenString, err
}

func CheckToken(tokenString string) (j.Claims, error) {
	token, err := j.Parse(tokenString, func(token *j.Token) (interface{}, error) {
		return _signingKey, nil
	})
	if err != nil {
		return nil, err
	}
	return token.Claims, err
}

func errorJSON(str string) string {
	return `{"error":"` + str + `"}`
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if len(tokenString) == 0 {
			http.Error(w, errorJSON("Missing Authorization Header"), http.StatusUnauthorized)
			return
		}

		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
		claims, err := CheckToken(tokenString)

		if err != nil {
			http.Error(w, errorJSON("Forbidden"), http.StatusForbidden)
			return
		}
		webfmwk.GetLogger().Infof("Authenticated user %s\n", claims.(jwt.MapClaims)["user"].(string))

		next.ServeHTTP(w, r)
	})
}
