package gnss

import (
	"fmt"
	"strings"
	"time"

	"github.com/streamingfast/gnss-controller/device/neom9n"
	"github.com/streamingfast/hivemapper-data-logger/data"
)

type EventFeed struct {
	subscriptions data.Subscriptions
}

func NewEventFeed() *EventFeed {
	return &EventFeed{
		subscriptions: make(data.Subscriptions),
	}
}

func (e *EventFeed) Subscribe(name string) *data.Subscription {
	sub := &data.Subscription{
		IncomingEvents: make(chan data.Event),
	}
	e.subscriptions[name] = sub
	return sub
}

func (e *EventFeed) HandleData(d *neom9n.Data) {
	e.emit(NewGnssEvent(d))
}

func (e *EventFeed) emit(event data.Event) {
	event.SetTime(time.Now())
	for _, subscription := range e.subscriptions {
		subscription.IncomingEvents <- event
	}
}

type GnssEvent struct {
	*data.BaseEvent
	Data *neom9n.Data `json:"data"`
}

func NewGnssEvent(d *neom9n.Data) *GnssEvent {
	return &GnssEvent{
		BaseEvent: data.NewBaseEvent("GNSS_EVENT", nil), // no x, y, z GForces for gnss data
		Data:      d,
	}
}

func (g *GnssEvent) String() string {
	var sb strings.Builder
	sb.WriteString("GNSS\n")
	sb.WriteString(fmt.Sprintf("\tLatitude: %.2f\n", g.Data.Latitude))
	sb.WriteString(fmt.Sprintf("\tLongitude: %.2f\n", g.Data.Longitude))
	sb.WriteString(fmt.Sprintf("\tAltitude: %.2f\n", g.Data.Altitude))
	sb.WriteString(fmt.Sprintf("\tHeading: %.2f\n", g.Data.Heading))
	sb.WriteString(fmt.Sprintf("\tSpeed: %.2f\n", g.Data.Speed))
	return sb.String()
}
