package http

import (
	"context"

	"github.com/gin-gonic/gin"
)

type HTTPServer struct {
	engine *gin.Engine
}

func NewHTTPServer(engine *gin.Engine) *HTTPServer {
	return &HTTPServer{
		engine: engine,
	}
}

func (s *HTTPServer) Run(ctx context.Context) {
	if err := s.engine.Run(":8080"); err != nil {
		panic(err)
	}
}
