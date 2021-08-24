package app

import (
	"database/sql"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/zlobste/fake-wallet/internal/app/utils"
	"github.com/zlobste/fake-wallet/internal/config"
	"net/http"
)

type App interface {
	Run() error
}

type app struct {
	log    *logrus.Logger
	config config.Config
	db     *sql.DB
	auth   utils.Auth
}

func New(cfg config.Config) App {
	return &app{
		config: cfg,
		log:    cfg.Logging(),
		db:     cfg.DB(),
		auth:   utils.NewAuth(cfg.JWTSecret()),
	}
}

func (a *app) Run() error {
	defer func() {
		if rvr := recover(); rvr != nil {
			a.log.Error("app panicked\n", rvr)
		}
	}()

	a.log.WithField("port", a.config.Listener()).Info("Starting server")
	if err := http.ListenAndServe(a.config.Listener(), a.router()); err != nil {
		return errors.Wrap(err, "listener failed")
	}
	return nil
}
