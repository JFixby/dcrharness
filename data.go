package dcrharness

import (
	"github.com/decred/dcrd/chaincfg"
	"github.com/decred/dcrd/dcrutil"
	"github.com/jfixby/coinharness"
)

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


