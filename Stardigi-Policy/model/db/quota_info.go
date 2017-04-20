package db

type QuotaInfo struct {
	MarathonName string
	AppId        string
	RuleType     string
	MaxThreshold float64
	MinThreshold float64
}
