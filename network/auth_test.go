package network

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuth(t *testing.T) {
	tst := func(in string) string {
		ret, err := AuthAnswer(in)
		if err != nil {
			return err.Error()
		}
		return ret
	}
	assert.Equal(t, "20180117-c1ebb26f", tst("20180117-e7b5e83e"))
	assert.Equal(t, "Wrong size", tst("2018011-e7b5e83e"))
	assert.Equal(t, "Wrong size", tst("20180117-e7b5e83"))
	assert.Equal(t, "Wrong parts", tst("20180117e7b5e83e"))
	assert.Equal(t, "strconv.ParseInt: parsing \"2180a17\": invalid syntax", tst("2a180a17-e7b5e83e"))
}
