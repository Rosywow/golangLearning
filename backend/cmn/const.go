package cmn

// //go:generate go install golang.org/x/tools/cmd/stringer@latest

//level "0": 无组/角色/数据限制, 可访问全部数据
//level "2": 机构#角色级别, 实现了不同角色授权，但不控制数据范围
//level "4": 机构#角色$ID, 实现了不同角色授权，可控制 creator || all
//level "8": 机构.DEPT#角色$ID, 实现了不同角色授权，可控制 creator || GRPs */
//AccessControlLevel string `json:"access_control_level,omitempty"`

type CAccessControlLevel int64

const (
	//CACLUnlimit allow access any data
	CACLUnlimit CAccessControlLevel = 0

	//CACLScopeByOrgRole allow access organization role scope data
	CACLScopeByOrgRole CAccessControlLevel = 2

	//CACLScopeByOrgRoleID allow access organize
	CACLScopeByOrgRoleID CAccessControlLevel = 4

	CACLScopeByOrgDeptRoleID CAccessControlLevel = 8
)

//go:generate stringer -type=CDomain
type CDomain int64

//常量命名规范
//C+分类+名称, 期中分类与名称使用驼峰命名法
//例如，表示角色(domain)中的系统(sys)角色，使用 CDomainSys作为名称
const (
	//CDomainSys domain constants
	CDomainSys                 CDomain = 322   //sys  系统
	CDomainSysAdmin            CDomain = 333   //sys^admin  系统.管理
	CDomainSysMaintain         CDomain = 366   //sys^maintain  系统.运维
	CDomainSysUser             CDomain = 377   //sys^user  系统.用户
	CDomainSysAnonymous        CDomain = 388   //sys^anonymous  系统.匿名
	CDomainSysPromotion        CDomain = 561   //sys^promotion  系统.运营
	CDomainSysSale             CDomain = 563   //sys^sale  系统.销售
	CDomainSysTrial            CDomain = 566   //sys^trial  系统.测试
	CDomainAdmin               CDomain = 567   //^admin  管理
	CDomainMaintain            CDomain = 569   //^maintain  运维
	CDomainUser                CDomain = 671   //^user  用户
	CDomainAnonymous           CDomain = 673   //^anonymous  匿名
	CDomainPromotion           CDomain = 675   //^promotion  运营
	CDomainSale                CDomain = 677   //^sale  销售
	CDomainTrial               CDomain = 679   //^trial  测试
	CDomainQNear               CDomain = 1077  //qnear  近邻科技
	CDomainQNearAdmin          CDomain = 1079  //qnear^admin  近邻科技.管理
	CDomainAbilityIdx          CDomain = 1177  //abilityIdx  能力索引
	CDomainAbilityIdxAdmin     CDomain = 1179  //abilityIdx^admin  能力索引.管理
	CDomainForeseeLab          CDomain = 1277  //foreseeLab  IT双创精英孵化实训室
	CDomainForeseeLabAdmin     CDomain = 1279  //foreseeLab^admin  IT双创精英孵化实训室.管理
	CDomainRecruitMgr          CDomain = 1377  //recruitMgr  人才引进
	CDomainRecruitMgrAdmin     CDomain = 1379  //recruitMgr^admin  人才引进.管理
	CDomainJXDD                CDomain = 1477  //jxdd  教学督导
	CDomainJXDDAdmin           CDomain = 1479  //jxdd^admin  教学督导.管理
	CDomainDonate              CDomain = 1577  //donate  校友会小额捐献
	CDomainDonateAdmin         CDomain = 1579  //donate^admin  校友会小额捐献.管理
	CDomainXKB                 CDomain = 10002 //xkb  校快保
	CDomainXKBAdmin            CDomain = 10004 //xkb^admin  校快保.管理
	CDomainXKBSale             CDomain = 10006 //xkb^sale  校快保.销售经理
	CDomainXKBSchoolAdmin      CDomain = 10008 //xkb.school^admin  校快保.学校管理员
	CDomainXKBSchoolStatistics CDomain = 10010 //xkb.school^statistics  校快保.学校统计员
	CDomainXKBUser             CDomain = 10012 //xkb^user  校快保.客户
	CDomainXKBPromotion        CDomain = 10016 //xkb^promotion  校快保.运营
	CDomainXKBFE               CDomain = 10020 //xkb^fe  校快保.前台
	CDomainXKBBE               CDomain = 10030 //xkb^be  校快保.后台
)

var roleToName = map[CDomain]string{
	CDomainSysAdmin:            "sys^admin",
	CDomainSysMaintain:         "sys^maintain",
	CDomainSysUser:             "sys^user",
	CDomainSysAnonymous:        "sys^anonymous",
	CDomainSysPromotion:        "sys^promotion",
	CDomainSysSale:             "sys^sale",
	CDomainSysTrial:            "sys^trial",
	CDomainAdmin:               "^admin",
	CDomainMaintain:            "^maintain",
	CDomainUser:                "^user",
	CDomainAnonymous:           "^anonymous",
	CDomainPromotion:           "^promotion",
	CDomainSale:                "^sale",
	CDomainTrial:               "^trial",
	CDomainQNear:               "qnear",
	CDomainQNearAdmin:          "qnear^admin",
	CDomainAbilityIdx:          "abilityIdx",
	CDomainAbilityIdxAdmin:     "abilityIdx^admin",
	CDomainForeseeLab:          "foreseeLab",
	CDomainForeseeLabAdmin:     "foreseeLab^admin",
	CDomainRecruitMgr:          "recruitMgr",
	CDomainRecruitMgrAdmin:     "recruitMgr^admin",
	CDomainJXDD:                "jxdd",
	CDomainJXDDAdmin:           "jxdd^admin",
	CDomainDonate:              "donate",
	CDomainDonateAdmin:         "donate^admin",
	CDomainXKB:                 "xkb",
	CDomainXKBAdmin:            "xkb^admin",
	CDomainXKBSale:             "xkb^sale",
	CDomainXKBSchoolAdmin:      "xkb.school^admin",
	CDomainXKBSchoolStatistics: "xkb.school^statistics",
	CDomainXKBUser:             "xkb^user",
	CDomainXKBPromotion:        "xkb^promotion",
	CDomainXKBFE:               "xkb^fe",
	CDomainXKBBE:               "xkb^be",
}

func RoleName(roleID CDomain) (s string) {
	s, _ = roleToName[roleID]
	if s == "" {
		s = "未知"
	}
	return
}

const (
	CUserDefaultCreator    CDomain = 1000
	CUserDefaultMaintainer CDomain = 1000
	CUserDefaultDomainID   CDomain = CDomainSysAdmin
)

const (
	CSysUserByName  = "tUserByName"  //TUser.Account->指向sysUser(TUser)
	CSysUserByID    = "tUserByID"    //TUser.ID-> 指向sysUser(TUser).account
	CSysUserByEmail = "tUserByEmail" //TUser.Email->指向sysUser(TUser).account
	CSysUserByTel   = "TUserByTel"   //TUser.MobilePhone->指向sysUser(TUser).account

	CWxUserByUnionID = "tWeChatUserByUnionID" //TWxUser.UnionID->指向wxUser(TWxUser)
	CWxUserByID      = "tWeChatUserByID"      //TUser.ID->TWxUser.UnionID

	CWxUserByOpenID = "tWeChatUserByOpenID" //TWxUser.MpOpenID->TUser.ID, TWxUser.WxOpenID->TUser.ID
)

// 用户访问的模块类型
const (
	//CFuncUnDetermined 未知类型
	CFuncUnDetermined = iota

	//CFuncApi 函数
	CFuncApi

	//CFuncNonAdminFileServe 后台管理员模块
	CFuncNonAdminFileServe

	//CFuncAdminFileServe 前台普通用户模块
	CFuncAdminFileServe
)

var reqFnTypeStr = map[int]string{
	CFuncUnDetermined:      "未知",
	CFuncApi:               "函数",
	CFuncNonAdminFileServe: "前台",
	CFuncAdminFileServe:    "后台",
}

func ReqFnTypeString(reqFnType int) (s string) {
	s, _ = reqFnTypeStr[reqFnType]
	if s == "" {
		s = "未知"
	}
	return
}
