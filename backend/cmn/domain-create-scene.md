# domain 权限相关场景用例

## domain 语法

```
org.dept.section.group^role@userID

    org: 机构，如广州大学，广州大洋教育科技有限公司
    dept: 部门，如保卫处，财务部
    section: 科室，如经费预算科，运维组
    group: 组，如车辆服务组，美工组
    role: 角色，如经理，处长，组长
    userID: 用户在系统中的编号，`    
```
它的实际意义由用户自己定义，在保存数据时保存在业务数据表的domain列中

### 样例数据
广州大学.保卫处.门岗.南门^班长@1022
