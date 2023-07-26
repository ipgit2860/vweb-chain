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

func CmdMintNft() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mint-nft [class-id] [id] [uri] [uri-hash] [data]",
		Short: "Broadcast message mint-nft",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argClassId := args[0]
			argId := args[1]
			argUri := args[2]
			argUriHash := args[3]
			nftData, err := readJsonFile[types.NftData](args[4])
			if nftData == nil || err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgMintNft(
				clientCtx.GetFromAddress().String(),
				argClassId,
				argId,
				argUri,
				argUriHash,
				*nftData,
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
