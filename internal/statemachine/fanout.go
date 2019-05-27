package statemachine

import (
  "github.com/lavalamp-/ipv666/internal/fanout"
  "github.com/spf13/viper"
)

func fanOutSlash64s() error {
  return fanout.Slash64s(viper.GetString("PingScanBandwidth"))
}

func fanOutNybbleAdjacent() error {
  return fanout.NybbleAdjacent(viper.GetString("PingScanBandwidth"))
}
