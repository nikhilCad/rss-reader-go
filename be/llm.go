package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// To add this later
func generateSummaryWithLlama(articleText string) (string, error) {
	payload := map[string]interface{}{
		"model":      "your-model-name", // e.g., "llama-2-7b-chat"
		"prompt":     fmt.Sprintf("Summarize this news article in 2-3 sentences:\n\n%s", articleText),
		"max_tokens": 200,
	}
	body, _ := json.Marshal(payload)
	resp, err := http.Post("http://localhost:8083/v1/completions", "application/json", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var result struct {
		Choices []struct {
			Text string `json:"text"`
		} `json:"choices"`
	}
	data, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(data, &result)
	if len(result.Choices) > 0 {
		return result.Choices[0].Text, nil
	}
	return "", fmt.Errorf("no summary returned")
}
