package modules

import (
	"github.com/forbole/callisto/v4/modules/actions"
	"github.com/forbole/callisto/v4/modules/types"

	"github.com/forbole/juno/v6/modules/pruning"
	"github.com/forbole/juno/v6/modules/telemetry"

	messagetype "github.com/forbole/callisto/v4/modules/message_type"
	"github.com/forbole/callisto/v4/modules/slashing"

	"github.com/cosmos/cosmos-sdk/codec"
	jmodules "github.com/forbole/juno/v6/modules"
	"github.com/forbole/juno/v6/modules/messages"
	"github.com/forbole/juno/v6/modules/registrar"
	juno "github.com/forbole/juno/v6/types"

	"github.com/forbole/callisto/v4/utils"

	"github.com/forbole/callisto/v4/database"
	"github.com/forbole/callisto/v4/modules/auth"
	"github.com/forbole/callisto/v4/modules/bank"
	"github.com/forbole/callisto/v4/modules/consensus"
	"github.com/forbole/callisto/v4/modules/distribution"
	"github.com/forbole/callisto/v4/modules/feegrant"

	"github.com/forbole/callisto/v4/modules/gov"
	"github.com/forbole/callisto/v4/modules/mint"
	"github.com/forbole/callisto/v4/modules/modules"
	"github.com/forbole/callisto/v4/modules/pricefeed"
	"github.com/forbole/callisto/v4/modules/staking"
)

// UniqueAddressesParser returns a wrapper around the given parser that removes all duplicated addresses
func UniqueAddressesParser(parser messages.MessageAddressesParser) messages.MessageAddressesParser {
	return func(tx *juno.Transaction) ([]string, error) {
		addresses, err := parser(tx)
		if err != nil {
			return nil, err
		}

		return utils.RemoveDuplicateValues(addresses), nil
	}
}

// --------------------------------------------------------------------------------------------------------------------

var (
	_ registrar.Registrar = &Registrar{}
)

// Registrar represents the modules.Registrar that allows to register all modules that are supported by BigDipper
type Registrar struct {
	parser messages.MessageAddressesParser
	cdc    codec.Codec
}

// NewRegistrar allows to build a new Registrar instance
func NewRegistrar(parser messages.MessageAddressesParser, cdc codec.Codec) *Registrar {
	return &Registrar{
		parser: UniqueAddressesParser(parser),
		cdc:    cdc,
	}
}

// BuildModules implements modules.Registrar
func (r *Registrar) BuildModules(ctx registrar.Context) jmodules.Modules {
	cdc := r.cdc.Marshaler
	db := database.Cast(ctx.Database)

	sources, err := types.BuildSources(ctx.JunoConfig.Node, r.cdc)
	if err != nil {
		panic(err)
	}

	actionsModule := actions.NewModule(ctx.JunoConfig, r.cdc, sources)
	authModule := auth.NewModule(r.parser, r.cdc, db)
	bankModule := bank.NewModule(r.parser, sources.BankSource, r.cdc, db)
	consensusModule := consensus.NewModule(db)
	distrModule := distribution.NewModule(sources.DistrSource, r.cdc, db)
	feegrantModule := feegrant.NewModule(r.cdc, db)
	messagetypeModule := messagetype.NewModule(r.parser, cdc, db)
	mintModule := mint.NewModule(sources.MintSource, r.cdc, db)
	slashingModule := slashing.NewModule(sources.SlashingSource, r.cdc, db)
	stakingModule := staking.NewModule(sources.StakingSource, slashingModule, r.cdc, db)
	govModule := gov.NewModule(sources.GovSource, authModule, distrModule, mintModule, slashingModule, stakingModule, r.cdc, db)

	return []jmodules.Module{
		messages.NewModule(r.parser, ctx.Database),
		telemetry.NewModule(ctx.JunoConfig),
		pruning.NewModule(ctx.JunoConfig, db, ctx.Logger),

		actionsModule,
		authModule,
		bankModule,
		consensusModule,
		distrModule,
		feegrantModule,
		govModule,
		mintModule,
		messagetypeModule,
		modules.NewModule(ctx.JunoConfig.Chain, db),
		pricefeed.NewModule(ctx.JunoConfig, r.cdc, db),
		slashingModule,
		stakingModule,
	}
}
