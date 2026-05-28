package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/MohdAdamSiThuThetNaing/automotive-aftermarket/internal/builder"
	"github.com/MohdAdamSiThuThetNaing/automotive-aftermarket/internal/cache"
	"github.com/MohdAdamSiThuThetNaing/automotive-aftermarket/internal/loader"
	"github.com/MohdAdamSiThuThetNaing/automotive-aftermarket/internal/repository"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)

	defer cancel()

	dbURL := os.Getenv("DATABASE_URL")

	if dbURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	dbConfig, err := pgxpool.ParseConfig(dbURL)

	if err != nil {
		log.Fatalf("parse database config: %v", err)
	}

	dbConfig.MaxConns = 10
	dbConfig.MinConns = 2
	dbConfig.MaxConnLifetime = 30 * time.Minute
	dbConfig.MaxConnIdleTime = 5 * time.Minute

	db, err := pgxpool.NewWithConfig(ctx, dbConfig)

	if err != nil {
		log.Fatalf("connect database: %v", err)
	}

	defer db.Close()
	if err := db.Ping(ctx); err != nil {
		log.Fatalf("ping database: %v", err)
	}

	log.Printf( "connected to PostgreSQL max_conns=%d", db.Stat().MaxConns(), )

	// Initialize cache
	translationCache := cache.NewTranslationCache(
		5 * time.Minute,
	)

	// Initialize translation loader
	translationLoader := loader.NewPostgresLoader(
		db,
		translationCache,
	)

	// Initialize repositories
	productRepository := repository.NewPostgresProductRepository(
		db,
	)

	// Initialize document builder
	productDocumentBuilder := builder.NewProductDocumentBuilder(
		productRepository,
		translationLoader,
	)

	// Example product_ID
	productID := os.Getenv("PRODUCT_ID")

	if productID == "" {
		log.Fatal("PRODUCT_ID is required")
	}

	// Elasticsearch document
	document, err := productDocumentBuilder.Build(
		ctx,
		productID,
		[]string{
			"en",
			"th",
		},
	)

	if err != nil {
		log.Fatalf("build product document: %v", err)
	}

	jsonData, err := json.MarshalIndent(
		document,
		"",
		"  ",
	)

	if err != nil {
		log.Fatalf(
			"marshal document: %v",
			err,
		)
	}

	fmt.Println(
		"Generated Product Document:",
	)

	fmt.Println(string(jsonData))

}