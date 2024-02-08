package config

import (
	"github.com/go-anon/simple/configs/env"
)

var (
	Exchange = env.Get("Exchange", "")

	Symbol = env.Get("Symbol", "")
)
