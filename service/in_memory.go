package service

import (
	"encoding/json"
	"github.com/twinj/uuid"
	"golang.org/x/crypto/bcrypt"
	"os"
	"time"
)

type storedUser struct {
	User
	Id         string    `json:"id"`
	SaltedHash string    `json:"saltedHash"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}
type inMemoryUserRepository struct {
	usersByEmail map[string]*storedUser
}

// NewUser adds a user to the repo.
func (imr *inMemoryUserRepository) NewUser(email string, handle string, password string) (string, error) {
	if email == "" {
		return "", newErrRepository("email is required")
	} else if handle == "" {
		return "", newErrRepository("handle is required")
	} else if password == "" {
		return "", newErrRepository("password is required")
	}

	id := uuid.NewV4().String()

	_, ok := imr.usersByEmail[email]
	if ok {
		return "", newErrRepository("user already exists")
	}

	saltedHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", newErrRepository("unable to generate password")
	}

	createdAt := time.Now()
	updatedAt := createdAt

	imr.usersByEmail[email] = &storedUser{User{Email: email, Username: handle}, id, string(saltedHash),
		createdAt, updatedAt}

	return id, nil
}

// Authenticate compares the given email and password combination against the salted hash in the repo.
func (imr *inMemoryUserRepository) Authenticate(email string, password string) (User, error) {
	if email == "" {
		return User{}, newErrRepository("email is required")
	} else if password == "" {
		return User{}, newErrRepository("password is required")
	}

	user, ok := imr.usersByEmail[email]
	if !ok {
		return User{}, newErrRepository("user not found")
	}

	if bcrypt.CompareHashAndPassword([]byte(user.SaltedHash), []byte(password)) != nil {
		return User{}, nil
	}

	return User{user.Id, user.Email, user.Username}, nil
}

// MakeInMemoryRepository constructs an in memory backed UserRepository from the given configuration.
func MakeInMemoryRepository(config Configuration) (UserRepository, error) {
	var err error

	usersByEmail, err := loadInitInMemoryDataset(config.GetInitDataSet())

	return &inMemoryUserRepository{usersByEmail}, err
}

func loadInitInMemoryDataset(dataset string) (map[string]*storedUser, error) {
	if dataset == "" {
		return make(map[string]*storedUser), nil
	}

	var err error
	storedUsers := make([]storedUser, 0)

	if err != nil {
		return nil, err
	}

	jsonBytes, err := os.ReadFile(dataset)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(jsonBytes, &storedUsers)

	if err != nil {
		return nil, err
	}

	usersByEmail := make(map[string]*storedUser)

	for index, storedUser := range storedUsers {
		usersByEmail[storedUser.Email] = &storedUsers[index]
	}

	return usersByEmail, err
}
