package flightaware

import "math"

const windowSize = 2
const maxFault = 0.2

func filterByDuration(r []flightHistoryRaw) []flightHistoryRaw {
	res := make([]flightHistoryRaw, 0)
	var sum uint64
	var avg uint64
	var count uint64 = 0
	length := len(r)
	for i := 0; i < length-windowSize+1; i++ {
		count = 0
		sum = 0

		for j := -windowSize; j < windowSize; j++ {
			if (i+j) < 0 || (i+j) > length {
				continue
			}
			sum += uint64(r[i+j].Duration)
			count++
		}
		avg = sum / count
		if math.Abs(1-(float64(r[i].Duration)/float64(avg))) < maxFault {
			res = append(res, r[i])
		}
	}
	return res
}
