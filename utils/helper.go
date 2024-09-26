package utils

import "encoding/json"

func JSONUnmarshal(data []byte, v interface{}) error {
	if err := json.Unmarshal([]byte(data), &v); err != nil {
		return err
	}
	return nil
}

func JSONMarshal(data interface{}) ([]byte, error) {
	result, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func ToByte(i interface{}) []byte {
	byte_, _ := JSONMarshal(i)
	return byte_
}

func Dump(i interface{}) string {
	return string(ToByte(i))
}
