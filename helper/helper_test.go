package helper

import (
	"database/sql/driver"
	"testing"

	"github.com/stretchr/testify/suite"
)

type HelperTestSuite struct {
	suite.Suite
	content  string
	prefix   string
	messages []string
}

func (suite *HelperTestSuite) TestValidateContentShouldTrue() {
	suite.content = "s.choose test, test2"
	suite.prefix = "s."

	res, ok := ValidateContent(suite.content, suite.prefix)
	suite.Equal("choose test, test2", res)
	suite.True(ok)
}

func (suite *HelperTestSuite) TestValidateContentShouldReturnFalseWhenLenContentLessThanLenPrefix() {
	suite.content = "s"
	suite.prefix = "s."

	res, ok := ValidateContent(suite.content, suite.prefix)
	suite.Equal("", res)
	suite.False(ok)
}

func (suite *HelperTestSuite) TestValidateContentShouldReturnFalseWhenPrefixNotMatch() {
	suite.content = "d.test"
	suite.prefix = "s."

	res, ok := ValidateContent(suite.content, suite.prefix)
	suite.Equal("", res)
	suite.False(ok)
}

func (suite *HelperTestSuite) TestGetRandomMessagesShouldSuccess() {
	suite.messages = []string{"test"}
	res := GetRandomMessage(suite.messages)

	suite.Equal("test", res)
}

func (suite *HelperTestSuite) TestReplaceQueryStringShouldSuccess() {
	query := "(test)"
	res := ReplaceQueryString(query)

	suite.Equal("\\(test\\)", res)
}

func (suite *HelperTestSuite) TestReplaceDiscordIDShouldSuccess() {
	id := "<@&749933518917992549>"
	res := ReplaceDiscordID(id)

	suite.Equal("749933518917992549", res)
}

func (suite *HelperTestSuite) TestTransformSliceInterfaceToSliceValueShouldSuccess() {
	suite.messages = []string{"test", "test2"}

	i := make([]interface{}, len(suite.messages))
	for key, value := range suite.messages {
		i[key] = value
	}

	res := TransformSliceInterfaceToSliceValue(i)

	suite.Equal([]driver.Value{"test", "test2"}, res)
}

func (suite *HelperTestSuite) TestGetMapKeysShouldSuccess() {
	maps := make(map[string][]string)
	maps["test"] = []string{"1"}
	res := GetMapKeys(maps)

	suite.Equal([]string{"test"}, res)
}

func TestHelperTestSuite(t *testing.T) {
	suite.Run(t, new(HelperTestSuite))
}
