package promwrapper

import (
	"math/big"

	ocr3types "github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
)

var _ ocr3types.MercuryPluginFactory = &promFactory{}

type promFactory struct {
	wrapped   ocr3types.MercuryPluginFactory
	name      string
	chainType string
	chainID   *big.Int
}

func (p *promFactory) NewMercuryPlugin(config ocr3types.MercuryPluginConfig) (ocr3types.MercuryPlugin, ocr3types.MercuryPluginInfo, error) {
	plugin, info, err := p.wrapped.NewMercuryPlugin(config)
	if err != nil {
		return nil, ocr3types.MercuryPluginInfo{}, err
	}

	prom := New(plugin, p.name, p.chainType, p.chainID, config, nil)
	return prom, info, nil
}

func NewPromFactory(wrapped ocr3types.MercuryPluginFactory, name, chainType string, chainID *big.Int) ocr3types.MercuryPluginFactory {
	return &promFactory{
		wrapped:   wrapped,
		name:      name,
		chainType: chainType,
		chainID:   chainID,
	}
}
