package main

import (
	"context"
	"net/http"

	"github.com/Rindrics/go_todo_app/auth"
	"github.com/Rindrics/go_todo_app/clock"
	"github.com/Rindrics/go_todo_app/config"
	"github.com/Rindrics/go_todo_app/handler"
	"github.com/Rindrics/go_todo_app/service"
	"github.com/Rindrics/go_todo_app/store"
	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
)

func NewMux(ctx context.Context, cfg *config.Config) (http.Handler, func(), error) {
	mux := chi.NewRouter()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json; charset=utf-8")
		_, _ = w.Write([]byte(`{"status": "ok"}`))
	})
	v := validator.New()
	db, cleanup, err := store.New(ctx, cfg)
	if err != nil {
		return nil, cleanup, err
	}
	clocker := clock.RealClocker{}
	r := &store.Repository{Clocker: clocker}
	rcli, err := store.NewKVS(ctx, *cfg)
	if err != nil {
		return nil, cleanup, err
	}
	jwter, err := auth.NewJWTer(rcli, clocker)
	if err != nil {
		return nil, cleanup, err
	}
	l := &handler.Login{
		Service: &service.Login{
			DB:             db,
			Repo:           r,
			TokenGenerator: jwter,
		},
		Validator: v,
	}
	mux.Post("/login", l.ServeHTTP)

	at := &handler.AddTask{
		Service:   &service.AddTask{DB: db, Repo: r},
		Validator: v,
	}
	lt := &handler.ListTask{
		Service: &service.ListTask{DB: db, Repo: r},
	}
	mux.Route("/tasks", func(r chi.Router) {
		r.Use(handler.AuthMiddleware(jwter))
		r.Post("/", at.ServeHTTP)
		r.Get("/", lt.ServeHTTP)
	})
	ru := &handler.RegisterUser{
		Service:   &service.RegisterUser{DB: db, Repo: r},
		Validator: v,
	}
	mux.Post("/register", ru.ServeHTTP)

	mux.Route("/admin", func(r chi.Router) {
		r.Use(handler.AuthMiddleware(jwter), handler.AdminMiddleware)
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-type", "application/json; charset=utf-8")
			_, _ = w.Write([]byte(`{"message":"admin only"}`))
		})
	})

	return mux, cleanup, nil
}
