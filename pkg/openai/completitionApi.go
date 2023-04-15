package openai

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type ModelCompletition string

const (
	TextDavinchi3 ModelCompletition = "text-davinci-003"
	TextDavinchi2 ModelCompletition = "text-davinci-002"
)

type CompletionResponse struct {
	Choices []struct {
		Text string `json:"text"`
	} `json:"choices"`
	Model string `json:"model"`
}

func (c *Client) getCompletion(prompt string, maxTokens int, temperature float32, model ModelCompletition, stop []string) (string, error) {
	apiURL := "https://api.openai.com/v1/completions"
	data := Data{
		Model:       string(model),
		Temperature: temperature,
		Stop:        stop,
		Prompt:      prompt,
		MaxTokens:   maxTokens,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("API request failed")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var completionResponse CompletionResponse
	err = json.Unmarshal(body, &completionResponse)
	if err != nil {
		return "", err
	}

	return completionResponse.Choices[0].Text, nil
}
