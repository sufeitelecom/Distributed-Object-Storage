package objects

import (
	"net/http"
	"strings"
	"io"
	"github.com/sufeitelecom/distributed-object-storage/objectstream"
	"github.com/sufeitelecom/distributed-object-storage/heartbeat"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/sufeitelecom/distributed-object-storage/locate"
)

func ApiHandler(w http.ResponseWriter,r *http.Request)  {
	method := r.Method //获取http请求的动作，根据动作进行相应处理
	if method == http.MethodPut {
		apiput(w,r)
		return
	} else if method == http.MethodGet{
		apiget(w,r)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func apiput(w http.ResponseWriter,r *http.Request)  {
	object := strings.Split(r.URL.EscapedPath(),"/")[2]
	c,err := storeobject(r.Body,object)
	if err != nil{
		log.Errorf("store object error %v",err)
	}
	w.WriteHeader(c)
}

func apiget(w http.ResponseWriter,r *http.Request)  {
	object := strings.Split(r.URL.EscapedPath(),"/")[2]
	stream,err := getstream(object)
	if err != nil{
		log.Errorf("Getstream error %v",err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	io.Copy(w,stream)
}

func storeobject(r io.Reader,object string) (int,error) {
	stream,err := putstream(object)
	if err != nil{
		return http.StatusServiceUnavailable,err
	}

	io.Copy(stream,r)
	err = stream.Close()
	if err != nil{
		return http.StatusInternalServerError,err
	}
	return http.StatusOK,nil
}

func putstream(object string) (*objectstream.PutStream,error)  {
	server := heartbeat.ChooseRandomServer()
	if server == ""{
		return nil,fmt.Errorf("cannot find any servers")
	}
	return objectstream.NewPutStream(server,object),nil
}

func getstream(object string)  (*objectstream.GetStream,error) {
	server := locate.Location(object)
	if server == ""{
		return nil,fmt.Errorf("Cannot find object %s.",object)
	}
	return objectstream.NewGetStream(server,object)
}