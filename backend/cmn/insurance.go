package cmn

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"database/sql"
	"encoding/hex"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/types"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
	"w2w.io/null"

	"encoding/json"
	"w2w.io/excelize"
)

type insurancesData struct {
	list     []*TInsuranceTypes
	jsonData []byte
	mapData  map[int64]*TInsuranceTypes
	lock     sync.RWMutex
}

const (

	//内部识别使用的清单类型，不编码
	contestListType   = "比赛人员清单"
	traineeListType   = "实习生清单"
	studentListType   = "校方清单"
	staffListType     = "教工清单"
	schoolBusListType = "校车清单"
	groundListType    = "场地俱乐部清单"
	poolListType      = "游泳池清单"

	//====比赛保险
	insuranceContest = 10000
	//比赛二级类别（非实际的险种）
	insuranceContest2  = 10002
	insuranceContest4  = 10004
	insuranceContest6  = 10006
	insuranceContest8  = 10008
	insuranceContest10 = 10010
	insuranceContest14 = 10014
	insuranceContest12 = 10012 //各类活动

	insuranceSchoolSeries = 10020 //===校方系列(非实际险种)
	insuranceSchool       = 10022 //校方责任险
	insuranceStaff        = 10024 //教工责任险
	insuranceTrainee      = 10026 //实习生责任险
	insuranceSchoolBus    = 10028 //校车责任险
	insuranceCanteen      = 10030 //食堂责任险

	insuranceAccident = 10040 //学生意外责任险

	insuranceOrganization = 10060 //比赛/活动组织方责任险
	insuranceClub         = 10070 //俱乐部/场地责任险
	insuranceSwimmingPool = 10080 //游泳池责任险

	//订单状态
	oStatusDraft     = "00" //草稿
	oStatusChecked   = "04" //检查通过
	oStatusBargain   = "08" //申请议价
	oStatusRefused   = "12" //拒保
	oStatusReadyPay  = "16" //等待支付
	oStatusPaying    = "18" //开始支付/支付中（生成支付二维码/地址）
	oStatusPaid      = "20" //已支付
	oStatusReturn    = "24" //退保
	oStatusCancelled = "28" //作废

	//订单相关文件标签
	labelBusinessLicense = "营业执照"
	labelCreditCode      = "统一社会信用代码"
	labelCreditCodePic1  = "统一社会信用代码证书"
	labelCreditCodePic2  = "工会信用代码证书"

	labelInvoice          = "发票"
	labelInvBorrowScan    = "预借发票申请函盖章扫描件"
	labelList             = "投保清单"      //后端生成
	labelListExcel        = "投保清单excel" //后端生成
	labelListScan         = "投保清单盖章扫描件" //用户上传
	labelPolicyForm       = "投保单"       //后端生成
	labelPolicyFormExcel  = "投保单excel"  //后端生成
	labelPolicyFormScan   = "投保单盖章扫描件"  //用户上传
	labelPolicyTpl        = "投保单模板"
	labelListTpl          = "投保清单模板"
	labelPolicyExample    = "投保单示例" //价格设置前端上传,用于给用户投保展示
	labelApplications     = "批改申请书"
	labelApplicationsScan = "批改申请书盖章扫描件"
	labelEndorsement      = "批单"

	labelInvBorrow    = "预借发票申请函"
	labelTransferAuth = "转账授权说明"

	//续保新单
	prevPolicyNoNew = "新单"

	onlinePayTimeLimit = 23 //比赛开始前几点后不能在线支付
)

func getTraits(orderID int64) (traitSlice []string, err error) {
	if orderID <= 0 {
		err = fmt.Errorf("非法id:%d", orderID)
		z.Error(err.Error())
		return
	}
	s := `select traits from t_order where id = $1`

	row := sqlxDB.QueryRowx(s, orderID)
	var traits null.String
	err = row.Scan(&traits)
	if err == sql.ErrNoRows {
		err = nil
		return
	}
	if err != nil {
		z.Error(err.Error())
		return
	}

	return strings.Split(traits.String, ","), nil
}

func inSlice(e string, s []string) bool {
	for _, v := range s {
		if e == v {
			return true
		}
	}
	return false
}

//检查smallint 范围
func checkSmallInt(num int64) error {
	if num > (1<<15 - 1) {
		return fmt.Errorf("%d 数值过大，超出范围限制", num)
	}
	return nil
}

const insuranceTypeQuery = `select id,ref_id,data_type,
		parent_id,
		parent_name,
		org_id,
		org_name,
		layout_order,
		insurer,
		name,
		alias,
		pay_type,
		pay_name,
		pay_channel,
		rule_batch,
		unit_price,
		price,
		price_config,
		max_insure_in_year,
		insured_in_month,
		insured_start_time,
		allow_start,
		allow_end,
		indate_start,
		indate_end,
		age_limit,
		bank_account,
		bank_account_name,
		bank_name,
		bank_id,
		floor_price,
		define_level,
		layout_level,
		list_tpl,
		files,
		pic,
		sudden_death_description,
		description,
		auto_fill,
		enable_import_list,
		have_dinner_num,
		invoice_title_update_times,
		receipt_account,
		contact,
		contact_qr_code,other_files,
		transfer_auth_files,
		remind_days,
		mail,
		order_repeat_limit,
		group_by_max_day,
		web_description,
		mobile_description,
		auto_fill_param,
		interval,
		insured_end_time,
		create_time,
		update_time,
		addi,
		status,
		time_status
	from v_insurance_type`

//-----------

/*
 * @Author: zhangzihan
 * @Last Modified by: zhangzihan
 * @Last Modified time: 2020-10-09 00:53:15
 * @Comment: 保险类型表的控制,包含险种显示控制/获取,基础设置等
 */

type getInsuranceOption struct {
	ReqType int `json:"ReqType"`
	//ReqType=2 来自首页的请求，返回险种id,name,parent_id,pic
	//ReqType=4 来自险种资源管理/我的保单的请求，返回险种id,name
	//ReqType=6 来自订单/保单筛选的请求，返回险种

	InsuranceID string `json:"InsuranceID"` //请求险种的ID
}

//ReceiptAccount 对公账号设置
type ReceiptAccount struct {
	Bank        string `json:"Bank"`        //开户行
	BankNum     string `json:"BankNum"`     //行号
	AccountName string `json:"AccountName"` //户名
	Account     string `json:"Account"`     //账号
}

//Contact 短信联系人
type Contact struct {
	Phone       string `json:"Phone"`       //联系电话
	ContactName string `json:"ContactName"` //联系名
}

//Mail 邮寄地址
type Mail struct {
	Phone    string `json:"Phone"`    //邮寄联系电话
	Receiver string `json:"Receiver"` //收件人
	Address  string `json:"Address"`  //邮寄地址
}

//Underwriter 承保公司账号
type Underwriter struct {
	StartTime int64  `json:"StartTime"` //开始时间
	EndTime   int64  `json:"EndTime"`   //结束时间
	Company   string `json:"Company"`   //承保公司
	Account   string `json:"Account"`   //承保公司账号
	Password  string `json:"Password"`  //密码
}

//AutoFillParam 自动录单参数
type AutoFillParam struct {
	Company  string `json:"Company"`
	Reviewer string `json:"Reviewer"`
	Phone    string `json:"Phone"`
}

// TODO:删除该接口，前端判断显隐
func insuranceTypesList(ctx context.Context) {
	q := GetCtxValue(ctx)
	z.Info("---->" + FncName())
	q.Stop = true

	switch strings.ToLower(q.R.Method) {
	case "delete":
		id := q.R.URL.Query().Get("id")
		if id == "" {
			q.Err = fmt.Errorf("撤消报错申请必须提供报错申请编号")
			z.Error(q.Err.Error())
			q.RespErr()
			return
		}

		var pid int64
		//---------系统预设的不能删
		pid, q.Err = strconv.ParseInt(id, 10, 64)
		if pid < 20000 {
			q.Err = fmt.Errorf("无效编号: %d", pid)
			z.Error(q.Err.Error())
			q.RespErr()
			return
		}

		s := `delete from t_insurance_types where id=$1`
		var r sql.Result
		r, q.Err = sqlxDB.Exec(s, pid)
		if q.Err != nil {
			z.Error(q.Err.Error())
			q.RespErr()
			return
		}
		var d int64
		d, q.Err = r.RowsAffected()
		if q.Err != nil {
			z.Error(q.Err.Error())
			q.RespErr()
			return
		}

		s = fmt.Sprintf(`{"RowAffected":%d}`, d)

		z.Info(s)
		q.Msg.Data = []byte(s)
		q.Resp()
		return
	case "get":
		/*
			const req = {
			data: {
				ReqType: path,
			},
			action: 'select',
			};
		*/
		qry := q.R.URL.Query().Get("q")
		if qry == "" || qry == "undefined" || qry == "null" || qry == `""` {
			q.Err = fmt.Errorf("the param 'q' is empty")
			z.Error(q.Err.Error())
			q.RespErr()
			return
		}

		z.Info(qry)
		var req ReqProto
		q.Err = json.Unmarshal([]byte(qry), &req)
		if q.Err != nil {
			z.Error(q.Err.Error())
			q.RespErr()
			return
		}
		if strings.ToLower(req.Action) != "select" {
			q.Err = fmt.Errorf("please specify select as query action")
			z.Error(q.Err.Error())
			q.RespErr()
			return
		}

		if len(req.Data) == 0 {
			q.Err = fmt.Errorf("不指定data，你想干啥子哦？")
			z.Error(q.Err.Error())
			q.RespErr()
			return
		}

		var data getInsuranceOption
		q.Err = json.Unmarshal(req.Data, &data)
		if q.Err != nil {
			z.Error(q.Err.Error())
			q.RespErr()
			return
		}

		var showInfos []*TInsuranceTypes

		//layout_level 保险显示层次 0:显示在第一级 4:显示在第二级 8:显示在第三级 20:不显示
		//layout_order 保险显示顺序
		//define_level 保险实际层次: 0:仅供显示使用的险种 2:实际险种
		switch data.ReqType {
		case 2:

			//首页无需登录
			//查询:递归查询显示层级为0的险种,过滤不显示的险种
			for _, v := range getInsuranceList() {
				if v.LayoutLevel.Int64 != 20 {
					showInfos = append(showInfos, v)
				}
			}
			if showInfos == nil {
				q.Err = fmt.Errorf("险种数据读取获取失败")
				z.Error(q.Err.Error())
				q.RespErr()
				return
			}
			break
		case 4:
			//查询:递归查询显示层级为0的险种及其子险种
			//返回时过滤掉define_level不为2的(即去掉非真实的险种)

			//var authType, authData string
			//authType, authData, _, q.Err = createAuthFilter(ctx, "ID")
			//if q.Err != nil {
			//	q.RespErr()
			//	return
			//}
			//z.Info(fmt.Sprintf("authType: %s, authData: %s", authType, authData))

			for _, v := range getInsuranceList() {
				if v.DefineLevel.Int64 == 2 {
					showInfos = append(showInfos, v)
				}
			}
			if showInfos == nil {
				q.Err = fmt.Errorf("险种数据读取获取失败")
				z.Error(q.Err.Error())
				q.RespErr()
				return
			}
			break
		case 6:

			//查询:递归查询显示层级为0的险种及其子险种
			//返回时过滤掉显示层级大于等于8的险种(如比赛/活动保险的子险种,学意险的学校种类)

			//var authType, authData string
			//authType, authData, _, q.Err = createAuthFilter(ctx, "ID")
			//if q.Err != nil {
			//	z.Error(q.Err.Error())
			//	q.RespErr()
			//	return
			//}
			//z.Info(fmt.Sprintf("authType: %s, authData: %s", authType, authData))

			for _, v := range getInsuranceList() {
				if v.LayoutLevel.Int64 < 8 {
					showInfos = append(showInfos, v)
				}
			}
			if showInfos == nil {
				q.Err = fmt.Errorf("险种数据读取获取失败")
				z.Error(q.Err.Error())
				q.RespErr()
				return
			}
			break
		default:
			q.Err = fmt.Errorf("没有传入指定的reqType")
			z.Error(q.Err.Error())
			q.RespErr()
			return
		}

		q.Msg.RowCount = int64(len(showInfos))
		var buf []byte
		buf, q.Err = json.Marshal(&showInfos)
		if q.Err != nil {
			z.Error(q.Err.Error())
			q.RespErr()
			return
		}
		q.Msg.Data = types.JSONText(string(buf))
		q.Resp()
		return
	case "post":
	case "put":
	}
}

func insuranceTypes(ctx context.Context) {
	q := GetCtxValue(ctx)

	z.Info("---->" + FncName())
	q.Stop = true

	//var authType, authData string
	//authType, authData, _, q.Err = createAuthFilter(ctx, "")
	//if q.Err != nil {
	//	q.RespErr()
	//	return
	//}
	//z.Info(fmt.Sprintf("authType: %s, authData: %s", authType, authData))

	switch strings.ToLower(q.R.Method) {
	case "delete":
		idSet := q.R.URL.Query().Get("id")
		if idSet == "" {
			q.Err = fmt.Errorf("撤消报错申请必须提供报错申请编号")
			z.Error(q.Err.Error())
			q.RespErr()
			return
		}
		values := make([]interface{}, 0)
		var typeID int64

		for _, v := range strings.Split(idSet, ",") {
			typeID, q.Err = strconv.ParseInt(v, 10, 64)
			if typeID <= 15000 {
				q.Err = fmt.Errorf("无效编号: %d", typeID)
				z.Error(q.Err.Error())
				q.RespErr()
				return
			}
			values = append(values, typeID)
		}
		format := ""
		for n := 1; n <= len(values); n++ {
			format += "$" + strconv.Itoa(n) + ","
		}
		format = format[:len(format)-1]
		s := fmt.Sprintf(`delete from t_insurance_types where id in (%s)`, format)
		z.Info(s)
		var r sql.Result
		r, q.Err = sqlxDB.Exec(s, values...)
		if q.Err != nil {
			z.Error(q.Err.Error())
			q.RespErr()
			return
		}
		var d int64
		d, q.Err = r.RowsAffected()
		if q.Err != nil {
			z.Error(q.Err.Error())
			q.RespErr()
			return
		}

		s = fmt.Sprintf(`{"RowAffected":%d}`, d)
		q.Msg.Data = []byte(s)
		q.Resp()

	case "get":
		qry := q.R.URL.Query().Get("q")
		if qry == "" {
			q.Err = fmt.Errorf("the param 'q' is empty")
			z.Error(q.Err.Error())
			q.RespErr()
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
			q.Err = fmt.Errorf("please specify select as insurance query action")
			z.Error(q.Err.Error())
			q.RespErr()
			return
		}

		if urlParam := q.R.URL.Query().Get("rule"); urlParam == "true" {
			var s TVInsuranceType
			s.TableMap = &s
			q.Err = DML(&s.Filter, &req)
			if q.Err != nil {
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
			q.Msg.RowCount = s.RowCount
			q.Msg.Data = types.JSONText(v)
			q.Resp()
			return
		}

		// var s TInsuranceTypes
		var s TInsuranceTypes
		s.TableMap = &s
		q.Err = DML(&s.Filter, &req)
		if q.Err != nil {
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

		var insuranceSets []TInsuranceTypes
		jsonText := types.JSONText(v)
		q.Err = json.Unmarshal(jsonText, &insuranceSets)
		if q.Err != nil {
			z.Error(q.Err.Error())
			q.RespErr()
			return
		}
		//对账号做解密处理
		if q.SysUser.ID.Int64 <= 20000 {
			for index, element := range insuranceSets {
				if len(element.Underwriter) != 0 && string(element.Underwriter) != "{}" && string(element.Underwriter) != "[{}]" {
					insuranceSets[index].Underwriter, q.Err = decodeUnderwriter(element.Underwriter)
					if q.Err != nil {
						q.RespErr()
						return
					}
				}
			}
		}
		jsonText, q.Err = json.Marshal(insuranceSets)
		if q.Err != nil {
			z.Error(q.Err.Error())
			q.RespErr()
			return
		}
		q.Msg.RowCount = s.RowCount
		q.Msg.Data = jsonText
		q.Resp()
		return

	case "post":
		if !q.IsAdmin {
			q.Err = fmt.Errorf("非管理员,不可修改险种信息")
			z.Error(q.Err.Error())
			q.RespErr()
			return
		}
		if plan := q.R.URL.Query().Get("addPlan"); plan == "true" {
			addPlan(ctx)
			return
		}

		typeID := q.R.URL.Query().Get("typeID")
		var insuranceTypeID int64
		insuranceTypeID, q.Err = strconv.ParseInt(typeID, 10, 64)
		if q.Err != nil {
			z.Error(q.Err.Error())
			q.RespErr()
			return
		}
		//--get pdf file
		fDesc := &fileOwnDesc{}
		var fd multipart.File
		var fileHeader *multipart.FileHeader
		fd, fileHeader, q.Err = q.R.FormFile("listTpl")
		if q.Err == http.ErrMissingFile {
			q.Err = fmt.Errorf("没有上传'listTpl'内容")
			z.Error(q.Err.Error())
			q.RespErr()
			return
		}
		if q.Err != nil {
			z.Error(q.Err.Error())
			q.RespErr()
			return
		}
		defer fd.Close()
		if fileHeader.Filename == "" {
			q.Err = fmt.Errorf("文件名不能为空")
			z.Error(q.Err.Error())
			q.RespErr()
			return
		}
		var buff []byte
		buff, q.Err = ioutil.ReadAll(fd)
		if q.Err != nil {
			z.Error(q.Err.Error())
			q.RespErr()
			return
		}
		z.Info("上传了文件" + fileHeader.Filename)
		fDesc.SN = null.IntFrom(0)
		fDesc.Name = null.StringFrom(fileHeader.Filename)
		fDesc.Label = null.StringFrom(labelListTpl)
		fDesc.OwnerType = null.StringFrom("insuranceTypes")
		fDesc.LinkID = null.IntFrom(insuranceTypeID)
		fDesc.Item = null.StringFrom("files")
		fDesc.reservedFile = true
		q.Err = saveFileBytes(buff, fDesc)
		if q.Err != nil {
			z.Error(q.Err.Error())
			q.RespErr()
			return
		}
		var fdExists, pathExists bool
		fdExists, pathExists, q.Err = fileDescUpdate(ctx, fDesc)
		if q.Err != nil {
			q.RespErr()
			return
		}

		if fdExists && pathExists {
			q.Err = getIDIfRepeatUpload(fDesc)
			if q.Err != nil {
				q.RespErr()
				return
			}
		}
		//--赋值到全局变量
		var f *excelize.File
		f, q.Err = excelize.OpenFile(fileStorePath + fDesc.MD5.String)
		if q.Err != nil {
			z.Error(q.Err.Error())
			q.RespErr()
			return
		}
		q.Err = initListMap(f, insuranceTypeID)
		if q.Err != nil {
			q.RespErr()
			return
		}
		//--删除原文件位置
		q.Err = deleteFileFromTableField(fDesc)
		if q.Err != nil && q.Err.Error() == "没有文件可以删除" {
			q.Err = nil
		}
		if q.Err != nil {
			q.RespErr()
			return
		}

		var filesField string
		filesField, q.Err = updateTableField(fDesc)
		if q.Err != nil {
			z.Error(q.Err.Error())
			q.RespErr()
			return
		}
		//---更改list_tpl
		s := "update t_insurance_types set list_tpl = $1 where id = $2"
		var stmt *sqlx.Stmt
		stmt, q.Err = sqlxDB.Preparex(s)
		if q.Err != nil {
			z.Error(q.Err.Error())
			q.RespErr()
			return
		}
		defer stmt.Close()
		var result sql.Result
		result, q.Err = stmt.Exec(fileStorePath+fDesc.MD5.String, insuranceTypeID)
		if q.Err != nil {
			z.Error(q.Err.Error())
			q.RespErr()
			return
		}

		var d int64
		if d, q.Err = result.RowsAffected(); q.Err != nil {
			z.Error(q.Err.Error())
			q.RespErr()
			return
		}
		if d <= 0 {
			q.Err = fmt.Errorf("清单路径没有修改成功")
			z.Error(q.Err.Error())
			q.RespErr()
			return
		}

		q.Msg.Msg = filesField
		q.Msg.Data = types.JSONText(fmt.Sprintf(`{"RowAffected":%d}`, 1))
		q.Resp()

	case "put":

		if !q.IsAdmin {
			q.Err = fmt.Errorf("非管理员,不可修改险种信息")
			z.Error(q.Err.Error())
			q.RespErr()
			return
		}
		var buf []byte
		buf, q.Err = ioutil.ReadAll(q.R.Body)
		if q.Err != nil {
			z.Error(q.Err.Error())
			q.RespErr()
			return
		}
		defer q.R.Body.Close()

		if len(buf) == 0 {
			q.Err = fmt.Errorf("Call /api/order by post with empty body")
			z.Error(q.Err.Error())
			q.RespErr()
			return
		}
		z.Info(string(buf))
		var req ReqProto
		q.Err = json.Unmarshal(buf, &req)
		if q.Err != nil {
			z.Error(q.Err.Error())
			q.RespErr()
			return
		}

		if strings.ToLower(req.Action) != "update" {
			q.Err = fmt.Errorf("req.Action is not update")
			z.Error(q.Err.Error())
			q.RespErr()
			return
		}
		if len(req.Data) == 0 {
			q.Err = fmt.Errorf("不指定data，你想干啥子哦？")
			z.Error(q.Err.Error())
			q.RespErr()
			return
		}

		//req.AuthFilter = authFilter
		var i TInsuranceTypes

		q.Err = json.Unmarshal([]byte(req.Data), &i)
		if q.Err != nil {
			z.Error(q.Err.Error())
			q.RespErr()
			return
		}

		if !i.ID.Valid || i.ID.Int64 == 0 {
			q.Err = fmt.Errorf("你需要指定要修改哪个险种")
			z.Error(q.Err.Error())
			q.RespErr()
			return
		}
		req.Filter = map[string]interface{}{
			"ID": map[string]interface{}{"EQ": i.ID.Int64},
		}

		fields := insuranceSetting[i.ID.Int64]
		if fields == nil {
			q.Err = fmt.Errorf("Unknown types of insurance")
			z.Error(q.Err.Error())
			q.RespErr()
			return
		}

		q.Err = invalidUnusedNullValue(&i, fields)
		if q.Err != nil {
			q.RespErr()
			return
		}

		// s := reflect.ValueOf(&i).Elem()

		var funcToCheckSetting = map[string]func() (err error){
			"InvoiceTitleUpdateTimes": i.checkInvoiceTitleUpdateTimes,
			"ReceiptAccount":          i.checkReceiptAccount,
			"Contact":                 i.checkContact,
			"OrderRepeatLimit":        i.checkOrderRepeatLimit,
			"RemindDays":              i.checkRemindDays,
			"Underwriter":             i.checkUnderwriter,
			"Mail":                    i.checkMail,
			"ContactQrCode":           i.checkContactQrCode,
			"CheckAutoFillParam":      i.checkAutoFillParam,
			"HaveDinnerNum":           i.checkHaveDinnerNum,
			"EnableImportList":        i.checkEnableImportList,
			"Alias":                   i.checkAlias,
			"AgeLimit":                i.checkAgeLimit,
			// "SuddenDeathDescription":  i.checkSuddenDeathDescription,
		}

		for _, updateField := range fields {
			// field := s.FieldByName(updateField).Interface()
			//获取对应的检测函数
			checkFunc := funcToCheckSetting[updateField]
			if checkFunc == nil {
				q.Err = fmt.Errorf("缺少对%s字段的检验函数", updateField)
				z.Error(q.Err.Error())
				q.RespErr()
				return
			}

			// q.Err = checkFunc(field)
			q.Err = checkFunc()
			if q.Err != nil {
				q.RespErr()
				return
			}
		}

		// 对账号做加密处理
		if len(i.Underwriter) != 0 {
			i.Underwriter, q.Err = encodeUnderwriter(i.Underwriter)
			if q.Err != nil {
				q.RespErr()
				return
			}
		}

		req.Data, q.Err = MarshalJSON(i)
		if q.Err != nil {
			z.Error(q.Err.Error())
			q.RespErr()
			return
		}

		i.TableMap = &i

		q.Err = DML(&i.Filter, &req)
		if q.Err != nil {
			q.RespErr()
			return
		}

		q.Err = refreshInsuranceTypes()
		if q.Err != nil {
			q.Err = fmt.Errorf("数据库更新成功了,但后端保险参数刷新失败。请立刻联系数据库管理员,不要擅自操作!!!错误信息:%s", q.Err.Error())
			z.Error(q.Err.Error())
			q.RespErr()
			return
		}

		q.Err = refreshBaseParam("保险参数")
		if q.Err != nil {
			q.Err = fmt.Errorf("数据库更新成功了,但保险参数刷新失败。请立刻联系数据库管理员,不要擅自操作!!!错误信息:%s", q.Err.Error())
			z.Error(q.Err.Error())
			q.RespErr()
			return
		}

		q.Msg.Data = types.JSONText(fmt.Sprintf(`{"RowAffected":%d}`, 1))
		q.Resp()

		return
	}

}

//传入sql语句 获取险种集合
func insuranceTypesShow(s string) (showInfos []*TInsuranceTypes, err error) {
	if s == "" {
		err = fmt.Errorf("sql  is null")
		z.Error(err.Error())
		return
	}
	z.Info("---->" + FncName())
	var stmt *sqlx.Stmt
	stmt, err = sqlxDB.Preparex(s)
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
		var v TInsuranceTypes
		err = rows.StructScan(&v)
		if err != nil {
			z.Error(err.Error())
			return
		}
		showInfos = append(showInfos, &v)
	}
	return
}

//如传入10012(各类活动),会返回比赛/活动保险,(但是remark为活动类)
//传入10024(教工责任险),仍会返回教工责任险
func getParentInsurance(insuranceID int64) (t *TInsuranceTypes, err error) {
	if insuranceID == 0 {
		err = fmt.Errorf("insuracneID is empty")
		z.Error(err.Error())
		return
	}
	//递归查询父级险种,直到查询到第一个真实险种,返回id和备注
	s := `WITH RECURSIVE r AS (
		SELECT * FROM t_insurance_types WHERE id = $1
	  union   ALL
		SELECT t.* FROM t_insurance_types t, r WHERE t.id = r.parent_id
	)
	select r.id,r.name,r.parent_id,r.define_level,t.remark Remark
	from r,t_insurance_types t where r.define_level = 2 and t.id = $1`
	var stmt *sqlx.Stmt
	stmt, err = sqlxDB.Preparex(s)
	if err != nil {
		z.Error(err.Error())
		return
	}
	defer stmt.Close()
	var temp TInsuranceTypes
	err = stmt.QueryRowx(insuranceID).StructScan(&temp)
	if err == sql.ErrNoRows {
		err = fmt.Errorf("未知的险种ID：%d", insuranceID)
		z.Error(err.Error())
		return
	}
	if err != nil {
		z.Error(err.Error())
		return
	}

	return &temp, nil
}

var insuranceSetting = map[int64][]string{

	insuranceContest: { //比赛活动保险
		"Alias",
		"InvoiceTitleUpdateTimes", //发票抬头修改
		"ReceiptAccount",          //对公账户设置
		"Contact",                 //联系人设置
		"OrderRepeatLimit",        //订单份数限制
		"RemindDays",              //自动催款天数
		"Underwriter",             //承保公司账号
		// "ContactQrCode",           //缴费联系人二维码
		"CheckAutoFillParam",
		"AgeLimit",
	},
	insuranceSchool: { //校方
		"InvoiceTitleUpdateTimes", //发票抬头修改
		"ReceiptAccount",          //对公账户设置
		"RemindDays",              //自动催款天数
		"Mail",                    //邮寄地址
		"ContactQrCode",           //缴费联系人二维码
		"EnableImportList",
		"AgeLimit",
	},
	insuranceStaff: { //教工
		"InvoiceTitleUpdateTimes", //发票抬头修改
		"ReceiptAccount",          //对公账户设置
		"RemindDays",              //自动催款天数
		"Mail",                    //邮寄地址
		"ContactQrCode",           //缴费联系人二维码
		"AgeLimit",
	},
	insuranceTrainee: { //实习生,比校方其他系列多了议价联系人
		"InvoiceTitleUpdateTimes", //发票抬头修改
		"ReceiptAccount",          //对公账户设置
		"Contact",                 //联系人设置
		"RemindDays",              //自动催款天数
		"Mail",                    //邮寄地址
		"ContactQrCode",           //缴费联系人二维码
		"AgeLimit",
	},
	insuranceSchoolBus: { //校车
		"InvoiceTitleUpdateTimes", //发票抬头修改
		"ReceiptAccount",          //对公账户设置
		"RemindDays",              //自动催款天数
		"Mail",                    //邮寄地址
		"ContactQrCode",           //缴费联系人二维码
	},
	insuranceCanteen: { //食品卫生
		"InvoiceTitleUpdateTimes", //发票抬头修改
		"ReceiptAccount",          //对公账户设置
		"RemindDays",              //自动催款天数
		"Mail",                    //邮寄地址
		"ContactQrCode",           //缴费联系人二维码
		"HaveDinnerNum",
	},
	insuranceOrganization: {
		"RemindDays",
		"ContactQrCode",  //缴费联系人二维码
		"ReceiptAccount", //对公账户设置
		"Underwriter",    //承保公司账号
		"Contact",        //联系人设置
		"Mail",
		// "SuddenDeathDescription", //猝死责任险描述
	},
	insuranceClub: {
		"RemindDays",
		"ContactQrCode",  //缴费联系人二维码
		"ReceiptAccount", //对公账户设置
		"Underwriter",    //承保公司账号
		"Contact",        //联系人设置
		"Mail",
		// "SuddenDeathDescription", //猝死责任险描述
	},
	insuranceSwimmingPool: {
		"RemindDays",
		"ContactQrCode",  //缴费联系人二维码
		"ReceiptAccount", //对公账户设置
		"Underwriter",    //承保公司账号
		"Contact",        //联系人设置
		"Mail",
	},
	insuranceAccident: {
		"AgeLimit",
	},
}

//检验发票抬头
func (it *TInsuranceTypes) checkInvoiceTitleUpdateTimes() (err error) {
	if !it.InvoiceTitleUpdateTimes.Valid {
		return
	}
	if it.InvoiceTitleUpdateTimes.Int64 < 0 {
		err = fmt.Errorf("发票抬头修改次数小于0是什么鬼")
		z.Error(err.Error())
		return
	}
	z.Info("InvoiceTitleUpdateTimes check success")
	return
}

//检验别名
func (it *TInsuranceTypes) checkAlias() (err error) {
	if !it.Alias.Valid {
		return
	}
	matched, err := regexp.MatchString("[`~!@#$^&*()=|{}':;',\\[\\].<>/?~！@#￥……&*（）——|{}【】‘；：”“'。，、？%]", it.Alias.String)
	if err != nil {
		z.Error(err.Error())
		return
	}
	if matched {
		err = fmt.Errorf("别名存在特殊字符")
		z.Error(err.Error())
		return
	}
	z.Info("Alias check success")
	return
}

//检验对公账号设置
func (it *TInsuranceTypes) checkReceiptAccount() (err error) {
	if len(it.ReceiptAccount) == 0 {
		return
	}
	if string(it.ReceiptAccount) == "{}" || string(it.ReceiptAccount) == "[{}]" {
		err = fmt.Errorf("ReceiptAccount is empty")
		z.Error(err.Error())
		return
	}
	var receipts []ReceiptAccount
	err = json.Unmarshal(it.ReceiptAccount, &receipts)
	if err != nil {
		z.Error(err.Error())
		return
	}
	for index, receipt := range receipts {
		if receipt.Bank == "" {
			err = fmt.Errorf("[对公账号设置]第%d个对公帐号设置的开户行为空", index+1)
			z.Error(err.Error())
			return
		}
		if receipt.BankNum == "" {
			err = fmt.Errorf("[对公账号设置]第%d个对公帐号设置的行号为空", index+1)
			z.Error(err.Error())
			return
		}
		if receipt.AccountName == "" {
			err = fmt.Errorf("[对公账号设置]第%d个对公帐号设置的户名为空", index+1)
			z.Error(err.Error())
			return
		}
		if receipt.Account == "" {
			err = fmt.Errorf("[对公账号设置]第%d个对公帐号设置的开户行的账号为空", index+1)
			z.Error(err.Error())
			return
		}
	}
	z.Info("ReceiptAccount check success")
	return
}

//检验短信联系人
func (it *TInsuranceTypes) checkContact() (err error) {
	if len(it.Contact) == 0 {
		return
	}
	if string(it.Contact) == "{}" || string(it.Contact) == "[{}]" {
		err = fmt.Errorf("Contact is empty")
		z.Error(err.Error())
		return
	}

	var cts []Contact
	err = json.Unmarshal(it.Contact, &cts)
	if err != nil {
		z.Error(err.Error())
		return
	}
	for index, ct := range cts {
		if ct.ContactName == "" {
			err = fmt.Errorf("[协议价短信联系人]第%d个的联系人姓名为空", index+1)
			z.Error(err.Error())
			return
		}

		if ct.Phone == "" {
			err = fmt.Errorf("[协议价短信联系人]设置的第%d个联系电话为空?", index+1)
			z.Error(err.Error())
			return
		}

		if !verifyTelNO(ct.Phone) {
			err = fmt.Errorf("[协议价短信联系人]的第%d个联系电话校验不通过", index+1)
			z.Error(err.Error())
			return
		}
	}

	z.Info("Contact check success")
	return
}

//检验订单份数
func (it *TInsuranceTypes) checkOrderRepeatLimit() (err error) {
	if !it.OrderRepeatLimit.Valid {
		return
	}
	if it.OrderRepeatLimit.Int64 <= 0 {
		err = fmt.Errorf("[最大订单份数]The number of orders is less than or equal to 0")
		z.Error(err.Error())
		return
	}
	z.Info("OrderRepeatLimit check success")
	return
}

//检验自动催款天数
func (it *TInsuranceTypes) checkRemindDays() (err error) {
	if !it.RemindDays.Valid {
		return
	}
	if it.RemindDays.Int64 <= 0 {
		err = fmt.Errorf("[自动催款天数]小于等于0?")
		z.Error(err.Error())
		return
	}

	z.Info("RemindDays check success")
	return
}

//检验是否开启清单录入
func (it *TInsuranceTypes) checkEnableImportList() (err error) {
	if !it.EnableImportList.Valid {
		return
	}
	z.Info("EnableImportList check success")
	return
}

//检验是否开启就餐人数
func (it *TInsuranceTypes) checkHaveDinnerNum() (err error) {
	if !it.HaveDinnerNum.Valid {
		return
	}
	z.Info("HaveDinnerNum check success")
	return
}

//检验承保公司账号
func (it *TInsuranceTypes) checkUnderwriter() (err error) {
	if len(it.Underwriter) == 0 {
		return
	}
	if string(it.Underwriter) == "{}" || string(it.Underwriter) == "[{}]" {
		err = fmt.Errorf("Underwriter is empty")
		z.Error(err.Error())
		return
	}
	var underwriters []Underwriter
	err = json.Unmarshal(it.Underwriter, &underwriters)
	if err != nil {
		z.Error(err.Error())
		return
	}

	for index, uw := range underwriters {
		if uw.StartTime <= 0 || uw.EndTime <= 0 {
			err = fmt.Errorf("[承保公司账号]设置的第%d个账号的开始时间或结束时间为0", index+1)
			z.Error(err.Error())
			return
		}
		if uw.EndTime-uw.StartTime <= 0 {
			err = fmt.Errorf("[承保公司账号]第%d个账号,开始时间和结束时间的时间差为0", index+1)
			z.Error(err.Error())
			return
		}

		if uw.Company == "" {
			err = fmt.Errorf("[承保公司账号]第%d个设置的账号所属公司为空", index+1)
			z.Error(err.Error())
			return
		}
		if uw.Account == "" {
			err = fmt.Errorf("[承保公司账号]第%d个设置的账号为空", index+1)
			z.Error(err.Error())
			return
		}

		if uw.Password == "" {
			err = fmt.Errorf("[承保公司账号]第%d个设置的账号密码为空", index+1)
			z.Error(err.Error())
			return
		}

		if v := index + 1; v < len(underwriters) {
			if uw.EndTime >= underwriters[v].StartTime {
				err = fmt.Errorf("[承保公司账号]账号%s与账号%s时间重叠", uw.Account, underwriters[v].Account)
				z.Error(err.Error())
				return
			}
		}
	}
	z.Info("Underwriter check success")
	return
}

func encodeUnderwriter(s types.JSONText) (encodeContent types.JSONText, err error) {

	if len(s) == 0 || string(s) == "{}" || string(s) == "[{}]" {
		err = fmt.Errorf("jsontext is empty")
		z.Error(err.Error())
		return
	}
	var underwriters []Underwriter
	err = json.Unmarshal(s, &underwriters)
	if err != nil {
		z.Error(err.Error())
		return
	}

	for index, underwriter := range underwriters {
		var encodePWD []byte

		//AES加密
		encodePWD, err = encryptAES([]byte(underwriter.Password), aesKey)
		if err != nil {
			err = fmt.Errorf("[承保公司账号]第%d个账号密码加密错误", index+1)
			z.Error(err.Error())
			return
		}

		//hex编码
		underwriters[index].Password = hex.EncodeToString(encodePWD)
		z.Info("密文(hex):" + underwriters[index].Password)
	}
	encodeContent, err = json.Marshal(underwriters)
	if err != nil {
		z.Error(err.Error())
		return
	}
	z.Info("Underwriter's password encoded success")
	return
}

func decodeUnderwriter(s types.JSONText) (decodeContent types.JSONText, err error) {
	if len(s) == 0 || string(s) == "{}" || string(s) == "[{}]" {
		err = fmt.Errorf("jsontext is empty")
		z.Error(err.Error())
		return
	}
	z.Info(string(s))

	var underwriters []Underwriter
	err = json.Unmarshal(s, &underwriters)
	if err != nil {
		z.Error(err.Error())
		return
	}

	for index, underwriter := range underwriters {
		var encodePWD, decodePWD []byte

		if underwriter.Password == "" {
			continue
		}

		//hex解码
		encodePWD, err = hex.DecodeString(underwriter.Password)
		if err != nil {
			z.Error(err.Error())
			return
		}

		//AES解密
		decodePWD, err = decryptAES(encodePWD, aesKey)
		if err != nil {
			err = fmt.Errorf("[承保公司账号]第%d个账号密码解密错误", index+1)
			z.Error(err.Error())
			return
		}
		underwriters[index].Password = string(decodePWD)
	}
	decodeContent, err = json.Marshal(underwriters)
	if err != nil {
		z.Error(err.Error())
		return
	}
	z.Info("Underwriter's password decoded success")
	return
}

//检验邮寄地址
func (it *TInsuranceTypes) checkMail() (err error) {
	if len(it.Mail) == 0 {
		return
	}
	if string(it.Mail) == "{}" || string(it.Mail) == "[{}]" {
		err = fmt.Errorf("Mail is empty")
		z.Error(err.Error())
		return
	}
	var mails []Mail
	err = json.Unmarshal(it.Mail, &mails)
	if err != nil {
		z.Error(err.Error())
		return
	}

	for index, mail := range mails {
		if mail.Receiver == "" {
			err = fmt.Errorf("[邮寄信息]第%d位收件人姓名为空", index+1)
			z.Error(err.Error())
			return
		}

		if mail.Phone == "" {
			err = fmt.Errorf("[邮寄信息]第%d位收件人电话为空", index+1)
			z.Error(err.Error())
			return
		}

		if !verifyTelNO(mail.Phone) {
			err = fmt.Errorf("[邮寄信息]第%d个设置的联系电话校验不通过", index+1)
			z.Error(err.Error())
			return
		}

		if mail.Address == "" {
			err = fmt.Errorf("[邮寄信息]第%d个邮寄地址为空", index+1)
			z.Error(err.Error())
			return
		}
	}

	z.Info("Underwriter check success")
	return
}

//检验年龄限制
func (it *TInsuranceTypes) checkAgeLimit() (err error) {
	if len(it.AgeLimit) == 0 {
		return
	}
	if string(it.AgeLimit) == "{}" || string(it.AgeLimit) == "[{}]" {
		err = fmt.Errorf("AgeLimit is empty")
		z.Error(err.Error())
		return
	}
	var ageLimit map[string]int
	err = json.Unmarshal(it.AgeLimit, &ageLimit)
	if err != nil {
		z.Error(err.Error())
		return
	}
	maleMax, ok := ageLimit["MaleMax"]
	if !ok {
		err = fmt.Errorf("缺少MaleMax字段")
		z.Error(err.Error())
		return
	}
	maleMin, ok := ageLimit["MaleMin"]
	if !ok {
		err = fmt.Errorf("缺少MaleMax字段")
		z.Error(err.Error())
		return
	}
	femaleMax, ok := ageLimit["FemaleMax"]
	if !ok {
		err = fmt.Errorf("缺少MaleMax字段")
		z.Error(err.Error())
		return
	}
	femaleMin, ok := ageLimit["FemaleMin"]
	if !ok {
		err = fmt.Errorf("缺少MaleMax字段")
		z.Error(err.Error())
		return
	}
	if maleMax < maleMin || femaleMax < femaleMin {
		err = fmt.Errorf("最高年龄限制不能小于最低年龄限制")
		z.Error(err.Error())
		return
	}

	z.Info("AgeLimit check success")
	return
}

//检验缴费联系人二维码
func (it *TInsuranceTypes) checkContactQrCode() (err error) {
	if !it.ContactQrCode.Valid {
		return
	}
	if it.ContactQrCode.String == "" {
		err = fmt.Errorf("ContactQrCode is empty")
		z.Error(err.Error())
		return
	}

	// TODO:对base64的检验
	z.Info("ContactQrCode check success")
	return
}

//检验自动化参数
func (it *TInsuranceTypes) checkAutoFillParam() (err error) {
	if len(it.AutoFillParam) == 0 {
		return
	}

	if string(it.AutoFillParam) == "{}" || string(it.AutoFillParam) == "[{}]" {
		err = fmt.Errorf("AutoFillParam is empty")
		z.Error(err.Error())
		return
	}
	var autoFillParam []AutoFillParam
	err = json.Unmarshal(it.AutoFillParam, &autoFillParam)
	if err != nil {
		z.Error(err.Error())
		return
	}

	for index, param := range autoFillParam {
		if param.Phone == "" {
			err = fmt.Errorf("[自动化参数]第%d个设置的手机号为空", index+1)
			z.Error(err.Error())
			return
		}

		if param.Reviewer == "" {
			err = fmt.Errorf("[自动化参数]第%d个设置的复核人姓名为空", index+1)
			z.Error(err.Error())
			return
		}

		if !verifyTelNO(param.Phone) {
			err = fmt.Errorf("[自动化参数]第%d个设置的联系电话校验不通过", index+1)
			z.Error(err.Error())
			return
		}
	}

	z.Info("AutoFillParam check success")
	return
}

// func (it *TInsuranceTypes) checkSuddenDeathDescription() (err error) {
// 	if len(it.SuddenDeathDescription) == 0 || string(it.SuddenDeathDescription) == "{}" {
// 		err = fmt.Errorf("猝死责任险描述为空")
// 		z.Error(err.Error())
// 		return
// 	}
// 	type sudden struct {
// 		true  string
// 		false string
// 	}
// 	var s sudden
// 	err = json.Unmarshal(it.SuddenDeathDescription, &s)
// 	if err != nil {
// 		z.Error(err.Error())
// 	}
// 	return
// }

//getCurrentUnderwriter:获取当前时间设置的承保公司账号,Underwriter为nil则未设置
func getCurrentUnderwriter(insuranceTypeID int64) (uw *Underwriter, err error) {
	if insuranceTypeID == 0 {
		err = fmt.Errorf("insuranceTypeID is 0")
		z.Error(err.Error())
		return
	}
	s := `with t as ( 
		select jsonb_array_elements(underwriter) underwriter,
	   to_timestamp((jsonb_array_elements(underwriter)->'StartTime')::bigint/1000) begin_time,
	   to_timestamp((jsonb_array_elements(underwriter)->'EndTime')::bigint/1000) end_time
	   from t_insurance_types 
	   where id = $1)
	   select underwriter from t where begin_time < now() and end_time > now()`

	var stmt *sqlx.Stmt
	stmt, err = sqlxDB.Preparex(s)
	if err != nil {
		z.Error(err.Error())
		return
	}
	defer stmt.Close()
	var jsonText []byte
	err = stmt.QueryRowx(insuranceTypeID).Scan(&jsonText)
	if err == sql.ErrNoRows {
		err = fmt.Errorf("当前时间未设置承保公司账号，请联系管理员")
		z.Error(err.Error())
		return
	}
	if err != nil {
		z.Error(err.Error())
		return
	}

	uw = &Underwriter{}
	err = json.Unmarshal(jsonText, &uw)
	if err != nil {
		z.Error(err.Error())
		return
	}
	if uw.Password != "" {
		var encodePWD, decodePWD []byte
		//hex解码
		encodePWD, err = hex.DecodeString(uw.Password)
		if err != nil {
			z.Error(err.Error())
			return
		}
		//AES解密
		decodePWD, err = decryptAES(encodePWD, aesKey)
		if err != nil {
			err = fmt.Errorf("承保公司账号账号密码解密错误")
			z.Error(err.Error())
			return
		}
		uw.Password = string(decodePWD)
	}

	return
}

//判断是否是比赛活动保险
func isInsuranceContest(id int64) bool {

	return id == insuranceContest2 || id == insuranceContest4 || id == insuranceContest6 || id == insuranceContest8 ||
		id == insuranceContest10 || id == insuranceContest14 || id == insuranceContest12
}

//刷新险种数据
func (i *insurancesData) refresh() (err error) {
	z.Info("---->" + FncName())
	i.lock.Lock()
	defer i.lock.Unlock()
	s := `select id, name, alias, data_type, parent_id, age_limit, rule_batch, org_id,
	pay_type, pay_channel, pay_name, bank_account, bank_account_name, bank_name, bank_id, 
	floor_price, price, price_config, define_level, layout_order, layout_level, list_tpl,
	files, pic, sudden_death_description, description, auto_fill, enable_import_list,
	have_dinner_num, invoice_title_update_times, receipt_account, contact, contact_qr_code, 
	underwriter, remind_days, mail, order_repeat_limit, group_by_max_day, web_description, 
	mobile_description, auto_fill_param, "interval", max_insure_in_year, insured_in_month, 
	insured_start_time, insured_end_time, allow_start, allow_end, indate_start, indate_end, addi, 
	domain_id, creator, create_time, updator, update_time, remark, status
	from t_insurance_types 
	where data_type <> '4' and data_type <> '6' and data_type <> '8'
	order by layout_order asc,id asc`

	var stmt *sqlx.Stmt
	stmt, err = sqlxDB.Preparex(s)
	if err != nil {
		z.Error(err.Error())
		return
	}
	defer stmt.Close()

	rows, err := stmt.Queryx()
	if err == sql.ErrNoRows {
		err = fmt.Errorf("无险种配置信息,请联系数据库管理员")
		z.Error(err.Error())
		return
	}
	if err != nil {
		z.Error(err.Error())
		return
	}
	defer rows.Close()
	var insurancesTemp = make([]*TInsuranceTypes, 0)
	var insurancesJSONTemp []byte
	var insurancesMapTemp = make(map[int64]*TInsuranceTypes)

	for rows.Next() {
		var insurance TInsuranceTypes
		err = rows.StructScan(&insurance)
		if err != nil {
			z.Error(err.Error())
			return
		}
		insurancesTemp = append(insurancesTemp, &insurance)
	}
	insurancesJSONTemp, err = json.Marshal(&insurancesTemp)
	if err != nil {
		z.Error(err.Error())
		return
	}

	var a []interface{}
	err = json.Unmarshal(insurancesJSONTemp, &a)
	if err != nil {
		z.Error(err.Error())
		return
	}

	a = removeEmptyArrayElement(a)
	insurancesJSONTemp, err = json.Marshal(&a)
	if err != nil {
		z.Error(err.Error())
		return
	}

	insurancesMapTemp = make(map[int64]*TInsuranceTypes)
	for _, v := range insurancesTemp {
		insurancesMapTemp[v.ID.Int64] = v
	}

	i.list = insurancesTemp
	i.jsonData = insurancesJSONTemp
	i.mapData = insurancesMapTemp
	return
}

func (i *insurancesData) getByID(id int64) *TInsuranceTypes {
	i.lock.RLock()
	defer i.lock.RUnlock()
	// 值拷贝
	t := *i.mapData[id]
	return &t
}

func (i *insurancesData) getList() []*TInsuranceTypes {
	return i.list
}

func (i *insurancesData) getMap() map[int64]*TInsuranceTypes {
	return i.mapData
}

func (i *insurancesData) getJSON() []byte {
	return i.jsonData
}

const defaultDoc = `
{
	"投保须知":	"",
	"保险条款":	"",
	"免责和退保声明":	"",
	"特别约定":	"",
	"time_set":{
		"投保须知":0,
		"保险条款":0,
		"免责和退保声明":0,
		"特别约定":0
	}
}
`

func (i *insurancesData) copyMap() map[int64]*TInsuranceTypes {
	i.lock.RLock()
	defer i.lock.RUnlock()
	mapData := make(map[int64]*TInsuranceTypes)
	for k, v := range i.mapData {
		mapData[k] = v
	}
	return mapData
}

func addPlan(ctx context.Context) {
	q := GetCtxValue(ctx)
	z.Info("---->" + FncName())
	q.Stop = true
	if !q.IsAdmin {
		q.Err = fmt.Errorf("非管理员,不可添加方案信息")
		z.Error(q.Err.Error())
		q.RespErr()
		return
	}

	var buf []byte

	str := q.R.FormValue("q")
	buf = []byte(str)
	if len(buf) == 0 {
		q.Err = fmt.Errorf("call api by post with empty body")
		z.Error(q.Err.Error())
		q.RespErr()
		return
	}
	var req ReqProto
	q.Err = json.Unmarshal(buf, &req)
	if q.Err != nil {
		z.Error(q.Err.Error())
		q.RespErr()
		return
	}

	//学意险多个机构的多个规则同时修改
	if urlParam := q.R.URL.Query().Get("accidentRule"); urlParam == "true" {
		var d []int64
		d, q.Err = accidentRule(&req)
		if q.Err != nil {
			q.RespErr()
			return
		}
		q.Msg.Data, q.Err = json.Marshal(d)
		if q.Err != nil {
			z.Error(q.Err.Error())
			q.RespErr()
			return
		}

		q.Resp()
		return
	}
	var pro TInsuranceTypesPro
	var i *TInsuranceTypes

	if gjson.Get(string(req.Data), "ContactQrCode").IsObject() {
		q.Err = fmt.Errorf("ContactQrCode只能是string不能是json对象")
		z.Error(q.Err.Error())
		q.RespErr()
		return
	}

	q.Err = json.Unmarshal(req.Data, &pro)
	if q.Err != nil {
		z.Error(q.Err.Error())
		q.RespErr()
		return
	}

	i = &pro.TInsuranceTypes

	if len(pro.Batch) == 0 && i.OrgID.Valid {
		pro.Batch = []int64{i.OrgID.Int64}
	}
	if len(pro.Batch) > 0 {
		orgSet := "("
		for _, o := range pro.Batch {
			orgSet += strconv.FormatInt(o, 10) + "),("
		}
		orgSet = orgSet[:len(orgSet)-2]
		z.Info(orgSet)
		s := fmt.Sprintf(`
		with org_set(org_id) as (values %s)
		SELECT org_id, s.name AS org_name, (
			SELECT id
			FROM v_insurance_type v 
			WHERE v.org_id = set.org_id and ref_id = $1
			) AS plan_id
		FROM org_set AS set 
		LEFT JOIN t_school s on org_id = s.id
		`, orgSet)
		z.Info(fmt.Sprintf(strings.ReplaceAll(s, "$1", "%d"), i.RefID.Int64))
		var stmt *sqlx.Stmt
		stmt, q.Err = sqlxDB.Preparex(s)
		if q.Err != nil {
			z.Error(q.Err.Error())
			q.RespErr()
			return
		}
		defer stmt.Close()
		var rows *sqlx.Rows
		rows, q.Err = stmt.Queryx(i.RefID.Int64)
		if q.Err != nil {
			z.Error(q.Err.Error())
			q.RespErr()
			return
		}
		defer rows.Close()
		var batchCheck []ruleAddCheck
		for rows.Next() {
			var r ruleAddCheck
			q.Err = rows.StructScan(&r)
			if q.Err != nil {
				z.Error(q.Err.Error())
				q.RespErr()
				return
			}
			batchCheck = append(batchCheck, r)
		}
		for _, v := range batchCheck {
			//check org_id exists
			if v.OrgID.Int64 != 0 && (!v.OrgName.Valid || v.OrgName.String == "") {
				q.Err = fmt.Errorf("机构 %d 不存在", v.OrgID.Int64)
				z.Error(q.Err.Error())
				q.RespErr()
				return
			}
			if v.OrgID.Int64 == pro.Batch[0] && strings.ToLower(req.Action) == "insert" && v.PlanID.Valid {
				//
				req.Action = "update"
				i.ID = v.PlanID
				z.Info(fmt.Sprintf("校验第一个机构重复,改为update操作-->> plan_id:%d", v.PlanID.Int64))
			}
		}
		if i.OrgID.Valid && i.OrgID.Int64 != pro.Batch[0] {
			q.Err = fmt.Errorf("批量设置时请将data.OrgID指定为Bacth[0]")
			z.Error(q.Err.Error())
			q.RespErr()
			return
		} else if !i.OrgID.Valid {
			i.OrgID = null.IntFrom(pro.Batch[0])
		}
	}
	z.Info(string(req.Data))

	var id int64
	if strings.ToLower(req.Action) == "update" {
		if !i.ID.Valid || i.ID.Int64 == 0 {
			q.Err = fmt.Errorf("please specify ID when update")
			z.Error(q.Err.Error())
			q.RespErr()
			return
		}
		id = i.ID.Int64
		if req.Filter == nil {
			req.Filter = map[string]interface{}{
				"ID": map[string]interface{}{"EQ": i.ID.Int64},
			}
		}
		i, q.Err = GetTInsuranceTypesByPk(sqlxDB, i.ID)
		if q.Err != nil {
			z.Error(q.Err.Error())
			q.RespErr()
			return
		}

	}

	if !i.ParentID.Valid {
		q.Err = fmt.Errorf("please specify ParentID")
		z.Error(q.Err.Error())
		q.RespErr()
		return
	}

	if !i.DataType.Valid || i.DataType.String == "" || i.DataType.String == "0" {
		q.Err = fmt.Errorf("please specify DataType")
		z.Error(q.Err.Error())
		q.RespErr()
		return
	}

	//DataType
	//4: 学意险方案/学意险区域投保规则/校方区域投保规则
	//6: 校方保险方案
	//8: 校方默认投保规则
	switch {
	//方案
	case i.DataType.String == "6" || (i.DataType.String == "4" && i.ParentID.Int64 == 10040):
		if !i.PayChannel.Valid || i.PayChannel.String == "" {
			q.Err = fmt.Errorf("please specify PayChannel")
			z.Error(q.Err.Error())
			q.RespErr()
			return
		}

		if !i.PayType.Valid {
			q.Err = fmt.Errorf("please specify PayType")
			z.Error(q.Err.Error())
			q.RespErr()
			return
		}

		switch i.PayType.String {
		case "在线支付":
			i.PayName = null.StringFrom("微信支付")
		case "线下支付":
			if !i.PayName.Valid || i.PayName.String == "" {
				q.Err = fmt.Errorf("please specify PayName")
				z.Error(q.Err.Error())
				q.RespErr()
				return
			}
			//if !i.ContactQrCode.Valid || i.ContactQrCode.String == "" {
			//	q.Err = fmt.Errorf("please specify ContactQrCode")
			//	z.Error(q.Err.Error())
			//	q.RespErr()
			//	return
			//}

		case "私对公转账":

			fallthrough
		case "公对公转账":
			if !i.PayName.Valid || i.PayName.String == "" {
				q.Err = fmt.Errorf("please specify PayName")
				z.Error(q.Err.Error())
				q.RespErr()
				return
			}
			if len(i.ReceiptAccount) < 2 {
				q.Err = fmt.Errorf("please specify ReceiptAccount")
				z.Error(q.Err.Error())
				q.RespErr()
				return
			}
			var receipts []ReceiptAccount
			q.Err = json.Unmarshal(i.ReceiptAccount, &receipts)
			if q.Err != nil {
				z.Error(q.Err.Error())
				q.RespErr()
				return
			}
			for _, receipt := range receipts {

				if receipt.Bank == "" {
					q.Err = fmt.Errorf("[对公账号设置]对公帐号设置的开户行为空")
					z.Error(q.Err.Error())
					q.RespErr()
					return
				}
				if receipt.BankNum == "" {
					q.Err = fmt.Errorf("[对公账号设置]对公帐号设置的行号为空")
					z.Error(q.Err.Error())
					q.RespErr()
					return
				}
				if receipt.AccountName == "" {
					q.Err = fmt.Errorf("[对公账号设置]对公帐号设置的户名为空")
					z.Error(q.Err.Error())
					q.RespErr()
					return
				}
				if receipt.Account == "" {
					q.Err = fmt.Errorf("[对公账号设置]对公帐号设置的开户行的账号为空")
					z.Error(q.Err.Error())
					q.RespErr()
					return
				}
			}
		default:
			q.Err = fmt.Errorf("未知 PayType :%s", i.PayType.String)
			z.Error(q.Err.Error())
			q.RespErr()
			return
		}

		switch i.ParentID.Int64 {
		case 10040:

			if !i.InsuredInMonth.Valid || i.InsuredInMonth.Int64 == 0 {
				q.Err = fmt.Errorf("please specify InsuredInMonth")
				z.Error(q.Err.Error())
				q.RespErr()
				return
			}
			if !i.Price.Valid || i.Price.Float64 == 0 {
				q.Err = fmt.Errorf("please specify Price")
				z.Error(q.Err.Error())
				q.RespErr()
				return
			}

			if i.AllowStart.Valid != i.AllowEnd.Valid {
				q.Err = fmt.Errorf("please specify AllowStart/AllowEnd at the same time")
				z.Error(q.Err.Error())
				q.RespErr()
				return
			}

			i.IndateStart = null.IntFrom(GetNowInMS())
		case 10022, 10024, 10026, 10028, 10030:
			if i.InsuredStartTime.Valid || i.InsuredEndTime.Valid {
				q.Err = fmt.Errorf("don't specify InsuredStartTime/InsuredEndTime when setting price-plan")
				z.Error(q.Err.Error())
				q.RespErr()
				return
			}
		default:
			q.Err = fmt.Errorf("未知的险种id： %d", i.ParentID.Int64)
			z.Error(q.Err.Error())
			q.RespErr()
			return
		}

		if !i.Insurer.Valid || i.Insurer.String == "" {
			q.Err = fmt.Errorf("please specify Insurer")
			z.Error(q.Err.Error())
			q.RespErr()
			return
		}
		if strings.ToLower(req.Action) == "insert" {
			i.Resource = []byte(defaultDoc)
		}

		//校方规则(学意险投保规则通过accidentRule()函数判断,所以此处只有校方系列的投保规则)
	case i.DataType.String == "8" || (i.DataType.String == "4" && i.ParentID.Int64 != 10040):
		rule := TVInsuranceType{}
		if strings.ToLower(req.Action) == "insert" {
			var base *TInsuranceTypes
			base, q.Err = GetTInsuranceTypesByPk(sqlxDB, i.RefID)
			if q.Err != nil {
				z.Error(q.Err.Error())
				q.RespErr()
				return
			}
			if base == nil {
				q.Err = fmt.Errorf("ref_id:%d 不存在", i.RefID.Int64)
				z.Error(q.Err.Error())
				q.RespErr()
				return
			}
			if base.Status.String != "0" {
				q.Err = fmt.Errorf("保险方案%d为禁用状态", i.RefID.Int64)
				z.Error(q.Err.Error())
				q.RespErr()
				return
			}
			if base.RefID.Valid {
				q.Err = fmt.Errorf("ref_id:%d是投保规则,不是方案", i.RefID.Int64)
				z.Error(q.Err.Error())
				q.RespErr()
				return
			}

			buf, q.Err = json.Marshal(base)
			if q.Err != nil {
				z.Error(q.Err.Error())
				q.RespErr()
				return
			}

			q.Err = json.Unmarshal(buf, &rule)
			if q.Err != nil {
				z.Error(q.Err.Error())
				q.RespErr()
				return
			}
			q.Err = json.Unmarshal(req.Data, &rule)
			if q.Err != nil {
				z.Error(q.Err.Error())
				q.RespErr()
				return
			}
		} else {
			q.Err = sqlxDB.QueryRowx(insuranceTypeQuery+` where id = $1`, i.ID.Int64).StructScan(&rule)
			if q.Err != nil {
				z.Error(q.Err.Error())
				q.RespErr()
				return
			}
		}
		//开始时间/结束时间
		if !i.AllowStart.Valid || !i.AllowEnd.Valid {
			q.Err = fmt.Errorf("please specify AllowStart/AllowEnd when setting region-rules")
			z.Error(q.Err.Error())
			q.RespErr()
			return
		}
		//存在org_id时,检查
		if i.OrgID.Valid && i.OrgID.Int64 != 0 {
			var org *TSchool
			org, q.Err = GetTSchoolByPk(sqlxDB, i.OrgID)
			if q.Err != nil {
				z.Error(q.Err.Error())
				q.RespErr()
				return
			}
			if org == nil {
				q.Err = fmt.Errorf("org_id:%d 不存在", i.OrgID.Int64)
				z.Error(q.Err.Error())
				q.RespErr()
				return
			}
			// if org.Status.String != "0" {
			// 	q.Err = fmt.Errorf("机构 %d 处于黑名单", i.OrgID.Int64)
			// 	z.Error(q.Err.Error())
			// 	q.RespErr()
			// 	return
			// }
		}
		if i.ParentID.Int64 == 10040 {

			//投保年限
			if !i.MaxInsureInYear.Valid || i.MaxInsureInYear.Int64 == 0 {
				q.Err = fmt.Errorf("please specify MaxInsureInYear")
				z.Error(q.Err.Error())
				q.RespErr()
				return
			}

			//具体规则校验
			var valid bool

			valid, q.Err = validateInsurePlan(&rule)
			if q.Err != nil {
				z.Error(q.Err.Error())
				q.RespErr()
				return
			}
			if !valid {
				q.Err = fmt.Errorf("校验失败")
				z.Error(q.Err.Error())
				q.RespErr()
				return
			}
		}
	}

	if i.Name == "" && i.RefID.Valid {
		i.Name = "区域投保规则"
	} else if i.Name == "" {
		i.Name = "未命名方案"
	}

	if !i.DataType.Valid || i.DataType.String == "0" || i.DataType.String == "" {
		i.DataType = null.StringFrom("4")
	}

	i.DefineLevel = null.IntFrom(0)
	i.LayoutLevel = null.IntFrom(0)

	i.TableMap = i

	if !i.Status.Valid {
		i.Status = null.StringFrom("0")
	} else if !inSlice(i.Status.String, []string{"0", "4"}) {
		q.Err = fmt.Errorf("status超出值域{'0','4'}")
		z.Error(q.Err.Error())
		q.RespErr()
		return
	}
	q.Err = DML(&i.Filter, &req)

	if q.Err != nil {
		if strings.Contains(q.Err.Error(), `duplicate key value violates unique constraint "idx_policy_plan"`) {
			errMsg := ""
			if i.OrgID.Int64 == 0 {
				errMsg = "已存在同方案的默认投保规则,请直接对原已有数据进行修改"
			} else {
				errMsg = "已存在同机构同方案的投保规则,请直接对原已有数据进行修改"
			}
			q.Err = fmt.Errorf(errMsg)
			q.Msg.Status = -7000
			q.RespErr()
			return
		}

		q.RespErr()
		return
	}
	_, ok := i.QryResult.(int64)
	if !ok {
		q.Err = fmt.Errorf("_, ok = i.filter.qryResult.(int64) should be ok while it's not")
		z.Error(q.Err.Error())
		q.RespErr()
		return
	}

	if strings.ToLower(req.Action) == "insert" {
		id, ok = i.QryResult.(int64)
		if !ok {
			q.Err = fmt.Errorf("_, ok = i.qryResult.(int64) should be ok while it's not")
			z.Error(q.Err.Error())
			q.RespErr()
			return
		}
		i.ID = null.IntFrom(id)
	}

	q.Resp()
	return
}

//从excel提取列所在位置，初始化到map，同时修改excel 去除定位字符串,初始化 清单模板有关的全局map
//每次启动后端/更新模板都会初始化新的map及文件路径到全局变量
func initListMap(f *excelize.File, insuranceType int64) (err error) {

	if f == nil {
		err = fmt.Errorf("nil arguments")
		z.Error(err.Error())
		return
	}

	modifyTplPath := fileStorePath + fmt.Sprintf("list_tpl_%d_%d", insuranceType, time.Now().UnixNano())

	//保存修改后的快照文件
	err = f.SaveAs(modifyTplPath)
	if err != nil {
		z.Error(err.Error())
	}

	return
}

// * if null.String/Int/Float/Bool/Time if not in useField then invalid it
// * p 是一个结构体的指针 useField为结构体p的字段名称，
// * 不包含在useField的null.String/Int/Float/Bool/Time类型会的Valid字段会设置为false

func invalidUnusedNullValue(p interface{}, useField []string) (err error) {

	if p == nil {
		err = fmt.Errorf("call invalidUnusedField with v==nil")
		z.Error(err.Error())
		return
	}

	if reflect.TypeOf(p).Kind() != reflect.Ptr {
		err = fmt.Errorf("please call invalidUnusedField with pointer to struct")
		z.Error(err.Error())
		return
	}

	pType := reflect.TypeOf(p).Elem()
	pValue := reflect.ValueOf(p).Elem()
	t := reflect.TypeOf(pValue.Interface()) //实际结构体类型
	if t.Kind() != reflect.Struct {
		err = fmt.Errorf("Call invalidUnusedField' pointer must to be struct ")
		z.Error(err.Error())
		return
	}

	//遍历结构体的字段，有使用的字段跳过，没有使用的值为无效或者0值
	for i := 0; i < pType.NumField(); i++ {
		f := pType.Field(i)

		isUsed := false
		for _, v := range useField {
			if f.Name == v {
				isUsed = true
				break
			}
		}
		if isUsed {
			continue
		}

		keyValue := pValue.Field(i)
		//跳过非导出类型
		if !keyValue.CanInterface() {
			continue
		}

		switch keyValue.Interface().(type) {
		case null.String:
			if keyValue.IsValid() && keyValue.CanSet() &&
				keyValue.FieldByName("Valid").IsValid() &&
				keyValue.FieldByName("Valid").CanSet() {
				keyValue.FieldByName("Valid").SetBool(false)
			}
		case null.Int:
			keyValue.FieldByName("Valid").SetBool(false)
			if keyValue.IsValid() && keyValue.CanSet() &&
				keyValue.FieldByName("Valid").IsValid() &&
				keyValue.FieldByName("Valid").CanSet() {
				keyValue.FieldByName("Valid").SetBool(false)
			}
		case null.Float:
			if keyValue.IsValid() && keyValue.CanSet() &&
				keyValue.FieldByName("Valid").IsValid() &&
				keyValue.FieldByName("Valid").CanSet() {
				keyValue.FieldByName("Valid").SetBool(false)
			}
		case null.Bool:
			if keyValue.IsValid() && keyValue.CanSet() &&
				keyValue.FieldByName("Valid").IsValid() &&
				keyValue.FieldByName("Valid").CanSet() {
				keyValue.FieldByName("Valid").SetBool(false)
			}
		case null.Time:
			if keyValue.IsValid() && keyValue.CanSet() &&
				keyValue.FieldByName("Valid").IsValid() &&
				keyValue.FieldByName("Valid").CanSet() {
				keyValue.FieldByName("Valid").SetBool(false)
			}
		default:
			//z.Info("skip")
		}
	}
	return
}

var aesKey = []byte{0x73, 0x68, 0x61, 0x62, 0x69, 0x7a, 0x68, 0x61, 0x6e, 0x67, 0x7a, 0x69, 0x68, 0x61, 0x6e, 0x2e}

//PKCS5Padding PKCS5填充数据
func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

//PKCS5UnPadding PKCS5去掉填充数据，读取真实数据
func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	if length < unpadding {
		return []byte("unpadding error")
	}
	return origData[:(length - unpadding)]
}

//AES加密
func encryptAES(origData []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	origData = PKCS5Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

//AES解密
func decryptAES(crypted []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	return origData, nil
}

//TInsuranceTypesPro 处理批量设置
type TInsuranceTypesPro struct {
	Batch     []int64           //机构批量处理
	BatchRule []TInsuranceTypes //学意险批量处理
	TInsuranceTypes
}

type ruleAddCheck struct {
	OrgID   null.Int    `db:"org_id"`
	OrgName null.String `db:"org_name"`
	PlanID  null.Int    `db:"plan_id"`
}

func accidentRule(req *ReqProto) (idSet []int64, err error) {

	if req == nil {
		return
	}
	var reqData map[string]interface{}
	idSet = make([]int64, 0)
	err = json.Unmarshal(req.Data, &reqData)
	if err != nil {
		z.Error(err.Error())
		return
	}

	tx, err := sqlxDB.Beginx()
	if err != nil {
		z.Error(err.Error())
		return
	}
	defer tx.Rollback()

	// {
	// 	"action":"insert",
	// 	"data":{
	// 		"BatchRule":[
	// 			{
	// 				"RefID":20011,
	// 				"MaxInsureInYear":3,
	// 				"Status":"0",
	// 			},
	// 			{
	// 				"RefID":20011,
	// 				"MaxInsureInYear":3,
	// 				"InsuredStartTime":{{timebegin}},
	// 				"Status":"4",
	// 			},
	// 		 ],
	// 		"Batch":[1012,1000],
	// 		"AllowStart": {{timebegin}},
	// 		"AllowEnd": {{timeend}},
	// 	}
	// }
	ShouldSpecifyFileds := []string{
		"Batch",
		"AllowStart",
		"AllowEnd",
	}
	for _, v := range ShouldSpecifyFileds {
		_, ok := reqData[v]
		if !ok {
			err = fmt.Errorf("%s should be specified", v)
			z.Error(err.Error())
			return
		}
	}
	var pro TInsuranceTypesPro

	err = json.Unmarshal(req.Data, &pro)
	if err != nil {
		z.Error(err.Error())
		return
	}
	orgSet := "("
	for _, o := range pro.Batch {
		orgSet += strconv.FormatInt(o, 10) + "),("
	}
	orgSet = orgSet[:len(orgSet)-2]

	formatOrg := ""
	formatOrgValue := make([]interface{}, 0)
	for n := 1; n <= len(pro.Batch); n++ {
		formatOrg += "$" + strconv.Itoa(n) + ","
		formatOrgValue = append(formatOrgValue, pro.Batch[n-1])
	}
	formatOrg = formatOrg[:len(formatOrg)-1]

	beforeUpdate := fmt.Sprintf(`delete from t_insurance_types where org_id in (%s) and parent_id = 10040;`, formatOrg)
	z.Info(beforeUpdate)
	var stmt *sqlx.Stmt
	stmt, err = tx.Preparex(beforeUpdate)
	if err != nil {
		z.Error(err.Error())
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(formatOrgValue...)
	if err != nil {
		z.Error(err.Error())
		return
	}
	timeNow := GetNowInMS()

	for _, rule := range pro.BatchRule {
		updateFields := []string{
			//------固定项
			"parent_id",
			"name",
			"data_type",
			"define_level",
			"layout_level",
			//------更改项
			"ref_id",
			"max_insure_in_year",
			"allow_start",
			"allow_end",
			"update_time",
			"rule_batch",
			"indate_start",
		}
		// i.Name = "区域投保规则"
		// i.DataType = null.StringFrom("4")
		// i.DefineLevel = null.IntFrom(0)
		// i.LayoutLevel = null.IntFrom(0)
		updateValues := make([]interface{}, 12)
		updateValues[0] = int64(10040)
		updateValues[1] = "区域投保规则"
		updateValues[2] = "4"
		updateValues[3] = int64(0)
		updateValues[4] = int64(0)

		updateValues[5] = rule.RefID.Int64
		updateValues[6] = rule.MaxInsureInYear.Int64
		updateValues[7] = pro.AllowStart.Int64
		updateValues[8] = pro.AllowEnd.Int64
		updateValues[9] = timeNow

		updateValues[10] = fmt.Sprintf("%s-%s",
			time.Unix(pro.AllowStart.Int64/1000, 0).Format("2006_01_02"),
			time.Unix(pro.AllowEnd.Int64/1000, 0).Format("2006_01_02")) //rule_batchs
		updateValues[11] = timeNow

		if !rule.RefID.Valid || rule.RefID.Int64 <= 0 {
			err = fmt.Errorf("please Specify RefID in BatchRule")
			z.Error(err.Error())
			return
		}
		if !rule.MaxInsureInYear.Valid || rule.MaxInsureInYear.Int64 <= 0 {
			err = fmt.Errorf("please Specify MaxInsureInYear in BatchRule")
			z.Error(err.Error())
			return
		}

		if rule.Status.Valid {
			if !inSlice(rule.Status.String, []string{"0", "4"}) {
				err = fmt.Errorf("please Specify Status in BatchRule")
				z.Error(err.Error())
				return
			}
			updateFields = append(updateFields, "status")
			updateValues = append(updateValues, rule.Status.String)
		}
		if rule.InsuredStartTime.Valid {
			updateFields = append(updateFields, "insured_start_time")
			updateValues = append(updateValues, rule.InsuredStartTime)
		}
		format := ""
		for n := 1; n <= len(updateFields); n++ {
			format += "$" + strconv.Itoa(n) + ","
		}
		format = format[:len(format)-1]
		updateF := strings.Join(updateFields, ",")

		s := fmt.Sprintf(`
		with va(org_id) as (values %s)
			INSERT INTO t_insurance_types (
				org_id,%s
			)
				select va.org_id,%s from va
			on conflict(org_id,ref_id) do update set (%s)=
			(%s)
			returning id
		`, orgSet, updateF, format, updateF, format)
		z.Info(s)
		stmt, err = tx.Preparex(s)
		if err != nil {
			z.Error(err.Error())
			return
		}
		defer stmt.Close()
		var result *sqlx.Rows
		result, err = stmt.Queryx(updateValues...)
		if err != nil {
			z.Error(err.Error())
			return
		}
		defer result.Close()

		for result.Next() {
			var r int64
			err = result.Scan(&r)
			if err != nil {
				z.Error(err.Error())
				return
			}
			idSet = append(idSet, r)
		}
	}
	tx.Commit()

	return
}

func validateInsurePlan(i *TVInsuranceType) (valid bool, err error) {
	if i == nil {
		err = fmt.Errorf("call validateInsurePlan with nil param")
		z.Error(err.Error())
		return
	}

	verifyBase := map[string]string{
		"Name":       "保险产品名称",
		"ParentID":   "隶属保险产品分类",
		"PayType":    "支付方式",
		"PayChannel": "支付渠道",
		"Price":      "价格",
		"Insurer":    "承保公司",
	}

	verifyTeam := map[string]string{
		"MaxInsureInYear":  "最长投保年限（年）",
		"InsuredInMonth":   "保障时长（月）",
		"InsuredStartTime": "起保日期",
		"InsuredEndTime":   "止保日期",
	}

	buf, err := MarshalJSON(i)
	if err != nil {
		return
	}

	plan := string(buf)
	if plan == "" {
		err = fmt.Errorf("%s(%d)无效配置", i.Name.String, i.ID.Int64)
		z.Error(err.Error())
		return
	}

	//校验基本信息
	for k, v := range verifyBase {
		r := gjson.Get(plan, k)
		switch r.Type {
		case gjson.String:
			if r.Str != "" {
				continue
			}

			if v == "PayChannel" && i.PayType.String != "在线支付" {
				continue
			}

			err = fmt.Errorf("%s(%d).%s无效配置", i.Name.String, i.ID.Int64, v)

			return
		case gjson.Number:
			if r.Num > 0 {
				continue
			}

			err = fmt.Errorf("%s(%d).%s无效配置", i.Name.String, i.ID.Int64, v)
			return
		default:
			if k == "PayChannel" && i.PayType.String != "在线支付" {
				continue
			}

			//学生意外险: 10040
			if k == "Price" && i.ParentID.Int64 != 10040 {
				continue
			}

			err = fmt.Errorf("%s(%d).%s无效配置", i.Name.String, i.ID.Int64, v)
			z.Error(err.Error())
			return
		}
	}

	maxInsureInYear := gjson.Get(plan, "MaxInsureInYear").Int()
	insuredInMonth := gjson.Get(plan, "InsuredInMonth").Int()
	insuredStartTime := gjson.Get(plan, "InsuredStartTime").Int()
	insuredEndTime := gjson.Get(plan, "InsuredEndTime").Int()

	if insuredStartTime > 0 && insuredEndTime > 0 &&
		insuredStartTime > insuredEndTime {
		err = fmt.Errorf("%s(%d)起保时间大于止保时间", i.Name.String, i.ID.Int64)
		z.Error(err.Error())
		return
	}

	if insuredStartTime == 0 && insuredEndTime > 0 {
		err = fmt.Errorf("%s(%d)有止保时间，则必须有起保时间", i.Name.String, i.ID.Int64)
		z.Error(err.Error())
		return
	}

	allowStart := gjson.Get(plan, "AllowStart").Int()
	allowEnd := gjson.Get(plan, "AllowEnd").Int()

	var isTeam bool
	if i.OrgID.Valid && i.OrgID.Int64 > 0 {
		isTeam = true
	}

	if !isTeam {
		//是否不符合散单规则
		verifyScattered := map[string]string{
			"MaxInsureInYear": "不能设置最长投保年限（年）[MaxInsureInYear]",
			"InsuredInMonth":  "保障时长（月）必须为12[InsuredInMonth]",

			"InsuredStartTime": "不能设置起保日期[InsuredStartTime]",
			"InsuredEndTime":   "不能设置止保日期[InsuredEndTime]",

			"AllowStart": "不能设置投保开始日期[AllowStart]",
			"AllowEnd":   "不能设置投保结束日期[AllowEnd]",
		}

		var hint string
		switch {
		case maxInsureInYear != 0:
			hint = verifyScattered["MaxInsureInYear"]

			//散单也允许设置起保日期了
		//case insuredStartTime != 0:
		//	hint = verifyScattered["InsuredStartTime"]

		case insuredEndTime != 0:
			hint = verifyScattered["InsuredEndTime"]

		case insuredInMonth != 12:
			hint = verifyScattered["InsuredInMonth"]

		case allowStart != 0:
			hint = verifyScattered["AllowStart"]

		case allowEnd != 0:
			hint = verifyScattered["AllowEnd"]
		}

		if hint != "" {
			err = fmt.Errorf("%s(%d)散单规则"+hint, i.Name.String, i.ID.Int64)
			z.Error(err.Error())
			return
		}

		//通过了，也不见得符合规则
	}

	if isTeam {
		//是否不符合团单规则
		teamHint := map[string]string{
			"AllowStart": "必须设置投保开始日期[AllowStart]",
			"AllowEnd":   "必须设置投保结束日期[AllowEnd]",
		}
		var hint string
		switch {
		case allowStart == 0:
			hint = teamHint["AllowStart"]

		case allowEnd == 0:
			hint = teamHint["AllowEnd"]
		}

		if hint != "" {
			err = fmt.Errorf("%s(%d)团单规则"+hint, i.Name.String, i.ID.Int64)
			z.Error(err.Error())
			return
		}

		//通过了，也不见得符合规则
	}

	//特定起保、止保时间
	if insuredStartTime > 0 && insuredEndTime > 0 && insuredEndTime > insuredStartTime {
		valid = true
		return
	}

	//非12个保障时间，有起保时间
	if insuredInMonth != 12 && insuredStartTime > 0 {
		valid = true
		return
	}

	//1年保障时间，指定起保日期
	if insuredInMonth == 12 && maxInsureInYear >= 1 &&
		insuredStartTime > 0 {
		valid = true
		return
	}

	//1年保障时间，不指定起保日期, 团单、散单皆可
	if insuredInMonth == 12 && maxInsureInYear >= 1 {
		valid = true
		return
	}

	//团单、散单皆可, 保障时间必须是1年
	if maxInsureInYear == 0 && insuredEndTime == 0 && insuredInMonth == 12 {
		valid = true
		return
	}

	//非12个保障时间，无起保时间
	if insuredInMonth != 12 {
		valid = true
		return
	}

	commonHint := verifyTeam["MaxInsureInYear"] + "," +
		verifyTeam["InsuredInMonth"] + "," +
		verifyTeam["InsuredStartTime"] + "," +
		verifyTeam["InsuredEndTime"]

	err = fmt.Errorf("%s(%d)的以下项目即不符合团单规则也不符合散单规则：\n\t"+
		commonHint, i.Name.String, i.ID.Int64)
	z.Error(err.Error())
	return
}

var (
	//--------以下参数在param.initGlobalParam()中会再赋值
	matchLevel1 int64 = 30
	matchLevel2 int64 = 61
	matchLevel3 int64 = 92
	matchLevel4 int64 = 182
	matchLevel5 int64 = 364

	organizationLevel1 int64 = 30
	organizationLevel2 int64 = 61
	organizationLevel3 int64 = 92
	organizationLevel4 int64 = 182
	organizationLevel5 int64 = 365

	contestReqStandard  float64 = 5000 //大于50元/人即可请求议价
	contestInsuredType1         = "学生/未成年人/教师"
	contestInsuredType2         = "成年人"
)

var regionsJSON types.JSONText
