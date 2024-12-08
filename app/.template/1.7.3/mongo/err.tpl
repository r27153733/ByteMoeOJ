package model

import (
	"errors"

	"github.com/r27153733/fastgozero/core/stores/mon"
)

var (
	ErrNotFound        = mon.ErrNotFound
	ErrInvalidObjectId = errors.New("invalid objectId")
)
