package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/stone1549/yapyapyap/auth/service"
	"log"
	"net/http"
)

func main() {
	flag.Parse()

	config, err := service.GetConfiguration()

	if err != nil {
		panic(fmt.Sprintf("Unable to load configuration: %s", err.Error()))
	}

	configMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			ctx := context.WithValue(request.Context(), "config", config)
			next.ServeHTTP(writer, request.WithContext(ctx))
		})
	}

	repo, err := service.NewUserRepository(config)

	if err != nil {
		panic(fmt.Sprintf("Unable to configure repository: %s", err.Error()))
	}

	repoMiddleWare := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), "repo", repo)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}

	tokenFactory, err := service.NewTokenFactory(config)

	if err != nil {
		panic(fmt.Sprintf("Unable to configure token factory: %s", err.Error()))
	}

	tokenMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), "tokenFactory", tokenFactory)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
	r := chi.NewRouter()

	// Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(configMiddleware)
	r.Use(repoMiddleWare)
	r.Use(tokenMiddleware)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(config.GetTimeout()))

	r.Route("/session", func(r chi.Router) {
		r.With(service.NewSessionMiddleware).Put("/", service.NewSession)
		r.With(service.JwtAuthMiddleware).With(service.RefreshSessionMiddleware).Get("/", service.RefreshSession)
	})

	r.Route("/user", func(r chi.Router) {
		r.With(service.NewUserMiddleware).Put("/", service.NewUser)
		r.With(service.GetUserMiddleware).Get("/{id}", service.GetUser)
		r.With(service.UpdateProfileMiddleware).Patch("/{id}", service.UpdateProfile)
	})

	err = http.ListenAndServe(":3333", r)

	if err != nil {
		log.Println(err)
	}
}
