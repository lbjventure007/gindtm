package routes

import (
	"database/sql"
	"dtmgin/common"
	"dtmgin/dtmutil"
	"dtmgin/request"
	"github.com/dtm-labs/client/dtmcli"
	"github.com/gin-gonic/gin"
)

func XATransin() gin.HandlerFunc {
	return dtmutil.WrapHandler(func(c *gin.Context) interface{} {
		return dtmcli.XaLocalTransaction(c.Request.URL.Query(), common.BusiConf, func(db *sql.DB, xa *dtmcli.Xa) error {
			req := &request.TranInReq{}
			err := c.BindJSON(req)
			if err != nil {
				return err
			}
			//	return dtmimp.ErrFailure
			//	//	return errors.New("条件不满足 不能转入")
			//	if err != nil {
			//		return err
			//	}
			return common.SagaAdjustBalance(db, req.UserId, req.Amount)
		})
	})
}

func XATransout() gin.HandlerFunc {
	return dtmutil.WrapHandler(func(c *gin.Context) interface{} {
		return dtmcli.XaLocalTransaction(c.Request.URL.Query(), common.BusiConf, func(db *sql.DB, xa *dtmcli.Xa) error {
			req := &request.TranInReq{}
			err := c.BindJSON(req)
			if err != nil {
				return err
			}
			return common.SagaAdjustBalance(db, req.UserId, req.Amount)
		})
	})
}
