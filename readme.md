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

## 架构说明

### 分层架构

在 Go Web 项目中，常见的分层架构包含表现层（Presentation Layer）、服务层（Service Layer）、数据访问层（Repository Layer）和数据存储层（如数据库）。各层的主要职责如下：

* 表现层：负责处理 HTTP 请求和响应，通常是控制器（Controller）或处理函数（Handler）的工作。
* 服务层：实现具体的业务逻辑，调用数据访问层的方法进行数据操作。
* 数据访问层：负责与数据存储层交互，如数据库查询、插入、更新和删除操作。
* 数据存储层：实际存储数据的地方，如 MySQL、Redis 等。

#### 依赖注入连接关系

1. 定义接口

```
type CategoryService interface {
    CreateCategory(ctx context.Context, category *v1.CategoryRequest) (*model.Category, error)
    GetAllCategory(req v1.GetCategoryRequest, ctx context.Context) ([]model.Category, error)
    UpdateCategory(ctx context.Context, id string, req *v1.CategoryRequest) error
}
```

这里定义了`CategoryService`接口，它规定了服务层应该具备的方法，表现层可以依赖这个接口来调用服务层的功能

2. 数据访问层

```
type CategoryRepository interface {
 CreateCategory(ctx context.Context, category *model.Category) (*model.Category, error)
 GetCategoryById(ctx context.Context, id string) (*model.Category, error)
 GetAllCategory(req v1.GetCategoryRequest, ctx context.Context) ([]model.Category, error)
 UpdateCategory(ctx context.Context, category *model.Category) error
}

```

数据访问层也定义了一个接口，具体的实现会在数据访问层的具体结构体中完成

3. 服务层实现

```
type categoryService struct {
    *Service
    categoryRepository repository.CategoryRepository
}

func NewCategoryService(
    service *Service,
    categoryRepository repository.CategoryRepository,
) CategoryService {
    return &categoryService{
        Service:            service,
        categoryRepository: categoryRepository,
    }
}

func (s *categoryService) CreateCategory(ctx context.Context, category *v1.CategoryRequest) (*model.Category, error) {
    // 将请求数据转换为模型数据
    modelCategory := &model.Category{
        // 赋值操作
    }
    // 调用数据访问层的方法
    return s.categoryRepository.CreateCategory(ctx, modelCategory)
}

func (s *categoryService) GetAllCategory(req v1.GetCategoryRequest, ctx context.Context) ([]model.Category, error) {
    return s.categoryRepository.GetAllCategory(ctx)
}

func (s *categoryService) UpdateCategory(ctx context.Context, id string, req *v1.CategoryRequest) error {
    // 将请求数据转换为模型数据
    modelCategory := &model.Category{
        // 赋值操作
    }
    return s.categoryRepository.UpdateCategory(ctx, id, modelCategory)
}
```

* `categoryService`结构体包含了一个`Service`指针和一个`categoryRepository`接口类型的字段。这体现了依赖注入的思想，`categoryService` 依赖于 `categoryRepository` 接口，而不是具体的实现。
* `NewCategoryService` 函数是一个工厂函数，用于创建 `categoryService` 实例。通过参数注入 `categoryRepository`，使得 `categoryService` 可以使用不同的 `CategoryRepository` 实现。
* 在 `categoryService` 的方法实现中，调用了 `categoryRepository` 的方法来完成具体的数据操作，实现了业务逻辑和数据访问的分离。

4. 表现层调用

```
type CategoryHandler struct {
 *Handler
 categoryService service.CategoryService
}

func NewCategoryHandler(
 handler *Handler,
 categoryService service.CategoryService,
) *CategoryHandler {
 return &CategoryHandler{
  Handler:         handler,
  categoryService: categoryService,
 }
}

func (h *CategoryHandler) GetAllCategory(ctx *gin.Context) {
 req := v1.GetCategoryRequest{}
 request.Assign(ctx, &req)
 res, err := h.categoryService.GetAllCategory(req, ctx)
 if err != nil {
  h.logger.WithContext(ctx).Error("categoryService.GetAllCategory", zap.Error(err))
  v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
  return
 }

 v1.HandleSuccess(ctx, res)
}

```

* `CategoryHandler` 结构体依赖于 `service.CategoryService` 接口，通过 `NewCategoryHandler` 函数注入 `categoryService` 实例。
* 在 `CreateCategory` 处理函数中，调用 `categoryService` 的 `CreateCategory` 方法来处理创建分类的请求。

#### 依赖注入的优势

* **解耦：** 各层之间通过接口进行依赖，降低了代码的耦合度。例如，当需要更换数据存储方式时，只需要实现一个新的 CategoryRepository 接口，并注入到 categoryService 中，而不需要修改 categoryService 的业务逻辑。

* **可测试性：** 可以方便地在单元测试中模拟接口的实现，独立测试各层的功能。例如，在测试 `categoryService` 时，可以使用模拟的 `CategoryRepository` 来避免依赖实际的数据库。

## 代码说明

### 方法的返回值

1. `GetByUsername`方法：

* 返回的是指针，因为是一个单一结构体对象，是值类型，返回指针可以避免拷贝整个结构体，提高性能

2. `GetUserAll`方法：

* 返回`[]model.User`切片，因为这是一个对象列表，返回切片本身已经是引用类型，不需要返回指针

```
func (r *userRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
 var user model.User
 if err := r.DB(ctx).Where("username = ?", username).First(&user).Error; err != nil {
  if errors.Is(err, gorm.ErrRecordNotFound) {
   return nil, nil
  }
  return nil, err
 }
 return &user, nil
}

func (r *userRepository) GetUserAll(req v1.GetAllUsersRequest, ctx context.Context) (results []model.User, err error) {
 // 执行查询
 qs := r.DB(ctx).Model(&model.User{}).Session(&gorm.Session{}).Scopes(db.FilterByQuery(req))
 if err := qs.Limit(req.Limit).Offset(req.Skip).Find(&results).Error; err != nil {
  return nil, err
 }
 return results, nil
}
```

### grom中的部分解释

* `r.DB(ctx)`：获取数据库链接
* `Model(&model.User{})`: 指定查询的模型为`model.User`
* `Session(&gorm.Session{}): 创建一个新的GORM会话，确保当前查询在一个独立的会话中进行，对于并发操作特别重要，避免不同查询之间的相互干扰
* `Scopes(db.FilterByQuery(req))`：应用查询范围
* `Limit(req.Limit)`: 限制查询结果的数量
* `Offset(req.Skip)`: 设置查询结果的偏移量，用于分页
* `Find(&results)`: 执行查询并将结果存储在results切片中
* `Error`: 检查查询是否发生错误，如果有错误，则返回nil和错误对象

```
func (r *userRepository) GetUserAll(req v1.GetAllUsersRequest, ctx context.Context) (results []model.User, err error) {
 // 执行查询
 qs := r.DB(ctx).Model(&model.User{}).Session(&gorm.Session{}).Scopes(db.FilterByQuery(req))
 if err := qs.Limit(req.Limit).Offset(req.Skip).Find(&results).Error; err != nil {
  return nil, err
 }
 return results, nil
}
```
