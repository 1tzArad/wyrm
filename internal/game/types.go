package game

type Direction string

const (
	UP    Direction = "up"
	DOWN  Direction = "down"
	LEFT  Direction = "left"
	RIGHT Direction = "right"
)

type ChatTypes string

const (
	GLOBAL  ChatTypes = "global"
	WHISPER ChatTypes = "whisper"
	PRIVATE ChatTypes = "private"
)
