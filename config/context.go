package config

import log "github.com/Sirupsen/logrus"

type RequestContext struct {
	RunID string
	TagID string
	Log   *log.Entry
	*SystemConfig
	*Cell
}
