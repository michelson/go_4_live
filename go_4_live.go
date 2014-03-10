package main

import(
  "fmt"
  "github.com/hypebeast/go-osc/osc"
  //"./live_connection"
  //"time"
  )

var Cfg = Config{host: "localhost", port: 8080}
var conn = NewLiveConnection(Cfg.host, Cfg.port) //client(&config)

func main(){
  fmt.Println(fmt.Sprintf("Creating OSC Client - IP: %s, Port: %d", Cfg.host, Cfg.port))
  conn.Send( "/live/tempo", int32(111) , string("hola"))
}

type Config struct{
  host string
  port int
}

// LIVE CONNECTION

type LiveConnection struct {
  host string
  port int
  connection osc.OscClient
  server osc.OscServer
}

func NewLiveConnection(host string, port int) *LiveConnection {
  //return &LiveConnection{host: host, port: port} //one liner
  p := new(LiveConnection)
  p.host = host
  p.port = port
  p.connection = *osc.NewOscClient(host, port)
  return p
}

//func (c *LiveConnection) Send(path string, msg *osc.OscMessage) {
func (c *LiveConnection) Send(path string, args ...interface{} ) {
  msg := osc.NewOscMessage(path)
  //fmt.Println(args)

  for _,element := range args {
    // element is the element from args for where we are
    msg.Append(element) //int32(111)
  }
  fmt.Println(msg.Arguments)
  fmt.Println(fmt.Sprintf("sent to %s", path))
  fmt.Println(c.connection)
  c.connection.Send(msg)
}


/*
func (c *LiveConnection) LivePath(arg *osc.OscMessage) {
  c.connection.Send("/live_path", arg)
}
*/



// LIVE SET
  //properties
  //childs

// LIVE OBJECT

// OBJECTS

// Clip

// Clip slot

// Device

// Device parameter

// Mixer device

// Track
