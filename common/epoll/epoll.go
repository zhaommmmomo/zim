package epoll

import (
	"golang.org/x/sys/unix"
	"net"
	"reflect"
	"sync"
	"sync/atomic"
)

type Epoll struct {
	Fd          int
	ConnCount   int32
	Connections sync.Map
}

type Connection struct {
	Fd   int
	Conn *net.Conn
}

type ConnectionEvent struct {
	C *Connection
	E *unix.EpollEvent
}

const (
	ReadEvent       = unix.EPOLLIN
	WriteEvent      = unix.EPOLLOUT
	CloseEvent      = unix.EPOLLHUP
	ReadWriteEvents = ReadEvent | WriteEvent
)

func NewEpoll() *Epoll {
	fd, err := unix.EpollCreate1(0)
	if err != nil {
		panic(err)
	}
	return &Epoll{
		Fd:          fd,
		ConnCount:   0,
		Connections: sync.Map{},
	}
}

func (e *Epoll) Add(c *Connection) error {
	c.Fd = socketFd(c.Conn)
	err := unix.EpollCtl(e.Fd, unix.EPOLL_CTL_ADD, c.Fd, &unix.EpollEvent{Events: ReadEvent | CloseEvent, Fd: int32(c.Fd)})
	if err != nil {
		return err
	}
	e.AddConn(c.Fd, c)
	return nil
}

func (e *Epoll) Del(c *Connection) error {
	e.DelConn(c.Fd)
	err := unix.EpollCtl(e.Fd, unix.EPOLL_CTL_DEL, c.Fd, nil)
	if err != nil {
		return err
	}
	return nil
}

func (e *Epoll) Wait(size, msec int) ([]*ConnectionEvent, error) {
	events := make([]unix.EpollEvent, size)
	n, err := unix.EpollWait(e.Fd, events, msec)
	if err != nil {
		return nil, err
	}
	var connections []*ConnectionEvent
	for i := 0; i < n; i++ {
		event := events[i]
		conn, _ := e.Connections.Load(int(event.Fd))
		connections = append(connections, &ConnectionEvent{
			E: &event,
			C: conn.(*Connection),
		})
	}
	return connections, nil
}

func (e *Epoll) AddConn(fd int, c *Connection) {
	atomic.AddInt32(&e.ConnCount, 1)
	e.Connections.Store(fd, c)
}

func (e *Epoll) DelConn(fd int) {
	if _, ok := e.Connections.LoadAndDelete(fd); ok {
		atomic.AddInt32(&e.ConnCount, -1)
	}
}

func socketFd(conn *net.Conn) int {
	tcpConn := reflect.Indirect(reflect.ValueOf(*conn)).FieldByName("conn")
	fdVal := tcpConn.FieldByName("fd")
	pfdVal := reflect.Indirect(fdVal).FieldByName("pfd")
	return int(pfdVal.FieldByName("Sysfd").Int())
}
