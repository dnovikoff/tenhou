package stats

import (
	"bufio"
	"fmt"
	"html"
	"io"
	"math"
	"path"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/facebookgo/stackerr"
	"github.com/golang/protobuf/ptypes"
	"go.uber.org/multierr"

	"github.com/dnovikoff/tenhou/genproto/stats"
)

type StatFileInfo struct {
	Source string
	Date   time.Time
}

func ParseStatFileName(filename string) (*StatFileInfo, error) {
	p := path.Base(filename)
	i := strings.Index(p, ".")
	if i != -1 {
		p = p[:i]
	}
	if len(p) < 3 {
		return nil, fmt.Errorf("Incorrect filename '%v'", filename)
	}
	source := p[:3]
	p = p[3:]
	const layout = "20060102"
	if len(p) > len(layout) {
		p = p[:len(layout)]
	}
	time, err := time.Parse("20060102", p)
	if err != nil {
		return nil, err
	}
	return &StatFileInfo{source, time}, nil
}

func ParseStatsForFile(reader io.Reader, filename string, f func(*stats.Record) error) error {
	info, err := ParseStatFileName(filename)
	if err != nil {
		return err
	}
	return ParseStats(reader, info, f)
}

func ParseStats(reader io.Reader, info *StatFileInfo, f func(*stats.Record) error) error {
	scanner := bufio.NewScanner(reader)
	parser := &statParser{time: info.Date}
	lparser := bindParser(parser, info.Source)
	if lparser == nil {
		return fmt.Errorf("unknown parser source: %v", info.Source)
	}
	for scanner.Scan() {
		str := scanner.Text()
		record := parser.Next()
		err := lparser.ParseLine(str)
		if err != nil {
			fmt.Printf("error parsing line '%v'\n", str)
			return err
		} else {
			err = f(record)
			if err != nil {
				fmt.Printf("error parsing line '%v'\n", str)
				return err
			}
		}
	}
	return scanner.Err()
}

const (
	dashSymbol = '－'
	timeFormat = "15:04"
)

type statParser struct {
	time   time.Time
	record *stats.Record
}

func (p *statParser) Next() *stats.Record {
	p.record = &stats.Record{}
	return p.record
}

func (p *statParser) Time(in string) error {
	t, err := time.Parse(timeFormat, in)
	if err != nil {
		return err
	}
	p.record.Time, err = ptypes.TimestampProto(
		p.time.Add(time.Nanosecond * time.Duration(t.UnixNano())),
	)
	return err
}

func (p *statParser) Duration(in string) error {
	dur, err := strconv.Atoi(in)
	if err != nil {
		return err
	}
	p.record.Duration = ptypes.DurationProto(time.Minute * time.Duration(dur))
	return nil
}

func (p *statParser) Info(in string) error {
	start := strings.Index(in, "\"")
	if start == -1 {
		return fmt.Errorf("no open quote found in '%v'", in)
	}
	end := strings.LastIndex(in, "\"")
	if end <= start {
		return fmt.Errorf("wrong quote in '%v'", in)
	}
	start = strings.LastIndex(in, "=")
	if start == -1 {
		return fmt.Errorf("no = found '%v'", in)
	}
	p.record.Id = in[start+1 : end]
	return nil
}

func (p *statParser) Lobby(in string) error {
	letter := in[0]
	number := in[1:]
	switch letter {
	case 'L':
	case 'C':
		p.record.IsChampionLobby = true
	default:
		return stackerr.Newf("Unknown first lobby letter '%v'", letter)
	}
	lobby, err := strconv.Atoi(number)
	if err != nil {
		return err
	}
	p.record.Number = int64(lobby)
	return nil
}

func sanitize(str string) string {
	if utf8.ValidString(str) {
		return str
	}
	v := make([]rune, 0, len(str))
	for i, r := range str {
		if r == utf8.RuneError {
			_, size := utf8.DecodeRuneInString(str[i:])
			if size == 1 {
				continue
			}
		}
		v = append(v, r)
	}
	return "_" + string(v)
}

func parsePlayer(in string, line *stats.Player) error {
	start := strings.LastIndex(in, "(")
	if start == -1 {
		return fmt.Errorf("No open brace found in '%v'", in)
	}
	end := strings.LastIndex(in, ")")
	if end < start {
		return fmt.Errorf("Wrong braces in '%v'", in)
	}
	var f float64
	mString := in[start+1 : end]
	var bonus int64
	_, err := fmt.Sscanf(mString, "%g,%d枚", &f, &bonus)
	if err != nil {
		f, err = strconv.ParseFloat(mString, 32)
		if err != nil {
			return err
		}
	} else {
		line.Coins = bonus
	}
	line.Score = int64(math.Floor(f*10 + .5))
	line.Name = sanitize(in[:start])
	return nil
}

func (p *statParser) Players(in string) error {
	brIndex := strings.Index(in, "<br>")
	if brIndex != -1 {
		in = in[:brIndex]
		in = html.UnescapeString(in)
	}
	in = strings.TrimSpace(in)
	s := strings.Split(in, " ")
	players := make([]*stats.Player, len(s))
	for k, v := range s {
		var p stats.Player
		err := parsePlayer(v, &p)
		if err != nil {
			return err
		}
		players[k] = &p
	}
	p.record.Players = players
	return nil
}

func (p *statParser) gameType(r rune) error {
	switch r {
	case '三':
		p.record.Type = stats.GameType_GAME_TYPE_3
	case '四':
		p.record.Type = stats.GameType_GAME_TYPE_4
	default:
		return fmt.Errorf("Unknown Game Type '%v'", string(r))
	}
	return nil
}

func (p *statParser) lobbyType(r rune) error {
	switch r {
	case '鳳':
		p.record.Lobby = stats.LobbyType_LOBBY_TYPE_PHOENIX
	case '特':
		p.record.Lobby = stats.LobbyType_LOBBY_TYPE_UPPERDAN
	case '上':
		p.record.Lobby = stats.LobbyType_LOBBY_TYPE_DAN
	case '般':
		p.record.Lobby = stats.LobbyType_LOBBY_TYPE_KU
	case '技':
		p.record.Lobby = stats.LobbyType_LOBBY_TYPE_DZ
	case '若':
		p.record.Lobby = stats.LobbyType_LOBBY_TYPE_X1
	case '銀':
		p.record.Lobby = stats.LobbyType_LOBBY_TYPE_X2
	case '琥':
		p.record.Lobby = stats.LobbyType_LOBBY_TYPE_X3
	case '－':
		p.record.Lobby = stats.LobbyType_LOBBY_TYPE_EXTERNAL
	default:
		return fmt.Errorf("unknown Lobby Type '%v'", string(r))
	}
	return nil
}

func (p *statParser) lenType(r rune) error {
	switch r {
	case '東':
		p.record.Length = stats.GameLength_GAME_LENGTH_EAST
	case '南':
		p.record.Length = stats.GameLength_GAME_LENGTH_SOUTH
	case dashSymbol:
		p.record.Length = stats.GameLength_GAME_LENGTH_ONE
	default:
		return fmt.Errorf("unknown Game Len '%v'", string(r))
	}
	return nil
}

func (p *statParser) tanyaoType(r rune) error {
	switch r {
	case '喰':
		p.record.Tanyao = stats.Tanyao_TANYAO_YES
	case dashSymbol:
		p.record.Tanyao = stats.Tanyao_TANYAO_NO
	default:
		return fmt.Errorf("unknown Tanyao type '%v'", string(r))
	}
	return nil
}

func (p *statParser) akkaType(r rune) error {
	switch r {
	case '赤':
		p.record.Akkas = stats.Akkas_AKKAS_YES
	case dashSymbol:
		p.record.Akkas = stats.Akkas_AKKAS_NO
	default:
		return fmt.Errorf("unknown Akka Len '%v'", string(r))
	}
	return nil
}

func (p *statParser) dzType(r rune) error {
	switch r {
	case '祝':
		p.record.IsDz = true
	case dashSymbol:
		p.record.IsDz = false
	default:
		return fmt.Errorf("unknown Dz Type '%v'", string(r))
	}
	return nil
}

func (p *statParser) five(r rune) error {
	switch r {
	case '５':
		p.record.NumberType = stats.NumberType_NUMBER_5
	case '２':
		p.record.NumberType = stats.NumberType_NUMBER_2
	case '０':
		p.record.NumberType = stats.NumberType_NUMBER_0
	default:
		return fmt.Errorf("unknown Number Type '%v'", string(r))
	}
	return nil
}

func (p *statParser) LobbyConfig(in string) error {
	runes := []rune(in)
	var allErrors error
	regErr := func(err error) {
		if err == nil {
			return
		}
		allErrors = multierr.Append(allErrors, err)
	}
	if len(runes) == 7 {
		regErr(p.five(runes[6]))
		runes = runes[:6]
	}
	switch len(runes) {
	case 6:
		regErr(p.dzType(runes[5]))
	case 5:
	default:
		return fmt.Errorf("wrong config string len %v for '%v'", len(runes), in)
	}

	regErr(p.gameType(runes[0]))
	regErr(p.lobbyType(runes[1]))
	regErr(p.lenType(runes[2]))
	regErr(p.tanyaoType(runes[3]))
	regErr(p.akkaType(runes[4]))
	return allErrors
}

type lineParsers []fieldParser

func (p lineParsers) ParseLine(in string) error {
	s := strings.SplitN(in, "|", len(p))
	for k, v := range s {
		s[k] = strings.TrimSpace(v)
	}
	return p.ParseArray(s)
}

func (p lineParsers) ParseArray(in []string) error {
	cnt := len(p)
	if len(in) != cnt {
		return stackerr.Newf("wrong number of args %v != %v : %v", len(in), cnt, in)
	}
	var retErr error
	for i, v := range p {
		err := v(in[i])
		if err != nil {
			retErr = multierr.Append(retErr, err)
		}
	}
	return retErr
}

type fieldParser func(string) error

func bindParser(p *statParser, prefix string) lineParsers {
	switch prefix {
	case "sca":
		return lineParsers{p.Lobby, p.Time, p.LobbyConfig, p.Players}
	case "scb":
		return lineParsers{p.Time, p.Duration, p.LobbyConfig, p.Players}
	case "scc":
		return lineParsers{p.Time, p.Duration, p.LobbyConfig, p.Info, p.Players}
	case "scd":
		return lineParsers{p.Time, p.Duration, p.LobbyConfig, p.Players}
	case "sce":
		return lineParsers{p.Time, p.Duration, p.LobbyConfig, p.Info, p.Players}
	case "scf":
		return lineParsers{p.Lobby, p.Time, p.LobbyConfig, p.Info, p.Players}
	}
	return nil
}
