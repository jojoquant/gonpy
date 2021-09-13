package engine

import "fmt"

type CtaEngine struct {
	// BaseEnginer
	BaseEngine
}

func (c *CtaEngine) Close() {
	fmt.Println("cta engine close")
}

func (c *CtaEngine) GetName() string {
	return c.Name
}

func (c *CtaEngine) SetEventEngine(eventEngine *EventEngine) {
	c.EventEngine = eventEngine
}

func NewCtaEngine(eventEngine *EventEngine) *CtaEngine {
	c := &CtaEngine{}
	c.Name = "CtaEngine"
	c.EventEngine = eventEngine
	return c
}
