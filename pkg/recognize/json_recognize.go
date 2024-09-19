package recognize

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

type JSONRecognizer struct {
	endpoint string
	template string
	imageURL string
}

func NewJSONRecognizer(endpoint, template, imageURL string) *JSONRecognizer {
	return &JSONRecognizer{
		endpoint: endpoint,
		template: template,
		imageURL: imageURL,
	}
}

func (r *JSONRecognizer) GetResponse() ([]byte, error) {
	// Prepare request from template
	jsonString := fmt.Sprintf(r.template, r.imageURL)
	body := []byte(jsonString)

	// Prepare HTTP requests
	request, err := http.NewRequest("POST", r.endpoint, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")

	// Send request
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err = io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("service responded with error: %s", body)
	}

	return body, nil
}
