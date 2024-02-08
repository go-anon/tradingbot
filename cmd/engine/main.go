package main

import (
	"context"
	"os"
	"sync"

	"go-anon/tradingbot/cmd/engine/internal/candles"
	"go-anon/tradingbot/cmd/engine/internal/orders"
	"go-anon/tradingbot/cmd/engine/internal/portfolio"
	"go-anon/tradingbot/cmd/engine/internal/strategies"

	"github.com/go-anon/simple/cmds"
	"github.com/go-anon/simple/logs"
)

func main() {
	clean := logs.SetByName(logs.Zap, logs.LevelDebug, "")
	defer clean()

	logs.Info("Starting...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg := sync.WaitGroup{}

	notifier := cmds.NewNotifier(nil)
	notifier.Listen(ctx)

	// ------------- Module Start-up -------------

	cnCandle := make(chan float64)           // Create channel for candles module to send candles
	cnCandlePortfolio := make(chan float64)  // Create channel for Portfolio module to receive candles
	cnCandleStrategies := make(chan float64) // Create channel for Strategies module to receive candles

	// Start Candles module
	wg.Add(1)
	candleModule := candles.New(cnCandle, true)
	cmds.Forever(ctx, "candles", candleModule.Start, candleModule.Panic, wg.Done)

	// Fan-out from Candle to Portfolio and Strategies
	go func() {
		for msg := range cnCandle {
			cnCandlePortfolio <- msg
			cnCandleStrategies <- msg
		}
	}()

	// Start Portfolio module
	wg.Add(1)
	portfolioModule := portfolio.New(cnCandlePortfolio, false)
	cmds.Forever(ctx, "portfolio", portfolioModule.Start, portfolioModule.Panic, wg.Done)

	cnSignal := make(chan float64) // Create channel for Strategies module to send signals
	// cnSignalOrders := make(chan float64)

	// Start Strategies module
	wg.Add(1)
	strategiesModule := strategies.New(cnCandleStrategies, cnSignal, false)
	cmds.Forever(ctx, "strategies", strategiesModule.Start, strategiesModule.Panic, wg.Done)

	// Start Orders module
	wg.Add(1)
	ordersModule := orders.New(cnSignal, true)
	cmds.Forever(ctx, "orders", ordersModule.Start, ordersModule.Panic, wg.Done)

	// ------------- End of module start-up -------------

	// Wait for a shutdown event: a system signal or normal termination
	select {
	case <-notifier.ShutdownChan():
		logs.Info("Received shutdown signal, performing graceful shutdown...")
		cancel()
		// Perform shutdown tasks in response to system signal
	case <-ctx.Done():
		logs.Info("Application is shutting down, no system signal received.")
		// Perform shutdown tasks due to normal application termination
	}

	wg.Wait()

	cleanup(candleModule.Clean, portfolioModule.Clean, strategiesModule.Clean, ordersModule.Clean)

	logs.Info("Shutting down...")
	os.Exit(1)
}

// Cleanup function to close resources, save states, etc.
func cleanup(cf ...func()) {
	logs.Info("Cleaning up resources...")

	if len(cf) > 0 {
		for _, f := range cf {
			f()
		}
	}
}
