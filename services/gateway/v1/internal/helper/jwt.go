package helper

import (
	"fmt"
	"gateway/v1/internal/constant"
	"log"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
)

func VerifyToken(tokenString string) (*constant.User, error) {
	// jwtSecret := os.Getenv("JWT_SECRET")
	jwt_data, err := os.ReadFile("/vault/secrets/jwt")
	if err != nil {
		log.Fatal("Error reading jwt secret file: ", err.Error())
	}
	parts := strings.SplitN(string(jwt_data), "=", 2)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(parts[1]), nil
	})
	if err != nil {
		log.Printf("%+v", errors.WithStack(err))
		return nil, err
	}
	var userClaim *constant.User
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		roleInterface, ok := claims["role"]
		if !ok {
			return nil, fmt.Errorf("role claim not found")
		}

		// Convert []interface{} to []int32
		roleSlice, ok := roleInterface.([]interface{})
		if !ok {
			return nil, fmt.Errorf("role claim is not an array")
		}

		roles := make([]int32, len(roleSlice))

		for _, val := range roleSlice {
			switch val := val.(type) {
			case float64:
				roles = append(roles, int32(val))
			case int:
				roles = append(roles, int32(val))
			case int32:
				roles = append(roles, val)
			default:
				return nil, fmt.Errorf("invalid role value type: %T", val)
			}
		}

		userClaim = &constant.User{
			ID:    claims["sub"].(string),
			Roles: roles,
		}
	} else {
		log.Printf("%+v :--- ", errors.WithStack(err))
		return nil, err
	}

	return userClaim, nil
}
