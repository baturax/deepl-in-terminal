package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var config Config

func main() {
	configFile()
	translate()
}

func translate() {
	apiKey := tokenCheck()
	url := "https://api-free.deepl.com/v2/translate"

	requestBody := map[string]any{
		"text":        []string{getInput()},
		"target_lang": config.TargetLanguage,
	}
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Println("JSON encode error:", err)
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	Err(err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "DeepL-Auth-Key "+apiKey)

	resp, err := http.DefaultClient.Do(req)
	Err(err)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	Err(err)

	fmt.Println("Response:")
	fmt.Println(string(body))
}

func getInput() string {

	args := os.Args[1:]
	return strings.Join(args, " ")
}

type Config struct {
	TargetLanguage string `json:"target_language"`
}

func configFile() string {
	home := os.Getenv("HOME")
	configPath := filepath.Join(home, ".config", "deepl-translator", "config.json")
	jsonFile, err := os.Open(configPath)
	Err(err)

	defer jsonFile.Close()

	fl, err := io.ReadAll(jsonFile)
	Err(err)
	json.Unmarshal(fl, &config)

	return config.TargetLanguage

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
