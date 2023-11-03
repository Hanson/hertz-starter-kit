#hertz-starter-kit

一个基于 hertz 的脚手架工具，开箱即用，对 gorm，protobuf 等内容做了些封装，方便快速开发。

# 特性支持

* 支持 proto 生成 model
* 通过 make 命令快捷生成 pb
* 对所有 proto model 的 string 做了默认 varchar(100)

# 环境
go install github.com/cloudwego/hertz/cmd/hz@latest
go install github.com/favadi/protoc-go-inject-tag@latest

# 下载运行
```
git clone https://github.com/Hanson/hertz-starter-kit
cd hertz-starter-kit
go run .
```

# 开发
可以全局搜索 hertz-starter-kit， 用你的包名替换掉

## db
项目默认使用 MySQL，如有其它需求需要修改，详情看 db 文件夹

* 所有 Model 都需要有 id, created_at, updated_at, deleted_at, 并为 int64
* proto 编写 Model 开头的结构体，创建的时候可自动移出前缀 model_, 详情看 db/naming.go:35
* 对所有 proto model 的 string 做了默认 varchar(100)， 详情看 db/migrator.go:29
* 对所有 proto text 改为 not null false
* 对 DeletedAt 进行默认 0
* db 使用可以用 db.NewInstance 或者 db.NewModel

## proto

* curd 的开发可以直接复制 idl/README.md,把你的 model 替换掉 demo 即可
* 执行命令 make update p=path 即可更新 pb，path 为 idl 里面的相对路径
