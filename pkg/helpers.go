package pkg

import (
	"fmt"
	"math/rand"
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/segmentio/ksuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	Production  = "prod"
	Development = "dev"

	IDIdentifier    = "id"
	EmailIdentifier = "email"

	AccessTokenDuration  = 30 * time.Minute
	RefreshTokenDuration = 6 * time.Hour
)

var (
	passwordRegex = regexp.MustCompile(`^[A-Za-z\d!@#$%^&*(),.?":{}|<>]{8,}$`)
)

var ValidatePassword validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if password, ok := fieldLevel.Field().Interface().(string); ok {
		return passwordRegex.MatchString(password)
	}
	return false
}

func GenerateID() string {
	return ksuid.New().String()
}

// HashPassword generates a bcrypt hash of the provided password. It returns the hashed password and an error if the hashing fails.
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedPassword), nil
}

// VerifyPassword compares a plaintext password with a hashed password and returns an error if they do not match.
func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func GenerateUsername() string {
	adjectives := []string{"happy", "cool", "super", "ninja", "mega"}
	nouns := []string{"tiger", "wolf", "coder", "hero", "star"}

	rand.Seed(time.Now().UnixNano())
	adj := adjectives[rand.Intn(len(adjectives))]
	noun := nouns[rand.Intn(len(nouns))]

	num := rand.Intn(90) + 10
	return fmt.Sprintf("%s%s%d", adj, noun, num)
}
