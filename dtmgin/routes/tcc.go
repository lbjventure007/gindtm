package routes

import (
	"database/sql"
	"dtmgin/common"
	"dtmgin/dtmutil"
	"dtmgin/request"
	"fmt"
	"github.com/dtm-labs/client/dtmcli/dtmimp"
	"github.com/gin-gonic/gin"
)

func TccTransOutTry() gin.HandlerFunc {
	return dtmutil.WrapHandler(func(c *gin.Context) interface{} {
		barrierFromGin := MustBarrierFromGin(c)
		return barrierFromGin.CallWithDB(pdbGet(), func(tx *sql.Tx) error {
			req := request.TranInReq{}
			err := c.BindJSON(&req)
			if err != nil {
				return dtmimp.ErrFailure
			}
			if req.Amount < 0 {
				return dtmimp.ErrFailure
			}
			affected, err1 := dtmimp.DBExec(common.BusiConf.Driver, tx,
				`update dtm_busi.user_account	set trading_balance=trading_balance+?
				where user_id=? and trading_balance + ? + balance >= 0`,
				-req.Amount, req.UserId, -req.Amount)
			if affected == 0 {
				return dtmimp.ErrFailure
			}
			return err1
		})
	})
}

func TccTransOutConfirm() gin.HandlerFunc {
	return dtmutil.WrapHandler(func(c *gin.Context) interface{} {
		barrierFromGin := MustBarrierFromGin(c)
		return barrierFromGin.CallWithDB(pdbGet(), func(tx *sql.Tx) error {
			req := request.TranInReq{}
			err := c.BindJSON(&req)
			if err != nil {
				return dtmimp.ErrFailure
			}
			if req.Amount < 0 {
				return dtmimp.ErrFailure
			}
			affected, err1 := dtmimp.DBExec(common.BusiConf.Driver, tx, `update dtm_busi.user_account
		set trading_balance=trading_balance-?,
		balance=balance+? where user_id=?`, -req.Amount, -req.Amount, req.UserId)
			if err1 == nil && affected == 0 {
				return fmt.Errorf("update user_account 0 rows")
			}
			return err
		})
	})
}

func TccTransOutRevert() gin.HandlerFunc {

	return dtmutil.WrapHandler(func(c *gin.Context) interface{} {
		barrierFromGin := MustBarrierFromGin(c)
		return barrierFromGin.CallWithDB(pdbGet(), func(tx *sql.Tx) error {
			req := request.TranInReq{}
			err := c.BindJSON(&req)
			if err != nil {
				return dtmimp.ErrFailure
			}
			affected, err1 := dtmimp.DBExec(common.BusiConf.Driver, tx,
				`update dtm_busi.user_account	set trading_balance=trading_balance+?
				where user_id=? and trading_balance + ? + balance >= 0`,
				req.Amount, req.UserId, req.Amount)
			if affected == 0 {
				return dtmimp.ErrFailure
			}
			return err1
		})
	})
}

func TccTransInTry() gin.HandlerFunc {
	return dtmutil.WrapHandler(func(c *gin.Context) interface{} {
		barrierFromGin := MustBarrierFromGin(c)
		return barrierFromGin.CallWithDB(pdbGet(), func(tx *sql.Tx) error {
			req := request.TranInReq{}
			err := c.BindJSON(&req)
			if err != nil {
				return dtmimp.ErrFailure
			}
			if req.Amount < 0 {
				return dtmimp.ErrFailure
			}
			affected, err1 := dtmimp.DBExec(common.BusiConf.Driver, tx,
				`update user.user_account	set trading_balance=trading_balance+?
				where user_id=? and trading_balance + ? + balance >= 0`,
				req.Amount, req.UserId, req.Amount)
			if affected == 0 {
				return dtmimp.ErrFailure
			}
			return err1
		})
	})
}

func TccTransInConfirm() gin.HandlerFunc {
	return dtmutil.WrapHandler(func(c *gin.Context) interface{} {
		barrierFromGin := MustBarrierFromGin(c)
		return barrierFromGin.CallWithDB(pdbGet(), func(tx *sql.Tx) error {
			req := request.TranInReq{}
			err := c.BindJSON(&req)
			if err != nil {
				return dtmimp.ErrFailure
			}
			if req.Amount < 0 {
				return dtmimp.ErrFailure
			}
			affected, err1 := dtmimp.DBExec(common.BusiConf.Driver, tx, `update user.user_account
		set trading_balance=trading_balance-?,
		balance=balance+? where user_id=?`, req.Amount, req.Amount, req.UserId)
			if err1 == nil && affected == 0 {
				return fmt.Errorf("update user_account 0 rows")
			}
			return err
		})
	})
}

func TccTransInRevert() gin.HandlerFunc {
	return dtmutil.WrapHandler(func(c *gin.Context) interface{} {
		barrierFromGin := MustBarrierFromGin(c)
		return barrierFromGin.CallWithDB(pdbGet(), func(tx *sql.Tx) error {
			req := request.TranInReq{}
			err := c.BindJSON(&req)
			if err != nil {
				return dtmimp.ErrFailure
			}
			affected, err1 := dtmimp.DBExec(common.BusiConf.Driver, tx,
				`update user.user_account	set trading_balance=trading_balance+?
				where user_id=? and trading_balance + ? + balance >= 0`,
				-req.Amount, req.UserId, -req.Amount)
			if affected == 0 {
				return dtmimp.ErrFailure
			}
			return err1
		})
	})

}
