package Controllers

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
)

var tokenName = "token"

type Claims struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Tipe     int    `json:"tipe"`
	jwt.StandardClaims
}

func generateToken(w http.ResponseWriter, email string, password string, tipe int) {
	var envErr = godotenv.Load()

	if envErr != nil {
		fmt.Println(envErr)
	}

	var jwtKey = []byte(os.Getenv("SECRETKEY"))

	tokenExpiryTime := time.Now().Add(5 * time.Minute)

	claims := &Claims{
		Email:    email,
		Password: password,
		Tipe:     tipe,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: tokenExpiryTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(jwtKey)
	if err != nil {
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    tokenName,
		Value:   signedToken,
		Expires: tokenExpiryTime,
		Path:    "/",
	})
}

func resetUserToken(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:    tokenName,
		Value:   "",
		Expires: time.Now(),
		Path:    "/",
	})
}

func Authenticate(next http.HandlerFunc, accessType int) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isValidToken := validateUserToken(r, accessType)
		if !isValidToken {
			PrintError(401, "Unathorized Access", w)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func validateUserToken(r *http.Request, accessType int) bool {
	isAccessTokenValid, _, _, tipe := validateTokenFromCookies(r)

	if isAccessTokenValid {
		isUserValid := tipe == accessType
		if isUserValid {
			return true
		}
	}
	return false
}

func validateTokenFromCookies(r *http.Request) (bool, string, string, int) {
	var envErr = godotenv.Load()

	if envErr != nil {
		fmt.Println(envErr)
	}

	var jwtKey = []byte(os.Getenv("SECRETKEY"))

	if cookie, err := r.Cookie(tokenName); err == nil {
		accessToken := cookie.Value
		accessClaims := &Claims{}
		parsedToken, err := jwt.ParseWithClaims(accessToken, accessClaims, func(accessToken *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err == nil && parsedToken.Valid {
			return true, accessClaims.Email, accessClaims.Password, accessClaims.Tipe
		}
	}
	return false, "", "", -1
}
