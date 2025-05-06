package rpc

type RpcInterface interface {

	// NewRpcQueue creates a new rpc queue.
	//
	// rpcQueueName is the name of the queue.
	//
	// handle is the function that will be called when a message is received.
	// data is the message in json string format, and the function should return a string as a json format.
	NewRpcQueue(rpcQueueName string, handle func(data string) string)

	// SendRpcReq sends a request to the rpc queue.
	//
	// rpcQueueName is the name of the queue.
	//
	// request is the request in json string format.
	Req(rpcQueueName string, request string) string
}
