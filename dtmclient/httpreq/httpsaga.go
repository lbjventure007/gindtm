package httpreq

import (
	"dtmclient/dtmutil"
	"dtmclient/request"
	"github.com/dtm-labs/client/dtmcli"
	"github.com/lithammer/shortuuid/v3"
	"log"
)

func HttpSaga() {

	req := &request.TranInReq{Amount: -16, UserId: 1}
	req1 := &request.TranInReq{Amount: 16, UserId: 2}

	saga := dtmcli.NewSaga(dtmutil.DefaultHTTPServer, shortuuid.New())
	saga.WaitResult = true

	saga.Add(ApiBaseUrl+"/sagatranout",
		ApiBaseUrl+"/sagatranoutrevert", req).
		Add(ApiBaseUrl+"/sagatranin",
			ApiBaseUrl+"/sagatraninrevert", req1)
	log.Println("saga busi trans submit")
	err := saga.Submit()
	log.Println("result gid is: %s", saga.Gid)
	log.Println(err)

}
