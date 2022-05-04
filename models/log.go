package models

// Ip represents data about a log record.
type Log struct {
	IP     string  `json:"IP"`   // json 回傳時IP 的 Key 對應到 IP
	Time   string  `json:"Time"` // json 回傳時Time 的 Key 對應到的 Time 要跟資料庫欄位相同
	URL    string  `json:"Url"`
	Status float64 `json:"Status"`
}
