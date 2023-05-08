package service

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"simple_web_games/internal/apperror"
	"simple_web_games/internal/controllers"
	"simple_web_games/internal/games"
	gameTtt "simple_web_games/internal/games/tic_tac_toe"
	"simple_web_games/internal/utils"
	"simple_web_games/pkg/logging"
	"sync"
)

type handler struct {
	joinList    sync.Map
	sessionList *sessionList
	logger      logging.Logger
}

type sessionList struct {
	mx sync.RWMutex
	m  map[string]games.GameField
}

func newSessionList() *sessionList {
	return &sessionList{
		mx: sync.RWMutex{},
		m:  make(map[string]games.GameField),
	}
}

func (c *sessionList) Load(key string) (games.GameField, bool) {
	c.mx.RLock()
	defer c.mx.RUnlock()
	val, ok := c.m[key]
	return val, ok
}

func (c *sessionList) Store(key string, value games.GameField) {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.m[key] = value
}

func New(logger logging.Logger) controllers.GameController {
	return &handler{
		joinList:    sync.Map{},
		sessionList: newSessionList(),
		logger:      logger,
	}
}

func (h *handler) Register(router *httprouter.Router) {
	router.GET("/list", h.GetList)
	router.POST("/create", sessionMiddleware(h.CreateGame))
	router.GET("/join/:session_id", sessionMiddleware(h.JoinGame))
	router.POST("/game/:game_id", sessionMiddleware(h.GameProgress))
}

func (h *handler) GetList(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	availableSessions := make(map[int]string)
	i := 1
	h.joinList.Range(func(k, _ any) bool {
		availableSessions[i] = k.(string)
		i++
		return true
	})
	jsonResp, err := json.Marshal(availableSessions)
	if err != nil {
		logger.WithError(err).Fatal("Error happened in JSON marshal")
	}
	w.Write(jsonResp)
}

func (h *handler) CreateGame(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var size gameTtt.RequestFieldSize
	resp := make(map[string]string)
	utils.JsonDecode(&h.logger, &r.Body, &size)
	session := r.Context().Value("session").(*sessions.Session)
	sessionId := uuid.New().String()
	gameId := uuid.New().String()
	session.Values["game_id"] = gameId
	session.Values["uid"] = 1
	if size.Size == 0 {
		size.Size = 3
	}
	field := gameTtt.NewGameField(&size.Size)
	h.sessionList.Store(gameId, field)
	err := session.Save(r, w)
	if err != nil {
		logger.Error(err)
		resp["message"] = "Failed to save session"
		w.WriteHeader(http.StatusTeapot)
	} else {
		resp["message"] = "Status Created"
		w.WriteHeader(http.StatusCreated)
	}
	jsonResp := utils.JsonEncode(logger, &resp)
	h.joinList.Store(sessionId, gameId)
	logger.Tracef("size: %d, uid: %s, gid: %s. Game created successfully.", size.Size, sessionId, gameId)
	w.Write(jsonResp)
}

func (h *handler) JoinGame(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	resp := make(map[string]string)
	session := r.Context().Value("session").(*sessions.Session)
	sessionId := ps.ByName("session_id")
	gameId, ok := h.joinList.LoadAndDelete(sessionId)
	if !ok {
		resp["message"] = "no such session exists"
		jsonResp := utils.JsonEncode(logger, &resp)
		w.Write(jsonResp)
		return
	}
	session.Values["game_id"] = gameId.(string)
	session.Values["uid"] = 2
	err := session.Save(r, w)
	if err != nil {
		log.Print(err)
	}
	resp["message"] = "connection to the session was successful"
	resp["game_id"] = gameId.(string)
	jsonResp := utils.JsonEncode(logger, &resp)
	w.Write(jsonResp)
}

func (h *handler) GameProgress(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var appError *apperror.AppError
	var jsonResp []byte
	resp := gameTtt.ResponseGame{}
	session := r.Context().Value("session").(*sessions.Session)
	gameId := ps.ByName("game_id")
	if session.Values["game_id"] == nil || gameId != session.Values["game_id"] {
		w.WriteHeader(http.StatusForbidden)
		resp.Message = "no such game exists"
		jsonResp = utils.JsonEncode(logger, &resp)
		w.Write(jsonResp)
		return
	}
	field, ok := h.sessionList.Load(gameId)
	if !ok {
		w.WriteHeader(http.StatusForbidden)
		resp.Message = "no such session exists"
		jsonResp = utils.JsonEncode(logger, &resp)
		w.Write(jsonResp)
		return
	}
	uid := session.Values["uid"].(int)
	var move gameTtt.RequestComb
	utils.JsonDecode(&h.logger, &r.Body, &move)
	winner, check, err := field.Progress(move.Row, move.Col, uid)
	if err != nil {
		if errors.As(err, &appError) {
			w.WriteHeader(http.StatusForbidden)
			w.Write(appError.Marshal())
			return
		}
		w.WriteHeader(http.StatusTeapot)
		logger.WithError(err).Fatal("unexpected error in the business-logic")
		w.Write(apperror.SystemError(err).Marshal())
		return
	}
	if !check {
		resp.Message = "move made"
	} else {
		resp.Message = "end of the game"
		resp.Winner = winner
		session.Options.MaxAge = -1
		session.Save(r, w)
	}
	w.WriteHeader(http.StatusOK)
	jsonResp = utils.JsonEncode(logger, &resp)
	logger.Trace(resp)
	w.Write(jsonResp)
}
