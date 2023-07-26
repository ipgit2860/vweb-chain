package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	"github.com/cosmos/gogoproto/proto"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCreateClass{}, "vnft/CreateClass", nil)
	cdc.RegisterConcrete(&ClassData{}, "vnft/ClassData", nil)
	cdc.RegisterConcrete(&MsgMintNft{}, "vnft/MintNft", nil)
	cdc.RegisterConcrete(&NftData{}, "vnft/NftData", nil)
	// cdc.RegisterConcrete(&MsgBurnNFT{}, "vnft/BurnNFT", nil)
	// this line is used by starport scaffolding # 2
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateClass{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgMintNft{},
	)
	// registry.RegisterImplementations((*sdk.Msg)(nil),
	// 	&MsgBurnNFT{},
	// )
	// this line is used by starport scaffolding # 3

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
	registry.RegisterImplementations((*proto.Message)(nil), &ClassData{})
	registry.RegisterImplementations((*proto.Message)(nil), &NftData{})
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
