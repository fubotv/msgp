package main

import (
	"bytes"
	"github.com/fubotv/msgp/msgp"
	"reflect"
	"testing"
)

func TestPolymorphicUnmarshal(t *testing.T) {
	city := &City{
		Name: "San Francisco",
		Hubs: []*TransportationHub{
			{
				Type: "bus",
				Hub: &BusDepot{
					NumBuses:         100,
					HasExpressRoutes: true,
				},
			},
			{
				Type: "trolley",
				Hub: &TrolleyTerminal{
					Name:            "Fisherman's Wharf",
					HasBalloonLoops: true,
					IsUnderground:   false,
				},
			},
			{
				Type: "airport",
				Hub: &Airport{
					Name:         "SFO",
					NumPlanes:    1000,
					LoungeRating: 4.5,
				},
			},
		},
	}

	var buf bytes.Buffer
	w := msgp.NewWriter(&buf)
	err := msgp.Encode(w, city)
	if err != nil {
		t.Fatal(err)
	}

	r := msgp.NewReader(&buf)
	deserialized := &City{}
	err = msgp.Decode(r, deserialized)
	if err != nil {
		t.Fatal(err)
	}

	if len(deserialized.Hubs) != 3 {
		t.Fatalf("expected 3 hubs, got %d", len(deserialized.Hubs))
	}

	for i := range deserialized.Hubs {
		if reflect.TypeOf(deserialized.Hubs[i].Hub) != reflect.TypeOf(city.Hubs[i].Hub) {
			t.Fatalf("expected deserialized.Hubs[%v].Hub to be a TrolleyTerminal, got %T", i, deserialized.Hubs[1].Hub)
		}
	}
}
