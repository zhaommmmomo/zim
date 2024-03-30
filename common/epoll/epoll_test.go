package epoll

import (
	"fmt"
	"net"
	"testing"
)

func TestSocketFd(t *testing.T) {
	udpConn := net.UDPConn{}
	tcpConn := net.TCPConn{}

	var udpConn1 net.Conn = &udpConn
	var tcpConn1 net.Conn = &tcpConn

	fmt.Println(SocketFd(&udpConn1))
	fmt.Println(SocketFd(&tcpConn1))
}
