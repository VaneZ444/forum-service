package usecase

import "errors"

var (
	ErrCategoryAlreadyExists = errors.New("category already exists")
	ErrCategoryNotFound      = errors.New("category not found")
	ErrTopicNotFound         = errors.New("topic not found")
	ErrPostNotFound          = errors.New("post not found")
	ErrCommentNotFound       = errors.New("comment not found")
	ErrTagNotFound           = errors.New("tag not found")
	ErrInvalidLimit          = errors.New("invalid limit")
	ErrInvalidOffset         = errors.New("invalid offset")
	ErrUpdateFailed          = errors.New("update failed")
	ErrDeleteFailed          = errors.New("delete failed")
)
