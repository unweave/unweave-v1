# Unweave

The architecture for Unweave is broken into:
- `Services` at the business logic level,
- `Router` at the API level, which forwards requests to the appropriate `Service`
- `Conductor` at the orchestration level. 

`Services` are responsible for the domain logic. Each Service defines a `Driver`interface 
that can be swapped out based on the provider. Drivers implement different behaviors
for the domain object. For example, An exec can be scheduled on a bare VM, a container 
inside a VM, in a container in a kubernetes pod, etc. Drivers allow different implementations
to coexist. Services implement domain level workflows while drivers implement behaviors. 

The `Router` parses the provider appropriate for serving a request and forwards the
request to the appropriate `Service` implementation.

The `Conductor` is responsible for managing the underlying implementation and 
orchestration of compute resources. It is independent of the `Services` and `Router`. It
orchestrates a Pool of nodes per provider and assigns incoming `Container` requests to a
suitable node. It also manages the lifecycle of the nodes.


## Observers and Informers

Services implement `Informer`s and allow registering `Observer`s. Observers correspond to a 
specific Informer and can implement different behaviors. For example, a `ExecStateInformer`
watches for state changes in an exec and informs the registered observers. A `ExecRunningEmailer`
Observer can register to the `ExecStateInformer` and send out an email when the exec is running.

**N.B. Currently, the implementation of observers and informers assumes exactly one
instance of each service. To support distributed systems in the future, we'll need to
extract the `Observer` and `Informer` interfaces into shared/thread-safe implementations.**