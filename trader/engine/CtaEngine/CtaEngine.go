package CtaEngine

import (
	"fmt"
	"gonpy/trader/engine"
)

type CtaEngine struct {
	// BaseEnginer
	engine.BaseEngine
}

func (c *CtaEngine) Close() {
	fmt.Println("cta engine close")
}

func (c *CtaEngine) GetName() string {
	return c.Name
}

func (c *CtaEngine) SetEventEngine(eventEngine *engine.EventEngine) {
	c.EventEngine = eventEngine
}

func NewCtaEngine(eventEngine *engine.EventEngine) *CtaEngine {
	c := &CtaEngine{}
	c.Name = "CtaEngine"
	c.EventEngine = eventEngine
	return c
}
