package exec

import "context"

type informer interface {
	inform(ctx context.Context)
	register(o observer)
	unregister(o observer)
	watch()
}

type stateInformer struct {
	store     Store
	driver    Driver
	observers map[string]observer
}

func NewStateInformer(store Store, driver Driver) *stateInformer {
	return &stateInformer{
		store:     store,
		driver:    driver,
		observers: make(map[string]observer),
	}
}

func (i *stateInformer) inform(ctx context.Context) {

}

func (i *stateInformer) register(o observer) {
	if _, ok := i.observers[o.id()]; ok {
		return
	}
	i.observers[o.id()] = o
}

func (i *stateInformer) unregister(o observer) {
	if _, ok := i.observers[o.id()]; !ok {
		return
	}
	delete(i.observers, o.id())
}

func (i *stateInformer) watch() {

}
