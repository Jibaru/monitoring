package scripts

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func generateToken(userID primitive.ObjectID, userEmail string, jwtSecret []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID.Hex(),
		"email":   userEmail,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func generateUsername() string {
	var adjectives = []string{
		"Cool", "Mighty", "Silent", "Fierce", "Crazy", "Lucky", "Charming", "Nimble",
	}

	var nouns = []string{
		"Panther", "Tiger", "Falcon", "Wizard", "Knight", "Dragon", "Shadow", "Phoenix",
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	adj := adjectives[r.Intn(len(adjectives))]
	noun := nouns[r.Intn(len(nouns))]
	number := r.Intn(100)
	return fmt.Sprintf("%s%s%d", adj, noun, number)
}

func isValidPassword(encryptedPassword string, incomingPassword string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(encryptedPassword), []byte(incomingPassword)); err != nil {
		return false
	}

	return true
}

func encryptPassword(plainPassword string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashed), err
}
