package erc20sdk

import (
	"github.com/mgintoki/erc20sdk/api/erc20client"
	"github.com/mgintoki/erc20sdk/chain/binance"
	"github.com/mgintoki/erc20sdk/chain/ethereum"
	"github.com/mgintoki/erc20sdk/chain/okex"
	"github.com/mgintoki/erc20sdk/errno"
	"github.com/mgintoki/multichain"
	"github.com/mgintoki/multichain/api/provider"
)

const (
	TypeETH  = multichain.TypeEthereum
	TypeBSC  = multichain.TypeBinance
	TypeOKEX = multichain.TypeOKEx
)

// NewClient 新建一个ERC20客户端
func NewClient(chainType uint, provider provider.CommonProvider) (erc20client.Client, error) {

	var cli erc20client.Client
	var err error

	switch chainType {
	case TypeETH:
		cli, err = ethereum.NewClient(provider)
	case TypeBSC:
		cli, err = binance.NewClient(provider)
	case TypeOKEX:
		cli, err = okex.NewClient(provider)
	}

	if err != nil {
		return nil, err
	}

	if cli == nil {
		return nil, errno.NotSupportChainType
	}

	return cli, nil
}
