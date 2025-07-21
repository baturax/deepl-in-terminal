package main

import (
	"bufio"
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

func main() {
	configFile()
	handleInput()
}

func handleInput() {
	if len(os.Args) < 2 {
		fmt.Print("Write a word: ")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		translate(input)
	} else if os.Args[1] == "--help" || os.Args[1] == "-h" {
		help()
	} else {
		input := strings.Join(os.Args[1:], " ")
		translate(input)
	}
}

func help() {
	fmt.Println(`Usage:
   --help, -h    Show this help
	 config is in: $HOME/.config/deepl-translator/config.json
	 example json: 
{
  "target_language": "TR"
}
	`)
}

type DeeplResponse struct {
	Translations []Translation `json:"translations"`
}

type Translation struct {
	Text string `json:"text"`
}

func translate(input string) {
	apiKey := tokenCheck()
	urlAPI := "https://api-free.deepl.com/v2/translate"

	requestBody := map[string]any{
		"text":        []string{input},
		"target_lang": config.TargetLanguage,
	}

	jsonData, err := json.Marshal(requestBody)
	Err(err)

	req, err := http.NewRequest("POST", urlAPI, bytes.NewBuffer(jsonData))
	Err(err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "DeepL-Auth-Key "+apiKey)

	resp, err := http.DefaultClient.Do(req)
	Err(err)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	Err(err)

	var deeplResp DeeplResponse
	err = json.Unmarshal(body, &deeplResp)
	Err(err)

	if len(deeplResp.Translations) > 0 {
		fmt.Println(deeplResp.Translations[0].Text)
	} else {
		fmt.Println("Çeviri alınamadı.")
	}
}

type Config struct {
	TargetLanguage string `json:"target_language"`
}

var config Config

func configFile() string {
	home := os.Getenv("HOME")
	configPath := filepath.Join(home, ".config", "deepl-translator", "config.json")
	jsonFile, err := os.Open(configPath)
	Err(err)
	defer jsonFile.Close()

	fl, err := io.ReadAll(jsonFile)
	Err(err)

	err = json.Unmarshal(fl, &config)
	Err(err)

	return config.TargetLanguage
}

func tokenCheck() string {
	home := os.Getenv("HOME")
	tokenPath := filepath.Join(home, ".cache", "deepltoken")
	tokenByte, err := os.ReadFile(tokenPath)
	Err(err)
	return strings.TrimSpace(string(tokenByte))
}

func Err(e error) {
	if e != nil {
		log.Fatalln(e)
	}
}
