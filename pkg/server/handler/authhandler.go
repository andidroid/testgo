package handler

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/andidroid/testgo/pkg/util"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

const (
	AUTHORIZATION_HEADER string = "Authorization"
	BASIC_SCHEMA         string = "Basic "
	BEARER_SCHEMA        string = "Bearer "
)

type Claims struct {
	//Username string `json:"username"`
	jwt.StandardClaims
}

var jwtSecret string
var jwtAudience string
var jwtIssuer string

func init() {
	jwtSecret = util.LookupEnv("JWT_SECRET", "password")
	jwtAudience = util.LookupEnv("JWT_AUDIENCE", "password")
	jwtIssuer = util.LookupEnv("JWT_ISSUER", "password")
}

func AuthMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {

		authHeader := c.Request.Header.Get(AUTHORIZATION_HEADER)
		log.Println("Authorization Header:", authHeader)
		// token := authHeader[len(BEARER_SCHEMA):]
		// log.Println("token:", token)

		if authHeader == "" {
			// log.Fatalf("no %s  %s", AUTHORIZATION_HEADER, authHeader)
			// c.JSON(http.StatusUnauthorized, gin.H{"error": "no Authorization header"})
			// return
		}

		// Confirm the request is sending Basic Authentication credentials.
		if !strings.HasPrefix(authHeader, BASIC_SCHEMA) && !strings.HasPrefix(authHeader, BEARER_SCHEMA) {
			// log.Fatalf("no %s", AUTHORIZATION_HEADER)
			// c.JSON(http.StatusUnauthorized, gin.H{"error": "no basic or bearer scheme in Authorization header"})
			// return
		}

		// Get the token from the request header
		// The first six characters are skipped - e.g. "Basic ".
		if strings.HasPrefix(authHeader, BASIC_SCHEMA) {
			str, err := base64.StdEncoding.DecodeString(authHeader[len(BASIC_SCHEMA):])
			if err != nil {
				log.Fatalf("Request to MongoDB failed: %s", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			creds := strings.Split(string(str), ":")

			log.Println(creds[0])
			return
		} else if strings.HasPrefix(authHeader, BEARER_SCHEMA) {

			tokenString := authHeader[len(BEARER_SCHEMA):]
			log.Println(tokenString)

			// token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// 	return []byte(os.Getenv("JWT_SECRET")), nil
			// })

			// https://github.com/PacktPublishing/Building-Distributed-Applications-in-Gin/blob/main/chapter04/api/handlers/auth.go
			claims := &Claims{}
			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(jwtSecret), nil
			})
			// if err != nil {
			// 	c.AbortWithStatus(http.StatusUnauthorized)
			// }
			// if !token.Valid {
			// 	c.AbortWithStatus(http.StatusUnauthorized)
			// }

			util.CheckErr(err)
			fmt.Println(token)
			fmt.Println(token.Valid)
			fmt.Println(token.Claims)
			fmt.Println(token.Claims.Valid())

			c.Set("roles", &token.Claims)

		} else {
			// log.Fatalf("no %s", AUTHORIZATION_HEADER)
			// c.JSON(http.StatusUnauthorized, gin.H{"error": "no basic or bearer scheme in Authorization header"})

			// return
		}

		// https://github.com/PacktPublishing/Building-Distributed-Applications-in-Gin/blob/main/chapter04/auth0/handlers/auth.go
		// var auth0Domain = "https://" + os.Getenv("AUTH0_DOMAIN") + "/"
		// client := auth0.NewJWKClient(auth0.JWKClientOptions{URI: auth0Domain + ".well-known/jwks.json"}, nil)
		// configuration := auth0.NewConfiguration(client, []string{os.Getenv("AUTH0_API_IDENTIFIER")}, auth0Domain, jose.RS256)
		// validator := auth0.NewValidator(configuration, nil)

		// _, err := validator.ValidateRequest(c.Request)

		// if err != nil {
		// 	c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
		// 	c.Abort()
		// 	return
		// }

		// var clientID = os.Getenv("GOOGLE_OAUTH2_CLIENT_ID")
		// var clientSecret = os.Getenv("GOOGLE_OAUTH2_CLIENT_SECRET")

		// ctx := context.Background()

		// provider, err := oidc.NewProvider(ctx, "https://accounts.google.com")
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// oidcConfig := &oidc.Config{
		// 	ClientID: clientID,
		// }
		// verifier := provider.Verifier(oidcConfig)
		// config := oauth2.Config{
		// 	ClientID:     clientID,
		// 	ClientSecret: clientSecret,
		// 	Endpoint:     provider.Endpoint(),
		// 	RedirectURL:  "http://127.0.0.1:5556/auth/google/callback",
		// 	Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
		// }
		// idToken, err := verifier.Verify(ctx, authHeader)
		// if err != nil {
		// 	http.Error(c.Writer, "Failed to verify ID Token: "+err.Error(), http.StatusInternalServerError)
		// 	return
		// }

		// Set example variable
		// c.Set("roles", &roles)

		c.Next()

	}
}

// https://auth0.com/blog/authentication-in-golang/
func getTokenString(token *jwt.Token) (interface{}, error) {
	// Verify 'aud' claim
	aud := jwtAudience
	checkAud := token.Claims.(jwt.MapClaims).VerifyAudience(aud, false)
	if !checkAud {
		return token, errors.New("Invalid audience.")
	}
	// Verify 'iss' claim
	iss := jwtIssuer
	checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
	if !checkIss {
		return token, errors.New("Invalid issuer.")
	}

	cert, err := getPemCert(token)
	if err != nil {
		panic(err.Error())
	}

	result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
	return result, nil
}

// SigningMethod: jwt.SigningMethodRS256,

func getPemCert(token *jwt.Token) (string, error) {
	cert := ""
	resp, err := http.Get("https://YOUR_DOMAIN/.well-known/jwks.json")

	if err != nil {
		return cert, err
	}
	defer resp.Body.Close()

	var jwks = Jwks{}
	err = json.NewDecoder(resp.Body).Decode(&jwks)

	if err != nil {
		return cert, err
	}

	for k, _ := range jwks.Keys {
		if token.Header["kid"] == jwks.Keys[k].Kid {
			cert = "-----BEGIN CERTIFICATE-----\n" + jwks.Keys[k].X5c[0] + "\n-----END CERTIFICATE-----"
		}
	}

	if cert == "" {
		err := errors.New("Unable to find appropriate key.")
		return cert, err
	}

	return cert, nil
}

type Jwks struct {
	Keys []JSONWebKeys `json:"keys"`
}

//http://127.0.0.1:8180/auth/realms/testgorealm/.well-known/openid-configuration
//http://127.0.0.1:8180/auth/realms/testgorealm/protocol/openid-connect/certs
type JSONWebKeys struct {
	Kty     string   `json:"kty"`
	Kid     string   `json:"kid"`
	Alg     string   `json:"alg"`
	Use     string   `json:"use"`
	N       string   `json:"n"`
	E       string   `json:"e"`
	X5c     []string `json:"x5c"`
	// X5t     string   `json:"x5t"`
	// X5tS256 string   `json:"x5t#S256"`
}
