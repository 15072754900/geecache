package geecache

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

type HTTPPool struct {
	self     string
	bashPath string
}

var defaultBasePath string = "/_geecache/"

func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self,
		bashPath: defaultBasePath,
	}
}

// 这里使用http://example.com/_geecache/; 开头的请求作为节点间通信的api

// 建立ServeHTTP服务，和日志服务

// v这里是一个slice，包含任意类型的元素

func (p *HTTPPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, v...))
}

func (p *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 判断是够有前缀
	if !strings.HasPrefix(r.URL.Path, p.bashPath) {
		panic("HTTPPool serving unexpected path: " + r.URL.Path)
	}
	p.Log("%s%s", r.Method, r.URL.Path)
	// 处理输入
	parts := strings.SplitN(r.URL.Path[len(p.bashPath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	// 进入缓存读取数据

	groupName := parts[0]
	key := parts[1]
	fmt.Println(parts, "--", groupName, "--", key, "--", len(parts))

	group := GetGroup(groupName)
	if group == nil {
		http.Error(w, "no such group: "+groupName, http.StatusNotFound)
		return
	}

	// 存在则进行获取实际数据
	view, err := group.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(view.ByteSlice())
}
