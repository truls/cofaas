package helloworld_exports

// #include "helloworld_exports.h"
import "C"

import "unsafe"

// hello-world
type HelloWorldHelloRequest struct {
  Name string
}

type HelloWorldHelloReply struct {
  Message string
}

var hello_world HelloWorld = nil
func SetHelloWorld(i HelloWorld) {
  hello_world = i
}
type HelloWorld interface {
  SayHello(reqyest HelloWorldHelloRequest) Result[HelloWorldHelloReply, int32] 
}
//export hello_world_say_hello
func HelloWorldSayHello(reqyest *C.hello_world_hello_request_t, ret *C.hello_world_hello_reply_t, err *C.int32_t) bool {
  defer C.hello_world_hello_request_free(reqyest)
  var lift_reqyest HelloWorldHelloRequest
  var lift_reqyest_Name string
  lift_reqyest_Name = C.GoStringN(reqyest.name.ptr, C.int(reqyest.name.len))
  lift_reqyest.Name = lift_reqyest_Name
  result := hello_world.SayHello(lift_reqyest)
  if result.IsOk() {
    lower_result_ptr := (*C.hello_world_hello_reply_t)(unsafe.Pointer(ret))
    var lower_result_val C.hello_world_hello_reply_t
    var lower_result_val_message C.helloworld_exports_string_t
    
    lower_result_val_message.ptr = C.CString(result.Unwrap().Message)
    lower_result_val_message.len = C.size_t(len(result.Unwrap().Message))
    lower_result_val.message = lower_result_val_message
    *lower_result_ptr = lower_result_val
  } else {
    lower_result_ptr := (*C.int32_t)(unsafe.Pointer(err))
    lower_result_val := C.int32_t(result.UnwrapErr())
    *lower_result_ptr = lower_result_val
  }
  return result.IsOk()
}
// helloworld-exports
var helloworld_exports HelloworldExports = nil
func SetHelloworldExports(i HelloworldExports) {
  helloworld_exports = i
}
type HelloworldExports interface {
  Main() 
}
//export helloworld_exports_main
func HelloworldExportsMain() {
  helloworld_exports.Main()
  
}
