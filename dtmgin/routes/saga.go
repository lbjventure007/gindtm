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

func Sagatransout() gin.HandlerFunc {
	return dtmutil.WrapHandler(func(c *gin.Context) interface{} {
		fmt.Println("--111---")

		barrier := common.MustBarrierFromGin(c)
		//return barrier.CallWithDB(pdbGet(), func(tx *sql.Tx) error {
		//	return common.SagaAdjustBalance(tx, req.UserId, req.Amount)
		//})
		return barrier.CallWithDB(pdbGet(), func(tx *sql.Tx) error {
			req := &request.TranInReq{}
			err := c.BindJSON(req)
			fmt.Println(err)
			fmt.Println("---1111-----1111", req.Amount)
			affetced, err1 := dtmimp.DBExec(common.BusiConf.Driver, tx, "update dtm_busi.user_account set balance = balance + ? where user_id = ?", req.Amount, req.UserId)
			fmt.Println("---11111-----2222", err1, affetced)
			if affetced == 0 {
				return dtmcli.ErrFailure
			}
			return err1
		})
	})
}

func SagatransOutRevert() gin.HandlerFunc {
	return dtmutil.WrapHandler(func(c *gin.Context) interface{} {
		fmt.Println("--3333---")

		req := &request.TranInReq{}
		err := c.BindJSON(req)
		fmt.Println(err)
		barrier := common.MustBarrierFromGin(c)

		fmt.Println("---33333-----22")
		//barrier.QueryPrepared(db).Error()
		return barrier.CallWithDB(pdbGet(), func(tx *sql.Tx) error {
			fmt.Println("---33333-----333", req.Amount)
			_, err = dtmimp.DBExec(common.BusiConf.Driver, tx, "update dtm_busi.user_account set balance = balance + ? where user_id = ?", -req.Amount, req.UserId)
			fmt.Println("---33333-----444", err)
			return err
		})
	})
}

func Sagatransin() gin.HandlerFunc {
	return dtmutil.WrapHandler(func(c *gin.Context) interface{} {
		fmt.Println("--2222---")

		barrier := common.MustBarrierFromGin(c)
		return barrier.CallWithDB(pdbGet(), func(tx *sql.Tx) error {
			//return errors.New("测试失败")
			req := &request.TranInReq{}
			err := c.BindJSON(req)

			if req.Amount > 10 {
				return dtmcli.ErrFailure
			}
			fmt.Println("---22222-----111", req.Amount, err)
			affetced, err1 := dtmimp.DBExec(common.BusiConf.Driver, tx, "update user.user_account set balance = balance + ? where user_id = ?", req.Amount, req.UserId)
			fmt.Println("---2222-----2222", err1, affetced)
			if affetced == 0 {
				return dtmcli.ErrFailure
			}

			return err1
			//return common.SagaAdjustBalance1(tx, req.UserId, req.Amount)
		})
	})
}
func SagatransInRevert() gin.HandlerFunc {
	return dtmutil.WrapHandler(func(c *gin.Context) interface{} {
		fmt.Println("--4444---")

		req := &request.TranInReq{}
		err := c.BindJSON(req)
		fmt.Println(err)
		barrier := common.MustBarrierFromGin(c)
		return barrier.CallWithDB(pdbGet(), func(tx *sql.Tx) error {
			value, ok := c.Get("trans_req")
			if ok {
				fmt.Println("0000--", value)
			} else {
				fmt.Println("11111")
			}
			fmt.Println("44444000----------")
			_, err = dtmimp.DBExec(common.BusiConf.Driver, tx, "update user.user_account set balance = balance + ? where user_id = ?", -req.Amount, req.UserId)

			return err
		})

	})
}
