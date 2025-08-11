package service

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/bdlm/log"
)

type NotificationPayload struct {
	Level                string `json:"level"`
	EmployeeAbbreviation string `json:"employeeAbbreviation"`
	Message              string `json:"message"`
}

type Notifier struct {
	connection string
}

func NewNotifier(connection string) *Notifier {
	return &Notifier{
		connection: connection,
	}
}

// SendMessage sends a warning message if an employee has 3 or more computers assigned to them.
func (n *Notifier) SendMessage(employeeAbbreviation string) {
	payload := NotificationPayload{
		Level:                "warning",
		EmployeeAbbreviation: employeeAbbreviation,
		Message:              "There are 3 or more computers assigned to the same employee.",
	}

	notification, err := json.Marshal(payload)
	if err != nil {
		log.Error("failed to send message: failed to marshal payload: " + err.Error())
		return
	}

	notifyURL := n.connection + "/api/notify"
	resp, err := http.Post(notifyURL, "application/json", bytes.NewReader(notification))
	if err != nil {
		log.Error("failed to send message: failed to send POST request to /api/notify: " + err.Error())
		return
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Error("failed to close response body: " + err.Error())
		}
	}()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("failed to read response body: " + err.Error())
		return
	}

	log.Infof("/api/notify responded with status %d: %s", resp.StatusCode, string(bodyBytes))
}
