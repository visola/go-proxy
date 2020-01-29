package event

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type testAccumulatingResultConsumer struct {
	consumed []HandleResult
}

func (r *testAccumulatingResultConsumer) Consume(toConsume HandleResult) {
	r.consumed = append(r.consumed, toConsume)
}

func TestConsumersGetCalledCorrectly(t *testing.T) {
	// Given a result consumer
	accumulatingConsumer := &testAccumulatingResultConsumer{}

	// When adding to the list of listeners
	AddRequestListener(&testAccumulatingResultConsumer{})
	AddRequestListener(accumulatingConsumer)
	AddRequestListener(&testAccumulatingResultConsumer{})

	// And an event is emitted
	result1 := HandleResult{
		ID: uuid.New().String(),
	}
	EmitHandleResult(result1)

	// Then it should receive it
	assert.Equal(t, 1, len(accumulatingConsumer.consumed))
	assert.Equal(t, result1.ID, accumulatingConsumer.consumed[0].ID)

	// When a new event is emitted
	result2 := HandleResult{
		ID: uuid.New().String(),
	}
	EmitHandleResult(result2)

	// Then it should receive it
	assert.Equal(t, 2, len(accumulatingConsumer.consumed))
	assert.Equal(t, result1.ID, accumulatingConsumer.consumed[0].ID)
	assert.Equal(t, result2.ID, accumulatingConsumer.consumed[1].ID)

	// Given it is removed from the list of listeners
	RemoveRequestListener(accumulatingConsumer)

	// When an event is emitted
	result3 := HandleResult{
		ID: uuid.New().String(),
	}
	EmitHandleResult(result3)

	// Then it should not receive it
	assert.Equal(t, 2, len(accumulatingConsumer.consumed))
	assert.Equal(t, result1.ID, accumulatingConsumer.consumed[0].ID)
	assert.Equal(t, result2.ID, accumulatingConsumer.consumed[1].ID)
}
