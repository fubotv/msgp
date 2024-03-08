package main

import (
	"github.com/tinylib/msgp/msgp"
)

//go:generate msgp

type BusDepot struct {
	NumBuses         int
	HasExpressRoutes bool
}

type TrolleyTerminal struct {
	Name            string
	HasBalloonLoops bool
	IsUnderground   bool
}

type Airport struct {
	Name         string
	NumPlanes    int
	LoungeRating float64
}

type TransportationHub struct {
	Type string      `msg:"Type"`            // can be "bus", "trolley", or "airport"
	Hub  interface{} `msg:"Hub,polymorphic"` // can be BusDepot, TrolleyTerminal, or Airport
}

type City struct {
	Name string
	Hubs []*TransportationHub
}

// ChooseType implements PolymorphicResolver to resolve the type of Hub
func (c *TransportationHub) ChooseType(field string) (msgp.MsgPackDeserializer, error) {
	switch c.Type {
	case "bus":
		return &BusDepot{}, nil
	case "trolley":
		return &TrolleyTerminal{}, nil
	case "airport":
		return &Airport{}, nil
	default:
		return nil, nil
	}
}

var _ msgp.PolymorphicResolver = (*TransportationHub)(nil)
