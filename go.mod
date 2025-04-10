module github.com/multiversx/mx-chain-scenario-go

go 1.20

replace (
	github.com/multiversx/mx-chain-core-go => github.com/multiversx/mx-chain-core-sovereign-go v1.2.25-0.20250410112225-9b4402144b11
	github.com/multiversx/mx-chain-vm-common-go => github.com/multiversx/mx-chain-vm-common-sovereign-go v1.5.17-0.20250410143856-c8c2b5eeaa7c
)

require (
	github.com/TwiN/go-color v1.1.0
	github.com/multiversx/mx-chain-core-go v1.2.24-0.20241119082458-e2451e147ab1
	github.com/multiversx/mx-chain-crypto-go v1.2.13-0.20250410124744-d21d37be8e32
	github.com/multiversx/mx-chain-logger-go v1.0.15
	github.com/multiversx/mx-chain-vm-common-go v1.5.13
	github.com/multiversx/mx-components-big-int v1.0.0
	github.com/stretchr/testify v1.8.4
	github.com/urfave/cli/v2 v2.27.1
	golang.org/x/crypto v0.17.0
)

require (
	filippo.io/edwards25519 v1.0.0 // indirect
	github.com/btcsuite/btcd/btcec/v2 v2.2.0 // indirect
	github.com/btcsuite/btcd/btcutil v1.1.3 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.0.1 // indirect
	github.com/denisbrodbeck/machineid v1.0.1 // indirect
	github.com/ethereum/go-ethereum v1.13.15 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/holiman/uint256 v1.2.4 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/mr-tron/base58 v1.2.0 // indirect
	github.com/multiversx/mx-sdk-abi-go v0.3.1-0.20240912062928-8502f4c3b37c // indirect
	github.com/pelletier/go-toml v1.9.3 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/xrash/smetrics v0.0.0-20201216005158-039620a65673 // indirect
	golang.org/x/sys v0.16.0 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
