package main


import (

    "github.com/sufeitelecom/distributed-object-storage/objects"
	"net/http"
	"os"
	log "github.com/sirupsen/logrus"
)

const 	Listen = ":12345"

func initLog()  {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
}

func main() {
	initLog()
	http.HandleFunc("/objects/",objects.Handler)
	log.Fatal(http.ListenAndServe(Listen,nil))
}