package pkgb

import "github.com/kuun/slog"

type slogPkgInfo struct {
}

var log = slog.GetLogger(slogPkgInfo{})

func Hello() {
	log.Debug("hello pkgb")
}
