package objectstream

import (
	"io"
	"net/http"
	"fmt"
)

type PutStream struct {
	write *io.PipeWriter
	c chan error
}

func NewPutStream(server string,object string) *PutStream {
	readio,writeio := io.Pipe()
	c := make(chan error)

	go func() {
		request,_ := http.NewRequest("PUT","http://"+server+"/objects/"+object,readio)
		client := http.Client{}
		r,err := client.Do(request)
		if err == nil && r.StatusCode != http.StatusOK{
			err = fmt.Errorf("dataserver return http code %d",err)
		}
		c <- err
	}()
	return &PutStream{writeio,c}
}

func (w *PutStream)Write(p []byte) (n int,err error)  {
	return w.write.Write(p)
}

//调用关闭，是为了让管道另一端的reader能收到io.EOF,否则在goroutine中运行的client.Do(request)始终阻塞无法返回
func (w * PutStream)Close() error {
	w.write.Close()
	return <-w.c
}
