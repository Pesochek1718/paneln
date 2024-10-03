package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	scp "github.com/bramvdbogaerde/go-scp"
	"golang.org/x/crypto/ssh"
	"golang.org/x/net/proxy"
)

func uploadDir(session *ssh.Client, localFolderPath string, remoteFolderPath string, sessionServer *ssh.Session) error {
	f, _ := os.Open(localFolderPath)
	client, err := scp.NewClientBySSH(session)
	if err != nil {
		fmt.Println("Error creating new SSH session from existing connection", err)
	}

	err = client.CopyFromFile(context.Background(), *f, remoteFolderPath, "0655")

	if err != nil {
		fmt.Println("Error while copying file ", err)
	}

	return nil
}

func createDirectoryServer(sessionServer *ssh.Session, directoryPath string) {
	cmd := "mkdir -p " + directoryPath

	if err := sessionServer.Run(cmd); err != nil {
		log.Fatalf("Failed to run mkdir: %v", err)
	}
}

func createFileServer(clientSession *ssh.Client, localFolderPath string, remoteFolderPath string) {
	f, _ := os.Open(localFolderPath)
	client, err := scp.NewClientBySSH(clientSession)
	if err != nil {
		fmt.Println("Error creating new SSH session from existing connection", err)
	}
	err = client.CopyFromFile(context.Background(), *f, remoteFolderPath, "0655")

	if err != nil {
		fmt.Println("Error while copying file ", err)
	}

}

func proxiedSSHClient(proxyAddress, sshServerAddress string, sshConfig *ssh.ClientConfig) (*ssh.Client, error) {
	dialer, err := proxy.SOCKS5("tcp", proxyAddress, nil, proxy.Direct)
	if err != nil {
		return nil, fmt.Errorf("Ошибка при создании прокси: %s", err)
	}

	conn, err := dialer.Dial("tcp", sshServerAddress)
	if err != nil {
		return nil, fmt.Errorf("Ошибка при подключении к серверу SSH: %s", err)
	}

	c, chans, reqs, err := ssh.NewClientConn(conn, sshServerAddress, sshConfig)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("Ошибка при создании SSH-соединения: %s", err)
	}

	return ssh.NewClient(c, chans, reqs), nil
}

func publicKeyFile(file string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("Ошибка при чтении приватного ключа: %s", err)
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		log.Fatalf("Ошибка при парсинге приватного ключа: %s", err)
	}
	return ssh.PublicKeys(key)
}

func createFirstNode(client *ssh.Client) {
	// if client == nil {
	// 	log.Fatalf("Client is nil")
	// }

	session, err := client.NewSession()
	if err != nil {
		client.Close()
		log.Fatalf("Ошибка создания сессии")
	}

	cmd := "wget -O install.sh 'https://nknx.org/api/v1/fast-deploy/install/72ecb02d-e6ac-4e31-a107-7c045fff4f8a/linux-amd64/My-Node-1'; bash install.sh"

	if err := session.Run(cmd); err != nil {
		log.Fatalf("Failed create first node: %v", err)
	}

	createFileServer(client, "D:/project/nkn.ovh/test/test/86/wallet.json", "/home/nknx/nkn-commercial/services/nkn-node/wallet.json")
	createFileServer(client, "D:/project/nkn.ovh/test/test/86/wallet.pswd", "/home/nknx/nkn-commercial/services/nkn-node/wallet.pswd")

	fmt.Println("Первая нода созданна и файлы генерации загруженны")

}

func createSecondNode(client *ssh.Client) {

	add_repository := "sudo add-apt-repository -y ppa:mysteriumnetwork/node"
	if err := runCommand(client, add_repository); err != nil {
		log.Fatalf("Failed create second node: %v", err)
	}

	update := "sudo apt-get update"
	if err := runCommand(client, update); err != nil {
		log.Fatalf("Failed create second node: %v", err)
	}

	install_myst := "sudo DEBIAN_FRONTEND=noninteractive apt install -y myst"
	if err := runCommand(client, install_myst); err != nil {
		log.Fatalf("Failed create second node: %v", err)
	}
	log.Printf("Команды второй ноды выполненные")
	localFilePath := "D:/project/nkn.ovh/test/test/86/mysterium-node"
	remoteFolderPath := "/var/lib/mysterium-node"

	err := visit(localFilePath, remoteFolderPath, client)

	if err != nil {
		log.Fatalf("Ошибка при загрузке папки: %s", err)
	}

	log.Printf("Файлы второй ноды созданные и добавленные генерации")

	reboot_server := "sudo reboot now"
	if err := runCommand(client, reboot_server); err != nil {
		log.Fatalf("Ошибка при перезагрузки сервера: %s", err)
	}

	log.Printf("Сервер перезагружен")

}
func runCommand(client *ssh.Client, command string) error {
	session, err := client.NewSession()
	if err != nil {
		client.Close()
		log.Fatalf("Ошибка создания сессии")
	}
	defer session.Close()

	var stdoutBuf, stderrBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	session.Stderr = &stderrBuf

	if err := session.Run(command); err != nil {
		return fmt.Errorf("failed to run command '%s': %v, stderr: %s", command, err, stderrBuf.String())
	}
	return nil
}
