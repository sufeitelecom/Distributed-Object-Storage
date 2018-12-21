package objectstream

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"strings"
)

type PutStream struct {
	Server string
	Uuid string
}

func NewPutStream(server string,hash string,size int64) (*PutStream , error) {
	request,e := http.NewRequest("POST","http://"+server+"/temp/"+hash,nil)
	if e != nil{
		return nil,e
	}
	request.Header.Set("size",fmt.Sprintf("%d",size))
	client := http.Client{}
	response,e := client.Do(request)
	if e != nil{
		return nil,e
	}
	uuid,e := ioutil.ReadAll(response.Body)
	if e != nil{
		return nil,e
	}
	return &PutStream{Server:server,Uuid:string(uuid)},nil
}

func (w *PutStream)Write(p []byte) (n int,err error)  {
	request,e := http.NewRequest("PATCH","http://"+w.Server+"/temp/"+w.Uuid,strings.NewReader(string(p)))
	if e != nil{
		return 0,e
	}
	client := http.Client{}
	r,e := client.Do(request)
	if e != nil{
		return  0,e
	}
	if r.StatusCode != http.StatusOK{
		return 0,fmt.Errorf("dataserver return http code %d",r.StatusCode)
	}
	return len(p),nil
}

func (w *PutStream)Commit(good bool)  {
	methed := "DELETE"
	if good {
		methed = "PUT"
	}
	request,_:= http.NewRequest(methed,"http://"+w.Server+"/temp/"+w.Uuid,nil)
	client := http.Client{}
	client.Do(request)
}

