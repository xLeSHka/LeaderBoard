package service

import (
	"context"
	"sort"
)

type Repo interface {
	AddRow(ctx context.Context, LeaderBoardID, RowID string) (bool, error)
	AddLeaderBoard(ctx context.Context, LeaderBoardID string) (bool, error)
	Get(ctx context.Context, LeaderBoardID string, SortBy string) ([]map[string]any, error)
	Remove(ctx context.Context, LeaderBoardID, RowID string) (bool, error)
	Set(ctx context.Context, LeaderBoardID, RowID string, data map[string]any) (bool, error)
	Delete(ctx context.Context, LeaderBoardID string) (bool, error)
}

type Service struct {
	Repo Repo
}

func New(repo Repo) *Service {
	return &Service{repo}
}
func (s *Service) AddRow(ctx context.Context, LeaderBoardID, RowID string) (bool, error) {
	return s.Repo.AddRow(ctx, LeaderBoardID, RowID)
}
func (s *Service) AddLeaderBoard(ctx context.Context, LeaderBoardID string) (bool, error) {
	return s.Repo.AddLeaderBoard(ctx, LeaderBoardID)
}
func (s *Service) Get(ctx context.Context, LeaderBoardID string, SortBy string) ([]map[string]any, error) {
	data, err := s.Repo.Get(ctx, LeaderBoardID, SortBy)
	if err != nil {
		return nil, err
	}
	sort.Slice(data, func(i, j int) bool {
		switch data[i][SortBy].(type) {
		case string:
			a, _ := data[i][SortBy].(string)
			b, _ := data[j][SortBy].(string)
			return a > b
		case int:
			a, _ := data[i][SortBy].(int)
			b, _ := data[j][SortBy].(int)
			return a > b
		case float64:
			a, _ := data[i][SortBy].(float64)
			b, _ := data[j][SortBy].(float64)
			return a > b
		default:
			return false
		}
	})
	return data, nil
}
func (s *Service) Remove(ctx context.Context, LeaderBoardID, RowID string) (bool, error) {
	return s.Repo.Remove(ctx, LeaderBoardID, RowID)
}
func (s *Service) Set(ctx context.Context, LeaderBoardID, RowID string, data map[string]any) (bool, error) {
	return s.Repo.Set(ctx, LeaderBoardID, RowID, data)
}
func (s *Service) Delete(ctx context.Context, LeaderBoardID string) (bool, error) {
	return s.Repo.Delete(ctx, LeaderBoardID)
}
