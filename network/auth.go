package network

import (
	"fmt"
	"strconv"
	"strings"
)

var translationTable = []int{63006, 9570, 49216, 45888, 9822, 23121, 59830, 51114, 54831, 4189, 580, 5203, 42174, 59972,
	55457, 59009, 59347, 64456, 8673, 52710, 49975, 2006, 62677, 3463, 17754, 5357}

func AuthAnswer(input string) (ret string, err error) {
	parts := strings.Split(input, "-")
	if len(parts) != 2 {
		err = fmt.Errorf("Wrong parts")
		return
	}
	first := parts[0]
	second := parts[1]

	if len(first) != 8 || len(second) != 8 {
		err = fmt.Errorf("Wrong size")
		return
	}

	i1, err := strconv.ParseInt("2"+first[2:8], 10, 64)
	if err != nil {
		return
	}
	i2, err := strconv.ParseInt("2"+first[7:8], 10, 64)
	if err != nil {
		return
	}
	index := i1 % (12 - i2) * 2

	a16, err := strconv.ParseUint(second[0:4], 16, 64)
	if err != nil {
		return
	}
	b16, err := strconv.ParseUint(second[4:8], 16, 64)
	if err != nil {
		return
	}

	a := translationTable[index] ^ int(a16)
	b := translationTable[index+1] ^ int(b16)

	ret = first + "-" + fmt.Sprintf("%2x%2x", a, b)
	return
}
