// package errmsg provides business error codes and messages.
// 注意：这里是“业务码”，不是 HTTP status code。
// HTTP 状态码建议在网关/handler 层决定（200/400/401/500...）。

package errmsg

import (
	"errors"
	"fmt"
)

// Code 表示业务错误码（不是 HTTP 状态码）。
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

// BizError 表示业务错误，包含业务码/对外消息，以及可选的底层原因（用于排查）。
type BizError struct {
	Code  Code
	Msg   string
	Cause error
}

// Message 根据业务码返回默认文案；如果业务码不存在，则返回通用“内部错误”文案。
func Message(code Code) string {
	if s, ok := msg[code]; ok {
		return s
	}
	return "内部错误"
}

// Error 实现 error 接口。
// 主要用于日志/调试输出；对外回包一般只使用 Code/Msg，不应直接透出 Cause。
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

// New 根据业务码创建 BizError，消息使用 Message(code) 的默认文案。
func New(code Code) *BizError {
	return &BizError{Code: code, Msg: Message(code)}
}

// NewMsg 根据业务码创建 BizError，并使用自定义消息文案。
// 适合参数校验等场景；注意不要把敏感信息（密码、内部异常细节）透出给用户。
func NewMsg(code Code, m string) *BizError {
	return &BizError{Code: code, Msg: m}
}

// Wrap 用业务码包装底层错误 cause，并使用默认文案 Message(code)。
// 如果 cause 为 nil，则退化为 New(code)。
func Wrap(code Code, cause error) *BizError {
	if cause == nil {
		return New(code)
	}
	return &BizError{Code: code, Msg: Message(code), Cause: cause}
}

// Unwrap 返回底层原因错误，用于支持 errors.Is / errors.As 的错误链解析。
func (e *BizError) Unwrap() error {
	return e.Cause
}

// From 将任意 error 统一转换成（业务码，消息）。
// - 如果 err 是 BizError（或被 BizError 包装），则提取其 Code/Msg（Msg 为空则用默认文案）
// - 否则认为是系统内部错误，返回 CodeInternal 及默认内部错误文案
func From(err error) (Code, string) {
	if err == nil {
		return CodeOK, Message(CodeOK)
	}
	var be *BizError
	if errors.As(err, &be) {
		// be.Msg 允许是自定义 msg
		if be.Msg != "" {
			return be.Code, be.Msg
		}
		return be.Code, Message(be.Code)
	}
	return CodeInternal, Message(CodeInternal)
}

// Newf 创建带格式化消息的 BizError（自定义文案）。
// 用于需要携带少量上下文信息的场景（仍需避免泄露敏感信息）。
func Newf(code Code, format string, args ...any) *BizError {
	return &BizError{Code: code, Msg: fmt.Sprintf(format, args...)}
}

// Wrapf 用业务码包装底层错误 cause，并使用格式化消息作为对外文案。
// 如果 cause 为 nil，则退化为 Newf(code, ...)。
func Wrapf(code Code, cause error, format string, args ...any) *BizError {
	if cause == nil {
		return Newf(code, format, args...)
	}
	return &BizError{Code: code, Msg: fmt.Sprintf(format, args...), Cause: cause}
}

// IsBiz 判断 err 是否为 BizError（或错误链中包含 BizError）。
// 常用于入口层决定日志级别：业务错误一般 warn/info；系统错误一般 error。
func IsBiz(err error) bool {
	var be *BizError
	return errors.As(err, &be)
}

// CodeOf 获取 err 对应的业务码：
// - err 为 nil：返回 CodeOK
// - err 为 BizError（或链中包含 BizError）：返回其 Code
// - 否则：返回 CodeInternal
func CodeOf(err error) Code {
	c, _ := From(err)
	return c
}
