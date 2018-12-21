package temp

import (
	"net/http"
	"strings"
	"os"
)

/*
Request DELETE /temp/<uuid> 数据不一致删除
 */
func delete(w http.ResponseWriter,r *http.Request) {
	uuid := strings.Split(r.URL.EscapedPath(),"/")[2]
	infofile := os.Getenv("STORAGE_ROOT")+"/temp/"+uuid
	datafile := infofile + ".dat"
	os.Remove(infofile)
	os.Remove(datafile)
}