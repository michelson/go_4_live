package main

import(
  "fmt"
  "github.com/hypebeast/go-osc/osc"
  //"./live_connection"
  "time"
  "reflect"
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
  //for {
    //conn.Send( "/live_path", true, int32(111), string("hola"))
    //conn.LivePath(string("goto live_set"))
    var live_set = NewLiveSet()
    fmt.Println("liveset", live_set)
    time.Sleep(time.Second * 1)
  //}
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
    //fmt.Println("Received message from ", msg)
    server_response = append(server_response, msg )
    //@response << [msg.address, msg.args]
  })

  //ss
  conn.server.AddMsgHandler("/slot1", func(msg *osc.OscMessage) {
    //fmt.Println("Received message from ", msg)
    server_response = append(server_response, msg )
    //@response << [msg.address, msg.args]
  })

  //s*
  conn.server.AddMsgHandler("/slot1", func(msg *osc.OscMessage) {
    //fmt.Println("Received message from ", msg)
    server_response = append(server_response, msg )
    //@response << [msg.address, msg.args]
  })
  //ssi
  conn.server.AddMsgHandler("/slot3", func(msg *osc.OscMessage) {
    //fmt.Println("Received message from ", msg)
    server_response = append(server_response, msg )
    //@response << [msg.address, msg.args]
  })

  fmt.Printf("Listening on %s:%d\n", conn.host, conn.receive_port)
  err := conn.server.ListenAndServe()
  if err != nil {
    fmt.Println("Error")
  }
}

func (c *LiveConnection) Send(path string, cbk bool, args ...interface{} ) *osc.OscMessage {
  msg := osc.NewOscMessage(path)

  for _,element := range args {
    // element is the element from args for where we are
    msg.Append(element) //int32(111)
  }
  //fmt.Println(msg.Arguments)
  fmt.Println(fmt.Sprintf("sent to %s", path))
  //fmt.Println(c.connection)
  c.connection.Send(msg)
  //if cbk {
    return c.get_callback() 
  //}
}

func (c *LiveConnection) SetLiveObjectPath(arg string) {
  c.Send("/set_live_object", false, arg)
}

func (c *LiveConnection) LivePath(arg string) *osc.OscMessage {
  res := c.Send("/live_path", true, arg)
  return res
}

func (c *LiveConnection) LiveObject(arg string) *osc.OscMessage {
  res := c.Send("/live_object", true, arg)
  return res
}

func (c *LiveConnection) get_callback() *osc.OscMessage {

  for len(server_response) <= 0 {
    time.Sleep(time.Millisecond / 100)
  }

  fmt.Println("Server response:", server_response[0] )
  return_value := server_response[0]
  server_response = nil // clear array
  return return_value
}

// LIVE SET
type LiveSet struct {
  objects []interface{}
}

func NewLiveSet() *LiveSet {
  p := new(LiveSet)
  p.GetMasterTrack()
  p.GetTracks()
  //p.Tracks()
  //-->p.GetDevicesParameters()
  ////p.devices.map {|d| d.get_parameters }
  //-->p.GetClipSlotsClips()
  ////p.clip_slots.map {|c| c.get_clip }
  fmt.Println( "Scan Complete" )
  return p
}

func (c *LiveSet) AddObject(obj *Track) {
  c.objects = append(c.objects, obj )
  fmt.Println(c.objects)
}

func (c *LiveSet) GetMasterTrack() {
  fmt.Println("Get maser track")
  res := conn.LivePath( string("goto live_set master_track") )
  id := res.Arguments[1].(int32)
  c.AddObject(NewTrack(id, true, 1, c))
}

// esto tiene que devolver numero
func (c *LiveSet) TrackCount() int32 {
  conn.LivePath( string("goto live_set") )
  res := conn.LivePath( string("getcount tracks") ).Arguments[2].(int32)
  //fmt.Println("RES", res)
  return res
}

func (c *LiveSet) GetTracks() {
  fmt.Println("Track count", c.TrackCount())
  for i := int32(0); i < c.TrackCount(); i++ {
    track_id := conn.LivePath(fmt.Sprintf("goto live_set tracks %d", i)).Arguments[1].(int32) //[0][1][1]
    fmt.Println("Track", track_id)
    c.AddObject(NewTrack(track_id, false, i, c))
  }
}

func (c *LiveSet) Tracks(){
  var elems []interface{}

  for _,element := range c.objects {
    //compare type track
    if reflect.TypeOf(element).Elem() == reflect.TypeOf(Track{}) {
      fmt.Println(reflect.TypeOf(element).Elem())
      elems = append(elems, element )
    }
  }
}

func (c *LiveSet) Devices(){
  var elems []interface{}

  for _,element := range c.objects {
    //compare type track
    if reflect.TypeOf(element).Elem() == reflect.TypeOf(Device{}) {
      fmt.Println(reflect.TypeOf(element).Elem())
      elems = append(elems, element )
    }
  }
}

func (c *LiveSet) DeviceParameters(){ 
  var elems []interface{}

  for _,element := range c.objects {
    //compare type track
    if reflect.TypeOf(element).Elem() == reflect.TypeOf(DeviceParameter{}) {
      fmt.Println(reflect.TypeOf(element).Elem())
      elems = append(elems, element )
    }
  }
}

func (c *LiveSet) ClipSlots(){
  var elems []interface{}

  for _,element := range c.objects {
    //compare type track
    if reflect.TypeOf(element).Elem() == reflect.TypeOf(ClipSlot{}) {
      fmt.Println(reflect.TypeOf(element).Elem())
      elems = append(elems, element )
    }
  }
}

func (c *LiveSet) Clips(){
  var elems []interface{}

  for _,element := range c.objects {
    //compare type track
    if reflect.TypeOf(element).Elem() == reflect.TypeOf(Clip{}) {
      fmt.Println(reflect.TypeOf(element).Elem())
      elems = append(elems, element )
    }
  }
}
// LIVE OBJECT

type LiveObject struct {
  id int32
  order int32
  //LiveSet
}

// OBJECTS

// Track
type Track struct {
  LiveObject
  is_master bool
  live_set *LiveSet
}

func NewTrack(id int32, is_master bool, order int32, live_set *LiveSet ) *Track {
  p := new(Track)
  p.id = id
  p.is_master = is_master
  p.live_set = live_set
  return p
}



// Clip
type Clip struct {
  LiveObject
}

// Clip slot
type ClipSlot struct {
  LiveObject
}

// Device
type Device struct {
  LiveObject
}

// Device parameter
type DeviceParameter struct {
  LiveObject
}

// Mixer device
type MixerDevice struct {
  LiveObject
}
