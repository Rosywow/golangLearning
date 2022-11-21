package cmn

import (
	"context"
	"regexp"
)

//ServeEndPoint define the service
type ServeEndPoint struct {
	Developer *ModuleAuthor `json:"developer"`

	//Path required, the service url must be unique
	Path string `json:"path,omitempty"`

	//Fn process function
	Fn func(context.Context) `json:"fn,omitempty"`

	/* Priority execute order in stack: 0 is highest, lower than 100 or
	larger than 10000 is utility service and ep.match always return true */
	Priority int `json:"priority,omitempty"`

	//PathMatcher required, the url path regexp matcher
	PathMatcher *regexp.Regexp `json:"path_matcher,omitempty"`

	//PathPattern required, the url path regular expression
	PathPattern string `json:"path_pattern,omitempty"`

	//IsFileServe is static html file service,
	// true: as the file service
	// false: call fn for service
	IsFileServe bool `json:"is_file_serve,omitempty"`

	AllowDirectoryList bool `json:"allow_directory_list"`
	//DocRoot static html file service root directory
	DocRoot string `json:"doc_root,omitempty"`

	//PageRoute 是否支持前端页面路由，即angular/vue/svelte等的前端路由,如果
	//  支持: 如果请求的路径未发现则返回路径及上级路径包含的index.html,
	// 			例如，请求的是 /a/b/c/d,如果没有发现d或d/index.html，则
	//			依次返回先找到的/a/b/c/index.html,/a/b/index.html,/a/index.html
	//  不支持: 如果请求的路径未发现则返回状态404
	PageRoute bool `json:"page_route,omitempty"`

	//WhiteList if true then no authorization/authentication needed
	WhiteList bool `json:"white_list,omitempty"`

	//LoginPath redirect to log in when  needed
	LoginPath string `json:"login_path,omitempty"`

	//Name required, the api name for debug only
	Name string `json:"name,omitempty"`

	MaintainerID int64 `json:"maintainer_id,omitempty"`

	//该功能属于的域(业务域/子系统/客户)
	DomainID int64 `json:"domain_id,omitempty"`

	//level "0": 无组/角色/数据限制, 可访问全部数据
	//level "2": 机构#角色级别, 实现了不同角色授权，但不控制数据范围
	//level "4": 机构#角色$ID, 实现了不同角色授权，可控制 creator || all
	//level "8": 机构.DEPT#角色$ID, 实现了不同角色授权，可控制 creator || GRPs */
	AccessControlLevel string `json:"access_control_level,omitempty"`

	//该功能默认属于的域(业务域/子系统/客户)
	DefaultDomain int64 `json:"default_domain,omitempty"`
}
