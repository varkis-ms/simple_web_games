package utils

import (
	"encoding/json"
	"io"
)

func JsonDecode(body *io.ReadCloser, pattern interface{}) error {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(*body)
	req, err := io.ReadAll(*body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(req, pattern)
	if err != nil {
		return err
	}
	return nil
}

func JsonEncode(resp interface{}) ([]byte, error) {
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}
	return jsonResp, nil
}
