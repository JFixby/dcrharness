package dcrharness

import (
	"fmt"
	"github.com/decred/dcrd/chaincfg"
	"github.com/decred/dcrd/chaincfg/chainhash"
	"github.com/decred/dcrd/dcrjson"
	"github.com/decred/dcrd/dcrutil"
	"github.com/decred/dcrd/rpcclient"
	"github.com/jfixby/coinharness"
	"github.com/jfixby/pin"
	"io/ioutil"
)

type RPCClientFactory struct {
}

func (f *RPCClientFactory) NewRPCConnection(config coinharness.RPCConnectionConfig, handlers coinharness.RPCClientNotificationHandlers) (coinharness.RPCClient, error) {
	var h *rpcclient.NotificationHandlers
	if handlers != nil {
		h = handlers.(*rpcclient.NotificationHandlers)
	}

	file := config.CertificateFile
	fmt.Println("reading: " + file)
	cert, err := ioutil.ReadFile(file)
	pin.CheckTestSetupMalfunction(err)

	cfg := &rpcclient.ConnConfig{
		Host:                 config.Host,
		Endpoint:             config.Endpoint,
		User:                 config.User,
		Pass:                 config.Pass,
		Certificates:         cert,
		DisableAutoReconnect: true,
		HTTPPostMode:         false,
	}

	return NewRPCClient(cfg, h)
}

type RPCClient struct {
	rpc *rpcclient.Client
}

func (c *RPCClient) ListUnspent() ([]*coinharness.Unspent, error) {
	result, err := c.rpc.ListUnspent()
	if err != nil {
		return nil, err
	}
	var r []*coinharness.Unspent
	for _, e := range result {
		x := &coinharness.Unspent{}
		x.Address = e.Address
		x.Spendable = e.Spendable
		x.TxType = e.TxType
		x.TxID = e.TxID
		x.Confirmations = e.Confirmations
		r = append(r, x)
	}
	return r, nil
}

func (c *RPCClient) AddNode(args *coinharness.AddNodeArguments) error {
	return c.rpc.AddNode(args.TargetAddr, args.Command.(rpcclient.AddNodeCommand))
}

func (c *RPCClient) LoadTxFilter(reload bool, addr []coinharness.Address) error {
	addresses := []dcrutil.Address{}
	for _, e := range addr {
		addresses = append(addresses, e.Internal().(dcrutil.Address))
	}
	return c.rpc.LoadTxFilter(reload, addresses, nil)
}

func (c *RPCClient) SubmitBlock(block coinharness.Block) error {
	return c.rpc.SubmitBlock(block.(*dcrutil.Block), nil)
}

func (c *RPCClient) Disconnect() {
	c.rpc.Disconnect()
}

func (c *RPCClient) Shutdown() {
	c.rpc.Shutdown()
}

func (c *RPCClient) NotifyBlocks() error {
	return c.rpc.NotifyBlocks()
}

func (c *RPCClient) GetBlockCount() (int64, error) {
	return c.rpc.GetBlockCount()
}

func (c *RPCClient) Generate(blocks uint32) (result []coinharness.Hash, e error) {
	list, e := c.rpc.Generate(blocks)
	if e != nil {
		return nil, e
	}
	for _, el := range list {
		result = append(result, el)
	}
	return result, nil
}

func (c *RPCClient) Internal() interface{} {
	return c.rpc
}

func (c *RPCClient) GetRawMempool(command interface{}) (result []coinharness.Hash, e error) {
	list, e := c.rpc.GetRawMempool(command.(dcrjson.GetRawMempoolTxTypeCmd))
	if e != nil {
		return nil, e
	}
	for _, el := range list {
		result = append(result, el)
	}
	return result, nil
}

func (c *RPCClient) SendRawTransaction(tx coinharness.CreatedTransactionTx, allowHighFees bool) (result coinharness.Hash, e error) {
	txx := TransactionTxToRaw(tx)
	r, e := c.rpc.SendRawTransaction(txx, allowHighFees)
	return r, e
}

func (c *RPCClient) GetBlock(hash coinharness.Hash) (coinharness.Block, error) {
	return c.rpc.GetBlock(hash.(*chainhash.Hash))
}

func (c *RPCClient) GetPeerInfo() ([]coinharness.PeerInfo, error) {
	pif, err := c.rpc.GetPeerInfo()
	if err != nil {
		return nil, err
	}

	l := []coinharness.PeerInfo{}
	for _, i := range pif {
		inf := coinharness.PeerInfo{}
		inf.Addr = i.Addr
		l = append(l, inf)

	}
	return l, nil
}

func NewRPCClient(config *rpcclient.ConnConfig, handlers *rpcclient.NotificationHandlers) (coinharness.RPCClient, error) {
	legacy, err := rpcclient.New(config, handlers)
	if err != nil {
		return nil, err
	}

	result := &RPCClient{rpc: legacy}
	return result, nil
}

func (c *RPCClient) GetNewAddress(account string) (coinharness.Address, error) {
	legacy, err := c.rpc.GetNewAddress(account)
	if err != nil {
		return nil, err
	}

	result := &Address{Address: legacy}
	return result, nil
}

func (c *RPCClient) ValidateAddress(address coinharness.Address) (*coinharness.ValidateAddressResult, error) {
	legacy, err := c.rpc.ValidateAddress(address.Internal().(dcrutil.Address))
	// *dcrjson.ValidateAddressWalletResult
	if err != nil {
		return nil, err
	}
	result := &coinharness.ValidateAddressResult{
		Address:      legacy.Address,
		Account:      legacy.Account,
		IsValid:      legacy.IsValid,
		IsMine:       legacy.IsMine,
		IsCompressed: legacy.IsCompressed,
	}
	return result, nil
}

func (c *RPCClient) GetBalance(account string) (*coinharness.GetBalanceResult, error) {
	legacy, err := c.rpc.GetBalance(account)
	// *dcrjson.ValidateAddressWalletResult
	if err != nil {
		return nil, err
	}
	result := &coinharness.GetBalanceResult{
		BlockHash:      legacy.BlockHash,
		TotalSpendable: coinharness.CoinsAmountFromFloat(legacy.TotalSpendable),
	}
	return result, nil
}

func (c *RPCClient) GetBestBlock() (coinharness.Hash, int64, error) {
	return c.rpc.GetBestBlock()
}

func (c *RPCClient) CreateNewAccount(account string) error {
	return c.rpc.CreateNewAccount(account)
}

func (c *RPCClient) WalletLock() error {
	return c.rpc.WalletLock()
}

func (c *RPCClient) WalletInfo() (*coinharness.WalletInfoResult, error) {
	r, err := c.rpc.WalletInfo()
	if err != nil {
		return nil, err
	}
	result := &coinharness.WalletInfoResult{
		Unlocked:        r.Unlocked,
		DaemonConnected: r.DaemonConnected,
		Voting:          r.DaemonConnected,
	}
	return result, nil
}

func (c *RPCClient) WalletUnlock(passphrase string, timeoutSecs int64) error {
	return c.rpc.WalletPassphrase(passphrase, timeoutSecs)
}

func (c *RPCClient) GetBuildVersion() (coinharness.BuildVersion, error) {
	//legacy, err := c.rpc.GetBuildVersion()
	//if err != nil {
	//	return nil, err
	//}
	//return legacy, nil
	return nil, fmt.Errorf("decred does not support this feature (GetBuildVersion)")
}

type Address struct {
	Address dcrutil.Address
}

func (c *Address) String() string {
	return c.Address.String()
}

func (c *Address) Internal() interface{} {
	return c.Address
}

func (c *Address) IsForNet(net coinharness.Network) bool {
	return c.Address.IsForNet(net.Params().(*chaincfg.Params))
}
