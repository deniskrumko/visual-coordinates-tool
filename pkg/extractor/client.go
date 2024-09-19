package extractor

import (
	"encoding/json"
	"fmt"

	"github.com/go-bongo/go-dotaccess"
	"github.com/mitchellh/mapstructure"
)

type Extractor struct {
	fieldPath string
}

func NewExtractor(fieldPath string) (*Extractor, error) {
	if fieldPath == "" {
		return nil, fmt.Errorf("field path is empty")
	}

	return &Extractor{
		fieldPath: fieldPath,
	}, nil
}

func (e *Extractor) Extract(response []byte) ([][]int, error) {
	// convert response to map
	var responseMap map[string]any

	err := json.Unmarshal(response, &responseMap)
	if err != nil {
		return nil, err
	}

	// extract coordinates
	coordinates, err := extractCoordinates(responseMap, e.fieldPath)
	if err != nil {
		return nil, err
	}

	return coordinates, nil
}

func extractCoordinates(data map[string]any, fieldPath string) ([][]int, error) {
	val, err := dotaccess.Get(data, fieldPath)
	if err != nil {
		return nil, fmt.Errorf("can't extract field path: %w", err)
	}

	type Coords struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
	}

	coordinatePairs, ok := val.([]any)
	if !ok {
		return nil, fmt.Errorf("response is not an array: %v", val)
	}

	var result [][]int
	for _, coordData := range coordinatePairs {
		var coord Coords
		if err := mapstructure.Decode(coordData, &coord); err != nil {
			return nil, fmt.Errorf("can't decode coordinate: %w", err)
		}

		result = append(result, []int{int(coord.X), int(coord.Y)})
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("response has no coordinates")
	}

	return result, nil
}
