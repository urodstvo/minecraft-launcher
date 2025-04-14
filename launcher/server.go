package launcher

import (
	"embed"
	"io/fs"
	"log"
	"net"
	"net/http"
)

//go:embed ui/*
var embeddedFiles embed.FS

type Server struct{}

func (s *Server) Run() string {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}

	// Создаём файловую систему на основе встроенных файлов
	subFS, _ := fs.Sub(embeddedFiles, "ui")
	fileServer := http.FileServer(http.FS(subFS))

	go func() {
		defer ln.Close()
		http.Handle("/", fileServer)
		log.Fatal(http.Serve(ln, nil))
	}()

	return "http://" + ln.Addr().String()
}
