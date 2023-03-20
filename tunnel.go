package gosshtunnel

import (
	"encoding/json"
	"io"
	"log"
	"net"
	"os"
	"sync"

	"golang.org/x/crypto/ssh"
)

// username         = "root"
// password         = "password"
// serverAddrString = "192.168.1.100:22"
// localAddrString  = "localhost:9000"
// remoteAddrString = "localhost:9999"

type Config struct {
	Username         string
	ServerAddrString string
	LocalAddrString  string
	RemoteAddrString string
	PrivateKeyPath   string
	HostKeyPath      string
	sshConfig        *ssh.ClientConfig
}

func tunnel(sshClient *ssh.Client, localConn net.Conn, config Config) {
	// Setup the connection to the service running on SSH server
	sshConn, err := sshClient.Dial("tcp", config.RemoteAddrString)
	if err != nil {
		localConn.Close()
		log.Println("WARN:", err)
		return
	}
	wg := sync.WaitGroup{}
	wg.Add(2)
	// Copy local stream to ssh
	go func() {
		_, err = io.Copy(sshConn, localConn)
		if err != nil {
			log.Println("WARN:", err)
		}
		wg.Done()
	}()

	// Copy ssh stream to local
	go func() {
		_, err = io.Copy(localConn, sshConn)
		if err != nil {
			log.Println("WARN:", err)
		}
		wg.Done()
	}()
	wg.Wait()
	sshConn.Close()
}

func GetConfig() (config Config) {
	databytes, err := os.ReadFile("./config.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(databytes, &config)
	if err != nil {
		panic(err)
	}
	log.Println("received config...")
	log.Println("checking for private key in", config.PrivateKeyPath, "...")
	key, err := os.ReadFile(config.PrivateKeyPath)
	if err != nil {
		panic(err)
	}
	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		panic(err)
	}
	log.Println("creating ssh config...")
	config.sshConfig = &ssh.ClientConfig{
		User: config.Username,
		Auth: []ssh.AuthMethod{
			// Use the PublicKeys method for remote authentication.
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	return config
}

func StartTunnel(config Config) {
	// Setup localListener to listen for incoming requests
	log.Println("Tunnel started and listening on", config.LocalAddrString)
	localListener, err := net.Listen("tcp", config.LocalAddrString)
	if err != nil {
		panic(err)
	}
	// Setup client to the SSH server
	sshClient, err := ssh.Dial("tcp", config.ServerAddrString, config.sshConfig)
	if err != nil {
		panic(err)
	}

	for {
		// Setup localConn (type net.Conn)
		localConn, err := localListener.Accept()
		if err != nil {
			panic(err)
		}
		go tunnel(sshClient, localConn, config)
	}
}

// TODO
// func (c Config) TrustedHostKeyCallback() ssh.HostKeyCallback {
// 	key, err := os.ReadFile(c.HostKeyPath)
// 	if err != nil {
// 		panic("unable to read private key: %v", err)
// 	}
// 	return func(_ string, _ net.Addr, k ssh.PublicKey) error {
// 		k.Verify(key, )
// 		return nil
// 	}
// }
