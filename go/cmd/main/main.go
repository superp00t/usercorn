package main

import (
	"github.com/superp00t/usercorn/go/cmd"

	_ "github.com/superp00t/usercorn/go/cmd/run"

	_ "github.com/superp00t/usercorn/go/cmd/cfg"
	_ "github.com/superp00t/usercorn/go/cmd/cgc"
	_ "github.com/superp00t/usercorn/go/cmd/com"
	_ "github.com/superp00t/usercorn/go/cmd/fuzz"
	_ "github.com/superp00t/usercorn/go/cmd/imgtrace"
	_ "github.com/superp00t/usercorn/go/cmd/repl"
	_ "github.com/superp00t/usercorn/go/cmd/shellcode"
	_ "github.com/superp00t/usercorn/go/cmd/trace"
)

func main() { cmd.Main() }
