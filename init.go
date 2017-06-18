package service

import (
	"github.com/rai-project/config"
	"github.com/rai-project/logger"
	"github.com/rai-project/tracer"
	"github.com/rai-project/tracer/zipkin"
	"github.com/sirupsen/logrus"
)

var (
	log *logrus.Entry
)

func init() {
	config.AfterInit(func() {
		log = logger.New().WithField("pkg", "service")
	})

}
