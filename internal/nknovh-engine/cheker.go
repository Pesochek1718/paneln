package nknovh_engine

import (
	"database/sql"
	"io/ioutil"
	"log"
	"math"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

func chekBusyIp(o *NKNOVH, ip string) bool {
	chekIP := true

	getIPNodes, err := o.sql.db["main"].Query("SELECT * FROM nodes WHERE ip = ?", ip)
	if err != nil {
		panic(err.Error())
	}
	defer getIPNodes.Close()

	if getIPNodes.Next() {
		chekIP = false
	}

	getIPWaitNodes, err := o.sql.db["main"].Query("SELECT * FROM wait_nodes WHERE ip = ?", ip)
	if err != nil {
		panic(err.Error())
	}
	defer getIPWaitNodes.Close()

	if getIPWaitNodes.Next() {
		chekIP = false
	}

	return chekIP
}

func getAllNodes(o *NKNOVH) (*sql.Rows, *sql.Rows) {
	rows, err := o.sql.db["main"].Query("SELECT * FROM nodes")
	if err != nil {
		panic(err.Error())
	}

	rows_wait_nodes, err := o.sql.db["main"].Query("SELECT * FROM wait_nodes")
	if err != nil {
		panic(err.Error())
	}

	return rows, rows_wait_nodes
}
func getDirectories() ([]string, error) {
	directory := "/home/nodes/"
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

func chekNameNodesAdd(o *NKNOVH) string {
	rowsNodes, waitNodes := getAllNodes(o)
	allDirectories, err := getDirectories()
	if err != nil {
		log.Fatalf("Ошибка при получении списка папок: %v", err)
	}
	defer rowsNodes.Close()
	defer waitNodes.Close()

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
		var ssh_key string
		var user string
		var password string
		var done bool
		var wait int
		if err := waitNodes.Scan(&id, &name, &ip, &ssh_key, &user, &password, &done, &wait); err != nil {
			panic(err.Error())
		}

		if _, ok := directoryMap[name]; ok {
			delete(directoryMap, name)
		}
	}

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
