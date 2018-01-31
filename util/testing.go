package util

import (
	"archive/zip"
	"bufio"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/stretchr/testify/require"
)

type TestingT interface {
	require.TestingT
	Log(args ...interface{})
	Logf(format string, args ...interface{})
}

func processZip(t TestingT, zipFileName string, prefix string, fn func(line string)) {
	reader, err := zip.OpenReader(zipFileName)
	require.NoError(t, err)
	defer reader.Close()

	for _, f := range reader.File {
		if f.FileInfo().IsDir() || !strings.HasSuffix(f.Name, ".log") {
			continue
		}
		t.Log("Processing " + zipFileName + ": " + f.Name)
		d, err := f.Open()
		require.NoError(t, err)
		processReadCloser(t, d, prefix, fn)
	}
}

func ProcessExampleLogs(t TestingT, prefix string, f func(line string)) {
	dir := "../network/examples/"
	infos, err := ioutil.ReadDir(dir)
	require.NoError(t, err)
	tested := false
	for _, info := range infos {
		name := info.Name()
		if strings.HasSuffix(name, ".log") {
			t.Log("Processing " + name)
			file, err := os.Open(dir + name)
			require.NoError(t, err)
			n := processReadCloser(t, file, prefix, f)
			if n > 0 {
				tested = true
			}
		} else if strings.HasSuffix(name, ".zip") {
			processZip(t, dir+name, prefix, f)
		}
	}
	require.True(t, tested)
}

func processReadCloser(t require.TestingT, r io.ReadCloser, prefix string, fn func(line string)) int {
	defer r.Close()
	lines := readExampleLines(t, r, prefix)
	for _, line := range lines {
		fn(line)
	}
	return len(lines)
}

func FixLine(in string) string {
	return strings.Replace(in, " />", "/>", -1)
}

func CompareLines(t require.TestingT, expected, actual string) {
	expected, actual = FixLine(expected), FixLine(actual)
	require.Equal(t, expected, actual, actual)
}

func readExampleLines(t require.TestingT, r io.Reader, prefix string) (ret []string) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, prefix) {
			continue
		}
		message := line[len(prefix):]
		message = strings.TrimSpace(message)
		ret = append(ret, message)
	}

	require.NoError(t, scanner.Err())
	return
}
