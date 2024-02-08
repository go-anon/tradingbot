package strategies

import (
	"context"
	"fmt"
	"time"

	"go-anon/tradingbot/internal/module"

	"github.com/go-anon/simple/logs"
)

type Module struct {
	module.Module

	candle <-chan float64
	signal chan<- float64

	count      int
	test_panic bool
}

func New(candle <-chan float64, signal chan<- float64, test_panic bool) *Module {
	return &Module{
		Module:     *module.New("Strategies"),
		candle:     candle,
		signal:     signal,
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

			// enviando señales
		case m.signal <- float64(m.count) * 1.2:

			time.Sleep(time.Second * 3)

			m.count++
			if m.test_panic && (m.count%3) == 0 {
				panic(fmt.Sprintf("%s fall!", m.Name()))
			}

			// recibiendo velas
		case v = <-m.candle:
			m.print(v)
		}
	}
}

func (m *Module) Stop() {
	m.Module.Stop()
	close(m.signal)
}

func (m *Module) print(v float64) {
	logs.Debug("%s received: %f", m.Name(), v)
}
