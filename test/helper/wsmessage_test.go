package helper_test

import (
	"encoding/json"
	"real-time-chat-api/internal/helper"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWSMessage_Structure(t *testing.T) {
	msg := helper.WSMessage{
		Event:   "send_message",
		Payload: json.RawMessage(`{"content":"Hello World"}`),
	}

	assert.Equal(t, "send_message", msg.Event)
	assert.NotNil(t, msg.Payload)
}

func TestWSMessage_MarshalJSON(t *testing.T) {
	msg := helper.WSMessage{
		Event:   "new_message",
		Payload: json.RawMessage(`{"content":"Test"}`),
	}

	data, err := json.Marshal(msg)
	assert.NoError(t, err)
	assert.Contains(t, string(data), "new_message")
	assert.Contains(t, string(data), "Test")
}

func TestWSMessage_UnmarshalJSON(t *testing.T) {
	raw := `{"event":"user_joined","payload":{"user_id":"123"}}`

	var msg helper.WSMessage
	err := json.Unmarshal([]byte(raw), &msg)
	assert.NoError(t, err)
	assert.Equal(t, "user_joined", msg.Event)
	assert.NotNil(t, msg.Payload)
}

func TestWSMessage_UnmarshalJSON_InvalidJSON(t *testing.T) {
	raw := `{invalid json}`

	var msg helper.WSMessage
	err := json.Unmarshal([]byte(raw), &msg)
	assert.Error(t, err)
}

func TestWSMessage_EmptyEvent(t *testing.T) {
	msg := helper.WSMessage{
		Event:   "",
		Payload: json.RawMessage(`{}`),
	}

	assert.Equal(t, "", msg.Event)
}
