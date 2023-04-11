package main

import (
	"context"
	"errors"
	"fmt"
	conf "github.com/ardanlabs/conf/v2"
	"github.com/cpustejovsky/mongotest/server"
	"github.com/cpustejovsky/mongotest/store"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var build = "develop"

func main() {
	// =========================================================================
	// Configuration
	cfg := struct {
		conf.Version
		Web struct {
			APIHost         string        `conf:"default:0.0.0.0:80"`
			ReadTimeout     time.Duration `conf:"default:5s"`
			WriteTimeout    time.Duration `conf:"default:10s"`
			IdleTimeout     time.Duration `conf:"default:120s"`
			ShutdownTimeout time.Duration `conf:"default:20s"`
			MongoUri        string        `conf:"default:mongodb://mongo:27017,mask"`
		}
	}{
		Version: conf.Version{
			Build: build,
			Desc:  "cpustejovsky MIT license",
		},
	}

	const prefix = "DEFAULT"
	help, err := conf.Parse(prefix, &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			os.Exit(0)
		}
		fmt.Println("error parsing config:\t", err)
		os.Exit(1)
	}

	// =========================================================================
	// Database Setup
	log.Println(cfg.Web.MongoUri)
	clientOptions := options.Client().
		ApplyURI(cfg.Web.MongoUri)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(ctx)
	log.Println("Successfully connected to database; URI:\t", cfg.Web.MongoUri)
	serverErrors := make(chan error)
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	store := store.NewAnimalStore(client, "mongotest")
	animalSrv := server.New(store)

	go func() {
		log.Printf("Starting server on %s", cfg.Web.APIHost)
		serverErrors <- http.ListenAndServe(cfg.Web.APIHost, animalSrv)
	}()
	select {
	case err := <-serverErrors:
		fmt.Errorf("server error: %w", err)
		os.Exit(1)

	case sig := <-shutdown:
		log.Println("shutdown started with signal:\t", sig)
		defer log.Println("shutdown complete with signal:\t", sig)

		// Give outstanding requests a deadline for completion.
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancel()
		log.Println("Cancelled with context error:\t", ctx.Err())
	}

	log.Println("stopping service")
}
