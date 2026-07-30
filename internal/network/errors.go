package network

import "errors"

var (
	ErrPlayerNotConnected = errors.New("player not connected!")
	ErrPlayerChannelFull  = errors.New("player channel full")
)
