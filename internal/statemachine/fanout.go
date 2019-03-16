package statemachine

import (
  "github.com/lavalamp-/ipv666/internal/fanout"
  "github.com/spf13/viper"
)

func fanOut() (error) {
  _, err := fanout.FanOut(viper.GetString("PingScanBandwidth"))
  return err
}
