package handlers

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"

	"drop.plus.or.kr/config"
	"drop.plus.or.kr/models"
)

func HandleUploadFile(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(config.AppConfig.MaxUploadSize)
	if err != nil {
		redirectWithMessage(w, r, "파일 업로드 실패 (용량 초과)")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		redirectWithMessage(w, r, "파일을 찾을 수 없습니다")
		return
	}
	defer file.Close()

	newUUID := uuid.New().String()
	uploader, ok := r.Context().Value(UserEmailKey).(string)
	if !ok || uploader == "" {
		redirectWithMessage(w, r, "인증이 필요합니다.")
		return
	}

	savePath := filepath.Join(config.AppConfig.DataDir, "uploads", newUUID)
	dst, err := os.Create(savePath)
	if err != nil {
		redirectWithMessage(w, r, "파일 저장 실패")
		return
	}
	defer dst.Close()

	size, err := io.Copy(dst, file)
	if err != nil {
		redirectWithMessage(w, r, "파일 복사 실패")
		return
	}

	query := `INSERT INTO files (uuid, original_name, uploaded_by, size, uploaded_at) VALUES (?, ?, ?, ?, ?)`
	_, err = models.DB.Exec(query, newUUID, header.Filename, uploader, size, time.Now())
	if err != nil {
		os.Remove(savePath)
		redirectWithMessage(w, r, "DB 저장 실패")
		return
	}

	redirectWithMessage(w, r, "파일이 성공적으로 업로드 되었습니다")
}

func HandleDeleteFile(w http.ResponseWriter, r *http.Request) {
	fileUUID := r.PathValue("uuid")

	uploader, ok := r.Context().Value(UserEmailKey).(string)
	if !ok || uploader == "" {
		redirectWithMessage(w, r, "인증이 필요합니다.")
		return
	}

	res, err := models.DB.Exec(`DELETE FROM files WHERE uuid = ? AND uploaded_by = ?`, fileUUID, uploader)
	if err != nil {
		redirectWithMessage(w, r, "삭제 실패")
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		redirectWithMessage(w, r, "파일이 없거나 삭제 권한이 없습니다.")
		return
	}

	os.Remove(filepath.Join(config.AppConfig.DataDir, "uploads", fileUUID))

	redirectWithMessage(w, r, "파일이 삭제되었습니다")
}

func HandleDownloadFile(w http.ResponseWriter, r *http.Request) {
	fileUUID := r.PathValue("uuid")

	var originalName string
	err := models.DB.QueryRow(`SELECT original_name FROM files WHERE uuid = ?`, fileUUID).Scan(&originalName)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	filePath := filepath.Join(config.AppConfig.DataDir, "uploads", fileUUID)
	
	models.DB.Exec(`UPDATE files SET last_used_at = CURRENT_TIMESTAMP WHERE uuid = ?`, fileUUID)
	
	if r.URL.Query().Get("download") == "true" {
		w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, originalName))
	} else {
		w.Header().Set("Content-Disposition", fmt.Sprintf(`inline; filename="%s"`, originalName))
	}

	http.ServeFile(w, r, filePath)
}

func redirectWithMessage(w http.ResponseWriter, r *http.Request, msg string) {
	http.Redirect(w, r, "/?msg="+url.QueryEscape(msg), http.StatusSeeOther)
}
