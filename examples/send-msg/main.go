package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"
)

type system struct {
	src       string
	directory []string
	msgch     chan []byte
}

type sms struct {
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Text     string `json:"text"`
}

func (sys *system) sendSMS(sms *sms) error {
	fmt.Printf("Sending message \"%s\" from %s to %s\n",
		sms.Text,
		sms.Sender,
		sms.Receiver,
	)

	// simulate network delay for sending msg
	time.Sleep(50 * time.Millisecond)

	// message, exception, err := twilio.SendSMS(sms.Sender, sms.Receiver, sms.Text, "", "")
	// if err != nil {
	// 	return err
	// }
	// if exception != nil {
	// 	return fmt.Errorf("exception: %s", exception.Message)
	// }

	payload, err := json.Marshal(sms)
	if err != nil {
		return err
	}

	sys.msgch <- payload
	return nil
}

func (sys *system) sendMessages(dir []string) {

	// Create a Twilio client
	// accountSID := "your_twilio_account_sid"
	// authToken := "your_twilio_auth_token"
	// twilio := gotwilio.NewTwilioClient(accountSID, authToken)

	var wg sync.WaitGroup
	for _, num := range dir {
		wg.Add(1)
		go func(receiver string) {
			defer wg.Done()
			sms := &sms{
				Sender:   sys.src,
				Receiver: receiver,
				Text:     "initial test",
			}
			if err := sys.sendSMS(sms); err != nil {
				log.Fatalf("error sending message to %s: %v", receiver, err)
			}
		}(num)
	}

	// Wait for all goroutines to finish
	wg.Wait()
}

func (sys *system) recvSMS() <-chan sms {
	recvCh := make(chan sms)
	go func() {
		defer close(recvCh)
		payload := <-sys.msgch
		//var msg gotwilio.SmsResponse
		var msg sms
		if err := json.Unmarshal(payload, &msg); err != nil {
			log.Printf("error decoding message: %v", err)
			return
		}

		recvCh <- msg
	}()

	return recvCh
}

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	dir := []string{"0649571937", "0623789027", "0634181723", "0676542312", "06234787125"}
	system := &system{
		src:       "127.0.0.1",
		directory: dir,
		msgch:     make(chan []byte),
	}

	system.run(ctx)
}

func (system *system) run(ctx context.Context) {
	defer close(system.msgch)

	// Schedule sending msges else deadlock
	go system.sendMessages(system.directory)

	for i := 0; i < len(system.directory); i++ {
		log.Printf("iteration:%d", i+1)
		select {
		case <-ctx.Done():
			log.Fatal("operation timeout")
		case msg := <-system.recvSMS():
			fmt.Printf("received message \"%s\" from %s to %s\n",
				msg.Text,
				msg.Sender,
				msg.Receiver,
			)
		}
	}

	fmt.Println("All messages sent successfully")
}
