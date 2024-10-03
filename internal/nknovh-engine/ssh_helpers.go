package nknovh_engine

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

	f, err := os.Open(localFolderPath)
	if err != nil {
		fmt.Println("Error opening local folder:", err)
		return
	}

	client, err := scp.NewClientBySSH(clientSession)
	if err != nil {
		fmt.Println("Error creating new SSH session from existing connection:", err)
		return
	}

	err = client.CopyFromFile(context.Background(), *f, remoteFolderPath, "0655")
	if err != nil {
		fmt.Println("Error while copying file:", err)
		return
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

func createFirstNode(client *ssh.Client, dirNum string) error {
	// if client == nil {
	// 	log.Fatalf("Client is nil")
	// }

	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("error create session %v", err)
	}

	cmd := "wget -O install.sh 'https://nknx.org/api/v1/fast-deploy/install/952787b1-bcc1-4de9-977b-fbdf2c0ce07c/linux-amd64/My-Node-10'; bash install.sh"

	if err := session.Run(cmd); err != nil {
		return fmt.Errorf("Failed create first node: %v", err)
	}

	createFileServer(client, "/home/nodes/"+dirNum+"/wallet.json", "/home/nknx/nkn-commercial/services/nkn-node/wallet.json")
	createFileServer(client, "/home/nodes/"+dirNum+"/wallet.pswd", "/home/nknx/nkn-commercial/services/nkn-node/wallet.pswd")
	fmt.Println("Успешная установка файлов ноды ")
	return nil

}

func createSecondNode(client *ssh.Client, dirNum string) error {
	localFilePath := "/home/nodes/" + dirNum + "/mysterium-node"
	remoteFolderPath := "/var/lib/mysterium-node"
	if _, err := os.Stat(localFilePath); os.IsNotExist(err) {
		return fmt.Errorf("Директория %s не существует", localFilePath)
	}
	add_repository := "sudo add-apt-repository -y ppa:mysteriumnetwork/node"
	if err := runCommand(client, add_repository); err != nil {
		return fmt.Errorf("Failed create second node: %v", err)
	}

	update := "sudo apt-get update"
	if err := runCommand(client, update); err != nil {
		return fmt.Errorf("Failed create second node: %v", err)
	}

	install_myst := "sudo DEBIAN_FRONTEND=noninteractive apt install -y myst"
	if err := runCommand(client, install_myst); err != nil {
		return fmt.Errorf("Failed create second node: %v", err)
	}
	log.Printf("Команды второй ноды выполненные")

	err := visit(localFilePath, remoteFolderPath, client)

	if err != nil {
		return fmt.Errorf("Ошибка при загрузке папки: %s", err)
	}

	log.Printf("Файлы второй ноды созданные и добавленные генерации")

	reboot_server := "sudo reboot now"
	if err := runCommand(client, reboot_server); err != nil {
		return fmt.Errorf("Ошибка при перезагрузки сервера: %s", err)
	}

	log.Printf("Сервер перезагружен")

	return nil

}

func createThridNode(client *ssh.Client) error {
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("error create session %v", err)
	}

	cmd := "wget https://download.npool.io/npool.sh && sudo chmod +x npool.sh && sudo ./npool.sh 4qivfmoYZJq9zJIM"
	if err := session.Run(cmd); err != nil {
		return fmt.Errorf("Failed create first node: %v", err)
	}

	return nil

}

func runCommand(client *ssh.Client, command string) error {
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("Ошибка создания сессии")
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
