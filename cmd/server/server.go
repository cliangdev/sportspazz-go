package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	firebase "firebase.google.com/go/v4"
	"github.com/gorilla/mux"
	rest_api "github.com/sportspazz/api/rest"
	web "github.com/sportspazz/api/web"
	"github.com/sportspazz/middleware"
	"github.com/sportspazz/service/user"
	"gorm.io/gorm"
)

type Server struct {
	host        string
	port        string
	db          *gorm.DB
	firebaseApp *firebase.App
}

func NewServer(host, port string, db *gorm.DB, firebaseApp *firebase.App) *Server {
	return &Server{
		host:        host,
		port:        port,
		db:          db,
		firebaseApp: firebaseApp,
	}
}

func (s *Server) Run() error {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	router := mux.NewRouter()
	subRouter := router.PathPrefix("/api/v1").Subrouter()

	// middlewares
	router.Use(
		middleware.LoggerMiddleWare(logger),
		middleware.ContentTypeHeaderMiddleWare,
	)
	// REST API handler
	firebaseClient, err := s.firebaseApp.Auth(context.Background())
	if err != nil {
		logger.Error("Cannot ping database", slog.Any("err", err))
		os.Exit(1)
	}

	store := user.NewUserStoe(s.db, logger)
	userService := user.NewUserService(store, firebaseClient, logger)
	userHandler := rest_api.NewUserHandler(userService)
	userHandler.RegisterRoutes(subRouter)

	// Templ HTMX handler
	registerHandler := web.NewRegisterHandler(userService, logger)
	registerHandler.RegisterRoutes(router)

	loginHandler := web.NewLoginHandler(userService, logger)
	loginHandler.RegisterRoutes(router)

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("public")))

	router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		path, _ := route.GetPathTemplate()
		methods, _ := route.GetMethods()

		logger.Info(fmt.Sprintf("...%v: %s", methods, path))
		return nil
	})

	return http.ListenAndServe(fmt.Sprintf("%s:%s", s.host, s.port), router)
}
