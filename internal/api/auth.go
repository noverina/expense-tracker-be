package api

import (
	"crypto/rsa"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var authColl *mongo.Collection

func InitAuth() {
	collName := os.Getenv("AUTH_COLL")
	authColl = GetDB().Collection(collName)
}

func loadPrivateKey() (*rsa.PrivateKey, error) {
	filePath := os.Getenv("PRIVATE_KEY")
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	keyData, err := os.ReadFile(dir + "/" + filePath)
	if err != nil {
		return nil, err
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(keyData)
	if err != nil {
		return nil, err
	}

	return key, nil
}

func loadPublicKey() (*rsa.PublicKey, error) {
	filePath := os.Getenv("PUBLIC_KEY")
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	keyData, err := os.ReadFile(dir + "/" + filePath)
	if err != nil {
		return nil, err
	}

	key, err := jwt.ParseRSAPublicKeyFromPEM(keyData)
	if err != nil {
		return nil, err
	}

	return key, nil
}

func validateJWT(tokenString string) (*jwt.Token, error) {
	publicKey, err := loadPublicKey()
	if err != nil {
		return nil, err
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("invalid signing method", token.Header["alg"])
		}
		return publicKey, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return token, nil
}

func generateJWT(identifier string, isAccess bool) (string, error) {
	privateKey, err := loadPrivateKey()
	if err != nil {
		return "", err
	}

	expiryAccess, err := strconv.Atoi(os.Getenv("ACCESS_EXPIRY"))
	if err != nil {
		return "", err
	}
	expiryRefresh, err := strconv.Atoi(os.Getenv("REFRESH_EXPIRY"))
	if err != nil {
		return "", err
	}
	var expiry time.Duration
	if isAccess {
		expiry = time.Hour * time.Duration(expiryAccess)
	} else {
		expiry = time.Hour * time.Duration(expiryRefresh)
	}

	claims := jwt.MapClaims{
		"iss": "expense-tracker-be",
		"aud": identifier,
		"iat": time.Now(),
		"exp": time.Now().Add(expiry).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// TODO this function is only for dev/testing purpose! inject data via db later
func Register(c *gin.Context, identifier string, secretKey string) (int, error) {
	secretKeyByte := []byte(secretKey)

	hashedSecretKey, err := bcrypt.GenerateFromPassword(secretKeyByte, 12)
	if err != nil {
		LogError("unable to encrypt password", "err", err)
		return http.StatusInternalServerError, err
	}

	accessToken, err := generateJWT(identifier, true)
	if err != nil {
		LogError("unable to generate access jwt", "err", err)
		return http.StatusInternalServerError, err
	}
	refreshToken, err := generateJWT(identifier, false)
	if err != nil {
		LogError("unable to generate refresh jwt", "err", err)
		return http.StatusInternalServerError, err
	}

	entity := Client{
		Identifier: identifier,
		SecretKey:  string(hashedSecretKey),
		Access:     accessToken,
		Refresh:    refreshToken,
	}

	if _, err := authColl.InsertOne(c, entity); err != nil {
		LogError("unable to insert", "err", err, "id", identifier)
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func Access(c *gin.Context, identifier string, secretKey string) (*AuthInfo, int, error) {
	filter := bson.M{"_id": identifier}

	var result Client
	if err := authColl.FindOne(c, filter).Decode(&result); err != nil && err != mongo.ErrNoDocuments {
		LogError("query failed", "err", err)
		return nil, http.StatusInternalServerError, err
	} else if err == mongo.ErrNoDocuments {
		return nil, http.StatusUnauthorized, fmt.Errorf("invalid credential")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(result.SecretKey), []byte(secretKey)); err != nil {
		return nil, http.StatusUnauthorized, fmt.Errorf("invalid credential")
	}

	accessToken, err := generateJWT(identifier, true)
	if err != nil {
		LogError("unable to generate access jwt", "err", err)
		return nil, http.StatusInternalServerError, err
	}
	refreshToken, err := generateJWT(identifier, false)
	if err != nil {
		LogError("unable to generate refresh jwt", "err", err)
		return nil, http.StatusInternalServerError, err
	}

	entity := bson.M{
		"$set": bson.M{
			"access":  accessToken,
			"refresh": refreshToken,
		},
	}

	if _, err := authColl.UpdateByID(c, result.Identifier, entity); err != nil {
		LogError("unable to update", "err", err, "id", result.Identifier)
		return nil, http.StatusInternalServerError, err
	}

	output := AuthInfo{AccessToken: accessToken, RefreshToken: refreshToken}
	return &output, http.StatusOK, nil

}

func Refresh(c *gin.Context, tokenString string) (string, int, error) {
	if _, err := validateJWT(tokenString); err != nil {
		return "", http.StatusBadRequest, fmt.Errorf("invalid token")
	} else {
		publicKey, err := loadPublicKey()
		if err != nil {
			return "", http.StatusInternalServerError, err
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("invalid signing method", token.Header["alg"])
			}
			return publicKey, nil
		})
		audience, err := token.Claims.GetAudience()
		if err != nil {
			return "", http.StatusInternalServerError, err
		}

		refreshedToken, err := generateJWT(audience[0], false)
		if err != nil {
			LogError("unable to generate jwt", "err", err)
			return "", http.StatusInternalServerError, nil
		}
		entity := bson.M{
			"$set": bson.M{
				"access": refreshedToken,
			},
		}

		if _, err := authColl.UpdateByID(c, audience[0], entity); err != nil {
			LogError("unable to update", "err", err, "id", audience[0])
			return "", http.StatusInternalServerError, err
		}

		return refreshedToken, http.StatusOK, nil
	}

}

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, HttpResponse{
				IsError: true,
				Message: "authorization header required",
			})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		_, err := validateJWT(tokenString)
		if err != nil {
			LogWarn("invalid jwt", "err", err)
			c.JSON(http.StatusUnauthorized, HttpResponse{
				IsError: true,
				Message: "invalid jwt",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
