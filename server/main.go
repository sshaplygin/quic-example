package main

import (
	"github.com/gin-gonic/gin"
	"github.com/lucas-clemente/quic-go/http3"
	"log"
	"net/http"
	"os"
	"path"
)

func main() {
	currentPath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	//logWriter, err := utils.GetQLOGWriter()
	//if err != nil {
	//	log.Fatal(err)
	//}

	//tracer := qlog.NewTracer(logWriter)
	//conf := &quic_go.Config{
	//	Tracer: tracer,
	//}

	router := gin.New()

	router.GET("/", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{"message": "OK"})
	})

	server := http3.Server{
		//QuicConfig: conf,
		Addr:    "127.0.0.1:8080",
		Handler: router,
	}

	err = server.ListenAndServeTLS(path.Join(currentPath, "cert.pem"), path.Join(currentPath, "private.key"))
	if err != nil {
		log.Printf("Server error: %v", err)
	}
}
