package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"

	"drop.plus.or.kr/models"
)

type FileData struct {
	UUID         string
	OriginalName string
	Size         string
	UploadedAt   string
}

func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}

	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func HandleIndexPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Printf("Template parsing error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	uploader := "testuser@drop.plus.or.kr" // Dummy auth

	rows, err := models.DB.Query(`SELECT uuid, original_name, size, uploaded_at FROM files WHERE uploaded_by = ? ORDER BY uploaded_at DESC`, uploader)
	
	var files []FileData
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var f FileData
			var size int64
			var uploadedAt string
			if err := rows.Scan(&f.UUID, &f.OriginalName, &size, &uploadedAt); err == nil {
				f.Size = formatSize(size)
				if len(uploadedAt) >= 16 {
					f.UploadedAt = uploadedAt[:16]
				} else {
					f.UploadedAt = uploadedAt
				}
				files = append(files, f)
			}
		}
		if err = rows.Err(); err != nil {
			log.Printf("Rows iteration error: %v", err)
		}
	} else {
		log.Printf("DB query error: %v", err)
	}

	msg := r.URL.Query().Get("msg")
	if msg != "" {
		msg, _ = url.QueryUnescape(msg)
	}

	data := struct {
		UserEmail string
		Message   string
		Files     []FileData
		HasPages  bool
	}{
		UserEmail: uploader,
		Message:   msg,
		Files:     files,
		HasPages:  false,
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Template execution error: %v", err)
	}
}
