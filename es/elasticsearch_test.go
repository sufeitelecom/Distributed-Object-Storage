package es

import (
	"testing"
	"strings"
	"time"
)


func TestPutAndGetMatadata(t *testing.T) {
	ori := Metadata{
		Name:"sufeitest",
		Version:1,
		Hash:"qwertyuiop",
		Size:1,
	}

	err := PutMetadata(ori.Name,ori.Version,ori.Size,ori.Hash)
	if err != nil{
		t.Errorf("put meta fail %v",err)
		return
	}

	time.Sleep(2*time.Second)
	now,err := GetMetadata(ori.Name,0)
	if err != nil{
		t.Errorf("get meta fail %v",err)
		return
	}
	time.Sleep(2*time.Second)
	t.Logf("[添加后]原始元数据为 : %v,读取的元数据为：%v.",ori,now)
	if strings.Compare(now.Name,ori.Name) != 0 {
		t.Errorf("test fail")
		return
	}

	DelMetadata(now.Name,now.Version)
	time.Sleep(2*time.Second)
	now2,err := GetMetadata(ori.Name,0)
	t.Logf("[删除后]原始元数据为 : %v,读取的元数据为：%v.",ori,now2)

	return
}