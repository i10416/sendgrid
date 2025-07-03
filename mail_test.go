package sendgrid

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewEmail(t *testing.T) {
	email := NewEmail("test@example.com", "Test User")
	assert.Equal(t, "test@example.com", email.Email)
	assert.Equal(t, "Test User", email.Name)
}

func TestNewContent(t *testing.T) {
	content := NewContent("text/plain", "Hello, World!")
	assert.Equal(t, "text/plain", content.Type)
	assert.Equal(t, "Hello, World!", content.Value)
}

func TestNewPersonalization(t *testing.T) {
	p := NewPersonalization()
	assert.NotNil(t, p)
	assert.Empty(t, p.To)
	assert.Empty(t, p.Cc)
	assert.Empty(t, p.Bcc)
}

func TestPersonalizationAddTo(t *testing.T) {
	p := NewPersonalization()
	email := NewEmail("to@example.com", "To User")
	p.AddTo(email)

	assert.Len(t, p.To, 1)
	assert.Equal(t, "to@example.com", p.To[0].Email)
	assert.Equal(t, "To User", p.To[0].Name)
}

func TestPersonalizationAddCc(t *testing.T) {
	p := NewPersonalization()
	email := NewEmail("cc@example.com", "CC User")
	p.AddCc(email)

	assert.Len(t, p.Cc, 1)
	assert.Equal(t, "cc@example.com", p.Cc[0].Email)
	assert.Equal(t, "CC User", p.Cc[0].Name)
}

func TestPersonalizationAddBcc(t *testing.T) {
	p := NewPersonalization()
	email := NewEmail("bcc@example.com", "BCC User")
	p.AddBcc(email)

	assert.Len(t, p.Bcc, 1)
	assert.Equal(t, "bcc@example.com", p.Bcc[0].Email)
	assert.Equal(t, "BCC User", p.Bcc[0].Name)
}

func TestPersonalizationSetSendAt(t *testing.T) {
	p := NewPersonalization()
	sendTime := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)
	p.SetSendAt(sendTime)

	assert.Equal(t, sendTime.Unix(), p.SendAt)
}

func TestNewInputSendMail(t *testing.T) {
	mail := NewInputSendMail()
	assert.NotNil(t, mail)
	assert.Empty(t, mail.Personalizations)
	assert.Empty(t, mail.Content)
	assert.Empty(t, mail.Attachments)
}

func TestInputSendMailSetFrom(t *testing.T) {
	mail := NewInputSendMail()
	from := NewEmail("from@example.com", "From User")
	mail.SetFrom(from)

	assert.Equal(t, "from@example.com", mail.From.Email)
	assert.Equal(t, "From User", mail.From.Name)
}

func TestInputSendMailSetSubject(t *testing.T) {
	mail := NewInputSendMail()
	subject := "Test Subject"
	mail.SetSubject(subject)

	assert.Equal(t, subject, mail.Subject)
}

func TestInputSendMailAddPersonalization(t *testing.T) {
	mail := NewInputSendMail()
	p := NewPersonalization()
	p.AddTo(NewEmail("to@example.com", "To User"))
	mail.AddPersonalization(p)

	assert.Len(t, mail.Personalizations, 1)
	assert.Len(t, mail.Personalizations[0].To, 1)
	assert.Equal(t, "to@example.com", mail.Personalizations[0].To[0].Email)
}

func TestInputSendMailAddContent(t *testing.T) {
	mail := NewInputSendMail()
	content := NewContent("text/plain", "Hello, World!")
	mail.AddContent(content)

	assert.Len(t, mail.Content, 1)
	assert.Equal(t, "text/plain", mail.Content[0].Type)
	assert.Equal(t, "Hello, World!", mail.Content[0].Value)
}

func TestInputSendMailAddAttachment(t *testing.T) {
	mail := NewInputSendMail()
	attachment := &Attachment{
		Content:  "dGVzdA==", // base64 encoded "test"
		Type:     "text/plain",
		Filename: "test.txt",
	}
	mail.AddAttachment(attachment)

	assert.Len(t, mail.Attachments, 1)
	assert.Equal(t, "dGVzdA==", mail.Attachments[0].Content)
	assert.Equal(t, "text/plain", mail.Attachments[0].Type)
	assert.Equal(t, "test.txt", mail.Attachments[0].Filename)
}

func TestInputSendMailSetTemplateID(t *testing.T) {
	mail := NewInputSendMail()
	templateID := "d-123456789"
	mail.SetTemplateID(templateID)

	assert.Equal(t, templateID, mail.TemplateID)
}

func TestInputSendMailSetSendAt(t *testing.T) {
	mail := NewInputSendMail()
	sendTime := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)
	mail.SetSendAt(sendTime)

	assert.Equal(t, sendTime.Unix(), mail.SendAt)
}

func TestInputSendMailAddCategory(t *testing.T) {
	mail := NewInputSendMail()
	category := "newsletter"
	mail.AddCategory(category)

	assert.Len(t, mail.Categories, 1)
	assert.Equal(t, category, mail.Categories[0])
}

// TestSendMail tests the SendMail function structure and helper functions
// Actual API testing is done in integration tests

// TestSendMailStructure tests building a complex email structure
func TestSendMailStructure(t *testing.T) {
	mail := NewInputSendMail()

	// Set from
	mail.SetFrom(NewEmail("from@example.com", "From User"))
	mail.SetSubject("Complex Test Email")

	// Add multiple personalizations
	p1 := NewPersonalization()
	p1.AddTo(NewEmail("to1@example.com", "To User 1"))
	p1.AddCc(NewEmail("cc1@example.com", "CC User 1"))
	p1.Subject = "Personalized Subject 1"
	p1.DynamicTemplateData = map[string]interface{}{
		"name":    "User 1",
		"product": "Product A",
	}
	mail.AddPersonalization(p1)

	p2 := NewPersonalization()
	p2.AddTo(NewEmail("to2@example.com", "To User 2"))
	p2.Subject = "Personalized Subject 2"
	p2.DynamicTemplateData = map[string]interface{}{
		"name":    "User 2",
		"product": "Product B",
	}
	mail.AddPersonalization(p2)

	// Add content
	mail.AddContent(NewContent("text/plain", "Hello, {{name}}! Check out {{product}}."))
	mail.AddContent(NewContent("text/html", "<h1>Hello, {{name}}!</h1><p>Check out {{product}}.</p>"))

	// Add attachment
	attachment := &Attachment{
		Content:     "VGhpcyBpcyBhIHRlc3QgYXR0YWNobWVudA==", // base64 encoded text
		Type:        "text/plain",
		Filename:    "test.txt",
		Disposition: "attachment",
	}
	mail.AddAttachment(attachment)

	// Add categories
	mail.AddCategory("newsletter")
	mail.AddCategory("test")

	// Set custom args
	mail.CustomArgs = map[string]string{
		"campaign_id": "12345",
		"user_id":     "67890",
	}

	// Set mail settings
	mail.MailSettings = &MailSettings{
		SandBoxMode: &Setting{Enable: Bool(true)},
	}

	// Set tracking settings
	mail.TrackingSettings = &TrackingSettings{
		ClickTracking: &ClickTrackingSetting{
			Enable:     Bool(true),
			EnableText: Bool(true),
		},
		OpenTracking: &OpenTrackingSetting{
			Enable: Bool(true),
		},
	}

	// Verify structure
	assert.Equal(t, "from@example.com", mail.From.Email)
	assert.Equal(t, "Complex Test Email", mail.Subject)
	assert.Len(t, mail.Personalizations, 2)
	assert.Len(t, mail.Content, 2)
	assert.Len(t, mail.Attachments, 1)
	assert.Len(t, mail.Categories, 2)
	assert.NotNil(t, mail.MailSettings)
	assert.NotNil(t, mail.TrackingSettings)
}

func TestSendMail(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/mail/send", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusAccepted)
		if _, err := fmt.Fprint(w, `{"message":"success"}`); err != nil {
			t.Fatal(err)
		}
	})

	mail := NewInputSendMail()
	mail.SetFrom(NewEmail("from@example.com", "From User"))
	mail.SetSubject("Test Email")

	p := NewPersonalization()
	p.AddTo(NewEmail("to@example.com", "To User"))
	mail.AddPersonalization(p)

	mail.AddContent(NewContent("text/plain", "Hello, World!"))

	result, err := client.SendMail(context.TODO(), mail)
	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestSendMail_Failed(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/mail/send", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		if _, err := fmt.Fprint(w, `{"errors":[{"message":"Invalid email address"}]}`); err != nil {
			t.Fatal(err)
		}
	})

	mail := NewInputSendMail()
	mail.SetFrom(NewEmail("invalid-email", "From User"))
	mail.SetSubject("Test Email")

	p := NewPersonalization()
	p.AddTo(NewEmail("to@example.com", "To User"))
	mail.AddPersonalization(p)

	mail.AddContent(NewContent("text/plain", "Hello, World!"))

	_, err := client.SendMail(context.TODO(), mail)
	assert.Error(t, err)
}

func TestSendMail_WithTemplate(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/mail/send", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		w.WriteHeader(http.StatusAccepted)
		if _, err := fmt.Fprint(w, `{"message":"success"}`); err != nil {
			t.Fatal(err)
		}
	})

	mail := NewInputSendMail()
	mail.SetFrom(NewEmail("from@example.com", "From User"))
	mail.SetTemplateID("d-123456789")

	p := NewPersonalization()
	p.AddTo(NewEmail("to@example.com", "To User"))
	p.DynamicTemplateData = map[string]interface{}{
		"name":    "John Doe",
		"product": "Amazing Product",
	}
	mail.AddPersonalization(p)

	result, err := client.SendMail(context.TODO(), mail)
	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestSendMail_WithAttachment(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/mail/send", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		w.WriteHeader(http.StatusAccepted)
		if _, err := fmt.Fprint(w, `{"message":"success"}`); err != nil {
			t.Fatal(err)
		}
	})

	mail := NewInputSendMail()
	mail.SetFrom(NewEmail("from@example.com", "From User"))
	mail.SetSubject("Email with Attachment")

	p := NewPersonalization()
	p.AddTo(NewEmail("to@example.com", "To User"))
	mail.AddPersonalization(p)

	mail.AddContent(NewContent("text/plain", "Please find the attachment."))

	attachment := &Attachment{
		Content:     "VGhpcyBpcyBhIHRlc3QgYXR0YWNobWVudA==", // base64 encoded
		Type:        "text/plain",
		Filename:    "test.txt",
		Disposition: "attachment",
		ContentID:   "test-content-id",
	}
	mail.AddAttachment(attachment)

	result, err := client.SendMail(context.TODO(), mail)
	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestSendMail_WithScheduledSend(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/mail/send", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		w.WriteHeader(http.StatusAccepted)
		if _, err := fmt.Fprint(w, `{"message":"success"}`); err != nil {
			t.Fatal(err)
		}
	})

	mail := NewInputSendMail()
	mail.SetFrom(NewEmail("from@example.com", "From User"))
	mail.SetSubject("Scheduled Email")

	// Set send time to 1 hour from now
	sendTime := time.Now().Add(time.Hour)
	mail.SetSendAt(sendTime)

	p := NewPersonalization()
	p.AddTo(NewEmail("to@example.com", "To User"))
	mail.AddPersonalization(p)

	mail.AddContent(NewContent("text/plain", "This is a scheduled email."))

	result, err := client.SendMail(context.TODO(), mail)
	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestSendMail_WithMultipleRecipients(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/mail/send", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		w.WriteHeader(http.StatusAccepted)
		if _, err := fmt.Fprint(w, `{"message":"success"}`); err != nil {
			t.Fatal(err)
		}
	})

	mail := NewInputSendMail()
	mail.SetFrom(NewEmail("from@example.com", "From User"))
	mail.SetSubject("Multiple Recipients Email")

	// First personalization
	p1 := NewPersonalization()
	p1.AddTo(NewEmail("to1@example.com", "To User 1"))
	p1.AddCc(NewEmail("cc1@example.com", "CC User 1"))
	p1.AddBcc(NewEmail("bcc1@example.com", "BCC User 1"))
	mail.AddPersonalization(p1)

	// Second personalization
	p2 := NewPersonalization()
	p2.AddTo(NewEmail("to2@example.com", "To User 2"))
	mail.AddPersonalization(p2)

	mail.AddContent(NewContent("text/plain", "Hello to all recipients!"))
	mail.AddContent(NewContent("text/html", "<h1>Hello to all recipients!</h1>"))

	result, err := client.SendMail(context.TODO(), mail)
	assert.NoError(t, err)
	assert.NotNil(t, result)
}
