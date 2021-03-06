package utils

import (
	"github.com/labstack/echo"
	"net/http"
	"fmt"
	"encoding/json"
	"github.com/beewit/beekit/utils/enum"
)

type ResultParam struct {
	Ret  int64       `json:"ret"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

const (
	//成功
	SUCCESS_CODE = 200
	//错误
	ERROR_CODE = 400
	//授权失败
	AUTH_FAIL = 401
	//登陆无效
	LOGIN_INVALID_CODE = 402
	//失败
	FAIL_DATA = 403
	//数据为空
	NULL_DATA = 404
	//不是会员
	NOT_MEMBER = 500
	//会员续费通知
	MEMBER_RENEW = 501

	AUTH_WECHAT_FAIL = 5401
)

func ToResultParam(b []byte) ResultParam {
	var rp ResultParam
	err := json.Unmarshal(b[:], &rp)
	if err != nil {
		println(err.Error())
		return ResultParam{}
	}
	return rp
}

func SuccessRespone(c echo.Context, data string) error {
	return c.HTML(http.StatusOK, data)
}

func Success(c echo.Context, msg string, data interface{}) error {
	return Result(c, SUCCESS_CODE, msg, data)
}

func SuccessNullMsg(c echo.Context, data interface{}) error {
	return Result(c, SUCCESS_CODE, "", data)
}
func SuccessNull(c echo.Context, msg string) error {
	return Result(c, SUCCESS_CODE, msg, nil)
}

func Error(c echo.Context, msg string, data interface{}) error {
	return Result(c, ERROR_CODE, msg, data)
}

func ErrorNull(c echo.Context, msg string) error {
	return Result(c, ERROR_CODE, msg, nil)
}

func NullData(c echo.Context) error {
	return Result(c, NULL_DATA, "暂无数据", nil)
}

func AuthFail(c echo.Context, msg string) error {
	return Result(c, AUTH_FAIL, msg, nil)
}

func AuthFailNull(c echo.Context) error {
	return Result(c, AUTH_FAIL, "未登录或登陆已失效", nil)
}

func AuthWechatFailNull(c echo.Context) error {
	return Result(c, AUTH_WECHAT_FAIL, "微信未能获取授权", nil)
}

func ResultApi(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusOK, data)
}

func Result(c echo.Context, ret int64, msg string, data interface{}) error {
	resultMap := map[string]interface{}{
		"ret":  ret,
		"msg":  msg,
		"data": data,
	}
	return c.JSON(http.StatusOK, resultMap)
}

func Alert(c echo.Context, tip string) error {
	return RedirectAndAlert(c, tip, "")
}

func RedirectAndAlert(c echo.Context, tip, url string) error {
	var js string
	if tip != "" {
		js += fmt.Sprintf("alert('%v');", tip)
	}
	js += fmt.Sprintf("parent.location.href = '%v';", url)
	return ResultHtml(c, fmt.Sprintf("<script>%v</script>", js))
}

func Redirect(c echo.Context, url string) error {
	return c.Redirect(http.StatusMovedPermanently, url)
}

func ResultHtml(c echo.Context, html string) error {
	return c.HTML(http.StatusOK, html)
}

func ResultString(c echo.Context, str string) error {
	return c.String(http.StatusOK, str)
}

func ActionLogs(c echo.Context, t string, accountId int64) map[string]interface{} {
	return ActionLogsMap(c, t, "", accountId)
}

func ActionLogsMap(c echo.Context, t, remarks string, accountId int64) map[string]interface{} {
	actionLog := map[string]interface{}{}
	if c.FormValue("actionSource") == "" {
		actionLog["source"] = enum.ACTION_PC
	} else {
		actionLog["source"] = c.FormValue("actionSource")
	}
	if remarks == "" {
		actionLog["remarks"] = c.FormValue("actionRemarks")
	} else {
		actionLog["remarks"] = remarks
	}
	actionLog["id"] = ID()
	actionLog["account_id"] = accountId
	actionLog["user_agent"] = c.Request().UserAgent()
	actionLog["type"] = t
	actionLog["ct_ip"] = c.RealIP()
	actionLog["ct_time"] = CurrentTime()
	actionLog["device"] = c.FormValue("actionDevice")
	return actionLog
}
