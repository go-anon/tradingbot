package candles

import (
	"context"
	"fmt"
	"time"

	"go-anon/tradingbot/internal/module"

	"github.com/go-anon/simple/logs"
)

type Module struct {
	module.Module

	candle chan<- float64

	count      int
	test_panic bool
}

func New(candle chan<- float64, test_panic bool) *Module {
	return &Module{
		Module:     *module.New("Candles"),
		candle:     candle,
		test_panic: test_panic,
	}
}

func (m *Module) Start(ctx context.Context) {
	logs.Info("Module ( %s ) is running...", m.Name())

	for {
		select {
		case <-ctx.Done():
			m.Stop()
			return

			// enviando velas
		case m.candle <- float64(m.count) * 1.1:

			time.Sleep(time.Second * 2)

			m.count++
			if m.test_panic && (m.count%4) == 0 {
				panic(fmt.Sprintf("%s fall!", m.Name()))
			}
		}
	}
}

func (m *Module) Stop() {
	m.Module.Stop()
	close(m.candle)
}
