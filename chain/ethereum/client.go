package ethereum

import (
	"github.com/mgintoki/erc20sdk/api/erc20client"
	"github.com/mgintoki/erc20sdk/artifacts"
	"github.com/mgintoki/erc20sdk/errno"
	"github.com/mgintoki/erc20sdk/tools"
	"github.com/mgintoki/multichain"
	mcCli "github.com/mgintoki/multichain/api/client"
	"github.com/mgintoki/multichain/api/provider"
	"github.com/mgintoki/multichain/api/tx"
	"github.com/mgintoki/multichain/api/txbuilder"
	"github.com/mgintoki/multichain/chain/ethereum"
	"math/big"
)

type Client struct {
	mc              mcCli.Client //multiChain erc20client
	defaultAccount  string
	contractAddress string
	abi             string
	provider        provider.CommonProvider
	ctb             *ethereum.ContractTxBuilder
	tb              *ethereum.TxBuilder
	ae              *ethereum.AddressEncoder
}

func NewClient(provider provider.CommonProvider) (*Client, error) {
	cli := &Client{}
	mCli, err := multichain.NewClient(multichain.TypeEthereum, provider)
	if err != nil {
		return nil, err
	}
	cli.mc = mCli
	cli.abi = artifacts.DefaultABI
	cli.provider = provider

	cli.ctb, err = ethereum.NewContractTxBuilder(provider)
	if err != nil {
		return nil, err
	}

	cli.tb, err = ethereum.NewTxBuilder(provider)
	if err != nil {
		return nil, err
	}

	if err := cli.SetProvider(provider); err != nil {
		return nil, err
	} else {
		return cli, nil
	}
}

func (c *Client) SetProvider(provider provider.CommonProvider) error {
	return c.mc.SetProvider(provider)
}

func (c *Client) SetPrivate(privateHex string) error {
	return c.mc.SetPrivate(privateHex)
}

func (c *Client) GetAccount() string {
	return c.defaultAccount
}
func (c *Client) SetAccount(account string) {
	c.defaultAccount = account
}

func (c *Client) SetContractAddress(addr string) {
	c.contractAddress = addr
}

func (c *Client) Name() (string, error) {
	queryRes, err := c.mc.QueryContract(mcCli.CallContractParam{
		From:            "",
		ContractAddress: c.contractAddress,
		Abi:             c.abi,
		CalledFunc:      artifacts.FuncName,
		Params:          nil,
	})
	if err != nil {
		return "", err
	}

	if res, ok := queryRes.DecodeRes["0"].(string); !ok {
		return "", errno.InvalidTypeAssert.Add("raw queryRes is : " + tools.FastMarshal(queryRes))
	} else {
		return res, nil
	}
}

func (c *Client) Symbol() (string, error) {
	queryRes, err := c.mc.QueryContract(mcCli.CallContractParam{
		ContractAddress: c.contractAddress,
		Abi:             artifacts.DefaultABI,
		CalledFunc:      artifacts.FuncSymbol,
		Params:          nil,
	})
	if err != nil {
		return "", err
	}

	if res, ok := queryRes.DecodeRes["0"].(string); !ok {
		return "", errno.InvalidTypeAssert.Add("raw queryRes is : " + tools.FastMarshal(queryRes))
	} else {
		return res, nil
	}
}

func (c *Client) Decimals() (uint8, error) {
	queryRes, err := c.mc.QueryContract(mcCli.CallContractParam{
		ContractAddress: c.contractAddress,
		Abi:             artifacts.DefaultABI,
		CalledFunc:      artifacts.FuncDecimals,
		Params:          nil,
	})
	if err != nil {
		return 0, err
	}

	if res, ok := queryRes.DecodeRes["0"].(uint8); !ok {
		return 0, errno.InvalidTypeAssert.Add("raw queryRes is : " + tools.FastMarshal(queryRes))
	} else {
		return res, nil
	}
}

func (c *Client) TotalSupply() (*big.Int, error) {
	queryRes, err := c.mc.QueryContract(mcCli.CallContractParam{
		ContractAddress: c.contractAddress,
		Abi:             artifacts.DefaultABI,
		CalledFunc:      artifacts.FuncTotalSupply,
		Params:          nil,
	})
	if err != nil {
		return nil, err
	}

	if res, ok := queryRes.DecodeRes["0"].(*big.Int); !ok {
		return nil, errno.InvalidTypeAssert.Add("raw queryRes is : " + tools.FastMarshal(queryRes))
	} else {
		return res, nil
	}
}

func (c *Client) BalanceOf(addr string) (*big.Int, error) {
	queryRes, err := c.mc.QueryContract(mcCli.CallContractParam{
		ContractAddress: c.contractAddress,
		Abi:             artifacts.DefaultABI,
		CalledFunc:      artifacts.FuncBalanceOf,
		Params:          []interface{}{c.ae.HexToAddress(addr)},
	})
	if err != nil {
		return nil, err
	}

	if res, ok := queryRes.DecodeRes["0"].(*big.Int); !ok {
		return nil, errno.InvalidTypeAssert.Add("raw queryRes is : " + tools.FastMarshal(queryRes))
	} else {
		return res, nil
	}
}

func (c *Client) Allowance(owner string, spender string) (*big.Int, error) {

	queryRes, err := c.mc.QueryContract(mcCli.CallContractParam{
		ContractAddress: c.contractAddress,
		Abi:             artifacts.DefaultABI,
		CalledFunc:      artifacts.FuncAllowance,
		Params:          []interface{}{c.ae.HexToAddress(owner), c.ae.HexToAddress(spender)},
	})
	if err != nil {
		return nil, err
	}

	if res, ok := queryRes.DecodeRes["0"].(*big.Int); !ok {
		return nil, errno.InvalidTypeAssert.Add("raw queryRes is : " + tools.FastMarshal(queryRes))
	} else {
		return res, nil
	}
}

func (c *Client) DeployContract(param erc20client.DeployERC20Param) (tx tx.Tx, err error) {

	abi := param.Abi
	byteCode := param.ByteCode

	if abi == "" || byteCode == "" {
		abi = artifacts.DefaultABI
		byteCode = artifacts.DefaultByteCode
	}

	return c.ctb.BuildDeployTx(txbuilder.BuildDeployTxReq{
		From:     c.GetAccount(),
		Abi:      abi,
		ByteCode: byteCode,
		Params:   param.Params,
	})

}

//todo 支持跨链项
func (c *Client) Transfer(to string, value *big.Int, option *erc20client.CrossChainOption) (tx tx.Tx, err error) {
	return c.ctb.BuildInvokeTx(txbuilder.BuildInvokeTxReq{
		From:            c.GetAccount(),
		Abi:             c.abi,
		Method:          artifacts.FuncTransfer,
		ContractAddress: c.contractAddress,
		Params:          []interface{}{c.ae.HexToAddress(to), value},
	})
}

func (c *Client) TransferFrom(from, to string, value *big.Int) (tx tx.Tx, err error) {

	return c.ctb.BuildInvokeTx(txbuilder.BuildInvokeTxReq{
		From:            c.GetAccount(),
		Abi:             c.abi,
		Method:          artifacts.FuncTransferFrom,
		ContractAddress: c.contractAddress,
		Params:          []interface{}{c.ae.HexToAddress(from), c.ae.HexToAddress(to), value},
	})
}

func (c *Client) Approve(spender string, value *big.Int) (tx tx.Tx, err error) {

	return c.ctb.BuildInvokeTx(txbuilder.BuildInvokeTxReq{
		From:            c.GetAccount(),
		Abi:             c.abi,
		Method:          artifacts.FuncApprove,
		ContractAddress: c.contractAddress,
		Params:          []interface{}{c.ae.HexToAddress(spender), value},
	})
}

func (c *Client) Mint(address string, amount *big.Int) (tx tx.Tx, err error) {
	return c.ctb.BuildInvokeTx(txbuilder.BuildInvokeTxReq{
		From:            c.GetAccount(),
		Abi:             c.abi,
		Method:          artifacts.FuncMint,
		ContractAddress: c.contractAddress,
		Params:          []interface{}{c.ae.HexToAddress(address), amount},
	})
}

func (c *Client) Burn(address string, amount *big.Int) (tx tx.Tx, err error) {
	return c.ctb.BuildInvokeTx(txbuilder.BuildInvokeTxReq{
		From:            c.GetAccount(),
		Abi:             c.abi,
		Method:          artifacts.FuncBurn,
		ContractAddress: c.contractAddress,
		Params:          []interface{}{c.ae.HexToAddress(address), amount},
	})
}
