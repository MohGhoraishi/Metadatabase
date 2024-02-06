package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/muesli/termenv"
)

// types for db manipulation
type db_table struct {
	id   int
	name string
}

type tble_table struct {
	id    int
	name  string
	db_id int
}

type attribute_table struct {
	id       int
	name     string
	a_type   string
	table_id int
}

type constraints_table struct {
	id           int
	name         string
	c_type       string
	c_detail     string
	attribute_id int
}

// state functions
func main() {
	termenv.AltScreen()
	home(1, false)
}

func home(counter int, with_err bool) {
	termenv.ClearScreen()
	println("--Not a Database-- \n \n")
	if with_err {
		println("invalid number. Choose another number \n")
	} else {
		println("Choose a number \n")
	}
	println("1. Create a new database \n")
	println("2. Select an existing database \n")
	println("3. Test backend database connection \n")
	println("4. Exit \n")
	var input int
	fmt.Scanln(&input)
	switch input {
	case 1:
		db_create()
	case 2:
		db_select()
	case 3:
		db_test()
	case 4:
		termenv.ClearScreen()
		os.Exit(0)
	default:
		home(1, true)
	}
}

func db_test() {
	termenv.ClearScreen()
	println("--Database Connection-- (to quit - :q)\n \n")
	println("Connection status: \n")
	db, err := sql.Open("mysql", "root:Gooli2002@/metadatabase")
	if err != nil {
		panic(err)
	}
	pingErr := db.Ping()
	if pingErr != nil {
		panic(pingErr)
	}
	println("Success\n")
	println("Press any key to return to main menu.")
	var in string
	fmt.Scanln(&in)
	check_exit(in)
	home(1, false)
}

func db_create() {
	termenv.ClearScreen()
	db, err := sql.Open("mysql", "root:Gooli2002@/metadatabase")
	if err != nil {
		panic(err)
	}
	println("--Create Database-- (to quit - :q)\n \n")
	println("Enter database name: \n")
	var dbname string
	fmt.Scanln(&dbname)
	check_exit(dbname)
	database := db_table{name: dbname}
	result, err := db.Exec("INSERT INTO db (database_name) VALUES (?)", database.name)
	if err != nil {
		println("ERROR: ")
		fmt.Println(err)
		println("\nPress any key to try again.")
		var in string
		fmt.Scan(&in)
		db_create()
	}
	id, err := result.LastInsertId()
	if err != nil {
		println("ERROR at %d: ", id)
		fmt.Println(err)
		println("\nPress any key to try again.")
		var in string
		fmt.Scan(&in)
		db_create()
	}
	println("Great Success!\nPress any key to return to main menu.")
	var in string
	fmt.Scanln(&in)
	check_exit(in)
	home(1, false)
}

func db_select() {
	termenv.ClearScreen()
	db, err := sql.Open("mysql", "root:Gooli2002@/metadatabase")
	if err != nil {
		panic(err)
	}
	var dbs []db_table
	rows, err := db.Query("SELECT * FROM db")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var dbt db_table
		if err := rows.Scan(&dbt.id, &dbt.name); err != nil {
			panic(err)
		}
		dbs = append(dbs, dbt)
	}
	println("--Available Databases--\n\n")
	for i := 0; i < len(dbs); i++ {
		name := dbs[i].name
		id := dbs[i].id
		fmt.Printf("%d. name: %s, id: %d \n", i+1, name, id)
	}
	println("\nSelect a database with it's corresponding number.")
	var in string
	fmt.Scanln(&in)
	check_exit(in)
	pos, err := strconv.Atoi(in)
	if err != nil {
		db_select()
	}
	if (pos > 0) && (pos <= len(dbs)) {
		db_id := dbs[pos-1].id
		db_name := dbs[pos-1].name
		table_choice(db_id, db_name)
	} else {
		println("Number out of bounds, try again.\n")
		var in string
		fmt.Scanln(&in)
		check_exit(in)
		db_select()
	}
}

func table_choice(db_id int, db_name string) {
	termenv.ClearScreen()
	fmt.Printf("--Table Options-- (current db = %s)\n\n", db_name)
	println("1. Create new table\n")
	println("2. Select existing table\n")
	println("\nChoose a number:")
	var input int
	fmt.Scanln(&input)
	switch input {
	case 1:
		table_create(db_id, db_name)
	case 2:
		table_select(db_id, db_name)
	default:
		println("Invalid number, press any button to try again.\n")
		var in string
		fmt.Scanln(&in)
		check_exit(in)
		table_choice(db_id, db_name)
	}
}

func table_create(db_id int, db_name string) {
	termenv.ClearScreen()
	fmt.Printf("--Create Table-- (current db = %s)\n\n", db_name)
	db, err := sql.Open("mysql", "root:Gooli2002@/metadatabase")
	if err != nil {
		panic(err)
	}
	println("Enter table name: \n")
	var tblename string
	fmt.Scanln(&tblename)
	check_exit(tblename)
	table := tble_table{name: tblename, db_id: db_id}
	result, err := db.Exec("INSERT INTO tble (table_name, db_id) VALUES (?, ?)", table.name, table.db_id)
	id, err := result.LastInsertId()
	if err != nil {
		println("ERROR at %d: ", id)
		fmt.Println(err)
		println("\nPress any key to try again.")
		var in string
		fmt.Scan(&in)
		table_create(db_id, db_name)
	}
	println("Great Success!\nPress any key to return to the table option menu.")
	var in string
	fmt.Scanln(&in)
	check_exit(in)
	table_choice(db_id, db_name)
}

func table_select(db_id int, db_name string) {
	termenv.ClearScreen()
	fmt.Printf("--Available Tables-- (current db = %s)\n\n", db_name)
	db, err := sql.Open("mysql", "root:Gooli2002@/metadatabase")
	if err != nil {
		panic(err)
	}
	var tbls []tble_table
	rows, err := db.Query("SELECT table_id, table_name FROM tble WHERE db_id = ?", db_id)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var tbl tble_table
		if err := rows.Scan(&tbl.id, &tbl.name); err != nil {
			panic(err)
		}
		tbls = append(tbls, tbl)
	}
	if len(tbls) == 0 {
		println("No tables available, press any key to return to menu.\n")
		var in string
		fmt.Scanln(&in)
		check_exit(in)
		table_choice(db_id, db_name)
	}
	for i := 0; i < len(tbls); i++ {
		name := tbls[i].name
		id := tbls[i].id
		fmt.Printf("%d. name: %s, id: %d \n", i+1, name, id)
	}
	println("\nSelect a table with it's corresponding number.")
	var in string
	fmt.Scanln(&in)
	check_exit(in)
	pos, err := strconv.Atoi(in)
	if err != nil {
		table_select(db_id, db_name)
	}
	if (pos > 0) && (pos <= len(tbls)) {
		tbl_id := tbls[pos-1].id
		tbl_name := tbls[pos-1].name
		table_more_choice(db_id, db_name, tbl_id, tbl_name)
	} else {
		println("Number out of bounds, try again.\n")
		var in string
		fmt.Scanln(&in)
		check_exit(in)
		table_select(db_id, db_name)
	}
}

func table_more_choice(db_id int, db_name string, tbl_id int, tbl_name string) {
	termenv.ClearScreen()
	fmt.Printf("--Modify Table-- (current db = %s) (current table = %s)\n\n", db_name, tbl_name)
	println("1. Add new column\n")
	println("2. Select existing Column\n")
	println("3. Insert data\n")
	println("4. Exit\n")
	println("\nChoose a number:")
	var input int
	fmt.Scanln(&input)
	switch input {
	case 1:
		column_create(db_id, db_name, tbl_id, tbl_name)
	case 2:
		column_select(db_id, db_name, tbl_id, tbl_name)
	case 3:
		//TODO:
		//data_insert(db_id, db_name, tbl_id, tbl_name)
	case 4:
		termenv.ClearScreen()
		os.Exit(0)
	default:
		println("Invalid number, press any button to try again.\n")
		var in string
		fmt.Scanln(&in)
		check_exit(in)
		table_more_choice(db_id, db_name, tbl_id, tbl_name)
	}
}

func column_create(db_id int, db_name string, tbl_id int, tbl_name string) {
	termenv.ClearScreen()
	fmt.Printf("--Create Column-- (current db = %s) (current table = %s)\n\n", db_name, tbl_name)
	db, err := sql.Open("mysql", "root:Gooli2002@/metadatabase")
	if err != nil {
		panic(err)
	}
	println("Enter column name: \n")
	var column_name string
	fmt.Scanln(&column_name)
	check_exit(column_name)
	println("Enter column type: ('INT', 'CHAR', 'FLOAT', 'VARCHAR', 'BINARY', 'BOOL') \n")
	var column_type string
	fmt.Scanln(&column_type)
	check_exit(column_type)
	column := attribute_table{name: column_name, a_type: column_type, table_id: tbl_id}
	result, err := db.Exec("INSERT INTO attribute (attribute_name, attribute_type, table_id) VALUES (?, ?, ?)", column.name, column.a_type, column.table_id)
	id, err := result.LastInsertId()
	if err != nil {
		println("ERROR at %d: ", id)
		fmt.Println(err)
		println("\nPress any key to try again.")
		var in string
		fmt.Scan(&in)
		column_create(db_id, db_name, tbl_id, tbl_name)
	}
	println("Great Success!\nPress any key to return to the table option menu.")
	var in string
	fmt.Scanln(&in)
	check_exit(in)
	table_more_choice(db_id, db_name, tbl_id, tbl_name)
}

func column_select(db_id int, db_name string, tbl_id int, tbl_name string) {
	termenv.ClearScreen()
	fmt.Printf("--Available Columns-- (current db = %s) (current table = %s)\n\n", db_name, tbl_name)
	db, err := sql.Open("mysql", "root:Gooli2002@/metadatabase")
	if err != nil {
		panic(err)
	}
	var columns []attribute_table
	rows, err := db.Query("SELECT attribute_id, attribute_name, attribute_type FROM attribute WHERE table_id = ?", tbl_id)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var clmn attribute_table
		if err := rows.Scan(&clmn.id, &clmn.name, &clmn.a_type); err != nil {
			panic(err)
		}
		columns = append(columns, clmn)
	}
	if len(columns) == 0 {
		println("No columns available, press any key to return to menu.\n")
		var in string
		fmt.Scanln(&in)
		check_exit(in)
		table_more_choice(db_id, db_name, tbl_id, tbl_name)
	}
	for i := 0; i < len(columns); i++ {
		name := columns[i].name
		id := columns[i].id
		value_type := columns[i].a_type
		fmt.Printf("%d. name: %s, type: %s, id: %d \n", i+1, name, value_type, id)
	}
	println("\nSelect a column with it's corresponding number.")
	var in string
	fmt.Scanln(&in)
	check_exit(in)
	pos, err := strconv.Atoi(in)
	if err != nil {
		column_select(db_id, db_name, tbl_id, tbl_name)
	}
	if (pos > 0) && (pos <= len(columns)) {
		column_id := columns[pos-1].id
		column_name := columns[pos-1].name
		column_choice(db_id, db_name, tbl_id, tbl_name, column_id, column_name)
	} else {
		println("Number out of bounds, try again.\n")
		var in string
		fmt.Scanln(&in)
		check_exit(in)
		column_select(db_id, db_name, tbl_id, tbl_name)
	}
}

func column_choice(db_id int, db_name string, tbl_id int, tbl_name string, column_id int, column_name string) {
	termenv.ClearScreen()
	fmt.Printf("--Modify Column-- (current db = %s) (current table = %s) (current column = %s)\n\n", db_name, tbl_name, column_name)
	println("1. Create new constraint\n")
	println("2. View constraints\n")
	println("3. Exit\n")
	println("\nChoose a number:")
	var input int
	fmt.Scanln(&input)
	switch input {
	case 1:
		constraint_create(db_id, db_name, tbl_id, tbl_name, column_id, column_name)
	case 2:
		constraint_select(db_id, db_name, tbl_id, tbl_name, column_id, column_name)
	case 3:
		termenv.ClearScreen()
		os.Exit(0)
	default:
		println("Invalid number, press any button to try again.\n")
		var in string
		fmt.Scanln(&in)
		check_exit(in)
		column_choice(db_id, db_name, tbl_id, tbl_name, column_id, column_name)
	}
}

func constraint_create(db_id int, db_name string, tbl_id int, tbl_name string, column_id int, column_name string) {
	termenv.ClearScreen()
	fmt.Printf("--Create Constraint-- (current db = %s) (current table = %s) (current column = %s)\n\n", db_name, tbl_name, column_name)
	db, err := sql.Open("mysql", "root:Gooli2002@/metadatabase")
	if err != nil {
		panic(err)
	}
	println("Enter constraint name:")
	var constraint_name string
	fmt.Scanln(&constraint_name)
	check_exit(constraint_name)
	reader := bufio.NewReader(os.Stdin)
	println("Enter constraint type: ('NOT NULL', 'PRIMARY KEY', 'FOREIGN KEY', 'UNIQUE') \n")
	constraint_type, _ := reader.ReadString('\n')
	constraint_type = strings.TrimSpace(constraint_type)
	check_exit(constraint_type)
	var constraint_detail string
	constraint_detail = "NULL"

	if constraint_type == "FOREIGN KEY" {
		println("References:\n")
		fmt.Scanln(&constraint_detail)
	}
	constraints := constraints_table{name: constraint_name, c_type: constraint_type, c_detail: constraint_detail, attribute_id: column_id}
	result, err := db.Exec("INSERT INTO constraints (constraint_name, constraint_type, constraint_detail, attribute_id) VALUES (?, ?, ?, ?)", constraints.name, constraints.c_type, constraints.c_detail, constraints.attribute_id)
	id, err := result.LastInsertId()
	if err != nil {
		println("ERROR at %d: ", id)
		fmt.Println(err)
		println("\nPress any key to try again.")
		var in string
		fmt.Scan(&in)
		constraint_create(db_id, db_name, tbl_id, tbl_name, column_id, column_name)
	}
	println("Great Success!\nPress any key to return to the column option menu.")
	var in string
	fmt.Scanln(&in)
	check_exit(in)
	column_choice(db_id, db_name, tbl_id, tbl_name, column_id, column_name)
}

func constraint_select(db_id int, db_name string, tbl_id int, tbl_name string, column_id int, column_name string) {
	termenv.ClearScreen()
	fmt.Printf("--Create Constraint-- (current db = %s) (current table = %s) (current column = %s)\n\n", db_name, tbl_name, column_name)
	db, err := sql.Open("mysql", "root:Gooli2002@/metadatabase")
	if err != nil {
		panic(err)
	}
	var constraints []constraints_table
	rows, err := db.Query("SELECT constraint_id, constraint_name, constraint_type, constraint_detail FROM constraints WHERE attribute_id = ?", column_id)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var cst constraints_table
		if err := rows.Scan(&cst.id, &cst.name, &cst.c_type, &cst.c_detail); err != nil {
			panic(err)
		}
		constraints = append(constraints, cst)
	}
	if len(constraints) == 0 {
		println("No constraints available, press any key to return to menu.\n")
		var in string
		fmt.Scanln(&in)
		check_exit(in)
		column_choice(db_id, db_name, tbl_id, tbl_name, column_id, column_name)
	}
	for i := 0; i < len(constraints); i++ {
		name := constraints[i].name
		id := constraints[i].id
		value_type := constraints[i].c_type
		detail := constraints[i].c_detail
		if detail == "NULL" {
			fmt.Printf("%d. name: %s, type: %s, id: %d \n", i+1, name, value_type, id)
		} else {
			fmt.Printf("%d. name: %s, type: %s, references: %s, id: %d \n", i+1, name, value_type, detail, id)
		}
	}
	println("\nPress any button to go back.")
	var in string
	fmt.Scanln(&in)
	check_exit(in)
	column_choice(db_id, db_name, tbl_id, tbl_name, column_id, column_name)
}

//Utility functions

func check_exit(str string) {
	if str == ":q" {
		termenv.ClearScreen()
		os.Exit(0)
	}
}
