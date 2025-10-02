package health

import (
	"context"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/rs/zerolog/log"
)

// LivenessHandler implements Kubernetes liveness probe
// Returns 200 if the application is running and hasn't deadlocked
func LivenessHandler(healthy *int32) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if atomic.LoadInt32(healthy) == 1 {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
			return
		}
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("Service Unavailable"))
	}
}

// ReadinessHandler implements Kubernetes readiness probe
// Returns 200 if the application is ready to serve traffic
// For a proxy service like this one, this includes checking connectivity to GitHub API
func ReadinessHandler(ready *int32, client *http.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// First check if the service is marked as ready
		if atomic.LoadInt32(ready) == 0 {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("Service Not Ready"))
			return
		}

		// For a GitHub API proxy, we verify we can reach GitHub
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		req, err := http.NewRequestWithContext(ctx, "GET", "https://api.github.com", nil)
		if err != nil {
			log.Debug().Err(err).Msg("failed to create GitHub API request for readiness check")
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("GitHub API Unreachable"))
			return
		}

		resp, err := client.Do(req)
		if err != nil {
			log.Debug().Err(err).Msg("failed to reach GitHub API for readiness check")
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("GitHub API Unreachable"))
			return
		}
		defer resp.Body.Close()

		// GitHub API should respond with some 2xx or 4xx status (4xx is still fine, means GitHub API is up)
		if resp.StatusCode >= 200 && resp.StatusCode < 500 {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Ready"))
			return
		}

		log.Debug().Int("status", resp.StatusCode).Msg("GitHub API returned unexpected status for readiness check")
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("GitHub API Error"))
	}
}
