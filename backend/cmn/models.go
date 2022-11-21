package cmn

import (
	"database/sql"
	"regexp"

	"github.com/jmoiron/sqlx/types"
	"github.com/pkg/errors"
	"w2w.io/null"
)

/*TAccountOprLog 用户关键信息变更日志，谁，在什么时间变更了数据，变更前数据是什么样子 represents kuser.t_account_opr_log */
type TAccountOprLog struct {
	ID         null.Int       `json:"ID,omitempty" db:"id,true,integer"`                    /* id 操作编号 */
	UserID     null.Int       `json:"UserID,omitempty" db:"user_id,false,bigint"`           /* user_id 被变更用户编号 */
	Original   types.JSONText `json:"Original,omitempty" db:"original,false,jsonb"`         /* original 原账号数据 */
	CreateTime null.Int       `json:"CreateTime,omitempty" db:"create_time,false,bigint"`   /* create_time 生成时间 */
	Creator    null.Int       `json:"Creator,omitempty" db:"creator,false,bigint"`          /* creator 本数据创建者 */
	DomainID   null.Int       `json:"DomainID,omitempty" db:"domain_id,false,bigint"`       /* domain_id 数据隶属 */
	Addi       types.JSONText `json:"Addi,omitempty" db:"addi,false,jsonb"`                 /* addi 附加信息 */
	Remark     null.String    `json:"Remark,omitempty" db:"remark,false,character varying"` /* remark 备注 */
	Filter                    // build DML where clause
}

//TAccountOprLogFields full field list for default query
var TAccountOprLogFields = []string{
	"ID",
	"UserID",
	"Original",
	"CreateTime",
	"Creator",
	"DomainID",
	"Addi",
	"Remark",
}

//Fields return all fields of struct.
func (r *TAccountOprLog) Fields() []string {
	return TAccountOprLogFields
}

//GetTableName return the associated db table name.
func (r *TAccountOprLog) GetTableName() string {
	var viewNamePattern = regexp.MustCompile(`(?i)^t_v_[a-z0-9_]+$`)
	tableName := "t_account_opr_log"
	if viewNamePattern.MatchString(tableName) {
		return tableName[2:]
	}
	return tableName
}

// Create inserts the TAccountOprLog to the database.
func (r *TAccountOprLog) Create(db Queryer) error {
	err := db.QueryRow(
		`INSERT INTO t_account_opr_log (user_id, original, create_time, creator, domain_id, addi, remark) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
		&r.UserID, &r.Original, &r.CreateTime, &r.Creator, &r.DomainID, &r.Addi, &r.Remark).Scan(&r.ID)
	if err != nil {
		return errors.Wrap(err, "failed to insert t_account_opr_log")
	}
	return nil
}

// GetTAccountOprLogByPk select the TAccountOprLog from the database.
func GetTAccountOprLogByPk(db Queryer, pk0 null.Int) (*TAccountOprLog, error) {
	var r TAccountOprLog
	err := db.QueryRow(
		`SELECT id, user_id, original, create_time, creator, domain_id, addi, remark FROM t_account_opr_log WHERE id = $1`,
		pk0).Scan(&r.ID, &r.UserID, &r.Original, &r.CreateTime, &r.Creator, &r.DomainID, &r.Addi, &r.Remark)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select t_account_opr_log")
	}
	return &r, nil
}

/*TAge 年龄表 represents kuser.t_age */
type TAge struct {
	ID              null.Int       `json:"ID,omitempty" db:"id,true,integer"`                             /* id 条目编号 */
	InsuranceTypeID null.Int       `json:"InsuranceTypeID,omitempty" db:"insurance_type_id,false,bigint"` /* insurance_type_id 保险类型 */
	SchoolType      null.String    `json:"SchoolType,omitempty" db:"school_type,false,character varying"` /* school_type 学校类别，可以多个类别，用空格间隔开识别 */
	Province        null.String    `json:"Province,omitempty" db:"province,false,character varying"`      /* province 省 */
	City            null.String    `json:"City,omitempty" db:"city,false,character varying"`              /* city 市 */
	District        null.String    `json:"District,omitempty" db:"district,false,character varying"`      /* district 区/县 */
	Enabled         null.Bool      `json:"Enabled,omitempty" db:"enabled,false,boolean"`                  /* enabled 是否开启限制年龄；true开启限制年龄，false关闭限制年龄 */
	MaleMax         null.Int       `json:"MaleMax,omitempty" db:"male_max,false,smallint"`                /* male_max 男性年龄最大值 */
	MaleMin         null.Int       `json:"MaleMin,omitempty" db:"male_min,false,smallint"`                /* male_min 男性年龄最小值 */
	FemaleMin       null.Int       `json:"FemaleMin,omitempty" db:"female_min,false,smallint"`            /* female_min 女性年龄最小值 */
	FemaleMax       null.Int       `json:"FemaleMax,omitempty" db:"female_max,false,smallint"`            /* female_max 女性年龄最大值 */
	DomainID        null.Int       `json:"DomainID,omitempty" db:"domain_id,false,bigint"`                /* domain_id 数据属主 */
	Addi            types.JSONText `json:"Addi,omitempty" db:"addi,false,jsonb"`                          /* addi 备用字段 */
	Creator         null.Int       `json:"Creator,omitempty" db:"creator,false,bigint"`                   /* creator 创建者用户ID */
	CreateTime      null.Int       `json:"CreateTime,omitempty" db:"create_time,false,bigint"`            /* create_time 创建时间 */
	UpdateTime      null.Int       `json:"UpdateTime,omitempty" db:"update_time,false,bigint"`            /* update_time 更新时间 */
	Updator         null.Int       `json:"Updator,omitempty" db:"updator,false,bigint"`                   /* updator 更新人 */
	Remark          null.String    `json:"Remark,omitempty" db:"remark,false,character varying"`          /* remark 备注 */
	Status          null.String    `json:"Status,omitempty" db:"status,false,character varying"`          /* status 0:有效, 2: 删除 */
	Filter                         // build DML where clause
}

//TAgeFields full field list for default query
var TAgeFields = []string{
	"ID",
	"InsuranceTypeID",
	"SchoolType",
	"Province",
	"City",
	"District",
	"Enabled",
	"MaleMax",
	"MaleMin",
	"FemaleMin",
	"FemaleMax",
	"DomainID",
	"Addi",
	"Creator",
	"CreateTime",
	"UpdateTime",
	"Updator",
	"Remark",
	"Status",
}

//Fields return all fields of struct.
func (r *TAge) Fields() []string {
	return TAgeFields
}

//GetTableName return the associated db table name.
func (r *TAge) GetTableName() string {
	var viewNamePattern = regexp.MustCompile(`(?i)^t_v_[a-z0-9_]+$`)
	tableName := "t_age"
	if viewNamePattern.MatchString(tableName) {
		return tableName[2:]
	}
	return tableName
}

// Create inserts the TAge to the database.
func (r *TAge) Create(db Queryer) error {
	err := db.QueryRow(
		`INSERT INTO t_age (insurance_type_id, school_type, province, city, district, enabled, male_max, male_min, female_min, female_max, domain_id, addi, creator, create_time, update_time, updator, remark, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18) RETURNING id`,
		&r.InsuranceTypeID, &r.SchoolType, &r.Province, &r.City, &r.District, &r.Enabled, &r.MaleMax, &r.MaleMin, &r.FemaleMin, &r.FemaleMax, &r.DomainID, &r.Addi, &r.Creator, &r.CreateTime, &r.UpdateTime, &r.Updator, &r.Remark, &r.Status).Scan(&r.ID)
	if err != nil {
		return errors.Wrap(err, "failed to insert t_age")
	}
	return nil
}

// GetTAgeByPk select the TAge from the database.
func GetTAgeByPk(db Queryer, pk0 null.Int) (*TAge, error) {
	var r TAge
	err := db.QueryRow(
		`SELECT id, insurance_type_id, school_type, province, city, district, enabled, male_max, male_min, female_min, female_max, domain_id, addi, creator, create_time, update_time, updator, remark, status FROM t_age WHERE id = $1`,
		pk0).Scan(&r.ID, &r.InsuranceTypeID, &r.SchoolType, &r.Province, &r.City, &r.District, &r.Enabled, &r.MaleMax, &r.MaleMin, &r.FemaleMin, &r.FemaleMax, &r.DomainID, &r.Addi, &r.Creator, &r.CreateTime, &r.UpdateTime, &r.Updator, &r.Remark, &r.Status)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select t_age")
	}
	return &r, nil
}

/*TAPI 接口信息表 represents kuser.t_api */
type TAPI struct {
	ID                 null.Int    `json:"ID,omitempty" db:"id,true,integer"`                                              /* id 编码 */
	Name               string      `json:"Name,omitempty" db:"name,false,character varying"`                               /* name 接口名称 */
	ExposePath         null.String `json:"ExposePath,omitempty" db:"expose_path,false,character varying"`                  /* expose_path 访问路径 */
	Maintainer         null.Int    `json:"Maintainer,omitempty" db:"maintainer,false,bigint"`                              /* maintainer 维护者 */
	AccessControlLevel string      `json:"AccessControlLevel,omitempty" db:"access_control_level,false,character varying"` /* access_control_level 访问控制实现层级
	level 0: 无组/角色/数据限制
	level 2: 机构#角色级别, 实现了不同角色授权，但不控制数据范围
	level 4: 机构#角色$ID, 实现了不同角色授权，可控制 creator || all
	level 8: 机构.DEPT#角色$ID, 实现了不同角色授权，可控制 creator || GRPs */
	Regenerator null.Int       `json:"Regenerator,omitempty" db:"regenerator,false,bigint"`  /* regenerator 更新者 */
	UpdateTime  null.Int       `json:"UpdateTime,omitempty" db:"update_time,false,bigint"`   /* update_time 帐号信息更新时间 */
	Creator     null.Int       `json:"Creator,omitempty" db:"creator,false,bigint"`          /* creator 本数据创建者 */
	CreateTime  null.Int       `json:"CreateTime,omitempty" db:"create_time,false,bigint"`   /* create_time 生成时间 */
	DomainID    null.Int       `json:"DomainID,omitempty" db:"domain_id,false,bigint"`       /* domain_id 数据隶属 */
	Addi        types.JSONText `json:"Addi,omitempty" db:"addi,false,jsonb"`                 /* addi 附加信息 */
	Remark      null.String    `json:"Remark,omitempty" db:"remark,false,character varying"` /* remark 备注 */
	Status      null.String    `json:"Status,omitempty" db:"status,false,character varying"` /* status 状态，00：草稿，01：有效，02：作废 */
	Filter                     // build DML where clause
}

//TAPIFields full field list for default query
var TAPIFields = []string{
	"ID",
	"Name",
	"ExposePath",
	"Maintainer",
	"AccessControlLevel",
	"Regenerator",
	"UpdateTime",
	"Creator",
	"CreateTime",
	"DomainID",
	"Addi",
	"Remark",
	"Status",
}

//Fields return all fields of struct.
func (r *TAPI) Fields() []string {
	return TAPIFields
}

//GetTableName return the associated db table name.
func (r *TAPI) GetTableName() string {
	var viewNamePattern = regexp.MustCompile(`(?i)^t_v_[a-z0-9_]+$`)
	tableName := "t_api"
	if viewNamePattern.MatchString(tableName) {
		return tableName[2:]
	}
	return tableName
}

// Create inserts the TAPI to the database.
func (r *TAPI) Create(db Queryer) error {
	err := db.QueryRow(
		`INSERT INTO t_api (name, expose_path, maintainer, access_control_level, regenerator, update_time, creator, create_time, domain_id, addi, remark, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id`,
		&r.Name, &r.ExposePath, &r.Maintainer, &r.AccessControlLevel, &r.Regenerator, &r.UpdateTime, &r.Creator, &r.CreateTime, &r.DomainID, &r.Addi, &r.Remark, &r.Status).Scan(&r.ID)
	if err != nil {
		return errors.Wrap(err, "failed to insert t_api")
	}
	return nil
}

// GetTAPIByPk select the TAPI from the database.
func GetTAPIByPk(db Queryer, pk0 null.Int) (*TAPI, error) {
	var r TAPI
	err := db.QueryRow(
		`SELECT id, name, expose_path, maintainer, access_control_level, regenerator, update_time, creator, create_time, domain_id, addi, remark, status FROM t_api WHERE id = $1`,
		pk0).Scan(&r.ID, &r.Name, &r.ExposePath, &r.Maintainer, &r.AccessControlLevel, &r.Regenerator, &r.UpdateTime, &r.Creator, &r.CreateTime, &r.DomainID, &r.Addi, &r.Remark, &r.Status)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select t_api")
	}
	return &r, nil
}

/*TBlacklist 黑名单表 represents kuser.t_blacklist */
type TBlacklist struct {
	ID          null.Int       `json:"ID,omitempty" db:"id,true,integer"`                      /* id 拒保黑名单编号 */
	OrderID     null.Int       `json:"OrderID,omitempty" db:"order_id,false,bigint"`           /* order_id 来源订单号 */
	Type        string         `json:"Type,omitempty" db:"type,false,character varying"`       /* type 黑名单类型（投保人，统一社会信用代码（投保人），投保联系人手机号码） */
	Content     string         `json:"Content,omitempty" db:"content,false,character varying"` /* content 黑名单内容 */
	Regenerator null.Int       `json:"Regenerator,omitempty" db:"regenerator,false,bigint"`    /* regenerator 更新者 */
	UpdateTime  null.Int       `json:"UpdateTime,omitempty" db:"update_time,false,bigint"`     /* update_time 更新时间 */
	Creator     null.Int       `json:"Creator,omitempty" db:"creator,false,bigint"`            /* creator 创建者用户ID */
	CreateTime  null.Int       `json:"CreateTime,omitempty" db:"create_time,false,bigint"`     /* create_time 创建时间 */
	DomainID    null.Int       `json:"DomainID,omitempty" db:"domain_id,false,bigint"`         /* domain_id 数据属主 */
	Addi        types.JSONText `json:"Addi,omitempty" db:"addi,false,jsonb"`                   /* addi 附加数据 */
	Status      null.String    `json:"Status,omitempty" db:"status,false,character varying"`   /* status 状态 0:有效, 2: 无效 */
	Filter                     // build DML where clause
}

//TBlacklistFields full field list for default query
var TBlacklistFields = []string{
	"ID",
	"OrderID",
	"Type",
	"Content",
	"Regenerator",
	"UpdateTime",
	"Creator",
	"CreateTime",
	"DomainID",
	"Addi",
	"Status",
}

//Fields return all fields of struct.
func (r *TBlacklist) Fields() []string {
	return TBlacklistFields
}

//GetTableName return the associated db table name.
func (r *TBlacklist) GetTableName() string {
	var viewNamePattern = regexp.MustCompile(`(?i)^t_v_[a-z0-9_]+$`)
	tableName := "t_blacklist"
	if viewNamePattern.MatchString(tableName) {
		return tableName[2:]
	}
	return tableName
}

// Create inserts the TBlacklist to the database.
func (r *TBlacklist) Create(db Queryer) error {
	err := db.QueryRow(
		`INSERT INTO t_blacklist (order_id, type, content, regenerator, update_time, creator, create_time, domain_id, addi, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`,
		&r.OrderID, &r.Type, &r.Content, &r.Regenerator, &r.UpdateTime, &r.Creator, &r.CreateTime, &r.DomainID, &r.Addi, &r.Status).Scan(&r.ID)
	if err != nil {
		return errors.Wrap(err, "failed to insert t_blacklist")
	}
	return nil
}

// GetTBlacklistByPk select the TBlacklist from the database.
func GetTBlacklistByPk(db Queryer, pk0 null.Int) (*TBlacklist, error) {
	var r TBlacklist
	err := db.QueryRow(
		`SELECT id, order_id, type, content, regenerator, update_time, creator, create_time, domain_id, addi, status FROM t_blacklist WHERE id = $1`,
		pk0).Scan(&r.ID, &r.OrderID, &r.Type, &r.Content, &r.Regenerator, &r.UpdateTime, &r.Creator, &r.CreateTime, &r.DomainID, &r.Addi, &r.Status)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select t_blacklist")
	}
	return &r, nil
}

/*TDegree 知识能力领域等级表 represents kuser.t_degree */
type TDegree struct {
	ID     null.Int    `json:"ID,omitempty" db:"id,true,integer"`                    /* id 编号 */
	Level  null.Int    `json:"Level,omitempty" db:"level,false,integer"`             /* level 等级 */
	Name   null.String `json:"Name,omitempty" db:"name,false,character varying"`     /* name 等级名称 */
	Limn   null.String `json:"Limn,omitempty" db:"limn,false,character varying"`     /* limn 等级描述 */
	Status null.String `json:"Status,omitempty" db:"status,false,character varying"` /* status 可用，禁用 */
	Filter             // build DML where clause
}

//TDegreeFields full field list for default query
var TDegreeFields = []string{
	"ID",
	"Level",
	"Name",
	"Limn",
	"Status",
}

//Fields return all fields of struct.
func (r *TDegree) Fields() []string {
	return TDegreeFields
}

//GetTableName return the associated db table name.
func (r *TDegree) GetTableName() string {
	var viewNamePattern = regexp.MustCompile(`(?i)^t_v_[a-z0-9_]+$`)
	tableName := "t_degree"
	if viewNamePattern.MatchString(tableName) {
		return tableName[2:]
	}
	return tableName
}

// Create inserts the TDegree to the database.
func (r *TDegree) Create(db Queryer) error {
	err := db.QueryRow(
		`INSERT INTO t_degree (level, name, limn, status) VALUES ($1, $2, $3, $4) RETURNING id`,
		&r.Level, &r.Name, &r.Limn, &r.Status).Scan(&r.ID)
	if err != nil {
		return errors.Wrap(err, "failed to insert t_degree")
	}
	return nil
}

// GetTDegreeByPk select the TDegree from the database.
func GetTDegreeByPk(db Queryer, pk0 null.Int) (*TDegree, error) {
	var r TDegree
	err := db.QueryRow(
		`SELECT id, level, name, limn, status FROM t_degree WHERE id = $1`,
		pk0).Scan(&r.ID, &r.Level, &r.Name, &r.Limn, &r.Status)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select t_degree")
	}
	return &r, nil
}

/*TDomain 用户组织结构定义，格式为：机构[部门.科室.组]#角色 represents kuser.t_domain */
type TDomain struct {
	ID          null.Int       `json:"ID,omitempty" db:"id,true,integer"`                    /* id 编码 */
	Name        string         `json:"Name,omitempty" db:"name,false,character varying"`     /* name 域名称 */
	Domain      string         `json:"Domain,omitempty" db:"domain,false,character varying"` /* domain 机构[部门.科室.组]^角色!userID */
	Priority    null.Int       `json:"Priority,omitempty" db:"priority,false,smallint"`      /* priority 0: 超级管理员, 可做任何事; 3: 普通管理员, 可做所属子系统的任何事; 5: 业务员, 可做一般的管理任务; 7: 普通用户, 只能访问自己的数据; 9: 匿名用户, 只能访问白名单功能。 */
	DomainID    null.Int       `json:"DomainID,omitempty" db:"domain_id,false,bigint"`       /* domain_id 数据隶属 */
	Regenerator null.Int       `json:"Regenerator,omitempty" db:"regenerator,false,bigint"`  /* regenerator 更新者 */
	UpdateTime  null.Int       `json:"UpdateTime,omitempty" db:"update_time,false,bigint"`   /* update_time 帐号信息更新时间 */
	Creator     null.Int       `json:"Creator,omitempty" db:"creator,false,bigint"`          /* creator 本数据创建者 */
	CreateTime  null.Int       `json:"CreateTime,omitempty" db:"create_time,false,bigint"`   /* create_time 生成时间 */
	Addi        types.JSONText `json:"Addi,omitempty" db:"addi,false,jsonb"`                 /* addi 附加信息 */
	Remark      null.String    `json:"Remark,omitempty" db:"remark,false,character varying"` /* remark 备注 */
	Status      null.String    `json:"Status,omitempty" db:"status,false,character varying"` /* status 状态，00：草稿，01：有效，02：作废 */
	Filter                     // build DML where clause
}

//TDomainFields full field list for default query
var TDomainFields = []string{
	"ID",
	"Name",
	"Domain",
	"Priority",
	"DomainID",
	"Regenerator",
	"UpdateTime",
	"Creator",
	"CreateTime",
	"Addi",
	"Remark",
	"Status",
}

//Fields return all fields of struct.
func (r *TDomain) Fields() []string {
	return TDomainFields
}

//GetTableName return the associated db table name.
func (r *TDomain) GetTableName() string {
	var viewNamePattern = regexp.MustCompile(`(?i)^t_v_[a-z0-9_]+$`)
	tableName := "t_domain"
	if viewNamePattern.MatchString(tableName) {
		return tableName[2:]
	}
	return tableName
}

// Create inserts the TDomain to the database.
func (r *TDomain) Create(db Queryer) error {
	err := db.QueryRow(
		`INSERT INTO t_domain (name, domain, priority, domain_id, regenerator, update_time, creator, create_time, addi, remark, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id`,
		&r.Name, &r.Domain, &r.Priority, &r.DomainID, &r.Regenerator, &r.UpdateTime, &r.Creator, &r.CreateTime, &r.Addi, &r.Remark, &r.Status).Scan(&r.ID)
	if err != nil {
		return errors.Wrap(err, "failed to insert t_domain")
	}
	return nil
}

// GetTDomainByPk select the TDomain from the database.
func GetTDomainByPk(db Queryer, pk0 null.Int) (*TDomain, error) {
	var r TDomain
	err := db.QueryRow(
		`SELECT id, name, domain, priority, domain_id, regenerator, update_time, creator, create_time, addi, remark, status FROM t_domain WHERE id = $1`,
		pk0).Scan(&r.ID, &r.Name, &r.Domain, &r.Priority, &r.DomainID, &r.Regenerator, &r.UpdateTime, &r.Creator, &r.CreateTime, &r.Addi, &r.Remark, &r.Status)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select t_domain")
	}
	return &r, nil
}

/*TDomainAPI 用户、接口、数据访问控制表 represents kuser.t_domain_api */
type TDomainAPI struct {
	ID          null.Int    `json:"ID,omitempty" db:"id,true,integer"`                               /* id 权限编码 */
	API         null.Int    `json:"API,omitempty" db:"api,false,bigint"`                             /* api 接口/功能编码 */
	Domain      null.Int    `json:"Domain,omitempty" db:"domain,false,bigint"`                       /* domain 组、角色编码 */
	GrantSource null.String `json:"GrantSource,omitempty" db:"grant_source,false,character varying"` /* grant_source grant:数据权限由t_relation中left_type:t_domain.id与left_type:t_user.id获得的数据决定,或data_scope中数据决定，但data_scope与t_relation只能存在一种，如果data_scope有效，则忽略t_relation;

	cousin:忽略data_scope与t_relation, 授权数据由被过虑数据的domain_id决定,即被过虑数据的domain_id 与登录用户的t_user.domain_id相同或级别更低的数据，例如
	    用户的t_user.domain为xkb^admin而数据的domain为xkb.school^admin，则用户可以获得该数据

	self: 被过虑数据的creator 与登录用户的t_user.id相同

	api: 由功能(api)自己决定  */
	DataAccessMode null.String    `json:"DataAccessMode,omitempty" db:"data_access_mode,false,character varying"` /* data_access_mode 数据访问类型, full:可读写, read: 只读, write: 写, partial: 部分写/混合 */
	DataScope      types.JSONText `json:"DataScope,omitempty" db:"data_scope,false,jsonb"`                        /* data_scope 当grant_source是grant时,以json数据方式提供数据授权范围格式为:
	  {"granter":"t_user.id","grantee":"t_school.id","data":[1234,456,789]}
	granter: 代表数据拥有者, t_user.id代表用户, t_domain.id代表角色,t_api.id代表功能
	grantee: 代表拥有的数据,t_school.id代表可以访问的机构列表。
	授权数据如果存储在t_relation中则各项分别对应如下
	granter对应left_type, left_key对应t_user_domain.sys_user或t_domain_api.domain
	grantee对应right_type, right_key对应right_type的意义 */
	DomainID    null.Int       `json:"DomainID,omitempty" db:"domain_id,false,bigint"`       /* domain_id 数据领域归属 */
	Regenerator null.Int       `json:"Regenerator,omitempty" db:"regenerator,false,bigint"`  /* regenerator 更新者 */
	UpdateTime  null.Int       `json:"UpdateTime,omitempty" db:"update_time,false,bigint"`   /* update_time 帐号信息更新时间 */
	Creator     null.Int       `json:"Creator,omitempty" db:"creator,false,bigint"`          /* creator 本数据创建者 */
	CreateTime  null.Int       `json:"CreateTime,omitempty" db:"create_time,false,bigint"`   /* create_time 生成时间 */
	Addi        types.JSONText `json:"Addi,omitempty" db:"addi,false,jsonb"`                 /* addi 附加信息 */
	Remark      null.String    `json:"Remark,omitempty" db:"remark,false,character varying"` /* remark 备注 */
	Status      null.String    `json:"Status,omitempty" db:"status,false,character varying"` /* status 状态，00：草稿，01：有效，02：作废 */
	Filter                     // build DML where clause
}

//TDomainAPIFields full field list for default query
var TDomainAPIFields = []string{
	"ID",
	"API",
	"Domain",
	"GrantSource",
	"DataAccessMode",
	"DataScope",
	"DomainID",
	"Regenerator",
	"UpdateTime",
	"Creator",
	"CreateTime",
	"Addi",
	"Remark",
	"Status",
}

//Fields return all fields of struct.
func (r *TDomainAPI) Fields() []string {
	return TDomainAPIFields
}

//GetTableName return the associated db table name.
func (r *TDomainAPI) GetTableName() string {
	var viewNamePattern = regexp.MustCompile(`(?i)^t_v_[a-z0-9_]+$`)
	tableName := "t_domain_api"
	if viewNamePattern.MatchString(tableName) {
		return tableName[2:]
	}
	return tableName
}

// Create inserts the TDomainAPI to the database.
func (r *TDomainAPI) Create(db Queryer) error {
	err := db.QueryRow(
		`INSERT INTO t_domain_api (api, domain, grant_source, data_access_mode, data_scope, domain_id, regenerator, update_time, creator, create_time, addi, remark, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) RETURNING id`,
		&r.API, &r.Domain, &r.GrantSource, &r.DataAccessMode, &r.DataScope, &r.DomainID, &r.Regenerator, &r.UpdateTime, &r.Creator, &r.CreateTime, &r.Addi, &r.Remark, &r.Status).Scan(&r.ID)
	if err != nil {
		return errors.Wrap(err, "failed to insert t_domain_api")
	}
	return nil
}

// GetTDomainAPIByPk select the TDomainAPI from the database.
func GetTDomainAPIByPk(db Queryer, pk0 null.Int) (*TDomainAPI, error) {
	var r TDomainAPI
	err := db.QueryRow(
		`SELECT id, api, domain, grant_source, data_access_mode, data_scope, domain_id, regenerator, update_time, creator, create_time, addi, remark, status FROM t_domain_api WHERE id = $1`,
		pk0).Scan(&r.ID, &r.API, &r.Domain, &r.GrantSource, &r.DataAccessMode, &r.DataScope, &r.DomainID, &r.Regenerator, &r.UpdateTime, &r.Creator, &r.CreateTime, &r.Addi, &r.Remark, &r.Status)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select t_domain_api")
	}
	return &r, nil
}

/*TExpertise 知识能力领域表 represents kuser.t_expertise */
type TExpertise struct {
	ID         null.Int    `json:"ID,omitempty" db:"id,true,integer"`                    /* id 编号 */
	Belongto   null.Int    `json:"Belongto,omitempty" db:"belongto,false,bigint"`        /* belongto 上级expertise */
	Name       null.String `json:"Name,omitempty" db:"name,false,character varying"`     /* name 知识能力领域名称 */
	Limn       null.String `json:"Limn,omitempty" db:"limn,false,character varying"`     /* limn 描述 */
	Creator    null.Int    `json:"Creator,omitempty" db:"creator,false,bigint"`          /* creator 创建者 */
	CreateTime null.Int    `json:"CreateTime,omitempty" db:"create_time,false,bigint"`   /* create_time 创建时间 */
	UpdateTime null.Int    `json:"UpdateTime,omitempty" db:"update_time,false,bigint"`   /* update_time 更新时间 */
	Status     null.String `json:"Status,omitempty" db:"status,false,character varying"` /* status 可用，禁用 */
	Filter                 // build DML where clause
}

//TExpertiseFields full field list for default query
var TExpertiseFields = []string{
	"ID",
	"Belongto",
	"Name",
	"Limn",
	"Creator",
	"CreateTime",
	"UpdateTime",
	"Status",
}

//Fields return all fields of struct.
func (r *TExpertise) Fields() []string {
	return TExpertiseFields
}

//GetTableName return the associated db table name.
func (r *TExpertise) GetTableName() string {
	var viewNamePattern = regexp.MustCompile(`(?i)^t_v_[a-z0-9_]+$`)
	tableName := "t_expertise"
	if viewNamePattern.MatchString(tableName) {
		return tableName[2:]
	}
	return tableName
}

// Create inserts the TExpertise to the database.
func (r *TExpertise) Create(db Queryer) error {
	err := db.QueryRow(
		`INSERT INTO t_expertise (belongto, name, limn, creator, create_time, update_time, status) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
		&r.Belongto, &r.Name, &r.Limn, &r.Creator, &r.CreateTime, &r.UpdateTime, &r.Status).Scan(&r.ID)
	if err != nil {
		return errors.Wrap(err, "failed to insert t_expertise")
	}
	return nil
}

// GetTExpertiseByPk select the TExpertise from the database.
func GetTExpertiseByPk(db Queryer, pk0 null.Int) (*TExpertise, error) {
	var r TExpertise
	err := db.QueryRow(
		`SELECT id, belongto, name, limn, creator, create_time, update_time, status FROM t_expertise WHERE id = $1`,
		pk0).Scan(&r.ID, &r.Belongto, &r.Name, &r.Limn, &r.Creator, &r.CreateTime, &r.UpdateTime, &r.Status)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select t_expertise")
	}
	return &r, nil
}

/*TExternalDomainConf 外部系统访问标识 represents kuser.t_external_domain_conf */
type TExternalDomainConf struct {
	ID      null.Int       `json:"ID,omitempty" db:"id,true,integer"`                       /* id 编号 */
	AppID   string         `json:"AppID,omitempty" db:"app_id,false,character varying"`     /* app_id 外部应用标识，如微信公众号appID */
	AppType string         `json:"AppType,omitempty" db:"app_type,false,character varying"` /* app_type wx_mp: 微信公众号, wx_open: 微信开放平台, ali: 阿里, ui: 联保 */
	AppName string         `json:"AppName,omitempty" db:"app_name,false,character varying"` /* app_name 应用名称，如校快保2019，近邻科技 */
	Tokens  types.JSONText `json:"Tokens,omitempty" db:"tokens,false,jsonb"`                /* tokens 例如，联保：{"appID":"xkbtest",	"appSecret":"123456",
		"branchID":"ba331eb1851d4d8bb5e838dfbf9e09d7",
		"userID":"e1d9441f16284d99a8c2732aedca5753"
	} */
	Regenerator null.Int       `json:"Regenerator,omitempty" db:"regenerator,false,bigint"`  /* regenerator 更新者 */
	UpdateTime  null.Int       `json:"UpdateTime,omitempty" db:"update_time,false,bigint"`   /* update_time 帐号信息更新时间 */
	Creator     null.Int       `json:"Creator,omitempty" db:"creator,false,bigint"`          /* creator 本数据创建者 */
	CreateTime  null.Int       `json:"CreateTime,omitempty" db:"create_time,false,bigint"`   /* create_time 生成时间 */
	DomainID    null.Int       `json:"DomainID,omitempty" db:"domain_id,false,bigint"`       /* domain_id 数据隶属 */
	Addi        types.JSONText `json:"Addi,omitempty" db:"addi,false,jsonb"`                 /* addi 附加信息 */
	Remark      null.String    `json:"Remark,omitempty" db:"remark,false,character varying"` /* remark 备注 */
	Status      null.String    `json:"Status,omitempty" db:"status,false,character varying"` /* status 状态，00：草稿，02：有效，04: 停用，06：作废 */
	Filter                     // build DML where clause
}

//TExternalDomainConfFields full field list for default query
var TExternalDomainConfFields = []string{
	"ID",
	"AppID",
	"AppType",
	"AppName",
	"Tokens",
	"Regenerator",
	"UpdateTime",
	"Creator",
	"CreateTime",
	"DomainID",
	"Addi",
	"Remark",
	"Status",
}

//Fields return all fields of struct.
func (r *TExternalDomainConf) Fields() []string {
	return TExternalDomainConfFields
}

//GetTableName return the associated db table name.
func (r *TExternalDomainConf) GetTableName() string {
	var viewNamePattern = regexp.MustCompile(`(?i)^t_v_[a-z0-9_]+$`)
	tableName := "t_external_domain_conf"
	if viewNamePattern.MatchString(tableName) {
		return tableName[2:]
	}
	return tableName
}

// Create inserts the TExternalDomainConf to the database.
func (r *TExternalDomainConf) Create(db Queryer) error {
	err := db.QueryRow(
		`INSERT INTO t_external_domain_conf (app_id, app_type, app_name, tokens, regenerator, update_time, creator, create_time, domain_id, addi, remark, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id`,
		&r.AppID, &r.AppType, &r.AppName, &r.Tokens, &r.Regenerator, &r.UpdateTime, &r.Creator, &r.CreateTime, &r.DomainID, &r.Addi, &r.Remark, &r.Status).Scan(&r.ID)
	if err != nil {
		return errors.Wrap(err, "failed to insert t_external_domain_conf")
	}
	return nil
}

// GetTExternalDomainConfByPk select the TExternalDomainConf from the database.
func GetTExternalDomainConfByPk(db Queryer, pk0 null.Int) (*TExternalDomainConf, error) {
	var r TExternalDomainConf
	err := db.QueryRow(
		`SELECT id, app_id, app_type, app_name, tokens, regenerator, update_time, creator, create_time, domain_id, addi, remark, status FROM t_external_domain_conf WHERE id = $1`,
		pk0).Scan(&r.ID, &r.AppID, &r.AppType, &r.AppName, &r.Tokens, &r.Regenerator, &r.UpdateTime, &r.Creator, &r.CreateTime, &r.DomainID, &r.Addi, &r.Remark, &r.Status)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select t_external_domain_conf")
	}
	return &r, nil
}

/*TExternalDomainUser 第三方平台用户标识 represents kuser.t_external_domain_user */
type TExternalDomainUser struct {
	ID                null.Int       `json:"ID,omitempty" db:"id,true,integer"`                                             /* id 编号 */
	UserID            null.Int       `json:"UserID,omitempty" db:"user_id,false,bigint"`                                    /* user_id 系统用户编号 */
	BusinessDomainID  string         `json:"BusinessDomainID,omitempty" db:"business_domain_id,false,character varying"`    /* business_domain_id 业务域甲方，如微信公众号appID，广州大学财务系统 */
	UserDomainID      string         `json:"UserDomainID,omitempty" db:"user_domain_id,false,character varying"`            /* user_domain_id 业务域乙方，如微信公众号openID，广州大学财务系统用户账号 */
	UserDomainUnionID null.String    `json:"UserDomainUnionID,omitempty" db:"user_domain_union_id,false,character varying"` /* user_domain_union_id 业务域乙方唯一标识，如微信unionID，广州大学教工编号/学号 */
	ApplyTo           null.String    `json:"ApplyTo,omitempty" db:"apply_to,false,character varying"`                       /* apply_to 该ID用途，如用于支付，标识用户 */
	DomainType        string         `json:"DomainType,omitempty" db:"domain_type,false,character varying"`                 /* domain_type wx_mp: 微信公众号, wx_open: 微信开放平台, ali: 阿里, ui: 联保 */
	Creator           null.Int       `json:"Creator,omitempty" db:"creator,false,bigint"`                                   /* creator 本数据创建者 */
	CreateTime        null.Int       `json:"CreateTime,omitempty" db:"create_time,false,bigint"`                            /* create_time 生成时间 */
	DomainID          null.Int       `json:"DomainID,omitempty" db:"domain_id,false,bigint"`                                /* domain_id 数据隶属 */
	Addi              types.JSONText `json:"Addi,omitempty" db:"addi,false,jsonb"`                                          /* addi 附加信息 */
	Status            null.String    `json:"Status,omitempty" db:"status,false,character varying"`                          /* status 状态，00：草稿，02：有效，04：禁用，06：作废 */
	Filter                           // build DML where clause
}

//TExternalDomainUserFields full field list for default query
var TExternalDomainUserFields = []string{
	"ID",
	"UserID",
	"BusinessDomainID",
	"UserDomainID",
	"UserDomainUnionID",
	"ApplyTo",
	"DomainType",
	"Creator",
	"CreateTime",
	"DomainID",
	"Addi",
	"Status",
}

//Fields return all fields of struct.
func (r *TExternalDomainUser) Fields() []string {
	return TExternalDomainUserFields
}

//GetTableName return the associated db table name.
func (r *TExternalDomainUser) GetTableName() string {
	var viewNamePattern = regexp.MustCompile(`(?i)^t_v_[a-z0-9_]+$`)
	tableName := "t_external_domain_user"
	if viewNamePattern.MatchString(tableName) {
		return tableName[2:]
	}
	return tableName
}

// Create inserts the TExternalDomainUser to the database.
func (r *TExternalDomainUser) Create(db Queryer) error {
	err := db.QueryRow(
		`INSERT INTO t_external_domain_user (user_id, business_domain_id, user_domain_id, user_domain_union_id, apply_to, domain_type, creator, create_time, domain_id, addi, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id`,
		&r.UserID, &r.BusinessDomainID, &r.UserDomainID, &r.UserDomainUnionID, &r.ApplyTo, &r.DomainType, &r.Creator, &r.CreateTime, &r.DomainID, &r.Addi, &r.Status).Scan(&r.ID)
	if err != nil {
		return errors.Wrap(err, "failed to insert t_external_domain_user")
	}
	return nil
}

// GetTExternalDomainUserByPk select the TExternalDomainUser from the database.
func GetTExternalDomainUserByPk(db Queryer, pk0 null.Int) (*TExternalDomainUser, error) {
	var r TExternalDomainUser
	err := db.QueryRow(
		`SELECT id, user_id, business_domain_id, user_domain_id, user_domain_union_id, apply_to, domain_type, creator, create_time, domain_id, addi, status FROM t_external_domain_user WHERE id = $1`,
		pk0).Scan(&r.ID, &r.UserID, &r.BusinessDomainID, &r.UserDomainID, &r.UserDomainUnionID, &r.ApplyTo, &r.DomainType, &r.Creator, &r.CreateTime, &r.DomainID, &r.Addi, &r.Status)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select t_external_domain_user")
	}
	return &r, nil
}

/*TFile 文件描述表 represents kuser.t_file */
type TFile struct {
	ID           null.Int       `json:"ID,omitempty" db:"id,true,integer"`                                 /* id 编号 */
	FileOid      null.Int       `json:"FileOid,omitempty" db:"file_oid,false,oid"`                         /* file_oid 文件数据库对象编号 */
	FileName     string         `json:"FileName,omitempty" db:"file_name,false,character varying"`         /* file_name 文件名 */
	Path         string         `json:"Path,omitempty" db:"path,false,character varying"`                  /* path 服务端存储文件路径 */
	BelongtoPath string         `json:"BelongtoPath,omitempty" db:"belongto_path,false,character varying"` /* belongto_path 分类(以路径方式表述) */
	Digest       string         `json:"Digest,omitempty" db:"digest,false,character varying"`              /* digest sha512 digest */
	Size         null.Int       `json:"Size,omitempty" db:"size,false,bigint"`                             /* size 文件大小 */
	CreateTime   null.Int       `json:"CreateTime,omitempty" db:"create_time,false,bigint"`                /* create_time 创建时间 */
	Creator      null.Int       `json:"Creator,omitempty" db:"creator,false,bigint"`                       /* creator 上传者 */
	DomainID     null.Int       `json:"DomainID,omitempty" db:"domain_id,false,bigint"`                    /* domain_id 数据隶属 */
	Count        null.Int       `json:"Count,omitempty" db:"count,false,integer"`                          /* count 文件引用计数 */
	Belongto     null.Int       `json:"Belongto,omitempty" db:"belongto,false,bigint"`                     /* belongto 隶属的对象编号 */
	Limn         null.String    `json:"Limn,omitempty" db:"limn,false,character varying"`                  /* limn 文件作用描述 */
	OriginPath   null.String    `json:"OriginPath,omitempty" db:"origin_path,false,character varying"`     /* origin_path 用户上传文件路径 */
	OriginName   null.String    `json:"OriginName,omitempty" db:"origin_name,false,character varying"`     /* origin_name 用户上传文件名 */
	Addi         types.JSONText `json:"Addi,omitempty" db:"addi,false,jsonb"`                              /* addi 附加信息 */
	Status       null.String    `json:"Status,omitempty" db:"status,false,character varying"`              /* status 0:有效, 2: 丢失 */
	Filter                      // build DML where clause
}

//TFileFields full field list for default query
var TFileFields = []string{
	"ID",
	"FileOid",
	"FileName",
	"Path",
	"BelongtoPath",
	"Digest",
	"Size",
	"CreateTime",
	"Creator",
	"DomainID",
	"Count",
	"Belongto",
	"Limn",
	"OriginPath",
	"OriginName",
	"Addi",
	"Status",
}

//Fields return all fields of struct.
func (r *TFile) Fields() []string {
	return TFileFields
}

//GetTableName return the associated db table name.
func (r *TFile) GetTableName() string {
	var viewNamePattern = regexp.MustCompile(`(?i)^t_v_[a-z0-9_]+$`)
	tableName := "t_file"
	if viewNamePattern.MatchString(tableName) {
		return tableName[2:]
	}
	return tableName
}

// Create inserts the TFile to the database.
func (r *TFile) Create(db Queryer) error {
	err := db.QueryRow(
		`INSERT INTO t_file (file_oid, file_name, path, belongto_path, digest, size, create_time, creator, domain_id, count, belongto, limn, origin_path, origin_name, addi, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16) RETURNING id`,
		&r.FileOid, &r.FileName, &r.Path, &r.BelongtoPath, &r.Digest, &r.Size, &r.CreateTime, &r.Creator, &r.DomainID, &r.Count, &r.Belongto, &r.Limn, &r.OriginPath, &r.OriginName, &r.Addi, &r.Status).Scan(&r.ID)
	if err != nil {
		return errors.Wrap(err, "failed to insert t_file")
	}
	return nil
}

// GetTFileByPk select the TFile from the database.
func GetTFileByPk(db Queryer, pk0 null.Int) (*TFile, error) {
	var r TFile
	err := db.QueryRow(
		`SELECT id, file_oid, file_name, path, belongto_path, digest, size, create_time, creator, domain_id, count, belongto, limn, origin_path, origin_name, addi, status FROM t_file WHERE id = $1`,
		pk0).Scan(&r.ID, &r.FileOid, &r.FileName, &r.Path, &r.BelongtoPath, &r.Digest, &r.Size, &r.CreateTime, &r.Creator, &r.DomainID, &r.Count, &r.Belongto, &r.Limn, &r.OriginPath, &r.OriginName, &r.Addi, &r.Status)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select t_file")
	}
	return &r, nil
}

/*TImportData excel导入表 represents kuser.t_import_data */
type TImportData struct {
	ID          null.Int       `json:"ID,omitempty" db:"id,true,integer"`                             /* id 参数编号 */
	Name        null.String    `json:"Name,omitempty" db:"name,false,character varying"`              /* name 导入数据名称 */
	Category    string         `json:"Category,omitempty" db:"category,false,character varying"`      /* category 导入数据类型 */
	Key         string         `json:"Key,omitempty" db:"key,false,character varying"`                /* key 数据唯一标识, 由struct中的key简单连接组成 */
	EntityID    null.String    `json:"EntityID,omitempty" db:"entity_id,false,character varying"`     /* entity_id 表示本条数据的逻辑标识，如，保单号，身份证号 */
	Struct      types.JSONText `json:"Struct,omitempty" db:"struct,false,jsonb"`                      /* struct excel数据结构 */
	Base        types.JSONText `json:"Base,omitempty" db:"base,false,jsonb"`                          /* base 表中的非重复信息 */
	Data        types.JSONText `json:"Data,omitempty" db:"data,false,jsonb"`                          /* data 数据 */
	File        types.JSONText `json:"File,omitempty" db:"file,false,jsonb"`                          /* file 导入文件信息 */
	FileDigest  string         `json:"FileDigest,omitempty" db:"file_digest,false,character varying"` /* file_digest file_digest */
	DomainID    null.Int       `json:"DomainID,omitempty" db:"domain_id,false,bigint"`                /* domain_id 数据隶属 */
	CreateTime  null.Int       `json:"CreateTime,omitempty" db:"create_time,false,bigint"`            /* create_time 生成时间 */
	Creator     null.Int       `json:"Creator,omitempty" db:"creator,false,bigint"`                   /* creator 本数据创建者 */
	Regenerator null.Int       `json:"Regenerator,omitempty" db:"regenerator,false,bigint"`           /* regenerator 更新者 */
	UpdateTime  null.Int       `json:"UpdateTime,omitempty" db:"update_time,false,bigint"`            /* update_time 帐号信息更新时间 */
	Addi        types.JSONText `json:"Addi,omitempty" db:"addi,false,jsonb"`                          /* addi 附加信息 */
	Remark      null.String    `json:"Remark,omitempty" db:"remark,false,character varying"`          /* remark 备注 */
	Status      null.String    `json:"Status,omitempty" db:"status,false,character varying"`          /* status 状态，00：草稿，01：有效，02：作废 */
	Filter                     // build DML where clause
}

//TImportDataFields full field list for default query
var TImportDataFields = []string{
	"ID",
	"Name",
	"Category",
	"Key",
	"EntityID",
	"Struct",
	"Base",
	"Data",
	"File",
	"FileDigest",
	"DomainID",
	"CreateTime",
	"Creator",
	"Regenerator",
	"UpdateTime",
	"Addi",
	"Remark",
	"Status",
}

//Fields return all fields of struct.
func (r *TImportData) Fields() []string {
	return TImportDataFields
}

//GetTableName return the associated db table name.
func (r *TImportData) GetTableName() string {
	var viewNamePattern = regexp.MustCompile(`(?i)^t_v_[a-z0-9_]+$`)
	tableName := "t_import_data"
	if viewNamePattern.MatchString(tableName) {
		return tableName[2:]
	}
	return tableName
}

// Create inserts the TImportData to the database.
func (r *TImportData) Create(db Queryer) error {
	err := db.QueryRow(
		`INSERT INTO t_import_data (name, category, key, entity_id, struct, base, data, file, file_digest, domain_id, create_time, creator, regenerator, update_time, addi, remark, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17) RETURNING id`,
		&r.Name, &r.Category, &r.Key, &r.EntityID, &r.Struct, &r.Base, &r.Data, &r.File, &r.FileDigest, &r.DomainID, &r.CreateTime, &r.Creator, &r.Regenerator, &r.UpdateTime, &r.Addi, &r.Remark, &r.Status).Scan(&r.ID)
	if err != nil {
		return errors.Wrap(err, "failed to insert t_import_data")
	}
	return nil
}

// GetTImportDataByPk select the TImportData from the database.
func GetTImportDataByPk(db Queryer, pk0 null.Int) (*TImportData, error) {
	var r TImportData
	err := db.QueryRow(
		`SELECT id, name, category, key, entity_id, struct, base, data, file, file_digest, domain_id, create_time, creator, regenerator, update_time, addi, remark, status FROM t_import_data WHERE id = $1`,
		pk0).Scan(&r.ID, &r.Name, &r.Category, &r.Key, &r.EntityID, &r.Struct, &r.Base, &r.Data, &r.File, &r.FileDigest, &r.DomainID, &r.CreateTime, &r.Creator, &r.Regenerator, &r.UpdateTime, &r.Addi, &r.Remark, &r.Status)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select t_import_data")
	}
	return &r, nil
}

/*TInsurancePolicy 保险单 represents kuser.t_insurance_policy */
type TInsurancePolicy struct {
	ID                      null.Int       `json:"ID,omitempty" db:"id,true,integer"`                                              /* id 编号 */
	Sn                      null.String    `json:"Sn,omitempty" db:"sn,false,character varying"`                                   /* sn 保单号 */
	SnCreator               null.Int       `json:"SnCreator,omitempty" db:"sn_creator,false,bigint"`                               /* sn_creator 保单号上传者ID */
	Name                    string         `json:"Name,omitempty" db:"name,false,character varying"`                               /* name 保险种类 */
	OrderID                 null.Int       `json:"OrderID,omitempty" db:"order_id,false,bigint"`                                   /* order_id 订单编号 */
	Policy                  string         `json:"Policy,omitempty" db:"policy,false,character varying"`                           /* policy 保险合同条款 */
	Start                   null.Int       `json:"Start,omitempty" db:"start,false,bigint"`                                        /* start 起保时间 */
	Cease                   null.Int       `json:"Cease,omitempty" db:"cease,false,bigint"`                                        /* cease 终保时间 */
	Year                    null.Int       `json:"Year,omitempty" db:"year,false,smallint"`                                        /* year 保单年份 */
	Duration                null.Int       `json:"Duration,omitempty" db:"duration,false,bigint"`                                  /* duration 保障期限 */
	Premium                 null.Float     `json:"Premium,omitempty" db:"premium,false,double precision"`                          /* premium 保费金额 */
	ThirdPartyPremium       null.Float     `json:"ThirdPartyPremium,omitempty" db:"third_party_premium,false,double precision"`    /* third_party_premium 第三方保费金额 */
	ThirdPartyAccount       null.String    `json:"ThirdPartyAccount,omitempty" db:"third_party_account,false,character varying"`   /* third_party_account 自动录单账号 */
	PayTime                 null.Int       `json:"PayTime,omitempty" db:"pay_time,false,bigint"`                                   /* pay_time 支付时间 */
	PayChannel              null.String    `json:"PayChannel,omitempty" db:"pay_channel,false,character varying"`                  /* pay_channel 校快保，泰合，近邻，人保，太平洋保险 */
	PayType                 null.String    `json:"PayType,omitempty" db:"pay_type,false,character varying"`                        /* pay_type 支付方式: 对公转账/在线支付/线下支付 */
	UnitPrice               null.Float     `json:"UnitPrice,omitempty" db:"unit_price,false,double precision"`                     /* unit_price 单价 */
	OrgID                   null.Int       `json:"OrgID,omitempty" db:"org_id,false,bigint"`                                       /* org_id 关联机构编号 */
	OrgManagerID            null.Int       `json:"OrgManagerID,omitempty" db:"org_manager_id,false,bigint"`                        /* org_manager_id 关联机构管理人 */
	PolicyholderType        null.String    `json:"PolicyholderType,omitempty" db:"policyholder_type,false,character varying"`      /* policyholder_type 投保人类型 */
	Policyholder            types.JSONText `json:"Policyholder,omitempty" db:"policyholder,false,jsonb"`                           /* policyholder 投保人 */
	PolicyholderID          null.Int       `json:"PolicyholderID,omitempty" db:"policyholder_id,false,bigint"`                     /* policyholder_id 投保人编码 */
	InsuranceType           null.String    `json:"InsuranceType,omitempty" db:"insurance_type,false,character varying"`            /* insurance_type 保险类型: 学生意外伤害险，活动/比赛险(旅游险),食品卫生责任险，教工责任险,校方责任险,实习生责任险,校车责任险,游泳池责任险 */
	InsuranceTypeID         null.Int       `json:"InsuranceTypeID,omitempty" db:"insurance_type_id,false,bigint"`                  /* insurance_type_id 保险产品编码 */
	PolicyScheme            types.JSONText `json:"PolicyScheme,omitempty" db:"policy_scheme,false,jsonb"`                          /* policy_scheme 保险方案 */
	ActivityName            null.String    `json:"ActivityName,omitempty" db:"activity_name,false,character varying"`              /* activity_name 活动名称 */
	ActivityCategory        null.String    `json:"ActivityCategory,omitempty" db:"activity_category,false,character varying"`      /* activity_category 活动类型 */
	ActivityDesc            null.String    `json:"ActivityDesc,omitempty" db:"activity_desc,false,character varying"`              /* activity_desc 活动描述 */
	ActivityLocation        null.String    `json:"ActivityLocation,omitempty" db:"activity_location,false,character varying"`      /* activity_location 活动地点 */
	ActivityDateSet         null.String    `json:"ActivityDateSet,omitempty" db:"activity_date_set,false,character varying"`       /* activity_date_set 具体活动日期，英文逗号隔开 */
	InsuredCount            null.Int       `json:"InsuredCount,omitempty" db:"insured_count,false,smallint"`                       /* insured_count 总数量/保障人数/车辆数 */
	CompulsoryStudentNum    null.Int       `json:"CompulsoryStudentNum,omitempty" db:"compulsory_student_num,false,bigint"`        /* compulsory_student_num 义务教育学生人数（校方） */
	NonCompulsoryStudentNum null.Int       `json:"NonCompulsoryStudentNum,omitempty" db:"non_compulsory_student_num,false,bigint"` /* non_compulsory_student_num 非义务教育人数（校方） */
	Contact                 types.JSONText `json:"Contact,omitempty" db:"contact,false,jsonb"`                                     /* contact 联系人 */
	FeeScheme               types.JSONText `json:"FeeScheme,omitempty" db:"fee_scheme,false,jsonb"`                                /* fee_scheme 计费标准/单价 */
	CarServiceTarget        null.String    `json:"CarServiceTarget,omitempty" db:"car_service_target,false,character varying"`     /* car_service_target 校车服务对象 */
	Same                    null.Bool      `json:"Same,omitempty" db:"same,false,boolean"`                                         /* same 投保人与被保险人是同一人 */
	Relation                null.String    `json:"Relation,omitempty" db:"relation,false,character varying"`                       /* relation 投保人与被保险人关系 */
	Insured                 types.JSONText `json:"Insured,omitempty" db:"insured,false,jsonb"`                                     /* insured 被保险人 */
	InsuredID               null.Int       `json:"InsuredID,omitempty" db:"insured_id,false,bigint"`                               /* insured_id 被保险人编号 */
	HaveInsuredList         null.Bool      `json:"HaveInsuredList,omitempty" db:"have_insured_list,false,boolean"`                 /* have_insured_list 有被保险对象清单 */
	InsuredGroupByDay       null.Bool      `json:"InsuredGroupByDay,omitempty" db:"insured_group_by_day,false,boolean"`            /* insured_group_by_day 被保险对象按日期分组 */
	InsuredType             null.String    `json:"InsuredType,omitempty" db:"insured_type,false,character varying"`                /* insured_type 被保险人类型: 学生，非学生 */
	InsuredList             types.JSONText `json:"InsuredList,omitempty" db:"insured_list,false,jsonb"`                            /* insured_list 被保险对象清单 */
	Indate                  null.Int       `json:"Indate,omitempty" db:"indate,false,bigint"`                                      /* indate 有效期(天) */
	Jurisdiction            null.String    `json:"Jurisdiction,omitempty" db:"jurisdiction,false,character varying"`               /* jurisdiction 司法管辖权 */
	DisputeHandling         null.String    `json:"DisputeHandling,omitempty" db:"dispute_handling,false,character varying"`        /* dispute_handling 争议处理 */
	PrevPolicyNo            null.String    `json:"PrevPolicyNo,omitempty" db:"prev_policy_no,false,character varying"`             /* prev_policy_no 续保保单号 */
	InsureBase              null.String    `json:"InsureBase,omitempty" db:"insure_base,false,character varying"`                  /* insure_base 承保基础 */
	BlanketInsureCode       null.String    `json:"BlanketInsureCode,omitempty" db:"blanket_insure_code,false,character varying"`   /* blanket_insure_code 统保代码 */
	CustomType              null.String    `json:"CustomType,omitempty" db:"custom_type,false,character varying"`                  /* custom_type 场地使用性质:internal, open, both */
	TrainProjects           null.String    `json:"TrainProjects,omitempty" db:"train_projects,false,character varying"`            /* train_projects 训练项目 */
	BusinessLocations       types.JSONText `json:"BusinessLocations,omitempty" db:"business_locations,false,jsonb"`                /* business_locations 承保地址/区域范围/游泳池场地地址 */
	ArbitralAgency          null.String    `json:"ArbitralAgency,omitempty" db:"arbitral_agency,false,character varying"`          /* arbitral_agency 仲裁机构 */
	PoolNum                 null.Int       `json:"PoolNum,omitempty" db:"pool_num,false,smallint"`                                 /* pool_num 游泳池个数 */
	OpenPoolNum             null.Int       `json:"OpenPoolNum,omitempty" db:"open_pool_num,false,smallint"`                        /* open_pool_num 对外开放游泳池数量 */
	HeatedPoolNum           null.Int       `json:"HeatedPoolNum,omitempty" db:"heated_pool_num,false,smallint"`                    /* heated_pool_num 恒温游泳池数量 */
	TrainingPoolNum         null.Int       `json:"TrainingPoolNum,omitempty" db:"training_pool_num,false,smallint"`                /* training_pool_num 培训游泳池数量 */
	InnerArea               null.Float     `json:"InnerArea,omitempty" db:"inner_area,false,double precision"`                     /* inner_area 室内面积 */
	OuterArea               null.Float     `json:"OuterArea,omitempty" db:"outer_area,false,double precision"`                     /* outer_area 室外面积 */
	PoolName                null.String    `json:"PoolName,omitempty" db:"pool_name,false,character varying"`                      /* pool_name 游泳池名称(英文逗号分隔)  */
	HaveDinnerNum           null.Bool      `json:"HaveDinnerNum,omitempty" db:"have_dinner_num,false,boolean"`                     /* have_dinner_num 是否开启就餐人数 */
	DinnerNum               null.Int       `json:"DinnerNum,omitempty" db:"dinner_num,false,integer"`                              /* dinner_num 用餐人数 */
	CanteenNum              null.Int       `json:"CanteenNum,omitempty" db:"canteen_num,false,integer"`                            /* canteen_num 食堂个数 */
	ShopNum                 null.Int       `json:"ShopNum,omitempty" db:"shop_num,false,integer"`                                  /* shop_num 商店个数 */
	HaveRides               null.Bool      `json:"HaveRides,omitempty" db:"have_rides,false,boolean"`                              /* have_rides 营业场所是否有游泳池外游乐设施、机械性游乐设施等 */
	HaveExplosive           null.Bool      `json:"HaveExplosive,omitempty" db:"have_explosive,false,boolean"`                      /* have_explosive 营业场所是否有制造、销售、储存易燃易爆危险品 */
	Area                    null.Int       `json:"Area,omitempty" db:"area,false,integer"`                                         /* area 营业场所总面积（平方米） */
	TrafficNum              null.Int       `json:"TrafficNum,omitempty" db:"traffic_num,false,integer"`                            /* traffic_num 每日客流量（人） */
	TemperatureType         null.String    `json:"TemperatureType,omitempty" db:"temperature_type,false,character varying"`        /* temperature_type 泳池性质:恒温、常温 */
	IsIndoor                null.String    `json:"IsIndoor,omitempty" db:"is_indoor,false,character varying"`                      /* is_indoor 泳池特性:室内、室外 */
	Extra                   types.JSONText `json:"Extra,omitempty" db:"extra,false,jsonb"`                                         /* extra 附加信息:
	附加条款
	企业经营描述
	相关保险情况
	保险公司提示
	保险销售事项确认书
	保险公司信息：经办人/工号、代理点代码、展业方式
	产险销售人员：姓名、职业证号 */
	BankAccount      types.JSONText `json:"BankAccount,omitempty" db:"bank_account,false,jsonb"`                        /* bank_account 对公帐号信息：户名、所在银行、账号 */
	PayContact       null.String    `json:"PayContact,omitempty" db:"pay_contact,false,character varying"`              /* pay_contact 线下支付联系人：微信二维码，base64 */
	HaveSuddenDeath  null.Bool      `json:"HaveSuddenDeath,omitempty" db:"have_sudden_death,false,boolean"`             /* have_sudden_death 是否开启猝死责任险 */
	SuddenDeathTerms null.String    `json:"SuddenDeathTerms,omitempty" db:"sudden_death_terms,false,character varying"` /* sudden_death_terms 猝死条款内容：附加猝死保险责任每人限额5万元，累计限额5万元。附加猝死责任保险条款（经法院判决、仲裁机构裁决或根据县级以上政府及县级以上政府有关部门的行政决定书或者调解证明等材料，需由被保险人承担的经济赔偿责任，由保险人负责赔偿） */
	SpecAgreement    null.String    `json:"SpecAgreement,omitempty" db:"spec_agreement,false,character varying"`        /* spec_agreement 特别约定 */
	RemindersNum     null.Int       `json:"RemindersNum,omitempty" db:"reminders_num,false,smallint"`                   /* reminders_num 催款次数 */
	IsEntryPolicy    null.Bool      `json:"IsEntryPolicy,omitempty" db:"is_entry_policy,false,boolean"`                 /* is_entry_policy 保单是否已录入承保公司系统 */
	IsAdminPay       null.Bool      `json:"IsAdminPay,omitempty" db:"is_admin_pay,false,boolean"`                       /* is_admin_pay 管理员是否支付 */
	PolicyEnrollTime null.Int       `json:"PolicyEnrollTime,omitempty" db:"policy_enroll_time,false,bigint"`            /* policy_enroll_time 录单时间 */
	ZeroPayStatus    null.String    `json:"ZeroPayStatus,omitempty" db:"zero_pay_status,false,character varying"`       /* zero_pay_status 0元实缴状态, 00: 未撤单, 02: 待实缴 04: 已实缴 06: 原保单已支付，不实缴 */
	ExternalStatus   null.String    `json:"ExternalStatus,omitempty" db:"external_status,false,character varying"`      /* external_status 保单外部状态, 00: 待撤单, 02:撤单成功, 04:撤单失败 */
	CancelDesc       null.String    `json:"CancelDesc,omitempty" db:"cancel_desc,false,character varying"`              /* cancel_desc 撤单类型,04 重新录单 08撤销 20 拒保 24 退保 */
	Favorite         null.Bool      `json:"Favorite,omitempty" db:"favorite,false,boolean"`                             /* favorite 收藏 */
	Creator          null.Int       `json:"Creator,omitempty" db:"creator,false,bigint"`                                /* creator 创建者用户ID */
	CreateTime       null.Int       `json:"CreateTime,omitempty" db:"create_time,false,bigint"`                         /* create_time 创建时间 */
	Regenerator      null.Int       `json:"Regenerator,omitempty" db:"regenerator,false,bigint"`                        /* regenerator 更新者 */
	UpdateTime       null.Int       `json:"UpdateTime,omitempty" db:"update_time,false,bigint"`                         /* update_time 更新时间 */
	DomainID         null.Int       `json:"DomainID,omitempty" db:"domain_id,false,bigint"`                             /* domain_id 数据属主 */
	Addi             types.JSONText `json:"Addi,omitempty" db:"addi,false,jsonb"`                                       /* addi 附加数据 */
	Remark           null.String    `json:"Remark,omitempty" db:"remark,false,character varying"`                       /* remark 备注 */
	Status           null.String    `json:"Status,omitempty" db:"status,false,character varying"`                       /* status 一期，0：受理中，2：在保，4：过保, 6: 作废。二期，00: 正常, 04: 重新录单, 08: 撤消, 12: 续保, 16: 已重新录单, 20: 退保, 24: 拒保 */
	Filter                          // build DML where clause
}

//TInsurancePolicyFields full field list for default query
var TInsurancePolicyFields = []string{
	"ID",
	"Sn",
	"SnCreator",
	"Name",
	"OrderID",
	"Policy",
	"Start",
	"Cease",
	"Year",
	"Duration",
	"Premium",
	"ThirdPartyPremium",
	"ThirdPartyAccount",
	"PayTime",
	"PayChannel",
	"PayType",
	"UnitPrice",
	"OrgID",
	"OrgManagerID",
	"PolicyholderType",
	"Policyholder",
	"PolicyholderID",
	"InsuranceType",
	"InsuranceTypeID",
	"PolicyScheme",
	"ActivityName",
	"ActivityCategory",
	"ActivityDesc",
	"ActivityLocation",
	"ActivityDateSet",
	"InsuredCount",
	"CompulsoryStudentNum",
	"NonCompulsoryStudentNum",
	"Contact",
	"FeeScheme",
	"CarServiceTarget",
	"Same",
	"Relation",
	"Insured",
	"InsuredID",
	"HaveInsuredList",
	"InsuredGroupByDay",
	"InsuredType",
	"InsuredList",
	"Indate",
	"Jurisdiction",
	"DisputeHandling",
	"PrevPolicyNo",
	"InsureBase",
	"BlanketInsureCode",
	"CustomType",
	"TrainProjects",
	"BusinessLocations",
	"ArbitralAgency",
	"PoolNum",
	"OpenPoolNum",
	"HeatedPoolNum",
	"TrainingPoolNum",
	"InnerArea",
	"OuterArea",
	"PoolName",
	"HaveDinnerNum",
	"DinnerNum",
	"CanteenNum",
	"ShopNum",
	"HaveRides",
	"HaveExplosive",
	"Area",
	"TrafficNum",
	"TemperatureType",
	"IsIndoor",
	"Extra",
	"BankAccount",
	"PayContact",
	"HaveSuddenDeath",
	"SuddenDeathTerms",
	"SpecAgreement",
	"RemindersNum",
	"IsEntryPolicy",
	"IsAdminPay",
	"PolicyEnrollTime",
	"ZeroPayStatus",
	"ExternalStatus",
	"CancelDesc",
	"Favorite",
	"Creator",
	"CreateTime",
	"Regenerator",
	"UpdateTime",
	"DomainID",
	"Addi",
	"Remark",
	"Status",
}

//Fields return all fields of struct.
func (r *TInsurancePolicy) Fields() []string {
	return TInsurancePolicyFields
}

//GetTableName return the associated db table name.
func (r *TInsurancePolicy) GetTableName() string {
	var viewNamePattern = regexp.MustCompile(`(?i)^t_v_[a-z0-9_]+$`)
	tableName := "t_insurance_policy"
	if viewNamePattern.MatchString(tableName) {
		return tableName[2:]
	}
	return tableName
}

// Create inserts the TInsurancePolicy to the database.
func (r *TInsurancePolicy) Create(db Queryer) error {
	err := db.QueryRow(
		`INSERT INTO t_insurance_policy (sn, sn_creator, name, order_id, policy, start, cease, year, duration, premium, third_party_premium, third_party_account, pay_time, pay_channel, pay_type, unit_price, org_id, org_manager_id, policyholder_type, policyholder, policyholder_id, insurance_type, insurance_type_id, policy_scheme, activity_name, activity_category, activity_desc, activity_location, activity_date_set, insured_count, compulsory_student_num, non_compulsory_student_num, contact, fee_scheme, car_service_target, same, relation, insured, insured_id, have_insured_list, insured_group_by_day, insured_type, insured_list, indate, jurisdiction, dispute_handling, prev_policy_no, insure_base, blanket_insure_code, custom_type, train_projects, business_locations, arbitral_agency, pool_num, open_pool_num, heated_pool_num, training_pool_num, inner_area, outer_area, pool_name, have_dinner_num, dinner_num, canteen_num, shop_num, have_rides, have_explosive, area, traffic_num, temperature_type, is_indoor, extra, bank_account, pay_contact, have_sudden_death, sudden_death_terms, spec_agreement, reminders_num, is_entry_policy, is_admin_pay, policy_enroll_time, zero_pay_status, external_status, cancel_desc, favorite, creator, create_time, regenerator, update_time, domain_id, addi, remark, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36, $37, $38, $39, $40, $41, $42, $43, $44, $45, $46, $47, $48, $49, $50, $51, $52, $53, $54, $55, $56, $57, $58, $59, $60, $61, $62, $63, $64, $65, $66, $67, $68, $69, $70, $71, $72, $73, $74, $75, $76, $77, $78, $79, $80, $81, $82, $83, $84, $85, $86, $87, $88, $89, $90, $91, $92) RETURNING id`,
		&r.Sn, &r.SnCreator, &r.Name, &r.OrderID, &r.Policy, &r.Start, &r.Cease, &r.Year, &r.Duration, &r.Premium, &r.ThirdPartyPremium, &r.ThirdPartyAccount, &r.PayTime, &r.PayChannel, &r.PayType, &r.UnitPrice, &r.OrgID, &r.OrgManagerID, &r.PolicyholderType, &r.Policyholder, &r.PolicyholderID, &r.InsuranceType, &r.InsuranceTypeID, &r.PolicyScheme, &r.ActivityName, &r.ActivityCategory, &r.ActivityDesc, &r.ActivityLocation, &r.ActivityDateSet, &r.InsuredCount, &r.CompulsoryStudentNum, &r.NonCompulsoryStudentNum, &r.Contact, &r.FeeScheme, &r.CarServiceTarget, &r.Same, &r.Relation, &r.Insured, &r.InsuredID, &r.HaveInsuredList, &r.InsuredGroupByDay, &r.InsuredType, &r.InsuredList, &r.Indate, &r.Jurisdiction, &r.DisputeHandling, &r.PrevPolicyNo, &r.InsureBase, &r.BlanketInsureCode, &r.CustomType, &r.TrainProjects, &r.BusinessLocations, &r.ArbitralAgency, &r.PoolNum, &r.OpenPoolNum, &r.HeatedPoolNum, &r.TrainingPoolNum, &r.InnerArea, &r.OuterArea, &r.PoolName, &r.HaveDinnerNum, &r.DinnerNum, &r.CanteenNum, &r.ShopNum, &r.HaveRides, &r.HaveExplosive, &r.Area, &r.TrafficNum, &r.TemperatureType, &r.IsIndoor, &r.Extra, &r.BankAccount, &r.PayContact, &r.HaveSuddenDeath, &r.SuddenDeathTerms, &r.SpecAgreement, &r.RemindersNum, &r.IsEntryPolicy, &r.IsAdminPay, &r.PolicyEnrollTime, &r.ZeroPayStatus, &r.ExternalStatus, &r.CancelDesc, &r.Favorite, &r.Creator, &r.CreateTime, &r.Regenerator, &r.UpdateTime, &r.DomainID, &r.Addi, &r.Remark, &r.Status).Scan(&r.ID)
	if err != nil {
		return errors.Wrap(err, "failed to insert t_insurance_policy")
	}
	return nil
}

// GetTInsurancePolicyByPk select the TInsurancePolicy from the database.
func GetTInsurancePolicyByPk(db Queryer, pk0 null.Int) (*TInsurancePolicy, error) {
	var r TInsurancePolicy
	err := db.QueryRow(
		`SELECT id, sn, sn_creator, name, order_id, policy, start, cease, year, duration, premium, third_party_premium, third_party_account, pay_time, pay_channel, pay_type, unit_price, org_id, org_manager_id, policyholder_type, policyholder, policyholder_id, insurance_type, insurance_type_id, policy_scheme, activity_name, activity_category, activity_desc, activity_location, activity_date_set, insured_count, compulsory_student_num, non_compulsory_student_num, contact, fee_scheme, car_service_target, same, relation, insured, insured_id, have_insured_list, insured_group_by_day, insured_type, insured_list, indate, jurisdiction, dispute_handling, prev_policy_no, insure_base, blanket_insure_code, custom_type, train_projects, business_locations, arbitral_agency, pool_num, open_pool_num, heated_pool_num, training_pool_num, inner_area, outer_area, pool_name, have_dinner_num, dinner_num, canteen_num, shop_num, have_rides, have_explosive, area, traffic_num, temperature_type, is_indoor, extra, bank_account, pay_contact, have_sudden_death, sudden_death_terms, spec_agreement, reminders_num, is_entry_policy, is_admin_pay, policy_enroll_time, zero_pay_status, external_status, cancel_desc, favorite, creator, create_time, regenerator, update_time, domain_id, addi, remark, status FROM t_insurance_policy WHERE id = $1`,
		pk0).Scan(&r.ID, &r.Sn, &r.SnCreator, &r.Name, &r.OrderID, &r.Policy, &r.Start, &r.Cease, &r.Year, &r.Duration, &r.Premium, &r.ThirdPartyPremium, &r.ThirdPartyAccount, &r.PayTime, &r.PayChannel, &r.PayType, &r.UnitPrice, &r.OrgID, &r.OrgManagerID, &r.PolicyholderType, &r.Policyholder, &r.PolicyholderID, &r.InsuranceType, &r.InsuranceTypeID, &r.PolicyScheme, &r.ActivityName, &r.ActivityCategory, &r.ActivityDesc, &r.ActivityLocation, &r.ActivityDateSet, &r.InsuredCount, &r.CompulsoryStudentNum, &r.NonCompulsoryStudentNum, &r.Contact, &r.FeeScheme, &r.CarServiceTarget, &r.Same, &r.Relation, &r.Insured, &r.InsuredID, &r.HaveInsuredList, &r.InsuredGroupByDay, &r.InsuredType, &r.InsuredList, &r.Indate, &r.Jurisdiction, &r.DisputeHandling, &r.PrevPolicyNo, &r.InsureBase, &r.BlanketInsureCode, &r.CustomType, &r.TrainProjects, &r.BusinessLocations, &r.ArbitralAgency, &r.PoolNum, &r.OpenPoolNum, &r.HeatedPoolNum, &r.TrainingPoolNum, &r.InnerArea, &r.OuterArea, &r.PoolName, &r.HaveDinnerNum, &r.DinnerNum, &r.CanteenNum, &r.ShopNum, &r.HaveRides, &r.HaveExplosive, &r.Area, &r.TrafficNum, &r.TemperatureType, &r.IsIndoor, &r.Extra, &r.BankAccount, &r.PayContact, &r.HaveSuddenDeath, &r.SuddenDeathTerms, &r.SpecAgreement, &r.RemindersNum, &r.IsEntryPolicy, &r.IsAdminPay, &r.PolicyEnrollTime, &r.ZeroPayStatus, &r.ExternalStatus, &r.CancelDesc, &r.Favorite, &r.Creator, &r.CreateTime, &r.Regenerator, &r.UpdateTime, &r.DomainID, &r.Addi, &r.Remark, &r.Status)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select t_insurance_policy")
	}
	return &r, nil
}

/*TInsuranceTypes 保险类型表 represents kuser.t_insurance_types */
type TInsuranceTypes struct {
	ID                      null.Int       `json:"ID,omitempty" db:"id,true,integer"`                                                /* id 保险产品id */
	RefID                   null.Int       `json:"RefID,omitempty" db:"ref_id,false,bigint"`                                         /* ref_id 机构引用的保险方案编码，即本表 org_id=0 and data_type="4" and parent_id=各险种ID的数据 */
	Name                    string         `json:"Name,omitempty" db:"name,false,character varying"`                                 /* name 保险产品名称 */
	Alias                   null.String    `json:"Alias,omitempty" db:"alias,false,character varying"`                               /* alias 别名 */
	DataType                null.String    `json:"DataType,omitempty" db:"data_type,false,character varying"`                        /* data_type 0: 保险产品分类, 2: 保险产品定义, 4: 投保规则, 6: 保险方案, 8: 默认投保规则/方案 */
	ParentID                null.Int       `json:"ParentID,omitempty" db:"parent_id,false,bigint"`                                   /* parent_id 隶属保险产品分类，0：表示无上级分类 */
	AgeLimit                types.JSONText `json:"AgeLimit,omitempty" db:"age_limit,false,jsonb"`                                    /* age_limit 年龄限制 */
	RuleBatch               null.String    `json:"RuleBatch,omitempty" db:"rule_batch,false,character varying"`                      /* rule_batch 规则批次 */
	OrgID                   null.Int       `json:"OrgID,omitempty" db:"org_id,false,bigint"`                                         /* org_id 投保规则、方案所属的机构编码 */
	PayType                 null.String    `json:"PayType,omitempty" db:"pay_type,false,character varying"`                          /* pay_type 支付方式: 对公转账/在线支付/线下支付 */
	PayChannel              null.String    `json:"PayChannel,omitempty" db:"pay_channel,false,character varying"`                    /* pay_channel 校快保，泰合，近邻，人保，太平洋保险 */
	PayName                 null.String    `json:"PayName,omitempty" db:"pay_name,false,character varying"`                          /* pay_name 支付项显示的名称 */
	BankAccount             null.String    `json:"BankAccount,omitempty" db:"bank_account,false,character varying"`                  /* bank_account 收款银行账号 */
	BankAccountName         null.String    `json:"BankAccountName,omitempty" db:"bank_account_name,false,character varying"`         /* bank_account_name 收款户名 */
	BankName                null.String    `json:"BankName,omitempty" db:"bank_name,false,character varying"`                        /* bank_name 开户行名称 */
	BankID                  null.String    `json:"BankID,omitempty" db:"bank_id,false,character varying"`                            /* bank_id 开户行行号 */
	FloorPrice              null.Float     `json:"FloorPrice,omitempty" db:"floor_price,false,double precision"`                     /* floor_price 首页显示的最低价 */
	UnitPrice               null.Float     `json:"UnitPrice,omitempty" db:"unit_price,false,double precision"`                       /* unit_price 单价 */
	Price                   null.Float     `json:"Price,omitempty" db:"price,false,double precision"`                                /* price 价格（分） */
	PriceConfig             types.JSONText `json:"PriceConfig,omitempty" db:"price_config,false,jsonb"`                              /* price_config 价格方案 */
	DefineLevel             null.Int       `json:"DefineLevel,omitempty" db:"define_level,false,smallint"`                           /* define_level 保险产品实际层次 */
	LayoutOrder             null.Int       `json:"LayoutOrder,omitempty" db:"layout_order,false,smallint"`                           /* layout_order 保险产品显示顺序 */
	LayoutLevel             null.Int       `json:"LayoutLevel,omitempty" db:"layout_level,false,smallint"`                           /* layout_level 保险产品显示层次 */
	ListTpl                 null.String    `json:"ListTpl,omitempty" db:"list_tpl,false,character varying"`                          /* list_tpl 清单模板 */
	Files                   types.JSONText `json:"Files,omitempty" db:"files,false,jsonb"`                                           /* files 清单模板 */
	Resource                types.JSONText `json:"Resource,omitempty" db:"resource,false,jsonb"`                                     /* resource 资源 */
	Pic                     null.String    `json:"Pic,omitempty" db:"pic,false,character varying"`                                   /* pic 关联图片 */
	SuddenDeathDescription  types.JSONText `json:"SuddenDeathDescription,omitempty" db:"sudden_death_description,false,jsonb"`       /* sudden_death_description 猝死责任险描述 */
	Description             null.String    `json:"Description,omitempty" db:"description,false,character varying"`                   /* description 首页描述 */
	AutoFill                null.String    `json:"AutoFill,omitempty" db:"auto_fill,false,character varying"`                        /* auto_fill 第三方录单(比赛险-人保录单), 0: 不自动录单，2：自动录单 */
	EnableImportList        null.Bool      `json:"EnableImportList,omitempty" db:"enable_import_list,false,boolean"`                 /* enable_import_list 允许录入清单 */
	HaveDinnerNum           null.Bool      `json:"HaveDinnerNum,omitempty" db:"have_dinner_num,false,boolean"`                       /* have_dinner_num 是否开启就餐人数 */
	InvoiceTitleUpdateTimes null.Int       `json:"InvoiceTitleUpdateTimes,omitempty" db:"invoice_title_update_times,false,smallint"` /* invoice_title_update_times 发票抬头修改次数设置 */
	ReceiptAccount          types.JSONText `json:"ReceiptAccount,omitempty" db:"receipt_account,false,jsonb"`                        /* receipt_account 对公账号设置,例:{"户名":"广州校快保科技有限公司 ",
	"开户行":"中国银行",
	"账号":"45641857894861548979"} */
	TransferAuthFiles types.JSONText `json:"TransferAuthFiles,omitempty" db:"transfer_auth_files,false,jsonb"`     /* transfer_auth_files 转账授权说明文件 */
	Contact           types.JSONText `json:"Contact,omitempty" db:"contact,false,jsonb"`                           /* contact 协议价短信联系人,例:{"联系人":"张鸣","联系电话":18311706633} */
	ContactQrCode     null.String    `json:"ContactQrCode,omitempty" db:"contact_qr_code,false,character varying"` /* contact_qr_code 缴费联系人设置,存放二维码 */
	OtherFiles        types.JSONText `json:"OtherFiles,omitempty" db:"other_files,false,jsonb"`                    /* other_files 其它相关文件 */
	Insurer           null.String    `json:"Insurer,omitempty" db:"insurer,false,character varying"`               /* insurer 承保公司,用于方案规则 */
	Underwriter       types.JSONText `json:"Underwriter,omitempty" db:"underwriter,false,jsonb"`                   /* underwriter 承保公司 */
	RemindDays        null.Int       `json:"RemindDays,omitempty" db:"remind_days,false,smallint"`                 /* remind_days 自动催款天数 */
	Mail              types.JSONText `json:"Mail,omitempty" db:"mail,false,jsonb"`                                 /* mail 邮寄地址设置,例:{"收件人":"张鸣",
	"联系电话":18311706633,
	"邮寄地址":"广东省广州市番禺区大学城外环西路303号校快保科技有限公司"} */
	OrderRepeatLimit  null.Int       `json:"OrderRepeatLimit,omitempty" db:"order_repeat_limit,false,smallint"`           /* order_repeat_limit 最大订单份数 */
	GroupByMaxDay     null.Int       `json:"GroupByMaxDay,omitempty" db:"group_by_max_day,false,smallint"`                /* group_by_max_day 允许最多按天分组数 */
	WebDescription    null.String    `json:"WebDescription,omitempty" db:"web_description,false,character varying"`       /* web_description PC页面描述 */
	MobileDescription null.String    `json:"MobileDescription,omitempty" db:"mobile_description,false,character varying"` /* mobile_description 移动端页面描述 */
	AutoFillParam     types.JSONText `json:"AutoFillParam,omitempty" db:"auto_fill_param,false,jsonb"`                    /* auto_fill_param 存放各个险的特定参数 */
	Interval          null.Int       `json:"Interval,omitempty" db:"interval,false,bigint"`                               /* interval 间隔时间 */
	MaxInsureInYear   null.Int       `json:"MaxInsureInYear,omitempty" db:"max_insure_in_year,false,smallint"`            /* max_insure_in_year 最长投保年限（年） */
	InsuredInMonth    null.Int       `json:"InsuredInMonth,omitempty" db:"insured_in_month,false,smallint"`               /* insured_in_month 保障时长（月） */
	InsuredStartTime  null.Int       `json:"InsuredStartTime,omitempty" db:"insured_start_time,false,bigint"`             /* insured_start_time 起保日期 */
	InsuredEndTime    null.Int       `json:"InsuredEndTime,omitempty" db:"insured_end_time,false,bigint"`                 /* insured_end_time 止保日期 */
	AllowStart        null.Int       `json:"AllowStart,omitempty" db:"allow_start,false,bigint"`                          /* allow_start 投保开始日期 */
	AllowEnd          null.Int       `json:"AllowEnd,omitempty" db:"allow_end,false,bigint"`                              /* allow_end 投保结束日期 */
	IndateStart       null.Int       `json:"IndateStart,omitempty" db:"indate_start,false,bigint"`                        /* indate_start 规则起效日期 */
	IndateEnd         null.Int       `json:"IndateEnd,omitempty" db:"indate_end,false,bigint"`                            /* indate_end 规则失效日期 */
	Creator           null.String    `json:"Creator,omitempty" db:"creator,false,character varying"`                      /* creator 创建者 */
	CreateTime        null.Int       `json:"CreateTime,omitempty" db:"create_time,false,bigint"`                          /* create_time 创建时间 */
	Regenerator       null.Int       `json:"Regenerator,omitempty" db:"regenerator,false,bigint"`                         /* regenerator 更新者 */
	UpdateTime        null.Int       `json:"UpdateTime,omitempty" db:"update_time,false,bigint"`                          /* update_time 更新时间 */
	Updator           null.String    `json:"Updator,omitempty" db:"updator,false,character varying"`                      /* updator 更新者 */
	Addi              types.JSONText `json:"Addi,omitempty" db:"addi,false,jsonb"`                                        /* addi 备用字段 */
	DomainID          null.Int       `json:"DomainID,omitempty" db:"domain_id,false,bigint"`                              /* domain_id 数据属主 */
	Remark            null.String    `json:"Remark,omitempty" db:"remark,false,character varying"`                        /* remark 备注 */
	Status            null.String    `json:"Status,omitempty" db:"status,false,character varying"`                        /* status 状态, 0: 正常，2:等待推出, 4：禁用，6：作废 */
	Filter                           // build DML where clause
}

//TInsuranceTypesFields full field list for default query
var TInsuranceTypesFields = []string{
	"ID",
	"RefID",
	"Name",
	"Alias",
	"DataType",
	"ParentID",
	"AgeLimit",
	"RuleBatch",
	"OrgID",
	"PayType",
	"PayChannel",
	"PayName",
	"BankAccount",
	"BankAccountName",
	"BankName",
	"BankID",
	"FloorPrice",
	"UnitPrice",
	"Price",
	"PriceConfig",
	"DefineLevel",
	"LayoutOrder",
	"LayoutLevel",
	"ListTpl",
	"Files",
	"Resource",
	"Pic",
	"SuddenDeathDescription",
	"Description",
	"AutoFill",
	"EnableImportList",
	"HaveDinnerNum",
	"InvoiceTitleUpdateTimes",
	"ReceiptAccount",
	"TransferAuthFiles",
	"Contact",
	"ContactQrCode",
	"OtherFiles",
	"Insurer",
	"Underwriter",
	"RemindDays",
	"Mail",
	"OrderRepeatLimit",
	"GroupByMaxDay",
	"WebDescription",
	"MobileDescription",
	"AutoFillParam",
	"Interval",
	"MaxInsureInYear",
	"InsuredInMonth",
	"InsuredStartTime",
	"InsuredEndTime",
	"AllowStart",
	"AllowEnd",
	"IndateStart",
	"IndateEnd",
	"Creator",
	"CreateTime",
	"Regenerator",
	"UpdateTime",
	"Updator",
	"Addi",
	"DomainID",
	"Remark",
	"Status",
}

//Fields return all fields of struct.
func (r *TInsuranceTypes) Fields() []string {
	return TInsuranceTypesFields
}

//GetTableName return the associated db table name.
func (r *TInsuranceTypes) GetTableName() string {
	var viewNamePattern = regexp.MustCompile(`(?i)^t_v_[a-z0-9_]+$`)
	tableName := "t_insurance_types"
	if viewNamePattern.MatchString(tableName) {
		return tableName[2:]
	}
	return tableName
}

// Create inserts the TInsuranceTypes to the database.
func (r *TInsuranceTypes) Create(db Queryer) error {
	err := db.QueryRow(
		`INSERT INTO t_insurance_types (ref_id, name, alias, data_type, parent_id, age_limit, rule_batch, org_id, pay_type, pay_channel, pay_name, bank_account, bank_account_name, bank_name, bank_id, floor_price, unit_price, price, price_config, define_level, layout_order, layout_level, list_tpl, files, resource, pic, sudden_death_description, description, auto_fill, enable_import_list, have_dinner_num, invoice_title_update_times, receipt_account, transfer_auth_files, contact, contact_qr_code, other_files, insurer, underwriter, remind_days, mail, order_repeat_limit, group_by_max_day, web_description, mobile_description, auto_fill_param, interval, max_insure_in_year, insured_in_month, insured_start_time, insured_end_time, allow_start, allow_end, indate_start, indate_end, creator, create_time, regenerator, update_time, updator, addi, domain_id, remark, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36, $37, $38, $39, $40, $41, $42, $43, $44, $45, $46, $47, $48, $49, $50, $51, $52, $53, $54, $55, $56, $57, $58, $59, $60, $61, $62, $63, $64) RETURNING id`,
		&r.RefID, &r.Name, &r.Alias, &r.DataType, &r.ParentID, &r.AgeLimit, &r.RuleBatch, &r.OrgID, &r.PayType, &r.PayChannel, &r.PayName, &r.BankAccount, &r.BankAccountName, &r.BankName, &r.BankID, &r.FloorPrice, &r.UnitPrice, &r.Price, &r.PriceConfig, &r.DefineLevel, &r.LayoutOrder, &r.LayoutLevel, &r.ListTpl, &r.Files, &r.Resource, &r.Pic, &r.SuddenDeathDescription, &r.Description, &r.AutoFill, &r.EnableImportList, &r.HaveDinnerNum, &r.InvoiceTitleUpdateTimes, &r.ReceiptAccount, &r.TransferAuthFiles, &r.Contact, &r.ContactQrCode, &r.OtherFiles, &r.Insurer, &r.Underwriter, &r.RemindDays, &r.Mail, &r.OrderRepeatLimit, &r.GroupByMaxDay, &r.WebDescription, &r.MobileDescription, &r.AutoFillParam, &r.Interval, &r.MaxInsureInYear, &r.InsuredInMonth, &r.InsuredStartTime, &r.InsuredEndTime, &r.AllowStart, &r.AllowEnd, &r.IndateStart, &r.IndateEnd, &r.Creator, &r.CreateTime, &r.Regenerator, &r.UpdateTime, &r.Updator, &r.Addi, &r.DomainID, &r.Remark, &r.Status).Scan(&r.ID)
	if err != nil {
		return errors.Wrap(err, "failed to insert t_insurance_types")
	}
	return nil
}

// GetTInsuranceTypesByPk select the TInsuranceTypes from the database.
func GetTInsuranceTypesByPk(db Queryer, pk0 null.Int) (*TInsuranceTypes, error) {
	var r TInsuranceTypes
	err := db.QueryRow(
		`SELECT id, ref_id, name, alias, data_type, parent_id, age_limit, rule_batch, org_id, pay_type, pay_channel, pay_name, bank_account, bank_account_name, bank_name, bank_id, floor_price, unit_price, price, price_config, define_level, layout_order, layout_level, list_tpl, files, resource, pic, sudden_death_description, description, auto_fill, enable_import_list, have_dinner_num, invoice_title_update_times, receipt_account, transfer_auth_files, contact, contact_qr_code, other_files, insurer, underwriter, remind_days, mail, order_repeat_limit, group_by_max_day, web_description, mobile_description, auto_fill_param, interval, max_insure_in_year, insured_in_month, insured_start_time, insured_end_time, allow_start, allow_end, indate_start, indate_end, creator, create_time, regenerator, update_time, updator, addi, domain_id, remark, status FROM t_insurance_types WHERE id = $1`,
		pk0).Scan(&r.ID, &r.RefID, &r.Name, &r.Alias, &r.DataType, &r.ParentID, &r.AgeLimit, &r.RuleBatch, &r.OrgID, &r.PayType, &r.PayChannel, &r.PayName, &r.BankAccount, &r.BankAccountName, &r.BankName, &r.BankID, &r.FloorPrice, &r.UnitPrice, &r.Price, &r.PriceConfig, &r.DefineLevel, &r.LayoutOrder, &r.LayoutLevel, &r.ListTpl, &r.Files, &r.Resource, &r.Pic, &r.SuddenDeathDescription, &r.Description, &r.AutoFill, &r.EnableImportList, &r.HaveDinnerNum, &r.InvoiceTitleUpdateTimes, &r.ReceiptAccount, &r.TransferAuthFiles, &r.Contact, &r.ContactQrCode, &r.OtherFiles, &r.Insurer, &r.Underwriter, &r.RemindDays, &r.Mail, &r.OrderRepeatLimit, &r.GroupByMaxDay, &r.WebDescription, &r.MobileDescription, &r.AutoFillParam, &r.Interval, &r.MaxInsureInYear, &r.InsuredInMonth, &r.InsuredStartTime, &r.InsuredEndTime, &r.AllowStart, &r.AllowEnd, &r.IndateStart, &r.IndateEnd, &r.Creator, &r.CreateTime, &r.Regenerator, &r.UpdateTime, &r.Updator, &r.Addi, &r.DomainID, &r.Remark, &r.Status)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select t_insurance_types")
	}
	return &r, nil
}

/*TInsureAttach 保单附件 represents kuser.t_insure_attach */
type TInsureAttach struct {
	ID             null.Int       `json:"ID,omitempty" db:"id,true,integer"`                           /* id 编号 */
	TUID           null.Int       `json:"TUID,omitempty" db:"t_u_id,false,bigint"`                     /* t_u_id 用户内部编号 */
	SchoolID       null.Int       `json:"SchoolID,omitempty" db:"school_id,false,bigint"`              /* school_id 学校编号 */
	Grade          null.String    `json:"Grade,omitempty" db:"grade,false,character varying"`          /* grade 年级 */
	Year           null.Int       `json:"Year,omitempty" db:"year,false,smallint"`                     /* year 保单年份 */
	Batch          null.String    `json:"Batch,omitempty" db:"batch,false,character varying"`          /* batch 批次 */
	PolicyNo       null.String    `json:"PolicyNo,omitempty" db:"policy_no,false,character varying"`   /* policy_no 保单号 */
	InsurePolicyID null.Int       `json:"InsurePolicyID,omitempty" db:"insure_policy_id,false,bigint"` /* insure_policy_id 系统保单编号 */
	Others         types.JSONText `json:"Others,omitempty" db:"others,false,jsonb"`                    /* others 其它 */
	Files          types.JSONText `json:"Files,omitempty" db:"files,false,jsonb"`                      /* files 保单附件 */
	Creator        null.Int       `json:"Creator,omitempty" db:"creator,false,bigint"`                 /* creator 创建者用户ID */
	CreateTime     null.Int       `json:"CreateTime,omitempty" db:"create_time,false,bigint"`          /* create_time 创建时间 */
	Regenerator    null.Int       `json:"Regenerator,omitempty" db:"regenerator,false,bigint"`         /* regenerator 更新者 */
	UpdateTime     null.Int       `json:"UpdateTime,omitempty" db:"update_time,false,bigint"`          /* update_time 修改时间 */
	Addi           types.JSONText `json:"Addi,omitempty" db:"addi,false,jsonb"`                        /* addi 附加数据 */
	DomainID       null.Int       `json:"DomainID,omitempty" db:"domain_id,false,bigint"`              /* domain_id 数据属主 */
	Remark         null.String    `json:"Remark,omitempty" db:"remark,false,character varying"`        /* remark 备注 */
	Status         null.String    `json:"Status,omitempty" db:"status,false,character varying"`        /* status 状态 */
	Filter                        // build DML where clause
}

//TInsureAttachFields full field list for default query
var TInsureAttachFields = []string{
	"ID",
	"TUID",
	"SchoolID",
	"Grade",
	"Year",
	"Batch",
	"PolicyNo",
	"InsurePolicyID",
	"Others",
	"Files",
	"Creator",
	"CreateTime",
	"Regenerator",
	"UpdateTime",
	"Addi",
	"DomainID",
	"Remark",
	"Status",
}

//Fields return all fields of struct.
func (r *TInsureAttach) Fields() []string {
	return TInsureAttachFields
}

//GetTableName return the associated db table name.
func (r *TInsureAttach) GetTableName() string {
	var viewNamePattern = regexp.MustCompile(`(?i)^t_v_[a-z0-9_]+$`)
	tableName := "t_insure_attach"
	if viewNamePattern.MatchString(tableName) {
		return tableName[2:]
	}
	return tableName
}

// Create inserts the TInsureAttach to the database.
func (r *TInsureAttach) Create(db Queryer) error {
	err := db.QueryRow(
		`INSERT INTO t_insure_attach (t_u_id, school_id, grade, year, batch, policy_no, insure_policy_id, others, files, creator, create_time, regenerator, update_time, addi, domain_id, remark, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17) RETURNING id`,
		&r.TUID, &r.SchoolID, &r.Grade, &r.Year, &r.Batch, &r.PolicyNo, &r.InsurePolicyID, &r.Others, &r.Files, &r.Creator, &r.CreateTime, &r.Regenerator, &r.UpdateTime, &r.Addi, &r.DomainID, &r.Remark, &r.Status).Scan(&r.ID)
	if err != nil {
		return errors.Wrap(err, "failed to insert t_insure_attach")
	}
	return nil
}

// GetTInsureAttachByPk select the TInsureAttach from the database.
func GetTInsureAttachByPk(db Queryer, pk0 null.Int) (*TInsureAttach, error) {
	var r TInsureAttach
	err := db.QueryRow(
		`SELECT id, t_u_id, school_id, grade, year, batch, policy_no, insure_policy_id, others, files, creator, create_time, regenerator, update_time, addi, domain_id, remark, status FROM t_insure_attach WHERE id = $1`,
		pk0).Scan(&r.ID, &r.TUID, &r.SchoolID, &r.Grade, &r.Year, &r.Batch, &r.PolicyNo, &r.InsurePolicyID, &r.Others, &r.Files, &r.Creator, &r.CreateTime, &r.Regenerator, &r.UpdateTime, &r.Addi, &r.DomainID, &r.Remark, &r.Status)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select t_insure_attach")
	}
	return &r, nil
}

/*TInsuredDetail 清单表，校车的校车信息 校车承运人存在同一行，但前端分开显示 represents kuser.t_insured_detail */
type TInsuredDetail struct {
	ID                    null.Int       `json:"ID,omitempty" db:"id,true,integer"`                                           /* id 主键 */
	Type                  null.String    `json:"Type,omitempty" db:"type,false,character varying"`                            /* type 清单类别 */
	SubType               null.String    `json:"SubType,omitempty" db:"sub_type,false,character varying"`                     /* sub_type 清单子类型 */
	OrderID               null.Int       `json:"OrderID,omitempty" db:"order_id,false,bigint"`                                /* order_id 订单号（内部订单号) */
	PolicyID              null.String    `json:"PolicyID,omitempty" db:"policy_id,false,character varying"`                   /* policy_id 保单号（在生成保单的时候才填写） */
	Name                  null.String    `json:"Name,omitempty" db:"name,false,character varying"`                            /* name 姓名（比赛、教工、实习生、校车承运人） */
	IDCardNo              null.String    `json:"IDCardNo,omitempty" db:"id_card_no,false,character varying"`                  /* id_card_no 证件号码（比赛、教工、实习生) */
	Gender                null.String    `json:"Gender,omitempty" db:"gender,false,character varying"`                        /* gender 性别 （比赛、教工、实习生 */
	Birthday              null.Int       `json:"Birthday,omitempty" db:"birthday,false,bigint"`                               /* birthday 出生日期（比赛、教工、实习生) */
	Role                  null.String    `json:"Role,omitempty" db:"role,false,character varying"`                            /* role 职位(工作类型)（教工) */
	Org                   null.String    `json:"Org,omitempty" db:"org,false,character varying"`                              /* org 所属机构（教工) */
	Class                 null.String    `json:"Class,omitempty" db:"class,false,character varying"`                          /* class 班别（实习生) */
	GroupDay              null.Int       `json:"GroupDay,omitempty" db:"group_day,false,bigint"`                              /* group_day 所在比赛/活动日期 */
	LicensePlateNo        null.String    `json:"LicensePlateNo,omitempty" db:"license_plate_no,false,character varying"`      /* license_plate_no 车牌号码（校车信息 校车承运人) */
	Brand                 null.String    `json:"Brand,omitempty" db:"brand,false,character varying"`                          /* brand 厂牌类型（校车信息) */
	DriverSeatNumber      null.Int       `json:"DriverSeatNumber,omitempty" db:"driver_seat_number,false,smallint"`           /* driver_seat_number 司机座位（校车承运人） */
	ApprovedPassengersNum null.Int       `json:"ApprovedPassengersNum,omitempty" db:"approved_passengers_num,false,smallint"` /* approved_passengers_num 核定客载人数 */
	SeatNum               null.Int       `json:"SeatNum,omitempty" db:"seat_num,false,smallint"`                              /* seat_num 座位数（校车信息 校车承运人) */
	RoadGrade             null.String    `json:"RoadGrade,omitempty" db:"road_grade,false,character varying"`                 /* road_grade 运营公路等级（校车信息) */
	DriverLicense         null.String    `json:"DriverLicense,omitempty" db:"driver_license,false,character varying"`         /* driver_license 驾驶证-图片 */
	DrivingLicense        null.String    `json:"DrivingLicense,omitempty" db:"driving_license,false,character varying"`       /* driving_license 行驶证-图片（校车信息) */
	Action                null.String    `json:"Action,omitempty" db:"action,false,character varying"`                        /* action 修改类型：2:新增 4:删除 6:修改（此处需要对清单ID进行比对，才可以得出） */
	ErrMsg                null.String    `json:"ErrMsg,omitempty" db:"err_msg,false,character varying"`                       /* err_msg 错误原因 */
	Province              null.String    `json:"Province,omitempty" db:"province,false,character varying"`                    /* province 省 */
	City                  null.String    `json:"City,omitempty" db:"city,false,character varying"`                            /* city 市 */
	District              null.String    `json:"District,omitempty" db:"district,false,character varying"`                    /* district 区 */
	Addr                  null.String    `json:"Addr,omitempty" db:"addr,false,character varying"`                            /* addr 地址 */
	TrainItem             null.String    `json:"TrainItem,omitempty" db:"train_item,false,character varying"`                 /* train_item 训练项目（英文逗号分隔） */
	OtherItem             null.String    `json:"OtherItem,omitempty" db:"other_item,false,character varying"`                 /* other_item 其它项目（英文逗号分隔） */
	FieldType             null.String    `json:"FieldType,omitempty" db:"field_type,false,character varying"`                 /* field_type 场地类型 */
	Area                  null.Float     `json:"Area,omitempty" db:"area,false,double precision"`                             /* area 场地面积 */
	Creator               null.Int       `json:"Creator,omitempty" db:"creator,false,bigint"`                                 /* creator 创建者 */
	CreateTime            null.Int       `json:"CreateTime,omitempty" db:"create_time,false,bigint"`                          /* create_time 创建时间 */
	Regenerator           null.Int       `json:"Regenerator,omitempty" db:"regenerator,false,bigint"`                         /* regenerator 更新者 */
	UpdateTime            null.Int       `json:"UpdateTime,omitempty" db:"update_time,false,bigint"`                          /* update_time 更新时间 */
	Updator               null.Int       `json:"Updator,omitempty" db:"updator,false,bigint"`                                 /* updator 更新人 */
	DomainID              null.Int       `json:"DomainID,omitempty" db:"domain_id,false,bigint"`                              /* domain_id 数据隶属 */
	Addi                  types.JSONText `json:"Addi,omitempty" db:"addi,false,jsonb"`                                        /* addi 附加信息 */
	Remark                null.String    `json:"Remark,omitempty" db:"remark,false,character varying"`                        /* remark 备注 （实习生) */
	Status                null.String    `json:"Status,omitempty" db:"status,false,character varying"`                        /* status 状态：0:有效 2:错误  4.拒保 */
	Filter                               // build DML where clause
}

//TInsuredDetailFields full field list for default query
var TInsuredDetailFields = []string{
	"ID",
	"Type",
	"SubType",
	"OrderID",
	"PolicyID",
	"Name",
	"IDCardNo",
	"Gender",
	"Birthday",
	"Role",
	"Org",
	"Class",
	"GroupDay",
	"LicensePlateNo",
	"Brand",
	"DriverSeatNumber",
	"ApprovedPassengersNum",
	"SeatNum",
	"RoadGrade",
	"DriverLicense",
	"DrivingLicense",
	"Action",
	"ErrMsg",
	"Province",
	"City",
	"District",
	"Addr",
	"TrainItem",
	"OtherItem",
	"FieldType",
	"Area",
	"Creator",
	"CreateTime",
	"Regenerator",
	"UpdateTime",
	"Updator",
	"DomainID",
	"Addi",
	"Remark",
	"Status",
}

//Fields return all fields of struct.
func (r *TInsuredDetail) Fields() []string {
	return TInsuredDetailFields
}

//GetTableName return the associated db table name.
func (r *TInsuredDetail) GetTableName() string {
	var viewNamePattern = regexp.MustCompile(`(?i)^t_v_[a-z0-9_]+$`)
	tableName := "t_insured_detail"
	if viewNamePattern.MatchString(tableName) {
		return tableName[2:]
	}
	return tableName
}

// Create inserts the TInsuredDetail to the database.
func (r *TInsuredDetail) Create(db Queryer) error {
	err := db.QueryRow(
		`INSERT INTO t_insured_detail (type, sub_type, order_id, policy_id, name, id_card_no, gender, birthday, role, org, class, group_day, license_plate_no, brand, driver_seat_number, approved_passengers_num, seat_num, road_grade, driver_license, driving_license, action, err_msg, province, city, district, addr, train_item, other_item, field_type, area, creator, create_time, regenerator, update_time, updator, domain_id, addi, remark, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36, $37, $38, $39) RETURNING id`,
		&r.Type, &r.SubType, &r.OrderID, &r.PolicyID, &r.Name, &r.IDCardNo, &r.Gender, &r.Birthday, &r.Role, &r.Org, &r.Class, &r.GroupDay, &r.LicensePlateNo, &r.Brand, &r.DriverSeatNumber, &r.ApprovedPassengersNum, &r.SeatNum, &r.RoadGrade, &r.DriverLicense, &r.DrivingLicense, &r.Action, &r.ErrMsg, &r.Province, &r.City, &r.District, &r.Addr, &r.TrainItem, &r.OtherItem, &r.FieldType, &r.Area, &r.Creator, &r.CreateTime, &r.Regenerator, &r.UpdateTime, &r.Updator, &r.DomainID, &r.Addi, &r.Remark, &r.Status).Scan(&r.ID)
	if err != nil {
		return errors.Wrap(err, "failed to insert t_insured_detail")
	}
	return nil
}

// GetTInsuredDetailByPk select the TInsuredDetail from the database.
func GetTInsuredDetailByPk(db Queryer, pk0 null.Int) (*TInsuredDetail, error) {
	var r TInsuredDetail
	err := db.QueryRow(
		`SELECT id, type, sub_type, order_id, policy_id, name, id_card_no, gender, birthday, role, org, class, group_day, license_plate_no, brand, driver_seat_number, approved_passengers_num, seat_num, road_grade, driver_license, driving_license, action, err_msg, province, city, district, addr, train_item, other_item, field_type, area, creator, create_time, regenerator, update_time, updator, domain_id, addi, remark, status FROM t_insured_detail WHERE id = $1`,
		pk0).Scan(&r.ID, &r.Type, &r.SubType, &r.OrderID, &r.PolicyID, &r.Name, &r.IDCardNo, &r.Gender, &r.Birthday, &r.Role, &r.Org, &r.Class, &r.GroupDay, &r.LicensePlateNo, &r.Brand, &r.DriverSeatNumber, &r.ApprovedPassengersNum, &r.SeatNum, &r.RoadGrade, &r.DriverLicense, &r.DrivingLicense, &r.Action, &r.ErrMsg, &r.Province, &r.City, &r.District, &r.Addr, &r.TrainItem, &r.OtherItem, &r.FieldType, &r.Area, &r.Creator, &r.CreateTime, &r.Regenerator, &r.UpdateTime, &r.Updator, &r.DomainID, &r.Addi, &r.Remark, &r.Status)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select t_insured_detail")
	}
	return &r, nil
}

/*TInsuredTerms 保险条款 represents kuser.t_insured_terms */
type TInsuredTerms struct {
	ID              null.Int       `json:"ID,omitempty" db:"id,true,integer"`                             /* id 编号 */
	InsuranceTypeID null.Int       `json:"InsuranceTypeID,omitempty" db:"insurance_type_id,false,bigint"` /* insurance_type_id 保险类型 */
	Topic           null.String    `json:"Topic,omitempty" db:"topic,false,character varying"`            /* topic 标题 */
	ParentID        null.Int       `json:"ParentID,omitempty" db:"parent_id,false,bigint"`                /* parent_id 父级ID */
	Level           null.Int       `json:"Level,omitempty" db:"level,false,smallint"`                     /* level 级别 */
	Content         null.String    `json:"Content,omitempty" db:"content,false,character varying"`        /* content 内容 */
	Updator         null.String    `json:"Updator,omitempty" db:"updator,false,character varying"`        /* updator 更新者 */
	Regenerator     null.Int       `json:"Regenerator,omitempty" db:"regenerator,false,bigint"`           /* regenerator 更新者 */
	UpdateTime      null.Int       `json:"UpdateTime,omitempty" db:"update_time,false,bigint"`            /* update_time 更新时间 */
	CreateTime      null.Int       `json:"CreateTime,omitempty" db:"create_time,false,bigint"`            /* create_time 创建时间 */
	Creator         null.String    `json:"Creator,omitempty" db:"creator,false,character varying"`        /* creator 创建者账号 */
	DomainID        null.Int       `json:"DomainID,omitempty" db:"domain_id,false,bigint"`                /* domain_id 数据属主 */
	Addi            types.JSONText `json:"Addi,omitempty" db:"addi,false,jsonb"`                          /* addi 附加数据 */
	Remark          null.String    `json:"Remark,omitempty" db:"remark,false,character varying"`          /* remark 备注 */
	Status          null.String    `json:"Status,omitempty" db:"status,false,character varying"`          /* status 状态0:有效, 2:修改，4删除 */
	Filter                         // build DML where clause
}

//TInsuredTermsFields full field list for default query
var TInsuredTermsFields = []string{
	"ID",
	"InsuranceTypeID",
	"Topic",
	"ParentID",
	"Level",
	"Content",
	"Updator",
	"Regenerator",
	"UpdateTime",
	"CreateTime",
	"Creator",
	"DomainID",
	"Addi",
	"Remark",
	"Status",
}

//Fields return all fields of struct.
func (r *TInsuredTerms) Fields() []string {
	return TInsuredTermsFields
}

//GetTableName return the associated db table name.
func (r *TInsuredTerms) GetTableName() string {
	var viewNamePattern = regexp.MustCompile(`(?i)^t_v_[a-z0-9_]+$`)
	tableName := "t_insured_terms"
	if viewNamePattern.MatchString(tableName) {
		return tableName[2:]
	}
	return tableName
}

// Create inserts the TInsuredTerms to the database.
func (r *TInsuredTerms) Create(db Queryer) error {
	err := db.QueryRow(
		`INSERT INTO t_insured_terms (insurance_type_id, topic, parent_id, level, content, updator, regenerator, update_time, create_time, creator, domain_id, addi, remark, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14) RETURNING id`,
		&r.InsuranceTypeID, &r.Topic, &r.ParentID, &r.Level, &r.Content, &r.Updator, &r.Regenerator, &r.UpdateTime, &r.CreateTime, &r.Creator, &r.DomainID, &r.Addi, &r.Remark, &r.Status).Scan(&r.ID)
	if err != nil {
		return errors.Wrap(err, "failed to insert t_insured_terms")
	}
	return nil
}

// GetTInsuredTermsByPk select the TInsuredTerms from the database.
func GetTInsuredTermsByPk(db Queryer, pk0 null.Int) (*TInsuredTerms, error) {
	var r TInsuredTerms
	err := db.QueryRow(
		`SELECT id, insurance_type_id, topic, parent_id, level, content, updator, regenerator, update_time, create_time, creator, domain_id, addi, remark, status FROM t_insured_terms WHERE id = $1`,
		pk0).Scan(&r.ID, &r.InsuranceTypeID, &r.Topic, &r.ParentID, &r.Level, &r.Content, &r.Updator, &r.Regenerator, &r.UpdateTime, &r.CreateTime, &r.Creator, &r.DomainID, &r.Addi, &r.Remark, &r.Status)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select t_insured_terms")
	}
	return &r, nil
}

/*TJudge 鉴定邀请表 represents kuser.t_judge */
type TJudge struct {
	ID          null.Int    `json:"ID,omitempty" db:"id,true,integer"`                    /* id 评价编号 */
	DeveloperID null.Int    `json:"DeveloperID,omitempty" db:"developer_id,false,bigint"` /* developer_id 被鉴定人 */
	ProofID     null.Int    `json:"ProofID,omitempty" db:"proof_id,false,bigint"`         /* proof_id 鉴定项 */
	WitnessID   null.Int    `json:"WitnessID,omitempty" db:"witness_id,false,bigint"`     /* witness_id 鉴定人 */
	ApplyTime   null.Int    `json:"ApplyTime,omitempty" db:"apply_time,false,bigint"`     /* apply_time 申请鉴定时间 */
	JudgeTime   null.Int    `json:"JudgeTime,omitempty" db:"judge_time,false,bigint"`     /* judge_time 鉴定时间 */
	Status      null.String `json:"Status,omitempty" db:"status,false,character varying"` /* status invited,已发出邀请
	judged,已评价、签定
	rejected,拒绝评价
	expired,邀请已过期/拒绝评价 */
	Filter // build DML where clause
}

//TJudgeFields full field list for default query
var TJudgeFields = []string{
	"ID",
	"DeveloperID",
	"ProofID",
	"WitnessID",
	"ApplyTime",
	"JudgeTime",
	"Status",
}

//Fields return all fields of struct.
func (r *TJudge) Fields() []string {
	return TJudgeFields
}

//GetTableName return the associated db table name.
func (r *TJudge) GetTableName() string {
	var viewNamePattern = regexp.MustCompile(`(?i)^t_v_[a-z0-9_]+$`)
	tableName := "t_judge"
	if viewNamePattern.MatchString(tableName) {
		return tableName[2:]
	}
	return tableName
}

// Create inserts the TJudge to the database.
func (r *TJudge) Create(db Queryer) error {
	err := db.QueryRow(
		`INSERT INTO t_judge (developer_id, proof_id, witness_id, apply_time, judge_time, status) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		&r.DeveloperID, &r.ProofID, &r.WitnessID, &r.ApplyTime, &r.JudgeTime, &r.Status).Scan(&r.ID)
	if err != nil {
		return errors.Wrap(err, "failed to insert t_judge")
	}
	return nil
}

// GetTJudgeByPk select the TJudge from the database.
func GetTJudgeByPk(db Queryer, pk0 null.Int) (*TJudge, error) {
	var r TJudge
	err := db.QueryRow(
		`SELECT id, developer_id, proof_id, witness_id, apply_time, judge_time, status FROM t_judge WHERE id = $1`,
		pk0).Scan(&r.ID, &r.DeveloperID, &r.ProofID, &r.WitnessID, &r.ApplyTime, &r.JudgeTime, &r.Status)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select t_judge")
	}
	return &r, nil
}

/*TLog t_log represents kuser.t_log */
type TLog struct {
	ID            null.Int    `json:"ID,omitempty" db:"id,true,integer"`                                    /* id 编号 */
	Grade         null.String `json:"Grade,omitempty" db:"grade,false,character varying"`                   /* grade 等级 */
	Msg           null.String `json:"Msg,omitempty" db:"msg,false,character varying"`                       /* msg 消息 */
	Caller        null.String `json:"Caller,omitempty" db:"caller,false,character varying"`                 /* caller 位置 */
	Stacktrace    null.String `json:"Stacktrace,omitempty" db:"stacktrace,false,character varying"`         /* stacktrace 栈 */
	Namespace     null.String `json:"Namespace,omitempty" db:"namespace,false,character varying"`           /* namespace 模块 */
	LoginUserName null.String `json:"LoginUserName,omitempty" db:"login_user_name,false,character varying"` /* login_user_name 用户名 */
	LoginUserID   null.Int    `json:"LoginUserID,omitempty" db:"login_user_id,false,bigint"`                /* login_user_id 用户编码 */
	DomainID      null.Int    `json:"DomainID,omitempty" db:"domain_id,false,bigint"`                       /* domain_id 数据隶属 */
	Creator       null.Int    `json:"Creator,omitempty" db:"creator,false,bigint"`                          /* creator 本数据创建者 */
	CreateTime    null.Int    `json:"CreateTime,omitempty" db:"create_time,false,bigint"`                   /* create_time 生成时间 */
	Filter                    // build DML where clause
}

//TLogFields full field list for default query
var TLogFields = []string{
	"ID",
	"Grade",
	"Msg",
	"Caller",
	"Stacktrace",
	"Namespace",
	"LoginUserName",
	"LoginUserID",
	"DomainID",
	"Creator",
	"CreateTime",
}

//Fields return all fields of struct.
func (r *TLog) Fields() []string {
	return TLogFields
}

//GetTableName return the associated db table name.
func (r *TLog) GetTableName() string {
	var viewNamePattern = regexp.MustCompile(`(?i)^t_v_[a-z0-9_]+$`)
	tableName := "t_log"
	if viewNamePattern.MatchString(tableName) {
		return tableName[2:]
	}
	return tableName
}

// Create inserts the TLog to the database.
func (r *TLog) Create(db Queryer) error {
	err := db.QueryRow(
		`INSERT INTO t_log (grade, msg, caller, stacktrace, namespace, login_user_name, login_user_id, domain_id, creator, create_time) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`,
		&r.Grade, &r.Msg, &r.Caller, &r.Stacktrace, &r.Namespace, &r.LoginUserName, &r.LoginUserID, &r.DomainID, &r.Creator, &r.CreateTime).Scan(&r.ID)
	if err != nil {
		return errors.Wrap(err, "failed to insert t_log")
	}
	return nil
}

// GetTLogByPk select the TLog from the database.
func GetTLogByPk(db Queryer, pk0 null.Int) (*TLog, error) {
	var r TLog
	err := db.QueryRow(
		`SELECT id, grade, msg, caller, stacktrace, namespace, login_user_name, login_user_id, domain_id, creator, create_time FROM t_log WHERE id = $1`,
		pk0).Scan(&r.ID, &r.Grade, &r.Msg, &r.Caller, &r.Stacktrace, &r.Namespace, &r.LoginUserName, &r.LoginUserID, &r.DomainID, &r.Creator, &r.CreateTime)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select t_log")
	}
	return &r, nil
}

/*TMistakeCorrect 报错 represents kuser.t_mistake_correct */
type TMistakeCorrect struct {
	ID                       null.Int       `json:"ID,omitempty" db:"id,true,integer"`                                                /* id id */
	OrderID                  null.Int       `json:"OrderID,omitempty" db:"order_id,false,bigint"`                                     /* order_id order_id */
	Policyholder             types.JSONText `json:"Policyholder,omitempty" db:"policyholder,false,jsonb"`                             /* policyholder 投保人 */
	Contact                  types.JSONText `json:"Contact,omitempty" db:"contact,false,jsonb"`                                       /* contact 投保联系人 */
	PolicyholderID           null.Int       `json:"PolicyholderID,omitempty" db:"policyholder_id,false,bigint"`                       /* policyholder_id 投保人编号 */
	OfficialNameP            null.String    `json:"OfficialNameP,omitempty" db:"official_name_p,false,character varying"`             /* official_name_p 投保人姓名 */
	IDCardTypeP              null.String    `json:"IDCardTypeP,omitempty" db:"id_card_type_p,false,character varying"`                /* id_card_type_p 投保人证件类型 */
	IDCardNoP                null.String    `json:"IDCardNoP,omitempty" db:"id_card_no_p,false,character varying"`                    /* id_card_no_p 投保人身份证号码 */
	GenderP                  null.String    `json:"GenderP,omitempty" db:"gender_p,false,character varying"`                          /* gender_p 投保人性别 */
	BirthdayP                null.Int       `json:"BirthdayP,omitempty" db:"birthday_p,false,bigint"`                                 /* birthday_p 投保人出生日期 */
	Insured                  types.JSONText `json:"Insured,omitempty" db:"insured,false,jsonb"`                                       /* insured 被保险人 */
	InsuredID                null.Int       `json:"InsuredID,omitempty" db:"insured_id,false,bigint"`                                 /* insured_id 被保险人编号 */
	OfficialName             null.String    `json:"OfficialName,omitempty" db:"official_name,false,character varying"`                /* official_name 姓名 */
	IDCardType               null.String    `json:"IDCardType,omitempty" db:"id_card_type,false,character varying"`                   /* id_card_type 证件类型 */
	IDCardNo                 null.String    `json:"IDCardNo,omitempty" db:"id_card_no,false,character varying"`                       /* id_card_no 身份证号码 */
	Gender                   null.String    `json:"Gender,omitempty" db:"gender,false,character varying"`                             /* gender 性别 */
	Birthday                 null.Int       `json:"Birthday,omitempty" db:"birthday,false,bigint"`                                    /* birthday 出生日期 */
	ClearList                null.Bool      `json:"ClearList,omitempty" db:"clear_list,false,boolean"`                                /* clear_list 清除清单 */
	InsuredList              types.JSONText `json:"InsuredList,omitempty" db:"insured_list,false,jsonb"`                              /* insured_list 被保险人清单 */
	InsuredCount             null.Int       `json:"InsuredCount,omitempty" db:"insured_count,false,smallint"`                         /* insured_count 被保险人数 */
	CommenceDate             null.Int       `json:"CommenceDate,omitempty" db:"commence_date,false,bigint"`                           /* commence_date 起保日期 */
	ExpiryDate               null.Int       `json:"ExpiryDate,omitempty" db:"expiry_date,false,bigint"`                               /* expiry_date 止保日期 */
	Indate                   null.Int       `json:"Indate,omitempty" db:"indate,false,bigint"`                                        /* indate 保险期间 */
	ChargeMode               null.String    `json:"ChargeMode,omitempty" db:"charge_mode,false,character varying"`                    /* charge_mode 计费方式 */
	ModifyType               null.String    `json:"ModifyType,omitempty" db:"modify_type,false,character varying"`                    /* modify_type 修改类型：2:普通修改 4:修改发票抬头 6:增减被保险人 */
	ActivityName             null.String    `json:"ActivityName,omitempty" db:"activity_name,false,character varying"`                /* activity_name 比赛/活动名称 */
	ActivityLocation         null.String    `json:"ActivityLocation,omitempty" db:"activity_location,false,character varying"`        /* activity_location 比赛地点 */
	ActivityDateSet          null.String    `json:"ActivityDateSet,omitempty" db:"activity_date_set,false,character varying"`         /* activity_date_set 具体活动日期（英文逗号分隔多个日期） */
	InsuredType              null.String    `json:"InsuredType,omitempty" db:"insured_type,false,character varying"`                  /* insured_type 参赛人员类型（比赛，学生教师/成年人） */
	SchoolbusCompany         null.String    `json:"SchoolbusCompany,omitempty" db:"schoolbus_company,false,character varying"`        /* schoolbus_company 校车服务单位(校车) */
	GuaranteeItem            types.JSONText `json:"GuaranteeItem,omitempty" db:"guarantee_item,false,jsonb"`                          /* guarantee_item 保障项目 */
	ConfirmGuaranteeStarTime null.Int       `json:"ConfirmGuaranteeStarTime,omitempty" db:"confirm_guarantee_star_time,false,bigint"` /* confirm_guarantee_star_time 确认保障开始时间状态(校车, 食堂) */
	NonCompulsoryStudentNum  null.Int       `json:"NonCompulsoryStudentNum,omitempty" db:"non_compulsory_student_num,false,bigint"`   /* non_compulsory_student_num 非义务教育人数(校方) */
	CompulsoryStudentNum     null.Int       `json:"CompulsoryStudentNum,omitempty" db:"compulsory_student_num,false,bigint"`          /* compulsory_student_num 义务教育人数(校方) */
	DinnerNum                null.Int       `json:"DinnerNum,omitempty" db:"dinner_num,false,bigint"`                                 /* dinner_num 就餐人数(食堂) */
	SchoolEnrolmentTotal     null.Int       `json:"SchoolEnrolmentTotal,omitempty" db:"school_enrolment_total,false,bigint"`          /* school_enrolment_total 注册学生人数(食堂) */
	ShopNum                  null.Int       `json:"ShopNum,omitempty" db:"shop_num,false,bigint"`                                     /* shop_num 小卖铺数量(食堂) */
	CanteenNum               null.Int       `json:"CanteenNum,omitempty" db:"canteen_num,false,bigint"`                               /* canteen_num 食堂数量(食堂) */
	ActivityDesc             null.String    `json:"ActivityDesc,omitempty" db:"activity_desc,false,character varying"`                /* activity_desc 简述(比赛) */
	InvoiceHeader            null.String    `json:"InvoiceHeader,omitempty" db:"invoice_header,false,character varying"`              /* invoice_header 发票抬头 */
	DisputeHandling          null.String    `json:"DisputeHandling,omitempty" db:"dispute_handling,false,character varying"`          /* dispute_handling 争议处理 */
	HaveSuddenDeath          null.Bool      `json:"HaveSuddenDeath,omitempty" db:"have_sudden_death,false,boolean"`                   /* have_sudden_death 启用猝死责任险 */
	PrevPolicyNo             null.String    `json:"PrevPolicyNo,omitempty" db:"prev_policy_no,false,character varying"`               /* prev_policy_no 续保保单号 */
	RevokedPolicyNo          null.String    `json:"RevokedPolicyNo,omitempty" db:"revoked_policy_no,false,character varying"`         /* revoked_policy_no 撤保保单号 */
	PoolName                 null.String    `json:"PoolName,omitempty" db:"pool_name,false,character varying"`                        /* pool_name 游泳池名称 */
	HaveExplosive            null.Bool      `json:"HaveExplosive,omitempty" db:"have_explosive,false,boolean"`                        /* have_explosive 危险易爆 */
	HaveRides                null.Bool      `json:"HaveRides,omitempty" db:"have_rides,false,boolean"`                                /* have_rides 机械性游乐设施 */
	InnerArea                null.Float     `json:"InnerArea,omitempty" db:"inner_area,false,double precision"`                       /* inner_area 室内面积 */
	OuterArea                null.Float     `json:"OuterArea,omitempty" db:"outer_area,false,double precision"`                       /* outer_area 室外面积 */
	TrafficNum               null.Int       `json:"TrafficNum,omitempty" db:"traffic_num,false,integer"`                              /* traffic_num 每日客流量 */
	TemperatureType          null.String    `json:"TemperatureType,omitempty" db:"temperature_type,false,character varying"`          /* temperature_type 常温池 */
	OpenPoolNum              null.Int       `json:"OpenPoolNum,omitempty" db:"open_pool_num,false,smallint"`                          /* open_pool_num 对外数量 */
	HeatedPoolNum            null.Int       `json:"HeatedPoolNum,omitempty" db:"heated_pool_num,false,smallint"`                      /* heated_pool_num 恒温池数量 */
	TrainingPoolNum          null.Int       `json:"TrainingPoolNum,omitempty" db:"training_pool_num,false,smallint"`                  /* training_pool_num 训练池数量 */
	PoolNum                  null.Int       `json:"PoolNum,omitempty" db:"pool_num,false,smallint"`                                   /* pool_num 泳池数量 */
	CustomType               null.String    `json:"CustomType,omitempty" db:"custom_type,false,character varying"`                    /* custom_type 场地使用性质:internal, open, both */
	Same                     null.Bool      `json:"Same,omitempty" db:"same,false,boolean"`                                           /* same 被保险人同投保人 */
	ArbitralAgency           null.String    `json:"ArbitralAgency,omitempty" db:"arbitral_agency,false,character varying"`            /* arbitral_agency 仲裁机构 */
	EndorsementStatus        null.String    `json:"EndorsementStatus,omitempty" db:"endorsement_status,false,character varying"`      /* endorsement_status 批单状态: 00未生成批单 04 已生成批改申请书 08用户上传批改申请书 12管理员上传批单 */
	ApplicationFiles         types.JSONText `json:"ApplicationFiles,omitempty" db:"application_files,false,jsonb"`                    /* application_files 批改申请书 */
	Amount                   null.Float     `json:"Amount,omitempty" db:"amount,false,double precision"`                              /* amount 更正后金额 */
	InsuredGroupByDay        null.Bool      `json:"InsuredGroupByDay,omitempty" db:"insured_group_by_day,false,boolean"`              /* insured_group_by_day 按天录入被保险人 */
	RefusedReason            null.String    `json:"RefusedReason,omitempty" db:"refused_reason,false,character varying"`              /* refused_reason 拒绝理由 */
	PayType                  null.String    `json:"PayType,omitempty" db:"pay_type,false,character varying"`                          /* pay_type 支付方式: 对公转账/在线支付/线下支付 */
	NeedBalance              null.Bool      `json:"NeedBalance,omitempty" db:"need_balance,false,boolean"`                            /* need_balance 需要录入差价 */
	FeeScheme                types.JSONText `json:"FeeScheme,omitempty" db:"fee_scheme,false,jsonb"`                                  /* fee_scheme 计费标准/单价 */
	HaveNegotiatedPrice      null.Bool      `json:"HaveNegotiatedPrice,omitempty" db:"have_negotiated_price,false,boolean"`           /* have_negotiated_price 是否使用协议价 */
	PolicyScheme             types.JSONText `json:"PolicyScheme,omitempty" db:"policy_scheme,false,jsonb"`                            /* policy_scheme 保险方案 */
	PlanID                   null.Int       `json:"PlanID,omitempty" db:"plan_id,false,bigint"`                                       /* plan_id plan_id */
	CorrectLevel             null.String    `json:"CorrectLevel,omitempty" db:"correct_level,false,character varying"`                /* correct_level 更正等级 */
	CorrectLog               types.JSONText `json:"CorrectLog,omitempty" db:"correct_log,false,jsonb"`                                /* correct_log 更正记录 */
	PolicyRegen              null.Bool      `json:"PolicyRegen,omitempty" db:"policy_regen,false,boolean"`                            /* policy_regen 重新生成保单 */
	Files                    types.JSONText `json:"Files,omitempty" db:"files,false,jsonb"`                                           /* files 附加文件 */
	FilesToRemove            null.String    `json:"FilesToRemove,omitempty" db:"files_to_remove,false,character varying"`             /* files_to_remove 待删除文件的digest,逗号分隔  */
	Creator                  null.Int       `json:"Creator,omitempty" db:"creator,false,bigint"`                                      /* creator 创建者用户ID */
	CreateTime               null.Int       `json:"CreateTime,omitempty" db:"create_time,false,bigint"`                               /* create_time 创建时间 */
	Regenerator              null.Int       `json:"Regenerator,omitempty" db:"regenerator,false,bigint"`                              /* regenerator 更新者 */
	UpdateTime               null.Int       `json:"UpdateTime,omitempty" db:"update_time,false,bigint"`                               /* update_time 修改时间 */
	DomainID                 null.Int       `json:"DomainID,omitempty" db:"domain_id,false,bigint"`                                   /* domain_id 数据属主 */
	Addi                     types.JSONText `json:"Addi,omitempty" db:"addi,false,jsonb"`                                             /* addi 附加数据 */
	Remark                   null.String    `json:"Remark,omitempty" db:"remark,false,character varying"`                             /* remark 备注 */
	Status                   null.String    `json:"Status,omitempty" db:"status,false,character varying"`                             /* status 状态,0: 草稿, 2: 受理中，4:同意， 6:拒绝 */
	Filter                                  // build DML where clause
}

//TMistakeCorrectFields full field list for default query
var TMistakeCorrectFields = []string{
	"ID",
	"OrderID",
	"Policyholder",
	"Contact",
	"PolicyholderID",
	"OfficialNameP",
	"IDCardTypeP",
	"IDCardNoP",
	"GenderP",
	"BirthdayP",
	"Insured",
	"InsuredID",
	"OfficialName",
	"IDCardType",
	"IDCardNo",
	"Gender",
	"Birthday",
	"ClearList",
	"InsuredList",
	"InsuredCount",
	"CommenceDate",
	"ExpiryDate",
	"Indate",
	"ChargeMode",
	"ModifyType",
	"ActivityName",
	"ActivityLocation",
	"ActivityDateSet",
	"InsuredType",
	"SchoolbusCompany",
	"GuaranteeItem",
	"ConfirmGuaranteeStarTime",
	"NonCompulsoryStudentNum",
	"CompulsoryStudentNum",
	"DinnerNum",
	"SchoolEnrolmentTotal",
	"ShopNum",
	"CanteenNum",
	"ActivityDesc",
	"InvoiceHeader",
	"DisputeHandling",
	"HaveSuddenDeath",
	"PrevPolicyNo",
	"RevokedPolicyNo",
	"PoolName",
	"HaveExplosive",
	"HaveRides",
	"InnerArea",
	"OuterArea",
	"TrafficNum",
	"TemperatureType",
	"OpenPoolNum",
	"HeatedPoolNum",
	"TrainingPoolNum",
	"PoolNum",
	"CustomType",
	"Same",
	"ArbitralAgency",
	"EndorsementStatus",
	"ApplicationFiles",
	"Amount",
	"InsuredGroupByDay",
	"RefusedReason",
	"PayType",
	"NeedBalance",
	"FeeScheme",
	"HaveNegotiatedPrice",
	"PolicyScheme",
	"PlanID",
	"CorrectLevel",
	"CorrectLog",
	"PolicyRegen",
	"Files",
	"FilesToRemove",
	"Creator",
	"CreateTime",
	"Regenerator",
	"UpdateTime",
	"DomainID",
	"Addi",
	"Remark",
	"Status",
}

//Fields return all fields of struct.
func (r *TMistakeCorrect) Fields() []string {
	return TMistakeCorrectFields
}

//GetTableName return the associated db table name.
func (r *TMistakeCorrect) GetTableName() string {
	var viewNamePattern = regexp.MustCompile(`(?i)^t_v_[a-z0-9_]+$`)
	tableName := "t_mistake_correct"
	if viewNamePattern.MatchString(tableName) {
		return tableName[2:]
	}
	return tableName
}

// Create inserts the TMistakeCorrect to the database.
func (r *TMistakeCorrect) Create(db Queryer) error {
	err := db.QueryRow(
		`INSERT INTO t_mistake_correct (order_id, policyholder, contact, policyholder_id, official_name_p, id_card_type_p, id_card_no_p, gender_p, birthday_p, insured, insured_id, official_name, id_card_type, id_card_no, gender, birthday, clear_list, insured_list, insured_count, commence_date, expiry_date, indate, charge_mode, modify_type, activity_name, activity_location, activity_date_set, insured_type, schoolbus_company, guarantee_item, confirm_guarantee_star_time, non_compulsory_student_num, compulsory_student_num, dinner_num, school_enrolment_total, shop_num, canteen_num, activity_desc, invoice_header, dispute_handling, have_sudden_death, prev_policy_no, revoked_policy_no, pool_name, have_explosive, have_rides, inner_area, outer_area, traffic_num, temperature_type, open_pool_num, heated_pool_num, training_pool_num, pool_num, custom_type, same, arbitral_agency, endorsement_status, application_files, amount, insured_group_by_day, refused_reason, pay_type, need_balance, fee_scheme, have_negotiated_price, policy_scheme, plan_id, correct_level, correct_log, policy_regen, files, files_to_remove, creator, create_time, regenerator, update_time, domain_id, addi, remark, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36, $37, $38, $39, $40, $41, $42, $43, $44, $45, $46, $47, $48, $49, $50, $51, $52, $53, $54, $55, $56, $57, $58, $59, $60, $61, $62, $63, $64, $65, $66, $67, $68, $69, $70, $71, $72, $73, $74, $75, $76, $77, $78, $79, $80, $81) RETURNING id`,
		&r.OrderID, &r.Policyholder, &r.Contact, &r.PolicyholderID, &r.OfficialNameP, &r.IDCardTypeP, &r.IDCardNoP, &r.GenderP, &r.BirthdayP, &r.Insured, &r.InsuredID, &r.OfficialName, &r.IDCardType, &r.IDCardNo, &r.Gender, &r.Birthday, &r.ClearList, &r.InsuredList, &r.InsuredCount, &r.CommenceDate, &r.ExpiryDate, &r.Indate, &r.ChargeMode, &r.ModifyType, &r.ActivityName, &r.ActivityLocation, &r.ActivityDateSet, &r.InsuredType, &r.SchoolbusCompany, &r.GuaranteeItem, &r.ConfirmGuaranteeStarTime, &r.NonCompulsoryStudentNum, &r.CompulsoryStudentNum, &r.DinnerNum, &r.SchoolEnrolmentTotal, &r.ShopNum, &r.CanteenNum, &r.ActivityDesc, &r.InvoiceHeader, &r.DisputeHandling, &r.HaveSuddenDeath, &r.PrevPolicyNo, &r.RevokedPolicyNo, &r.PoolName, &r.HaveExplosive, &r.HaveRides, &r.InnerArea, &r.OuterArea, &r.TrafficNum, &r.TemperatureType, &r.OpenPoolNum, &r.HeatedPoolNum, &r.TrainingPoolNum, &r.PoolNum, &r.CustomType, &r.Same, &r.ArbitralAgency, &r.EndorsementStatus, &r.ApplicationFiles, &r.Amount, &r.InsuredGroupByDay, &r.RefusedReason, &r.PayType, &r.NeedBalance, &r.FeeScheme, &r.HaveNegotiatedPrice, &r.PolicyScheme, &r.PlanID, &r.CorrectLevel, &r.CorrectLog, &r.PolicyRegen, &r.Files, &r.FilesToRemove, &r.Creator, &r.CreateTime, &r.Regenerator, &r.UpdateTime, &r.DomainID, &r.Addi, &r.Remark, &r.Status).Scan(&r.ID)
	if err != nil {
		return errors.Wrap(err, "failed to insert t_mistake_correct")
	}
	return nil
}

// GetTMistakeCorrectByPk select the TMistakeCorrect from the database.
func GetTMistakeCorrectByPk(db Queryer, pk0 null.Int) (*TMistakeCorrect, error) {
	var r TMistakeCorrect
	err := db.QueryRow(
		`SELECT id, order_id, policyholder, contact, policyholder_id, official_name_p, id_card_type_p, id_card_no_p, gender_p, birthday_p, insured, insured_id, official_name, id_card_type, id_card_no, gender, birthday, clear_list, insured_list, insured_count, commence_date, expiry_date, indate, charge_mode, modify_type, activity_name, activity_location, activity_date_set, insured_type, schoolbus_company, guarantee_item, confirm_guarantee_star_time, non_compulsory_student_num, compulsory_student_num, dinner_num, school_enrolment_total, shop_num, canteen_num, activity_desc, invoice_header, dispute_handling, have_sudden_death, prev_policy_no, revoked_policy_no, pool_name, have_explosive, have_rides, inner_area, outer_area, traffic_num, temperature_type, open_pool_num, heated_pool_num, training_pool_num, pool_num, custom_type, same, arbitral_agency, endorsement_status, application_files, amount, insured_group_by_day, refused_reason, pay_type, need_balance, fee_scheme, have_negotiated_price, policy_scheme, plan_id, correct_level, correct_log, policy_regen, files, files_to_remove, creator, create_time, regenerator, update_time, domain_id, addi, remark, status FROM t_mistake_correct WHERE id = $1`,
		pk0).Scan(&r.ID, &r.OrderID, &r.Policyholder, &r.Contact, &r.PolicyholderID, &r.OfficialNameP, &r.IDCardTypeP, &r.IDCardNoP, &r.GenderP, &r.BirthdayP, &r.Insured, &r.InsuredID, &r.OfficialName, &r.IDCardType, &r.IDCardNo, &r.Gender, &r.Birthday, &r.ClearList, &r.InsuredList, &r.InsuredCount, &r.CommenceDate, &r.ExpiryDate, &r.Indate, &r.ChargeMode, &r.ModifyType, &r.ActivityName, &r.ActivityLocation, &r.ActivityDateSet, &r.InsuredType, &r.SchoolbusCompany, &r.GuaranteeItem, &r.ConfirmGuaranteeStarTime, &r.NonCompulsoryStudentNum, &r.CompulsoryStudentNum, &r.DinnerNum, &r.SchoolEnrolmentTotal, &r.ShopNum, &r.CanteenNum, &r.ActivityDesc, &r.InvoiceHeader, &r.DisputeHandling, &r.HaveSuddenDeath, &r.PrevPolicyNo, &r.RevokedPolicyNo, &r.PoolName, &r.HaveExplosive, &r.HaveRides, &r.InnerArea, &r.OuterArea, &r.TrafficNum, &r.TemperatureType, &r.OpenPoolNum, &r.HeatedPoolNum, &r.TrainingPoolNum, &r.PoolNum, &r.CustomType, &r.Same, &r.ArbitralAgency, &r.EndorsementStatus, &r.ApplicationFiles, &r.Amount, &r.InsuredGroupByDay, &r.RefusedReason, &r.PayType, &r.NeedBalance, &r.FeeScheme, &r.HaveNegotiatedPrice, &r.PolicyScheme, &r.PlanID, &r.CorrectLevel, &r.CorrectLog, &r.PolicyRegen, &r.Files, &r.FilesToRemove, &r.Creator, &r.CreateTime, &r.Regenerator, &r.UpdateTime, &r.DomainID, &r.Addi, &r.Remark, &r.Status)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select t_mistake_correct")
	}
	return &r, nil
}

/*TMsg 消息：包含新闻，私信，广告，通知等 represents kuser.t_msg */
type TMsg struct {
	ID          null.Int       `json:"ID,omitempty" db:"id,true,integer"`                        /* id 参数编号 */
	Author      null.String    `json:"Author,omitempty" db:"author,false,character varying"`     /* author 作者 */
	Title       null.String    `json:"Title,omitempty" db:"title,false,character varying"`       /* title 标题 */
	Subtitle    null.String    `json:"Subtitle,omitempty" db:"subtitle,false,character varying"` /* subtitle 副标题 */
	Keyword     null.String    `json:"Keyword,omitempty" db:"keyword,false,character varying"`   /* keyword 关键字 */
	Belong      null.Int       `json:"Belong,omitempty" db:"belong,false,bigint"`                /* belong 属于 */
	Channel     types.JSONText `json:"Channel,omitempty" db:"channel,false,jsonb"`               /* channel 频道 */
	Type        types.JSONText `json:"Type,omitempty" db:"type,false,jsonb"`                     /* type 内容类型：搞笑，新闻， */
	Domain      types.JSONText `json:"Domain,omitempty" db:"domain,false,jsonb"`                 /* domain 领域：教育/游戏/电子等 */
	Quality     null.Int       `json:"Quality,omitempty" db:"quality,false,integer"`             /* quality 内容质量 */
	Viewed      null.Int       `json:"Viewed,omitempty" db:"viewed,false,integer"`               /* viewed 阅读次数 */
	Score       types.JSONText `json:"Score,omitempty" db:"score,false,jsonb"`                   /* score 读者评分 */
	Prosecute   types.JSONText `json:"Prosecute,omitempty" db:"prosecute,false,jsonb"`           /* prosecute 举报 */
	AssentNum   null.Int       `json:"AssentNum,omitempty" db:"assent_num,false,integer"`        /* assent_num 赞同数 */
	OpposeNum   null.Int       `json:"OpposeNum,omitempty" db:"oppose_num,false,integer"`        /* oppose_num 反对数 */
	Source      null.String    `json:"Source,omitempty" db:"source,false,character varying"`     /* source 来源 */
	Tags        null.String    `json:"Tags,omitempty" db:"tags,false,character varying"`         /* tags 标签 */
	FacePicNum  null.Int       `json:"FacePicNum,omitempty" db:"face_pic_num,false,smallint"`    /* face_pic_num 封面图片数 */
	Content     types.JSONText `json:"Content,omitempty" db:"content,false,jsonb"`               /* content 内容 */
	Files       types.JSONText `json:"Files,omitempty" db:"files,false,jsonb"`                   /* files 附加文件 */
	Creator     null.Int       `json:"Creator,omitempty" db:"creator,false,bigint"`              /* creator 本数据创建者 */
	CreateTime  null.Int       `json:"CreateTime,omitempty" db:"create_time,false,bigint"`       /* create_time 生成时间 */
	Regenerator null.Int       `json:"Regenerator,omitempty" db:"regenerator,false,bigint"`      /* regenerator 更新者 */
	UpdateTime  null.Int       `json:"UpdateTime,omitempty" db:"update_time,false,bigint"`       /* update_time 修改时间 */
	DomainID    null.Int       `json:"DomainID,omitempty" db:"domain_id,false,bigint"`           /* domain_id 数据隶属 */
	Addi        types.JSONText `json:"Addi,omitempty" db:"addi,false,jsonb"`                     /* addi 附加信息 */
	Remark      null.String    `json:"Remark,omitempty" db:"remark,false,character varying"`     /* remark 备注 */
	Status      null.String    `json:"Status,omitempty" db:"status,false,character varying"`     /* status 状态，00：草稿，01：有效，02：作废 */
	Filter                     // build DML where clause
}

//TMsgFields full field list for default query
var TMsgFields = []string{
	"ID",
	"Author",
	"Title",
	"Subtitle",
	"Keyword",
	"Belong",
	"Channel",
	"Type",
	"Domain",
	"Quality",
	"Viewed",
	"Score",
	"Prosecute",
	"AssentNum",
	"OpposeNum",
	"Source",
	"Tags",
	"FacePicNum",
	"Content",
	"Files",
	"Creator",
	"CreateTime",
	"Regenerator",
	"UpdateTime",
	"DomainID",
	"Addi",
	"Remark",
	"Status",
}

//Fields return all fields of struct.
func (r *TMsg) Fields() []string {
	return TMsgFields
}

//GetTableName return the associated db table name.
func (r *TMsg) GetTableName() string {
	var viewNamePattern = regexp.MustCompile(`(?i)^t_v_[a-z0-9_]+$`)
	tableName := "t_msg"
	if viewNamePattern.MatchString(tableName) {
		return tableName[2:]
	}
	return tableName
}

// Create inserts the TMsg to the database.
func (r *TMsg) Create(db Queryer) error {
	err := db.QueryRow(
		`INSERT INTO t_msg (author, title, subtitle, keyword, belong, channel, type, domain, quality, viewed, score, prosecute, assent_num, oppose_num, source, tags, face_pic_num, content, files, creator, create_time, regenerator, update_time, domain_id, addi, remark, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27) RETURNING id`,
		&r.Author, &r.Title, &r.Subtitle, &r.Keyword, &r.Belong, &r.Channel, &r.Type, &r.Domain, &r.Quality, &r.Viewed, &r.Score, &r.Prosecute, &r.AssentNum, &r.OpposeNum, &r.Source, &r.Tags, &r.FacePicNum, &r.Content, &r.Files, &r.Creator, &r.CreateTime, &r.Regenerator, &r.UpdateTime, &r.DomainID, &r.Addi, &r.Remark, &r.Status).Scan(&r.ID)
	if err != nil {
		return errors.Wrap(err, "failed to insert t_msg")
	}
	return nil
}

// GetTMsgByPk select the TMsg from the database.
func GetTMsgByPk(db Queryer, pk0 null.Int) (*TMsg, error) {
	var r TMsg
	err := db.QueryRow(
		`SELECT id, author, title, subtitle, keyword, belong, channel, type, domain, quality, viewed, score, prosecute, assent_num, oppose_num, source, tags, face_pic_num, content, files, creator, create_time, regenerator, update_time, domain_id, addi, remark, status FROM t_msg WHERE id = $1`,
		pk0).Scan(&r.ID, &r.Author, &r.Title, &r.Subtitle, &r.Keyword, &r.Belong, &r.Channel, &r.Type, &r.Domain, &r.Quality, &r.Viewed, &r.Score, &r.Prosecute, &r.AssentNum, &r.OpposeNum, &r.Source, &r.Tags, &r.FacePicNum, &r.Content, &r.Files, &r.Creator, &r.CreateTime, &r.Regenerator, &r.UpdateTime, &r.DomainID, &r.Addi, &r.Remark, &r.Status)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select t_msg")
	}
	return &r, nil
}

/*TNegotiatedPrice 协议价表 represents kuser.t_negotiated_price */
type TNegotiatedPrice struct {
	ID              null.Int       `json:"ID,omitempty" db:"id,true,integer"`                             /* id 议价表设置编号 */
	Keyword         null.String    `json:"Keyword,omitempty" db:"keyword,false,character varying"`        /* keyword 关键词 */
	CommenceDate    null.Int       `json:"CommenceDate,omitempty" db:"commence_date,false,bigint"`        /* commence_date 开始日期（纳秒） */
	Location        null.String    `json:"Location,omitempty" db:"location,false,character varying"`      /* location 地点 */
	Province        null.String    `json:"Province,omitempty" db:"province,false,character varying"`      /* province 省 */
	City            null.String    `json:"City,omitempty" db:"city,false,character varying"`              /* city 市 */
	District        null.String    `json:"District,omitempty" db:"district,false,character varying"`      /* district 区 */
	PriceType       string         `json:"PriceType,omitempty" db:"price_type,false,character varying"`   /* price_type 议价类型（协议价/会议价） */
	Price           null.Int       `json:"Price,omitempty" db:"price,false,integer"`                      /* price 议价价格 */
	InsuranceTypeID null.Int       `json:"InsuranceTypeID,omitempty" db:"insurance_type_id,false,bigint"` /* insurance_type_id 保险类型 */
	MatchTimes      null.Int       `json:"MatchTimes,omitempty" db:"match_times,false,integer"`           /* match_times 匹配次数 */
	Indate          null.Int       `json:"Indate,omitempty" db:"indate,false,bigint"`                     /* indate 保险期间 */
	Creator         null.Int       `json:"Creator,omitempty" db:"creator,false,bigint"`                   /* creator 创建者用户ID */
	CreateTime      null.Int       `json:"CreateTime,omitempty" db:"create_time,false,bigint"`            /* create_time 创建时间 */
	Regenerator     null.Int       `json:"Regenerator,omitempty" db:"regenerator,false,bigint"`           /* regenerator 更新者 */
	UpdateTime      null.Int       `json:"UpdateTime,omitempty" db:"update_time,false,bigint"`            /* update_time 更新时间 */
	Updator         null.Int       `json:"Updator,omitempty" db:"updator,false,bigint"`                   /* updator 更新者ID */
	DomainID        null.Int       `json:"DomainID,omitempty" db:"domain_id,false,bigint"`                /* domain_id 数据属主 */
	Addi            types.JSONText `json:"Addi,omitempty" db:"addi,false,jsonb"`                          /* addi 附加数据 */
	Remark          null.String    `json:"Remark,omitempty" db:"remark,false,character varying"`          /* remark 备注 */
	Status          null.String    `json:"Status,omitempty" db:"status,false,character varying"`          /* status 0:有效, 2: 删除 */
	Filter                         // build DML where clause
}

//TNegotiatedPriceFields full field list for default query
var TNegotiatedPriceFields = []string{
	"ID",
	"Keyword",
	"CommenceDate",
	"Location",
	"Province",
	"City",
	"District",
	"PriceType",
	"Price",
	"InsuranceTypeID",
	"MatchTimes",
	"Indate",
	"Creator",
	"CreateTime",
	"Regenerator",
	"UpdateTime",
	"Updator",
	"DomainID",
	"Addi",
	"Remark",
	"Status",
}

//Fields return all fields of struct.
func (r *TNegotiatedPrice) Fields() []string {
	return TNegotiatedPriceFields
}

//GetTableName return the associated db table name.
func (r *TNegotiatedPrice) GetTableName() string {
	var viewNamePattern = regexp.MustCompile(`(?i)^t_v_[a-z0-9_]+$`)
	tableName := "t_negotiated_price"
	if viewNamePattern.MatchString(tableName) {
		return tableName[2:]
	}
	return tableName
}

// Create inserts the TNegotiatedPrice to the database.
func (r *TNegotiatedPrice) Create(db Queryer) error {
	err := db.QueryRow(
		`INSERT INTO t_negotiated_price (keyword, commence_date, location, province, city, district, price_type, price, insurance_type_id, match_times, indate, creator, create_time, regenerator, update_time, updator, domain_id, addi, remark, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20) RETURNING id`,
		&r.Keyword, &r.CommenceDate, &r.Location, &r.Province, &r.City, &r.District, &r.PriceType, &r.Price, &r.InsuranceTypeID, &r.MatchTimes, &r.Indate, &r.Creator, &r.CreateTime, &r.Regenerator, &r.UpdateTime, &r.Updator, &r.DomainID, &r.Addi, &r.Remark, &r.Status).Scan(&r.ID)
	if err != nil {
		return errors.Wrap(err, "failed to insert t_negotiated_price")
	}
	return nil
}

// GetTNegotiatedPriceByPk select the TNegotiatedPrice from the database.
func GetTNegotiatedPriceByPk(db Queryer, pk0 null.Int) (*TNegotiatedPrice, error) {
	var r TNegotiatedPrice
	err := db.QueryRow(
		`SELECT id, keyword, commence_date, location, province, city, district, price_type, price, insurance_type_id, match_times, indate, creator, create_time, regenerator, update_time, updator, domain_id, addi, remark, status FROM t_negotiated_price WHERE id = $1`,
		pk0).Scan(&r.ID, &r.Keyword, &r.CommenceDate, &r.Location, &r.Province, &r.City, &r.District, &r.PriceType, &r.Price, &r.InsuranceTypeID, &r.MatchTimes, &r.Indate, &r.Creator, &r.CreateTime, &r.Regenerator, &r.UpdateTime, &r.Updator, &r.DomainID, &r.Addi, &r.Remark, &r.Status)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select t_negotiated_price")
	}
	return &r, nil
}

/*TOrder 订单 represents kuser.t_order */
type TOrder struct {
	ID                      null.Int       `json:"ID,omitempty" db:"id,true,integer"`                                              /* id 订单编号 */
	TradeNo                 null.String    `json:"TradeNo,omitempty" db:"trade_no,false,character varying"`                        /* trade_no 外部订单号 */
	PayOrderNo              null.String    `json:"PayOrderNo,omitempty" db:"pay_order_no,false,character varying"`                 /* pay_order_no 外部支付订单号 */
	TransactionID           null.String    `json:"TransactionID,omitempty" db:"transaction_id,false,character varying"`            /* transaction_id 支付平台订单号 */
	Batch                   null.String    `json:"Batch,omitempty" db:"batch,false,character varying"`                             /* batch 批次编号 */
	PayTime                 null.Int       `json:"PayTime,omitempty" db:"pay_time,false,bigint"`                                   /* pay_time 支付时间 */
	PayType                 null.String    `json:"PayType,omitempty" db:"pay_type,false,character varying"`                        /* pay_type 支付方式: 对公转账/在线支付/线下支付 */
	PayChannel              null.String    `json:"PayChannel,omitempty" db:"pay_channel,false,character varying"`                  /* pay_channel 校快保，泰合，近邻，人保，太平洋保险 */
	PayName                 null.String    `json:"PayName,omitempty" db:"pay_name,false,character varying"`                        /* pay_name 支付方式名称 */
	PayAccountInfo          types.JSONText `json:"PayAccountInfo,omitempty" db:"pay_account_info,false,jsonb"`                     /* pay_account_info 在线支付信息（支付后填写用于核查） */
	Refundable              null.Bool      `json:"Refundable,omitempty" db:"refundable,false,boolean"`                             /* refundable 是否支持在线退款 */
	UnitPrice               null.Float     `json:"UnitPrice,omitempty" db:"unit_price,false,double precision"`                     /* unit_price 单价 */
	RefundDesc              null.String    `json:"RefundDesc,omitempty" db:"refund_desc,false,character varying"`                  /* refund_desc 退款原因 */
	Amount                  null.Float     `json:"Amount,omitempty" db:"amount,false,double precision"`                            /* amount 应收金额 */
	ActualAmount            null.Float     `json:"ActualAmount,omitempty" db:"actual_amount,false,double precision"`               /* actual_amount 实收金额 */
	Balance                 null.Float     `json:"Balance,omitempty" db:"balance,false,double precision"`                          /* balance 累计差额 */
	BalanceList             types.JSONText `json:"BalanceList,omitempty" db:"balance_list,false,jsonb"`                            /* balance_list 差额详表 */
	InsureOrderNo           null.String    `json:"InsureOrderNo,omitempty" db:"insure_order_no,false,character varying"`           /* insure_order_no 外部保险系统订单号 */
	RefundNo                null.String    `json:"RefundNo,omitempty" db:"refund_no,false,character varying"`                      /* refund_no 退款单号 */
	Refund                  null.Float     `json:"Refund,omitempty" db:"refund,false,double precision"`                            /* refund (待)退款金额 */
	RefundTime              null.Int       `json:"RefundTime,omitempty" db:"refund_time,false,bigint"`                             /* refund_time 退款时间 */
	ConfirmRefund           null.Bool      `json:"ConfirmRefund,omitempty" db:"confirm_refund,false,boolean"`                      /* confirm_refund 确认退款 */
	AgencyID                null.Int       `json:"AgencyID,omitempty" db:"agency_id,false,bigint"`                                 /* agency_id 表示用于统计的机构编号, 因为目前org_id,policyholder_id,insured_id都用来表示机构编号 */
	OrgID                   null.Int       `json:"OrgID,omitempty" db:"org_id,false,bigint"`                                       /* org_id 关联机构编号 */
	OrgManagerID            null.Int       `json:"OrgManagerID,omitempty" db:"org_manager_id,false,bigint"`                        /* org_manager_id 关联机构管理人 */
	InsuranceType           string         `json:"InsuranceType,omitempty" db:"insurance_type,false,character varying"`            /* insurance_type 保险类型: 学生意外伤害险，活动/比赛险(旅游险),食品卫生责任险，教工责任险,校方责任险,实习生责任险,校车责任险,游泳池责任险 */
	InsuranceTypeID         null.Int       `json:"InsuranceTypeID,omitempty" db:"insurance_type_id,false,bigint"`                  /* insurance_type_id 保险产品编码，对应t_insurance_types.id */
	InsurancePoliceID       null.Int       `json:"InsurancePoliceID,omitempty" db:"insurance_police_id,false,bigint"`              /* insurance_police_id 保单编号[学意险] */
	PlanID                  null.Int       `json:"PlanID,omitempty" db:"plan_id,false,bigint"`                                     /* plan_id 保险方案编码，对应t_insurance_types.id */
	PlanName                null.String    `json:"PlanName,omitempty" db:"plan_name,false,character varying"`                      /* plan_name 保险方案名称(前端暂存) */
	Insurer                 null.String    `json:"Insurer,omitempty" db:"insurer,false,character varying"`                         /* insurer 保险方案承保公司(前端暂存) */
	PolicyScheme            types.JSONText `json:"PolicyScheme,omitempty" db:"policy_scheme,false,jsonb"`                          /* policy_scheme 保险方案 */
	PolicyDoc               null.String    `json:"PolicyDoc,omitempty" db:"policy_doc,false,character varying"`                    /* policy_doc 保险条款 */
	ActivityName            null.String    `json:"ActivityName,omitempty" db:"activity_name,false,character varying"`              /* activity_name 活动名称 */
	ActivityCategory        null.String    `json:"ActivityCategory,omitempty" db:"activity_category,false,character varying"`      /* activity_category 活动类型 */
	ActivityDesc            null.String    `json:"ActivityDesc,omitempty" db:"activity_desc,false,character varying"`              /* activity_desc 活动描述 */
	ActivityLocation        null.String    `json:"ActivityLocation,omitempty" db:"activity_location,false,character varying"`      /* activity_location 活动地点 */
	ActivityDateSet         null.String    `json:"ActivityDateSet,omitempty" db:"activity_date_set,false,character varying"`       /* activity_date_set 具体活动日期(英文逗号隔开) */
	CopiesNum               null.Int       `json:"CopiesNum,omitempty" db:"copies_num,false,smallint"`                             /* copies_num 订单份数 */
	InsuredCount            null.Int       `json:"InsuredCount,omitempty" db:"insured_count,false,smallint"`                       /* insured_count 总数量/保障人数/车辆数 */
	CompulsoryStudentNum    null.Int       `json:"CompulsoryStudentNum,omitempty" db:"compulsory_student_num,false,bigint"`        /* compulsory_student_num 义务教育学生人数（校方） */
	NonCompulsoryStudentNum null.Int       `json:"NonCompulsoryStudentNum,omitempty" db:"non_compulsory_student_num,false,bigint"` /* non_compulsory_student_num 非义务教育人数（校方） */
	Contact                 types.JSONText `json:"Contact,omitempty" db:"contact,false,jsonb"`                                     /* contact 联系人 */
	FeeScheme               types.JSONText `json:"FeeScheme,omitempty" db:"fee_scheme,false,jsonb"`                                /* fee_scheme 计费标准/单价 */
	CarServiceTarget        null.String    `json:"CarServiceTarget,omitempty" db:"car_service_target,false,character varying"`     /* car_service_target 校车服务对象 */
	Policyholder            types.JSONText `json:"Policyholder,omitempty" db:"policyholder,false,jsonb"`                           /* policyholder 投保人 */
	PolicyholderType        null.String    `json:"PolicyholderType,omitempty" db:"policyholder_type,false,character varying"`      /* policyholder_type 投保人类型：个人，机构 */
	PolicyholderID          null.Int       `json:"PolicyholderID,omitempty" db:"policyholder_id,false,bigint"`                     /* policyholder_id 投保人编号 */
	Same                    null.Bool      `json:"Same,omitempty" db:"same,false,boolean"`                                         /* same 投保人与被保险人是同一人 */
	Relation                null.String    `json:"Relation,omitempty" db:"relation,false,character varying"`                       /* relation 投保人与被保险人关系 */
	Insured                 types.JSONText `json:"Insured,omitempty" db:"insured,false,jsonb"`                                     /* insured 被保险人 */
	InsuredID               null.Int       `json:"InsuredID,omitempty" db:"insured_id,false,bigint"`                               /* insured_id 被保险人编号 */
	HealthSurvey            types.JSONText `json:"HealthSurvey,omitempty" db:"health_survey,false,jsonb"`                          /* health_survey 健康调查结果 */
	OrgName                 null.String    `json:"OrgName,omitempty" db:"org_name,false,character varying"`                        /* org_name 学校名称，用户输入的需要新建的学校名称 */
	OrgCategory             null.String    `json:"OrgCategory,omitempty" db:"org_category,false,character varying"`                /* org_category 学校类别: 幼儿园、小学等 */
	HaveInsuredList         null.Bool      `json:"HaveInsuredList,omitempty" db:"have_insured_list,false,boolean"`                 /* have_insured_list 有被保险对象清单 */
	InsuredGroupByDay       null.Bool      `json:"InsuredGroupByDay,omitempty" db:"insured_group_by_day,false,boolean"`            /* insured_group_by_day 被保险对象按日期分组 */
	InsuredType             null.String    `json:"InsuredType,omitempty" db:"insured_type,false,character varying"`                /* insured_type 被保险人类型: 学生，非学生 */
	InsuredList             types.JSONText `json:"InsuredList,omitempty" db:"insured_list,false,jsonb"`                            /* insured_list 被保险对象清单 */
	CommenceDate            null.Int       `json:"CommenceDate,omitempty" db:"commence_date,false,bigint"`                         /* commence_date 起保日(毫秒) */
	ExpiryDate              null.Int       `json:"ExpiryDate,omitempty" db:"expiry_date,false,bigint"`                             /* expiry_date 止保日(毫秒) */
	Indate                  null.Int       `json:"Indate,omitempty" db:"indate,false,bigint"`                                      /* indate 有效期(天) */
	Sign                    null.String    `json:"Sign,omitempty" db:"sign,false,character varying"`                               /* sign 用户签名 */
	Jurisdiction            null.String    `json:"Jurisdiction,omitempty" db:"jurisdiction,false,character varying"`               /* jurisdiction 司法管辖权 */
	DisputeHandling         null.String    `json:"DisputeHandling,omitempty" db:"dispute_handling,false,character varying"`        /* dispute_handling 争议处理 */
	PrevPolicyNo            null.String    `json:"PrevPolicyNo,omitempty" db:"prev_policy_no,false,character varying"`             /* prev_policy_no 续保保单号 */
	InsureBase              null.String    `json:"InsureBase,omitempty" db:"insure_base,false,character varying"`                  /* insure_base 承保基础 */
	BlanketInsureCode       null.String    `json:"BlanketInsureCode,omitempty" db:"blanket_insure_code,false,character varying"`   /* blanket_insure_code 统保代码 */
	CustomType              null.String    `json:"CustomType,omitempty" db:"custom_type,false,character varying"`                  /* custom_type 场地使用性质:internal, open, both */
	TrainProjects           null.String    `json:"TrainProjects,omitempty" db:"train_projects,false,character varying"`            /* train_projects 训练项目 */
	BusinessLocations       types.JSONText `json:"BusinessLocations,omitempty" db:"business_locations,false,jsonb"`                /* business_locations 承保地址/区域范围/游泳池场地地址 */
	OpenPoolNum             null.Int       `json:"OpenPoolNum,omitempty" db:"open_pool_num,false,smallint"`                        /* open_pool_num 对外开放游泳池数量 */
	HeatedPoolNum           null.Int       `json:"HeatedPoolNum,omitempty" db:"heated_pool_num,false,smallint"`                    /* heated_pool_num 恒温游泳池数量 */
	TrainingPoolNum         null.Int       `json:"TrainingPoolNum,omitempty" db:"training_pool_num,false,smallint"`                /* training_pool_num 培训游泳池数量 */
	PoolNum                 null.Int       `json:"PoolNum,omitempty" db:"pool_num,false,smallint"`                                 /* pool_num 游泳池数量 */
	DinnerNum               null.Int       `json:"DinnerNum,omitempty" db:"dinner_num,false,integer"`                              /* dinner_num 用餐人数 */
	HaveDinnerNum           null.Bool      `json:"HaveDinnerNum,omitempty" db:"have_dinner_num,false,boolean"`                     /* have_dinner_num 是否开启就餐人数 */
	CanteenNum              null.Int       `json:"CanteenNum,omitempty" db:"canteen_num,false,integer"`                            /* canteen_num 食堂个数 */
	ShopNum                 null.Int       `json:"ShopNum,omitempty" db:"shop_num,false,integer"`                                  /* shop_num 商店个数 */
	HaveRides               null.Bool      `json:"HaveRides,omitempty" db:"have_rides,false,boolean"`                              /* have_rides 营业场所是否有游泳池外游乐设施、机械性游乐设施等 */
	HaveExplosive           null.Bool      `json:"HaveExplosive,omitempty" db:"have_explosive,false,boolean"`                      /* have_explosive 营业场所是否有制造、销售、储存易燃易爆危险品 */
	InnerArea               null.Float     `json:"InnerArea,omitempty" db:"inner_area,false,double precision"`                     /* inner_area 室内面积 */
	OuterArea               null.Float     `json:"OuterArea,omitempty" db:"outer_area,false,double precision"`                     /* outer_area 室外面积 */
	PoolName                null.String    `json:"PoolName,omitempty" db:"pool_name,false,character varying"`                      /* pool_name 游泳池名称(英文逗号分隔)  */
	ArbitralAgency          null.String    `json:"ArbitralAgency,omitempty" db:"arbitral_agency,false,character varying"`          /* arbitral_agency 仲裁机构 */
	TrafficNum              null.Int       `json:"TrafficNum,omitempty" db:"traffic_num,false,integer"`                            /* traffic_num 每日客流量（人） */
	TemperatureType         null.String    `json:"TemperatureType,omitempty" db:"temperature_type,false,character varying"`        /* temperature_type 泳池性质:恒温、常温 */
	IsIndoor                null.String    `json:"IsIndoor,omitempty" db:"is_indoor,false,character varying"`                      /* is_indoor 泳池特性:室内、室外 */
	Extra                   types.JSONText `json:"Extra,omitempty" db:"extra,false,jsonb"`                                         /* extra 附加信息:
	附加条款
	企业经营描述
	相关保险情况
	保险公司提示
	保险销售事项确认书
	保险公司信息：经办人/工号、代理点代码、展业方式
	产险销售人员：姓名、职业证号 */
	BankAccount             types.JSONText `json:"BankAccount,omitempty" db:"bank_account,false,jsonb"`                                      /* bank_account 对公帐号信息：户名、所在银行、账号 */
	PayContact              null.String    `json:"PayContact,omitempty" db:"pay_contact,false,character varying"`                            /* pay_contact 线下支付联系人：微信二维码，base64 */
	SuddenDeathTerms        null.String    `json:"SuddenDeathTerms,omitempty" db:"sudden_death_terms,false,character varying"`               /* sudden_death_terms 猝死条款内容：附加猝死保险责任每人限额5万元，累计限额5万元。附加猝死责任保险条款（经法院判决、仲裁机构裁决或根据县级以上政府及县级以上政府有关部门的行政决定书或者调解证明等材料，需由被保险人承担的经济赔偿责任，由保险人负责赔偿） */
	SpecAgreement           null.String    `json:"SpecAgreement,omitempty" db:"spec_agreement,false,character varying"`                      /* spec_agreement 特别约定 */
	HaveNegotiatedPrice     null.Bool      `json:"HaveNegotiatedPrice,omitempty" db:"have_negotiated_price,false,boolean"`                   /* have_negotiated_price 是否使用协议价 */
	HaveRenewalReminder     null.Bool      `json:"HaveRenewalReminder,omitempty" db:"have_renewal_reminder,false,boolean"`                   /* have_renewal_reminder 是否有续保通知 */
	LockStatus              null.String    `json:"LockStatus,omitempty" db:"lock_status,false,character varying"`                            /* lock_status 锁定状态:0(或留空):未锁定 2:未解锁 4:已解锁 */
	InsuranceCompany        null.String    `json:"InsuranceCompany,omitempty" db:"insurance_company,false,character varying"`                /* insurance_company 承保公司 */
	InsuranceCompanyAccount null.String    `json:"InsuranceCompanyAccount,omitempty" db:"insurance_company_account,false,character varying"` /* insurance_company_account 承保公司账号 */
	ChargeMode              null.String    `json:"ChargeMode,omitempty" db:"charge_mode,false,character varying"`                            /* charge_mode 购买方式：按天购买，按月购买。学意险模式
	0: 团单, 特定起保、止保时间,
	2: 团单, 非12个保障时间，有起保时间
	4: 团单, 1年保障时间，指定起保日期
	6: 团单, 1年保障时间，不指定起保日期
	8: 散单, 1年保障时间，不指定起保日期 */
	CanRevokeOrder     null.Bool   `json:"CanRevokeOrder,omitempty" db:"can_revoke_order,false,boolean"`         /* can_revoke_order 是否允许撤销订单 */
	CanPublicTransfers null.Bool   `json:"CanPublicTransfers,omitempty" db:"can_public_transfers,false,boolean"` /* can_public_transfers 是否允许对公转账 */
	IsReminder         null.Bool   `json:"IsReminder,omitempty" db:"is_reminder,false,boolean"`                  /* is_reminder 是否开启自动催款 */
	GroundNum          null.Int    `json:"GroundNum,omitempty" db:"ground_num,false,smallint"`                   /* ground_num 场地个数 */
	RemindersNum       null.Int    `json:"RemindersNum,omitempty" db:"reminders_num,false,smallint"`             /* reminders_num 催款次数 */
	ReminderTimes      null.String `json:"ReminderTimes,omitempty" db:"reminder_times,false,character varying"`  /* reminder_times 催款时间,用英文逗号分隔 */
	RefusedReason      null.String `json:"RefusedReason,omitempty" db:"refused_reason,false,character varying"`  /* refused_reason 拒保理由 */
	UnpaidReason       null.String `json:"UnpaidReason,omitempty" db:"unpaid_reason,false,character varying"`    /* unpaid_reason 未收款理由 */
	AdminReceived      null.Bool   `json:"AdminReceived,omitempty" db:"admin_received,false,boolean"`            /* admin_received 管理员已收件 */
	UserReceived       null.Bool   `json:"UserReceived,omitempty" db:"user_received,false,boolean"`              /* user_received 用户已收件 */
	HaveSuddenDeath    null.Bool   `json:"HaveSuddenDeath,omitempty" db:"have_sudden_death,false,boolean"`       /* have_sudden_death 是否开启猝死责任险 */
	HaveConfirmDate    null.Bool   `json:"HaveConfirmDate,omitempty" db:"have_confirm_date,false,boolean"`       /* have_confirm_date 是否确认保障时间 */
	IsInvoice          null.Bool   `json:"IsInvoice,omitempty" db:"is_invoice,false,boolean"`                    /* is_invoice 是否开具发票, false|0: 未开具发票, true|1: 已开具发票 */
	InvVisible         null.String `json:"InvVisible,omitempty" db:"inv_visible,false,character varying"`        /* inv_visible 用户是否可见发票, 0: 发票用户不可见, 2: 发票用户可见 */
	InvBorrow          null.String `json:"InvBorrow,omitempty" db:"inv_borrow,false,character varying"`          /* inv_borrow 00: 未生成发票, 30: 可预借发票, 34: 已生成预借发票申请函, 38: 用户已下载预借发票申请函, 42: 用户已上传盖章预借发票申请函 */
	InvTitle           null.String `json:"InvTitle,omitempty" db:"inv_title,false,character varying"`            /* inv_title 发票抬头修改状态, 16: 申请改发票抬头, 20: 已上传新发票, 24: 已下载新发票 */
	Traits             null.String `json:"Traits,omitempty" db:"traits,false,character varying"`                 /* traits 特殊订单，标志为字符串数组，值域为: {"allowIntraday","ignoreAmountLimit","ignorePayDeadline","ignoreOrderDeadline","ignoreAgeLimit"},
	allowIntraday:当天起保,ignoreAmountLimit:允许小于三人投保,ignorePayDeadline: 允许超时支付,ignoreOrderDeadline：允许超过截止时间投保,ignoreAgeLimit：允许超龄投保 */
	Files       types.JSONText `json:"Files,omitempty" db:"files,false,jsonb"`                          /* files 附加文件 */
	InvStatus   null.String    `json:"InvStatus,omitempty" db:"inv_status,false,character varying"`     /* inv_status 发票状态, 00: 未生成发票, 04: 发票已上传, 08: 发票已下载, 12: 发票已快递 */
	OrderStatus null.String    `json:"OrderStatus,omitempty" db:"order_status,false,character varying"` /* order_status 订单状态, 00: 草稿, 04: 用户投保, 08: 申请议价, 12: 拒保, 16: 可支付, 18: 开始支付, 20: 已支付, 24: 退保, 28: 作废 */
	UpdStatus   null.String    `json:"UpdStatus,omitempty" db:"upd_status,false,character varying"`     /* upd_status 订单更正状态, 00: 未更正, 02: 用户撤消申请, 04: 用户申请更正, 08: 接受更正,16: 更新订单, 20: 生成批改申请书, 24: 用户已下载申请书, 28: 用户已上传盖章批改申请书, 36: 管理员上传批改单, 40: 用户已下载批改单, 44: 批改单已快递给用户, 48: 申请被拒绝 */
	Creator     null.Int       `json:"Creator,omitempty" db:"creator,false,bigint"`                     /* creator 创建者用户ID */
	CreateTime  null.Int       `json:"CreateTime,omitempty" db:"create_time,false,bigint"`              /* create_time 创建时间 */
	Regenerator null.Int       `json:"Regenerator,omitempty" db:"regenerator,false,bigint"`             /* regenerator 更新者 */
	UpdateTime  null.Int       `json:"UpdateTime,omitempty" db:"update_time,false,bigint"`              /* update_time 更新时间 */
	Updator     null.Int       `json:"Updator,omitempty" db:"updator,false,bigint"`                     /* updator 更新者 */
	DomainID    null.Int       `json:"DomainID,omitempty" db:"domain_id,false,bigint"`                  /* domain_id 数据属主 */
	Addi        types.JSONText `json:"Addi,omitempty" db:"addi,false,jsonb"`                            /* addi 附加数据 */
	Remark      null.String    `json:"Remark,omitempty" db:"remark,false,character varying"`            /* remark 备注 */
	Status      null.String    `json:"Status,omitempty" db:"status,false,character varying"`            /* status 0: 未支付, 2: 已支付，4: 已生成保单, 6: 已作废 */
	Filter                     // build DML where clause
}

//TOrderFields full field list for default query
var TOrderFields = []string{
	"ID",
	"TradeNo",
	"PayOrderNo",
	"TransactionID",
	"Batch",
	"PayTime",
	"PayType",
	"PayChannel",
	"PayName",
	"PayAccountInfo",
	"Refundable",
	"UnitPrice",
	"RefundDesc",
	"Amount",
	"ActualAmount",
	"Balance",
	"BalanceList",
	"InsureOrderNo",
	"RefundNo",
	"Refund",
	"RefundTime",
	"ConfirmRefund",
	"AgencyID",
	"OrgID",
	"OrgManagerID",
	"InsuranceType",
	"InsuranceTypeID",
	"InsurancePoliceID",
	"PlanID",
	"PlanName",
	"Insurer",
	"PolicyScheme",
	"PolicyDoc",
	"ActivityName",
	"ActivityCategory",
	"ActivityDesc",
	"ActivityLocation",
	"ActivityDateSet",
	"CopiesNum",
	"InsuredCount",
	"CompulsoryStudentNum",
	"NonCompulsoryStudentNum",
	"Contact",
	"FeeScheme",
	"CarServiceTarget",
	"Policyholder",
	"PolicyholderType",
	"PolicyholderID",
	"Same",
	"Relation",
	"Insured",
	"InsuredID",
	"HealthSurvey",
	"OrgName",
	"OrgCategory",
	"HaveInsuredList",
	"InsuredGroupByDay",
	"InsuredType",
	"InsuredList",
	"CommenceDate",
	"ExpiryDate",
	"Indate",
	"Sign",
	"Jurisdiction",
	"DisputeHandling",
	"PrevPolicyNo",
	"InsureBase",
	"BlanketInsureCode",
	"CustomType",
	"TrainProjects",
	"BusinessLocations",
	"OpenPoolNum",
	"HeatedPoolNum",
	"TrainingPoolNum",
	"PoolNum",
	"DinnerNum",
	"HaveDinnerNum",
	"CanteenNum",
	"ShopNum",
	"HaveRides",
	"HaveExplosive",
	"InnerArea",
	"OuterArea",
	"PoolName",
	"ArbitralAgency",
	"TrafficNum",
	"TemperatureType",
	"IsIndoor",
	"Extra",
	"BankAccount",
	"PayContact",
	"SuddenDeathTerms",
	"SpecAgreement",
	"HaveNegotiatedPrice",
	"HaveRenewalReminder",
	"LockStatus",
	"InsuranceCompany",
	"InsuranceCompanyAccount",
	"ChargeMode",
	"CanRevokeOrder",
	"CanPublicTransfers",
	"IsReminder",
	"GroundNum",
	"RemindersNum",
	"ReminderTimes",
	"RefusedReason",
	"UnpaidReason",
	"AdminReceived",
	"UserReceived",
	"HaveSuddenDeath",
	"HaveConfirmDate",
	"IsInvoice",
	"InvVisible",
	"InvBorrow",
	"InvTitle",
	"Traits",
	"Files",
	"InvStatus",
	"OrderStatus",
	"UpdStatus",
	"Creator",
	"CreateTime",
	"Regenerator",
	"UpdateTime",
	"Updator",
	"DomainID",
	"Addi",
	"Remark",
	"Status",
}

//Fields return all fields of struct.
func (r *TOrder) Fields() []string {
	return TOrderFields
}

//GetTableName return the associated db table name.
func (r *TOrder) GetTableName() string {
	var viewNamePattern = regexp.MustCompile(`(?i)^t_v_[a-z0-9_]+$`)
	tableName := "t_order"
	if viewNamePattern.MatchString(tableName) {
		return tableName[2:]
	}
	return tableName
}

// Create inserts the TOrder to the database.
func (r *TOrder) Create(db Queryer) error {
	err := db.QueryRow(
		`INSERT INTO t_order (trade_no, pay_order_no, transaction_id, batch, pay_time, pay_type, pay_channel, pay_name, pay_account_info, refundable, unit_price, refund_desc, amount, actual_amount, balance, balance_list, insure_order_no, refund_no, refund, refund_time, confirm_refund, agency_id, org_id, org_manager_id, insurance_type, insurance_type_id, insurance_police_id, plan_id, plan_name, insurer, policy_scheme, policy_doc, activity_name, activity_category, activity_desc, activity_location, activity_date_set, copies_num, insured_count, compulsory_student_num, non_compulsory_student_num, contact, fee_scheme, car_service_target, policyholder, policyholder_type, policyholder_id, same, relation, insured, insured_id, health_survey, org_name, org_category, have_insured_list, insured_group_by_day, insured_type, insured_list, commence_date, expiry_date, indate, sign, jurisdiction, dispute_handling, prev_policy_no, insure_base, blanket_insure_code, custom_type, train_projects, business_locations, open_pool_num, heated_pool_num, training_pool_num, pool_num, dinner_num, have_dinner_num, canteen_num, shop_num, have_rides, have_explosive, inner_area, outer_area, pool_name, arbitral_agency, traffic_num, temperature_type, is_indoor, extra, bank_account, pay_contact, sudden_death_terms, spec_agreement, have_negotiated_price, have_renewal_reminder, lock_status, insurance_company, insurance_company_account, charge_mode, can_revoke_order, can_public_transfers, is_reminder, ground_num, reminders_num, reminder_times, refused_reason, unpaid_reason, admin_received, user_received, have_sudden_death, have_confirm_date, is_invoice, inv_visible, inv_borrow, inv_title, traits, files, inv_status, order_status, upd_status, creator, create_time, regenerator, update_time, updator, domain_id, addi, remark, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36, $37, $38, $39, $40, $41, $42, $43, $44, $45, $46, $47, $48, $49, $50, $51, $52, $53, $54, $55, $56, $57, $58, $59, $60, $61, $62, $63, $64, $65, $66, $67, $68, $69, $70, $71, $72, $73, $74, $75, $76, $77, $78, $79, $80, $81, $82, $83, $84, $85, $86, $87, $88, $89, $90, $91, $92, $93, $94, $95, $96, $97, $98, $99, $100, $101, $102, $103, $104, $105, $106, $107, $108, $109, $110, $111, $112, $113, $114, $115, $116, $117, $118, $119, $120, $121, $122, $123, $124, $125, $126, $127, $128) RETURNING id`,
		&r.TradeNo, &r.PayOrderNo, &r.TransactionID, &r.Batch, &r.PayTime, &r.PayType, &r.PayChannel, &r.PayName, &r.PayAccountInfo, &r.Refundable, &r.UnitPrice, &r.RefundDesc, &r.Amount, &r.ActualAmount, &r.Balance, &r.BalanceList, &r.InsureOrderNo, &r.RefundNo, &r.Refund, &r.RefundTime, &r.ConfirmRefund, &r.AgencyID, &r.OrgID, &r.OrgManagerID, &r.InsuranceType, &r.InsuranceTypeID, &r.InsurancePoliceID, &r.PlanID, &r.PlanName, &r.Insurer, &r.PolicyScheme, &r.PolicyDoc, &r.ActivityName, &r.ActivityCategory, &r.ActivityDesc, &r.ActivityLocation, &r.ActivityDateSet, &r.CopiesNum, &r.InsuredCount, &r.CompulsoryStudentNum, &r.NonCompulsoryStudentNum, &r.Contact, &r.FeeScheme, &r.CarServiceTarget, &r.Policyholder, &r.PolicyholderType, &r.PolicyholderID, &r.Same, &r.Relation, &r.Insured, &r.InsuredID, &r.HealthSurvey, &r.OrgName, &r.OrgCategory, &r.HaveInsuredList, &r.InsuredGroupByDay, &r.InsuredType, &r.InsuredList, &r.CommenceDate, &r.ExpiryDate, &r.Indate, &r.Sign, &r.Jurisdiction, &r.DisputeHandling, &r.PrevPolicyNo, &r.InsureBase, &r.BlanketInsureCode, &r.CustomType, &r.TrainProjects, &r.BusinessLocations, &r.OpenPoolNum, &r.HeatedPoolNum, &r.TrainingPoolNum, &r.PoolNum, &r.DinnerNum, &r.HaveDinnerNum, &r.CanteenNum, &r.ShopNum, &r.HaveRides, &r.HaveExplosive, &r.InnerArea, &r.OuterArea, &r.PoolName, &r.ArbitralAgency, &r.TrafficNum, &r.TemperatureType, &r.IsIndoor, &r.Extra, &r.BankAccount, &r.PayContact, &r.SuddenDeathTerms, &r.SpecAgreement, &r.HaveNegotiatedPrice, &r.HaveRenewalReminder, &r.LockStatus, &r.InsuranceCompany, &r.InsuranceCompanyAccount, &r.ChargeMode, &r.CanRevokeOrder, &r.CanPublicTransfers, &r.IsReminder, &r.GroundNum, &r.RemindersNum, &r.ReminderTimes, &r.RefusedReason, &r.UnpaidReason, &r.AdminReceived, &r.UserReceived, &r.HaveSuddenDeath, &r.HaveConfirmDate, &r.IsInvoice, &r.InvVisible, &r.InvBorrow, &r.InvTitle, &r.Traits, &r.Files, &r.InvStatus, &r.OrderStatus, &r.UpdStatus, &r.Creator, &r.CreateTime, &r.Regenerator, &r.UpdateTime, &r.Updator, &r.DomainID, &r.Addi, &r.Remark, &r.Status).Scan(&r.ID)
	if err != nil {
		return errors.Wrap(err, "failed to insert t_order")
	}
	return nil
}

// GetTOrderByPk select the TOrder from the database.
func GetTOrderByPk(db Queryer, pk0 null.Int) (*TOrder, error) {
	var r TOrder
	err := db.QueryRow(
		`SELECT id, trade_no, pay_order_no, transaction_id, batch, pay_time, pay_type, pay_channel, pay_name, pay_account_info, refundable, unit_price, refund_desc, amount, actual_amount, balance, balance_list, insure_order_no, refund_no, refund, refund_time, confirm_refund, agency_id, org_id, org_manager_id, insurance_type, insurance_type_id, insurance_police_id, plan_id, plan_name, insurer, policy_scheme, policy_doc, activity_name, activity_category, activity_desc, activity_location, activity_date_set, copies_num, insured_count, compulsory_student_num, non_compulsory_student_num, contact, fee_scheme, car_service_target, policyholder, policyholder_type, policyholder_id, same, relation, insured, insured_id, health_survey, org_name, org_category, have_insured_list, insured_group_by_day, insured_type, insured_list, commence_date, expiry_date, indate, sign, jurisdiction, dispute_handling, prev_policy_no, insure_base, blanket_insure_code, custom_type, train_projects, business_locations, open_pool_num, heated_pool_num, training_pool_num, pool_num, dinner_num, have_dinner_num, canteen_num, shop_num, have_rides, have_explosive, inner_area, outer_area, pool_name, arbitral_agency, traffic_num, temperature_type, is_indoor, extra, bank_account, pay_contact, sudden_death_terms, spec_agreement, have_negotiated_price, have_renewal_reminder, lock_status, insurance_company, insurance_company_account, charge_mode, can_revoke_order, can_public_transfers, is_reminder, ground_num, reminders_num, reminder_times, refused_reason, unpaid_reason, admin_received, user_received, have_sudden_death, have_confirm_date, is_invoice, inv_visible, inv_borrow, inv_title, traits, files, inv_status, order_status, upd_status, creator, create_time, regenerator, update_time, updator, domain_id, addi, remark, status FROM t_order WHERE id = $1`,
		pk0).Scan(&r.ID, &r.TradeNo, &r.PayOrderNo, &r.TransactionID, &r.Batch, &r.PayTime, &r.PayType, &r.PayChannel, &r.PayName, &r.PayAccountInfo, &r.Refundable, &r.UnitPrice, &r.RefundDesc, &r.Amount, &r.ActualAmount, &r.Balance, &r.BalanceList, &r.InsureOrderNo, &r.RefundNo, &r.Refund, &r.RefundTime, &r.ConfirmRefund, &r.AgencyID, &r.OrgID, &r.OrgManagerID, &r.InsuranceType, &r.InsuranceTypeID, &r.InsurancePoliceID, &r.PlanID, &r.PlanName, &r.Insurer, &r.PolicyScheme, &r.PolicyDoc, &r.ActivityName, &r.ActivityCategory, &r.ActivityDesc, &r.ActivityLocation, &r.ActivityDateSet, &r.CopiesNum, &r.InsuredCount, &r.CompulsoryStudentNum, &r.NonCompulsoryStudentNum, &r.Contact, &r.FeeScheme, &r.CarServiceTarget, &r.Policyholder, &r.PolicyholderType, &r.PolicyholderID, &r.Same, &r.Relation, &r.Insured, &r.InsuredID, &r.HealthSurvey, &r.OrgName, &r.OrgCategory, &r.HaveInsuredList, &r.InsuredGroupByDay, &r.InsuredType, &r.InsuredList, &r.CommenceDate, &r.ExpiryDate, &r.Indate, &r.Sign, &r.Jurisdiction, &r.DisputeHandling, &r.PrevPolicyNo, &r.InsureBase, &r.BlanketInsureCode, &r.CustomType, &r.TrainProjects, &r.BusinessLocations, &r.OpenPoolNum, &r.HeatedPoolNum, &r.TrainingPoolNum, &r.PoolNum, &r.DinnerNum, &r.HaveDinnerNum, &r.CanteenNum, &r.ShopNum, &r.HaveRides, &r.HaveExplosive, &r.InnerArea, &r.OuterArea, &r.PoolName, &r.ArbitralAgency, &r.TrafficNum, &r.TemperatureType, &r.IsIndoor, &r.Extra, &r.BankAccount, &r.PayContact, &r.SuddenDeathTerms, &r.SpecAgreement, &r.HaveNegotiatedPrice, &r.HaveRenewalReminder, &r.LockStatus, &r.InsuranceCompany, &r.InsuranceCompanyAccount, &r.ChargeMode, &r.CanRevokeOrder, &r.CanPublicTransfers, &r.IsReminder, &r.GroundNum, &r.RemindersNum, &r.ReminderTimes, &r.RefusedReason, &r.UnpaidReason, &r.AdminReceived, &r.UserReceived, &r.HaveSuddenDeath, &r.HaveConfirmDate, &r.IsInvoice, &r.InvVisible, &r.InvBorrow, &r.InvTitle, &r.Traits, &r.Files, &r.InvStatus, &r.OrderStatus, &r.UpdStatus, &r.Creator, &r.CreateTime, &r.Regenerator, &r.UpdateTime, &r.Updator, &r.DomainID, &r.Addi, &r.Remark, &r.Status)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select t_order")
	}
	return &r, nil
}

/*TParam 提供用户设置参数
belongTo value scope
系统一级:1000-1990
系统二级:2000-2990
系统三级:3000-3990
系统四级:4000-4990

应用系统一级:11000-11990
应用系统二级:12000-12990
应用系统三级:13000-13990
应用系统四级:14000-14990
预置参数ID只使用偶数 represents kuser.t_param */
type TParam struct {
	ID         null.Int       `json:"ID,omitempty" db:"id,true,integer"`                         /* id 参数编号 */
	Belongto   null.Int       `json:"Belongto,omitempty" db:"belongto,false,bigint"`             /* belongto 类属 */
	Name       string         `json:"Name,omitempty" db:"name,false,character varying"`          /* name 参数名称 */
	Value      null.String    `json:"Value,omitempty" db:"value,false,character varying"`        /* value 参数值 */
	DataType   null.String    `json:"DataType,omitempty" db:"data_type,false,character varying"` /* data_type 数据类型, string,number,bool,nil */
	CreateTime null.Int       `json:"CreateTime,omitempty" db:"create_time,false,bigint"`        /* create_time 生成时间 */
	Creator    null.Int       `json:"Creator,omitempty" db:"creator,false,bigint"`               /* creator 本数据创建者 */
	DomainID   null.Int       `json:"DomainID,omitempty" db:"domain_id,false,bigint"`            /* domain_id 数据隶属 */
	Addi       types.JSONText `json:"Addi,omitempty" db:"addi,false,jsonb"`                      /* addi 附加信息 */
	Remark     null.String    `json:"Remark,omitempty" db:"remark,false,character varying"`      /* remark 备注 */
	Status     null.String    `json:"Status,omitempty" db:"status,false,character varying"`      /* status 状态，00：草稿，01：有效，02：作废 */
	Filter                    // build DML where clause
}

//TParamFields full field list for default query
var TParamFields = []string{
	"ID",
	"Belongto",
	"Name",
	"Value",
	"DataType",
	"CreateTime",
	"Creator",
	"DomainID",
	"Addi",
	"Remark",
	"Status",
}

//Fields return all fields of struct.
func (r *TParam) Fields() []string {
	return TParamFields
}

//GetTableName return the associated db table name.
func (r *TParam) GetTableName() string {
	var viewNamePattern = regexp.MustCompile(`(?i)^t_v_[a-z0-9_]+$`)
	tableName := "t_param"
	if viewNamePattern.MatchString(tableName) {
		return tableName[2:]
	}
	return tableName
}

// Create inserts the TParam to the database.
func (r *TParam) Create(db Queryer) error {
	err := db.QueryRow(
		`INSERT INTO t_param (belongto, name, value, data_type, create_time, creator, domain_id, addi, remark, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`,
		&r.Belongto, &r.Name, &r.Value, &r.DataType, &r.CreateTime, &r.Creator, &r.DomainID, &r.Addi, &r.Remark, &r.Status).Scan(&r.ID)
	if err != nil {
		return errors.Wrap(err, "failed to insert t_param")
	}
	return nil
}

// GetTParamByPk select the TParam from the database.
func GetTParamByPk(db Queryer, pk0 null.Int) (*TParam, error) {
	var r TParam
	err := db.QueryRow(
		`SELECT id, belongto, name, value, data_type, create_time, creator, domain_id, addi, remark, status FROM t_param WHERE id = $1`,
		pk0).Scan(&r.ID, &r.Belongto, &r.Name, &r.Value, &r.DataType, &r.CreateTime, &r.Creator, &r.DomainID, &r.Addi, &r.Remark, &r.Status)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select t_param")
	}
	return &r, nil
}

/*TPayAccount 支付账号信息 represents kuser.t_pay_account */
type TPayAccount struct {
	ID          null.Int       `json:"ID,omitempty" db:"id,true,integer"`                      /* id 编号 */
	Type        null.String    `json:"Type,omitempty" db:"type,false,character varying"`       /* type wx_mp: 微信公众号, wx_open: 微信开放平台, ali: 阿里, ui: 联保 */
	Name        null.String    `json:"Name,omitempty" db:"name,false,character varying"`       /* name 名称：校快保，泰合，联保，近邻，人保，太平洋保险，人寿 */
	AppID       null.String    `json:"AppID,omitempty" db:"app_id,false,character varying"`    /* app_id 关联应用ID: 微信公众号 */
	Account     null.String    `json:"Account,omitempty" db:"account,false,character varying"` /* account 账号，微信支付商户号，支付宝商户号 */
	Key         null.String    `json:"Key,omitempty" db:"key,false,character varying"`         /* key 密钥 */
	Cert        null.String    `json:"Cert,omitempty" db:"cert,false,character varying"`       /* cert 证书 */
	Refundable  null.Bool      `json:"Refundable,omitempty" db:"refundable,false,boolean"`     /* refundable 是否支持退款 */
	Regenerator null.Int       `json:"Regenerator,omitempty" db:"regenerator,false,bigint"`    /* regenerator 更新者 */
	UpdateTime  null.Int       `json:"UpdateTime,omitempty" db:"update_time,false,bigint"`     /* update_time 帐号信息更新时间 */
	Creator     null.Int       `json:"Creator,omitempty" db:"creator,false,bigint"`            /* creator 本数据创建者 */
	CreateTime  null.Int       `json:"CreateTime,omitempty" db:"create_time,false,bigint"`     /* create_time 生成时间 */
	DomainID    null.Int       `json:"DomainID,omitempty" db:"domain_id,false,bigint"`         /* domain_id 数据隶属 */
	Addi        types.JSONText `json:"Addi,omitempty" db:"addi,false,jsonb"`                   /* addi 附加信息 */
	Remark      null.String    `json:"Remark,omitempty" db:"remark,false,character varying"`   /* remark 备注 */
	Status      null.String    `json:"Status,omitempty" db:"status,false,character varying"`   /* status 状态，00：草稿，02：有效，04: 停用，06：作废 */
	Filter                     // build DML where clause
}

//TPayAccountFields full field list for default query
var TPayAccountFields = []string{
	"ID",
	"Type",
	"Name",
	"AppID",
	"Account",
	"Key",
	"Cert",
	"Refundable",
	"Regenerator",
	"UpdateTime",
	"Creator",
	"CreateTime",
	"DomainID",
	"Addi",
	"Remark",
	"Status",
}

//Fields return all fields of struct.
func (r *TPayAccount) Fields() []string {
	return TPayAccountFields
}

//GetTableName return the associated db table name.
func (r *TPayAccount) GetTableName() string {
	var viewNamePattern = regexp.MustCompile(`(?i)^t_v_[a-z0-9_]+$`)
	tableName := "t_pay_account"
	if viewNamePattern.MatchString(tableName) {
		return tableName[2:]
	}
	return tableName
}

// Create inserts the TPayAccount to the database.
func (r *TPayAccount) Create(db Queryer) error {
	err := db.QueryRow(
		`INSERT INTO t_pay_account (type, name, app_id, account, key, cert, refundable, regenerator, update_time, creator, create_time, domain_id, addi, remark, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15) RETURNING id`,
		&r.Type, &r.Name, &r.AppID, &r.Account, &r.Key, &r.Cert, &r.Refundable, &r.Regenerator, &r.UpdateTime, &r.Creator, &r.CreateTime, &r.DomainID, &r.Addi, &r.Remark, &r.Status).Scan(&r.ID)
	if err != nil {
		return errors.Wrap(err, "failed to insert t_pay_account")
	}
	return nil
}

// GetTPayAccountByPk select the TPayAccount from the database.
func GetTPayAccountByPk(db Queryer, pk0 null.Int) (*TPayAccount, error) {
	var r TPayAccount
	err := db.QueryRow(
		`SELECT id, type, name, app_id, account, key, cert, refundable, regenerator, update_time, creator, create_time, domain_id, addi, remark, status FROM t_pay_account WHERE id = $1`,
		pk0).Scan(&r.ID, &r.Type, &r.Name, &r.AppID, &r.Account, &r.Key, &r.Cert, &r.Refundable, &r.Regenerator, &r.UpdateTime, &r.Creator, &r.CreateTime, &r.DomainID, &r.Addi, &r.Remark, &r.Status)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select t_pay_account")
	}
	return &r, nil
}

/*TPayment 缴费表(用于对公转账自动化) represents kuser.t_payment */
type TPayment struct {
	ID             null.Int       `json:"ID,omitempty" db:"id,true,integer"`                                    /* id 编码 */
	Batch          null.String    `json:"Batch,omitempty" db:"batch,false,character varying"`                   /* batch 批次 */
	PolicyNo       null.String    `json:"PolicyNo,omitempty" db:"policy_no,false,character varying"`            /* policy_no 保单号 */
	TransferNo     null.String    `json:"TransferNo,omitempty" db:"transfer_no,false,character varying"`        /* transfer_no 转账流水号 */
	TransferAmount null.Float     `json:"TransferAmount,omitempty" db:"transfer_amount,false,double precision"` /* transfer_amount 金额 */
	Creator        null.Int       `json:"Creator,omitempty" db:"creator,false,bigint"`                          /* creator 创建者用户ID */
	CreateTime     null.Int       `json:"CreateTime,omitempty" db:"create_time,false,bigint"`                   /* create_time 创建时间 */
	Regenerator    null.Int       `json:"Regenerator,omitempty" db:"regenerator,false,bigint"`                  /* regenerator 更新者 */
	UpdateTime     null.Int       `json:"UpdateTime,omitempty" db:"update_time,false,bigint"`                   /* update_time 修改时间 */
	DomainID       null.Int       `json:"DomainID,omitempty" db:"domain_id,false,bigint"`                       /* domain_id 数据属主 */
	Addi           types.JSONText `json:"Addi,omitempty" db:"addi,false,jsonb"`                                 /* addi 附加数据 */
	Remark         null.String    `json:"Remark,omitempty" db:"remark,false,character varying"`                 /* remark 备注 */
	Status         null.String    `json:"Status,omitempty" db:"status,false,character varying"`                 /* status 状态, 未缴费: 0, 已缴费: 2, 作废: 4 */
	Filter                        // build DML where clause
}

//TPaymentFields full field list for default query
var TPaymentFields = []string{
	"ID",
	"Batch",
	"PolicyNo",
	"TransferNo",
	"TransferAmount",
	"Creator",
	"CreateTime",
	"Regenerator",
	"UpdateTime",
	"DomainID",
	"Addi",
	"Remark",
	"Status",
}

//Fields return all fields of struct.
func (r *TPayment) Fields() []string {
	return TPaymentFields
}

//GetTableName return the associated db table name.
func (r *TPayment) GetTableName() string {
	var viewNamePattern = regexp.MustCompile(`(?i)^t_v_[a-z0-9_]+$`)
	tableName := "t_payment"
	if viewNamePattern.MatchString(tableName) {
		return tableName[2:]
	}
	return tableName
}

// Create inserts the TPayment to the database.
func (r *TPayment) Create(db Queryer) error {
	err := db.QueryRow(
		`INSERT INTO t_payment (batch, policy_no, transfer_no, transfer_amount, creator, create_time, regenerator, update_time, domain_id, addi, remark, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id`,
		&r.Batch, &r.PolicyNo, &r.TransferNo, &r.TransferAmount, &r.Creator, &r.CreateTime, &r.Regenerator, &r.UpdateTime, &r.DomainID, &r.Addi, &r.Remark, &r.Status).Scan(&r.ID)
	if err != nil {
		return errors.Wrap(err, "failed to insert t_payment")
	}
	return nil
}

// GetTPaymentByPk select the TPayment from the database.
func GetTPaymentByPk(db Queryer, pk0 null.Int) (*TPayment, error) {
	var r TPayment
	err := db.QueryRow(
		`SELECT id, batch, policy_no, transfer_no, transfer_amount, creator, create_time, regenerator, update_time, domain_id, addi, remark, status FROM t_payment WHERE id = $1`,
		pk0).Scan(&r.ID, &r.Batch, &r.PolicyNo, &r.TransferNo, &r.TransferAmount, &r.Creator, &r.CreateTime, &r.Regenerator, &r.UpdateTime, &r.DomainID, &r.Addi, &r.Remark, &r.Status)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select t_payment")
	}
	return &r, nil
}

/*TPrice 价格设置表
represents kuser.t_price */
type TPrice struct {
	ID                 null.Int       `json:"ID,omitempty" db:"id,true,integer"`                                 /* id 价格id */
	Title              null.String    `json:"Title,omitempty" db:"title,false,character varying"`                /* title 标题 */
	Category           null.String    `json:"Category,omitempty" db:"category,false,character varying"`          /* category 类型 */
	InsuranceTypeID    null.Int       `json:"InsuranceTypeID,omitempty" db:"insurance_type_id,false,bigint"`     /* insurance_type_id 险种id */
	OrgName            null.String    `json:"OrgName,omitempty" db:"org_name,false,character varying"`           /* org_name 投保单位 */
	Province           null.String    `json:"Province,omitempty" db:"province,false,character varying"`          /* province 省 */
	City               null.String    `json:"City,omitempty" db:"city,false,character varying"`                  /* city 市 */
	District           null.String    `json:"District,omitempty" db:"district,false,character varying"`          /* district 区/县 */
	GuaranteedProjects types.JSONText `json:"GuaranteedProjects,omitempty" db:"guaranteed_projects,false,jsonb"` /* guaranteed_projects 保障项目 */
	ExtraProjects      types.JSONText `json:"ExtraProjects,omitempty" db:"extra_projects,false,jsonb"`           /* extra_projects 附加条款 */
	PriceConfig        types.JSONText `json:"PriceConfig,omitempty" db:"price_config,false,jsonb"`               /* price_config 价格配置 */
	Files              types.JSONText `json:"Files,omitempty" db:"files,false,jsonb"`                            /* files 模板路径  */
	IsDefault          null.Bool      `json:"IsDefault,omitempty" db:"is_default,false,boolean"`                 /* is_default 是否默认值 */
	CreatorID          null.Int       `json:"CreatorID,omitempty" db:"creator_id,false,bigint"`                  /* creator_id 创建者id */
	CreateTime         null.Int       `json:"CreateTime,omitempty" db:"create_time,false,bigint"`                /* create_time 创建时间 */
	Regenerator        null.Int       `json:"Regenerator,omitempty" db:"regenerator,false,bigint"`               /* regenerator 更新者 */
	UpdateTime         null.Int       `json:"UpdateTime,omitempty" db:"update_time,false,bigint"`                /* update_time 修改时间 */
	Updator            null.Int       `json:"Updator,omitempty" db:"updator,false,bigint"`                       /* updator 更新者 */
	DomainID           null.Int       `json:"DomainID,omitempty" db:"domain_id,false,bigint"`                    /* domain_id 数据属主 */
	Addi               types.JSONText `json:"Addi,omitempty" db:"addi,false,jsonb"`                              /* addi 备用字段 */
	Remark             null.String    `json:"Remark,omitempty" db:"remark,false,character varying"`              /* remark 备注 */
	Status             null.String    `json:"Status,omitempty" db:"status,false,character varying"`              /* status 状态，0：有效，2：无效 */
	Filter                            // build DML where clause
}

//TPriceFields full field list for default query
var TPriceFields = []string{
	"ID",
	"Title",
	"Category",
	"InsuranceTypeID",
	"OrgName",
	"Province",
	"City",
	"District",
	"GuaranteedProjects",
	"ExtraProjects",
	"PriceConfig",
	"Files",
	"IsDefault",
	"CreatorID",
	"CreateTime",
	"Regenerator",
	"UpdateTime",
	"Updator",
	"DomainID",
	"Addi",
	"Remark",
	"Status",
}

//Fields return all fields of struct.
func (r *TPrice) Fields() []string {
	return TPriceFields
}

//GetTableName return the associated db table name.
func (r *TPrice) GetTableName() string {
	var viewNamePattern = regexp.MustCompile(`(?i)^t_v_[a-z0-9_]+$`)
	tableName := "t_price"
	if viewNamePattern.MatchString(tableName) {
		return tableName[2:]
	}
	return tableName
}

// Create inserts the TPrice to the database.
func (r *TPrice) Create(db Queryer) error {
	err := db.QueryRow(
		`INSERT INTO t_price (title, category, insurance_type_id, org_name, province, city, district, guaranteed_projects, extra_projects, price_config, files, is_default, creator_id, create_time, regenerator, update_time, updator, domain_id, addi, remark, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21) RETURNING id`,
		&r.Title, &r.Category, &r.InsuranceTypeID, &r.OrgName, &r.Province, &r.City, &r.District, &r.GuaranteedProjects, &r.ExtraProjects, &r.PriceConfig, &r.Files, &r.IsDefault, &r.CreatorID, &r.CreateTime, &r.Regenerator, &r.UpdateTime, &r.Updator, &r.DomainID, &r.Addi, &r.Remark, &r.Status).Scan(&r.ID)
	if err != nil {
		return errors.Wrap(err, "failed to insert t_price")
	}
	return nil
}

// GetTPriceByPk select the TPrice from the database.
func GetTPriceByPk(db Queryer, pk0 null.Int) (*TPrice, error) {
	var r TPrice
	err := db.QueryRow(
		`SELECT id, title, category, insurance_type_id, org_name, province, city, district, guaranteed_projects, extra_projects, price_config, files, is_default, creator_id, create_time, regenerator, update_time, updator, domain_id, addi, remark, status FROM t_price WHERE id = $1`,
		pk0).Scan(&r.ID, &r.Title, &r.Category, &r.InsuranceTypeID, &r.OrgName, &r.Province, &r.City, &r.District, &r.GuaranteedProjects, &r.ExtraProjects, &r.PriceConfig, &r.Files, &r.IsDefault, &r.CreatorID, &r.CreateTime, &r.Regenerator, &r.UpdateTime, &r.Updator, &r.DomainID, &r.Addi, &r.Remark, &r.Status)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select t_price")
	}
	return &r, nil
}

/*TPrj 项目信息表 represents kuser.t_prj */
type TPrj struct {
	ID         null.Int    `json:"ID,omitempty" db:"id,true,integer"`                    /* id 编号 */
	Name       null.String `json:"Name,omitempty" db:"name,false,character varying"`     /* name 项目名称 */
	Limn       null.String `json:"Limn,omitempty" db:"limn,false,character varying"`     /* limn 项目描述 */
	Price      null.Float  `json:"Price,omitempty" db:"price,false,numeric"`             /* price 报价 */
	Cycle      null.Int    `json:"Cycle,omitempty" db:"cycle,false,integer"`             /* cycle 期望周期，以自然日为单位 */
	Issuer     null.Int    `json:"Issuer,omitempty" db:"issuer,false,bigint"`            /* issuer 发布者编号，四方 */
	CreateTime null.Int    `json:"CreateTime,omitempty" db:"create_time,false,bigint"`   /* create_time 创建时间 */
	IssueTime  null.Int    `json:"IssueTime,omitempty" db:"issue_time,false,bigint"`     /* issue_time 发布时间 */
	Deadline   null.Int    `json:"Deadline,omitempty" db:"deadline,false,bigint"`        /* deadline 截止时间 */
	Remark     null.String `json:"Remark,omitempty" db:"remark,false,character varying"` /* remark 备注 */
	Status     null.String `json:"Status,omitempty" db:"status,false,character varying"` /* status draft,未发布
	isuued,已发布
	cancelled,取消
	engaged,确定了承接人
	signed,签订合同
	finished,项目完成 */
	Filter // build DML where clause
}

//TPrjFields full field list for default query
var TPrjFields = []string{
	"ID",
	"Name",
	"Limn",
	"Price",
	"Cycle",
	"Issuer",
	"CreateTime",
	"IssueTime",
	"Deadline",
	"Remark",
	"Status",
}

//Fields return all fields of struct.
func (r *TPrj) Fields() []string {
	return TPrjFields
}

//GetTableName return the associated db table name.
func (r *TPrj) GetTableName() string {
	var viewNamePattern = regexp.MustCompile(`(?i)^t_v_[a-z0-9_]+$`)
	tableName := "t_prj"
	if viewNamePattern.MatchString(tableName) {
		return tableName[2:]
	}
	return tableName
}

// Create inserts the TPrj to the database.
func (r *TPrj) Create(db Queryer) error {
	err := db.QueryRow(
		`INSERT INTO t_prj (name, limn, price, cycle, issuer, create_time, issue_time, deadline, remark, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`,
		&r.Name, &r.Limn, &r.Price, &r.Cycle, &r.Issuer, &r.CreateTime, &r.IssueTime, &r.Deadline, &r.Remark, &r.Status).Scan(&r.ID)
	if err != nil {
		return errors.Wrap(err, "failed to insert t_prj")
	}
	return nil
}

// GetTPrjByPk select the TPrj from the database.
func GetTPrjByPk(db Queryer, pk0 null.Int) (*TPrj, error) {
	var r TPrj
	err := db.QueryRow(
		`SELECT id, name, limn, price, cycle, issuer, create_time, issue_time, deadline, remark, status FROM t_prj WHERE id = $1`,
		pk0).Scan(&r.ID, &r.Name, &r.Limn, &r.Price, &r.Cycle, &r.Issuer, &r.CreateTime, &r.IssueTime, &r.Deadline, &r.Remark, &r.Status)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select t_prj")
	}
	return &r, nil
}

/*TProof 人才知识能力领域说明表 represents kuser.t_proof */
type TProof struct {
	ID          null.Int    `json:"ID,omitempty" db:"id,true,integer"`                    /* id 编号 */
	UserID      null.Int    `json:"UserID,omitempty" db:"user_id,false,bigint"`           /* user_id 用户编号 */
	ExpertiseID null.Int    `json:"ExpertiseID,omitempty" db:"expertise_id,false,bigint"` /* expertise_id 知识能力领域编号 */
	Limn        null.String `json:"Limn,omitempty" db:"limn,false,character varying"`     /* limn 能力描述 */
	CreateTime  null.Int    `json:"CreateTime,omitempty" db:"create_time,false,bigint"`   /* create_time 创建时间 */
	UpdateTime  null.Int    `json:"UpdateTime,omitempty" db:"update_time,false,bigint"`   /* update_time 更新时间 */
	Filter                  // build DML where clause
}

//TProofFields full field list for default query
var TProofFields = []string{
	"ID",
	"UserID",
	"ExpertiseID",
	"Limn",
	"CreateTime",
	"UpdateTime",
}

//Fields return all fields of struct.
func (r *TProof) Fields() []string {
	return TProofFields
}

//GetTableName return the associated db table name.
func (r *TProof) GetTableName() string {
	var viewNamePattern = regexp.MustCompile(`(?i)^t_v_[a-z0-9_]+$`)
	tableName := "t_proof"
	if viewNamePattern.MatchString(tableName) {
		return tableName[2:]
	}
	return tableName
}

// Create inserts the TProof to the database.
func (r *TProof) Create(db Queryer) error {
	err := db.QueryRow(
		`INSERT INTO t_proof (user_id, expertise_id, limn, create_time, update_time) VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		&r.UserID, &r.ExpertiseID, &r.Limn, &r.CreateTime, &r.UpdateTime).Scan(&r.ID)
	if err != nil {
		return errors.Wrap(err, "failed to insert t_proof")
	}
	return nil
}

// GetTProofByPk select the TProof from the database.
func GetTProofByPk(db Queryer, pk0 null.Int) (*TProof, error) {
	var r TProof
	err := db.QueryRow(
		`SELECT id, user_id, expertise_id, limn, create_time, update_time FROM t_proof WHERE id = $1`,
		pk0).Scan(&r.ID, &r.UserID, &r.ExpertiseID, &r.Limn, &r.CreateTime, &r.UpdateTime)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select t_proof")
	}
	return &r, nil
}

/*TProve 知识能力领域鉴定、证明表 represents kuser.t_prove */
type TProve struct {
	ID         null.Int    `json:"ID,omitempty" db:"id,true,integer"`                          /* id 编号 */
	ProofID    null.Int    `json:"ProofID,omitempty" db:"proof_id,false,bigint"`               /* proof_id 被鉴定材料编号 */
	Judgement  null.String `json:"Judgement,omitempty" db:"judgement,false,character varying"` /* judgement 鉴定结论 */
	Creator    null.Int    `json:"Creator,omitempty" db:"creator,false,bigint"`                /* creator 鉴定者 */
	CreateTime null.Int    `json:"CreateTime,omitempty" db:"create_time,false,bigint"`         /* create_time 鉴定时间 */
	UpdateTime null.Int    `json:"UpdateTime,omitempty" db:"update_time,false,bigint"`         /* update_time 鉴定更新时间 */
	Filter                 // build DML where clause
}

//TProveFields full field list for default query
var TProveFields = []string{
	"ID",
	"ProofID",
	"Judgement",
	"Creator",
	"CreateTime",
	"UpdateTime",
}

//Fields return all fields of struct.
func (r *TProve) Fields() []string {
	return TProveFields
}

//GetTableName return the associated db table name.
func (r *TProve) GetTableName() string {
	var viewNamePattern = regexp.MustCompile(`(?i)^t_v_[a-z0-9_]+$`)
	tableName := "t_prove"
	if viewNamePattern.MatchString(tableName) {
		return tableName[2:]
	}
	return tableName
}

// Create inserts the TProve to the database.
func (r *TProve) Create(db Queryer) error {
	err := db.QueryRow(
		`INSERT INTO t_prove (proof_id, judgement, creator, create_time, update_time) VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		&r.ProofID, &r.Judgement, &r.Creator, &r.CreateTime, &r.UpdateTime).Scan(&r.ID)
	if err != nil {
		return errors.Wrap(err, "failed to insert t_prove")
	}
	return nil
}

// GetTProveByPk select the TProve from the database.
func GetTProveByPk(db Queryer, pk0 null.Int) (*TProve, error) {
	var r TProve
	err := db.QueryRow(
		`SELECT id, proof_id, judgement, creator, create_time, update_time FROM t_prove WHERE id = $1`,
		pk0).Scan(&r.ID, &r.ProofID, &r.Judgement, &r.Creator, &r.CreateTime, &r.UpdateTime)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select t_prove")
	}
	return &r, nil
}

/*TQualification 人才资质表 represents kuser.t_qualification */
type TQualification struct {
	ID          null.Int `json:"ID,omitempty" db:"id,true,integer"`                    /* id 资质证明编号 */
	UserID      null.Int `json:"UserID,omitempty" db:"user_id,false,bigint"`           /* user_id 用户编号 */
	ExpertiseID null.Int `json:"ExpertiseID,omitempty" db:"expertise_id,false,bigint"` /* expertise_id 专长编号 */
	CreateTime  null.Int `json:"CreateTime,omitempty" db:"create_time,false,bigint"`   /* create_time 创建时间 */
	Filter               // build DML where clause
}

//TQualificationFields full field list for default query
var TQualificationFields = []string{
	"ID",
	"UserID",
	"ExpertiseID",
	"CreateTime",
}

//Fields return all fields of struct.
func (r *TQualification) Fields() []string {
	return TQualificationFields
}

//GetTableName return the associated db table name.
func (r *TQualification) GetTableName() string {
	var viewNamePattern = regexp.MustCompile(`(?i)^t_v_[a-z0-9_]+$`)
	tableName := "t_qualification"
	if viewNamePattern.MatchString(tableName) {
		return tableName[2:]
	}
	return tableName
}

// Create inserts the TQualification to the database.
func (r *TQualification) Create(db Queryer) error {
	err := db.QueryRow(
		`INSERT INTO t_qualification (user_id, expertise_id, create_time) VALUES ($1, $2, $3) RETURNING id`,
		&r.UserID, &r.ExpertiseID, &r.CreateTime).Scan(&r.ID)
	if err != nil {
		return errors.Wrap(err, "failed to insert t_qualification")
	}
	return nil
}

// GetTQualificationByPk select the TQualification from the database.
func GetTQualificationByPk(db Queryer, pk0 null.Int) (*TQualification, error) {
	var r TQualification
	err := db.QueryRow(
		`SELECT id, user_id, expertise_id, create_time FROM t_qualification WHERE id = $1`,
		pk0).Scan(&r.ID, &r.UserID, &r.ExpertiseID, &r.CreateTime)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select t_qualification")
	}
	return &r, nil
}

/*TRegion 区域列表 represents kuser.t_region */
type TRegion struct {
	ID              null.Int       `json:"ID,omitempty" db:"id,true,integer"`                                        /* id 区域编号 */
	RegionName      null.String    `json:"RegionName,omitempty" db:"region_name,false,character varying"`            /* region_name 区域名称 */
	Code            null.Int       `json:"Code,omitempty" db:"code,false,bigint"`                                    /* code 区域行政编码 */
	RegionShortName null.String    `json:"RegionShortName,omitempty" db:"region_short_name,false,character varying"` /* region_short_name 区域缩写 */
	ParentID        null.Int       `json:"ParentID,omitempty" db:"parent_id,false,bigint"`                           /* parent_id 区域的父级id, 定义: 省级没有父级, parent_id为0; 市级的父级是省; 区县的父级是市 */
	Level           null.Int       `json:"Level,omitempty" db:"level,false,bigint"`                                  /* level 地区级别: 2-省,4-市,6-区/县 */
	UpdateTime      null.Int       `json:"UpdateTime,omitempty" db:"update_time,false,bigint"`                       /* update_time 可能以后存在着一年更新一次区域表 */
	Creator         null.Int       `json:"Creator,omitempty" db:"creator,false,bigint"`                              /* creator 本数据创建者 */
	DomainID        null.Int       `json:"DomainID,omitempty" db:"domain_id,false,bigint"`                           /* domain_id 数据隶属 */
	Addi            types.JSONText `json:"Addi,omitempty" db:"addi,false,jsonb"`                                     /* addi 附加信息 */
	Remark          null.String    `json:"Remark,omitempty" db:"remark,false,character varying"`                     /* remark 备注 */
	Status          null.String    `json:"Status,omitempty" db:"status,false,character varying"`                     /* status 0:有效, 2: 删除, 过了一段时间有些区域可能会被删除 */
	Filter                         // build DML where clause
}

//TRegionFields full field list for default query
var TRegionFields = []string{
	"ID",
	"RegionName",
	"Code",
	"RegionShortName",
	"ParentID",
	"Level",
	"UpdateTime",
	"Creator",
	"DomainID",
	"Addi",
	"Remark",
	"Status",
}

//Fields return all fields of struct.
func (r *TRegion) Fields() []string {
	return TRegionFields
}

//GetTableName return the associated db table name.
func (r *TRegion) GetTableName() string {
	var viewNamePattern = regexp.MustCompile(`(?i)^t_v_[a-z0-9_]+$`)
	tableName := "t_region"
	if viewNamePattern.MatchString(tableName) {
		return tableName[2:]
	}
	return tableName
}

// Create inserts the TRegion to the database.
func (r *TRegion) Create(db Queryer) error {
	err := db.QueryRow(
		`INSERT INTO t_region (region_name, code, region_short_name, parent_id, level, update_time, creator, domain_id, addi, remark, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id`,
		&r.RegionName, &r.Code, &r.RegionShortName, &r.ParentID, &r.Level, &r.UpdateTime, &r.Creator, &r.DomainID, &r.Addi, &r.Remark, &r.Status).Scan(&r.ID)
	if err != nil {
		return errors.Wrap(err, "failed to insert t_region")
	}
	return nil
}

// GetTRegionByPk select the TRegion from the database.
func GetTRegionByPk(db Queryer, pk0 null.Int) (*TRegion, error) {
	var r TRegion
	err := db.QueryRow(
		`SELECT id, region_name, code, region_short_name, parent_id, level, update_time, creator, domain_id, addi, remark, status FROM t_region WHERE id = $1`,
		pk0).Scan(&r.ID, &r.RegionName, &r.Code, &r.RegionShortName, &r.ParentID, &r.Level, &r.UpdateTime, &r.Creator, &r.DomainID, &r.Addi, &r.Remark, &r.Status)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select t_region")
	}
	return &r, nil
}

/*TRelation 描述两个实体间关系
校快保，描述销售/学校管理员/学校统计员与学校间的对应关系 represents kuser.t_relation */
type TRelation struct {
	ID             null.Int       `json:"ID,omitempty" db:"id,true,integer"`                                      /* id 编号 */
	LeftID         null.Int       `json:"LeftID,omitempty" db:"left_id,false,bigint"`                             /* left_id 左编号 */
	LeftType       null.String    `json:"LeftType,omitempty" db:"left_type,false,character varying"`              /* left_type 左类型，用户编号: t_user.id */
	LeftKey        null.String    `json:"LeftKey,omitempty" db:"left_key,false,character varying"`                /* left_key 左识别标识 */
	LeftKeyType    null.String    `json:"LeftKeyType,omitempty" db:"left_key_type,false,character varying"`       /* left_key_type 左识别标识类型，帐号: account, 邮箱: email, 手机: tel, 微信公众号openID: mp_open_id, 微信开放平台openID: wx_open_id */
	Kind           string         `json:"Kind,omitempty" db:"kind,false,character varying"`                       /* kind 关系类型 */
	RightID        null.Int       `json:"RightID,omitempty" db:"right_id,false,bigint"`                           /* right_id 目标资源编号，如学校ID */
	RightType      null.String    `json:"RightType,omitempty" db:"right_type,false,character varying"`            /* right_type 右类型，如t_school.id */
	RightValueType null.String    `json:"RightValueType,omitempty" db:"right_value_type,false,character varying"` /* right_value_type 右值数据类型，默认为int8,则值存储于right_id,其它类型则存储于right_value中 */
	RightValue     null.String    `json:"RightValue,omitempty" db:"right_value,false,character varying"`          /* right_value 右值(非int8类型) */
	RuleArea       null.String    `json:"RuleArea,omitempty" db:"rule_area,false,character varying"`              /* rule_area 管辖地区 */
	Creator        null.Int       `json:"Creator,omitempty" db:"creator,false,bigint"`                            /* creator 创建者用户ID */
	CreateTime     null.Int       `json:"CreateTime,omitempty" db:"create_time,false,bigint"`                     /* create_time 创建时间 */
	Regenerator    null.Int       `json:"Regenerator,omitempty" db:"regenerator,false,bigint"`                    /* regenerator 更新者 */
	UpdateTime     null.Int       `json:"UpdateTime,omitempty" db:"update_time,false,bigint"`                     /* update_time 修改时间 */
	DomainID       null.Int       `json:"DomainID,omitempty" db:"domain_id,false,bigint"`                         /* domain_id 数据属主 */
	Addi           types.JSONText `json:"Addi,omitempty" db:"addi,false,jsonb"`                                   /* addi 附加数据 */
	Remark         null.String    `json:"Remark,omitempty" db:"remark,false,character varying"`                   /* remark 备注 */
	Status         null.String    `json:"Status,omitempty" db:"status,false,character varying"`                   /* status 状态 */
	Filter                        // build DML where clause
}

//TRelationFields full field list for default query
var TRelationFields = []string{
	"ID",
	"LeftID",
	"LeftType",
	"LeftKey",
	"LeftKeyType",
	"Kind",
	"RightID",
	"RightType",
	"RightValueType",
	"RightValue",
	"RuleArea",
	"Creator",
	"CreateTime",
	"Regenerator",
	"UpdateTime",
	"DomainID",
	"Addi",
	"Remark",
	"Status",
}

//Fields return all fields of struct.
func (r *TRelation) Fields() []string {
	return TRelationFields
}

//GetTableName return the associated db table name.
func (r *TRelation) GetTableName() string {
	var viewNamePattern = regexp.MustCompile(`(?i)^t_v_[a-z0-9_]+$`)
	tableName := "t_relation"
	if viewNamePattern.MatchString(tableName) {
		return tableName[2:]
	}
	return tableName
}

// Create inserts the TRelation to the database.
func (r *TRelation) Create(db Queryer) error {
	err := db.QueryRow(
		`INSERT INTO t_relation (left_id, left_type, left_key, left_key_type, kind, right_id, right_type, right_value_type, right_value, rule_area, creator, create_time, regenerator, update_time, domain_id, addi, remark, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18) RETURNING id`,
		&r.LeftID, &r.LeftType, &r.LeftKey, &r.LeftKeyType, &r.Kind, &r.RightID, &r.RightType, &r.RightValueType, &r.RightValue, &r.RuleArea, &r.Creator, &r.CreateTime, &r.Regenerator, &r.UpdateTime, &r.DomainID, &r.Addi, &r.Remark, &r.Status).Scan(&r.ID)
	if err != nil {
		return errors.Wrap(err, "failed to insert t_relation")
	}
	return nil
}

// GetTRelationByPk select the TRelation from the database.
func GetTRelationByPk(db Queryer, pk0 null.Int) (*TRelation, error) {
	var r TRelation
	err := db.QueryRow(
		`SELECT id, left_id, left_type, left_key, left_key_type, kind, right_id, right_type, right_value_type, right_value, rule_area, creator, create_time, regenerator, update_time, domain_id, addi, remark, status FROM t_relation WHERE id = $1`,
		pk0).Scan(&r.ID, &r.LeftID, &r.LeftType, &r.LeftKey, &r.LeftKeyType, &r.Kind, &r.RightID, &r.RightType, &r.RightValueType, &r.RightValue, &r.RuleArea, &r.Creator, &r.CreateTime, &r.Regenerator, &r.UpdateTime, &r.DomainID, &r.Addi, &r.Remark, &r.Status)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select t_relation")
	}
	return &r, nil
}

/*TRelationHistory 描述两个实体间关系
校快保，描述销售/学校管理员/学校统计员与学校间的对应关系 represents kuser.t_relation_history */
type TRelationHistory struct {
	ID             null.Int       `json:"ID,omitempty" db:"id,false,integer"`                                     /* id 编号 */
	LeftID         null.Int       `json:"LeftID,omitempty" db:"left_id,false,bigint"`                             /* left_id 左编号 */
	LeftType       null.String    `json:"LeftType,omitempty" db:"left_type,false,character varying"`              /* left_type 左类型，用户编号: t_user.id */
	LeftKey        null.String    `json:"LeftKey,omitempty" db:"left_key,false,character varying"`                /* left_key 左识别标识 */
	LeftKeyType    null.String    `json:"LeftKeyType,omitempty" db:"left_key_type,false,character varying"`       /* left_key_type 左识别标识类型，帐号: account, 邮箱: email, 手机: tel, 微信公众号openID: mp_open_id, 微信开放平台openID: wx_open_id */
	Kind           null.String    `json:"Kind,omitempty" db:"kind,false,character varying"`                       /* kind 关系类型 */
	RightID        null.Int       `json:"RightID,omitempty" db:"right_id,false,bigint"`                           /* right_id 目标资源编号，如学校ID */
	RightType      null.String    `json:"RightType,omitempty" db:"right_type,false,character varying"`            /* right_type school: 学校 */
	RightValueType null.String    `json:"RightValueType,omitempty" db:"right_value_type,false,character varying"` /* right_value_type right_value_type */
	RightValue     null.String    `json:"RightValue,omitempty" db:"right_value,false,character varying"`          /* right_value right_value */
	RuleArea       null.String    `json:"RuleArea,omitempty" db:"rule_area,false,character varying"`              /* rule_area 管辖地区 */
	Creator        null.Int       `json:"Creator,omitempty" db:"creator,false,bigint"`                            /* creator 创建者用户ID */
	CreateTime     null.Int       `json:"CreateTime,omitempty" db:"create_time,false,bigint"`                     /* create_time 创建时间 */
	Regenerator    null.Int       `json:"Regenerator,omitempty" db:"regenerator,false,bigint"`                    /* regenerator 更新者 */
	UpdateTime     null.Int       `json:"UpdateTime,omitempty" db:"update_time,false,bigint"`                     /* update_time 修改时间 */
	DomainID       null.Int       `json:"DomainID,omitempty" db:"domain_id,false,bigint"`                         /* domain_id 数据属主 */
	Addi           types.JSONText `json:"Addi,omitempty" db:"addi,false,jsonb"`                                   /* addi 附加数据 */
	Remark         null.String    `json:"Remark,omitempty" db:"remark,false,character varying"`                   /* remark 备注 */
	Status         null.String    `json:"Status,omitempty" db:"status,false,character varying"`                   /* status 状态 */
	Filter                        // build DML where clause
}

//TRelationHistoryFields full field list for default query
var TRelationHistoryFields = []string{
	"ID",
	"LeftID",
	"LeftType",
	"LeftKey",
	"LeftKeyType",
	"Kind",
	"RightID",
	"RightType",
	"RightValueType",
	"RightValue",
	"RuleArea",
	"Creator",
	"CreateTime",
	"Regenerator",
	"UpdateTime",
	"DomainID",
	"Addi",
	"Remark",
	"Status",
}

//Fields return all fields of struct.
func (r *TRelationHistory) Fields() []string {
	return TRelationHistoryFields
}

//GetTableName return the associated db table name.
func (r *TRelationHistory) GetTableName() string {
	var viewNamePattern = regexp.MustCompile(`(?i)^t_v_[a-z0-9_]+$`)
	tableName := "t_relation_history"
	if viewNamePattern.MatchString(tableName) {
		return tableName[2:]
	}
	return tableName
}

// Create inserts the TRelationHistory to the database.
func (r *TRelationHistory) Create(db Queryer) error {
	_, err := db.Exec(
		`INSERT INTO t_relation_history (id, left_id, left_type, left_key, left_key_type, kind, right_id, right_type, right_value_type, right_value, rule_area, creator, create_time, regenerator, update_time, domain_id, addi, remark, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)`,
		&r.ID, &r.LeftID, &r.LeftType, &r.LeftKey, &r.LeftKeyType, &r.Kind, &r.RightID, &r.RightType, &r.RightValueType, &r.RightValue, &r.RuleArea, &r.Creator, &r.CreateTime, &r.Regenerator, &r.UpdateTime, &r.DomainID, &r.Addi, &r.Remark, &r.Status)
	if err != nil {
		return errors.Wrap(err, "failed to insert t_relation_history")
	}
	return nil
}

// GetTRelationHistoryByPk select the TRelationHistory from the database.
func GetTRelationHistoryByPk(db Queryer) (*TRelationHistory, error) {
	var r TRelationHistory
	err := db.QueryRow(
		`SELECT id, left_id, left_type, left_key, left_key_type, kind, right_id, right_type, right_value_type, right_value, rule_area, creator, create_time, regenerator, update_time, domain_id, addi, remark, status FROM t_relation_history WHERE `,
	).Scan(&r.ID, &r.LeftID, &r.LeftType, &r.LeftKey, &r.LeftKeyType, &r.Kind, &r.RightID, &r.RightType, &r.RightValueType, &r.RightValue, &r.RuleArea, &r.Creator, &r.CreateTime, &r.Regenerator, &r.UpdateTime, &r.DomainID, &r.Addi, &r.Remark, &r.Status)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select t_relation_history")
	}
	return &r, nil
}

/*TReportClaims 报案理赔 represents kuser.t_report_claims */
type TReportClaims struct {
	ID                     null.Int       `json:"ID,omitempty" db:"id,true,integer"`                                            /* id 编号 */
	InformantID            null.Int       `json:"InformantID,omitempty" db:"informant_id,false,bigint"`                         /* informant_id 报案人编号 */
	Informant              types.JSONText `json:"Informant,omitempty" db:"informant,false,jsonb"`                               /* informant 报案人 */
	InsuredID              null.Int       `json:"InsuredID,omitempty" db:"insured_id,false,bigint"`                             /* insured_id 被保险人编号 */
	Insured                types.JSONText `json:"Insured,omitempty" db:"insured,false,jsonb"`                                   /* insured 被保险人 */
	InsuranceType          null.String    `json:"InsuranceType,omitempty" db:"insurance_type,false,character varying"`          /* insurance_type 保险类型 */
	InsuranceTypeID        null.Int       `json:"InsuranceTypeID,omitempty" db:"insurance_type_id,false,bigint"`                /* insurance_type_id 险种类别 */
	InsurancePolicySn      null.String    `json:"InsurancePolicySn,omitempty" db:"insurance_policy_sn,false,character varying"` /* insurance_policy_sn 保单号 */
	InsurancePolicyID      null.Int       `json:"InsurancePolicyID,omitempty" db:"insurance_policy_id,false,bigint"`            /* insurance_policy_id 保单编号  */
	InsurancePolicyStart   null.Int       `json:"InsurancePolicyStart,omitempty" db:"insurance_policy_start,false,bigint"`      /* insurance_policy_start 起保时间 */
	InsurancePolicyCease   null.Int       `json:"InsurancePolicyCease,omitempty" db:"insurance_policy_cease,false,bigint"`      /* insurance_policy_cease 脱保时间 */
	ReportSn               null.String    `json:"ReportSn,omitempty" db:"report_sn,false,character varying"`                    /* report_sn 报案号 */
	InsuredChannel         string         `json:"InsuredChannel,omitempty" db:"insured_channel,false,character varying"`        /* insured_channel 投保渠道[mp,web] */
	InsuredOrg             null.String    `json:"InsuredOrg,omitempty" db:"insured_org,false,character varying"`                /* insured_org 投保机构 */
	Treatment              null.String    `json:"Treatment,omitempty" db:"treatment,false,character varying"`                   /* treatment 治疗方式 */
	Hospital               null.String    `json:"Hospital,omitempty" db:"hospital,false,character varying"`                     /* hospital 就诊医院 */
	InjuredLocation        null.String    `json:"InjuredLocation,omitempty" db:"injured_location,false,character varying"`      /* injured_location 受伤地点 */
	InjuredPart            null.String    `json:"InjuredPart,omitempty" db:"injured_part,false,character varying"`              /* injured_part 受伤部位 */
	Reason                 null.String    `json:"Reason,omitempty" db:"reason,false,character varying"`                         /* reason 受伤原因 */
	InjuredDesc            null.String    `json:"InjuredDesc,omitempty" db:"injured_desc,false,character varying"`              /* injured_desc 受伤过程描述 */
	CreditCode             null.String    `json:"CreditCode,omitempty" db:"credit_code,false,character varying"`                /* credit_code 统一社会信用代码 */
	BankAccountType        null.String    `json:"BankAccountType,omitempty" db:"bank_account_type,false,character varying"`     /* bank_account_type 银行账户类型 */
	BankAccountName        null.String    `json:"BankAccountName,omitempty" db:"bank_account_name,false,character varying"`     /* bank_account_name 银行账户名 */
	BankName               null.String    `json:"BankName,omitempty" db:"bank_name,false,character varying"`                    /* bank_name 开户行 */
	BankAccountID          null.String    `json:"BankAccountID,omitempty" db:"bank_account_id,false,character varying"`         /* bank_account_id 银行卡号/账号 */
	BankCardPic            types.JSONText `json:"BankCardPic,omitempty" db:"bank_card_pic,false,jsonb"`                         /* bank_card_pic 银行卡/存折照片 */
	InjuredIDPic           types.JSONText `json:"InjuredIDPic,omitempty" db:"injured_id_pic,false,jsonb"`                       /* injured_id_pic 被保险人身份证照片 */
	GuardianIDPic          types.JSONText `json:"GuardianIDPic,omitempty" db:"guardian_id_pic,false,jsonb"`                     /* guardian_id_pic 监护人身份证照片 */
	OrgLicPic              types.JSONText `json:"OrgLicPic,omitempty" db:"org_lic_pic,false,jsonb"`                             /* org_lic_pic 营业执照照片 */
	RelationProvePic       types.JSONText `json:"RelationProvePic,omitempty" db:"relation_prove_pic,false,jsonb"`               /* relation_prove_pic 与被保险人关系证明照片 */
	BillsPic               types.JSONText `json:"BillsPic,omitempty" db:"bills_pic,false,jsonb"`                                /* bills_pic 门诊费用清单照片 */
	HospitalizedBillsPic   types.JSONText `json:"HospitalizedBillsPic,omitempty" db:"hospitalized_bills_pic,false,jsonb"`       /* hospitalized_bills_pic 住院费用清单照片 */
	InvoicePic             types.JSONText `json:"InvoicePic,omitempty" db:"invoice_pic,false,jsonb"`                            /* invoice_pic 医疗费用发票照片 */
	MedicalRecordPic       types.JSONText `json:"MedicalRecordPic,omitempty" db:"medical_record_pic,false,jsonb"`               /* medical_record_pic 病历照片 */
	DignosticInspectionPic types.JSONText `json:"DignosticInspectionPic,omitempty" db:"dignostic_inspection_pic,false,jsonb"`   /* dignostic_inspection_pic 检验检查报告照片 */
	DischargeAbstractPic   types.JSONText `json:"DischargeAbstractPic,omitempty" db:"discharge_abstract_pic,false,jsonb"`       /* discharge_abstract_pic 出院小结照片 */
	OtherPic               types.JSONText `json:"OtherPic,omitempty" db:"other_pic,false,jsonb"`                                /* other_pic 其它资料照片 */
	CourierSnPic           types.JSONText `json:"CourierSnPic,omitempty" db:"courier_sn_pic,false,jsonb"`                       /* courier_sn_pic 快递单号照片 */
	PaidNoticePic          types.JSONText `json:"PaidNoticePic,omitempty" db:"paid_notice_pic,false,jsonb"`                     /* paid_notice_pic 保险金给付通知书 */
	ClaimApplyPic          types.JSONText `json:"ClaimApplyPic,omitempty" db:"claim_apply_pic,false,jsonb"`                     /* claim_apply_pic 索赔申请书 */
	EquityTransferFile     types.JSONText `json:"EquityTransferFile,omitempty" db:"equity_transfer_file,false,jsonb"`           /* equity_transfer_file 权益转让书
	 */
	MatchProgrammePic        types.JSONText `json:"MatchProgrammePic,omitempty" db:"match_programme_pic,false,jsonb"`               /* match_programme_pic 已有投保单位盖章的比赛秩序册 */
	PolicyFile               types.JSONText `json:"PolicyFile,omitempty" db:"policy_file,false,jsonb"`                              /* policy_file 保单文件 */
	AddiPic                  types.JSONText `json:"AddiPic,omitempty" db:"addi_pic,false,jsonb"`                                    /* addi_pic 补充资料照片 */
	CourierSn                null.String    `json:"CourierSn,omitempty" db:"courier_sn,false,character varying"`                    /* courier_sn 快递单号 */
	ReplyAddr                null.String    `json:"ReplyAddr,omitempty" db:"reply_addr,false,character varying"`                    /* reply_addr 资料回寄地址 */
	InjuredTime              null.Int       `json:"InjuredTime,omitempty" db:"injured_time,false,bigint"`                           /* injured_time 受伤时间 */
	ReportTime               null.Int       `json:"ReportTime,omitempty" db:"report_time,false,bigint"`                             /* report_time 报案时间 */
	ReplyTime                null.Int       `json:"ReplyTime,omitempty" db:"reply_time,false,bigint"`                               /* reply_time 回复时间 */
	ClaimsMatAddTime         null.Int       `json:"ClaimsMatAddTime,omitempty" db:"claims_mat_add_time,false,bigint"`               /* claims_mat_add_time 索赔资料提交时间 */
	MatReturnDate            null.Int       `json:"MatReturnDate,omitempty" db:"mat_return_date,false,bigint"`                      /* mat_return_date 发票寄回时间 */
	CloseDate                null.Int       `json:"CloseDate,omitempty" db:"close_date,false,bigint"`                               /* close_date 结案日期 */
	FaceAmount               null.Float     `json:"FaceAmount,omitempty" db:"face_amount,false,double precision"`                   /* face_amount 发票金额 */
	MediAssureAmount         null.Float     `json:"MediAssureAmount,omitempty" db:"medi_assure_amount,false,double precision"`      /* medi_assure_amount 医保统筹金额 */
	ThirdPayAmount           null.Float     `json:"ThirdPayAmount,omitempty" db:"third_pay_amount,false,double precision"`          /* third_pay_amount 第三方赔付金额 */
	ClaimAmount              null.Float     `json:"ClaimAmount,omitempty" db:"claim_amount,false,double precision"`                 /* claim_amount 赔付金额 */
	OccurrReason             null.String    `json:"OccurrReason,omitempty" db:"occurr_reason,false,character varying"`              /* occurr_reason 出险原因 */
	TreatmentResult          null.String    `json:"TreatmentResult,omitempty" db:"treatment_result,false,character varying"`        /* treatment_result 治疗结果 */
	DiseaseDiagnosisPic      types.JSONText `json:"DiseaseDiagnosisPic,omitempty" db:"disease_diagnosis_pic,false,jsonb"`           /* disease_diagnosis_pic 诊断证明 */
	DisabilityCertificate    types.JSONText `json:"DisabilityCertificate,omitempty" db:"disability_certificate,false,jsonb"`        /* disability_certificate 残疾证明 */
	DeathCertificate         types.JSONText `json:"DeathCertificate,omitempty" db:"death_certificate,false,jsonb"`                  /* death_certificate 死亡证明 */
	StudentStatusCertificate types.JSONText `json:"StudentStatusCertificate,omitempty" db:"student_status_certificate,false,jsonb"` /* student_status_certificate 学籍证明 */
	RefuseDesc               null.String    `json:"RefuseDesc,omitempty" db:"refuse_desc,false,character varying"`                  /* refuse_desc 拒绝理由 */
	DomainID                 null.Int       `json:"DomainID,omitempty" db:"domain_id,false,bigint"`                                 /* domain_id 数据属主 */
	Creator                  null.Int       `json:"Creator,omitempty" db:"creator,false,bigint"`                                    /* creator 创建者用户ID */
	CreateTime               null.Int       `json:"CreateTime,omitempty" db:"create_time,false,bigint"`                             /* create_time 创建时间 */
	Regenerator              null.Int       `json:"Regenerator,omitempty" db:"regenerator,false,bigint"`                            /* regenerator 更新者 */
	UpdateTime               null.Int       `json:"UpdateTime,omitempty" db:"update_time,false,bigint"`                             /* update_time 修改时间 */
	Addi                     types.JSONText `json:"Addi,omitempty" db:"addi,false,jsonb"`                                           /* addi 附加 */
	Remark                   null.String    `json:"Remark,omitempty" db:"remark,false,character varying"`                           /* remark 备注 */
	Status                   null.String    `json:"Status,omitempty" db:"status,false,character varying"`                           /* status 状态:，2: 已报案，等待上传索赔资料，4: 受理中, 6: 等待补充资料, 8: 已结案, 10: 撤销报案, 12: 拒赔 */
	Filter                                  // build DML where clause
}

//TReportClaimsFields full field list for default query
var TReportClaimsFields = []string{
	"ID",
	"InformantID",
	"Informant",
	"InsuredID",
	"Insured",
	"InsuranceType",
	"InsuranceTypeID",
	"InsurancePolicySn",
	"InsurancePolicyID",
	"InsurancePolicyStart",
	"InsurancePolicyCease",
	"ReportSn",
	"InsuredChannel",
	"InsuredOrg",
	"Treatment",
	"Hospital",
	"InjuredLocation",
	"InjuredPart",
	"Reason",
	"InjuredDesc",
	"CreditCode",
	"BankAccountType",
	"BankAccountName",
	"BankName",
	"BankAccountID",
	"BankCardPic",
	"InjuredIDPic",
	"GuardianIDPic",
	"OrgLicPic",
	"RelationProvePic",
	"BillsPic",
	"HospitalizedBillsPic",
	"InvoicePic",
	"MedicalRecordPic",
	"DignosticInspectionPic",
	"DischargeAbstractPic",
	"OtherPic",
	"CourierSnPic",
	"PaidNoticePic",
	"ClaimApplyPic",
	"EquityTransferFile",
	"MatchProgrammePic",
	"PolicyFile",
	"AddiPic",
	"CourierSn",
	"ReplyAddr",
	"InjuredTime",
	"ReportTime",
	"ReplyTime",
	"ClaimsMatAddTime",
	"MatReturnDate",
	"CloseDate",
	"FaceAmount",
	"MediAssureAmount",
	"ThirdPayAmount",
	"ClaimAmount",
	"OccurrReason",
	"TreatmentResult",
	"DiseaseDiagnosisPic",
	"DisabilityCertificate",
	"DeathCertificate",
	"StudentStatusCertificate",
	"RefuseDesc",
	"DomainID",
	"Creator",
	"CreateTime",
	"Regenerator",
	"UpdateTime",
	"Addi",
	"Remark",
	"Status",
}

//Fields return all fields of struct.
func (r *TReportClaims) Fields() []string {
	return TReportClaimsFields
}

//GetTableName return the associated db table name.
func (r *TReportClaims) GetTableName() string {
	var viewNamePattern = regexp.MustCompile(`(?i)^t_v_[a-z0-9_]+$`)
	tableName := "t_report_claims"
	if viewNamePattern.MatchString(tableName) {
		return tableName[2:]
	}
	return tableName
}

// Create inserts the TReportClaims to the database.
func (r *TReportClaims) Create(db Queryer) error {
	err := db.QueryRow(
		`INSERT INTO t_report_claims (informant_id, informant, insured_id, insured, insurance_type, insurance_type_id, insurance_policy_sn, insurance_policy_id, insurance_policy_start, insurance_policy_cease, report_sn, insured_channel, insured_org, treatment, hospital, injured_location, injured_part, reason, injured_desc, credit_code, bank_account_type, bank_account_name, bank_name, bank_account_id, bank_card_pic, injured_id_pic, guardian_id_pic, org_lic_pic, relation_prove_pic, bills_pic, hospitalized_bills_pic, invoice_pic, medical_record_pic, dignostic_inspection_pic, discharge_abstract_pic, other_pic, courier_sn_pic, paid_notice_pic, claim_apply_pic, equity_transfer_file, match_programme_pic, policy_file, addi_pic, courier_sn, reply_addr, injured_time, report_time, reply_time, claims_mat_add_time, mat_return_date, close_date, face_amount, medi_assure_amount, third_pay_amount, claim_amount, occurr_reason, treatment_result, disease_diagnosis_pic, disability_certificate, death_certificate, student_status_certificate, refuse_desc, domain_id, creator, create_time, regenerator, update_time, addi, remark, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36, $37, $38, $39, $40, $41, $42, $43, $44, $45, $46, $47, $48, $49, $50, $51, $52, $53, $54, $55, $56, $57, $58, $59, $60, $61, $62, $63, $64, $65, $66, $67, $68, $69, $70) RETURNING id`,
		&r.InformantID, &r.Informant, &r.InsuredID, &r.Insured, &r.InsuranceType, &r.InsuranceTypeID, &r.InsurancePolicySn, &r.InsurancePolicyID, &r.InsurancePolicyStart, &r.InsurancePolicyCease, &r.ReportSn, &r.InsuredChannel, &r.InsuredOrg, &r.Treatment, &r.Hospital, &r.InjuredLocation, &r.InjuredPart, &r.Reason, &r.InjuredDesc, &r.CreditCode, &r.BankAccountType, &r.BankAccountName, &r.BankName, &r.BankAccountID, &r.BankCardPic, &r.InjuredIDPic, &r.GuardianIDPic, &r.OrgLicPic, &r.RelationProvePic, &r.BillsPic, &r.HospitalizedBillsPic, &r.InvoicePic, &r.MedicalRecordPic, &r.DignosticInspectionPic, &r.DischargeAbstractPic, &r.OtherPic, &r.CourierSnPic, &r.PaidNoticePic, &r.ClaimApplyPic, &r.EquityTransferFile, &r.MatchProgrammePic, &r.PolicyFile, &r.AddiPic, &r.CourierSn, &r.ReplyAddr, &r.InjuredTime, &r.ReportTime, &r.ReplyTime, &r.ClaimsMatAddTime, &r.MatReturnDate, &r.CloseDate, &r.FaceAmount, &r.MediAssureAmount, &r.ThirdPayAmount, &r.ClaimAmount, &r.OccurrReason, &r.TreatmentResult, &r.DiseaseDiagnosisPic, &r.DisabilityCertificate, &r.DeathCertificate, &r.StudentStatusCertificate, &r.RefuseDesc, &r.DomainID, &r.Creator, &r.CreateTime, &r.Regenerator, &r.UpdateTime, &r.Addi, &r.Remark, &r.Status).Scan(&r.ID)
	if err != nil {
		return errors.Wrap(err, "failed to insert t_report_claims")
	}
	return nil
}

// GetTReportClaimsByPk select the TReportClaims from the database.
func GetTReportClaimsByPk(db Queryer, pk0 null.Int) (*TReportClaims, error) {
	var r TReportClaims
	err := db.QueryRow(
		`SELECT id, informant_id, informant, insured_id, insured, insurance_type, insurance_type_id, insurance_policy_sn, insurance_policy_id, insurance_policy_start, insurance_policy_cease, report_sn, insured_channel, insured_org, treatment, hospital, injured_location, injured_part, reason, injured_desc, credit_code, bank_account_type, bank_account_name, bank_name, bank_account_id, bank_card_pic, injured_id_pic, guardian_id_pic, org_lic_pic, relation_prove_pic, bills_pic, hospitalized_bills_pic, invoice_pic, medical_record_pic, dignostic_inspection_pic, discharge_abstract_pic, other_pic, courier_sn_pic, paid_notice_pic, claim_apply_pic, equity_transfer_file, match_programme_pic, policy_file, addi_pic, courier_sn, reply_addr, injured_time, report_time, reply_time, claims_mat_add_time, mat_return_date, close_date, face_amount, medi_assure_amount, third_pay_amount, claim_amount, occurr_reason, treatment_result, disease_diagnosis_pic, disability_certificate, death_certificate, student_status_certificate, refuse_desc, domain_id, creator, create_time, regenerator, update_time, addi, remark, status FROM t_report_claims WHERE id = $1`,
		pk0).Scan(&r.ID, &r.InformantID, &r.Informant, &r.InsuredID, &r.Insured, &r.InsuranceType, &r.InsuranceTypeID, &r.InsurancePolicySn, &r.InsurancePolicyID, &r.InsurancePolicyStart, &r.InsurancePolicyCease, &r.ReportSn, &r.InsuredChannel, &r.InsuredOrg, &r.Treatment, &r.Hospital, &r.InjuredLocation, &r.InjuredPart, &r.Reason, &r.InjuredDesc, &r.CreditCode, &r.BankAccountType, &r.BankAccountName, &r.BankName, &r.BankAccountID, &r.BankCardPic, &r.InjuredIDPic, &r.GuardianIDPic, &r.OrgLicPic, &r.RelationProvePic, &r.BillsPic, &r.HospitalizedBillsPic, &r.InvoicePic, &r.MedicalRecordPic, &r.DignosticInspectionPic, &r.DischargeAbstractPic, &r.OtherPic, &r.CourierSnPic, &r.PaidNoticePic, &r.ClaimApplyPic, &r.EquityTransferFile, &r.MatchProgrammePic, &r.PolicyFile, &r.AddiPic, &r.CourierSn, &r.ReplyAddr, &r.InjuredTime, &r.ReportTime, &r.ReplyTime, &r.ClaimsMatAddTime, &r.MatReturnDate, &r.CloseDate, &r.FaceAmount, &r.MediAssureAmount, &r.ThirdPayAmount, &r.ClaimAmount, &r.OccurrReason, &r.TreatmentResult, &r.DiseaseDiagnosisPic, &r.DisabilityCertificate, &r.DeathCertificate, &r.StudentStatusCertificate, &r.RefuseDesc, &r.DomainID, &r.Creator, &r.CreateTime, &r.Regenerator, &r.UpdateTime, &r.Addi, &r.Remark, &r.Status)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select t_report_claims")
	}
	return &r, nil
}

/*TResource 资源列表 represents kuser.t_resource */
type TResource struct {
	ID              null.Int       `json:"ID,omitempty" db:"id,true,integer"`                             /* id 资源编号 */
	InsuranceTypeID null.Int       `json:"InsuranceTypeID,omitempty" db:"insurance_type_id,false,bigint"` /* insurance_type_id 保险类型ID */
	Name            null.String    `json:"Name,omitempty" db:"name,false,character varying"`              /* name 资源名称 */
	Content         null.String    `json:"Content,omitempty" db:"content,false,character varying"`        /* content 资源内容 */
	Link            types.JSONText `json:"Link,omitempty" db:"link,false,jsonb"`                          /* link 链接 */
	Picture         types.JSONText `json:"Picture,omitempty" db:"picture,false,jsonb"`                    /* picture 图片 */
	Tag             null.String    `json:"Tag,omitempty" db:"tag,false,character varying"`                /* tag 标签 */
	IsTop           null.Bool      `json:"IsTop,omitempty" db:"is_top,false,boolean"`                     /* is_top 是否首页显示：用户进入智能客服后直接显示 */
	IsPolicy        null.Bool      `json:"IsPolicy,omitempty" db:"is_policy,false,boolean"`               /* is_policy 判断是否是保险条款 */
	Updator         null.String    `json:"Updator,omitempty" db:"updator,false,character varying"`        /* updator 更新者 */
	Regenerator     null.Int       `json:"Regenerator,omitempty" db:"regenerator,false,bigint"`           /* regenerator 更新者 */
	UpdateTime      null.Int       `json:"UpdateTime,omitempty" db:"update_time,false,bigint"`            /* update_time 更新时间 */
	Creator         null.String    `json:"Creator,omitempty" db:"creator,false,character varying"`        /* creator 创建者账号 */
	CreateTime      null.Int       `json:"CreateTime,omitempty" db:"create_time,false,bigint"`            /* create_time 创建时间 */
	DomainID        null.Int       `json:"DomainID,omitempty" db:"domain_id,false,bigint"`                /* domain_id 数据属主 */
	Addi            types.JSONText `json:"Addi,omitempty" db:"addi,false,jsonb"`                          /* addi 附加数据 */
	Remark          null.String    `json:"Remark,omitempty" db:"remark,false,character varying"`          /* remark 备注 */
	Status          null.String    `json:"Status,omitempty" db:"status,false,character varying"`          /* status 状态0:有效, 2:修改，4删除 */
	Filter                         // build DML where clause
}

//TResourceFields full field list for default query
var TResourceFields = []string{
	"ID",
	"InsuranceTypeID",
	"Name",
	"Content",
	"Link",
	"Picture",
	"Tag",
	"IsTop",
	"IsPolicy",
	"Updator",
	"Regenerator",
	"UpdateTime",
	"Creator",
	"CreateTime",
	"DomainID",
	"Addi",
	"Remark",
	"Status",
}

//Fields return all fields of struct.
func (r *TResource) Fields() []string {
	return TResourceFields
}

//GetTableName return the associated db table name.
func (r *TResource) GetTableName() string {
	var viewNamePattern = regexp.MustCompile(`(?i)^t_v_[a-z0-9_]+$`)
	tableName := "t_resource"
	if viewNamePattern.MatchString(tableName) {
		return tableName[2:]
	}
	return tableName
}

// Create inserts the TResource to the database.
func (r *TResource) Create(db Queryer) error {
	err := db.QueryRow(
		`INSERT INTO t_resource (insurance_type_id, name, content, link, picture, tag, is_top, is_policy, updator, regenerator, update_time, creator, create_time, domain_id, addi, remark, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17) RETURNING id`,
		&r.InsuranceTypeID, &r.Name, &r.Content, &r.Link, &r.Picture, &r.Tag, &r.IsTop, &r.IsPolicy, &r.Updator, &r.Regenerator, &r.UpdateTime, &r.Creator, &r.CreateTime, &r.DomainID, &r.Addi, &r.Remark, &r.Status).Scan(&r.ID)
	if err != nil {
		return errors.Wrap(err, "failed to insert t_resource")
	}
	return nil
}

// GetTResourceByPk select the TResource from the database.
func GetTResourceByPk(db Queryer, pk0 null.Int) (*TResource, error) {
	var r TResource
	err := db.QueryRow(
		`SELECT id, insurance_type_id, name, content, link, picture, tag, is_top, is_policy, updator, regenerator, update_time, creator, create_time, domain_id, addi, remark, status FROM t_resource WHERE id = $1`,
		pk0).Scan(&r.ID, &r.InsuranceTypeID, &r.Name, &r.Content, &r.Link, &r.Picture, &r.Tag, &r.IsTop, &r.IsPolicy, &r.Updator, &r.Regenerator, &r.UpdateTime, &r.Creator, &r.CreateTime, &r.DomainID, &r.Addi, &r.Remark, &r.Status)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select t_resource")
	}
	return &r, nil
}

/*TScanTdc 请求二维码记录 represents kuser.t_scan_tdc */
type TScanTdc struct {
	ID         null.Int    `json:"ID,omitempty" db:"id,true,integer"`                             /* id 二维码编号 */
	TdcID      null.Int    `json:"TdcID,omitempty" db:"tdc_id,false,bigint"`                      /* tdc_id 二维码编号 */
	ExternalID null.String `json:"ExternalID,omitempty" db:"external_id,false,character varying"` /* external_id 外部平台ID */
	ReqTime    null.Int    `json:"ReqTime,omitempty" db:"req_time,false,bigint"`                  /* req_time 请求二维码时间 */
	ReqSrc     null.String `json:"ReqSrc,omitempty" db:"req_src,false,character varying"`         /* req_src 请求来源 */
	Filter                 // build DML where clause
}

//TScanTdcFields full field list for default query
var TScanTdcFields = []string{
	"ID",
	"TdcID",
	"ExternalID",
	"ReqTime",
	"ReqSrc",
}

//Fields return all fields of struct.
func (r *TScanTdc) Fields() []string {
	return TScanTdcFields
}

//GetTableName return the associated db table name.
func (r *TScanTdc) GetTableName() string {
	var viewNamePattern = regexp.MustCompile(`(?i)^t_v_[a-z0-9_]+$`)
	tableName := "t_scan_tdc"
	if viewNamePattern.MatchString(tableName) {
		return tableName[2:]
	}
	return tableName
}

// Create inserts the TScanTdc to the database.
func (r *TScanTdc) Create(db Queryer) error {
	err := db.QueryRow(
		`INSERT INTO t_scan_tdc (tdc_id, external_id, req_time, req_src) VALUES ($1, $2, $3, $4) RETURNING id`,
		&r.TdcID, &r.ExternalID, &r.ReqTime, &r.ReqSrc).Scan(&r.ID)
	if err != nil {
		return errors.Wrap(err, "failed to insert t_scan_tdc")
	}
	return nil
}

// GetTScanTdcByPk select the TScanTdc from the database.
func GetTScanTdcByPk(db Queryer, pk0 null.Int) (*TScanTdc, error) {
	var r TScanTdc
	err := db.QueryRow(
		`SELECT id, tdc_id, external_id, req_time, req_src FROM t_scan_tdc WHERE id = $1`,
		pk0).Scan(&r.ID, &r.TdcID, &r.ExternalID, &r.ReqTime, &r.ReqSrc)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select t_scan_tdc")
	}
	return &r, nil
}

/*TSchool 学校信息表，包含了销售经理，学校管理员，投保规则 represents kuser.t_school */
type TSchool struct {
	ID                      null.Int       `json:"ID,omitempty" db:"id,true,integer"`                                              /* id 学校编号 */
	Name                    string         `json:"Name,omitempty" db:"name,false,character varying"`                               /* name 名称 */
	OrgCode                 null.String    `json:"OrgCode,omitempty" db:"org_code,false,character varying"`                        /* org_code 机构代码 */
	Faculty                 types.JSONText `json:"Faculty,omitempty" db:"faculty,false,jsonb"`                                     /* faculty 学院 */
	Branches                types.JSONText `json:"Branches,omitempty" db:"branches,false,jsonb"`                                   /* branches 校区 */
	Category                string         `json:"Category,omitempty" db:"category,false,character varying"`                       /* category 类别:幼儿园，小学，初中，高中，大学 */
	Contact                 null.String    `json:"Contact,omitempty" db:"contact,false,character varying"`                         /* contact 联系人 */
	PostCode                null.String    `json:"PostCode,omitempty" db:"post_code,false,character varying"`                      /* post_code 邮编 */
	Phone                   null.String    `json:"Phone,omitempty" db:"phone,false,character varying"`                             /* phone 联系电话 */
	Addr                    null.String    `json:"Addr,omitempty" db:"addr,false,character varying"`                               /* addr 详细地址 */
	Province                null.String    `json:"Province,omitempty" db:"province,false,character varying"`                       /* province 省 */
	City                    null.String    `json:"City,omitempty" db:"city,false,character varying"`                               /* city 市 */
	District                null.String    `json:"District,omitempty" db:"district,false,character varying"`                       /* district 区/县 */
	Street                  null.String    `json:"Street,omitempty" db:"street,false,character varying"`                           /* street 街道/片区 */
	DataSyncTarget          null.String    `json:"DataSyncTarget,omitempty" db:"data_sync_target,false,character varying"`         /* data_sync_target 数据同步类型 */
	SaleManagers            types.JSONText `json:"SaleManagers,omitempty" db:"sale_managers,false,jsonb"`                          /* sale_managers 销售 */
	SchoolManagers          types.JSONText `json:"SchoolManagers,omitempty" db:"school_managers,false,jsonb"`                      /* school_managers 学校管理员 */
	PurchaseRule            types.JSONText `json:"PurchaseRule,omitempty" db:"purchase_rule,false,jsonb"`                          /* purchase_rule 投保规则 */
	BusinessDomain          null.String    `json:"BusinessDomain,omitempty" db:"business_domain,false,character varying"`          /* business_domain 营业性质：文体体育、广告、事业单位、政府机关、其它 */
	SchoolCategory          null.String    `json:"SchoolCategory,omitempty" db:"school_category,false,character varying"`          /* school_category 学校性质：民办，公办 */
	AllowBackdating         types.JSONText `json:"AllowBackdating,omitempty" db:"allow_backdating,false,jsonb"`                    /* allow_backdating 允许倒签 */
	UseCreditCode           null.Bool      `json:"UseCreditCode,omitempty" db:"use_credit_code,false,boolean"`                     /* use_credit_code 使用信用代码 */
	CreditCode              null.String    `json:"CreditCode,omitempty" db:"credit_code,false,character varying"`                  /* credit_code 统一社会信用代码 */
	CreditCodePic           types.JSONText `json:"CreditCodePic,omitempty" db:"credit_code_pic,false,jsonb"`                       /* credit_code_pic 统一社会信用代码证书，base64图片 */
	InvoiceTitle            null.String    `json:"InvoiceTitle,omitempty" db:"invoice_title,false,character varying"`              /* invoice_title 发票抬头 */
	IsCompulsory            null.Bool      `json:"IsCompulsory,omitempty" db:"is_compulsory,false,boolean"`                        /* is_compulsory 单位性质,true: 是义务教育，false: 不是非义务教育 */
	RegNum                  null.Int       `json:"RegNum,omitempty" db:"reg_num,false,integer"`                                    /* reg_num 注册人数 */
	CompulsoryStudentNum    null.Int       `json:"CompulsoryStudentNum,omitempty" db:"compulsory_student_num,false,bigint"`        /* compulsory_student_num 义务教育学生人数（校方） */
	NonCompulsoryStudentNum null.Int       `json:"NonCompulsoryStudentNum,omitempty" db:"non_compulsory_student_num,false,bigint"` /* non_compulsory_student_num 非义务教育人数（校方） */
	DinnerNum               null.Int       `json:"DinnerNum,omitempty" db:"dinner_num,false,integer"`                              /* dinner_num 用餐人数 */
	CanteenNum              null.Int       `json:"CanteenNum,omitempty" db:"canteen_num,false,integer"`                            /* canteen_num 食堂个数 */
	ShopNum                 null.Int       `json:"ShopNum,omitempty" db:"shop_num,false,integer"`                                  /* shop_num 商店个数 */
	Files                   types.JSONText `json:"Files,omitempty" db:"files,false,jsonb"`                                         /* files 附加文件 */
	ContactRole             null.String    `json:"ContactRole,omitempty" db:"contact_role,false,character varying"`                /* contact_role 联系人职位 */
	IsSchool                null.Bool      `json:"IsSchool,omitempty" db:"is_school,false,boolean"`                                /* is_school 是学校否 */
	Updator                 null.Int       `json:"Updator,omitempty" db:"updator,false,bigint"`                                    /* updator 更新人 */
	Creator                 null.Int       `json:"Creator,omitempty" db:"creator,false,bigint"`                                    /* creator 创建者用户ID */
	CreateTime              null.Int       `json:"CreateTime,omitempty" db:"create_time,false,bigint"`                             /* create_time 创建时间 */
	Regenerator             null.Int       `json:"Regenerator,omitempty" db:"regenerator,false,bigint"`                            /* regenerator 更新者 */
	UpdateTime              null.Int       `json:"UpdateTime,omitempty" db:"update_time,false,bigint"`                             /* update_time 更新时间 */
	DomainID                null.Int       `json:"DomainID,omitempty" db:"domain_id,false,bigint"`                                 /* domain_id 数据属主 */
	Addi                    types.JSONText `json:"Addi,omitempty" db:"addi,false,jsonb"`                                           /* addi 附加数据 */
	Remark                  null.String    `json:"Remark,omitempty" db:"remark,false,character varying"`                           /* remark 备注 */
	Status                  null.String    `json:"Status,omitempty" db:"status,false,character varying"`                           /* status 状态, '0': 未启用, '2': 启用, '6': 作废 */
	Filter                                 // build DML where clause
}

//TSchoolFields full field list for default query
var TSchoolFields = []string{
	"ID",
	"Name",
	"OrgCode",
	"Faculty",
	"Branches",
	"Category",
	"Contact",
	"PostCode",
	"Phone",
	"Addr",
	"Province",
	"City",
	"District",
	"Street",
	"DataSyncTarget",
	"SaleManagers",
	"SchoolManagers",
	"PurchaseRule",
	"BusinessDomain",
	"SchoolCategory",
	"AllowBackdating",
	"UseCreditCode",
	"CreditCode",
	"CreditCodePic",
	"InvoiceTitle",
	"IsCompulsory",
	"RegNum",
	"CompulsoryStudentNum",
	"NonCompulsoryStudentNum",
	"DinnerNum",
	"CanteenNum",
	"ShopNum",
	"Files",
	"ContactRole",
	"IsSchool",
	"Updator",
	"Creator",
	"CreateTime",
	"Regenerator",
	"UpdateTime",
	"DomainID",
	"Addi",
	"Remark",
	"Status",
}

//Fields return all fields of struct.
func (r *TSchool) Fields() []string {
	return TSchoolFields
}

//GetTableName return the associated db table name.
func (r *TSchool) GetTableName() string {
	var viewNamePattern = regexp.MustCompile(`(?i)^t_v_[a-z0-9_]+$`)
	tableName := "t_school"
	if viewNamePattern.MatchString(tableName) {
		return tableName[2:]
	}
	return tableName
}

// Create inserts the TSchool to the database.
func (r *TSchool) Create(db Queryer) error {
	err := db.QueryRow(
		`INSERT INTO t_school (name, org_code, faculty, branches, category, contact, post_code, phone, addr, province, city, district, street, data_sync_target, sale_managers, school_managers, purchase_rule, business_domain, school_category, allow_backdating, use_credit_code, credit_code, credit_code_pic, invoice_title, is_compulsory, reg_num, compulsory_student_num, non_compulsory_student_num, dinner_num, canteen_num, shop_num, files, contact_role, is_school, updator, creator, create_time, regenerator, update_time, domain_id, addi, remark, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36, $37, $38, $39, $40, $41, $42, $43) RETURNING id`,
		&r.Name, &r.OrgCode, &r.Faculty, &r.Branches, &r.Category, &r.Contact, &r.PostCode, &r.Phone, &r.Addr, &r.Province, &r.City, &r.District, &r.Street, &r.DataSyncTarget, &r.SaleManagers, &r.SchoolManagers, &r.PurchaseRule, &r.BusinessDomain, &r.SchoolCategory, &r.AllowBackdating, &r.UseCreditCode, &r.CreditCode, &r.CreditCodePic, &r.InvoiceTitle, &r.IsCompulsory, &r.RegNum, &r.CompulsoryStudentNum, &r.NonCompulsoryStudentNum, &r.DinnerNum, &r.CanteenNum, &r.ShopNum, &r.Files, &r.ContactRole, &r.IsSchool, &r.Updator, &r.Creator, &r.CreateTime, &r.Regenerator, &r.UpdateTime, &r.DomainID, &r.Addi, &r.Remark, &r.Status).Scan(&r.ID)
	if err != nil {
		return errors.Wrap(err, "failed to insert t_school")
	}
	return nil
}

// GetTSchoolByPk select the TSchool from the database.
func GetTSchoolByPk(db Queryer, pk0 null.Int) (*TSchool, error) {
	var r TSchool
	err := db.QueryRow(
		`SELECT id, name, org_code, faculty, branches, category, contact, post_code, phone, addr, province, city, district, street, data_sync_target, sale_managers, school_managers, purchase_rule, business_domain, school_category, allow_backdating, use_credit_code, credit_code, credit_code_pic, invoice_title, is_compulsory, reg_num, compulsory_student_num, non_compulsory_student_num, dinner_num, canteen_num, shop_num, files, contact_role, is_school, updator, creator, create_time, regenerator, update_time, domain_id, addi, remark, status FROM t_school WHERE id = $1`,
		pk0).Scan(&r.ID, &r.Name, &r.OrgCode, &r.Faculty, &r.Branches, &r.Category, &r.Contact, &r.PostCode, &r.Phone, &r.Addr, &r.Province, &r.City, &r.District, &r.Street, &r.DataSyncTarget, &r.SaleManagers, &r.SchoolManagers, &r.PurchaseRule, &r.BusinessDomain, &r.SchoolCategory, &r.AllowBackdating, &r.UseCreditCode, &r.CreditCode, &r.CreditCodePic, &r.InvoiceTitle, &r.IsCompulsory, &r.RegNum, &r.CompulsoryStudentNum, &r.NonCompulsoryStudentNum, &r.DinnerNum, &r.CanteenNum, &r.ShopNum, &r.Files, &r.ContactRole, &r.IsSchool, &r.Updator, &r.Creator, &r.CreateTime, &r.Regenerator, &r.UpdateTime, &r.DomainID, &r.Addi, &r.Remark, &r.Status)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select t_school")
	}
	return &r, nil
}

/*TSpecialOrder t_special_order represents kuser.t_special_order */
type TSpecialOrder struct {
	ID            null.Int       `json:"ID,omitempty" db:"id,true,integer"`                                   /* id 订单编号 */
	IDCardNo      string         `json:"IDCardNo,omitempty" db:"id_card_no,false,character varying"`          /* id_card_no 身份证号码 */
	Name          string         `json:"Name,omitempty" db:"name,false,character varying"`                    /* name 姓名 */
	Grade         null.String    `json:"Grade,omitempty" db:"grade,false,character varying"`                  /* grade 年级 */
	District      null.String    `json:"District,omitempty" db:"district,false,character varying"`            /* district 校区 */
	Project       string         `json:"Project,omitempty" db:"project,false,character varying"`              /* project 收费项目 */
	Amount        null.Float     `json:"Amount,omitempty" db:"amount,false,double precision"`                 /* amount 应收金额 */
	PayTime       null.Int       `json:"PayTime,omitempty" db:"pay_time,false,bigint"`                        /* pay_time 支付时间 */
	OpenID        string         `json:"OpenID,omitempty" db:"open_id,false,character varying"`               /* open_id open id */
	TradeNo       null.String    `json:"TradeNo,omitempty" db:"trade_no,false,character varying"`             /* trade_no 外部订单号 */
	TransactionID null.String    `json:"TransactionID,omitempty" db:"transaction_id,false,character varying"` /* transaction_id 支付平台订单号 */
	RefundNo      null.String    `json:"RefundNo,omitempty" db:"refund_no,false,character varying"`           /* refund_no 退款单号 */
	RefundTime    null.Int       `json:"RefundTime,omitempty" db:"refund_time,false,bigint"`                  /* refund_time 退款时间 */
	Creator       null.Int       `json:"Creator,omitempty" db:"creator,false,bigint"`                         /* creator 创建者用户ID */
	CreateTime    null.Int       `json:"CreateTime,omitempty" db:"create_time,false,bigint"`                  /* create_time 创建时间 */
	Regenerator   null.Int       `json:"Regenerator,omitempty" db:"regenerator,false,bigint"`                 /* regenerator 更新者 */
	UpdateTime    null.Int       `json:"UpdateTime,omitempty" db:"update_time,false,bigint"`                  /* update_time 更新时间 */
	Updater       null.Int       `json:"Updater,omitempty" db:"updater,false,bigint"`                         /* updater 更新者 */
	DomainID      null.Int       `json:"DomainID,omitempty" db:"domain_id,false,bigint"`                      /* domain_id 数据属主 */
	Remark        null.String    `json:"Remark,omitempty" db:"remark,false,character varying"`                /* remark 备注 */
	Addi          types.JSONText `json:"Addi,omitempty" db:"addi,false,jsonb"`                                /* addi 附加数据 */
	Status        null.String    `json:"Status,omitempty" db:"status,false,character varying"`                /* status 状态,0：未支付，2：已支付，4：超时，6：作废 */
	Filter                       // build DML where clause
}

//TSpecialOrderFields full field list for default query
var TSpecialOrderFields = []string{
	"ID",
	"IDCardNo",
	"Name",
	"Grade",
	"District",
	"Project",
	"Amount",
	"PayTime",
	"OpenID",
	"TradeNo",
	"TransactionID",
	"RefundNo",
	"RefundTime",
	"Creator",
	"CreateTime",
	"Regenerator",
	"UpdateTime",
	"Updater",
	"DomainID",
	"Remark",
	"Addi",
	"Status",
}

//Fields return all fields of struct.
func (r *TSpecialOrder) Fields() []string {
	return TSpecialOrderFields
}

//GetTableName return the associated db table name.
func (r *TSpecialOrder) GetTableName() string {
	var viewNamePattern = regexp.MustCompile(`(?i)^t_v_[a-z0-9_]+$`)
	tableName := "t_special_order"
	if viewNamePattern.MatchString(tableName) {
		return tableName[2:]
	}
	return tableName
}

// Create inserts the TSpecialOrder to the database.
func (r *TSpecialOrder) Create(db Queryer) error {
	err := db.QueryRow(
		`INSERT INTO t_special_order (id_card_no, name, grade, district, project, amount, pay_time, open_id, trade_no, transaction_id, refund_no, refund_time, creator, create_time, regenerator, update_time, updater, domain_id, remark, addi, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21) RETURNING id`,
		&r.IDCardNo, &r.Name, &r.Grade, &r.District, &r.Project, &r.Amount, &r.PayTime, &r.OpenID, &r.TradeNo, &r.TransactionID, &r.RefundNo, &r.RefundTime, &r.Creator, &r.CreateTime, &r.Regenerator, &r.UpdateTime, &r.Updater, &r.DomainID, &r.Remark, &r.Addi, &r.Status).Scan(&r.ID)
	if err != nil {
		return errors.Wrap(err, "failed to insert t_special_order")
	}
	return nil
}

// GetTSpecialOrderByPk select the TSpecialOrder from the database.
func GetTSpecialOrderByPk(db Queryer, pk0 null.Int) (*TSpecialOrder, error) {
	var r TSpecialOrder
	err := db.QueryRow(
		`SELECT id, id_card_no, name, grade, district, project, amount, pay_time, open_id, trade_no, transaction_id, refund_no, refund_time, creator, create_time, regenerator, update_time, updater, domain_id, remark, addi, status FROM t_special_order WHERE id = $1`,
		pk0).Scan(&r.ID, &r.IDCardNo, &r.Name, &r.Grade, &r.District, &r.Project, &r.Amount, &r.PayTime, &r.OpenID, &r.TradeNo, &r.TransactionID, &r.RefundNo, &r.RefundTime, &r.Creator, &r.CreateTime, &r.Regenerator, &r.UpdateTime, &r.Updater, &r.DomainID, &r.Remark, &r.Addi, &r.Status)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select t_special_order")
	}
	return &r, nil
}

/*TSysVer 应用版本包含业务模型、前端、后端、配置文件等

1、业务模型版本在模型生成时建立；
2、后端模型版本在每次后端启动时建立或更新；
3、配置文件版本在每次后端启动时建立或更新；
4、前端版本在每次后端启动时建立或更新；

 represents kuser.t_sys_ver */
type TSysVer struct {
	ID         null.Int       `json:"ID,omitempty" db:"id,true,integer"`                             /* id 编码 */
	Tag        null.String    `json:"Tag,omitempty" db:"tag,false,character varying"`                /* tag 标识 */
	Name       null.String    `json:"Name,omitempty" db:"name,false,character varying"`              /* name 名称 */
	Ver        null.String    `json:"Ver,omitempty" db:"ver,false,character varying"`                /* ver 版本 */
	CreateTime null.String    `json:"CreateTime,omitempty" db:"create_time,false,character varying"` /* create_time 创建时间 */
	UpdateTime null.String    `json:"UpdateTime,omitempty" db:"update_time,false,character varying"` /* update_time 更新时间 */
	Addi       types.JSONText `json:"Addi,omitempty" db:"addi,false,jsonb"`                          /* addi 附加 */
	Remark     null.String    `json:"Remark,omitempty" db:"remark,false,character varying"`          /* remark 备注 */
	Status     null.String    `json:"Status,omitempty" db:"status,false,character varying"`          /* status 状态 */
	Filter                    // build DML where clause
}

//TSysVerFields full field list for default query
var TSysVerFields = []string{
	"ID",
	"Tag",
	"Name",
	"Ver",
	"CreateTime",
	"UpdateTime",
	"Addi",
	"Remark",
	"Status",
}

//Fields return all fields of struct.
func (r *TSysVer) Fields() []string {
	return TSysVerFields
}

//GetTableName return the associated db table name.
func (r *TSysVer) GetTableName() string {
	var viewNamePattern = regexp.MustCompile(`(?i)^t_v_[a-z0-9_]+$`)
	tableName := "t_sys_ver"
	if viewNamePattern.MatchString(tableName) {
		return tableName[2:]
	}
	return tableName
}

// Create inserts the TSysVer to the database.
func (r *TSysVer) Create(db Queryer) error {
	err := db.QueryRow(
		`INSERT INTO t_sys_ver (tag, name, ver, create_time, update_time, addi, remark, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`,
		&r.Tag, &r.Name, &r.Ver, &r.CreateTime, &r.UpdateTime, &r.Addi, &r.Remark, &r.Status).Scan(&r.ID)
	if err != nil {
		return errors.Wrap(err, "failed to insert t_sys_ver")
	}
	return nil
}

// GetTSysVerByPk select the TSysVer from the database.
func GetTSysVerByPk(db Queryer, pk0 null.Int) (*TSysVer, error) {
	var r TSysVer
	err := db.QueryRow(
		`SELECT id, tag, name, ver, create_time, update_time, addi, remark, status FROM t_sys_ver WHERE id = $1`,
		pk0).Scan(&r.ID, &r.Tag, &r.Name, &r.Ver, &r.CreateTime, &r.UpdateTime, &r.Addi, &r.Remark, &r.Status)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select t_sys_ver")
	}
	return &r, nil
}

/*TTdc two-dimension-code represents kuser.t_tdc */
type TTdc struct {
	ID         null.Int       `json:"ID,omitempty" db:"id,true,integer"`                         /* id 编码 */
	TdcID      null.String    `json:"TdcID,omitempty" db:"tdc_id,false,character varying"`       /* tdc_id 二维码标识 */
	Name       null.String    `json:"Name,omitempty" db:"name,false,character varying"`          /* name 名称 */
	Issuer     null.Int       `json:"Issuer,omitempty" db:"issuer,false,bigint"`                 /* issuer 发布者 */
	IssueTime  null.Int       `json:"IssueTime,omitempty" db:"issue_time,false,bigint"`          /* issue_time 发布时间 */
	Limn       null.String    `json:"Limn,omitempty" db:"limn,false,character varying"`          /* limn 描述 */
	Data       types.JSONText `json:"Data,omitempty" db:"data,false,jsonb"`                      /* data 附加数据 */
	Expiration null.Int       `json:"Expiration,omitempty" db:"expiration,false,bigint"`         /* expiration 过期时间 */
	Type       null.String    `json:"Type,omitempty" db:"type,false,character varying"`          /* type 类型 */
	GotoView   null.String    `json:"GotoView,omitempty" db:"goto_view,false,character varying"` /* goto_view 扫描后的目标页面 */
	Requested  null.Int       `json:"Requested,omitempty" db:"requested,false,integer"`          /* requested 使用次数 */
	Accepted   null.Int       `json:"Accepted,omitempty" db:"accepted,false,smallint"`           /* accepted 成功使用次数 */
	Remark     null.String    `json:"Remark,omitempty" db:"remark,false,character varying"`      /* remark 备注 */
	Status     null.String    `json:"Status,omitempty" db:"status,false,character varying"`      /* status enabled,有效
	disabled,无效
	expired,过期、无效 */
	Filter // build DML where clause
}

//TTdcFields full field list for default query
var TTdcFields = []string{
	"ID",
	"TdcID",
	"Name",
	"Issuer",
	"IssueTime",
	"Limn",
	"Data",
	"Expiration",
	"Type",
	"GotoView",
	"Requested",
	"Accepted",
	"Remark",
	"Status",
}

//Fields return all fields of struct.
func (r *TTdc) Fields() []string {
	return TTdcFields
}

//GetTableName return the associated db table name.
func (r *TTdc) GetTableName() string {
	var viewNamePattern = regexp.MustCompile(`(?i)^t_v_[a-z0-9_]+$`)
	tableName := "t_tdc"
	if viewNamePattern.MatchString(tableName) {
		return tableName[2:]
	}
	return tableName
}

// Create inserts the TTdc to the database.
func (r *TTdc) Create(db Queryer) error {
	err := db.QueryRow(
		`INSERT INTO t_tdc (tdc_id, name, issuer, issue_time, limn, data, expiration, type, goto_view, requested, accepted, remark, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) RETURNING id`,
		&r.TdcID, &r.Name, &r.Issuer, &r.IssueTime, &r.Limn, &r.Data, &r.Expiration, &r.Type, &r.GotoView, &r.Requested, &r.Accepted, &r.Remark, &r.Status).Scan(&r.ID)
	if err != nil {
		return errors.Wrap(err, "failed to insert t_tdc")
	}
	return nil
}

// GetTTdcByPk select the TTdc from the database.
func GetTTdcByPk(db Queryer, pk0 null.Int) (*TTdc, error) {
	var r TTdc
	err := db.QueryRow(
		`SELECT id, tdc_id, name, issuer, issue_time, limn, data, expiration, type, goto_view, requested, accepted, remark, status FROM t_tdc WHERE id = $1`,
		pk0).Scan(&r.ID, &r.TdcID, &r.Name, &r.Issuer, &r.IssueTime, &r.Limn, &r.Data, &r.Expiration, &r.Type, &r.GotoView, &r.Requested, &r.Accepted, &r.Remark, &r.Status)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select t_tdc")
	}
	return &r, nil
}

/*TUndertaker 项目承接表 represents kuser.t_undertaker */
type TUndertaker struct {
	ID            null.Int    `json:"ID,omitempty" db:"id,true,integer"`                                   /* id 承接编号 */
	PrjID         null.Int    `json:"PrjID,omitempty" db:"prj_id,false,bigint"`                            /* prj_id 项目编号 */
	DeveloperID   null.Int    `json:"DeveloperID,omitempty" db:"developer_id,false,bigint"`                /* developer_id 开发者编号 */
	DeveloperType null.String `json:"DeveloperType,omitempty" db:"developer_type,false,character varying"` /* developer_type undertaker,承接者
	invitee,邀请承接
	apply,申请承接 */
	CreateTime null.Int    `json:"CreateTime,omitempty" db:"create_time,false,bigint"`   /* create_time 创建时间 */
	UpdateTime null.Int    `json:"UpdateTime,omitempty" db:"update_time,false,bigint"`   /* update_time 更新时间 */
	Remark     null.String `json:"Remark,omitempty" db:"remark,false,character varying"` /* remark 备注 */
	Status     null.String `json:"Status,omitempty" db:"status,false,character varying"` /* status 状态
	inviting,邀请承接中
	applying,申请承接中
	accepted,已接受
	rejected,被拒绝
	signed, 已签定合同 */
	Filter // build DML where clause
}

//TUndertakerFields full field list for default query
var TUndertakerFields = []string{
	"ID",
	"PrjID",
	"DeveloperID",
	"DeveloperType",
	"CreateTime",
	"UpdateTime",
	"Remark",
	"Status",
}

//Fields return all fields of struct.
func (r *TUndertaker) Fields() []string {
	return TUndertakerFields
}

//GetTableName return the associated db table name.
func (r *TUndertaker) GetTableName() string {
	var viewNamePattern = regexp.MustCompile(`(?i)^t_v_[a-z0-9_]+$`)
	tableName := "t_undertaker"
	if viewNamePattern.MatchString(tableName) {
		return tableName[2:]
	}
	return tableName
}

// Create inserts the TUndertaker to the database.
func (r *TUndertaker) Create(db Queryer) error {
	err := db.QueryRow(
		`INSERT INTO t_undertaker (prj_id, developer_id, developer_type, create_time, update_time, remark, status) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
		&r.PrjID, &r.DeveloperID, &r.DeveloperType, &r.CreateTime, &r.UpdateTime, &r.Remark, &r.Status).Scan(&r.ID)
	if err != nil {
		return errors.Wrap(err, "failed to insert t_undertaker")
	}
	return nil
}

// GetTUndertakerByPk select the TUndertaker from the database.
func GetTUndertakerByPk(db Queryer, pk0 null.Int) (*TUndertaker, error) {
	var r TUndertaker
	err := db.QueryRow(
		`SELECT id, prj_id, developer_id, developer_type, create_time, update_time, remark, status FROM t_undertaker WHERE id = $1`,
		pk0).Scan(&r.ID, &r.PrjID, &r.DeveloperID, &r.DeveloperType, &r.CreateTime, &r.UpdateTime, &r.Remark, &r.Status)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select t_undertaker")
	}
	return &r, nil
}

/*TUser t_user represents kuser.t_user */
type TUser struct {
	ID             null.Int    `json:"ID,omitempty" db:"id,true,integer"`                                      /* id 用户内部编号 */
	ExternalIDType null.String `json:"ExternalIDType,omitempty" db:"external_id_type,false,character varying"` /* external_id_type 用户外部标识类型 */
	ExternalID     null.String `json:"ExternalID,omitempty" db:"external_id,false,character varying"`          /* external_id 用户外部标识 */
	Category       string      `json:"Category,omitempty" db:"category,false,character varying"`               /* category 用户分类 */
	Type           null.String `json:"Type,omitempty" db:"type,false,character varying"`                       /* type 用户类型,
	00:匿名用户, 0000-0001，未提供外部可识别标识用户，未付费，不可识别与联系
	02:注册用户, 0000-0010，具备可识别信息
	04:试用用户, 0000-0100，帐号有过期时间，使用了付费功能，具备可识别信息
	08:机构上帝, 0000-1000，帐号有特定管理功能，具备可识别信息

	10:测试用户, 0001-0000，用来测试的用户，具备可识别信息
	80:系统上帝, 1000-0000，系统管理员功能，具备可识别信息
	　　　　  */
	Language        null.String    `json:"Language,omitempty" db:"language,false,character varying"`          /* language 用户喜好语言 */
	Country         null.String    `json:"Country,omitempty" db:"country,false,character varying"`            /* country 国家 */
	Province        null.String    `json:"Province,omitempty" db:"province,false,character varying"`          /* province 省份 */
	City            null.String    `json:"City,omitempty" db:"city,false,character varying"`                  /* city 城市 */
	Addr            null.String    `json:"Addr,omitempty" db:"addr,false,character varying"`                  /* addr 详细地址 */
	OfficialName    null.String    `json:"OfficialName,omitempty" db:"official_name,false,character varying"` /* official_name 姓名 */
	IDCardType      null.String    `json:"IDCardType,omitempty" db:"id_card_type,false,character varying"`    /* id_card_type 证件类型 */
	IDCardNo        null.String    `json:"IDCardNo,omitempty" db:"id_card_no,false,character varying"`        /* id_card_no 身份证号码 */
	MobilePhone     null.String    `json:"MobilePhone,omitempty" db:"mobile_phone,false,character varying"`   /* mobile_phone 手机号码 */
	Email           null.String    `json:"Email,omitempty" db:"email,false,character varying"`                /* email 电子邮件 */
	Account         null.String    `json:"Account,omitempty" db:"account,false,character varying"`            /* account 登录账号 */
	Gender          null.String    `json:"Gender,omitempty" db:"gender,false,character varying"`              /* gender 性别 */
	Birthday        null.Int       `json:"Birthday,omitempty" db:"birthday,false,bigint"`                     /* birthday 出生日期 */
	Nickname        null.String    `json:"Nickname,omitempty" db:"nickname,false,character varying"`          /* nickname 呢称 */
	Avatar          []byte         `json:"Avatar,omitempty" db:"avatar,false,bytea"`                          /* avatar 头像 */
	AvatarType      null.String    `json:"AvatarType,omitempty" db:"avatar_type,false,character varying"`     /* avatar_type 头像类型, LINK: URL链接，B64: BASE64编码图片 */
	DevID           null.String    `json:"DevID,omitempty" db:"dev_id,false,character varying"`               /* dev_id 终端设备标识,iOS/Android DeviceId */
	DevUserID       null.String    `json:"DevUserID,omitempty" db:"dev_user_id,false,character varying"`      /* dev_user_id 终端用户标识,google Account, iTunes Account */
	DevAccount      null.String    `json:"DevAccount,omitempty" db:"dev_account,false,character varying"`     /* dev_account 与设备关联的用于C2DM/APNS 的Android/iOS帐号 */
	Cert            null.String    `json:"Cert,omitempty" db:"cert,false,character varying"`                  /* cert 证书 */
	UserToken       null.String    `json:"UserToken,omitempty" db:"user_token,false,character varying"`       /* user_token crypt('pwd',gen_salt('bf')) */
	Role            null.Int       `json:"Role,omitempty" db:"role,false,bigint"`                             /* role 最近用户使用角色编号 */
	Grp             null.Int       `json:"Grp,omitempty" db:"grp,false,bigint"`                               /* grp 最近用户使用组编号 */
	IP              null.String    `json:"IP,omitempty" db:"ip,false,character varying"`                      /* ip 最近IP */
	Port            null.String    `json:"Port,omitempty" db:"port,false,character varying"`                  /* port 最近端口 */
	AuthFailedCount null.Int       `json:"AuthFailedCount,omitempty" db:"auth_failed_count,false,integer"`    /* auth_failed_count 登录失败次数 */
	LockDuration    null.Int       `json:"LockDuration,omitempty" db:"lock_duration,false,integer"`           /* lock_duration 需要锁定时长，以秒计 */
	VisitCount      null.Int       `json:"VisitCount,omitempty" db:"visit_count,false,integer"`               /* visit_count 访问计数 */
	AttackCount     null.Int       `json:"AttackCount,omitempty" db:"attack_count,false,integer"`             /* attack_count 攻击次数 */
	LockReason      null.String    `json:"LockReason,omitempty" db:"lock_reason,false,character varying"`     /* lock_reason 锁定原因 */
	LogonTime       null.Int       `json:"LogonTime,omitempty" db:"logon_time,false,bigint"`                  /* logon_time 最近登录时间 */
	BeginLockTime   null.Int       `json:"BeginLockTime,omitempty" db:"begin_lock_time,false,bigint"`         /* begin_lock_time 开始锁定时间 */
	Creator         null.Int       `json:"Creator,omitempty" db:"creator,false,bigint"`                       /* creator 创建者 */
	CreateTime      null.Int       `json:"CreateTime,omitempty" db:"create_time,false,bigint"`                /* create_time 创建时间 */
	Regenerator     null.Int       `json:"Regenerator,omitempty" db:"regenerator,false,bigint"`               /* regenerator 更新者 */
	UpdateTime      null.Int       `json:"UpdateTime,omitempty" db:"update_time,false,bigint"`                /* update_time 帐号信息更新时间 */
	DomainID        null.Int       `json:"DomainID,omitempty" db:"domain_id,false,bigint"`                    /* domain_id 数据隶属 */
	Addi            types.JSONText `json:"Addi,omitempty" db:"addi,false,jsonb"`                              /* addi 用户定制数据 */
	Remark          null.String    `json:"Remark,omitempty" db:"remark,false,character varying"`              /* remark 备注 */
	Status          null.String    `json:"Status,omitempty" db:"status,false,character varying"`              /* status 状态,00: 有效, 02: 禁止登录, 04: 锁定, 06: 攻击者, 08: 过期 */
	Filter                         // build DML where clause
}

//TUserFields full field list for default query
var TUserFields = []string{
	"ID",
	"ExternalIDType",
	"ExternalID",
	"Category",
	"Type",
	"Language",
	"Country",
	"Province",
	"City",
	"Addr",
	"OfficialName",
	"IDCardType",
	"IDCardNo",
	"MobilePhone",
	"Email",
	"Account",
	"Gender",
	"Birthday",
	"Nickname",
	"Avatar",
	"AvatarType",
	"DevID",
	"DevUserID",
	"DevAccount",
	"Cert",
	"UserToken",
	"Role",
	"Grp",
	"IP",
	"Port",
	"AuthFailedCount",
	"LockDuration",
	"VisitCount",
	"AttackCount",
	"LockReason",
	"LogonTime",
	"BeginLockTime",
	"Creator",
	"CreateTime",
	"Regenerator",
	"UpdateTime",
	"DomainID",
	"Addi",
	"Remark",
	"Status",
}

//Fields return all fields of struct.
func (r *TUser) Fields() []string {
	return TUserFields
}

//GetTableName return the associated db table name.
func (r *TUser) GetTableName() string {
	var viewNamePattern = regexp.MustCompile(`(?i)^t_v_[a-z0-9_]+$`)
	tableName := "t_user"
	if viewNamePattern.MatchString(tableName) {
		return tableName[2:]
	}
	return tableName
}

// Create inserts the TUser to the database.
func (r *TUser) Create(db Queryer) error {
	err := db.QueryRow(
		`INSERT INTO t_user (external_id_type, external_id, category, type, language, country, province, city, addr, official_name, id_card_type, id_card_no, mobile_phone, email, account, gender, birthday, nickname, avatar, avatar_type, dev_id, dev_user_id, dev_account, cert, user_token, role, grp, ip, port, auth_failed_count, lock_duration, visit_count, attack_count, lock_reason, logon_time, begin_lock_time, creator, create_time, regenerator, update_time, domain_id, addi, remark, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36, $37, $38, $39, $40, $41, $42, $43, $44) RETURNING id`,
		&r.ExternalIDType, &r.ExternalID, &r.Category, &r.Type, &r.Language, &r.Country, &r.Province, &r.City, &r.Addr, &r.OfficialName, &r.IDCardType, &r.IDCardNo, &r.MobilePhone, &r.Email, &r.Account, &r.Gender, &r.Birthday, &r.Nickname, &r.Avatar, &r.AvatarType, &r.DevID, &r.DevUserID, &r.DevAccount, &r.Cert, &r.UserToken, &r.Role, &r.Grp, &r.IP, &r.Port, &r.AuthFailedCount, &r.LockDuration, &r.VisitCount, &r.AttackCount, &r.LockReason, &r.LogonTime, &r.BeginLockTime, &r.Creator, &r.CreateTime, &r.Regenerator, &r.UpdateTime, &r.DomainID, &r.Addi, &r.Remark, &r.Status).Scan(&r.ID)
	if err != nil {
		return errors.Wrap(err, "failed to insert t_user")
	}
	return nil
}

// GetTUserByPk select the TUser from the database.
func GetTUserByPk(db Queryer, pk0 null.Int) (*TUser, error) {
	var r TUser
	err := db.QueryRow(
		`SELECT id, external_id_type, external_id, category, type, language, country, province, city, addr, official_name, id_card_type, id_card_no, mobile_phone, email, account, gender, birthday, nickname, avatar, avatar_type, dev_id, dev_user_id, dev_account, cert, user_token, role, grp, ip, port, auth_failed_count, lock_duration, visit_count, attack_count, lock_reason, logon_time, begin_lock_time, creator, create_time, regenerator, update_time, domain_id, addi, remark, status FROM t_user WHERE id = $1`,
		pk0).Scan(&r.ID, &r.ExternalIDType, &r.ExternalID, &r.Category, &r.Type, &r.Language, &r.Country, &r.Province, &r.City, &r.Addr, &r.OfficialName, &r.IDCardType, &r.IDCardNo, &r.MobilePhone, &r.Email, &r.Account, &r.Gender, &r.Birthday, &r.Nickname, &r.Avatar, &r.AvatarType, &r.DevID, &r.DevUserID, &r.DevAccount, &r.Cert, &r.UserToken, &r.Role, &r.Grp, &r.IP, &r.Port, &r.AuthFailedCount, &r.LockDuration, &r.VisitCount, &r.AttackCount, &r.LockReason, &r.LogonTime, &r.BeginLockTime, &r.Creator, &r.CreateTime, &r.Regenerator, &r.UpdateTime, &r.DomainID, &r.Addi, &r.Remark, &r.Status)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select t_user")
	}
	return &r, nil
}

/*TUserDegree 用户等级表 represents kuser.t_user_degree */
type TUserDegree struct {
	ID       null.Int `json:"ID,omitempty" db:"id,true,integer"`              /* id 用户能力等级编号 */
	UserID   null.Int `json:"UserID,omitempty" db:"user_id,false,bigint"`     /* user_id 用户编号 */
	DegreeID null.Int `json:"DegreeID,omitempty" db:"degree_id,false,bigint"` /* degree_id 能力等级编号 */
	Filter            // build DML where clause
}

//TUserDegreeFields full field list for default query
var TUserDegreeFields = []string{
	"ID",
	"UserID",
	"DegreeID",
}

//Fields return all fields of struct.
func (r *TUserDegree) Fields() []string {
	return TUserDegreeFields
}

//GetTableName return the associated db table name.
func (r *TUserDegree) GetTableName() string {
	var viewNamePattern = regexp.MustCompile(`(?i)^t_v_[a-z0-9_]+$`)
	tableName := "t_user_degree"
	if viewNamePattern.MatchString(tableName) {
		return tableName[2:]
	}
	return tableName
}

// Create inserts the TUserDegree to the database.
func (r *TUserDegree) Create(db Queryer) error {
	err := db.QueryRow(
		`INSERT INTO t_user_degree (user_id, degree_id) VALUES ($1, $2) RETURNING id`,
		&r.UserID, &r.DegreeID).Scan(&r.ID)
	if err != nil {
		return errors.Wrap(err, "failed to insert t_user_degree")
	}
	return nil
}

// GetTUserDegreeByPk select the TUserDegree from the database.
func GetTUserDegreeByPk(db Queryer, pk0 null.Int) (*TUserDegree, error) {
	var r TUserDegree
	err := db.QueryRow(
		`SELECT id, user_id, degree_id FROM t_user_degree WHERE id = $1`,
		pk0).Scan(&r.ID, &r.UserID, &r.DegreeID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select t_user_degree")
	}
	return &r, nil
}

/*TUserDomain 用户，组，角色配置表 represents kuser.t_user_domain */
type TUserDomain struct {
	ID          null.Int    `json:"ID,omitempty" db:"id,true,integer"`                               /* id 编码 */
	SysUser     null.Int    `json:"SysUser,omitempty" db:"sys_user,false,bigint"`                    /* sys_user 系统用户编码 */
	IDOnDomain  null.String `json:"IDOnDomain,omitempty" db:"id_on_domain,false,character varying"`  /* id_on_domain 基于用户域的用户编码，如广州大学员工号，后勤部员工号，采购组采购员编号，保卫科保安员工号 */
	Domain      null.Int    `json:"Domain,omitempty" db:"domain,false,bigint"`                       /* domain 域, 格式: 机构[部门.科室.组[^角色[!userID]]]，[option]表示可选 */
	GrantSource null.String `json:"GrantSource,omitempty" db:"grant_source,false,character varying"` /* grant_source grant:数据权限由t_relation中left_type:t_domain.id与left_type:t_user.id获得的数据决定,或data_scope中数据决定，但data_scope与t_relation只能存在一种，如果data_scope有效，则忽略t_relation;

	cousin:忽略data_scope与t_relation, 授权数据由被过虑数据的domain_id决定,即被过虑数据的domain_id 与登录用户的t_user.domain_id相同或级别更低的数据，例如
	    用户的t_user.domain为xkb^admin而数据的domain为xkb.school^admin，则用户可以获得该数据

	self: 被过虑数据的creator 与登录用户的t_user.id相同

	api: 由功能(api)自己决定  */
	DataAccessMode null.String    `json:"DataAccessMode,omitempty" db:"data_access_mode,false,character varying"` /* data_access_mode 数据访问类型, full:可读写, read: 只读, write: 写, partial: 部分写/混合 */
	DataScope      types.JSONText `json:"DataScope,omitempty" db:"data_scope,false,jsonb"`                        /* data_scope 当grant_source是grant时,以json数据方式提供数据授权范围格式为:
	  {"granter":"t_user.id","grantee":"t_school.id","data":[1234,456,789]}
	granter: 代表数据拥有者, t_user.id代表用户, t_domain.id代表角色,t_api.id代表功能
	grantee: 代表拥有的数据,t_school.id代表可以访问的机构列表。
	授权数据如果存储在t_relation中则各项分别对应如下
	granter对应left_type, left_key对应t_user_domain.sys_user或t_domain_api.domain
	grantee对应right_type, right_key对应right_type的意义 */
	DomainID    null.Int       `json:"DomainID,omitempty" db:"domain_id,false,bigint"`       /* domain_id 数据隶属 */
	Creator     null.Int       `json:"Creator,omitempty" db:"creator,false,bigint"`          /* creator 本数据创建者 */
	CreateTime  null.Int       `json:"CreateTime,omitempty" db:"create_time,false,bigint"`   /* create_time 生成时间 */
	UpdateTime  null.Int       `json:"UpdateTime,omitempty" db:"update_time,false,bigint"`   /* update_time 帐号信息更新时间 */
	Regenerator null.Int       `json:"Regenerator,omitempty" db:"regenerator,false,bigint"`  /* regenerator 更新者 */
	Addi        types.JSONText `json:"Addi,omitempty" db:"addi,false,jsonb"`                 /* addi 附加信息 */
	Remark      null.String    `json:"Remark,omitempty" db:"remark,false,character varying"` /* remark 备注 */
	Status      null.String    `json:"Status,omitempty" db:"status,false,character varying"` /* status 状态，00：草稿，01：有效，02：作废 */
	Filter                     // build DML where clause
}

//TUserDomainFields full field list for default query
var TUserDomainFields = []string{
	"ID",
	"SysUser",
	"IDOnDomain",
	"Domain",
	"GrantSource",
	"DataAccessMode",
	"DataScope",
	"DomainID",
	"Creator",
	"CreateTime",
	"UpdateTime",
	"Regenerator",
	"Addi",
	"Remark",
	"Status",
}

//Fields return all fields of struct.
func (r *TUserDomain) Fields() []string {
	return TUserDomainFields
}

//GetTableName return the associated db table name.
func (r *TUserDomain) GetTableName() string {
	var viewNamePattern = regexp.MustCompile(`(?i)^t_v_[a-z0-9_]+$`)
	tableName := "t_user_domain"
	if viewNamePattern.MatchString(tableName) {
		return tableName[2:]
	}
	return tableName
}

// Create inserts the TUserDomain to the database.
func (r *TUserDomain) Create(db Queryer) error {
	err := db.QueryRow(
		`INSERT INTO t_user_domain (sys_user, id_on_domain, domain, grant_source, data_access_mode, data_scope, domain_id, creator, create_time, update_time, regenerator, addi, remark, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14) RETURNING id`,
		&r.SysUser, &r.IDOnDomain, &r.Domain, &r.GrantSource, &r.DataAccessMode, &r.DataScope, &r.DomainID, &r.Creator, &r.CreateTime, &r.UpdateTime, &r.Regenerator, &r.Addi, &r.Remark, &r.Status).Scan(&r.ID)
	if err != nil {
		return errors.Wrap(err, "failed to insert t_user_domain")
	}
	return nil
}

// GetTUserDomainByPk select the TUserDomain from the database.
func GetTUserDomainByPk(db Queryer, pk0 null.Int) (*TUserDomain, error) {
	var r TUserDomain
	err := db.QueryRow(
		`SELECT id, sys_user, id_on_domain, domain, grant_source, data_access_mode, data_scope, domain_id, creator, create_time, update_time, regenerator, addi, remark, status FROM t_user_domain WHERE id = $1`,
		pk0).Scan(&r.ID, &r.SysUser, &r.IDOnDomain, &r.Domain, &r.GrantSource, &r.DataAccessMode, &r.DataScope, &r.DomainID, &r.Creator, &r.CreateTime, &r.UpdateTime, &r.Regenerator, &r.Addi, &r.Remark, &r.Status)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select t_user_domain")
	}
	return &r, nil
}

/*TVAa t_v_aa represents kuser.t_v_aa */
type TVAa struct {
	DomainAPIID      null.Int    `json:"DomainAPIID,omitempty" db:"domain_api_id,false,integer"`                      /* domain_api_id domain_api_id */
	DaGrantData      null.String `json:"DaGrantData,omitempty" db:"da_grant_data,false,text"`                         /* da_grant_data da_grant_data */
	DaGrantType      null.String `json:"DaGrantType,omitempty" db:"da_grant_type,false,text"`                         /* da_grant_type da_grant_type */
	DaGrantSource    null.String `json:"DaGrantSource,omitempty" db:"da_grant_source,false,character varying"`        /* da_grant_source da_grant_source */
	DaDataAccessMode null.String `json:"DaDataAccessMode,omitempty" db:"da_data_access_mode,false,character varying"` /* da_data_access_mode da_data_access_mode */
	UserDomainID     null.Int    `json:"UserDomainID,omitempty" db:"user_domain_id,false,integer"`                    /* user_domain_id user_domain_id */
	UdGrantData      null.String `json:"UdGrantData,omitempty" db:"ud_grant_data,false,text"`                         /* ud_grant_data ud_grant_data */
	UdGrantType      null.String `json:"UdGrantType,omitempty" db:"ud_grant_type,false,text"`                         /* ud_grant_type ud_grant_type */
	UdGrantSource    null.String `json:"UdGrantSource,omitempty" db:"ud_grant_source,false,character varying"`        /* ud_grant_source ud_grant_source */
	UdDataAccessMode null.String `json:"UdDataAccessMode,omitempty" db:"ud_data_access_mode,false,character varying"` /* ud_data_access_mode ud_data_access_mode */
	UserID           null.Int    `json:"UserID,omitempty" db:"user_id,false,integer"`                                 /* user_id user_id */
	UserName         null.String `json:"UserName,omitempty" db:"user_name,false,character varying"`                   /* user_name user_name */
	MobilePhone      null.String `json:"MobilePhone,omitempty" db:"mobile_phone,false,character varying"`             /* mobile_phone mobile_phone */
	DomainName       null.String `json:"DomainName,omitempty" db:"domain_name,false,character varying"`               /* domain_name domain_name */
	DomainID         null.Int    `json:"DomainID,omitempty" db:"domain_id,false,integer"`                             /* domain_id domain_id */
	Domain           null.String `json:"Domain,omitempty" db:"domain,false,character varying"`                        /* domain domain */
	Priority         null.Int    `json:"Priority,omitempty" db:"priority,false,smallint"`                             /* priority priority */
	APIID            null.Int    `json:"APIID,omitempty" db:"api_id,false,integer"`                                   /* api_id api_id */
	APIName          null.String `json:"APIName,omitempty" db:"api_name,false,character varying"`                     /* api_name api_name */
	API              null.String `json:"API,omitempty" db:"api,false,character varying"`                              /* api api */
	Filter                       // build DML where clause
}

//TVAaFields full field list for default query
var TVAaFields = []string{
	"DomainAPIID",
	"DaGrantData",
	"DaGrantType",
	"DaGrantSource",
	"DaDataAccessMode",
	"UserDomainID",
	"UdGrantData",
	"UdGrantType",
	"UdGrantSource",
	"UdDataAccessMode",
	"UserID",
	"UserName",
	"MobilePhone",
	"DomainName",
	"DomainID",
	"Domain",
	"Priority",
	"APIID",
	"APIName",
	"API",
}

//Fields return all fields of struct.
func (r *TVAa) Fields() []string {
	return TVAaFields
}

//GetTableName return the associated db table name.
func (r *TVAa) GetTableName() string {
	var viewNamePattern = regexp.MustCompile(`(?i)^t_v_[a-z0-9_]+$`)
	tableName := "t_v_aa"
	if viewNamePattern.MatchString(tableName) {
		return tableName[2:]
	}
	return tableName
}

// Create inserts the TVAa to the database.
func (r *TVAa) Create(db Queryer) error {
	_, err := db.Exec(
		`INSERT INTO t_v_aa (domain_api_id, da_grant_data, da_grant_type, da_grant_source, da_data_access_mode, user_domain_id, ud_grant_data, ud_grant_type, ud_grant_source, ud_data_access_mode, user_id, user_name, mobile_phone, domain_name, domain_id, domain, priority, api_id, api_name, api) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)`,
		&r.DomainAPIID, &r.DaGrantData, &r.DaGrantType, &r.DaGrantSource, &r.DaDataAccessMode, &r.UserDomainID, &r.UdGrantData, &r.UdGrantType, &r.UdGrantSource, &r.UdDataAccessMode, &r.UserID, &r.UserName, &r.MobilePhone, &r.DomainName, &r.DomainID, &r.Domain, &r.Priority, &r.APIID, &r.APIName, &r.API)
	if err != nil {
		return errors.Wrap(err, "failed to insert t_v_aa")
	}
	return nil
}

// GetTVAaByPk select the TVAa from the database.
func GetTVAaByPk(db Queryer) (*TVAa, error) {
	var r TVAa
	err := db.QueryRow(
		`SELECT domain_api_id, da_grant_data, da_grant_type, da_grant_source, da_data_access_mode, user_domain_id, ud_grant_data, ud_grant_type, ud_grant_source, ud_data_access_mode, user_id, user_name, mobile_phone, domain_name, domain_id, domain, priority, api_id, api_name, api FROM t_v_aa WHERE `,
	).Scan(&r.DomainAPIID, &r.DaGrantData, &r.DaGrantType, &r.DaGrantSource, &r.DaDataAccessMode, &r.UserDomainID, &r.UdGrantData, &r.UdGrantType, &r.UdGrantSource, &r.UdDataAccessMode, &r.UserID, &r.UserName, &r.MobilePhone, &r.DomainName, &r.DomainID, &r.Domain, &r.Priority, &r.APIID, &r.APIName, &r.API)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select t_v_aa")
	}
	return &r, nil
}

/*TVAPIDomain t_v_api_domain represents kuser.t_v_api_domain */
type TVAPIDomain struct {
	ID             null.Int       `json:"ID,omitempty" db:"id,false,integer"`                                     /* id id */
	APIID          null.Int       `json:"APIID,omitempty" db:"api_id,false,integer"`                              /* api_id api_id */
	APIName        null.String    `json:"APIName,omitempty" db:"api_name,false,character varying"`                /* api_name api_name */
	ExposePath     null.String    `json:"ExposePath,omitempty" db:"expose_path,false,character varying"`          /* expose_path expose_path */
	DomainID       null.Int       `json:"DomainID,omitempty" db:"domain_id,false,integer"`                        /* domain_id domain_id */
	DomainName     null.String    `json:"DomainName,omitempty" db:"domain_name,false,character varying"`          /* domain_name domain_name */
	Priority       null.Int       `json:"Priority,omitempty" db:"priority,false,smallint"`                        /* priority priority */
	GrantSource    null.String    `json:"GrantSource,omitempty" db:"grant_source,false,character varying"`        /* grant_source grant_source */
	DataAccessMode null.String    `json:"DataAccessMode,omitempty" db:"data_access_mode,false,character varying"` /* data_access_mode data_access_mode */
	DataScope      types.JSONText `json:"DataScope,omitempty" db:"data_scope,false,jsonb"`                        /* data_scope data_scope */
	Filter                        // build DML where clause
}

//TVAPIDomainFields full field list for default query
var TVAPIDomainFields = []string{
	"ID",
	"APIID",
	"APIName",
	"ExposePath",
	"DomainID",
	"DomainName",
	"Priority",
	"GrantSource",
	"DataAccessMode",
	"DataScope",
}

//Fields return all fields of struct.
func (r *TVAPIDomain) Fields() []string {
	return TVAPIDomainFields
}

//GetTableName return the associated db table name.
func (r *TVAPIDomain) GetTableName() string {
	var viewNamePattern = regexp.MustCompile(`(?i)^t_v_[a-z0-9_]+$`)
	tableName := "t_v_api_domain"
	if viewNamePattern.MatchString(tableName) {
		return tableName[2:]
	}
	return tableName
}

// Create inserts the TVAPIDomain to the database.
func (r *TVAPIDomain) Create(db Queryer) error {
	_, err := db.Exec(
		`INSERT INTO t_v_api_domain (id, api_id, api_name, expose_path, domain_id, domain_name, priority, grant_source, data_access_mode, data_scope) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		&r.ID, &r.APIID, &r.APIName, &r.ExposePath, &r.DomainID, &r.DomainName, &r.Priority, &r.GrantSource, &r.DataAccessMode, &r.DataScope)
	if err != nil {
		return errors.Wrap(err, "failed to insert t_v_api_domain")
	}
	return nil
}

// GetTVAPIDomainByPk select the TVAPIDomain from the database.
func GetTVAPIDomainByPk(db Queryer) (*TVAPIDomain, error) {
	var r TVAPIDomain
	err := db.QueryRow(
		`SELECT id, api_id, api_name, expose_path, domain_id, domain_name, priority, grant_source, data_access_mode, data_scope FROM t_v_api_domain WHERE `,
	).Scan(&r.ID, &r.APIID, &r.APIName, &r.ExposePath, &r.DomainID, &r.DomainName, &r.Priority, &r.GrantSource, &r.DataAccessMode, &r.DataScope)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select t_v_api_domain")
	}
	return &r, nil
}

/*TVAuthenticate t_v_authenticate represents kuser.t_v_authenticate */
type TVAuthenticate struct {
	GrantID        null.Int    `json:"GrantID,omitempty" db:"grant_id,false,integer"`                          /* grant_id grant_id */
	GrantData      null.String `json:"GrantData,omitempty" db:"grant_data,false,text"`                         /* grant_data grant_data */
	GrantType      null.String `json:"GrantType,omitempty" db:"grant_type,false,text"`                         /* grant_type grant_type */
	GrantSource    null.String `json:"GrantSource,omitempty" db:"grant_source,false,character varying"`        /* grant_source grant_source */
	DataAccessMode null.String `json:"DataAccessMode,omitempty" db:"data_access_mode,false,character varying"` /* data_access_mode data_access_mode */
	UserID         null.Int    `json:"UserID,omitempty" db:"user_id,false,integer"`                            /* user_id user_id */
	UserName       null.String `json:"UserName,omitempty" db:"user_name,false,character varying"`              /* user_name user_name */
	MobilePhone    null.String `json:"MobilePhone,omitempty" db:"mobile_phone,false,character varying"`        /* mobile_phone mobile_phone */
	DomainName     null.String `json:"DomainName,omitempty" db:"domain_name,false,character varying"`          /* domain_name domain_name */
	DomainID       null.Int    `json:"DomainID,omitempty" db:"domain_id,false,integer"`                        /* domain_id domain_id */
	Domain         null.String `json:"Domain,omitempty" db:"domain,false,character varying"`                   /* domain domain */
	Priority       null.Int    `json:"Priority,omitempty" db:"priority,false,smallint"`                        /* priority priority */
	APIID          null.Int    `json:"APIID,omitempty" db:"api_id,false,integer"`                              /* api_id api_id */
	APIName        null.String `json:"APIName,omitempty" db:"api_name,false,character varying"`                /* api_name api_name */
	API            null.String `json:"API,omitempty" db:"api,false,character varying"`                         /* api api */
	Filter                     // build DML where clause
}

//TVAuthenticateFields full field list for default query
var TVAuthenticateFields = []string{
	"GrantID",
	"GrantData",
	"GrantType",
	"GrantSource",
	"DataAccessMode",
	"UserID",
	"UserName",
	"MobilePhone",
	"DomainName",
	"DomainID",
	"Domain",
	"Priority",
	"APIID",
	"APIName",
	"API",
}

//Fields return all fields of struct.
func (r *TVAuthenticate) Fields() []string {
	return TVAuthenticateFields
}

//GetTableName return the associated db table name.
func (r *TVAuthenticate) GetTableName() string {
	var viewNamePattern = regexp.MustCompile(`(?i)^t_v_[a-z0-9_]+$`)
	tableName := "t_v_authenticate"
	if viewNamePattern.MatchString(tableName) {
		return tableName[2:]
	}
	return tableName
}

// Create inserts the TVAuthenticate to the database.
func (r *TVAuthenticate) Create(db Queryer) error {
	_, err := db.Exec(
		`INSERT INTO t_v_authenticate (grant_id, grant_data, grant_type, grant_source, data_access_mode, user_id, user_name, mobile_phone, domain_name, domain_id, domain, priority, api_id, api_name, api) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`,
		&r.GrantID, &r.GrantData, &r.GrantType, &r.GrantSource, &r.DataAccessMode, &r.UserID, &r.UserName, &r.MobilePhone, &r.DomainName, &r.DomainID, &r.Domain, &r.Priority, &r.APIID, &r.APIName, &r.API)
	if err != nil {
		return errors.Wrap(err, "failed to insert t_v_authenticate")
	}
	return nil
}

// GetTVAuthenticateByPk select the TVAuthenticate from the database.
func GetTVAuthenticateByPk(db Queryer) (*TVAuthenticate, error) {
	var r TVAuthenticate
	err := db.QueryRow(
		`SELECT grant_id, grant_data, grant_type, grant_source, data_access_mode, user_id, user_name, mobile_phone, domain_name, domain_id, domain, priority, api_id, api_name, api FROM t_v_authenticate WHERE `,
	).Scan(&r.GrantID, &r.GrantData, &r.GrantType, &r.GrantSource, &r.DataAccessMode, &r.UserID, &r.UserName, &r.MobilePhone, &r.DomainName, &r.DomainID, &r.Domain, &r.Priority, &r.APIID, &r.APIName, &r.API)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select t_v_authenticate")
	}
	return &r, nil
}

/*TVDomainAPI t_v_domain_api represents kuser.t_v_domain_api */
type TVDomainAPI struct {
	ID                 null.Int       `json:"ID,omitempty" db:"id,false,integer"`                                             /* id id */
	AuthDomainID       null.Int       `json:"AuthDomainID,omitempty" db:"auth_domain_id,false,integer"`                       /* auth_domain_id auth_domain_id */
	DomainName         null.String    `json:"DomainName,omitempty" db:"domain_name,false,character varying"`                  /* domain_name domain_name */
	Domain             null.String    `json:"Domain,omitempty" db:"domain,false,character varying"`                           /* domain domain */
	Priority           null.Int       `json:"Priority,omitempty" db:"priority,false,smallint"`                                /* priority priority */
	APIID              null.Int       `json:"APIID,omitempty" db:"api_id,false,integer"`                                      /* api_id api_id */
	APIName            null.String    `json:"APIName,omitempty" db:"api_name,false,character varying"`                        /* api_name api_name */
	ExposePath         null.String    `json:"ExposePath,omitempty" db:"expose_path,false,character varying"`                  /* expose_path expose_path */
	AccessControlLevel null.String    `json:"AccessControlLevel,omitempty" db:"access_control_level,false,character varying"` /* access_control_level access_control_level */
	DomainID           null.Int       `json:"DomainID,omitempty" db:"domain_id,false,bigint"`                                 /* domain_id domain_id */
	GrantSource        null.String    `json:"GrantSource,omitempty" db:"grant_source,false,character varying"`                /* grant_source grant_source */
	DataAccessMode     null.String    `json:"DataAccessMode,omitempty" db:"data_access_mode,false,character varying"`         /* data_access_mode data_access_mode */
	DataScope          types.JSONText `json:"DataScope,omitempty" db:"data_scope,false,jsonb"`                                /* data_scope data_scope */
	CreateTime         null.Int       `json:"CreateTime,omitempty" db:"create_time,false,bigint"`                             /* create_time create_time */
	Remark             null.String    `json:"Remark,omitempty" db:"remark,false,character varying"`                           /* remark remark */
	Addi               types.JSONText `json:"Addi,omitempty" db:"addi,false,jsonb"`                                           /* addi addi */
	Creator            null.Int       `json:"Creator,omitempty" db:"creator,false,bigint"`                                    /* creator creator */
	Status             null.String    `json:"Status,omitempty" db:"status,false,character varying"`                           /* status status */
	Filter                            // build DML where clause
}

//TVDomainAPIFields full field list for default query
var TVDomainAPIFields = []string{
	"ID",
	"AuthDomainID",
	"DomainName",
	"Domain",
	"Priority",
	"APIID",
	"APIName",
	"ExposePath",
	"AccessControlLevel",
	"DomainID",
	"GrantSource",
	"DataAccessMode",
	"DataScope",
	"CreateTime",
	"Remark",
	"Addi",
	"Creator",
	"Status",
}

//Fields return all fields of struct.
func (r *TVDomainAPI) Fields() []string {
	return TVDomainAPIFields
}

//GetTableName return the associated db table name.
func (r *TVDomainAPI) GetTableName() string {
	var viewNamePattern = regexp.MustCompile(`(?i)^t_v_[a-z0-9_]+$`)
	tableName := "t_v_domain_api"
	if viewNamePattern.MatchString(tableName) {
		return tableName[2:]
	}
	return tableName
}

// Create inserts the TVDomainAPI to the database.
func (r *TVDomainAPI) Create(db Queryer) error {
	_, err := db.Exec(
		`INSERT INTO t_v_domain_api (id, auth_domain_id, domain_name, domain, priority, api_id, api_name, expose_path, access_control_level, domain_id, grant_source, data_access_mode, data_scope, create_time, remark, addi, creator, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)`,
		&r.ID, &r.AuthDomainID, &r.DomainName, &r.Domain, &r.Priority, &r.APIID, &r.APIName, &r.ExposePath, &r.AccessControlLevel, &r.DomainID, &r.GrantSource, &r.DataAccessMode, &r.DataScope, &r.CreateTime, &r.Remark, &r.Addi, &r.Creator, &r.Status)
	if err != nil {
		return errors.Wrap(err, "failed to insert t_v_domain_api")
	}
	return nil
}

// GetTVDomainAPIByPk select the TVDomainAPI from the database.
func GetTVDomainAPIByPk(db Queryer) (*TVDomainAPI, error) {
	var r TVDomainAPI
	err := db.QueryRow(
		`SELECT id, auth_domain_id, domain_name, domain, priority, api_id, api_name, expose_path, access_control_level, domain_id, grant_source, data_access_mode, data_scope, create_time, remark, addi, creator, status FROM t_v_domain_api WHERE `,
	).Scan(&r.ID, &r.AuthDomainID, &r.DomainName, &r.Domain, &r.Priority, &r.APIID, &r.APIName, &r.ExposePath, &r.AccessControlLevel, &r.DomainID, &r.GrantSource, &r.DataAccessMode, &r.DataScope, &r.CreateTime, &r.Remark, &r.Addi, &r.Creator, &r.Status)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select t_v_domain_api")
	}
	return &r, nil
}

/*TVDomainUser t_v_domain_user represents kuser.t_v_domain_user */
type TVDomainUser struct {
	ID         null.Int    `json:"ID,omitempty" db:"id,false,integer"`                            /* id id */
	DomainID   null.Int    `json:"DomainID,omitempty" db:"domain_id,false,integer"`               /* domain_id domain_id */
	DomainName null.String `json:"DomainName,omitempty" db:"domain_name,false,character varying"` /* domain_name domain_name */
	Priority   null.Int    `json:"Priority,omitempty" db:"priority,false,smallint"`               /* priority priority */
	UserID     null.Int    `json:"UserID,omitempty" db:"user_id,false,integer"`                   /* user_id user_id */
	UserName   null.String `json:"UserName,omitempty" db:"user_name,false,character varying"`     /* user_name user_name */
	Filter                 // build DML where clause
}

//TVDomainUserFields full field list for default query
var TVDomainUserFields = []string{
	"ID",
	"DomainID",
	"DomainName",
	"Priority",
	"UserID",
	"UserName",
}

//Fields return all fields of struct.
func (r *TVDomainUser) Fields() []string {
	return TVDomainUserFields
}

//GetTableName return the associated db table name.
func (r *TVDomainUser) GetTableName() string {
	var viewNamePattern = regexp.MustCompile(`(?i)^t_v_[a-z0-9_]+$`)
	tableName := "t_v_domain_user"
	if viewNamePattern.MatchString(tableName) {
		return tableName[2:]
	}
	return tableName
}

// Create inserts the TVDomainUser to the database.
func (r *TVDomainUser) Create(db Queryer) error {
	_, err := db.Exec(
		`INSERT INTO t_v_domain_user (id, domain_id, domain_name, priority, user_id, user_name) VALUES ($1, $2, $3, $4, $5, $6)`,
		&r.ID, &r.DomainID, &r.DomainName, &r.Priority, &r.UserID, &r.UserName)
	if err != nil {
		return errors.Wrap(err, "failed to insert t_v_domain_user")
	}
	return nil
}

// GetTVDomainUserByPk select the TVDomainUser from the database.
func GetTVDomainUserByPk(db Queryer) (*TVDomainUser, error) {
	var r TVDomainUser
	err := db.QueryRow(
		`SELECT id, domain_id, domain_name, priority, user_id, user_name FROM t_v_domain_user WHERE `,
	).Scan(&r.ID, &r.DomainID, &r.DomainName, &r.Priority, &r.UserID, &r.UserName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select t_v_domain_user")
	}
	return &r, nil
}

/*TVInsurancePolicy t_v_insurance_policy represents kuser.t_v_insurance_policy */
type TVInsurancePolicy struct {
	ID                    null.Int       `json:"ID,omitempty" db:"id,false,integer"`                                         /* id id */
	OrderID               null.Int       `json:"OrderID,omitempty" db:"order_id,false,bigint"`                               /* order_id order_id */
	Sn                    null.String    `json:"Sn,omitempty" db:"sn,false,character varying"`                               /* sn sn */
	Name                  null.String    `json:"Name,omitempty" db:"name,false,character varying"`                           /* name name */
	Policy                null.String    `json:"Policy,omitempty" db:"policy,false,character varying"`                       /* policy policy */
	Start                 null.Int       `json:"Start,omitempty" db:"start,false,bigint"`                                    /* start start */
	Cease                 null.Int       `json:"Cease,omitempty" db:"cease,false,bigint"`                                    /* cease cease */
	Year                  null.Int       `json:"Year,omitempty" db:"year,false,smallint"`                                    /* year year */
	Duration              null.Int       `json:"Duration,omitempty" db:"duration,false,bigint"`                              /* duration duration */
	Premium               null.Float     `json:"Premium,omitempty" db:"premium,false,double precision"`                      /* premium premium */
	PolicyScheme          types.JSONText `json:"PolicyScheme,omitempty" db:"policy_scheme,false,jsonb"`                      /* policy_scheme policy_scheme */
	CreateTime            null.Int       `json:"CreateTime,omitempty" db:"create_time,false,bigint"`                         /* create_time create_time */
	UpdateTime            null.Int       `json:"UpdateTime,omitempty" db:"update_time,false,bigint"`                         /* update_time update_time */
	SnCreator             null.Int       `json:"SnCreator,omitempty" db:"sn_creator,false,bigint"`                           /* sn_creator sn_creator */
	Creator               null.Int       `json:"Creator,omitempty" db:"creator,false,bigint"`                                /* creator creator */
	Addi                  types.JSONText `json:"Addi,omitempty" db:"addi,false,jsonb"`                                       /* addi addi */
	Remark                null.String    `json:"Remark,omitempty" db:"remark,false,character varying"`                       /* remark remark */
	Status                null.String    `json:"Status,omitempty" db:"status,false,character varying"`                       /* status status */
	InsuranceTypeID       null.Int       `json:"InsuranceTypeID,omitempty" db:"insurance_type_id,false,bigint"`              /* insurance_type_id insurance_type_id */
	IsEntryPolicy         null.Bool      `json:"IsEntryPolicy,omitempty" db:"is_entry_policy,false,boolean"`                 /* is_entry_policy is_entry_policy */
	Favorite              null.Bool      `json:"Favorite,omitempty" db:"favorite,false,boolean"`                             /* favorite favorite */
	TradeNo               null.String    `json:"TradeNo,omitempty" db:"trade_no,false,character varying"`                    /* trade_no trade_no */
	PayOrderNo            null.String    `json:"PayOrderNo,omitempty" db:"pay_order_no,false,character varying"`             /* pay_order_no pay_order_no */
	InsureOrderNo         null.String    `json:"InsureOrderNo,omitempty" db:"insure_order_no,false,character varying"`       /* insure_order_no insure_order_no */
	OrgID                 null.Int       `json:"OrgID,omitempty" db:"org_id,false,bigint"`                                   /* org_id org_id */
	PlanID                null.Int       `json:"PlanID,omitempty" db:"plan_id,false,bigint"`                                 /* plan_id plan_id */
	Batch                 null.String    `json:"Batch,omitempty" db:"batch,false,character varying"`                         /* batch batch */
	InsuredID             null.Int       `json:"InsuredID,omitempty" db:"insured_id,false,bigint"`                           /* insured_id insured_id */
	PolicyholderID        null.Int       `json:"PolicyholderID,omitempty" db:"policyholder_id,false,bigint"`                 /* policyholder_id policyholder_id */
	OrderCreateTime       null.Int       `json:"OrderCreateTime,omitempty" db:"order_create_time,false,bigint"`              /* order_create_time order_create_time */
	PayTime               null.Int       `json:"PayTime,omitempty" db:"pay_time,false,bigint"`                               /* pay_time pay_time */
	PayType               null.String    `json:"PayType,omitempty" db:"pay_type,false,character varying"`                    /* pay_type pay_type */
	Amount                null.Float     `json:"Amount,omitempty" db:"amount,false,double precision"`                        /* amount amount */
	UnitPrice             null.Float     `json:"UnitPrice,omitempty" db:"unit_price,false,double precision"`                 /* unit_price unit_price */
	CommenceDate          null.Int       `json:"CommenceDate,omitempty" db:"commence_date,false,bigint"`                     /* commence_date commence_date */
	ExpiryDate            null.Int       `json:"ExpiryDate,omitempty" db:"expiry_date,false,bigint"`                         /* expiry_date expiry_date */
	Indate                null.Int       `json:"Indate,omitempty" db:"indate,false,bigint"`                                  /* indate indate */
	ChargeMode            null.String    `json:"ChargeMode,omitempty" db:"charge_mode,false,character varying"`              /* charge_mode charge_mode */
	Relation              null.String    `json:"Relation,omitempty" db:"relation,false,character varying"`                   /* relation relation */
	InsuranceType         null.String    `json:"InsuranceType,omitempty" db:"insurance_type,false,character varying"`        /* insurance_type insurance_type */
	PolicyDoc             null.String    `json:"PolicyDoc,omitempty" db:"policy_doc,false,character varying"`                /* policy_doc policy_doc */
	Same                  null.Bool      `json:"Same,omitempty" db:"same,false,boolean"`                                     /* same same */
	OrderStatus           null.String    `json:"OrderStatus,omitempty" db:"order_status,false,character varying"`            /* order_status order_status */
	OrderCreator          null.Int       `json:"OrderCreator,omitempty" db:"order_creator,false,bigint"`                     /* order_creator order_creator */
	Insurer               null.String    `json:"Insurer,omitempty" db:"insurer,false,character varying"`                     /* insurer insurer */
	IOfficialName         null.String    `json:"IOfficialName,omitempty" db:"i_official_name,false,character varying"`       /* i_official_name i_official_name */
	IIDCardType           null.String    `json:"IIDCardType,omitempty" db:"i_id_card_type,false,character varying"`          /* i_id_card_type i_id_card_type */
	IIDCardNo             null.String    `json:"IIDCardNo,omitempty" db:"i_id_card_no,false,character varying"`              /* i_id_card_no i_id_card_no */
	IMobilePhone          null.String    `json:"IMobilePhone,omitempty" db:"i_mobile_phone,false,character varying"`         /* i_mobile_phone i_mobile_phone */
	IGender               null.String    `json:"IGender,omitempty" db:"i_gender,false,character varying"`                    /* i_gender i_gender */
	IBirthday             null.Int       `json:"IBirthday,omitempty" db:"i_birthday,false,bigint"`                           /* i_birthday i_birthday */
	IAddi                 types.JSONText `json:"IAddi,omitempty" db:"i_addi,false,jsonb"`                                    /* i_addi i_addi */
	HOfficialName         null.String    `json:"HOfficialName,omitempty" db:"h_official_name,false,character varying"`       /* h_official_name h_official_name */
	HIDCardType           null.String    `json:"HIDCardType,omitempty" db:"h_id_card_type,false,character varying"`          /* h_id_card_type h_id_card_type */
	HIDCardNo             null.String    `json:"HIDCardNo,omitempty" db:"h_id_card_no,false,character varying"`              /* h_id_card_no h_id_card_no */
	HMobilePhone          null.String    `json:"HMobilePhone,omitempty" db:"h_mobile_phone,false,character varying"`         /* h_mobile_phone h_mobile_phone */
	HAddi                 types.JSONText `json:"HAddi,omitempty" db:"h_addi,false,jsonb"`                                    /* h_addi h_addi */
	Subdistrict           null.String    `json:"Subdistrict,omitempty" db:"subdistrict,false,character varying"`             /* subdistrict subdistrict */
	Faculty               null.String    `json:"Faculty,omitempty" db:"faculty,false,character varying"`                     /* faculty faculty */
	Grade                 null.String    `json:"Grade,omitempty" db:"grade,false,character varying"`                         /* grade grade */
	Class                 null.String    `json:"Class,omitempty" db:"class,false,character varying"`                         /* class class */
	XCreateTime           null.Int       `json:"XCreateTime,omitempty" db:"x_create_time,false,bigint"`                      /* x_create_time x_create_time */
	School                null.String    `json:"School,omitempty" db:"school,false,character varying"`                       /* school school */
	SFaculty              types.JSONText `json:"SFaculty,omitempty" db:"s_faculty,false,jsonb"`                              /* s_faculty s_faculty */
	SBranches             types.JSONText `json:"SBranches,omitempty" db:"s_branches,false,jsonb"`                            /* s_branches s_branches */
	SCategory             null.String    `json:"SCategory,omitempty" db:"s_category,false,character varying"`                /* s_category s_category */
	Province              null.String    `json:"Province,omitempty" db:"province,false,character varying"`                   /* province province */
	City                  null.String    `json:"City,omitempty" db:"city,false,character varying"`                           /* city city */
	District              null.String    `json:"District,omitempty" db:"district,false,character varying"`                   /* district district */
	DataSyncTarget        null.String    `json:"DataSyncTarget,omitempty" db:"data_sync_target,false,character varying"`     /* data_sync_target data_sync_target */
	SaleManagers          types.JSONText `json:"SaleManagers,omitempty" db:"sale_managers,false,jsonb"`                      /* sale_managers sale_managers */
	SchoolManagers        types.JSONText `json:"SchoolManagers,omitempty" db:"school_managers,false,jsonb"`                  /* school_managers school_managers */
	PurchaseRule          types.JSONText `json:"PurchaseRule,omitempty" db:"purchase_rule,false,jsonb"`                      /* purchase_rule purchase_rule */
	SCreateTime           null.Int       `json:"SCreateTime,omitempty" db:"s_create_time,false,bigint"`                      /* s_create_time s_create_time */
	InsuranceTypeParentID null.Int       `json:"InsuranceTypeParentID,omitempty" db:"insurance_type_parent_id,false,bigint"` /* insurance_type_parent_id insurance_type_parent_id */
	InsureAttachID        null.Int       `json:"InsureAttachID,omitempty" db:"insure_attach_id,false,integer"`               /* insure_attach_id insure_attach_id */
	PolicyNo              null.String    `json:"PolicyNo,omitempty" db:"policy_no,false,character varying"`                  /* policy_no policy_no */
	Others                types.JSONText `json:"Others,omitempty" db:"others,false,jsonb"`                                   /* others others */
	Files                 types.JSONText `json:"Files,omitempty" db:"files,false,jsonb"`                                     /* files files */
	AttachAddi            types.JSONText `json:"AttachAddi,omitempty" db:"attach_addi,false,jsonb"`                          /* attach_addi attach_addi */
	AttachCreateTime      null.Int       `json:"AttachCreateTime,omitempty" db:"attach_create_time,false,bigint"`            /* attach_create_time attach_create_time */
	AttachUpdateTime      null.Int       `json:"AttachUpdateTime,omitempty" db:"attach_update_time,false,bigint"`            /* attach_update_time attach_update_time */
	AttachCreator         null.Int       `json:"AttachCreator,omitempty" db:"attach_creator,false,bigint"`                   /* attach_creator attach_creator */
	PolicyUploadStatus    null.String    `json:"PolicyUploadStatus,omitempty" db:"policy_upload_status,false,text"`          /* policy_upload_status policy_upload_status */
	InvoiceUploadStatus   null.String    `json:"InvoiceUploadStatus,omitempty" db:"invoice_upload_status,false,text"`        /* invoice_upload_status invoice_upload_status */
	Filter                               // build DML where clause
}

//TVInsurancePolicyFields full field list for default query
var TVInsurancePolicyFields = []string{
	"ID",
	"OrderID",
	"Sn",
	"Name",
	"Policy",
	"Start",
	"Cease",
	"Year",
	"Duration",
	"Premium",
	"PolicyScheme",
	"CreateTime",
	"UpdateTime",
	"SnCreator",
	"Creator",
	"Addi",
	"Remark",
	"Status",
	"InsuranceTypeID",
	"IsEntryPolicy",
	"Favorite",
	"TradeNo",
	"PayOrderNo",
	"InsureOrderNo",
	"OrgID",
	"PlanID",
	"Batch",
	"InsuredID",
	"PolicyholderID",
	"OrderCreateTime",
	"PayTime",
	"PayType",
	"Amount",
	"UnitPrice",
	"CommenceDate",
	"ExpiryDate",
	"Indate",
	"ChargeMode",
	"Relation",
	"InsuranceType",
	"PolicyDoc",
	"Same",
	"OrderStatus",
	"OrderCreator",
	"Insurer",
	"IOfficialName",
	"IIDCardType",
	"IIDCardNo",
	"IMobilePhone",
	"IGender",
	"IBirthday",
	"IAddi",
	"HOfficialName",
	"HIDCardType",
	"HIDCardNo",
	"HMobilePhone",
	"HAddi",
	"Subdistrict",
	"Faculty",
	"Grade",
	"Class",
	"XCreateTime",
	"School",
	"SFaculty",
	"SBranches",
	"SCategory",
	"Province",
	"City",
	"District",
	"DataSyncTarget",
	"SaleManagers",
	"SchoolManagers",
	"PurchaseRule",
	"SCreateTime",
	"InsuranceTypeParentID",
	"InsureAttachID",
	"PolicyNo",
	"Others",
	"Files",
	"AttachAddi",
	"AttachCreateTime",
	"AttachUpdateTime",
	"AttachCreator",
	"PolicyUploadStatus",
	"InvoiceUploadStatus",
}

//Fields return all fields of struct.
func (r *TVInsurancePolicy) Fields() []string {
	return TVInsurancePolicyFields
}

//GetTableName return the associated db table name.
func (r *TVInsurancePolicy) GetTableName() string {
	var viewNamePattern = regexp.MustCompile(`(?i)^t_v_[a-z0-9_]+$`)
	tableName := "t_v_insurance_policy"
	if viewNamePattern.MatchString(tableName) {
		return tableName[2:]
	}
	return tableName
}

// Create inserts the TVInsurancePolicy to the database.
func (r *TVInsurancePolicy) Create(db Queryer) error {
	_, err := db.Exec(
		`INSERT INTO t_v_insurance_policy (id, order_id, sn, name, policy, start, cease, year, duration, premium, policy_scheme, create_time, update_time, sn_creator, creator, addi, remark, status, insurance_type_id, is_entry_policy, favorite, trade_no, pay_order_no, insure_order_no, org_id, plan_id, batch, insured_id, policyholder_id, order_create_time, pay_time, pay_type, amount, unit_price, commence_date, expiry_date, indate, charge_mode, relation, insurance_type, policy_doc, same, order_status, order_creator, insurer, i_official_name, i_id_card_type, i_id_card_no, i_mobile_phone, i_gender, i_birthday, i_addi, h_official_name, h_id_card_type, h_id_card_no, h_mobile_phone, h_addi, subdistrict, faculty, grade, class, x_create_time, school, s_faculty, s_branches, s_category, province, city, district, data_sync_target, sale_managers, school_managers, purchase_rule, s_create_time, insurance_type_parent_id, insure_attach_id, policy_no, others, files, attach_addi, attach_create_time, attach_update_time, attach_creator, policy_upload_status, invoice_upload_status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36, $37, $38, $39, $40, $41, $42, $43, $44, $45, $46, $47, $48, $49, $50, $51, $52, $53, $54, $55, $56, $57, $58, $59, $60, $61, $62, $63, $64, $65, $66, $67, $68, $69, $70, $71, $72, $73, $74, $75, $76, $77, $78, $79, $80, $81, $82, $83, $84, $85)`,
		&r.ID, &r.OrderID, &r.Sn, &r.Name, &r.Policy, &r.Start, &r.Cease, &r.Year, &r.Duration, &r.Premium, &r.PolicyScheme, &r.CreateTime, &r.UpdateTime, &r.SnCreator, &r.Creator, &r.Addi, &r.Remark, &r.Status, &r.InsuranceTypeID, &r.IsEntryPolicy, &r.Favorite, &r.TradeNo, &r.PayOrderNo, &r.InsureOrderNo, &r.OrgID, &r.PlanID, &r.Batch, &r.InsuredID, &r.PolicyholderID, &r.OrderCreateTime, &r.PayTime, &r.PayType, &r.Amount, &r.UnitPrice, &r.CommenceDate, &r.ExpiryDate, &r.Indate, &r.ChargeMode, &r.Relation, &r.InsuranceType, &r.PolicyDoc, &r.Same, &r.OrderStatus, &r.OrderCreator, &r.Insurer, &r.IOfficialName, &r.IIDCardType, &r.IIDCardNo, &r.IMobilePhone, &r.IGender, &r.IBirthday, &r.IAddi, &r.HOfficialName, &r.HIDCardType, &r.HIDCardNo, &r.HMobilePhone, &r.HAddi, &r.Subdistrict, &r.Faculty, &r.Grade, &r.Class, &r.XCreateTime, &r.School, &r.SFaculty, &r.SBranches, &r.SCategory, &r.Province, &r.City, &r.District, &r.DataSyncTarget, &r.SaleManagers, &r.SchoolManagers, &r.PurchaseRule, &r.SCreateTime, &r.InsuranceTypeParentID, &r.InsureAttachID, &r.PolicyNo, &r.Others, &r.Files, &r.AttachAddi, &r.AttachCreateTime, &r.AttachUpdateTime, &r.AttachCreator, &r.PolicyUploadStatus, &r.InvoiceUploadStatus)
	if err != nil {
		return errors.Wrap(err, "failed to insert t_v_insurance_policy")
	}
	return nil
}

// GetTVInsurancePolicyByPk select the TVInsurancePolicy from the database.
func GetTVInsurancePolicyByPk(db Queryer) (*TVInsurancePolicy, error) {
	var r TVInsurancePolicy
	err := db.QueryRow(
		`SELECT id, order_id, sn, name, policy, start, cease, year, duration, premium, policy_scheme, create_time, update_time, sn_creator, creator, addi, remark, status, insurance_type_id, is_entry_policy, favorite, trade_no, pay_order_no, insure_order_no, org_id, plan_id, batch, insured_id, policyholder_id, order_create_time, pay_time, pay_type, amount, unit_price, commence_date, expiry_date, indate, charge_mode, relation, insurance_type, policy_doc, same, order_status, order_creator, insurer, i_official_name, i_id_card_type, i_id_card_no, i_mobile_phone, i_gender, i_birthday, i_addi, h_official_name, h_id_card_type, h_id_card_no, h_mobile_phone, h_addi, subdistrict, faculty, grade, class, x_create_time, school, s_faculty, s_branches, s_category, province, city, district, data_sync_target, sale_managers, school_managers, purchase_rule, s_create_time, insurance_type_parent_id, insure_attach_id, policy_no, others, files, attach_addi, attach_create_time, attach_update_time, attach_creator, policy_upload_status, invoice_upload_status FROM t_v_insurance_policy WHERE `,
	).Scan(&r.ID, &r.OrderID, &r.Sn, &r.Name, &r.Policy, &r.Start, &r.Cease, &r.Year, &r.Duration, &r.Premium, &r.PolicyScheme, &r.CreateTime, &r.UpdateTime, &r.SnCreator, &r.Creator, &r.Addi, &r.Remark, &r.Status, &r.InsuranceTypeID, &r.IsEntryPolicy, &r.Favorite, &r.TradeNo, &r.PayOrderNo, &r.InsureOrderNo, &r.OrgID, &r.PlanID, &r.Batch, &r.InsuredID, &r.PolicyholderID, &r.OrderCreateTime, &r.PayTime, &r.PayType, &r.Amount, &r.UnitPrice, &r.CommenceDate, &r.ExpiryDate, &r.Indate, &r.ChargeMode, &r.Relation, &r.InsuranceType, &r.PolicyDoc, &r.Same, &r.OrderStatus, &r.OrderCreator, &r.Insurer, &r.IOfficialName, &r.IIDCardType, &r.IIDCardNo, &r.IMobilePhone, &r.IGender, &r.IBirthday, &r.IAddi, &r.HOfficialName, &r.HIDCardType, &r.HIDCardNo, &r.HMobilePhone, &r.HAddi, &r.Subdistrict, &r.Faculty, &r.Grade, &r.Class, &r.XCreateTime, &r.School, &r.SFaculty, &r.SBranches, &r.SCategory, &r.Province, &r.City, &r.District, &r.DataSyncTarget, &r.SaleManagers, &r.SchoolManagers, &r.PurchaseRule, &r.SCreateTime, &r.InsuranceTypeParentID, &r.InsureAttachID, &r.PolicyNo, &r.Others, &r.Files, &r.AttachAddi, &r.AttachCreateTime, &r.AttachUpdateTime, &r.AttachCreator, &r.PolicyUploadStatus, &r.InvoiceUploadStatus)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select t_v_insurance_policy")
	}
	return &r, nil
}

/*TVInsurancePolicy2 t_v_insurance_policy2 represents kuser.t_v_insurance_policy2 */
type TVInsurancePolicy2 struct {
	ID                      null.Int       `json:"ID,omitempty" db:"id,false,integer"`                                                       /* id id */
	Sn                      null.String    `json:"Sn,omitempty" db:"sn,false,character varying"`                                             /* sn sn */
	SnCreator               null.Int       `json:"SnCreator,omitempty" db:"sn_creator,false,bigint"`                                         /* sn_creator sn_creator */
	Name                    null.String    `json:"Name,omitempty" db:"name,false,character varying"`                                         /* name name */
	OrderID                 null.Int       `json:"OrderID,omitempty" db:"order_id,false,bigint"`                                             /* order_id order_id */
	Policy                  null.String    `json:"Policy,omitempty" db:"policy,false,character varying"`                                     /* policy policy */
	Start                   null.Int       `json:"Start,omitempty" db:"start,false,bigint"`                                                  /* start start */
	Cease                   null.Int       `json:"Cease,omitempty" db:"cease,false,bigint"`                                                  /* cease cease */
	Year                    null.Int       `json:"Year,omitempty" db:"year,false,smallint"`                                                  /* year year */
	Duration                null.Int       `json:"Duration,omitempty" db:"duration,false,bigint"`                                            /* duration duration */
	Premium                 null.Float     `json:"Premium,omitempty" db:"premium,false,double precision"`                                    /* premium premium */
	ThirdPartyPremium       null.Float     `json:"ThirdPartyPremium,omitempty" db:"third_party_premium,false,double precision"`              /* third_party_premium third_party_premium */
	InsuranceTypeID         null.Int       `json:"InsuranceTypeID,omitempty" db:"insurance_type_id,false,bigint"`                            /* insurance_type_id insurance_type_id */
	PlanID                  null.Int       `json:"PlanID,omitempty" db:"plan_id,false,bigint"`                                               /* plan_id plan_id */
	CreateTime              null.Int       `json:"CreateTime,omitempty" db:"create_time,false,bigint"`                                       /* create_time create_time */
	UpdateTime              null.Int       `json:"UpdateTime,omitempty" db:"update_time,false,bigint"`                                       /* update_time update_time */
	PayTime                 null.Int       `json:"PayTime,omitempty" db:"pay_time,false,bigint"`                                             /* pay_time pay_time */
	PayChannel              null.String    `json:"PayChannel,omitempty" db:"pay_channel,false,character varying"`                            /* pay_channel pay_channel */
	PayType                 null.String    `json:"PayType,omitempty" db:"pay_type,false,character varying"`                                  /* pay_type pay_type */
	UnitPrice               null.Float     `json:"UnitPrice,omitempty" db:"unit_price,false,double precision"`                               /* unit_price unit_price */
	ExternalStatus          null.String    `json:"ExternalStatus,omitempty" db:"external_status,false,character varying"`                    /* external_status external_status */
	OrgID                   null.Int       `json:"OrgID,omitempty" db:"org_id,false,bigint"`                                                 /* org_id org_id */
	Insurer                 null.String    `json:"Insurer,omitempty" db:"insurer,false,character varying"`                                   /* insurer insurer */
	OrgManagerID            null.Int       `json:"OrgManagerID,omitempty" db:"org_manager_id,false,bigint"`                                  /* org_manager_id org_manager_id */
	InsuranceType           null.String    `json:"InsuranceType,omitempty" db:"insurance_type,false,character varying"`                      /* insurance_type insurance_type */
	PolicyScheme            types.JSONText `json:"PolicyScheme,omitempty" db:"policy_scheme,false,jsonb"`                                    /* policy_scheme policy_scheme */
	ActivityName            null.String    `json:"ActivityName,omitempty" db:"activity_name,false,character varying"`                        /* activity_name activity_name */
	ActivityCategory        null.String    `json:"ActivityCategory,omitempty" db:"activity_category,false,character varying"`                /* activity_category activity_category */
	ActivityDesc            null.String    `json:"ActivityDesc,omitempty" db:"activity_desc,false,character varying"`                        /* activity_desc activity_desc */
	ActivityLocation        null.String    `json:"ActivityLocation,omitempty" db:"activity_location,false,character varying"`                /* activity_location activity_location */
	ActivityDateSet         null.String    `json:"ActivityDateSet,omitempty" db:"activity_date_set,false,character varying"`                 /* activity_date_set activity_date_set */
	InsuredCount            null.Int       `json:"InsuredCount,omitempty" db:"insured_count,false,smallint"`                                 /* insured_count insured_count */
	CompulsoryStudentNum    null.Int       `json:"CompulsoryStudentNum,omitempty" db:"compulsory_student_num,false,bigint"`                  /* compulsory_student_num compulsory_student_num */
	NonCompulsoryStudentNum null.Int       `json:"NonCompulsoryStudentNum,omitempty" db:"non_compulsory_student_num,false,bigint"`           /* non_compulsory_student_num non_compulsory_student_num */
	Contact                 types.JSONText `json:"Contact,omitempty" db:"contact,false,jsonb"`                                               /* contact contact */
	FeeScheme               types.JSONText `json:"FeeScheme,omitempty" db:"fee_scheme,false,jsonb"`                                          /* fee_scheme fee_scheme */
	CarServiceTarget        null.String    `json:"CarServiceTarget,omitempty" db:"car_service_target,false,character varying"`               /* car_service_target car_service_target */
	PolicyEnrollTime        null.Int       `json:"PolicyEnrollTime,omitempty" db:"policy_enroll_time,false,bigint"`                          /* policy_enroll_time policy_enroll_time */
	Policyholder            types.JSONText `json:"Policyholder,omitempty" db:"policyholder,false,jsonb"`                                     /* policyholder policyholder */
	PolicyholderType        null.String    `json:"PolicyholderType,omitempty" db:"policyholder_type,false,character varying"`                /* policyholder_type policyholder_type */
	PolicyholderID          null.Int       `json:"PolicyholderID,omitempty" db:"policyholder_id,false,bigint"`                               /* policyholder_id policyholder_id */
	OrgName                 null.String    `json:"OrgName,omitempty" db:"org_name,false,text"`                                               /* org_name org_name */
	OrgProvince             null.String    `json:"OrgProvince,omitempty" db:"org_province,false,text"`                                       /* org_province org_province */
	OrgCity                 null.String    `json:"OrgCity,omitempty" db:"org_city,false,text"`                                               /* org_city org_city */
	OrgDistrict             null.String    `json:"OrgDistrict,omitempty" db:"org_district,false,text"`                                       /* org_district org_district */
	OrgSchoolCategory       null.String    `json:"OrgSchoolCategory,omitempty" db:"org_school_category,false,text"`                          /* org_school_category org_school_category */
	OrgIsCompulsory         null.String    `json:"OrgIsCompulsory,omitempty" db:"org_is_compulsory,false,text"`                              /* org_is_compulsory org_is_compulsory */
	OrgIsSchool             null.String    `json:"OrgIsSchool,omitempty" db:"org_is_school,false,text"`                                      /* org_is_school org_is_school */
	InsuredProvince         null.String    `json:"InsuredProvince,omitempty" db:"insured_province,false,text"`                               /* insured_province insured_province */
	InsuredCity             null.String    `json:"InsuredCity,omitempty" db:"insured_city,false,text"`                                       /* insured_city insured_city */
	InsuredDistrict         null.String    `json:"InsuredDistrict,omitempty" db:"insured_district,false,text"`                               /* insured_district insured_district */
	InsuredSchoolCategory   null.String    `json:"InsuredSchoolCategory,omitempty" db:"insured_school_category,false,text"`                  /* insured_school_category insured_school_category */
	InsuredIsCompulsory     null.Bool      `json:"InsuredIsCompulsory,omitempty" db:"insured_is_compulsory,false,boolean"`                   /* insured_is_compulsory insured_is_compulsory */
	InsuredName             null.String    `json:"InsuredName,omitempty" db:"insured_name,false,text"`                                       /* insured_name insured_name */
	InsuredCategory         null.String    `json:"InsuredCategory,omitempty" db:"insured_category,false,text"`                               /* insured_category insured_category */
	DriverSeatSum           null.Int       `json:"DriverSeatSum,omitempty" db:"driver_seat_sum,false,bigint"`                                /* driver_seat_sum driver_seat_sum */
	SeatSum                 null.Int       `json:"SeatSum,omitempty" db:"seat_sum,false,bigint"`                                             /* seat_sum seat_sum */
	Same                    null.Bool      `json:"Same,omitempty" db:"same,false,boolean"`                                                   /* same same */
	Relation                null.String    `json:"Relation,omitempty" db:"relation,false,character varying"`                                 /* relation relation */
	Insured                 types.JSONText `json:"Insured,omitempty" db:"insured,false,jsonb"`                                               /* insured insured */
	InsuredID               null.Int       `json:"InsuredID,omitempty" db:"insured_id,false,bigint"`                                         /* insured_id insured_id */
	HaveInsuredList         null.Bool      `json:"HaveInsuredList,omitempty" db:"have_insured_list,false,boolean"`                           /* have_insured_list have_insured_list */
	InsuredGroupByDay       null.Bool      `json:"InsuredGroupByDay,omitempty" db:"insured_group_by_day,false,boolean"`                      /* insured_group_by_day insured_group_by_day */
	InsuredType             null.String    `json:"InsuredType,omitempty" db:"insured_type,false,character varying"`                          /* insured_type insured_type */
	InsuredList             types.JSONText `json:"InsuredList,omitempty" db:"insured_list,false,jsonb"`                                      /* insured_list insured_list */
	Indate                  null.Int       `json:"Indate,omitempty" db:"indate,false,bigint"`                                                /* indate indate */
	Jurisdiction            null.String    `json:"Jurisdiction,omitempty" db:"jurisdiction,false,character varying"`                         /* jurisdiction jurisdiction */
	DisputeHandling         null.String    `json:"DisputeHandling,omitempty" db:"dispute_handling,false,character varying"`                  /* dispute_handling dispute_handling */
	PrevPolicyNo            null.String    `json:"PrevPolicyNo,omitempty" db:"prev_policy_no,false,character varying"`                       /* prev_policy_no prev_policy_no */
	InsureBase              null.String    `json:"InsureBase,omitempty" db:"insure_base,false,character varying"`                            /* insure_base insure_base */
	BlanketInsureCode       null.String    `json:"BlanketInsureCode,omitempty" db:"blanket_insure_code,false,character varying"`             /* blanket_insure_code blanket_insure_code */
	CustomType              null.String    `json:"CustomType,omitempty" db:"custom_type,false,character varying"`                            /* custom_type custom_type */
	TrainProjects           null.String    `json:"TrainProjects,omitempty" db:"train_projects,false,character varying"`                      /* train_projects train_projects */
	BusinessLocations       types.JSONText `json:"BusinessLocations,omitempty" db:"business_locations,false,jsonb"`                          /* business_locations business_locations */
	PoolNum                 null.Int       `json:"PoolNum,omitempty" db:"pool_num,false,smallint"`                                           /* pool_num pool_num */
	HaveDinnerNum           null.Bool      `json:"HaveDinnerNum,omitempty" db:"have_dinner_num,false,boolean"`                               /* have_dinner_num have_dinner_num */
	OpenPoolNum             null.Int       `json:"OpenPoolNum,omitempty" db:"open_pool_num,false,smallint"`                                  /* open_pool_num open_pool_num */
	HeatedPoolNum           null.Int       `json:"HeatedPoolNum,omitempty" db:"heated_pool_num,false,smallint"`                              /* heated_pool_num heated_pool_num */
	TrainingPoolNum         null.Int       `json:"TrainingPoolNum,omitempty" db:"training_pool_num,false,smallint"`                          /* training_pool_num training_pool_num */
	InnerArea               null.Float     `json:"InnerArea,omitempty" db:"inner_area,false,double precision"`                               /* inner_area inner_area */
	OuterArea               null.Float     `json:"OuterArea,omitempty" db:"outer_area,false,double precision"`                               /* outer_area outer_area */
	PoolName                null.String    `json:"PoolName,omitempty" db:"pool_name,false,character varying"`                                /* pool_name pool_name */
	ArbitralAgency          null.String    `json:"ArbitralAgency,omitempty" db:"arbitral_agency,false,character varying"`                    /* arbitral_agency arbitral_agency */
	DinnerNum               null.Int       `json:"DinnerNum,omitempty" db:"dinner_num,false,integer"`                                        /* dinner_num dinner_num */
	CanteenNum              null.Int       `json:"CanteenNum,omitempty" db:"canteen_num,false,integer"`                                      /* canteen_num canteen_num */
	ShopNum                 null.Int       `json:"ShopNum,omitempty" db:"shop_num,false,integer"`                                            /* shop_num shop_num */
	HaveRides               null.Bool      `json:"HaveRides,omitempty" db:"have_rides,false,boolean"`                                        /* have_rides have_rides */
	HaveExplosive           null.Bool      `json:"HaveExplosive,omitempty" db:"have_explosive,false,boolean"`                                /* have_explosive have_explosive */
	Area                    null.Int       `json:"Area,omitempty" db:"area,false,integer"`                                                   /* area area */
	TrafficNum              null.Int       `json:"TrafficNum,omitempty" db:"traffic_num,false,integer"`                                      /* traffic_num traffic_num */
	TemperatureType         null.String    `json:"TemperatureType,omitempty" db:"temperature_type,false,character varying"`                  /* temperature_type temperature_type */
	IsIndoor                null.String    `json:"IsIndoor,omitempty" db:"is_indoor,false,character varying"`                                /* is_indoor is_indoor */
	Extra                   types.JSONText `json:"Extra,omitempty" db:"extra,false,jsonb"`                                                   /* extra extra */
	BankAccount             types.JSONText `json:"BankAccount,omitempty" db:"bank_account,false,jsonb"`                                      /* bank_account bank_account */
	PayContact              null.String    `json:"PayContact,omitempty" db:"pay_contact,false,character varying"`                            /* pay_contact pay_contact */
	SuddenDeathTerms        null.String    `json:"SuddenDeathTerms,omitempty" db:"sudden_death_terms,false,character varying"`               /* sudden_death_terms sudden_death_terms */
	HaveSuddenDeath         null.Bool      `json:"HaveSuddenDeath,omitempty" db:"have_sudden_death,false,boolean"`                           /* have_sudden_death have_sudden_death */
	SpecAgreement           null.String    `json:"SpecAgreement,omitempty" db:"spec_agreement,false,character varying"`                      /* spec_agreement spec_agreement */
	ThirdPartyAccount       null.String    `json:"ThirdPartyAccount,omitempty" db:"third_party_account,false,character varying"`             /* third_party_account third_party_account */
	Creator                 null.Int       `json:"Creator,omitempty" db:"creator,false,bigint"`                                              /* creator creator */
	DomainID                null.Int       `json:"DomainID,omitempty" db:"domain_id,false,bigint"`                                           /* domain_id domain_id */
	Status                  null.String    `json:"Status,omitempty" db:"status,false,character varying"`                                     /* status status */
	IsEntryPolicy           null.Bool      `json:"IsEntryPolicy,omitempty" db:"is_entry_policy,false,boolean"`                               /* is_entry_policy is_entry_policy */
	IsAdminPay              null.Bool      `json:"IsAdminPay,omitempty" db:"is_admin_pay,false,boolean"`                                     /* is_admin_pay is_admin_pay */
	Favorite                null.Bool      `json:"Favorite,omitempty" db:"favorite,false,boolean"`                                           /* favorite favorite */
	CancelDesc              null.String    `json:"CancelDesc,omitempty" db:"cancel_desc,false,character varying"`                            /* cancel_desc cancel_desc */
	ZeroPayStatus           null.String    `json:"ZeroPayStatus,omitempty" db:"zero_pay_status,false,character varying"`                     /* zero_pay_status zero_pay_status */
	Others                  types.JSONText `json:"Others,omitempty" db:"others,false,jsonb"`                                                 /* others others */
	Files                   types.JSONText `json:"Files,omitempty" db:"files,false,jsonb"`                                                   /* files files */
	InsurePolicyID          null.Int       `json:"InsurePolicyID,omitempty" db:"insure_policy_id,false,bigint"`                              /* insure_policy_id insure_policy_id */
	AStatus                 null.String    `json:"AStatus,omitempty" db:"a_status,false,character varying"`                                  /* a_status a_status */
	PolicyStatus            null.String    `json:"PolicyStatus,omitempty" db:"policy_status,false,text"`                                     /* policy_status policy_status */
	InsuranceTypeParentID   null.Int       `json:"InsuranceTypeParentID,omitempty" db:"insurance_type_parent_id,false,bigint"`               /* insurance_type_parent_id insurance_type_parent_id */
	InsuranceDisplay        null.String    `json:"InsuranceDisplay,omitempty" db:"insurance_display,false,character varying"`                /* insurance_display insurance_display */
	ChargeMode              null.String    `json:"ChargeMode,omitempty" db:"charge_mode,false,character varying"`                            /* charge_mode charge_mode */
	InsuranceCompany        null.String    `json:"InsuranceCompany,omitempty" db:"insurance_company,false,character varying"`                /* insurance_company insurance_company */
	InsuranceCompanyAccount null.String    `json:"InsuranceCompanyAccount,omitempty" db:"insurance_company_account,false,character varying"` /* insurance_company_account insurance_company_account */
	OrderCreateTime         null.Int       `json:"OrderCreateTime,omitempty" db:"order_create_time,false,bigint"`                            /* order_create_time order_create_time */
	IsInvoice               null.Bool      `json:"IsInvoice,omitempty" db:"is_invoice,false,boolean"`                                        /* is_invoice is_invoice */
	InvBorrow               null.String    `json:"InvBorrow,omitempty" db:"inv_borrow,false,character varying"`                              /* inv_borrow inv_borrow */
	InvVisible              null.String    `json:"InvVisible,omitempty" db:"inv_visible,false,character varying"`                            /* inv_visible inv_visible */
	InvTitle                null.String    `json:"InvTitle,omitempty" db:"inv_title,false,character varying"`                                /* inv_title inv_title */
	InvStatus               null.String    `json:"InvStatus,omitempty" db:"inv_status,false,character varying"`                              /* inv_status inv_status */
	OFiles                  types.JSONText `json:"OFiles,omitempty" db:"o_files,false,jsonb"`                                                /* o_files o_files */
	PolicyUploadStatus      null.String    `json:"PolicyUploadStatus,omitempty" db:"policy_upload_status,false,text"`                        /* policy_upload_status policy_upload_status */
	InvoiceUploadStatus     null.String    `json:"InvoiceUploadStatus,omitempty" db:"invoice_upload_status,false,text"`                      /* invoice_upload_status invoice_upload_status */
	Filter                                 // build DML where clause
}

//TVInsurancePolicy2Fields full field list for default query
var TVInsurancePolicy2Fields = []string{
	"ID",
	"Sn",
	"SnCreator",
	"Name",
	"OrderID",
	"Policy",
	"Start",
	"Cease",
	"Year",
	"Duration",
	"Premium",
	"ThirdPartyPremium",
	"InsuranceTypeID",
	"PlanID",
	"CreateTime",
	"UpdateTime",
	"PayTime",
	"PayChannel",
	"PayType",
	"UnitPrice",
	"ExternalStatus",
	"OrgID",
	"Insurer",
	"OrgManagerID",
	"InsuranceType",
	"PolicyScheme",
	"ActivityName",
	"ActivityCategory",
	"ActivityDesc",
	"ActivityLocation",
	"ActivityDateSet",
	"InsuredCount",
	"CompulsoryStudentNum",
	"NonCompulsoryStudentNum",
	"Contact",
	"FeeScheme",
	"CarServiceTarget",
	"PolicyEnrollTime",
	"Policyholder",
	"PolicyholderType",
	"PolicyholderID",
	"OrgName",
	"OrgProvince",
	"OrgCity",
	"OrgDistrict",
	"OrgSchoolCategory",
	"OrgIsCompulsory",
	"OrgIsSchool",
	"InsuredProvince",
	"InsuredCity",
	"InsuredDistrict",
	"InsuredSchoolCategory",
	"InsuredIsCompulsory",
	"InsuredName",
	"InsuredCategory",
	"DriverSeatSum",
	"SeatSum",
	"Same",
	"Relation",
	"Insured",
	"InsuredID",
	"HaveInsuredList",
	"InsuredGroupByDay",
	"InsuredType",
	"InsuredList",
	"Indate",
	"Jurisdiction",
	"DisputeHandling",
	"PrevPolicyNo",
	"InsureBase",
	"BlanketInsureCode",
	"CustomType",
	"TrainProjects",
	"BusinessLocations",
	"PoolNum",
	"HaveDinnerNum",
	"OpenPoolNum",
	"HeatedPoolNum",
	"TrainingPoolNum",
	"InnerArea",
	"OuterArea",
	"PoolName",
	"ArbitralAgency",
	"DinnerNum",
	"CanteenNum",
	"ShopNum",
	"HaveRides",
	"HaveExplosive",
	"Area",
	"TrafficNum",
	"TemperatureType",
	"IsIndoor",
	"Extra",
	"BankAccount",
	"PayContact",
	"SuddenDeathTerms",
	"HaveSuddenDeath",
	"SpecAgreement",
	"ThirdPartyAccount",
	"Creator",
	"DomainID",
	"Status",
	"IsEntryPolicy",
	"IsAdminPay",
	"Favorite",
	"CancelDesc",
	"ZeroPayStatus",
	"Others",
	"Files",
	"InsurePolicyID",
	"AStatus",
	"PolicyStatus",
	"InsuranceTypeParentID",
	"InsuranceDisplay",
	"ChargeMode",
	"InsuranceCompany",
	"InsuranceCompanyAccount",
	"OrderCreateTime",
	"IsInvoice",
	"InvBorrow",
	"InvVisible",
	"InvTitle",
	"InvStatus",
	"OFiles",
	"PolicyUploadStatus",
	"InvoiceUploadStatus",
}

//Fields return all fields of struct.
func (r *TVInsurancePolicy2) Fields() []string {
	return TVInsurancePolicy2Fields
}

//GetTableName return the associated db table name.
func (r *TVInsurancePolicy2) GetTableName() string {
	var viewNamePattern = regexp.MustCompile(`(?i)^t_v_[a-z0-9_]+$`)
	tableName := "t_v_insurance_policy2"
	if viewNamePattern.MatchString(tableName) {
		return tableName[2:]
	}
	return tableName
}

// Create inserts the TVInsurancePolicy2 to the database.
func (r *TVInsurancePolicy2) Create(db Queryer) error {
	_, err := db.Exec(
		`INSERT INTO t_v_insurance_policy2 (id, sn, sn_creator, name, order_id, policy, start, cease, year, duration, premium, third_party_premium, insurance_type_id, plan_id, create_time, update_time, pay_time, pay_channel, pay_type, unit_price, external_status, org_id, insurer, org_manager_id, insurance_type, policy_scheme, activity_name, activity_category, activity_desc, activity_location, activity_date_set, insured_count, compulsory_student_num, non_compulsory_student_num, contact, fee_scheme, car_service_target, policy_enroll_time, policyholder, policyholder_type, policyholder_id, org_name, org_province, org_city, org_district, org_school_category, org_is_compulsory, org_is_school, insured_province, insured_city, insured_district, insured_school_category, insured_is_compulsory, insured_name, insured_category, driver_seat_sum, seat_sum, same, relation, insured, insured_id, have_insured_list, insured_group_by_day, insured_type, insured_list, indate, jurisdiction, dispute_handling, prev_policy_no, insure_base, blanket_insure_code, custom_type, train_projects, business_locations, pool_num, have_dinner_num, open_pool_num, heated_pool_num, training_pool_num, inner_area, outer_area, pool_name, arbitral_agency, dinner_num, canteen_num, shop_num, have_rides, have_explosive, area, traffic_num, temperature_type, is_indoor, extra, bank_account, pay_contact, sudden_death_terms, have_sudden_death, spec_agreement, third_party_account, creator, domain_id, status, is_entry_policy, is_admin_pay, favorite, cancel_desc, zero_pay_status, others, files, insure_policy_id, a_status, policy_status, insurance_type_parent_id, insurance_display, charge_mode, insurance_company, insurance_company_account, order_create_time, is_invoice, inv_borrow, inv_visible, inv_title, inv_status, o_files, policy_upload_status, invoice_upload_status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36, $37, $38, $39, $40, $41, $42, $43, $44, $45, $46, $47, $48, $49, $50, $51, $52, $53, $54, $55, $56, $57, $58, $59, $60, $61, $62, $63, $64, $65, $66, $67, $68, $69, $70, $71, $72, $73, $74, $75, $76, $77, $78, $79, $80, $81, $82, $83, $84, $85, $86, $87, $88, $89, $90, $91, $92, $93, $94, $95, $96, $97, $98, $99, $100, $101, $102, $103, $104, $105, $106, $107, $108, $109, $110, $111, $112, $113, $114, $115, $116, $117, $118, $119, $120, $121, $122, $123, $124, $125, $126)`,
		&r.ID, &r.Sn, &r.SnCreator, &r.Name, &r.OrderID, &r.Policy, &r.Start, &r.Cease, &r.Year, &r.Duration, &r.Premium, &r.ThirdPartyPremium, &r.InsuranceTypeID, &r.PlanID, &r.CreateTime, &r.UpdateTime, &r.PayTime, &r.PayChannel, &r.PayType, &r.UnitPrice, &r.ExternalStatus, &r.OrgID, &r.Insurer, &r.OrgManagerID, &r.InsuranceType, &r.PolicyScheme, &r.ActivityName, &r.ActivityCategory, &r.ActivityDesc, &r.ActivityLocation, &r.ActivityDateSet, &r.InsuredCount, &r.CompulsoryStudentNum, &r.NonCompulsoryStudentNum, &r.Contact, &r.FeeScheme, &r.CarServiceTarget, &r.PolicyEnrollTime, &r.Policyholder, &r.PolicyholderType, &r.PolicyholderID, &r.OrgName, &r.OrgProvince, &r.OrgCity, &r.OrgDistrict, &r.OrgSchoolCategory, &r.OrgIsCompulsory, &r.OrgIsSchool, &r.InsuredProvince, &r.InsuredCity, &r.InsuredDistrict, &r.InsuredSchoolCategory, &r.InsuredIsCompulsory, &r.InsuredName, &r.InsuredCategory, &r.DriverSeatSum, &r.SeatSum, &r.Same, &r.Relation, &r.Insured, &r.InsuredID, &r.HaveInsuredList, &r.InsuredGroupByDay, &r.InsuredType, &r.InsuredList, &r.Indate, &r.Jurisdiction, &r.DisputeHandling, &r.PrevPolicyNo, &r.InsureBase, &r.BlanketInsureCode, &r.CustomType, &r.TrainProjects, &r.BusinessLocations, &r.PoolNum, &r.HaveDinnerNum, &r.OpenPoolNum, &r.HeatedPoolNum, &r.TrainingPoolNum, &r.InnerArea, &r.OuterArea, &r.PoolName, &r.ArbitralAgency, &r.DinnerNum, &r.CanteenNum, &r.ShopNum, &r.HaveRides, &r.HaveExplosive, &r.Area, &r.TrafficNum, &r.TemperatureType, &r.IsIndoor, &r.Extra, &r.BankAccount, &r.PayContact, &r.SuddenDeathTerms, &r.HaveSuddenDeath, &r.SpecAgreement, &r.ThirdPartyAccount, &r.Creator, &r.DomainID, &r.Status, &r.IsEntryPolicy, &r.IsAdminPay, &r.Favorite, &r.CancelDesc, &r.ZeroPayStatus, &r.Others, &r.Files, &r.InsurePolicyID, &r.AStatus, &r.PolicyStatus, &r.InsuranceTypeParentID, &r.InsuranceDisplay, &r.ChargeMode, &r.InsuranceCompany, &r.InsuranceCompanyAccount, &r.OrderCreateTime, &r.IsInvoice, &r.InvBorrow, &r.InvVisible, &r.InvTitle, &r.InvStatus, &r.OFiles, &r.PolicyUploadStatus, &r.InvoiceUploadStatus)
	if err != nil {
		return errors.Wrap(err, "failed to insert t_v_insurance_policy2")
	}
	return nil
}

// GetTVInsurancePolicy2ByPk select the TVInsurancePolicy2 from the database.
func GetTVInsurancePolicy2ByPk(db Queryer) (*TVInsurancePolicy2, error) {
	var r TVInsurancePolicy2
	err := db.QueryRow(
		`SELECT id, sn, sn_creator, name, order_id, policy, start, cease, year, duration, premium, third_party_premium, insurance_type_id, plan_id, create_time, update_time, pay_time, pay_channel, pay_type, unit_price, external_status, org_id, insurer, org_manager_id, insurance_type, policy_scheme, activity_name, activity_category, activity_desc, activity_location, activity_date_set, insured_count, compulsory_student_num, non_compulsory_student_num, contact, fee_scheme, car_service_target, policy_enroll_time, policyholder, policyholder_type, policyholder_id, org_name, org_province, org_city, org_district, org_school_category, org_is_compulsory, org_is_school, insured_province, insured_city, insured_district, insured_school_category, insured_is_compulsory, insured_name, insured_category, driver_seat_sum, seat_sum, same, relation, insured, insured_id, have_insured_list, insured_group_by_day, insured_type, insured_list, indate, jurisdiction, dispute_handling, prev_policy_no, insure_base, blanket_insure_code, custom_type, train_projects, business_locations, pool_num, have_dinner_num, open_pool_num, heated_pool_num, training_pool_num, inner_area, outer_area, pool_name, arbitral_agency, dinner_num, canteen_num, shop_num, have_rides, have_explosive, area, traffic_num, temperature_type, is_indoor, extra, bank_account, pay_contact, sudden_death_terms, have_sudden_death, spec_agreement, third_party_account, creator, domain_id, status, is_entry_policy, is_admin_pay, favorite, cancel_desc, zero_pay_status, others, files, insure_policy_id, a_status, policy_status, insurance_type_parent_id, insurance_display, charge_mode, insurance_company, insurance_company_account, order_create_time, is_invoice, inv_borrow, inv_visible, inv_title, inv_status, o_files, policy_upload_status, invoice_upload_status FROM t_v_insurance_policy2 WHERE `,
	).Scan(&r.ID, &r.Sn, &r.SnCreator, &r.Name, &r.OrderID, &r.Policy, &r.Start, &r.Cease, &r.Year, &r.Duration, &r.Premium, &r.ThirdPartyPremium, &r.InsuranceTypeID, &r.PlanID, &r.CreateTime, &r.UpdateTime, &r.PayTime, &r.PayChannel, &r.PayType, &r.UnitPrice, &r.ExternalStatus, &r.OrgID, &r.Insurer, &r.OrgManagerID, &r.InsuranceType, &r.PolicyScheme, &r.ActivityName, &r.ActivityCategory, &r.ActivityDesc, &r.ActivityLocation, &r.ActivityDateSet, &r.InsuredCount, &r.CompulsoryStudentNum, &r.NonCompulsoryStudentNum, &r.Contact, &r.FeeScheme, &r.CarServiceTarget, &r.PolicyEnrollTime, &r.Policyholder, &r.PolicyholderType, &r.PolicyholderID, &r.OrgName, &r.OrgProvince, &r.OrgCity, &r.OrgDistrict, &r.OrgSchoolCategory, &r.OrgIsCompulsory, &r.OrgIsSchool, &r.InsuredProvince, &r.InsuredCity, &r.InsuredDistrict, &r.InsuredSchoolCategory, &r.InsuredIsCompulsory, &r.InsuredName, &r.InsuredCategory, &r.DriverSeatSum, &r.SeatSum, &r.Same, &r.Relation, &r.Insured, &r.InsuredID, &r.HaveInsuredList, &r.InsuredGroupByDay, &r.InsuredType, &r.InsuredList, &r.Indate, &r.Jurisdiction, &r.DisputeHandling, &r.PrevPolicyNo, &r.InsureBase, &r.BlanketInsureCode, &r.CustomType, &r.TrainProjects, &r.BusinessLocations, &r.PoolNum, &r.HaveDinnerNum, &r.OpenPoolNum, &r.HeatedPoolNum, &r.TrainingPoolNum, &r.InnerArea, &r.OuterArea, &r.PoolName, &r.ArbitralAgency, &r.DinnerNum, &r.CanteenNum, &r.ShopNum, &r.HaveRides, &r.HaveExplosive, &r.Area, &r.TrafficNum, &r.TemperatureType, &r.IsIndoor, &r.Extra, &r.BankAccount, &r.PayContact, &r.SuddenDeathTerms, &r.HaveSuddenDeath, &r.SpecAgreement, &r.ThirdPartyAccount, &r.Creator, &r.DomainID, &r.Status, &r.IsEntryPolicy, &r.IsAdminPay, &r.Favorite, &r.CancelDesc, &r.ZeroPayStatus, &r.Others, &r.Files, &r.InsurePolicyID, &r.AStatus, &r.PolicyStatus, &r.InsuranceTypeParentID, &r.InsuranceDisplay, &r.ChargeMode, &r.InsuranceCompany, &r.InsuranceCompanyAccount, &r.OrderCreateTime, &r.IsInvoice, &r.InvBorrow, &r.InvVisible, &r.InvTitle, &r.InvStatus, &r.OFiles, &r.PolicyUploadStatus, &r.InvoiceUploadStatus)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select t_v_insurance_policy2")
	}
	return &r, nil
}

/*TVInsuranceType t_v_insurance_type represents kuser.t_v_insurance_type */
type TVInsuranceType struct {
	ID                      null.Int       `json:"ID,omitempty" db:"id,false,integer"`                                               /* id id */
	ParentID                null.Int       `json:"ParentID,omitempty" db:"parent_id,false,bigint"`                                   /* parent_id parent_id */
	DataType                null.String    `json:"DataType,omitempty" db:"data_type,false,character varying"`                        /* data_type data_type */
	ParentName              null.String    `json:"ParentName,omitempty" db:"parent_name,false,character varying"`                    /* parent_name parent_name */
	OrgName                 null.String    `json:"OrgName,omitempty" db:"org_name,false,character varying"`                          /* org_name org_name */
	OrgID                   null.Int       `json:"OrgID,omitempty" db:"org_id,false,bigint"`                                         /* org_id org_id */
	LayoutOrder             null.Int       `json:"LayoutOrder,omitempty" db:"layout_order,false,smallint"`                           /* layout_order layout_order */
	Insurer                 null.String    `json:"Insurer,omitempty" db:"insurer,false,character varying"`                           /* insurer insurer */
	RefID                   null.Int       `json:"RefID,omitempty" db:"ref_id,false,bigint"`                                         /* ref_id ref_id */
	Name                    null.String    `json:"Name,omitempty" db:"name,false,character varying"`                                 /* name name */
	Alias                   null.String    `json:"Alias,omitempty" db:"alias,false,character varying"`                               /* alias alias */
	PayType                 null.String    `json:"PayType,omitempty" db:"pay_type,false,character varying"`                          /* pay_type pay_type */
	PayName                 null.String    `json:"PayName,omitempty" db:"pay_name,false,character varying"`                          /* pay_name pay_name */
	PayChannel              null.String    `json:"PayChannel,omitempty" db:"pay_channel,false,character varying"`                    /* pay_channel pay_channel */
	RuleBatch               null.String    `json:"RuleBatch,omitempty" db:"rule_batch,false,character varying"`                      /* rule_batch rule_batch */
	UnitPrice               null.Float     `json:"UnitPrice,omitempty" db:"unit_price,false,double precision"`                       /* unit_price unit_price */
	Price                   null.Float     `json:"Price,omitempty" db:"price,false,double precision"`                                /* price price */
	PriceConfig             types.JSONText `json:"PriceConfig,omitempty" db:"price_config,false,jsonb"`                              /* price_config price_config */
	AllowStart              null.Int       `json:"AllowStart,omitempty" db:"allow_start,false,bigint"`                               /* allow_start allow_start */
	AllowEnd                null.Int       `json:"AllowEnd,omitempty" db:"allow_end,false,bigint"`                                   /* allow_end allow_end */
	TimeStatus              null.String    `json:"TimeStatus,omitempty" db:"time_status,false,text"`                                 /* time_status time_status */
	MaxInsureInYear         null.Int       `json:"MaxInsureInYear,omitempty" db:"max_insure_in_year,false,smallint"`                 /* max_insure_in_year max_insure_in_year */
	InsuredStartTime        null.Int       `json:"InsuredStartTime,omitempty" db:"insured_start_time,false,bigint"`                  /* insured_start_time insured_start_time */
	InsuredEndTime          null.Int       `json:"InsuredEndTime,omitempty" db:"insured_end_time,false,bigint"`                      /* insured_end_time insured_end_time */
	InsuredInMonth          null.Int       `json:"InsuredInMonth,omitempty" db:"insured_in_month,false,smallint"`                    /* insured_in_month insured_in_month */
	IndateStart             null.Int       `json:"IndateStart,omitempty" db:"indate_start,false,bigint"`                             /* indate_start indate_start */
	IndateEnd               null.Int       `json:"IndateEnd,omitempty" db:"indate_end,false,bigint"`                                 /* indate_end indate_end */
	AgeLimit                types.JSONText `json:"AgeLimit,omitempty" db:"age_limit,false,jsonb"`                                    /* age_limit age_limit */
	BankAccount             null.String    `json:"BankAccount,omitempty" db:"bank_account,false,character varying"`                  /* bank_account bank_account */
	BankAccountName         null.String    `json:"BankAccountName,omitempty" db:"bank_account_name,false,character varying"`         /* bank_account_name bank_account_name */
	BankName                null.String    `json:"BankName,omitempty" db:"bank_name,false,character varying"`                        /* bank_name bank_name */
	BankID                  null.String    `json:"BankID,omitempty" db:"bank_id,false,character varying"`                            /* bank_id bank_id */
	FloorPrice              null.Float     `json:"FloorPrice,omitempty" db:"floor_price,false,double precision"`                     /* floor_price floor_price */
	DefineLevel             null.Int       `json:"DefineLevel,omitempty" db:"define_level,false,smallint"`                           /* define_level define_level */
	LayoutLevel             null.Int       `json:"LayoutLevel,omitempty" db:"layout_level,false,smallint"`                           /* layout_level layout_level */
	ListTpl                 null.String    `json:"ListTpl,omitempty" db:"list_tpl,false,character varying"`                          /* list_tpl list_tpl */
	Files                   types.JSONText `json:"Files,omitempty" db:"files,false,jsonb"`                                           /* files files */
	Pic                     null.String    `json:"Pic,omitempty" db:"pic,false,character varying"`                                   /* pic pic */
	SuddenDeathDescription  types.JSONText `json:"SuddenDeathDescription,omitempty" db:"sudden_death_description,false,jsonb"`       /* sudden_death_description sudden_death_description */
	Description             null.String    `json:"Description,omitempty" db:"description,false,character varying"`                   /* description description */
	AutoFill                null.String    `json:"AutoFill,omitempty" db:"auto_fill,false,character varying"`                        /* auto_fill auto_fill */
	EnableImportList        null.Bool      `json:"EnableImportList,omitempty" db:"enable_import_list,false,boolean"`                 /* enable_import_list enable_import_list */
	HaveDinnerNum           null.Bool      `json:"HaveDinnerNum,omitempty" db:"have_dinner_num,false,boolean"`                       /* have_dinner_num have_dinner_num */
	InvoiceTitleUpdateTimes null.Int       `json:"InvoiceTitleUpdateTimes,omitempty" db:"invoice_title_update_times,false,smallint"` /* invoice_title_update_times invoice_title_update_times */
	TransferAuthFiles       types.JSONText `json:"TransferAuthFiles,omitempty" db:"transfer_auth_files,false,jsonb"`                 /* transfer_auth_files transfer_auth_files */
	ReceiptAccount          types.JSONText `json:"ReceiptAccount,omitempty" db:"receipt_account,false,jsonb"`                        /* receipt_account receipt_account */
	ContactQrCode           null.String    `json:"ContactQrCode,omitempty" db:"contact_qr_code,false,character varying"`             /* contact_qr_code contact_qr_code */
	OtherFiles              types.JSONText `json:"OtherFiles,omitempty" db:"other_files,false,jsonb"`                                /* other_files other_files */
	Contact                 types.JSONText `json:"Contact,omitempty" db:"contact,false,jsonb"`                                       /* contact contact */
	Underwriter             types.JSONText `json:"Underwriter,omitempty" db:"underwriter,false,jsonb"`                               /* underwriter underwriter */
	RemindDays              null.Int       `json:"RemindDays,omitempty" db:"remind_days,false,smallint"`                             /* remind_days remind_days */
	Mail                    types.JSONText `json:"Mail,omitempty" db:"mail,false,jsonb"`                                             /* mail mail */
	OrderRepeatLimit        null.Int       `json:"OrderRepeatLimit,omitempty" db:"order_repeat_limit,false,smallint"`                /* order_repeat_limit order_repeat_limit */
	GroupByMaxDay           null.Int       `json:"GroupByMaxDay,omitempty" db:"group_by_max_day,false,smallint"`                     /* group_by_max_day group_by_max_day */
	WebDescription          null.String    `json:"WebDescription,omitempty" db:"web_description,false,character varying"`            /* web_description web_description */
	MobileDescription       null.String    `json:"MobileDescription,omitempty" db:"mobile_description,false,character varying"`      /* mobile_description mobile_description */
	AutoFillParam           types.JSONText `json:"AutoFillParam,omitempty" db:"auto_fill_param,false,jsonb"`                         /* auto_fill_param auto_fill_param */
	Interval                null.Int       `json:"Interval,omitempty" db:"interval,false,bigint"`                                    /* interval interval */
	Addi                    types.JSONText `json:"Addi,omitempty" db:"addi,false,jsonb"`                                             /* addi addi */
	Resource                types.JSONText `json:"Resource,omitempty" db:"resource,false,jsonb"`                                     /* resource resource */
	CreateTime              null.Int       `json:"CreateTime,omitempty" db:"create_time,false,bigint"`                               /* create_time create_time */
	UpdateTime              null.Int       `json:"UpdateTime,omitempty" db:"update_time,false,bigint"`                               /* update_time update_time */
	Status                  null.String    `json:"Status,omitempty" db:"status,false,character varying"`                             /* status status */
	Filter                                 // build DML where clause
}

//TVInsuranceTypeFields full field list for default query
var TVInsuranceTypeFields = []string{
	"ID",
	"ParentID",
	"DataType",
	"ParentName",
	"OrgName",
	"OrgID",
	"LayoutOrder",
	"Insurer",
	"RefID",
	"Name",
	"Alias",
	"PayType",
	"PayName",
	"PayChannel",
	"RuleBatch",
	"UnitPrice",
	"Price",
	"PriceConfig",
	"AllowStart",
	"AllowEnd",
	"TimeStatus",
	"MaxInsureInYear",
	"InsuredStartTime",
	"InsuredEndTime",
	"InsuredInMonth",
	"IndateStart",
	"IndateEnd",
	"AgeLimit",
	"BankAccount",
	"BankAccountName",
	"BankName",
	"BankID",
	"FloorPrice",
	"DefineLevel",
	"LayoutLevel",
	"ListTpl",
	"Files",
	"Pic",
	"SuddenDeathDescription",
	"Description",
	"AutoFill",
	"EnableImportList",
	"HaveDinnerNum",
	"InvoiceTitleUpdateTimes",
	"TransferAuthFiles",
	"ReceiptAccount",
	"ContactQrCode",
	"OtherFiles",
	"Contact",
	"Underwriter",
	"RemindDays",
	"Mail",
	"OrderRepeatLimit",
	"GroupByMaxDay",
	"WebDescription",
	"MobileDescription",
	"AutoFillParam",
	"Interval",
	"Addi",
	"Resource",
	"CreateTime",
	"UpdateTime",
	"Status",
}

//Fields return all fields of struct.
func (r *TVInsuranceType) Fields() []string {
	return TVInsuranceTypeFields
}

//GetTableName return the associated db table name.
func (r *TVInsuranceType) GetTableName() string {
	var viewNamePattern = regexp.MustCompile(`(?i)^t_v_[a-z0-9_]+$`)
	tableName := "t_v_insurance_type"
	if viewNamePattern.MatchString(tableName) {
		return tableName[2:]
	}
	return tableName
}

// Create inserts the TVInsuranceType to the database.
func (r *TVInsuranceType) Create(db Queryer) error {
	_, err := db.Exec(
		`INSERT INTO t_v_insurance_type (id, parent_id, data_type, parent_name, org_name, org_id, layout_order, insurer, ref_id, name, alias, pay_type, pay_name, pay_channel, rule_batch, unit_price, price, price_config, allow_start, allow_end, time_status, max_insure_in_year, insured_start_time, insured_end_time, insured_in_month, indate_start, indate_end, age_limit, bank_account, bank_account_name, bank_name, bank_id, floor_price, define_level, layout_level, list_tpl, files, pic, sudden_death_description, description, auto_fill, enable_import_list, have_dinner_num, invoice_title_update_times, transfer_auth_files, receipt_account, contact_qr_code, other_files, contact, underwriter, remind_days, mail, order_repeat_limit, group_by_max_day, web_description, mobile_description, auto_fill_param, interval, addi, resource, create_time, update_time, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36, $37, $38, $39, $40, $41, $42, $43, $44, $45, $46, $47, $48, $49, $50, $51, $52, $53, $54, $55, $56, $57, $58, $59, $60, $61, $62, $63)`,
		&r.ID, &r.ParentID, &r.DataType, &r.ParentName, &r.OrgName, &r.OrgID, &r.LayoutOrder, &r.Insurer, &r.RefID, &r.Name, &r.Alias, &r.PayType, &r.PayName, &r.PayChannel, &r.RuleBatch, &r.UnitPrice, &r.Price, &r.PriceConfig, &r.AllowStart, &r.AllowEnd, &r.TimeStatus, &r.MaxInsureInYear, &r.InsuredStartTime, &r.InsuredEndTime, &r.InsuredInMonth, &r.IndateStart, &r.IndateEnd, &r.AgeLimit, &r.BankAccount, &r.BankAccountName, &r.BankName, &r.BankID, &r.FloorPrice, &r.DefineLevel, &r.LayoutLevel, &r.ListTpl, &r.Files, &r.Pic, &r.SuddenDeathDescription, &r.Description, &r.AutoFill, &r.EnableImportList, &r.HaveDinnerNum, &r.InvoiceTitleUpdateTimes, &r.TransferAuthFiles, &r.ReceiptAccount, &r.ContactQrCode, &r.OtherFiles, &r.Contact, &r.Underwriter, &r.RemindDays, &r.Mail, &r.OrderRepeatLimit, &r.GroupByMaxDay, &r.WebDescription, &r.MobileDescription, &r.AutoFillParam, &r.Interval, &r.Addi, &r.Resource, &r.CreateTime, &r.UpdateTime, &r.Status)
	if err != nil {
		return errors.Wrap(err, "failed to insert t_v_insurance_type")
	}
	return nil
}

// GetTVInsuranceTypeByPk select the TVInsuranceType from the database.
func GetTVInsuranceTypeByPk(db Queryer) (*TVInsuranceType, error) {
	var r TVInsuranceType
	err := db.QueryRow(
		`SELECT id, parent_id, data_type, parent_name, org_name, org_id, layout_order, insurer, ref_id, name, alias, pay_type, pay_name, pay_channel, rule_batch, unit_price, price, price_config, allow_start, allow_end, time_status, max_insure_in_year, insured_start_time, insured_end_time, insured_in_month, indate_start, indate_end, age_limit, bank_account, bank_account_name, bank_name, bank_id, floor_price, define_level, layout_level, list_tpl, files, pic, sudden_death_description, description, auto_fill, enable_import_list, have_dinner_num, invoice_title_update_times, transfer_auth_files, receipt_account, contact_qr_code, other_files, contact, underwriter, remind_days, mail, order_repeat_limit, group_by_max_day, web_description, mobile_description, auto_fill_param, interval, addi, resource, create_time, update_time, status FROM t_v_insurance_type WHERE `,
	).Scan(&r.ID, &r.ParentID, &r.DataType, &r.ParentName, &r.OrgName, &r.OrgID, &r.LayoutOrder, &r.Insurer, &r.RefID, &r.Name, &r.Alias, &r.PayType, &r.PayName, &r.PayChannel, &r.RuleBatch, &r.UnitPrice, &r.Price, &r.PriceConfig, &r.AllowStart, &r.AllowEnd, &r.TimeStatus, &r.MaxInsureInYear, &r.InsuredStartTime, &r.InsuredEndTime, &r.InsuredInMonth, &r.IndateStart, &r.IndateEnd, &r.AgeLimit, &r.BankAccount, &r.BankAccountName, &r.BankName, &r.BankID, &r.FloorPrice, &r.DefineLevel, &r.LayoutLevel, &r.ListTpl, &r.Files, &r.Pic, &r.SuddenDeathDescription, &r.Description, &r.AutoFill, &r.EnableImportList, &r.HaveDinnerNum, &r.InvoiceTitleUpdateTimes, &r.TransferAuthFiles, &r.ReceiptAccount, &r.ContactQrCode, &r.OtherFiles, &r.Contact, &r.Underwriter, &r.RemindDays, &r.Mail, &r.OrderRepeatLimit, &r.GroupByMaxDay, &r.WebDescription, &r.MobileDescription, &r.AutoFillParam, &r.Interval, &r.Addi, &r.Resource, &r.CreateTime, &r.UpdateTime, &r.Status)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select t_v_insurance_type")
	}
	return &r, nil
}

/*TVInsureAttach t_v_insure_attach represents kuser.t_v_insure_attach */
type TVInsureAttach struct {
	ID                  null.Int       `json:"ID,omitempty" db:"id,false,integer"`                                  /* id id */
	OrgID               null.Int       `json:"OrgID,omitempty" db:"org_id,false,bigint"`                            /* org_id org_id */
	School              null.String    `json:"School,omitempty" db:"school,false,character varying"`                /* school school */
	Category            null.String    `json:"Category,omitempty" db:"category,false,character varying"`            /* category category */
	Insurer             null.String    `json:"Insurer,omitempty" db:"insurer,false,character varying"`              /* insurer insurer */
	InsuranceType       null.String    `json:"InsuranceType,omitempty" db:"insurance_type,false,character varying"` /* insurance_type insurance_type */
	Year                null.Int       `json:"Year,omitempty" db:"year,false,smallint"`                             /* year year */
	Batch               null.String    `json:"Batch,omitempty" db:"batch,false,character varying"`                  /* batch batch */
	Grade               null.String    `json:"Grade,omitempty" db:"grade,false,character varying"`                  /* grade grade */
	Files               types.JSONText `json:"Files,omitempty" db:"files,false,jsonb"`                              /* files files */
	PolicyNo            null.String    `json:"PolicyNo,omitempty" db:"policy_no,false,character varying"`           /* policy_no policy_no */
	Others              types.JSONText `json:"Others,omitempty" db:"others,false,jsonb"`                            /* others others */
	PolicyUploadStatus  null.String    `json:"PolicyUploadStatus,omitempty" db:"policy_upload_status,false,text"`   /* policy_upload_status policy_upload_status */
	InvoiceUploadStatus null.String    `json:"InvoiceUploadStatus,omitempty" db:"invoice_upload_status,false,text"` /* invoice_upload_status invoice_upload_status */
	Addi                types.JSONText `json:"Addi,omitempty" db:"addi,false,jsonb"`                                /* addi addi */
	CreateTime          null.Int       `json:"CreateTime,omitempty" db:"create_time,false,bigint"`                  /* create_time create_time */
	UpdateTime          null.Int       `json:"UpdateTime,omitempty" db:"update_time,false,bigint"`                  /* update_time update_time */
	Creator             null.Int       `json:"Creator,omitempty" db:"creator,false,bigint"`                         /* creator creator */
	Filter                             // build DML where clause
}

//TVInsureAttachFields full field list for default query
var TVInsureAttachFields = []string{
	"ID",
	"OrgID",
	"School",
	"Category",
	"Insurer",
	"InsuranceType",
	"Year",
	"Batch",
	"Grade",
	"Files",
	"PolicyNo",
	"Others",
	"PolicyUploadStatus",
	"InvoiceUploadStatus",
	"Addi",
	"CreateTime",
	"UpdateTime",
	"Creator",
}

//Fields return all fields of struct.
func (r *TVInsureAttach) Fields() []string {
	return TVInsureAttachFields
}

//GetTableName return the associated db table name.
func (r *TVInsureAttach) GetTableName() string {
	var viewNamePattern = regexp.MustCompile(`(?i)^t_v_[a-z0-9_]+$`)
	tableName := "t_v_insure_attach"
	if viewNamePattern.MatchString(tableName) {
		return tableName[2:]
	}
	return tableName
}

// Create inserts the TVInsureAttach to the database.
func (r *TVInsureAttach) Create(db Queryer) error {
	_, err := db.Exec(
		`INSERT INTO t_v_insure_attach (id, org_id, school, category, insurer, insurance_type, year, batch, grade, files, policy_no, others, policy_upload_status, invoice_upload_status, addi, create_time, update_time, creator) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)`,
		&r.ID, &r.OrgID, &r.School, &r.Category, &r.Insurer, &r.InsuranceType, &r.Year, &r.Batch, &r.Grade, &r.Files, &r.PolicyNo, &r.Others, &r.PolicyUploadStatus, &r.InvoiceUploadStatus, &r.Addi, &r.CreateTime, &r.UpdateTime, &r.Creator)
	if err != nil {
		return errors.Wrap(err, "failed to insert t_v_insure_attach")
	}
	return nil
}

// GetTVInsureAttachByPk select the TVInsureAttach from the database.
func GetTVInsureAttachByPk(db Queryer) (*TVInsureAttach, error) {
	var r TVInsureAttach
	err := db.QueryRow(
		`SELECT id, org_id, school, category, insurer, insurance_type, year, batch, grade, files, policy_no, others, policy_upload_status, invoice_upload_status, addi, create_time, update_time, creator FROM t_v_insure_attach WHERE `,
	).Scan(&r.ID, &r.OrgID, &r.School, &r.Category, &r.Insurer, &r.InsuranceType, &r.Year, &r.Batch, &r.Grade, &r.Files, &r.PolicyNo, &r.Others, &r.PolicyUploadStatus, &r.InvoiceUploadStatus, &r.Addi, &r.CreateTime, &r.UpdateTime, &r.Creator)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select t_v_insure_attach")
	}
	return &r, nil
}

/*TVInsuredSchool t_v_insured_school represents kuser.t_v_insured_school */
type TVInsuredSchool struct {
	ID        null.Int       `json:"ID,omitempty" db:"id,false,integer"`                          /* id id */
	Name      null.String    `json:"Name,omitempty" db:"name,false,character varying"`            /* name name */
	Category  null.String    `json:"Category,omitempty" db:"category,false,character varying"`    /* category category */
	Province  null.String    `json:"Province,omitempty" db:"province,false,character varying"`    /* province province */
	City      null.String    `json:"City,omitempty" db:"city,false,character varying"`            /* city city */
	District  null.String    `json:"District,omitempty" db:"district,false,character varying"`    /* district district */
	Street    null.String    `json:"Street,omitempty" db:"street,false,character varying"`        /* street street */
	IsSchool  null.Bool      `json:"IsSchool,omitempty" db:"is_school,false,boolean"`             /* is_school is_school */
	OrgStatus null.String    `json:"OrgStatus,omitempty" db:"org_status,false,character varying"` /* org_status org_status */
	AllowTime types.JSONText `json:"AllowTime,omitempty" db:"allow_time,false,jsonb"`             /* allow_time allow_time */
	Filter                   // build DML where clause
}

//TVInsuredSchoolFields full field list for default query
var TVInsuredSchoolFields = []string{
	"ID",
	"Name",
	"Category",
	"Province",
	"City",
	"District",
	"Street",
	"IsSchool",
	"OrgStatus",
	"AllowTime",
}

//Fields return all fields of struct.
func (r *TVInsuredSchool) Fields() []string {
	return TVInsuredSchoolFields
}

//GetTableName return the associated db table name.
func (r *TVInsuredSchool) GetTableName() string {
	var viewNamePattern = regexp.MustCompile(`(?i)^t_v_[a-z0-9_]+$`)
	tableName := "t_v_insured_school"
	if viewNamePattern.MatchString(tableName) {
		return tableName[2:]
	}
	return tableName
}

// Create inserts the TVInsuredSchool to the database.
func (r *TVInsuredSchool) Create(db Queryer) error {
	_, err := db.Exec(
		`INSERT INTO t_v_insured_school (id, name, category, province, city, district, street, is_school, org_status, allow_time) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		&r.ID, &r.Name, &r.Category, &r.Province, &r.City, &r.District, &r.Street, &r.IsSchool, &r.OrgStatus, &r.AllowTime)
	if err != nil {
		return errors.Wrap(err, "failed to insert t_v_insured_school")
	}
	return nil
}

// GetTVInsuredSchoolByPk select the TVInsuredSchool from the database.
func GetTVInsuredSchoolByPk(db Queryer) (*TVInsuredSchool, error) {
	var r TVInsuredSchool
	err := db.QueryRow(
		`SELECT id, name, category, province, city, district, street, is_school, org_status, allow_time FROM t_v_insured_school WHERE `,
	).Scan(&r.ID, &r.Name, &r.Category, &r.Province, &r.City, &r.District, &r.Street, &r.IsSchool, &r.OrgStatus, &r.AllowTime)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select t_v_insured_school")
	}
	return &r, nil
}

/*TVInsurer t_v_insurer represents kuser.t_v_insurer */
type TVInsurer struct {
	ID       null.Int    `json:"ID,omitempty" db:"id,false,integer"`                     /* id id */
	Name     null.String `json:"Name,omitempty" db:"name,false,character varying"`       /* name name */
	RefID    null.Int    `json:"RefID,omitempty" db:"ref_id,false,bigint"`               /* ref_id ref_id */
	ParentID null.Int    `json:"ParentID,omitempty" db:"parent_id,false,bigint"`         /* parent_id parent_id */
	Insurer  null.String `json:"Insurer,omitempty" db:"insurer,false,character varying"` /* insurer insurer */
	Filter               // build DML where clause
}

//TVInsurerFields full field list for default query
var TVInsurerFields = []string{
	"ID",
	"Name",
	"RefID",
	"ParentID",
	"Insurer",
}

//Fields return all fields of struct.
func (r *TVInsurer) Fields() []string {
	return TVInsurerFields
}

//GetTableName return the associated db table name.
func (r *TVInsurer) GetTableName() string {
	var viewNamePattern = regexp.MustCompile(`(?i)^t_v_[a-z0-9_]+$`)
	tableName := "t_v_insurer"
	if viewNamePattern.MatchString(tableName) {
		return tableName[2:]
	}
	return tableName
}

// Create inserts the TVInsurer to the database.
func (r *TVInsurer) Create(db Queryer) error {
	_, err := db.Exec(
		`INSERT INTO t_v_insurer (id, name, ref_id, parent_id, insurer) VALUES ($1, $2, $3, $4, $5)`,
		&r.ID, &r.Name, &r.RefID, &r.ParentID, &r.Insurer)
	if err != nil {
		return errors.Wrap(err, "failed to insert t_v_insurer")
	}
	return nil
}

// GetTVInsurerByPk select the TVInsurer from the database.
func GetTVInsurerByPk(db Queryer) (*TVInsurer, error) {
	var r TVInsurer
	err := db.QueryRow(
		`SELECT id, name, ref_id, parent_id, insurer FROM t_v_insurer WHERE `,
	).Scan(&r.ID, &r.Name, &r.RefID, &r.ParentID, &r.Insurer)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select t_v_insurer")
	}
	return &r, nil
}

/*TVManagerSchool t_v_manager_school represents kuser.t_v_manager_school */
type TVManagerSchool struct {
	UserID     null.Int       `json:"UserID,omitempty" db:"user_id,false,bigint"`                    /* user_id user_id */
	Name       null.String    `json:"Name,omitempty" db:"name,false,character varying"`              /* name name */
	Tel        null.String    `json:"Tel,omitempty" db:"tel,false,character varying"`                /* tel tel */
	SchoolID   null.Int       `json:"SchoolID,omitempty" db:"school_id,false,bigint"`                /* school_id school_id */
	SchoolName null.String    `json:"SchoolName,omitempty" db:"school_name,false,character varying"` /* school_name school_name */
	UserRole   null.String    `json:"UserRole,omitempty" db:"user_role,false,text"`                  /* user_role user_role */
	RelType    null.String    `json:"RelType,omitempty" db:"rel_type,false,text"`                    /* rel_type rel_type */
	Addi       types.JSONText `json:"Addi,omitempty" db:"addi,false,jsonb"`                          /* addi addi */
	Filter                    // build DML where clause
}

//TVManagerSchoolFields full field list for default query
var TVManagerSchoolFields = []string{
	"UserID",
	"Name",
	"Tel",
	"SchoolID",
	"SchoolName",
	"UserRole",
	"RelType",
	"Addi",
}

//Fields return all fields of struct.
func (r *TVManagerSchool) Fields() []string {
	return TVManagerSchoolFields
}

//GetTableName return the associated db table name.
func (r *TVManagerSchool) GetTableName() string {
	var viewNamePattern = regexp.MustCompile(`(?i)^t_v_[a-z0-9_]+$`)
	tableName := "t_v_manager_school"
	if viewNamePattern.MatchString(tableName) {
		return tableName[2:]
	}
	return tableName
}

// Create inserts the TVManagerSchool to the database.
func (r *TVManagerSchool) Create(db Queryer) error {
	_, err := db.Exec(
		`INSERT INTO t_v_manager_school (user_id, name, tel, school_id, school_name, user_role, rel_type, addi) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		&r.UserID, &r.Name, &r.Tel, &r.SchoolID, &r.SchoolName, &r.UserRole, &r.RelType, &r.Addi)
	if err != nil {
		return errors.Wrap(err, "failed to insert t_v_manager_school")
	}
	return nil
}

// GetTVManagerSchoolByPk select the TVManagerSchool from the database.
func GetTVManagerSchoolByPk(db Queryer) (*TVManagerSchool, error) {
	var r TVManagerSchool
	err := db.QueryRow(
		`SELECT user_id, name, tel, school_id, school_name, user_role, rel_type, addi FROM t_v_manager_school WHERE `,
	).Scan(&r.UserID, &r.Name, &r.Tel, &r.SchoolID, &r.SchoolName, &r.UserRole, &r.RelType, &r.Addi)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select t_v_manager_school")
	}
	return &r, nil
}

/*TVMistakeCorrect t_v_mistake_correct represents kuser.t_v_mistake_correct */
type TVMistakeCorrect struct {
	ID                    null.Int       `json:"ID,omitempty" db:"id,false,integer"`                                                    /* id id */
	OrderID               null.Int       `json:"OrderID,omitempty" db:"order_id,false,bigint"`                                          /* order_id order_id */
	OrgID                 null.Int       `json:"OrgID,omitempty" db:"org_id,false,bigint"`                                              /* org_id org_id */
	CommenceDate          null.Int       `json:"CommenceDate,omitempty" db:"commence_date,false,bigint"`                                /* commence_date commence_date */
	ExpiryDate            null.Int       `json:"ExpiryDate,omitempty" db:"expiry_date,false,bigint"`                                    /* expiry_date expiry_date */
	InsuranceType         null.String    `json:"InsuranceType,omitempty" db:"insurance_type,false,character varying"`                   /* insurance_type insurance_type */
	OfficialName          null.String    `json:"OfficialName,omitempty" db:"official_name,false,character varying"`                     /* official_name official_name */
	IDCardType            null.String    `json:"IDCardType,omitempty" db:"id_card_type,false,character varying"`                        /* id_card_type id_card_type */
	IDCardNo              null.String    `json:"IDCardNo,omitempty" db:"id_card_no,false,character varying"`                            /* id_card_no id_card_no */
	Gender                null.String    `json:"Gender,omitempty" db:"gender,false,character varying"`                                  /* gender gender */
	Birthday              null.Int       `json:"Birthday,omitempty" db:"birthday,false,bigint"`                                         /* birthday birthday */
	School                null.String    `json:"School,omitempty" db:"school,false,character varying"`                                  /* school school */
	SchoolID              null.Int       `json:"SchoolID,omitempty" db:"school_id,false,bigint"`                                        /* school_id school_id */
	SchoolType            null.String    `json:"SchoolType,omitempty" db:"school_type,false,character varying"`                         /* school_type school_type */
	InsuredID             null.Int       `json:"InsuredID,omitempty" db:"insured_id,false,bigint"`                                      /* insured_id insured_id */
	OriginalOfficialName  null.String    `json:"OriginalOfficialName,omitempty" db:"original_official_name,false,character varying"`    /* original_official_name original_official_name */
	OriginalIDCardType    null.String    `json:"OriginalIDCardType,omitempty" db:"original_id_card_type,false,character varying"`       /* original_id_card_type original_id_card_type */
	OriginalIDCardNo      null.String    `json:"OriginalIDCardNo,omitempty" db:"original_id_card_no,false,character varying"`           /* original_id_card_no original_id_card_no */
	OriginalGender        null.String    `json:"OriginalGender,omitempty" db:"original_gender,false,character varying"`                 /* original_gender original_gender */
	OriginalBirthday      null.Int       `json:"OriginalBirthday,omitempty" db:"original_birthday,false,bigint"`                        /* original_birthday original_birthday */
	OfficialNameP         null.String    `json:"OfficialNameP,omitempty" db:"official_name_p,false,character varying"`                  /* official_name_p official_name_p */
	IDCardTypeP           null.String    `json:"IDCardTypeP,omitempty" db:"id_card_type_p,false,character varying"`                     /* id_card_type_p id_card_type_p */
	IDCardNoP             null.String    `json:"IDCardNoP,omitempty" db:"id_card_no_p,false,character varying"`                         /* id_card_no_p id_card_no_p */
	GenderP               null.String    `json:"GenderP,omitempty" db:"gender_p,false,character varying"`                               /* gender_p gender_p */
	BirthdayP             null.Int       `json:"BirthdayP,omitempty" db:"birthday_p,false,bigint"`                                      /* birthday_p birthday_p */
	PolicyholderID        null.Int       `json:"PolicyholderID,omitempty" db:"policyholder_id,false,bigint"`                            /* policyholder_id policyholder_id */
	OriginalOfficialNameP null.String    `json:"OriginalOfficialNameP,omitempty" db:"original_official_name_p,false,character varying"` /* original_official_name_p original_official_name_p */
	OriginalIDCardTypeP   null.String    `json:"OriginalIDCardTypeP,omitempty" db:"original_id_card_type_p,false,character varying"`    /* original_id_card_type_p original_id_card_type_p */
	OriginalIDCardNoP     null.String    `json:"OriginalIDCardNoP,omitempty" db:"original_id_card_no_p,false,character varying"`        /* original_id_card_no_p original_id_card_no_p */
	OriginalGenderP       null.String    `json:"OriginalGenderP,omitempty" db:"original_gender_p,false,character varying"`              /* original_gender_p original_gender_p */
	OriginalBirthdayP     null.Int       `json:"OriginalBirthdayP,omitempty" db:"original_birthday_p,false,bigint"`                     /* original_birthday_p original_birthday_p */
	Addi                  types.JSONText `json:"Addi,omitempty" db:"addi,false,jsonb"`                                                  /* addi addi */
	CreateTime            null.Int       `json:"CreateTime,omitempty" db:"create_time,false,bigint"`                                    /* create_time create_time */
	UpdateTime            null.Int       `json:"UpdateTime,omitempty" db:"update_time,false,bigint"`                                    /* update_time update_time */
	Creator               null.Int       `json:"Creator,omitempty" db:"creator,false,bigint"`                                           /* creator creator */
	Remark                null.String    `json:"Remark,omitempty" db:"remark,false,character varying"`                                  /* remark remark */
	Status                null.String    `json:"Status,omitempty" db:"status,false,character varying"`                                  /* status status */
	Filter                               // build DML where clause
}

//TVMistakeCorrectFields full field list for default query
var TVMistakeCorrectFields = []string{
	"ID",
	"OrderID",
	"OrgID",
	"CommenceDate",
	"ExpiryDate",
	"InsuranceType",
	"OfficialName",
	"IDCardType",
	"IDCardNo",
	"Gender",
	"Birthday",
	"School",
	"SchoolID",
	"SchoolType",
	"InsuredID",
	"OriginalOfficialName",
	"OriginalIDCardType",
	"OriginalIDCardNo",
	"OriginalGender",
	"OriginalBirthday",
	"OfficialNameP",
	"IDCardTypeP",
	"IDCardNoP",
	"GenderP",
	"BirthdayP",
	"PolicyholderID",
	"OriginalOfficialNameP",
	"OriginalIDCardTypeP",
	"OriginalIDCardNoP",
	"OriginalGenderP",
	"OriginalBirthdayP",
	"Addi",
	"CreateTime",
	"UpdateTime",
	"Creator",
	"Remark",
	"Status",
}

//Fields return all fields of struct.
func (r *TVMistakeCorrect) Fields() []string {
	return TVMistakeCorrectFields
}

//GetTableName return the associated db table name.
func (r *TVMistakeCorrect) GetTableName() string {
	var viewNamePattern = regexp.MustCompile(`(?i)^t_v_[a-z0-9_]+$`)
	tableName := "t_v_mistake_correct"
	if viewNamePattern.MatchString(tableName) {
		return tableName[2:]
	}
	return tableName
}

// Create inserts the TVMistakeCorrect to the database.
func (r *TVMistakeCorrect) Create(db Queryer) error {
	_, err := db.Exec(
		`INSERT INTO t_v_mistake_correct (id, order_id, org_id, commence_date, expiry_date, insurance_type, official_name, id_card_type, id_card_no, gender, birthday, school, school_id, school_type, insured_id, original_official_name, original_id_card_type, original_id_card_no, original_gender, original_birthday, official_name_p, id_card_type_p, id_card_no_p, gender_p, birthday_p, policyholder_id, original_official_name_p, original_id_card_type_p, original_id_card_no_p, original_gender_p, original_birthday_p, addi, create_time, update_time, creator, remark, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36, $37)`,
		&r.ID, &r.OrderID, &r.OrgID, &r.CommenceDate, &r.ExpiryDate, &r.InsuranceType, &r.OfficialName, &r.IDCardType, &r.IDCardNo, &r.Gender, &r.Birthday, &r.School, &r.SchoolID, &r.SchoolType, &r.InsuredID, &r.OriginalOfficialName, &r.OriginalIDCardType, &r.OriginalIDCardNo, &r.OriginalGender, &r.OriginalBirthday, &r.OfficialNameP, &r.IDCardTypeP, &r.IDCardNoP, &r.GenderP, &r.BirthdayP, &r.PolicyholderID, &r.OriginalOfficialNameP, &r.OriginalIDCardTypeP, &r.OriginalIDCardNoP, &r.OriginalGenderP, &r.OriginalBirthdayP, &r.Addi, &r.CreateTime, &r.UpdateTime, &r.Creator, &r.Remark, &r.Status)
	if err != nil {
		return errors.Wrap(err, "failed to insert t_v_mistake_correct")
	}
	return nil
}

// GetTVMistakeCorrectByPk select the TVMistakeCorrect from the database.
func GetTVMistakeCorrectByPk(db Queryer) (*TVMistakeCorrect, error) {
	var r TVMistakeCorrect
	err := db.QueryRow(
		`SELECT id, order_id, org_id, commence_date, expiry_date, insurance_type, official_name, id_card_type, id_card_no, gender, birthday, school, school_id, school_type, insured_id, original_official_name, original_id_card_type, original_id_card_no, original_gender, original_birthday, official_name_p, id_card_type_p, id_card_no_p, gender_p, birthday_p, policyholder_id, original_official_name_p, original_id_card_type_p, original_id_card_no_p, original_gender_p, original_birthday_p, addi, create_time, update_time, creator, remark, status FROM t_v_mistake_correct WHERE `,
	).Scan(&r.ID, &r.OrderID, &r.OrgID, &r.CommenceDate, &r.ExpiryDate, &r.InsuranceType, &r.OfficialName, &r.IDCardType, &r.IDCardNo, &r.Gender, &r.Birthday, &r.School, &r.SchoolID, &r.SchoolType, &r.InsuredID, &r.OriginalOfficialName, &r.OriginalIDCardType, &r.OriginalIDCardNo, &r.OriginalGender, &r.OriginalBirthday, &r.OfficialNameP, &r.IDCardTypeP, &r.IDCardNoP, &r.GenderP, &r.BirthdayP, &r.PolicyholderID, &r.OriginalOfficialNameP, &r.OriginalIDCardTypeP, &r.OriginalIDCardNoP, &r.OriginalGenderP, &r.OriginalBirthdayP, &r.Addi, &r.CreateTime, &r.UpdateTime, &r.Creator, &r.Remark, &r.Status)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select t_v_mistake_correct")
	}
	return &r, nil
}

/*TVMistakeCorrect2 t_v_mistake_correct2 represents kuser.t_v_mistake_correct2 */
type TVMistakeCorrect2 struct {
	ID                              null.Int       `json:"ID,omitempty" db:"id,false,integer"`                                                              /* id id */
	OrderID                         null.Int       `json:"OrderID,omitempty" db:"order_id,false,bigint"`                                                    /* order_id order_id */
	InsuranceTypeID                 null.Int       `json:"InsuranceTypeID,omitempty" db:"insurance_type_id,false,bigint"`                                   /* insurance_type_id insurance_type_id */
	InsuranceTypeParentID           null.Int       `json:"InsuranceTypeParentID,omitempty" db:"insurance_type_parent_id,false,bigint"`                      /* insurance_type_parent_id insurance_type_parent_id */
	OrgID                           null.Int       `json:"OrgID,omitempty" db:"org_id,false,bigint"`                                                        /* org_id org_id */
	HaveDinnerNum                   null.Bool      `json:"HaveDinnerNum,omitempty" db:"have_dinner_num,false,boolean"`                                      /* have_dinner_num have_dinner_num */
	CommenceDate                    null.Int       `json:"CommenceDate,omitempty" db:"commence_date,false,bigint"`                                          /* commence_date commence_date */
	NewCommenceDate                 null.Int       `json:"NewCommenceDate,omitempty" db:"new_commence_date,false,bigint"`                                   /* new_commence_date new_commence_date */
	ExpiryDate                      null.Int       `json:"ExpiryDate,omitempty" db:"expiry_date,false,bigint"`                                              /* expiry_date expiry_date */
	NewExpiryDate                   null.Int       `json:"NewExpiryDate,omitempty" db:"new_expiry_date,false,bigint"`                                       /* new_expiry_date new_expiry_date */
	HaveInsuredList                 null.Bool      `json:"HaveInsuredList,omitempty" db:"have_insured_list,false,boolean"`                                  /* have_insured_list have_insured_list */
	ModifyType                      null.String    `json:"ModifyType,omitempty" db:"modify_type,false,character varying"`                                   /* modify_type modify_type */
	InsuranceType                   null.String    `json:"InsuranceType,omitempty" db:"insurance_type,false,character varying"`                             /* insurance_type insurance_type */
	ActivityCategory                null.String    `json:"ActivityCategory,omitempty" db:"activity_category,false,character varying"`                       /* activity_category activity_category */
	PlanID                          null.Int       `json:"PlanID,omitempty" db:"plan_id,false,bigint"`                                                      /* plan_id plan_id */
	OriginalPlanID                  null.Int       `json:"OriginalPlanID,omitempty" db:"original_plan_id,false,bigint"`                                     /* original_plan_id original_plan_id */
	PlanName                        null.String    `json:"PlanName,omitempty" db:"plan_name,false,character varying"`                                       /* plan_name plan_name */
	OrgName                         null.String    `json:"OrgName,omitempty" db:"org_name,false,text"`                                                      /* org_name org_name */
	OrgAddr                         null.String    `json:"OrgAddr,omitempty" db:"org_addr,false,text"`                                                      /* org_addr org_addr */
	OrgCreditCode                   null.String    `json:"OrgCreditCode,omitempty" db:"org_credit_code,false,text"`                                         /* org_credit_code org_credit_code */
	OrgContact                      null.String    `json:"OrgContact,omitempty" db:"org_contact,false,text"`                                                /* org_contact org_contact */
	OrgPhone                        null.String    `json:"OrgPhone,omitempty" db:"org_phone,false,text"`                                                    /* org_phone org_phone */
	OrgContactRole                  null.String    `json:"OrgContactRole,omitempty" db:"org_contact_role,false,text"`                                       /* org_contact_role org_contact_role */
	OrgCreditCodePic                null.String    `json:"OrgCreditCodePic,omitempty" db:"org_credit_code_pic,false,text"`                                  /* org_credit_code_pic org_credit_code_pic */
	OrgSchoolCategory               null.String    `json:"OrgSchoolCategory,omitempty" db:"org_school_category,false,text"`                                 /* org_school_category org_school_category */
	OrgProvince                     null.String    `json:"OrgProvince,omitempty" db:"org_province,false,text"`                                              /* org_province org_province */
	OrgCity                         null.String    `json:"OrgCity,omitempty" db:"org_city,false,text"`                                                      /* org_city org_city */
	OrgDistrict                     null.String    `json:"OrgDistrict,omitempty" db:"org_district,false,text"`                                              /* org_district org_district */
	InsuredName                     null.String    `json:"InsuredName,omitempty" db:"insured_name,false,text"`                                              /* insured_name insured_name */
	InsuredProvince                 null.String    `json:"InsuredProvince,omitempty" db:"insured_province,false,text"`                                      /* insured_province insured_province */
	InsuredCity                     null.String    `json:"InsuredCity,omitempty" db:"insured_city,false,text"`                                              /* insured_city insured_city */
	InsuredDistrict                 null.String    `json:"InsuredDistrict,omitempty" db:"insured_district,false,text"`                                      /* insured_district insured_district */
	InsuredIsCompulsory             null.Bool      `json:"InsuredIsCompulsory,omitempty" db:"insured_is_compulsory,false,boolean"`                          /* insured_is_compulsory insured_is_compulsory */
	InsuredCategory                 null.String    `json:"InsuredCategory,omitempty" db:"insured_category,false,text"`                                      /* insured_category insured_category */
	InsuredSchoolCategory           null.String    `json:"InsuredSchoolCategory,omitempty" db:"insured_school_category,false,text"`                         /* insured_school_category insured_school_category */
	OfficialName                    null.String    `json:"OfficialName,omitempty" db:"official_name,false,character varying"`                               /* official_name official_name */
	IDCardType                      null.String    `json:"IDCardType,omitempty" db:"id_card_type,false,character varying"`                                  /* id_card_type id_card_type */
	IDCardNo                        null.String    `json:"IDCardNo,omitempty" db:"id_card_no,false,character varying"`                                      /* id_card_no id_card_no */
	Gender                          null.String    `json:"Gender,omitempty" db:"gender,false,character varying"`                                            /* gender gender */
	Birthday                        null.Int       `json:"Birthday,omitempty" db:"birthday,false,bigint"`                                                   /* birthday birthday */
	School                          null.String    `json:"School,omitempty" db:"school,false,character varying"`                                            /* school school */
	SchoolID                        null.Int       `json:"SchoolID,omitempty" db:"school_id,false,bigint"`                                                  /* school_id school_id */
	SchoolType                      null.String    `json:"SchoolType,omitempty" db:"school_type,false,character varying"`                                   /* school_type school_type */
	OriginalOfficialName            null.String    `json:"OriginalOfficialName,omitempty" db:"original_official_name,false,character varying"`              /* original_official_name original_official_name */
	OriginalIDCardType              null.String    `json:"OriginalIDCardType,omitempty" db:"original_id_card_type,false,character varying"`                 /* original_id_card_type original_id_card_type */
	OriginalIDCardNo                null.String    `json:"OriginalIDCardNo,omitempty" db:"original_id_card_no,false,character varying"`                     /* original_id_card_no original_id_card_no */
	OriginalGender                  null.String    `json:"OriginalGender,omitempty" db:"original_gender,false,character varying"`                           /* original_gender original_gender */
	OriginalBirthday                null.Int       `json:"OriginalBirthday,omitempty" db:"original_birthday,false,bigint"`                                  /* original_birthday original_birthday */
	OfficialNameP                   null.String    `json:"OfficialNameP,omitempty" db:"official_name_p,false,character varying"`                            /* official_name_p official_name_p */
	IDCardTypeP                     null.String    `json:"IDCardTypeP,omitempty" db:"id_card_type_p,false,character varying"`                               /* id_card_type_p id_card_type_p */
	IDCardNoP                       null.String    `json:"IDCardNoP,omitempty" db:"id_card_no_p,false,character varying"`                                   /* id_card_no_p id_card_no_p */
	GenderP                         null.String    `json:"GenderP,omitempty" db:"gender_p,false,character varying"`                                         /* gender_p gender_p */
	BirthdayP                       null.Int       `json:"BirthdayP,omitempty" db:"birthday_p,false,bigint"`                                                /* birthday_p birthday_p */
	OriginalOfficialNameP           null.String    `json:"OriginalOfficialNameP,omitempty" db:"original_official_name_p,false,character varying"`           /* original_official_name_p original_official_name_p */
	OriginalIDCardTypeP             null.String    `json:"OriginalIDCardTypeP,omitempty" db:"original_id_card_type_p,false,character varying"`              /* original_id_card_type_p original_id_card_type_p */
	OriginalIDCardNoP               null.String    `json:"OriginalIDCardNoP,omitempty" db:"original_id_card_no_p,false,character varying"`                  /* original_id_card_no_p original_id_card_no_p */
	OriginalGenderP                 null.String    `json:"OriginalGenderP,omitempty" db:"original_gender_p,false,character varying"`                        /* original_gender_p original_gender_p */
	OriginalBirthdayP               null.Int       `json:"OriginalBirthdayP,omitempty" db:"original_birthday_p,false,bigint"`                               /* original_birthday_p original_birthday_p */
	Insurer                         null.String    `json:"Insurer,omitempty" db:"insurer,false,character varying"`                                          /* insurer insurer */
	OriginalActivityName            null.String    `json:"OriginalActivityName,omitempty" db:"original_activity_name,false,character varying"`              /* original_activity_name original_activity_name */
	ActivityName                    null.String    `json:"ActivityName,omitempty" db:"activity_name,false,character varying"`                               /* activity_name activity_name */
	OriginalDesc                    null.String    `json:"OriginalDesc,omitempty" db:"original_desc,false,character varying"`                               /* original_desc original_desc */
	ActivityDesc                    null.String    `json:"ActivityDesc,omitempty" db:"activity_desc,false,character varying"`                               /* activity_desc activity_desc */
	OriginalActivityLocation        null.String    `json:"OriginalActivityLocation,omitempty" db:"original_activity_location,false,character varying"`      /* original_activity_location original_activity_location */
	ActivityLocation                null.String    `json:"ActivityLocation,omitempty" db:"activity_location,false,character varying"`                       /* activity_location activity_location */
	OriginalActivityDateSet         null.String    `json:"OriginalActivityDateSet,omitempty" db:"original_activity_date_set,false,character varying"`       /* original_activity_date_set original_activity_date_set */
	ActivityDateSet                 null.String    `json:"ActivityDateSet,omitempty" db:"activity_date_set,false,character varying"`                        /* activity_date_set activity_date_set */
	OriginalIndate                  null.Int       `json:"OriginalIndate,omitempty" db:"original_indate,false,bigint"`                                      /* original_indate original_indate */
	Indate                          null.Int       `json:"Indate,omitempty" db:"indate,false,bigint"`                                                       /* indate indate */
	OriginalPolicyholder            types.JSONText `json:"OriginalPolicyholder,omitempty" db:"original_policyholder,false,jsonb"`                           /* original_policyholder original_policyholder */
	Policyholder                    types.JSONText `json:"Policyholder,omitempty" db:"policyholder,false,jsonb"`                                            /* policyholder policyholder */
	PolicyholderID                  null.Int       `json:"PolicyholderID,omitempty" db:"policyholder_id,false,bigint"`                                      /* policyholder_id policyholder_id */
	OriginalInsured                 types.JSONText `json:"OriginalInsured,omitempty" db:"original_insured,false,jsonb"`                                     /* original_insured original_insured */
	Insured                         types.JSONText `json:"Insured,omitempty" db:"insured,false,jsonb"`                                                      /* insured insured */
	InsuredID                       null.Int       `json:"InsuredID,omitempty" db:"insured_id,false,bigint"`                                                /* insured_id insured_id */
	OriginalInsuredGroupByDay       null.Bool      `json:"OriginalInsuredGroupByDay,omitempty" db:"original_insured_group_by_day,false,boolean"`            /* original_insured_group_by_day original_insured_group_by_day */
	InsuredGroupByDay               null.Bool      `json:"InsuredGroupByDay,omitempty" db:"insured_group_by_day,false,boolean"`                             /* insured_group_by_day insured_group_by_day */
	OriginalChargeMode              null.String    `json:"OriginalChargeMode,omitempty" db:"original_charge_mode,false,character varying"`                  /* original_charge_mode original_charge_mode */
	ChargeMode                      null.String    `json:"ChargeMode,omitempty" db:"charge_mode,false,character varying"`                                   /* charge_mode charge_mode */
	OriginalAmount                  null.Float     `json:"OriginalAmount,omitempty" db:"original_amount,false,double precision"`                            /* original_amount original_amount */
	Amount                          null.Float     `json:"Amount,omitempty" db:"amount,false,double precision"`                                             /* amount amount */
	OriginalInsuredCount            null.Int       `json:"OriginalInsuredCount,omitempty" db:"original_insured_count,false,smallint"`                       /* original_insured_count original_insured_count */
	InsuredCount                    null.Int       `json:"InsuredCount,omitempty" db:"insured_count,false,smallint"`                                        /* insured_count insured_count */
	OriginalInsuredType             null.String    `json:"OriginalInsuredType,omitempty" db:"original_insured_type,false,character varying"`                /* original_insured_type original_insured_type */
	InsuredType                     null.String    `json:"InsuredType,omitempty" db:"insured_type,false,character varying"`                                 /* insured_type insured_type */
	OriginalInsuredList             types.JSONText `json:"OriginalInsuredList,omitempty" db:"original_insured_list,false,jsonb"`                            /* original_insured_list original_insured_list */
	InsuredList                     types.JSONText `json:"InsuredList,omitempty" db:"insured_list,false,jsonb"`                                             /* insured_list insured_list */
	InsertInsuredList               types.JSONText `json:"InsertInsuredList,omitempty" db:"insert_insured_list,false,jsonb"`                                /* insert_insured_list insert_insured_list */
	DeleteInsuredList               types.JSONText `json:"DeleteInsuredList,omitempty" db:"delete_insured_list,false,jsonb"`                                /* delete_insured_list delete_insured_list */
	UpdateInsuredList               types.JSONText `json:"UpdateInsuredList,omitempty" db:"update_insured_list,false,jsonb"`                                /* update_insured_list update_insured_list */
	RequireUpdateInsuredList        types.JSONText `json:"RequireUpdateInsuredList,omitempty" db:"require_update_insured_list,false,jsonb"`                 /* require_update_insured_list require_update_insured_list */
	OriginalNonCompulsoryStudentNum null.Int       `json:"OriginalNonCompulsoryStudentNum,omitempty" db:"original_non_compulsory_student_num,false,bigint"` /* original_non_compulsory_student_num original_non_compulsory_student_num */
	NonCompulsoryStudentNum         null.Int       `json:"NonCompulsoryStudentNum,omitempty" db:"non_compulsory_student_num,false,bigint"`                  /* non_compulsory_student_num non_compulsory_student_num */
	OriginalCompulsoryStudentNum    null.Int       `json:"OriginalCompulsoryStudentNum,omitempty" db:"original_compulsory_student_num,false,bigint"`        /* original_compulsory_student_num original_compulsory_student_num */
	CompulsoryStudentNum            null.Int       `json:"CompulsoryStudentNum,omitempty" db:"compulsory_student_num,false,bigint"`                         /* compulsory_student_num compulsory_student_num */
	OriginalCanteenNum              null.Int       `json:"OriginalCanteenNum,omitempty" db:"original_canteen_num,false,integer"`                            /* original_canteen_num original_canteen_num */
	CanteenNum                      null.Int       `json:"CanteenNum,omitempty" db:"canteen_num,false,bigint"`                                              /* canteen_num canteen_num */
	OriginalShopNum                 null.Int       `json:"OriginalShopNum,omitempty" db:"original_shop_num,false,integer"`                                  /* original_shop_num original_shop_num */
	ShopNum                         null.Int       `json:"ShopNum,omitempty" db:"shop_num,false,bigint"`                                                    /* shop_num shop_num */
	OriginalDinnerNum               null.Int       `json:"OriginalDinnerNum,omitempty" db:"original_dinner_num,false,integer"`                              /* original_dinner_num original_dinner_num */
	DinnerNum                       null.Int       `json:"DinnerNum,omitempty" db:"dinner_num,false,bigint"`                                                /* dinner_num dinner_num */
	OriginalPayType                 null.String    `json:"OriginalPayType,omitempty" db:"original_pay_type,false,character varying"`                        /* original_pay_type original_pay_type */
	PayType                         null.String    `json:"PayType,omitempty" db:"pay_type,false,character varying"`                                         /* pay_type pay_type */
	OriginalFeeScheme               types.JSONText `json:"OriginalFeeScheme,omitempty" db:"original_fee_scheme,false,jsonb"`                                /* original_fee_scheme original_fee_scheme */
	FeeScheme                       types.JSONText `json:"FeeScheme,omitempty" db:"fee_scheme,false,jsonb"`                                                 /* fee_scheme fee_scheme */
	NeedBalance                     null.Bool      `json:"NeedBalance,omitempty" db:"need_balance,false,boolean"`                                           /* need_balance need_balance */
	OrderStatus                     null.String    `json:"OrderStatus,omitempty" db:"order_status,false,character varying"`                                 /* order_status order_status */
	OriginalDisputeHandling         null.String    `json:"OriginalDisputeHandling,omitempty" db:"original_dispute_handling,false,character varying"`        /* original_dispute_handling original_dispute_handling */
	DisputeHandling                 null.String    `json:"DisputeHandling,omitempty" db:"dispute_handling,false,character varying"`                         /* dispute_handling dispute_handling */
	OriginalHaveSuddenDeath         null.Bool      `json:"OriginalHaveSuddenDeath,omitempty" db:"original_have_sudden_death,false,boolean"`                 /* original_have_sudden_death original_have_sudden_death */
	HaveSuddenDeath                 null.Bool      `json:"HaveSuddenDeath,omitempty" db:"have_sudden_death,false,boolean"`                                  /* have_sudden_death have_sudden_death */
	OriginalPrevPolicyNo            null.String    `json:"OriginalPrevPolicyNo,omitempty" db:"original_prev_policy_no,false,character varying"`             /* original_prev_policy_no original_prev_policy_no */
	RevokedPolicyNo                 null.String    `json:"RevokedPolicyNo,omitempty" db:"revoked_policy_no,false,character varying"`                        /* revoked_policy_no revoked_policy_no */
	PrevPolicyNo                    null.String    `json:"PrevPolicyNo,omitempty" db:"prev_policy_no,false,character varying"`                              /* prev_policy_no prev_policy_no */
	OriginalPoolName                null.String    `json:"OriginalPoolName,omitempty" db:"original_pool_name,false,character varying"`                      /* original_pool_name original_pool_name */
	PoolName                        null.String    `json:"PoolName,omitempty" db:"pool_name,false,character varying"`                                       /* pool_name pool_name */
	OriginalHaveExplosive           null.Bool      `json:"OriginalHaveExplosive,omitempty" db:"original_have_explosive,false,boolean"`                      /* original_have_explosive original_have_explosive */
	HaveExplosive                   null.Bool      `json:"HaveExplosive,omitempty" db:"have_explosive,false,boolean"`                                       /* have_explosive have_explosive */
	OriginalHaveRides               null.Bool      `json:"OriginalHaveRides,omitempty" db:"original_have_rides,false,boolean"`                              /* original_have_rides original_have_rides */
	HaveRides                       null.Bool      `json:"HaveRides,omitempty" db:"have_rides,false,boolean"`                                               /* have_rides have_rides */
	OriginalInnerArea               null.Float     `json:"OriginalInnerArea,omitempty" db:"original_inner_area,false,double precision"`                     /* original_inner_area original_inner_area */
	InnerArea                       null.Float     `json:"InnerArea,omitempty" db:"inner_area,false,double precision"`                                      /* inner_area inner_area */
	OriginalOuterArea               null.Float     `json:"OriginalOuterArea,omitempty" db:"original_outer_area,false,double precision"`                     /* original_outer_area original_outer_area */
	OuterArea                       null.Float     `json:"OuterArea,omitempty" db:"outer_area,false,double precision"`                                      /* outer_area outer_area */
	OriginalTrafficNum              null.Int       `json:"OriginalTrafficNum,omitempty" db:"original_traffic_num,false,integer"`                            /* original_traffic_num original_traffic_num */
	TrafficNum                      null.Int       `json:"TrafficNum,omitempty" db:"traffic_num,false,integer"`                                             /* traffic_num traffic_num */
	OriginalTemperatureType         null.String    `json:"OriginalTemperatureType,omitempty" db:"original_temperature_type,false,character varying"`        /* original_temperature_type original_temperature_type */
	TemperatureType                 null.String    `json:"TemperatureType,omitempty" db:"temperature_type,false,character varying"`                         /* temperature_type temperature_type */
	OriginalOpenPoolNum             null.Int       `json:"OriginalOpenPoolNum,omitempty" db:"original_open_pool_num,false,smallint"`                        /* original_open_pool_num original_open_pool_num */
	OpenPoolNum                     null.Int       `json:"OpenPoolNum,omitempty" db:"open_pool_num,false,smallint"`                                         /* open_pool_num open_pool_num */
	OriginalHeatedPoolNum           null.Int       `json:"OriginalHeatedPoolNum,omitempty" db:"original_heated_pool_num,false,smallint"`                    /* original_heated_pool_num original_heated_pool_num */
	HeatedPoolNum                   null.Int       `json:"HeatedPoolNum,omitempty" db:"heated_pool_num,false,smallint"`                                     /* heated_pool_num heated_pool_num */
	OriginalTrainingPoolNum         null.Int       `json:"OriginalTrainingPoolNum,omitempty" db:"original_training_pool_num,false,smallint"`                /* original_training_pool_num original_training_pool_num */
	TrainingPoolNum                 null.Int       `json:"TrainingPoolNum,omitempty" db:"training_pool_num,false,smallint"`                                 /* training_pool_num training_pool_num */
	OriginalPoolNum                 null.Int       `json:"OriginalPoolNum,omitempty" db:"original_pool_num,false,smallint"`                                 /* original_pool_num original_pool_num */
	PoolNum                         null.Int       `json:"PoolNum,omitempty" db:"pool_num,false,smallint"`                                                  /* pool_num pool_num */
	OriginalCustomType              null.String    `json:"OriginalCustomType,omitempty" db:"original_custom_type,false,character varying"`                  /* original_custom_type original_custom_type */
	CustomType                      null.String    `json:"CustomType,omitempty" db:"custom_type,false,character varying"`                                   /* custom_type custom_type */
	OriginalSame                    null.Bool      `json:"OriginalSame,omitempty" db:"original_same,false,boolean"`                                         /* original_same original_same */
	Same                            null.Bool      `json:"Same,omitempty" db:"same,false,boolean"`                                                          /* same same */
	OriginalArbitralAgency          null.String    `json:"OriginalArbitralAgency,omitempty" db:"original_arbitral_agency,false,character varying"`          /* original_arbitral_agency original_arbitral_agency */
	ArbitralAgency                  null.String    `json:"ArbitralAgency,omitempty" db:"arbitral_agency,false,character varying"`                           /* arbitral_agency arbitral_agency */
	EndorsementStatus               null.String    `json:"EndorsementStatus,omitempty" db:"endorsement_status,false,character varying"`                     /* endorsement_status endorsement_status */
	ApplicationFiles                types.JSONText `json:"ApplicationFiles,omitempty" db:"application_files,false,jsonb"`                                   /* application_files application_files */
	Balance                         null.Float     `json:"Balance,omitempty" db:"balance,false,double precision"`                                           /* balance balance */
	BalanceList                     types.JSONText `json:"BalanceList,omitempty" db:"balance_list,false,jsonb"`                                             /* balance_list balance_list */
	HaveNegotiatedPrice             null.Bool      `json:"HaveNegotiatedPrice,omitempty" db:"have_negotiated_price,false,boolean"`                          /* have_negotiated_price have_negotiated_price */
	Sn                              null.String    `json:"Sn,omitempty" db:"sn,false,text"`                                                                 /* sn sn */
	PolicyRegen                     null.Bool      `json:"PolicyRegen,omitempty" db:"policy_regen,false,boolean"`                                           /* policy_regen policy_regen */
	ClearList                       null.Bool      `json:"ClearList,omitempty" db:"clear_list,false,boolean"`                                               /* clear_list clear_list */
	FilesToRemove                   null.String    `json:"FilesToRemove,omitempty" db:"files_to_remove,false,character varying"`                            /* files_to_remove files_to_remove */
	OriginalPolicyScheme            types.JSONText `json:"OriginalPolicyScheme,omitempty" db:"original_policy_scheme,false,jsonb"`                          /* original_policy_scheme original_policy_scheme */
	PolicyScheme                    types.JSONText `json:"PolicyScheme,omitempty" db:"policy_scheme,false,jsonb"`                                           /* policy_scheme policy_scheme */
	InvoiceHeader                   null.String    `json:"InvoiceHeader,omitempty" db:"invoice_header,false,character varying"`                             /* invoice_header invoice_header */
	CorrectLevel                    null.String    `json:"CorrectLevel,omitempty" db:"correct_level,false,character varying"`                               /* correct_level correct_level */
	CorrectLog                      types.JSONText `json:"CorrectLog,omitempty" db:"correct_log,false,jsonb"`                                               /* correct_log correct_log */
	OriginalFiles                   types.JSONText `json:"OriginalFiles,omitempty" db:"original_files,false,jsonb"`                                         /* original_files original_files */
	Files                           types.JSONText `json:"Files,omitempty" db:"files,false,jsonb"`                                                          /* files files */
	RefusedReason                   null.String    `json:"RefusedReason,omitempty" db:"refused_reason,false,character varying"`                             /* refused_reason refused_reason */
	Addi                            types.JSONText `json:"Addi,omitempty" db:"addi,false,jsonb"`                                                            /* addi addi */
	CreateTime                      null.Int       `json:"CreateTime,omitempty" db:"create_time,false,bigint"`                                              /* create_time create_time */
	UpdateTime                      null.Int       `json:"UpdateTime,omitempty" db:"update_time,false,bigint"`                                              /* update_time update_time */
	Creator                         null.Int       `json:"Creator,omitempty" db:"creator,false,bigint"`                                                     /* creator creator */
	Remark                          null.String    `json:"Remark,omitempty" db:"remark,false,character varying"`                                            /* remark remark */
	Status                          null.String    `json:"Status,omitempty" db:"status,false,character varying"`                                            /* status status */
	Filter                                         // build DML where clause
}

//TVMistakeCorrect2Fields full field list for default query
var TVMistakeCorrect2Fields = []string{
	"ID",
	"OrderID",
	"InsuranceTypeID",
	"InsuranceTypeParentID",
	"OrgID",
	"HaveDinnerNum",
	"CommenceDate",
	"NewCommenceDate",
	"ExpiryDate",
	"NewExpiryDate",
	"HaveInsuredList",
	"ModifyType",
	"InsuranceType",
	"ActivityCategory",
	"PlanID",
	"OriginalPlanID",
	"PlanName",
	"OrgName",
	"OrgAddr",
	"OrgCreditCode",
	"OrgContact",
	"OrgPhone",
	"OrgContactRole",
	"OrgCreditCodePic",
	"OrgSchoolCategory",
	"OrgProvince",
	"OrgCity",
	"OrgDistrict",
	"InsuredName",
	"InsuredProvince",
	"InsuredCity",
	"InsuredDistrict",
	"InsuredIsCompulsory",
	"InsuredCategory",
	"InsuredSchoolCategory",
	"OfficialName",
	"IDCardType",
	"IDCardNo",
	"Gender",
	"Birthday",
	"School",
	"SchoolID",
	"SchoolType",
	"OriginalOfficialName",
	"OriginalIDCardType",
	"OriginalIDCardNo",
	"OriginalGender",
	"OriginalBirthday",
	"OfficialNameP",
	"IDCardTypeP",
	"IDCardNoP",
	"GenderP",
	"BirthdayP",
	"OriginalOfficialNameP",
	"OriginalIDCardTypeP",
	"OriginalIDCardNoP",
	"OriginalGenderP",
	"OriginalBirthdayP",
	"Insurer",
	"OriginalActivityName",
	"ActivityName",
	"OriginalDesc",
	"ActivityDesc",
	"OriginalActivityLocation",
	"ActivityLocation",
	"OriginalActivityDateSet",
	"ActivityDateSet",
	"OriginalIndate",
	"Indate",
	"OriginalPolicyholder",
	"Policyholder",
	"PolicyholderID",
	"OriginalInsured",
	"Insured",
	"InsuredID",
	"OriginalInsuredGroupByDay",
	"InsuredGroupByDay",
	"OriginalChargeMode",
	"ChargeMode",
	"OriginalAmount",
	"Amount",
	"OriginalInsuredCount",
	"InsuredCount",
	"OriginalInsuredType",
	"InsuredType",
	"OriginalInsuredList",
	"InsuredList",
	"InsertInsuredList",
	"DeleteInsuredList",
	"UpdateInsuredList",
	"RequireUpdateInsuredList",
	"OriginalNonCompulsoryStudentNum",
	"NonCompulsoryStudentNum",
	"OriginalCompulsoryStudentNum",
	"CompulsoryStudentNum",
	"OriginalCanteenNum",
	"CanteenNum",
	"OriginalShopNum",
	"ShopNum",
	"OriginalDinnerNum",
	"DinnerNum",
	"OriginalPayType",
	"PayType",
	"OriginalFeeScheme",
	"FeeScheme",
	"NeedBalance",
	"OrderStatus",
	"OriginalDisputeHandling",
	"DisputeHandling",
	"OriginalHaveSuddenDeath",
	"HaveSuddenDeath",
	"OriginalPrevPolicyNo",
	"RevokedPolicyNo",
	"PrevPolicyNo",
	"OriginalPoolName",
	"PoolName",
	"OriginalHaveExplosive",
	"HaveExplosive",
	"OriginalHaveRides",
	"HaveRides",
	"OriginalInnerArea",
	"InnerArea",
	"OriginalOuterArea",
	"OuterArea",
	"OriginalTrafficNum",
	"TrafficNum",
	"OriginalTemperatureType",
	"TemperatureType",
	"OriginalOpenPoolNum",
	"OpenPoolNum",
	"OriginalHeatedPoolNum",
	"HeatedPoolNum",
	"OriginalTrainingPoolNum",
	"TrainingPoolNum",
	"OriginalPoolNum",
	"PoolNum",
	"OriginalCustomType",
	"CustomType",
	"OriginalSame",
	"Same",
	"OriginalArbitralAgency",
	"ArbitralAgency",
	"EndorsementStatus",
	"ApplicationFiles",
	"Balance",
	"BalanceList",
	"HaveNegotiatedPrice",
	"Sn",
	"PolicyRegen",
	"ClearList",
	"FilesToRemove",
	"OriginalPolicyScheme",
	"PolicyScheme",
	"InvoiceHeader",
	"CorrectLevel",
	"CorrectLog",
	"OriginalFiles",
	"Files",
	"RefusedReason",
	"Addi",
	"CreateTime",
	"UpdateTime",
	"Creator",
	"Remark",
	"Status",
}

//Fields return all fields of struct.
func (r *TVMistakeCorrect2) Fields() []string {
	return TVMistakeCorrect2Fields
}

//GetTableName return the associated db table name.
func (r *TVMistakeCorrect2) GetTableName() string {
	var viewNamePattern = regexp.MustCompile(`(?i)^t_v_[a-z0-9_]+$`)
	tableName := "t_v_mistake_correct2"
	if viewNamePattern.MatchString(tableName) {
		return tableName[2:]
	}
	return tableName
}

// Create inserts the TVMistakeCorrect2 to the database.
func (r *TVMistakeCorrect2) Create(db Queryer) error {
	_, err := db.Exec(
		`INSERT INTO t_v_mistake_correct2 (id, order_id, insurance_type_id, insurance_type_parent_id, org_id, have_dinner_num, commence_date, new_commence_date, expiry_date, new_expiry_date, have_insured_list, modify_type, insurance_type, activity_category, plan_id, original_plan_id, plan_name, org_name, org_addr, org_credit_code, org_contact, org_phone, org_contact_role, org_credit_code_pic, org_school_category, org_province, org_city, org_district, insured_name, insured_province, insured_city, insured_district, insured_is_compulsory, insured_category, insured_school_category, official_name, id_card_type, id_card_no, gender, birthday, school, school_id, school_type, original_official_name, original_id_card_type, original_id_card_no, original_gender, original_birthday, official_name_p, id_card_type_p, id_card_no_p, gender_p, birthday_p, original_official_name_p, original_id_card_type_p, original_id_card_no_p, original_gender_p, original_birthday_p, insurer, original_activity_name, activity_name, original_desc, activity_desc, original_activity_location, activity_location, original_activity_date_set, activity_date_set, original_indate, indate, original_policyholder, policyholder, policyholder_id, original_insured, insured, insured_id, original_insured_group_by_day, insured_group_by_day, original_charge_mode, charge_mode, original_amount, amount, original_insured_count, insured_count, original_insured_type, insured_type, original_insured_list, insured_list, insert_insured_list, delete_insured_list, update_insured_list, require_update_insured_list, original_non_compulsory_student_num, non_compulsory_student_num, original_compulsory_student_num, compulsory_student_num, original_canteen_num, canteen_num, original_shop_num, shop_num, original_dinner_num, dinner_num, original_pay_type, pay_type, original_fee_scheme, fee_scheme, need_balance, order_status, original_dispute_handling, dispute_handling, original_have_sudden_death, have_sudden_death, original_prev_policy_no, revoked_policy_no, prev_policy_no, original_pool_name, pool_name, original_have_explosive, have_explosive, original_have_rides, have_rides, original_inner_area, inner_area, original_outer_area, outer_area, original_traffic_num, traffic_num, original_temperature_type, temperature_type, original_open_pool_num, open_pool_num, original_heated_pool_num, heated_pool_num, original_training_pool_num, training_pool_num, original_pool_num, pool_num, original_custom_type, custom_type, original_same, same, original_arbitral_agency, arbitral_agency, endorsement_status, application_files, balance, balance_list, have_negotiated_price, sn, policy_regen, clear_list, files_to_remove, original_policy_scheme, policy_scheme, invoice_header, correct_level, correct_log, original_files, files, refused_reason, addi, create_time, update_time, creator, remark, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36, $37, $38, $39, $40, $41, $42, $43, $44, $45, $46, $47, $48, $49, $50, $51, $52, $53, $54, $55, $56, $57, $58, $59, $60, $61, $62, $63, $64, $65, $66, $67, $68, $69, $70, $71, $72, $73, $74, $75, $76, $77, $78, $79, $80, $81, $82, $83, $84, $85, $86, $87, $88, $89, $90, $91, $92, $93, $94, $95, $96, $97, $98, $99, $100, $101, $102, $103, $104, $105, $106, $107, $108, $109, $110, $111, $112, $113, $114, $115, $116, $117, $118, $119, $120, $121, $122, $123, $124, $125, $126, $127, $128, $129, $130, $131, $132, $133, $134, $135, $136, $137, $138, $139, $140, $141, $142, $143, $144, $145, $146, $147, $148, $149, $150, $151, $152, $153, $154, $155, $156, $157, $158, $159, $160, $161, $162, $163, $164, $165)`,
		&r.ID, &r.OrderID, &r.InsuranceTypeID, &r.InsuranceTypeParentID, &r.OrgID, &r.HaveDinnerNum, &r.CommenceDate, &r.NewCommenceDate, &r.ExpiryDate, &r.NewExpiryDate, &r.HaveInsuredList, &r.ModifyType, &r.InsuranceType, &r.ActivityCategory, &r.PlanID, &r.OriginalPlanID, &r.PlanName, &r.OrgName, &r.OrgAddr, &r.OrgCreditCode, &r.OrgContact, &r.OrgPhone, &r.OrgContactRole, &r.OrgCreditCodePic, &r.OrgSchoolCategory, &r.OrgProvince, &r.OrgCity, &r.OrgDistrict, &r.InsuredName, &r.InsuredProvince, &r.InsuredCity, &r.InsuredDistrict, &r.InsuredIsCompulsory, &r.InsuredCategory, &r.InsuredSchoolCategory, &r.OfficialName, &r.IDCardType, &r.IDCardNo, &r.Gender, &r.Birthday, &r.School, &r.SchoolID, &r.SchoolType, &r.OriginalOfficialName, &r.OriginalIDCardType, &r.OriginalIDCardNo, &r.OriginalGender, &r.OriginalBirthday, &r.OfficialNameP, &r.IDCardTypeP, &r.IDCardNoP, &r.GenderP, &r.BirthdayP, &r.OriginalOfficialNameP, &r.OriginalIDCardTypeP, &r.OriginalIDCardNoP, &r.OriginalGenderP, &r.OriginalBirthdayP, &r.Insurer, &r.OriginalActivityName, &r.ActivityName, &r.OriginalDesc, &r.ActivityDesc, &r.OriginalActivityLocation, &r.ActivityLocation, &r.OriginalActivityDateSet, &r.ActivityDateSet, &r.OriginalIndate, &r.Indate, &r.OriginalPolicyholder, &r.Policyholder, &r.PolicyholderID, &r.OriginalInsured, &r.Insured, &r.InsuredID, &r.OriginalInsuredGroupByDay, &r.InsuredGroupByDay, &r.OriginalChargeMode, &r.ChargeMode, &r.OriginalAmount, &r.Amount, &r.OriginalInsuredCount, &r.InsuredCount, &r.OriginalInsuredType, &r.InsuredType, &r.OriginalInsuredList, &r.InsuredList, &r.InsertInsuredList, &r.DeleteInsuredList, &r.UpdateInsuredList, &r.RequireUpdateInsuredList, &r.OriginalNonCompulsoryStudentNum, &r.NonCompulsoryStudentNum, &r.OriginalCompulsoryStudentNum, &r.CompulsoryStudentNum, &r.OriginalCanteenNum, &r.CanteenNum, &r.OriginalShopNum, &r.ShopNum, &r.OriginalDinnerNum, &r.DinnerNum, &r.OriginalPayType, &r.PayType, &r.OriginalFeeScheme, &r.FeeScheme, &r.NeedBalance, &r.OrderStatus, &r.OriginalDisputeHandling, &r.DisputeHandling, &r.OriginalHaveSuddenDeath, &r.HaveSuddenDeath, &r.OriginalPrevPolicyNo, &r.RevokedPolicyNo, &r.PrevPolicyNo, &r.OriginalPoolName, &r.PoolName, &r.OriginalHaveExplosive, &r.HaveExplosive, &r.OriginalHaveRides, &r.HaveRides, &r.OriginalInnerArea, &r.InnerArea, &r.OriginalOuterArea, &r.OuterArea, &r.OriginalTrafficNum, &r.TrafficNum, &r.OriginalTemperatureType, &r.TemperatureType, &r.OriginalOpenPoolNum, &r.OpenPoolNum, &r.OriginalHeatedPoolNum, &r.HeatedPoolNum, &r.OriginalTrainingPoolNum, &r.TrainingPoolNum, &r.OriginalPoolNum, &r.PoolNum, &r.OriginalCustomType, &r.CustomType, &r.OriginalSame, &r.Same, &r.OriginalArbitralAgency, &r.ArbitralAgency, &r.EndorsementStatus, &r.ApplicationFiles, &r.Balance, &r.BalanceList, &r.HaveNegotiatedPrice, &r.Sn, &r.PolicyRegen, &r.ClearList, &r.FilesToRemove, &r.OriginalPolicyScheme, &r.PolicyScheme, &r.InvoiceHeader, &r.CorrectLevel, &r.CorrectLog, &r.OriginalFiles, &r.Files, &r.RefusedReason, &r.Addi, &r.CreateTime, &r.UpdateTime, &r.Creator, &r.Remark, &r.Status)
	if err != nil {
		return errors.Wrap(err, "failed to insert t_v_mistake_correct2")
	}
	return nil
}

// GetTVMistakeCorrect2ByPk select the TVMistakeCorrect2 from the database.
func GetTVMistakeCorrect2ByPk(db Queryer) (*TVMistakeCorrect2, error) {
	var r TVMistakeCorrect2
	err := db.QueryRow(
		`SELECT id, order_id, insurance_type_id, insurance_type_parent_id, org_id, have_dinner_num, commence_date, new_commence_date, expiry_date, new_expiry_date, have_insured_list, modify_type, insurance_type, activity_category, plan_id, original_plan_id, plan_name, org_name, org_addr, org_credit_code, org_contact, org_phone, org_contact_role, org_credit_code_pic, org_school_category, org_province, org_city, org_district, insured_name, insured_province, insured_city, insured_district, insured_is_compulsory, insured_category, insured_school_category, official_name, id_card_type, id_card_no, gender, birthday, school, school_id, school_type, original_official_name, original_id_card_type, original_id_card_no, original_gender, original_birthday, official_name_p, id_card_type_p, id_card_no_p, gender_p, birthday_p, original_official_name_p, original_id_card_type_p, original_id_card_no_p, original_gender_p, original_birthday_p, insurer, original_activity_name, activity_name, original_desc, activity_desc, original_activity_location, activity_location, original_activity_date_set, activity_date_set, original_indate, indate, original_policyholder, policyholder, policyholder_id, original_insured, insured, insured_id, original_insured_group_by_day, insured_group_by_day, original_charge_mode, charge_mode, original_amount, amount, original_insured_count, insured_count, original_insured_type, insured_type, original_insured_list, insured_list, insert_insured_list, delete_insured_list, update_insured_list, require_update_insured_list, original_non_compulsory_student_num, non_compulsory_student_num, original_compulsory_student_num, compulsory_student_num, original_canteen_num, canteen_num, original_shop_num, shop_num, original_dinner_num, dinner_num, original_pay_type, pay_type, original_fee_scheme, fee_scheme, need_balance, order_status, original_dispute_handling, dispute_handling, original_have_sudden_death, have_sudden_death, original_prev_policy_no, revoked_policy_no, prev_policy_no, original_pool_name, pool_name, original_have_explosive, have_explosive, original_have_rides, have_rides, original_inner_area, inner_area, original_outer_area, outer_area, original_traffic_num, traffic_num, original_temperature_type, temperature_type, original_open_pool_num, open_pool_num, original_heated_pool_num, heated_pool_num, original_training_pool_num, training_pool_num, original_pool_num, pool_num, original_custom_type, custom_type, original_same, same, original_arbitral_agency, arbitral_agency, endorsement_status, application_files, balance, balance_list, have_negotiated_price, sn, policy_regen, clear_list, files_to_remove, original_policy_scheme, policy_scheme, invoice_header, correct_level, correct_log, original_files, files, refused_reason, addi, create_time, update_time, creator, remark, status FROM t_v_mistake_correct2 WHERE `,
	).Scan(&r.ID, &r.OrderID, &r.InsuranceTypeID, &r.InsuranceTypeParentID, &r.OrgID, &r.HaveDinnerNum, &r.CommenceDate, &r.NewCommenceDate, &r.ExpiryDate, &r.NewExpiryDate, &r.HaveInsuredList, &r.ModifyType, &r.InsuranceType, &r.ActivityCategory, &r.PlanID, &r.OriginalPlanID, &r.PlanName, &r.OrgName, &r.OrgAddr, &r.OrgCreditCode, &r.OrgContact, &r.OrgPhone, &r.OrgContactRole, &r.OrgCreditCodePic, &r.OrgSchoolCategory, &r.OrgProvince, &r.OrgCity, &r.OrgDistrict, &r.InsuredName, &r.InsuredProvince, &r.InsuredCity, &r.InsuredDistrict, &r.InsuredIsCompulsory, &r.InsuredCategory, &r.InsuredSchoolCategory, &r.OfficialName, &r.IDCardType, &r.IDCardNo, &r.Gender, &r.Birthday, &r.School, &r.SchoolID, &r.SchoolType, &r.OriginalOfficialName, &r.OriginalIDCardType, &r.OriginalIDCardNo, &r.OriginalGender, &r.OriginalBirthday, &r.OfficialNameP, &r.IDCardTypeP, &r.IDCardNoP, &r.GenderP, &r.BirthdayP, &r.OriginalOfficialNameP, &r.OriginalIDCardTypeP, &r.OriginalIDCardNoP, &r.OriginalGenderP, &r.OriginalBirthdayP, &r.Insurer, &r.OriginalActivityName, &r.ActivityName, &r.OriginalDesc, &r.ActivityDesc, &r.OriginalActivityLocation, &r.ActivityLocation, &r.OriginalActivityDateSet, &r.ActivityDateSet, &r.OriginalIndate, &r.Indate, &r.OriginalPolicyholder, &r.Policyholder, &r.PolicyholderID, &r.OriginalInsured, &r.Insured, &r.InsuredID, &r.OriginalInsuredGroupByDay, &r.InsuredGroupByDay, &r.OriginalChargeMode, &r.ChargeMode, &r.OriginalAmount, &r.Amount, &r.OriginalInsuredCount, &r.InsuredCount, &r.OriginalInsuredType, &r.InsuredType, &r.OriginalInsuredList, &r.InsuredList, &r.InsertInsuredList, &r.DeleteInsuredList, &r.UpdateInsuredList, &r.RequireUpdateInsuredList, &r.OriginalNonCompulsoryStudentNum, &r.NonCompulsoryStudentNum, &r.OriginalCompulsoryStudentNum, &r.CompulsoryStudentNum, &r.OriginalCanteenNum, &r.CanteenNum, &r.OriginalShopNum, &r.ShopNum, &r.OriginalDinnerNum, &r.DinnerNum, &r.OriginalPayType, &r.PayType, &r.OriginalFeeScheme, &r.FeeScheme, &r.NeedBalance, &r.OrderStatus, &r.OriginalDisputeHandling, &r.DisputeHandling, &r.OriginalHaveSuddenDeath, &r.HaveSuddenDeath, &r.OriginalPrevPolicyNo, &r.RevokedPolicyNo, &r.PrevPolicyNo, &r.OriginalPoolName, &r.PoolName, &r.OriginalHaveExplosive, &r.HaveExplosive, &r.OriginalHaveRides, &r.HaveRides, &r.OriginalInnerArea, &r.InnerArea, &r.OriginalOuterArea, &r.OuterArea, &r.OriginalTrafficNum, &r.TrafficNum, &r.OriginalTemperatureType, &r.TemperatureType, &r.OriginalOpenPoolNum, &r.OpenPoolNum, &r.OriginalHeatedPoolNum, &r.HeatedPoolNum, &r.OriginalTrainingPoolNum, &r.TrainingPoolNum, &r.OriginalPoolNum, &r.PoolNum, &r.OriginalCustomType, &r.CustomType, &r.OriginalSame, &r.Same, &r.OriginalArbitralAgency, &r.ArbitralAgency, &r.EndorsementStatus, &r.ApplicationFiles, &r.Balance, &r.BalanceList, &r.HaveNegotiatedPrice, &r.Sn, &r.PolicyRegen, &r.ClearList, &r.FilesToRemove, &r.OriginalPolicyScheme, &r.PolicyScheme, &r.InvoiceHeader, &r.CorrectLevel, &r.CorrectLog, &r.OriginalFiles, &r.Files, &r.RefusedReason, &r.Addi, &r.CreateTime, &r.UpdateTime, &r.Creator, &r.Remark, &r.Status)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select t_v_mistake_correct2")
	}
	return &r, nil
}

/*TVMistakeCorrectShow t_v_mistake_correct_show represents kuser.t_v_mistake_correct_show */
type TVMistakeCorrectShow struct {
	ID                       null.Int       `json:"ID,omitempty" db:"id,false,integer"`                                              /* id id */
	OrderID                  null.Int       `json:"OrderID,omitempty" db:"order_id,false,bigint"`                                    /* order_id order_id */
	InsuranceTypeID          null.Int       `json:"InsuranceTypeID,omitempty" db:"insurance_type_id,false,bigint"`                   /* insurance_type_id insurance_type_id */
	InsuranceTypeParentID    null.Int       `json:"InsuranceTypeParentID,omitempty" db:"insurance_type_parent_id,false,bigint"`      /* insurance_type_parent_id insurance_type_parent_id */
	OrgID                    null.Int       `json:"OrgID,omitempty" db:"org_id,false,bigint"`                                        /* org_id org_id */
	CommenceDate             null.Int       `json:"CommenceDate,omitempty" db:"commence_date,false,bigint"`                          /* commence_date commence_date */
	ExpiryDate               null.Int       `json:"ExpiryDate,omitempty" db:"expiry_date,false,bigint"`                              /* expiry_date expiry_date */
	ModifyType               null.String    `json:"ModifyType,omitempty" db:"modify_type,false,character varying"`                   /* modify_type modify_type */
	HaveInsuredList          null.Bool      `json:"HaveInsuredList,omitempty" db:"have_insured_list,false,boolean"`                  /* have_insured_list have_insured_list */
	InsuranceType            null.String    `json:"InsuranceType,omitempty" db:"insurance_type,false,character varying"`             /* insurance_type insurance_type */
	ActivityCategory         null.String    `json:"ActivityCategory,omitempty" db:"activity_category,false,character varying"`       /* activity_category activity_category */
	PlanID                   null.Int       `json:"PlanID,omitempty" db:"plan_id,false,bigint"`                                      /* plan_id plan_id */
	Insurer                  null.String    `json:"Insurer,omitempty" db:"insurer,false,text"`                                       /* insurer insurer */
	PlanName                 null.String    `json:"PlanName,omitempty" db:"plan_name,false,text"`                                    /* plan_name plan_name */
	OrgName                  null.String    `json:"OrgName,omitempty" db:"org_name,false,text"`                                      /* org_name org_name */
	OrgAddr                  null.String    `json:"OrgAddr,omitempty" db:"org_addr,false,text"`                                      /* org_addr org_addr */
	OrgCreditCode            null.String    `json:"OrgCreditCode,omitempty" db:"org_credit_code,false,text"`                         /* org_credit_code org_credit_code */
	OrgContact               null.String    `json:"OrgContact,omitempty" db:"org_contact,false,text"`                                /* org_contact org_contact */
	OrgPhone                 null.String    `json:"OrgPhone,omitempty" db:"org_phone,false,text"`                                    /* org_phone org_phone */
	OrgContactRole           null.String    `json:"OrgContactRole,omitempty" db:"org_contact_role,false,text"`                       /* org_contact_role org_contact_role */
	OrgCreditCodePic         null.String    `json:"OrgCreditCodePic,omitempty" db:"org_credit_code_pic,false,text"`                  /* org_credit_code_pic org_credit_code_pic */
	OrgSchoolCategory        null.String    `json:"OrgSchoolCategory,omitempty" db:"org_school_category,false,text"`                 /* org_school_category org_school_category */
	OrgProvince              null.String    `json:"OrgProvince,omitempty" db:"org_province,false,text"`                              /* org_province org_province */
	OrgCity                  null.String    `json:"OrgCity,omitempty" db:"org_city,false,text"`                                      /* org_city org_city */
	OrgDistrict              null.String    `json:"OrgDistrict,omitempty" db:"org_district,false,text"`                              /* org_district org_district */
	InsuredName              null.String    `json:"InsuredName,omitempty" db:"insured_name,false,text"`                              /* insured_name insured_name */
	InsuredProvince          null.String    `json:"InsuredProvince,omitempty" db:"insured_province,false,text"`                      /* insured_province insured_province */
	InsuredCity              null.String    `json:"InsuredCity,omitempty" db:"insured_city,false,text"`                              /* insured_city insured_city */
	InsuredDistrict          null.String    `json:"InsuredDistrict,omitempty" db:"insured_district,false,text"`                      /* insured_district insured_district */
	InsuredIsCompulsory      null.Bool      `json:"InsuredIsCompulsory,omitempty" db:"insured_is_compulsory,false,boolean"`          /* insured_is_compulsory insured_is_compulsory */
	InsuredCategory          null.String    `json:"InsuredCategory,omitempty" db:"insured_category,false,text"`                      /* insured_category insured_category */
	InsuredSchoolCategory    null.String    `json:"InsuredSchoolCategory,omitempty" db:"insured_school_category,false,text"`         /* insured_school_category insured_school_category */
	OriginalInsuredList      types.JSONText `json:"OriginalInsuredList,omitempty" db:"original_insured_list,false,jsonb"`            /* original_insured_list original_insured_list */
	InsuredList              types.JSONText `json:"InsuredList,omitempty" db:"insured_list,false,jsonb"`                             /* insured_list insured_list */
	InsertInsuredList        types.JSONText `json:"InsertInsuredList,omitempty" db:"insert_insured_list,false,jsonb"`                /* insert_insured_list insert_insured_list */
	DeleteInsuredList        types.JSONText `json:"DeleteInsuredList,omitempty" db:"delete_insured_list,false,jsonb"`                /* delete_insured_list delete_insured_list */
	UpdateInsuredList        types.JSONText `json:"UpdateInsuredList,omitempty" db:"update_insured_list,false,jsonb"`                /* update_insured_list update_insured_list */
	RequireUpdateInsuredList types.JSONText `json:"RequireUpdateInsuredList,omitempty" db:"require_update_insured_list,false,jsonb"` /* require_update_insured_list require_update_insured_list */
	PolicyScheme             types.JSONText `json:"PolicyScheme,omitempty" db:"policy_scheme,false,jsonb"`                           /* policy_scheme policy_scheme */
	ActivityName             null.String    `json:"ActivityName,omitempty" db:"activity_name,false,character varying"`               /* activity_name activity_name */
	ActivityDesc             null.String    `json:"ActivityDesc,omitempty" db:"activity_desc,false,character varying"`               /* activity_desc activity_desc */
	ActivityLocation         null.String    `json:"ActivityLocation,omitempty" db:"activity_location,false,character varying"`       /* activity_location activity_location */
	ActivityDateSet          null.String    `json:"ActivityDateSet,omitempty" db:"activity_date_set,false,character varying"`        /* activity_date_set activity_date_set */
	Indate                   null.Int       `json:"Indate,omitempty" db:"indate,false,bigint"`                                       /* indate indate */
	Policyholder             types.JSONText `json:"Policyholder,omitempty" db:"policyholder,false,jsonb"`                            /* policyholder policyholder */
	PolicyholderID           null.Int       `json:"PolicyholderID,omitempty" db:"policyholder_id,false,bigint"`                      /* policyholder_id policyholder_id */
	Insured                  types.JSONText `json:"Insured,omitempty" db:"insured,false,jsonb"`                                      /* insured insured */
	InsuredID                null.Int       `json:"InsuredID,omitempty" db:"insured_id,false,bigint"`                                /* insured_id insured_id */
	InsuredGroupByDay        null.Bool      `json:"InsuredGroupByDay,omitempty" db:"insured_group_by_day,false,boolean"`             /* insured_group_by_day insured_group_by_day */
	ChargeMode               null.String    `json:"ChargeMode,omitempty" db:"charge_mode,false,character varying"`                   /* charge_mode charge_mode */
	Amount                   null.Float     `json:"Amount,omitempty" db:"amount,false,double precision"`                             /* amount amount */
	InsuredCount             null.Int       `json:"InsuredCount,omitempty" db:"insured_count,false,smallint"`                        /* insured_count insured_count */
	InsuredType              null.String    `json:"InsuredType,omitempty" db:"insured_type,false,character varying"`                 /* insured_type insured_type */
	NonCompulsoryStudentNum  null.Int       `json:"NonCompulsoryStudentNum,omitempty" db:"non_compulsory_student_num,false,bigint"`  /* non_compulsory_student_num non_compulsory_student_num */
	CompulsoryStudentNum     null.Int       `json:"CompulsoryStudentNum,omitempty" db:"compulsory_student_num,false,bigint"`         /* compulsory_student_num compulsory_student_num */
	CanteenNum               null.Int       `json:"CanteenNum,omitempty" db:"canteen_num,false,bigint"`                              /* canteen_num canteen_num */
	ShopNum                  null.Int       `json:"ShopNum,omitempty" db:"shop_num,false,bigint"`                                    /* shop_num shop_num */
	DinnerNum                null.Int       `json:"DinnerNum,omitempty" db:"dinner_num,false,bigint"`                                /* dinner_num dinner_num */
	PayType                  null.String    `json:"PayType,omitempty" db:"pay_type,false,character varying"`                         /* pay_type pay_type */
	FeeScheme                types.JSONText `json:"FeeScheme,omitempty" db:"fee_scheme,false,jsonb"`                                 /* fee_scheme fee_scheme */
	DisputeHandling          null.String    `json:"DisputeHandling,omitempty" db:"dispute_handling,false,character varying"`         /* dispute_handling dispute_handling */
	HaveSuddenDeath          null.Bool      `json:"HaveSuddenDeath,omitempty" db:"have_sudden_death,false,boolean"`                  /* have_sudden_death have_sudden_death */
	PrevPolicyNo             null.String    `json:"PrevPolicyNo,omitempty" db:"prev_policy_no,false,character varying"`              /* prev_policy_no prev_policy_no */
	PoolName                 null.String    `json:"PoolName,omitempty" db:"pool_name,false,character varying"`                       /* pool_name pool_name */
	HaveExplosive            null.Bool      `json:"HaveExplosive,omitempty" db:"have_explosive,false,boolean"`                       /* have_explosive have_explosive */
	HaveRides                null.Bool      `json:"HaveRides,omitempty" db:"have_rides,false,boolean"`                               /* have_rides have_rides */
	InnerArea                null.Float     `json:"InnerArea,omitempty" db:"inner_area,false,double precision"`                      /* inner_area inner_area */
	OuterArea                null.Float     `json:"OuterArea,omitempty" db:"outer_area,false,double precision"`                      /* outer_area outer_area */
	TrafficNum               null.Int       `json:"TrafficNum,omitempty" db:"traffic_num,false,integer"`                             /* traffic_num traffic_num */
	TemperatureType          null.String    `json:"TemperatureType,omitempty" db:"temperature_type,false,character varying"`         /* temperature_type temperature_type */
	OpenPoolNum              null.Int       `json:"OpenPoolNum,omitempty" db:"open_pool_num,false,smallint"`                         /* open_pool_num open_pool_num */
	HeatedPoolNum            null.Int       `json:"HeatedPoolNum,omitempty" db:"heated_pool_num,false,smallint"`                     /* heated_pool_num heated_pool_num */
	TrainingPoolNum          null.Int       `json:"TrainingPoolNum,omitempty" db:"training_pool_num,false,smallint"`                 /* training_pool_num training_pool_num */
	PoolNum                  null.Int       `json:"PoolNum,omitempty" db:"pool_num,false,smallint"`                                  /* pool_num pool_num */
	CustomType               null.String    `json:"CustomType,omitempty" db:"custom_type,false,character varying"`                   /* custom_type custom_type */
	Same                     null.Bool      `json:"Same,omitempty" db:"same,false,boolean"`                                          /* same same */
	ArbitralAgency           null.String    `json:"ArbitralAgency,omitempty" db:"arbitral_agency,false,character varying"`           /* arbitral_agency arbitral_agency */
	Files                    types.JSONText `json:"Files,omitempty" db:"files,false,jsonb"`                                          /* files files */
	HaveNegotiatedPrice      null.Bool      `json:"HaveNegotiatedPrice,omitempty" db:"have_negotiated_price,false,boolean"`          /* have_negotiated_price have_negotiated_price */
	EndorsementStatus        null.String    `json:"EndorsementStatus,omitempty" db:"endorsement_status,false,character varying"`     /* endorsement_status endorsement_status */
	ApplicationFiles         types.JSONText `json:"ApplicationFiles,omitempty" db:"application_files,false,jsonb"`                   /* application_files application_files */
	NeedBalance              null.Bool      `json:"NeedBalance,omitempty" db:"need_balance,false,boolean"`                           /* need_balance need_balance */
	Balance                  null.Float     `json:"Balance,omitempty" db:"balance,false,double precision"`                           /* balance balance */
	BalanceList              types.JSONText `json:"BalanceList,omitempty" db:"balance_list,false,jsonb"`                             /* balance_list balance_list */
	OrderStatus              null.String    `json:"OrderStatus,omitempty" db:"order_status,false,character varying"`                 /* order_status order_status */
	Sn                       null.String    `json:"Sn,omitempty" db:"sn,false,text"`                                                 /* sn sn */
	RevokedPolicyNo          null.String    `json:"RevokedPolicyNo,omitempty" db:"revoked_policy_no,false,character varying"`        /* revoked_policy_no revoked_policy_no */
	PolicyRegen              null.Bool      `json:"PolicyRegen,omitempty" db:"policy_regen,false,boolean"`                           /* policy_regen policy_regen */
	ClearList                null.Bool      `json:"ClearList,omitempty" db:"clear_list,false,boolean"`                               /* clear_list clear_list */
	FilesToRemove            null.String    `json:"FilesToRemove,omitempty" db:"files_to_remove,false,character varying"`            /* files_to_remove files_to_remove */
	InvoiceHeader            null.String    `json:"InvoiceHeader,omitempty" db:"invoice_header,false,character varying"`             /* invoice_header invoice_header */
	CorrectLevel             null.String    `json:"CorrectLevel,omitempty" db:"correct_level,false,character varying"`               /* correct_level correct_level */
	CorrectLog               types.JSONText `json:"CorrectLog,omitempty" db:"correct_log,false,jsonb"`                               /* correct_log correct_log */
	RefusedReason            null.String    `json:"RefusedReason,omitempty" db:"refused_reason,false,character varying"`             /* refused_reason refused_reason */
	Addi                     types.JSONText `json:"Addi,omitempty" db:"addi,false,jsonb"`                                            /* addi addi */
	CreateTime               null.Int       `json:"CreateTime,omitempty" db:"create_time,false,bigint"`                              /* create_time create_time */
	UpdateTime               null.Int       `json:"UpdateTime,omitempty" db:"update_time,false,bigint"`                              /* update_time update_time */
	Creator                  null.Int       `json:"Creator,omitempty" db:"creator,false,bigint"`                                     /* creator creator */
	Remark                   null.String    `json:"Remark,omitempty" db:"remark,false,character varying"`                            /* remark remark */
	Status                   null.String    `json:"Status,omitempty" db:"status,false,character varying"`                            /* status status */
	Filter                                  // build DML where clause
}

//TVMistakeCorrectShowFields full field list for default query
var TVMistakeCorrectShowFields = []string{
	"ID",
	"OrderID",
	"InsuranceTypeID",
	"InsuranceTypeParentID",
	"OrgID",
	"CommenceDate",
	"ExpiryDate",
	"ModifyType",
	"HaveInsuredList",
	"InsuranceType",
	"ActivityCategory",
	"PlanID",
	"Insurer",
	"PlanName",
	"OrgName",
	"OrgAddr",
	"OrgCreditCode",
	"OrgContact",
	"OrgPhone",
	"OrgContactRole",
	"OrgCreditCodePic",
	"OrgSchoolCategory",
	"OrgProvince",
	"OrgCity",
	"OrgDistrict",
	"InsuredName",
	"InsuredProvince",
	"InsuredCity",
	"InsuredDistrict",
	"InsuredIsCompulsory",
	"InsuredCategory",
	"InsuredSchoolCategory",
	"OriginalInsuredList",
	"InsuredList",
	"InsertInsuredList",
	"DeleteInsuredList",
	"UpdateInsuredList",
	"RequireUpdateInsuredList",
	"PolicyScheme",
	"ActivityName",
	"ActivityDesc",
	"ActivityLocation",
	"ActivityDateSet",
	"Indate",
	"Policyholder",
	"PolicyholderID",
	"Insured",
	"InsuredID",
	"InsuredGroupByDay",
	"ChargeMode",
	"Amount",
	"InsuredCount",
	"InsuredType",
	"NonCompulsoryStudentNum",
	"CompulsoryStudentNum",
	"CanteenNum",
	"ShopNum",
	"DinnerNum",
	"PayType",
	"FeeScheme",
	"DisputeHandling",
	"HaveSuddenDeath",
	"PrevPolicyNo",
	"PoolName",
	"HaveExplosive",
	"HaveRides",
	"InnerArea",
	"OuterArea",
	"TrafficNum",
	"TemperatureType",
	"OpenPoolNum",
	"HeatedPoolNum",
	"TrainingPoolNum",
	"PoolNum",
	"CustomType",
	"Same",
	"ArbitralAgency",
	"Files",
	"HaveNegotiatedPrice",
	"EndorsementStatus",
	"ApplicationFiles",
	"NeedBalance",
	"Balance",
	"BalanceList",
	"OrderStatus",
	"Sn",
	"RevokedPolicyNo",
	"PolicyRegen",
	"ClearList",
	"FilesToRemove",
	"InvoiceHeader",
	"CorrectLevel",
	"CorrectLog",
	"RefusedReason",
	"Addi",
	"CreateTime",
	"UpdateTime",
	"Creator",
	"Remark",
	"Status",
}

//Fields return all fields of struct.
func (r *TVMistakeCorrectShow) Fields() []string {
	return TVMistakeCorrectShowFields
}

//GetTableName return the associated db table name.
func (r *TVMistakeCorrectShow) GetTableName() string {
	var viewNamePattern = regexp.MustCompile(`(?i)^t_v_[a-z0-9_]+$`)
	tableName := "t_v_mistake_correct_show"
	if viewNamePattern.MatchString(tableName) {
		return tableName[2:]
	}
	return tableName
}

// Create inserts the TVMistakeCorrectShow to the database.
func (r *TVMistakeCorrectShow) Create(db Queryer) error {
	_, err := db.Exec(
		`INSERT INTO t_v_mistake_correct_show (id, order_id, insurance_type_id, insurance_type_parent_id, org_id, commence_date, expiry_date, modify_type, have_insured_list, insurance_type, activity_category, plan_id, insurer, plan_name, org_name, org_addr, org_credit_code, org_contact, org_phone, org_contact_role, org_credit_code_pic, org_school_category, org_province, org_city, org_district, insured_name, insured_province, insured_city, insured_district, insured_is_compulsory, insured_category, insured_school_category, original_insured_list, insured_list, insert_insured_list, delete_insured_list, update_insured_list, require_update_insured_list, policy_scheme, activity_name, activity_desc, activity_location, activity_date_set, indate, policyholder, policyholder_id, insured, insured_id, insured_group_by_day, charge_mode, amount, insured_count, insured_type, non_compulsory_student_num, compulsory_student_num, canteen_num, shop_num, dinner_num, pay_type, fee_scheme, dispute_handling, have_sudden_death, prev_policy_no, pool_name, have_explosive, have_rides, inner_area, outer_area, traffic_num, temperature_type, open_pool_num, heated_pool_num, training_pool_num, pool_num, custom_type, same, arbitral_agency, files, have_negotiated_price, endorsement_status, application_files, need_balance, balance, balance_list, order_status, sn, revoked_policy_no, policy_regen, clear_list, files_to_remove, invoice_header, correct_level, correct_log, refused_reason, addi, create_time, update_time, creator, remark, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36, $37, $38, $39, $40, $41, $42, $43, $44, $45, $46, $47, $48, $49, $50, $51, $52, $53, $54, $55, $56, $57, $58, $59, $60, $61, $62, $63, $64, $65, $66, $67, $68, $69, $70, $71, $72, $73, $74, $75, $76, $77, $78, $79, $80, $81, $82, $83, $84, $85, $86, $87, $88, $89, $90, $91, $92, $93, $94, $95, $96, $97, $98, $99, $100)`,
		&r.ID, &r.OrderID, &r.InsuranceTypeID, &r.InsuranceTypeParentID, &r.OrgID, &r.CommenceDate, &r.ExpiryDate, &r.ModifyType, &r.HaveInsuredList, &r.InsuranceType, &r.ActivityCategory, &r.PlanID, &r.Insurer, &r.PlanName, &r.OrgName, &r.OrgAddr, &r.OrgCreditCode, &r.OrgContact, &r.OrgPhone, &r.OrgContactRole, &r.OrgCreditCodePic, &r.OrgSchoolCategory, &r.OrgProvince, &r.OrgCity, &r.OrgDistrict, &r.InsuredName, &r.InsuredProvince, &r.InsuredCity, &r.InsuredDistrict, &r.InsuredIsCompulsory, &r.InsuredCategory, &r.InsuredSchoolCategory, &r.OriginalInsuredList, &r.InsuredList, &r.InsertInsuredList, &r.DeleteInsuredList, &r.UpdateInsuredList, &r.RequireUpdateInsuredList, &r.PolicyScheme, &r.ActivityName, &r.ActivityDesc, &r.ActivityLocation, &r.ActivityDateSet, &r.Indate, &r.Policyholder, &r.PolicyholderID, &r.Insured, &r.InsuredID, &r.InsuredGroupByDay, &r.ChargeMode, &r.Amount, &r.InsuredCount, &r.InsuredType, &r.NonCompulsoryStudentNum, &r.CompulsoryStudentNum, &r.CanteenNum, &r.ShopNum, &r.DinnerNum, &r.PayType, &r.FeeScheme, &r.DisputeHandling, &r.HaveSuddenDeath, &r.PrevPolicyNo, &r.PoolName, &r.HaveExplosive, &r.HaveRides, &r.InnerArea, &r.OuterArea, &r.TrafficNum, &r.TemperatureType, &r.OpenPoolNum, &r.HeatedPoolNum, &r.TrainingPoolNum, &r.PoolNum, &r.CustomType, &r.Same, &r.ArbitralAgency, &r.Files, &r.HaveNegotiatedPrice, &r.EndorsementStatus, &r.ApplicationFiles, &r.NeedBalance, &r.Balance, &r.BalanceList, &r.OrderStatus, &r.Sn, &r.RevokedPolicyNo, &r.PolicyRegen, &r.ClearList, &r.FilesToRemove, &r.InvoiceHeader, &r.CorrectLevel, &r.CorrectLog, &r.RefusedReason, &r.Addi, &r.CreateTime, &r.UpdateTime, &r.Creator, &r.Remark, &r.Status)
	if err != nil {
		return errors.Wrap(err, "failed to insert t_v_mistake_correct_show")
	}
	return nil
}

// GetTVMistakeCorrectShowByPk select the TVMistakeCorrectShow from the database.
func GetTVMistakeCorrectShowByPk(db Queryer) (*TVMistakeCorrectShow, error) {
	var r TVMistakeCorrectShow
	err := db.QueryRow(
		`SELECT id, order_id, insurance_type_id, insurance_type_parent_id, org_id, commence_date, expiry_date, modify_type, have_insured_list, insurance_type, activity_category, plan_id, insurer, plan_name, org_name, org_addr, org_credit_code, org_contact, org_phone, org_contact_role, org_credit_code_pic, org_school_category, org_province, org_city, org_district, insured_name, insured_province, insured_city, insured_district, insured_is_compulsory, insured_category, insured_school_category, original_insured_list, insured_list, insert_insured_list, delete_insured_list, update_insured_list, require_update_insured_list, policy_scheme, activity_name, activity_desc, activity_location, activity_date_set, indate, policyholder, policyholder_id, insured, insured_id, insured_group_by_day, charge_mode, amount, insured_count, insured_type, non_compulsory_student_num, compulsory_student_num, canteen_num, shop_num, dinner_num, pay_type, fee_scheme, dispute_handling, have_sudden_death, prev_policy_no, pool_name, have_explosive, have_rides, inner_area, outer_area, traffic_num, temperature_type, open_pool_num, heated_pool_num, training_pool_num, pool_num, custom_type, same, arbitral_agency, files, have_negotiated_price, endorsement_status, application_files, need_balance, balance, balance_list, order_status, sn, revoked_policy_no, policy_regen, clear_list, files_to_remove, invoice_header, correct_level, correct_log, refused_reason, addi, create_time, update_time, creator, remark, status FROM t_v_mistake_correct_show WHERE `,
	).Scan(&r.ID, &r.OrderID, &r.InsuranceTypeID, &r.InsuranceTypeParentID, &r.OrgID, &r.CommenceDate, &r.ExpiryDate, &r.ModifyType, &r.HaveInsuredList, &r.InsuranceType, &r.ActivityCategory, &r.PlanID, &r.Insurer, &r.PlanName, &r.OrgName, &r.OrgAddr, &r.OrgCreditCode, &r.OrgContact, &r.OrgPhone, &r.OrgContactRole, &r.OrgCreditCodePic, &r.OrgSchoolCategory, &r.OrgProvince, &r.OrgCity, &r.OrgDistrict, &r.InsuredName, &r.InsuredProvince, &r.InsuredCity, &r.InsuredDistrict, &r.InsuredIsCompulsory, &r.InsuredCategory, &r.InsuredSchoolCategory, &r.OriginalInsuredList, &r.InsuredList, &r.InsertInsuredList, &r.DeleteInsuredList, &r.UpdateInsuredList, &r.RequireUpdateInsuredList, &r.PolicyScheme, &r.ActivityName, &r.ActivityDesc, &r.ActivityLocation, &r.ActivityDateSet, &r.Indate, &r.Policyholder, &r.PolicyholderID, &r.Insured, &r.InsuredID, &r.InsuredGroupByDay, &r.ChargeMode, &r.Amount, &r.InsuredCount, &r.InsuredType, &r.NonCompulsoryStudentNum, &r.CompulsoryStudentNum, &r.CanteenNum, &r.ShopNum, &r.DinnerNum, &r.PayType, &r.FeeScheme, &r.DisputeHandling, &r.HaveSuddenDeath, &r.PrevPolicyNo, &r.PoolName, &r.HaveExplosive, &r.HaveRides, &r.InnerArea, &r.OuterArea, &r.TrafficNum, &r.TemperatureType, &r.OpenPoolNum, &r.HeatedPoolNum, &r.TrainingPoolNum, &r.PoolNum, &r.CustomType, &r.Same, &r.ArbitralAgency, &r.Files, &r.HaveNegotiatedPrice, &r.EndorsementStatus, &r.ApplicationFiles, &r.NeedBalance, &r.Balance, &r.BalanceList, &r.OrderStatus, &r.Sn, &r.RevokedPolicyNo, &r.PolicyRegen, &r.ClearList, &r.FilesToRemove, &r.InvoiceHeader, &r.CorrectLevel, &r.CorrectLog, &r.RefusedReason, &r.Addi, &r.CreateTime, &r.UpdateTime, &r.Creator, &r.Remark, &r.Status)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select t_v_mistake_correct_show")
	}
	return &r, nil
}

/*TVOrder t_v_order represents kuser.t_v_order */
type TVOrder struct {
	ID             null.Int       `json:"ID,omitempty" db:"id,false,integer"`                                     /* id id */
	TradeNo        null.String    `json:"TradeNo,omitempty" db:"trade_no,false,character varying"`                /* trade_no trade_no */
	OrgID          null.Int       `json:"OrgID,omitempty" db:"org_id,false,bigint"`                               /* org_id org_id */
	InsuredID      null.Int       `json:"InsuredID,omitempty" db:"insured_id,false,bigint"`                       /* insured_id insured_id */
	PolicyholderID null.Int       `json:"PolicyholderID,omitempty" db:"policyholder_id,false,bigint"`             /* policyholder_id policyholder_id */
	Sign           null.String    `json:"Sign,omitempty" db:"sign,false,character varying"`                       /* sign sign */
	InsureOrderNo  null.String    `json:"InsureOrderNo,omitempty" db:"insure_order_no,false,character varying"`   /* insure_order_no insure_order_no */
	PayOrderNo     null.String    `json:"PayOrderNo,omitempty" db:"pay_order_no,false,character varying"`         /* pay_order_no pay_order_no */
	Batch          null.String    `json:"Batch,omitempty" db:"batch,false,character varying"`                     /* batch batch */
	CreateTime     null.Int       `json:"CreateTime,omitempty" db:"create_time,false,bigint"`                     /* create_time create_time */
	PayTime        null.Int       `json:"PayTime,omitempty" db:"pay_time,false,bigint"`                           /* pay_time pay_time */
	PayType        null.String    `json:"PayType,omitempty" db:"pay_type,false,character varying"`                /* pay_type pay_type */
	Amount         null.Float     `json:"Amount,omitempty" db:"amount,false,double precision"`                    /* amount amount */
	UnitPrice      null.Float     `json:"UnitPrice,omitempty" db:"unit_price,false,double precision"`             /* unit_price unit_price */
	CommenceDate   null.Int       `json:"CommenceDate,omitempty" db:"commence_date,false,bigint"`                 /* commence_date commence_date */
	ExpiryDate     null.Int       `json:"ExpiryDate,omitempty" db:"expiry_date,false,bigint"`                     /* expiry_date expiry_date */
	Indate         null.Int       `json:"Indate,omitempty" db:"indate,false,bigint"`                              /* indate indate */
	ConfirmRefund  null.Bool      `json:"ConfirmRefund,omitempty" db:"confirm_refund,false,boolean"`              /* confirm_refund confirm_refund */
	Relation       null.String    `json:"Relation,omitempty" db:"relation,false,character varying"`               /* relation relation */
	InsuranceType  null.String    `json:"InsuranceType,omitempty" db:"insurance_type,false,character varying"`    /* insurance_type insurance_type */
	PlanID         null.Int       `json:"PlanID,omitempty" db:"plan_id,false,bigint"`                             /* plan_id plan_id */
	PlanName       null.String    `json:"PlanName,omitempty" db:"plan_name,false,character varying"`              /* plan_name plan_name */
	Insurer        null.String    `json:"Insurer,omitempty" db:"insurer,false,character varying"`                 /* insurer insurer */
	PolicyDoc      null.String    `json:"PolicyDoc,omitempty" db:"policy_doc,false,character varying"`            /* policy_doc policy_doc */
	HealthSurvey   types.JSONText `json:"HealthSurvey,omitempty" db:"health_survey,false,jsonb"`                  /* health_survey health_survey */
	Same           null.Bool      `json:"Same,omitempty" db:"same,false,boolean"`                                 /* same same */
	OrderFiles     types.JSONText `json:"OrderFiles,omitempty" db:"order_files,false,jsonb"`                      /* order_files order_files */
	Addi           types.JSONText `json:"Addi,omitempty" db:"addi,false,jsonb"`                                   /* addi addi */
	Status         null.String    `json:"Status,omitempty" db:"status,false,character varying"`                   /* status status */
	Creator        null.Int       `json:"Creator,omitempty" db:"creator,false,bigint"`                            /* creator creator */
	IOfficialName  null.String    `json:"IOfficialName,omitempty" db:"i_official_name,false,character varying"`   /* i_official_name i_official_name */
	IIDCardType    null.String    `json:"IIDCardType,omitempty" db:"i_id_card_type,false,character varying"`      /* i_id_card_type i_id_card_type */
	IIDCardNo      null.String    `json:"IIDCardNo,omitempty" db:"i_id_card_no,false,character varying"`          /* i_id_card_no i_id_card_no */
	IMobilePhone   null.String    `json:"IMobilePhone,omitempty" db:"i_mobile_phone,false,character varying"`     /* i_mobile_phone i_mobile_phone */
	IGender        null.String    `json:"IGender,omitempty" db:"i_gender,false,character varying"`                /* i_gender i_gender */
	IBirthday      null.Int       `json:"IBirthday,omitempty" db:"i_birthday,false,bigint"`                       /* i_birthday i_birthday */
	IAddi          types.JSONText `json:"IAddi,omitempty" db:"i_addi,false,jsonb"`                                /* i_addi i_addi */
	HOfficialName  null.String    `json:"HOfficialName,omitempty" db:"h_official_name,false,character varying"`   /* h_official_name h_official_name */
	HIDCardType    null.String    `json:"HIDCardType,omitempty" db:"h_id_card_type,false,character varying"`      /* h_id_card_type h_id_card_type */
	HIDCardNo      null.String    `json:"HIDCardNo,omitempty" db:"h_id_card_no,false,character varying"`          /* h_id_card_no h_id_card_no */
	HMobilePhone   null.String    `json:"HMobilePhone,omitempty" db:"h_mobile_phone,false,character varying"`     /* h_mobile_phone h_mobile_phone */
	HAddi          types.JSONText `json:"HAddi,omitempty" db:"h_addi,false,jsonb"`                                /* h_addi h_addi */
	Subdistrict    null.String    `json:"Subdistrict,omitempty" db:"subdistrict,false,character varying"`         /* subdistrict subdistrict */
	Faculty        null.String    `json:"Faculty,omitempty" db:"faculty,false,character varying"`                 /* faculty faculty */
	Grade          null.String    `json:"Grade,omitempty" db:"grade,false,character varying"`                     /* grade grade */
	Class          null.String    `json:"Class,omitempty" db:"class,false,character varying"`                     /* class class */
	XCreateTime    null.Int       `json:"XCreateTime,omitempty" db:"x_create_time,false,bigint"`                  /* x_create_time x_create_time */
	School         null.String    `json:"School,omitempty" db:"school,false,character varying"`                   /* school school */
	SFaculty       types.JSONText `json:"SFaculty,omitempty" db:"s_faculty,false,jsonb"`                          /* s_faculty s_faculty */
	SBranches      types.JSONText `json:"SBranches,omitempty" db:"s_branches,false,jsonb"`                        /* s_branches s_branches */
	SCategory      null.String    `json:"SCategory,omitempty" db:"s_category,false,character varying"`            /* s_category s_category */
	Province       null.String    `json:"Province,omitempty" db:"province,false,character varying"`               /* province province */
	City           null.String    `json:"City,omitempty" db:"city,false,character varying"`                       /* city city */
	District       null.String    `json:"District,omitempty" db:"district,false,character varying"`               /* district district */
	DataSyncTarget null.String    `json:"DataSyncTarget,omitempty" db:"data_sync_target,false,character varying"` /* data_sync_target data_sync_target */
	SaleManagers   types.JSONText `json:"SaleManagers,omitempty" db:"sale_managers,false,jsonb"`                  /* sale_managers sale_managers */
	SchoolManagers types.JSONText `json:"SchoolManagers,omitempty" db:"school_managers,false,jsonb"`              /* school_managers school_managers */
	PurchaseRule   types.JSONText `json:"PurchaseRule,omitempty" db:"purchase_rule,false,jsonb"`                  /* purchase_rule purchase_rule */
	SCreateTime    null.Int       `json:"SCreateTime,omitempty" db:"s_create_time,false,bigint"`                  /* s_create_time s_create_time */
	Filter                        // build DML where clause
}

//TVOrderFields full field list for default query
var TVOrderFields = []string{
	"ID",
	"TradeNo",
	"OrgID",
	"InsuredID",
	"PolicyholderID",
	"Sign",
	"InsureOrderNo",
	"PayOrderNo",
	"Batch",
	"CreateTime",
	"PayTime",
	"PayType",
	"Amount",
	"UnitPrice",
	"CommenceDate",
	"ExpiryDate",
	"Indate",
	"ConfirmRefund",
	"Relation",
	"InsuranceType",
	"PlanID",
	"PlanName",
	"Insurer",
	"PolicyDoc",
	"HealthSurvey",
	"Same",
	"OrderFiles",
	"Addi",
	"Status",
	"Creator",
	"IOfficialName",
	"IIDCardType",
	"IIDCardNo",
	"IMobilePhone",
	"IGender",
	"IBirthday",
	"IAddi",
	"HOfficialName",
	"HIDCardType",
	"HIDCardNo",
	"HMobilePhone",
	"HAddi",
	"Subdistrict",
	"Faculty",
	"Grade",
	"Class",
	"XCreateTime",
	"School",
	"SFaculty",
	"SBranches",
	"SCategory",
	"Province",
	"City",
	"District",
	"DataSyncTarget",
	"SaleManagers",
	"SchoolManagers",
	"PurchaseRule",
	"SCreateTime",
}

//Fields return all fields of struct.
func (r *TVOrder) Fields() []string {
	return TVOrderFields
}

//GetTableName return the associated db table name.
func (r *TVOrder) GetTableName() string {
	var viewNamePattern = regexp.MustCompile(`(?i)^t_v_[a-z0-9_]+$`)
	tableName := "t_v_order"
	if viewNamePattern.MatchString(tableName) {
		return tableName[2:]
	}
	return tableName
}

// Create inserts the TVOrder to the database.
func (r *TVOrder) Create(db Queryer) error {
	_, err := db.Exec(
		`INSERT INTO t_v_order (id, trade_no, org_id, insured_id, policyholder_id, sign, insure_order_no, pay_order_no, batch, create_time, pay_time, pay_type, amount, unit_price, commence_date, expiry_date, indate, confirm_refund, relation, insurance_type, plan_id, plan_name, insurer, policy_doc, health_survey, same, order_files, addi, status, creator, i_official_name, i_id_card_type, i_id_card_no, i_mobile_phone, i_gender, i_birthday, i_addi, h_official_name, h_id_card_type, h_id_card_no, h_mobile_phone, h_addi, subdistrict, faculty, grade, class, x_create_time, school, s_faculty, s_branches, s_category, province, city, district, data_sync_target, sale_managers, school_managers, purchase_rule, s_create_time) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36, $37, $38, $39, $40, $41, $42, $43, $44, $45, $46, $47, $48, $49, $50, $51, $52, $53, $54, $55, $56, $57, $58, $59)`,
		&r.ID, &r.TradeNo, &r.OrgID, &r.InsuredID, &r.PolicyholderID, &r.Sign, &r.InsureOrderNo, &r.PayOrderNo, &r.Batch, &r.CreateTime, &r.PayTime, &r.PayType, &r.Amount, &r.UnitPrice, &r.CommenceDate, &r.ExpiryDate, &r.Indate, &r.ConfirmRefund, &r.Relation, &r.InsuranceType, &r.PlanID, &r.PlanName, &r.Insurer, &r.PolicyDoc, &r.HealthSurvey, &r.Same, &r.OrderFiles, &r.Addi, &r.Status, &r.Creator, &r.IOfficialName, &r.IIDCardType, &r.IIDCardNo, &r.IMobilePhone, &r.IGender, &r.IBirthday, &r.IAddi, &r.HOfficialName, &r.HIDCardType, &r.HIDCardNo, &r.HMobilePhone, &r.HAddi, &r.Subdistrict, &r.Faculty, &r.Grade, &r.Class, &r.XCreateTime, &r.School, &r.SFaculty, &r.SBranches, &r.SCategory, &r.Province, &r.City, &r.District, &r.DataSyncTarget, &r.SaleManagers, &r.SchoolManagers, &r.PurchaseRule, &r.SCreateTime)
	if err != nil {
		return errors.Wrap(err, "failed to insert t_v_order")
	}
	return nil
}

// GetTVOrderByPk select the TVOrder from the database.
func GetTVOrderByPk(db Queryer) (*TVOrder, error) {
	var r TVOrder
	err := db.QueryRow(
		`SELECT id, trade_no, org_id, insured_id, policyholder_id, sign, insure_order_no, pay_order_no, batch, create_time, pay_time, pay_type, amount, unit_price, commence_date, expiry_date, indate, confirm_refund, relation, insurance_type, plan_id, plan_name, insurer, policy_doc, health_survey, same, order_files, addi, status, creator, i_official_name, i_id_card_type, i_id_card_no, i_mobile_phone, i_gender, i_birthday, i_addi, h_official_name, h_id_card_type, h_id_card_no, h_mobile_phone, h_addi, subdistrict, faculty, grade, class, x_create_time, school, s_faculty, s_branches, s_category, province, city, district, data_sync_target, sale_managers, school_managers, purchase_rule, s_create_time FROM t_v_order WHERE `,
	).Scan(&r.ID, &r.TradeNo, &r.OrgID, &r.InsuredID, &r.PolicyholderID, &r.Sign, &r.InsureOrderNo, &r.PayOrderNo, &r.Batch, &r.CreateTime, &r.PayTime, &r.PayType, &r.Amount, &r.UnitPrice, &r.CommenceDate, &r.ExpiryDate, &r.Indate, &r.ConfirmRefund, &r.Relation, &r.InsuranceType, &r.PlanID, &r.PlanName, &r.Insurer, &r.PolicyDoc, &r.HealthSurvey, &r.Same, &r.OrderFiles, &r.Addi, &r.Status, &r.Creator, &r.IOfficialName, &r.IIDCardType, &r.IIDCardNo, &r.IMobilePhone, &r.IGender, &r.IBirthday, &r.IAddi, &r.HOfficialName, &r.HIDCardType, &r.HIDCardNo, &r.HMobilePhone, &r.HAddi, &r.Subdistrict, &r.Faculty, &r.Grade, &r.Class, &r.XCreateTime, &r.School, &r.SFaculty, &r.SBranches, &r.SCategory, &r.Province, &r.City, &r.District, &r.DataSyncTarget, &r.SaleManagers, &r.SchoolManagers, &r.PurchaseRule, &r.SCreateTime)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select t_v_order")
	}
	return &r, nil
}

/*TVOrder2 t_v_order2 represents kuser.t_v_order2 */
type TVOrder2 struct {
	ID                         null.Int       `json:"ID,omitempty" db:"id,false,integer"`                                                       /* id id */
	TradeNo                    null.String    `json:"TradeNo,omitempty" db:"trade_no,false,character varying"`                                  /* trade_no trade_no */
	PayOrderNo                 null.String    `json:"PayOrderNo,omitempty" db:"pay_order_no,false,character varying"`                           /* pay_order_no pay_order_no */
	InsureOrderNo              null.String    `json:"InsureOrderNo,omitempty" db:"insure_order_no,false,character varying"`                     /* insure_order_no insure_order_no */
	Batch                      null.String    `json:"Batch,omitempty" db:"batch,false,character varying"`                                       /* batch batch */
	CreateTime                 null.Int       `json:"CreateTime,omitempty" db:"create_time,false,bigint"`                                       /* create_time create_time */
	PayTime                    null.Int       `json:"PayTime,omitempty" db:"pay_time,false,bigint"`                                             /* pay_time pay_time */
	PayChannel                 null.String    `json:"PayChannel,omitempty" db:"pay_channel,false,character varying"`                            /* pay_channel pay_channel */
	PayType                    null.String    `json:"PayType,omitempty" db:"pay_type,false,character varying"`                                  /* pay_type pay_type */
	PayName                    null.String    `json:"PayName,omitempty" db:"pay_name,false,character varying"`                                  /* pay_name pay_name */
	UnitPrice                  null.Float     `json:"UnitPrice,omitempty" db:"unit_price,false,double precision"`                               /* unit_price unit_price */
	Amount                     null.Float     `json:"Amount,omitempty" db:"amount,false,double precision"`                                      /* amount amount */
	Balance                    null.Float     `json:"Balance,omitempty" db:"balance,false,double precision"`                                    /* balance balance */
	BalanceList                types.JSONText `json:"BalanceList,omitempty" db:"balance_list,false,jsonb"`                                      /* balance_list balance_list */
	OrgID                      null.Int       `json:"OrgID,omitempty" db:"org_id,false,bigint"`                                                 /* org_id org_id */
	HaveSuddenDeath            null.Bool      `json:"HaveSuddenDeath,omitempty" db:"have_sudden_death,false,boolean"`                           /* have_sudden_death have_sudden_death */
	GroundNum                  null.Int       `json:"GroundNum,omitempty" db:"ground_num,false,smallint"`                                       /* ground_num ground_num */
	PlanType                   null.String    `json:"PlanType,omitempty" db:"plan_type,false,text"`                                             /* plan_type plan_type */
	RemindersNum               null.Int       `json:"RemindersNum,omitempty" db:"reminders_num,false,smallint"`                                 /* reminders_num reminders_num */
	OrgManagerID               null.Int       `json:"OrgManagerID,omitempty" db:"org_manager_id,false,bigint"`                                  /* org_manager_id org_manager_id */
	InsuranceType              null.String    `json:"InsuranceType,omitempty" db:"insurance_type,false,character varying"`                      /* insurance_type insurance_type */
	HaveDinnerNum              null.Bool      `json:"HaveDinnerNum,omitempty" db:"have_dinner_num,false,boolean"`                               /* have_dinner_num have_dinner_num */
	HaveConfirmDate            null.Bool      `json:"HaveConfirmDate,omitempty" db:"have_confirm_date,false,boolean"`                           /* have_confirm_date have_confirm_date */
	InsuranceTypeID            null.Int       `json:"InsuranceTypeID,omitempty" db:"insurance_type_id,false,bigint"`                            /* insurance_type_id insurance_type_id */
	HealthSurvey               types.JSONText `json:"HealthSurvey,omitempty" db:"health_survey,false,jsonb"`                                    /* health_survey health_survey */
	PlanID                     null.Int       `json:"PlanID,omitempty" db:"plan_id,false,bigint"`                                               /* plan_id plan_id */
	PlanName                   null.String    `json:"PlanName,omitempty" db:"plan_name,false,character varying"`                                /* plan_name plan_name */
	Insurer                    null.String    `json:"Insurer,omitempty" db:"insurer,false,character varying"`                                   /* insurer insurer */
	PolicyScheme               types.JSONText `json:"PolicyScheme,omitempty" db:"policy_scheme,false,jsonb"`                                    /* policy_scheme policy_scheme */
	PolicyDoc                  null.String    `json:"PolicyDoc,omitempty" db:"policy_doc,false,character varying"`                              /* policy_doc policy_doc */
	ActivityName               null.String    `json:"ActivityName,omitempty" db:"activity_name,false,character varying"`                        /* activity_name activity_name */
	ActivityCategory           null.String    `json:"ActivityCategory,omitempty" db:"activity_category,false,character varying"`                /* activity_category activity_category */
	ActivityDesc               null.String    `json:"ActivityDesc,omitempty" db:"activity_desc,false,character varying"`                        /* activity_desc activity_desc */
	ActivityLocation           null.String    `json:"ActivityLocation,omitempty" db:"activity_location,false,character varying"`                /* activity_location activity_location */
	ActivityDateSet            null.String    `json:"ActivityDateSet,omitempty" db:"activity_date_set,false,character varying"`                 /* activity_date_set activity_date_set */
	CopiesNum                  null.Int       `json:"CopiesNum,omitempty" db:"copies_num,false,smallint"`                                       /* copies_num copies_num */
	InsuredCount               null.Int       `json:"InsuredCount,omitempty" db:"insured_count,false,smallint"`                                 /* insured_count insured_count */
	CompulsoryStudentNum       null.Int       `json:"CompulsoryStudentNum,omitempty" db:"compulsory_student_num,false,bigint"`                  /* compulsory_student_num compulsory_student_num */
	NonCompulsoryStudentNum    null.Int       `json:"NonCompulsoryStudentNum,omitempty" db:"non_compulsory_student_num,false,bigint"`           /* non_compulsory_student_num non_compulsory_student_num */
	Contact                    types.JSONText `json:"Contact,omitempty" db:"contact,false,jsonb"`                                               /* contact contact */
	FeeScheme                  types.JSONText `json:"FeeScheme,omitempty" db:"fee_scheme,false,jsonb"`                                          /* fee_scheme fee_scheme */
	CarServiceTarget           null.String    `json:"CarServiceTarget,omitempty" db:"car_service_target,false,character varying"`               /* car_service_target car_service_target */
	Policyholder               types.JSONText `json:"Policyholder,omitempty" db:"policyholder,false,jsonb"`                                     /* policyholder policyholder */
	FirstInsuredIDCardNo       null.String    `json:"FirstInsuredIDCardNo,omitempty" db:"first_insured_id_card_no,false,text"`                  /* first_insured_id_card_no first_insured_id_card_no */
	OrgAddr                    null.String    `json:"OrgAddr,omitempty" db:"org_addr,false,text"`                                               /* org_addr org_addr */
	OrgCreditCode              null.String    `json:"OrgCreditCode,omitempty" db:"org_credit_code,false,text"`                                  /* org_credit_code org_credit_code */
	OrgContact                 null.String    `json:"OrgContact,omitempty" db:"org_contact,false,text"`                                         /* org_contact org_contact */
	OrgPhone                   null.String    `json:"OrgPhone,omitempty" db:"org_phone,false,text"`                                             /* org_phone org_phone */
	OrgContactRole             null.String    `json:"OrgContactRole,omitempty" db:"org_contact_role,false,text"`                                /* org_contact_role org_contact_role */
	OrgCreditCodePic           null.String    `json:"OrgCreditCodePic,omitempty" db:"org_credit_code_pic,false,text"`                           /* org_credit_code_pic org_credit_code_pic */
	OrgSchoolCategory          null.String    `json:"OrgSchoolCategory,omitempty" db:"org_school_category,false,text"`                          /* org_school_category org_school_category */
	OrgCompulsoryStudentNum    null.String    `json:"OrgCompulsoryStudentNum,omitempty" db:"org_compulsory_student_num,false,text"`             /* org_compulsory_student_num org_compulsory_student_num */
	OrgNonCompulsoryStudentNum null.String    `json:"OrgNonCompulsoryStudentNum,omitempty" db:"org_non_compulsory_student_num,false,text"`      /* org_non_compulsory_student_num org_non_compulsory_student_num */
	PolicyholderType           null.String    `json:"PolicyholderType,omitempty" db:"policyholder_type,false,character varying"`                /* policyholder_type policyholder_type */
	PolicyholderID             null.Int       `json:"PolicyholderID,omitempty" db:"policyholder_id,false,bigint"`                               /* policyholder_id policyholder_id */
	Same                       null.Bool      `json:"Same,omitempty" db:"same,false,boolean"`                                                   /* same same */
	Relation                   null.String    `json:"Relation,omitempty" db:"relation,false,character varying"`                                 /* relation relation */
	Insured                    types.JSONText `json:"Insured,omitempty" db:"insured,false,jsonb"`                                               /* insured insured */
	InsuredName                null.String    `json:"InsuredName,omitempty" db:"insured_name,false,text"`                                       /* insured_name insured_name */
	InsuredProvince            null.String    `json:"InsuredProvince,omitempty" db:"insured_province,false,text"`                               /* insured_province insured_province */
	InsuredCity                null.String    `json:"InsuredCity,omitempty" db:"insured_city,false,text"`                                       /* insured_city insured_city */
	InsuredDistrict            null.String    `json:"InsuredDistrict,omitempty" db:"insured_district,false,text"`                               /* insured_district insured_district */
	InsuredIsCompulsory        null.Bool      `json:"InsuredIsCompulsory,omitempty" db:"insured_is_compulsory,false,boolean"`                   /* insured_is_compulsory insured_is_compulsory */
	InsuredCategory            null.String    `json:"InsuredCategory,omitempty" db:"insured_category,false,text"`                               /* insured_category insured_category */
	InsuredSchoolCategory      null.String    `json:"InsuredSchoolCategory,omitempty" db:"insured_school_category,false,text"`                  /* insured_school_category insured_school_category */
	InsuredPostCode            null.String    `json:"InsuredPostCode,omitempty" db:"insured_post_code,false,text"`                              /* insured_post_code insured_post_code */
	InsuredPhone               null.String    `json:"InsuredPhone,omitempty" db:"insured_phone,false,text"`                                     /* insured_phone insured_phone */
	PolicySchemeTitle          null.String    `json:"PolicySchemeTitle,omitempty" db:"policy_scheme_title,false,text"`                          /* policy_scheme_title policy_scheme_title */
	OrgBusinessDomain          types.JSONText `json:"OrgBusinessDomain,omitempty" db:"org_business_domain,false,jsonb"`                         /* org_business_domain org_business_domain */
	InsuredBusinessDomain      types.JSONText `json:"InsuredBusinessDomain,omitempty" db:"insured_business_domain,false,jsonb"`                 /* insured_business_domain insured_business_domain */
	InsuredID                  null.Int       `json:"InsuredID,omitempty" db:"insured_id,false,bigint"`                                         /* insured_id insured_id */
	HaveInsuredList            null.Bool      `json:"HaveInsuredList,omitempty" db:"have_insured_list,false,boolean"`                           /* have_insured_list have_insured_list */
	InsuredGroupByDay          null.Bool      `json:"InsuredGroupByDay,omitempty" db:"insured_group_by_day,false,boolean"`                      /* insured_group_by_day insured_group_by_day */
	InsuredType                null.String    `json:"InsuredType,omitempty" db:"insured_type,false,character varying"`                          /* insured_type insured_type */
	InsuredList                types.JSONText `json:"InsuredList,omitempty" db:"insured_list,false,jsonb"`                                      /* insured_list insured_list */
	CommenceDate               null.Int       `json:"CommenceDate,omitempty" db:"commence_date,false,bigint"`                                   /* commence_date commence_date */
	ExpiryDate                 null.Int       `json:"ExpiryDate,omitempty" db:"expiry_date,false,bigint"`                                       /* expiry_date expiry_date */
	Indate                     null.Int       `json:"Indate,omitempty" db:"indate,false,bigint"`                                                /* indate indate */
	Sign                       null.String    `json:"Sign,omitempty" db:"sign,false,character varying"`                                         /* sign sign */
	Jurisdiction               null.String    `json:"Jurisdiction,omitempty" db:"jurisdiction,false,character varying"`                         /* jurisdiction jurisdiction */
	DisputeHandling            null.String    `json:"DisputeHandling,omitempty" db:"dispute_handling,false,character varying"`                  /* dispute_handling dispute_handling */
	PrevPolicyNo               null.String    `json:"PrevPolicyNo,omitempty" db:"prev_policy_no,false,character varying"`                       /* prev_policy_no prev_policy_no */
	ReminderTimes              null.String    `json:"ReminderTimes,omitempty" db:"reminder_times,false,character varying"`                      /* reminder_times reminder_times */
	InsureBase                 null.String    `json:"InsureBase,omitempty" db:"insure_base,false,character varying"`                            /* insure_base insure_base */
	BlanketInsureCode          null.String    `json:"BlanketInsureCode,omitempty" db:"blanket_insure_code,false,character varying"`             /* blanket_insure_code blanket_insure_code */
	CustomType                 null.String    `json:"CustomType,omitempty" db:"custom_type,false,character varying"`                            /* custom_type custom_type */
	TrainProjects              null.String    `json:"TrainProjects,omitempty" db:"train_projects,false,character varying"`                      /* train_projects train_projects */
	BusinessLocations          types.JSONText `json:"BusinessLocations,omitempty" db:"business_locations,false,jsonb"`                          /* business_locations business_locations */
	TrainingPoolNum            null.Int       `json:"TrainingPoolNum,omitempty" db:"training_pool_num,false,smallint"`                          /* training_pool_num training_pool_num */
	HeatedPoolNum              null.Int       `json:"HeatedPoolNum,omitempty" db:"heated_pool_num,false,smallint"`                              /* heated_pool_num heated_pool_num */
	OpenPoolNum                null.Int       `json:"OpenPoolNum,omitempty" db:"open_pool_num,false,smallint"`                                  /* open_pool_num open_pool_num */
	PoolNum                    null.Int       `json:"PoolNum,omitempty" db:"pool_num,false,smallint"`                                           /* pool_num pool_num */
	DinnerNum                  null.Int       `json:"DinnerNum,omitempty" db:"dinner_num,false,integer"`                                        /* dinner_num dinner_num */
	CanteenNum                 null.Int       `json:"CanteenNum,omitempty" db:"canteen_num,false,integer"`                                      /* canteen_num canteen_num */
	ShopNum                    null.Int       `json:"ShopNum,omitempty" db:"shop_num,false,integer"`                                            /* shop_num shop_num */
	HaveRides                  null.Bool      `json:"HaveRides,omitempty" db:"have_rides,false,boolean"`                                        /* have_rides have_rides */
	HaveExplosive              null.Bool      `json:"HaveExplosive,omitempty" db:"have_explosive,false,boolean"`                                /* have_explosive have_explosive */
	Area                       null.Float     `json:"Area,omitempty" db:"area,false,double precision"`                                          /* area area */
	TrafficNum                 null.Int       `json:"TrafficNum,omitempty" db:"traffic_num,false,integer"`                                      /* traffic_num traffic_num */
	TemperatureType            null.String    `json:"TemperatureType,omitempty" db:"temperature_type,false,character varying"`                  /* temperature_type temperature_type */
	IsIndoor                   null.String    `json:"IsIndoor,omitempty" db:"is_indoor,false,character varying"`                                /* is_indoor is_indoor */
	Extra                      types.JSONText `json:"Extra,omitempty" db:"extra,false,jsonb"`                                                   /* extra extra */
	BankAccount                types.JSONText `json:"BankAccount,omitempty" db:"bank_account,false,jsonb"`                                      /* bank_account bank_account */
	PayContact                 null.String    `json:"PayContact,omitempty" db:"pay_contact,false,character varying"`                            /* pay_contact pay_contact */
	SuddenDeathTerms           null.String    `json:"SuddenDeathTerms,omitempty" db:"sudden_death_terms,false,character varying"`               /* sudden_death_terms sudden_death_terms */
	SpecAgreement              null.String    `json:"SpecAgreement,omitempty" db:"spec_agreement,false,character varying"`                      /* spec_agreement spec_agreement */
	InnerArea                  null.Float     `json:"InnerArea,omitempty" db:"inner_area,false,double precision"`                               /* inner_area inner_area */
	OuterArea                  null.Float     `json:"OuterArea,omitempty" db:"outer_area,false,double precision"`                               /* outer_area outer_area */
	PoolName                   null.String    `json:"PoolName,omitempty" db:"pool_name,false,character varying"`                                /* pool_name pool_name */
	ArbitralAgency             null.String    `json:"ArbitralAgency,omitempty" db:"arbitral_agency,false,character varying"`                    /* arbitral_agency arbitral_agency */
	ConfirmRefund              null.Bool      `json:"ConfirmRefund,omitempty" db:"confirm_refund,false,boolean"`                                /* confirm_refund confirm_refund */
	InsuredCreditCode          null.String    `json:"InsuredCreditCode,omitempty" db:"insured_credit_code,false,text"`                          /* insured_credit_code insured_credit_code */
	InsuredAddr                null.String    `json:"InsuredAddr,omitempty" db:"insured_addr,false,text"`                                       /* insured_addr insured_addr */
	Creator                    null.Int       `json:"Creator,omitempty" db:"creator,false,bigint"`                                              /* creator creator */
	DomainID                   null.Int       `json:"DomainID,omitempty" db:"domain_id,false,bigint"`                                           /* domain_id domain_id */
	OrderFiles                 types.JSONText `json:"OrderFiles,omitempty" db:"order_files,false,jsonb"`                                        /* order_files order_files */
	Addi                       types.JSONText `json:"Addi,omitempty" db:"addi,false,jsonb"`                                                     /* addi addi */
	Remark                     null.String    `json:"Remark,omitempty" db:"remark,false,character varying"`                                     /* remark remark */
	Status                     null.String    `json:"Status,omitempty" db:"status,false,character varying"`                                     /* status status */
	OrderStatus                null.String    `json:"OrderStatus,omitempty" db:"order_status,false,character varying"`                          /* order_status order_status */
	HaveNegotiatedPrice        null.Bool      `json:"HaveNegotiatedPrice,omitempty" db:"have_negotiated_price,false,boolean"`                   /* have_negotiated_price have_negotiated_price */
	LockStatus                 null.String    `json:"LockStatus,omitempty" db:"lock_status,false,character varying"`                            /* lock_status lock_status */
	InsuranceCompany           null.String    `json:"InsuranceCompany,omitempty" db:"insurance_company,false,character varying"`                /* insurance_company insurance_company */
	InsuranceCompanyAccount    null.String    `json:"InsuranceCompanyAccount,omitempty" db:"insurance_company_account,false,character varying"` /* insurance_company_account insurance_company_account */
	ActualAmount               null.Float     `json:"ActualAmount,omitempty" db:"actual_amount,false,double precision"`                         /* actual_amount actual_amount */
	CanRevokeOrder             null.Bool      `json:"CanRevokeOrder,omitempty" db:"can_revoke_order,false,boolean"`                             /* can_revoke_order can_revoke_order */
	CanPublicTransfers         null.Bool      `json:"CanPublicTransfers,omitempty" db:"can_public_transfers,false,boolean"`                     /* can_public_transfers can_public_transfers */
	IsReminder                 null.Bool      `json:"IsReminder,omitempty" db:"is_reminder,false,boolean"`                                      /* is_reminder is_reminder */
	Traits                     null.String    `json:"Traits,omitempty" db:"traits,false,character varying"`                                     /* traits traits */
	IsInvoice                  null.Bool      `json:"IsInvoice,omitempty" db:"is_invoice,false,boolean"`                                        /* is_invoice is_invoice */
	InvBorrow                  null.String    `json:"InvBorrow,omitempty" db:"inv_borrow,false,character varying"`                              /* inv_borrow inv_borrow */
	InvVisible                 null.String    `json:"InvVisible,omitempty" db:"inv_visible,false,character varying"`                            /* inv_visible inv_visible */
	InvTitle                   null.String    `json:"InvTitle,omitempty" db:"inv_title,false,character varying"`                                /* inv_title inv_title */
	InvStatus                  null.String    `json:"InvStatus,omitempty" db:"inv_status,false,character varying"`                              /* inv_status inv_status */
	UpdStatus                  null.String    `json:"UpdStatus,omitempty" db:"upd_status,false,character varying"`                              /* upd_status upd_status */
	DriverSeatNum              null.Int       `json:"DriverSeatNum,omitempty" db:"driver_seat_num,false,bigint"`                                /* driver_seat_num driver_seat_num */
	ApprovedPassengersNum      null.Int       `json:"ApprovedPassengersNum,omitempty" db:"approved_passengers_num,false,bigint"`                /* approved_passengers_num approved_passengers_num */
	RefusedReason              null.String    `json:"RefusedReason,omitempty" db:"refused_reason,false,character varying"`                      /* refused_reason refused_reason */
	UnpaidReason               null.String    `json:"UnpaidReason,omitempty" db:"unpaid_reason,false,character varying"`                        /* unpaid_reason unpaid_reason */
	Updator                    null.Int       `json:"Updator,omitempty" db:"updator,false,bigint"`                                              /* updator updator */
	UpdateTime                 null.Int       `json:"UpdateTime,omitempty" db:"update_time,false,bigint"`                                       /* update_time update_time */
	HaveRenewalReminder        null.Bool      `json:"HaveRenewalReminder,omitempty" db:"have_renewal_reminder,false,boolean"`                   /* have_renewal_reminder have_renewal_reminder */
	ChargeMode                 null.String    `json:"ChargeMode,omitempty" db:"charge_mode,false,character varying"`                            /* charge_mode charge_mode */
	AdminReceived              null.Bool      `json:"AdminReceived,omitempty" db:"admin_received,false,boolean"`                                /* admin_received admin_received */
	UserReceived               null.Bool      `json:"UserReceived,omitempty" db:"user_received,false,boolean"`                                  /* user_received user_received */
	HavePolicy                 null.Bool      `json:"HavePolicy,omitempty" db:"have_policy,false,boolean"`                                      /* have_policy have_policy */
	InsuranceTypeParentID      null.Int       `json:"InsuranceTypeParentID,omitempty" db:"insurance_type_parent_id,false,bigint"`               /* insurance_type_parent_id insurance_type_parent_id */
	CanInvTitleModify          null.Bool      `json:"CanInvTitleModify,omitempty" db:"can_inv_title_modify,false,boolean"`                      /* can_inv_title_modify can_inv_title_modify */
	InsuranceDisplay           null.String    `json:"InsuranceDisplay,omitempty" db:"insurance_display,false,character varying"`                /* insurance_display insurance_display */
	UserCorrectTimes           null.Int       `json:"UserCorrectTimes,omitempty" db:"user_correct_times,false,bigint"`                          /* user_correct_times user_correct_times */
	CorrectLevel               null.String    `json:"CorrectLevel,omitempty" db:"correct_level,false,character varying"`                        /* correct_level correct_level */
	OrgName                    null.String    `json:"OrgName,omitempty" db:"org_name,false,text"`                                               /* org_name org_name */
	OrgProvince                null.String    `json:"OrgProvince,omitempty" db:"org_province,false,text"`                                       /* org_province org_province */
	OrgCity                    null.String    `json:"OrgCity,omitempty" db:"org_city,false,text"`                                               /* org_city org_city */
	OrgDistrict                null.String    `json:"OrgDistrict,omitempty" db:"org_district,false,text"`                                       /* org_district org_district */
	OrgIsCompulsory            null.String    `json:"OrgIsCompulsory,omitempty" db:"org_is_compulsory,false,text"`                              /* org_is_compulsory org_is_compulsory */
	OrgIsSchool                null.String    `json:"OrgIsSchool,omitempty" db:"org_is_school,false,text"`                                      /* org_is_school org_is_school */
	ReminderTimesCount         null.Int       `json:"ReminderTimesCount,omitempty" db:"reminder_times_count,false,integer"`                     /* reminder_times_count reminder_times_count */
	IOfficialName              null.String    `json:"IOfficialName,omitempty" db:"i_official_name,false,character varying"`                     /* i_official_name i_official_name */
	IIDCardType                null.String    `json:"IIDCardType,omitempty" db:"i_id_card_type,false,character varying"`                        /* i_id_card_type i_id_card_type */
	IIDCardNo                  null.String    `json:"IIDCardNo,omitempty" db:"i_id_card_no,false,character varying"`                            /* i_id_card_no i_id_card_no */
	IMobilePhone               null.String    `json:"IMobilePhone,omitempty" db:"i_mobile_phone,false,character varying"`                       /* i_mobile_phone i_mobile_phone */
	IGender                    null.String    `json:"IGender,omitempty" db:"i_gender,false,character varying"`                                  /* i_gender i_gender */
	IBirthday                  null.Int       `json:"IBirthday,omitempty" db:"i_birthday,false,bigint"`                                         /* i_birthday i_birthday */
	IAddi                      types.JSONText `json:"IAddi,omitempty" db:"i_addi,false,jsonb"`                                                  /* i_addi i_addi */
	HOfficialName              null.String    `json:"HOfficialName,omitempty" db:"h_official_name,false,character varying"`                     /* h_official_name h_official_name */
	HIDCardType                null.String    `json:"HIDCardType,omitempty" db:"h_id_card_type,false,character varying"`                        /* h_id_card_type h_id_card_type */
	HIDCardNo                  null.String    `json:"HIDCardNo,omitempty" db:"h_id_card_no,false,character varying"`                            /* h_id_card_no h_id_card_no */
	HMobilePhone               null.String    `json:"HMobilePhone,omitempty" db:"h_mobile_phone,false,character varying"`                       /* h_mobile_phone h_mobile_phone */
	HAddi                      types.JSONText `json:"HAddi,omitempty" db:"h_addi,false,jsonb"`                                                  /* h_addi h_addi */
	Subdistrict                null.String    `json:"Subdistrict,omitempty" db:"subdistrict,false,character varying"`                           /* subdistrict subdistrict */
	Faculty                    null.String    `json:"Faculty,omitempty" db:"faculty,false,character varying"`                                   /* faculty faculty */
	Grade                      null.String    `json:"Grade,omitempty" db:"grade,false,character varying"`                                       /* grade grade */
	Class                      null.String    `json:"Class,omitempty" db:"class,false,character varying"`                                       /* class class */
	XCreateTime                null.Int       `json:"XCreateTime,omitempty" db:"x_create_time,false,bigint"`                                    /* x_create_time x_create_time */
	School                     null.String    `json:"School,omitempty" db:"school,false,character varying"`                                     /* school school */
	SFaculty                   types.JSONText `json:"SFaculty,omitempty" db:"s_faculty,false,jsonb"`                                            /* s_faculty s_faculty */
	SBranches                  types.JSONText `json:"SBranches,omitempty" db:"s_branches,false,jsonb"`                                          /* s_branches s_branches */
	SCategory                  null.String    `json:"SCategory,omitempty" db:"s_category,false,character varying"`                              /* s_category s_category */
	Province                   null.String    `json:"Province,omitempty" db:"province,false,character varying"`                                 /* province province */
	City                       null.String    `json:"City,omitempty" db:"city,false,character varying"`                                         /* city city */
	District                   null.String    `json:"District,omitempty" db:"district,false,character varying"`                                 /* district district */
	DataSyncTarget             null.String    `json:"DataSyncTarget,omitempty" db:"data_sync_target,false,character varying"`                   /* data_sync_target data_sync_target */
	SaleManagers               types.JSONText `json:"SaleManagers,omitempty" db:"sale_managers,false,jsonb"`                                    /* sale_managers sale_managers */
	SchoolManagers             types.JSONText `json:"SchoolManagers,omitempty" db:"school_managers,false,jsonb"`                                /* school_managers school_managers */
	PurchaseRule               types.JSONText `json:"PurchaseRule,omitempty" db:"purchase_rule,false,jsonb"`                                    /* purchase_rule purchase_rule */
	SCreateTime                null.Int       `json:"SCreateTime,omitempty" db:"s_create_time,false,bigint"`                                    /* s_create_time s_create_time */
	Difference                 null.Float     `json:"Difference,omitempty" db:"difference,false,double precision"`                              /* difference difference */
	PolicyNo                   null.String    `json:"PolicyNo,omitempty" db:"policy_no,false,text"`                                             /* policy_no policy_no */
	FeeStatus                  null.String    `json:"FeeStatus,omitempty" db:"fee_status,false,text"`                                           /* fee_status fee_status */
	Filter                                    // build DML where clause
}

//TVOrder2Fields full field list for default query
var TVOrder2Fields = []string{
	"ID",
	"TradeNo",
	"PayOrderNo",
	"InsureOrderNo",
	"Batch",
	"CreateTime",
	"PayTime",
	"PayChannel",
	"PayType",
	"PayName",
	"UnitPrice",
	"Amount",
	"Balance",
	"BalanceList",
	"OrgID",
	"HaveSuddenDeath",
	"GroundNum",
	"PlanType",
	"RemindersNum",
	"OrgManagerID",
	"InsuranceType",
	"HaveDinnerNum",
	"HaveConfirmDate",
	"InsuranceTypeID",
	"HealthSurvey",
	"PlanID",
	"PlanName",
	"Insurer",
	"PolicyScheme",
	"PolicyDoc",
	"ActivityName",
	"ActivityCategory",
	"ActivityDesc",
	"ActivityLocation",
	"ActivityDateSet",
	"CopiesNum",
	"InsuredCount",
	"CompulsoryStudentNum",
	"NonCompulsoryStudentNum",
	"Contact",
	"FeeScheme",
	"CarServiceTarget",
	"Policyholder",
	"FirstInsuredIDCardNo",
	"OrgAddr",
	"OrgCreditCode",
	"OrgContact",
	"OrgPhone",
	"OrgContactRole",
	"OrgCreditCodePic",
	"OrgSchoolCategory",
	"OrgCompulsoryStudentNum",
	"OrgNonCompulsoryStudentNum",
	"PolicyholderType",
	"PolicyholderID",
	"Same",
	"Relation",
	"Insured",
	"InsuredName",
	"InsuredProvince",
	"InsuredCity",
	"InsuredDistrict",
	"InsuredIsCompulsory",
	"InsuredCategory",
	"InsuredSchoolCategory",
	"InsuredPostCode",
	"InsuredPhone",
	"PolicySchemeTitle",
	"OrgBusinessDomain",
	"InsuredBusinessDomain",
	"InsuredID",
	"HaveInsuredList",
	"InsuredGroupByDay",
	"InsuredType",
	"InsuredList",
	"CommenceDate",
	"ExpiryDate",
	"Indate",
	"Sign",
	"Jurisdiction",
	"DisputeHandling",
	"PrevPolicyNo",
	"ReminderTimes",
	"InsureBase",
	"BlanketInsureCode",
	"CustomType",
	"TrainProjects",
	"BusinessLocations",
	"TrainingPoolNum",
	"HeatedPoolNum",
	"OpenPoolNum",
	"PoolNum",
	"DinnerNum",
	"CanteenNum",
	"ShopNum",
	"HaveRides",
	"HaveExplosive",
	"Area",
	"TrafficNum",
	"TemperatureType",
	"IsIndoor",
	"Extra",
	"BankAccount",
	"PayContact",
	"SuddenDeathTerms",
	"SpecAgreement",
	"InnerArea",
	"OuterArea",
	"PoolName",
	"ArbitralAgency",
	"ConfirmRefund",
	"InsuredCreditCode",
	"InsuredAddr",
	"Creator",
	"DomainID",
	"OrderFiles",
	"Addi",
	"Remark",
	"Status",
	"OrderStatus",
	"HaveNegotiatedPrice",
	"LockStatus",
	"InsuranceCompany",
	"InsuranceCompanyAccount",
	"ActualAmount",
	"CanRevokeOrder",
	"CanPublicTransfers",
	"IsReminder",
	"Traits",
	"IsInvoice",
	"InvBorrow",
	"InvVisible",
	"InvTitle",
	"InvStatus",
	"UpdStatus",
	"DriverSeatNum",
	"ApprovedPassengersNum",
	"RefusedReason",
	"UnpaidReason",
	"Updator",
	"UpdateTime",
	"HaveRenewalReminder",
	"ChargeMode",
	"AdminReceived",
	"UserReceived",
	"HavePolicy",
	"InsuranceTypeParentID",
	"CanInvTitleModify",
	"InsuranceDisplay",
	"UserCorrectTimes",
	"CorrectLevel",
	"OrgName",
	"OrgProvince",
	"OrgCity",
	"OrgDistrict",
	"OrgIsCompulsory",
	"OrgIsSchool",
	"ReminderTimesCount",
	"IOfficialName",
	"IIDCardType",
	"IIDCardNo",
	"IMobilePhone",
	"IGender",
	"IBirthday",
	"IAddi",
	"HOfficialName",
	"HIDCardType",
	"HIDCardNo",
	"HMobilePhone",
	"HAddi",
	"Subdistrict",
	"Faculty",
	"Grade",
	"Class",
	"XCreateTime",
	"School",
	"SFaculty",
	"SBranches",
	"SCategory",
	"Province",
	"City",
	"District",
	"DataSyncTarget",
	"SaleManagers",
	"SchoolManagers",
	"PurchaseRule",
	"SCreateTime",
	"Difference",
	"PolicyNo",
	"FeeStatus",
}

//Fields return all fields of struct.
func (r *TVOrder2) Fields() []string {
	return TVOrder2Fields
}

//GetTableName return the associated db table name.
func (r *TVOrder2) GetTableName() string {
	var viewNamePattern = regexp.MustCompile(`(?i)^t_v_[a-z0-9_]+$`)
	tableName := "t_v_order2"
	if viewNamePattern.MatchString(tableName) {
		return tableName[2:]
	}
	return tableName
}

// Create inserts the TVOrder2 to the database.
func (r *TVOrder2) Create(db Queryer) error {
	_, err := db.Exec(
		`INSERT INTO t_v_order2 (id, trade_no, pay_order_no, insure_order_no, batch, create_time, pay_time, pay_channel, pay_type, pay_name, unit_price, amount, balance, balance_list, org_id, have_sudden_death, ground_num, plan_type, reminders_num, org_manager_id, insurance_type, have_dinner_num, have_confirm_date, insurance_type_id, health_survey, plan_id, plan_name, insurer, policy_scheme, policy_doc, activity_name, activity_category, activity_desc, activity_location, activity_date_set, copies_num, insured_count, compulsory_student_num, non_compulsory_student_num, contact, fee_scheme, car_service_target, policyholder, first_insured_id_card_no, org_addr, org_credit_code, org_contact, org_phone, org_contact_role, org_credit_code_pic, org_school_category, org_compulsory_student_num, org_non_compulsory_student_num, policyholder_type, policyholder_id, same, relation, insured, insured_name, insured_province, insured_city, insured_district, insured_is_compulsory, insured_category, insured_school_category, insured_post_code, insured_phone, policy_scheme_title, org_business_domain, insured_business_domain, insured_id, have_insured_list, insured_group_by_day, insured_type, insured_list, commence_date, expiry_date, indate, sign, jurisdiction, dispute_handling, prev_policy_no, reminder_times, insure_base, blanket_insure_code, custom_type, train_projects, business_locations, training_pool_num, heated_pool_num, open_pool_num, pool_num, dinner_num, canteen_num, shop_num, have_rides, have_explosive, area, traffic_num, temperature_type, is_indoor, extra, bank_account, pay_contact, sudden_death_terms, spec_agreement, inner_area, outer_area, pool_name, arbitral_agency, confirm_refund, insured_credit_code, insured_addr, creator, domain_id, order_files, addi, remark, status, order_status, have_negotiated_price, lock_status, insurance_company, insurance_company_account, actual_amount, can_revoke_order, can_public_transfers, is_reminder, traits, is_invoice, inv_borrow, inv_visible, inv_title, inv_status, upd_status, driver_seat_num, approved_passengers_num, refused_reason, unpaid_reason, updator, update_time, have_renewal_reminder, charge_mode, admin_received, user_received, have_policy, insurance_type_parent_id, can_inv_title_modify, insurance_display, user_correct_times, correct_level, org_name, org_province, org_city, org_district, org_is_compulsory, org_is_school, reminder_times_count, i_official_name, i_id_card_type, i_id_card_no, i_mobile_phone, i_gender, i_birthday, i_addi, h_official_name, h_id_card_type, h_id_card_no, h_mobile_phone, h_addi, subdistrict, faculty, grade, class, x_create_time, school, s_faculty, s_branches, s_category, province, city, district, data_sync_target, sale_managers, school_managers, purchase_rule, s_create_time, difference, policy_no, fee_status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36, $37, $38, $39, $40, $41, $42, $43, $44, $45, $46, $47, $48, $49, $50, $51, $52, $53, $54, $55, $56, $57, $58, $59, $60, $61, $62, $63, $64, $65, $66, $67, $68, $69, $70, $71, $72, $73, $74, $75, $76, $77, $78, $79, $80, $81, $82, $83, $84, $85, $86, $87, $88, $89, $90, $91, $92, $93, $94, $95, $96, $97, $98, $99, $100, $101, $102, $103, $104, $105, $106, $107, $108, $109, $110, $111, $112, $113, $114, $115, $116, $117, $118, $119, $120, $121, $122, $123, $124, $125, $126, $127, $128, $129, $130, $131, $132, $133, $134, $135, $136, $137, $138, $139, $140, $141, $142, $143, $144, $145, $146, $147, $148, $149, $150, $151, $152, $153, $154, $155, $156, $157, $158, $159, $160, $161, $162, $163, $164, $165, $166, $167, $168, $169, $170, $171, $172, $173, $174, $175, $176, $177, $178, $179, $180, $181, $182, $183, $184, $185, $186, $187, $188, $189, $190)`,
		&r.ID, &r.TradeNo, &r.PayOrderNo, &r.InsureOrderNo, &r.Batch, &r.CreateTime, &r.PayTime, &r.PayChannel, &r.PayType, &r.PayName, &r.UnitPrice, &r.Amount, &r.Balance, &r.BalanceList, &r.OrgID, &r.HaveSuddenDeath, &r.GroundNum, &r.PlanType, &r.RemindersNum, &r.OrgManagerID, &r.InsuranceType, &r.HaveDinnerNum, &r.HaveConfirmDate, &r.InsuranceTypeID, &r.HealthSurvey, &r.PlanID, &r.PlanName, &r.Insurer, &r.PolicyScheme, &r.PolicyDoc, &r.ActivityName, &r.ActivityCategory, &r.ActivityDesc, &r.ActivityLocation, &r.ActivityDateSet, &r.CopiesNum, &r.InsuredCount, &r.CompulsoryStudentNum, &r.NonCompulsoryStudentNum, &r.Contact, &r.FeeScheme, &r.CarServiceTarget, &r.Policyholder, &r.FirstInsuredIDCardNo, &r.OrgAddr, &r.OrgCreditCode, &r.OrgContact, &r.OrgPhone, &r.OrgContactRole, &r.OrgCreditCodePic, &r.OrgSchoolCategory, &r.OrgCompulsoryStudentNum, &r.OrgNonCompulsoryStudentNum, &r.PolicyholderType, &r.PolicyholderID, &r.Same, &r.Relation, &r.Insured, &r.InsuredName, &r.InsuredProvince, &r.InsuredCity, &r.InsuredDistrict, &r.InsuredIsCompulsory, &r.InsuredCategory, &r.InsuredSchoolCategory, &r.InsuredPostCode, &r.InsuredPhone, &r.PolicySchemeTitle, &r.OrgBusinessDomain, &r.InsuredBusinessDomain, &r.InsuredID, &r.HaveInsuredList, &r.InsuredGroupByDay, &r.InsuredType, &r.InsuredList, &r.CommenceDate, &r.ExpiryDate, &r.Indate, &r.Sign, &r.Jurisdiction, &r.DisputeHandling, &r.PrevPolicyNo, &r.ReminderTimes, &r.InsureBase, &r.BlanketInsureCode, &r.CustomType, &r.TrainProjects, &r.BusinessLocations, &r.TrainingPoolNum, &r.HeatedPoolNum, &r.OpenPoolNum, &r.PoolNum, &r.DinnerNum, &r.CanteenNum, &r.ShopNum, &r.HaveRides, &r.HaveExplosive, &r.Area, &r.TrafficNum, &r.TemperatureType, &r.IsIndoor, &r.Extra, &r.BankAccount, &r.PayContact, &r.SuddenDeathTerms, &r.SpecAgreement, &r.InnerArea, &r.OuterArea, &r.PoolName, &r.ArbitralAgency, &r.ConfirmRefund, &r.InsuredCreditCode, &r.InsuredAddr, &r.Creator, &r.DomainID, &r.OrderFiles, &r.Addi, &r.Remark, &r.Status, &r.OrderStatus, &r.HaveNegotiatedPrice, &r.LockStatus, &r.InsuranceCompany, &r.InsuranceCompanyAccount, &r.ActualAmount, &r.CanRevokeOrder, &r.CanPublicTransfers, &r.IsReminder, &r.Traits, &r.IsInvoice, &r.InvBorrow, &r.InvVisible, &r.InvTitle, &r.InvStatus, &r.UpdStatus, &r.DriverSeatNum, &r.ApprovedPassengersNum, &r.RefusedReason, &r.UnpaidReason, &r.Updator, &r.UpdateTime, &r.HaveRenewalReminder, &r.ChargeMode, &r.AdminReceived, &r.UserReceived, &r.HavePolicy, &r.InsuranceTypeParentID, &r.CanInvTitleModify, &r.InsuranceDisplay, &r.UserCorrectTimes, &r.CorrectLevel, &r.OrgName, &r.OrgProvince, &r.OrgCity, &r.OrgDistrict, &r.OrgIsCompulsory, &r.OrgIsSchool, &r.ReminderTimesCount, &r.IOfficialName, &r.IIDCardType, &r.IIDCardNo, &r.IMobilePhone, &r.IGender, &r.IBirthday, &r.IAddi, &r.HOfficialName, &r.HIDCardType, &r.HIDCardNo, &r.HMobilePhone, &r.HAddi, &r.Subdistrict, &r.Faculty, &r.Grade, &r.Class, &r.XCreateTime, &r.School, &r.SFaculty, &r.SBranches, &r.SCategory, &r.Province, &r.City, &r.District, &r.DataSyncTarget, &r.SaleManagers, &r.SchoolManagers, &r.PurchaseRule, &r.SCreateTime, &r.Difference, &r.PolicyNo, &r.FeeStatus)
	if err != nil {
		return errors.Wrap(err, "failed to insert t_v_order2")
	}
	return nil
}

// GetTVOrder2ByPk select the TVOrder2 from the database.
func GetTVOrder2ByPk(db Queryer) (*TVOrder2, error) {
	var r TVOrder2
	err := db.QueryRow(
		`SELECT id, trade_no, pay_order_no, insure_order_no, batch, create_time, pay_time, pay_channel, pay_type, pay_name, unit_price, amount, balance, balance_list, org_id, have_sudden_death, ground_num, plan_type, reminders_num, org_manager_id, insurance_type, have_dinner_num, have_confirm_date, insurance_type_id, health_survey, plan_id, plan_name, insurer, policy_scheme, policy_doc, activity_name, activity_category, activity_desc, activity_location, activity_date_set, copies_num, insured_count, compulsory_student_num, non_compulsory_student_num, contact, fee_scheme, car_service_target, policyholder, first_insured_id_card_no, org_addr, org_credit_code, org_contact, org_phone, org_contact_role, org_credit_code_pic, org_school_category, org_compulsory_student_num, org_non_compulsory_student_num, policyholder_type, policyholder_id, same, relation, insured, insured_name, insured_province, insured_city, insured_district, insured_is_compulsory, insured_category, insured_school_category, insured_post_code, insured_phone, policy_scheme_title, org_business_domain, insured_business_domain, insured_id, have_insured_list, insured_group_by_day, insured_type, insured_list, commence_date, expiry_date, indate, sign, jurisdiction, dispute_handling, prev_policy_no, reminder_times, insure_base, blanket_insure_code, custom_type, train_projects, business_locations, training_pool_num, heated_pool_num, open_pool_num, pool_num, dinner_num, canteen_num, shop_num, have_rides, have_explosive, area, traffic_num, temperature_type, is_indoor, extra, bank_account, pay_contact, sudden_death_terms, spec_agreement, inner_area, outer_area, pool_name, arbitral_agency, confirm_refund, insured_credit_code, insured_addr, creator, domain_id, order_files, addi, remark, status, order_status, have_negotiated_price, lock_status, insurance_company, insurance_company_account, actual_amount, can_revoke_order, can_public_transfers, is_reminder, traits, is_invoice, inv_borrow, inv_visible, inv_title, inv_status, upd_status, driver_seat_num, approved_passengers_num, refused_reason, unpaid_reason, updator, update_time, have_renewal_reminder, charge_mode, admin_received, user_received, have_policy, insurance_type_parent_id, can_inv_title_modify, insurance_display, user_correct_times, correct_level, org_name, org_province, org_city, org_district, org_is_compulsory, org_is_school, reminder_times_count, i_official_name, i_id_card_type, i_id_card_no, i_mobile_phone, i_gender, i_birthday, i_addi, h_official_name, h_id_card_type, h_id_card_no, h_mobile_phone, h_addi, subdistrict, faculty, grade, class, x_create_time, school, s_faculty, s_branches, s_category, province, city, district, data_sync_target, sale_managers, school_managers, purchase_rule, s_create_time, difference, policy_no, fee_status FROM t_v_order2 WHERE `,
	).Scan(&r.ID, &r.TradeNo, &r.PayOrderNo, &r.InsureOrderNo, &r.Batch, &r.CreateTime, &r.PayTime, &r.PayChannel, &r.PayType, &r.PayName, &r.UnitPrice, &r.Amount, &r.Balance, &r.BalanceList, &r.OrgID, &r.HaveSuddenDeath, &r.GroundNum, &r.PlanType, &r.RemindersNum, &r.OrgManagerID, &r.InsuranceType, &r.HaveDinnerNum, &r.HaveConfirmDate, &r.InsuranceTypeID, &r.HealthSurvey, &r.PlanID, &r.PlanName, &r.Insurer, &r.PolicyScheme, &r.PolicyDoc, &r.ActivityName, &r.ActivityCategory, &r.ActivityDesc, &r.ActivityLocation, &r.ActivityDateSet, &r.CopiesNum, &r.InsuredCount, &r.CompulsoryStudentNum, &r.NonCompulsoryStudentNum, &r.Contact, &r.FeeScheme, &r.CarServiceTarget, &r.Policyholder, &r.FirstInsuredIDCardNo, &r.OrgAddr, &r.OrgCreditCode, &r.OrgContact, &r.OrgPhone, &r.OrgContactRole, &r.OrgCreditCodePic, &r.OrgSchoolCategory, &r.OrgCompulsoryStudentNum, &r.OrgNonCompulsoryStudentNum, &r.PolicyholderType, &r.PolicyholderID, &r.Same, &r.Relation, &r.Insured, &r.InsuredName, &r.InsuredProvince, &r.InsuredCity, &r.InsuredDistrict, &r.InsuredIsCompulsory, &r.InsuredCategory, &r.InsuredSchoolCategory, &r.InsuredPostCode, &r.InsuredPhone, &r.PolicySchemeTitle, &r.OrgBusinessDomain, &r.InsuredBusinessDomain, &r.InsuredID, &r.HaveInsuredList, &r.InsuredGroupByDay, &r.InsuredType, &r.InsuredList, &r.CommenceDate, &r.ExpiryDate, &r.Indate, &r.Sign, &r.Jurisdiction, &r.DisputeHandling, &r.PrevPolicyNo, &r.ReminderTimes, &r.InsureBase, &r.BlanketInsureCode, &r.CustomType, &r.TrainProjects, &r.BusinessLocations, &r.TrainingPoolNum, &r.HeatedPoolNum, &r.OpenPoolNum, &r.PoolNum, &r.DinnerNum, &r.CanteenNum, &r.ShopNum, &r.HaveRides, &r.HaveExplosive, &r.Area, &r.TrafficNum, &r.TemperatureType, &r.IsIndoor, &r.Extra, &r.BankAccount, &r.PayContact, &r.SuddenDeathTerms, &r.SpecAgreement, &r.InnerArea, &r.OuterArea, &r.PoolName, &r.ArbitralAgency, &r.ConfirmRefund, &r.InsuredCreditCode, &r.InsuredAddr, &r.Creator, &r.DomainID, &r.OrderFiles, &r.Addi, &r.Remark, &r.Status, &r.OrderStatus, &r.HaveNegotiatedPrice, &r.LockStatus, &r.InsuranceCompany, &r.InsuranceCompanyAccount, &r.ActualAmount, &r.CanRevokeOrder, &r.CanPublicTransfers, &r.IsReminder, &r.Traits, &r.IsInvoice, &r.InvBorrow, &r.InvVisible, &r.InvTitle, &r.InvStatus, &r.UpdStatus, &r.DriverSeatNum, &r.ApprovedPassengersNum, &r.RefusedReason, &r.UnpaidReason, &r.Updator, &r.UpdateTime, &r.HaveRenewalReminder, &r.ChargeMode, &r.AdminReceived, &r.UserReceived, &r.HavePolicy, &r.InsuranceTypeParentID, &r.CanInvTitleModify, &r.InsuranceDisplay, &r.UserCorrectTimes, &r.CorrectLevel, &r.OrgName, &r.OrgProvince, &r.OrgCity, &r.OrgDistrict, &r.OrgIsCompulsory, &r.OrgIsSchool, &r.ReminderTimesCount, &r.IOfficialName, &r.IIDCardType, &r.IIDCardNo, &r.IMobilePhone, &r.IGender, &r.IBirthday, &r.IAddi, &r.HOfficialName, &r.HIDCardType, &r.HIDCardNo, &r.HMobilePhone, &r.HAddi, &r.Subdistrict, &r.Faculty, &r.Grade, &r.Class, &r.XCreateTime, &r.School, &r.SFaculty, &r.SBranches, &r.SCategory, &r.Province, &r.City, &r.District, &r.DataSyncTarget, &r.SaleManagers, &r.SchoolManagers, &r.PurchaseRule, &r.SCreateTime, &r.Difference, &r.PolicyNo, &r.FeeStatus)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select t_v_order2")
	}
	return &r, nil
}

/*TVOrderSum t_v_order_sum represents kuser.t_v_order_sum */
type TVOrderSum struct {
	School       null.String `json:"School,omitempty" db:"school,false,character varying"`             /* school school */
	OrgID        null.Int    `json:"OrgID,omitempty" db:"org_id,false,bigint"`                         /* org_id org_id */
	Batch        null.String `json:"Batch,omitempty" db:"batch,false,character varying"`               /* batch batch */
	OrderNumber  null.Int    `json:"OrderNumber,omitempty" db:"order_number,false,bigint"`             /* order_number order_number */
	OrderAmount  null.Float  `json:"OrderAmount,omitempty" db:"order_amount,false,double precision"`   /* order_amount order_amount */
	CancelNumber null.Int    `json:"CancelNumber,omitempty" db:"cancel_number,false,bigint"`           /* cancel_number cancel_number */
	CancelAmount null.Float  `json:"CancelAmount,omitempty" db:"cancel_amount,false,double precision"` /* cancel_amount cancel_amount */
	Filter                   // build DML where clause
}

//TVOrderSumFields full field list for default query
var TVOrderSumFields = []string{
	"School",
	"OrgID",
	"Batch",
	"OrderNumber",
	"OrderAmount",
	"CancelNumber",
	"CancelAmount",
}

//Fields return all fields of struct.
func (r *TVOrderSum) Fields() []string {
	return TVOrderSumFields
}

//GetTableName return the associated db table name.
func (r *TVOrderSum) GetTableName() string {
	var viewNamePattern = regexp.MustCompile(`(?i)^t_v_[a-z0-9_]+$`)
	tableName := "t_v_order_sum"
	if viewNamePattern.MatchString(tableName) {
		return tableName[2:]
	}
	return tableName
}

// Create inserts the TVOrderSum to the database.
func (r *TVOrderSum) Create(db Queryer) error {
	_, err := db.Exec(
		`INSERT INTO t_v_order_sum (school, org_id, batch, order_number, order_amount, cancel_number, cancel_amount) VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		&r.School, &r.OrgID, &r.Batch, &r.OrderNumber, &r.OrderAmount, &r.CancelNumber, &r.CancelAmount)
	if err != nil {
		return errors.Wrap(err, "failed to insert t_v_order_sum")
	}
	return nil
}

// GetTVOrderSumByPk select the TVOrderSum from the database.
func GetTVOrderSumByPk(db Queryer) (*TVOrderSum, error) {
	var r TVOrderSum
	err := db.QueryRow(
		`SELECT school, org_id, batch, order_number, order_amount, cancel_number, cancel_amount FROM t_v_order_sum WHERE `,
	).Scan(&r.School, &r.OrgID, &r.Batch, &r.OrderNumber, &r.OrderAmount, &r.CancelNumber, &r.CancelAmount)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select t_v_order_sum")
	}
	return &r, nil
}

/*TVParam t_v_param represents kuser.t_v_param */
type TVParam struct {
	ParentID   null.Int    `json:"ParentID,omitempty" db:"parent_id,false,bigint"`                /* parent_id parent_id */
	ParentName null.String `json:"ParentName,omitempty" db:"parent_name,false,character varying"` /* parent_name parent_name */
	ID         null.Int    `json:"ID,omitempty" db:"id,false,integer"`                            /* id id */
	Name       null.String `json:"Name,omitempty" db:"name,false,character varying"`              /* name name */
	DataType   null.String `json:"DataType,omitempty" db:"data_type,false,character varying"`     /* data_type data_type */
	Value      null.String `json:"Value,omitempty" db:"value,false,character varying"`            /* value value */
	Remark     null.String `json:"Remark,omitempty" db:"remark,false,character varying"`          /* remark remark */
	Status     null.String `json:"Status,omitempty" db:"status,false,character varying"`          /* status status */
	Filter                 // build DML where clause
}

//TVParamFields full field list for default query
var TVParamFields = []string{
	"ParentID",
	"ParentName",
	"ID",
	"Name",
	"DataType",
	"Value",
	"Remark",
	"Status",
}

//Fields return all fields of struct.
func (r *TVParam) Fields() []string {
	return TVParamFields
}

//GetTableName return the associated db table name.
func (r *TVParam) GetTableName() string {
	var viewNamePattern = regexp.MustCompile(`(?i)^t_v_[a-z0-9_]+$`)
	tableName := "t_v_param"
	if viewNamePattern.MatchString(tableName) {
		return tableName[2:]
	}
	return tableName
}

// Create inserts the TVParam to the database.
func (r *TVParam) Create(db Queryer) error {
	_, err := db.Exec(
		`INSERT INTO t_v_param (parent_id, parent_name, id, name, data_type, value, remark, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		&r.ParentID, &r.ParentName, &r.ID, &r.Name, &r.DataType, &r.Value, &r.Remark, &r.Status)
	if err != nil {
		return errors.Wrap(err, "failed to insert t_v_param")
	}
	return nil
}

// GetTVParamByPk select the TVParam from the database.
func GetTVParamByPk(db Queryer) (*TVParam, error) {
	var r TVParam
	err := db.QueryRow(
		`SELECT parent_id, parent_name, id, name, data_type, value, remark, status FROM t_v_param WHERE `,
	).Scan(&r.ParentID, &r.ParentName, &r.ID, &r.Name, &r.DataType, &r.Value, &r.Remark, &r.Status)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select t_v_param")
	}
	return &r, nil
}

/*TVPayment t_v_payment represents kuser.t_v_payment */
type TVPayment struct {
	ID                null.Int       `json:"ID,omitempty" db:"id,false,integer"`                                          /* id id */
	Batch             null.String    `json:"Batch,omitempty" db:"batch,false,character varying"`                          /* batch batch */
	PolicyNo          null.String    `json:"PolicyNo,omitempty" db:"policy_no,false,character varying"`                   /* policy_no policy_no */
	TransferNo        null.String    `json:"TransferNo,omitempty" db:"transfer_no,false,character varying"`               /* transfer_no transfer_no */
	TransferAmount    null.Float     `json:"TransferAmount,omitempty" db:"transfer_amount,false,double precision"`        /* transfer_amount transfer_amount */
	Creator           null.Int       `json:"Creator,omitempty" db:"creator,false,bigint"`                                 /* creator creator */
	CreateTime        null.Int       `json:"CreateTime,omitempty" db:"create_time,false,bigint"`                          /* create_time create_time */
	Regenerator       null.Int       `json:"Regenerator,omitempty" db:"regenerator,false,bigint"`                         /* regenerator regenerator */
	UpdateTime        null.Int       `json:"UpdateTime,omitempty" db:"update_time,false,bigint"`                          /* update_time update_time */
	DomainID          null.Int       `json:"DomainID,omitempty" db:"domain_id,false,bigint"`                              /* domain_id domain_id */
	Addi              types.JSONText `json:"Addi,omitempty" db:"addi,false,jsonb"`                                        /* addi addi */
	Remark            null.String    `json:"Remark,omitempty" db:"remark,false,character varying"`                        /* remark remark */
	Status            null.String    `json:"Status,omitempty" db:"status,false,character varying"`                        /* status status */
	IsAdminPay        null.Bool      `json:"IsAdminPay,omitempty" db:"is_admin_pay,false,boolean"`                        /* is_admin_pay is_admin_pay */
	Premium           null.Float     `json:"Premium,omitempty" db:"premium,false,double precision"`                       /* premium premium */
	ThirdPartyPremium null.Float     `json:"ThirdPartyPremium,omitempty" db:"third_party_premium,false,double precision"` /* third_party_premium third_party_premium */
	PolicyholderName  null.String    `json:"PolicyholderName,omitempty" db:"policyholder_name,false,text"`                /* policyholder_name policyholder_name */
	Filter                           // build DML where clause
}

//TVPaymentFields full field list for default query
var TVPaymentFields = []string{
	"ID",
	"Batch",
	"PolicyNo",
	"TransferNo",
	"TransferAmount",
	"Creator",
	"CreateTime",
	"Regenerator",
	"UpdateTime",
	"DomainID",
	"Addi",
	"Remark",
	"Status",
	"IsAdminPay",
	"Premium",
	"ThirdPartyPremium",
	"PolicyholderName",
}

//Fields return all fields of struct.
func (r *TVPayment) Fields() []string {
	return TVPaymentFields
}

//GetTableName return the associated db table name.
func (r *TVPayment) GetTableName() string {
	var viewNamePattern = regexp.MustCompile(`(?i)^t_v_[a-z0-9_]+$`)
	tableName := "t_v_payment"
	if viewNamePattern.MatchString(tableName) {
		return tableName[2:]
	}
	return tableName
}

// Create inserts the TVPayment to the database.
func (r *TVPayment) Create(db Queryer) error {
	_, err := db.Exec(
		`INSERT INTO t_v_payment (id, batch, policy_no, transfer_no, transfer_amount, creator, create_time, regenerator, update_time, domain_id, addi, remark, status, is_admin_pay, premium, third_party_premium, policyholder_name) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)`,
		&r.ID, &r.Batch, &r.PolicyNo, &r.TransferNo, &r.TransferAmount, &r.Creator, &r.CreateTime, &r.Regenerator, &r.UpdateTime, &r.DomainID, &r.Addi, &r.Remark, &r.Status, &r.IsAdminPay, &r.Premium, &r.ThirdPartyPremium, &r.PolicyholderName)
	if err != nil {
		return errors.Wrap(err, "failed to insert t_v_payment")
	}
	return nil
}

// GetTVPaymentByPk select the TVPayment from the database.
func GetTVPaymentByPk(db Queryer) (*TVPayment, error) {
	var r TVPayment
	err := db.QueryRow(
		`SELECT id, batch, policy_no, transfer_no, transfer_amount, creator, create_time, regenerator, update_time, domain_id, addi, remark, status, is_admin_pay, premium, third_party_premium, policyholder_name FROM t_v_payment WHERE `,
	).Scan(&r.ID, &r.Batch, &r.PolicyNo, &r.TransferNo, &r.TransferAmount, &r.Creator, &r.CreateTime, &r.Regenerator, &r.UpdateTime, &r.DomainID, &r.Addi, &r.Remark, &r.Status, &r.IsAdminPay, &r.Premium, &r.ThirdPartyPremium, &r.PolicyholderName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select t_v_payment")
	}
	return &r, nil
}

/*TVRegion t_v_region represents kuser.t_v_region */
type TVRegion struct {
	Province null.String `json:"Province,omitempty" db:"province,false,character varying"` /* province province */
	City     null.String `json:"City,omitempty" db:"city,false,character varying"`         /* city city */
	District null.String `json:"District,omitempty" db:"district,false,character varying"` /* district district */
	Street   null.String `json:"Street,omitempty" db:"street,false,character varying"`     /* street street */
	Filter               // build DML where clause
}

//TVRegionFields full field list for default query
var TVRegionFields = []string{
	"Province",
	"City",
	"District",
	"Street",
}

//Fields return all fields of struct.
func (r *TVRegion) Fields() []string {
	return TVRegionFields
}

//GetTableName return the associated db table name.
func (r *TVRegion) GetTableName() string {
	var viewNamePattern = regexp.MustCompile(`(?i)^t_v_[a-z0-9_]+$`)
	tableName := "t_v_region"
	if viewNamePattern.MatchString(tableName) {
		return tableName[2:]
	}
	return tableName
}

// Create inserts the TVRegion to the database.
func (r *TVRegion) Create(db Queryer) error {
	_, err := db.Exec(
		`INSERT INTO t_v_region (province, city, district, street) VALUES ($1, $2, $3, $4)`,
		&r.Province, &r.City, &r.District, &r.Street)
	if err != nil {
		return errors.Wrap(err, "failed to insert t_v_region")
	}
	return nil
}

// GetTVRegionByPk select the TVRegion from the database.
func GetTVRegionByPk(db Queryer) (*TVRegion, error) {
	var r TVRegion
	err := db.QueryRow(
		`SELECT province, city, district, street FROM t_v_region WHERE `,
	).Scan(&r.Province, &r.City, &r.District, &r.Street)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select t_v_region")
	}
	return &r, nil
}

/*TVReportClaims t_v_report_claims represents kuser.t_v_report_claims */
type TVReportClaims struct {
	ID                       null.Int       `json:"ID,omitempty" db:"id,false,integer"`                                             /* id id */
	InformantID              null.Int       `json:"InformantID,omitempty" db:"informant_id,false,bigint"`                           /* informant_id informant_id */
	Informant                types.JSONText `json:"Informant,omitempty" db:"informant,false,jsonb"`                                 /* informant informant */
	InsuredID                null.Int       `json:"InsuredID,omitempty" db:"insured_id,false,bigint"`                               /* insured_id insured_id */
	Insured                  types.JSONText `json:"Insured,omitempty" db:"insured,false,jsonb"`                                     /* insured insured */
	InsuranceType            null.String    `json:"InsuranceType,omitempty" db:"insurance_type,false,character varying"`            /* insurance_type insurance_type */
	InsurancePolicySn        null.String    `json:"InsurancePolicySn,omitempty" db:"insurance_policy_sn,false,character varying"`   /* insurance_policy_sn insurance_policy_sn */
	InsurancePolicyID        null.Int       `json:"InsurancePolicyID,omitempty" db:"insurance_policy_id,false,bigint"`              /* insurance_policy_id insurance_policy_id */
	OrderID                  null.Int       `json:"OrderID,omitempty" db:"order_id,false,bigint"`                                   /* order_id order_id */
	OrgID                    null.Int       `json:"OrgID,omitempty" db:"org_id,false,bigint"`                                       /* org_id org_id */
	PlanID                   null.Int       `json:"PlanID,omitempty" db:"plan_id,false,bigint"`                                     /* plan_id plan_id */
	PolicyNo                 null.String    `json:"PolicyNo,omitempty" db:"policy_no,false,character varying"`                      /* policy_no policy_no */
	InsurancePolicyStart     null.Int       `json:"InsurancePolicyStart,omitempty" db:"insurance_policy_start,false,bigint"`        /* insurance_policy_start insurance_policy_start */
	InsurancePolicyCease     null.Int       `json:"InsurancePolicyCease,omitempty" db:"insurance_policy_cease,false,bigint"`        /* insurance_policy_cease insurance_policy_cease */
	ReportSn                 null.String    `json:"ReportSn,omitempty" db:"report_sn,false,character varying"`                      /* report_sn report_sn */
	InsuredChannel           null.String    `json:"InsuredChannel,omitempty" db:"insured_channel,false,character varying"`          /* insured_channel insured_channel */
	InsuredOrg               null.String    `json:"InsuredOrg,omitempty" db:"insured_org,false,character varying"`                  /* insured_org insured_org */
	Treatment                null.String    `json:"Treatment,omitempty" db:"treatment,false,character varying"`                     /* treatment treatment */
	Hospital                 null.String    `json:"Hospital,omitempty" db:"hospital,false,character varying"`                       /* hospital hospital */
	InjuredLocation          null.String    `json:"InjuredLocation,omitempty" db:"injured_location,false,character varying"`        /* injured_location injured_location */
	InjuredPart              null.String    `json:"InjuredPart,omitempty" db:"injured_part,false,character varying"`                /* injured_part injured_part */
	Reason                   null.String    `json:"Reason,omitempty" db:"reason,false,character varying"`                           /* reason reason */
	InjuredDesc              null.String    `json:"InjuredDesc,omitempty" db:"injured_desc,false,character varying"`                /* injured_desc injured_desc */
	CreditCode               null.String    `json:"CreditCode,omitempty" db:"credit_code,false,character varying"`                  /* credit_code credit_code */
	BankAccountType          null.String    `json:"BankAccountType,omitempty" db:"bank_account_type,false,character varying"`       /* bank_account_type bank_account_type */
	BankAccountName          null.String    `json:"BankAccountName,omitempty" db:"bank_account_name,false,character varying"`       /* bank_account_name bank_account_name */
	BankName                 null.String    `json:"BankName,omitempty" db:"bank_name,false,character varying"`                      /* bank_name bank_name */
	BankAccountID            null.String    `json:"BankAccountID,omitempty" db:"bank_account_id,false,character varying"`           /* bank_account_id bank_account_id */
	BankCardPic              types.JSONText `json:"BankCardPic,omitempty" db:"bank_card_pic,false,jsonb"`                           /* bank_card_pic bank_card_pic */
	InjuredIDPic             types.JSONText `json:"InjuredIDPic,omitempty" db:"injured_id_pic,false,jsonb"`                         /* injured_id_pic injured_id_pic */
	GuardianIDPic            types.JSONText `json:"GuardianIDPic,omitempty" db:"guardian_id_pic,false,jsonb"`                       /* guardian_id_pic guardian_id_pic */
	OrgLicPic                types.JSONText `json:"OrgLicPic,omitempty" db:"org_lic_pic,false,jsonb"`                               /* org_lic_pic org_lic_pic */
	RelationProvePic         types.JSONText `json:"RelationProvePic,omitempty" db:"relation_prove_pic,false,jsonb"`                 /* relation_prove_pic relation_prove_pic */
	BillsPic                 types.JSONText `json:"BillsPic,omitempty" db:"bills_pic,false,jsonb"`                                  /* bills_pic bills_pic */
	HospitalizedBillsPic     types.JSONText `json:"HospitalizedBillsPic,omitempty" db:"hospitalized_bills_pic,false,jsonb"`         /* hospitalized_bills_pic hospitalized_bills_pic */
	InvoicePic               types.JSONText `json:"InvoicePic,omitempty" db:"invoice_pic,false,jsonb"`                              /* invoice_pic invoice_pic */
	MedicalRecordPic         types.JSONText `json:"MedicalRecordPic,omitempty" db:"medical_record_pic,false,jsonb"`                 /* medical_record_pic medical_record_pic */
	DignosticInspectionPic   types.JSONText `json:"DignosticInspectionPic,omitempty" db:"dignostic_inspection_pic,false,jsonb"`     /* dignostic_inspection_pic dignostic_inspection_pic */
	DischargeAbstractPic     types.JSONText `json:"DischargeAbstractPic,omitempty" db:"discharge_abstract_pic,false,jsonb"`         /* discharge_abstract_pic discharge_abstract_pic */
	OtherPic                 types.JSONText `json:"OtherPic,omitempty" db:"other_pic,false,jsonb"`                                  /* other_pic other_pic */
	CourierSnPic             types.JSONText `json:"CourierSnPic,omitempty" db:"courier_sn_pic,false,jsonb"`                         /* courier_sn_pic courier_sn_pic */
	PaidNoticePic            types.JSONText `json:"PaidNoticePic,omitempty" db:"paid_notice_pic,false,jsonb"`                       /* paid_notice_pic paid_notice_pic */
	ClaimApplyPic            types.JSONText `json:"ClaimApplyPic,omitempty" db:"claim_apply_pic,false,jsonb"`                       /* claim_apply_pic claim_apply_pic */
	EquityTransferFile       types.JSONText `json:"EquityTransferFile,omitempty" db:"equity_transfer_file,false,jsonb"`             /* equity_transfer_file equity_transfer_file */
	MatchProgrammePic        types.JSONText `json:"MatchProgrammePic,omitempty" db:"match_programme_pic,false,jsonb"`               /* match_programme_pic match_programme_pic */
	PolicyFile               types.JSONText `json:"PolicyFile,omitempty" db:"policy_file,false,jsonb"`                              /* policy_file policy_file */
	AddiPic                  types.JSONText `json:"AddiPic,omitempty" db:"addi_pic,false,jsonb"`                                    /* addi_pic addi_pic */
	CourierSn                null.String    `json:"CourierSn,omitempty" db:"courier_sn,false,character varying"`                    /* courier_sn courier_sn */
	ReplyAddr                null.String    `json:"ReplyAddr,omitempty" db:"reply_addr,false,character varying"`                    /* reply_addr reply_addr */
	InjuredTime              null.Int       `json:"InjuredTime,omitempty" db:"injured_time,false,bigint"`                           /* injured_time injured_time */
	ReportTime               null.Int       `json:"ReportTime,omitempty" db:"report_time,false,bigint"`                             /* report_time report_time */
	ReplyTime                null.Int       `json:"ReplyTime,omitempty" db:"reply_time,false,bigint"`                               /* reply_time reply_time */
	ClaimsMatAddTime         null.Int       `json:"ClaimsMatAddTime,omitempty" db:"claims_mat_add_time,false,bigint"`               /* claims_mat_add_time claims_mat_add_time */
	MatReturnDate            null.Int       `json:"MatReturnDate,omitempty" db:"mat_return_date,false,bigint"`                      /* mat_return_date mat_return_date */
	CloseDate                null.Int       `json:"CloseDate,omitempty" db:"close_date,false,bigint"`                               /* close_date close_date */
	FaceAmount               null.Float     `json:"FaceAmount,omitempty" db:"face_amount,false,double precision"`                   /* face_amount face_amount */
	MediAssureAmount         null.Float     `json:"MediAssureAmount,omitempty" db:"medi_assure_amount,false,double precision"`      /* medi_assure_amount medi_assure_amount */
	ThirdPayAmount           null.Float     `json:"ThirdPayAmount,omitempty" db:"third_pay_amount,false,double precision"`          /* third_pay_amount third_pay_amount */
	ClaimAmount              null.Float     `json:"ClaimAmount,omitempty" db:"claim_amount,false,double precision"`                 /* claim_amount claim_amount */
	RefuseDesc               null.String    `json:"RefuseDesc,omitempty" db:"refuse_desc,false,character varying"`                  /* refuse_desc refuse_desc */
	Addi                     types.JSONText `json:"Addi,omitempty" db:"addi,false,jsonb"`                                           /* addi addi */
	Creator                  null.Int       `json:"Creator,omitempty" db:"creator,false,bigint"`                                    /* creator creator */
	DomainID                 null.Int       `json:"DomainID,omitempty" db:"domain_id,false,bigint"`                                 /* domain_id domain_id */
	Remark                   null.String    `json:"Remark,omitempty" db:"remark,false,character varying"`                           /* remark remark */
	Status                   null.String    `json:"Status,omitempty" db:"status,false,character varying"`                           /* status status */
	School                   null.String    `json:"School,omitempty" db:"school,false,character varying"`                           /* school school */
	SchoolType               null.String    `json:"SchoolType,omitempty" db:"school_type,false,character varying"`                  /* school_type school_type */
	Grade                    null.String    `json:"Grade,omitempty" db:"grade,false,character varying"`                             /* grade grade */
	Class                    null.String    `json:"Class,omitempty" db:"class,false,character varying"`                             /* class class */
	OfficialName             null.String    `json:"OfficialName,omitempty" db:"official_name,false,character varying"`              /* official_name official_name */
	Gender                   null.String    `json:"Gender,omitempty" db:"gender,false,character varying"`                           /* gender gender */
	IDCardType               null.String    `json:"IDCardType,omitempty" db:"id_card_type,false,character varying"`                 /* id_card_type id_card_type */
	IDCardNo                 null.String    `json:"IDCardNo,omitempty" db:"id_card_no,false,character varying"`                     /* id_card_no id_card_no */
	Birthday                 null.Int       `json:"Birthday,omitempty" db:"birthday,false,bigint"`                                  /* birthday birthday */
	WOffcialName             null.String    `json:"WOffcialName,omitempty" db:"w_offcial_name,false,text"`                          /* w_offcial_name w_offcial_name */
	WIDCardType              null.String    `json:"WIDCardType,omitempty" db:"w_id_card_type,false,text"`                           /* w_id_card_type w_id_card_type */
	WIDCardNo                null.String    `json:"WIDCardNo,omitempty" db:"w_id_card_no,false,text"`                               /* w_id_card_no w_id_card_no */
	InsuredOffcialName       null.String    `json:"InsuredOffcialName,omitempty" db:"insured_offcial_name,false,text"`              /* insured_offcial_name insured_offcial_name */
	MOffcialName             null.String    `json:"MOffcialName,omitempty" db:"m_offcial_name,false,character varying"`             /* m_offcial_name m_offcial_name */
	MGender                  null.String    `json:"MGender,omitempty" db:"m_gender,false,character varying"`                        /* m_gender m_gender */
	MIDCardType              null.String    `json:"MIDCardType,omitempty" db:"m_id_card_type,false,character varying"`              /* m_id_card_type m_id_card_type */
	MIDCardNo                null.String    `json:"MIDCardNo,omitempty" db:"m_id_card_no,false,character varying"`                  /* m_id_card_no m_id_card_no */
	MMobilePhone             null.String    `json:"MMobilePhone,omitempty" db:"m_mobile_phone,false,character varying"`             /* m_mobile_phone m_mobile_phone */
	WMOffcialName            null.String    `json:"WMOffcialName,omitempty" db:"w_m_offcial_name,false,text"`                       /* w_m_offcial_name w_m_offcial_name */
	WMMobilePhone            null.String    `json:"WMMobilePhone,omitempty" db:"w_m_mobile_phone,false,text"`                       /* w_m_mobile_phone w_m_mobile_phone */
	InsuranceTypeID          null.Int       `json:"InsuranceTypeID,omitempty" db:"insurance_type_id,false,bigint"`                  /* insurance_type_id insurance_type_id */
	InsuranceTypeParentID    null.Int       `json:"InsuranceTypeParentID,omitempty" db:"insurance_type_parent_id,false,bigint"`     /* insurance_type_parent_id insurance_type_parent_id */
	Sn                       null.String    `json:"Sn,omitempty" db:"sn,false,character varying"`                                   /* sn sn */
	OccurrReason             null.String    `json:"OccurrReason,omitempty" db:"occurr_reason,false,character varying"`              /* occurr_reason occurr_reason */
	TreatmentResult          null.String    `json:"TreatmentResult,omitempty" db:"treatment_result,false,character varying"`        /* treatment_result treatment_result */
	DiseaseDiagnosisPic      types.JSONText `json:"DiseaseDiagnosisPic,omitempty" db:"disease_diagnosis_pic,false,jsonb"`           /* disease_diagnosis_pic disease_diagnosis_pic */
	DisabilityCertificate    types.JSONText `json:"DisabilityCertificate,omitempty" db:"disability_certificate,false,jsonb"`        /* disability_certificate disability_certificate */
	DeathCertificate         types.JSONText `json:"DeathCertificate,omitempty" db:"death_certificate,false,jsonb"`                  /* death_certificate death_certificate */
	StudentStatusCertificate types.JSONText `json:"StudentStatusCertificate,omitempty" db:"student_status_certificate,false,jsonb"` /* student_status_certificate student_status_certificate */
	Filter                                  // build DML where clause
}

//TVReportClaimsFields full field list for default query
var TVReportClaimsFields = []string{
	"ID",
	"InformantID",
	"Informant",
	"InsuredID",
	"Insured",
	"InsuranceType",
	"InsurancePolicySn",
	"InsurancePolicyID",
	"OrderID",
	"OrgID",
	"PlanID",
	"PolicyNo",
	"InsurancePolicyStart",
	"InsurancePolicyCease",
	"ReportSn",
	"InsuredChannel",
	"InsuredOrg",
	"Treatment",
	"Hospital",
	"InjuredLocation",
	"InjuredPart",
	"Reason",
	"InjuredDesc",
	"CreditCode",
	"BankAccountType",
	"BankAccountName",
	"BankName",
	"BankAccountID",
	"BankCardPic",
	"InjuredIDPic",
	"GuardianIDPic",
	"OrgLicPic",
	"RelationProvePic",
	"BillsPic",
	"HospitalizedBillsPic",
	"InvoicePic",
	"MedicalRecordPic",
	"DignosticInspectionPic",
	"DischargeAbstractPic",
	"OtherPic",
	"CourierSnPic",
	"PaidNoticePic",
	"ClaimApplyPic",
	"EquityTransferFile",
	"MatchProgrammePic",
	"PolicyFile",
	"AddiPic",
	"CourierSn",
	"ReplyAddr",
	"InjuredTime",
	"ReportTime",
	"ReplyTime",
	"ClaimsMatAddTime",
	"MatReturnDate",
	"CloseDate",
	"FaceAmount",
	"MediAssureAmount",
	"ThirdPayAmount",
	"ClaimAmount",
	"RefuseDesc",
	"Addi",
	"Creator",
	"DomainID",
	"Remark",
	"Status",
	"School",
	"SchoolType",
	"Grade",
	"Class",
	"OfficialName",
	"Gender",
	"IDCardType",
	"IDCardNo",
	"Birthday",
	"WOffcialName",
	"WIDCardType",
	"WIDCardNo",
	"InsuredOffcialName",
	"MOffcialName",
	"MGender",
	"MIDCardType",
	"MIDCardNo",
	"MMobilePhone",
	"WMOffcialName",
	"WMMobilePhone",
	"InsuranceTypeID",
	"InsuranceTypeParentID",
	"Sn",
	"OccurrReason",
	"TreatmentResult",
	"DiseaseDiagnosisPic",
	"DisabilityCertificate",
	"DeathCertificate",
	"StudentStatusCertificate",
}

//Fields return all fields of struct.
func (r *TVReportClaims) Fields() []string {
	return TVReportClaimsFields
}

//GetTableName return the associated db table name.
func (r *TVReportClaims) GetTableName() string {
	var viewNamePattern = regexp.MustCompile(`(?i)^t_v_[a-z0-9_]+$`)
	tableName := "t_v_report_claims"
	if viewNamePattern.MatchString(tableName) {
		return tableName[2:]
	}
	return tableName
}

// Create inserts the TVReportClaims to the database.
func (r *TVReportClaims) Create(db Queryer) error {
	_, err := db.Exec(
		`INSERT INTO t_v_report_claims (id, informant_id, informant, insured_id, insured, insurance_type, insurance_policy_sn, insurance_policy_id, order_id, org_id, plan_id, policy_no, insurance_policy_start, insurance_policy_cease, report_sn, insured_channel, insured_org, treatment, hospital, injured_location, injured_part, reason, injured_desc, credit_code, bank_account_type, bank_account_name, bank_name, bank_account_id, bank_card_pic, injured_id_pic, guardian_id_pic, org_lic_pic, relation_prove_pic, bills_pic, hospitalized_bills_pic, invoice_pic, medical_record_pic, dignostic_inspection_pic, discharge_abstract_pic, other_pic, courier_sn_pic, paid_notice_pic, claim_apply_pic, equity_transfer_file, match_programme_pic, policy_file, addi_pic, courier_sn, reply_addr, injured_time, report_time, reply_time, claims_mat_add_time, mat_return_date, close_date, face_amount, medi_assure_amount, third_pay_amount, claim_amount, refuse_desc, addi, creator, domain_id, remark, status, school, school_type, grade, class, official_name, gender, id_card_type, id_card_no, birthday, w_offcial_name, w_id_card_type, w_id_card_no, insured_offcial_name, m_offcial_name, m_gender, m_id_card_type, m_id_card_no, m_mobile_phone, w_m_offcial_name, w_m_mobile_phone, insurance_type_id, insurance_type_parent_id, sn, occurr_reason, treatment_result, disease_diagnosis_pic, disability_certificate, death_certificate, student_status_certificate) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36, $37, $38, $39, $40, $41, $42, $43, $44, $45, $46, $47, $48, $49, $50, $51, $52, $53, $54, $55, $56, $57, $58, $59, $60, $61, $62, $63, $64, $65, $66, $67, $68, $69, $70, $71, $72, $73, $74, $75, $76, $77, $78, $79, $80, $81, $82, $83, $84, $85, $86, $87, $88, $89, $90, $91, $92, $93, $94)`,
		&r.ID, &r.InformantID, &r.Informant, &r.InsuredID, &r.Insured, &r.InsuranceType, &r.InsurancePolicySn, &r.InsurancePolicyID, &r.OrderID, &r.OrgID, &r.PlanID, &r.PolicyNo, &r.InsurancePolicyStart, &r.InsurancePolicyCease, &r.ReportSn, &r.InsuredChannel, &r.InsuredOrg, &r.Treatment, &r.Hospital, &r.InjuredLocation, &r.InjuredPart, &r.Reason, &r.InjuredDesc, &r.CreditCode, &r.BankAccountType, &r.BankAccountName, &r.BankName, &r.BankAccountID, &r.BankCardPic, &r.InjuredIDPic, &r.GuardianIDPic, &r.OrgLicPic, &r.RelationProvePic, &r.BillsPic, &r.HospitalizedBillsPic, &r.InvoicePic, &r.MedicalRecordPic, &r.DignosticInspectionPic, &r.DischargeAbstractPic, &r.OtherPic, &r.CourierSnPic, &r.PaidNoticePic, &r.ClaimApplyPic, &r.EquityTransferFile, &r.MatchProgrammePic, &r.PolicyFile, &r.AddiPic, &r.CourierSn, &r.ReplyAddr, &r.InjuredTime, &r.ReportTime, &r.ReplyTime, &r.ClaimsMatAddTime, &r.MatReturnDate, &r.CloseDate, &r.FaceAmount, &r.MediAssureAmount, &r.ThirdPayAmount, &r.ClaimAmount, &r.RefuseDesc, &r.Addi, &r.Creator, &r.DomainID, &r.Remark, &r.Status, &r.School, &r.SchoolType, &r.Grade, &r.Class, &r.OfficialName, &r.Gender, &r.IDCardType, &r.IDCardNo, &r.Birthday, &r.WOffcialName, &r.WIDCardType, &r.WIDCardNo, &r.InsuredOffcialName, &r.MOffcialName, &r.MGender, &r.MIDCardType, &r.MIDCardNo, &r.MMobilePhone, &r.WMOffcialName, &r.WMMobilePhone, &r.InsuranceTypeID, &r.InsuranceTypeParentID, &r.Sn, &r.OccurrReason, &r.TreatmentResult, &r.DiseaseDiagnosisPic, &r.DisabilityCertificate, &r.DeathCertificate, &r.StudentStatusCertificate)
	if err != nil {
		return errors.Wrap(err, "failed to insert t_v_report_claims")
	}
	return nil
}

// GetTVReportClaimsByPk select the TVReportClaims from the database.
func GetTVReportClaimsByPk(db Queryer) (*TVReportClaims, error) {
	var r TVReportClaims
	err := db.QueryRow(
		`SELECT id, informant_id, informant, insured_id, insured, insurance_type, insurance_policy_sn, insurance_policy_id, order_id, org_id, plan_id, policy_no, insurance_policy_start, insurance_policy_cease, report_sn, insured_channel, insured_org, treatment, hospital, injured_location, injured_part, reason, injured_desc, credit_code, bank_account_type, bank_account_name, bank_name, bank_account_id, bank_card_pic, injured_id_pic, guardian_id_pic, org_lic_pic, relation_prove_pic, bills_pic, hospitalized_bills_pic, invoice_pic, medical_record_pic, dignostic_inspection_pic, discharge_abstract_pic, other_pic, courier_sn_pic, paid_notice_pic, claim_apply_pic, equity_transfer_file, match_programme_pic, policy_file, addi_pic, courier_sn, reply_addr, injured_time, report_time, reply_time, claims_mat_add_time, mat_return_date, close_date, face_amount, medi_assure_amount, third_pay_amount, claim_amount, refuse_desc, addi, creator, domain_id, remark, status, school, school_type, grade, class, official_name, gender, id_card_type, id_card_no, birthday, w_offcial_name, w_id_card_type, w_id_card_no, insured_offcial_name, m_offcial_name, m_gender, m_id_card_type, m_id_card_no, m_mobile_phone, w_m_offcial_name, w_m_mobile_phone, insurance_type_id, insurance_type_parent_id, sn, occurr_reason, treatment_result, disease_diagnosis_pic, disability_certificate, death_certificate, student_status_certificate FROM t_v_report_claims WHERE `,
	).Scan(&r.ID, &r.InformantID, &r.Informant, &r.InsuredID, &r.Insured, &r.InsuranceType, &r.InsurancePolicySn, &r.InsurancePolicyID, &r.OrderID, &r.OrgID, &r.PlanID, &r.PolicyNo, &r.InsurancePolicyStart, &r.InsurancePolicyCease, &r.ReportSn, &r.InsuredChannel, &r.InsuredOrg, &r.Treatment, &r.Hospital, &r.InjuredLocation, &r.InjuredPart, &r.Reason, &r.InjuredDesc, &r.CreditCode, &r.BankAccountType, &r.BankAccountName, &r.BankName, &r.BankAccountID, &r.BankCardPic, &r.InjuredIDPic, &r.GuardianIDPic, &r.OrgLicPic, &r.RelationProvePic, &r.BillsPic, &r.HospitalizedBillsPic, &r.InvoicePic, &r.MedicalRecordPic, &r.DignosticInspectionPic, &r.DischargeAbstractPic, &r.OtherPic, &r.CourierSnPic, &r.PaidNoticePic, &r.ClaimApplyPic, &r.EquityTransferFile, &r.MatchProgrammePic, &r.PolicyFile, &r.AddiPic, &r.CourierSn, &r.ReplyAddr, &r.InjuredTime, &r.ReportTime, &r.ReplyTime, &r.ClaimsMatAddTime, &r.MatReturnDate, &r.CloseDate, &r.FaceAmount, &r.MediAssureAmount, &r.ThirdPayAmount, &r.ClaimAmount, &r.RefuseDesc, &r.Addi, &r.Creator, &r.DomainID, &r.Remark, &r.Status, &r.School, &r.SchoolType, &r.Grade, &r.Class, &r.OfficialName, &r.Gender, &r.IDCardType, &r.IDCardNo, &r.Birthday, &r.WOffcialName, &r.WIDCardType, &r.WIDCardNo, &r.InsuredOffcialName, &r.MOffcialName, &r.MGender, &r.MIDCardType, &r.MIDCardNo, &r.MMobilePhone, &r.WMOffcialName, &r.WMMobilePhone, &r.InsuranceTypeID, &r.InsuranceTypeParentID, &r.Sn, &r.OccurrReason, &r.TreatmentResult, &r.DiseaseDiagnosisPic, &r.DisabilityCertificate, &r.DeathCertificate, &r.StudentStatusCertificate)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select t_v_report_claims")
	}
	return &r, nil
}

/*TVUserDomain t_v_user_domain represents kuser.t_v_user_domain */
type TVUserDomain struct {
	ID             null.Int       `json:"ID,omitempty" db:"id,false,integer"`                                     /* id id */
	UserID         null.Int       `json:"UserID,omitempty" db:"user_id,false,integer"`                            /* user_id user_id */
	UserName       null.String    `json:"UserName,omitempty" db:"user_name,false,character varying"`              /* user_name user_name */
	MobilePhone    null.String    `json:"MobilePhone,omitempty" db:"mobile_phone,false,character varying"`        /* mobile_phone mobile_phone */
	Email          null.String    `json:"Email,omitempty" db:"email,false,character varying"`                     /* email email */
	IDCardNo       null.String    `json:"IDCardNo,omitempty" db:"id_card_no,false,character varying"`             /* id_card_no id_card_no */
	IDCardType     null.String    `json:"IDCardType,omitempty" db:"id_card_type,false,character varying"`         /* id_card_type id_card_type */
	ExternalID     null.String    `json:"ExternalID,omitempty" db:"external_id,false,character varying"`          /* external_id external_id */
	ExternalIDType null.String    `json:"ExternalIDType,omitempty" db:"external_id_type,false,character varying"` /* external_id_type external_id_type */
	AuthDomainID   null.Int       `json:"AuthDomainID,omitempty" db:"auth_domain_id,false,integer"`               /* auth_domain_id auth_domain_id */
	Priority       null.Int       `json:"Priority,omitempty" db:"priority,false,smallint"`                        /* priority priority */
	DomainName     null.String    `json:"DomainName,omitempty" db:"domain_name,false,character varying"`          /* domain_name domain_name */
	Domain         null.String    `json:"Domain,omitempty" db:"domain,false,character varying"`                   /* domain domain */
	GrantSource    null.String    `json:"GrantSource,omitempty" db:"grant_source,false,character varying"`        /* grant_source grant_source */
	DataAccessMode null.String    `json:"DataAccessMode,omitempty" db:"data_access_mode,false,character varying"` /* data_access_mode data_access_mode */
	DataScope      types.JSONText `json:"DataScope,omitempty" db:"data_scope,false,jsonb"`                        /* data_scope data_scope */
	DomainID       null.Int       `json:"DomainID,omitempty" db:"domain_id,false,bigint"`                         /* domain_id domain_id */
	CreateTime     null.Int       `json:"CreateTime,omitempty" db:"create_time,false,bigint"`                     /* create_time create_time */
	Remark         null.String    `json:"Remark,omitempty" db:"remark,false,character varying"`                   /* remark remark */
	Addi           types.JSONText `json:"Addi,omitempty" db:"addi,false,jsonb"`                                   /* addi addi */
	Creator        null.Int       `json:"Creator,omitempty" db:"creator,false,bigint"`                            /* creator creator */
	Status         null.String    `json:"Status,omitempty" db:"status,false,character varying"`                   /* status status */
	Filter                        // build DML where clause
}

//TVUserDomainFields full field list for default query
var TVUserDomainFields = []string{
	"ID",
	"UserID",
	"UserName",
	"MobilePhone",
	"Email",
	"IDCardNo",
	"IDCardType",
	"ExternalID",
	"ExternalIDType",
	"AuthDomainID",
	"Priority",
	"DomainName",
	"Domain",
	"GrantSource",
	"DataAccessMode",
	"DataScope",
	"DomainID",
	"CreateTime",
	"Remark",
	"Addi",
	"Creator",
	"Status",
}

//Fields return all fields of struct.
func (r *TVUserDomain) Fields() []string {
	return TVUserDomainFields
}

//GetTableName return the associated db table name.
func (r *TVUserDomain) GetTableName() string {
	var viewNamePattern = regexp.MustCompile(`(?i)^t_v_[a-z0-9_]+$`)
	tableName := "t_v_user_domain"
	if viewNamePattern.MatchString(tableName) {
		return tableName[2:]
	}
	return tableName
}

// Create inserts the TVUserDomain to the database.
func (r *TVUserDomain) Create(db Queryer) error {
	_, err := db.Exec(
		`INSERT INTO t_v_user_domain (id, user_id, user_name, mobile_phone, email, id_card_no, id_card_type, external_id, external_id_type, auth_domain_id, priority, domain_name, domain, grant_source, data_access_mode, data_scope, domain_id, create_time, remark, addi, creator, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22)`,
		&r.ID, &r.UserID, &r.UserName, &r.MobilePhone, &r.Email, &r.IDCardNo, &r.IDCardType, &r.ExternalID, &r.ExternalIDType, &r.AuthDomainID, &r.Priority, &r.DomainName, &r.Domain, &r.GrantSource, &r.DataAccessMode, &r.DataScope, &r.DomainID, &r.CreateTime, &r.Remark, &r.Addi, &r.Creator, &r.Status)
	if err != nil {
		return errors.Wrap(err, "failed to insert t_v_user_domain")
	}
	return nil
}

// GetTVUserDomainByPk select the TVUserDomain from the database.
func GetTVUserDomainByPk(db Queryer) (*TVUserDomain, error) {
	var r TVUserDomain
	err := db.QueryRow(
		`SELECT id, user_id, user_name, mobile_phone, email, id_card_no, id_card_type, external_id, external_id_type, auth_domain_id, priority, domain_name, domain, grant_source, data_access_mode, data_scope, domain_id, create_time, remark, addi, creator, status FROM t_v_user_domain WHERE `,
	).Scan(&r.ID, &r.UserID, &r.UserName, &r.MobilePhone, &r.Email, &r.IDCardNo, &r.IDCardType, &r.ExternalID, &r.ExternalIDType, &r.AuthDomainID, &r.Priority, &r.DomainName, &r.Domain, &r.GrantSource, &r.DataAccessMode, &r.DataScope, &r.DomainID, &r.CreateTime, &r.Remark, &r.Addi, &r.Creator, &r.Status)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select t_v_user_domain")
	}
	return &r, nil
}

/*TVUserDomainAPI t_v_user_domain_api represents kuser.t_v_user_domain_api */
type TVUserDomainAPI struct {
	UserID                   null.Int       `json:"UserID,omitempty" db:"user_id,false,integer"`                                                  /* user_id user_id */
	OfficialName             null.String    `json:"OfficialName,omitempty" db:"official_name,false,character varying"`                            /* official_name official_name */
	UserName                 null.String    `json:"UserName,omitempty" db:"user_name,false,character varying"`                                    /* user_name user_name */
	Role                     null.Int       `json:"Role,omitempty" db:"role,false,bigint"`                                                        /* role role */
	MobilePhone              null.String    `json:"MobilePhone,omitempty" db:"mobile_phone,false,character varying"`                              /* mobile_phone mobile_phone */
	APIID                    null.Int       `json:"APIID,omitempty" db:"api_id,false,integer"`                                                    /* api_id api_id */
	APIName                  null.String    `json:"APIName,omitempty" db:"api_name,false,character varying"`                                      /* api_name api_name */
	APIExposePath            null.String    `json:"APIExposePath,omitempty" db:"api_expose_path,false,character varying"`                         /* api_expose_path api_expose_path */
	DomainName               null.String    `json:"DomainName,omitempty" db:"domain_name,false,character varying"`                                /* domain_name domain_name */
	DomainID                 null.Int       `json:"DomainID,omitempty" db:"domain_id,false,integer"`                                              /* domain_id domain_id */
	Domain                   null.String    `json:"Domain,omitempty" db:"domain,false,character varying"`                                         /* domain domain */
	Priority                 null.Int       `json:"Priority,omitempty" db:"priority,false,smallint"`                                              /* priority priority */
	UserDomainID             null.Int       `json:"UserDomainID,omitempty" db:"user_domain_id,false,integer"`                                     /* user_domain_id user_domain_id */
	UserDomainGrantSource    null.String    `json:"UserDomainGrantSource,omitempty" db:"user_domain_grant_source,false,character varying"`        /* user_domain_grant_source user_domain_grant_source */
	UserDomainDataAccessMode null.String    `json:"UserDomainDataAccessMode,omitempty" db:"user_domain_data_access_mode,false,character varying"` /* user_domain_data_access_mode user_domain_data_access_mode */
	UserDomainDataScope      types.JSONText `json:"UserDomainDataScope,omitempty" db:"user_domain_data_scope,false,jsonb"`                        /* user_domain_data_scope user_domain_data_scope */
	UserDomainDataScopeData  null.String    `json:"UserDomainDataScopeData,omitempty" db:"user_domain_data_scope_data,false,text"`                /* user_domain_data_scope_data user_domain_data_scope_data */
	UserDomainDataScopeType  null.String    `json:"UserDomainDataScopeType,omitempty" db:"user_domain_data_scope_type,false,text"`                /* user_domain_data_scope_type user_domain_data_scope_type */
	IDOnDomain               null.String    `json:"IDOnDomain,omitempty" db:"id_on_domain,false,character varying"`                               /* id_on_domain id_on_domain */
	UserDomainCreateTime     null.Int       `json:"UserDomainCreateTime,omitempty" db:"user_domain_create_time,false,bigint"`                     /* user_domain_create_time user_domain_create_time */
	DomainAPIID              null.Int       `json:"DomainAPIID,omitempty" db:"domain_api_id,false,integer"`                                       /* domain_api_id domain_api_id */
	DomainAPIGrantSource     null.String    `json:"DomainAPIGrantSource,omitempty" db:"domain_api_grant_source,false,character varying"`          /* domain_api_grant_source domain_api_grant_source */
	DomainAPIDataAccessMode  null.String    `json:"DomainAPIDataAccessMode,omitempty" db:"domain_api_data_access_mode,false,character varying"`   /* domain_api_data_access_mode domain_api_data_access_mode */
	DomainAPIDataScope       types.JSONText `json:"DomainAPIDataScope,omitempty" db:"domain_api_data_scope,false,jsonb"`                          /* domain_api_data_scope domain_api_data_scope */
	DomainAPIDataScopeData   null.String    `json:"DomainAPIDataScopeData,omitempty" db:"domain_api_data_scope_data,false,text"`                  /* domain_api_data_scope_data domain_api_data_scope_data */
	DomainAPIDataScopeType   null.String    `json:"DomainAPIDataScopeType,omitempty" db:"domain_api_data_scope_type,false,text"`                  /* domain_api_data_scope_type domain_api_data_scope_type */
	DomainAPICreateTime      null.Int       `json:"DomainAPICreateTime,omitempty" db:"domain_api_create_time,false,bigint"`                       /* domain_api_create_time domain_api_create_time */
	Filter                                  // build DML where clause
}

//TVUserDomainAPIFields full field list for default query
var TVUserDomainAPIFields = []string{
	"UserID",
	"OfficialName",
	"UserName",
	"Role",
	"MobilePhone",
	"APIID",
	"APIName",
	"APIExposePath",
	"DomainName",
	"DomainID",
	"Domain",
	"Priority",
	"UserDomainID",
	"UserDomainGrantSource",
	"UserDomainDataAccessMode",
	"UserDomainDataScope",
	"UserDomainDataScopeData",
	"UserDomainDataScopeType",
	"IDOnDomain",
	"UserDomainCreateTime",
	"DomainAPIID",
	"DomainAPIGrantSource",
	"DomainAPIDataAccessMode",
	"DomainAPIDataScope",
	"DomainAPIDataScopeData",
	"DomainAPIDataScopeType",
	"DomainAPICreateTime",
}

//Fields return all fields of struct.
func (r *TVUserDomainAPI) Fields() []string {
	return TVUserDomainAPIFields
}

//GetTableName return the associated db table name.
func (r *TVUserDomainAPI) GetTableName() string {
	var viewNamePattern = regexp.MustCompile(`(?i)^t_v_[a-z0-9_]+$`)
	tableName := "t_v_user_domain_api"
	if viewNamePattern.MatchString(tableName) {
		return tableName[2:]
	}
	return tableName
}

// Create inserts the TVUserDomainAPI to the database.
func (r *TVUserDomainAPI) Create(db Queryer) error {
	_, err := db.Exec(
		`INSERT INTO t_v_user_domain_api (user_id, official_name, user_name, role, mobile_phone, api_id, api_name, api_expose_path, domain_name, domain_id, domain, priority, user_domain_id, user_domain_grant_source, user_domain_data_access_mode, user_domain_data_scope, user_domain_data_scope_data, user_domain_data_scope_type, id_on_domain, user_domain_create_time, domain_api_id, domain_api_grant_source, domain_api_data_access_mode, domain_api_data_scope, domain_api_data_scope_data, domain_api_data_scope_type, domain_api_create_time) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27)`,
		&r.UserID, &r.OfficialName, &r.UserName, &r.Role, &r.MobilePhone, &r.APIID, &r.APIName, &r.APIExposePath, &r.DomainName, &r.DomainID, &r.Domain, &r.Priority, &r.UserDomainID, &r.UserDomainGrantSource, &r.UserDomainDataAccessMode, &r.UserDomainDataScope, &r.UserDomainDataScopeData, &r.UserDomainDataScopeType, &r.IDOnDomain, &r.UserDomainCreateTime, &r.DomainAPIID, &r.DomainAPIGrantSource, &r.DomainAPIDataAccessMode, &r.DomainAPIDataScope, &r.DomainAPIDataScopeData, &r.DomainAPIDataScopeType, &r.DomainAPICreateTime)
	if err != nil {
		return errors.Wrap(err, "failed to insert t_v_user_domain_api")
	}
	return nil
}

// GetTVUserDomainAPIByPk select the TVUserDomainAPI from the database.
func GetTVUserDomainAPIByPk(db Queryer) (*TVUserDomainAPI, error) {
	var r TVUserDomainAPI
	err := db.QueryRow(
		`SELECT user_id, official_name, user_name, role, mobile_phone, api_id, api_name, api_expose_path, domain_name, domain_id, domain, priority, user_domain_id, user_domain_grant_source, user_domain_data_access_mode, user_domain_data_scope, user_domain_data_scope_data, user_domain_data_scope_type, id_on_domain, user_domain_create_time, domain_api_id, domain_api_grant_source, domain_api_data_access_mode, domain_api_data_scope, domain_api_data_scope_data, domain_api_data_scope_type, domain_api_create_time FROM t_v_user_domain_api WHERE `,
	).Scan(&r.UserID, &r.OfficialName, &r.UserName, &r.Role, &r.MobilePhone, &r.APIID, &r.APIName, &r.APIExposePath, &r.DomainName, &r.DomainID, &r.Domain, &r.Priority, &r.UserDomainID, &r.UserDomainGrantSource, &r.UserDomainDataAccessMode, &r.UserDomainDataScope, &r.UserDomainDataScopeData, &r.UserDomainDataScopeType, &r.IDOnDomain, &r.UserDomainCreateTime, &r.DomainAPIID, &r.DomainAPIGrantSource, &r.DomainAPIDataAccessMode, &r.DomainAPIDataScope, &r.DomainAPIDataScopeData, &r.DomainAPIDataScopeType, &r.DomainAPICreateTime)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select t_v_user_domain_api")
	}
	return &r, nil
}

/*TVXkbSchoolLayout t_v_xkb_school_layout represents kuser.t_v_xkb_school_layout */
type TVXkbSchoolLayout struct {
	School     null.String    `json:"School,omitempty" db:"school,false,character varying"` /* school school */
	Schoolid   null.Int       `json:"Schoolid,omitempty" db:"schoolid,false,integer"`       /* schoolid schoolid */
	SchoolAddi types.JSONText `json:"SchoolAddi,omitempty" db:"school_addi,false,jsonb"`    /* school_addi school_addi */
	Grade      null.String    `json:"Grade,omitempty" db:"grade,false,character varying"`   /* grade grade */
	Gradeid    null.Int       `json:"Gradeid,omitempty" db:"gradeid,false,integer"`         /* gradeid gradeid */
	GradeAddi  types.JSONText `json:"GradeAddi,omitempty" db:"grade_addi,false,jsonb"`      /* grade_addi grade_addi */
	Class      null.String    `json:"Class,omitempty" db:"class,false,character varying"`   /* class class */
	Classid    null.Int       `json:"Classid,omitempty" db:"classid,false,integer"`         /* classid classid */
	ClassAddi  types.JSONText `json:"ClassAddi,omitempty" db:"class_addi,false,jsonb"`      /* class_addi class_addi */
	Filter                    // build DML where clause
}

//TVXkbSchoolLayoutFields full field list for default query
var TVXkbSchoolLayoutFields = []string{
	"School",
	"Schoolid",
	"SchoolAddi",
	"Grade",
	"Gradeid",
	"GradeAddi",
	"Class",
	"Classid",
	"ClassAddi",
}

//Fields return all fields of struct.
func (r *TVXkbSchoolLayout) Fields() []string {
	return TVXkbSchoolLayoutFields
}

//GetTableName return the associated db table name.
func (r *TVXkbSchoolLayout) GetTableName() string {
	var viewNamePattern = regexp.MustCompile(`(?i)^t_v_[a-z0-9_]+$`)
	tableName := "t_v_xkb_school_layout"
	if viewNamePattern.MatchString(tableName) {
		return tableName[2:]
	}
	return tableName
}

// Create inserts the TVXkbSchoolLayout to the database.
func (r *TVXkbSchoolLayout) Create(db Queryer) error {
	_, err := db.Exec(
		`INSERT INTO t_v_xkb_school_layout (school, schoolid, school_addi, grade, gradeid, grade_addi, class, classid, class_addi) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		&r.School, &r.Schoolid, &r.SchoolAddi, &r.Grade, &r.Gradeid, &r.GradeAddi, &r.Class, &r.Classid, &r.ClassAddi)
	if err != nil {
		return errors.Wrap(err, "failed to insert t_v_xkb_school_layout")
	}
	return nil
}

// GetTVXkbSchoolLayoutByPk select the TVXkbSchoolLayout from the database.
func GetTVXkbSchoolLayoutByPk(db Queryer) (*TVXkbSchoolLayout, error) {
	var r TVXkbSchoolLayout
	err := db.QueryRow(
		`SELECT school, schoolid, school_addi, grade, gradeid, grade_addi, class, classid, class_addi FROM t_v_xkb_school_layout WHERE `,
	).Scan(&r.School, &r.Schoolid, &r.SchoolAddi, &r.Grade, &r.Gradeid, &r.GradeAddi, &r.Class, &r.Classid, &r.ClassAddi)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select t_v_xkb_school_layout")
	}
	return &r, nil
}

/*TVXkbUser t_v_xkb_user represents kuser.t_v_xkb_user */
type TVXkbUser struct {
	ID             null.Int       `json:"ID,omitempty" db:"id,false,integer"`                                     /* id id */
	Account        null.String    `json:"Account,omitempty" db:"account,false,character varying"`                 /* account account */
	OfficialName   null.String    `json:"OfficialName,omitempty" db:"official_name,false,character varying"`      /* official_name official_name */
	IDCardType     null.String    `json:"IDCardType,omitempty" db:"id_card_type,false,character varying"`         /* id_card_type id_card_type */
	IDCardNo       null.String    `json:"IDCardNo,omitempty" db:"id_card_no,false,character varying"`             /* id_card_no id_card_no */
	ExternalID     null.String    `json:"ExternalID,omitempty" db:"external_id,false,character varying"`          /* external_id external_id */
	ExternalIDType null.String    `json:"ExternalIDType,omitempty" db:"external_id_type,false,character varying"` /* external_id_type external_id_type */
	Gender         null.String    `json:"Gender,omitempty" db:"gender,false,character varying"`                   /* gender gender */
	Birthday       null.Int       `json:"Birthday,omitempty" db:"birthday,false,bigint"`                          /* birthday birthday */
	Category       null.String    `json:"Category,omitempty" db:"category,false,character varying"`               /* category category */
	Type           null.String    `json:"Type,omitempty" db:"type,false,character varying"`                       /* type type */
	Province       null.String    `json:"Province,omitempty" db:"province,false,character varying"`               /* province province */
	City           null.String    `json:"City,omitempty" db:"city,false,character varying"`                       /* city city */
	Addr           null.String    `json:"Addr,omitempty" db:"addr,false,character varying"`                       /* addr addr */
	MobilePhone    null.String    `json:"MobilePhone,omitempty" db:"mobile_phone,false,character varying"`        /* mobile_phone mobile_phone */
	Email          null.String    `json:"Email,omitempty" db:"email,false,character varying"`                     /* email email */
	Nickname       null.String    `json:"Nickname,omitempty" db:"nickname,false,character varying"`               /* nickname nickname */
	Avatar         []byte         `json:"Avatar,omitempty" db:"avatar,false,bytea"`                               /* avatar avatar */
	AvatarType     null.String    `json:"AvatarType,omitempty" db:"avatar_type,false,character varying"`          /* avatar_type avatar_type */
	Role           null.Int       `json:"Role,omitempty" db:"role,false,bigint"`                                  /* role role */
	Grp            null.Int       `json:"Grp,omitempty" db:"grp,false,bigint"`                                    /* grp grp */
	Addi           types.JSONText `json:"Addi,omitempty" db:"addi,false,jsonb"`                                   /* addi addi */
	Remark         null.String    `json:"Remark,omitempty" db:"remark,false,character varying"`                   /* remark remark */
	Status         null.String    `json:"Status,omitempty" db:"status,false,character varying"`                   /* status status */
	CreateTime     null.Int       `json:"CreateTime,omitempty" db:"create_time,false,bigint"`                     /* create_time create_time */
	Creator        null.Int       `json:"Creator,omitempty" db:"creator,false,bigint"`                            /* creator creator */
	GradeCount     null.String    `json:"GradeCount,omitempty" db:"grade_count,false,text"`                       /* grade_count grade_count */
	ClassCount     null.String    `json:"ClassCount,omitempty" db:"class_count,false,text"`                       /* class_count class_count */
	SchoolID       null.Int       `json:"SchoolID,omitempty" db:"school_id,false,bigint"`                         /* school_id school_id */
	PurchaseRule   types.JSONText `json:"PurchaseRule,omitempty" db:"purchase_rule,false,jsonb"`                  /* purchase_rule purchase_rule */
	SchoolType     null.String    `json:"SchoolType,omitempty" db:"school_type,false,character varying"`          /* school_type school_type */
	School         null.String    `json:"School,omitempty" db:"school,false,character varying"`                   /* school school */
	Subdistrict    null.String    `json:"Subdistrict,omitempty" db:"subdistrict,false,character varying"`         /* subdistrict subdistrict */
	Faculty        null.String    `json:"Faculty,omitempty" db:"faculty,false,character varying"`                 /* faculty faculty */
	Grade          null.String    `json:"Grade,omitempty" db:"grade,false,character varying"`                     /* grade grade */
	GradeSn        null.String    `json:"GradeSn,omitempty" db:"grade_sn,false,text"`                             /* grade_sn grade_sn */
	Class          null.String    `json:"Class,omitempty" db:"class,false,character varying"`                     /* class class */
	ClassSn        null.String    `json:"ClassSn,omitempty" db:"class_sn,false,text"`                             /* class_sn class_sn */
	UnionID        null.String    `json:"UnionID,omitempty" db:"union_id,false,character varying"`                /* union_id union_id */
	WxOpenID       null.String    `json:"WxOpenID,omitempty" db:"wx_open_id,false,character varying"`             /* wx_open_id wx_open_id */
	MpOpenID       null.String    `json:"MpOpenID,omitempty" db:"mp_open_id,false,character varying"`             /* mp_open_id mp_open_id */
	Filter                        // build DML where clause
}

//TVXkbUserFields full field list for default query
var TVXkbUserFields = []string{
	"ID",
	"Account",
	"OfficialName",
	"IDCardType",
	"IDCardNo",
	"ExternalID",
	"ExternalIDType",
	"Gender",
	"Birthday",
	"Category",
	"Type",
	"Province",
	"City",
	"Addr",
	"MobilePhone",
	"Email",
	"Nickname",
	"Avatar",
	"AvatarType",
	"Role",
	"Grp",
	"Addi",
	"Remark",
	"Status",
	"CreateTime",
	"Creator",
	"GradeCount",
	"ClassCount",
	"SchoolID",
	"PurchaseRule",
	"SchoolType",
	"School",
	"Subdistrict",
	"Faculty",
	"Grade",
	"GradeSn",
	"Class",
	"ClassSn",
	"UnionID",
	"WxOpenID",
	"MpOpenID",
}

//Fields return all fields of struct.
func (r *TVXkbUser) Fields() []string {
	return TVXkbUserFields
}

//GetTableName return the associated db table name.
func (r *TVXkbUser) GetTableName() string {
	var viewNamePattern = regexp.MustCompile(`(?i)^t_v_[a-z0-9_]+$`)
	tableName := "t_v_xkb_user"
	if viewNamePattern.MatchString(tableName) {
		return tableName[2:]
	}
	return tableName
}

// Create inserts the TVXkbUser to the database.
func (r *TVXkbUser) Create(db Queryer) error {
	_, err := db.Exec(
		`INSERT INTO t_v_xkb_user (id, account, official_name, id_card_type, id_card_no, external_id, external_id_type, gender, birthday, category, type, province, city, addr, mobile_phone, email, nickname, avatar, avatar_type, role, grp, addi, remark, status, create_time, creator, grade_count, class_count, school_id, purchase_rule, school_type, school, subdistrict, faculty, grade, grade_sn, class, class_sn, union_id, wx_open_id, mp_open_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36, $37, $38, $39, $40, $41)`,
		&r.ID, &r.Account, &r.OfficialName, &r.IDCardType, &r.IDCardNo, &r.ExternalID, &r.ExternalIDType, &r.Gender, &r.Birthday, &r.Category, &r.Type, &r.Province, &r.City, &r.Addr, &r.MobilePhone, &r.Email, &r.Nickname, &r.Avatar, &r.AvatarType, &r.Role, &r.Grp, &r.Addi, &r.Remark, &r.Status, &r.CreateTime, &r.Creator, &r.GradeCount, &r.ClassCount, &r.SchoolID, &r.PurchaseRule, &r.SchoolType, &r.School, &r.Subdistrict, &r.Faculty, &r.Grade, &r.GradeSn, &r.Class, &r.ClassSn, &r.UnionID, &r.WxOpenID, &r.MpOpenID)
	if err != nil {
		return errors.Wrap(err, "failed to insert t_v_xkb_user")
	}
	return nil
}

// GetTVXkbUserByPk select the TVXkbUser from the database.
func GetTVXkbUserByPk(db Queryer) (*TVXkbUser, error) {
	var r TVXkbUser
	err := db.QueryRow(
		`SELECT id, account, official_name, id_card_type, id_card_no, external_id, external_id_type, gender, birthday, category, type, province, city, addr, mobile_phone, email, nickname, avatar, avatar_type, role, grp, addi, remark, status, create_time, creator, grade_count, class_count, school_id, purchase_rule, school_type, school, subdistrict, faculty, grade, grade_sn, class, class_sn, union_id, wx_open_id, mp_open_id FROM t_v_xkb_user WHERE `,
	).Scan(&r.ID, &r.Account, &r.OfficialName, &r.IDCardType, &r.IDCardNo, &r.ExternalID, &r.ExternalIDType, &r.Gender, &r.Birthday, &r.Category, &r.Type, &r.Province, &r.City, &r.Addr, &r.MobilePhone, &r.Email, &r.Nickname, &r.Avatar, &r.AvatarType, &r.Role, &r.Grp, &r.Addi, &r.Remark, &r.Status, &r.CreateTime, &r.Creator, &r.GradeCount, &r.ClassCount, &r.SchoolID, &r.PurchaseRule, &r.SchoolType, &r.School, &r.Subdistrict, &r.Faculty, &r.Grade, &r.GradeSn, &r.Class, &r.ClassSn, &r.UnionID, &r.WxOpenID, &r.MpOpenID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select t_v_xkb_user")
	}
	return &r, nil
}

/*TWxUser 微信开放接口用户信息 represents kuser.t_wx_user */
type TWxUser struct {
	ID             null.Int       `json:"ID,omitempty" db:"id,true,bigint"`                                      /* id 编号 */
	Subscribe      null.Int       `json:"Subscribe,omitempty" db:"subscribe,false,integer"`                      /* subscribe 是否订阅 */
	SubscribeTime  null.Int       `json:"SubscribeTime,omitempty" db:"subscribe_time,false,integer"`             /* subscribe_time 订阅时间 */
	WxOpenID       null.String    `json:"WxOpenID,omitempty" db:"wx_open_id,false,character varying"`            /* wx_open_id 微信公众号openID */
	MpOpenID       null.String    `json:"MpOpenID,omitempty" db:"mp_open_id,false,character varying"`            /* mp_open_id 微信开放平台openID */
	PayOpenID      null.String    `json:"PayOpenID,omitempty" db:"pay_open_id,false,character varying"`          /* pay_open_id 用于微信支付的关系公众号openID */
	UnionID        null.String    `json:"UnionID,omitempty" db:"union_id,false,character varying"`               /* union_id 联合ID */
	GroupID        null.Int       `json:"GroupID,omitempty" db:"group_id,false,integer"`                         /* group_id 组编码 */
	OpenID         null.String    `json:"OpenID,omitempty" db:"open_id,false,character varying"`                 /* open_id openID */
	TagIDList      null.String    `json:"TagIDList,omitempty" db:"tag_id_list,false,character varying"`          /* tag_id_list 标签编码组 */
	Nickname       null.String    `json:"Nickname,omitempty" db:"nickname,false,character varying"`              /* nickname 昵称 */
	Sex            null.Int       `json:"Sex,omitempty" db:"sex,false,integer"`                                  /* sex 性别 */
	Language       null.String    `json:"Language,omitempty" db:"language,false,character varying"`              /* language 语言 */
	City           null.String    `json:"City,omitempty" db:"city,false,character varying"`                      /* city 城市 */
	Province       null.String    `json:"Province,omitempty" db:"province,false,character varying"`              /* province 省份 */
	Country        null.String    `json:"Country,omitempty" db:"country,false,character varying"`                /* country 国家 */
	HeadImgURL     null.String    `json:"HeadImgURL,omitempty" db:"head_img_url,false,character varying"`        /* head_img_url 头像 */
	Privilege      null.String    `json:"Privilege,omitempty" db:"privilege,false,character varying"`            /* privilege 权限 */
	QrScene        null.Int       `json:"QrScene,omitempty" db:"qr_scene,false,integer"`                         /* qr_scene 二维码 */
	SubscribeScene null.String    `json:"SubscribeScene,omitempty" db:"subscribe_scene,false,character varying"` /* subscribe_scene 订阅场景 */
	QrSceneStr     null.String    `json:"QrSceneStr,omitempty" db:"qr_scene_str,false,character varying"`        /* qr_scene_str 一维码 */
	ErrCode        null.Int       `json:"ErrCode,omitempty" db:"err_code,false,integer"`                         /* err_code 错误编码 */
	ErrMsg         null.String    `json:"ErrMsg,omitempty" db:"err_msg,false,character varying"`                 /* err_msg 错误信息 */
	Creator        null.Int       `json:"Creator,omitempty" db:"creator,false,bigint"`                           /* creator 本数据创建者 */
	CreateTime     null.Int       `json:"CreateTime,omitempty" db:"create_time,false,bigint"`                    /* create_time 生成时间 */
	Regenerator    null.Int       `json:"Regenerator,omitempty" db:"regenerator,false,bigint"`                   /* regenerator 更新者 */
	UpdateTime     null.Int       `json:"UpdateTime,omitempty" db:"update_time,false,bigint"`                    /* update_time 帐号信息更新时间 */
	DomainID       null.Int       `json:"DomainID,omitempty" db:"domain_id,false,bigint"`                        /* domain_id 数据隶属 */
	Remark         null.String    `json:"Remark,omitempty" db:"remark,false,character varying"`                  /* remark 备注 */
	Addi           types.JSONText `json:"Addi,omitempty" db:"addi,false,jsonb"`                                  /* addi 附加信息 */
	Status         null.String    `json:"Status,omitempty" db:"status,false,character varying"`                  /* status 状态,00: 有效, 02: 禁止登录, 04: 锁定, 06: 攻击者, 08: 过期 */
	Filter                        // build DML where clause
}

//TWxUserFields full field list for default query
var TWxUserFields = []string{
	"ID",
	"Subscribe",
	"SubscribeTime",
	"WxOpenID",
	"MpOpenID",
	"PayOpenID",
	"UnionID",
	"GroupID",
	"OpenID",
	"TagIDList",
	"Nickname",
	"Sex",
	"Language",
	"City",
	"Province",
	"Country",
	"HeadImgURL",
	"Privilege",
	"QrScene",
	"SubscribeScene",
	"QrSceneStr",
	"ErrCode",
	"ErrMsg",
	"Creator",
	"CreateTime",
	"Regenerator",
	"UpdateTime",
	"DomainID",
	"Remark",
	"Addi",
	"Status",
}

//Fields return all fields of struct.
func (r *TWxUser) Fields() []string {
	return TWxUserFields
}

//GetTableName return the associated db table name.
func (r *TWxUser) GetTableName() string {
	var viewNamePattern = regexp.MustCompile(`(?i)^t_v_[a-z0-9_]+$`)
	tableName := "t_wx_user"
	if viewNamePattern.MatchString(tableName) {
		return tableName[2:]
	}
	return tableName
}

// Create inserts the TWxUser to the database.
func (r *TWxUser) Create(db Queryer) error {
	_, err := db.Exec(
		`INSERT INTO t_wx_user (id, subscribe, subscribe_time, wx_open_id, mp_open_id, pay_open_id, union_id, group_id, open_id, tag_id_list, nickname, sex, language, city, province, country, head_img_url, privilege, qr_scene, subscribe_scene, qr_scene_str, err_code, err_msg, creator, create_time, regenerator, update_time, domain_id, remark, addi, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31)`,
		&r.ID, &r.Subscribe, &r.SubscribeTime, &r.WxOpenID, &r.MpOpenID, &r.PayOpenID, &r.UnionID, &r.GroupID, &r.OpenID, &r.TagIDList, &r.Nickname, &r.Sex, &r.Language, &r.City, &r.Province, &r.Country, &r.HeadImgURL, &r.Privilege, &r.QrScene, &r.SubscribeScene, &r.QrSceneStr, &r.ErrCode, &r.ErrMsg, &r.Creator, &r.CreateTime, &r.Regenerator, &r.UpdateTime, &r.DomainID, &r.Remark, &r.Addi, &r.Status)
	if err != nil {
		return errors.Wrap(err, "failed to insert t_wx_user")
	}
	return nil
}

// GetTWxUserByPk select the TWxUser from the database.
func GetTWxUserByPk(db Queryer, pk0 null.Int) (*TWxUser, error) {
	var r TWxUser
	err := db.QueryRow(
		`SELECT id, subscribe, subscribe_time, wx_open_id, mp_open_id, pay_open_id, union_id, group_id, open_id, tag_id_list, nickname, sex, language, city, province, country, head_img_url, privilege, qr_scene, subscribe_scene, qr_scene_str, err_code, err_msg, creator, create_time, regenerator, update_time, domain_id, remark, addi, status FROM t_wx_user WHERE id = $1`,
		pk0).Scan(&r.ID, &r.Subscribe, &r.SubscribeTime, &r.WxOpenID, &r.MpOpenID, &r.PayOpenID, &r.UnionID, &r.GroupID, &r.OpenID, &r.TagIDList, &r.Nickname, &r.Sex, &r.Language, &r.City, &r.Province, &r.Country, &r.HeadImgURL, &r.Privilege, &r.QrScene, &r.SubscribeScene, &r.QrSceneStr, &r.ErrCode, &r.ErrMsg, &r.Creator, &r.CreateTime, &r.Regenerator, &r.UpdateTime, &r.DomainID, &r.Remark, &r.Addi, &r.Status)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select t_wx_user")
	}
	return &r, nil
}

/*TXkbUser 校快保补充用户信息 represents kuser.t_xkb_user */
type TXkbUser struct {
	ID          null.Int       `json:"ID,omitempty" db:"id,true,bigint"`                               /* id 编号 */
	SchoolID    null.Int       `json:"SchoolID,omitempty" db:"school_id,false,bigint"`                 /* school_id 学校 */
	Subdistrict null.String    `json:"Subdistrict,omitempty" db:"subdistrict,false,character varying"` /* subdistrict 校区 */
	Faculty     null.String    `json:"Faculty,omitempty" db:"faculty,false,character varying"`         /* faculty 学院 */
	Grade       null.String    `json:"Grade,omitempty" db:"grade,false,character varying"`             /* grade 年级 */
	Class       null.String    `json:"Class,omitempty" db:"class,false,character varying"`             /* class 班级 */
	DomainID    null.Int       `json:"DomainID,omitempty" db:"domain_id,false,bigint"`                 /* domain_id 数据属主 */
	Creator     null.Int       `json:"Creator,omitempty" db:"creator,false,bigint"`                    /* creator 创建者用户ID */
	CreateTime  null.Int       `json:"CreateTime,omitempty" db:"create_time,false,bigint"`             /* create_time 创建时间 */
	Regenerator null.Int       `json:"Regenerator,omitempty" db:"regenerator,false,bigint"`            /* regenerator 更新者 */
	UpdateTime  null.Int       `json:"UpdateTime,omitempty" db:"update_time,false,bigint"`             /* update_time 修改时间 */
	Addi        types.JSONText `json:"Addi,omitempty" db:"addi,false,jsonb"`                           /* addi 附加信息 */
	Remark      null.String    `json:"Remark,omitempty" db:"remark,false,character varying"`           /* remark 备注 */
	Status      null.String    `json:"Status,omitempty" db:"status,false,character varying"`           /* status 状态，00：草稿，01：有效，02：作废 */
	Filter                     // build DML where clause
}

//TXkbUserFields full field list for default query
var TXkbUserFields = []string{
	"ID",
	"SchoolID",
	"Subdistrict",
	"Faculty",
	"Grade",
	"Class",
	"DomainID",
	"Creator",
	"CreateTime",
	"Regenerator",
	"UpdateTime",
	"Addi",
	"Remark",
	"Status",
}

//Fields return all fields of struct.
func (r *TXkbUser) Fields() []string {
	return TXkbUserFields
}

//GetTableName return the associated db table name.
func (r *TXkbUser) GetTableName() string {
	var viewNamePattern = regexp.MustCompile(`(?i)^t_v_[a-z0-9_]+$`)
	tableName := "t_xkb_user"
	if viewNamePattern.MatchString(tableName) {
		return tableName[2:]
	}
	return tableName
}

// Create inserts the TXkbUser to the database.
func (r *TXkbUser) Create(db Queryer) error {
	_, err := db.Exec(
		`INSERT INTO t_xkb_user (id, school_id, subdistrict, faculty, grade, class, domain_id, creator, create_time, regenerator, update_time, addi, remark, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`,
		&r.ID, &r.SchoolID, &r.Subdistrict, &r.Faculty, &r.Grade, &r.Class, &r.DomainID, &r.Creator, &r.CreateTime, &r.Regenerator, &r.UpdateTime, &r.Addi, &r.Remark, &r.Status)
	if err != nil {
		return errors.Wrap(err, "failed to insert t_xkb_user")
	}
	return nil
}

// GetTXkbUserByPk select the TXkbUser from the database.
func GetTXkbUserByPk(db Queryer, pk0 null.Int) (*TXkbUser, error) {
	var r TXkbUser
	err := db.QueryRow(
		`SELECT id, school_id, subdistrict, faculty, grade, class, domain_id, creator, create_time, regenerator, update_time, addi, remark, status FROM t_xkb_user WHERE id = $1`,
		pk0).Scan(&r.ID, &r.SchoolID, &r.Subdistrict, &r.Faculty, &r.Grade, &r.Class, &r.DomainID, &r.Creator, &r.CreateTime, &r.Regenerator, &r.UpdateTime, &r.Addi, &r.Remark, &r.Status)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select t_xkb_user")
	}
	return &r, nil
}

// Queryer database/sql compatible query interface
type Queryer interface {
	Exec(string, ...interface{}) (sql.Result, error)
	Query(string, ...interface{}) (*sql.Rows, error)
	QueryRow(string, ...interface{}) *sql.Row
}
