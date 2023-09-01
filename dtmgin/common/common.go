package common

import (
	"github.com/dtm-labs/client/dtmcli"
	"github.com/dtm-labs/client/dtmcli/dtmimp"
	"github.com/dtm-labs/logger"
	"github.com/gin-gonic/gin"
)

var BusiPort = 9090
var BusiConf = dtmcli.DBConf{
	Driver:   "mysql",
	Host:     "localhost",
	Port:     3306,
	User:     "root",
	Password: "1234qwer",
	//Db:       "dtm_busi",
}

func SagaAdjustBalance(db dtmcli.DB, uid int64, amount float64) error {

	_, err := dtmimp.DBExec(BusiConf.Driver, db, "update dtm_busi.user_account set balance = balance + ? where user_id = ?", amount, uid)
	return err
}

func SagaAdjustBalance1(db dtmcli.DB, uid int64, amount float64) error {

	_, err := dtmimp.DBExec(BusiConf.Driver, db, "update user.user_account set balance = balance + ? where user_id = ?", amount, uid)
	return err
}

func MustBarrierFromGin(c *gin.Context) *dtmcli.BranchBarrier {
	ti, err := dtmcli.BarrierFromQuery(c.Request.URL.Query())
	logger.FatalIfError(err)
	return ti
}
