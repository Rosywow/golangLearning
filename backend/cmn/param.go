package cmn

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/types"

	"w2w.io/null"
)

type selectOptions struct {
	Name  string
	Value string
}
type value struct {
	ID   null.Int    `json:"id,omitempty"`
	Name null.String `json:"name,omitempty"`
}

type paramDef struct {
	belongto int64
	name     string
	values   []value
}

// var cacheParamMap = make(map[string]interface{})
// var cacheParam types.JSONText
// var cacheParamLock sync.RWMutex

type TParamData struct {
	// 单纯读TParam表的数据
	data map[int64]*TParam
	lock sync.RWMutex
}

func (t *TParamData) getByID(id int64) (p *TParam) {
	if len(t.data) == 0 {
		err := t.refresh()
		if err != nil {
			z.Error(err.Error())
			return
		}
	}
	t.lock.RLock()
	p = t.data[id]
	t.lock.RUnlock()
	return p
}

func (t *TParamData) copyParam() map[int64]*TParam {
	t.lock.RLock()
	defer t.lock.RUnlock()
	p := make(map[int64]*TParam)
	for k, v := range t.data {
		p[k] = v
	}
	return p
}

func (t *TParamData) refresh() (err error) {
	z.Info("---->" + FncName())
	t.lock.Lock()
	z.Info("正在获取系统参数")
	defer t.lock.Unlock()
	sql := `select * from t_param`
	var stmt *sqlx.Stmt
	stmt, err = sqlxDB.Preparex(sql)
	if err != nil {
		z.Error(err.Error())
		return
	}
	defer stmt.Close()
	var rows *sqlx.Rows
	rows, err = stmt.Queryx()
	if err != nil {
		z.Error(err.Error())
		return
	}
	defer rows.Close()

	for rows.Next() {
		var v TParam
		err = rows.StructScan(&v)
		if err != nil {
			z.Error(err.Error())
			return
		}
		t.data[v.ID.Int64] = &v
	}

	//---价格计算相关变量初始化
	matchLevel1 = t.data[13300].valueToInt64()
	matchLevel2 = t.data[13302].valueToInt64()
	matchLevel3 = t.data[13304].valueToInt64()
	matchLevel4 = t.data[13306].valueToInt64()
	matchLevel5 = t.data[13308].valueToInt64()
	contestReqStandard = float64(t.data[13320].valueToInt())
	if t.data[13380] != nil && t.data[13382] != nil {
		contestInsuredType1 = t.data[13380].Name
		contestInsuredType2 = t.data[13382].Name
	}
	//------
	return
}

func (t *TParamData) getByName(name string) (p *TParam) {
	t.lock.RLock()
	defer t.lock.RUnlock()
	for _, v := range t.data {
		if v.Name == name {
			return v
		}
	}
	return
}

func (p *TParam) valueToInt64() (v int64) {
	if p == nil {
		return
	}
	v, _ = strconv.ParseInt(p.Value.String, 10, 64)
	return
}

func (p *TParam) valueToInt() (v int) {
	if p == nil {
		z.Info("nil")
		return
	}
	v, _ = strconv.Atoi(p.Value.String)
	return
}

//将当前参数的所有子属参数Name提取成字符串切片
func (p *TParam) sonNameToStringSlice() (s []string) {
	if p == nil {
		return
	}
	gParam := getSysParam()
	for _, v := range gParam {
		if v.Belongto.Int64 == p.ID.Int64 {
			s = append(s, v.Name)
		}
	}
	return
}

// const paramURI = `(?i)^/api/param(/.*)?$`
// const paramURILen = len(paramURI)

// var rParamURI = regexp.MustCompile(paramURI)

func param(ctx context.Context) {
	q := GetCtxValue(ctx)
	// if (len(q.R.URL.Path)+12) < paramURILen ||
	// 	!rParamURI.MatchString(q.R.URL.Path) {
	// 	return
	// }
	z.Info("---->" + FncName())
	q.Stop = true
	switch strings.ToLower(q.R.Method) {
	case "delete":
	case "get":
		qry := q.R.URL.Query().Get("q")
		if qry == "PayChannel" {
			s := `select name from t_pay_account`
			r, err := sqlxDB.Queryx(s)
			if err != nil {
				z.Error(err.Error())
				return
			}
			defer r.Close()
			channels := make([]string, 0)
			for r.Next() {
				var c string
				err = r.Scan(&c)
				if err != nil {
					z.Error(err.Error())
					return
				}
				channels = append(channels, c)
			}

			q.Msg.Data, q.Err = json.Marshal(channels)
			if q.Err != nil {
				z.Error(q.Err.Error())
				q.RespErr()
				return
			}
			q.Resp()
			return
		}

		if qry == "" || qry == "undefined" || qry == "null" || qry == `""` {
			baseParam(ctx)
			return
		}

		if qry == "refresh" {
			//------保险数据重新获取

			z.Warn("开始刷新系统参数,将对系统参数进行加锁,刷新结束前不允许读取")

			//---险种数据刷新
			q.Err = refreshInsuranceTypes()
			if q.Err != nil {
				z.Error(q.Err.Error())
				q.RespErr()
				return
			}
			//---参数刷新
			q.Err = refreshBaseParam("all")
			if q.Err != nil {
				z.Error(q.Err.Error())
				q.RespErr()
				return
			}
			z.Info("系统参数刷新结束")
			q.Msg.Data = getCacheParam()
			q.Resp()
			return
		}

		var req ReqProto
		q.Err = json.Unmarshal([]byte(qry), &req)
		if q.Err != nil {
			z.Error(q.Err.Error())
			q.RespErr()
			return
		}

		if strings.ToLower(req.Action) != "select" {
			q.Err = fmt.Errorf("please specify select as school query action")
			z.Error(q.Err.Error())
			q.RespErr()
			return
		}

		var s TParam
		s.TableMap = &s
		q.Err = DML(&s.Filter, &req)
		if q.Err != nil {
			z.Error(q.Err.Error())
			q.RespErr()
			return
		}
		v, ok := s.QryResult.(string)
		if !ok {
			q.Err = fmt.Errorf("s.qryResult should be string, but it isn't")
			z.Error(q.Err.Error())
			q.RespErr()
			return
		}
		z.Info(v)
		q.Msg.RowCount = s.RowCount
		q.Msg.Data = types.JSONText(v)
		q.Resp()
	case "post":
	case "put":
	}

}

func baseParam(ctx context.Context) {
	q := GetCtxValue(ctx)
	z.Info("---->" + FncName())
	q.Msg.Data = getCacheParam()
	q.Resp()
}

func getParam(
	id int64, name string, parentID int64, parentName string,
	valuesOnly bool) (val interface{}, err error) {

	var filterExpr []string
	var values []interface{}

	for {
		if id > 0 {
			filterExpr = append(filterExpr, fmt.Sprintf(`id=$%d`, len(filterExpr)+1))
			values = append(values, id)
			break
		}
		if name != "" {
			filterExpr = append(filterExpr, fmt.Sprintf(`name=$%d`, len(filterExpr)+1))
			values = append(values, name)
		}

		if parentID > 0 {
			filterExpr = append(filterExpr, fmt.Sprintf(`parent_id=$%d`, len(filterExpr)+1))
			values = append(values, parentID)
			break
		}

		if parentName != "" {
			filterExpr = append(filterExpr, fmt.Sprintf(`parent_name=$%d`, len(filterExpr)+1))
			values = append(values, parentName)
			break
		}
		if len(filterExpr) > 0 {
			break
		}
		err = fmt.Errorf("please provide atleast one condition to retrieve parameter")
		z.Error(err.Error())
		return
	}

	s := `select id,name,value,data_type,remark,status from v_param where ` +
		strings.Join(filterExpr, " and ") +
		" order by id ASC"

	stmt, err := sqlxDB.Preparex(s)
	if err != nil {
		z.Error(err.Error())
		return
	}
	defer stmt.Close()

	var params []TVParam
	var rows *sqlx.Rows
	rows, err = stmt.Queryx(values...)
	defer rows.Close()
	for rows.Next() {
		var param TVParam
		err = rows.StructScan(&param)
		if err != nil {
			z.Error(err.Error())
			return
		}
		params = append(params, param)
	}

	if len(params) == 0 {
		err = fmt.Errorf("can't find param with %s", strings.Join(filterExpr, " and "))
		return
	}

	var dst []string
	if valuesOnly {
		if len(params) == 1 {
			val = params[0].Value.String
			return
		}

		for _, v := range params {
			dst = append(dst, v.Value.String)
		}
		val = dst
		return
	}

	val = params

	return
}
