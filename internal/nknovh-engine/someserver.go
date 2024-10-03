package nknovh_engine

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"time"

	"golang.org/x/crypto/ssh"
)

func stringToBytes(str string) []byte {
	return []byte(str)
}
func someServer(host string, user string, password string, proxyAddress string, sleepHours int, dirNum string, q *WSQuery, c *CLIENT, o *NKNOVH, keySsh string, id int) (err error, r WSReply) {
	// _, errAddNodes := o.sql.db["main"].Exec("INSERT INTO wait_nodes (name, ip) VALUES (?, ?)", dirNum, host)
	// if errAddNodes != nil {
	// 	return errAddNodes, WSReply{}
	// }

	// sleepDuration := time.Duration(sleepHours) * time.Minute
	// time.Sleep(sleepDuration)

	client, err := createClient(host, user, password, proxyAddress, keySsh)
	if err != nil {
		// log.Printf("Ошибка при создании сессии: %s", err)
		// ip, port := getIP()
		// log.Printf("Новый прокси: %s", ip+":"+port)
		// someServer(host, user, password, ip+":"+port, sleepHours, dirNum, q, c, o)
		fmt.Println("Ошибка клиента", host)
		return err, WSReply{}
	}
	if client != nil {
		defer client.Close()
	}

	fmt.Println("Установка первой ноды ", host)
	timer := time.NewTimer(2000 * time.Second)
	defer timer.Stop()

	errCh := make(chan error, 1)
	go func() {
		errCh <- createFirstNode(client, dirNum)
	}()

	select {
	case err := <-errCh:
		if err != nil {
			fmt.Println("Ошибка createFirstNode", host)
			return err, WSReply{}
		}
	case <-timer.C:
		fmt.Println("Превышено время ожидания для createFirstNode ", host)
		return nil, WSReply{}
	}
	// err = createSecondNode(client, dirNum)
	// if err != nil {
	// 	deleteWaitNodes(o, id)
	// 	return err, WSReply{}
	// }

	q.Value = make(map[string]interface{})
	q.Value["Ip"] = host
	q.Value["Multiple"] = false
	q.Value["Name"] = dirNum
	fmt.Println("Установка apiAddNodes", host)
	if err, reply := o.apiAddNodes(q, c); err != nil {
		fmt.Println("Ошибка apiAddNodes", host)
		return err, WSReply{}
	} else {
		return nil, reply
	}

}

func deleteWaitNodes(o *NKNOVH, id int) {
	_, errDeleteWait := o.sql.db["main"].Exec("DELETE FROM wait_nodes WHERE id = ?", id)
	if errDeleteWait != nil {
		log.Printf("Ошибка удаления")
	}
}

func createClient(host string, user string, password string, proxyAddress string, keySsh string) (*ssh.Client, error) {
	var sshConfig *ssh.ClientConfig
	if keySsh == "" {
		sshConfig = &ssh.ClientConfig{
			User: user,
			Auth: []ssh.AuthMethod{
				ssh.Password(password),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		}
	} else {
		byteSshKey := stringToBytes(keySsh)
		signer, err := ssh.ParsePrivateKey(byteSshKey)
		if err != nil {
			return nil, err
		}
		sshConfig = &ssh.ClientConfig{
			User: user,
			Auth: []ssh.AuthMethod{
				ssh.PublicKeys(signer),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		}
	}

	client, err := proxiedSSHClient(proxyAddress, host+":22", sshConfig)
	if err != nil {
		var client *ssh.Client
		var lastErr error
		for attempt := 0; attempt < 2; attempt++ {
			time.Sleep(time.Minute)
			fmt.Println("Подключения через новый прокси к ", host)
			ip, port := getIP()
			proxyAddress := ip + ":" + port
			client, lastErr = proxiedSSHClient(proxyAddress, host+":22", sshConfig)
			if lastErr == nil {
				return client, nil
			} else {
				continue
			}
		}

		return nil, fmt.Errorf("Ошибка при создании клиента: %s", err)
	}

	return client, nil
}

func visit(localFolderPath string, remoteFolderPath string, client *ssh.Client) error {
	session, err := client.NewSession()
	if err != nil {
		client.Close()
		return nil
	}

	defer session.Close()
	files, err := ioutil.ReadDir(localFolderPath)
	if err != nil {
		return err
	}
	createDirectoryServer(session, remoteFolderPath)
	for _, file := range files {
		filePath := filepath.Join(localFolderPath, file.Name())
		remoteFilePath := filepath.Join(remoteFolderPath, file.Name())
		remoteFilePath = filepath.ToSlash(remoteFilePath)

		if file.IsDir() {
			err = visit(filePath, remoteFilePath, client)
			if err != nil {
				return err
			}
		} else {
			createFileServer(client, filePath, remoteFilePath)
		}
	}
	return nil
}

func createThirdScript(host string, user string, password string, proxyAddress string, sleepHours int, o *NKNOVH, keySsh string, id int) (err error, r WSReply) {

	// _, errAddNodes := o.sql.db["main"].Exec("INSERT INTO wait_nodes (ip, name) VALUES (?, ?)", host, "None")
	// if errAddNodes != nil {
	// 	return errAddNodes, WSReply{}
	// }
	fmt.Println("Создания третего скрипта: ", host)

	client, err := createClient(host, user, password, proxyAddress, keySsh)
	if err != nil {
		fmt.Println("Ошибка ", host)
		return err, WSReply{}
	}
	if client != nil {
		fmt.Println("Ошибка ", host)
		defer client.Close()
	}

	err = createThridNode(client)
	if err != nil {
		fmt.Println("Ошибка ", host)
		return err, WSReply{}
	}
	return err, WSReply{}

}

func chekNodesCreate(q *WSQuery, c *CLIENT, o *NKNOVH) {
	fmt.Println("Start")
	for {
		rows_wait_nodes, err := o.sql.db["main"].Query("SELECT * FROM wait_nodes WHERE done = 1 ORDER BY id LIMIT 1")
		if err != nil {
			fmt.Println("Error querying wait_nodes:", err)
			time.Sleep(time.Minute)
			continue
		}
		defer rows_wait_nodes.Close()

		for rows_wait_nodes.Next() {
			var id int
			var name string
			var ip string
			var ssh_key string
			var user string
			var password string
			var done bool
			var wait int
			err := rows_wait_nodes.Scan(&id, &name, &ip, &ssh_key, &user, &password, &done, &wait)
			if err != nil {
				deleteWaitNodes(o, id)
				fmt.Println("Error scanning row:", err)
				continue
			}
			fmt.Println("Start: ", ip)
			if done == true {
				ipProxy, portProxy := getIP()
				if name == "None" {
					fmt.Println("createThirdScript: ", ip)
					createThirdScript(ip, user, password, ipProxy+":"+portProxy, wait, o, ssh_key, id)
				} else {
					fmt.Println("someServer: ", ip)
					someServer(ip, user, password, ipProxy+":"+portProxy, wait, name, q, c, o, ssh_key, id)
				}
				fmt.Println("Успешно добавлен сервер: ", ip)
				deleteWaitNodes(o, id)
			}
		}
		if err := rows_wait_nodes.Err(); err != nil {
			fmt.Println("Error iterating through rows:", err)
		}
		time.Sleep(time.Minute)
	}
}

func addNewCreateServer(host string, user string, password string, dirNum string, keySsh string, o *NKNOVH, wait int) (int64, error) {
	result, errAddNodes := o.sql.db["main"].Exec("INSERT INTO wait_nodes (name, ip, ssh_key, user, password, done, wait) VALUES (?, ?, ?, ?, ?, ?, ?)", dirNum, host, keySsh, user, password, false, wait)
	if errAddNodes != nil {
		return 0, errAddNodes
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func workingCreateserver(host string, user string, password string, dirNum string, keySsh string, o *NKNOVH, wait int) {
	id, err := addNewCreateServer(host, user, password, dirNum, keySsh, o, wait)
	if err != nil {
		fmt.Println("Error adding server:", err)
	}
	fmt.Println("Ожидания")
	sleepDuration := time.Duration(wait) * time.Minute
	time.Sleep(sleepDuration)
	_, err = o.sql.db["main"].Exec("UPDATE wait_nodes SET done = ? WHERE id = ?", true, id)
	if err != nil {
		fmt.Println("Error updating done status:", err)
	}
	fmt.Println("Завершения ожидания")
}
