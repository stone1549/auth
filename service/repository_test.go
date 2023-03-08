package service_test

import (
	"crypto/rsa"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/stone1549/yapyapyap/auth/service"
	"io/ioutil"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

// assert fails the test if the condition is false.
func assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		tb.FailNow()
	}
}

type configuration int

const (
	inMemoryEmpty configuration = 0
	inMemorySmall configuration = iota
	pgEmpty       configuration = iota
	pgSmall       configuration = iota
	inMemoryRsa   configuration = iota
)

func (c configuration) GetLifeCycle() service.LifeCycle {
	return service.DevLifeCycle
}

func (c configuration) GetRepoType() service.UserRepositoryType {
	switch c {
	case pgSmall:
		fallthrough
	case pgEmpty:
		return service.PostgreSqlRepo
	case inMemorySmall:
		fallthrough
	case inMemoryEmpty:
		return service.InMemoryRepo
	case inMemoryRsa:
		return service.InMemoryRepo
	default:
		return service.InMemoryRepo
	}
}

func (c configuration) GetTimeout() time.Duration {
	return 60 * time.Second
}

func (c configuration) GetPort() int {
	return 3333
}

func (c configuration) GetInitDataSet() string {
	switch c {
	case inMemoryEmpty:
		fallthrough
	case pgEmpty:
		return ""
	case inMemorySmall:
		fallthrough
	case pgSmall:
		return "../data/small_set.json"
	case inMemoryRsa:
		return ""
	default:
		return ""
	}
}

func (c configuration) GetPgUrl() string {
	switch c {
	case inMemoryEmpty:
		fallthrough
	case inMemorySmall:
		return ""
	case pgEmpty:
		fallthrough
	case pgSmall:
		return "postgres://test:test@localhost:5432/postgres?sslmode=disable"
	case inMemoryRsa:
		return ""
	default:
		return ""
	}
}

func (c configuration) GetTokenSecretKey() string {
	switch c {
	case inMemorySmall:
		fallthrough
	case inMemoryEmpty:
		fallthrough
	case pgEmpty:
		fallthrough
	case pgSmall:
		return "SECRET!"
	case inMemoryRsa:
		return ""
	default:
		return ""
	}
}

func (c configuration) GetTokenPrivateKey() *rsa.PrivateKey {
	switch c {
	case inMemorySmall:
		fallthrough
	case inMemoryEmpty:
		fallthrough
	case pgEmpty:
		fallthrough
	case pgSmall:
		return nil
	case inMemoryRsa:
		signBytes, err := ioutil.ReadFile("../data/sample.key")
		if err != nil {
			return nil
		}

		privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)
		if err != nil {
			return nil
		}

		return privateKey
	default:
		return nil
	}
}

func (c configuration) GetTokenPublicKey() *rsa.PublicKey {
	switch c {
	case inMemorySmall:
		fallthrough
	case inMemoryEmpty:
		fallthrough
	case pgEmpty:
		fallthrough
	case pgSmall:
		return nil
	case inMemoryRsa:
		signBytes, err := ioutil.ReadFile("../data/sample.pub")
		if err != nil {
			return nil
		}

		publicKey, err := jwt.ParseRSAPublicKeyFromPEM(signBytes)
		if err != nil {
			return nil
		}

		return publicKey
	default:
		return nil
	}
}

// TestNewUserRepository_ImSuccessEmpty ensures an empty in memory repo can be constructed
func TestNewUserRepository_ImSuccessEmpty(t *testing.T) {
	_, err := service.NewUserRepository(inMemoryEmpty)
	ok(t, err)
}

// TestNewUserRepository_ImSuccessSmall ensures a prepopulated memory repo can be constructed
func TestNewUserRepository_ImSuccessSmall(t *testing.T) {
	_, err := service.NewUserRepository(inMemorySmall)
	ok(t, err)
}

// TestNewUserRepository_PgSuccessEmpty ensures an empty PG repo can be constructed
func TestNewUserRepository_PgSuccessEmpty(t *testing.T) {
	_, err := service.NewUserRepository(pgEmpty)
	ok(t, err)
}
