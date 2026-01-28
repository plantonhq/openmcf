package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/plantonhq/openmcf/app/backend/internal/database"
	"github.com/plantonhq/openmcf/app/backend/internal/service"
	"github.com/sirupsen/logrus"

	cloudresourcev1connect "github.com/plantonhq/openmcf/apis/org/openmcf/app/cloudresource/v1/cloudresourcev1connect"
	credentialv1connect "github.com/plantonhq/openmcf/apis/org/openmcf/app/credential/v1/credentialv1connect"
	stackupdatev1connect "github.com/plantonhq/openmcf/apis/org/openmcf/app/stackupdate/v1/stackupdatev1connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

// Server wraps the HTTP server and dependencies.
type Server struct {
	httpServer *http.Server
	mongo      *database.MongoDB
}

// Config holds server configuration.
type Config struct {
	Port    string
	MongoDB *database.MongoDB
}

// corsMiddleware wraps an HTTP handler with CORS headers.
// Supports Connect RPC and gRPC-Web protocols.
//
// CORS is enabled by default for direct Docker deployments.
// If you're using a reverse proxy (Caddy/nginx) that handles CORS,
// set ENABLE_CORS=false to avoid duplicate headers.
func corsMiddleware(next http.Handler) http.Handler {
	// Check if CORS should be disabled (when using reverse proxy)
	enableCORS := os.Getenv("ENABLE_CORS")
	if enableCORS == "false" {
		logrus.Info("CORS middleware disabled (reverse proxy handles CORS)")
		return next
	}

	logrus.Info("CORS middleware enabled (handling browser requests)")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers for allowed origins
		origin := r.Header.Get("Origin")
		if origin != "" {
			// Check CORS_ALLOWED_ORIGINS env var for custom origins, or use defaults
			allowedOriginsEnv := os.Getenv("CORS_ALLOWED_ORIGINS")
			var allowedOrigins []string

			if allowedOriginsEnv != "" {
				// Parse comma-separated origins from env var
				for _, origin := range splitAndTrim(allowedOriginsEnv, ",") {
					allowedOrigins = append(allowedOrigins, origin)
				}
			} else {
				// Default: Allow same-origin requests (common Docker setup)
				// Users can customize via CORS_ALLOWED_ORIGINS env var
				allowedOrigins = []string{
					"http://localhost:3000",
					"http://127.0.0.1:3000",
					"http://localhost:3001",
					"http://127.0.0.1:3001",
				}
			}

			allowed := false
			for _, allowedOrigin := range allowedOrigins {
				if origin == allowedOrigin || allowedOrigin == "*" {
					allowed = true
					break
				}
			}

			if allowed {
				// Set proper origin (or * if configured)
				if origin != "" {
					w.Header().Set("Access-Control-Allow-Origin", origin)
				}
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
				// Connect RPC and gRPC-Web specific headers
				w.Header().Set("Access-Control-Allow-Headers",
					"Content-Type, Authorization, X-Requested-With, "+
						"Connect-Protocol-Version, Connect-Timeout-Ms, "+
						"grpc-timeout, keep-alive, "+
						"x-accept-content-transfer-encoding, x-accept-response-streaming, "+
						"x-grpc-web, x-user-agent")
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				w.Header().Set("Access-Control-Max-Age", "3600")
				// Expose both gRPC-Web and Connect RPC headers (required for Connect-RPC)
				w.Header().Set("Access-Control-Expose-Headers",
					"Content-Encoding, Connect-Content-Encoding, "+
						"Grpc-Encoding, Grpc-Accept-Encoding, "+
						"Grpc-Status, Grpc-Message, Grpc-Status-Details-Bin, "+
						"Connect-Protocol-Version, Connect-Accept-Encoding")
			}
		}

		// Handle preflight OPTIONS request
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// NewServer creates a new server instance.
func NewServer(cfg *Config) *Server {
	// Create repositories
	cloudResourceRepo := database.NewCloudResourceRepository(cfg.MongoDB)
	stackUpdateRepo := database.NewStackUpdateRepository(cfg.MongoDB)
	stackUpdateStreamingResponseRepo := database.NewStackUpdateStreamingResponseRepository(cfg.MongoDB)
	credentialRepo := database.NewCredentialRepository(cfg.MongoDB)

	// Create credential resolver
	credentialResolver := service.NewCredentialResolver(credentialRepo)

	// Create services
	stackUpdateService := service.NewStackUpdateService(stackUpdateRepo, cloudResourceRepo, stackUpdateStreamingResponseRepo, credentialResolver)
	cloudResourceService := service.NewCloudResourceService(cloudResourceRepo, stackUpdateService)
	credentialService := service.NewCredentialService(credentialRepo)

	mux := http.NewServeMux()

	// Register the CloudResourceCommandController (Create, Update, Delete, Apply)
	cloudResourceCommandPath, cloudResourceCommandHandler := cloudresourcev1connect.NewCloudResourceCommandControllerHandler(cloudResourceService)
	mux.Handle(cloudResourceCommandPath, corsMiddleware(cloudResourceCommandHandler))

	// Register the CloudResourceQueryController (List, Get, Count)
	cloudResourceQueryPath, cloudResourceQueryHandler := cloudresourcev1connect.NewCloudResourceQueryControllerHandler(cloudResourceService)
	mux.Handle(cloudResourceQueryPath, corsMiddleware(cloudResourceQueryHandler))

	// Register the StackUpdateCommandController (DeployCloudResource)
	stackUpdateCommandPath, stackUpdateCommandHandler := stackupdatev1connect.NewStackUpdateCommandControllerHandler(stackUpdateService)
	mux.Handle(stackUpdateCommandPath, corsMiddleware(stackUpdateCommandHandler))

	// Register the StackUpdateQueryController (GetStackUpdate, ListStackUpdates, StreamStackUpdateOutput)
	stackUpdateQueryPath, stackUpdateQueryHandler := stackupdatev1connect.NewStackUpdateQueryControllerHandler(stackUpdateService)
	mux.Handle(stackUpdateQueryPath, corsMiddleware(stackUpdateQueryHandler))

	// Register the CredentialCommandController
	credentialCommandPath, credentialCommandHandler := credentialv1connect.NewCredentialCommandControllerHandler(credentialService)
	mux.Handle(credentialCommandPath, corsMiddleware(credentialCommandHandler))

	// Register the CredentialQueryController
	credentialQueryPath, credentialQueryHandler := credentialv1connect.NewCredentialQueryControllerHandler(credentialService)
	mux.Handle(credentialQueryPath, corsMiddleware(credentialQueryHandler))

	// Add health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := fmt.Sprintf(`{"status":"ok","service":"openmcf-backend","timestamp":"%s"}`,
			time.Now().Format(time.RFC3339))
		w.Write([]byte(response))
	})

	// Create HTTP server with h2c (HTTP/2 Cleartext) for gRPC
	httpServer := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: h2c.NewHandler(mux, &http2.Server{}),
	}

	return &Server{
		httpServer: httpServer,
		mongo:      cfg.MongoDB,
	}
}

// Start starts the server and blocks until it stops.
func (s *Server) Start() error {
	logrus.WithField("addr", s.httpServer.Addr).Info("Starting gRPC server")
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server failed: %w", err)
	}
	return nil
}

// Shutdown gracefully shuts down the server.
func (s *Server) Shutdown(ctx context.Context) error {
	logrus.Info("Shutting down server")
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown server: %w", err)
	}
	if err := s.mongo.Disconnect(ctx); err != nil {
		return fmt.Errorf("failed to disconnect MongoDB: %w", err)
	}
	return nil
}

// splitAndTrim splits a string by delimiter and trims whitespace from each part.
func splitAndTrim(s, delimiter string) []string {
	parts := []string{}
	for _, part := range split(s, delimiter) {
		trimmed := trim(part)
		if trimmed != "" {
			parts = append(parts, trimmed)
		}
	}
	return parts
}

// split is a helper to split strings (avoiding importing strings package)
func split(s, sep string) []string {
	var result []string
	var current string
	for i := 0; i < len(s); i++ {
		if i+len(sep) <= len(s) && s[i:i+len(sep)] == sep {
			result = append(result, current)
			current = ""
			i += len(sep) - 1
		} else {
			current += string(s[i])
		}
	}
	result = append(result, current)
	return result
}

// trim removes leading/trailing whitespace
func trim(s string) string {
	start := 0
	end := len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t' || s[start] == '\n' || s[start] == '\r') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\n' || s[end-1] == '\r') {
		end--
	}
	return s[start:end]
}
