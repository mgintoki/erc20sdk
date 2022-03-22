package erc20client

import (
	"github.com/mgintoki/multichain/api/provider"
	"github.com/mgintoki/multichain/api/tx"
	"math/big"
)

// CrossChainOption 跨链参数
// 当前只支持 ethereum 跨 weelink
// todo 支持 weelink 跨 ethereum
type CrossChainOption struct {
	GravityContract      string `json:"gravity_contract"`       // 跨链合约地址
	DestinationChainType string `json:"destination_chain_type"` // 目标链的类型
	DestinationChainId   string `json:"destination_chain_id"`   //目标链的id
	//DestinationAccount   string `json:"destination_account"` // 目标链接收代币的账户
}

// DeployERC20Param ERC20合约部署参数
// 如果不传 Abi 或是 ByteCode, 将会使用系统内置的ERC20 Abi 和 ByteCode
// Params 必须和 Abi 中合约的构造函数参数完全一致
// 系统内置的ERC20合约构造函数包括五个参数 名称(string) 符号(string) 精度(uint8) 初始发行量(*big.int),是否可铸币（bool）
type DeployERC20Param struct {
	Abi      string        //合约abi
	ByteCode string        //合约字节码
	Params   []interface{} //合约初始化参数
}

// Client 定义了通用的ERC20接口
//
// Client 分为四个部分
// 第一部分是 Client 的设置接口，用来管理 Client
// 第二部分是 ERC20 的标准接口定义, 见 https://eips.ethereum.org/EIPS/eip-20
// 第三部分是对于ERC20的扩展接口,用于增强ERC20代币的能力，本SDK的内置ERC20合约已支持对应的扩展接口
// 第四部分是部署合约接口，会构建一个部署ERC20合约的交易
//
// 注意：下文中交易的from账户 指的是，发起交易的账户地址。
// 修改链上状态是通过统一的发送交易执行的，交易构建好之后，交易的发起者会使用自己的私钥对账户签名（私钥推出公钥，公钥推出账户，所以私钥和交易的from字段必须能够对应）
// 签名完成之后，会对交易进行发送
//
// 注意：本SDK对于调用合约（修改链上状态）的接口，返回值是构建好的交易（返回值为 tx.Tx ) 需要实例化一个 multichain client 实例，用来发送交易
// 发送交易相关使用请参考 multichain sdk或是本sdk的 erc20_test.go 文件
//
// 如果本地持有私钥，则实例化multichain client 之后，对 multichain client 设置私钥，并调用sendTx接口发送交易，sendTx接口会在内部使用设置的私钥签名交易并发送
//
// 如果本地无私钥，实例化构建好交易之后，可导出交易需要签名的hash，在别处使用相同的交易算法签名之后，将签名注入到交易中，生成签名后的交易，
// 然后实例化 multichain client 使用 sendSignedTx接口发送交易
type Client interface {

	// 第一部分 设置Client
	//
	// SetProvider 设置链服务提供商
	SetProvider(provider provider.CommonProvider) error //

	// SetContractAddress 设置操作的合约地址
	SetContractAddress(addr string)

	// GetAccount 获取设置的账户
	GetAccount() string

	// SetAccount 设置账户（发送交易前必须设置过账户，该账户作为交易发起方。查询合约则无需设置账户）
	SetAccount(account string)

	// 第二部分 ERC20标准接口  https://eips.ethereum.org/EIPS/eip-20
	//
	// Name 查询代币名称
	Name() (string, error)

	// Symbol 查询代币符号
	Symbol() (string, error)

	// 查询某个地址的代币余额
	BalanceOf(addr string) (*big.Int, error)

	// Transfer 从 SetAccount 设置的账户转账 value 对应的数额到 to 对应的账户
	// SetAccount 之后，需要调用multichain client 发送交易接口
	// 可选跨链参数, 返回构建好的交易
	Transfer(to string, value *big.Int, option *CrossChainOption) (tx tx.Tx, err error)

	// TransferFrom 从授权地址from转出value对应的代币到to地址
	// 注意: 该from和交易的from不能是同一个账户，此处的from是，对交易的from设置过代币限额（allowance）的账户
	// from账户对于交易from账户的 allowance必须大于value值
	TransferFrom(from, to string, value *big.Int) (tx tx.Tx, err error)

	// Approve 交易的from账户对于某个spender账户设置spender可以花费的代币金额
	Approve(spender string, value *big.Int) (tx tx.Tx, err error)

	// Allowance 查询spender对于owner的可用金额
	Allowance(owner string, spender string) (*big.Int, error)

	// 第三部分 ERC20扩展接口
	// 注意：使用扩展接口的时候，请确认部署的合约中支持这些扩展接口（SDK内置合约已支持）
	//
	// Decimals 查询代币精度
	Decimals() (uint8, error)

	// TotalSupply 查询代币总发行量
	TotalSupply() (*big.Int, error)

	// Mint 给address地址新铸造amount数量的代币，需要管理员权限
	Mint(address string, amount *big.Int) (tx tx.Tx, err error)

	// Burn 销毁指定地址指定数量的代币，要求参数地址是交易发起者，或是已经授予交易发起者大于amount数量的allowance
	Burn(address string, amount *big.Int) (tx tx.Tx, err error)

	// 第四部分 部署合约接口
	//
	// 见 DeployERC20Param 定义
	// 部署完成合约后，可调用 multichain client 的 获取交易详情接口，拿到部署合约生成的合约地址
	// 部署合约的人会被视为管理员，拥有发币权限
	DeployContract(param DeployERC20Param) (tx tx.Tx, err error)
}
