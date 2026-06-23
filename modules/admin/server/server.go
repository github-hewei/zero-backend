package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/241x/zero-kit/logger"
	"github.com/241x/zero-web/config"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// HTTPServer HTTP 服务
type HTTPServer struct {
	httpServ *http.Server
	logger   logger.Logger
	handler  *gin.Engine
	db       *gorm.DB
}

// NewHTTPServer 创建 HTTP 服务
func NewHTTPServer(cfg config.ServerConfig, handler *gin.Engine, logger logger.Logger, db *gorm.DB) *HTTPServer {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	httpServ := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	return &HTTPServer{
		httpServ: httpServ,
		logger:   logger,
		handler:  handler,
		db:       db,
	}
}

// Run 启动 HTTP 服务
func (s *HTTPServer) Run() {
	sig := make(chan os.Signal, 1)

	// SIGINT: `ctrl+c`, SIGTERM: `kill` Command
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		if err := s.httpServ.ListenAndServe(); err != nil {
			s.logger.Err(err, "Listen And Serve")
		}
	}()

	s.logger.Info("Server Start", "listen", s.httpServ.Addr)

	<-sig
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := s.httpServ.Shutdown(ctx); err != nil {
		s.logger.Err(err, "Server Shutdown")
	}
}
