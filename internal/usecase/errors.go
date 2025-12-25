package usecase

import "errors"

var (
	ErrEmptyTitle   = errors.New("task title is empty")
	ErrTaskNotFound = errors.New("task not found")
)
