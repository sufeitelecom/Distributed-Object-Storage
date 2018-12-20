package main

import (
	"github.com/sufeitelecom/distributed-object-storage/heartbeat"
	"os"
	log "github.com/sirupsen/logrus"
	"net/http"
	"github.com/sufeitelecom/distributed-object-storage/objects"
	"github.com/sufeitelecom/distributed-object-storage/locate"
	"github.com/sufeitelecom/distributed-object-storage/version"
)

func main()  {
	//初始化日志系统
	initLog()
	checkEnv()


	//开启一个线程完成上报心跳以及服务发现
	go heartbeat.ListenHeartbeat()


	/*
	注册http处理函数，如果有客户端访问本机的http服务，
	url以/objects/开头，那么将由objects.Handler函数处理
	url以/locate/开头，那么将由locate.handler函数处理
	*/
	http.HandleFunc("/objects/",objects.ApiHandler)
	http.HandleFunc("/locate/",locate.Handler)
	http.HandleFunc("/versions/",version.Handler)
	/*
	http.ListenAndServe开始监听,监听端口由环境变量指定，这样方便同一台机器启动多个数据服务
	*/
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"),nil))
}

func initLog()  {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
}

func checkEnv()  {
	str := os.Getenv("LISTEN_ADDRESS")
	if str == ""{
		log.Fatalf("please set environment variables LISTEN_ADDRESS")
	}
	str = os.Getenv("RABBITMQ_SERVER")
	if str == ""{
		log.Fatalf("please set environment variables RABBITMQ_SERVER")
	}
}
