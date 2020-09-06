package helper

import (
	"database/sql/driver"
	"math/rand"
	"strings"
)

// ValidateContent will validate the content
func ValidateContent(content, prefix string) (string, bool) {
	if len(content) <= len(prefix) {
		return "", false
	}
	if strings.ToLower(content[:len(prefix)]) != prefix {
		return "", false
	}

	return content[len(prefix):], true
}

// GetRandomMessage will take a list of messages and it will return one random message
func GetRandomMessage(messages []string) string {
	result := rand.Intn(len(messages))
	return messages[result]
}

// ReplaceQueryString is the helper for making a query string
func ReplaceQueryString(s string) string {
	replacer := strings.NewReplacer(`(`, `\(`,
		`)`, `\)`,
		`$`, `\$`)

	return replacer.Replace(s)
}

// ReplaceDiscordID is the helper for making a query string
func ReplaceDiscordID(s string) string {
	replacer := strings.NewReplacer(`<`, ``,
		`&`, ``,
		`>`, ``,
		`@`, ``,
		`!`, ``)

	return replacer.Replace(s)
}

// TransformSliceInterfaceToSliceValue is the helper for making an argument for DB query
func TransformSliceInterfaceToSliceValue(i []interface{}) (res []driver.Value) {
	for _, v := range i {
		res = append(res, v)
	}

	return
}

// GetMapKeys will return all key in a map
func GetMapKeys(slices map[string][]string) (res []string) {
	for key := range slices {
		res = append(res, key)
	}
	return
}
