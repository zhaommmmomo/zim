package epoll

import (
	"golang.org/x/sys/unix"
	"net"
	"reflect"
	"sync"
)

type Epoll struct {
	Fd          int
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
		Fd: fd,
	}
}

func (e *Epoll) Add(c *Connection) error {
	fd := socketFd(c.Conn)
	err := unix.EpollCtl(e.Fd, unix.EPOLL_CTL_ADD, fd, &unix.EpollEvent{Events: ReadEvent | CloseEvent, Fd: int32(fd)})
	if err != nil {
		return err
	}
	c.Fd = fd
	e.Connections.Store(fd, c)
	return nil
}

func (e *Epoll) Del(c *Connection) error {
	err := unix.EpollCtl(e.Fd, unix.EPOLL_CTL_DEL, c.Fd, nil)
	if err != nil {
		return err
	}
	e.Connections.Delete(c.Fd)
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

func socketFd(conn *net.Conn) int {
	tcpConn := reflect.Indirect(reflect.ValueOf(*conn)).FieldByName("conn")
	fdVal := tcpConn.FieldByName("fd")
	pfdVal := reflect.Indirect(fdVal).FieldByName("pfd")
	return int(pfdVal.FieldByName("Sysfd").Int())
}
