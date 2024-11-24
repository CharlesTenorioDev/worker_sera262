package groq

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/sera_backend/internal/config"
	"github.com/sera_backend/internal/config/logger"
	"github.com/sera_backend/internal/dto"
)

type GroqlientInterface interface {
	DoRequest(method, endpoint string, payload dto.GeneratePayloadGroq) (*http.Response, error)
}

type ClientGroq struct {
	apiKey  string
	baseUrl string
	cliente *http.Client
}

type GroqResponse struct {
	Questions []struct {
		Content string   `json:"content"`
		Options []string `json:"options,omitempty"`
		Answer  string   `json:"answer"`
		Type    string   `json:"type"`
	} `json:"questions"`
}

var _ GroqlientInterface = (*ClientGroq)(nil)

func NewClient(cfg *config.Config) *ClientGroq {
	return &ClientGroq{
		apiKey:  cfg.GroqConfig.API_KEY,
		baseUrl: cfg.GroqConfig.URL,
		cliente: &http.Client{
			Timeout: time.Duration(cfg.AsaasConfig.ASAAS_TIMEOUT) * time.Second,
			Transport: &http.Transport{
				ForceAttemptHTTP2:   false,
				MaxConnsPerHost:     1,
				MaxIdleConns:        1,
				MaxIdleConnsPerHost: 1,
				TLSHandshakeTimeout: time.Duration(10) * time.Second,
			},
		},
	}
}

func (c *ClientGroq) DoRequest(method, endpoint string, payload dto.GeneratePayloadGroq) (*http.Response, error) {
	url := c.baseUrl + endpoint

	// Convert payload to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url, bytes.NewReader(payloadBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("access_token", c.apiKey)

	logger.Info("Sending request to URL: " + url)

	resp, err := c.cliente.Do(req)
	if err != nil {
		return nil, err
	}

	logger.Info("Response Status Code: " + strconv.Itoa(resp.StatusCode))
	return resp, nil
}
