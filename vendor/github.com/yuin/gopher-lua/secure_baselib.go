package lua

func OpenSecureBase(L *LState) int {
	global := L.Get(GlobalsIndex).(*LTable)
	L.SetGlobal("_G", global)
	basemod := L.RegisterModule("_G", secureBaseFuncs)
	global.RawSetString("ipairs", L.NewClosure(baseIpairs, L.NewFunction(ipairsaux)))
	global.RawSetString("pairs", L.NewClosure(basePairs, L.NewFunction(pairsaux)))
	L.Push(basemod)
	return 1
}

var secureBaseFuncs = map[string]LGFunction{
	"error":    baseError,
	"tonumber": baseToNumber,
	"tostring": baseToString,
	"type":     baseType,
	"next":     baseNext,
	"unpack":   baseUnpack,
	// loadlib
	"module":  loModule,
	"require": loRequire,
}
