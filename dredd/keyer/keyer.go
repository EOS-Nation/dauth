package keyer

import (
	"time"
)

const BLACK_LIST_USER_KEY_PREFIX = "BLU"
const DELIMITER = "_"

//DCD:[USER_ID]:[DATE]
func DocumentConsumptionDaily(userID string, date time.Time) string {
	return "DCD" + DELIMITER + userID + DELIMITER + date.Format("20060102")
}

//DCM:[USER_ID]:[DATETIME]
func DocumentConsumptionMinutely(userID string, date time.Time) string {
	return "DCM" + DELIMITER + userID + DELIMITER + date.Format("20060102-1504")
}

//DCD:[USER_ID]:[DATE]
func DocumentConsumptionLast30Days(userID string, date time.Time) []string {
	var out []string
	for i := 0; i < 30; i++ {
		d := date.Add((-24 * time.Hour) * time.Duration(i))
		out = append(out, DocumentConsumptionDaily(userID, d))
	}
	return out
}

//DCD:[USER_ID]:[DATETIME]
func DocumentConsumptionLastWindow(userID string, windowSize int, date time.Time) []string {
	var out []string
	for i := 0; i < windowSize; i++ {
		d := date.Add((-1 * time.Minute) * time.Duration(i))
		out = append(out, DocumentConsumptionMinutely(userID, d))
	}
	return out
}

//DCCP:[USER_ID]
func CurrentPeriodDocumentConsumption(userID string) string {
	return "DCCP" + DELIMITER + userID
}

////IPC:[USER_ID]:[API_KEY]:[DATE]
//func IPAddressCountByApiKeyKey(userID string, apiKey string, date time.Time) string {
//	return "IPC:" + userID + DELIMITER + apiKey + DELIMITER + date.Format("20060102")
//}

//BLU:[USER_ID]
func UserIDBlackListKey(userID string) string {
	return BLACK_LIST_USER_KEY_PREFIX + DELIMITER + userID
}

//BLU:VERSION
func UserBlackListVersionKey() string {
	return "BLU" + DELIMITER + "VERSION"
}

////BLAK_:[USER_ID]:[API_KEY]
//func APIKeyBlackListKey(userID string, apiKey string) string {
//	return "BLAK:" + userID + DELIMITER + apiKey
//}

////MSTAT:[USER_ID]
//func MonthlyStatsKey(userID string) string {
//	return "MSTAT:" + userID
//}

//QT:[USER_ID]
func UserQuotaCachePrefix(userID string) string {
	return "QT" + DELIMITER + userID
}

//QT:[USER_ID]:[API_KEY_ID]
func UserQuotaCacheKey(userID string, apiKeyID string) string {
	return UserQuotaCachePrefix(userID) + DELIMITER + apiKeyID
}
