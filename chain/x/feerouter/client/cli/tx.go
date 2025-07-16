package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"

	"github.com/Crocodile-ark/gxrchaind/x/feerouter/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	// Add transaction commands here if needed
	// For now, feerouter module only has automatic functions

	return cmd
}