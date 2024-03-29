package gateway

import (
	"context"
	"errors"
	"github.com/zhaommmmomo/zim/common/config"
	"github.com/zhaommmmomo/zim/common/domain"
	"github.com/zhaommmmomo/zim/common/epoll"
	"github.com/zhaommmmomo/zim/common/log"
	"github.com/zhaommmmomo/zim/common/trace"
	"golang.org/x/sys/unix"
	"io"
	"net"
	"runtime"
	"strconv"
	"time"
)

type reactorManager struct {
	ln          *net.Listener
	mainReactor *reactor
	subReactors []*reactor
	// 负载均衡器
}

type reactor struct {
	name string
	// 获取连接的通道
	connChan chan epoll.Connection
	epoll    *epoll.Epoll
}

const (
	MAX_EPOLL_NUM = 100
)

func initReactor(ln *net.Listener) (*reactorManager, error) {
	// 初始化 reactor manager
	m := &reactorManager{
		ln: ln,
	}
	// 选择负载均衡器

	// 激活 sub reactors
	m.activeSubReactors()
	// 激活 main reactor
	m.activeMainReactor()
	return m, nil
}

func (m *reactorManager) activeSubReactors() {
	// 获取 sub reactor 的数量
	epollNum := config.GetGatewayEpollNum()
	if epollNum <= 0 || epollNum > MAX_EPOLL_NUM {
		epollNum = runtime.NumCPU()
	}
	for i := 0; i < epollNum; i++ {
		r := &reactor{
			name:     "subReactor-" + strconv.Itoa(i),
			connChan: make(chan epoll.Connection),
			epoll:    epoll.NewEpoll(),
		}
		m.subReactors = append(m.subReactors, r)
		go m.runSubReactor(r)
	}
}

func (m *reactorManager) activeMainReactor() {
	// 构建 main reactor
	m.mainReactor = &reactor{
		name: "mainReactor",
	}
	go m.runMainReactor()
}

func (m *reactorManager) runSubReactor(r *reactor) {
	ctx := trace.NewCustomCtxWithTraceId(r.name)
	log.DebugCtx(ctx, "start runSubReactor...")
	go func() {
		for {
			select {
			case conn := <-r.connChan:
				// 将连接添加到 epoll 中
				err := r.epoll.Add(&conn)
				log.DebugCtx(ctx, "gateway add conn to epoll",
					log.Any("remoteAddr", (*conn.Conn).RemoteAddr()))
				if err != nil {
					log.ErrorCtx(ctx, "gateway add conn to epoll has err",
						log.Any("remoteAddr", (*conn.Conn).RemoteAddr()), log.Err(err))
					continue
				}
			}
		}
	}()

	msec := -1
	// 循环 wait 获取数据
	for {
		connectionEvents, err := r.epoll.Wait(config.GetGatewayEpollWaitQueueSize(), msec)
		if err != nil && !errors.Is(err, unix.EINTR) {
			log.ErrorCtx(ctx, "gateway epoll wait error", log.Err(err))
			return
		}
		if len(connectionEvents) <= 0 {
			msec = -1
			continue
		}
		msec = 0
		handEvents(ctx, connectionEvents, r)
	}
}

func handEvents(ctx *context.Context, events []*epoll.ConnectionEvent, r *reactor) {
	for _, event := range events {
		if event.E.Events&epoll.CloseEvent != 0 {
			// 如果是连接断开事件
		}
		if event.E.Events&epoll.ReadEvent != 0 {
			// 如果是读事件
			doHandleReadEvent(ctx, event.C, r)
		}
	}
}

func doHandleReadEvent(ctx *context.Context, c *epoll.Connection, r *reactor) {
	// 设置读取超时
	_ = (*c.Conn).SetReadDeadline(time.Now().Add(time.Duration(120) * time.Second))

	// 读取对应连接中的消息包
	m := readConnData(ctx, c.Conn, r)
	log.DebugCtx(ctx, "gateway read msg", log.Any("data", m))
	// 将读取到的数据通过 work pool 发送到 state server
}

func readConnData(ctx *context.Context, conn *net.Conn, r *reactor) *domain.Message {
	m, err := decoder(ctx, *conn)
	if err != nil {
		if errors.Is(err, io.EOF) {
			// 如果读取数据中途连接断开，清理对应的连接信息
			(*conn).Close()
		}
		return nil
	}
	return m
}

func (m *reactorManager) runMainReactor() {
	ctx := trace.NewCustomCtxWithTraceId(m.mainReactor.name)
	log.DebugCtx(ctx, "start runMainReactor...")
	var i = -1
	for {
		// 不断循环 accept 接口
		// 如果有新的连接，将对应的连接通过负载均衡算法传递到sub reactor的连接channel中
		conn, err := (*m.ln).Accept()
		log.DebugCtx(ctx, "gateway accept conn", log.Any("remoteAddr", conn.RemoteAddr()))
		// 是否限流
		if acceptLimit() {
			_ = conn.Close()
			continue
		}
		if err != nil {
			log.ErrorCtx(ctx, "gateway accept fail", log.Err(err))
			continue
		}
		i++
		if i > 1 {
			i = 0
		}
		// 通知 sub reactor 有新连接建立
		m.subReactors[i].connChan <- epoll.Connection{
			Conn: &conn,
		}
	}
}

func acceptLimit() bool {
	return false
}