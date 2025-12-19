package db

import "time"

// The struct on how a secret file is stored in the SQLite database
type SecretFile struct {
	ID        string    `gorm:"type:TEXT;primaryKey"` // UUIDv4
	Name      string    `gorm:"type:BYTE"`            // The encrypted filename
	Value     []byte    `gorm:"type:BYTE"`            // The encrypted file bytes
	CreatedAt time.Time `gorm:"autoCreateTime;index"`
	ExpiresAt time.Time `gorm:"index"`
}

// Adds a new secret file to the SQLite database
func (db *Database) AddSecretFile(secret *SecretFile) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	return db.sqlite.Create(secret).Error
}

// Retrieves and then deletes a secret file from the SQLite database by its ID
func (db *Database) GetAndDeleteSecretFileByID(id string) (*SecretFile, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	secret := &SecretFile{}
	if err := db.sqlite.First(secret, "ID = ?", id).Error; err != nil {
		return nil, err
	}

	if err := db.sqlite.Delete(&SecretFile{}, "ID = ?", id).Error; err != nil {
		return secret, err
	}

	return secret, nil
}
