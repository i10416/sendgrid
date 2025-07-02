package sendgrid

import (
	"context"
	"time"
)

// Email represents an email address with optional name
type Email struct {
	Email string `json:"email"`
	Name  string `json:"name,omitempty"`
}

// Content represents email content with type and value
type Content struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

// Attachment represents an email attachment
type Attachment struct {
	Content     string `json:"content"`
	Type        string `json:"type,omitempty"`
	Filename    string `json:"filename"`
	Disposition string `json:"disposition,omitempty"`
	ContentID   string `json:"content_id,omitempty"`
}

// Personalization contains personalized data for recipients
type Personalization struct {
	To                  []*Email               `json:"to,omitempty"`
	Cc                  []*Email               `json:"cc,omitempty"`
	Bcc                 []*Email               `json:"bcc,omitempty"`
	Subject             string                 `json:"subject,omitempty"`
	Headers             map[string]string      `json:"headers,omitempty"`
	Substitutions       map[string]string      `json:"substitutions,omitempty"`
	DynamicTemplateData map[string]interface{} `json:"dynamic_template_data,omitempty"`
	CustomArgs          map[string]string      `json:"custom_args,omitempty"`
	SendAt              int64                  `json:"send_at,omitempty"`
}

// MailSettings contains various mail settings
type MailSettings struct {
	BypassListManagement        *Setting          `json:"bypass_list_management,omitempty"`
	BypassSpamManagement        *Setting          `json:"bypass_spam_management,omitempty"`
	BypassBounceManagement      *Setting          `json:"bypass_bounce_management,omitempty"`
	BypassUnsubscribeManagement *Setting          `json:"bypass_unsubscribe_management,omitempty"`
	Footer                      *FooterSetting    `json:"footer,omitempty"`
	SandBoxMode                 *Setting          `json:"sandbox_mode,omitempty"`
	SpamCheck                   *SpamCheckSetting `json:"spam_check,omitempty"`
}

// Setting represents a boolean setting
type Setting struct {
	Enable *bool `json:"enable,omitempty"`
}

// FooterSetting represents footer settings
type FooterSetting struct {
	Enable *bool  `json:"enable,omitempty"`
	Text   string `json:"text,omitempty"`
	HTML   string `json:"html,omitempty"`
}

// SpamCheckSetting represents spam check settings
type SpamCheckSetting struct {
	Enable    *bool  `json:"enable,omitempty"`
	Threshold int    `json:"threshold,omitempty"`
	PostToURL string `json:"post_to_url,omitempty"`
}

// TrackingSettings contains various tracking settings
type TrackingSettings struct {
	ClickTracking        *ClickTrackingSetting        `json:"click_tracking,omitempty"`
	OpenTracking         *OpenTrackingSetting         `json:"open_tracking,omitempty"`
	SubscriptionTracking *SubscriptionTrackingSetting `json:"subscription_tracking,omitempty"`
	GoogleAnalytics      *GoogleAnalyticsSetting      `json:"ganalytics,omitempty"`
}

// ClickTrackingSetting represents click tracking settings
type ClickTrackingSetting struct {
	Enable     *bool `json:"enable,omitempty"`
	EnableText *bool `json:"enable_text,omitempty"`
}

// OpenTrackingSetting represents open tracking settings
type OpenTrackingSetting struct {
	Enable          *bool  `json:"enable,omitempty"`
	SubstitutionTag string `json:"substitution_tag,omitempty"`
}

// SubscriptionTrackingSetting represents subscription tracking settings
type SubscriptionTrackingSetting struct {
	Enable          *bool  `json:"enable,omitempty"`
	Text            string `json:"text,omitempty"`
	HTML            string `json:"html,omitempty"`
	SubstitutionTag string `json:"substitution_tag,omitempty"`
}

// GoogleAnalyticsSetting represents Google Analytics settings
type GoogleAnalyticsSetting struct {
	Enable      *bool  `json:"enable,omitempty"`
	UTMSource   string `json:"utm_source,omitempty"`
	UTMMedium   string `json:"utm_medium,omitempty"`
	UTMTerm     string `json:"utm_term,omitempty"`
	UTMContent  string `json:"utm_content,omitempty"`
	UTMCampaign string `json:"utm_campaign,omitempty"`
}

// ReplyTo represents reply-to settings
type ReplyTo struct {
	Email string `json:"email"`
	Name  string `json:"name,omitempty"`
}

// ReplyToList represents a list of reply-to addresses
type ReplyToList struct {
	Email string `json:"email"`
	Name  string `json:"name,omitempty"`
}

// InputSendMail represents the request body for sending mail
type InputSendMail struct {
	Personalizations []*Personalization `json:"personalizations"`
	From             *Email             `json:"from"`
	ReplyTo          *ReplyTo           `json:"reply_to,omitempty"`
	ReplyToList      []*ReplyToList     `json:"reply_to_list,omitempty"`
	Subject          string             `json:"subject,omitempty"`
	Content          []*Content         `json:"content,omitempty"`
	Attachments      []*Attachment      `json:"attachments,omitempty"`
	TemplateID       string             `json:"template_id,omitempty"`
	Headers          map[string]string  `json:"headers,omitempty"`
	Categories       []string           `json:"categories,omitempty"`
	CustomArgs       map[string]string  `json:"custom_args,omitempty"`
	SendAt           int64              `json:"send_at,omitempty"`
	BatchID          string             `json:"batch_id,omitempty"`
	ASM              *ASM               `json:"asm,omitempty"`
	IPPoolName       string             `json:"ip_pool_name,omitempty"`
	MailSettings     *MailSettings      `json:"mail_settings,omitempty"`
	TrackingSettings *TrackingSettings  `json:"tracking_settings,omitempty"`
}

// ASM represents Advanced Suppression Manager settings
type ASM struct {
	GroupID         int   `json:"group_id"`
	GroupsToDisplay []int `json:"groups_to_display,omitempty"`
}

// OutputSendMail represents the response from sending mail
type OutputSendMail struct {
	MessageID string `json:"message-id,omitempty"`
}

// NewEmail creates a new Email struct
func NewEmail(email, name string) *Email {
	return &Email{
		Email: email,
		Name:  name,
	}
}

// NewContent creates a new Content struct
func NewContent(contentType, value string) *Content {
	return &Content{
		Type:  contentType,
		Value: value,
	}
}

// NewPersonalization creates a new Personalization struct
func NewPersonalization() *Personalization {
	return &Personalization{}
}

// AddTo adds a recipient to personalization
func (p *Personalization) AddTo(email *Email) {
	p.To = append(p.To, email)
}

// AddCc adds a CC recipient to personalization
func (p *Personalization) AddCc(email *Email) {
	p.Cc = append(p.Cc, email)
}

// AddBcc adds a BCC recipient to personalization
func (p *Personalization) AddBcc(email *Email) {
	p.Bcc = append(p.Bcc, email)
}

// SetSendAt sets the send time for personalization
func (p *Personalization) SetSendAt(sendAt time.Time) {
	p.SendAt = sendAt.Unix()
}

// NewInputSendMail creates a new InputSendMail struct
func NewInputSendMail() *InputSendMail {
	return &InputSendMail{}
}

// SetFrom sets the from email address
func (m *InputSendMail) SetFrom(from *Email) {
	m.From = from
}

// SetSubject sets the email subject
func (m *InputSendMail) SetSubject(subject string) {
	m.Subject = subject
}

// AddPersonalization adds a personalization to the mail
func (m *InputSendMail) AddPersonalization(personalization *Personalization) {
	m.Personalizations = append(m.Personalizations, personalization)
}

// AddContent adds content to the mail
func (m *InputSendMail) AddContent(content *Content) {
	m.Content = append(m.Content, content)
}

// AddAttachment adds an attachment to the mail
func (m *InputSendMail) AddAttachment(attachment *Attachment) {
	m.Attachments = append(m.Attachments, attachment)
}

// SetTemplateID sets the template ID for the mail
func (m *InputSendMail) SetTemplateID(templateID string) {
	m.TemplateID = templateID
}

// SetSendAt sets the send time for the mail
func (m *InputSendMail) SetSendAt(sendAt time.Time) {
	m.SendAt = sendAt.Unix()
}

// AddCategory adds a category to the mail
func (m *InputSendMail) AddCategory(category string) {
	m.Categories = append(m.Categories, category)
}

// SendMail sends an email using SendGrid's mail/send API
// see: https://www.twilio.com/docs/sendgrid/api-reference/mail-send/mail-send
func (c *Client) SendMail(ctx context.Context, input *InputSendMail) (*OutputSendMail, error) {
	path := "/mail/send"

	req, err := c.NewRequest("POST", path, input)
	if err != nil {
		return nil, err
	}

	r := new(OutputSendMail)
	if err := c.Do(ctx, req, &r); err != nil {
		return nil, err
	}

	return r, nil
}
