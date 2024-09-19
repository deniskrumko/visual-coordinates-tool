package recognize

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
)

type FormDataRecognizer struct {
	endpoint  string
	fieldName string
	content   io.Reader
}

func NewFormDataRecognizer(endpoint, template string, content io.Reader) *FormDataRecognizer {
	return &FormDataRecognizer{
		endpoint:  endpoint,
		fieldName: template,
		content:   content,
	}
}

func (r *FormDataRecognizer) GetResponse() ([]byte, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add formdata file
	part, err := writer.CreateFormFile(r.fieldName, r.fieldName)
	if err != nil {
		return nil, err
	}

	if _, err = io.Copy(part, r.content); err != nil {
		return nil, err
	}

	// Close multipart writer
	if err := writer.Close(); err != nil {
		return nil, err
	}

	// Make request
	request, err := http.NewRequest("POST", r.endpoint, body)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", writer.FormDataContentType())

	// Perform request
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Read response
	respBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}
