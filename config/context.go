package config

import (
	log "github.com/Sirupsen/logrus"
	"github.com/oklog/ulid"
)

type RequestContext struct {
	RunID ulid.ULID
	Log   *log.Entry
	*SystemConfig
	*Cell
}
