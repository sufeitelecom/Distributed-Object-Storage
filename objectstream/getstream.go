package objectstream

import (
	"io"
	"fmt"
	"net/http"
)

type GetStream struct {
	reader io.Reader
}

func NewGetStream(server,object string) (*GetStream,error) {
	if server == "" || object == ""{
		return nil,fmt.Errorf("invalid server %s object %s",server,object)
	}
	return newGetStream("http://"+server+"/objects/"+object)
}

func newGetStream(url string)  (*GetStream,error){
	r,err := http.Get(url)
	if err != nil{
		return nil,err
	}

	if r.StatusCode != http.StatusOK{
		return nil,fmt.Errorf("dataserver return http code %d",r.StatusCode)
	}

	return &GetStream{r.Body},nil
}

func (r *GetStream)Read(p []byte) (b int,err error)  {
	return r.reader.Read(p)
}