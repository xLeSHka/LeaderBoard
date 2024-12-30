package server

import (
	"context"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Service interface {
	AddRow(c echo.Context) error
	GetLeaderBoard(c echo.Context) error
	AddLeaderBoard(c echo.Context) error
	Set(c echo.Context) error
	Delete(c echo.Context) error
	Remove(c echo.Context) error
}

type Server struct {
	server *echo.Echo
}

func New(srv Service) (*Server, error) {
	e := echo.New()
	e.Use(middleware.LoggerWithConfig(
		middleware.LoggerConfig{
			Format:           `{"time":"${time_rfc3339}, "host":"${host}", "method":"${method}", "uri":"${uri}", "status":${status}, "error":"${error}}` + "\n",
			CustomTimeFormat: "2006-01-02 15:04:05",
		},
	))
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowMethods: []string{"POST", "OPTIONS", "GET", "DELETE", "PATCH"},
	}))
	server := &Server{e}
	server.server.POST("/:leaderboardID/addRow", srv.AddRow)
	server.server.POST("/:leaderboardID/addLeaderBoard", srv.AddLeaderBoard)
	server.server.PATCH("/:leaderboardID/set", srv.Set)
	server.server.DELETE("/:leaderboardID/deleteLeaderBoard", srv.Delete)
	server.server.DELETE("/:leaderboardID/remove", srv.Remove)
	server.server.GET("/:leaderboardID/", srv.GetLeaderBoard)
	return server, nil
}
func (s *Server) Start(port int) error {
	return s.server.Start(fmt.Sprintf(":%d", port))
}
func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
