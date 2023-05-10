package games

type RequestFieldSize struct {
	Size int `json:"size,omitempty"`
}

type RequestComb struct {
	Row int `json:"row"`
	Col int `json:"col"`
}

type ResponseCreate struct {
	SessionId string `json:"session_id,omitempty"`
	Message   string `json:"message,omitempty"`
	Err       string `json:"err,omitempty"`
}

type ResponseJoin struct {
	GameId  string `json:"game_id,omitempty"`
	Message string `json:"message,omitempty"`
	Err     string `json:"err,omitempty"`
}

type ResponseGame struct {
	Winner  int    `json:"winner,omitempty"`
	Message string `json:"message,omitempty"`
	Field   string `json:"field,omitempty"`
	Err     string `json:"err,omitempty"`
}
