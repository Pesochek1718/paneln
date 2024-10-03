package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

func chekBusyIp(ip string) bool {
	chekIP := true
	db, err := sql.Open("mysql", "caesar:caesar@tcp(127.0.0.1:3306)/nkn")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	getIPNodes, err := db.Query("SELECT * FROM nodes WHERE ip = ?", ip)
	if err != nil {
		panic(err.Error())
	}
	defer getIPNodes.Close()

	if getIPNodes.Next() {
		chekIP = false
	}

	getIPWaitNodes, err := db.Query("SELECT * FROM wait_nodes WHERE ip = ?", ip)
	if err != nil {
		panic(err.Error())
	}
	defer getIPWaitNodes.Close()

	if getIPWaitNodes.Next() {
		chekIP = false
	}

	return chekIP
}

func getAllNodes() (*sql.Rows, *sql.Rows) {
	db, err := sql.Open("mysql", "caesar:caesar@tcp(127.0.0.1:3306)/nkn")

	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	rows, err := db.Query("SELECT * FROM nodes")
	if err != nil {
		panic(err.Error())
	}

	rows_wait_nodes, err := db.Query("SELECT * FROM wait_nodes")
	if err != nil {
		panic(err.Error())
	}

	return rows, rows_wait_nodes
}
func getDirectories() ([]string, error) {
	directory := "../nodes"
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	var directories []string
	for _, file := range files {
		if file.IsDir() {
			directories = append(directories, file.Name())
		}
	}

	return directories, nil
}

func chekNameNodesAdd() string {
	rowsNodes, waitNodes := getAllNodes()
	allDirectories, err := getDirectories()
	if err != nil {
		log.Fatalf("Ошибка при получении списка папок: %v", err)
	}
	defer rowsNodes.Close()

	directoryMap := make(map[string]bool)

	for _, dir := range allDirectories {
		directoryMap[dir] = true
	}

	for rowsNodes.Next() {
		var id int
		var name string
		var hashID string
		var created string
		var dirty string
		var dirtyFCnt string
		var ip string
		if err := rowsNodes.Scan(&id, &hashID, &name, &ip, &created, &dirty, &dirtyFCnt); err != nil {
			panic(err.Error())
		}

		if _, ok := directoryMap[name]; ok {
			delete(directoryMap, name)
		}
	}

	for waitNodes.Next() {
		var id int
		var name string
		var ip string
		if err := waitNodes.Scan(&id, &name, &ip); err != nil {
			panic(err.Error())
		}

		if _, ok := directoryMap[name]; ok {
			delete(directoryMap, name)
		}
	}
	fmt.Println(directoryMap)
	newDir := getSmallestKey(directoryMap)

	return newDir

}

func getSmallestKey(m map[string]bool) string {
	smallest := math.MaxInt64
	smallestKey := ""
	for key := range m {
		val, err := strconv.Atoi(key)
		if err != nil {
			continue
		}
		if val < smallest {
			smallest = val
			smallestKey = key
		}
	}
	return smallestKey
}
