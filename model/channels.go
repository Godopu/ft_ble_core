package model

var trialChannel chan *Trial

func init() {
	trialChannel = make(chan *Trial)
}

func RequestTrial(trial *Trial) {
	trialChannel <- trial
}

func RetrieveTrialRequest() *Trial {
	trial, ok := <-trialChannel
	if !ok {
		return nil
	}

	return trial
}
