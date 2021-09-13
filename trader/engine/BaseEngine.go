package engine

type BaseEnginer interface {
	Close()
	GetName() string
	SetEventEngine(eventEngine *EventEngine)
}

type BaseEngine struct{
	Name        string
	EventEngine *EventEngine
}


