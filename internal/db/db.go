package db

import (
	"context"
	"os"
	"sync"
	"time"

	"github.com/glebarez/sqlite"
	logger "github.com/rtfmkiesel/kisslog"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
)

var log = logger.New("db")

type Database struct {
	mu     *sync.Mutex
	sqlite *gorm.DB

	// To stop the background task for cleaning up expired secrets

	ctxCleanup  context.Context
	stopCleanup context.CancelFunc
	wgCleanup   *sync.WaitGroup
}

// Opens and migrates the database, starts a background task for cleaning expired secrets
func MustOpen() *Database {
	dbPath, exists := os.LookupEnv("DB_PATH")
	if !exists {
		dbPath = "./onetim3.sqlite" // default path
	}
	log.Debug("using sqlite database %s", dbPath)

	sqliteDatabase, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: glogger.Default.LogMode(glogger.Silent), // Stops gorms logger
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Debug("sqlite databased opened")

	if err := sqliteDatabase.AutoMigrate(&SecretString{}, &SecretFile{}); err != nil {
		log.Fatal(err)
	}
	log.Debug("sqlite database migrated")

	ctx, cancelFunc := context.WithCancel(context.Background())
	db := &Database{
		mu:          new(sync.Mutex),
		sqlite:      sqliteDatabase,
		ctxCleanup:  ctx,
		stopCleanup: cancelFunc,
		wgCleanup:   new(sync.WaitGroup),
	}

	// Periodically remove expired secrets
	db.startCleanup()

	log.Info("sqlite database ready")
	return db
}

// Stops the background worker for cleaning expired secrets and closes the database
func (db *Database) Close() {
	db.stopCleanup()
	db.wgCleanup.Wait()

	log.Debug("closing sqlite db")
	database, err := db.sqlite.DB()
	if err != nil {
		log.Error(err)
		return
	}

	if err := database.Close(); err != nil {
		log.Error(err)
		return
	}

	log.Info("sqlite database closed")
}

// Removes expired secrets based on model from the database
func (db *Database) startCleanup() {
	log.Debug("starting cleanup task for expired secrets")

	deleteExpiredSecrets := func(model any, name string) {
		db.mu.Lock()
		defer db.mu.Unlock()

		res := db.sqlite.Where("expires_at <= ?", time.Now()).
			Limit(1000).
			Delete(model)
		if res.Error != nil {
			log.Error("failed to delete expired %s: %s", name, res.Error)
		} else if res.RowsAffected > 0 {
			log.Info("deleted %d expired %s", res.RowsAffected, name)
		}
	}

	db.wgCleanup.Go(func() {
		ticker := time.NewTicker(time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-db.ctxCleanup.Done():
				log.Debug("stopping cleanup task")
				return
			case <-ticker.C:
				log.Debug("checking for expired secrets")
				deleteExpiredSecrets(&SecretString{}, "strings")
				deleteExpiredSecrets(&SecretFile{}, "files")
			}
		}
	})
}
