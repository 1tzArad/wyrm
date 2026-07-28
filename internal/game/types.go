package game

type Direction string

const (
	UP    Direction = "up"
	DOWN  Direction = "down"
	LEFT  Direction = "left"
	RIGHT Direction = "right"
)

func (d Direction) IsValid() bool {
	switch d {
	case UP, DOWN, LEFT, RIGHT:
		return true
	default:
		return false
	}
}

type ChatTypes string

const (
	GLOBAL  ChatTypes = "global"
	WHISPER ChatTypes = "whisper"
	PRIVATE ChatTypes = "private"
)

func (t ChatTypes) IsValid() bool {
	switch t {
	case GLOBAL, WHISPER, PRIVATE:
		return true
	default:
		return false
	}
}
