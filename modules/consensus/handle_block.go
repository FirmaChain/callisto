package consensus

import (
	"fmt"

	juno "github.com/forbole/juno/v6/types"

	"github.com/rs/zerolog/log"

	cmttypes "github.com/cometbft/cometbft/rpc/core/types"
)

// HandleBlock implements modules.Module
func (m *Module) HandleBlock(
	b *cmttypes.ResultBlock, _ *cmttypes.ResultBlockResults, _ []*juno.Transaction, _ *cmttypes.ResultValidators,
) error {
	err := m.updateBlockTimeFromGenesis(b)
	if err != nil {
		log.Error().Str("module", "consensus").Int64("height", b.Block.Height).
			Err(err).Msg("error while updating block time from genesis")
	}

	return nil
}

// updateBlockTimeFromGenesis insert average block time from genesis
func (m *Module) updateBlockTimeFromGenesis(block *cmttypes.ResultBlock) error {
	log.Trace().Str("module", "consensus").Int64("height", block.Block.Height).
		Msg("updating block time from genesis")

	genesis, err := m.db.GetGenesis()
	if err != nil {
		return fmt.Errorf("error while getting genesis: %s", err)
	}
	if genesis == nil {
		return fmt.Errorf("genesis table is empty")
	}

	newBlockTime := block.Block.Time.Sub(genesis.Time).Seconds() / float64(block.Block.Height-genesis.InitialHeight)
	return m.db.SaveAverageBlockTimeGenesis(newBlockTime, block.Block.Height)
}
