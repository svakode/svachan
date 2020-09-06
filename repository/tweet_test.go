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

type TweetTestSuite struct {
	suite.Suite
	twitterID, username, channelID string

	mockPostgre sqlmock.Sqlmock

	db         *sql.DB
	repository Repository
}

func (suite *TweetTestSuite) SetupTest() {
	suite.twitterID = "twitter-id"
	suite.username = "@username"
	suite.channelID = "channel-id"
	suite.db, suite.mockPostgre, _ = sqlmock.New()
}

func (suite *TweetTestSuite) TestAddStreamShouldSuccess() {
	defer suite.db.Close()

	queryInsert, argsInsert := prepareAddStreamQuery(suite.twitterID, suite.username, suite.channelID)

	suite.mockPostgre.ExpectExec(queryInsert).
		WithArgs(argsInsert...).
		WillReturnResult(sqlmock.NewResult(1, 1))

	suite.repository = NewRepository(suite.db)
	err := suite.repository.AddStream(suite.twitterID, suite.username, suite.channelID)

	suite.Nil(err)
}

func (suite *TweetTestSuite) TestAddStreamShouldReturnErrorWhenExecError() {
	defer suite.db.Close()

	expectedErr := errors.New(dictionary.DBError)

	queryInsert, argsInsert := prepareAddStreamQuery(suite.twitterID, suite.username, suite.channelID)

	suite.mockPostgre.ExpectExec(queryInsert).
		WithArgs(argsInsert...).
		WillReturnError(errors.New("some-error"))

	suite.repository = NewRepository(suite.db)
	err := suite.repository.AddStream(suite.twitterID, suite.username, suite.channelID)

	suite.Equal(expectedErr, err)
}

func (suite *TweetTestSuite) TestAddStreamShouldReturnErrorWhenNoRowsAffected() {
	defer suite.db.Close()

	expectedErr := fmt.Errorf(dictionary.AlreadyStreamingError, suite.twitterID)

	queryInsert, argsInsert := prepareAddStreamQuery(suite.twitterID, suite.username, suite.channelID)

	suite.mockPostgre.ExpectExec(queryInsert).
		WithArgs(argsInsert...).
		WillReturnResult(sqlmock.NewResult(1, 0))

	suite.repository = NewRepository(suite.db)
	err := suite.repository.AddStream(suite.twitterID, suite.username, suite.channelID)

	suite.Equal(expectedErr, err)
}

func (suite *TweetTestSuite) TestRemoveStreamShouldSuccess() {
	defer suite.db.Close()

	queryDelete, argsDelete := prepareDeleteStreamQuery(suite.twitterID, suite.channelID)

	suite.mockPostgre.ExpectExec(queryDelete).
		WithArgs(argsDelete...).
		WillReturnResult(sqlmock.NewResult(1, 1))

	suite.repository = NewRepository(suite.db)
	err := suite.repository.RemoveStream(suite.twitterID, suite.channelID)

	suite.Nil(err)
}

func (suite *TweetTestSuite) TestRemoveStreamShouldReturnErrorWhenExecError() {
	defer suite.db.Close()

	expectedErr := errors.New(dictionary.DBError)

	queryDelete, argsDelete := prepareDeleteStreamQuery(suite.twitterID, suite.channelID)

	suite.mockPostgre.ExpectExec(queryDelete).
		WithArgs(argsDelete...).
		WillReturnError(errors.New(dictionary.DBError))

	suite.repository = NewRepository(suite.db)
	err := suite.repository.RemoveStream(suite.twitterID, suite.channelID)

	suite.Equal(expectedErr, err)
}

func (suite *TweetTestSuite) TestRemoveStreamShouldReturnErrorWhenNoRowsAffected() {
	defer suite.db.Close()

	expectedErr := fmt.Errorf(dictionary.StreamingNotFoundError, suite.twitterID)

	queryDelete, argsDelete := prepareDeleteStreamQuery(suite.twitterID, suite.channelID)

	suite.mockPostgre.ExpectExec(queryDelete).
		WithArgs(argsDelete...).
		WillReturnResult(sqlmock.NewResult(1, 0))

	suite.repository = NewRepository(suite.db)
	err := suite.repository.RemoveStream(suite.twitterID, suite.channelID)

	suite.Equal(expectedErr, err)
}

func (suite *TweetTestSuite) TestGetStreamsShouldSuccess() {
	defer suite.db.Close()

	expectedRes := make(map[string][]string)
	expectedRes[suite.twitterID] = []string{suite.channelID}

	querySelect, argsSelect := prepareGetStreamsQuery()

	suite.mockPostgre.ExpectQuery(querySelect).
		WithArgs(argsSelect...).
		WillReturnRows(sqlmock.NewRows([]string{"twitter_id", "channel_id"}).AddRow(suite.twitterID, suite.channelID))

	suite.repository = NewRepository(suite.db)
	res, err := suite.repository.GetStreams()

	suite.Equal(expectedRes, res)
	suite.Nil(err)
}

func (suite *TweetTestSuite) TestGetStreamsShouldReturnErrorWhenExecError() {
	defer suite.db.Close()

	expectedErr := errors.New(dictionary.DBError)
	querySelect, argsSelect := prepareGetStreamsQuery()

	suite.mockPostgre.ExpectQuery(querySelect).
		WithArgs(argsSelect...).
		WillReturnError(errors.New("some-error"))

	suite.repository = NewRepository(suite.db)
	res, err := suite.repository.GetStreams()

	suite.Nil(res)
	suite.Equal(expectedErr, err)
}

func (suite *TweetTestSuite) TestGetStreamsByChannelShouldSuccess() {
	defer suite.db.Close()

	expectedRes := []string{suite.username}

	querySelect, argsSelect := prepareGetStreamsByChannelQuery(suite.channelID)

	suite.mockPostgre.ExpectQuery(querySelect).
		WithArgs(argsSelect...).
		WillReturnRows(sqlmock.NewRows([]string{"twitter_username"}).AddRow(suite.username))

	suite.repository = NewRepository(suite.db)
	res, err := suite.repository.GetStreamsByChannel(suite.channelID)

	suite.Equal(expectedRes, res)
	suite.Nil(err)
}

func (suite *TweetTestSuite) TestGetStreamsByChannelShouldReturnErrorWhenExecError() {
	defer suite.db.Close()

	expectedErr := errors.New(dictionary.DBError)
	querySelect, argsSelect := prepareGetStreamsByChannelQuery(suite.channelID)

	suite.mockPostgre.ExpectQuery(querySelect).
		WithArgs(argsSelect...).
		WillReturnError(errors.New("some-error"))

	suite.repository = NewRepository(suite.db)
	res, err := suite.repository.GetStreamsByChannel(suite.channelID)

	suite.Nil(res)
	suite.Equal(expectedErr, err)
}

func prepareGetStreamsQuery() (string, []driver.Value) {
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
	sb.Select("twitter_id", "channel_id")
	sb.From("twitter_streams")

	querySelect, argsSelect := sb.Build()
	querySelect = helper.ReplaceQueryString(querySelect)
	newArgsSelect := helper.TransformSliceInterfaceToSliceValue(argsSelect)

	return querySelect, newArgsSelect
}

func prepareGetStreamsByChannelQuery(channelID string) (string, []driver.Value) {
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
	sb.Select("twitter_username")
	sb.From("twitter_streams")
	sb.Where(
		sb.Equal("channel_id", channelID),
	)

	querySelect, argsSelect := sb.Build()
	querySelect = helper.ReplaceQueryString(querySelect)
	newArgsSelect := helper.TransformSliceInterfaceToSliceValue(argsSelect)

	return querySelect, newArgsSelect
}

func prepareAddStreamQuery(twitterID, twitterUsername, channelID string) (string, []driver.Value) {
	ib := sqlbuilder.PostgreSQL.NewInsertBuilder()
	ib.InsertInto("twitter_streams")
	ib.Cols("twitter_id", "twitter_username", "channel_id")
	ib.Values(twitterID, twitterUsername, channelID)

	queryInsert, argsInsert := ib.Build()
	queryInsert = helper.ReplaceQueryString(queryInsert)
	newArgsInsert := helper.TransformSliceInterfaceToSliceValue(argsInsert)

	return queryInsert, newArgsInsert
}

func prepareDeleteStreamQuery(twitterID, channelID string) (string, []driver.Value) {
	delB := sqlbuilder.PostgreSQL.NewDeleteBuilder()
	delB.DeleteFrom("twitter_streams")
	delB.Where(
		delB.Equal("twitter_id", twitterID),
		delB.Equal("channel_id", channelID),
	)

	query, args := delB.Build()
	query = helper.ReplaceQueryString(query)
	newArgs := helper.TransformSliceInterfaceToSliceValue(args)

	return query, newArgs
}

func TestTweetTestSuite(t *testing.T) {
	suite.Run(t, new(TweetTestSuite))
}
