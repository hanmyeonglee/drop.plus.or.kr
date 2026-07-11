package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

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
	tmplName := "templates/index.html"
	ua := strings.ToLower(r.UserAgent())
	if strings.Contains(ua, "mobile") || strings.Contains(ua, "android") || strings.Contains(ua, "iphone") {
		tmplName = "templates/index_mobile.html"
	}
	
	tmpl, err := template.ParseFiles(tmplName)
	if err != nil {
		log.Printf("Template parsing error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	uploader, ok := r.Context().Value(UserEmailKey).(string)
	if !ok || uploader == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	pageStr := r.URL.Query().Get("page")
	currentPage, err := strconv.Atoi(pageStr)
	if err != nil || currentPage < 1 {
		currentPage = 1
	}
	itemsPerPage := 5

	var totalFiles int
	err = models.DB.QueryRow(`SELECT COUNT(*) FROM files WHERE uploaded_by = ?`, uploader).Scan(&totalFiles)
	if err != nil {
		log.Printf("Failed to count files: %v", err)
		totalFiles = 0
	}

	totalPages := (totalFiles + itemsPerPage - 1) / itemsPerPage
	if currentPage > totalPages && totalPages > 0 {
		currentPage = totalPages
	}

	offset := (currentPage - 1) * itemsPerPage
	if offset < 0 {
		offset = 0
	}

	rows, err := models.DB.Query(`SELECT uuid, original_name, size, uploaded_at FROM files WHERE uploaded_by = ? ORDER BY uploaded_at DESC LIMIT ? OFFSET ?`, uploader, itemsPerPage, offset)

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

	var pages []int
	startPage := currentPage - 2
	endPage := currentPage + 2

	if startPage < 1 {
		endPage += (1 - startPage)
		startPage = 1
	}
	if endPage > totalPages {
		startPage -= (endPage - totalPages)
		endPage = totalPages
		if startPage < 1 {
			startPage = 1
		}
	}
	for i := startPage; i <= endPage; i++ {
		pages = append(pages, i)
	}

	prevPage := 0
	if currentPage > 1 {
		prevPage = currentPage - 1
	}

	nextPage := 0
	if currentPage < totalPages {
		nextPage = currentPage + 1
	}

	data := struct {
		UserEmail   string
		Message     string
		Files       []FileData
		HasPages    bool
		CurrentPage int
		PrevPage    int
		NextPage    int
		Pages       []int
	}{
		UserEmail:   uploader,
		Message:     msg,
		Files:       files,
		HasPages:    totalPages > 1,
		CurrentPage: currentPage,
		PrevPage:    prevPage,
		NextPage:    nextPage,
		Pages:       pages,
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Template execution error: %v", err)
	}
}
