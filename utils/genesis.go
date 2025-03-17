package utils

import (
	"fmt"

	tmjson "github.com/cometbft/cometbft/libs/json"
	tmos "github.com/cometbft/cometbft/libs/os"
	cmtypes "github.com/cometbft/cometbft/types"
	"github.com/forbole/juno/v6/node"
	"github.com/forbole/juno/v6/types/config"
)

// ReadGenesis reads the genesis data based on the given config
func ReadGenesis(config config.Config, node node.Node) (*cmtypes.GenesisDoc, error) {
	if config.Parser.GenesisFilePath != "" {
		return readGenesisFromFilePath(config.Parser.GenesisFilePath)
	}

	return readGenesisFromNode(node)
}

func readGenesisFromFilePath(path string) (*cmtypes.GenesisDoc, error) {
	bz, err := tmos.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read genesis file: %s", err)
	}

	var genDoc cmtypes.GenesisDoc
	err = tmjson.Unmarshal(bz, &genDoc)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal genesis doc: %s", err)
	}

	return &genDoc, nil
}

func readGenesisFromNode(node node.Node) (*cmtypes.GenesisDoc, error) {
	response, err := node.Genesis()
	if err != nil {
		return nil, fmt.Errorf("failed to get genesis: %s", err)
	}

	return response.Genesis, nil
}
