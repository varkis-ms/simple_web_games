package tic_tac_toe

type RequestFieldSize struct {
	Size int `json:"size,omitempty"`
}

type RequestComb struct {
	Row int `json:"row"`
	Col int `json:"col"`
}

type ResponseGame struct {
	Winner  int
	Message string
}
