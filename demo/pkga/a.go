package pkga

import "github.com/kuun/slog"

var log = slog.GetLogger()

func Hello() {
	log.Debug("hello pkga")
}
