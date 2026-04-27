package config

import "fmt"

func RunMigrations() {
	db := DB

	queries := []string{
		`CREATE INDEX IF NOT EXISTS idx_rides_driver_id ON rides(driver_id);`,
		`CREATE INDEX IF NOT EXISTS idx_rides_status ON rides(status);`,
	}

	for _, q := range queries {
		err := db.Exec(q).Error
		if err != nil {
			fmt.Println("Migration failed:", err)
		}
	}
}