//Package document management
package document

//author:{"name":"document","tel":"13580452503","email":"KManager@GMail.com"}
//annotation:document-service

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
	var staticDocRoot string
	if viper.IsSet("webServe.static.document") {
		staticDocRoot = viper.GetString("webServe.static.document")
	}
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

	z.Info("this is static.Enroll called")
	cmn.AddService(&cmn.ServeEndPoint{
		Path: "/document",
		Name: "document",

		IsFileServe: true,

		AllowDirectoryList: true,

		DocRoot: staticDocRoot,

		Developer: developer,
		WhiteList: true,

		DomainID:      int64(cmn.CDomainSys),
		DefaultDomain: int64(cmn.CDomainSys),
	})
}
