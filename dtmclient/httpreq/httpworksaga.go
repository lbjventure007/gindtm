package httpreq

import (
	"dtmclient/dtmutil"
	"dtmclient/request"
	"github.com/dtm-labs/client/dtmcli"
	"github.com/dtm-labs/client/dtmcli/dtmimp"
	"github.com/dtm-labs/client/workflow"
	"github.com/lithammer/shortuuid/v3"
	"log"
)

//func HttpWorkSaga() {
//
//	wfName := "wf_saga_test"
//
//	err := workflow.Register(wfName, func(wf *workflow.Workflow, data []byte) error {
//
//		fmt.Println("regin in")
//
//		req := request.TranReqWork{}
//		err := json.Unmarshal(data, &req)
//		fmt.Println("regin in", req, err)
//		if err != nil {
//			return err
//		}
//
//		_, err1 := wf.NewBranch().OnRollback(func(bb *dtmcli.BranchBarrier) error {
//			_, err2 := wf.NewRequest().SetBody(req).Post(ApiBaseUrl + "/workSagaTransInRevert")
//			return err2
//		}).NewRequest().SetBody(req).Post(ApiBaseUrl + "/workSagaTransIn")
//		if err1 != nil {
//			return err1
//		}
//
//		_, err3 := wf.NewBranch().OnRollback(func(bb *dtmcli.BranchBarrier) error {
//			_, err2 := wf.NewRequest().SetBody(req).Post(ApiBaseUrl + "/workSagaTransOutRevert")
//			return err2
//		}).NewRequest().SetBody(req).Post(ApiBaseUrl + "/workSagaTransOut")
//		if err3 != nil {
//			return err3
//		}
//
//		return nil
//	})
//
//	fmt.Println(err)
//
//	reqParam := request.TranReqWork{}
//	reqParam.FromUserId = 1
//	reqParam.FromAmount = 10
//	reqParam.ToUserId = 2
//	reqParam.ToAmount = 10
//
//	gid := shortuuid.New()
//
//	//data, err := json.Marshal(&reqParam)
//	//if err != nil {
//	//	fmt.Println(err)
//	//	return
//	//}
//	//fmt.Println("gid data err ", gid, data, err)
//
//	err = workflow.Execute(wfName, gid, dtmimp.MustMarshal(&reqParam))
//	if err != nil {
//		fmt.Println("错误", err)
//		return
//	}
//	fmt.Println("success")
//}

func MustUnmarshalReqHTTP(data []byte) *request.TranReqWork {
	var req request.TranReqWork
	dtmimp.MustUnmarshal(data, &req)
	return &req
}
func HttpWorkSaga() {
	workflow.InitHTTP(dtmutil.DefaultHTTPServer, "http://localhost:9090/api/workflowResume")
	wfName := "wf_saga_barrier"
	err := workflow.Register(wfName, func(wf *workflow.Workflow, data []byte) error {
		req := MustUnmarshalReqHTTP(data)
		_, err := wf.NewBranch().OnRollback(func(bb *dtmcli.BranchBarrier) error {
			_, err := wf.NewRequest().SetBody(req).Post(ApiBaseUrl + "/workSagaTransInRevert")
			return err
		}).NewRequest().SetBody(req).Post(ApiBaseUrl + "/workSagaTransIn")

		if err != nil {
			return err
		}

		_, err = wf.NewBranch().OnRollback(func(bb *dtmcli.BranchBarrier) error {
			_, errs := wf.NewRequest().SetBody(req).Post(ApiBaseUrl + "/workSagaTransOutRevert")
			return errs
		}).NewRequest().SetBody(req).Post(ApiBaseUrl + "/workSagaTransOut")

		if err != nil {
			return err
		}
		return nil
	})

	req1 := request.TranReqWork{}
	req1.FromUserId = 1
	req1.FromAmount = 5
	req1.ToUserId = 2
	req1.ToAmount = 5
	gid := shortuuid.New()
	//	workflow.SetProtocolForTest("http")

	execute2, err := workflow.Execute2(wfName, gid, dtmimp.MustMarshal(&req1))

	log.Printf("gid %s result is: %v err is %v", gid, execute2, err)

}
