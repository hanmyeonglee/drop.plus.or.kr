package handlers

import (
	"fmt"
	"net/http"
)

func HandleIndexPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "메인 페이지 렌더링 (로그인 상태에 따라 다름) (구현 예정)")
}
