package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/vnft module sentinel errors
var (
	ErrSample                 = sdkerrors.Register(ModuleName, 1100, "sample error")
	ErrFailedToSaveClass      = sdkerrors.Register(ModuleName, 3, "Failed to save class")
	ErrFailedToMarshalData    = sdkerrors.Register(ModuleName, 4, "Failed to marshal data")
	ErrNftClassNotFound       = sdkerrors.Register(ModuleName, 5, "NFT Class not found")
	ErrFailedToUnmarshalData  = sdkerrors.Register(ModuleName, 6, "Failed to unmarshal data")
	ErrFailedToMintNFT        = sdkerrors.Register(ModuleName, 7, "Failed to mint NFT")
	ErrFailedToGetUserAddress = sdkerrors.Register(ModuleName, 8, "Failed to obtain user address")
	ErrFailedToCreateNftId    = sdkerrors.Register(ModuleName, 9, "Failed to Create NFT ID")
	ErrOwnerAddressUsedBefore = sdkerrors.Register(ModuleName, 10, "Owner Address is used before")
	ErrNftNotFound            = sdkerrors.Register(ModuleName, 11, "NFT was not found")
	ErrFailedToBurnNFT        = sdkerrors.Register(ModuleName, 12, "Failed to Burn the NFT")
	ErrNftNotBurnable         = sdkerrors.Register(ModuleName, 13, "NFT is not burnable")
)
