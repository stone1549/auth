package service_test

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stone1549/yapyapyap/auth/service"
	"testing"
	"time"
)

func mockExpectExecTimes(mock sqlmock.Sqlmock, sqlRegexStr string, times int) {
	for i := 0; i < times; i++ {
		mock.ExpectExec(sqlRegexStr).WillReturnResult(sqlmock.NewResult(int64(i), 1))
	}
}

func makeAndTestPgSmallRepo() (*sql.DB, sqlmock.Sqlmock, service.UserRepository, error) {
	var err error
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, nil, err
	}

	mock.ExpectBegin()
	mockExpectExecTimes(mock, "INSERT INTO login", 5)
	mock.ExpectCommit()
	repo, err := service.MakePostgresqlUserRespository(pgSmall, db)

	return db, mock, repo, err
}

// TestMakePostgresqlUserRespository_Ds ensures that a dataset can be loaded when a pg repo is constructed.
func TestMakePostgresqlUserRespository_Ds(t *testing.T) {
	db, mock, _, err := makeAndTestPgSmallRepo()
	defer db.Close()
	ok(t, err)
	ok(t, mock.ExpectationsWereMet())
}

// TestMakePostgresqlUserRespository ensures that an empty pg repo can be constructed.
func TestMakePostgresqlUserRespository(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	_, err = service.MakePostgresqlUserRespository(pgEmpty, db)
	ok(t, err)
	ok(t, mock.ExpectationsWereMet())
}

func getProductColumns() []string {
	columns := make([]string, 0)
	columns = append(columns, "id")
	columns = append(columns, "email")
	columns = append(columns, "salted_hash")
	columns = append(columns, "created_at")
	columns = append(columns, "updated_at")
	return columns
}

func addExpectedUserId1Row(rows *sqlmock.Rows) *sqlmock.Rows {
	createdAt, _ := time.Parse("2006-01-15T15:20:59", "2017-01-01T00:00:00Z")
	updatedAt, _ := time.Parse("2006-01-15T15:20:59", "2018-01-01T00:00:20Z")
	return rows.AddRow(
		"1",
		"user@justinstone.net",
		"$2a$10$EAd78WN5SS0uUP7AlEIs..fKXBs8Zj.bZrdXNVQEn.fx72MJGwXmy",
		createdAt,
		updatedAt,
	)
}

func addExpectedUserId2Row(rows *sqlmock.Rows) *sqlmock.Rows {
	createdAt, _ := time.Parse("2006-01-15T15:20:59", "2017-01-01T00:00:01Z")
	updatedAt, _ := time.Parse("2006-01-15T15:20:59", "2018-01-01T00:00:19Z")
	return rows.AddRow(
		"2",
		"user2@justinstone.net",
		"$2a$10$EMnu7T8Nhs2o3yVv.U73mO0HY5MdjYyHFpHU8kjLUWLYOGvJJ1RrG",
		createdAt,
		updatedAt,
	)
}

func addExpectedUserId3Row(rows *sqlmock.Rows) *sqlmock.Rows {
	createdAt, _ := time.Parse("2006-01-15T15:20:59", "2017-01-01T00:00:02Z")
	updatedAt, _ := time.Parse("2006-01-15T15:20:59", "2018-01-01T00:00:18Z")
	return rows.AddRow(
		"3",
		"user3@justinstone.net",
		"$2a$10$XdLPi4OEE1HhSgdtZekQAu.0P0W.sPn4KcojbRZr2hOfkBDFSQI0a",
		createdAt,
		updatedAt,
	)
}

func addExpectedUserId4Row(rows *sqlmock.Rows) *sqlmock.Rows {
	createdAt, _ := time.Parse("2006-01-15T15:20:59", "2017-01-01T00:00:03Z")
	updatedAt, _ := time.Parse("2006-01-15T15:20:59", "2018-01-01T00:00:17Z")
	return rows.AddRow(
		"4",
		"user4@justinstone.net",
		"$2a$10$wjp0yIJYZ0FP/DhHbQ3kwugsPEklf15zw/oWx.0sTyjoSDIsXaF7a",
		createdAt,
		updatedAt,
	)
}

func addExpectedUserId5Row(rows *sqlmock.Rows) *sqlmock.Rows {
	createdAt, _ := time.Parse("2006-01-15T15:20:59", "2017-01-01T00:00:04Z")
	updatedAt, _ := time.Parse("2006-01-15T15:20:59", "2018-01-01T00:00:16Z")
	return rows.AddRow(
		"5",
		"user5@justinstone.net",
		"$2a$10$rg4W8bMMqjvh4Q9AbA89qOvdh40P8tfsIhI9Dvv.IPGKcFxlwxn6C",
		createdAt,
		updatedAt,
	)
}
