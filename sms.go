package twilio

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type SMSResponse struct {
	AccountSid  string  `json:"account_sid"`
	ApiVersion  string  `json:"api_version"`
	Body        string  `json:"body"`
	DateCreated string  `json:"date_created,omitempty"`
	DateSent    string  `json:"date_sent,omitempty"`
	DateUpdated string  `json:"date_updated,omitempty"`
	Direction   string  `json:"direction"`
	From        string  `json:"from"`
	Price       float32 `json:"price,omitempty"`
	Sid         string  `json:"sid"`
	Status      string  `json:"status"`
	To          string  `json:"to"`
	Uri         string  `json:"uri"`
}

type SMSListResponse struct {
	Start           int    `json:"start"`
	Total           int    `json:"total"`
	NumPages        int    `json:"num_pages"`
	Page            int    `json:"page"`
	PageSize        int    `json:"page_size"`
	End             int    `json:"end"`
	Uri             string `json:"uri"`
	FirstPageUri    string `json:"first_page_uri"`
	LastPageUri     string `json:"last_page_uri"`
	NextPageUri     string `json:"next_page_uri"`
	PreviousPageUri string `json:"previous_page_uri"`
	SMSMessages     []SMSResponse
}

type SMSParams struct {
	StatusCallback string
	ApplicationSid string
}

func (t *Twilio) smsEndpoint() string {
	return fmt.Sprintf("%s/Accounts/%s/SMS/Messages", t.BaseUrl, t.AccountSid)
}

// Simple version of Send SMS with no optional parameters support.
func (t *Twilio) SimpleSendSMS(from, to, body string) (*SMSResponse, error) {
	return t.SendSMS(from, to, body, SMSParams{})
}

// Send SMS with more verbose options. It's support optional parameters.
//	StatusCallback : A URL that Twilio will POST to when your message is processed.
//	ApplicationSid : Twilio will POST `SMSSid` as well as other statuses to the URL in the `SMSStatusCallback` property of this application
func (t *Twilio) SendSMS(from, to, body string, p SMSParams) (s *SMSResponse, err error) {
	endpoint := fmt.Sprintf("%s.%s", t.smsEndpoint(), apiFormat)
	params := url.Values{}
	params.Set("From", from)
	params.Set("To", to)
	params.Set("Body", body)

	if p.StatusCallback != "" {
		params.Set("StatusCallback", p.StatusCallback)
	}

	if p.ApplicationSid != "" {
		params.Set("ApplicationSid", p.ApplicationSid)
	}

	b, status, err := t.post(endpoint, params)
	if err != nil {
		return
	}

	if status != http.StatusCreated {
		e := new(Exception)
		err = json.Unmarshal(b, &e)
		if err != nil {
			return
		}

		return nil, e
	}

	err = json.Unmarshal(b, &s)
	if err != nil {
		return nil, err
	}

	return
}

func (t *Twilio) GetSMS(sid string) (s *SMSResponse, err error) {
	endpoint := fmt.Sprintf("%s/%s.%s", t.smsEndpoint(), sid, apiFormat)

	b, status, err := t.get(endpoint, url.Values{})
	if err != nil {
		return
	}

	if status != http.StatusOK {
		e := new(Exception)
		err = json.Unmarshal(b, &e)
		if err != nil {
			return
		}

		return nil, e
	}

	err = json.Unmarshal(b, &s)
	if err != nil {
		return nil, err
	}

	return
}

// Returns a list of SMS messages associates with your account. It's support list filters via `map[string]string`:
//	"To" : Only show SMS messages to this phone number
//	"From" : Only show SMS messages from this phone number
//	"DateSent" : Only show SMS messages sent on this date (in GMT format), given as `YYYY-MM-DD`.
func (t *Twilio) ListSMS(filters map[string]string) (sl *SMSListResponse, err error) {
	endpoint := fmt.Sprintf("%s.%s", t.smsEndpoint(), apiFormat)
	params := url.Values{}

	for key, value := range filters {
		params.Set(key, value)
	}

	b, status, err := t.get(endpoint, params)
	if err != nil {
		return
	}

	if status != http.StatusOK {
		e := new(Exception)
		err = json.Unmarshal(b, &e)
		if err != nil {
			return
		}

		return nil, e
	}

	err = json.Unmarshal(b, &sl)
	if err != nil {
		return nil, err
	}

	return
}
