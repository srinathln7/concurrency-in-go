package main

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

func TestSendSMS(t *testing.T) {
	sys := &system{
		src:       "127.0.0.1",
		directory: []string{"1234567890"},
		msgch:     make(chan []byte),
	}

	smsData := &sms{
		Sender:   "127.0.0.1",
		Receiver: "1234567890",
		Text:     "Test message",
	}

	expectedPayload, err := json.Marshal(smsData)
	if err != nil {
		t.Fatalf("error marshaling expected payload: %v", err)
	}

	err = sys.sendSMS(smsData)
	if err != nil {
		t.Fatalf("error sending SMS: %v", err)
	}

	select {
	case <-time.After(100 * time.Millisecond):
		t.Fatal("timeout while waiting for message to be sent")
	case receivedPayload := <-sys.msgch:
		if !reflect.DeepEqual(receivedPayload, expectedPayload) {
			t.Errorf("received payload does not match the expected payload")
		}
	}
}
