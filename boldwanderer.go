package boldwanderer

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net"
	"os"
	"strconv"
	"strings"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

type Options struct {
	LogFormat string
	LogLevel  string
}

func Execute() int {
	options := parseArgs()

	logger, err := getLogger(options.LogLevel, options.LogFormat)
	if err != nil {
		slog.Error("getLogger", "error", err)
		return 1
	}

	slog.SetDefault(logger)

	err = run(options)
	if err != nil {
		slog.Error("run failed", "error", err)
		return 1
	}
	return 0
}

func parseArgs() Options {
	options := Options{}

	flag.StringVar(&options.LogLevel, "log-level", "info", "Log level (debug, info, warn, error), default: info")
	flag.StringVar(&options.LogFormat, "log-format", "", "Log format (text or json)")

	flag.Parse()

	return options
}

func run(options Options) error {
	var connStr string
	// Check if the positional argument is provided
	args := flag.Args()
	if len(args) > 0 {
		// Use the first positional argument as the connection string
		connStr = args[0]
	}

	// Check if the connection string is provided
	if connStr == "" {
		fmt.Println("Connection string not provided. Please use positional argument 'username@hostname:port'.")
	}

	// Split the connection string into username, hostname, and optional port
	parts := strings.Split(connStr, "@")
	if len(parts) != 2 {
		fmt.Println("Invalid connection string format. Please use 'username@hostname:port'.")
	}

	if len(parts) < 2 {
		return fmt.Errorf("can't parse host")
	}

	user := parts[0]
	hostWithPort := parts[1]

	var host string
	var port int

	// Split the host string to separate hostname and port
	hostParts := strings.Split(hostWithPort, ":")
	host = hostParts[0]

	// Convert string to integer using Atoi function
	port, err := strconv.Atoi(hostParts[1])
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	sshAgentConn, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK"))
	if err != nil {
		log.Println("Failed to connect to SSH agent:", err)
		return err
	}
	defer sshAgentConn.Close()

	agentClient := agent.NewClient(sshAgentConn)

	hostkeyCallback := ssh.InsecureIgnoreHostKey()

	conf := &ssh.ClientConfig{
		User:            user,
		HostKeyCallback: hostkeyCallback,
		Auth: []ssh.AuthMethod{
			// Use the private keys from the SSH agent for authentication
			ssh.PublicKeysCallback(agentClient.Signers),
		},
	}

	slog.Debug("stats", "user", user, "host", host, "port", port)

	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", host, port), conf)
	if err != nil {
		log.Println("Failed to establish SSH connection:", err)
		return err
	}
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil {
		log.Println("Failed to create SSH session:", err)
		return err
	}
	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b

	err = session.Run("ls")
	if err != nil {
		return err
	}

	fmt.Println(b.String())

	return nil
}
