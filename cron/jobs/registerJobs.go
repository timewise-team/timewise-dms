package jobs

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/robfig/cron/v3"
	"github.com/timewise-team/timewise-models/models"
	"gopkg.in/gomail.v2"
	"io/ioutil"
	"net/http"
	"time"
)

// RegisterJobs registers all cron jobs
func RegisterJobs() {
	c := cron.New()

	// Add a sample job
	_, err := c.AddFunc("@every 1m", func() { sendNotification() })
	if err != nil {
		fmt.Println("Error adding cron job:", err)
		return
	}

	c.Start()
	fmt.Println("Cron jobs started")

	// Keep the program running
	select {}
}

func sendNotification() {
	fmt.Println("Starting cron job: sendNotification at", time.Now())

	unsentNotifications, err := GetUnsentNotifications()
	if err != nil {
		fmt.Println("Error getting unsent notifications:", err)
		return
	}

	for _, notification := range unsentNotifications {
		if notification.NotifiedAt != nil && notification.NotifiedAt.Before(time.Now()) {
			// Send notification
			fmt.Printf("Sending notification ID %d to email %d", notification.ID, notification.UserEmail.Email)

			// Send email
			err := SendEmail(notification.UserEmail.Email, "Notification", notification.Message)
			if err != nil {
				fmt.Println("Error sending email:", err)
				continue
			}

			// Update notification to sent
			err = updateNotificationToSent(notification.ID)
			if err != nil {
				fmt.Println("Error updating notification to sent:", err)
			}
		}
	}
}

func GetUnsentNotifications() ([]models.TwNotifications, error) {
	resp, err := http.Get("https://dms.timewise.space/dbms/v1/notification")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var notifications []models.TwNotifications
	if err := json.Unmarshal(body, &notifications); err != nil {
		return nil, err
	}

	return notifications, nil
}

func updateNotificationToSent(notificationID int) error {
	url := fmt.Sprintf("https://dms.timewise.space/dbms/v1/notification/%d", notificationID)
	req, err := http.NewRequest(http.MethodPut, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update notification: status code %d", resp.StatusCode)
	}

	return nil
}

func SendEmail(to string, subject string, body string) error {
	dialer := ConfigSMTP()
	if dialer == nil {
		return errors.New("failed to configure SMTP dialer")
	}
	m := gomail.NewMessage()
	m.SetHeader("From", dialer.Username)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	if err := dialer.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
func ConfigSMTP() *gomail.Dialer {
	SmtpHost := "smtp.gmail.com"
	SmtpPort := 587
	SmtpEmail := "timewise.space@gmail.com"
	SmtpPassword := "dczt wlvd eisn cixf"
	return gomail.NewDialer(SmtpHost, SmtpPort, SmtpEmail, SmtpPassword)
}
