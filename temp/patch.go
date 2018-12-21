package temp

import (
	"net/http"
	"strings"
	"os"
	"io/ioutil"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io"
)

/*
Request PATCH /temp/<uuid>  Context object
 */
func patch(w http.ResponseWriter,r *http.Request)  {
	uuid := strings.Split(r.URL.EscapedPath(),"/")[2]
	info,e := readFile(uuid)
	if e !=nil{
		log.Errorf("Temp Patch read file error %v",e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	infofile := os.Getenv("STORAGE_ROOT")+"/temp/"+uuid
	datafile := infofile + ".dat"
	f,e := os.OpenFile(datafile,os.O_WRONLY|os.O_APPEND,0)
	if e != nil{
		log.Errorf("Temp Patch open file error %v",e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer f.Close()

	_,e = io.Copy(f,r.Body)
	if e != nil{
		log.Errorf("Temp Patch copy error %v",e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	stat,e := f.Stat()
	if e != nil{
		log.Errorf("Temp Patch file stat error %v",e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	actual := stat.Size()
	if actual > info.Size{
		os.Remove(infofile)
		os.Remove(datafile)
		log.Errorf("actual size %d is not equal tempinfo.size %d",actual,info.Size)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func readFile(uuid string) (*tempinfo,error){
	f,e := os.Open(os.Getenv("STORAGE_ROOT")+"/temp/"+uuid)
	if e != nil{
		return nil,e
	}
	defer f.Close()

	b,_ := ioutil.ReadAll(f)
	var info tempinfo
	json.Unmarshal(b,&info)
	return &info,nil
}