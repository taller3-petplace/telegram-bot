package sender

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

const (
	secretEnvVar                 = "SECRET"
	cryptographicAlgorithmEnvVar = "ALGORITHM"
	accessCodeEnvVar             = "ACCESS_CODE"
	jwtHeader                    = "Authorization"
)

// AccessControl middleware use by each endpoint to create a context.AppContext
func AccessControl() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := checkAccessToken(c.Request.Header.Get(jwtHeader))
		if err != nil {
			logrus.Error(err)
			errResponse := errorResponse{
				StatusCode: http.StatusUnauthorized,
				Message:    "Hi FBI? we have a badass trying to make requests against our app",
			}
			c.JSON(errResponse.StatusCode, errResponse)
			c.Abort()
			return
		}

		c.Next()
	}
}

// checkAccessToken checks if the JWT is valid and contains the correct access code
func checkAccessToken(tokenString string) error {
	if tokenString == "" {
		return fmt.Errorf("error token is missing")
	}

	secretKey := os.Getenv(secretEnvVar)
	if secretKey == "" {
		return fmt.Errorf("error secret key is missing")
	}

	algorithm := os.Getenv(cryptographicAlgorithmEnvVar)
	if algorithm == "" {
		return fmt.Errorf("error algorithm is missing")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Sanity check: the algorithm must be the same as the one in the env var
		if token.Method.Alg() != algorithm {
			return nil, fmt.Errorf("error unexpected signing method: %s", token.Method.Alg())
		}

		return []byte(secretKey), nil
	})

	if err != nil {
		return fmt.Errorf("error parsing JWT: %v", err)
	}

	// Access the claims if the signature is valid
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		var accessCode string
		accessCodeRaw, found := claims["access_code"]
		if found {
			accessCode = accessCodeRaw.(string)
		}
		if accessCode != os.Getenv(accessCodeEnvVar) {
			return fmt.Errorf("error invalid access code")
		}
		return nil
	}

	return fmt.Errorf("error invalid JWT")
}
