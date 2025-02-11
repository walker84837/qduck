package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// initial information sent by duck.ai
type ChatInformation struct {
	Role    string `json:"role"`
	Message string `json:"message"`
	Created int    `json:"created"`
	Id      string `json:"id"`
	Action  string `json:"action"`
	Model   string `json:"model"`
}

// actual message fragments from duck.ai
type MessageFragment struct {
	Message string `json:"message"`
	Created int    `json:"created"`
	Id      string `json:"id"`
	Action  string `json:"action"`
	Model   string `json:"model"`
}

type Args struct {
	Prompt string
	Model  string
}

func getVqdToken() (string, error) {
	req, err := http.NewRequest("GET", "https://duckduckgo.com/duckchat/v1/status", nil)

	if err != nil {
		return "", fmt.Errorf("failed to make a new HTTP request: %s", err)
	}

	req.Header.Set("X-Vqd-Accept", "1")
	req.Header.Set("Referer", "https://duckduckgo.com/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.6778.86 Safari/537.36")

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		return "", fmt.Errorf("failed to act on request: %s", err)
	}

	defer resp.Body.Close()

	vqd := resp.Header.Get("x-vqd-4")

	return vqd, nil
}

func prompt(input string) (string, error) {
	var vqd string = getVqdToken()

	bodystring := `{"model": "gpt-4o-mini", "messages": [{"role": "user","content": "%s"}]}`
	formattedstring := fmt.Sprintf(bodystring, input)
	jsondata := []byte(formattedstring)

	req, err := http.NewRequest("POST", "https://duckduckgo.com/duckchat/v1/chat", bytes.NewBuffer(jsondata))
	if err != nil {
		return "", fmt.Errorf("HANDLE ERROR: %s", err)
	}

	req.Header.Set("X-Vqd-4", vqd)
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Referer", "https://duckduckgo.com/")
	req.Header.Set("Origin", "https://duckduckgo.com")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.6778.86 Safari/537.36")
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("HANDLE ERROR: %s", err)		
	}
	defer resp.Body.Close()

	var responseBuilder strings.Builder
	scanner := bufio.NewScanner(resp.Body)

	lineNumber := 0
	var line string
	for scanner.Scan() {
		// first line seems to look like chat information
		lineNumber++
		if lineNumber == 1 {
			info := scanner.Text()

			var chatinfo ChatInformation
			json.Unmarshal([]byte(info), &chatinfo)
			fmt.Printf("Chat Information: %s\n", info)
		}

		line = scanner.Text()

		if !strings.HasPrefix(line, "data: ") {
			continue
		}

		jsonData := line[6:]

		// everything after this prefix is message fragments and valid json
		var frag MessageFragment
		err := json.Unmarshal([]byte(jsonData), &frag)
		if err != nil {
			continue
		}

		// add the message fragment to the final response from ai
		responseBuilder.WriteString(frag.Message)
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("Error reading stream: %v", err)
	}

	return responseBuilder.String(), nil
}

func main() {

}
