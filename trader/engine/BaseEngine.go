package engine

type BaseEnginer interface {
	Close()
	GetName() string
	SetEventEngine(eventEngine *EventEngine)
}


