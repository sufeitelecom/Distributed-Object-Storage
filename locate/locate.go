package locate

import (
	"github.com/sufeitelecom/distributed-object-storage/rabbitmq"
	"strconv"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"sync"
)

var objects  = make(map[string]int)
var mutex sync.Mutex

func StartLocate()  {
	mq := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	defer mq.Close()

	mq.Bind("dataservers")

	c := mq.Consume()
	for msg := range c {
		hash ,err := strconv.Unquote(string(msg.Body))
		if err != nil{
			log.Fatalf("Msg.body error %v",err)
		}

		if Locate(hash) {
			mq.Send(msg.ReplyTo,os.Getenv("LISTEN_ADDRESS"))
		}
	}
}

func CollectObject()  {
	files,_:= filepath.Glob(os.Getenv("STORAGE_ROOT")+"/objects/*")
	for i := range files{
		hash := filepath.Base(files[i])
		objects[hash] = 1
	}
}

func Locate(hash string) bool  {
	mutex.Lock()
	_,ok := objects[hash]
	mutex.Unlock()
	return ok
}

func Add(hash string)  {
	mutex.Lock()
	objects[hash] = 1
	mutex.Unlock()
}

func Del(hash string)  {
	mutex.Lock()
	delete(objects,hash)
	mutex.Unlock()
}