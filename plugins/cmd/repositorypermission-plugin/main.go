package main

import (
	"net/http"

	"github.com/krateoplatformops/azuredevops-provider-kog/pkg/handlers"
	"github.com/krateoplatformops/azuredevops-provider-kog/pkg/health"
	"github.com/krateoplatformops/azuredevops-provider-kog/pkg/server"
	repositorypermission "github.com/krateoplatformops/azuredevops-provider-kog/repositorypermission-plugin/handlers"
	"github.com/rs/zerolog/log"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title           Azure DevOps RepositoryPermission Plugin API for Krateo Operator Generator (KOG)
// @version         1.0
// @description     Simple wrapper around Azure DevOps API to provide consistency of API response for Krateo Operator Generator (KOG)
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

	// RepositoryPermission
	srv.Mux().Handle("GET /plugin/repositorypermission/{organization}/{projectId}/{repositoryId}", repositorypermission.GetRepositoryPermission(opts))
	srv.Mux().Handle("POST /plugin/repositorypermission/{organization}/{projectId}/{repositoryId}", repositorypermission.PostRepositoryPermission(opts))

	// Swagger UI
	srv.Mux().Handle("/swagger/", httpSwagger.WrapHandler)

	// Kubernetes health check endpoints
	srv.Mux().HandleFunc("GET /healthz", health.LivenessHandler(srv.Healthy()))
	srv.Mux().HandleFunc("GET /readyz", health.ReadinessHandler(srv.Ready(), opts.Client.(*http.Client)))

	srv.Run()
}
