package distributed_object_storage

import (
	log "github.com/Sirupsen/logrus"
	"os"
	"net/http"
)

const 	Listen = ":12345"

func initLog()  {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
}

func main() {
	initLog()
	http.HandleFunc("/objects/",objects.Handler())
	log.Fatal(http.ListenAndServe(Listen,nil))
}