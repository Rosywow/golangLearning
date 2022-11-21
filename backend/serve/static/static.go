//Package static management
package static

//annotation:static-service
//author:{"name":"static","tel":"18928776452","email":"XUnion@GMail.com"}

import (
	"encoding/json"
	"github.com/spf13/viper"
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

	var staticDocRoot string
	if viper.IsSet("webServe.static.docRoot") {
		staticDocRoot = viper.GetString("webServe.static.docRoot")
	}

	z.Info("this is static.Enroll called")
	cmn.AddService(&cmn.ServeEndPoint{
		Path: "/",
		Name: "rootPath",

		IsFileServe: true,

		PageRoute: true,

		DocRoot: staticDocRoot,

		Developer: developer,
		WhiteList: true,

		DomainID:      int64(cmn.CDomainSys),
		DefaultDomain: int64(cmn.CDomainSys),
	})
}
