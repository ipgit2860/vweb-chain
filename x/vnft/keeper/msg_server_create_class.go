package keeper

import (
	"context"
	"encoding/json"
	"fmt"

	"vcoa/x/vnft/types"

	sdkctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/nft"
)

func (k msgServer) CreateClass(goCtx context.Context, msg *types.MsgCreateClass) (*types.MsgCreateClassResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	userAddress, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, err
	}

	classId, err := types.NewClassId(userAddress.Bytes(), 0)
	if err != nil {
		return nil, err
	}

	for i := 1; k.nftKeeper.HasClass(ctx, classId); i++ {
		classId, err = types.NewClassId(userAddress.Bytes(), i)
		if err != nil {
			return nil, err
		}
	}

	databyte := []byte(msg.Data.Metadata)
	var meta types.JsonInput
	if json.Unmarshal(databyte, &meta) != nil {
		return nil, types.ErrFailedToMarshalData.Wrapf("Unmarshal %s", err.Error())
	}

	classData := types.ClassData{
		Metadata: meta,
	}

	fmt.Println("-----------------------------------------------------------------------------")
	fmt.Println()
	fmt.Println()
	fmt.Println("classId: ", classId)
	fmt.Println("classData: ", classData)
	fmt.Println()
	fmt.Println()
	fmt.Println("-----------------------------------------------------------------------------")

	classDataInAny, err := sdkctypes.NewAnyWithValue(&classData)
	if err != nil {
		return nil, types.ErrFailedToMarshalData.Wrapf("%s", err.Error())
	}

	except := nft.Class{
		Id:          classId,
		Name:        msg.Name,
		Symbol:      msg.Symbol,
		Description: msg.Description,
		Uri:         msg.Uri,
		UriHash:     msg.UriHash,
		Data:        classDataInAny,
	}

	err = k.nftKeeper.SaveClass(ctx, except)
	if err != nil {
		return nil, types.ErrFailedToSaveClass.Wrapf("%s", err.Error())
	}

	return &types.MsgCreateClassResponse{
		// Class: except,
	}, nil
}
