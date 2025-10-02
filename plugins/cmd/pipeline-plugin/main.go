package main

import (
	"net/http"

	pipeline "github.com/krateoplatformops/azuredevops-provider-kog/pipeline-plugin/handlers"
	"github.com/krateoplatformops/azuredevops-provider-kog/pkg/handlers"
	"github.com/krateoplatformops/azuredevops-provider-kog/pkg/health"
	"github.com/krateoplatformops/azuredevops-provider-kog/pkg/server"
	"github.com/rs/zerolog/log"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title           Azure DevOps Pipeline Plugin API for Krateo Operator Generator (KOG)
// @version         1.0
// @description     Simple wrapper around Azure DevOps API to provide consisentency of API response for Krateo Operator Generator (KOG)
// @termsOfService  http://swagger.io/terms/
// @contact.name    Krateo Support
// @contact.url     https://krateo.io
// @contact.email   contact@krateoplatformops.io
// @license.name    Apache 2.0
// @license.url     http://www.apache.org/licenses/LICENSE-2.0.html
// @host            localhost:8080
// @BasePath        /
// @schemes         http
func main() {
	srv := server.New()

	opts := handlers.HandlerOptions{
		Log:    &log.Logger,
		Client: http.DefaultClient,
	}

	// Pipeline
	srv.Mux().Handle("GET /api/{organization}/{project}/pipelines/{id}", pipeline.GetPipeline(opts))
	srv.Mux().Handle("PUT /api/{organization}/{project}/pipelines/{id}", pipeline.PutPipeline(opts))
	srv.Mux().Handle("DELETE /api/{organization}/{project}/pipelines/{id}", pipeline.DeletePipeline(opts))

	// Swagger UI
	srv.Mux().Handle("/swagger/", httpSwagger.WrapHandler)

	// Kubernetes health check endpoints
	srv.Mux().HandleFunc("GET /healthz", health.LivenessHandler(srv.Healthy()))
	srv.Mux().HandleFunc("GET /readyz", health.ReadinessHandler(srv.Ready(), opts.Client.(*http.Client)))

	srv.Run()
}
