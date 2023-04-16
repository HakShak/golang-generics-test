package alias_test

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
	return SDKMessage[T]{}
}

// This interface is needed so we can accept an instance of a concrete message
// otherwise this has to be a type alias, which I have yet to get to work because we can't reference
// methods within the generic function
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

func NewConcreteMessage() ConcreteMessage {
	return &SDKMessage[*GeneratedMessage]{
		GeneratedReference: &GeneratedMessage{},
	}
}

type ConcreteMessage *SDKMessage[*GeneratedMessage]

func TestBoilerPlate(t *testing.T) {
	theCallback := func(msg ConcreteMessage) {
		t.Logf("\nType: %T\nData: %+v", msg, msg)
		msg.GeneratedReference.GeneratedConcreteAccessor()
	}

	var _ SDKMessageInterface[*GeneratedMessage] = (*SDKMessage[*GeneratedMessage])(nil)                         // compiler agrees that all functions are implemented
	var _ SDKMessageInterface[*GeneratedMessage] = (ConcreteMessage)(nil)                                        // however, doesn't agree when it's an alias
	RegisterWithSDKFactory[ConcreteMessage, *GeneratedMessage](NewConcreteMessage(), theCallback)                // therefore the type parameter doesn't agree either
	RegisterWithSDKFactory[*SDKMessage[*GeneratedMessage], *GeneratedMessage](NewConcreteMessage(), theCallback) // But if we change the first type parameter away from the alias, the callback is a mismatch, which defeats the purpose altogether.
}
