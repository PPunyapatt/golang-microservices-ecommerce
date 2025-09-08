package helper

import (
	"auth-service/v1/internal/constant"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
)

func GenerateToken(userID string, role []int32) (*string, error) {
	lifetimeString := os.Getenv("JWT_LIFETIME")
	lifetime, err := strconv.Atoi(lifetimeString)
	if err != nil {
		log.Printf("%+v", errors.WithStack(err))
		return nil, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": jwt.NewNumericDate(time.Now().Add(time.Duration(lifetime) * time.Second)),
		// "iat":  jwt.NewNumericDate(time.Now()),
		// "nbf":  jwt.NewNumericDate(time.Now()),
		"iss":  os.Getenv("JWT_ISSUER"),
		"aud":  []string{"api"},
		"role": role,
		"sub":  userID,
	})

	jwt_data, err := os.ReadFile("/vault/secrets/jwt")
	if err != nil {
		log.Fatal("Error reading jwt secret file: ", err.Error())
	}
	parts := strings.SplitN(string(jwt_data), "=", 2)
	// jwtSecret := os.Getenv("JWT_SECRET")
	signedToken, err := token.SignedString([]byte(parts[1]))
	if err != nil {
		log.Printf("%+v", errors.WithStack(err))
		return nil, err
	}

	return &signedToken, nil
}

func VerifyToken(tokenString string) (*constant.User, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(jwtSecret), nil
	})
	if err != nil {
		log.Printf("%+v", errors.WithStack(err))
		return nil, err
	}
	var userClaim *constant.User
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if claims["exp"] != nil {
			exp, err := strconv.Atoi(claims["exp"].(string))
			if err != nil {
				log.Printf("%+v", errors.WithStack(err))
				return nil, err
			}
			if exp < int(time.Now().Unix()) {
				return nil, errors.New("token expired")
			}

			userClaim = &constant.User{
				ID: claims["sub"].(string),
				// Role: claims["role"].([]int),
			}
		}
	} else {
		log.Printf("%+v", errors.WithStack(err))
		return nil, err
	}

	return userClaim, nil
}
