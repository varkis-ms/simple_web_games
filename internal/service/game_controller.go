package service

import "github.com/julienschmidt/httprouter"

type GameController interface {
	Register(router *httprouter.Router)
}
