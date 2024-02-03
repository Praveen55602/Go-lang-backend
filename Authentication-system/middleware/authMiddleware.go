package middleware

import (
	helper "auth/helpers"
	"errors"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func Authenticate(c *gin.Context) {
	// Extract the token from the Authorization header
	token, refreshToken, tokenFetchErr := getTokensFromRequest(c)

	if tokenFetchErr != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": tokenFetchErr.Error()})
		c.Abort()
		return
	}

	claims, err := helper.ValidateToken(token)
	//ve, ok := err.(*jwt.ValidationError)

	if err != nil {
		//first we convert the error into the jwtValidationError which it will get converted because helper.Validate function is returning validation error only, then if the it's converted successfully then it bitwise and's it with a flag that is a contant for token expiry in jwt pakage. if the location of set bit in error is different then the flag the result will be 0 hence, token is expired.
		if ve, ok := err.(*jwt.ValidationError); refreshToken != "" && ok && ve.Errors&jwt.ValidationErrorExpired != 0 {
			//we'll generate a new token
			token, refreshToken, updateErr := helper.UpdateAllTokens(refreshToken)
			if updateErr != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": updateErr.Error()})
				c.Abort()
				return
			}
			c.Header("token", "Bearer "+token)
			c.Header("refresh-token", "Bearer "+refreshToken)
			c.JSON(http.StatusOK, gin.H{"message": "Token refreshed successfully"})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

	}

	//now setting up the payload info which is right now present only in the header that to encoded. into the context object as key value pair so that the user related info is available for the subsequent middleware functions or request handlers. and also it avoids the need to create global variables in the code to store user info
	c.Set("email", claims.Email)
	c.Set("first_name", claims.First_name)
	c.Set("last_name", claims.Last_name)
	c.Set("uid", claims.Uid)
	c.Set("user_type", claims.User_type)

	c.Next() // passing control to the next middle ware or next handler
}

func getTokensFromRequest(c *gin.Context) (string, string, error) {
	var authHeader = c.GetHeader("token")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
		c.Abort()
		return "", "", errors.New("token not found")
	}

	// Check if the Authorization header has the expected format ("Bearer <token>")
	var parts = strings.Fields(authHeader) // this splits the authe header at " "
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
		c.Abort()
		return "", "", errors.New("invalid token format")
	}

	// Extract the token from the header
	tokenString := parts[1]
	authHeader = c.GetHeader("refresh-token")
	if authHeader == "" {
		return tokenString, "", nil
	}

	// Check if the Authorization header has the expected format ("Bearer <token>")
	parts = strings.Fields(authHeader) // this splits the authe header at " "
	if len(parts) != 2 || parts[0] != "Bearer" {
		return tokenString, "", nil
	}
	refreshToken := parts[1]

	return tokenString, refreshToken, nil
}
