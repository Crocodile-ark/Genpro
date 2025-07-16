package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"

	"github.com/Crocodile-ark/gxrchaind/x/feerouter/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string) *cobra.Command {
	// Group feerouter queries under a subcommand
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		CmdQueryParams(),
		CmdQueryFeeStats(),
		CmdQueryLPPools(),
	)

	return cmd
}

// CmdQueryParams implements the params query command.
func CmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Args:  cobra.NoArgs,
		Short: "Query the current feerouter parameters",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Params(cmd.Context(), &types.QueryParamsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// CmdQueryFeeStats implements the fee stats query command.
func CmdQueryFeeStats() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fee-stats",
		Args:  cobra.NoArgs,
		Short: "Query the fee collection and distribution statistics",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.FeeStats(cmd.Context(), &types.QueryFeeStatsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// CmdQueryLPPools implements the LP pools query command.
func CmdQueryLPPools() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lp-pools",
		Args:  cobra.NoArgs,
		Short: "Query all registered LP pools",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.LPPools(cmd.Context(), &types.QueryLPPoolsRequest{
				Pagination: pageReq,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "LP pools")

	return cmd
}