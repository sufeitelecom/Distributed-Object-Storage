package heartbeat

import (
	"github.com/sufeitelecom/distributed-object-storage/rabbitmq"
	"os"
	"time"
)

/*
rabbitmq服务首先创建了两个交换机apiservers和dataservers
 */
func Startheartbeat()  {
	mq := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	defer mq.Close()

	for {
		mq.Publish("apiservers",os.Getenv("LISTEN_ADDRESS"))
		time.Sleep(5 * time.Second)
	}
}

