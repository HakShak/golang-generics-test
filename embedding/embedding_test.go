package embedding_test

import (
	"testing"
)

type GeneratedInterface interface {
	GeneratedName() string
}

type GeneratedMessage struct{}

func (gm *GeneratedMessage) GeneratedName() string {
	return "generated name"
}

func (gm *GeneratedMessage) GeneratedConcreteAccessor() {}

type SDKMessage[T GeneratedInterface] struct {
	GeneratedReference T
}

func (sm *SDKMessage[T]) SDKFunctionality() T {
	return *new(T)
}

func (sm *SDKMessage[T]) GeneratedName() string {
	return sm.GeneratedReference.GeneratedName()
}

func (sm *SDKMessage[T]) GetGeneratedMessage() T {
	return sm.GeneratedReference
}

func (sm *SDKMessage[T]) SetGeneratedMessage(msg T) {
	sm.GeneratedReference = msg
}

func (sm *SDKMessage[T]) FactoryNew() interface{} {
	panic("not implemented")
}

// This interface is needed so we can accept an instance of a concrete message
// otherwise this has to be a type alias, which I have yet to get to work
type SDKMessageInterface[T GeneratedInterface] interface {
	GeneratedInterface
	SDKFunctionality() T
	GetGeneratedMessage() T
	SetGeneratedMessage(msg T)
	FactoryNew() interface{}
}

func RegisterWithSDKFactory[T SDKMessageInterface[S], S GeneratedInterface](msg T, callback func(newMsg T)) {
	// Check the name for logic/key, just for example
	msg.GeneratedName()

	// transfer reference
	newGeneratedReference := msg.GetGeneratedMessage()

	// Do SDK things so the end user doesn't have to
	msg.SDKFunctionality()

	newT := msg.FactoryNew()
	concreteMessage := newT.(T)
	concreteMessage.SetGeneratedMessage(newGeneratedReference)
	callback(concreteMessage)
}

func NewConcreteMessage() *ConcreteMessage {
	return &ConcreteMessage{
		&SDKMessage[*GeneratedMessage]{
			GeneratedReference: &GeneratedMessage{},
		},
	}
}

type ConcreteMessage struct {
	*SDKMessage[*GeneratedMessage]
}

func (cm *ConcreteMessage) FactoryNew() interface{} {
	return NewConcreteMessage()
}

func TestBoilerPlate(t *testing.T) {
	theCallback := func(msg *ConcreteMessage) {
		t.Logf("\nType: %T\nData: %+v", msg, msg)
		msg.GeneratedReference.GeneratedConcreteAccessor()
	}

	// This works, but there are major drawbacks.
	// 1. Very verbose for the developer because type inference in this scenario simply doesn't work, so both type parameters have to be explicit.
	// 2. Each concrete message needs a factory. That's more boilerplate than already exists. So it defeats its own purpose. :(
	RegisterWithSDKFactory[*ConcreteMessage, *GeneratedMessage](NewConcreteMessage(), theCallback)
}
