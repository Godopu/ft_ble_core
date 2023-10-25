package model

type Trial struct {
	Alias          string `json:"alias"`
	EMG            []int  `json:"emg"`
	NumberOfTrials int    `json:"number_of_trials"`
	ViaCloud       bool   `json:"via_cloud"`
}

var Trials = make([]*Trial, 0, 100)

func RemoveItemFromTrials(i int) {
	Trials[i] = Trials[len(Trials)-1]
	Trials = Trials[:len(Trials)-1]
}
