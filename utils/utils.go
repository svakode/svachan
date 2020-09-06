package utils

import (
	"log"
	"math/rand"
	"os"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

const CharSet = "abcdedfghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

func checkKey(key string) {
	if !viper.IsSet(key) && os.Getenv(key) == "" {
		log.Print("can't find key ", key)
		debug.PrintStack()
		log.Fatalf("%s key is not set", key)
	}
}

func GetString(key string) string {
	value := os.Getenv(key)
	if value == "" {
		value = viper.GetString(key)
	}
	return value
}

func FatalGetString(key string) string {
	checkKey(key)
	value := os.Getenv(key)
	if value == "" {
		value = viper.GetString(key)
	}
	return value
}

func GetIntOrPanic(key string) int {
	checkKey(key)
	v, err := strconv.Atoi(FatalGetString(key))
	if err != nil {
		v, err = strconv.Atoi(os.Getenv(key))
		debug.PrintStack()
		log.Fatalf("Could not parse key: %s, Error: %s", key, err)
	}
	return v
}

// Contains tells whether a string is exist in a slice
func Contains(slice []string, str string) bool {
	for _, val := range slice {
		if str == val {
			return true
		}
	}
	return false
}

// Find tells the index of a string in a slice, -1 if not found
func Find(slices []string, s string) int {
	for i, slice := range slices {
		if s == slice {
			return i
		}
	}
	return -1
}

// RemoveFromSlice will remove an element from slice
func RemoveFromSlice(slices []string, idx int) []string {
	slices[len(slices)-1], slices[idx] = slices[idx], slices[len(slices)-1]
	return slices[:len(slices)-1]
}

// GetRandomString will generate random string with a certain length
func GetRandomString(length int) string {
	var output strings.Builder

	for i := 0; i < length; i++ {
		random := rand.Intn(len(CharSet))
		randomChar := CharSet[random]
		output.WriteString(string(randomChar))
	}

	return output.String()
}
