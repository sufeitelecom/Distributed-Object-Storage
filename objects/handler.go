package objects

import (
	"net/http"
	"os"
	"strings"
	log "github.com/sirupsen/logrus"
	"io"
)

const StoragePath = "/home/sufei/object_storage/objects/"

func Handler(w http.ResponseWriter,r *http.Request)  {
	method := r.Method
	if method == http.MethodPut {
		put(w,r)
		return
	} else if method == http.MethodGet{
		get(w,r)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func put(w http.ResponseWriter,r *http.Request)  {
	file,err := os.Create(StoragePath+strings.Split(r.URL.EscapedPath(),"/")[2])
	if err != nil{
		log.Errorf("Can't create the file %s, %v",r.URL.EscapedPath(),err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()

	io.Copy(file,r.Body)
}

func get(w http.ResponseWriter,r *http.Request)  {
	file,err := os.Open(StoragePath+strings.Split(r.URL.EscapedPath(),"/")[2])
	if err != nil{
		log.Errorf("Can't open the file %s, %v",r.URL.EscapedPath(),err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	defer file.Close()

	io.Copy(w,file)
}