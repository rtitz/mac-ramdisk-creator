package main

import (
  "fmt"

  "github.com/rtitz/mac-ramdisk-creator/variables"
)

func main () {
  fmt.Printf("%s %s\n\n", variables.AppName, variables.AppVersion)
}
