package await

import (
	"time"
)

type AwaitType string
type AwaitStep string

const (
	AwaitGift   AwaitType = "gift"
	AwaitResult AwaitType = "result"

	StepNone AwaitStep = ""
	Wait     AwaitStep = "wait"
	Received AwaitStep = "received"
)

type Await struct {
	Type AwaitType
	Step AwaitStep
	Time time.Time
}

var (
	awaitingMessage = make(map[int64]Await)
)

// SetAwaiting отмечает пользователя как ожидающего ввода и запоминает текущее время
func SetAwaiting(userID int64, seconds int, awaitType AwaitType) {
	awaitingMessage[userID] = Await{
		Type: awaitType,
		Step: Wait,
		Time: time.Now().Add(time.Duration(seconds) * time.Second),
	}
}

// GetAwaiting проверяет, ожидается ли сообщение от пользователя и не истек ли таймаут
func GetAwaiting(userID int64) (*Await, bool) {
	userAwait, exists := awaitingMessage[userID]
	if !exists {
		return nil, false
	}

	if time.Now().After(userAwait.Time) {
		delete(awaitingMessage, userID)
		return nil, false
	}

	return &userAwait, true
}

// CleanupOldAwaiting очищает все записи старше времени ожидания
func CleanupOldAwaiting() {
	for id, userAwait := range awaitingMessage {
		if time.Now().After(userAwait.Time) {
			delete(awaitingMessage, id)
		}
	}
}

// DeleteAwaiting удаляет запись вручную
func DeleteAwaiting(userID int64) {
	delete(awaitingMessage, userID)
}
