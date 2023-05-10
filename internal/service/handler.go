package service

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/julienschmidt/httprouter"
	"net/http"
	_ "simple_web_games/docs"
	"simple_web_games/internal/apperror"
	"simple_web_games/internal/games"
	gameTtt "simple_web_games/internal/games/tic_tac_toe"
	"simple_web_games/internal/utils"
	"simple_web_games/pkg/logging"
	"sync"
)

var gameName = "Tik-Tak-Toe"

type handler struct {
	joinList    sync.Map
	sessionList *utils.SyncMap
	logger      logging.Logger
}

func New(logger logging.Logger) GameController {
	return &handler{
		joinList:    sync.Map{},
		sessionList: utils.NewSyncMap(),
		logger:      logger,
	}
}

func (h *handler) Register(router *httprouter.Router) {
	router.GET("/list", h.GetList)
	router.POST("/create", sessionMiddleware(&h.logger, h.CreateGame))
	router.GET("/join/:session_id", sessionMiddleware(&h.logger, h.JoinGame))
	router.POST("/game/:game_id", sessionMiddleware(&h.logger, h.GameProgress))
	router.ServeFiles("/swagger/*filepath", http.Dir("docs"))
}

// GetList
// @Summary Get List
// @Tags Games
// @Description Getting a list of available game sessions
// @ID Get-list
// @Accept json
// @Produce json
// @Success 200 {object} object{session_id}
// @Router /list [get]
func (h *handler) GetList(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	availableSessions := make(map[string][2]string)
	h.joinList.Range(func(k, _ any) bool {
		availableSessions[k.(string)] = [2]string{gameName, "2"}
		return true
	})
	jsonResp, err := json.Marshal(availableSessions)
	if err != nil {
		h.logger.Error(err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
}

// CreateGame
// @Summary Create Game Session
// @Tags Games
// @Description Create a new game session
// @ID Create-game
// @Accept json
// @Produce json
// @Param input body games.RequestFieldSize true "Field size"
// @Success 201 {object} games.ResponseCreate "Status Created"
// @Failure 418 {object} apperror.AppError "Bad params"
// @Failure 418 {object} apperror.AppError "Failed to save session"
// @Router /create [post]
func (h *handler) CreateGame(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var size games.RequestFieldSize
	resp := games.ResponseCreate{}
	err := utils.JsonDecode(&r.Body, &size)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.logger.Error(err)
		w.Write(apperror.SystemError(err).JsonMarshal())
		return
	}
	session := r.Context().Value("session").(*sessions.Session)
	sessionId := uuid.New().String()
	gameId := uuid.New().String()
	session.Values["game_id"] = gameId
	session.Values["uid"] = 1
	if size.Size == 0 {
		size.Size = 3
	} else if size.Size > 10 {
		w.WriteHeader(http.StatusBadRequest)
		h.logger.Error(apperror.ErrBadSize)
		w.Write(apperror.ErrBadSize.JsonMarshal())
		return
	}
	field := gameTtt.NewGameField(size.Size)
	h.sessionList.Store(gameId, field)
	err = session.Save(r, w)
	if err != nil {
		w.WriteHeader(http.StatusTeapot)
		h.logger.Error(err)
		w.Write(apperror.SystemError(err).JsonMarshal())
		return
	}
	resp.Message = "Session Created"
	jsonResp, err := utils.JsonEncode(&resp)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.logger.Error(err)
		w.Write(apperror.SystemError(err).JsonMarshal())
	}
	w.WriteHeader(http.StatusCreated)
	resp.SessionId = sessionId
	h.joinList.Store(sessionId, gameId)
	w.Write(jsonResp)
}

// JoinGame
// @Summary Join to Game Session
// @Tags Games
// @Description Connect to the created game session
// @ID Join-game
// @Accept json
// @Produce json
// @Param Cookie header string  false "game_token"
// @Param session_id path string true "Session ID"
// @Success 200 {object} games.ResponseJoin "connection to the session was successful"
// @Failure 418 {object} games.ResponseJoin "no such session exists"
// @Failure 418 {object} apperror.AppError "unexpected error"
// @Router /join/{session_id} [get]
func (h *handler) JoinGame(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	resp := games.ResponseJoin{}
	session := r.Context().Value("session").(*sessions.Session)
	sessionId := ps.ByName("session_id")
	gameId, ok := h.joinList.LoadAndDelete(sessionId)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(apperror.ErrNoSession.JsonMarshal())
		return
	}
	session.Values["game_id"] = gameId.(string)
	session.Values["uid"] = 2
	err := session.Save(r, w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp.Err = err.Error()
		jsonResp, err := utils.JsonEncode(&resp)
		if err != nil {
			h.logger.Error(err)
			w.Write(apperror.SystemError(err).JsonMarshal())
			return
		}
		h.logger.Error(err)
		w.Write(jsonResp)
		return
	}
	resp.Message = "connection to the session was successful"
	resp.GameId = gameId.(string)
	jsonResp, err := utils.JsonEncode(&resp)
	if err != nil {
		h.logger.Error(err)
		w.Write(apperror.SystemError(err).JsonMarshal())
	}
	w.Write(jsonResp)
}

// GameProgress
// @Summary Game Progress
// @Tags Games
// @Description Making moves and determining the winner
// @ID Game-progress
// @Accept json
// @Produce json
// @Param Cookie header string  false "game_token"
// @Param game_id path string true "Game ID"
// @Param input body games.RequestComb true "User's move"
// @Success 200 {object} games.ResponseGame "move made"
// @Failure 403 {object} apperror.AppError "no such game exists"
// @Failure 403 {object} apperror.AppError "no such session exists"
// @Failure 418 {object} games.ResponseGame "unexpected error"
// @Router /game/{game_id} [post]
func (h *handler) GameProgress(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var appError *apperror.AppError
	resp := games.ResponseGame{}
	session := r.Context().Value("session").(*sessions.Session)
	gameId := ps.ByName("game_id")
	if session.Values["game_id"] == nil || gameId != session.Values["game_id"] {
		w.WriteHeader(http.StatusForbidden)
		w.Write(apperror.ErrNoGame.JsonMarshal())
		return
	}
	field, ok := h.sessionList.Load(gameId)
	if !ok {
		w.WriteHeader(http.StatusForbidden)
		w.Write(apperror.ErrNoSession.JsonMarshal())
		return
	}
	uid := session.Values["uid"].(int)
	var move games.RequestComb
	err := utils.JsonDecode(&r.Body, &move)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		h.logger.Error(err)
		w.Write(apperror.SystemError(err).JsonMarshal())
		return
	}
	winner, check, err := field.Progress(move.Row, move.Col, uid)
	if err != nil {
		if errors.As(err, &appError) {
			w.WriteHeader(http.StatusForbidden)
			w.Write(appError.JsonMarshal())
			return
		}
		w.WriteHeader(http.StatusTeapot)
		h.logger.WithError(err)
		w.Write(apperror.SystemError(err).JsonMarshal())
		return
	}
	if !check {
		resp.Message = "move made"
	} else {
		resp.Message = "end of the game"
		resp.Winner = winner
		session.Options.MaxAge = -1
		session.Save(r, w)
		h.sessionList.Delete(gameId)
	}
	resp.Field = field.GetPrettyField()
	w.WriteHeader(http.StatusOK)
	jsonResp, err := utils.JsonEncode(&resp)
	if err != nil {
		h.logger.Error(err)
		w.Write(apperror.SystemError(err).JsonMarshal())
	}
	w.Write(jsonResp)
}
