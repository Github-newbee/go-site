# go项目脚手架，以go-nunu项目为模板

## gitee 链接

`https://gitee.com/go-nunu/nunu-layout-advanced`

## 目录结构

```
├── api
│   └── v1
├── cmd
│   ├── migration
│   ├── server
│   │   ├── wire
│   │   │   ├── wire.go
│   │   │   └── wire_gen.go
│   │   └── main.go
│   └── task
├── config
├── deploy
├── docs
├── internal
│   ├── handler
│   ├── middleware
│   ├── model
│   ├── repository
│   ├── server
│   └── service
├── pkg
├── scripts
├── test
│   ├── mocks
│   └── server
├── web
├── Makefile
├── go.mod
└── go.sum
```

* cmd：该模块包含了应用的入口点，根据不同的命令进行不同的操作，例如启动服务器、执行数据库迁移等。每个子模块都有一个main.go文件作为入口文件，以及wire.go和wire_gen.go文件用于依赖注入。

* config：该模块包含了应用的配置文件，根据不同的环境（如开发环境和生产环境）提供不同的配置。

* deploy：该模块用于部署应用，包含了一些部署脚本和配置文件

* internal：该模块是应用的核心模块，包含了各种业务逻辑的实现。

  * `handler`：该子模块包含了处理HTTP请求的处理器，负责接收请求并调用相应的服务进行处理。

  * `job`：该子模块包含了后台任务的逻辑实现。

  * `model`：该子模块包含了数据模型的定义。

  * `repository`：该子模块包含了数据访问层的实现，负责与数据库进行交互。

  * `server`：该子模块包含了HTTP服务器的实现。

  * `service`：该子模块包含了业务逻辑的实现，负责处理具体的业务操作。
* `pkg`：该模块包含了一些通用的功能和工具。

* `scripts`：该模块包含了一些脚本文件，用于项目的构建、测试和部署等操作。

* `storage`：该模块用于存储文件或其他静态资源。

* `test`：该模块包含了各个模块的单元测试，按照模块划分子目录。

* `web`：该模块包含了前端相关的文件，如HTML、CSS和JavaScript等。

此外，还包含了一些其他的文件和目录，如授权文件、构建文件、README等。

## 要求

项目的基本环境

* Golang 1.19或以上版本
* Git
* Docker （可选）
* MySQL5.7或更高版本（可选）
* Redis（可选）

## 安装

国内用户使用GOPROXY加速go install

```
go env -w GO111MODULE = on
go env -w GOPROXY= https://goproxy.cn,direct
```

**一些go命令**

```
go run xxx.go
```

**主要功能**
用于直接编译并运行 Go 程序，它不会生成可执行文件，适合快速验证代码逻辑。

```
go mod tidy （拉取下拉的项目使用该命令安装依赖）
```

**主要功能**

* 移除`go.mod`文件中未使用的依赖
* 添加`go.mod`文件中缺失的依赖项
* 更新`go.sum`文件以匹配`go.mod`文件中的依赖项

```
go build 
```

**主要功能**

* 用于编译包或可执行程序，但不会将编译后的文件安装到特点目录，而是将可执行文件生成在当前目录下
* 如果只是想在本地测试代码编译情况，使用go build 更合适

```
go install [packages]
````

**主要功能**

* 安装可执行文件
* 安装第三方工具
* 安装本地开发的工具

## 依赖说明

### wire

wire 是 Go 语言的一个依赖注入工具，由 Google 开发并开源。依赖注入（Dependency Injection，简称 DI）是一种设计模式，它允许对象在创建时接收其依赖项，而不是自己创建这些依赖项，这样可以提高代码的可测试性、可维护性和可扩展性。

wire 的主要作用就是`自动生成依赖注入代码`，避免手动编写大量样板代码。在大型项目中，随着依赖关系变得复杂，手动管理和创建对象之间的依赖关系会变得非常繁琐且容易出错，wire 可以根据开发者定义的规则自动生成初始化代码，让依赖注入的过程更加高效和可靠。

执行生成的命令，需要在wire.go文件目录下执行
命令执行后会在目录下生成wire_gen.go文件

```
wire
```

### zap

是一个高性能的结构化日志库，由 Uber 开发。它提供了高效的日志记录功能，支持结构化日志和非结构化日志，并且在性能上进行了优化。

**主要特点**

* 高性能： `zap` 通过减少内存分配和垃圾回收，提高了日志记录的性能
* 结构化：支持结构化日志记录，可以方便记录键值对形式的日志信息
* 灵活性：提供了多种日志模式，包括同步和异步日志记录
* 可扩展性：支持自定义日志记录器和日志格式


### 面向对象编程

业务代码中，使用结构体和方法来实现面向对象编程

```
type UserHandler struct {
	userService service.UserService
}
```

* 这种结构设计遵循依赖注入模式
* 通过结构体字段来存储依赖服务和组件
* 通过NewUserHandler函数来创建UserHandler实例，并传递必要的依赖
* 在UserHandler的方法中，可以直接使用依赖服务和组件
* 这种设计使得UserHandler的实现更加灵活和可测试，同时也符合面向对象编程的原则

**方法接收器**

Go 通过在结构体上定义方法来组织代码

```
func (h *UserHandler) UpdateProfile(ctx *gin.Context) {
	userId := GetUserIdFromCtx(ctx)
}
```

这里的`h`是方法接收器，它允许方法访问和修改结构体的字段

*UserHandler 表示这个方法属于 UserHandler 结构体

**好处** 

* 模块化：将数据和行为封装在结构体中，使得代码更加模块化
* 封装：通过结构体和方法，将实现细节隐藏在结构体内部，外部只需关注接口和方法
* 可维护性：通过结构体和方法，将实现细节隐藏在结构体内部，外部只需关注接口和方法
* 依赖清晰：通过结构体字段明确声明了组件之间的依赖关系