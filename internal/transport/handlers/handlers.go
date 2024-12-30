package handlers

import (
	"context"
	"encoding/json"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/xLeSHka/LeaderBoard/internal/models"
	"go.uber.org/zap"
)

type Service interface {
	AddRow(ctx context.Context, LeaderBoardID, RowID string) (bool, error)
	AddLeaderBoard(ctx context.Context, LeaderBoardID string) (bool, error)
	Get(ctx context.Context, LeaderBoardID string, SortBy string) ([]map[string]any, error)
	Remove(ctx context.Context, LeaderBoardID, RowID string) (bool, error)
	Set(ctx context.Context, LeaderBoardID, RowID string, data map[string]any) (bool, error)
	Delete(ctx context.Context, LeaderBoardID string) (bool, error)
}
type LeaderBoardService struct {
	service Service
	l       zap.Logger
}

func New(srv Service, l zap.Logger) *LeaderBoardService {
	return &LeaderBoardService{service: srv, l: l}
}

func (s *LeaderBoardService) AddRow(c echo.Context) error {
	if c.Request().Header.Get("Content-Type") != "application/json" {
		return echo.NewHTTPError(400, "Content type not allowed")
	}
	req := models.AddRowRequest{}
	err := c.Bind(&req)
	if err != nil {
		s.l.Error(err.Error())
		return echo.NewHTTPError(400, "Body is not valid")
	}
	id := c.Param("leaderboardID")
	if id == "" {
		return echo.NewHTTPError(400, "leaderID is not valid")
	}

	suc, err := s.service.AddRow(c.Request().Context(), id, req.RowID)
	if err != nil {
		s.l.Error(err.Error())
		return echo.NewHTTPError(400, "failed add row to leaderboard")
	}
	return c.JSON(200, models.AddRowResponse{Success: suc})
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (s *LeaderBoardService) GetLeaderBoard(c echo.Context) error {
	id := c.Param("leaderboardID")
	if id == "" {
		return echo.NewHTTPError(400, "leaderID is not valid")
	}
	columnName := c.QueryParam("sortBy")
	if columnName == "" {
		return echo.NewHTTPError(400, "sortBy is not valid")
	}
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()
	for {
		leaderBoard, err := s.service.Get(c.Request().Context(), id, columnName)
		if err != nil && err != websocket.ErrCloseSent {
			s.l.Error(err.Error())
			return echo.NewHTTPError(400, "failed get leaderboard")
		}
		data, err := json.Marshal(leaderBoard)
		if err != nil && err != websocket.ErrCloseSent {
			return err
		}
		// Write
		err = ws.WriteMessage(websocket.TextMessage, data)
		if err != nil && err != websocket.ErrCloseSent {
			return err
		}
		_, _, err = ws.ReadMessage()
		if c, k := err.(*websocket.CloseError); k {
			if c.Code == 1000 {
				break
			}
			return err
		}
	}
	return nil
}

func (s *LeaderBoardService) AddLeaderBoard(c echo.Context) error {
	id := c.Param("leaderboardID")
	if id == "" {
		return echo.NewHTTPError(400, "leaderID is not valid")
	}

	suc, err := s.service.AddLeaderBoard(c.Request().Context(), id)
	if err != nil {
		s.l.Error(err.Error())
		return echo.NewHTTPError(400, "failed add leaderboard")
	}
	return c.JSON(200, models.AddLeaderBoardResponse{Success: suc})
}

func (s *LeaderBoardService) Set(c echo.Context) error {
	if c.Request().Header.Get("Content-Type") != "application/json" {
		return echo.NewHTTPError(400, "Content type not allowed")
	}
	req := models.SetRequest{}
	err := c.Bind(&req)
	if err != nil {
		s.l.Error(err.Error())
		return echo.NewHTTPError(400, "Body is not valid")
	}
	id := c.Param("leaderboardID")
	if id == "" {
		return echo.NewHTTPError(400, "leaderID is not valid")
	}

	suc, err := s.service.Set(c.Request().Context(), id, req.RowID, req.Data)
	if err != nil {
		s.l.Error(err.Error())
		return echo.NewHTTPError(400, "failed set row`s data")
	}
	return c.JSON(200, models.SetResponse{Success: suc})
}

func (s *LeaderBoardService) Delete(c echo.Context) error {
	id := c.Param("leaderboardID")
	if id == "" {
		return echo.NewHTTPError(400, "leaderID is not valid")
	}

	suc, err := s.service.Delete(c.Request().Context(), id)
	if err != nil {
		s.l.Error(err.Error())
		return echo.NewHTTPError(400, "failed set row`s data")
	}
	return c.JSON(200, models.DeleteResponse{Success: suc})
}

func (s *LeaderBoardService) Remove(c echo.Context) error {
	if c.Request().Header.Get("Content-Type") != "application/json" {
		return echo.NewHTTPError(400, "Content type not allowed")
	}
	req := models.RemoveRequest{}
	err := c.Bind(&req)
	if err != nil {
		s.l.Error(err.Error())
		return echo.NewHTTPError(400, "Body is not valid")
	}
	id := c.Param("leaderboardID")
	if id == "" {
		return echo.NewHTTPError(400, "leaderID is not valid")
	}

	suc, err := s.service.Remove(c.Request().Context(), id, req.RowID)
	if err != nil {
		s.l.Error(err.Error())
		return echo.NewHTTPError(400, "failed remove row")
	}
	return c.JSON(200, models.RemoveResponse{Success: suc})
}
