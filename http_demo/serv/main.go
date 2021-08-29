package main

import (
	"fmt"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.New()
	engine.Handle(http.MethodGet, "/test/context", ContextTest)

	listener, err := net.Listen("tcp", ":1789")
	if err != nil {
		panic(err)
	}
	if err := engine.RunListener(listener); err != nil {
		panic(err)
	}
}

func ContextTest(c *gin.Context) {
	ctx := c.Request.Context()
	fmt.Println("receive request")
	select {
	case <-ctx.Done():
		fmt.Println("context done")
	}
	c.JSON(http.StatusOK, "http response ok")
}
/*
curl 'http://localhost:1789/test/context'
如果不主动ctrl c 会一直存在
*/
