package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"time"

	"golang.org/x/crypto/ssh"
)

func someServer(host string, user string, password string, proxyAddress string, sleepHours int) {
	log.Printf("Создания сервера")
	client, err := createClient(host, user, password, proxyAddress)
	if err != nil {
		log.Printf("Ошибка при создании сессии: %s", err)
		ip, port := getIP()
		log.Printf("Новый прокси: %s", ip+":"+port)
		someServer(host, user, password, ip+":"+port, sleepHours)
	}
	defer client.Close()

	log.Printf("Сервер готов к работе, ожидаем")

	seconds := 0
	off_second := 10

	for {
		seconds++
		fmt.Println(seconds)
		time.Sleep(time.Second)
		if seconds >= off_second {
			break
		}
	}

	createFirstNode(client)
	createSecondNode(client)

}

func createClient(host string, user string, password string, proxyAddress string) (*ssh.Client, error) {
	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := proxiedSSHClient(proxyAddress, host+":22", sshConfig)
	if err != nil {
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
