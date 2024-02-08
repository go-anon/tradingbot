package portfolio

import (
	"context"
	"fmt"

	"go-anon/tradingbot/internal/module"

	"github.com/go-anon/simple/logs"
)

type Module struct {
	module.Module

	candle <-chan float64

	count      int
	test_panic bool
}

func New(candle <-chan float64, test_panic bool) *Module {
	return &Module{
		Module:     *module.New("Portfolio"),
		candle:     candle,
		test_panic: test_panic,
	}
}

func (m *Module) Start(ctx context.Context) {
	logs.Info("Module ( %s ) is running...", m.Name())

	v := float64(0)
	for {
		select {
		case <-ctx.Done():
			m.Stop()
			return

			// recibiendo velas
		case v = <-m.candle:
			m.print(v)

			m.count++
			if m.test_panic && (m.count%5) == 0 {
				panic(fmt.Sprintf("%s fall!", m.Name()))
			}
		}
	}
}

func (m *Module) print(v float64) {
	logs.Debug("%s received: %f", m.Name(), v)
}
