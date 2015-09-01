package main

import (
  "github.com/pwnall/mapreduced/mrd"
  "flag"
  "fmt"
  "os"
)

var controller mrd.Controller

func main() {
  var port int
  var bindAddress string
  var namePrefix string

  flag.StringVar(&bindAddress, "bind", "127.0.0.1",
      "Interface IP to bind to when listening to HTTP connections")
  flag.IntVar(&port, "port", 8912, "Port to listen to for HTTP connections")
  flag.StringVar(&namePrefix, "name-prefix", "mapreduced",
      "Prefix added to the names of all Docker objects created by this")
  flag.Parse()

  if err := controller.Init(namePrefix); err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
  if err := controller.CleanOldState(); err != nil {
    fmt.Println(err)
    os.Exit(1)
  }

  var builder TemplateBuilder
  if err := builder.Init(id, zipReaderAt, zipSize); err != nil {
    return err
  }

}
