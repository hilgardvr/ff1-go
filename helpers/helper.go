package helpers

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

func GenerateLoginCode() string {
	rand.Seed(time.Now().UnixMicro())
	code := rand.Intn(100000)
	str := strconv.Itoa(code)
	padded := fmt.Sprintf("%05s", str)
	return padded
}