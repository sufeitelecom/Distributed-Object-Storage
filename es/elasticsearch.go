package es

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"io/ioutil"
	"encoding/json"
	"net/url"
)

type Metadata struct {
	Name string
	Version int
	Size int64
	Hash string
}

type hit struct {
	Source Metadata `json:"_source"`
}

type searchResult struct {
	Hits struct {
		Total int
		Hits  []hit
	}
}

func PutMetadata(name string,version int,size int64,hash string) error {
	ctx := fmt.Sprintf(`{"name":"%s","version":%d,"size":%d,"hash":"%s"}`,name,version,size,hash)
	client := http.Client{}

	url := fmt.Sprintf("http://%s/metadata/objects/%s_%d?op_type=create",os.Getenv("ES_SERVER"),name,version)
	request,_ := http.NewRequest("PUT",url,strings.NewReader(ctx))
	request.Header.Set("Content-Type", "application/json")
	result,err := client.Do(request)
	if err != nil{
		return err
	}
	if result.StatusCode == http.StatusConflict {
		return PutMetadata(name,version+1,size,hash)
	}
	if result.StatusCode != http.StatusCreated {
		r,_ := ioutil.ReadAll(result.Body)
		return fmt.Errorf("fail to put metadata:%d %s",result.StatusCode,string(r))
	}
	return nil
}

func getMetadata(name string, version int) (meta Metadata,e error)  {
	url :=  fmt.Sprintf("http://%s/metadata/objects/%s_%d/_source",os.Getenv("ES_SERVER"),name,version)
	r,e := http.Get(url)
	if e != nil{
		return
	}
	if r.StatusCode != http.StatusOK{
		e = fmt.Errorf("fail to get %s_%d: %d",name,version,r.StatusCode)
		return
	}
	result,_ := ioutil.ReadAll(r.Body)
	json.Unmarshal(result,&meta)
	return
}

func SearchLatestVersion(name string) (meta Metadata,err error)  {
	url := fmt.Sprintf("http://%s/metadata/_search?q=name:%s&size=1&sort=version:desc",
		os.Getenv("ES_SERVER"), url.PathEscape(name))
	r , e:= http.Get(url)
	if e != nil{
		return
	}
	if r.StatusCode != http.StatusOK{
		e = fmt.Errorf("fail to search latest metadata:%d",r.StatusCode)
		return
	}
	result,_ := ioutil.ReadAll(r.Body)
	var sr searchResult
	json.Unmarshal(result,&sr)
	if len(sr.Hits.Hits) != 0{
		meta = sr.Hits.Hits[0].Source
	}
	return
}

func GetMetadata(name string,version int) (Metadata,error)  {
	if version == 0{
		return SearchLatestVersion(name)
	}
	return getMetadata(name,version)
}

func DelMetadata(name string, version int) {
	client := http.Client{}
	url := fmt.Sprintf("http://%s/metadata/objects/%s_%d",
		os.Getenv("ES_SERVER"), name, version)
	request, _ := http.NewRequest("DELETE", url, nil)
	client.Do(request)
}

func AddVersion(name, hash string, size int64) error {
	version, e := SearchLatestVersion(name)
	if e != nil {
		return e
	}
	return PutMetadata(name, version.Version+1, size, hash)
}

func SearchAllVersions(name string, from, size int) ([]Metadata, error) {
	url := fmt.Sprintf("http://%s/metadata/_search?sort=name,version&from=%d&size=%d",
		os.Getenv("ES_SERVER"), from, size)
	if name != "" {
		url += "&q=name:" + name
	}
	r, e := http.Get(url)
	if e != nil {
		return nil, e
	}
	metas := make([]Metadata, 0)
	result, _ := ioutil.ReadAll(r.Body)
	var sr searchResult
	json.Unmarshal(result, &sr)
	for i := range sr.Hits.Hits {
		metas = append(metas, sr.Hits.Hits[i].Source)
	}
	return metas, nil
}
