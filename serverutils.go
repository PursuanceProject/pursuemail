package main

import (
	"sync"

	log "github.com/Sirupsen/logrus"
	emailLib "github.com/jordan-wright/email"
)

func SendBulkEmail(emailAccounts []*EmailAccount, sendBulkEmailReq *SendBulkEmailRequest, emailPool *emailLib.Pool) (failedIds []string) {
	failedIdsChan := make(chan string)
	wg := new(sync.WaitGroup)
	for _, email := range emailAccounts {
		if sendBulkEmailReq.SecureOnly && !email.HasPubKey() {
			failedIds = append(failedIds, email.Id)
			continue
		}

		wg.Add(1)
		go func(email *EmailAccount) {
			defer wg.Done()
			err := email.Send(sendBulkEmailReq.EmailData, emailPool)
			if err != nil {
				log.Errorf("Error sending email: %v", err)
				failedIdsChan <- email.Id
			}
		}(email)
	}

	wg.Wait()
	close(failedIdsChan)
	for failedId := range failedIdsChan {
		failedIds = append(failedIds, failedId)
	}

	return failedIds
}
