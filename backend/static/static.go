package static

import (
	"embed"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"io"
	"io/fs"
	"mime"
	"net"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

//go:embed dist/*
var staticFiles embed.FS

var (
	globalData = make(map[string]interface{})
	mu         sync.RWMutex
)

func Set(key string, value interface{}) {
	mu.Lock()
	defer mu.Unlock()
	globalData[key] = value
}

func Get(key string) interface{} {
	mu.RLock()
	defer mu.RUnlock()
	return globalData[key]
}

func Delete(key string) {
	mu.RLock()
	defer mu.RUnlock()
	delete(globalData, key)
}

func StartFileServer() string {
	r := chi.NewRouter()

	files, _ := fs.Sub(staticFiles, "dist")

	r.HandleFunc("/*", func(w http.ResponseWriter, req *http.Request) {
		serveStaticOrIndex(w, req, files)
	})

	l, _ := net.Listen("tcp", "127.0.0.1:0")
	go func() {
		if err := http.Serve(l, r); err != nil {
			// 可以加日志
		}
	}()

	return l.Addr().String()
}

func serveStaticOrIndex(w http.ResponseWriter, req *http.Request, files fs.FS) {
	reqPath := strings.TrimPrefix(req.URL.Path, "/")

	if reqPath == "pxStore" {
		port := req.URL.Query().Get("port")
		secret := req.URL.Query().Get("secret")
		Set("port", port)
		Set("secret", secret)
		render.PlainText(w, req, "ok")
		return
	}

	f, err := files.Open(reqPath)
	if err != nil {
		// 打不开，返回 index.html
		serveIndexHtml(w, files)
		return
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil || stat.IsDir() {
		// 读取失败或是目录，也返回 index.html
		serveIndexHtml(w, files)
		return
	}

	// 是静态文件，正常返回
	setContentType(w, reqPath)
	w.Header().Set("Content-Length", strconv.FormatInt(stat.Size(), 10))
	io.Copy(w, f)
}

func serveIndexHtml(w http.ResponseWriter, files fs.FS) {
	f, err := files.Open("index.html")
	if err != nil {
		http.Error(w, "index.html not found", http.StatusInternalServerError)
		return
	}
	defer f.Close()

	stat, _ := f.Stat()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Content-Length", strconv.FormatInt(stat.Size(), 10))
	io.Copy(w, f)
}

func setContentType(w http.ResponseWriter, name string) {
	ext := filepath.Ext(name)
	if ext != "" {
		if ctype := mime.TypeByExtension(ext); ctype != "" {
			w.Header().Set("Content-Type", ctype)
			return
		}
	}
	w.Header().Set("Content-Type", "application/octet-stream")
}
