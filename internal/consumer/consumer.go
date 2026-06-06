package consumer

import (
	"context"
	"sync"

	"github.com/fnuritdinov/firstService/internal/service/eventBus"
	"github.com/fnuritdinov/firstService/internal/service/eventLogs"
	"github.com/fnuritdinov/firstService/pkg/logger"
)

func StartAuditConsumer(ctx context.Context, wg *sync.WaitGroup, bus *eventBus.Bus, logger logger.Logger) {
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			select {
			case <-ctx.Done():
				return
			case event := <-bus.Subscribe():
				eventLogs.Audit(event.UserID, event.Type, "message from bus.Subscribe")
			}
		}
	}()
}
