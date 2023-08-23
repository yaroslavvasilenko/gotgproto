package storage

import (
	"os"
	"time"

	"github.com/AnimeKaizoku/cacher"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var SESSION *gorm.DB

func Load(sessionName string, inMemory bool) error {
	loadCache(inMemory)
	if inMemory {
		return nil
	}

	// Create a new file if it doesn't exist
	if _, err := os.Stat(sessionName); os.IsNotExist(err) {
		file, createErr := os.Create(sessionName)
		if createErr != nil {
			return createErr
		}
		defer file.Close()
	}

	db, err := gorm.Open(sqlite.Open(sessionName), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return err
	}
	SESSION = db
	dB, _ := db.DB()
	dB.SetMaxOpenConns(100)

	// Create tables if they don't exist
	_ = SESSION.AutoMigrate(&Session{}, &Peer{})
	return nil
}

func loadCache(inMemory bool) {
	var opts *cacher.NewCacherOpts
	if inMemory {
		storeInMemory = true
		opts = nil
	} else {
		opts = &cacher.NewCacherOpts{
			TimeToLive:    6 * time.Hour,
			CleanInterval: 24 * time.Hour,
			Revaluate:     true,
		}
	}
	peerCache = cacher.NewCacher[int64, *Peer](opts)
}
