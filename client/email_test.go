package client

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	email "github.com/microapis/email-api"
	emailClient "github.com/microapis/email-api/client"
	"github.com/microapis/messages-core/message"

	"github.com/oklog/ulid"
)

func before() (string, string, error) {
	host := os.Getenv("HOST")
	if host == "" {
		err := fmt.Errorf(fmt.Sprintf("Create: missing env variable HOST, failed with %s value", host))
		return "", "", err
	}

	port := os.Getenv("PORT")
	if port == "" {
		err := fmt.Errorf(fmt.Sprintf("Create: missing env variable PORT, failed with %s value", port))
		return "", "", err
	}

	return host, port, nil
}

// TestSend ...
func TestSend(t *testing.T) {
	host, port, err := before()
	if err != nil {
		t.Errorf(err.Error())
	}

	es, err := emailClient.New(host + ":" + port)
	if err != nil {
		log.Fatalln(err)
	}

	randomUUID := uuid.New()

	e := &email.Message{
		From:     "fake_email_" + randomUUID.String() + "@pensionatebien.cl",
		FromName: "fake_user",
		To:       []string{"camilo.aku@gmail.com"},
		Subject:  "fake_subject",
		Text:     "fake_text",
		Provider: "sendgrid",
	}

	id, err := es.Send(e, 3)
	if err != nil {
		t.Errorf("TestSend: err %v", err)
		return
	}

	_, err = ulid.Parse(id)
	if err != nil {
		t.Errorf("TestSend: err %v", err)
		return
	}

	msg, err := es.Get(id)
	if err != nil {
		t.Errorf("TestSend: err %v", err)
		return
	}

	if msg.Status != message.Pending {
		t.Errorf("TestSend: invalid status=%s", msg.Status)
		return
	}

	time.Sleep(time.Duration(5) * time.Second)

	msg, err = es.Get(id)
	if err != nil {
		t.Errorf("TestSend: err %v", err)
		return
	}

	if msg.Status != message.Sent {
		t.Errorf("TestSend: invalid status=%s", msg.Status)
		return
	}
}

// TestCancel ...
func TestCancel(t *testing.T) {
	host, port, err := before()
	if err != nil {
		t.Errorf(err.Error())
	}

	es, err := emailClient.New(host + ":" + port)
	if err != nil {
		log.Fatalln(err)
	}

	randomUUID := uuid.New()

	e := &email.Message{
		From:     "fake_email_" + randomUUID.String() + "@pensionatebien.cl",
		FromName: "fake_user",
		To:       []string{"camilo.aku@gmail.com"},
		Subject:  "fake_subject",
		Text:     "fake_text",
		Provider: "sendgrid",
	}

	id, err := es.Send(e, 3)
	if err != nil {
		t.Errorf("TestCancel: err %v", err)
		return
	}

	err = es.Cancel(id)
	if err != nil {
		t.Errorf("TestCancel: err %v", err)
		return
	}

	time.Sleep(time.Duration(5) * time.Second)

	msg, err := es.Get(id)
	if err != nil {
		t.Errorf("TestCancel: err %v", err)
		return
	}

	if msg.Status != message.Cancelled {
		t.Errorf("TestCancel: invalid status=%s", msg.Status)
		return
	}
}

// TestUpdate ...
func TestUpdate(t *testing.T) {
	host, port, err := before()
	if err != nil {
		t.Errorf(err.Error())
	}

	es, err := emailClient.New(host + ":" + port)
	if err != nil {
		log.Fatalln(err)
	}

	randomUUID := uuid.New()

	e := &email.Message{
		From:     "fake_email_" + randomUUID.String() + "@pensionatebien.cl",
		FromName: "fake_user",
		To:       []string{"camilo.aku@gmail.com"},
		Subject:  "fake_subject",
		Text:     "fake_text",
		Provider: "sendgrid",
	}

	id, err := es.Send(e, 3)
	if err != nil {
		t.Errorf("TestUpdate: err %v", err)
		return
	}

	expectedFrom := "changed_fake_from@pensionatebien.cl"
	e.From = expectedFrom

	expectedFromName := "changed_from_name"
	e.FromName = expectedFromName

	expectedTo := []string{"changes_fake_to@gmail.com"}
	e.To = expectedTo

	expectedSubject := "changed_fake_subject"
	e.Subject = expectedSubject

	expectedText := "changed_fake_text"
	e.Text = expectedText

	expectedHTML := "changed_fake_html"
	e.HTML = expectedHTML

	expectedProvider := "sendgrid"
	e.Provider = expectedProvider

	err = es.Update(id, e)
	if err != nil {
		t.Errorf("TestUpdate: err %v", err)
		return
	}

	msg, err := es.Get(id)
	if err != nil {
		t.Errorf("TestUpdate: err %v", err)
		return
	}

	if msg.From != expectedFrom {
		t.Errorf("TestUpdate: invalid text=%s in content", msg.From)
		return
	}

	if msg.FromName != expectedFromName {
		t.Errorf("TestUpdate: invalid text=%s in content", msg.FromName)
		return
	}

	// TODO(ca): add support test to attribute with multiples emails
	if msg.To[0] != expectedTo[0] {
		t.Errorf("TestUpdate: invalid text=%s in content", msg.To[0])
		return
	}

	if msg.Subject != expectedSubject {
		t.Errorf("TestUpdate: invalid text=%s in content", msg.Subject)
		return
	}

	if msg.Text != expectedText {
		t.Errorf("TestUpdate: invalid text=%s in content", msg.Text)
		return
	}

	if msg.HTML != expectedHTML {
		t.Errorf("TestUpdate: invalid text=%s in content", msg.HTML)
		return
	}

	if msg.Provider != expectedProvider {
		t.Errorf("TestUpdate: invalid text=%s in content", msg.Provider)
		return
	}
}

// TODO(ca): add support test to attribute with multiples emails
// TODO(ca): implement test that verify when the message was sent and its status is "sent", so the message information cannot be changed.
// TODO(ca): implement test that verify when update de message but not the status attribute.
