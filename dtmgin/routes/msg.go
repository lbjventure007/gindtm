package routes

import (
	"database/sql"
	"dtmgin/common"
	"dtmgin/dtmutil"
	"dtmgin/request"
	"fmt"
	"github.com/dtm-labs/client/dtmcli"
	"github.com/dtm-labs/client/dtmcli/dtmimp"
	"github.com/gin-gonic/gin"
)

func MsgTransIn() gin.HandlerFunc {
	return dtmutil.WrapHandler(func(c *gin.Context) interface{} {
		barrier := common.MustBarrierFromGin(c)
		fmt.Println("msg in")
		return barrier.CallWithDB(pdbGet(), func(tx *sql.Tx) error {
			//return errors.New("测试失败")
			req := &request.TranInReq{}
			err := c.BindJSON(req)
			fmt.Println("msg in 1")
			if req.Amount > 15 || err != nil {
				fmt.Println("msg in 2")
				return dtmcli.ErrFailure
			}

			affetced, err1 := dtmimp.DBExec(common.BusiConf.Driver, tx,
				"update user.user_account set balance = balance + ? where user_id = ?",
				req.Amount, req.UserId)

			fmt.Println("msg in 3", err1, affetced)
			if affetced == 0 {
				return dtmcli.ErrFailure
			}

			return err1
			//return common.SagaAdjustBalance1(tx, req.UserId, req.Amount)
		})
	})
}

func MsgTransOut() gin.HandlerFunc {
	return dtmutil.WrapHandler(func(c *gin.Context) interface{} {
		barrier := common.MustBarrierFromGin(c)
		fmt.Println("msg out 1")
		return barrier.CallWithDB(pdbGet(), func(tx *sql.Tx) error {
			//return errors.New("测试失败")
			req := &request.TranInReq{}
			err := c.BindJSON(req)
			fmt.Println("msg out 2")
			if req.Amount > 15 || err != nil {
				fmt.Println("msg out 3")
				return dtmcli.ErrFailure
			}

			affetced, err1 := dtmimp.DBExec(common.BusiConf.Driver, tx,
				"update dtm_busi.user_account set balance = balance + ? where user_id = ?",
				-req.Amount, req.UserId)

			fmt.Println("msg out 4", err1, affetced)
			if affetced == 0 {
				return dtmcli.ErrFailure
			}

			return err1
			//return common.SagaAdjustBalance1(tx, req.UserId, req.Amount)
		})
	})
}

func MsgQuery() gin.HandlerFunc {
	fmt.Println("msg query 1")
	return dtmutil.WrapHandler(func(c *gin.Context) interface{} {
		fmt.Println("msg query 2", c.Query("gid"))
		bb := MustBarrierFromGin(c)
		db := dbGet().ToSQLDB()
		return bb.QueryPrepared(db)
	})
}
