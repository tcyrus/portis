package main

import (
	"database/sql"
	"encoding/csv"
	"flag"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/syohex/go-texttable"
	"log"
	"net/http"
	"path"
	"runtime"
	"strconv"
)

const (
	GET_DATA string = "SELECT name, port, protocol, description FROM ports WHERE "
	DELETE_TABLE string = "DROP TABLE IF EXISTS ports"
	CREATE_TABLE string = `CREATE TABLE ports (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		name VARCHAR(50),
		port VARCHAR(12),
		protocol VARCHAR(3),
		description VARCHAR(100)
	)`
	INSERT_TABLE string = "INSERT INTO ports (id, name, port, protocol, description) VALUES (?, ?, ?, ?, ?)"
	IANA_URL string = "http://www.iana.org/assignments/service-names-port-numbers/service-names-port-numbers.csv"
)

// This function checks if the string can be cast into an integer;
// if a string can be cast into an integer, then it is a valid port.
func validPort(v string) bool {
	_, err := strconv.Atoi(v)
	return (err == nil)
}

// Gets Absolute Path for ports.db
func dbPath() string {
	// andrewbrookins.com/tech/golang-get-directory-of-the-current-file/
	_, dbPath, _, _ := runtime.Caller(1)
	return path.Join(path.Dir(dbPath), "ports.db")
}

// This function creates the SQL query depending on the specified port and
// the --like option. It then calls makeTable.
func getPorts(port string, like bool) {
	db, err := sql.Open("sqlite3", dbPath())
	if err != nil {
    	log.Fatal(err)
	}
	defer db.Close()

	WHERE_FIELD := "name"
	if validPort(port) {
		WHERE_FIELD = "port"
	}

	if like {
		port = "%" + port + "%"
	}

	rows, err := db.Query(GET_DATA + WHERE_FIELD + " LIKE ?", port)
	if err != nil {
    	log.Fatal(err)
	}
	defer rows.Close()

	makeTable(rows)
}

// Update SQLite Database
func updatePorts() {
	resp, err := http.Get(IANA_URL)
	if err != nil {
    	log.Fatal(err)
	}
	defer resp.Body.Close()
	fmt.Println("Downloaded CSV")

	reader := csv.NewReader(resp.Body)
	reader.FieldsPerRecord = -1

	rawCSVdata, err := reader.ReadAll()
	if err != nil {
    	log.Fatal(err)
	}
	fmt.Println("Parsed CSV")

	db, err := sql.Open("sqlite3", dbPath())
	if err != nil {
    	log.Fatal(err)
	}
	defer db.Close()

	db.Exec(DELETE_TABLE)
	fmt.Println("Deleted Table")

	db.Exec(CREATE_TABLE)
	fmt.Println("Added Table Schema")

	for _, each := range rawCSVdata[1:] {
		db.Exec(INSERT_TABLE, nil, each[0], each[1], each[2], each[3])
	}
	fmt.Println("Database Update Complete")
}

// This function returns a pretty table used to display the port results.
func makeTable(rows *sql.Rows) {
	tbl := &texttable.TextTable{}
	tbl.SetHeader("Name", "Port", "Protocol", "Description")

	var name, port, protocol, description string

	for rows.Next() {
		rows.Scan(&name, &port, &protocol, &description)
		tbl.AddRow(name, port, protocol, description)
	}

	fmt.Println(tbl.Draw())
}

func main() {
	like := flag.Bool("like", false, "search ports containing the pattern")

	flag.Parse()

	port := flag.Arg(0)

	if port != "" {
		getPorts(port, *like)
	} else {
		updatePorts()
	}
}
