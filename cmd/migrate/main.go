package main

import (
	"context"
	"log"
	"time"

	"ppharma/backend/internal/config"
	mongowrap "ppharma/backend/support-pkg/db/mongo"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Migration failed: config load error: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Printf("Connecting to MongoDB at: %s", cfg.DB.DBURI)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.DB.DBURI))
	if err != nil {
		log.Fatalf("Migration failed: unable to connect to MongoDB: %v", err)
	}
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatalf("Migration cleanup failed: %v", err)
		}
	}()

	db := client.Database(cfg.DB.DBName)

	log.Println("Running MongoDB migrations...")

	// ensure explicit collection creations (optional as indexing also enforces this organically)
	collections := []string{
		"users",
		"customers",
		"common",
		"doctors",
		"consultations",
		"products",
		"orders",
		"payments",
	}

	for _, col := range collections {
		// MongoDB ignores creation if it actively exists natively 
		db.CreateCollection(ctx, col)
		log.Printf("Checked/Created collection: %s", col)
	}

	log.Println("Applying strict indexing rules...")
	if err := mongowrap.EnsureIndexes(ctx, db); err != nil {
		log.Fatalf("Migration failed while enforcing indexes: %v", err)
	}

	log.Println("Database migration completed successfully!")
}
