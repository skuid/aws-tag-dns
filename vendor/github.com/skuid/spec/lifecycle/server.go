package lifecycle

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/skuid/spec"
	"go.uber.org/zap"
)

var term = syscall.SIGTERM

// ShutdownOnTerm accepts an *http.Server and will gracefully shut it down
// when a SIGTERM is received, after ShutdownTimer seconds (default 15)
func ShutdownOnTerm(srv *http.Server) {
	// subscribe to SIGTERM signal
	c := make(chan os.Signal)
	signal.Notify(c, term)

	go func() {
		<-c
		Ready = false
		Shutdown = true

		spec.Logger.Info("Received SIGTERM! Beginning shutdown", zap.Int64("timeout", ShutdownTimer))
		time.Sleep(time.Duration(ShutdownTimer) * time.Second)

		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			<-ctx.Done()
			spec.Logger.Fatal("Error Shutting down", zap.Error(ctx.Err()))
		}

		select {
		case <-time.After(6 * time.Second):
			spec.Logger.Fatal("Server did not shut down in time, exiting")
			os.Exit(1)
		case <-ctx.Done():
			err := srv.Close()
			if err != nil && err != http.ErrServerClosed {
				spec.Logger.Error("Error shutting down", zap.Error(err))
				os.Exit(1)
			} else {
				spec.Logger.Info("Shut down successfully")
				os.Exit(0)
			}
		}
	}()
}
