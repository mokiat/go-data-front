package common_test_help

import "github.com/mokiat/go-data-front/common"

type EventHandlerTracker struct {
	Events []common.Event
}

func (h *EventHandlerTracker) Handle(event common.Event) error {
	h.Events = append(h.Events, event)
	return nil
}
