package main

import (
	"net/http"

	buildpermission "github.com/krateoplatformops/azuredevops-provider-kog/buildpermission-plugin/handlers"
	"github.com/krateoplatformops/azuredevops-provider-kog/pkg/handlers"
	"github.com/krateoplatformops/azuredevops-provider-kog/pkg/health"
	"github.com/krateoplatformops/azuredevops-provider-kog/pkg/server"
	"github.com/rs/zerolog/log"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title           Azure DevOps BuildPermission Plugin API for Krateo Operator Generator (KOG)
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

	// BuildPermission
	srv.Mux().Handle("GET /plugin/buildpermission/{organization}/{projectId}", buildpermission.GetBuildPermission(opts))
	srv.Mux().Handle("POST /plugin/buildpermission/{organization}/{projectId}", buildpermission.PostBuildPermission(opts))

	// Swagger UI
	srv.Mux().Handle("/swagger/", httpSwagger.WrapHandler)

	// Kubernetes health check endpoints
	srv.Mux().HandleFunc("GET /healthz", health.LivenessHandler(srv.Healthy()))
	srv.Mux().HandleFunc("GET /readyz", health.ReadinessHandler(srv.Ready(), opts.Client.(*http.Client)))

	srv.Run()
}
