package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type Config struct {
	message  string
	api_key  string
	endpoint string
}

func NewConfig(message string, api_key string, endpoint string) Config {
	return Config{
		message,
		api_key,
		endpoint,
	}
}

func GetResponse(config *Config) (string, error) {
	if config.message == "" {
		return "", errors.New("Invalid Message")
	}

	url := config.endpoint + config.api_key

	requestBody := map[string]any{
		"contents": []map[string]any{
			{
				"parts": []map[string]string{
					{"text": config.message},
				},
			},
		},
		"generationConfig": map[string]any{
			"temperature": 0.7,
			"topK":        40,
			"topP":        0.95,
		},
	}

	// make json string from request body
	jsonData, err := json.Marshal(requestBody)

	if err != nil {
		return "", err
	}

	// get buffer from json
	buffer := bytes.NewBuffer(jsonData)
	req, err := http.NewRequest("post", url, buffer)
	if err != nil {
		return "", err
	}

	// set the data sent to be json data (change the headers)
	req.Header.Add("Content-Type", "application/json")

	// initialize a new client and send the request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	var data GeminiResponse
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&data); err != nil {
		return "", err
	}

	result := data.Candidates[0].Content.Parts[0].Text
	return result, nil
}

type GeminiResponse struct {
	Candidates []Candidate `json:"candidates"`
}

type Candidate struct {
	Content       Content        `json:"content"`
	FinishReason  string         `json:"finishReason"`
	SafetyRatings []SafetyRating `json:"safetyRatings"`
}

type Content struct {
	Parts []Part `json:"parts"`
}

type Part struct {
	Text string `json:"text"`
}

type SafetyRating struct {
	Category    string `json:"category"`
	Probability string `json:"probability"`
}
