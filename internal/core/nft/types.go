package nft

import (
	"errors"
	"log/slog"

	"github.com/google/nftables"
)

type NFTOpts struct {
	CreateIfNotExists bool
	Logger            *slog.Logger
}

type IfaceIndex int

// NFTCtx holds the nftables connection and all QOSM table structures.
type NFTCtx struct {
	conn *nftables.Conn
	qosmTable

	Logger *slog.Logger
}

// qosmTable represents the nftables table with all QOSM chains and sets.
type qosmTable struct {
	*nftables.Table
	qosmChains
	qosmSets
}

// qosmChains holds the output and forward chains for QOSM.
type qosmChains struct {
	outputChain  qosmChain
	forwardChain qosmChain
}

// qosmChain represents an nftables chain with its associated rules.
type qosmChain struct {
	*nftables.Chain
	Rules map[IfaceIndex]qosmRules
}

// qosmSets holds the nftables ip sets for high and low priority traffic.
type qosmSets struct {
	highPrioSet *nftables.Set
	lowPrioSet  *nftables.Set
}

// qosmRules holds the nftables rules for high and low priority traffic.
type qosmRules struct {
	highPrioRule *nftables.Rule
	lowPrioRule  *nftables.Rule
}

type RuleStats struct {
	PacketCount uint64
	ByteCount   uint64
}

type InterfaceStats struct {
	IfIndex  int
	HighPrio RuleStats
	LowPrio  RuleStats
}

const (
	TABLENAME        = "qosmtable"
	OUTPUTCHAINNAME  = "output"
	FORWARDCHAINNAME = "forward"
)

const (
	HIGHPRIORULENAME  = "high_prio_rule"
	HIGHPRIOIPSETNAME = "high_prio_ips"
	HIGHPRIOMARK      = 10
)

const (
	LOWPRIORULENAME  = "low_prio_rule"
	LOWPRIOIPSETNAME = "low_prio_ips"
	LOWPRIOMARK      = 20
)

type ruleParams struct {
	table       *nftables.Table
	chain       *nftables.Chain
	ipSet       *nftables.Set
	oifaceIndex int
	mark        int
	ruleName    string
}

var (
	ErrNotFound      = errors.New("nft object not found")
	ErrTableNotFound = errors.New("qosm table not found")
	ErrChainNotFound = errors.New("qosm chains not found")
)

type ErrSetNotFound struct {
	Name string
}

func (e ErrSetNotFound) Error() string {
	return "nft set " + e.Name + " not found"
}

type ErrRuleNotFound struct {
	Name string
}

func (e ErrRuleNotFound) Error() string {
	return "nft chain " + e.Name + " not found"
}
