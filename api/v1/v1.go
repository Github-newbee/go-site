package v1

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func HandleSuccess(ctx *gin.Context, data interface{}) {
	if data == nil {
		data = map[string]interface{}{}
	}
	resp := Response{Code: errorCodeMap[ErrSuccess], Message: ErrSuccess.Error(), Data: data}
	if _, ok := errorCodeMap[ErrSuccess]; !ok {
		resp = Response{Code: 0, Message: "", Data: data}
	}
	ctx.JSON(http.StatusOK, resp)
}

func HandleError(ctx *gin.Context, httpCode int, err error, data interface{}) {
	if data == nil {
		data = map[string]string{}
	}
	var resp Response
	if _, ok := errorCodeMap[err]; ok {
		resp = Response{Code: errorCodeMap[err], Message: err.Error(), Data: data}
	} else {
		switch e := err.(type) {
		case *mysql.MySQLError:
			errCode, errMsg := sqlErr(e)
			resp = Response{Code: errCode, Message: errMsg, Data: data}
		default:
			resp = Response{Code: 500, Message: "unknown error", Data: data}
		}
	}
	fmt.Println("-----详细错误", resp)
	fmt.Printf("错误的类型%T \n", err)
	// if _, ok := errorCodeMap[err]; !ok {
	// 	resp = Response{Code: 500, Message: "unknown error", Data: data}
	// }
	ctx.JSON(httpCode, resp)
}

func sqlErr(err *mysql.MySQLError) (int, string) {
	// 根据错误码返回对应的错误码和消息
	switch err.Number {
	case 1062:
		return int(err.Number), fmt.Sprintf("数据已存在: %s", err.Message)
	case 1406:
		return int(err.Number), fmt.Sprintf("参数过长: %s", strings.Split(err.Message, "'")[1])
	case 1452:
		return int(err.Number), fmt.Sprintf("外键约束错误: %s", err.Message)
	default:
		return 500, "DB error"
	}
}

type Error struct {
	Code    int
	Message string
}

var errorCodeMap = map[error]int{}

func newError(code int, msg string) error {
	err := errors.New(msg)
	errorCodeMap[err] = code
	return err
}
func (e Error) Error() string {
	return e.Message
}
