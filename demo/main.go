package main

import (
	"github.com/kuun/slog"
	"github.com/kuun/slog/demo/pkga"
	"github.com/kuun/slog/demo/pkgb"
)

var log = slog.GetLogger()

func main() {
	log.Debug("hello slog")
	pkga.Hello()
	pkgb.Hello()
	slog.Close()
}
