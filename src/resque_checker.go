/*
 * resque-tester.go
 * Resque Tester
 *
 * Created by Artyom Fedenko on 2015-05-13.
 * Copyright 2015 Artyom Fedenko. All rights reserved.
 */


package main

import  (
  "fmt"
  "net/http"
  "os"
  "io"
  "sync"
  "time"
  "strings"
  "io/ioutil"
)

var wg sync.WaitGroup

// REQUEST URL
var RESQUE_URL string  = os.Args[1]

// LOG FILE PATHES
const LOG_NAME = "../log/resque_log"

func main() {

  createFileInNotExists(LOG_NAME)

  fmt.Println("Starting all night scripting")

  //Add more threads if needed
  wg.Add(1)
  go startCheck(RESQUE_URL, LOG_NAME)

  wg.Wait()
}

func startCheck(url string, logName string) {
  for {
    checkResquePage(url, logName)
    time.Sleep(1000 * time.Millisecond)
  }
  defer wg.Done()
}

func checkResquePage(url string, logName string) {
  var output string = ""

  fmt.Println("Checking " + url)

  resp, err := http.Get(url)

  printErr(err)

  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)

  printErr(err)

  if strings.Contains(string(body), "Can't connect to MongoDB!") {
    output = "Error | " + time.Now().String() + " | " + url + "\n"
  } else {
    output = "Success | " + time.Now().String() + " | " + url + "\n"
  }

  f, err := os.OpenFile(logName, os.O_RDWR|os.O_APPEND, 0660);

  printErr(err)

  io.WriteString(f, output)
}

func printErr(err error) {
  if err != nil {
    fmt.Println(err)
  }
}

func createFileInNotExists(filename string) {
  if _, err := os.Stat(filename); os.IsNotExist(err) {
    fmt.Printf("Creating file: %s", filename)
    os.Create(filename)
    return
  }
}
