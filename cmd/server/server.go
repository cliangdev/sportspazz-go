package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	firebase "firebase.google.com/go/v4"
	"github.com/gorilla/mux"
	"github.com/sportspazz/api/client"
	rest_api "github.com/sportspazz/api/rest"
	web "github.com/sportspazz/api/web"
	"github.com/sportspazz/middleware"
	"github.com/sportspazz/service/poi"
	"github.com/sportspazz/service/user"
	"gorm.io/gorm"
)

type Server struct {
	host         string
	port         string
	db           *gorm.DB
	firebaseApp  *firebase.App
	firebaseRest *client.FirebaseClient
}

func NewServer(host, port string, db *gorm.DB, firebaseApp *firebase.App, firebaseRest *client.FirebaseClient) *Server {
	return &Server{
		host:         host,
		port:         port,
		db:           db,
		firebaseApp:  firebaseApp,
		firebaseRest: firebaseRest,
	}
}

func (s *Server) Run() error {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	router := mux.NewRouter()
	subRouter := router.PathPrefix("/api/v1").Subrouter()

	// middlewares
	firebaseAdminClient, err := s.firebaseApp.Auth(context.Background())
	if err != nil {
		logger.Error("Cannot initialize Firebase admin client", slog.Any("err", err))
		os.Exit(1)
	}
	router.Use(
		middleware.LoggerMiddleWare(logger),
		middleware.ContentTypeHeaderMiddleWare,
		middleware.AuthenticateMiddleWare(s.firebaseRest, logger),
	)

	// REST API handler
	userStore := user.NewUserStore(s.db, logger)
	userService := user.NewUserService(userStore, firebaseAdminClient, logger)
	userHandler := rest_api.NewUserHandler(userService)
	userHandler.RegisterRoutes(subRouter)

	poiStore := poi.NewPoiStore(s.db, logger)
	poiService := poi.NewPoiService(poiStore, logger)
	poiHandler := rest_api.NewPoiHandler(poiService, s.firebaseRest)
	poiHandler.RegisterRoutes(subRouter)

	// HTML handler
	homeHandler := web.NewHomeHandler(logger)
	homeHandler.RegisterRoutes(router)

	registerHandler := web.NewRegisterHandler(userService, logger)
	registerHandler.RegisterRoutes(router)

	loginHandler := web.NewLoginHandler(userService, s.firebaseRest, logger)
	loginHandler.RegisterRoutes(router)

	whereToPlay := web.NewWhereToPlayHandler(logger, poiService)
	whereToPlay.RegisterRoutes(router)

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("public")))

	router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		path, _ := route.GetPathTemplate()
		methods, _ := route.GetMethods()

		logger.Info(fmt.Sprintf("...%v: %s", methods, path))
		return nil
	})

	return http.ListenAndServe(fmt.Sprintf("%s:%s", s.host, s.port), router)
}
