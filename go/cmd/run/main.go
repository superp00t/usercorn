package run

import (
	"os"

	"github.com/superp00t/usercorn/go/cmd"
)

func Main(args []string) {
	os.Exit(cmd.NewUsercornCmd().Run(args, os.Environ()))
}

func init() { cmd.Register("run", "execute a binary", Main) }
