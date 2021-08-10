package http

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"example/core"
)

type Options struct {
	Addr string
	Core *core.Core
	Log  *log.Logger
}

type Server struct {
	addr string
	echo *echo.Echo
	core *core.Core
	log  *log.Logger
}

type operation struct {
	ID        uint64         `json:"id"`
	Time      time.Time      `json:"time"`
	Direction core.Direction `json:"dir"`
	Amount    uint64         `json:"amount"`
}

type response struct {
	OK      bool        `json:"ok"`
	Payload interface{} `json:"payload,omitempty"`
}

func New(options *Options) *Server {
	s := &Server{
		addr: options.Addr,
		core: options.Core,
		echo: echo.New(),
		log:  options.Log,
	}
	s.echo.HTTPErrorHandler = func(err error, c echo.Context) {
		s.replyErr(c, err)
	}
	s.echo.POST("/add_wallet", s.handleAddWallet)
	s.echo.POST("/deposit_money", s.handleDepositMoney)
	s.echo.POST("/transfer_money", s.handleTransferMoney)
	s.echo.GET("/get_operations", s.handleGetOperations)
	return s
}

func (s *Server) Serve() error {
	return s.echo.Start(s.addr)
}

func (s *Server) replyErr(c echo.Context, err error) error {
	s.log.Printf("error: %v", err)
	return c.JSON(http.StatusOK, response{OK: false})
}

func (s *Server) replyOK(c echo.Context, payload interface{}) error {
	return c.JSON(http.StatusOK, response{OK: true, Payload: payload})
}

func (s *Server) handleAddWallet(c echo.Context) error {
	ctx := c.Request().Context()
	name, err := s.getName(c, "name")
	if err != nil {
		return s.replyErr(c, err)
	}
	id, err := s.core.AddWallet(ctx, name)
	if err != nil {
		s.log.Printf("error: %v", err)
		return s.replyErr(c, err)
	}
	return s.replyOK(c, id)
}

func (s *Server) handleDepositMoney(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := s.getUint(c, "id")
	if err != nil {
		return s.replyErr(c, err)
	}
	amount, err := s.getUint(c, "amount")
	if err != nil || amount == 0 {
		return s.replyErr(c, err)
	}
	err = s.core.DepositMoney(ctx, id, amount)
	if err != nil {
		return s.replyErr(c, err)
	}
	return s.replyOK(c, nil)
}

func (s *Server) handleTransferMoney(c echo.Context) error {
	ctx := c.Request().Context()
	from, err := s.getUint(c, "from")
	if err != nil {
		return s.replyErr(c, err)
	}
	to, err := s.getUint(c, "to")
	if err != nil {
		return s.replyErr(c, err)
	}
	if from == to {
		return s.replyErr(c, fmt.Errorf("%w: 'to' and 'from' equal", errParam))
	}
	amount, err := s.getUint(c, "amount")
	if err != nil || amount == 0 {
		return s.replyErr(c, err)
	}
	err = s.core.TransferMoney(ctx, from, to, amount)
	if err != nil {
		return s.replyErr(c, err)
	}
	return s.replyOK(c, nil)
}

func (s *Server) handleGetOperations(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := s.getUint(c, "id")
	if err != nil {
		return s.replyErr(c, err)
	}
	from, err := s.getTime(c, "from", time.Now().Add(-24*time.Hour))
	if err != nil {
		return s.replyErr(c, err)
	}
	to, err := s.getTime(c, "to", time.Now())
	if err != nil {
		return s.replyErr(c, err)
	}
	if to.Before(from) {
		return s.replyErr(c, fmt.Errorf("%w: 'to' is before 'from'", errParam))
	}
	dir, err := s.getDirection(c, "dir")
	if err != nil {
		return s.replyErr(c, err)
	}
	ops, err := s.core.OperationsHistory(ctx, id, from, to, dir)
	if err != nil {
		return s.replyErr(c, err)
	}
	operations := make([]operation, len(ops))
	for i, op := range ops {
		operations[i] = operation(op)
	}
	return s.replyOK(c, operations)
}
