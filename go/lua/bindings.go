package lua

import (
	"strconv"

	lua "github.com/lunixbochs/luaish"
)

func (L *LuaRepl) printFunc(_ *lua.LState) int {
	L.PrettyPrint(L.getArgs(), false)
	return 0
}

func (L *LuaRepl) intFunc(_ *lua.LState) int {
	switch v := L.CheckAny(1).(type) {
	case lua.LString:
		n, err := strconv.ParseInt(string(v), 0, 64)
		if err == nil {
			L.Push(lua.LInt(n))
			return 1
		}
	case lua.LFloat:
		L.Push(lua.LInt(v))
		return 1
	case lua.LInt:
		L.Push(v)
		return 1
	}
	return 0
}

func (L *LuaRepl) loadBindings() error {
	// populate predefined globals as "_builtins" so they can be skipped in help()/dir()
	builtins := L.NewTable()
	g := L.GetGlobal("_G")
	if s, ok := g.(*lua.LTable); ok {
		s.ForEach(func(k, v lua.LValue) {
			builtins.RawSet(k, lua.LTrue)
		})
	}
	L.SetGlobal("_builtins", builtins)

	print := L.NewFunction(L.printFunc)
	L.SetGlobal("print", print)

	toint := L.NewFunction(L.intFunc)
	L.SetGlobal("int", toint)

	if err := bindCpu(L); err != nil {
		return err
	} else if err := bindCallConv(L); err != nil {
		return err
	} else if err := bindUsercorn(L); err != nil {
		return err
	} else if err := L.DoString(sugarRc); err != nil {
		return err
	} else if err := L.DoString(cmdRc); err != nil {
		return err
	}
	return nil
}
