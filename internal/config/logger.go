package config

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"sync"
)

type Logger interface {
	Logging() *logrus.Logger
}

type logger struct {
	level string

	cache struct {
		log *logrus.Logger
	}

	sync.Once
}

func NewLogger(level string) Logger {
	return &logger{
		level: level,
	}
}

func (l *logger) Logging() *logrus.Logger {
	l.Once.Do(func() {
		level, err := logrus.ParseLevel(l.level)
		if err != nil {
			panic(errors.Wrapf(err, "failed to parse logging level %s", l.level))
		}

		l.cache.log = logrus.New()
		l.cache.log.SetLevel(level)
	})
	return l.cache.log
}