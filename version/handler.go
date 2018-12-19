package version

import (
	"net/http"
	"strings"
	"github.com/sufeitelecom/distributed-object-storage/es"
	log "github.com/sirupsen/logrus"
	"encoding/json"
)

func Handler(w http.ResponseWriter,r *http.Request)  {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	from := 0
	size := 100
	name := strings.Split(r.URL.EscapedPath(),"/")[2]
	for {
		metas,err := es.SearchAllVersions(name,from,size)
		if err != nil{
			log.Errorf("es.SearchAllVersions error :%v",err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		for i := range metas{
			b,_ := json.Marshal(metas[i])
			w.Write(b)
			w.Write([]byte("\n"))
		}
		if len(metas) != size{
			return
		}
		from += size
	}
}
