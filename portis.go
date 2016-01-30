package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/syohex/go-texttable"
	"path"
	"runtime"
	"strconv"
)

const BASE_SQL = "SELECT name, port, protocol, description FROM ports WHERE "

// This function checks if the string can be cast into an integer;
// if a string can be cast into an integer, then it is a valid port.
func validPort(v string) bool {
	_, err := strconv.Atoi(v)
	return (err == nil)
}

// This function creates the SQL query depending on the specified port and
// the --like option. It then calls makeTable.
func getPorts(port string, like bool) {
	// andrewbrookins.com/tech/golang-get-directory-of-the-current-file/
	_, filename, _, _ := runtime.Caller(1)
	filename = path.Join(path.Dir(filename), "ports.db")

	db, _ := sql.Open("sqlite3", filename)
	defer db.Close()

	WHERE_FIELD := "name"
	if validPort(port) {
		WHERE_FIELD = "port"
	}

	if like {
		port = fmt.Sprintf("%%%s%%", port)
	}

	rows, _ := db.Query(BASE_SQL + WHERE_FIELD + " LIKE ?", port)
	defer rows.Close()

	makeTable(rows)
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
	}
}
