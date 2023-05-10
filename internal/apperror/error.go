package apperror

import "encoding/json"

var (
	ErrNoSession = New(nil, "there is no such session exists",
		"", "SWG-0001")
	ErrNoGame = New(nil, "there is no such game",
		"", "SWG-0002")
	ErrBadSize = New(nil, "field size should not exceed 10x10",
		"", "SWG-0003")
	ErrWrongUser = New(nil, "it's not your turn now, wait for your opponent's move",
		"", "SWG-0004")
	ErrIncorrectMove = New(nil, "there is no such cell in the field of the given dimension",
		"", "SWG-0005")
	ErrFilledCell = New(nil, "this cell is already occupied",
		"", "SWG-0006")
	ErrBadSession = New(nil, "session exists but could not be decoded",
		"", "SWG-0007")
)

type AppError struct {
	Err              error  `json:"-"`
	Message          string `json:"message,omitempty"`
	DeveloperMessage string `json:"developer_message,omitempty"`
	Code             string `json:"code,omitempty"`
}

func (e *AppError) Error() string { return e.Message }

func (e *AppError) JsonMarshal() []byte {
	marshal, err := json.Marshal(e)
	if err != nil {
		return nil
	}
	return marshal
}

func New(err error, message, developerMessage, code string) *AppError {
	return &AppError{
		Err:              err,
		Message:          message,
		DeveloperMessage: developerMessage,
		Code:             code,
	}
}

func SystemError(err error) *AppError {
	return New(err, "internal system error", err.Error(), "SWG-0000")
}
