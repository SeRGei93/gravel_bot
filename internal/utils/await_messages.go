package utils

import (
	"sync"
	"time"
)

var (
	awaitingMessage     = make(map[int64]time.Time)
	awaitingMessageLock sync.Mutex
)

// SetAwaiting отмечает пользователя как ожидающего ввода и запоминает текущее время
func SetAwaiting(userID int64, seconds int) {
	awaitingMessageLock.Lock()
	defer awaitingMessageLock.Unlock()
	awaitingMessage[userID] = time.Now().Add(time.Duration(seconds) * time.Second)
}

// IsAwaiting проверяет, ожидается ли сообщение от пользователя и не истек ли таймаут
func IsAwaiting(userID int64) bool {
	awaitingMessageLock.Lock()
	defer awaitingMessageLock.Unlock()

	end, exists := awaitingMessage[userID]
	if !exists {
		return false
	}

	if time.Now().After(end) {
		delete(awaitingMessage, userID)
		return false
	}

	return true
}

// CleanupOldAwaiting очищает все записи старше времени ожидания
func CleanupOldAwaiting() {
	awaitingMessageLock.Lock()
	defer awaitingMessageLock.Unlock()
	for id, t := range awaitingMessage {
		if time.Now().After(t) {
			delete(awaitingMessage, id)
		}
	}
}

// DeleteAwaiting удаляет запись вручную
func DeleteAwaiting(userID int64) {
	awaitingMessageLock.Lock()
	defer awaitingMessageLock.Unlock()
	delete(awaitingMessage, userID)
}
