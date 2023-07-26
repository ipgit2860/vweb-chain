package vnft

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"vcoa/testutil/sample"
	vnftsimulation "vcoa/x/vnft/simulation"
	"vcoa/x/vnft/types"
)

// avoid unused import issue
var (
	_ = sample.AccAddress
	_ = vnftsimulation.FindAccount
	_ = simulation.MsgEntryKind
	_ = baseapp.Paramspace
	_ = rand.Rand{}
)

const (
	opWeightMsgCreateClass = "op_weight_msg_create_class"
	// TODO: Determine the simulation weight value
	defaultWeightMsgCreateClass int = 100

	opWeightMsgMintNft = "op_weight_msg_mint_nft"
	// TODO: Determine the simulation weight value
	defaultWeightMsgMintNft int = 100

	// this line is used by starport scaffolding # simapp/module/const
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	vnftGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&vnftGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ sdk.StoreDecoderRegistry) {}

// ProposalContents doesn't return any content functions for governance proposals.
func (AppModule) ProposalContents(_ module.SimulationState) []simtypes.WeightedProposalContent {
	return nil
}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	var weightMsgCreateClass int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgCreateClass, &weightMsgCreateClass, nil,
		func(_ *rand.Rand) {
			weightMsgCreateClass = defaultWeightMsgCreateClass
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateClass,
		vnftsimulation.SimulateMsgCreateClass(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgMintNft int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgMintNft, &weightMsgMintNft, nil,
		func(_ *rand.Rand) {
			weightMsgMintNft = defaultWeightMsgMintNft
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgMintNft,
		vnftsimulation.SimulateMsgMintNft(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{
		simulation.NewWeightedProposalMsg(
			opWeightMsgCreateClass,
			defaultWeightMsgCreateClass,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				vnftsimulation.SimulateMsgCreateClass(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgMintNft,
			defaultWeightMsgMintNft,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				vnftsimulation.SimulateMsgMintNft(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		// this line is used by starport scaffolding # simapp/module/OpMsg
	}
}
