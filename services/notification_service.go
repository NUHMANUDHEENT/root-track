package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"roottrack-backend/models"
	"roottrack-backend/repositories"
	"roottrack-backend/utils"
	"strings"
	"time"
)

type NotificationService struct {
	RoutineRepo repositories.RoutineRepository
}

func (s *NotificationService) CheckAndNotify(user *models.User) {
	if user.ExpoPushToken == "" {
		fmt.Printf("User %s has no push token, skipping notification\n", user.ID)
		return
	}

	today := time.Now().UTC()
	weekday := strings.ToLower(today.Weekday().String())

	routines, err := s.RoutineRepo.GetActiveByUser(user.ID, today)
	if err != nil {
		fmt.Printf("Failed to fetch routines for user %s: %v\n", user.ID, err)
		return
	}

	todayRoutinesCount := 0
	for _, routine := range routines {
		switch routine.Recurrence {
		case "none":
			if utils.SameDate(routine.StartDate, today) {
				todayRoutinesCount++
			}
		case "daily":
			todayRoutinesCount++
		case "weekly":
			for _, day := range routine.DaysOfWeek {
				if strings.ToLower(day) == weekday {
					todayRoutinesCount++
					break
				}
			}
		}
	}

	if todayRoutinesCount > 0 {
		message := fmt.Sprintf("You have %d tasks scheduled for today!", todayRoutinesCount)
		err := s.SendPushNotification(user.ExpoPushToken, "Today's Tasks", message)
		if err != nil {
			fmt.Printf("Failed to send push notification to user %s: %v\n", user.ID, err)
		} else {
			fmt.Printf("Success: Sent notification to user %s (%d tasks)\n", user.ID, todayRoutinesCount)
		}
	} else {
		fmt.Printf("User %s has no tasks today, skipping notification\n", user.ID)
	}
}


func (s *NotificationService) SendPushNotification(token, title, body string) error {
	url := "https://exp.host/--/api/v2/push/send"

	data := map[string]interface{}{
		"to":    token,
		"title": title,
		"body":  body,
		"sound": "default",
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("expo API returned status: %s", resp.Status)
	}
	return nil
}
