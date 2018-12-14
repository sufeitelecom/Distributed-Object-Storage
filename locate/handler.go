package locate

import (
	"net/http"
	"github.com/sufeitelecom/distributed-object-storage/rabbitmq"
	"time"
	"strconv"
	"strings"
	"encoding/json"
	"os"
)

func Handler(w http.ResponseWriter,r *http.Request)  {
	m := r.Method
	if m != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	info := Location(strings.Split(r.URL.EscapedPath(),"/")[2])
	if len(info) == 0{
		w.WriteHeader(http.StatusNotFound)
		return
	}
	b,_ := json.Marshal(info)
	w.Write(b)
}

func Location(name string) string  {
	mq := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	mq.Publish("dataservers",name)
	c := mq.Consume()

	go func() {
		time.Sleep(time.Second)
		mq.Close()
	}()

	msg := <-c
	s,_ := strconv.Unquote(string(msg.Body))
	return s
}

func Exist(name string) bool {
	return Location(name) != ""
}