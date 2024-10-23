package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/deniskrumko/visual-coordinates-tool/pkg/extractor"
	"github.com/deniskrumko/visual-coordinates-tool/pkg/recognize"
	"github.com/go-chi/chi/v5"
)

const (
	FieldEndpoint             = "endpoint"
	FieldRequestIsJson        = "requestIsJson"
	FieldRequestJsonTemplate  = "requestJsonTemplate"
	FieldRequestFormDataField = "requestFormdataField"
	FieldResponseXYField      = "responseXYField"
	FieldImageURL             = "imageUrl"
	FieldFormFile             = "imageFile"
	CheckboxEnabled           = "on"
)

type Recognizer interface {
	GetResponse() ([]byte, error)
}

type Coordinates struct {
	Coordinates [][]int
	Response    map[string]any
}

func getRecognizeRouter() (chi.Router, error) {
	r := chi.NewRouter()

	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		isJsonRequest := r.FormValue(FieldRequestIsJson) == CheckboxEnabled

		// If request is not JSON - read image bytes using request data
		var imageContent io.Reader
		if !isJsonRequest {
			result, err := getFileContent(r)
			if err != nil {
				errorResponse(w, nil, fmt.Errorf("can't get file content: %w", err))
				return
			}

			imageContent = result
		}

		// Measure execution time
		start := time.Now()

		// Make request to specified service
		coordinates, err := recognizeCoordinates(
			isJsonRequest,
			getFormValue(r, FieldEndpoint),
			getFormValue(r, FieldRequestJsonTemplate),
			getFormValue(r, FieldRequestFormDataField),
			getFormValue(r, FieldResponseXYField),
			getFormValue(r, FieldImageURL),
			imageContent,
		)

		// Measure execution time
		if err != nil {
			var responseMap map[string]any
			if coordinates != nil {
				responseMap = coordinates.Response
			}
			errorResponse(w, responseMap, err)
		} else {
			successResponse(w, struct {
				Coordinates   [][]int        `json:"coordinates"`
				Response      map[string]any `json:"response"`
				ExecutionTime int            `json:"executionTime"`
			}{
				coordinates.Coordinates,
				coordinates.Response,
				int(time.Since(start).Milliseconds()),
			})
		}
	})

	return r, nil
}

// Recognize coordinates using data from request
func recognizeCoordinates(
	isJSON bool,
	endpoint string,
	jsonTemplate string,
	formdateField string,
	xyField string,
	imageURL string,
	imageContent io.Reader,
) (*Coordinates, error) {
	var recognizer Recognizer
	if isJSON {
		recognizer = recognize.NewJSONRecognizer(endpoint, jsonTemplate, imageURL)
	} else {
		recognizer = recognize.NewFormDataRecognizer(endpoint, formdateField, imageContent)
	}

	response, err := recognizer.GetResponse()
	if err != nil {
		return nil, fmt.Errorf("error getting response from service: %w", err)
	}

	ext, err := extractor.NewExtractor(xyField)
	if err != nil {
		return nil, fmt.Errorf("error creating response extractor: %w", err)
	}

	var responseMap map[string]any
	if err := json.Unmarshal(response, &responseMap); err != nil {
		return nil, err
	}

	coordinates, err := ext.Extract(responseMap)
	if err != nil {
		return &Coordinates{
			Response: responseMap,
		}, fmt.Errorf("error extracting coordinates: %w", err)
	}

	return &Coordinates{
		Coordinates: coordinates,
		Response:    responseMap,
	}, nil
}

func getFormValue(r *http.Request, name string) string {
	value := r.FormValue(name)
	return strings.Trim(value, " ")
}

// Get image content from request
func getFileContent(r *http.Request) (io.Reader, error) {
	// If request contains URL – download image from it
	if fileUrl := r.FormValue(FieldImageURL); fileUrl != "" {
		response, err := http.Get(fileUrl)
		if err != nil {
			return nil, fmt.Errorf("can't get image from url: %w", err)
		}

		return response.Body, nil
	}

	// Otherwise – get raw image from form
	file, _, err := r.FormFile(FieldFormFile)
	if err != nil {
		return nil, fmt.Errorf("can't get image from form: %w", err)
	}

	return file, nil
}
