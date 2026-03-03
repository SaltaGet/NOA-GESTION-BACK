package utils

import "encoding/json"

func ModelToString(model any) *string {
	data, err := json.Marshal(model)
	if err != nil {
		return nil
	}
	stringData := string(data)
	return &stringData
}