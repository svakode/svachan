package repository

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/huandu/go-sqlbuilder"
	"github.com/stretchr/testify/suite"

	"github.com/svakode/svachan/dictionary"
	"github.com/svakode/svachan/helper"
)

type MemberTestSuite struct {
	suite.Suite
	username, email string
	usernames       []interface{}

	mockPostgre sqlmock.Sqlmock

	db         *sql.DB
	repository Repository
}

func (suite *MemberTestSuite) SetupTest() {
	suite.username = "@username"
	suite.usernames = []interface{}{"username", "username2"}
	suite.email = "email@gmail.com"
	suite.db, suite.mockPostgre, _ = sqlmock.New()
}

func (suite *MemberTestSuite) TestSetMemberEmailUpdateShouldSuccess() {
	defer suite.db.Close()

	queryInsert, argsInsert := prepareSetMemberEmailInsertQuery(suite.username, suite.email)
	queryUpdate, argsUpdate := prepareSetMemberEmailUpdateQuery(suite.username, suite.email)

	suite.mockPostgre.ExpectExec(queryInsert).
		WithArgs(argsInsert...).
		WillReturnResult(sqlmock.NewResult(1, 0))
	suite.mockPostgre.ExpectExec(queryUpdate).
		WithArgs(argsUpdate...).
		WillReturnResult(sqlmock.NewResult(1, 1))

	suite.repository = NewRepository(suite.db)
	err := suite.repository.SetMemberEmail(suite.username, suite.email)

	suite.Nil(err)
}

func (suite *MemberTestSuite) TestSetMemberEmailInsertShouldSuccess() {
	defer suite.db.Close()

	queryInsert, argsInsert := prepareSetMemberEmailInsertQuery(suite.username, suite.email)

	suite.mockPostgre.ExpectExec(queryInsert).
		WithArgs(argsInsert...).
		WillReturnResult(sqlmock.NewResult(1, 1))

	suite.repository = NewRepository(suite.db)
	err := suite.repository.SetMemberEmail(suite.username, suite.email)

	suite.Nil(err)
}

func (suite *MemberTestSuite) TestSetMemberEmailInsertShouldReturnError() {
	expectedErr := errors.New(dictionary.DBError)
	defer suite.db.Close()

	queryInsert, argsInsert := prepareSetMemberEmailInsertQuery(suite.username, suite.email)

	suite.mockPostgre.ExpectExec(queryInsert).
		WithArgs(argsInsert...).
		WillReturnError(errors.New("some-error"))

	suite.repository = NewRepository(suite.db)
	err := suite.repository.SetMemberEmail(suite.username, suite.email)

	suite.Equal(expectedErr, err)
}

func (suite *MemberTestSuite) TestSetMemberEmailUpdateShouldReturnError() {
	expectedErr := errors.New(dictionary.DBError)
	defer suite.db.Close()

	queryInsert, argsInsert := prepareSetMemberEmailInsertQuery(suite.username, suite.email)
	queryUpdate, argsUpdate := prepareSetMemberEmailUpdateQuery(suite.username, suite.email)

	suite.mockPostgre.ExpectExec(queryInsert).
		WithArgs(argsInsert...).
		WillReturnResult(sqlmock.NewResult(1, 0))
	suite.mockPostgre.ExpectExec(queryUpdate).
		WithArgs(argsUpdate...).
		WillReturnError(errors.New("some-error"))

	suite.repository = NewRepository(suite.db)
	err := suite.repository.SetMemberEmail(suite.username, suite.email)

	suite.Equal(expectedErr, err)
}

func (suite *MemberTestSuite) TestGetMemberEmailShouldSuccess() {
	defer suite.db.Close()

	query, args := prepareGetMemberEmailQuery(suite.username)

	suite.mockPostgre.ExpectQuery(query).
		WithArgs(args...).
		WillReturnRows(sqlmock.NewRows([]string{"email"}).AddRow(suite.email))

	suite.repository = NewRepository(suite.db)
	res, err := suite.repository.GetMemberEmail(suite.username)

	suite.Nil(err)
	suite.Equal(suite.email, res)
}

func (suite *MemberTestSuite) TestGetMemberEmailShouldReturnErrorWhenNotFound() {
	expectedErr := fmt.Errorf(dictionary.MemberNotFoundMessage)
	defer suite.db.Close()

	query, args := prepareGetMemberEmailQuery(suite.username)

	suite.mockPostgre.ExpectQuery(query).
		WithArgs(args...).
		WillReturnRows(sqlmock.NewRows([]string{"email"}))

	suite.repository = NewRepository(suite.db)
	res, err := suite.repository.GetMemberEmail(suite.username)

	suite.Empty(res)
	suite.Equal(expectedErr, err)
}

func (suite *MemberTestSuite) TestGetMemberEmailShouldReturnError() {
	expectedErr := errors.New(dictionary.DBError)
	defer suite.db.Close()

	query, args := prepareGetMemberEmailQuery(suite.username)

	suite.mockPostgre.ExpectQuery(query).
		WithArgs(args...).
		WillReturnError(errors.New("some-error"))

	suite.repository = NewRepository(suite.db)
	res, err := suite.repository.GetMemberEmail(suite.username)

	suite.Empty(res)
	suite.Equal(expectedErr, err)
}

func (suite *MemberTestSuite) TestGetMembersEmailShouldReturnErrorWhenNotFound() {
	expectedErr := fmt.Errorf(dictionary.MemberNotFoundMessage)
	defer suite.db.Close()

	query, args := prepareGetMembersEmailQuery(suite.usernames)

	suite.mockPostgre.ExpectQuery(query).
		WithArgs(args...).
		WillReturnRows(sqlmock.NewRows([]string{"username", "email"}))

	suite.repository = NewRepository(suite.db)
	_, err := suite.repository.GetMembersEmail(suite.usernames)

	suite.Equal(expectedErr, err)
}

func (suite *MemberTestSuite) TestGetMembersEmailShouldReturnErrorWhenDBError() {
	expectedErr := fmt.Errorf(dictionary.DBError)
	defer suite.db.Close()

	query, args := prepareGetMembersEmailQuery(suite.usernames)

	suite.mockPostgre.ExpectQuery(query).
		WithArgs(args...).
		WillReturnError(errors.New("some-error"))

	suite.repository = NewRepository(suite.db)
	_, err := suite.repository.GetMembersEmail(suite.usernames)

	suite.Equal(expectedErr, err)
}

func (suite *MemberTestSuite) TestGetMembersEmailShouldSuccess() {
	defer suite.db.Close()

	query, args := prepareGetMembersEmailQuery(suite.usernames)

	suite.mockPostgre.ExpectQuery(query).
		WithArgs(args...).
		WillReturnRows(sqlmock.NewRows([]string{"username", "email"}).AddRow(suite.username, suite.email))

	suite.repository = NewRepository(suite.db)
	res, err := suite.repository.GetMembersEmail(suite.usernames)

	suite.Nil(err)
	suite.Equal(suite.email, res[suite.username])
}

func prepareSetMemberEmailInsertQuery(username, email string) (string, []driver.Value) {
	ib := sqlbuilder.PostgreSQL.NewInsertBuilder()
	ib.InsertInto("members")
	ib.Cols("username", "email")
	ib.Values(username, email)

	query, args := ib.Build()
	query = helper.ReplaceQueryString(query)
	newArgs := helper.TransformSliceInterfaceToSliceValue(args)

	return query, newArgs
}

func prepareSetMemberEmailUpdateQuery(username, email string) (string, []driver.Value) {
	ub := sqlbuilder.PostgreSQL.NewUpdateBuilder()
	ub.Update("members")
	ub.Set("email", email)
	ub.Where(
		ub.Equal("username", username),
	)

	query, args := ub.Build()
	query = helper.ReplaceQueryString(query)
	newArgs := helper.TransformSliceInterfaceToSliceValue(args)

	return query, newArgs
}

func prepareGetMemberEmailQuery(username string) (string, []driver.Value) {
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
	sb.Select("email")
	sb.From("members")
	sb.Where(sb.Equal("username", username))

	query, args := sb.Build()
	query = helper.ReplaceQueryString(query)
	newArgs := helper.TransformSliceInterfaceToSliceValue(args)

	return query, newArgs
}

func prepareGetMembersEmailQuery(usernames []interface{}) (string, []driver.Value) {
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
	sb.Select("username", "email")
	sb.From("members")
	sb.Where(sb.In("username", usernames...))

	query, args := sb.Build()
	query = helper.ReplaceQueryString(query)
	newArgs := helper.TransformSliceInterfaceToSliceValue(args)

	return query, newArgs
}

func TestMemberTestSuite(t *testing.T) {
	suite.Run(t, new(MemberTestSuite))
}
