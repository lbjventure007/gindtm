package routes

import (
	"database/sql"
	"dtmgin/common"
	"dtmgin/dtmutil"
	"dtmgin/request"
	"fmt"
	"github.com/dtm-labs/client/dtmcli"
	"github.com/dtm-labs/client/dtmcli/dtmimp"
	"github.com/dtm-labs/client/workflow"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
)

var API_BASE_URL = "http://localhost:9090/api/"

func WorkflowInit() string {
	return API_BASE_URL + "/workflowResume"
}

func WorkflowResume() gin.HandlerFunc {
	return dtmutil.WrapHandler(func(ctx *gin.Context) interface{} {
		data, err := ioutil.ReadAll(ctx.Request.Body)
		log.Println(err)
		return workflow.ExecuteByQS(ctx.Request.URL.Query(), data)
	})
}

func WorkSagaTransOut() gin.HandlerFunc {
	return dtmutil.WrapHandler(func(c *gin.Context) interface{} {
		fmt.Println("--111---")

		barrier := common.MustBarrierFromGin(c)
		//return barrier.CallWithDB(pdbGet(), func(tx *sql.Tx) error {
		//	return common.SagaAdjustBalance(tx, req.UserId, req.Amount)
		//})
		return barrier.CallWithDB(pdbGet(), func(tx *sql.Tx) error {
			req := &request.TranReqWork{}
			err := c.BindJSON(req)
			if req.FromAmount > 1 {
				fmt.Println("amount over 10 fail")
				return dtmimp.ErrFailure
			}
			fmt.Println(err)
			fmt.Println("---1111-----1111", req.FromAmount)
			affetced, err1 := dtmimp.DBExec(common.BusiConf.Driver, tx, "update dtm_busi.user_account set balance = balance + ? where user_id = ?", req.FromAmount, req.FromUserId)
			fmt.Println("---11111-----2222", err1, affetced)
			if affetced == 0 {
				return dtmcli.ErrFailure
			}
			return err1
		})
	})
}

func WorkSagaTransOutRevert() gin.HandlerFunc {
	return dtmutil.WrapHandler(func(c *gin.Context) interface{} {
		fmt.Println("--3333---")

		req := &request.TranReqWork{}
		err := c.BindJSON(req)
		fmt.Println(err)
		barrier := common.MustBarrierFromGin(c)

		fmt.Println("---33333-----22")
		//barrier.QueryPrepared(db).Error()
		return barrier.CallWithDB(pdbGet(), func(tx *sql.Tx) error {

			fmt.Println("---33333-----333", req.FromAmount)
			affected, err := dtmimp.DBExec(common.BusiConf.Driver, tx, "update dtm_busi.user_account set balance = balance + ? where user_id = ?", -req.FromAmount, req.FromUserId)
			fmt.Println("---33333-----444", affected, err)
			if affected == 0 {
				return dtmimp.ErrFailure
			}

			return err
		})
	})
}

func WorkSagaTransIn() gin.HandlerFunc {
	return dtmutil.WrapHandler(func(c *gin.Context) interface{} {
		fmt.Println("--2222---")

		barrier := common.MustBarrierFromGin(c)
		return barrier.CallWithDB(pdbGet(), func(tx *sql.Tx) error {
			//return errors.New("测试失败")
			req := &request.TranReqWork{}
			err := c.BindJSON(req)

			if req.ToAmount > 10 {
				return dtmcli.ErrFailure
			}
			fmt.Println("---22222-----111", req.ToAmount, err)
			affetced, err1 := dtmimp.DBExec(common.BusiConf.Driver, tx, "update user.user_account set balance = balance + ? where user_id = ?", req.ToAmount, req.ToUserId)
			fmt.Println("---2222-----2222", err1, affetced)
			if affetced == 0 {
				return dtmcli.ErrFailure
			}

			return err1
			//return common.SagaAdjustBalance1(tx, req.UserId, req.Amount)
		})
	})
}
func WorkSagaTransInRevert() gin.HandlerFunc {
	return dtmutil.WrapHandler(func(c *gin.Context) interface{} {
		fmt.Println("--4444---")
		//return errors.New("测试失败")
		//	return dtmimp.ErrFailure
		req := &request.TranReqWork{}
		err := c.BindJSON(req)
		fmt.Println(err)
		barrier := common.MustBarrierFromGin(c)
		return barrier.CallWithDB(pdbGet(), func(tx *sql.Tx) error {

			fmt.Println("44444000----------")
			affecetd, err := dtmimp.DBExec(common.BusiConf.Driver, tx, "update user.user_account set balance = balance + ? where user_id = ?", -req.ToAmount, req.ToUserId)
			if affecetd == 0 {
				fmt.Println("in revert fail")
				return dtmimp.ErrFailure
			}
			return err
		})

	})
}
