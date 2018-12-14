package objects

import (
	"net/http"
	"os"
	"strings"
	log "github.com/sirupsen/logrus"
	"io"
)



func DataHandler(w http.ResponseWriter,r *http.Request)  {
	method := r.Method //获取http请求的动作，根据动作进行相应处理
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
	//创建相关文件
	file,err := os.Create(os.Getenv("STORAGE_ROOT")+"/objects/"+strings.Split(r.URL.EscapedPath(),"/")[2])
	if err != nil{
		log.Errorf("Can't create the file %s, %v",r.URL.EscapedPath(),err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()
	//写入内容
	io.Copy(file,r.Body)
}

func get(w http.ResponseWriter,r *http.Request)  {
	file,err := os.Open(os.Getenv("STORAGE_ROOT")+"/objects/"+strings.Split(r.URL.EscapedPath(),"/")[2])
	if err != nil{
		log.Errorf("Can't open the file %s, %v",r.URL.EscapedPath(),err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	defer file.Close()

	io.Copy(w,file)
}