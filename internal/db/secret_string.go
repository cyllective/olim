package db

import "time"

// The struct on how a secret string is stored in the SQLite database
type SecretString struct {
	ID        string    `gorm:"type:TEXT;primaryKey"` // UUIDv4
	Value     []byte    `gorm:"type:BYTE"`            // The encrypted string
	CreatedAt time.Time `gorm:"autoCreateTime;index"`
	ExpiresAt time.Time `gorm:"index"`
}

// Adds a new secret to the SQLite database
func (db *Database) AddSecretString(secret *SecretString) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	return db.sqlite.Create(secret).Error
}

// Retrieves and then deletes a secret from the SQLite database by its ID
func (db *Database) GetAndDeleteSecretStringByID(id string) (*SecretString, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	secret := &SecretString{}
	if err := db.sqlite.First(secret, "ID = ?", id).Error; err != nil {
		return nil, err
	}

	if err := db.sqlite.Delete(&SecretString{}, "ID = ?", id).Error; err != nil {
		return secret, err
	}

	return secret, nil
}
