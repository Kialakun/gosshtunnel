# gosshtunnel – Simple TCP Tunnel over SSH

`gosshtunnel` is a lightweight Go package and CLI tool that forwards local TCP traffic through an SSH server to a remote target. It reads a JSON configuration file, establishes an SSH connection using a private key, and listens on a local port – forwarding each incoming connection over the SSH tunnel.

---

## Features

- **SSH tunnelling** – forward any TCP port from local machine to a remote host via an SSH server.
- **Key‑based authentication** – uses a private key file (PEM format) for secure SSH login.
- **JSON configuration** – simple, human‑readable config file.
- **Concurrent connections** – each incoming connection is handled in its own goroutine.
- **Bidirectional streaming** – copies data in both directions between local and remote.
