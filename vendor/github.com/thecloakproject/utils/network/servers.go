// Steve Phillips / elimisteve
// 2013.01.13

package network

import (
	"fmt"
	"net"
	"log"
)

var (
	DEBUG = false
)

// TCPServer creates a TCP server to listen for remote connections and
// pass them to the given handler
func TCPServer(listenIPandPort string, maxConns int, handler func(net.Conn)) error {
	// Create TCP connection listener
	tcpAddr, err := net.ResolveTCPAddr("tcp4", listenIPandPort)
	if err != nil {
		return fmt.Errorf("Error calling net.ResolveTCPAddr: " + err.Error())
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return fmt.Errorf("Error calling net.ListenTCP: " + err.Error())
	}

	if DEBUG {
		log.Printf("%s maxConns == %d\n", listenIPandPort, maxConns)
	}

	// Semaphore
	activeConns := make(chan int, maxConns)

	for {
		// Every time someone connects and the number of active connections
		// <= maxConns, handle the connection

		activeConns <- 1
		if DEBUG {
			log.Printf("Added 1 to semaphore. Accepting connections...\n")
		}

		// Accept new connections
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("TCPServer: Error accepting TCP traffic: %v", err)
			<-activeConns
			continue
		}

		log.Printf("* New connection: %s\n\n", conn.RemoteAddr())

		// Handle
		go func() {
			handler(conn)
			if DEBUG {
				log.Printf("handler for %s returned\n", conn.RemoteAddr())
			}
			<-activeConns
		}()
	}
	return nil
}

// TCPServerSimple creates a TCP server to listen for remote
// connections and pass them to the given handler
//
// TODO: Decide whether `handler` should take a net.Conn or
// ReadWriteCloser
func TCPServerSimple(listenIPandPort string, handler func(net.Conn)) error {
	// Create TCP connection listener
	tcpAddr, err := net.ResolveTCPAddr("tcp4", listenIPandPort)
	if err != nil {
		return fmt.Errorf("Error calling net.ResolveTCPAddr: " + err.Error())
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return fmt.Errorf("Error calling net.ListenTCP: " + err.Error())
	}

	//
	for {
		// Accept new connections
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("TCPServer: Error accepting TCP traffic: %v", err)
			continue
		}
		go handler(conn)
	}
	return nil
}
