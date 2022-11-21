package authmgmt

//author: {"name":"auth","email":"XUnion@GMail.com"}
//annotation:auth-mgmt-service

import (
	"context"
	"encoding/json"
	"go.uber.org/zap"
	"w2w.io/cmn"
)

var z *zap.Logger

func init() {
	//Setup package scope variables, just like logger, db connector, configure parameters, etc.
	cmn.PackageStarters = append(cmn.PackageStarters, func() {
		z = cmn.GetLogger()
		z.Info("auth mgmt zLogger settled")
	})
}

func Enroll(author string) {
	z.Info("auth.Enroll called")
	var developer *cmn.ModuleAuthor
	if author != "" {
		var d cmn.ModuleAuthor
		err := json.Unmarshal([]byte(author), &d)
		if err != nil {
			z.Error(err.Error())
			return
		}
		developer = &d
	}

	cmn.AddService(&cmn.ServeEndPoint{
		Fn: authMgmt,

		Path: "/auth-mgr",
		Name: "authMgmt",

		Developer: developer,
		WhiteList: true,

		DomainID:      int64(cmn.CDomainSys),
		DefaultDomain: int64(cmn.CDomainSys),
	})
}

// authMgmt authenticate/authorization management
func authMgmt(ctx context.Context) {
	q := cmn.GetCtxValue(ctx)
	z.Info("---->" + cmn.FncName())
	q.Msg.Msg = cmn.FncName()
	q.Resp()
}
