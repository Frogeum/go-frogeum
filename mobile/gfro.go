// Copyright 2016 The go-frogeum Authors
// This file is part of the go-frogeum library.
//
// The go-frogeum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-frogeum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-frogeum library. If not, see <http://www.gnu.org/licenses/>.

// Contains all the wrappers from the node package to support client side node
// management on mobile platforms.

package gfro

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/frogeum/go-frogeum/core"
	"github.com/frogeum/go-frogeum/eth/downloader"
	"github.com/frogeum/go-frogeum/eth/ethconfig"
	"github.com/frogeum/go-frogeum/ethclient"
	"github.com/frogeum/go-frogeum/ethstats"
	"github.com/frogeum/go-frogeum/internal/debug"
	"github.com/frogeum/go-frogeum/les"
	"github.com/frogeum/go-frogeum/node"
	"github.com/frogeum/go-frogeum/p2p"
	"github.com/frogeum/go-frogeum/p2p/nat"
	"github.com/frogeum/go-frogeum/params"
)

// NodeConfig represents the collection of configuration values to fine tune the Gfro
// node embedded into a mobile process. The available values are a subset of the
// entire API provided by go-frogeum to reduce the maintenance surface and dev
// complexity.
type NodeConfig struct {
	// Bootstrap nodes used to establish connectivity with the rest of the network.
	BootstrapNodes *Enodes

	// MaxPeers is the maximum number of peers that can be connected. If this is
	// set to zero, then only the configured static and trusted peers can connect.
	MaxPeers int

	// FrogeumEnabled specifies whether the node should run the Frogeum protocol.
	FrogeumEnabled bool

	// FrogeumNetworkID is the network identifier used by the Frogeum protocol to
	// decide if remote peers should be accepted or not.
	FrogeumNetworkID int64 // uint64 in truth, but Java can't handle that...

	// FrogeumGenesis is the genesis JSON to use to seed the blockchain with. An
	// empty genesis state is equivalent to using the mainnet's state.
	FrogeumGenesis string

	// FrogeumDatabaseCache is the system memory in MB to allocate for database caching.
	// A minimum of 16MB is always reserved.
	FrogeumDatabaseCache int

	// FrogeumNetStats is a netstats connection string to use to report various
	// chain, transaction and node stats to a monitoring server.
	//
	// It has the form "nodename:secret@host:port"
	FrogeumNetStats string

	// Listening address of pprof server.
	PprofAddress string
}

// defaultNodeConfig contains the default node configuration values to use if all
// or some fields are missing from the user's specified list.
var defaultNodeConfig = &NodeConfig{
	BootstrapNodes:       FoundationBootnodes(),
	MaxPeers:             25,
	FrogeumEnabled:       true,
	FrogeumNetworkID:     1213,
	FrogeumDatabaseCache: 16,
}

// NewNodeConfig creates a new node option set, initialized to the default values.
func NewNodeConfig() *NodeConfig {
	config := *defaultNodeConfig
	return &config
}

// AddBootstrapNode adds an additional bootstrap node to the node config.
func (conf *NodeConfig) AddBootstrapNode(node *Enode) {
	conf.BootstrapNodes.Append(node)
}

// EncodeJSON encodes a NodeConfig into a JSON data dump.
func (conf *NodeConfig) EncodeJSON() (string, error) {
	data, err := json.Marshal(conf)
	return string(data), err
}

// String returns a printable representation of the node config.
func (conf *NodeConfig) String() string {
	return encodeOrError(conf)
}

// Node represents a Gfro Frogeum node instance.
type Node struct {
	node *node.Node
}

// NewNode creates and configures a new Gfro node.
func NewNode(datadir string, config *NodeConfig) (stack *Node, _ error) {
	// If no or partial configurations were specified, use defaults
	if config == nil {
		config = NewNodeConfig()
	}
	if config.MaxPeers == 0 {
		config.MaxPeers = defaultNodeConfig.MaxPeers
	}
	if config.BootstrapNodes == nil || config.BootstrapNodes.Size() == 0 {
		config.BootstrapNodes = defaultNodeConfig.BootstrapNodes
	}

	if config.PprofAddress != "" {
		debug.StartPProf(config.PprofAddress, true)
	}

	// Create the empty networking stack
	nodeConf := &node.Config{
		Name:        clientIdentifier,
		Version:     params.VersionWithMeta,
		DataDir:     datadir,
		KeyStoreDir: filepath.Join(datadir, "keystore"), // Mobile should never use internal keystores!
		P2P: p2p.Config{
			NoDiscovery:      true,
			DiscoveryV5:      true,
			BootstrapNodesV5: config.BootstrapNodes.nodes,
			ListenAddr:       ":0",
			NAT:              nat.Any(),
			MaxPeers:         config.MaxPeers,
		},
	}

	rawStack, err := node.New(nodeConf)
	if err != nil {
		return nil, err
	}

	debug.Memsize.Add("node", rawStack)

	var genesis *core.Genesis
	if config.FrogeumGenesis != "" {
		// Parse the user supplied genesis spec if not mainnet
		genesis = new(core.Genesis)
		if err := json.Unmarshal([]byte(config.FrogeumGenesis), genesis); err != nil {
			return nil, fmt.Errorf("invalid genesis spec: %v", err)
		}
		// If we have the Longcat testnet, hard code the chain configs too
		if config.FrogeumGenesis == LongcatGenesis() {
			genesis.Config = params.LongcatChainConfig
			if config.FrogeumNetworkID == 1213 {
				config.FrogeumNetworkID = 777
			}
		}
		// If we have the Rinkeby testnet, hard code the chain configs too
		if config.FrogeumGenesis == RinkebyGenesis() {
			genesis.Config = params.RinkebyChainConfig
			if config.FrogeumNetworkID == 1213 {
				config.FrogeumNetworkID = 4
			}
		}
		// If we have the Goerli testnet, hard code the chain configs too
		if config.FrogeumGenesis == GoerliGenesis() {
			genesis.Config = params.GoerliChainConfig
			if config.FrogeumNetworkID == 1213 {
				config.FrogeumNetworkID = 5
			}
		}
	}
	// Register the Frogeum protocol if requested
	if config.FrogeumEnabled {
		ethConf := ethconfig.Defaults
		ethConf.Genesis = genesis
		ethConf.SyncMode = downloader.LightSync
		ethConf.NetworkId = uint64(config.FrogeumNetworkID)
		ethConf.DatabaseCache = config.FrogeumDatabaseCache
		lesBackend, err := les.New(rawStack, &ethConf)
		if err != nil {
			return nil, fmt.Errorf("frogeum init: %v", err)
		}
		// If netstats reporting is requested, do it
		if config.FrogeumNetStats != "" {
			if err := ethstats.New(rawStack, lesBackend.ApiBackend, lesBackend.Engine(), config.FrogeumNetStats); err != nil {
				return nil, fmt.Errorf("netstats init: %v", err)
			}
		}
	}
	return &Node{rawStack}, nil
}

// Close terminates a running node along with all it's services, tearing internal state
// down. It is not possible to restart a closed node.
func (n *Node) Close() error {
	return n.node.Close()
}

// Start creates a live P2P node and starts running it.
func (n *Node) Start() error {
	// TODO: recreate the node so it can be started multiple times
	return n.node.Start()
}

// Stop terminates a running node along with all its services. If the node was not started,
// an error is returned. It is not possible to restart a stopped node.
//
// Deprecated: use Close()
func (n *Node) Stop() error {
	return n.node.Close()
}

// GetFrogeumClient retrieves a client to access the Frogeum subsystem.
func (n *Node) GetFrogeumClient() (client *FrogeumClient, _ error) {
	rpc, err := n.node.Attach()
	if err != nil {
		return nil, err
	}
	return &FrogeumClient{ethclient.NewClient(rpc)}, nil
}

// GetNodeInfo gathers and returns a collection of metadata known about the host.
func (n *Node) GetNodeInfo() *NodeInfo {
	return &NodeInfo{n.node.Server().NodeInfo()}
}

// GetPeersInfo returns an array of metadata objects describing connected peers.
func (n *Node) GetPeersInfo() *PeerInfos {
	return &PeerInfos{n.node.Server().PeersInfo()}
}
