package memwallet

import (
	"github.com/decred/dcrd/chaincfg"
	"github.com/decred/dcrd/dcrutil"
	"github.com/decred/dcrd/hdkeychain"
	"github.com/decred/dcrd/wire"
	"github.com/jfixby/coinharness"
	"github.com/jfixby/dcrharness"
	"github.com/jfixby/pin"
)

// WalletFactory produces a new InMemoryWallet-instance upon request
type WalletFactory struct {
}

// NewWallet creates and returns a fully initialized instance of the InMemoryWallet.
func (f *WalletFactory) NewWallet(cfg *coinharness.TestWalletConfig) coinharness.Wallet {
	pin.AssertNotNil("ActiveNet", cfg.ActiveNet)
	w, e := newMemWallet(cfg.ActiveNet, cfg.Seed)
	pin.CheckTestSetupMalfunction(e)
	return w
}

func newMemWallet(net coinharness.Network, harnessHDSeed coinharness.Seed) (*InMemoryWallet, error) {
	hdRoot, err := hdkeychain.NewMaster(harnessHDSeed.([]byte)[:], net.Params().(*chaincfg.Params))
	if err != nil {
		return nil, nil
	}

	// The first child key from the hd root is reserved as the coinbase
	// generation address.
	coinbaseChild, err := hdRoot.Child(0)
	if err != nil {
		return nil, err
	}
	coinbaseKey, err := coinbaseChild.ECPrivKey()
	if err != nil {
		return nil, err
	}
	coinbaseAddr, err := keyToAddr(coinbaseKey, net)
	if err != nil {
		return nil, err
	}

	// Track the coinbase generation address to ensure we properly track
	// newly generated coins we can spend.
	addrs := make(map[uint32]dcrutil.Address)
	addrs[0] = coinbaseAddr

	clientFac := &dcrharness.RPCClientFactory{}

	return &InMemoryWallet{
		net:               net,
		coinbaseKey:       coinbaseKey,
		coinbaseAddr:      coinbaseAddr,
		hdIndex:           1,
		hdRoot:            hdRoot,
		addrs:             addrs,
		utxos:             make(map[coinharness.OutPoint]*utxo),
		chainUpdateSignal: make(chan string),
		reorgJournal:      make(map[int64]*undoEntry),
		RPCClientFactory:  clientFac,
	}, nil
}
