package openai

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type ModelChat string

const (
	GPT35Turbo ModelChat = "gpt-3.5-turbo"
)

type ChatCompletionResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
		Index        int    `json:"index"`
	} `json:"choices"`
}

func (c *Client) getChatCompletion(messages []Message, temperature float32, model ModelChat, stop []string) (string, error) {
	apiURL := "https://api.openai.com/v1/chat/completions"
	data := Data{
		Model:       string(model),
		Temperature: temperature,
		Stop:        stop,
		Messages:    messages,
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
		body, _ := ioutil.ReadAll(resp.Body)
		return "", errors.New("API request failed" + string(body))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var chatCompletionResponse ChatCompletionResponse
	err = json.Unmarshal(body, &chatCompletionResponse)
	if err != nil {
		return "", err
	}

	return chatCompletionResponse.Choices[0].Message.Content, nil
}
