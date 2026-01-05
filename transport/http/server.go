package http

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"

	"gitlab.com/hmajid2301/goroutinely/internal/config"
)

type Server struct {
	logger *slog.Logger
	server *http.Server
	static embed.FS
}

func NewServer(logger *slog.Logger, static embed.FS, conf config.ServerConfig) *Server {
	s := &Server{
		logger: logger,
		static: static,
	}

	mux := http.NewServeMux()
	s.setupRoutes(mux)

	s.server = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		Handler: mux,
	}

	return s
}

func (s *Server) setupRoutes(mux *http.ServeMux) {
	// Health check
	mux.HandleFunc("GET /health", s.handleHealth)

	// Static files
	fsys, err := fs.Sub(s.static, "web/static")
	if err != nil {
		s.logger.Error("failed to create static file system", "error", err)
	} else {
		mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(fsys))))
	}

	// Home page
	mux.HandleFunc("GET /", s.handleHome)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (s *Server) handleHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Goroutinely - Habit Tracker</title>
    <link rel="stylesheet" href="/static/css/output.css">
</head>
<body class="bg-base text-text">
    <div class="container mx-auto px-4 py-8">
        <h1 class="text-3xl font-bold mb-4">Goroutinely</h1>
        <p class="text-lg">A cozy habit tracker for the self-hosted community 🌱</p>
    </div>
</body>
</html>
	`))
}

func (s *Server) Serve(ctx context.Context) error {
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
