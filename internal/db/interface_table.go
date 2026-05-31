package db

import (
	"database/sql"
	"errors"
	"fmt"
)

var ErrNotExists = errors.New("requested object does not exist in the database")

type Interface struct {
	IfaceIndex int
	Name       string
	Enabled    bool
}

func AddInterface(db *sql.DB, iface Interface) error {
	_, err := db.Exec(
		`
		INSERT OR REPLACE INTO interfaces (
			name,
			if_index,
			enabled
		)
		VALUES (?, ?, ?)
	`,
		iface.Name,
		iface.IfaceIndex,
		iface.Enabled,
	)

	return err
}

func CheckInterfaceExists(db *sql.DB, name string) (bool, error) {
	var exists bool

	err := db.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM interfaces WHERE  name = ?
		)
	`, name).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func CheckInterfaceExistsByIndex(db *sql.DB, index int) (bool, error) {
	var exists bool

	err := db.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM interfaces WHERE  if_index = ?
		)
	`, index).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func GetAllInterfaces(db *sql.DB) ([]Interface, error) {
	rows, err := db.Query(`
		SELECT if_index, name, enabled
		FROM interfaces
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var interfaces []Interface
	var enabled int

	for rows.Next() {
		var iface Interface
		err = rows.Scan(&iface.IfaceIndex, &iface.Name, &enabled)
		if err != nil {
			return nil, err
		}

		iface.Enabled = enabled == 1
		interfaces = append(interfaces, iface)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return interfaces, nil
}

func InterfaceByName(db *sql.DB, name string) (Interface, error) {
	row := db.QueryRow(`
		SELECT if_index, name, enabled 
		FROM interfaces
		WHERE name = ?
	`, name)

	var iface Interface
	var enabled int

	err := row.Scan(&iface.IfaceIndex, &iface.Name, &enabled)
	if err != nil {
		return Interface{}, err
	}
	iface.Enabled = enabled == 1

	return iface, nil
}

func InterfaceByIndex(db *sql.DB, index int) (Interface, error) {
	row := db.QueryRow(`
		SELECT if_index, name, enabled
		FROM interfaces
		WHERE if_index = ?
	`, index)

	var iface Interface
	var enabled int

	err := row.Scan(&iface.IfaceIndex, &iface.Name, &enabled)
	if err != nil {
		return Interface{}, err
	}
	iface.Enabled = enabled == 1

	return iface, nil
}

func DeleteInterfaceByName(db *sql.DB, name string) error {
	_, err := db.Exec(`
		DELETE FROM interfaces
		WHERE name = ?
	`, name)

	return err
}

func DeleteInterfaceByIndex(db *sql.DB, index int) error {
	_, err := db.Exec(`
		DELETE FROM interfaces
		WHERE if_index = ?
	`, index)

	return err
}

func DisableInterface(db *sql.DB, name string) error {
	return updateField(db, name, "enabled", false)
}

func EnableInterface(db *sql.DB, name string) error {
	return updateField(db, name, "enabled", true)
}

func InterfaceEnabled(db *sql.DB, name string) (bool, error) {
	enabled, err := getField(db, name, "enabled")
	if err != nil {
		if errors.Is(err, ErrNotExists) {
			return false, nil
		}
		return false, err
	}

	return enabled.(int64) == 1, nil
}

func getField(db *sql.DB, ifaceName string, field string) (any, error) {
	allowed := map[string]struct{}{
		"if_index": {},
		"name":     {},
		"enabled":  {},
	}
	if _, ok := allowed[field]; !ok {
		return nil, fmt.Errorf("db: unknown interfaces table field: %v", field)
	}

	query := fmt.Sprintf(`
		SELECT %s 
		FROM interfaces
		WHERE name = ?
	`, field)

	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}

	row := stmt.QueryRow(ifaceName)

	var value any
	err = row.Scan(&value)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotExists
	}

	return value, nil
}

func updateField(db *sql.DB, ifaceName string, field string, value any) error {
	allowed := map[string]struct{}{
		"if_index": {},
		"name":     {},
		"enabled":  {},
	}
	if _, ok := allowed[field]; !ok {
		return fmt.Errorf("db: unknown interfaces table field: %v", field)
	}

	query := fmt.Sprintf(`
		UPDATE interfaces
		SET %s = ?
		WHERE name = ?
	`, field)

	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(value, ifaceName)

	return err
}
