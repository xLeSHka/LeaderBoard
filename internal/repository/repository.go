package repository

import (
	"context"
	"reflect"

	"github.com/redis/go-redis/v9"
)

type Repository struct {
	db *redis.Client
}

func New(db *redis.Client) *Repository {
	return &Repository{db}
}
func (r *Repository) AddRow(ctx context.Context, LeaderBoardID, RowID string) (bool, error) {
	err := r.IsMember(ctx, "ids", LeaderBoardID)
	if err != nil {
		return false, err
	}
	err = r.IsMember(ctx, LeaderBoardID, RowID)
	if err == nil {
		return false, ErrLeaderBoardExists
	}
	if err != ErrNotFound {
		return false, err
	}
	err = r.AddToList(ctx, LeaderBoardID, RowID)
	if err != nil {
		return false, err
	}
	return true, nil
}
func (r *Repository) AddLeaderBoard(ctx context.Context, LeaderBoardID string) (bool, error) {
	err := r.IsMember(ctx, "ids", LeaderBoardID)
	if err == nil {
		return false, ErrLeaderBoardExists
	}
	if err != ErrNotFound {
		return false, err
	}
	err = r.AddToList(ctx, "ids", LeaderBoardID)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *Repository) Get(ctx context.Context, LeaderBoardID string, SortBy string) ([]map[string]any, error) {
	err := r.IsMember(ctx, "ids", LeaderBoardID)
	if err != nil {
		return nil, err
	}
	rows, err := r.Members(ctx, LeaderBoardID)
	if err != nil {
		return nil, err
	}
	res := make([]map[string]any, 0)
	for _, rowID := range rows {
		pairs, err := r.db.HGetAll(ctx, LeaderBoardID+"_"+rowID+"_types").Result()
		if err != nil {
			return nil, err
		}

		row := make(map[string]interface{})
		for k, t := range pairs {
			var val interface{}
			resp := r.db.HGet(ctx, LeaderBoardID+"_"+rowID, k)
			if resp.Err() != nil {
				return nil, resp.Err()
			}
			switch t {
			case "int":
				val, _ = resp.Int()
			case "float64":
				val, _ = resp.Float64()
			case "string":
				val, _ = resp.Result()
			default:
				val, _ = resp.Result()
			}
			row[k] = val
		}
		res = append(res, row)
	}
	return res, nil
}
func (r *Repository) Remove(ctx context.Context, LeaderBoardID, RowID string) (bool, error) {
	err := r.IsMember(ctx, "ids", LeaderBoardID)
	if err != nil {
		return false, err
	}
	err = r.IsMember(ctx, LeaderBoardID, RowID)
	if err != nil {
		return false, err
	}

	keys, err := r.db.HKeys(ctx, LeaderBoardID+"_"+RowID).Result()
	if err != nil {
		return false, err
	}
	_, err = r.db.HDel(ctx, LeaderBoardID+"_"+RowID, keys...).Result()
	if err != nil {
		return false, err
	}
	keys2, err := r.db.HKeys(ctx, LeaderBoardID+"_"+RowID+"_types").Result()
	if err != nil {
		return false, err
	}
	_, err = r.db.HDel(ctx, LeaderBoardID+"_"+RowID+"_types", keys2...).Result()
	if err != nil {
		return false, err
	}
	_, err = r.db.SRem(ctx, LeaderBoardID, RowID).Result()
	if err != nil {
		return false, err
	}
	return true, nil
}
func (r *Repository) Set(ctx context.Context, LeaderBoardID, RowID string, data map[string]any) (bool, error) {
	err := r.IsMember(ctx, "ids", LeaderBoardID)
	if err != nil {
		return false, err
	}
	err = r.IsMember(ctx, LeaderBoardID, RowID)
	if err != nil {
		return false, err
	}
	_, err = r.db.HSet(ctx, LeaderBoardID+"_"+RowID, data).Result()
	if err != nil {
		return false, err
	}
	pairs := make(map[string]string)
	for k, v := range data {
		pairs[k] = reflect.TypeOf(v).Name()
	}
	_, err = r.db.HSet(ctx, LeaderBoardID+"_"+RowID+"_types", pairs).Result()
	if err != nil {
		return false, err
	}
	return true, nil
}
func (r *Repository) Delete(ctx context.Context, LeaderBoardID string) (bool, error) {
	err := r.IsMember(ctx, "ids", LeaderBoardID)
	if err != nil {
		return false, err
	}

	rows, err := r.Members(ctx, LeaderBoardID)
	if err != nil && err != ErrNoMembers {
		return false, err
	}
	if err == nil {
		for _, row := range rows {
			_, err := r.Remove(ctx, LeaderBoardID, row)
			if err != nil {
				return false, err
			}
		}

	}
	_, err = r.db.SRem(ctx, "ids", LeaderBoardID).Result()
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *Repository) AddToList(ctx context.Context, key, val string) error {
	_, err := r.db.SAdd(ctx, key, val).Result()
	if err != nil {
		return err
	}
	return nil
}
func (r *Repository) IsMember(ctx context.Context, key, val string) error {
	is, err := r.db.SIsMember(ctx, key, val).Result()
	if err != nil {
		return err
	}
	if !is {
		return ErrNotFound
	}
	return nil
}
func (r *Repository) Members(ctx context.Context, key string) ([]string, error) {
	rows, err := r.db.SMembers(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, ErrNoMembers
	}
	return rows, nil
}
