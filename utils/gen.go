package utils

import (
	"fmt"
	"math/rand"
	"strconv"
)

var Encoding string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-"
var CodeEncoding string = "abcdefghijklmnopqrstuvwxyz0123456789"

func GenID(n int) string {
	var ID string

	for i := 0; i < n; i++ {
		ID += string(Encoding[rand.Intn(64)])
	}

	return ID
}

func GenCode(n int) string {
	var code string

	for i := 0; i < n; i++ {
		code += strconv.Itoa(rand.Intn(10))
	}
	fmt.Println(code)

	return code
}
