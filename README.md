Decred harness
=======
[![ISC License](http://img.shields.io/badge/license-ISC-blue.svg)](http://copyfree.org)

 - [memwallet](https://github.com/jfixby/dcrharness/tree/master/memwallet)
 Offers a simple in-memory HD wallet capable of properly syncing to the
 generated chain, creating new addresses, and crafting fully signed transactions
 paying to an arbitrary set of outputs.

 - [nodecls](https://github.com/jfixby/dcrharness/tree/master/nodecls)
 Provides wrapper that launches a new `dcrd`-instance using command-line call.

 ## Build
 ```
 set GO111MODULE=on
 go build ./...
 go clean -testcache
 go test ./...
```
 ## License
 This code is licensed under the [copyfree](http://copyfree.org) ISC License.