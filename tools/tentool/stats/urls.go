package stats

const (
	MainPageURL = "http://tenhou.net/sc/raw/"
	ListURL     = "http://tenhou.net/sc/raw/list.cgi"
	ListOldURL  = "http://tenhou.net/sc/raw/list.cgi?old"
)

func MakeFullURL(short string) string {
	return MainPageURL + short
}

func MakeFullURLs(short []string) []string {
	if short == nil {
		return nil
	}
	x := make([]string, len(short))
	for k, v := range short {
		x[k] = MakeFullURL(v)
	}
	return x
}
