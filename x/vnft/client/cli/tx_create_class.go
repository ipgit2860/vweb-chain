package cli

import (
	"strconv"

	"vcoa/x/vnft/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdCreateClass() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-class [name] [symbol] [description] [uri] [uri-hash] [data]",
		Short: "Broadcast message create-class",
		Args:  cobra.ExactArgs(6),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argName := args[0]
			argSymbol := args[1]
			argDescription := args[2]
			argUri := args[3]
			argUriHash := args[4]
			classData, err := readJsonFile[types.ClassData](args[5])
			if classData == nil || err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreateClass(
				clientCtx.GetFromAddress().String(),
				argName,
				argSymbol,
				argDescription,
				argUri,
				argUriHash,
				*classData,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
