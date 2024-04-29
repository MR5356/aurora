package response

const (
	CodeSuccess      = "00000"
	CodeNotFound     = "B0001"
	CodeParamsError  = "B0002"
	CodeNoPermission = "B0003"
	CodeServerError  = "A0001"

	MessageSuccess      = "success"
	MessageNotFound     = "not found"
	MessageParamsError  = "params error"
	MessageNoPermission = "no permission"
	MessageServerError  = "server error"

	MessageUnknown = "unknown error"
)

var msgMap = map[string]string{
	CodeSuccess:      MessageSuccess,
	CodeNotFound:     MessageNotFound,
	CodeParamsError:  MessageParamsError,
	CodeServerError:  MessageServerError,
	CodeNoPermission: MessageNoPermission,
}

func message(code string) string {
	if msg, ok := msgMap[code]; ok {
		return msg
	}
	return MessageUnknown
}
