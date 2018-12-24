package temp

import (
	"net/http"
	"os/exec"
	"strings"
	"strconv"
	"os"
	"encoding/json"
	log "github.com/sirupsen/logrus"
)

type tempinfo struct {
	Uuid string
	Name string
	Size int64
}

/*
Request  POST /temp/<hash>  Header Size
Response uuid
 */
func post(w http.ResponseWriter,r *http.Request) {
	output, _ := exec.Command("dbus-uuidgen").Output()
	uuid := strings.TrimSuffix(string(output), "\n")  //注意生成的uuid包含\n后缀，而在url中该字符别翻译为%OA,造成无法删除临时问题
	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	size, e := strconv.ParseInt(r.Header.Get("size"), 0, 64)
	if e != nil{
		log.Errorf("Temp/<hash> post parse_size error %v",e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	t := tempinfo{Uuid:uuid,Name:name,Size:size}
	e = t.writeToFile()
	if e!= nil{
		log.Errorf("Temp/<hash> post write to file error %v",e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	os.Create(os.Getenv("STORAGE_ROOT")+"/temp/"+t.Uuid+".dat")
	w.Write([]byte(t.Uuid))
}

func (t *tempinfo)writeToFile() error  {
	f,e := os.Create(os.Getenv("STORAGE_ROOT")+"/temp/"+t.Uuid)
	if e != nil{
		return e
	}
	defer f.Close()
	b,_ := json.Marshal(t)
	f.Write(b)
	return nil
}