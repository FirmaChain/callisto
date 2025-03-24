package main

//TODO: fix import firmachain
/*import (
	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
	nfttypes "github.com/firmachain/firmachain/x/nft/types"
)

func FirmaChainNFTMessagesParser(_ codec.Codec, cosmosMsg sdk.Msg) ([]string, error) {
	switch msg := cosmosMsg.(type) {
	case *nfttypes.MsgMint:
		return []string{msg.Owner}, nil

	case *nfttypes.MsgBurn:
		return []string{msg.Owner}, nil

	case *nfttypes.MsgTransfer:
		return []string{msg.Owner, msg.ToAddress}, nil
	}

	return nil, nil
}*/
