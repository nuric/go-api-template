package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nuric/go-api-template/utils"
	"github.com/rs/zerolog/log"
)

func main() {
	mux := http.NewServeMux()
	// ---------------------------
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		utils.Encode(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	// ---------------------------
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	go func() {
		log.Info().Str("httpAddr", server.Addr).Msg("HTTPAPI.Serve")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Failed to listen on port 8080")
		}
	}()
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	sig := <-quit
	log.Info().Str("signal", sig.String()).Msg("Shutting down server...")
	// The default kubernetes grace period is 30 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	if err := server.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("Server forced to shutdown")
	}
	cancel()
	log.Info().Msg("Server gracefully stopped")

}
