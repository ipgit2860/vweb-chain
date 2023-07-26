package keeper

import (
	"context"
	"encoding/json"

	"vcoa/x/vnft/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) BurnNFT(goCtx context.Context, msg *types.MsgBurnNFT) (*types.MsgBurnNFTResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check nft exists
	exists := k.nftKeeper.HasNFT(ctx, msg.ClassId, msg.NftId)
	if !exists {
		return nil, types.ErrNftNotFound.Wrapf("Class %s NFT %s does not exist", msg.ClassId, msg.NftId)
	}

	// Check user is owner
	owner := k.nftKeeper.GetOwner(ctx, msg.ClassId, msg.NftId)
	if err := k.assertBech32EqualsAccAddress(msg.Creator, owner); err != nil {
		return nil, err
	}

	// Check class is set to burnable
	class, hasClass := k.nftKeeper.GetClass(ctx, msg.ClassId)
	if !hasClass {
		return nil, types.ErrNftClassNotFound.Wrapf("Class id %s not found", msg.ClassId)
	}
	// Unmarshal class data
	var classData types.ClassData
	if err := classData.Unmarshal(class.Data.Value); err != nil {
		return nil, types.ErrFailedToUnmarshalData.Wrapf(err.Error())
	}

	type format struct {
		CField    []string `protobuf:"bytes,1,opt,name=cField,proto3" json:"cField,omitempty"`
		CType     []string `protobuf:"bytes,2,opt,name=cType,proto3" json:"cType,omitempty"`
		COptional []bool   `protobuf:"bytes,3,opt,name=cOptional,proto3" json:"cOptional,omitempty"`
	}

	type metadata struct {
		Version      int             `protobuf:"bytes,1,opt,name=version,proto3" json:"version,omitempty"`
		Format       types.JsonInput `protobuf:"bytes,2,opt,name=format,proto3" json:"format,omitempty"`
		MultiAuthor  bool            `protobuf:"bytes,3,opt,name=multiAuthor,proto3" json:"multiAuthor,omitempty"`
		Burnable     bool            `protobuf:"bytes,4,opt,name=burnable,proto3" json:"burnable,omitempty"`
		Transferable bool            `protobuf:"bytes,5,opt,name=transferable,proto3" json:"transferable,omitempty"`
	}

	var classMeta metadata
	json.Unmarshal([]byte(classData.Metadata), &classMeta)

	if !classMeta.Burnable {
		return nil, types.ErrNftNotBurnable.Wrapf("NFT of class %s is not burnable", class.Id)
	}

	// Burn NFT
	err := k.nftKeeper.Burn(ctx, msg.ClassId, msg.NftId)
	if err != nil {
		return nil, types.ErrFailedToBurnNFT.Wrapf("%s", err.Error())
	}

	// Emit event
	ctx.EventManager().EmitTypedEvent(&types.EventBurnNFT{
		ClassId: class.Id,
		NftId:   msg.NftId,
		Owner:   owner.String(),
	})

	return &types.MsgBurnNFTResponse{}, nil
}
