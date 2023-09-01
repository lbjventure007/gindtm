package httpreq

import (
	"database/sql"
	"dtmclient/common"
	"dtmclient/dtmutil"
	"dtmclient/request"
	"fmt"
	"github.com/dtm-labs/client/dtmcli"
	"github.com/dtm-labs/client/dtmcli/dtmimp"
	"github.com/lithammer/shortuuid/v3"
)

func HttpMsg() {

	gid := shortuuid.New()
	reqIn := request.TranInReq{}
	reqIn.UserId = 1
	reqIn.Amount = 5

	reqOut := request.TranInReq{}
	reqOut.UserId = 2
	reqOut.Amount = 5

	msg := dtmcli.NewMsg(dtmutil.DefaultHTTPServer, gid).
		Add(ApiBaseUrl+"/msgTransIn", reqIn)
	//Add(ApiBaseUrl+"/msgTransOut", reqOut)
	msg.WaitResult = true

	//msg.Prepare(ApiBaseUrl + "/msgQuery")
	err := msg.DoAndSubmitDB(ApiBaseUrl+"/msgQuery", dtmutil.DbGet(common.BusiConf).ToSQLDB(), func(tx *sql.Tx) error {

		if reqOut.Amount > 5 {
			return dtmimp.ErrFailure
		}
		affected, err := dtmimp.DBExec(common.BusiConf.Driver, tx, "update dtm_busi.user_account set balance = balance + ? where user_id = ?", -reqOut.Amount, reqOut.UserId)
		fmt.Println("-----1111", err, affected)

		if affected == 0 {
			return dtmimp.ErrFailure
		}

		return err
	})

	//err := msg.Submit()

	fmt.Println("msg---", err, gid)
}
