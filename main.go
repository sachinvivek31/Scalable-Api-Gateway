package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chi_middleware "github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/sachinvivek31/api-gateway/internal/config"
	my_middleware "github.com/sachinvivek31/api-gateway/internal/middleware"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}

	r := chi.NewRouter()

	// Global Middlewares
	r.Use(my_middleware.RequestID)
	r.Use(chi_middleware.Logger)
	r.Use(chi_middleware.Recoverer)

	// EXPOSE METRICS ENDPOINT
	r.Handle("/metrics", promhttp.Handler())

	for _, svc := range cfg.Services {
		targetURL, _ := url.Parse(svc.Target)
		proxy := httputil.NewSingleHostReverseProxy(targetURL)

		// Set up the Proxy Director
		originalDirector := proxy.Director
		proxy.Director = func(req *http.Request) {
			originalDirector(req)
			if reqID, ok := req.Context().Value(my_middleware.RequestIDKey).(string); ok {
				req.Header.Set("X-Request-ID", reqID)
			}
			req.Host = targetURL.Host
		}

		// CHAINING LOGIC: Start with the proxy
		var handler http.Handler = proxy

		// 1. Apply Rate Limiting (we could customize this per service using svc.RateLimit)
		handler = my_middleware.RateLimiter(handler)

		// 2. Apply Auth only if required by config
		if svc.RequiresAuth {
			handler = my_middleware.Authenticate(handler)
		}

		r.Mount(svc.Prefix, handler)
		fmt.Printf("Mapped %s (Auth: %v) -> %s\n", svc.Prefix, svc.RequiresAuth, svc.Target)
	}

	// Standard Graceful Shutdown logic below...
	server := &http.Server{Addr: fmt.Sprintf(":%d", cfg.Server.Port), Handler: r}
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		fmt.Printf("Gateway live on :%d\n", cfg.Server.Port)
		server.ListenAndServe()
	}()

	<-stop
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	server.Shutdown(ctx)
	fmt.Println("Gateway gracefully stopped.")
}