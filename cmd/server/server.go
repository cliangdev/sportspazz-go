package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"cloud.google.com/go/storage"
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
	port           string
	db             *gorm.DB
	firebaseApp    *firebase.App
	firebaseClient *client.FirebaseClient
	storageClient  *storage.Client
	bucket         string
}

func NewServer(
	port string,
	db *gorm.DB,
	firebaseApp *firebase.App,
	firebaseClient *client.FirebaseClient,
	storageClient *storage.Client,
	bucket string) *Server {

	return &Server{
		port:           port,
		db:             db,
		firebaseApp:    firebaseApp,
		firebaseClient: firebaseClient,
		storageClient:  storageClient,
		bucket:         bucket,
	}
}

func (s *Server) Run() error {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	router := mux.NewRouter()
	subRouter := router.PathPrefix("/api/v1").Subrouter()

	ctx := context.Background()
	// middlewares
	firebaseAdminClient, err := s.firebaseApp.Auth(ctx)
	if err != nil {
		logger.Error("Cannot initialize Firebase admin client", slog.Any("err", err))
		os.Exit(1)
	}
	router.Use(
		middleware.LoggerMiddleWare(logger),
		middleware.ContentTypeHeaderMiddleWare,
		middleware.AuthenticateMiddleWare(s.firebaseClient, logger),
	)

	// REST API handler
	userStore := user.NewUserStore(s.db, logger)
	userService := user.NewUserService(userStore, firebaseAdminClient, logger)
	userHandler := rest_api.NewUserHandler(userService)
	userHandler.RegisterRoutes(subRouter)

	poiStore := poi.NewPoiStore(s.db, logger)
	poiService := poi.NewPoiService(poiStore, logger)
	poiHandler := rest_api.NewPoiHandler(poiService, s.firebaseClient)
	poiHandler.RegisterRoutes(subRouter)

	// HTML handler
	homeHandler := web.NewHomeHandler(logger)
	homeHandler.RegisterRoutes(router)

	registerHandler := web.NewRegisterHandler(userService, logger)
	registerHandler.RegisterRoutes(router)

	loginHandler := web.NewLoginHandler(userService, s.firebaseClient, logger)
	loginHandler.RegisterRoutes(router)

	whereToPlay := web.NewWhereToPlayHandler(logger, poiService, s.storageClient, s.bucket)
	whereToPlay.RegisterRoutes(router)

	fs := http.FileServer(http.Dir("./static"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		path, _ := route.GetPathTemplate()
		methods, _ := route.GetMethods()

		logger.Info(fmt.Sprintf("...%v: %s", methods, path))
		return nil
	})

	return http.ListenAndServe(fmt.Sprintf(":%s", s.port), router)
}
