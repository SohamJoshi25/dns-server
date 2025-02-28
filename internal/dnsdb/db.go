package dnsdb

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/sohamjoshi25/dns-server/internal/dnslookup"
)

var DB *sql.DB

func init() {
	db, err := sql.Open("postgres", "user=postgres password=mypassword dbname=postgres sslmode=disable")

	if err != nil {
		fmt.Println("Could not connect to the database:", err)
		os.Exit(1)
	}

	query := `
	CREATE TABLE IF NOT EXISTS dns_records (
		id SERIAL PRIMARY KEY,
		domain VARCHAR(255) NOT NULL,
		type VARCHAR(10) NOT NULL,
		answer TEXT NOT NULL
	);
	`
	_, err = db.Exec(query)
	if err != nil {
		fmt.Println("Failed to create table:", err)
	}

	DB = db
}

func QueryDatabase(domain string, rrType uint16) ([]string, error) {
	var answers []string

	typeStr, okType := dnslookup.RRTypeMap[rrType]

	if !okType {
		return nil, fmt.Errorf("Unknown RR Type")
	}

	query := `SELECT answer FROM dns_records WHERE domain=$1 AND type=$2`
	rows, err := DB.Query(query, domain, typeStr)

	if err != nil {
		fmt.Printf("Database query error: %v\n", err)
		return nil, err
	}

	defer rows.Close()

	found := false
	for rows.Next() {
		var answer string
		if err := rows.Scan(&answer); err != nil {
			fmt.Printf("Failed to scan row: %v\n", err)
			continue
		}
		answers = append(answers, answer)
		found = true
	}

	if !found {
		return nil, fmt.Errorf("No Record Found")
	}

	return answers, nil
}

func InsertRecord(domain, rtype, answer string) {
	ResetSequence()
	query := `INSERT INTO dns_records (domain, type, answer) VALUES ($1, $2, $3)`
	_, err := DB.Exec(query, domain, rtype, answer)
	if err != nil {
		fmt.Println("Failed to insert DNS record:", err, query, domain, rtype, answer)
	} else {
		fmt.Printf("Record inserted successfully\n   Domain: %s\n   Answer:%s\n   Type:%s\n", domain, answer, rtype)
	}
	defer DB.Close()
}

func ResetSequence() {
	query := `
		ALTER SEQUENCE dns_records_id_seq RESTART WITH 1;
		UPDATE dns_records SET id = DEFAULT;
	`
	_, err := DB.Exec(query)
	if err != nil {
		fmt.Println("Failed to reset sequence:", err)
	}
}

func GetAllRecords() {
	query := `SELECT id, domain, type, answer FROM dns_records`
	rows, err := DB.Query(query)
	if err != nil {
		fmt.Println("Failed to fetch DNS records:", err)
		return
	}
	defer rows.Close()
	defer DB.Close()

	for rows.Next() {
		var id int
		var domain, rtype, answer string

		if err := rows.Scan(&id, &domain, &rtype, &answer); err != nil {
			fmt.Println("Failed to scan row:", err)
			continue
		}

		fmt.Println(id, domain, rtype, answer)
	}

}

func DeleteRecordByID(id int) {
	query := `DELETE FROM dns_records WHERE id=$1`
	res, err := DB.Exec(query, id)
	if err != nil {
		fmt.Println("Failed to delete DNS record:", err)
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		fmt.Printf("No record found with ID %d\n", id)
		return
	}

	fmt.Printf("Record with ID %d deleted successfully\n", id)
	defer DB.Close()
}
