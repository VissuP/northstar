package lua

func OpenSecureOs(L *LState) int {
	osmod := L.RegisterModule(OsLibName, secureOsFuncs)
	L.Push(osmod)
	return 1
}

var secureOsFuncs = map[string]LGFunction{
	"clock":    osClock,
	"difftime": osDiffTime,
	"date":     osDate,
	"time":     osTime,
}
