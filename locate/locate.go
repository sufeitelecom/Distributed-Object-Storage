package locate

import (
	"github.com/sufeitelecom/distributed-object-storage/rabbitmq"
	"strconv"
	log "github.com/sirupsen/logrus"
	"os"
)

func StartLocate()  {
	mq := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	defer mq.Close()

	mq.Bind("dataservers")

	c := mq.Consume()
	for msg := range c {
		object ,err := strconv.Unquote(string(msg.Body))
		if err != nil{
			log.Fatalf("Msg.body error %v",err)
		}

		if Locate(os.Getenv("STORAGE_ROOT")+"/objects/"+object) {
			mq.Send(msg.ReplyTo,os.Getenv("LISTEN_ADDRESS"))
		}
	}
}

func Locate(name string) bool  {
	_,err := os.Stat(name)
	return !os.IsNotExist(err)
}