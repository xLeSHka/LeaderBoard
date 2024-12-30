package repository

import "errors"

var (
	ErrNotFound          = errors.New("value in set not found")
	ErrLeaderBoardExists = errors.New("Leader board alredy exist")
	ErrNoMembers         = errors.New("no members of set")
)
