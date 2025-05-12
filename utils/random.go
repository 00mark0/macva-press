package utils

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// Predefined list of sensible categories
var predefinedCategories = []string{
	"Sports",
	"Technology",
	"Health",
	"Education",
	"Finance",
	"Entertainment",
	"Travel",
	"Food",
	"Fashion",
	"Science",
}

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

// RandomInt generate a random integer between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomString generates a random string of length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// RandomUser generates a random username
func RandomUser() string {
	return RandomString(6)
}

// RandomAmount generates a random amount
func RandomAmount() int64 {
	return RandomInt(0, 1000)
}

func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(6))
}

// RandomCategory generates a random category name from the predefined list
func RandomCategory() string {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	return predefinedCategories[rand.Intn(len(predefinedCategories))]
}

// RandomCategoryList generates a shuffled list of unique categories
func RandomCategoryList(n int) []string {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	shuffled := make([]string, len(predefinedCategories))
	copy(shuffled, predefinedCategories)
	rand.Shuffle(len(shuffled), func(i, j int) { shuffled[i], shuffled[j] = shuffled[j], shuffled[i] })

	if n > len(shuffled) {
		return shuffled
	}
	return shuffled[:n]
}

// RandomPlacement returns a random ad placement between "footer", "side", and "top"
func RandomPlacement() string {
	placements := []string{"footer", "side", "top"}
	return placements[rand.Intn(len(placements))]
}
