package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strconv"

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

	uploader, ok := r.Context().Value(UserEmailKey).(string)
	if !ok || uploader == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// 1. Pagination 변수 설정
	pageStr := r.URL.Query().Get("page")
	currentPage, err := strconv.Atoi(pageStr)
	if err != nil || currentPage < 1 {
		currentPage = 1
	}
	itemsPerPage := 5 // 한 페이지당 보여줄 파일 수 (테스트를 위해 5개로 설정)

	// 2. 전체 파일 수 카운트
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

	// 3. DB 조회 (Limit & Offset 적용)
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

	// 4. 페이지네이션 배열 계산 (현재 페이지 기준 앞뒤 2개씩 노출)
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
