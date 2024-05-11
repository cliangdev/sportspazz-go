package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	firebase "firebase.google.com/go/v4"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/sportspazz/cmd/server"
	"github.com/sportspazz/configs"
	"google.golang.org/api/option"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true",
		configs.Envs.DBUser,
		configs.Envs.DBPassword,
		configs.Envs.DBAddress,
		configs.Envs.DBName,
	)

	db, _ := gorm.Open(mysql.New(mysql.Config{
		DSN: dsn,
	}), &gorm.Config{})

	if err := db.Exec("SELECT 1").Error; err != nil {
		logger.Error("cannot ping database", slog.Any("err", err))
		os.Exit(1)
	}
	logger.Info("Connected to database", slog.String("database", configs.Envs.DBName))

	m, err := migrate.New("file://"+configs.Envs.DBMigrationDir, "mysql://"+dsn)
	if err != nil {
		logger.Error("cannot initialize db migration", slog.Any("err", err))
		os.Exit(1)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logger.Error("error applying migrations", slog.Any("err", err))
		os.Exit(1)
	}

	opt := option.WithCredentialsFile(configs.Envs.FirebaseServiceAccountJson)
	firebaseApp, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		logger.Error("error initializing firebase app", slog.Any("err", err))
		os.Exit(1)
	}

	killSig := make(chan os.Signal, 1)
	signal.Notify(killSig, os.Interrupt, syscall.SIGTERM)

	go func() {
		server := server.NewServer(configs.Envs.Host, configs.Envs.Port, db, firebaseApp)
		if err := server.Run(); err != nil {
			logger.Error("Cannot start server", slog.Any("err", err))
		}
	}()

	logger.Info("Server started", slog.String("port", configs.Envs.Port))
	<-killSig

}
