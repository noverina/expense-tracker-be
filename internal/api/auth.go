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

type Client struct {
	Identifier string `json:"identifier" bson:"_id"`
	SecretKey  string `json:"secret_key" bson:"secret_key"`
	Token      string `json:"token" bson:"token"`
	Role       string `json:"role" bson:"role"`
	Exp        int64  `json:"exp" bson:"exp"`
}

type Claims struct {
	jwt.RegisteredClaims
	Role string
}

var authColl *mongo.Collection

func InitAuth() {
	collName := os.Getenv("AUTH_COLL")
	authColl = GetDB().Collection(collName)
}

func findById(c *gin.Context, id string) (*Client, error) {
	filter := bson.M{"_id": id}

	var result Client
	err := authColl.FindOne(c, filter).Decode(&result)
	if err != nil {
		LogError("query failed", "err", err)
		return nil, err
	}

	return &result, nil
}

func loadPrivateKey() (*rsa.PrivateKey, error) {
	filePath := os.Getenv("PRIVATE_KEY")

	keyData, err := os.ReadFile(filePath)
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

	keyData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	key, err := jwt.ParseRSAPublicKeyFromPEM(keyData)
	if err != nil {
		return nil, err
	}

	return key, nil
}

func validateJWT(tokenString string) (*jwt.Token, *Claims, error) {
	publicKey, err := loadPublicKey()
	if err != nil {
		return nil, nil, err
	}
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("invalid signing method", token.Header["alg"])
		}
		return publicKey, nil
	})
	if err != nil {
		return nil, nil, err
	}

	if !token.Valid {
		return nil, nil, ErrInvalidToken
	}

	return token, claims, nil
}

func generateJWT(identifier string, role string) (string, error) {
	privateKey, err := loadPrivateKey()
	if err != nil {
		return "", err
	}

	expiry, err := strconv.Atoi(os.Getenv("TOKEN_EXPIRY"))
	if err != nil {
		return "", err
	}
	exp := time.Second * time.Duration(expiry)

	claims := jwt.MapClaims{
		"iss":  "expense-tracker-be",
		"aud":  identifier,
		"role": role,
		"iat":  time.Now().Unix(),
		"exp":  time.Now().Add(exp).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func GenerateToken(c *gin.Context, identifier string, secretKey string) (string, int, error) {
	filter := bson.M{"_id": identifier}

	var result Client
	if err := authColl.FindOne(c, filter).Decode(&result); err != nil && err != mongo.ErrNoDocuments {
		LogError("query failed", "err", err)
		return "", http.StatusInternalServerError, err
	} else if err == mongo.ErrNoDocuments {
		return "", http.StatusUnauthorized, fmt.Errorf("invalid credential")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(result.SecretKey), []byte(secretKey)); err != nil {
		return "", http.StatusUnauthorized, fmt.Errorf("invalid credential")
	}

	token, err := generateJWT(identifier, result.Role)
	if err != nil {
		LogError("unable to generate jwt", "err", err)
		return "", http.StatusInternalServerError, err
	}

	expiry, err := strconv.Atoi(os.Getenv("TOKEN_EXPIRY"))
	if err != nil {
		return "", http.StatusInternalServerError, err
	}
	exp := time.Second * time.Duration(expiry)
	entity := bson.M{
		"$set": bson.M{
			"token": token,
			"exp":   time.Now().Add(exp).Unix(),
		},
	}

	if _, err := authColl.UpdateByID(c, result.Identifier, entity); err != nil {
		LogError("unable to update", "err", err, "id", result.Identifier)
		return "", http.StatusInternalServerError, err
	}

	return token, http.StatusOK, nil
}

func InvalidateToken(c *gin.Context, token string) (int, error) {
	filter := bson.M{"token": token}

	var result Client
	if err := authColl.FindOne(c, filter).Decode(&result); err != nil && err != mongo.ErrNoDocuments {
		LogError("query failed", "err", err)
		return http.StatusInternalServerError, err
	} else if err == mongo.ErrNoDocuments {
		return http.StatusUnauthorized, fmt.Errorf("invalid token")
	}

	entity := bson.M{
		"$set": bson.M{
			"token": "token has been invalidated. please generate a new one",
			"exp":   0,
		},
	}

	if _, err := authColl.UpdateByID(c, result.Identifier, entity); err != nil {
		LogError("unable to update", "err", err, "id", result.Identifier)
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
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

		_, claims, err := validateJWT(tokenString)
		if err != nil {
			LogWarn("jwt auth error", "err", err)
			c.JSON(http.StatusUnauthorized, HttpResponse{
				IsError: true,
				Message: "invalid jwt",
			})
			c.Abort()
			return
		}

		aud := claims.Audience
		dbToken, err := findById(c, aud[0])
		if tokenString != dbToken.Token || err != nil {
			if err == nil {
				err = fmt.Errorf("token doesn't match value in db")
			}
			LogWarn("jwt auth error", "err", err)
			c.JSON(http.StatusUnauthorized, HttpResponse{
				IsError: true,
				Message: "invalid jwt",
			})
			c.Abort()
			return
		}

		c.Set("role", claims.Role)
		c.Next()
	}
}

func RoleAuthMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || (role != requiredRole && role != "admin") {
			LogWarn("jwt auth error", "err", "no permission")
			c.JSON(http.StatusForbidden, HttpResponse{
				IsError: true,
				Message: "you don't have permission to access this resource",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
