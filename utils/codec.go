package utils

import (
	"sync"

	evidencetypes "cosmossdk.io/x/evidence/types"
	"cosmossdk.io/x/feegrant"
	nft "cosmossdk.io/x/nft"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/std"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	consensustypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govv1types "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govv1beta1types "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/cosmos/cosmos-sdk/x/group"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	paramsproposaltypes "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/gogoproto/proto"
	packetforwardtypes "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v8/packetforward/types"
	icacontrollertypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/controller/types"
	icahosttypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/host/types"
	icatypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	ibctypes "github.com/cosmos/ibc-go/v8/modules/core/types"
	ibctm "github.com/cosmos/ibc-go/v8/modules/light-clients/07-tendermint"
)

var (
	once sync.Once
	cdc  *codec.ProtoCodec
)

func GetCodec() codec.Codec {
	once.Do(func() {
		interfaceRegistry := codectypes.NewInterfaceRegistry()
		authtypes.RegisterInterfaces(interfaceRegistry)
		authz.RegisterInterfaces(interfaceRegistry)
		banktypes.RegisterInterfaces(interfaceRegistry)
		stakingtypes.RegisterInterfaces(interfaceRegistry)
		minttypes.RegisterInterfaces(interfaceRegistry)
		distrtypes.RegisterInterfaces(interfaceRegistry)
		govv1types.RegisterInterfaces(interfaceRegistry)
		govv1beta1types.RegisterInterfaces(interfaceRegistry)
		paramsproposaltypes.RegisterInterfaces(interfaceRegistry)
		crisistypes.RegisterInterfaces(interfaceRegistry)
		slashingtypes.RegisterInterfaces(interfaceRegistry)
		feegrant.RegisterInterfaces(interfaceRegistry)
		group.RegisterInterfaces(interfaceRegistry)
		ibctypes.RegisterInterfaces(interfaceRegistry)
		ibctm.RegisterInterfaces(interfaceRegistry)
		packetforwardtypes.RegisterInterfaces(interfaceRegistry)
		icacontrollertypes.RegisterInterfaces(interfaceRegistry)
		icahosttypes.RegisterInterfaces(interfaceRegistry)
		icatypes.RegisterInterfaces(interfaceRegistry)
		upgradetypes.RegisterInterfaces(interfaceRegistry)
		evidencetypes.RegisterInterfaces(interfaceRegistry)
		ibctransfertypes.RegisterInterfaces(interfaceRegistry)
		vestingtypes.RegisterInterfaces(interfaceRegistry)
		consensustypes.RegisterInterfaces(interfaceRegistry)
		wasmtypes.RegisterInterfaces(interfaceRegistry)
		nft.RegisterInterfaces(interfaceRegistry)

		std.RegisterInterfaces(interfaceRegistry)

		cdc = codec.NewProtoCodec(interfaceRegistry)
	})
	return cdc
}

// UnpackMessage unpacks a message from a byte slice
func UnpackMessage[T proto.Message](cdc codec.Codec, bz []byte, ptr T) T {
	var any codectypes.Any
	cdc.MustUnmarshalJSON(bz, &any)
	var cosmosMsg sdk.Msg
	if err := cdc.UnpackAny(&any, &cosmosMsg); err != nil {
		panic(err)
	}
	return cosmosMsg.(T)
}
