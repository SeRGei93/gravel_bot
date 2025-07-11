package await

const AdminDialogID = 6969539582 // ID администратора

var currentDialogUserID int64 = 0 // Текущий пользователь в диалоге

// NewDialog начинает новый диалог с пользователем (может быть вызвано только админом)
func NewDialog(adminID, userID int64) bool {
	if adminID != AdminDialogID {
		return false // Не админ — нельзя начать диалог
	}
	currentDialogUserID = userID
	return true
}

// GetDialog возвращает текущий диалог (админ может запросить)
func GetDialog(adminID int64) int64 {
	if adminID != AdminDialogID {
		return 0 // Не админ — нельзя получить диалог
	}
	return currentDialogUserID
}

// EndDialog завершает текущий диалог (может быть вызвано только админом)
func EndDialog(adminID int64) bool {
	if adminID != AdminDialogID {
		return false // Не админ — нельзя завершить диалог
	}
	currentDialogUserID = 0
	return true
}
