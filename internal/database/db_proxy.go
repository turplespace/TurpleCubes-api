package database

import (
	"database/sql"
	"fmt"
)

type Proxy struct {
	ID        int
	CubeID    int
	Domain    string
	Port      int
	Type      string
	Default   bool
	CreatedAt string
}

func GetProxyByID(id int) (*Proxy, error) {
	db_path, _ := GetPath()
	db, err := sql.Open("sqlite3", db_path)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}
	defer db.Close()
	query := `SELECT id, cube_id, domain, port, type, "default", created_at FROM proxy WHERE id = ?`
	row := db.QueryRow(query, id)

	var proxy Proxy
	err = row.Scan(&proxy.ID, &proxy.CubeID, &proxy.Domain, &proxy.Port, &proxy.Type, &proxy.Default, &proxy.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get proxy by id: %v", err)
	}

	return &proxy, nil
}

func GetProxiesByCubeID(cubeID int) ([]Proxy, error) {
	db_path, _ := GetPath()
	db, err := sql.Open("sqlite3", db_path)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}
	defer db.Close()
	query := `SELECT id, cube_id, domain, port, type, "default", created_at FROM proxy WHERE cube_id = ?`
	rows, err := db.Query(query, cubeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get proxies by cube_id: %v", err)
	}
	defer rows.Close()

	var proxies []Proxy
	for rows.Next() {
		var proxy Proxy
		err := rows.Scan(&proxy.ID, &proxy.CubeID, &proxy.Domain, &proxy.Port, &proxy.Type, &proxy.Default, &proxy.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan proxy: %v", err)
		}
		proxies = append(proxies, proxy)
	}

	return proxies, nil
}

func AddProxy(cubeID int, domain string, port int, proxyType string, isDefault bool) (int64, error) {
	db_path, _ := GetPath()
	db, err := sql.Open("sqlite3", db_path)
	if err != nil {
		return 0, fmt.Errorf("failed to open database: %v", err)
	}
	defer db.Close()
	query := `INSERT INTO proxy (cube_id, domain, port, type, "default", created_at) VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP)`
	result, err := db.Exec(query, cubeID, domain, port, proxyType, isDefault)
	if err != nil {
		return 0, fmt.Errorf("failed to add proxy: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id: %v", err)
	}

	return id, nil
}

func EditProxyByID(id int, domain string, port int, proxyType string, isDefault bool) error {
	db_path, _ := GetPath()
	db, err := sql.Open("sqlite3", db_path)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}
	defer db.Close()
	query := `UPDATE proxy SET domain = ?, port = ?, type = ?, "default" = ? WHERE id = ?`
	_, err = db.Exec(query, domain, port, proxyType, isDefault, id)
	if err != nil {
		return fmt.Errorf("failed to edit proxy by id: %v", err)
	}

	return nil
}

func DeleteProxyByID(id int) error {
	db_path, _ := GetPath()
	db, err := sql.Open("sqlite3", db_path)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}
	defer db.Close()
	query := `DELETE FROM proxy WHERE id = ?`
	_, err = db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete proxy by id: %v", err)
	}

	return nil
}

func DeleteProxiesByCubeID(cubeID int) error {
	db_path, _ := GetPath()
	db, err := sql.Open("sqlite3", db_path)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}
	defer db.Close()
	query := `DELETE FROM proxy WHERE cube_id = ?`
	_, err = db.Exec(query, cubeID)
	if err != nil {
		return fmt.Errorf("failed to delete proxies by cube_id: %v", err)
	}

	return nil
}
