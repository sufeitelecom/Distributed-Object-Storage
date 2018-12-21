package objects

import (
	"net/http"
	"os"
	"strings"
	log "github.com/sirupsen/logrus"
	"io"
	"net/url"
	"github.com/sufeitelecom/distributed-object-storage/tools"
	"github.com/sufeitelecom/distributed-object-storage/locate"
)



func DataHandler(w http.ResponseWriter,r *http.Request)  {
	method := r.Method //获取http请求的动作，根据动作进行相应处理
	if method == http.MethodGet{
		get(w,r)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}



func get(w http.ResponseWriter,r *http.Request)  {
	hash := strings.Split(r.URL.EscapedPath(),"/")[2]
	file := checkfie(hash)
	if file == ""{
		w.WriteHeader(http.StatusNotFound)
		return
	}
	f,_ := os.Open(file)
	defer f.Close()
	io.Copy(w,f)
}

func checkfie(hash string) string  {
	filename := os.Getenv("STORAGE_ROOT")+"/objects/"+hash
	file,err := os.Open(filename)
	if err != nil{
		log.Errorf("Can't open the file %s, %v",hash,err)
		return ""
	}
	d := url.PathEscape(tools.CalculateHash(file))
	file.Close()
	if d != hash{
		log.Errorf("object hash mismatch,remove")
		locate.Del(hash)
		os.Remove(filename)
		return ""
	}
	return filename
}