package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"cloud.google.com/go/storage"
	firebase "firebase.google.com/go/v4"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/sportspazz/api/client"
	"github.com/sportspazz/cmd/server"
	"github.com/sportspazz/configs"
	"google.golang.org/api/option"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&TimeZone=UTC",
		configs.Envs.DBUser,
		configs.Envs.DBPassword,
		configs.Envs.DBHost,
		configs.Envs.DBPort,
		configs.Envs.DBName,
	)

	db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err := db.Exec("SELECT 1").Error; err != nil {
		logger.Error("cannot ping database", slog.Any("err", err))
		os.Exit(1)
	}
	logger.Info("Connected to database", slog.String("database", configs.Envs.DBName))

	m, err := migrate.New("file://"+configs.Envs.DBMigrationDir, dsn)
	if err != nil {
		logger.Error("cannot initialize db migration", slog.Any("err", err))
		os.Exit(1)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logger.Error("error applying migrations", slog.Any("err", err))
		os.Exit(1)
	}

	ctx := context.Background()
	firebaseApp, err := firebase.NewApp(ctx, nil, option.WithCredentialsFile(configs.Envs.GCPApiKey))
	if err != nil {
		logger.Error("error initializing firebase app", slog.Any("err", err))
		os.Exit(1)
	}
	firebaseRest := client.NewFirebaseClient(configs.Envs.FirebaseApiKey, configs.Envs.FirebaseProjectID, logger)
	storageClient, err := storage.NewClient(ctx, option.WithCredentialsFile(configs.Envs.GCPApiKey))
	if err != nil {
		logger.Error("error initializing cloud storage", slog.Any("err", err))
		os.Exit(1)
	}

	killSig := make(chan os.Signal, 1)
	signal.Notify(killSig, os.Interrupt, syscall.SIGTERM)

	go func() {
		server := server.NewServer(configs.Envs.Host, configs.Envs.Port, db, firebaseApp, firebaseRest, storageClient)
		if err := server.Run(); err != nil {
			logger.Error("Cannot start server", slog.Any("err", err))
		}
	}()

	logger.Info("Server started", slog.String("port", configs.Envs.Port))
	<-killSig

}
