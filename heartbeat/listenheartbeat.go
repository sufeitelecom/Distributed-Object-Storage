package heartbeat

import (
	"github.com/sufeitelecom/distributed-object-storage/rabbitmq"
	"os"
	"time"
	"sync"
	"strconv"
	"log"
	"math/rand"
)

var dataServers  = make(map[string]time.Time)
var mutex sync.Mutex

func ListenHeartbeat()  {
	mq := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	defer mq.Close()

	mq.Bind("apiservers")
	c := mq.Consume()

	go removeExpireDataserver()
	for msg := range c{
		dataserver,err := strconv.Unquote(string(msg.Body))
		if err != nil{
			log.Fatalf("Msg.body error %v",err)
		}

		mutex.Lock()
		dataServers[dataserver] = time.Now()
		mutex.Unlock()
	}
}

func removeExpireDataserver()  {
	for{
		time.Sleep(5*time.Second)
		mutex.Lock()
		for s,t := range dataServers{
			if t.Add(10 *time.Second).Before(time.Now()){
				delete(dataServers,s)
			}
		}
		mutex.Unlock()
	}
}
func GetDataServers() []string {
	mutex.Lock()
	defer mutex.Unlock()

	ds := make([]string,0)
	for s,_ := range dataServers{
		ds =append(ds,s)
	}
	return ds
}

func ChooseRandomServer() string {
	ds := GetDataServers()
	n := len(ds)
	if n == 0{
		return ""
	}
	return ds[rand.Intn(n)]
}