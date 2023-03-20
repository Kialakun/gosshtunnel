package main

import tun "github.com/Kialakun/gosshtunnel"

func main() {
	// retreive configuration from file
	config := tun.GetConfig()
	// start tunnelling
	tun.StartTunnel(config)
}
