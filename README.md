# `portis` - a command to search port names and numbers

It often happens that we need to find the default port number for a specific service, or what service is listening on a given port.

This application is based on [ncrocfer/whatportis](https://github.com/ncrocfer/whatportis).

The difference with this version is:
* It has a shorter name
* It's written in Go
* It's faster
* You can update the database

## Install

```
$ go get github.com/tcyrus/portis
```

## Usage

This tool allows you to find what port is associated with a service:

```
$ portis redis
+-------+------+----------+---------------------------------------+
| Name  | Port | Protocol | Description                           |
+-------+------+----------+---------------------------------------+
| redis | 6379 | tcp      | An advanced key-value cache and store |
+-------+------+----------+---------------------------------------+
```

Or, conversely, what service is associated with a port number:

```
$ portis 5432
+------------+------+----------+---------------------+
| Name       | Port | Protocol | Description         |
+------------+------+----------+---------------------+
| postgresql | 5432 | tcp      | PostgreSQL Database |
| postgresql | 5432 | udp      | PostgreSQL Database |
+------------+------+----------+---------------------+
```


## Notes

- You can search a pattern without knowing the exact name by adding the `--like` option:

```
$ portis --like mysql
+----------------+-------+----------+-----------------------------------+
| Name           | Port  | Protocol | Description                       |
+----------------+-------+----------+-----------------------------------+
| mysql-cluster  |  1186 | tcp      | MySQL Cluster Manager             |
| mysql-cluster  |  1186 | udp      | MySQL Cluster Manager             |
| mysql-cm-agent |  1862 | tcp      | MySQL Cluster Manager Agent       |
| mysql-cm-agent |  1862 | udp      | MySQL Cluster Manager Agent       |
| mysql-im       |  2273 | tcp      | MySQL Instance Manager            |
| mysql-im       |  2273 | udp      | MySQL Instance Manager            |
| mysql          |  3306 | tcp      | MySQL                             |
| mysql          |  3306 | udp      | MySQL                             |
| mysql-proxy    |  6446 | tcp      | MySQL Proxy                       |
| mysql-proxy    |  6446 | udp      | MySQL Proxy                       |
| mysqlx         | 33060 | tcp      | MySQL Database Extended Interface |
+----------------+-------+----------+-----------------------------------+
```

- If you have no parameters, the database will update:

```
$ portis
Downloaded CSV
Parsed CSV
Deleted Table
Added Table Schema
Database Update Complete
```

- Why not use `grep <port> /etc/services`? Simply because I want a portable command that display the output in a nice format (a pretty table).

- The tool uses the [iana.org](http://www.iana.org/assignments/port-numbers) website to get the official list of ports.
