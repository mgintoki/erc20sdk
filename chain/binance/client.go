package binance

import (
	"github.com/mgintoki/erc20sdk/chain/ethereum"
	"github.com/mgintoki/multichain/api/provider"
)

type Client = ethereum.Client

func NewClient(provider provider.CommonProvider) (*Client, error) {
	return ethereum.NewClient(provider)
}
