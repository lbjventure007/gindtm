package routes

import (
	"bytes"
	"dtmgin/common"
	"fmt"

	"github.com/dtm-labs/client/dtmcli/dtmimp"
	"github.com/dtm-labs/client/workflow"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
)

func GetApp() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	app := gin.New()
	app.Use(gin.Recovery())
	app.Use(func(c *gin.Context) {
		body := ""
		if c.Request.Body != nil {
			rb, err := c.GetRawData()
			dtmimp.E2P(err)
			if len(rb) > 0 {
				body = string(rb)
				c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(rb))
			}
		}
		log.Println("begin %s %s body: %s", c.Request.Method, c.Request.URL, body)
		c.Next()
	})
	app.GET("/api/ping", func(c *gin.Context) { c.JSON(200, map[string]interface{}{"msg": "pong"}) })

	//xa
	app.POST("/api/xatransin", XATransin())
	app.POST("/api/xatransout", XATransout())

	//saga in
	app.POST("/api/sagatranin", Sagatransin())
	app.POST("/api/sagatraninrevert", SagatransInRevert())

	//saga out
	app.POST("/api/sagatranout", Sagatransout())
	app.POST("/api/sagatranoutrevert", SagatransOutRevert())

	//tcc out
	app.POST("/api/tccTransOutTry", TccTransOutTry())
	app.POST("/api/tccTransOutConfirm", TccTransOutConfirm())
	app.POST("/api/tccTransOutRevert", TccTransOutRevert())

	// tcc in
	app.POST("/api/tccTransInTry", TccTransInTry())
	app.POST("/api/tccTransInConfirm", TccTransInConfirm())
	app.POST("/api/tccTransInRevert", TccTransInRevert())

	// msg trans
	app.POST("/api/msgTransOut", MsgTransOut())
	app.POST("/api/msgTransIn", MsgTransIn())
	app.POST("/api/msgQuery", MsgQuery())

	// workflow trans

	app.POST("/api/workSagaTransOut", WorkSagaTransOut())
	app.POST("/api/workSagaTransOutRevert", WorkSagaTransOutRevert())
	app.POST("/api/workSagaTransInRevert", WorkSagaTransInRevert())
	app.POST("/api/workSagaTransIn", WorkSagaTransIn())

	app.POST("/api/workflowResume", func(ctx *gin.Context) {

		fmt.Println("----api---work-----")
		log.Printf("workflowResume")
		data, err := ioutil.ReadAll(ctx.Request.Body)
		log.Println("workflowResume err:", err)
		workflow.ExecuteByQS(ctx.Request.URL.Query(), data)
	})
	//callback := WorkflowInit()

	return app
}

func Run(app *gin.Engine) {
	err := app.Run(fmt.Sprintf(":%d", common.BusiPort))
	if err != nil {

		log.Fatal(err)
	}
}
