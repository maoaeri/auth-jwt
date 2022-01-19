package helper

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func RespondError(w http.ResponseWriter, code int, message string) {
	RespondJSON(w, code, map[string]string{"error": message})
}

func RespondJSON(w http.ResponseWriter, code int, payload interface{}) {
	res, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(res)
}

func Response(w http.ResponseWriter, message string) {
	w.Write([]byte(message))
}

//generate hash from password
func GenerateHash(password string) string {
	hashpass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		log.Fatalln("Error in generating password")
	}
	hashstring := fmt.Sprintf("%s", hashpass)
	return hashstring
}

//compare password and hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

//Generate token string
func GenerateTokenPairString(email, role string) map[string]string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"authorized": true,
		"role":       role,
		"email":      email,
		"exp":        time.Now().Add(time.Minute * 15).Unix(),
	})

	tokenString, err := token.SignedString([]byte(GetEnvVar("SECRETKEY"))) //key type is byte
	if err != nil {
		fmt.Print("Error in generating token.")
		log.Fatalln(err)
	}

	refreshtoken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * 7 * 24).Unix(),
	})

	rtokenString, err := refreshtoken.SignedString([]byte(GetEnvVar("SECRETKEY")))
	if err != nil {
		fmt.Print("Error in generating token.")
		log.Fatalln(err)
	}
	return map[string]string{
		"token":        tokenString,
		"refreshtoken": rtokenString,
	}
}

func ExtractEmail(r *http.Request) (string, error) {
	cookie, err := r.Cookie("token")
	var tokenString string
	if err != nil {
		return "", err
	} else {
		tokenString = cookie.Value
	}

	if tokenString == "" {
		var err error
		err = errors.New("There is no token :(")
		return "", err
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		jwtTokenSecret := GetEnvVar("SECRETKEY")

		return []byte(jwtTokenSecret), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		email := fmt.Sprintf("%s", claims["email"]) //convert to string
		return email, nil
	} else {
		return "", nil
	}
}

//get environment variables
func GetEnvVar(name string) string {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	value := os.Getenv(name)
	return value
}
