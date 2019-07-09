package main

import (
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lucas-clemente/quic-go"
	"github.com/lucas-clemente/quic-go/http3"
	"github.com/lucas-clemente/quic-go/quictrace"
)

func main() {

	router := gin.New()

	router.GET("/", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{"message": "OK"})
	})
	// _, filename, _, ok := runtime.Caller(0)
	// if !ok {
	// 	panic("Failed to get current frame")
	// }
	// currentPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	currentPath, err := os.Getwd()
	// currentPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(currentPath)
	tracer := quictrace.NewTracer()
	quicConf := &quic.Config{QuicTracer: tracer}

	server := http3.Server{
		Server: &http.Server{
			Addr:           ":8080",
			Handler:        router,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		},
		QuicConfig: quicConf,
	}

	err = server.ListenAndServeTLS(path.Join(currentPath, "cert.pem"), path.Join(currentPath, "priv.key"))

	if err != nil {
		log.Printf("Server error: %v", err)
	}
}
