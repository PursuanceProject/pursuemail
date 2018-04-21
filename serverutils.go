package main

import (
	"sync"

	log "github.com/Sirupsen/logrus"
	emailLib "github.com/jordan-wright/email"
)

func SendBulkEmail(emailAccounts []*EmailAccount, sendBulkEmailReq *SendBulkEmailRequest, emailPool *emailLib.Pool) (failedIds []string) {
	failedIdsChan := make(chan string)
	var total int
	wg := new(sync.WaitGroup)
	for i, email := range emailAccounts {
		if sendBulkEmailReq.SecureOnly && !email.HasPubKey() {
			failedIds = append(failedIds, email.Id)
			continue
		}

		total++
		wg.Add(1)
		go func(i int, email *EmailAccount) {
			var maybeFailedId string
			log.Debugf("Sending bulk email #%v", i+1)
			err := email.Send(sendBulkEmailReq.EmailData, emailPool)
			if err != nil {
				log.Errorf("Error sending (instance of bulk) email: %v", err)
				maybeFailedId = email.Id
			}
			wg.Done()
			failedIdsChan <- maybeFailedId
		}(i, email)
	}

	log.Debugf("SendBulkEmailRequest: waiting for %v email(s) to send", total)
	wg.Wait()

	var failedId string
	for i := 0; i < total; i++ {
		failedId = <-failedIdsChan
		if failedId != "" {
			failedIds = append(failedIds, failedId)
		}
	}

	return failedIds
}
