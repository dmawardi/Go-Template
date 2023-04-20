package auth

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/dmawardi/Go-Template/internal/config"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
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
	UserID string `json:"userID"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.StandardClaims
}

// Setup RBAC enforcer based using gorm client. Connects to DB and builds base policy
func EnforcerSetup(db *gorm.DB) (*casbin.Enforcer, error) {
	// Grab environment variables for connection
	var DB_PORT string = os.Getenv("DB_PORT")

	// Build adapter
	adapter, err := gormadapter.NewAdapterByDB(db)
	// If error
	if err != nil {
		log.Fatal("Couldn't build adapter for enforcer: ", err, "\nDB PORT", DB_PORT)
		return nil, err
	}

	// Build path to policy model
	rbacModelPath := buildPathToPolicyModel()

	// Initialize RBAC Authorization
	enforcer, err := casbin.NewEnforcer(rbacModelPath, adapter)

	// If error
	if err != nil {
		log.Fatal("Couldn't build RBAC enforcer: ", err)
		return nil, err
	}

	// Create default policies if not already detected within system
	SetupCasbinPolicy(enforcer, DefaultPolicyList)

	// else
	return enforcer, nil
}

// Generates a JSON web token based on user's details
func GenerateJWT(userID int, email, roleName string) (string, error) {
	// Build expiration time
	expirationTime := time.Now().Add(12 * time.Hour)

	// Build claims to be stored in token
	claims := &AuthToken{
		Email: email,
		// Convert ID to string
		UserID: fmt.Sprint(userID),
		Role:   roleName,
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
func ValidateAndParseToken(w http.ResponseWriter, r *http.Request) (tokenData *AuthToken, err error) {
	// Grab request header
	header := r.Header
	// Extract token string from Authorization header by removing prefix "Bearer "
	_, tokenString, _ := strings.Cut(header.Get("Authorization"), " ")

	if tokenString == "" {
		err := errors.New("Authentication Token not detected")
		return nil, err
	}
	// Parse token string and claims. Filter through auth token
	token, err := jwt.ParseWithClaims(
		tokenString,
		&AuthToken{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(JWTKey), nil
		},
	)
	if err != nil {
		err = errors.New("couldn't parse token")
		return &AuthToken{}, err
	}

	// Extract claims from parsed tocken
	claims, ok := token.Claims.(*AuthToken)
	// If failed
	if !ok {
		err = errors.New("couldn't parse claims")
		return &AuthToken{}, err
	}
	// If successful but expired
	if claims.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("token expired")
		return &AuthToken{}, err
	}
	// else return claims
	return claims, nil
}

// Takes the http method and returns a string based on it
// for authorization assessment
func ActionFromMethod(httpMethod string) string {
	switch httpMethod {
	case "GET":
		return "read"
	case "POST":
		return "create"
	case "PUT":
		return "update"
	case "DELETE":
		return "delete"
	default:
		return ""
	}
}

// Set up policy settings in DB for casbin rules
func SetupCasbinPolicy(enforcer *casbin.Enforcer, sliceOfPolicies []policySet) {
	for _, policy := range sliceOfPolicies {

		// if enforcer does not already have policy
		if hasPolicy := enforcer.HasPolicy(policy.subject, policy.object, policy.action); !hasPolicy {
			// create policy
			enforcer.AddPolicy(policy.subject, policy.object, policy.action)
		}
	}
}

// Extracts user id from authentication token
func ExtractIdFromToken(w http.ResponseWriter, r *http.Request) (*int, error) {
	// Validate and parse the token
	tokenData, err := ValidateAndParseToken(w, r)
	// If error detected
	if err != nil {
		return nil, err
	}
	// Convert to int
	userId, err := strconv.Atoi(tokenData.UserID)
	if err != nil {
		return nil, err
	}

	return &userId, nil
}

// Returns a string that is agnostic for test and production usage
func buildPathToPolicyModel() string {
	// generate path
	dirPath, err := os.Getwd()
	if err != nil {
		log.Fatal("Could not get working directory")
	}

	// Split path to remove excess path when running tests
	splitPath := strings.Split(dirPath, "internal")

	// Grab initial part of path and join with path from project root directory
	rbacModelPath := splitPath[0] + "/internal/auth/rbac_model.conf"
	return rbacModelPath
}
