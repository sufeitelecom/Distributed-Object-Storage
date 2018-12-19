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
	"github.com/sufeitelecom/distributed-object-storage/es"
	"github.com/sufeitelecom/distributed-object-storage/tools"
	"net/url"
	"strconv"
)

func ApiHandler(w http.ResponseWriter,r *http.Request)  {
	method := r.Method //获取http请求的动作，根据动作进行相应处理
	if method == http.MethodPut {
		apiput(w,r)
		return
	} else if method == http.MethodGet{
		apiget(w,r)
		return
	} else if method == http.MethodDelete{
		apidelete(w,r)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}
func apidelete(w http.ResponseWriter,r *http.Request)  {
	name := strings.Split(r.URL.EscapedPath(),"/")[2]
	version,err := es.SearchLatestVersion(name)
	if err != nil{
		log.Errorf("es.SearchLatestVersion error %v",err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = es.PutMetadata(name,version.Version+1,0,"")
	if err != nil{
		log.Errorf("es.PutMetadata error %v",err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func apiput(w http.ResponseWriter,r *http.Request)  {
	hash := tools.GetHashFromHeader(r.Header)
	if hash == ""{
		log.Errorf("missing object hash in degist header")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	c,err := storeobject(r.Body,url.PathEscape(hash))
	if err != nil{
		log.Errorf("store object error %v",err)
		w.WriteHeader(c)
		return
	}
	if c != http.StatusOK{
		w.WriteHeader(c)
		return
	}
	name := strings.Split(r.URL.EscapedPath(),"/")[2]
	size := tools.GetSizeFromHeader(r.Header)
	err = es.AddVersion(name,hash,size)
	if err != nil {
		log.Errorf("es.AddVersion error %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func apiget(w http.ResponseWriter,r *http.Request)  {
	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	versionId := r.URL.Query()["version"]
	version := 0
	var e error
	if len(versionId) != 0 {
		version, e = strconv.Atoi(versionId[0])
		if e != nil {
			log.Println(e)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	meta, e := es.GetMatadata(name, version)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if meta.Hash == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	object := url.PathEscape(meta.Hash)
	stream, e := getstream(object)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	io.Copy(w, stream)
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

