# Node Split

## Problem

We currently have both nodes - worker and head nodes - living in the same codebase.
Biggest benefit of this setup is that we have a single executable - `node` which can operate as both the head and worker node, depending on passed CLI flags.
Both worker and head nodes are implemented on the same type - `Node`, which operates differently depending on the set node role.

Problems with this approach:

1. Functionality for both node roles is compiled and "exists" at the same time. Technically, worker node has the capability (functions/methods) to issue a roll call, and the head node has the capability to invoke executor and execute a function using Blockless Runtime. This does not happen because there's a number of checks performed in code that dictate how the node operates. 

2. Data for both nodes exists at the same time. There are fields that are needed for worker node operations, whereas we have other that are head-node specific. However, since the node is implemented as a single type, they both exist in the `Node` type, and it's a question of mental accounting and discipline to not reference a field that does not make sense for a given node role. Especially with long-running/detached executions, nodes started diverging even more, with more and more node-specific fields.

3. Mixed and uneeded dependencies. Worker node requires the `executor` component to invoke Blockless Runtime, and the `fstore` component to keep accounting of installed functions. Head node needs none of these. We currently must handle all of these and make sure that they are present for the given node type, while absent for the other.

4. Message type handling. We currently have a mixed bag of messages. For example, we use the same messages `MsgExecute` and `MsgExecuteResponse` for both head nodes and worker nodes executions and execution responses. These messages have different content for different roles, and we have to shoehorn data into messages that might not fit either use case fully. For example, head node returns execution results from worker nodes in the format of `{ "peer1": ..., "peer2: ..., "peer3": ...}`; so the worker nodes must also fit into this mold and introduce unnecessary indirection by sending its single result as `{ "peer1": ..." }`.

## Proposed Changes

Splitting the `Node` type to two different types - `Worker` and `HeadNode` and a third one that acts as a mandatory dependency - `core`. All three get implemented in separate packages and can evolve independendtly from each other.

From the user perspective nothing changes because we still have the `node` executable as the one single thing.

Node Core is an interface that provides all shared node functionality. It gets embedded in both `Worker` and `HeadNode` types. Node core provides logging, telemetry and networking faculties.

Head node for one becomes super slim now - it needs to take care of roll calls and executions and that's basically it.
Worker node inherits a large chunk of the codebase - function install, executions, attributes etc...

Both types only hold their respective states.

By having these types separate and specialized, the code becomes clearer.
There's no code for the `Worker` type that shouldn't be there and should not be executed, and you don't have to check everything constantly to prevent misuse.

Also, messages that are invalid for a given node can be safely ignored.
For example, it's wrong for a head node to process a roll call message.
No need to check it in code - head node does not have a handler for it => message gets ignored.

Second - we can split some messages and have cleaner and clearer models.
Execution process can is split to two messages: `MsgExecute` is the execution request as sent to the head node.
The head node, instead of using the same message for the worker, sends a `MsgWorkOrder` to the worker node.
Bigger benefit though comes from the `*Response` types - `MsgExecuteResponse` can be formatted for what it is, a collection of execution results, vs `MsgWorkOrderResponse` can be what it is - a response from a single node.

It's quite possible that even a large part of the node main loop, the `Run(context.Context) error` method - can be abstracted away and handled by the node core, as well as health ping, peer discovery and a lot of other things.


```go
type Core interface {
	Logger
	Network
	Telemetry
}

type Logger interface {
	Log() *zerolog.Logger
}

type Network interface {
	Host() *host.Host
	Messaging
}

type Messaging interface {
	Connected(peer.ID) bool

	Send(context.Context, peer.ID, blockless.Message) error
	SendToMany(context.Context, []peer.ID, blockless.Message, bool) error

	JoinTopic(string) error
	Subscribe(context.Context, string) error
	Publish(context.Context, blockless.Message) error
	PublishToTopic(context.Context, string, blockless.Message) error
}

type Telemetry interface {
	Tracer() *tracing.Tracer
	Metrics() *metrics.Metrics
}
```

