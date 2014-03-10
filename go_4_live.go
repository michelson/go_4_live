package main

import(
  "fmt"
  "github.com/hypebeast/go-osc/osc"
  //"./live_connection"
  "time"
  )

var Cfg = Config{host: "localhost", send_port: 7402, receive_port: 7403}
var conn = NewLiveConnection(Cfg.host, Cfg.send_port, Cfg.receive_port) //client(&config)
var channel chan string = make(chan string)
var server_response []*osc.OscMessage  //[]interface{}

func main(){

  go conn.ListenAndServe()
  go pinger()
  var input string
  fmt.Scanln(&input)
}

func pinger() {
    for {
      conn.Send( "/live_path", true, int32(111), string("hola"))
      time.Sleep(time.Second * 1)
    }
}

type Config struct{
  host string
  send_port int
  receive_port int
}

// LIVE CONNECTION

type LiveConnection struct {
  host string
  send_port int
  receive_port int
  connection osc.OscClient
  server osc.OscServer
}

func NewLiveConnection(host string, send_port int, receive_port int) *LiveConnection {
  //return &LiveConnection{host: host, port: port} //one liner
  p := new(LiveConnection)
  p.host = host
  p.send_port = send_port
  p.receive_port = receive_port
  p.server = *osc.NewOscServer(host, receive_port)
  fmt.Println(fmt.Sprintf("Creating OSC Server - IP: %s, Port: %d", host, receive_port))
  p.connection = *osc.NewOscClient(host, send_port)
  fmt.Println(fmt.Sprintf("Creating OSC Client - IP: %s, Port: %d", host, send_port))
  return p
}

func (conn *LiveConnection) ListenAndServe(){

  //si
  conn.server.AddMsgHandler("/slot1", func(msg *osc.OscMessage) {
    fmt.Println("Received message from ", msg)
    server_response = append(server_response, msg )
    //@response << [msg.address, msg.args]
  })

  //ss
  conn.server.AddMsgHandler("/slot1", func(msg *osc.OscMessage) {
    fmt.Println("Received message from ", msg)
    server_response = append(server_response, msg )
    //@response << [msg.address, msg.args]
  })

  //s*
  conn.server.AddMsgHandler("/slot1", func(msg *osc.OscMessage) {
    fmt.Println("Received message from ", msg)
    server_response = append(server_response, msg )
    //@response << [msg.address, msg.args]
  })
  //ssi
  conn.server.AddMsgHandler("/slot3", func(msg *osc.OscMessage) {
    fmt.Println("Received message from ", msg)
    server_response = append(server_response, msg )
    //@response << [msg.address, msg.args]
  })

  fmt.Printf("Listening on %s:%d\n", conn.host, conn.receive_port)
  err := conn.server.ListenAndServe()
  if err != nil {
    fmt.Println("Error")
  }
}

//func (c *LiveConnection) Send(path string, msg *osc.OscMessage) {
func (c *LiveConnection) Send(path string, cbk bool, args ...interface{} ) {
  msg := osc.NewOscMessage(path)
  //fmt.Println(args)

  for _,element := range args {
    // element is the element from args for where we are
    msg.Append(element) //int32(111)
  }
  //fmt.Println(msg.Arguments)
  fmt.Println(fmt.Sprintf("sent to %s", path))
  //fmt.Println(c.connection)
  c.connection.Send(msg)
  if cbk {
    c.get_callback()
  }
}

func (c *LiveConnection) get_callback() *osc.OscMessage {

  for len(server_response) <= 0 {
    time.Sleep(1 * time.Millisecond)
  }

  fmt.Println("Server response:", server_response[0])
  return_value := server_response[0]
  server_response = nil // clear array
  return return_value
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
