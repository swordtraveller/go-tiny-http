package http

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type File interface {
	io.Closer
	io.Reader
	io.Seeker
	Readdir(count int) ([]fs.FileInfo, error)
	Stat() (fs.FileInfo, error)
}

type FileSystem interface {
	Open(name string) (File, error)
}

type fileHandler struct {
	root FileSystem
}

type Dir string

func (f fileHandler) ServeHTTP(resp ResponseWriter, req *Request) {
	// If Path does not have a prefix "/", add it
	// 如果没有前缀"/"，则补充
	if !strings.HasPrefix(req.URL.Path, "/") {
		req.URL.Path = "/" + req.URL.Path
	}

	// Wrapped new Open() to add root in front of Path
	// 包装了新的Open方法用于在Path的前面加上root
	file, err := f.root.Open(req.URL.Path)
	if err != nil {
		Error(resp, "Not Found", StatusNotFound)
		return
	}
	defer func(file File) {
		err := file.Close()
		if err != nil {
			return
		}
	}(file)

	// Set the Content-Type response header according to the file type
	// 根据文件类型设置Content-Type响应头
	resp.Header().Set("Content-Type", "text/plain; charset=utf-8")
	if strings.HasSuffix(req.URL.Path, "css") {
		resp.Header().Set("Content-Type", "text/css; charset=utf-8")
	}

	// send file
	_, err = io.Copy(resp, file)
	if err != nil {
		Error(resp, "File Error", StatusInternalServerError)
	}
}

// Open returns the file "d/name"
func (d Dir) Open(name string) (File, error) {
	dir := string(d)
	fullPath := filepath.Join(dir, name)

	// Check if the file exists
	// 检查文件是否存在
	_, err := os.Stat(fullPath)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(fullPath)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func FileServer(root FileSystem) Handler {
	return &fileHandler{root}
}
