# api comment sample

 please document api by reference it 

## 1 包的说明存放的位置

包(package)如果包含有多个文件或该包的API代码比较多, 则单独建立一个与包同名的文件，作为包说明文档。 否则，package说明放在与package同名的文件内。


## 2 pkgsite 使用条件

### *** 包的名称必须是```x.y```这种格式

### 2.1 不要使用 go mod vendor 或添加 -list=false

go.mod 所在目录及所有递归子目录中，不能有vendor目录, 或在命令行中添加```-list=false```参数 

### 2.2 必须使用git 

go.mod 所在目录或上级目录必须是由git管理的目录，即使用了git作为版本管理的目录。

## 3 pkgsite使用样例

### 3.1 命令模式
```bash
# 如果没有使用 go mod vendor, 以下命令会显示当前目录下的注释文档
pkgsite -http=:6666

# 如果使用了 go mod vendor 则使用以下格式
pkgsite -http=:6666 -list=false
```

#### 3.1.1 关于-list参数
如果go.mod的内容如下：

```
module w2w.io

go 1.18

require (
  github.com/99designs/gqlgen v0.17.5
  github.com/asdine/storm/v3 v3.2.1
  github.com/clbanning/mxj v1.8.4
)	
```


无-list=false则显示所有go.mod中的所有package, 即
```
w2w.io
github.com/99designs/gqlgen v0.17.5
github.com/asdine/storm/v3 v3.2.1
github.com/clbanning/mxj v1.8.4  
```
有，则仅显示
```
w2w.io
```
### 3.2 如果package为w2w.io, 则浏览器地址为
``` http://locahost:6666/w2w.io ```
