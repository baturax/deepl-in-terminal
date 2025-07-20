package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	tokenCheck()
}

func tokenCheck() string {
	home := os.Getenv("HOME")
	cache := filepath.Join(home, ".cache")
	tokenPath := filepath.Join(cache, "deepltoken")
	tokenByte, err := os.ReadFile(tokenPath)
	Err(err)
	token := strings.TrimSpace(string(tokenByte))
	return token
}

func Err(e error) {
	if e != nil {
		log.Fatalln(e)
	}
}
