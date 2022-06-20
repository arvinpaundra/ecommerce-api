package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// Generate token
func CreateToken(customer_id uint32) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["customer_id"] = customer_id
	claims["exp"] = time.Now().Add(time.Hour * 3).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(os.Getenv("API_SECRET")))
}

// Check token validity
func TokenValid(r *http.Request) error {
	tokenString := ExtractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("API_SECRET")), nil
	})

	if err != nil {
		return err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		Pretty(claims)
	}

	return nil
}

// Get token from query params or header authorization
func ExtractToken(r *http.Request) string {
	keys := r.URL.Query()
	token := keys.Get("token")

	if token != "" {
		return token
	}

	bearerToken := r.Header.Get("Authorization")

	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}

	return ""
}

// Extract the payload property id of customer from token
func ExtractTokenID(r *http.Request) (uint32, error) {
	tokenString := ExtractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("API_SECRET")), nil
	})

	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid {
		customer_id, err := strconv.ParseUint(fmt.Sprintf("%.0f", claims["customer_id"]), 10, 32)

		if err != nil {
			return 0, err
		}

		return uint32(customer_id), nil
	}

	return 0, nil
}

// Format display for payload claims in the terminal
func Pretty(data interface{}) {
	b, err := json.MarshalIndent(data, "", " ")

	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println(string(b))
}
