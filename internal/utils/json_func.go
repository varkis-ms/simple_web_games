package utils

import (
	"encoding/json"
	"io"
	"simple_web_games/pkg/logging"
)

func JsonDecode(logger *logging.Logger, body *io.ReadCloser, pattern interface{}) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Error(err)
		}
	}(*body)
	req, err := io.ReadAll(*body)
	if err != nil {
		logger.Error(err)
	}
	err = json.Unmarshal(req, pattern)
	if err != nil {
		logger.Errorf("can't unmarshal: %v", err)
	}
}

func JsonEncode(logger *logging.Logger, resp interface{}) []byte {
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		logger.WithError(err).Fatal("Error happened in JSON marshal")
	}
	return jsonResp
}
