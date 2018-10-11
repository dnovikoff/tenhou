package stats

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseMainPage(t *testing.T) {
	assert.Equal(t, []string{
		"scraw2006.zip",
		"scraw2007.zip",
		"scraw2008.zip",
		"scraw2009.zip",
		"scraw2010.zip",
		"scraw2011.zip",
		"scraw2012.zip",
		"scraw2013.zip",
		"scraw2014.zip",
		"scraw2015.zip",
		"scraw2016.zip",
		"scraw2017.zip",
	}, ParseMain(loadTestData(t, "raw.html")))
}

func TestParseList(t *testing.T) {
	assert.Equal(t, []ListItem{
		{File: "sca20180926.log.gz", Size: 35919},
		{File: "sca20180927.log.gz", Size: 33557},
		{File: "sca20180928.log.gz", Size: 35839},
	}, MustParseList(loadTestData(t, "list.js")))
}

func loadTestData(t require.TestingT, filename string) string {
	data, err := ioutil.ReadFile("test_data/" + filename)
	require.NoError(t, err)
	return string(data)
}
