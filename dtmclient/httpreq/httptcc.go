package httpreq

import (
	"dtmclient/dtmutil"
	"dtmclient/request"
	"fmt"
	"github.com/dtm-labs/client/dtmcli"
	"github.com/go-resty/resty/v2"
	"github.com/lithammer/shortuuid/v3"
)

func HttpTcc() {
	gid := shortuuid.New()
	fmt.Println("---11")
	err := dtmcli.TccGlobalTransaction(dtmutil.DefaultHTTPServer, gid, func(tcc *dtmcli.Tcc) (*resty.Response, error) {

		reqIn := request.TranInReq{}
		reqIn.Amount = 13
		reqIn.UserId = 1

		branch, err := tcc.CallBranch(&reqIn, ApiBaseUrl+"/tccTransOutTry", ApiBaseUrl+"/tccTransOutConfirm", ApiBaseUrl+"/tccTransOutRevert")

		if err != nil {
			return branch, err
		}
		reqOut := request.TranInReq{}
		reqOut.Amount = 13
		reqOut.UserId = 2

		return tcc.CallBranch(&reqOut, ApiBaseUrl+"/tccTransInTry", ApiBaseUrl+"/tccTransInConfirm", ApiBaseUrl+"/tccTransInRevert")
	})
	fmt.Println(err, gid, "----")

}
