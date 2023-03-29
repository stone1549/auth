package service

import (
	"database/sql"
	pg "github.com/lib/pq"
	"github.com/twinj/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	getUser           = "SELECT l.email, l.username, up.gender, up.age, up.topics  FROM login l JOIN user_profile up ON (l.id=up.user_id)  WHERE l.id=$1"
	updateProfile     = "UPDATE user_profile (gender, age, topics) VALUES ($1, $2, $3) WHERE user_id=$4"
	insertLogin       = "INSERT INTO login (id, email, username, salted_hash) VALUES ($1, $2, $3, $4)"
	insertUserProfile = "INSERT INTO user_profile (user_id, gender, age, topics) VALUES ($1, $2, $3, $4)"
	authenticate      = "SELECT l.salted_hash, l.id, l.username, up.gender, up.age, up.topics  FROM login l JOIN user_profile up ON (l.id=up.user_id)  WHERE l.email=$1"
	insertStoredLogin = "INSERT INTO login (id, email, username, salted_hash, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)"
)

type postgresqlUserRepository struct {
	db *sql.DB
}

// NewUser adds a user to the repo.
func (impr *postgresqlUserRepository) NewUser(
	email string,
	handle string,
	password string,
	gender Gender,
	age int,
	topics []string,
) (string, error) {
	if email == "" {
		return "", newErrRepository("email is required")
	} else if handle == "" {
		return "", newErrRepository("handle is required")
	} else if password == "" {
		return "", newErrRepository("password is required")
	} else if gender == "" {
		return "", newErrRepository("gender is required")
	} else if age == 0 {
		return "", newErrRepository("age is required")
	} else if topics == nil {
		return "", newErrRepository("topics is required")
	}

	id := uuid.NewV4().String()

	saltedHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", newErrRepository("unable to generate password")
	}

	tx, err := impr.db.Begin()

	if err != nil {
		_ = tx.Rollback()
		return "", err
	}

	_, err = tx.Exec(insertLogin, id, email, handle, saltedHash)

	if err != nil {
		_ = tx.Rollback()
		return "", err
	}

	_, err = tx.Exec(insertUserProfile, id, gender, age, pg.Array(topics))

	if err != nil {
		_ = tx.Rollback()
		return "", err
	}

	err = tx.Commit()
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
	var gender Gender
	var age int
	var topics []string

	err := row.Scan(&saltedHash, &id, &username, &gender, &age, &topics)

	if err == sql.ErrNoRows {
		return User{}, nil
	} else if err != nil {
		return User{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(saltedHash), []byte(password))

	if err != nil {
		return User{}, nil
	}

	return User{id, email, username, UserProfile{Gender: gender, Age: age, Topics: topics}}, nil
}

func (imr *postgresqlUserRepository) GetUser(id string) (User, error) {
	if id == "" {
		return User{}, newErrRepository("id is required")
	}
	row := imr.db.QueryRow(id)
	var email string
	var username string
	var gender Gender
	var age int
	var topics []string

	err := row.Scan(&email, &username, &gender, &age, &topics)

	if err == sql.ErrNoRows {
		return User{}, nil
	} else if err != nil {
		return User{}, err
	}

	return User{id, email, username, UserProfile{Gender: gender, Age: age, Topics: topics}}, nil
}

func (impr *postgresqlUserRepository) UpdateProfile(userId string, profile UserProfile) error {
	if userId == "" {
		return newErrRepository("userId is required")
	} else if profile.Gender == "" {
		return newErrRepository("gender is required")
	} else if profile.Age == 0 {
		return newErrRepository("age is required")
	} else if profile.Topics == nil {
		return newErrRepository("topics is required")
	}

	_, err := impr.db.Exec(updateProfile, profile.Gender, profile.Age, profile.Topics, userId)

	return err
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
