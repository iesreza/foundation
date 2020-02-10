package lib

import (
	"encoding/json"
)

func EncodeJSON(obj interface{}) []byte {
	b, _ := json.Marshal(obj)
	return b
}

func DecodeJSON(input []byte, v interface{}) error {
	err := json.Unmarshal(input, v)
	return err
}
