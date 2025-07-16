package main

import (
	"os"

	"github.com/cosmos/cosmos-sdk/server"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"

	"github.com/Crocodile-ark/gxrchaind/app"
	"github.com/Crocodile-ark/gxrchaind/cmd/gxrchaind/cmd"
)

func main() {
	// Set the default bond denomination to ugen before starting
	app.SetDefaultBondDenom()

	rootCmd, _ := cmd.NewRootCmd()

	if err := svrcmd.Execute(rootCmd, "", app.DefaultNodeHome); err != nil {
		switch e := err.(type) {
		case server.ErrorCode:
			os.Exit(e.Code)

		default:
			os.Exit(1)
		}
	}
}