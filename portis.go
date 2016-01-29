package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/tcyrus/portis/Godeps/_workspace/src/github.com/mattn/go-sqlite3"
	"github.com/tcyrus/portis/Godeps/_workspace/src/github.com/syohex/go-texttable"
	"path"
	"runtime"
	"strconv"
)

const BASE_SQL = "SELECT name, port, protocol, description FROM ports WHERE "

// This function checks if the string can be cast into an integer
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

	var where_field string = "name"
	if validPort(port) {
		where_field = "port"
	}

	var where_value string = "='%s'"
	if like {
		where_value = " LIKE '%%%s%%'"
	}

	var sql string = BASE_SQL + where_field + fmt.Sprintf(where_value, port)

	rows, _ := db.Query(sql)
	defer rows.Close()

	makeTable(rows)
}

// This function returns a pretty table used to display the port results.
func makeTable(rows *sql.Rows) {
	tbl := &texttable.TextTable{}
	tbl.SetHeader("Name", "Port", "Protocol", "Description")

	for rows.Next() {
		var name, port, protocol, description string
		rows.Scan(&name, &port, &protocol, &description)
		tbl.AddRow(name, port, protocol, description)
	}

	fmt.Println(tbl.Draw())
}

func main() {
	var like bool

	flag.BoolVar(&like, "like", false, "search ports containing the pattern")
	flag.Parse()

	if len(flag.Args()) > 0 {
		getPorts(flag.Args()[0], like)
	}
}
