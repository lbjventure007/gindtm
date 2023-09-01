package routes

import (
	"context"
	"database/sql"
	"dtmgin/common"
	"encoding/json"
	"fmt"
	"strings"
	sync "sync"
	"time"

	"dtmgin/dtmutil"
	"github.com/dtm-labs/client/dtmcli"
	"github.com/dtm-labs/client/dtmcli/dtmimp"
	"github.com/dtm-labs/client/dtmcli/logger"
	"github.com/dtm-labs/client/dtmgrpc"
	"github.com/dtm-labs/client/dtmgrpc/dtmgpb"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/go-resty/resty/v2"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func dbGet() *dtmutil.DB {
	return dtmutil.DbGet(common.BusiConf)
}

func pdbGet() *sql.DB {
	db, err := dtmimp.PooledDB(dtmimp.DBConf(common.BusiConf))
	logger.FatalIfError(err)
	return db
}

func txGet() *sql.Tx {
	db := pdbGet()
	tx, err := db.Begin()
	logger.FatalIfError(err)
	return tx
}

// ResetXaData will rollback all pending xa transaction
func ResetXaData() {
	if common.BusiConf.Driver != "mysql" {
		return
	}

	db := dbGet()
	type XaRow struct {
		Data string
	}
	xas := []XaRow{}
	db.Must().Raw("xa recover").Scan(&xas)
	for _, xa := range xas {
		db.Must().Exec(fmt.Sprintf("xa rollback '%s'", xa.Data))
	}
}

// MustBarrierFromGin 1
func MustBarrierFromGin(c *gin.Context) *dtmcli.BranchBarrier {
	ti, err := dtmcli.BarrierFromQuery(c.Request.URL.Query())
	logger.FatalIfError(err)
	return ti
}

// MustBarrierFromGrpc 1
func MustBarrierFromGrpc(ctx context.Context) *dtmcli.BranchBarrier {
	ti, err := dtmgrpc.BarrierFromGrpc(ctx)
	logger.FatalIfError(err)
	return (*dtmcli.BranchBarrier)(ti)
}

// string2DtmError translate string to dtm error
func string2DtmError(str string) error {
	return map[string]error{
		dtmcli.ResultFailure: dtmcli.ErrFailure,
		dtmcli.ResultOngoing: dtmcli.ErrOngoing,
		dtmcli.ResultSuccess: nil,
		"":                   nil,
	}[str]
}

// SetGrpcHeaderForHeadersYes interceptor to set head for HeadersYes
func SetGrpcHeaderForHeadersYes(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	if r, ok := req.(*dtmgpb.DtmRequest); ok && strings.HasSuffix(r.Gid, "HeadersYes") {
		logger.Debugf("writing test_header:test to ctx")
		md := metadata.New(map[string]string{"test_header": "test"})
		ctx = metadata.NewOutgoingContext(ctx, md)
	}
	return invoker(ctx, method, req, reply, cc, opts...)
}

// SetHTTPHeaderForHeadersYes interceptor to set head for HeadersYes
func SetHTTPHeaderForHeadersYes(c *resty.Client, r *resty.Request) error {
	if b, ok := r.Body.(*dtmimp.TransBase); ok && strings.HasSuffix(b.Gid, "HeadersYes") {
		logger.Debugf("set test_header for url: %s", r.URL)
		r.SetHeader("test_header", "yes")
	}
	return nil
}

// oldWrapHandler old wrap handler for test use of dtm
func oldWrapHandler(fn func(*gin.Context) (interface{}, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		began := time.Now()
		r, err := func() (r interface{}, rerr error) {
			defer dtmimp.P2E(&rerr)
			return fn(c)
		}()
		var b = []byte{}
		if resp, ok := r.(*resty.Response); ok { // if it is a response，the get the body
			b = resp.Body()
		} else if err == nil {
			b, err = json.Marshal(r)
		}

		if err != nil {
			logger.Errorf("%2dms 500 %s %s %s %s", time.Since(began).Milliseconds(), err.Error(), c.Request.Method, c.Request.RequestURI, string(b))
			c.JSON(500, map[string]interface{}{"code": 500, "message": err.Error()})
		} else {
			logger.Infof("%2dms 200 %s %s %s", time.Since(began).Milliseconds(), c.Request.Method, c.Request.RequestURI, string(b))
			c.Status(200)
			c.Writer.Header().Add("Content-Type", "application/json")
			_, err = c.Writer.Write(b)
			dtmimp.E2P(err)
		}
	}
}

var (
	rdb  *redis.Client
	once sync.Once
)

// RedisGet 1
func RedisGet() *redis.Client {
	once.Do(func() {
		logger.Debugf("connecting to client redis")
		rdb = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:6379", common.BusiConf.Host),
			Username: "root",
			Password: "",
		})
	})
	return rdb
}
