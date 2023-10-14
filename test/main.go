package main

import "github.com/kuun/slog"

type __logger struct {
}

var log = slog.GetLogger(__logger{})

func main() {
	defer slog.Close()
	log.Debug("debug")
	log.Info("info")
	log.Notice("notice")
	log.Warn("warn")
	log.Error("error")
	log.Fatal("fatal")
}
