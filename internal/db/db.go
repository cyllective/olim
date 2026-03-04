package db

import (
	"context"
	"os"
	"sync"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
)

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
		dbPath = "./olim.sqlite" // default path
	}
	log.Debug().Msgf("using sqlite database %s", dbPath)

	sqliteDatabase, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: glogger.Default.LogMode(glogger.Silent), // Stops gorms logger
	})
	if err != nil {
		log.Fatal().Err(err)
	}
	log.Debug().Msg("sqlite databased opened")

	if err := sqliteDatabase.AutoMigrate(&SecretString{}, &SecretFile{}); err != nil {
		log.Fatal().Err(err)
	}
	log.Debug().Msg("sqlite database migrated")

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

	log.Info().Msg("sqlite database ready")
	return db
}

// Stops the background worker for cleaning expired secrets and closes the database
func (db *Database) Close() {
	db.stopCleanup()
	db.wgCleanup.Wait()

	log.Debug().Msg("closing sqlite db")
	database, err := db.sqlite.DB()
	if err != nil {
		log.Error().Err(err)
		return
	}

	if err := database.Close(); err != nil {
		log.Error().Err(err)
		return
	}

	log.Info().Msg("sqlite database closed")
}

// Removes expired secrets based on model from the database
func (db *Database) startCleanup() {
	log.Debug().Msg("starting cleanup task for expired secrets")

	deleteExpiredSecrets := func(model any, name string) {
		db.mu.Lock()
		defer db.mu.Unlock()

		res := db.sqlite.Where("expires_at <= ?", time.Now()).
			Limit(1000).
			Delete(model)
		if res.Error != nil {
			log.Error().Msgf("failed to delete expired %s: %s", name, res.Error)
		} else if res.RowsAffected > 0 {
			log.Info().Msgf("deleted %d expired %s", res.RowsAffected, name)
		}
	}

	db.wgCleanup.Go(func() {
		ticker := time.NewTicker(time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-db.ctxCleanup.Done():
				log.Debug().Msg("stopping cleanup task")
				return
			case <-ticker.C:
				log.Debug().Msg("checking for expired secrets")
				deleteExpiredSecrets(&SecretString{}, "strings")
				deleteExpiredSecrets(&SecretFile{}, "files")
			}
		}
	})
}
