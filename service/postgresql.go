package service

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/twinj/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	insertLogin       = "INSERT INTO login (id, email, username, salted_hash) VALUES ($1, $2, $3, $4)"
	authenticate      = "SELECT salted_hash, id, username FROM login WHERE email=$1"
	insertStoredLogin = "INSERT INTO login (id, email, username, salted_hash, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)"
)

type postgresqlUserRepository struct {
	db *sql.DB
}

// NewUser adds a user to the repo.
func (impr *postgresqlUserRepository) NewUser(email string, handle string, password string) (string, error) {
	if email == "" {
		return "", newErrRepository("email is required")
	} else if handle == "" {
		return "", newErrRepository("handle is required")
	} else if password == "" {
		return "", newErrRepository("password is required")
	}

	id := uuid.NewV4().String()

	saltedHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", newErrRepository("unable to generate password")
	}

	_, err = impr.db.Exec(insertLogin, id, email, handle, saltedHash)

	return id, err
}

// Authenticate compares a given email and password combination against the salted hash in the repo.
func (impr *postgresqlUserRepository) Authenticate(email string, password string) (User, error) {
	if email == "" {
		return User{}, newErrRepository("email is required")
	}

	if password == "" {
		return User{}, newErrRepository("password is required")
	}

	row := impr.db.QueryRow(authenticate, email)
	var saltedHash string
	var id string
	var username string

	err := row.Scan(&saltedHash, &id, &username)

	if err == sql.ErrNoRows {
		return User{}, nil
	} else if err != nil {
		return User{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(saltedHash), []byte(password))

	if err != nil {
		return User{}, nil
	}

	return User{id, email, username}, nil
}

func loadInitPostgresqlData(db *sql.DB, dataset string) error {
	users, err := loadInitInMemoryDataset(dataset)

	if err != nil {
		return err
	}

	txn, err := db.Begin()

	if err != nil {
		return err
	}

	for id, user := range users {
		_, err = txn.Exec(insertStoredLogin, id, user.Email, user.Username, user.SaltedHash, user.CreatedAt, user.UpdatedAt)

		if err != nil {
			return err
		}
	}

	return txn.Commit()
}

// MakePostgresqlUserRespository constructs a PostgreSQL backed UserRepository from the given params.
func MakePostgresqlUserRespository(config Configuration, db *sql.DB) (UserRepository, error) {
	var err error
	if config.GetInitDataSet() != "" {
		err = loadInitPostgresqlData(db, config.GetInitDataSet())
	}

	if err != nil {
		return nil, err
	}

	return &postgresqlUserRepository{db}, nil
}
