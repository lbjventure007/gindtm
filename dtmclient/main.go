package main

import (
	"dtmclient/httpreq"
	"fmt"
)

func main() {

	//httpreq.HttpSaga()
	//httpreq.HttpMsg()
	fmt.Println("start")
	httpreq.HttpWorkSaga()
	fmt.Println("end")
}
