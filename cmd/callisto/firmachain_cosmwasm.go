package main

import (
	cdc "github.com/cosmos/cosmos-sdk/codec"

	wasmvmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/forbole/juno/v6/types"
)

func FirmaChainCosmWasmMessagesParser(tx *types.Transaction) ([]string, error) {
	// Create a slice to hold the addresses.
	var addresses []string

	// Loop over each message in the transaction body.
	for _, anyMsg := range tx.Body.Messages {
		// Unpack the Any message into a concrete sdk.Msg.
		anyMsgByte := anyMsg.GetBytes()
		var cosmosMsg sdk.Msg
		if err := cdc.NewLegacyAmino().Unmarshal(anyMsgByte, &cosmosMsg); err != nil {
			return nil, err
		}

		// Type switch to handle different message types.
		switch msg := cosmosMsg.(type) {
		case *wasmvmtypes.MsgInstantiateContract:
			addresses = append(addresses, msg.Sender, msg.Admin)
		case *wasmvmtypes.MsgStoreCode:
			addresses = append(addresses, msg.Sender)
		case *wasmvmtypes.MsgExecuteContract:
			addresses = append(addresses, msg.Sender, msg.Contract)
		case *wasmvmtypes.MsgUpdateAdmin:
			addresses = append(addresses, msg.Sender, msg.Contract, msg.NewAdmin)
		case *wasmvmtypes.MsgClearAdmin:
			addresses = append(addresses, msg.Sender, msg.Contract)
		}
	}
	return addresses, nil
}
