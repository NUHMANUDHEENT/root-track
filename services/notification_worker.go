package services

import (
	"fmt"
	"sync"
	"time"

	"roottrack-backend/models"
	"roottrack-backend/repositories"
)

type NotificationWorker struct {
	UserRepo        repositories.UserRepository
	NotificationSvc *NotificationService
	WorkerLimit     int
}

func (w *NotificationWorker) Start() {
	ticker := time.NewTicker(10 * time.Minute)

	fmt.Println("Notification worker started...")

	for {
		<-ticker.C
		fmt.Println("Running notification job...")
		w.RunJob()
	}
}

func (w *NotificationWorker) RunJob() {
	users, err := w.UserRepo.GetAllUsers()
	if err != nil {
		fmt.Println("Failed to fetch users:", err)
		return
	}

	// Limit concurrency
	sem := make(chan struct{}, w.WorkerLimit)
	var wg sync.WaitGroup

	for _, user := range users {
		wg.Add(1)
		sem <- struct{}{}

		go func(u models.User) {
			defer wg.Done()
			defer func() { <-sem }()

			w.NotificationSvc.CheckAndNotify(&u)
		}(user)
	}

	wg.Wait()
	fmt.Println("Notification job completed")
}
