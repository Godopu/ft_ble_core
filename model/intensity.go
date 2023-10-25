package model

type Intensity struct {
	Alias string `json:"alias"`
	EMS   []int  `json:"ems"`
	Tens  []int  `json:"tens"`
}

var Intensities = []Intensity{
	{
		"Intensity 1",
		[]int{1, 2, 3, 4},
		[]int{1, 2, 3, 4},
	},
	{
		"Intensity 1",
		[]int{1, 2, 3, 4},
		[]int{1, 2, 3, 4},
	},
}

func RemoveItemFromIntensities(i int) {
	Intensities[i] = Intensities[len(Intensities)-1]
	Intensities = Intensities[:len(Intensities)-1]
}
