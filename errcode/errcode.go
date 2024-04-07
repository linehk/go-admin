package errcode

const (
	Parse    int32 = 20000
	Database int32 = 20001
	Convert  int32 = 20002
)

var msg = map[int32]string{
	Parse:    "parse error",
	Database: "database error",
	Convert:  "convert error",
}

func Msg(e int32) string {
	return msg[e]
}
