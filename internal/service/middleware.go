package service

import (
	"context"
	"github.com/gorilla/sessions"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"simple_web_games/internal/apperror"
	"simple_web_games/internal/utils"
	"simple_web_games/pkg/logging"
)

var store *sessions.CookieStore

func SetupCookieStorage(cfg *utils.StorageConfig) {
	store = sessions.NewCookieStore([]byte(cfg.SecretKey))
	store.Options.MaxAge = 0
}

func sessionMiddleware(logger *logging.Logger, h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")
		session, err := store.Get(r, "game_token")
		if err != nil {
			session.Options.MaxAge = -1
			w.WriteHeader(http.StatusTeapot)
			logger.WithError(err).Error(apperror.ErrBadSession)
			w.Write(apperror.ErrBadSession.JsonMarshal())
			return
		}
		r = r.WithContext(context.WithValue(r.Context(), "session", session))
		h(w, r, ps)
	}
}
