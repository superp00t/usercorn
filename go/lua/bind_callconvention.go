package lua

import (
	lua "github.com/lunixbochs/luaish"
	"github.com/superp00t/usercorn/go/models"
)

var callConv = map[string]models.Callconvention{
	"__cdecl":    models.Cdecl,
	"__fastcall": models.Fastcall,
	"__thiscall": models.Thiscall,
}

func bindCallConv(L *LuaRepl) error {
	for k, v := range callConv {
		L.SetGlobal(k, lua.LInt(v))
	}
	return nil
}
