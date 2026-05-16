package filter

import (
	"fmt"
	"net/netip"

	"github.com/google/nftables"
)

// NFTCtx holds the nftables connection and all QOSM table structures.
type NFTCtx struct {
	conn *nftables.Conn
	QosmTable
}

// QosmTable represents the nftables table with all QOSM chains and sets.
type QosmTable struct {
	*nftables.Table
	QosmChains
	QosmSets
}

// QosmChain represents an nftables chain with its associated rules.
type QosmChain struct {
	*nftables.Chain
	QosmRules
}

// QosmChains holds the output and forward chains for QOSM.
type QosmChains struct {
	outputChain  QosmChain
	forwardChain QosmChain
}

// QosmSets holds the nftables ip sets for high and low priority traffic.
type QosmSets struct {
	highPrioSet *nftables.Set
	lowPrioSet  *nftables.Set
}

// QosmRules holds the nftables rules for high and low priority traffic.
type QosmRules struct {
	highPrioRule *nftables.Rule
	lowPrioRule  *nftables.Rule
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

// AddTargetsToHighPriority ip addresses to the high-priority IP set.
func (c *NFTCtx) AddTargetsToHighPriority(targets []netip.Prefix) error {
	return addIPsToQoSMIPSet(c.conn, c.highPrioSet, targets)
}

// AddTargetsToLowPriority adds ip addresses to the low-priority IP set.
func (c *NFTCtx) AddTargetsToLowPriority(targets []netip.Prefix) error {
	return addIPsToQoSMIPSet(c.conn, c.lowPrioSet, targets)
}

// DeleteTargetFromHighPriority removes the given ip addresses from the high-priority IP set.
func (c *NFTCtx) DeleteTargetFromHighPriority(targets []netip.Prefix) error {
	return deleteIPsFromQoSIPSet(c.conn, c.highPrioSet, targets)
}

// DeleteTargetFromLowPriority removes the given ip addresses from the low-priority IP set.
func (c *NFTCtx) DeleteTargetFromLowPriority(targets []netip.Prefix) error {
	return deleteIPsFromQoSIPSet(c.conn, c.lowPrioSet, targets)
}

// GetHighPrioIPs returns all IP addresses in the high-priority set.
func (c *NFTCtx) GetHighPrioIPs() ([]netip.Addr, error) {
	return getIPSetElements(c.conn, c.highPrioSet)
}

// GetLowPrioIPs returns all IP addresses in the low-priority set.
func (c *NFTCtx) GetLowPrioIPs() ([]netip.Addr, error) {
	return getIPSetElements(c.conn, c.lowPrioSet)
}

// DeleteTable removes the qosm nftables table from the system. The context becomes invalid after this operation.
func (c *NFTCtx) DeleteTable() error {
	fmt.Println("Deleting table")
	c.conn.DelTable(c.Table)
	return c.conn.Flush()
}
