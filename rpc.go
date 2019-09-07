package dcrharness

import (
	"fmt"
	"github.com/decred/dcrd/dcrjson"
	"github.com/decred/dcrd/rpcclient"
	"github.com/jfixby/coinharness"
	"github.com/jfixby/pin"
	"io/ioutil"
)

type DcrRPCClientFactory struct {
}

func (f *DcrRPCClientFactory) NewRPCConnection(config coinharness.RPCConnectionConfig, handlers coinharness.RPCClientNotificationHandlers) (coinharness.RPCClient, error) {
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

type DCRPCClient struct {
	rpc *rpcclient.Client
}

func (c *DCRPCClient) AddNode(args *coinharness.AddNodeArguments) error {
	return c.rpc.AddNode(args.TargetAddr, args.Command.(rpcclient.AddNodeCommand))
}

func (c *DCRPCClient) Disconnect() {
	c.rpc.Disconnect()
}

func (c *DCRPCClient) Shutdown() {
	c.rpc.Shutdown()
}

func (c *DCRPCClient) NotifyBlocks() error {
	return c.rpc.NotifyBlocks()
}

func (c *DCRPCClient) GetBlockCount() (int64, error) {
	return c.rpc.GetBlockCount()
}

func (c *DCRPCClient) Generate(blocks uint32) (result []coinharness.Hash, e error) {
	list, e := c.rpc.Generate(blocks)
	if e != nil {
		return nil, e
	}
	for _, el := range list {
		result = append(result, el)
	}
	return result, nil
}

func (c *DCRPCClient) Internal() interface{} {
	return c.rpc
}

func (c *DCRPCClient) GetRawMempool(command interface{}) (result []coinharness.Hash, e error) {
	list, e := c.rpc.GetRawMempool(command.(dcrjson.GetRawMempoolTxTypeCmd))
	if e != nil {
		return nil, e
	}
	for _, el := range list {
		result = append(result, el)
	}
	return result, nil
}

func (c *DCRPCClient) SendRawTransaction(tx coinharness.CreatedTransactionTx, allowHighFees bool) (result coinharness.Hash, e error) {
	txx := TransactionTxToRaw(tx)
	r, e := c.rpc.SendRawTransaction(txx, allowHighFees)
	return r, e
}

func (c *DCRPCClient) GetPeerInfo() ([]coinharness.PeerInfo, error) {
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

	result := &DCRPCClient{rpc: legacy}
	return result, nil
}

func (c *DCRPCClient) GetNewAddress(account string) (coinharness.Address, error) {
	return c.rpc.GetNewAddress(account)
}
