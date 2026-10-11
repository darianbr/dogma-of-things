package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/darian/dot"
	"github.com/darian/dot/ui/api"
	"github.com/darian/dot/ui/projections"
	"github.com/dogmatiq/runkit"
	"golang.org/x/sync/errgroup"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	db := projections.MustNewDB()
	defer db.Close()

	app := &dot.App{
		ReadDB: db,
	}

	e := &runkit.Engine{
		DB:  db,
		App: app,
	}

	// Set up API
	apiHandler := &api.Handler{
		CommandExecutor: e,
		ReadDB:          db,
	}

	mux := http.NewServeMux()
	apiHandler.SetupRoutes(mux)

	// Health check
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "ok")
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Println("DoT server starting on :8080")
	log.Println("API endpoints:")
	log.Println("  POST /api/devices     — Register a new device")
	log.Println("  GET  /api/devices     — List all devices")
	log.Println("  GET  /api/devices/:id — Get device by ID")

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			return fmt.Errorf("server error: %w", err)
		}
		return nil
	})

	g.Go(func() error {
		if err := e.Run(ctx); err != nil {
			log.Printf("engine error: %v", err)
		}
		return nil
	})

	g.Go(func() error {
		<-ctx.Done()
		log.Println("shutting down...")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		server.Shutdown(shutdownCtx)

		return nil
	})

	if err := g.Wait(); err != nil {
		log.Fatalf("error: %v", err)
	}
}
