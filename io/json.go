package pio

import "encoding/json"

func MarshalJSON(in interface{}) (string, error) {
	if in == nil {
		return "", nil
	}

	b, err := json.Marshal(in)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
