package main

import (
	"net/http"

	"github.com/krateoplatformops/azuredevops-provider-kog/pkg/handlers"
	"github.com/krateoplatformops/azuredevops-provider-kog/pkg/health"
	"github.com/krateoplatformops/azuredevops-provider-kog/pkg/server"
	variablegroup "github.com/krateoplatformops/azuredevops-provider-kog/variablegroup-plugin/handlers"
	"github.com/rs/zerolog/log"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title           Azure DevOps VariableGroup Plugin API for Krateo Operator Generator (KOG)
// @version         1.0
// @description     Simple wrapper around Azure DevOps API to handle VariableGroup delete operation with projectIds extraction
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

	srv.Mux().Handle("POST /plugin/{organization}/{project}/variablegroups/{groupId}", variablegroup.GetVariableGroup(opts))
	srv.Mux().Handle("POST /plugin/{organization}/{project}/variablegroups", variablegroup.FindByVariableGroup(opts))
	srv.Mux().Handle("DELETE /plugin/{organization}/{project}/variablegroups/{groupId}", variablegroup.DeleteVariableGroup(opts))

	// Swagger UI
	srv.Mux().Handle("/swagger/", httpSwagger.WrapHandler)

	// Kubernetes health check endpoints
	srv.Mux().HandleFunc("GET /healthz", health.LivenessHandler(srv.Healthy()))
	srv.Mux().HandleFunc("GET /readyz", health.ReadinessHandler(srv.Ready(), opts.Client.(*http.Client)))

	srv.Run()
}
