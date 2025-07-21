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

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var text = tview.NewTextView()
var app = tview.NewApplication()

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "--help" || os.Args[1] == "-h") {
		help()
		return
	}

	configFile()

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'q' || event.Rune() == 'Q' {
			app.Stop()
		}
		return event
	})

	text.SetDynamicColors(true).
		SetBorder(true).
		SetTitle("DeepL Translator - Press 'q' to quit")

	if len(os.Args) < 2 {
		fmt.Print("Write a word: ")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		go runTranslation(input)
	} else {
		input := strings.Join(os.Args[1:], " ")
		go runTranslation(input)
	}

	if err := app.SetRoot(text, true).Run(); err != nil {
		log.Fatal(err)
	}
}

type DeeplResponse struct {
	Translations []Translation `json:"translations"`
}
type Translation struct {
	Text string `json:"text"`
}

func runTranslation(input string) {
	apiKey := tokenCheck()
	urlAPI := "https://api-free.deepl.com/v2/translate"

	requestBody := map[string]any{
		"text":        []string{input},
		"target_lang": config.TargetLanguage,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		showError("JSON encode error: " + err.Error())
		return
	}

	req, err := http.NewRequest("POST", urlAPI, bytes.NewBuffer(jsonData))
	if err != nil {
		showError("Request error: " + err.Error())
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "DeepL-Auth-Key "+apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		showError("HTTP error: " + err.Error())
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		showError("Read error: " + err.Error())
		return
	}

	var deeplResp DeeplResponse
	err = json.Unmarshal(body, &deeplResp)
	if err != nil {
		showError("Unmarshal error: " + err.Error())
		return
	}

	if len(deeplResp.Translations) > 0 {
		app.QueueUpdateDraw(func() {
			text.SetText(deeplResp.Translations[0].Text)
		})
	} else {
		showError("No translation found.")
	}
}

func showError(message string) {
	app.QueueUpdateDraw(func() {
		text.SetText("[red]Error: " + message)
	})
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

	data, err := io.ReadAll(jsonFile)
	Err(err)

	err = json.Unmarshal(data, &config)
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

func help() {
	fmt.Println(`
-h, --help
	for help

Just write the word and over`)
}
