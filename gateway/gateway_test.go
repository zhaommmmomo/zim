package gateway

import (
	"encoding/binary"
	"fmt"
	"github.com/zhaommmmomo/zim/common/domain"
	"net"
	"os"
	"testing"
	"time"
)

func TestGatewayConn1(t *testing.T) {
	serverAddr := "127.0.0.1:8002"
	// 建立 TCP 连接
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Printf("Failed to connect to server: %s\n", err)
		os.Exit(1)
	}
	fmt.Println("---------connected 1---------")
	defer conn.Close()
	time.Sleep(time.Second * 30)
}

func TestGatewayConn2(t *testing.T) {
	serverAddr := "127.0.0.1:8002" // 服务器地址

	// 连接服务器
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Printf("Failed to connect to server: %s\n", err)
		return
	}
	defer conn.Close()

	for i := 0; i < 4; i++ {
		if i != 3 {
			time.Sleep(time.Second * 2)
		}
		vHeader := []byte(fmt.Sprintf("Variable Header %d", i))
		payload := []byte(fmt.Sprintf("Payload Data %d", i))
		// 构造消息
		message := &domain.Message{
			FHeader: &domain.FixedHeader{
				V:          1,
				Cmd:        2,
				VarHLen:    uint32(len(vHeader)),
				PayloadLen: uint32(len(payload)),
				Crc32sum:   0,
			},
			VHeader: vHeader,
			Payload: payload,
		}

		// 发送消息
		if err := sendMessage(conn, message); err != nil {
			fmt.Printf("Failed to send message: %s\n", err)
			return
		}
		fmt.Println("Message sent successfully!")
	}
	time.Sleep(time.Second * 30)
}

func TestServer(t *testing.T) {
	serverAddr := "0.0.0.0:8002" // 服务器监听地址
	// 监听指定地址的TCP连接
	listener, err := net.Listen("tcp", serverAddr)
	if err != nil {
		fmt.Printf("Failed to start server: %s\n", err)
		return
	}

	fmt.Println("Server started, waiting for connections...")

	// 接收并处理客户端连接
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Failed to accept connection: %s\n", err)
			continue
		}

		// 在新的 goroutine 中处理连接
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Printf("handleConnection: %v\n", conn.RemoteAddr())
	// 读取并解析消息
	p, err := readMessage(conn)
	if err != nil {
		fmt.Printf("Failed to read message: %s\n", err)
		return
	}

	// 打印消息
	fmt.Println("Received message:")
	fmt.Printf("FixedHeader: %+v\n", p.FHeader)
	fmt.Printf("VariableHeader: %s\n", string(p.VHeader))
	fmt.Printf("Payload: %s\n", string(p.Payload))
}

func readMessage(conn net.Conn) (*domain.Message, error) {
	// 读取消息头部
	headerBuf := make([]byte, 14)
	if _, err := conn.Read(headerBuf); err != nil {
		return nil, fmt.Errorf("failed to read message header: %s", err)
	}

	// 解析消息头部
	varHLen := binary.BigEndian.Uint32(headerBuf[2:6])
	payloadLen := binary.BigEndian.Uint32(headerBuf[6:10])

	// 读取变长头部
	varHeader := make([]byte, varHLen)
	if _, err := conn.Read(varHeader); err != nil {
		return nil, fmt.Errorf("failed to read variable header: %s", err)
	}

	// 读取载荷
	payload := make([]byte, payloadLen)
	if _, err := conn.Read(payload); err != nil {
		return nil, fmt.Errorf("failed to read payload: %s", err)
	}

	// 构造 Message 对象
	p := &domain.Message{
		FHeader: &domain.FixedHeader{
			V:          headerBuf[0],
			Cmd:        headerBuf[1],
			VarHLen:    varHLen,
			PayloadLen: payloadLen,
			Crc32sum:   binary.BigEndian.Uint32(headerBuf[10:14]),
		},
		VHeader: varHeader,
		Payload: payload,
	}

	return p, nil
}

func TestGatewayConn3(t *testing.T) {
	serverAddr := "127.0.0.1:8002"
	// 建立 TCP 连接
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Printf("Failed to connect to server: %s\n", err)
		os.Exit(1)
	}
	fmt.Println("---------connected 3---------")
	defer conn.Close()
	time.Sleep(time.Second * 300)
}

func sendMessage(conn net.Conn, message *domain.Message) error {
	if _, err := conn.Write(encoder(message)); err != nil {
		return fmt.Errorf("failed to send message: %s", err)
	}
	return nil
}
