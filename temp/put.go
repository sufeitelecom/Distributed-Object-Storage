package temp

import (
	"net/http"
	"strings"
	"os"
	log "github.com/sirupsen/logrus"
	"github.com/sufeitelecom/distributed-object-storage/locate"
)

/*
Request PUT /temp/<uuid> 临时数据转正
 */
func put(w http.ResponseWriter,r *http.Request)  {
	uuid := strings.Split(r.URL.EscapedPath(),"/")[2]
	info,e := readFile(uuid)
	if e !=nil{
		log.Errorf("Temp put read file error %v",e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	infofile := os.Getenv("STORAGE_ROOT")+"/temp/"+uuid
	datafile := infofile + ".dat"
	f,e := os.Open(datafile)
	if e != nil{
		log.Errorf("Temp PUT open file error %v",e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer f.Close()

	stat,e := f.Stat()
	if e != nil{
		log.Errorf("Temp PUT file stat error %v",e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	os.Remove(infofile)
	if info.Size != stat.Size(){
		os.Remove(datafile)
		log.Println("actual size %d is not equal tempinfo.size %d",stat.Size(),info.Size)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	os.Rename(datafile, os.Getenv("STORAGE_ROOT")+"/objects/"+info.Name)
	locate.Add(info.Name)

}