package pkgb

import "github.com/kuun/slog"

var log = slog.GetLogger()

func Hello() {
	log.Debug("hello pkgb")
}
