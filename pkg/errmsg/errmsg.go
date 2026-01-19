// package errno provides business error codes and messages.
// 注意：这里是“业务码”，不是 HTTP status code。
// HTTP 状态码建议在网关/handler 层决定（200/400/401/500...）。

package errno

import "fmt"

type Code int32

const (
	// ===== 通用 =====
	CodeOK           Code = 0
	CodeBadRequest   Code = 1000 // 参数/请求不合法
	CodeUnauthorized Code = 1001 // 未登录 / 鉴权失败
	CodeForbidden    Code = 1002 // 权限不足
	CodeInternal     Code = 9001 // 内部错误（兜底）
	CodeServerBusy   Code = 9000 // 服务繁忙/限流/依赖不可用

	// ===== 用户/登录 ===== (1100~1299)
	CodeUserAlreadyExists       Code = 1100
	CodeUserNotFound            Code = 1101
	CodeUsernameOrPasswordWrong Code = 1102
	CodeUserAlreadyLoggedIn     Code = 1103

	// ===== Token/JWT ===== (1200~1299)
	CodeTokenMissing       Code = 1200
	CodeTokenInvalid       Code = 1201
	CodeTokenExpired       Code = 1202
	CodeTokenRefreshFailed Code = 1203

	// ===== 社区 ===== (2000~2099)
	CodeCommunityAlreadyExists Code = 2000
	CodeCommunityNotFound      Code = 2001

	// ===== 帖子 ===== (3000~3099)
	CodePostAlreadyExists Code = 3000
	CodePostNotFound      Code = 3001

	// ===== 投票 ===== (4000~4099)
	CodeVoteRepeated    Code = 4000
	CodeVoteTimeExpired Code = 4001
)

var msg = map[Code]string{
	CodeOK:           "OK",
	CodeBadRequest:   "请求参数错误",
	CodeUnauthorized: "未登录或登录已失效",
	CodeForbidden:    "权限不足",
	CodeServerBusy:   "服务繁忙，请稍后再试",
	CodeInternal:     "内部错误",

	CodeUserAlreadyExists:       "用户名已存在",
	CodeUserNotFound:            "用户不存在",
	CodeUsernameOrPasswordWrong: "用户名或密码错误",
	CodeUserAlreadyLoggedIn:     "已登录",

	CodeTokenMissing:       "TOKEN缺失",
	CodeTokenInvalid:       "TOKEN无效",
	CodeTokenExpired:       "TOKEN已过期",
	CodeTokenRefreshFailed: "TOKEN刷新失败",

	CodeCommunityAlreadyExists: "社区已存在",
	CodeCommunityNotFound:      "社区不存在",

	CodePostAlreadyExists: "帖子已存在",
	CodePostNotFound:      "帖子不存在",

	CodeVoteRepeated:    "请勿重复投票",
	CodeVoteTimeExpired: "投票时间已过",
}

func Message(code Code) string {
	if s, ok := msg[code]; ok {
		return s
	}
	return msg[CodeInternal]
}

// BizError 业务错误
type BizError struct {
	Code  Code
	Msg   string
	Cause error
}

func (e *BizError) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.Cause != nil {
		// 给日志用：带 cause；对外回包一般只用 Code/Msg
		return fmt.Sprintf("code=%d msg=%s cause=%v", e.Code, e.Msg, e.Cause)
	}
	return fmt.Sprintf("code=%d msg=%s", e.Code, e.Msg)
}

func New(code Code) *BizError {
	return &BizError{Code: code, Msg: Message(code)}
}

func NewMsg(code Code, m string) *BizError {
	return &BizError{Code: code, Msg: m}
}

func Wrap(code Code, cause error) *BizError {
	return &BizError{Code: code, Msg: Message(code), Cause: cause}
}

func (e *BizError) Unwrap() error {
	return e.Cause
}
