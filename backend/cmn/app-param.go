package cmn

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/types"
)

var appParam = struct {
	// 险种数据
	insurances *insurancesData
	// 系统应用参数数据
	tParamData *TParamData
	// 提供给前端的应用参数
	cacheParam *cacheParam
}{
	&insurancesData{
		mapData: make(map[int64]*TInsuranceTypes),
	},
	&TParamData{
		data: make(map[int64]*TParam),
	},
	&cacheParam{
		mapData: make(map[string]interface{}),
	},
}

func refreshInsuranceTypes() error {
	return appParam.insurances.refresh()
}

func getInsuranceByID(id int64) *TInsuranceTypes {
	return appParam.insurances.getByID(id)
}

func getInsuranceList() []*TInsuranceTypes {
	return appParam.insurances.getList()
}

func getInsuranceMap() map[int64]*TInsuranceTypes {
	return appParam.insurances.getMap()
}

// 深拷贝
func getInsuranceCopyMap() map[int64]*TInsuranceTypes {
	return appParam.insurances.copyMap()
}

func getInsuranceJSON() []byte {
	return appParam.insurances.getJSON()
}

func initSysParam() error {
	return appParam.tParamData.refresh()
}

func getParamByID(id int64) (p *TParam) {
	return appParam.tParamData.getByID(id)
}

func getSysParam() map[int64]*TParam {
	return appParam.tParamData.copyParam()
}

func getParamByName(name string) (p *TParam) {
	return appParam.tParamData.getByName(name)
}

type cacheParam struct {
	mapData  map[string]interface{}
	jsonData types.JSONText
	lock     sync.RWMutex
}

func getCacheParam() types.JSONText {
	return appParam.cacheParam.getJSONData()
}

func (c *cacheParam) getJSONData() types.JSONText {
	if len(c.jsonData) == 0 {
		err := c.refresh("all")
		if err != nil {
			z.Error(err.Error())
			return nil
		}
	}
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.jsonData
}

func getCacheParamByKey(key string) interface{} {
	return appParam.cacheParam.getByKey(key)
}

func (c *cacheParam) getByKey(key string) interface{} {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.mapData[key]
}

func refreshBaseParam(refreshOption string) error {
	return appParam.cacheParam.refresh(refreshOption)
}

// 指定刷新参数内容，如果为all或者为空则刷新所有
func (c *cacheParam) refresh(refreshOption string) (err error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	var refreshFunc = map[string]func() error{

		"baseParam": c.refreshBaseParam,

		"保险参数": func() error {
			c.mapData["保险参数"] = getInsuranceCopyMap()
			return nil
		},

		"价格设置": func() error {

			return nil
		},

		"议价设置": func() error {

			return nil
		},

		"健康告知书": func() error {
			var buf []byte
			s := "select addi from t_param where id=14000"
			row := pgxConn.QueryRow(context.Background(), s)
			err := row.Scan(&buf)
			if err != nil {
				z.Error(err.Error())
				return err
			}
			c.mapData["健康告知书"] = json.RawMessage(buf)
			return nil
		},
	}
	if refreshOption == "all" || refreshOption == "" {
		for _, function := range refreshFunc {
			err = function()
			if err != nil {
				z.Error(err.Error())
				return
			}
		}
	}
	if function, ok := refreshFunc[refreshOption]; ok {
		err = function()
		if err != nil {
			z.Error(err.Error())
			return
		}
	}
	var buf []byte
	buf, err = json.Marshal(&c.mapData)
	if err != nil {
		z.Error(err.Error())
		return
	}
	str := string(buf)
	//班级类型中以_\d_作为注释,需过滤
	// 但是价格设置刚好有字段命中_\d_,再做转换处理

	str = strings.ReplaceAll(str, "day_1_count_GT_50_", "replace_label1")
	str = strings.ReplaceAll(str, "day_2_count_GT_50_", "replace_label2")
	str = strings.ReplaceAll(str, "day_EQ_30_", "replace_label3")
	str = strings.ReplaceAll(str, "_LTE_50_price", "replace_label4")
	str = strings.ReplaceAll(str, "_GT_50_price", "replace_label5")

	r := regexp.MustCompile(`_\d+_`)
	str = r.ReplaceAllString(str, "")

	str = strings.ReplaceAll(str, "replace_label1", "day_1_count_GT_50_")
	str = strings.ReplaceAll(str, "replace_label2", "day_2_count_GT_50_")
	str = strings.ReplaceAll(str, "replace_label3", "day_EQ_30_")
	str = strings.ReplaceAll(str, "replace_label4", "_LTE_50_price")
	str = strings.ReplaceAll(str, "replace_label5", "_GT_50_price")

	c.jsonData = types.JSONText(str)
	return
}

func (c *cacheParam) refreshBaseParam() (err error) {
	//---后端参数刷新
	err = initSysParam()
	if err != nil {
		z.Error(err.Error())
		return
	}
	//----------------------
	//班级类型
	s := `select * from v_xkb_school_layout`
	var stmt *sqlx.Stmt
	stmt, err = sqlxDB.Preparex(s)
	if err != nil {
		z.Error(err.Error())
		return
	}
	defer stmt.Close()

	classTypes := make(map[string]map[string][]string)

	var rows *sqlx.Rows
	rows, err = stmt.Queryx()
	if err != nil {
		z.Error(err.Error())
		return
	}
	defer rows.Close()

	var p *TVXkbSchoolLayout

	for rows.Next() {
		var c TVXkbSchoolLayout
		err = rows.StructScan(&c)
		if err != nil {
			z.Error(err.Error())
			return
		}
		//schooType := make(map[string]map[string][]string)
		if p == nil || c.School.String != p.School.String {
			classTypes[c.School.String] = make(map[string][]string)
			p = &c
		}

		grade := fmt.Sprintf("_%d_%s", c.Gradeid.Int64, c.Grade.String)
		classTypes[c.School.String][grade] =
			append(classTypes[c.School.String][grade], c.Class.String)
	}

	paramQry := []paramDef{
		{belongto: 12000, name: "学校类型"},
		{belongto: 12004, name: "证件类型"},
		{belongto: 12006, name: "数据同步目标"},
		{belongto: 12008, name: "保险类型"},
		{belongto: 12010, name: "治疗方式"},
		{belongto: 12012, name: "性别"},
		{belongto: 12014, name: "与被保险人关系"},
		{belongto: 12016, name: "账户类型"},
		{belongto: 12018, name: "保险时间"},
		{belongto: 12020, name: "比赛/活动保险保险期间"},
		{belongto: 12022, name: "比赛/活动保险参数"},
		{belongto: 12024, name: "投保单位性质"},
		{belongto: 12026, name: "投保联系人职位"},
		{belongto: 12030, name: "比赛/活动保险参与人员类型"},
		{belongto: 12032, name: "学校性质"},
		// {belongto: 12034, name: "收费标准(校方)"},
		{belongto: 12036, name: "筛选器-订单状态"},
		{belongto: 12038, name: "筛选器-保单状态"},
		{belongto: 12040, name: "支付方式(校方)"},
		{belongto: 12042, name: "筛选器-学校类型"},
		{belongto: 12044, name: "筛选器-缴费状态"},
		{belongto: 12046, name: "筛选器-付款方式"},
		{belongto: 12048, name: "地区选择器-默认值"},
		{belongto: 12050, name: "文件标签"},
		{belongto: 12052, name: "餐饮场所责任保险子类别"},
		{belongto: 12054, name: "争议处理"},
		{belongto: 12056, name: "支付方式(太平洋)"},
		{belongto: 12058, name: "是否首次投保"},
		{belongto: 12060, name: "俱乐部/场地责任保险子类别"},
		{belongto: 12062, name: "营业性质"},
		{belongto: 12064, name: "场地使用性质"},
		{belongto: 12066, name: "泳池性质"},
		{belongto: 12068, name: "比赛/活动组织方责任险保险期间"},
		{belongto: 12070, name: "文件路径"},
		{belongto: 12072, name: "议价类型"},
		{belongto: 12074, name: "训练项目"},
		{belongto: 12076, name: "场地类型"},
		{belongto: 12078, name: "联系客服"},
		{belongto: 12080, name: "治疗结果"},
		{belongto: 12082, name: "出险原因"},
		{belongto: 12084, name: "教职员工职位"},
		{belongto: 12086, name: "更正状态"},
		{belongto: 12088, name: "更正类型"},
	}

	for _, v := range paramQry {
		var val interface{}
		val, err = getParam(0, "", v.belongto, "", false)
		if err != nil {
			z.Error(err.Error())
			return
		}

		rows, ok := val.([]TVParam)
		if !ok {
			err = fmt.Errorf("value isn't a []TVParam, while it should be")
			z.Error(err.Error())
			return
		}

		var options []*selectOptions
		for _, v := range rows {
			var n selectOptions
			n.Name = v.Name.String
			n.Value = v.Value.String
			options = append(options, &n)
		}
		c.mapData[v.name] = &options
	}

	c.mapData["班级类型"] = &classTypes

	//-------地区因为基本不会变,也放这里
	c.mapData["地区"] = &regionsJSON

	schoolSeriesCategory, schoolSeriesCategoryParams, err := getSchoolSeriesCategory()
	if err != nil {
		return err
	}
	c.mapData["学校类型(校方)"] = schoolSeriesCategory
	c.mapData["收费标准(校方)"] = schoolSeriesCategoryParams

	return
}

type schoolSeriesCategory struct {
	Name   string
	Value  string
	Remark string
	//备注 义务教育--0 非义务教育--2 两者都有--4
	RemarkInt int64
}

func getSchoolSeriesCategory() (val interface{}, val2 interface{}, err error) {

	s := `select id,name,value,data_type,remark,status from v_param where parent_id = 12042 order by id asc;`
	z.Info(s)

	stmt, err := sqlxDB.Preparex(s)
	if err != nil {
		z.Error(err.Error())
		return
	}
	defer stmt.Close()

	type schoolCategoryParam struct {
		Name  string
		Value []string
	}

	var remarkToInt = map[string]int64{
		"义务教育":       0,
		"非义务教育":      2,
		"义务教育,非义务教育": 4,
	}

	var schoolCategoryParams = []schoolCategoryParam{
		{Name: "义务教育"},
		{Name: "非义务教育"},
	}
	schoolCategoryParams[0].Value, schoolCategoryParams[1].Value = make([]string, 0), make([]string, 0)
	var vals []schoolSeriesCategory
	var rows *sqlx.Rows
	rows, err = stmt.Queryx()
	defer rows.Close()
	for rows.Next() {
		var param TVParam
		err = rows.StructScan(&param)
		if err != nil {
			z.Error(err.Error())
			return
		}
		var category = schoolSeriesCategory{
			Name:      param.Name.String,
			Value:     param.Value.String,
			Remark:    param.Remark.String,
			RemarkInt: remarkToInt[param.Remark.String],
		}
		vals = append(vals, category)
		switch remarkToInt[param.Remark.String] {
		case 0:
			schoolCategoryParams[0].Value = append(schoolCategoryParams[0].Value, param.Name.String)
		case 2:
			schoolCategoryParams[1].Value = append(schoolCategoryParams[1].Value, param.Name.String)
		case 4:
			schoolCategoryParams[0].Value = append(schoolCategoryParams[0].Value, param.Name.String)
			schoolCategoryParams[1].Value = append(schoolCategoryParams[1].Value, param.Name.String)
		}

	}

	return vals, schoolCategoryParams, err
}
