package main

import (
	"context"
	"errors"
	"log"
	"net"

	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/go-mysql-org/go-mysql/server"
)

type AuthMethod struct {
	Username string `toml:"username"`
	Passwd   string `toml:"password"`
}
type MysqlProxyServer struct {
	ctx    context.Context
	cancel context.CancelFunc

	Host  string       `toml:"host"`
	Auths []AuthMethod `toml:"auths"`

	h server.Handler `toml:"-"`
}

func NewMysqlServer(h server.Handler) *MysqlProxyServer {
	ms := MysqlProxyServer{
		Host: "localhost:4000",
		h:    h,
		Auths: []AuthMethod{
			{Username: "root", Passwd: "password"},
		},
	}
	return &ms
}

func (s *MysqlProxyServer) Start(ctx context.Context) error {
	s.ctx, s.cancel = context.WithCancel(ctx)

	l, err := net.Listen("tcp", s.Host)
	if err != nil {
		panic(err)
	}

	provider := server.NewInMemoryProvider()
	if len(s.Auths) == 0 {
		s.Auths = append(s.Auths, AuthMethod{
			Username: "root",
			Passwd:   "",
		})
	}
	for _, auth := range s.Auths {
		provider.AddUser(auth.Username, auth.Passwd)
	}

	for {
		select {
		case <-s.ctx.Done():
			return s.ctx.Err()
		default:
			c, err := l.Accept()
			if err != nil {
				if err == context.Canceled {
					return s.ctx.Err()
				}
				log.Printf("Error accepting connection err:%v\n", err)
				continue
			}

			conn, err := server.NewCustomizedConn(c, server.NewDefaultServer(), provider, s.h)
			if err != nil {
				log.Printf("Error creating connection err:%v\n", err)
				continue
			}

			go s.handleCommand(conn)
		}
	}
}

func (s *MysqlProxyServer) handleCommand(conn *server.Conn) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Handler panic: %v\n", err)
			return
		}
	}()

	for {
		select {
		case <-s.ctx.Done():
			return
		default:
			if err := conn.HandleCommand(); err != nil {
				if !errors.Is(err, mysql.ErrBadConn) {
					log.Printf("Error handling command: %v\n", err)
				}
				return
			}
		}
	}
}
