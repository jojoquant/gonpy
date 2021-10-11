package converter

import (
	"gonpy/trader/engine"
	"gonpy/trader/object"
)

type OffsetConverter struct {
	MainEngine *engine.MainEngine
	Holdings   map[string]*PositionHolding
}

func (o *OffsetConverter)UpdatePosition(position *object.PositionData){

}

func (o *OffsetConverter)IsConvertRequired(vtSymbol string) bool{
	// if contract:= o.MainEngine
	return true
}