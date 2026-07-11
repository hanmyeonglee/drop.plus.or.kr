package handlers

import (
	"fmt"
	"net/http"
)

func HandleUploadFile(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "파일 업로드 처리 (구현 예정)")
}

func HandleDownloadFile(w http.ResponseWriter, r *http.Request) {
	uuid := r.PathValue("uuid")
	fmt.Fprintf(w, "파일 다운로드 처리: UUID=%s (구현 예정)\n", uuid)
}

func HandleDeleteFile(w http.ResponseWriter, r *http.Request) {
	uuid := r.PathValue("uuid")
	fmt.Fprintf(w, "파일 삭제 처리: UUID=%s (구현 예정)\n", uuid)
}
