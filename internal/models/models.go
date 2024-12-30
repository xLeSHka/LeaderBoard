package models

type AddRowRequest struct {
	RowID string `json:"rowID"`
}
type AddRowResponse struct {
	Success bool `json:"success"`
}

type LeaderBoardRequest struct {
}
type LeaderBoardResponse struct {
	LeaderBoard []map[string]any `json:"leaderBoard"`
}

type SetRequest struct {
	RowID string       `json:"rowID"`
	Data  map[string]any `json:"data"`
}
type SetResponse struct {
	Success bool `json:"success"`
}

type RemoveRequest struct {
	RowID string `json:"rowID"`
}
type RemoveResponse struct {
	Success bool `json:"success"`
}

type DeleteRequest struct {
}
type DeleteResponse struct {
	Success bool `json:"success"`
}
type AddLeaderBoardRequest struct {
}
type AddLeaderBoardResponse struct {
	Success bool `json:"success"`
}
