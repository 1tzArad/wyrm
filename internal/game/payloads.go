package game

type MovePayload struct {
	Direction Direction `json:"direction"`
}

type ChatPayload struct {
	Type   ChatTypes `json:"chat_type"`
	Target string    `json:"target"`
	Text   string    `json:"text"`
}

type ChatReceivePayload struct {
	PlayerUUID string `json:"player_uuid"`
	Text       string `json:"text"`
}
