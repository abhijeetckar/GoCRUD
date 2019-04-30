package main

import (
    "database/sql"
    "log"
    "net/http"
	"io/ioutil"
    _ "github.com/go-sql-driver/mysql"
    "encoding/json"
)

type Employee struct {
    Id    int
    Name  string
    City string
}

func dbConn() (db *sql.DB) {
    dbDriver := "mysql"
    dbUser := "root"
    dbPass := "people2018"
    dbName := "shaadi"
    db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
    if err != nil {
        panic(err.Error())
    }
    return db
}


func showAll(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    selDB, err := db.Query("SELECT * FROM employee ORDER BY id DESC")
    if err != nil {
        panic(err.Error())
    }
    emp := Employee{}
    res := []Employee{}
    for selDB.Next() {
        var id int
        var name, city string
        err = selDB.Scan(&id, &name, &city)
        if err != nil {
            panic(err.Error())
        }
        emp.Id = id
        emp.Name = name
        emp.City = city
        res = append(res, emp)
    }
    log.Println(res);
    defer db.Close()
}

func show(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    nId := r.URL.Query().Get("id")
    selDB, err := db.Query("SELECT * FROM Employee WHERE id=?", nId)
    if err != nil {
        panic(err.Error())
    }
    emp := Employee{}
    for selDB.Next() {
        var id int
        var name, city string
        err = selDB.Scan(&id, &name, &city)
        if err != nil {
            panic(err.Error())
        }
        emp.Id = id
        emp.Name = name
        emp.City = city
    }
    log.Println(emp)
    defer db.Close()
}


func insert(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    body,_ := ioutil.ReadAll(r.Body)

    emp := Employee{}
    json.Unmarshal(body, &emp)

    if r.Method == "POST" {
        name := emp.Name
        city := emp.City
        stmtIns, err := db.Prepare("INSERT INTO employee(name, city) VALUES(?,?)")
        if err != nil {
            panic(err.Error())
        }
        stmtIns.Exec(name, city)
        log.Println("INSERT: Name: " + name + " | City: " + city)
    }
    defer db.Close()
}

func update(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    body,_ := ioutil.ReadAll(r.Body)

    emp := Employee{}
    json.Unmarshal(body, &emp)

    if r.Method == "POST" {
        name := emp.Name
        city := emp.City
        id := emp.Id
        stmt, err := db.Prepare("UPDATE employee SET name=?, city=? WHERE id=?")
        if err != nil {
            panic(err.Error())
        }
        stmt.Exec(name, city, id)
        log.Println("UPDATE: Name: " + name + " | City: " + city)
    }
    defer db.Close()
}

func delete(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    emp := r.URL.Query().Get("id")
    stmt, err := db.Prepare("DELETE FROM employee WHERE id=?")
    if err != nil {
        panic(err.Error())
    }
    stmt.Exec(emp)
    log.Println("DELETE")
    defer db.Close()
}

func main() {
    log.Println("Server started on: http://localhost:8080")
    http.HandleFunc("/", showAll)
    http.HandleFunc("/show", show)
    http.HandleFunc("/insert", insert)
    http.HandleFunc("/update", update)
    http.HandleFunc("/delete", delete)
    http.ListenAndServe(":8080", nil)
}