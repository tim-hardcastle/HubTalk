package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/lmorg/readline"
)

const (
	VERSION = "0.1"
	RED     = "\033[31m"
	RESET   = "\033[0m"
	GREEN   = "\033[32m"
)

var (
	rline              *readline.Instance
	host               string
	username, password string
	OK                 = Green("ok")
)

type jsonRequest = struct {
	Body     string
	Username string
	Password string
}

type jsonResponse = struct {
	Body    string
	Service string
}

func main() {
	fmt.Print(Logo())

	rline = readline.NewInstance()
	rline.SetPrompt("ht → ")

	for {
		line, _ := rline.Readline()
		if line == "ht config" {
			config()
			continue
		}

		if line == "ht quit" {
			break
		}

		// Otherwise, we pass the line to the hub and see what we get back.

		do(line)

	}
}

func do(line string) {

	jRq := jsonRequest{Body: line, Username: username, Password: password}

	body, _ := json.Marshal(jRq)

	request, err := http.NewRequest("POST", host, bytes.NewBuffer(body))

	if err != nil {
		fmt.Println(err.Error())
	}

	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer response.Body.Close()

	rBody, _ := io.ReadAll(response.Body)
	var jRsp jsonResponse
	json.Unmarshal(rBody, &jRsp)
	if jRsp.Service == "" {
		rline.SetPrompt("→ ")
	} else {
		rline.SetPrompt(jRsp.Service + " → ")
	}
	fmt.Print(jRsp.Body)
}

func config() {
	rline.SetPrompt("Host: ")
	host, _ = rline.Readline()
	rline.SetPrompt("Username: ")
	username, _ = rline.Readline()
	rline.SetPrompt("Password: ")
	rline.PasswordMask = '▪'
	password, _ = rline.Readline()
	rline.PasswordMask = 0
	rline.SetPrompt("ht → ")
	fmt.Println(OK)
	do("hub services")
}

func Logo() string {
	var padding string
	if len(VERSION)%2 == 0 {
		padding = ","
	}
	titleText := " HubTalk" + padding + " version " + VERSION + " "
	loveHeart := Red("♥")
	leftMargin := "  "
	bar := strings.Repeat("═", len(titleText)/2)
	logoString := "\n" +
		leftMargin + "╔" + bar + loveHeart + bar + "╗\n" +
		leftMargin + "║" + titleText + "║\n" +
		leftMargin + "╚" + bar + loveHeart + bar + "╝\n\n"
	return logoString
}

func Red(s string) string {
	return RED + s + RESET
}

func Green(s string) string {
	return GREEN + s + RESET
}
