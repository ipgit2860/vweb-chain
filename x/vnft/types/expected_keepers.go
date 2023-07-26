package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	nft "github.com/cosmos/cosmos-sdk/x/nft"
)

type NftKeeper interface {
	SaveClass(ctx sdk.Context, class nft.Class) error
	GetClass(ctx sdk.Context, classID string) (nft.Class, bool)
	HasClass(ctx sdk.Context, classID string) bool
	GetBalance(ctx sdk.Context, classID string, owner sdk.AccAddress) (balance uint64)
	HasNFT(ctx sdk.Context, classID, id string) bool
	Mint(ctx sdk.Context, token nft.NFT, receiver sdk.AccAddress) error
	GetTotalSupply(ctx sdk.Context, classID string) uint64
	Burn(ctx sdk.Context, classID string, nftID string) error
	GetOwner(ctx sdk.Context, classID string, nftID string) sdk.AccAddress
	// Methods imported from nft should be defined here
}

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) types.AccountI
	// Methods imported from account should be defined here
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	// Methods imported from bank should be defined here
}
