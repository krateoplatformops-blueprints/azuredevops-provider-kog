package health

import (
	"context"
	"encoding/json"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/rs/zerolog/log"
)

type healthResponse struct {
	Status struct {
		Health string `json:"health"`
	} `json:"status"`
}

var (
	// Reference: https://learn.microsoft.com/en-us/rest/api/azure/devops/status/health/get
	azureDevOpsHealthURL = "https://status.dev.azure.com/_apis/status/health"
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
// For a proxy service like this one, this includes checking connectivity to Azure DevOps API
func ReadinessHandler(ready *int32, client *http.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// First check if the service is marked as ready
		if atomic.LoadInt32(ready) == 0 {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("Service Not Ready"))
			return
		}

		// For a Azure DevOps API proxy, we verify the Azure DevOps status endpoint is reachable and healthy
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		req, err := http.NewRequestWithContext(ctx, "GET", azureDevOpsHealthURL, nil)
		if err != nil {
			log.Debug().Err(err).Msg("failed to create Azure DevOps API request for readiness check")
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("Azure DevOps API Unreachable"))
			return
		}

		resp, err := client.Do(req)
		if err != nil {
			log.Debug().Err(err).Msg("failed to reach Azure DevOps API for readiness check")
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("Azure DevOps API Unreachable"))
			return
		}
		defer resp.Body.Close()

		// Check for 200 OK and that "status":{"health":"healthy"} is in the response
		if resp.StatusCode == http.StatusOK {
			var hr healthResponse
			if err := json.NewDecoder(resp.Body).Decode(&hr); err == nil && hr.Status.Health == "healthy" {
				//log.Debug().Msg("Azure DevOps API readiness check: healthy")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("OK"))
				return
			}
			log.Debug().Msg("Azure DevOps API readiness check: health not healthy or failed to parse response")
		} else {
			log.Debug().Int("status", resp.StatusCode).Msg("Azure DevOps API returned unexpected status for readiness check")
		}
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("Azure DevOps API Error"))
	}
}
