package common

func Percent(hr, mn, mx int) float64 {
	if mx == mn || hr <= mn {
		return 0
	} else if hr >= mx {
		return 1
	}

	return float64(hr-mn) / float64(mx-mn)
}

func PercentFull(hr, mn, mx int) float64 {
	return 2*Percent(hr, mn, mx) - 1
}
