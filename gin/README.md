# go-web-boilerplate/gin

Golang boilerplate web using gin server engine.

This boilerplate also using standard dependency injection and singleton for config, database, and redis.

## Run

To run, you can use this command:

```sh
# ask go to run main.go file with args: 127.0.0.1:2379 (connect to etcd service), /gin/dev (get config with this key from etcd service)
go run main.go 127.0.0.1:2379 /gin/dev
```

## Example

Since this are boilerplate, please rename the app or module name from `boilerplate/` to your application or module name.

### Config

For config, this boilerplate will use `etcd` to get initial config.

### Constant

Example of route `constant.go`.

```go
package constant

import "time"

const (
	DEFAULT = "default" // redis connection name: default
	SETTING = "setting" // redis connection name: setting

	MEMBER_TOKEN_KEY = "member:token:%d"
	MEMBER_TOKEN_TTL = 24 * time.Hour // time to live for redis key

	SETTING_LIST = "setting:list"
	SETTING_MAP  = "setting:map"
	SETTING_TTL  = 0 // time to live for redis key, set 0 for key without expiry time
)

...
```

### Route

Example of route `base.go`.

```go
package route

import (
	"boilerplate/server/controller"
	"boilerplate/server/middleware"

	"github.com/gin-gonic/gin"
)

type Route struct {
	...

	userController    *controller.UserController
  // add controller
}

func NewRoute(gin *gin.Engine, cfg *config.Config, db *database.Database, resp *response.Response, mdl *middleware.Middleware) *Route {
	ctrl := controller.NewController(cfg, db, mdl.Pagination, resp)

	route := &Route{
		...

		userController:    controller.NewUserController(ctrl),
		// add controller
	}

	route.addRoutes()

	return route
}

func (r *Route) addRoutes() {
  r.defaultRoutes()

	...

	r.authRoutes()
  // add routes
}

// add default route for get, post, health check, no route (404), and no method (405)
func (r *Route) defaultRoutes() {
	r.gin.GET("", r.controller.Default)
	r.gin.POST("", r.controller.Default)
	r.gin.GET("/health", r.controller.HealthCheck)
	r.gin.POST("/health", r.controller.HealthCheck)
	r.gin.NoRoute(r.controller.RouteNotFound)
	r.gin.NoMethod(r.controller.MethodNotAllowed)
}
```

Example of route `user.go`.

```go
package route

func (r *Route) authRoutes() {
  // set middleware in route group (possible to add multiple middleware)
	userRoutes := r.apiGroup.Group("/auth", r.middleware.Language())
	{
		...

    // set middleware in route (possible to add multiple middleware)
		userRoutes.GET("/user", r.middleware.JwtAuth(), r.userController.GetUserLogin)
    // add routes
	}
}

...
```

### Request

Example of request `user.go`.

```go
package request

// go `validator` doesn't support clean whitespace, so this boilerplate also add custom validation `empty_string` to handle it
// gin `binding` tag will help to validate the field
type RegisterInput struct {
	Username          string `json:"username" binding:"required,empty_string,min=3,max=45"`
	Password          string `json:"password" binding:"required,empty_string,min=10,max=255"`
	ConfirmedPassword string `json:"confirmed_password" binding:"required,empty_string,min=10,max=255,eqfield=Password"`
}
```

### Response

Example of request `user.go`.

```go
package response

import "time"

type LoginOutput struct {
	UserInfoOutput
	Token string `json:"token"`
}

type UserInfoOutput struct {
	Id           uint64     `json:"id"`
	Username     string     `json:"username"`
	RegisteredAt time.Time  `json:"registered_at"`
	LoginAt      *time.Time `json:"login_at"`
}
```

### Model

Example of model `user.go`.

```go
package model

import (
	"database/sql"
	"time"
)

type User struct {
	Base
	Id           uint64       `gorm:"type:bigint UNSIGNED AUTO_INCREMENT;primary_key" json:"id"`
	Username     string       `gorm:"type:varchar(45);not null;column:username;index:username_uq,unique" json:"username"`
	Password     string       `gorm:"type:varchar(255);not null;column:password" json:"-"`
	RegisteredAt time.Time    `gorm:"type:time;not null;column:registered_at" json:"registered_at"`
	LoginAt      sql.NullTime `gorm:"type:time;null;column:login_at" json:"login_at"`
}

// set model table name
// ** gorm are using users as table name, so if your user table are not users, you need to change it in here
func (u *User) TableName() string {
	return "user"
}
```

### Controller

Example of controller `user.go`

```go
package controller

import (
	"boilerplate/server/middleware"
	"boilerplate/server/request"
	"boilerplate/server/response"
	"boilerplate/service"
	cerror "boilerplate/utils/error"
	"boilerplate/utils/logger"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type UserController struct {
	*Controller
}

func NewUserController(c *Controller) *UserController {
	return &UserController{c}
}

func (uc *UserController) Register(c *gin.Context) {
	var registerInput request.RegisterInput

  // bind json to register request
	if err := c.ShouldBindBodyWith(&registerInput, binding.JSON); err != nil {
		uc.response.Json(c, http.StatusBadRequest, nil, cerror.Fail(cerror.FuncName(), "invalid_input_request", map[string]any{
			"validator": err,
			"input":     registerInput,
		}, err))
		return
	}

  // call service
	user, token, err := uc.userService.Register(registerInput.Username, registerInput.Password)

	if err != nil {
    // return error
		uc.response.Json(c, http.StatusBadRequest, nil, err)
		return
	}

  // return success
	uc.response.Json(c, http.StatusOK, response.LoginOutput{
		UserInfoOutput: response.UserInfoOutput{
			Id:           user.Id,
			Username:     user.Username,
			RegisteredAt: user.RegisteredAt,
			LoginAt:      &user.LoginAt.Time,
		},
		Token: token,
	}, nil)
}
```

### Service

Example of service `user.go`

```go
package service

import (
	"boilerplate/model"
	"boilerplate/repository"
	"boilerplate/server/middleware"
	cconstant "boilerplate/utils/constant"
	cerror "boilerplate/utils/error"
	"boilerplate/utils/token"
	"database/sql"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	*Service
}

func NewUserService(s *Service) *UserService {
	return &UserService{s}
}

func (us *UserService) Register(username, password string) (*model.User, string, error) {
	hashedPassword, err := us.HashPassword(password)
	if err != nil {
		return nil, "", err
	}

  // call repository
	user, err := us.userRepo.CreateUser(&model.User{
		Username:     username,
		Password:     hashedPassword,
		RegisteredAt: time.Now(),
		LoginAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	})

	if err != nil {
		return nil, "", err
	}

	token, err := us.GenerateToken(user.Id)

	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}
```

### Repository

Example of service `user.go`.

```go
package repository

import (
	"boilerplate/model"
	"boilerplate/server/middleware"
	cerror "boilerplate/utils/error"

	"gorm.io/gorm"
)

type UserRepository struct {
	*Repository
}

func NewUserRepository(r *Repository) *UserRepository {
	return &UserRepository{r}
}

// create user and also example to omit some field
func (repo *UserRepository) CreateUser(user *model.User) (*model.User, error) {
	if err := repo.db.MySQL.Orm.Omit("TotalRows").Create(&user).Error; err != nil {
		return nil, cerror.Fail(cerror.FuncName(), "failed_db_insert", nil, err)
	}

	return user, nil
}
```

### Pagination

Example of pagination.

Property `page` and `pageSize` are set from `middleware pagination`.

```go
// controller
func (uc *UserController) GetUsers(c *gin.Context) {
	res, err := uc.userService.GetUserList()

	if err != nil {
		uc.response.Json(c, http.StatusBadRequest, nil, err)
		return
	}

  // set response pagination from pagination
	uc.response.Pagination = uc.Pagination
	uc.response.Json(c, http.StatusOK, res, nil)
}

// repository
func (repo *UserRepository) GetUserList() ([]*model.User, error) {
	var userList []*model.User

	result := repo.db.MySQL.Orm.Select("id, username, password, point, registered_at, login_at")

  // add scopes into query when pagination used
	if repo.pagination.Page > 0 {
		result = result.Scopes(repo.Paginate)
	}

	result = result.Find(&userList)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return userList, nil
		}

		return nil, cerror.Fail(cerror.FuncName(), "failed_db_query", nil, result.Error)
	}

  // when pagination used, set pagination TotalRecord property
	if (len(userList) > 0) && (repo.pagination.Page > 0) {
		totalRecord, err := repo.GetTotalRecord(result, "id")

		if err != nil {
			return nil, err
		}

		repo.pagination.TotalRecord = totalRecord
	}

	return userList, nil
}
```

### Redis

Example of `redis` cache

```go
// generate token and cache to redis
func (us *UserService) GenerateToken(userId uint64) (tkn string, err error) {
	tkn, err = token.GenerateToken(userId, us.config.TokenConfig.TokenHourLifeSpan, []byte(us.config.TokenConfig.SecretKey))

	if err != nil {
		return "", err
	}

	key := fmt.Sprintf(cconstant.MEMBER_TOKEN_KEY, userId)

  // cache generated token to redis within one hour
	if err = us.redis.Set(key, tkn, cconstant.MEMBER_TOKEN_TTL); err != nil {
		return "", err
	}

	return tkn, nil
}

// cache settings data to redis
func (us *SettingService) CachedSettings() ([]*model.Setting, error) {
	var cachedSettings []*model.Setting
	var err error

  // get cached settings from different connection
	settings, err := us.redis.Connection(cconstant.SETTING).Get(cconstant.SETTING_LIST)
	if err != nil {
		logger.Sugar.Error(err)
		return nil, err
	}

	if !utils.IsEmptyString(settings) {
		if err = json.Unmarshal([]byte(settings), &cachedSettings); err != nil {
			return nil, cerror.Fail(cerror.FuncName(), "failed_unmarshal", nil, err)
		}
	}

	if len(cachedSettings) == 0 {
		cachedSettings, err = us.settingRepo.GetSettingList()

		if err != nil {
			return nil, err
		}

		jsonData, err := json.Marshal(cachedSettings)

		if err != nil {
			return nil, cerror.Fail(cerror.FuncName(), "failed_marshal", nil, err)
		}

    // cache settings to different connection
		if err = us.redis.Connection(cconstant.SETTING).Set(cconstant.SETTING_LIST, string(jsonData), cconstant.SETTING_TTL); err != nil {
			return nil, err
		}
	}

	return cachedSettings, nil
}

// remove redis key when setting data updated (rows > 0, means database really update setting data, if row == 0, means there's no such data or value of setting data remain same), possible to remove multiple key
func (us *SettingService) UpdateSettingByName(name string, value string) (int64, error) {
	rows, err := us.settingRepo.UpdateSettingByName(name, value)

	if err != nil {
		return rows, err
	}

	if rows > 0 {
		if err := us.db.Redis.Connection(cconstant.SETTING).Del([]string{cconstant.SETTING_LIST, cconstant.SETTING_MAP}...); err != nil {
			return rows, err
		}
	}

	return rows, err
}
```

### Locales

locales folder to provide i18n key value.

### Database

#### MySQL

Using `gorm` package to handle database connection and pooling.

`gorm` also support multiple connection and resolver.

Example

```json
  {
    ...

    "default": {
      "network": "tcp",
      "host": "127.0.0.1",
      "port": 3306,
      "username": "root",
      "password": "",
      "schema": "",
      "charset": "utf8mb4"
    },
    "connections": [
      {
        "datas": ["user"],
        "writes": [
          {
            "network": "tcp",
            "host": "127.0.0.1",
            "port": 3306,
            "username": "root",
            "password": "",
            "schema": "",
            "charset": "utf8mb4"
          }
        ],
        "reads": [
          {
            "network": "tcp",
            "host": "127.0.0.1",
            "port": 3306,
            "username": "root",
            "password": "",
            "schema": "",
            "charset": "utf8mb4"
          }
        ]
      }
    ]

    ...
  }
```

Property `default` are used to initialized `mysql` connection with `gorm`.

When `connections` property set in config, boilerplate will register it to `gorm` as db resolver, where `datas` field to determine which table are need to use these connections.

Db resolver from `gorm` will automatically determine which query use write or read connection.

## Contributing

See the [contributing guide](CONTRIBUTING.md) to learn how to contribute to the repository and the development workflow.

## License

MIT
