// Package qos
package qos

import (
	"fmt"
	"log/slog"
	"net/netip"

	"github.com/florianl/go-tc"
	"github.com/kakeetopius/qosm/internal/core/htb"
	"github.com/kakeetopius/qosm/internal/core/nft"
	"github.com/kakeetopius/qosm/internal/prio"
)

type QoSManager struct {
	TcConn     *tc.Tc
	HTBIfaces  map[int]htb.HTBIface
	Classifier *nft.NFT

	Logger *slog.Logger
}

func NewManager() (*QoSManager, error) {
	tcnl, err := tc.Open(&tc.Config{})
	if err != nil {
		return nil, err
	}

	htbCtx := QoSManager{}

	htbCtx.TcConn = tcnl

	return &htbCtx, nil
}

func (c *QoSManager) WithLogger(l *slog.Logger) {
	c.Logger = l
}

func (c *QoSManager) InitQoSClassifier(createIfNotExists bool) error {
	nftCtx, err := nft.NewNFTCtx(nft.NFTOpts{
		CreateIfNotExists: createIfNotExists,
		Logger:            c.Logger,
	})
	if err != nil {
		return err
	}
	c.Classifier = &nftCtx

	return nil
}

func (c *QoSManager) AddTargetsToPriority(target []netip.Prefix, priority prio.Priority) (err error) {
	if c.Classifier == nil {
		return fmt.Errorf(" HTB filter uninitialised")
	}

	switch priority {
	case prio.PRIORITYHIGH:
		err = c.Classifier.AddTargetsToHighPriority(target)
	case prio.PRIORITYLOW:
		err = c.Classifier.AddTargetsToLowPriority(target)
	default:
		return fmt.Errorf("unknown priority %v", priority)
	}
	if err != nil {
		return err
	}

	return nil
}

func (c *QoSManager) RemoveTargetsFromPriority(target []netip.Prefix, priority prio.Priority) (err error) {
	if c.Classifier == nil {
		return fmt.Errorf(" HTB filter uninitialised")
	}

	switch priority {
	case prio.PRIORITYHIGH:
		err = c.Classifier.DeleteTargetFromHighPriority(target)
	case prio.PRIORITYLOW:
		err = c.Classifier.DeleteTargetFromLowPriority(target)
	default:
		return fmt.Errorf("unknown priority %v", priority)
	}
	if err != nil {
		return err
	}

	return nil
}
