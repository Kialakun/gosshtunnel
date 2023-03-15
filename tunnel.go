package gosshtunnel

import (
	"encoding/json"
	"io"
	"log"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
)

// username         = "root"
// password         = "password"
// serverAddrString = "192.168.1.100:22"
// localAddrString  = "localhost:9000"
// remoteAddrString = "localhost:9999"

type Config struct {
	ServerAddrString string
	LocalAddrString  string
	RemoteAddrString string
	PrivateKeyPath   string
	sshConfig        *ssh.ClientConfig
}

func tunnel(localConn net.Conn, config Config) {
	// Setup connection to the SSH server
	sshClientConn, err := ssh.Dial("tcp", config.ServerAddrString, config.sshConfig)
	if err != nil {
		log.Fatalf("ssh.Dial failed: %s", err)
	}

	// Setup the connection to the service running on SSH server
	sshConn, err := sshClientConn.Dial("tcp", config.RemoteAddrString)

	// Copy local stream to ssh
	go func() {
		_, err = io.Copy(sshConn, localConn)
		if err != nil {
			log.Fatalf("io.Copy localConn -> sshConn failed: %v", err)
		}
	}()

	// Copy ssh stream to local
	go func() {
		_, err = io.Copy(localConn, sshConn)
		if err != nil {
			log.Fatalf("io.Copy sshConn -> localConn failed: %v", err)
		}
	}()
}

func GetConfig() (config Config) {
	var hostKey ssh.PublicKey
	databytes, err := os.ReadFile("./config.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(databytes, &config)
	if err != nil {
		panic(err)
	}
	key, err := os.ReadFile(config.PrivateKeyPath)
	if err != nil {
		log.Fatalf("unable to read private key: %v", err)
	}
	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("unable to parse private key: %v", err)
	}
	config.sshConfig = &ssh.ClientConfig{
		User: "user",
		Auth: []ssh.AuthMethod{
			// Use the PublicKeys method for remote authentication.
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.FixedHostKey(hostKey),
	}
	return config
}

func StartTunnel(config Config) {
	// Setup localListener to listen for incoming requests
	localListener, err := net.Listen("tcp", config.LocalAddrString)
	if err != nil {
		log.Fatalf("net.Listen failed: %v", err)
	}

	for {
		// Setup localConn (type net.Conn)
		localConn, err := localListener.Accept()
		if err != nil {
			log.Fatalf("listen.Accept failed: %v", err)
		}
		go tunnel(localConn, config)
	}
}
