// Package tc is used to manage traffic control settings, statistics etc of an interface.
package tc

import (
	"database/sql"
	"fmt"
	"net"

	"github.com/kakeetopius/qosm/internal/core/htb"
	"github.com/kakeetopius/qosm/internal/db"
)

func EnableTcOnInterface(iface net.Interface, htbCtx *htb.HTBCtx, dbConn *sql.DB) (err error) {
	defer func() {
		if err != nil {
			db.AddErrorLog(dbConn, err, "")
		} else {
			addSuccessLog(dbConn, iface.Name)
		}
	}()

	if htbCtx == nil {
		return fmt.Errorf("htb context not intialised")
	}
	err = htbCtx.InitHTBIface(iface.Name)
	if err != nil {
		return err
	}

	if htbCtx.NFTFilter == nil {
		return fmt.Errorf("nft filter not intialised")
	}

	err = htbCtx.NFTFilter.AddIfaceRules(iface.Index)
	if err != nil {
		return err
	}

	err = db.AddInterface(dbConn, db.Interface{
		Name:       iface.Name,
		IfaceIndex: iface.Index,
		Enabled:    true,
	})
	if err != nil {
		return err
	}

	return nil
}

func DisableTcOnInterface(iface net.Interface, htbCtx *htb.HTBCtx, dbConn *sql.DB) (err error) {
	defer func() {
		if err != nil {
			db.AddErrorLog(dbConn, err, "")
		} else {
			addDeleteLog(dbConn, iface.Name)
		}
	}()

	if htbCtx == nil {
		return fmt.Errorf("htb context not intialised")
	}
	err = htbCtx.FlushQdisc(iface.Index)
	if err != nil {
		return err
	}

	if htbCtx.NFTFilter != nil {
		err = htbCtx.NFTFilter.DeleteIfaceRules(iface.Index)
		if err != nil {
			return err
		}
	}

	err = db.DisableInterface(dbConn, iface.Name)
	if err != nil {
		return err
	}

	return nil
}

func InitSavedInterfaceSettings(dbConn *sql.DB, htbCtx *htb.HTBCtx) error {
	if htbCtx == nil {
		return fmt.Errorf("htb context not intialised")
	}
	if htbCtx.NFTFilter == nil {
		return fmt.Errorf("nft filter not intialised")
	}
	enableIfaces, err := db.GetEnabledInterfaces(dbConn)
	if err != nil {
		return err
	}

	for _, iface := range enableIfaces {
		err := htbCtx.InitHTBIface(iface.Name)
		if err != nil {
			return err
		}
		err = htbCtx.NFTFilter.AddIfaceRules(iface.IfaceIndex)
		if err != nil {
			return err
		}
	}

	return nil
}

func addSuccessLog(dbCon *sql.DB, iface string) error {
	return db.AddLog(
		dbCon,
		db.Log{
			EventType:   "TC",
			Description: fmt.Sprintf("enabled traffic control on interface %s", iface),
		},
	)
}

func addDeleteLog(dbCon *sql.DB, iface string) error {
	return db.AddLog(
		dbCon,
		db.Log{
			EventType:   "TC",
			Description: fmt.Sprintf("disabled traffic control on interface %s", iface),
		},
	)
}
