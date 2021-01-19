package ego

import (
	"github.com/ebar-go/ego/http"
)

// HttpServer 获取httpServer示例
func HttpServer() *http.Server {
	return http.New()
}
