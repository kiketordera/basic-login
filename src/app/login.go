package app

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// TokenClaims is the token that we store in the cookies
type TokenClaims struct {
	Username string `json:"usr,omitempty"`
	Role     string `json:"rol,omitempty"`
	jwt.StandardClaims
}

// This middleware ensures that a request will be aborted with an error if the user is not logged in
// This method checks if the user if logged in
func checkToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		// If there's an error or if the token is empty
		// the user is not logged in
		cookie, err := c.Cookie(projectName)
		if err != nil {
			setSession(c, "", "")
			return
		}

		token, err := jwt.ParseWithClaims(cookie, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("Unexpected signing method: " + token.Header["alg"].(string))
			}

			// Return signing key to check the token
			return tokenSigningKey, nil
		})

		if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
			c.Set("username", claims.Username)
			c.Set("role", claims.Role)
			if claims.Username != "" {
				c.Set("isLoggedIn", true)
			} else {
				c.Set("isLoggedIn", false)
			}

		} else {
			c.Set("username", "")
			c.Set("role", "")
			c.Set("isLoggedIn", false)
			fmt.Println("Redirecting to /log-in from check token")
			c.Redirect(http.StatusFound, "/log-in")
		}
	}
}

// Set the session in the cookie
func setSession(c *gin.Context, username string, role string) {
	// Define token claims
	claims := TokenClaims{
		username,
		role,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Unix() + int64(sessionTime),
			Issuer:    projectName,
		},
	}
	// Create a new token object, specifying signing method and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(tokenSigningKey)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	c.SetCookie(projectName, tokenString, sessionTime, "", "", false, true)
}

// Ensures that the user is logged in
func ensureLoggedIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		u, e := c.Get("username")
		if e {
			username := u.(string)
			if username == "" {
				fmt.Println("Redirecting to /log-in from ensureLoggedIn")
				c.Redirect(http.StatusFound, "/log-in")
			}
		} else {
			c.Redirect(http.StatusFound, "/log-in")
		}
	}
}

// Ensures that the user is NOT logged in
func ensureNotLoggedIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		u, e := c.Get("username")
		// If the User is logged in, we send him to the Dashboard
		if e {
			username := u.(string)
			if username != "" {
				//Redirection when is logged
				fmt.Println("Redirecting to /dashboard from ensureNotLoggedIn")
				c.Redirect(http.StatusFound, "/users")
			}
		}
	}
}

// Logs out of the app for the user
func (VS *DB) logout(c *gin.Context) {
	// Clear the cookie
	c.SetCookie(projectName, "", -1, "", "", false, true)
	c.Set("isLoggedIn", false)
	// Redirect to the home page
	fmt.Println("Redirecting to /log-in from logout")
	c.Redirect(http.StatusFound, "/")
}
