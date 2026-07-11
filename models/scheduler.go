package models

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"drop.plus.or.kr/config"
)

func StartAutoDeleteScheduler() {
	if config.AppConfig.AutoDeleteSeconds <= 0 {
		log.Println("[INFO] Auto delete is disabled (AUTO_DELETE_SECONDS <= 0)")
		return
	}

	log.Printf("[INFO] Auto delete scheduler started. Deleting files inactive for %d seconds.", config.AppConfig.AutoDeleteSeconds)

	ticker := time.NewTicker(6 * time.Hour)
	go func() {
		for range ticker.C {
			deleteOldFiles()
		}
	}()
	
	go deleteOldFiles()
}

func deleteOldFiles() {
	seconds := config.AppConfig.AutoDeleteSeconds
	modifier := fmt.Sprintf("-%d seconds", seconds)
	
	query := `SELECT uuid FROM files WHERE last_used_at < datetime('now', ?)`
	rows, err := DB.Query(query, modifier)
	if err != nil {
		log.Printf("[ERROR] Scheduler DB query failed: %v", err)
		return
	}
	defer rows.Close()

	var uuids []string
	for rows.Next() {
		var uuid string
		if err := rows.Scan(&uuid); err != nil {
			log.Printf("[WARN] Scan error in scheduler: %v", err)
			continue
		}
		uuids = append(uuids, uuid)
	}
	if err := rows.Err(); err != nil {
		log.Printf("[ERROR] Scheduler rows error: %v", err)
	}

	for _, uuid := range uuids {
		log.Printf("[INFO] Auto-deleting inactive file: %s", uuid)
		
		filePath := filepath.Join(config.AppConfig.DataDir, "uploads", uuid)
		os.Remove(filePath)
		
		DB.Exec(`DELETE FROM files WHERE uuid = ?`, uuid)
	}
}
