//Package user management
package user

//annotation:user-mgmt-service
//author:{"name":"user","tel":"18928776452","email":"XUnion@GMail.com"}

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
		z.Info("user zLogger settled")
	})
}

func Enroll(author string) {
	z.Info("user.Enroll called")

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
		Fn: user,

		Path: "/user",
		Name: "user",

		Developer: developer,
		WhiteList: true,

		DomainID:      int64(cmn.CDomainSys),
		DefaultDomain: int64(cmn.CDomainSys),
	})
}

func user(ctx context.Context) {
	q := cmn.GetCtxValue(ctx)
	z.Info("---->" + cmn.FncName())
	q.Msg.Msg = cmn.FncName()
	q.Resp()
}
