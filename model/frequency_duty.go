package model

type FrequencyDutyRatio struct {
	Frequency int `json:"frequency"`
	DutyRatio int `json:"DutyRatio"`
}

var FreqDuty *FrequencyDutyRatio = nil
