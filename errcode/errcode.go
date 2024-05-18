package errcode

const (
	Parse    int32 = 20000
	Database int32 = 20001
	Convert  int32 = 20002
	Validate int32 = 20003

	UsernameOccupy int32 = 30000
	UserNotExist   int32 = 30001

	RoleCodeOccupy int32 = 40000
	RoleNotExist   int32 = 40001

	MenuCodeOccupy int32 = 50000
	MenuNotExist   int32 = 50001
)

var msg = map[int32]string{
	Parse:    "parse error",
	Database: "database error",
	Convert:  "convert error",
	Validate: "validate error",

	UsernameOccupy: "username occupy",
	UserNotExist:   "user not exist",

	RoleCodeOccupy: "role code occupy",
	RoleNotExist:   "role not exist",

	MenuCodeOccupy: "menu code occupy",
	MenuNotExist:   "menu not exist",
}

func Msg(e int32) string {
	return msg[e]
}
