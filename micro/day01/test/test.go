package main

import (
	"context"
	"fmt"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/client/selector"
	"github.com/micro/go-micro/registry"
	myhttp "github.com/micro/go-plugins/client/http"
	"github.com/micro/go-plugins/registry/consul"

	"io/ioutil"
	"log"
	"net/http"
)

//callAPI 基本调用API方式
func callAPI(addr string,path string,method string)(string,error){
	req,_ := http.NewRequest(method,"http://"+addr+path,nil)
	client := http.DefaultClient
	res,err := client.Do(req)
	if err != nil{
		return "",err
	}
	defer res.Body.Close()
	buf,_ := ioutil.ReadAll(res.Body)
	return string(buf),nil
}

//callAPI2 使用插件的方式调用API
func callAPI2(s selector.Selector){
	myClient := myhttp.NewClient(
		client.Selector(s),
		client.ContentType("application/json"),
		)
	req := myClient.NewRequest("prodservice","/v1/prods",map[string]string{})
	var resp map[string]interface{}
	err := myClient.Call(context.Background(),req,&resp)
	if err != nil{
		log.Fatal(err)
	}
	fmt.Println(resp["data"])
}

func main() {
	consulReg := consul.NewRegistry(
		registry.Addrs("192.168.86.113:8500"),
	)

	// 基本调用方式
	//{
	//	getService,err := consulReg.GetService("prodservice")
	//	if err != nil{
	//		log.Fatal(err)
	//	}
	//	next := selector.Random(getService)
	//
	//	node,err := next()
	//	if err != nil{
	//		log.Fatal(err)
	//	}
	//
	//	fmt.Println(node.Id,node.Address,node.Metadata)
	//	callRes,err := callAPI(node.Address,"/v1/prods","POST")
	//	if err != nil{
	//		log.Fatal(err)
	//	}
	//	fmt.Println(callRes)
	//}

	//通过go-micro插件调用
	{
		mySelector := selector.NewSelector(
			selector.Registry(consulReg),
			selector.SetStrategy(selector.RoundRobin),
			)
		callAPI2(mySelector)
	}

}
