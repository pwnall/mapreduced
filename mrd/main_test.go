package mrd

// This file doesn't have any test cases. It only contains TestMain, which has
// setup and teardown code for all the tests.

import (
  "fmt"
  "os"
  "testing"
)

func TestMain(m *testing.M) {
  // Test archives will be stored in test_tmp.
  if err := os.MkdirAll("../test_tmp", 0777); err != nil {
    fmt.Printf("%v\n", err)
    os.Exit(1)
  }

  helloZip, err := os.Create("../test_tmp/hello.zip")
  if err != nil {
    fmt.Printf("%v\n", err)
    os.Exit(1)
  }
  if err := ZipDirectory("../testdata/hello", helloZip); err != nil {
    fmt.Printf("%v\n", err)
    os.Exit(1)
  }
  helloZip.Close()

  fibZip, err := os.Create("../test_tmp/fib.zip")
  if err != nil {
    fmt.Printf("%v\n", err)
    os.Exit(1)
  }
  if err := ZipDirectory("../testdata/fib", fibZip); err != nil {
    fmt.Printf("%v\n", err)
    os.Exit(1)
  }
  fibZip.Close()

  var controller Controller
  controller.Init("mapreduced_mrdtests")
  if err := controller.CleanOldState(); err != nil {
    fmt.Printf("%v\n", err)
    os.Exit(1)
  }

  result := m.Run()
  os.Exit(result)
}
