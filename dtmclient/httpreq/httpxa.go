package httpreq

import (
	"dtmclient/dtmutil"
	"dtmclient/request"
	"fmt"
	"github.com/dtm-labs/client/dtmcli"
	"github.com/go-resty/resty/v2"
	"github.com/lithammer/shortuuid/v3"
	"log"
)

func HttpXa() {
	gid := shortuuid.New()
	err := dtmcli.XaGlobalTransaction(dtmutil.DefaultHTTPServer, gid, func(xa *dtmcli.Xa) (*resty.Response, error) {

		req := request.TranInReq{}
		req.UserId = 1
		req.Amount = 10

		req1 := request.TranInReq{}
		req1.UserId = 2
		req1.Amount = -10
		resp, err := xa.CallBranch(&req, ApiBaseUrl+"/transout")

		fmt.Println("---1")
		if err != nil {
			fmt.Println("---2")
			return resp, err
		}
		fmt.Println("---3")
		return xa.CallBranch(&req1, ApiBaseUrl+"/transin")
	})
	fmt.Println("---4")
	log.Println(err, "---"+gid)
}
