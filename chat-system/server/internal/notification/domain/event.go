package domain

type Event struct {
	Type string
	Data interface{}
}

type PresenceData struct {
	Account string `json:"account"`
}

type SystemData struct {
	Content string `json:"content"`
}
