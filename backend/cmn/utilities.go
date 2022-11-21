package cmn

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/types"
	"github.com/tidwall/pretty"
	"github.com/tidwall/sjson"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
	"w2w.io/null"
)

var PathSeparator = '/'
var DualPathSeparator = "//"

var terminateSignal chan os.Signal

type RunningSession struct {
	Api        string
	BeginTime  int64
	RemoteAddr string
}

var RunningSessions = make(map[int64]RunningSession)

//GetTerminateSignal return terminateSignal
func GetTerminateSignal() chan os.Signal {
	return terminateSignal
}
func FncName() (name string) {
	pc, _, _, ok := runtime.Caller(1)
	if ok {
		name = runtime.FuncForPC(pc).Name()
	}
	return
}

func SetQuitChannel(c chan os.Signal) {
	if c == nil {
		return
	}

	terminateSignal = c
}

//Terminate be called when application want to exit with
//  user command or serious error occurred
func Terminate(_ int) {
	if terminateSignal == nil {
		//z.Error("call Terminate with nil terminateSignal")
		return
	}
	terminateSignal <- syscall.SIGINT
}

//JsonWrite set json key/value to file by path
// if fn doesn't exist then create it
func JsonWrite(fn string, path string, value interface{}) (err error) {
	switch {
	case fn == "":
		err = fmt.Errorf("call JsonWrite with empty fn")
	case path == "":
		err = fmt.Errorf("call JsonWrite with empty path")
	case value == nil:
		err = fmt.Errorf("call JsonWrite with nil value")
	}

	if err != nil {
		D.Error(err.Error())
		return
	}

	buf, err := ioutil.ReadFile(fn)
	if os.IsNotExist(err) {
		err = nil
	}

	if err != nil {
		D.Error(err.Error())
		return
	}
	if len(buf) == 0 {
		buf = []byte("{}")
	}

	rsl, err := sjson.Set(string(buf), path, value)
	if err != nil {
		D.Error(err.Error())
		return
	}

	err = ioutil.WriteFile(fn, pretty.Pretty([]byte(rsl)), 0644)
	if err != nil {
		D.Error(err.Error())
	}
	return
}

var buildVer string

// SetBuildVer app version
// format: r${svn revision}.b${Jenkins build number}(datetime)
// only call by main
func SetBuildVer(argBuildVer string) {
	buildVer = argBuildVer
}

// GetBuildVer get system build version
func GetBuildVer() string {
	return buildVer
}

//JSONBEscape escape postgresql jsonb string
func JSONBEscape(v string) string {
	//v = strings.Replace(v, `\`, `\\`, -1)
	v = strings.Replace(v, `'`, `''`, -1)
	return `'` + v + `'`
}

const escape = `0123456789\"\abort?xeUu`

//PqEscape using postgresql $Tag$content$Tag$ to escape string
func PqEscape(v string) string {
	// '\\' == 92, '\''=39
	dst := make([]rune, 0, len(v)*2)
	if strings.ContainsRune(v, 92) || strings.ContainsRune(v, 39) {
		dst = append(dst, 'E')
	}
	dst = append(dst, 39)

	preIsBackslash := false
	for _, c := range v {
		//fmt.Printf("%c", c)
		// ' escape to \'
		if c == 39 {
			dst = append(dst, 92, 39)
			continue
		}

		//it's backslash only
		if preIsBackslash && !strings.ContainsRune(escape, c) {
			dst = append(dst, 92, c)
			preIsBackslash = false
			continue
		}

		dst = append(dst, c)
		if c == 92 && !preIsBackslash { //
			preIsBackslash = true
			continue
		}
		preIsBackslash = false
	}
	dst = append(dst, 39)
	return string(dst)
}

//UtilCleanup release resource
func UtilCleanup() {
	if rootDB != nil {
		D.Info("close boltdb")
		_ = rootDB.Close()
		D.Info("boltdb closed")
	}

	if sqlxDB != nil {
		D.Info("close sqlxDB")
		_ = sqlxDB.Close()
		D.Info("sqlxDB closed")
	}

	if redisPool != nil {
		D.Info("close redis")
		_ = redisPool.Close()
		D.Info("redis closed")
	}

	if pgxConn != nil {

		if CancelWaitDbNotify != nil {
			CancelWaitDbNotify()
		}

		D.Info("close pgxConn")
		pgxConn.Close()
		D.Info("pgxConn closed")
	}
}

//GetNowInMS get current time in millisecond
func GetNowInMS() int64 {
	return time.Now().UnixNano() / 1e6
}

func QryDBStatus(ctx context.Context) {
	q := GetCtxValue(ctx)
	z.Info("---->" + FncName())

	q.Msg.Data = []byte(DbState(nil))
	q.Resp()
}

func DbState(db interface{}) (dbStatus string) {
	var dbStat string
	switch d := db.(type) {
	case *sqlx.DB:
		stat := d.Stats()
		dbStat = fmt.Sprintf("   sqlx, Idle: %d, InUse: %d, OpenConnections: %d",
			stat.Idle, stat.InUse, stat.OpenConnections)
	case *pgxpool.Pool:
		stat := d.Stat()
		dbStat = fmt.Sprintf("pgxpool, Idle: %d, InUse: %d, OpenConnections: %d",
			stat.IdleConns(), stat.AcquiredConns(), stat.TotalConns())
	default:
		appStartTimeLayout := "2006-01-02 15:04:05.000-0700"

		s1 := sqlxDB.Stats()
		s2 := pgxConn.Stat()
		dbStatus = fmt.Sprintf(`{"sysLaunchOn":"%s", "sqlx":{"idle":%d,"inUse":%d,"openConnections":%d},
			"pgpool":{"idle":%d,"inUse":%d,"openConnections":%d},"runningSessions":%d}`,
			AppStartTime.Format(appStartTimeLayout),
			s1.Idle, s1.InUse, s1.OpenConnections,
			s2.IdleConns(), s2.AcquiredConns(), s2.TotalConns(), len(RunningSessions))
		return
	}
	z.Info(dbStat)
	return
}

//InDebugMode

func DebugMode(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:    "debug",
		Value:   strconv.FormatBool(InDebugMode),
		Expires: time.Now().Add(60 * 2 * time.Second)})
}

//IsZero implement isZero
type IsZero interface {
	IsZero() bool
}

//MarshalJSON implement null/zero omitempty
func MarshalJSON(s interface{}) ([]byte, error) {
	v := reflect.ValueOf(s)

	if v.Kind() == reflect.Ptr {
		v = reflect.ValueOf(s).Elem()
	}
	if v.Kind() != reflect.Struct {
		err := errors.New("parameter should be struct or pointer to struct")
		z.Error(err.Error())
		return nil, err
	}

	buf := bytes.NewBufferString("{")

	// written fields count
	writtenFieldsCount := 0
	for i := 0; i < v.NumField(); i++ {
		valueField := v.Field(i)
		structField := v.Type().Field(i)

		fieldName := structField.Name
		//skipped := false

		jsonTag, ok := structField.Tag.Lookup("json")
		if !ok {
			continue
		}

		parts := strings.Split(jsonTag, ",")
		fieldName = parts[0]

		if fieldName == "-" {
			continue
		}

		if v, ok := valueField.Interface().(IsZero); ok {
			if v.IsZero() {
				continue
			}
		}
		k := valueField.Kind()

		//Is empty array/map/slice
		if (k == reflect.Array || k == reflect.Map || k == reflect.Slice) && valueField.Len() == 0 {
			//Yes it is
			continue
		}

		//Is zero/empty field
		if (k != reflect.Array && k != reflect.Map && k != reflect.Slice) &&
			valueField.Interface() == reflect.Zero(valueField.Type()).Interface() {
			//Yes it is
			continue
		}

		dbDataTypeTag, ok := structField.Tag.Lookup("db")
		if !ok {
			continue
		}

		parts = strings.Split(dbDataTypeTag, ",")
		if len(parts) < 2 {
			continue
		}
		// dbDataType := strings.ToLower(parts[2])

		fieldValue := valueField.Interface()

		// switch d := fieldValue.(type) {
		// case types.JSONText:
		// 	s = strings.ReplaceAll(string(d), " ", "")

		// 	//Is empty json string
		// 	if (s == "[]" || s == "{}") && (dbDataType == "json" || dbDataType == "jsonb") {
		// 		fieldValue = "null"
		// 	}
		// }

		if i > 0 && writtenFieldsCount > 0 {
			buf.WriteString(",")
		}

		b, err := json.Marshal(fieldValue)
		if err != nil {
			return nil, err
		}
		buf.WriteString(`"` + fieldName + `":` + string(b))
		writtenFieldsCount++

	}

	if _, err := buf.WriteString("}"); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (r *TWxUser) getErrCode() int {
	if r.ErrCode.Valid {
		return int(r.ErrCode.Int64)
	}
	return 0
}

func (r *TWxUser) getErrMsg() string {
	if r.ErrMsg.Valid {
		return r.ErrMsg.NullString.String
	}
	return ""
}

/*InvalidEmptyNullValue
if null.String/Int/Float/Bool/Time is ""/0/0/false/0 then invalid it */
func InvalidEmptyNullValue(p interface{}) (err error) {
	if p == nil {
		err = fmt.Errorf("call invalidEmptyNullValue with v==nil")
		z.Error(err.Error())
		return
	}

	if reflect.TypeOf(p).Kind() != reflect.Ptr {
		err = fmt.Errorf("please call invalidEmptyNullValue with pointer to struct")
		z.Error(err.Error())
		return
	}

	s := reflect.ValueOf(p).Elem()
	t := reflect.TypeOf(s.Interface())
	if t.Kind() != reflect.Struct {
		err = fmt.Errorf("Call invalidEmptyNullValue' pointer must to be struct")
		z.Error(err.Error())
		return
	}

	for i := 0; i < s.NumField(); i++ {
		if !s.Field(i).CanInterface() {
			continue
		}

		switch d := s.Field(i).Interface().(type) {
		case null.String:
			if d.Valid && d.String == "" {
				if s.Field(i).IsValid() && s.Field(i).CanSet() &&
					s.Field(i).FieldByName("Valid").IsValid() &&
					s.Field(i).FieldByName("Valid").CanSet() {
					s.Field(i).FieldByName("Valid").SetBool(false)
				}
			}
		// case null.Int:
		// 	if d.Valid && d.Int64 == 0 {
		// 		if s.Field(i).IsValid() && s.Field(i).CanSet() &&
		// 			s.Field(i).FieldByName("Valid").IsValid() &&
		// 			s.Field(i).FieldByName("Valid").CanSet() {
		// 			s.Field(i).FieldByName("Valid").SetBool(false)
		// 		}
		// 	}
		// case null.Float:
		// 	if d.Valid && d.Float64 == 0 {
		// 		if s.Field(i).IsValid() && s.Field(i).CanSet() &&
		// 			s.Field(i).FieldByName("Valid").IsValid() &&
		// 			s.Field(i).FieldByName("Valid").CanSet() {
		// 			s.Field(i).FieldByName("Valid").SetBool(false)
		// 		}
		// 	}
		// case null.Bool:
		// 	if d.Valid && !d.Bool {
		// 		if s.Field(i).IsValid() && s.Field(i).CanSet() &&
		// 			s.Field(i).FieldByName("Valid").IsValid() &&
		// 			s.Field(i).FieldByName("Valid").CanSet() {
		// 			s.Field(i).FieldByName("Valid").SetBool(false)
		// 		}
		// 	}
		case null.Time:
			if d.Valid && d.Time.Unix() == 0 {
				if s.Field(i).IsValid() && s.Field(i).CanSet() &&
					s.Field(i).FieldByName("Valid").IsValid() &&
					s.Field(i).FieldByName("Valid").CanSet() {
					s.Field(i).FieldByName("Valid").SetBool(false)
				}
			}
		}
	}

	return
}

func generateNonce() string {
	var letterRunes = []rune("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, 32)
	n := len(letterRunes)
	for i := range b {
		b[i] = letterRunes[rand.Intn(n)]
	}
	return string(b)
}

func clnAddr(r *http.Request) string {
	if r == nil {
		z.Error("call getIPAddress with nil r")
		return ""
	}

	for _, h := range []string{"X-Forwarded-For", "X-Real-Ip"} {
		addresses := strings.Split(r.Header.Get(h), ",")
		// march from right to left until we get a public address
		// that will be the address right before our proxy.
		for i := len(addresses) - 1; i >= 0; i-- {
			ip := strings.TrimSpace(addresses[i])
			// header can contain spaces too, strip those out.
			realIP := net.ParseIP(ip)
			if !realIP.IsGlobalUnicast() || isPrivateSubnet(realIP) {
				// bad address, go to next
				continue
			}
			return ip
		}
	}

	return strings.Split(r.RemoteAddr, ":")[0]
}

//tradeNoNonce 创建交易订单号
func tradeNoNonce(nonceLength int) string {
	if nonceLength <= 0 {
		return ""
	}

	var letterRunes = []rune("1234567890-abcdefghijklmnopqrstuvwxyz_ABCDEFGHIJKLMNOPQRSTUVWXYZ|*")
	b := make([]rune, nonceLength)
	n := len(letterRunes)
	for i := range b {
		b[i] = letterRunes[rand.Intn(n)]
	}
	return string(b)
}

func tradeNoWithID(orderID int64) string {

	x := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(
		base64.StdEncoding.EncodeToString(
			[]byte(fmt.Sprintf("%x", orderID))), "+", "_"), "/", "-"), "=", "*")

	return fmt.Sprintf("%02x%s%s", len(x), x, tradeNoNonce(30-len(x)))
}

func idFromTradeNo(tradeNo string) (id int64, err error) {
	if len(tradeNo) < 2 {
		err = fmt.Errorf("len(tradeNo) < 2, %s", tradeNo)
		z.Error(err.Error())
		return
	}

	var idB64Len int64
	idB64Len, err = strconv.ParseInt(tradeNo[:2], 16, 8)
	if err != nil {
		z.Error(err.Error())
		return
	}

	if idB64Len <= 0 || idB64Len > 20 {
		err = fmt.Errorf("invalid ID length: %d", idB64Len)
		z.Error(err.Error())
		return
	}

	s := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(
		tradeNo[2:(idB64Len+2)], "_", "+"), "-", "/"), "*", "=")

	buf, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		z.Error(err.Error())
		return
	}

	id, err = strconv.ParseInt(string(buf), 16, 64)
	if err != nil {
		z.Error(err.Error())
	}
	return
}

var telNoPattern = regexp.MustCompile(`^1[3-9][0-9]{9}$`)

func verifyTelNO(telNO string) bool {
	return telNoPattern.MatchString(telNO)
}

func removeEmptyArrayElement(a []interface{}) []interface{} {
	if len(a) == 0 {
		return nil
	}

	for i := len(a) - 1; i >= 0; i-- {
		e := a[i]
		if e == nil {
			a = append(a[:i], a[i+1:]...)
			continue
		}

		switch q := a[i].(type) {
		case string:
			if q == "" {
				a = append(a[:i], a[i+1:]...)
				continue
			}
		case map[string]interface{}:
			q = removeObjEmptyField(q)
			if q == nil || len(q) == 0 {
				a = append(a[:i], a[i+1:]...)
				continue
			}
		case []interface{}:
			q = removeEmptyArrayElement(q)
			if q == nil || len(q) == 0 {
				a = append(a[:i], a[i+1:]...)
				continue
			}
		}
	}

	if len(a) == 0 {
		return nil
	}

	return a
}

func removeObjEmptyField(m map[string]interface{}) map[string]interface{} {
	if m == nil {
		return nil
	}

	for k, v := range m {
		if v == nil {
			delete(m, k)
			continue
		}

		switch d := v.(type) {
		case string:
			if d == "" {
				delete(m, k)
				continue
			}

		case map[string]interface{}:
			r := removeObjEmptyField(d)
			if r == nil {
				delete(m, k)
				continue
			}

		case []interface{}:
			if len(d) == 0 {
				delete(m, k)
				continue
			}
			d = removeEmptyArrayElement(d)
			if d == nil || len(d) == 0 {
				delete(m, k)
				continue
			}
		}
	}

	if len(m) == 0 {
		return nil
	}
	return m
}

var externalDomainsConf map[string]*TExternalDomainConf
var payAccounts map[string]*TPayAccount

func LoadPayAccount() {
	if sqlxDB == nil {
		err := "sqlxDB is nil"
		z.Error(err)
		panic(err)
	}
	externalDomainsConf = make(map[string]*TExternalDomainConf)
	payAccounts = make(map[string]*TPayAccount)

	s := `select id,app_id, app_type, app_name, tokens, 
		creator, create_time, domain_id, addi, remark, status
		from t_external_domain_conf
		order by id asc`
	rows, err := sqlxDB.Queryx(s)
	if err != nil {
		z.Error(err.Error())
		panic(err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var c TExternalDomainConf
		err = rows.StructScan(&c)
		if err != nil {
			z.Error(err.Error())
			panic(err.Error())
		}
		externalDomainsConf[c.AppID+"#"+c.AppType] = &c
	}

	s = `select type, name, app_id, account, key, cert, 
	creator, domain_id, addi, remark, status, 
	create_time, update_time
	from t_pay_account
	order by id asc`
	rows, err = sqlxDB.Queryx(s)
	if err != nil {
		z.Error(err.Error())
		panic(err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var c TPayAccount
		err = rows.StructScan(&c)
		if err != nil {
			z.Error(err.Error())
			panic(err.Error())
		}
		payAccounts[c.Name.String] = &c
	}
}

func backEndVer(ctx context.Context) {
	q := GetCtxValue(ctx)
	q.Stop = true
	ver := GetBuildVer()
	short := ver[:8]
	q.Msg.Data = types.JSONText(fmt.Sprintf(`{"version":"%s","short":"%s"}`, ver, short))
	z.Info(string(q.Msg.Data))
	q.Resp()
}
