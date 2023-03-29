package service

import (
	"database/sql"
	"errors"
)

type errRepository struct {
	err error
}

func (er errRepository) Error() string {
	return er.err.Error()
}

func newErrRepository(msg string) error {
	return errRepository{errors.New(msg)}
}

// UserRepository represents a data source through which users can be managed.
type UserRepository interface {
	GetUser(id string) (User, error)
	UpdateProfile(userId string, profile UserProfile) error
	// NewUser adds a user to the repo.
	NewUser(email string, handle string, password string, gender Gender, age int, topics []string) (string, error)
	// Authenticate validates email and password combo with what is stored in the repo. Returns users unique id on
	// success and empty string on failure
	Authenticate(email string, password string) (User, error)
}

// NewUserRepository constructs a UserRepository from the given configuration.
func NewUserRepository(config Configuration) (UserRepository, error) {
	var err error
	var repo UserRepository
	var db *sql.DB
	switch config.GetRepoType() {
	case InMemoryRepo:
		repo, err = MakeInMemoryRepository(config)
	case PostgreSqlRepo:
		db, err = sql.Open("postgres", config.GetPgUrl())

		if err != nil {
			return nil, err
		}
		repo, err = MakePostgresqlUserRespository(config, db)
	default:
		err = newErrRepository("repository type unimplemented")
	}

	return repo, err
}
