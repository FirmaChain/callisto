package main

import (
	"github.com/cosmos/cosmos-sdk/codec"

	wasmvmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func FirmaChainCosmWasmMessagesParser(_ codec.Codec, cosmosMsg sdk.Msg) ([]string, error) {

	switch msg := cosmosMsg.(type) {
	case *wasmvmtypes.MsgInstantiateContract:
		var stringArray = []string{msg.Sender, msg.Admin}
		return stringArray, nil

	case *wasmvmtypes.MsgStoreCode:
		var stringArray = []string{msg.Sender}
		return stringArray, nil

	case *wasmvmtypes.MsgExecuteContract:
		var stringArray = []string{msg.Sender, msg.Contract}
		return stringArray, nil

	case *wasmvmtypes.MsgUpdateAdmin:
		var stringArray = []string{msg.Sender, msg.Contract, msg.NewAdmin}
		return stringArray, nil

	case *wasmvmtypes.MsgClearAdmin:
		var stringArray = []string{msg.Sender, msg.Contract}
		return stringArray, nil

	}

	return nil, nil
}
