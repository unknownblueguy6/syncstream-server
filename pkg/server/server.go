package server

import (
	"flag"
	"log/slog"
	"net/http"

	"syncstream-server/pkg/internal/request"
	"syncstream-server/pkg/internal/room"
)

var addr = flag.String("addr", "localhost:8080", "Address of Server")
var debug = flag.Bool("debug", false, "Enable debug logging.")

func Run() {
	// TODO make ci/cd pipeline work
	// TODO add validation for all request handlers
	flag.Parse()
	if *debug {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		slog.Debug("Debug Logging Enabled")
	}
	go room.Manager.Run()

	slog.Info("Starting Server at " + *addr)
	http.HandleFunc("POST /init", request.InitHandler)
	http.HandleFunc("POST /create", request.CreateHandler)
	http.HandleFunc("POST /joinToken", request.JoinTokenHandler)
	http.HandleFunc("/join", request.JoinHandler)
	err := http.ListenAndServe(*addr, nil)

	if err != nil {
		slog.Error("Unable to start HTTP server")
		return
	}
}
