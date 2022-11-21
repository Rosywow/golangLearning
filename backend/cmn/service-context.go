package cmn

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/sessions"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

type ModuleAuthor struct {
	Name  string `json:"name"`
	Tel   string `json:"tel"`
	Email string `json:"email"`
	Addi  string `json:"addi"`
}

type ctxKey string

func (v ctxKey) String() string {
	return string(v)
}

const (
	CUnknownCaller       = iota
	CPcBrowserCaller     = 1 << 0
	CAndroidWxCaller     = 1 << 1
	CIOSWxCaller         = 1 << 2
	CMobileBrowserCaller = 1 << 3
	CMacWxCaller         = 1 << 4
	CWinWxCaller         = 1 << 5

	CUnknownWxCaller = 1 << 7
)

func GetCallerTypeName(i int) string {
	switch i {
	case CUnknownCaller:
		return "unknown"

	case CUnknownWxCaller:
		return "unknownWx"

	case CAndroidWxCaller:
		return "androidWx"

	case CIOSWxCaller:
		return "iOSWx"

	case CMacWxCaller:
		return "macWx"

	case CWinWxCaller:
		return "winWx"

	case CPcBrowserCaller:
		return "pcBrowser"

	case CMobileBrowserCaller:
		return "mobileBrowser"

	default:
		return "unknown"
	}
}

type ServiceCtx struct {
	Err  error // error occurred during process
	Stop bool  // should run next process

	Attacker  bool // the requester is an attacker
	WhiteList bool // the request path in white list

	Ep *ServeEndPoint

	//stack *stack

	Responded bool // Dose response written

	Session *sessions.Session // gorilla cookie's

	Redis redis.Conn

	W http.ResponseWriter
	R *http.Request

	DomainList []string

	Domains []TDomain

	//角色中是否有管理员角色
	IsAdmin bool

	//是否在请求的URL中包含了admin=true
	ReqAdminFnc bool

	WxUser  *TWxUser
	SysUser *TUser
	Msg     *ReplyProto

	CallerType int

	UserAgent string

	WxLoginProcessed bool

	//xkb *xkbCtx
	//reqScope map[string]interface{} // session variables

	TouchTime time.Time

	Channel chan []byte

	RoutineID int
	BeginTime time.Time

	Tag map[string]interface{}

	//用户访问系统所使用的角色
	Role int64

	//用户访问的模块类型: 未知类型，函数，同未知类型，后台管理员模块，前台普通用户模块
	ReqFnType int
}

func (v *ServiceCtx) RespErr() {
	if v.Responded {
		z.Error("responded")
		return
	}

	// for test session without http context
	if v.W == nil {
		return
	}

	v.Responded = true
	v.Stop = true
	if v.Err == nil {
		v.Err = fmt.Errorf("v.err is nil")
	}

	if v.Msg == nil {
		v.Msg = &ReplyProto{
			API:    v.R.URL.Path,
			Method: v.R.Method,
		}
	}

	v.Msg.Msg = v.Err.Error()
	if v.Msg.Status >= 0 {
		v.Msg.Status = -1
	}

	//-410xx都是权限或账号错，需要清除session后重新与数据库同步
	if (v.Msg.Status / 100) == -410 {
		CleanSession(context.WithValue(context.Background(),
			QNearKey, v))
	}

	buf, err := json.Marshal(v.Msg)
	if err != nil {
		z.Error(err.Error())
		_, _ = fmt.Fprintf(v.W, err.Error())
		return
	}

	s := string(buf)
	if len(v.Msg.Data) > 0 {
		trial := fmt.Sprintf(`{"trial":%s}`, string(v.Msg.Data))
		t := make(map[string]interface{})

		v.Err = json.Unmarshal([]byte(trial), &t)
		if v.Err != nil {
			z.Error(trial)
			z.Error(v.Err.Error())
			v.RespErr()
			return
		}
		s = s[:len(buf)-1] + `,"data":` + string(v.Msg.Data) + "}"
	}

	v.W.Header().Add("Content-Type", "application/json")
	_, _ = fmt.Fprintf(v.W, s)
}

func (v *ServiceCtx) Resp() {
	if v.Responded {
		z.Error("responded")
		return
	}

	// for test session without http context
	if v.W == nil {
		return
	}

	if v.Msg == nil {
		v.Msg = &ReplyProto{
			API:    v.R.URL.Path,
			Method: v.R.Method,
		}
		v.Err = errors.New("v.Msg is nil")
		v.RespErr()
		return
	}

	buf, err := json.Marshal(v.Msg)
	if err != nil {
		z.Error(err.Error())
		_, _ = fmt.Fprintf(v.W, err.Error())
		return
	}

	s := string(buf)
	if len(v.Msg.Data) > 0 {
		trial := fmt.Sprintf(`{"trial":%s}`, string(v.Msg.Data))
		t := make(map[string]interface{})

		v.Err = json.Unmarshal([]byte(trial), &t)
		if v.Err != nil {
			v.Msg.Data = nil
			z.Error(trial)
			z.Error(v.Err.Error())
			v.RespErr()
			return
		}
		s = s[:len(buf)-1] + `,"data":` + string(v.Msg.Data) + "}"
	}

	v.W.Header().Add("Content-Type", "application/json")
	_, _ = fmt.Fprintf(v.W, "%s", s)

	v.Responded = true
}

const QNearKey = ctxKey("ServiceCtx")

func GetCtxValue(ctx context.Context) (q *ServiceCtx) {
	var err error
	f := ctx.Value(QNearKey)
	if f == nil {
		err = fmt.Errorf(`get nil from ctx.Value["%s"]`, QNearKey.String())
		z.Error(err.Error())
		panic(err.Error())
	}
	var ok bool
	q, ok = f.(*ServiceCtx)
	if !ok {
		err := fmt.Errorf("failed to type assertion for *ServiceCtx")
		z.Error(err.Error())
		panic(err.Error())
	}
	if q == nil {
		err := fmt.Errorf(`ctx.Value["%s"] should be non nil *ServiceCtx`, QNearKey.String())
		z.Error(err.Error())
		panic(err.Error())
	}
	return
}

func CleanSession(ctx context.Context) {
	q := GetCtxValue(ctx)
	userID, _ := q.Session.Values["ID"].(int64)
	if userID <= 0 {
		q.Err = fmt.Errorf("invalid session")
		z.Error(q.Err.Error())
		return
	}
	defer func() {
		z.Warn(fmt.Sprintf("%d 's session has been cleaned", userID))
	}()

	q.Session.Options.MaxAge = -1
	for k := range q.Session.Values {
		delete(q.Session.Values, k)
	}

	q.Err = q.Session.Save(q.R, q.W)
	if q.Err != nil {
		z.Error(q.Err.Error())
		return
	}

	q.Err = CleanCacheByUserID(userID)
	if q.Err != nil {
		z.Error(q.Err.Error())
	}
	if strings.ToLower(q.R.URL.Query().Get("erase")) == "true" {
		q.Err = EraseUser(userID)
	}
}

//EraseUser 抹除用户及其数据
//https://qnear.cn/api/dbStatus?xCleanSession=142857&erase=true
func EraseUser(userID int64) (err error) {
	s := []string{
		`delete from t_external_domain_user where user_id = $1`,
		`delete from t_insurance_policy where creator = $1`,
		`delete from t_order where creator = $1`,
		`delete from t_user_domain where sys_user = $1`,
		`delete from t_relation where left_id = $1 and left_type = 't_user.id'`,
		`delete from t_xkb_user where id = $1`,
		`delete from t_wx_user where id = $1`,
		`delete from t_user where id = $1`,
	}
	ctx := context.Background()
	tx, err := pgxConn.Begin(ctx)
	if err != nil {
		z.Error(err.Error())
		return
	}
	defer func() { _ = tx.Rollback(ctx) }()

	for _, v := range s {
		_, err = tx.Exec(ctx, v, userID)
		if err != nil {
			z.Error(err.Error())
			return
		}
	}

	func() { _ = tx.Commit(ctx) }()
	return
}

func CleanCacheByUserID(userID int64) (err error) {
	if userID <= 0 {
		err = fmt.Errorf("zero/Invalid userID")
		z.Error(err.Error())
		return
	}

	r := GetRedisConn()
	var keys []interface{}
	for {
		key := fmt.Sprintf("%s:%d", CSysUserByID, userID)
		var account string
		account, err = redis.String(r.Do("get", key))
		if err != nil {
			z.Error(err.Error())
			return
		}
		keys = append(keys, key)

		key = fmt.Sprintf("%s:%s", CSysUserByName, account)
		var userData string
		userData, err = redis.String(r.Do("JSON.GET", key, "."))
		if err != nil {
			z.Error(err.Error())
			return
		}

		rX := gjson.Get(userData, "MobilePhone")
		if rX.Exists() && rX.Num > 0 {
			userData, _ = sjson.Set(userData, "MobilePhone",
				fmt.Sprintf("%d", int64(rX.Num)))
		}

		keys = append(keys, key)

		var u TUser
		err = json.Unmarshal([]byte(userData), &u)
		if err != nil {
			z.Error(err.Error())
			return
		}

		key = fmt.Sprintf("%s:%s", CSysUserByTel, u.MobilePhone.String)

		account, err = redis.String(r.Do("get", key))
		if err != nil {
			z.Warn(err.Error())
		}
		if account != "" {
			keys = append(keys, key)
		}

		key = fmt.Sprintf("%s:%s", CSysUserByEmail, u.Email.String)
		account, err = redis.String(r.Do("get", key))
		if err != nil {
			z.Warn(err.Error())
		}
		if account != "" {
			keys = append(keys, key)
		}

		key = fmt.Sprintf("%s:%d", CWxUserByID, userID)
		var unionID string
		unionID, err = redis.String(r.Do("get", key))
		if err != nil {
			z.Error(err.Error())
			break
		}
		keys = append(keys, key)

		key = fmt.Sprintf("%s:%s", CWxUserByUnionID, unionID)
		var wxUserData string
		wxUserData, err = redis.String(r.Do("JSON.GET", key, "."))
		if err != nil {
			z.Error(err.Error())
			break
		}
		keys = append(keys, key)

		var x TWxUser
		err = json.Unmarshal([]byte(wxUserData), &x)
		if err != nil {
			z.Error(err.Error())
			return
		}

		key = fmt.Sprintf("%s:%s", CWxUserByOpenID, x.MpOpenID.String)
		account, err = redis.String(r.Do("get", key))
		if err != nil {
			z.Warn(err.Error())
		}
		if account != "" {
			keys = append(keys, key)
		}
		break
	}

	for k, v := range keys {
		z.Info(fmt.Sprintf("%d:%s", k, v))
	}

	var reply interface{}
	reply, err = r.Do("DEL", keys...)
	if err != nil {
		z.Error(err.Error())
		return
	}

	keysDropped, ok := reply.(int64)
	if !ok {
		err = fmt.Errorf("reply should be a int, %v", keysDropped)
		z.Error(err.Error())
		return
	}

	z.Info(fmt.Sprintf("user(%d) cache cleaned", userID))
	return
}

func XCleanSession(w http.ResponseWriter, _ *http.Request) {
	cookie := http.Cookie{
		Name:   "qNearSession",
		Value:  "",
		Domain: viper.GetString("webServe.serverName"),
		Path:   "/",
		MaxAge: -1,
	}

	http.SetCookie(w, &cookie)
}

var (
	Services = make(map[string]*ServeEndPoint)

	serviceMutex sync.Mutex
)

var rIsAPI = regexp.MustCompile(`(?i)^/api/(.*)?$`)

func AddService(ep *ServeEndPoint) (err error) {
	for {
		if ep == nil {
			err = errors.New("ep is nil")
			break
		}

		if ep.Path == "" {
			err = errors.New("ep.path empty")
			break
		}

		if ep.PathPattern == "" {
			ep.PathPattern = fmt.Sprintf(`(?i)^%s(/.*)?$`, ep.Path)
		}
		ep.PathMatcher = regexp.MustCompile(ep.PathPattern)

		if ep.IsFileServe {
			if ep.DocRoot == "" {
				err = errors.New("must specify docRoot when ep.isFileServe equal true")
				break
			}

			if ep.Fn == nil {
				ep.Fn = WebFS
			}
		} else {
			if ep.Fn == nil {
				err = errors.New("must specify fn when ep.isFileServe equal false")
				break
			}

			if !rIsAPI.MatchString(ep.Path) {
				ep.Path = strings.ReplaceAll("/api/"+ep.Path, "//", "/")
			}
		}

		if ep.Name == "" {
			err = errors.New("must specify apiName")
			break
		}

		if ep.DomainID == 0 {
			ep.DomainID = int64(CDomainSys)
		}
		if ep.DefaultDomain == 0 {
			ep.DefaultDomain = int64(CDomainSys)
		}

		if ep.AccessControlLevel == "" {
			ep.AccessControlLevel = "0"
		}
		_, ok := Services[ep.Path]
		if ok {
			err = errors.New(fmt.Sprintf("%s[%s] already exists", ep.Path, ep.Name))
		}
		break
	}

	if err != nil {
		z.Error(err.Error())
		return
	}

	z.Info(ep.Name + " added")

	serviceMutex.Lock()
	defer serviceMutex.Unlock()
	
	Services[ep.Path] = ep
	return
}

func buildURL(r *http.Request) (dst string) {
	if r == nil {
		return
	}
	scheme := "https:"
	host := viper.GetString("webServe.serverName")

	if r.URL.Scheme != "" {
		dst = r.URL.Scheme + "//"
	} else {
		dst = scheme + "//"
	}
	if r.URL.Host != "" {
		dst = dst + r.URL.Host
	} else {
		dst = dst + host
	}
	if r.URL.Path != "" {
		dst = dst + r.URL.Path
	}
	if r.URL.RawQuery != "" {
		dst = dst + "?" + r.URL.RawQuery
	}
	if r.URL.Fragment != "" {
		dst = dst + "#" + r.URL.Fragment
	}
	return
}

func guessIdxFile(f string) []string {
	f = filepath.Clean(f)

	var idxList []string
	for {
		f, _ = filepath.Split(f)
		idxList = append(idxList, f+"index.html")
		n := filepath.Clean(f)
		if f != "" && f != n {
			f = n
			continue
		}
		break
	}
	return idxList
}

const webFilePattern = `^(/*\S+)*/*\S+\.\S+$`

var rWebFilePattern = regexp.MustCompile(webFilePattern)

/*WebFS static file serve
1 if found the request file then return it,
2 if we can not find the target and the q.Ep.PageRoute is true then return the guessed index.html,
3 else return 404 not found */
func WebFS(ctx context.Context) {
	q := GetCtxValue(ctx)
	q.Responded = true
	q.Stop = true
	z.Info("---->" + FncName())

	if q.Ep == nil {
		q.Err = fmt.Errorf("call jsFS with nil q.Ep")
		z.Error(q.Err.Error())
		q.RespErr()
		return
	}

	if q.Ep.DocRoot == "" {
		q.Err = fmt.Errorf("call jsFS with empty q.Ep.docRoot")
		z.Error(q.Err.Error())
		q.RespErr()
		return
	}

	if len(q.R.URL.Path) < len(q.Ep.Path) {
		q.Err = fmt.Errorf("len(q.R.URL.Path) < len(q.Ep.path), it shouldn't happen")
		z.Error(q.Err.Error())
		q.RespErr()
		return
	}

	//---------
	if origin := q.R.Header.Get("Origin"); origin != "" {
		q.W.Header().Set("Access-Control-Allow-Origin", origin)
		q.W.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		q.W.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		q.W.Header().Set("Vary", "Origin")
		q.W.Header().Set("Access-Control-Allow-Credentials", "true")

		// Stop here if its Preflighted OPTIONS request
		if q.R.Method == "OPTIONS" {
			return
		}
	}

	//---------
	fsRoot := filepath.Clean(q.Ep.DocRoot)
	epPath := strings.TrimSuffix(q.Ep.Path, "/")
	urlPath := strings.TrimSuffix(q.R.URL.Path, "/")
	f := strings.TrimPrefix(urlPath[len(epPath):], "/")
	targetFileName := filepath.Clean(fsRoot + string(PathSeparator) + f)

	fileInfo, err := os.Stat(targetFileName)
	if os.IsNotExist(err) {
		z.Warn("InExistence: " + targetFileName)

		if rWebFilePattern.Match([]byte(f)) || !q.Ep.PageRoute {
			// missing the request specific file or non page route app.
			http.NotFound(q.W, q.R)
			return
		}

		// request is a path and q.Ep.PageRoute is true

		var idxHtmlFound bool
		idxList := guessIdxFile("/" + f)
		for _, v := range idxList {
			targetFileName = fsRoot + v
			_, err := os.Stat(targetFileName)
			if os.IsNotExist(err) {
				continue
			}

			if err != nil {
				q.Err = err
				z.Error(err.Error())
				q.Responded = false
				q.RespErr()
				return
			}
			idxHtmlFound = true
			break
		}

		if !idxHtmlFound {
			http.NotFound(q.W, q.R)
			return
		}
	}

	if err == nil && fileInfo.IsDir() && !q.Ep.AllowDirectoryList {
		idxHTML := filepath.Clean(targetFileName + string(PathSeparator) + "index.html")
		if _, err := os.Stat(idxHTML); os.IsNotExist(err) {
			q.W.WriteHeader(http.StatusForbidden)
			_, _ = q.W.Write([]byte("Access to the resource is forbidden!"))
			return
		}
	}

	http.ServeFile(q.W, q.R, targetFileName)
}
