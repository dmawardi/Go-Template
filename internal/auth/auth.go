package auth

import (
	"errors"
	"os"
	"time"

	"github.com/dmawardi/Go-Template/internal/config"
	"github.com/golang-jwt/jwt/v4"
)

var app *config.AppConfig

var JWTKey = []byte(os.Getenv("HMAC_SECRET"))

// JWTSecretKey := os.Getenv("HMAC_SECRET")
// var JWTKey = []byte("")

// Function called in main.go to connect app state to current file
func SetStateInAuth(a *config.AppConfig) {
	app = a
}

// Authorization

type AuthToken struct {
	Username string `json:"userID"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

// Generates a JSON web token based on user's details
func GenerateJWT(username, email, roleName string) (string, error) {
	// Build expiration time
	expirationTime := time.Now().Add(12 * time.Hour)

	// Build claims to be stored in token
	claims := &AuthToken{
		Email:    email,
		Username: username,
		Role:     roleName,
		StandardClaims: jwt.StandardClaims{
			// Set expiry
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Create new token using built claims and signing method
	authToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Decrypt token using key to generate string
	tokenString, err := authToken.SignedString(JWTKey)
	// If error
	if err != nil {
		return "", err
	}
	// else, return token string
	return tokenString, nil
}

// Validates and parses signed token
func ValidateAndParseToken(signedToken string) (tokenData interface{}, err error) {
	// Parse token and claims
	token, err := jwt.ParseWithClaims(
		signedToken,
		&AuthToken{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(JWTKey), nil
		},
	)
	if err != nil {
		err = errors.New("couldn't parse token")
		return nil, err
	}

	// Extract claims from parsed tocken
	claims, ok := token.Claims.(*AuthToken)
	// If failed
	if !ok {
		err = errors.New("couldn't parse claims")
		return nil, err
	}
	// If successful but expired
	if claims.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("token expired")
		return nil, err
	}
	// else return claims
	return claims, nil
}

// Authorizer takes user, action and asset and returns
// whether access should be provided based on auth policy
type Authorizer interface {
	HasPermission(userID, action, asset string) bool
}

// Takes the http method and returns a string based on it
// for authorization assessment
func ActionFromMethod(httpMethod string) string {
	switch httpMethod {
	case "GET":
		return "gather"
	case "POST":
		return "consume"
	case "DELETE":
		return "destroy"
	default:
		return ""
	}
}

// func (a *Authorizer) HasPermission(userID, action, asset string) bool {
// 	// Check for user
// 	user, err := app.DbClient.User.Get(app.Ctx, 8)
// 	if err != nil {
// 		// Unknown userID
// 		log.Print("Can't find user to check permissions. ID:", userID)
// 		return false
// 	}

// 	hasPermission, err := app.RBEnforcer.Enforce(user.Role, asset, action)
// 	if err != nil {
// 		log.Printf("User '%s' does not have permission to access '%s'", user.Username, asset)
// 	}

// 	return hasPermission

// 	// if hasPermission {
// 	// 	return true
// 	// }
// 	// // for _, role := range user.Roles {
// 	// // }

// 	// return false
// }
