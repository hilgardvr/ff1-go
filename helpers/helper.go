package helpers

import (
	"fmt"
	"math/rand"
	"strconv"
)

func GenerateLoginCode() string {
	code := rand.Intn(100000)
	str := strconv.Itoa(code)
	padded := fmt.Sprintf("%05s", str)
	return padded
}