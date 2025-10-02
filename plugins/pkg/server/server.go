
package server

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/krateoplatformops/plumbing/env"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Server struct {
	*http.Server
	mux     *http.ServeMux
	healthy int32
	ready   int32
}

func New() *Server {
	debugOn := flag.Bool("debug", env.Bool("DEBUG", true), "dump verbose output")
	port := flag.Int("port", env.Int("PORT", 8080), "port to listen on")
	noColor := flag.Bool("no-color", env.Bool("NO_COLOR", false), "disable color output")

	flag.Parse()

	mux := http.NewServeMux()

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debugOn {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:     os.Stderr,
		NoColor: *noColor,
	}).With().Timestamp().Logger()

	return &Server{
		Server: &http.Server{
			Addr:         fmt.Sprintf(":%d", *port),
			Handler:      mux,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 50 * time.Second,
			IdleTimeout:  30 * time.Second,
		},
		mux: mux,
	}
}

func (s *Server) Mux() *http.ServeMux {
	return s.mux
}

func (s *Server) Healthy() *int32 {
	return &s.healthy
}

func (s *Server) Ready() *int32 {
	return &s.ready
}

func (s *Server) Run() {
	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGKILL,
		syscall.SIGHUP,
		syscall.SIGQUIT,
	)
	defer stop()

	go func() {
		atomic.StoreInt32(&s.healthy, 1)
		atomic.StoreInt32(&s.ready, 1)

		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msgf("could not listen on %s", s.Addr)
		}
	}()

	log.Info().Msgf("server is ready to handle requests at %s", s.Addr)
	<-ctx.Done()

	stop()
	log.Info().Msg("server is shutting down gracefully, press Ctrl+C again to force")

	atomic.StoreInt32(&s.ready, 0)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	s.SetKeepAlivesEnabled(false)
	if err := s.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("server forced to shutdown")
	}

	atomic.StoreInt32(&s.healthy, 0)
	log.Info().Msg("server gracefully stopped")
}
