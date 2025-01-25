package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

func getvqd() string {
	req, err := http.NewRequest("GET", "https://duckduckgo.com/duckchat/v1/status", nil)

	if err != nil {
		fmt.Println(fmt.Sprintf("HANDLE ERROR: %s", err))
	}

	req.Header.Set("X-Vqd-Accept", "1")
	req.Header.Set("Referer", "https://duckduckgo.com/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.6778.86 Safari/537.36")

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println(fmt.Sprintf("HANDLE ERROR: %s", err))
	}

	defer resp.Body.Close()

	vqd := resp.Header.Get("x-vqd-4")

	return vqd
}

func prompt(input string) interface{} {
	var vqd string = getvqd()

	bodystring := `{"model": "gpt-4o-mini", "messages": [{"role": "user","content": "%s"}]}`

	formattedstring := fmt.Sprintf(bodystring, input)

	jsondata := []byte(formattedstring)

	req, err := http.NewRequest("POST", "https://duckduckgo.com/duckchat/v1/chat", bytes.NewBuffer(jsondata))

	if err != nil {
		fmt.Println(fmt.Sprintf("HANDLE ERROR: %s", err))
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
		fmt.Println(fmt.Sprintf("HANDLE ERROR: %s", err))
	}

	defer resp.Body.Close()

	byteslop, err := io.ReadAll(resp.Body)

	if err != nil {
		fmt.Println(fmt.Sprintf("HANDLE ERROR: %s", err))
	}

	aislop := string(byteslop)

	return aislop
}

func main() {
	fmt.Println(prompt("If you see this tell me LETS GOOOO"))
}
