package entity

type Status int

const (
	StatusActive Status = iota + 1
	StatusDeleted
	StatusHidden
)
