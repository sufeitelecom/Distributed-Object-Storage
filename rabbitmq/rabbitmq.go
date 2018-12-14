package rabbitmq

import (
	"github.com/streadway/amqp"
	log "github.com/sirupsen/logrus"
	"encoding/json"
)


type RabbitMQ struct {
	channel *amqp.Channel   //信道，一个连接可以对应多个信道
	Name string             //队列名字
	exchange string         //交换机名字
}

func New(s string) *RabbitMQ  {
	conn,err := amqp.Dial(s)
	if err != nil{
		log.Fatalf("Connect RabbitMQ %s error %v",s,err)
	}

	ch,err :=conn.Channel()
	if err != nil{
		log.Fatalf("Open channel error %v",err)
	}

	q,err := ch.QueueDeclare("",false,true,false,false,nil)
	if err != nil{
		log.Fatalf("Create queue error %v.",err)
	}

	mq := new(RabbitMQ)
	mq.Name = q.Name
	mq.channel = ch
	return mq
}

func (mq *RabbitMQ)Bind(exchange string)  {
	err := mq.channel.QueueBind(mq.Name,"",exchange,false,nil)
	if err != nil{
		log.Fatalf("Bind queue error %v",err)
	}
	mq.exchange = exchange
}

func (mq *RabbitMQ)Send(queue string,body interface{})  {
	str,err := json.Marshal(body)
	if err != nil{
		log.Fatalf("Json marshal error %v",err)
	}

	err = mq.channel.Publish("",queue,false,false,amqp.Publishing{
		ReplyTo:mq.Name,
		Body:[]byte(str),
	})
	if err != nil{
		log.Fatalf("Publish error %v",err)
	}
}

func (mq *RabbitMQ)Publish(exchange string,body interface{})  {
	str,err := json.Marshal(body)
	if err != nil{
		log.Fatalf("Json marshal error %v",err)
	}
	err = mq.channel.Publish(exchange,"",false,false,amqp.Publishing{
		ReplyTo:mq.Name,
		Body:[]byte(str),
	})
	if err != nil{
		log.Fatalf("Publish error %v",err)
	}
}

func (mq *RabbitMQ)Consume() <-chan amqp.Delivery  {
	c,err := mq.channel.Consume(mq.Name,"",true,false,false,false,nil)
	if err != nil{
		log.Fatalf("Consume error %v",err)
	}
	return c
}

func (mq *RabbitMQ)Close()  {
	mq.channel.Close()
}