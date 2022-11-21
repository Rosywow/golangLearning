package cmn

import (
	"archive/zip"
	"bytes"
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"

	"w2w.io/null"

	"github.com/jackc/pgx/v4"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/types"
)

var ownerItemToTable = map[string]string{
	"rptClaims.BankCardPic":              "t_report_claims.bank_card_pic",              //赔付银行卡/存折
	"rptClaims.InjuredIDPic":             "t_report_claims.injured_id_pic",             //出险人身份证
	"rptClaims.GuardianIDPic":            "t_report_claims.guardian_id_pic",            //监护人身份证
	"rptClaims.RelationProvePic":         "t_report_claims.relation_prove_pic",         //出险人与收款人关系证明
	"rptClaims.InvoicePic":               "t_report_claims.invoice_pic",                //医疗费用发票
	"rptClaims.BillsPic":                 "t_report_claims.bills_pic",                  //门诊费用清单
	"rptClaims.HospitalizedBillsPic":     "t_report_claims.hospitalized_bills_pic",     //住院费用清单
	"rptClaims.MedicalRecordPic":         "t_report_claims.medical_record_pic",         //门诊病历、门诊检查报告等
	"rptClaims.DiseaseDiagnosisPic":      "t_report_claims.disease_diagnosis_pic",      //疾病诊断证明
	"rptClaims.DischargeAbstractPic":     "t_report_claims.discharge_abstract_pic",     //出院证明
	"rptClaims.OtherPic":                 "t_report_claims.other_pic",                  //其他资料
	"rptClaims.PaidNoticePic":            "t_report_claims.paid_notice_pic",            // 保险金给付通知书
	"rptClaims.ClaimApplyPic":            "t_report_claims.claim_apply_pic",            // 索赔申请书
	"rptClaims.EquityTransferFile":       "t_report_claims.equity_transfer_file",       // 权益转让书
	"rptClaims.MatchProgrammePic":        "t_report_claims.match_programme_pic",        // 已有投保单位盖章的比赛秩序册
	"rptClaims.PolicyFile":               "t_report_claims.policy_file",                // 保单文件
	"rptClaims.DisabilityCertificate":    "t_report_claims.disability_certificate",     // 残疾证明
	"rptClaims.DeathCertificate":         "t_report_claims.death_certificate",          // 身故文件
	"rptClaims.StudentStatusCertificate": "t_report_claims.student_status_certificate", // 学籍证明

	"rptClaims.OrgLicPic": "t_report_claims.org_lic_pic",

	"rptClaims.CourierSnPic": "t_report_claims.courier_sn_pic",
	"rptClaims.AddiPic":      "t_report_claims.addi_pic",

	"rptClaims.DignosticInspectionPic": "t_report_claims.dignostic_inspection_pic",

	"insureAttach.files": "t_insure_attach.files",
	"price.files":        "t_price.files",

	"school.CreditCodePic": "t_school.credit_code_pic",

	"order.files": "t_order.files",

	"insuranceTypes.files":      "t_insurance_types.files",
	"insuranceTypes.OtherFiles": "t_insurance_types.other_files",

	"insuranceTypes.transfer_auth_files": "t_insurance_types.transfer_auth_files",

	"mistakeCorrect.files":            "t_mistake_correct.files",
	"mistakeCorrect.ApplicationFiles": "t_mistake_correct.application_files",
	"resource.picture":                "t_resource.picture", //客户服务图片

	"msg.files": "t_msg.files",
}

type fileDesc struct {
	SN     null.Int    `json:"SN,omitempty"`
	Label  null.String `json:"label,omitempty"`
	Name   null.String `json:"name,omitempty"`
	Digest null.String `json:"digest,omitempty"`

	//在文件表(t_file)中的id，用于删除时的文件标识
	FileID  null.Int `json:"fileID,omitempty"`
	FileOID null.Int `json:"fileOID,omitempty"`
}

type fileOwnDesc struct {
	SN null.Int `json:"SN,omitempty"`

	// 文件所属的表
	OwnerType null.String `json:"ownerType,omitempty"`

	// 文件所属表的列
	Item null.String `json:"item,omitempty"`

	// 文件所属的表的行
	LinkID null.Int `json:"linkID,omitempty"`

	MD5 null.String `json:"md5,omitempty"`

	//服务器文件存储路径，不含文件名
	Path null.String `json:"path,omitempty"`

	//客户上传的文件名，如果重复，则会被修改为: $MD5_$Name
	Name null.String `json:"name,omitempty"`

	Label null.String `json:"label,omitempty"`

	// 文件物理大小，字节为单位
	Size null.Int `json:"size,omitempty"`

	//t_file的ID
	FileID null.Int `json:"fileID,omitempty"`

	//保存于数据库中的pg_large object.loid
	FileOID null.Int `json:"fileOID,omitempty"`

	//是否保留原物理文件
	//于函数deleteFileFromTableField(*fileOwnDesc)使用,为true时将清除字段信息,但保留物理文件和file表中信息
	reservedFile bool
}

var fileStoreToDB bool
var fileStorePath string
var maxUploadFileSize = int64(2 * 1024 * 1024)

/*
	上传更新文件信息到数据库表中
	SN int，文件序号，如果未指定则添加到文件列表的最后，如果指定，则插入指定位置，原位
		置文件顺推，如原文件列表为0,1,2, 如果SN: 1则新文件插入1位置，原1变为2，原2变为3，即，0（旧),1(新),2（旧),3（旧)
	Item string，代表表中的哪一列
	OwnerType string,代表哪个表
	LinkID int, 代表表中的哪一行
	MD5 string，摘要进行重复性检测，如果重复，则只添加描述不复制文件本身，删除时，直到引用为零后才物理删除
	文件存储与更新算法：
	belongtoPath: 代表用户路径，由 ownerType/Item/LinkID/FileName 构成

	1. 获得文件重复度
	select digest,count(id) as row_count
		from t_file
		where digest=新文件的MD5
		group by digest;
		if row_count > 0 {
			fdExists=true
		}else{
			fdExists=false
		}
	2. belongtoPath是否存在
	select digest
		from t_file
		where belongto_path=新文件的belongtoPath
		group by digest;
		if digest!=""  {
			pathExists=true
		}else{
			pathExists=false
		}
	3. if fdExists && pathExists{
		文件，路径皆相同,什么也不做
	}
	4. if (fdExists && !pathExists) || (!fdExists && !pathExists){
		文件相同，路径不同,创建新记录
	}
	5. if !fdExists && pathExists{
		文件不相同，路径相同,创建新记录文件名改为 name+digest
	}
	6.创建/更新目标表、列、行信息
	7. if !fdExists
		保存文件


*/
func fileDescUpdate(ctx context.Context, f *fileOwnDesc) (fdExists, pathExists bool, err error) {
	q := GetCtxValue(ctx)

	if q.SysUser == nil || !q.SysUser.ID.Valid || q.SysUser.ID.Int64 <= 0 {
		err = fmt.Errorf("用户身份已过期,请重新登录")
		z.Error(err.Error())
		return
	}

	if f == nil ||
		!f.OwnerType.Valid || f.OwnerType.String == "" ||
		!f.Item.Valid || f.Item.String == "" ||
		!f.LinkID.Valid || f.LinkID.Int64 == 0 ||
		!f.Name.Valid || f.Name.String == "" ||
		!f.Path.Valid || f.Path.String == "" ||
		!f.MD5.Valid || f.MD5.String == "" ||
		!f.Size.Valid || f.Size.Int64 == 0 {
		err = fmt.Errorf("invalid fileOwnDesc:%+v", f)
		z.Error(err.Error())
		return
	}

	s := `select count(id) as row_count
		from t_file
		where digest=$1
		group by digest;`
	var stmt *sqlx.Stmt
	stmt, err = sqlxDB.Preparex(s)
	if err != nil {
		z.Error(err.Error())
		return
	}
	defer func() { _ = stmt.Close() }()
	row := stmt.QueryRowx(f.MD5)

	var rowCount null.Int
	//var belongToPath null.String

	err = row.Scan(&rowCount)
	if err == nil {
		fdExists = true
	}
	if err == sql.ErrNoRows {
		err = nil
	}
	if err != nil {
		z.Error(err.Error())
		return
	}

	s = `select digest,id
		from t_file
		where belongto_path=$1`
	stmt, err = sqlxDB.Preparex(s)
	if err != nil {
		z.Error(err.Error())
		return
	}
	defer func() { _ = stmt.Close() }()

	belongToPath := fmt.Sprintf("%s/%s/%d/%s",
		f.OwnerType.String, f.Item.String, f.LinkID.Int64, f.Name.String)

	row = stmt.QueryRowx(belongToPath)

	var digest null.String
	var fileID null.Int
	err = row.Scan(&digest, &fileID)
	if err == nil {
		pathExists = true
	}
	if err == sql.ErrNoRows {
		err = nil
	}
	if err != nil {
		z.Error(err.Error())
		return
	}
	if pathExists && (!digest.Valid || digest.String == "") {
		err = fmt.Errorf("digest is empty with belongtoPath=%s", belongToPath)
		z.Error(err.Error())
		return
	}
	if fdExists && pathExists {
		z.Info(fmt.Sprintf("用户重复上传同一文件:%s", f.Name.String))
		f.FileID = fileID
		return
	}
	if !fdExists && pathExists {
		f.Name = null.StringFrom(fmt.Sprintf("%s_%s", f.MD5.String, f.Name.String))
	}
	//注意,下面的path不包含文件，仅是文件系统存储路径
	s = `insert into t_file(path,file_name,digest,belongto_path,create_time,
			size,creator,file_oid) 
		values($1,$2,$3,$4,$5,$6,$7,$8)  
		RETURNING ID`
	{
		//var stmt *sql.Stmt
		stmt, err = sqlxDB.Preparex(s)
		if err != nil {
			z.Error(err.Error())
			return
		}
		defer func() { _ = stmt.Close() }()
		r := stmt.QueryRow(f.Path, f.Name, f.MD5, belongToPath,
			GetNowInMS(), f.Size, q.SysUser.ID.Int64, f.FileOID)
		var id int64
		err = r.Scan(&id)
		if err != nil && err != sql.ErrNoRows {
			z.Error(err.Error())
			return
		}

		if id > 0 {
			z.Info(fmt.Sprintf("insert successfully with id = %d", id))
		}
		f.FileID = null.IntFrom(id)
	}
	return
}

func getTableField(f *fileOwnDesc) (tableField string, err error) {
	if f == nil || !f.OwnerType.Valid || f.OwnerType.String == "" ||
		!f.Item.Valid || f.Item.String == "" {
		err = fmt.Errorf("call getTableField with invalid f")
		z.Error(err.Error())
		return
	}

	key := fmt.Sprintf("%s.%s", f.OwnerType.String, f.Item.String)
	var ok bool
	tableField, ok = ownerItemToTable[key]
	if !ok {
		err = fmt.Errorf("failed to find %s in ownerItemToTable", key)
		z.Error(err.Error())
	}
	return
}

func updateTableField(f *fileOwnDesc) (filesField string, err error) {
	errMsg := ""
	switch {
	case f == nil:
		errMsg = "f"
	case !f.OwnerType.Valid || f.OwnerType.String == "":
		errMsg = "ownerType"
	case !f.Item.Valid || f.Item.String == "":
		errMsg = "item"
	case !f.LinkID.Valid || f.LinkID.Int64 == 0:
		errMsg = "linkID"
	case !f.Name.Valid || f.Name.String == "":
		errMsg = "name"
	case !f.FileID.Valid || f.FileID.Int64 == 0:
		errMsg = "fileID"
	}

	if errMsg != "" {
		err = fmt.Errorf("invalid %s", errMsg)
		z.Error(err.Error())
		return
	}

	var tableField string
	tableField, err = getTableField(f)
	if err != nil {
		return
	}
	desc := strings.Split(tableField, ".")
	if len(desc) != 2 {
		err = fmt.Errorf("invalid value for key: %s in ownerItemToTable",
			f.OwnerType.String+"."+f.Item.String)
		z.Error(err.Error())
		return
	}
	s := fmt.Sprintf(`select %s from %s where id=$1`, desc[1], desc[0])
	var stmt *sqlx.Stmt
	stmt, err = sqlxDB.Preparex(s)
	if err != nil {
		z.Error(err.Error())
		return
	}
	defer func() { _ = stmt.Close() }()
	row := stmt.QueryRowx(f.LinkID.Int64)

	var fileList null.String
	err = row.Scan(&fileList)
	if err == sql.ErrNoRows {
		err = fmt.Errorf("不存在%s.%s.%d", f.OwnerType.String,
			f.Item.String, f.LinkID.Int64)
		z.Error(err.Error())
		return
	}
	if err != nil {
		z.Error(err.Error())
		return
	}

	jsonStr := "[]"
	if fileList.Valid && fileList.String != "" &&
		fileList.String != "{}" && strings.Index(fileList.String, "[") == 0 {
		jsonStr = fileList.String
	}

	var files []fileDesc
	err = json.Unmarshal([]byte(jsonStr), &files)
	if err != nil {
		z.Error(err.Error())
		return
	}
	fileNum := len(files)

	//newFile 新上传的文件
	newFile := fileDesc{
		FileID: null.IntFrom(f.FileID.Int64),   //必须有
		Name:   null.StringFrom(f.Name.String), // 必须有
		Digest: null.StringFrom(f.MD5.String),  // 必须有

		SN: null.NewInt(f.SN.Int64, f.SN.Valid), // 可以没有

		Label: null.NewString(f.Label.String, f.Label.Valid), // 可以没有
	}

	if f.FileOID.Valid && f.FileOID.Int64 > 0 {
		newFile.FileOID = f.FileOID
	}

	//确定该文件的序号，该序号主要用于前端的文件显示次序
	if !newFile.SN.Valid {
		//未指定SN则放到最后
		newFile.SN = null.IntFrom(int64(fileNum))
		files = append(files, newFile)
	} else if fileNum == 0 {
		//原来没有文件，指定SN为零，放到最前
		newFile.SN = null.IntFrom(0)
		files = append(files, newFile)
	} else {

		n := int(newFile.SN.Int64)
		if n >= fileNum {
			//指定SN大于现在文件数，则放到最后
			n = fileNum
		} else if n < 0 {
			//指定SN小于零，则放到最前
			n = 0
		}

		files = append(files[:n], append([]fileDesc{newFile}, files[n:]...)...)

		//重新按files数组的自然次序排序
		for i := 0; i < fileNum+1; i++ {
			files[i].SN = null.IntFrom(int64(i))
		}
	}

	var buf []byte
	buf, err = json.Marshal(&files)
	if err != nil {
		z.Error(err.Error())
		return
	}

	s = fmt.Sprintf(`update %s set %s=$1 where id=$2 returning %s`,
		desc[0], desc[1], desc[1])

	z.Info(fmt.Sprintf("%s %s %d", s, string(buf), f.LinkID.Int64))

	stmt, err = sqlxDB.Preparex(s)
	if err != nil {
		z.Error(err.Error())
		return
	}
	defer func() { _ = stmt.Close() }()

	var savedFilesDesc null.String
	err = stmt.QueryRow(string(buf), f.LinkID.Int64).Scan(&savedFilesDesc)
	if err != nil {
		z.Error(err.Error())
		return
	}

	if savedFilesDesc.Valid {
		filesField = savedFilesDesc.String
	}

	for _, v := range files {
		if v.Digest.String == f.MD5.String {
			f.SN = null.IntFrom(v.SN.Int64)
		}
	}
	return string(buf), nil
}

/*
	deleteFileFromTableField
	如果SN是负数，则删除该表/项/行/下的所有文件
	如果SN大于零且Name不为空，则删除Name相同的文件
	如果SN在当前文件序号中存在，则删除该表/项/行/序号对应的文件
	如果SN大于当前当前文件数，则删除最后一个

	删除后，将重新调整SN，以使SN保持连续
*/
func deleteFileFromTableField(f *fileOwnDesc) (err error) {
	if f == nil ||
		!f.OwnerType.Valid || f.OwnerType.String == "" ||
		!f.Item.Valid || f.Item.String == "" ||
		!f.LinkID.Valid || f.LinkID.Int64 == 0 || !f.SN.Valid {
		err = fmt.Errorf("call deleteFileFromTableField with invalid f")
		z.Error(err.Error())
		return
	}
	var tableField string
	tableField, err = getTableField(f)
	if err != nil {
		return
	}
	desc := strings.Split(tableField, ".")
	if len(desc) != 2 {
		err = fmt.Errorf("invalid value for key: %s in ownerItemToTable",
			f.OwnerType.String+"."+f.Item.String)
		z.Error(err.Error())
		return
	}
	s := fmt.Sprintf(`select %s from %s where id=$1`, desc[1], desc[0])
	var stmt *sqlx.Stmt
	stmt, err = sqlxDB.Preparex(s)
	if err != nil {
		z.Error(err.Error())
		return
	}
	defer func() { _ = stmt.Close() }()
	row := stmt.QueryRowx(f.LinkID.Int64)

	var fileList null.String
	err = row.Scan(&fileList)
	if err == sql.ErrNoRows {
		err = fmt.Errorf("不存在%s.%s.%d", f.OwnerType.String,
			f.Item.String, f.LinkID.Int64)
		z.Error(err.Error())
		return
	}
	if err != nil {
		z.Error(err.Error())
		return
	}

	jsonStr := "[]"
	if fileList.Valid && fileList.String != "" && fileList.String != "{}" {
		jsonStr = fileList.String
	}
	var files []fileDesc
	err = json.Unmarshal([]byte(jsonStr), &files)
	if err != nil {
		z.Error(err.Error())
		return
	}

	if len(files) <= 0 {
		err = fmt.Errorf("没有文件可以删除")
		z.Error(err.Error())
		return
	}

	var filesToDel []fileDesc
	fileSNToDel := f.SN.Int64

	if fileSNToDel < 0 {
		filesToDel = append(filesToDel, files...)
		files = []fileDesc{}
	} else if f.Name.Valid && f.Name.String != "" {
		//Delete file by Name
		for k, v := range files {
			if v.Name.String != f.Name.String {
				continue
			}
			filesToDel = append(filesToDel, v)
			files = append(files[:k], files[k+1:]...)
			break
		}
	} else {
		fileNum := len(files)
		if int(fileSNToDel) >= fileNum {
			filesToDel = append(filesToDel, files[fileNum-1])
			files = files[:fileNum-1]
		} else {
			filesToDel = append(filesToDel, files[fileSNToDel])
			//files = append(files[:i], files[i+1:]...)
			files = append(files[:fileSNToDel], files[fileSNToDel+1:]...)
		}
	}

	if len(filesToDel) <= 0 {
		err = fmt.Errorf("没有文件可以删除")
		z.Error(err.Error())
		return
	}

	for i := 0; i < len(files); i++ {
		files[i].SN = null.IntFrom(int64(i))
	}

	var buf []byte
	buf, err = json.Marshal(&files)
	if err != nil {
		z.Error(err.Error())
		return
	}
	s = fmt.Sprintf(`update %s set %s=$1 where id=$2`, desc[0], desc[1])
	var stmtX *sql.Stmt
	stmtX, err = sqlxDB.Prepare(s)
	if err != nil {
		z.Error(err.Error())
		return
	}
	defer func() { _ = stmtX.Close() }()
	var result sql.Result
	result, err = stmtX.Exec(string(buf), f.LinkID.Int64)
	if err != nil {
		z.Error(err.Error())
		return
	}
	var d int64
	if d, err = result.RowsAffected(); err != nil {
		z.Error(err.Error())
		return
	}
	if d <= 0 {
		err = fmt.Errorf("对%s/%s/%d/%s更新失败", f.OwnerType.String, f.Item.String, f.LinkID.Int64, f.Name.String)
		z.Error(err.Error())
		return
	}

	if !f.reservedFile {
		for _, e := range filesToDel {
			err = deleteFileByID(e.FileID.Int64)
			if err != nil {
				return
			}
		}
	}

	return
}

func deleteFileByID(fileID int64) (err error) {
	if fileID <= 0 {
		err = fmt.Errorf("call deleteFileByID with zero fileID")
		return
	}
	s := `select digest,path,count(id) as row_count
		from t_file
			where digest=(select digest from t_file where id=$1)
			group by digest,path`

	var stmt *sqlx.Stmt
	stmt, err = sqlxDB.Preparex(s)
	if err != nil {
		z.Error(err.Error())
		return
	}
	defer func() { _ = stmt.Close() }()

	row := stmt.QueryRowx(fileID)

	var digest null.String
	var path null.String
	var rowCount null.Int
	err = row.Scan(&digest, &path, &rowCount)
	if err == sql.ErrNoRows {
		err = fmt.Errorf("inexistent file with id=%d", fileID)
		z.Error(err.Error())
		return
	}
	if err != nil {
		z.Error(err.Error())
		return
	}
	if !digest.Valid || digest.String == "" {
		err = fmt.Errorf("invalid file digest with fileID=%d", fileID)
		z.Error(err.Error())
		return
	}

	s = `delete from t_file where id=$1`
	stmt, err = sqlxDB.Preparex(s)
	if err != nil {
		z.Error(err.Error())
		return
	}
	defer func() { _ = stmt.Close() }()

	var result sql.Result
	result, err = stmt.Exec(fileID)
	if err != nil {
		z.Error(err.Error())
		return
	}
	var d int64
	if d, err = result.RowsAffected(); err != nil {
		z.Error(err.Error())
		return
	}
	if d != 1 {
		err = fmt.Errorf("failed to delte fileInfo(id=%d)", fileID)
		z.Error(err.Error())
		return
	}

	rowCount.Int64--
	if rowCount.Int64 > 0 {
		return
	}

	// fileStorePath only use for save file to disk and save path to db then
	// when we want to delete/view file would use the 'path' column value from t_file table.

	if !path.Valid || path.String == "" {
		err = fmt.Errorf("file.id: %d is empty , delete failed", fileID)
		z.Error(err.Error())
		return
	}

	//fn := fileStorePath + digest.String
	fn := path.String + digest.String
	_, err = os.Stat(fn)
	if os.IsNotExist(err) {
		err = fmt.Errorf("%s inexistence, delete failed", fn)
		z.Error(err.Error())
		return
	}
	err = os.Remove(fn)
	if err != nil {
		z.Error(err.Error())
	}
	return
}

func fileView(ctx context.Context, view string) {
	q := GetCtxValue(ctx)
	z.Info("---->" + FncName())
	q.Stop = true

	if view == "" || len(view) != 32 {
		q.Err = fmt.Errorf("call fileView with empty/invalid idx")
		z.Error(q.Err.Error())
		q.RespErr()
		return
	}
	s := `select file_name,path,file_oid from t_file where digest=$1`
	var stmt *sqlx.Stmt
	stmt, q.Err = sqlxDB.Preparex(s)
	if q.Err != nil {
		z.Error(q.Err.Error())
		return
	}
	defer func() { _ = stmt.Close() }()
	row := stmt.QueryRowx(view)

	var fileName, path null.String
	var fileOID null.Int
	q.Err = row.Scan(&fileName, &path, &fileOID)

	if q.Err == sql.ErrNoRows {
		q.Err = fmt.Errorf("文件%s不存在", view)
		z.Error(q.Err.Error())
		q.RespErr()
		return
	}

	if q.Err != nil {
		z.Error(q.Err.Error())
		return
	}

	if !path.Valid || path.String == "" || !fileName.Valid || fileName.String == "" {
		q.Err = fmt.Errorf("查询时与'%s'相关的path、fileName无效", view)
		z.Error(q.Err.Error())
		q.RespErr()
		return
	}
	fn := path.String + view
	_, q.Err = os.Stat(fn)

	if os.IsNotExist(q.Err) {
		if fileOID.Int64 <= 0 {
			q.Err = fmt.Errorf("文件'%s'不存在", fn)
			z.Error(q.Err.Error())
			q.RespErr()
			return
		}

		var tx pgx.Tx
		tx, q.Err = pgxConn.Begin(ctx)
		if q.Err != nil {
			z.Error(q.Err.Error())
			q.RespErr()
			return
		}

		defer func() { _ = tx.Rollback(ctx) }()

		lo := tx.LargeObjects()

		var dbFile *pgx.LargeObject
		dbFile, q.Err = lo.Open(ctx, uint32(fileOID.Int64),
			pgx.LargeObjectModeRead|pgx.LargeObjectModeWrite)
		if q.Err != nil {
			z.Error(q.Err.Error())
			q.RespErr()
			_ = tx.Rollback(ctx)
			return
		}
		defer func() { _ = dbFile.Close() }()

		var fd *os.File
		fd, q.Err = os.OpenFile(fn, os.O_RDWR|os.O_CREATE, 0664)
		if q.Err != nil {
			z.Error(q.Err.Error())
			q.RespErr()
			return
		}
		defer func() { _ = fd.Close() }()

		var ubiety int64
		ubiety, q.Err = io.Copy(fd, dbFile)
		if q.Err != nil {
			z.Error(q.Err.Error())
			q.RespErr()
			return
		}
		if ubiety <= 0 {
			q.Err = fmt.Errorf("file(digest:%s)零长度文件拷贝", view)
			z.Error(q.Err.Error())
			q.RespErr()
			return
		}
	}

	fileName.String = url.QueryEscape(fileName.String) //防止IE下载文件名乱码
	q.W.Header().Set("Content-Disposition",
		fmt.Sprintf("inline; filename=%s", fileName.String))
	http.ServeFile(q.W, q.R, fn)
}

/*
1、文件上传接口
  接口: /api/file
  http method: POST
  参数：q
   {
      action: 'insert',
      data: {
        SN: 0, // 表示文件显示时的次序
        Item:   // 表示为哪一个项目的文件，如，身份证照片
        ownerType: 'rptClaims', //请填写此固定值
        linkID: 20012, // 报案理赔表的行编号，即ID
      },
    }
  参数file,
    作为文件内容，以post的body发送到服务端

  以下为item的值，取value,前端显示为name
    { value: 'BankCardPic', name: '银行卡/存折照片' },
    { value: 'InjuredIDPic', name: '被保险人身份证照片' },
    { value: 'GuardianIDPic', name: '监护人身份证照片' },
    { value: 'OrgLicPic', name: '营业执照照片' },
    { value: 'RelationProvePic', name: '与被保险人关系证明照片' },
    { value: 'BillsPic', name: '门诊/住院费用清单照片' },
    { value: 'InvoicePic', name: '医疗费用发票照片' },
    { value: 'MedicalRecordPic', name: '病历照片' },
    { value: 'DignosticInspectionPic', name: '检验检查报告照片' },
    { value: 'DischargeAbstractPic', name: '出院小结照片' },
    { value: 'OtherPic', name: '其它资料照片' },
    { value: 'CourierSnPic', name: '快递单号照片' },
		{ value: 'AddiPic', name: '补充照片' },
使用说明:
   1、SN可以不填写，则该文件添加到最后一个
   2、如果填写SN则按SN的值进行插入，例如，原来的文件列表为:1,a.doc;2,b.doc;3,c.doc
    此时，以SN:2,file:x.doc调用接口，则结果为：1,a.co;2,x.doc;3,b.doc,4:c.doc
   3、如果SN为零，则插入到第一个之前，原有的后推
   4、如果SN的值大于现存的文件个数，则新插入的文件放到最后一个
   5、不允许重复上传同一文件，如果上传同一上件后端会返回-3000错误，表示该文件已上传过了。
   6、允许上传同名，但内容不同的文件，该文件名将会被加前辍，以示区别同名的另一个文件，前辍无规则。

用例请参考: report-claims.component.ts, onFile(e)函数

2、文件删除接口(如无必要，请勿对接)
  接口: /api/file
  http method: delete
  参数：q
   {
      action: 'insert',
      data: {
        SN: 0, // 表示被删除文件的显示次序
        Item:   // 表示为哪一个项目的文件，如，身份证照片,取值请参考文件上传接口
        ownerType: 'rptClaims', //请填写此固定值
        linkID: 20012, // 报案理赔表的行编号，即ID
      },
    }
使用说明：
  1、每一个参数都必填；
  2、SN=-1表示删除所有项目下的文件，小心使用；
  3、SN大于现有文件数，则表示删除最后一个文件；

二、各个表中files的格式
1、格式为
	[{文件1信息},{文件2信息},...]

2、说明
	******* 前端绝对不要回传任何files列 ******
	后端必须把把前端传回的files列清除掉，对于业务需要校验files列的时候，则根据主键id查出数据库中的files进行校验。
  前端只能读后端返回的files数据，***不能修改***
  如果需要修改，只能通过/api/file接口进行，即 /api/file post为增加, /api/file delete为删除。

3、样例t_order.files

[
  {
    "SN": 0,
    "name": "转账授权说明.docx",
    "label": "转账授权说明",
    "digest": "76fff418a598078fee2e16ac4f525bef",
    "fileID": 21436
  },
  {
    "SN": 1,
    "name": "橙色3.jpg",
    "label": "统一社会信用代码证书",
    "digest": "bc8b146ce4bd154423534dc6cdc8c906",
    "fileID": 21437
  },
  {
    "SN": 2,
    "name": "8517_投保单_20200614_140152.pdf",
    "label": "投保单",
    "digest": "3e75a8525d50d532e79915e7306507fc",
    "fileID": 21438
  },
  {
    "SN": 3,
    "name": "8517_投保单excel_20200614_140152.xlsx",
    "label": "投保单excel",
    "digest": "7200fdc6344c97c2d99421ed476bb2ec",
    "fileID": 21439
  }
]

三、返回结果说明
失败: 上传文件中的任何一个失败，都会导致失败，q.Msg.Status < 0, 错误消息在q.Msg.Msg中,
	*** 注意重复上传的文件不提示，认为是上传成功 ***

成功: 上传成功的文件信息放在 q.Msg.Data中，格式如下
[{
	"SN": 0,
	"name": "8750e103ff2d8b19b07eb331e769236.jpg",
	"label": "工会信用代码证书",
	"digest": "ed8d87adccaf8cae51bcc840357ac805",
	"fileID": 25440
},{
	"SN": 1,
	"name": "19070_投保清单_20200916_090205.pdf",
	"label": "投保清单",
	"digest": "7b7607cf2b4f4916c0f9190fe8fa794c",
	"fileID": 25441
},{
	"SN": 2,
	"name": "19070_投保清单excel_20200916_090206.xlsx",
	"label": "投保清单excel",
	"digest": "d1229c5e20e48ad10cdedb2ebfad1864",
	"fileID": 25442
},{
	"SN": 3,
	"name": "19070_投保单_20200916_090206.pdf",
	"label": "投保单",
	"digest": "1e15c8eec3c3ba28ccfc95a4436c910c",
	"fileID": 25443,
	"new":true
},{
	"SN": 4,
	"name": "19070_投保单excel_20200916_090206.xlsx",
	"label": "投保单excel",
	"digest": "05b38cbc5dfeb95a061e6cadf6422a40",
	"fileID": 25444,
	"new":true
},{
	"SN": 5,
	"name": "6f313d231b6b77fd3338d64750f84c0.jpg",
	"label": "付款凭证回执",
	"digest": "5005e6419369002364731b924657bab3",
	"fileID": 25641,
	"new":true
}]

说明:
	1) q.Msg.Data中包含该linkID对应files列的所有的文件信息
	2) 数据中 "new":true, 表示刚刚成功上传的文件
*/
var fileDoor sync.Mutex

func qFile(ctx context.Context) {
	q := GetCtxValue(ctx)
	z.Info("---->" + FncName())

	fileDoor.Lock()
	defer func() {
		fileDoor.Unlock()
	}()

	q.Stop = true

	method := strings.ToLower(q.R.Method)
	view := q.R.URL.Query().Get("v")
	if method == "get" && view != "" {
		fileView(ctx, view)
		return
	}

	qry := q.R.URL.Query().Get("q")
	if qry == "" {
		q.Err = fmt.Errorf(`要提供文件信息，要不然不让你做啥哦`)
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

	if len(req.Data) == 0 {
		q.Err = fmt.Errorf(`请在data中指定文件信息,如{
			"SN":3,
			"label":"身份证正面",
			"item":"InjuredIDPic",
			"ownerType":"rptClaims",
			"linkID":12
		}`)
		z.Error(q.Err.Error())
		q.RespErr()
		return
	}
	var f fileOwnDesc
	q.Err = json.Unmarshal(req.Data, &f)
	if q.Err != nil {
		z.Error(q.Err.Error())
		q.RespErr()
		return
	}
	_ = InvalidEmptyNullValue(&f)
	if !f.LinkID.Valid || f.LinkID.Int64 == 0 ||
		!f.OwnerType.Valid || f.OwnerType.String == "" ||
		!f.Item.Valid || f.Item.String == "" {
		q.Err = fmt.Errorf(`请在data中指定文件用途标识,如,被保险人身份证照片{
			"SN": 0,
			"label":"身份证反面",
			"item": 'InjuredIDPic',
			"ownerType": 'rptClaims',
			linkID: 20223,
		}`)
		z.Error(q.Err.Error())
		q.RespErr()
		return
	}

	reqAction := strings.ToLower(req.Action)

	switch method {
	case "delete":
		if reqAction != "delete" {
			q.Err = fmt.Errorf("please specify delete as action")
			z.Error(q.Err.Error())
			q.RespErr()
			return
		}
		q.Err = deleteFileFromTableField(&f)
		if q.Err != nil {
			q.RespErr()
			return
		}
		q.Msg.Data = types.JSONText(`{"RowAffected":1}`)
		q.Resp()
		return

	case "post":
		if reqAction != "insert" {
			q.Err = fmt.Errorf("please specify insert as action")
			z.Error(q.Err.Error())
			q.RespErr()
			return
		}

		var filesField string
		filesField, q.Err = saveFiles(ctx, &f)
		if q.Err != nil {
			q.RespErr()
			return
		}
		q.Msg.Msg = filesField
		q.Msg.Data = []byte(filesField)
		q.Resp()
	}
}

//multiFilesUpload 是否支持单次请求上传多个文件
var multiFilesUpload bool

/* saveFiles 单请求多文件上传
失败: 上传文件中的任何一个失败，都会导致失败，q.Msg.Status < 0, 错误消息在q.Msg.Msg中,
	*** 注意重复上传的文件不提示，认为是上传成功 ***

成功: 上传成功的文件信息放在 q.Msg.Data中，格式如下
[{
	"SN": 0,
	"name": "8750e103ff2d8b19b07eb331e769236.jpg",
	"label": "工会信用代码证书",
	"digest": "ed8d87adccaf8cae51bcc840357ac805",
	"fileID": 25440
},{
	"SN": 1,
	"name": "19070_投保清单_20200916_090205.pdf",
	"label": "投保清单",
	"digest": "7b7607cf2b4f4916c0f9190fe8fa794c",
	"fileID": 25441
},{
	"SN": 2,
	"name": "19070_投保清单excel_20200916_090206.xlsx",
	"label": "投保清单excel",
	"digest": "d1229c5e20e48ad10cdedb2ebfad1864",
	"fileID": 25442
},{
	"SN": 3,
	"name": "19070_投保单_20200916_090206.pdf",
	"label": "投保单",
	"digest": "1e15c8eec3c3ba28ccfc95a4436c910c",
	"fileID": 25443,
	"new":true
},{
	"SN": 4,
	"name": "19070_投保单excel_20200916_090206.xlsx",
	"label": "投保单excel",
	"digest": "05b38cbc5dfeb95a061e6cadf6422a40",
	"fileID": 25444,
	"new":true
},{
	"SN": 5,
	"name": "6f313d231b6b77fd3338d64750f84c0.jpg",
	"label": "付款凭证回执",
	"digest": "5005e6419369002364731b924657bab3",
	"fileID": 25641,
	"new":true
}]

说明:
	1) q.Msg.Data中包含该linkID对应files列的所有的文件信息
	2) 数据中 "new":true, 表示刚刚成功上传的文件
	3) 如果单次请求中上传了多个文件，则这些文件的label相同
*/

func saveFiles(ctx context.Context, fd *fileOwnDesc) (filesField string, err error) {
	q := GetCtxValue(ctx)
	if fd == nil {
		err = fmt.Errorf("f is nil")
		z.Error(err.Error())
		return
	}
	err = q.R.ParseMultipartForm(1024 * 1024 * 32)
	if err != nil {
		z.Error(err.Error())
		return
	}
	formData := q.R.MultipartForm
	files := formData.File["file"]

	fileSN := int64(10000)
	if fd.SN.Valid {
		fileSN = fd.SN.Int64
	}

	var tx pgx.Tx
	tx, err = pgxConn.Begin(ctx)
	if err != nil {
		z.Error(err.Error())
		return
	}

	defer func() { _ = tx.Rollback(ctx) }()

	delta := fileSN
	var filesMD5 []string
	for i := range files {
		fileSN = int64(i) + delta

		srcFileInfo := files[i]
		if srcFileInfo.Filename == "" {
			err = fmt.Errorf("第%d个文件的名称为空", i)
			z.Error(err.Error())
			return
		}
		if srcFileInfo.Size == 0 {
			err = fmt.Errorf("第%d个文件(%s)的长度为零", i, srcFileInfo.Filename)
			z.Error(err.Error())
			return
		}

		if srcFileInfo.Size > maxUploadFileSize {
			err = fmt.Errorf("文件%s的大小为%4.2f兆，超过了系统的%4.2f兆限制",
				srcFileInfo.Filename,
				float64(srcFileInfo.Size)/(1024*1024),
				float64(maxUploadFileSize)/(1024*1024))
			z.Error(err.Error())
			return
		}

		var src multipart.File
		src, err = files[i].Open()
		if err != nil {
			err = fmt.Errorf("第%d个文件(%s)打开时出错: %s",
				i, srcFileInfo.Filename, err.Error())
			z.Error(err.Error())
			return
		}
		defer func() { _ = src.Close() }()

		hash := md5.New()
		_, err = io.Copy(hash, src)
		if err != nil {
			err = fmt.Errorf("第%d个文件(%s)摘要时出错: %s",
				i, srcFileInfo.Filename, err.Error())
			z.Error(err.Error())
			return
		}

		fd.MD5 = null.StringFrom(hex.EncodeToString(hash.Sum(nil)))
		fd.Path = null.StringFrom(fileStorePath)
		fd.Name = null.StringFrom(srcFileInfo.Filename)
		fd.Size = null.IntFrom(srcFileInfo.Size)
		fd.SN = null.IntFrom(fileSN)
		filesMD5 = append(filesMD5, fd.MD5.String)

		s := `select exists(select 1 from t_file where digest=$1)`
		r := sqlxDB.QueryRow(s, fd.MD5)
		var exists bool
		err = r.Scan(&exists)
		if err != nil {
			err = fmt.Errorf("查询第%d个文件(%s)(MD5: %s)时出错: %s",
				i, srcFileInfo.Filename, fd.MD5.String, err.Error())
			z.Error(err.Error())
			return
		}

		//如果文件在磁盘上不存在，则保存文件到磁盘
		if !exists {
			fn := fileStorePath + fd.MD5.String
			var dst *os.File
			dst, err = os.OpenFile(fn, os.O_WRONLY|os.O_CREATE, 0666)
			if err != nil {
				err = fmt.Errorf("后端打开第%d个文件(%s)(MD5: %s)时出错: %s",
					i, srcFileInfo.Filename, fd.MD5.String, err.Error())
				z.Error(err.Error())
				return
			}
			defer func() { _ = dst.Close() }()

			var ubiety int64
			ubiety, err = src.Seek(0, io.SeekStart)
			if err != nil {
				err = fmt.Errorf("倒带第%d个文件(%s)(MD5: %s)时出错: %s",
					i, srcFileInfo.Filename, fd.MD5.String, err.Error())
				z.Error(err.Error())
				return
			}
			if ubiety != 0 {
				err = fmt.Errorf("倒带第%d个文件(%s)(MD5: %s)时出错: 没倒到头",
					i, srcFileInfo.Filename, fd.MD5.String)
				z.Error(err.Error())
				return
			}

			ubiety, err = io.Copy(dst, src)
			if err != nil {
				err = fmt.Errorf("后端保存第%d个文件(%s)(MD5: %s)时出错: %s",
					i, srcFileInfo.Filename, fd.MD5.String, err.Error())
				z.Error(err.Error())
				return
			}

			if ubiety != srcFileInfo.Size {
				err = fmt.Errorf("写第%d个文件(%s)(MD5: %s)时出错: 没有写完整",
					i, srcFileInfo.Filename, fd.MD5.String)
				z.Error(err.Error())
				return
			}

			// 存文件到数据库pg_largeobject
			if fileStoreToDB {
				ubiety, err = src.Seek(0, io.SeekStart)
				if err != nil {
					err = fmt.Errorf("倒带第%d个文件(%s)(MD5: %s)时出错: %s",
						i, srcFileInfo.Filename, fd.MD5.String, err.Error())
					z.Error(err.Error())
					return
				}

				if ubiety != 0 {
					err = fmt.Errorf("倒带第%d个文件(%s)(MD5: %s)时出错: 没倒到头",
						i, srcFileInfo.Filename, fd.MD5.String)
					z.Error(err.Error())
					return
				}

				//-------

				lo := tx.LargeObjects()
				var oid uint32
				oid, err = lo.Create(ctx, 0)
				if err != nil {
					err = fmt.Errorf("后端保存第%d个文件(%s)(MD5: %s)时出错: %s",
						i, srcFileInfo.Filename, fd.MD5.String, err.Error())
					z.Error(err.Error())
					return
				}

				var dbFile *pgx.LargeObject
				dbFile, err = lo.Open(ctx, oid, pgx.LargeObjectModeRead|pgx.LargeObjectModeWrite)
				if err != nil {
					err = fmt.Errorf("后端保存第%d个文件(%s)(MD5: %s)时出错: %s",
						i, srcFileInfo.Filename, fd.MD5.String, err.Error())
					z.Error(err.Error())
					return
				}
				defer func() { _ = dbFile.Close() }()

				ubiety, err = io.Copy(dbFile, src)
				if err != nil {
					err = fmt.Errorf("后端保存第%d个文件(%s)(MD5: %s)时出错: %s",
						i, srcFileInfo.Filename, fd.MD5.String, err.Error())
					z.Error(err.Error())
					return
				}

				if ubiety != srcFileInfo.Size {
					err = fmt.Errorf("写第%d个文件(%s)(MD5: %s)时出错: 没有写完整",
						i, srcFileInfo.Filename, fd.MD5.String)
					z.Error(err.Error())
					return
				}
				fd.FileOID = null.IntFrom(int64(oid))
			}
		}

		//------------
		var fdExists, pathExists bool
		fdExists, pathExists, err = fileDescUpdate(ctx, fd) //t_file
		if err != nil {
			return
		}

		if fdExists && pathExists {
			continue
		}

		filesField, err = updateTableField(fd)
		if err != nil {
			return
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		z.Error(err.Error())
		return
	}

	if filesField == "" {
		var tableField string
		tableField, err = getTableField(fd)
		if err != nil {
			return
		}

		desc := strings.Split(tableField, ".")
		if len(desc) != 2 {
			err = fmt.Errorf("invalid value for key: %s in ownerItemToTable",
				fd.OwnerType.String+"."+fd.Item.String)
			z.Error(err.Error())
			return
		}
		s := fmt.Sprintf(`select %s from %s where id=$1`, desc[1], desc[0])

		var r null.String
		err = pgxConn.QueryRow(ctx, s, fd.LinkID.Int64).Scan(&r)
		if err != nil {
			z.Error(err.Error())
			return
		}

		filesField = r.String
		if filesField == "" {
			err = fmt.Errorf("上传后文件集竟然为空, ownerType: %s,item: %s, linkID: %d",
				fd.OwnerType.String, fd.Item.String, fd.LinkID.Int64)
			z.Error(err.Error())
			return
		}
	}

	for _, v := range filesMD5 {
		for i := 0; ; i++ {
			d := gjson.Get(filesField, fmt.Sprintf("%d.digest", i))
			if !d.Exists() {
				break
			}

			if d.Str != v {
				continue
			}

			filesField, err = sjson.Set(filesField, fmt.Sprintf("%d.new", i), true)
			if err != nil {
				z.Error(err.Error())
				return
			}

			filesField, err = sjson.Set(filesField, fmt.Sprintf("%d.linkID", i), fd.LinkID.Int64)
			if err != nil {
				z.Error(err.Error())
				return
			}

			filesField, err = sjson.Set(filesField, fmt.Sprintf("%d.ownerType", i), fd.OwnerType.String)
			if err != nil {
				z.Error(err.Error())
				return
			}

			filesField, err = sjson.Set(filesField, fmt.Sprintf("%d.item", i), fd.Item.String)
			if err != nil {
				z.Error(err.Error())
				return
			}

			break
		}
	}

	return
}

func saveFileBytes(buff []byte, f *fileOwnDesc) (err error) {
	z.Info("---->" + FncName())
	if len(buff) == 0 {
		err = fmt.Errorf("参数buff长度为空")
		z.Error(err.Error())
		return
	}
	if f == nil {
		err = fmt.Errorf("参数f(文件信息)为空")
		z.Error(err.Error())
		return
	}

	hash := md5.New()
	var written int
	written, err = hash.Write(buff)
	if err != nil {
		z.Error(err.Error())
		return
	}
	if written != len(buff) {
		err = fmt.Errorf("short written")
		z.Error(err.Error())
		return
	}

	f.MD5 = null.StringFrom(hex.EncodeToString(hash.Sum(nil)))
	f.Path = null.StringFrom(fileStorePath)
	f.Size = null.IntFrom(int64(len(buff)))

	s := `select exists(select 1 from t_file where digest=$1)`
	r := sqlxDB.QueryRow(s, f.MD5)
	var exists bool
	err = r.Scan(&exists)
	if err != nil {
		z.Error(err.Error())
		return
	}
	if exists {
		return nil
	}

	fn := fileStorePath + f.MD5.String
	var dst *os.File
	dst, err = os.OpenFile(fn, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		z.Error(err.Error())
		return
	}
	defer func() { _ = dst.Close() }()
	written, err = dst.Write(buff)
	if err != nil {
		z.Error(err.Error())
		return
	}
	if written != len(buff) {
		err = fmt.Errorf("short written")
		z.Error(err.Error())
		return
	}
	return
}

type fileToZip struct {
	FilePath string //文件路径（不支持文件夹）
	Header   string //文件头部信息，即在zip文档内的（路径+）文件名称，留空则使用默认文件名
}

func filesToZipBuffer(src []fileToZip) (buf *bytes.Buffer, err error) {

	//buf, _ := os.Create(dest)
	//defer d.Close()

	buf = &bytes.Buffer{}
	zw := zip.NewWriter(buf)
	defer func() { _ = zw.Close() }()

	for _, v := range src {
		var file *os.File
		file, err = os.Open(v.FilePath)

		if err != nil {
			if err != os.ErrNotExist {
				err = fmt.Errorf("找不到对应文件:%s 路径:%s", v.Header, v.FilePath)
				return
			}
			z.Error(err.Error())
			return
		}
		var info os.FileInfo
		info, err = file.Stat()
		if err != nil {
			z.Error(err.Error())
			return
		}
		if info.IsDir() {
			_ = file.Close()
			err = fmt.Errorf("不支持文件夹")
			z.Error(err.Error())
			return
		}
		var header *zip.FileHeader
		header, err = zip.FileInfoHeader(info)
		if err != nil {
			z.Error(err.Error())
			return
		}

		v.Header = strings.TrimRight(v.Header, "/")
		if v.Header == "" {
			v.Header = header.Name
		}
		header.Name = v.Header

		var writer io.Writer
		writer, err = zw.CreateHeader(header)
		if err != nil {
			z.Error(err.Error())
			return
		}
		_, err = io.Copy(writer, file)
		if err != nil {
			z.Error(err.Error())
			return
		}
		_ = file.Close()
	}
	return
}

//if fdExists && pathExists,调用此函数获取fileID
func getIDIfRepeatUpload(f *fileOwnDesc) (err error) {
	if f == nil ||
		!f.OwnerType.Valid || f.OwnerType.String == "" ||
		!f.Item.Valid || f.Item.String == "" ||
		!f.LinkID.Valid || f.LinkID.Int64 == 0 ||
		!f.Name.Valid || f.Name.String == "" ||
		!f.Path.Valid || f.Path.String == "" ||
		!f.MD5.Valid || f.MD5.String == "" ||
		!f.Size.Valid || f.Size.Int64 == 0 {
		err = fmt.Errorf("invalid fileOwnDesc")
		z.Error(err.Error())
		return
	}
	digest := f.MD5
	belongToPath := fmt.Sprintf("%s/%s/%d/%s",
		f.OwnerType.String, f.Item.String, f.LinkID.Int64, f.Name.String)
	s := `select id from t_file where digest= $1 and belongto_path = $2`
	stmt, err := sqlxDB.Preparex(s)
	if err != nil {
		z.Error(err.Error())
		return
	}
	defer func() { _ = stmt.Close() }()

	row := stmt.QueryRowx(digest, belongToPath)
	var id null.Int
	err = row.Scan(&id)
	if err == sql.ErrNoRows {
		err = fmt.Errorf("不存在的文件{digest:%s,belongToPath:%s},却调用了查询重复文件接口", digest.String, belongToPath)
		z.Error(err.Error())
		return
	}

	if err != nil {
		z.Error(err.Error())
		return
	}
	if !id.Valid || id.Int64 == 0 {
		err = fmt.Errorf("查询重复文件接口,却未获取文件id")
		z.Error(err.Error())
		return
	}
	f.FileID = id
	return
}

//fileDesc查表转换成fileOwnDesc,但是OwnerType/Item/LinkID/Label需手动指定
func getFileOwnDesc(f *fileDesc) (fDesc *fileOwnDesc, err error) {

	if f == nil || !f.FileID.Valid || f.FileID.Int64 == 0 {
		err = fmt.Errorf("invalid fileDesc")
		z.Error(err.Error())
		return
	}

	tFile, err := GetTFileByPk(sqlxDB, f.FileID)
	if err != nil {
		z.Error(err.Error())
		return
	}

	return &fileOwnDesc{
		MD5:    null.StringFrom(tFile.Digest),
		Path:   null.StringFrom(tFile.Path),
		Name:   null.StringFrom(tFile.FileName),
		Size:   tFile.Size,
		FileID: f.FileID,
	}, nil
}

type structWithFiles interface {
	GetTableName() string
	Fields() []string
}

func getFileDescByStruct(swf interface{}, id null.Int) (files []fileDesc, err error) {
	if swf == nil || !id.Valid || id.Int64 <= 0 {
		err = fmt.Errorf("struct is nil or id is invalid")
		z.Error(err.Error())
		return
	}

	u, ok := swf.(structWithFiles)
	if !ok {
		err = fmt.Errorf("could not parse as interface structWithFiles")
		z.Error(err.Error())
		return
	}
	if !inSlice("ID", u.Fields()) || !inSlice("Files", u.Fields()) {
		err = fmt.Errorf("missing ID/Files in %v", swf)
		z.Error(err.Error())
		return
	}

	s := fmt.Sprintf(`select coalesce(files,'{}') from %s where id = %d`, u.GetTableName(), id.Int64)

	var f types.JSONText
	err = sqlxDB.QueryRowx(s).Scan(&f)
	if err != nil {
		z.Error(err.Error())
		return
	}
	if len(f) <= 2 {
		return
	}
	err = f.Unmarshal(&files)
	if err != nil {
		z.Error(err.Error())
		return
	}

	return
}

/*
	getFileMD5 获取文件的MD5值
		buf 内存缓冲，如果非nil,则忽略src,fn
		src of.Open[File]()返回的文件描述符, 如果该值存在则不使用fn
		fn 文件名，如果src为nil则会使用fn作为文件名
		md5Sum 文件的MD5值
		err 返回出错信息
*/
func getFileMD5(buf []byte, src io.Reader, fn string) (md5Sum string, err error) {
	if buf != nil && len(buf) > 0 {
		src = bytes.NewReader(buf)
	}

	if src == nil {
		if fn == "" {
			err = fmt.Errorf("call getFileMD5 with empty buf and fn and src")
			z.Error(err.Error())
			return
		}

		var fd *os.File
		fd, err = os.Open(fn)
		if os.IsNotExist(err) {
			z.Error(fmt.Sprintf("%s 不存在", fn))
			return
		}

		if err != nil {
			z.Error(err.Error())
			return
		}
		defer func() { _ = fd.Close() }()
		src = fd
	}

	hash := md5.New()
	_, err = io.Copy(hash, src)
	if err != nil {
		z.Error(err.Error())
		return
	}
	md5Sum = hex.EncodeToString(hash.Sum(nil))
	return
}

func saveFile(ctx context.Context, fileInfo *multipart.FileHeader,
	fd *fileOwnDesc) (filesField string, err error) {

	fileSN := int64(10000)
	if fd.SN.Valid {
		fileSN = fd.SN.Int64
	}

	var tx pgx.Tx
	tx, err = pgxConn.Begin(ctx)
	if err != nil {
		z.Error(err.Error())
		return
	}

	defer func() { _ = tx.Rollback(ctx) }()

	var filesMD5 []string

	if fileInfo.Filename == "" {
		err = fmt.Errorf("文件的名称为空")
		z.Error(err.Error())
		return
	}
	if fileInfo.Size == 0 {
		err = fmt.Errorf("文件(%s)的长度为零", fileInfo.Filename)
		z.Error(err.Error())
		return
	}

	if fileInfo.Size > maxUploadFileSize {
		err = fmt.Errorf("文件%s的大小为%4.2f兆，超过了系统的%4.2f兆限制",
			fileInfo.Filename,
			float64(fileInfo.Size)/(1024*1024),
			float64(maxUploadFileSize)/(1024*1024))
		z.Error(err.Error())
		return
	}

	var src multipart.File
	src, err = fileInfo.Open()
	if err != nil {
		err = fmt.Errorf("文件(%s)打开时出错: %s",
			fileInfo.Filename, err.Error())
		z.Error(err.Error())
		return
	}
	defer func() { _ = src.Close() }()

	hash := md5.New()
	_, err = io.Copy(hash, src)
	if err != nil {
		err = fmt.Errorf("生成文件(%s)摘要时出错: %s",
			fileInfo.Filename, err.Error())
		z.Error(err.Error())
		return
	}

	fd.MD5 = null.StringFrom(hex.EncodeToString(hash.Sum(nil)))
	fd.Path = null.StringFrom(fileStorePath)
	fd.Name = null.StringFrom(fileInfo.Filename)
	fd.Size = null.IntFrom(fileInfo.Size)
	fd.SN = null.IntFrom(fileSN)
	filesMD5 = append(filesMD5, fd.MD5.String)

	s := `select exists(select 1 from t_file where digest=$1)`
	r := sqlxDB.QueryRow(s, fd.MD5)
	var exists bool
	err = r.Scan(&exists)
	if err != nil {
		err = fmt.Errorf("文件(%s)(MD5: %s)时出错: %s",
			fileInfo.Filename, fd.MD5.String, err.Error())
		z.Error(err.Error())
		return
	}

	//如果文件在磁盘上不存在，则保存文件到磁盘
	if !exists {
		fn := fileStorePath + fd.MD5.String
		var dst *os.File
		dst, err = os.OpenFile(fn, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			err = fmt.Errorf("后端打开文件(%s)(MD5: %s)时出错: %s",
				fileInfo.Filename, fd.MD5.String, err.Error())
			z.Error(err.Error())
			return
		}
		defer func() { _ = dst.Close() }()

		var ubiety int64
		ubiety, err = src.Seek(0, io.SeekStart)
		if err != nil {
			err = fmt.Errorf("倒带文件(%s)(MD5: %s)时出错: %s",
				fileInfo.Filename, fd.MD5.String, err.Error())
			z.Error(err.Error())
			return
		}
		if ubiety != 0 {
			err = fmt.Errorf("倒带文件(%s)(MD5: %s)时出错: 没倒到头",
				fileInfo.Filename, fd.MD5.String)
			z.Error(err.Error())
			return
		}

		ubiety, err = io.Copy(dst, src)
		if err != nil {
			err = fmt.Errorf("后端保存文件(%s)(MD5: %s)时出错: %s",
				fileInfo.Filename, fd.MD5.String, err.Error())
			z.Error(err.Error())
			return
		}

		if ubiety != fileInfo.Size {
			err = fmt.Errorf("写文件(%s)(MD5: %s)时出错: 没有写完整",
				fileInfo.Filename, fd.MD5.String)
			z.Error(err.Error())
			return
		}

		// 存文件到数据库pg_largeobject
		if fileStoreToDB {
			ubiety, err = src.Seek(0, io.SeekStart)
			if err != nil {
				err = fmt.Errorf("倒带文件(%s)(MD5: %s)时出错: %s",
					fileInfo.Filename, fd.MD5.String, err.Error())
				z.Error(err.Error())
				return
			}

			if ubiety != 0 {
				err = fmt.Errorf("倒带文件(%s)(MD5: %s)时出错: 没倒到头",
					fileInfo.Filename, fd.MD5.String)
				z.Error(err.Error())
				return
			}

			//-------

			lo := tx.LargeObjects()
			var oid uint32
			oid, err = lo.Create(ctx, 0)
			if err != nil {
				err = fmt.Errorf("后端保存文件(%s)(MD5: %s)时出错: %s",
					fileInfo.Filename, fd.MD5.String, err.Error())
				z.Error(err.Error())
				return
			}

			var dbFile *pgx.LargeObject
			dbFile, err = lo.Open(ctx, oid, pgx.LargeObjectModeRead|pgx.LargeObjectModeWrite)
			if err != nil {
				err = fmt.Errorf("后端保存文件(%s)(MD5: %s)时出错: %s",
					fileInfo.Filename, fd.MD5.String, err.Error())
				z.Error(err.Error())
				return
			}
			defer func() { _ = dbFile.Close() }()

			ubiety, err = io.Copy(dbFile, src)
			if err != nil {
				err = fmt.Errorf("后端保存文件(%s)(MD5: %s)时出错: %s",
					fileInfo.Filename, fd.MD5.String, err.Error())
				z.Error(err.Error())
				return
			}

			if ubiety != fileInfo.Size {
				err = fmt.Errorf("写文件(%s)(MD5: %s)时出错: 没有写完整",
					fileInfo.Filename, fd.MD5.String)
				z.Error(err.Error())
				return
			}
			fd.FileOID = null.IntFrom(int64(oid))
		}
	}

	//------------
	var fdExists, pathExists bool
	fdExists, pathExists, err = fileDescUpdate(ctx, fd) //t_file
	if err != nil || (fdExists && pathExists) {
		return
	}

	filesField, err = updateTableField(fd)
	if err != nil {
		return
	}

	err = tx.Commit(ctx)
	if err != nil {
		z.Error(err.Error())
		return
	}

	if filesField == "" {
		var tableField string
		tableField, err = getTableField(fd)
		if err != nil {
			return
		}

		desc := strings.Split(tableField, ".")
		if len(desc) != 2 {
			err = fmt.Errorf("invalid value for key: %s in ownerItemToTable",
				fd.OwnerType.String+"."+fd.Item.String)
			z.Error(err.Error())
			return
		}
		s := fmt.Sprintf(`select %s from %s where id=$1`, desc[1], desc[0])

		var r null.String
		err = pgxConn.QueryRow(ctx, s, fd.LinkID.Int64).Scan(&r)
		if err != nil {
			z.Error(err.Error())
			return
		}

		filesField = r.String
		if filesField == "" {
			err = fmt.Errorf("上传后文件集竟然为空, ownerType: %s,item: %s, linkID: %d",
				fd.OwnerType.String, fd.Item.String, fd.LinkID.Int64)
			z.Error(err.Error())
			return
		}
	}

	for _, v := range filesMD5 {
		for i := 0; ; i++ {
			d := gjson.Get(filesField, fmt.Sprintf("%d.digest", i))
			if !d.Exists() {
				break
			}

			if d.Str != v {
				continue
			}

			filesField, err = sjson.Set(filesField, fmt.Sprintf("%d.new", i), true)
			if err != nil {
				z.Error(err.Error())
				return
			}

			filesField, err = sjson.Set(filesField, fmt.Sprintf("%d.linkID", i), fd.LinkID.Int64)
			if err != nil {
				z.Error(err.Error())
				return
			}

			filesField, err = sjson.Set(filesField, fmt.Sprintf("%d.ownerType", i), fd.OwnerType.String)
			if err != nil {
				z.Error(err.Error())
				return
			}

			filesField, err = sjson.Set(filesField, fmt.Sprintf("%d.item", i), fd.Item.String)
			if err != nil {
				z.Error(err.Error())
				return
			}
			break
		}
	}

	return
}
