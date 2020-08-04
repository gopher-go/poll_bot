package poll_bot

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const message = `{
	"event":"message",
	"timestamp":1457764197627,
	"message_token":4912661846655238145,
	"sender":{
		 "id":"01234567890A=",
		 "name":"John McClane",
		 "avatar":"http://avatar.example.com",
		 "country":"UK",
		 "language":"en",
		 "api_version":1
	},
	"message":{
		 "type":"text",
		 "text":"a message to the service",
		 "media":"http://example.com",
		 "location":{
				"lat":50.76891,
				"lon":6.11499
		 },
		 "tracking_data":"tracking data"
	}
}`

func TestParseViberMessage(t *testing.T) {
	m, err := parseCallback([]byte(message))
	require.NoError(t, err)
	require.Equal(t, m.Event, "message")
	require.Equal(t, m.Message.Type, "text")
	require.Equal(t, m.Message.Text, "a message to the service")
	require.Equal(t, m.User.ID, "01234567890A=")
	require.Equal(t, m.User.Name, "John McClane")
	require.Equal(t, m.User.Country, "UK")
}

const subscribeMessage = `{
	"event":"subscribed",
	"timestamp":1457764197627,
	"user":{
		 "id":"01234567890A=",
		 "name":"John McClane",
		 "avatar":"http://avatar.example.com",
		 "country":"UK",
		 "language":"en",
		 "api_version":1
	},
	"message_token":4912661846655238145
}`

func TestParseViberSubscribe(t *testing.T) {
	m, err := parseCallback([]byte(subscribeMessage))
	require.NoError(t, err)
	require.Equal(t, m.Event, "subscribed")
	require.Equal(t, m.User.ID, "01234567890A=")
	require.Equal(t, m.User.Name, "John McClane")
	require.Equal(t, m.User.Country, "UK")
}

const deliveryCallback = `{"event":"delivered","timestamp":1595324677891,"chat_hostname":"SN-CHAT-02_","message_token":5466297578174182639,"user_id":"3SQNc4FPUQmysaM/AQEhXg=="}`

func TestParseViberDeliveryCallback(t *testing.T) {
	m, err := parseCallback([]byte(deliveryCallback))
	require.NoError(t, err)
	require.Equal(t, m.Event, "delivered")
	require.Equal(t, m.User.ID, "3SQNc4FPUQmysaM/AQEhXg==")
}

const seenCallback = `{"event":"seen","timestamp":1595324652811,"chat_hostname":"SN-CHAT-02_","message_token":5466297467809468529,"user_id":"3SQNc4FPUQmysaM/AQEhXg=="}`

func TestParseViberSeenCallback(t *testing.T) {
	m, err := parseCallback([]byte(seenCallback))
	require.NoError(t, err)
	require.Equal(t, m.Event, "seen")
	require.Equal(t, m.User.ID, "3SQNc4FPUQmysaM/AQEhXg==")
}

const conversationStarted = `{"event":"conversation_started","context":"context-source","timestamp":1595339165214,"chat_hostname":"SN-376_","message_token":5466358343372472277,"type":"open","user":{"id":"3SQNc4FPUQmysaM/AQEhXg==","name":"Georgy Buranov","avatar":"","language":"en-CA","country":"DE","api_version":8},"subscribed":false}`

func TestParseConversationStartedCallback(t *testing.T) {
	m, err := parseCallback([]byte(conversationStarted))
	require.NoError(t, err)
	require.Equal(t, "conversation_started", m.Event)
	require.Equal(t, "3SQNc4FPUQmysaM/AQEhXg==", m.User.ID)
	require.Equal(t, "Georgy Buranov", m.User.Name)
	require.Equal(t, "DE", m.User.Country)
	require.Equal(t, "context-source", m.Context)
}

const unsubscribedMessage = `{"event":"unsubscribed","timestamp":1595347885535,"chat_hostname":"SN-376_","user_id":"3SQNc4FPUQmysaM/AQEhXg==","message_token":5466394919049723652}`

func TestUnsubscribedMessage(t *testing.T) {
	m, err := parseCallback([]byte(unsubscribedMessage))
	require.NoError(t, err)
	require.Equal(t, m.Event, "unsubscribed")
	require.Equal(t, m.User.ID, "3SQNc4FPUQmysaM/AQEhXg==")
}

const webhookMessage = `{"event":"webhook","timestamp":1595347748566,"chat_hostname":"SN-CHAT-22_","message_token":5466394344559121456}`

func TestWebhookMessage(t *testing.T) {
	m, err := parseCallback([]byte(webhookMessage))
	require.NoError(t, err)
	require.Equal(t, m.Event, "webhook")
	require.Equal(t, m.User.ID, "")
}
