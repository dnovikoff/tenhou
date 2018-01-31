package tbase

type DrawType int

const (
	DrawUnknown DrawType = iota
	DrawEnd
	Draw9
	DrawReach4
	DrawRon3
	DrawKan4
	DrawWind4
	DrawNagashi
)

var DrawMap = map[string]DrawType{
	"":       DrawEnd,
	"yao9":   Draw9,
	"reach4": DrawReach4,
	"ron3":   DrawRon3,
	"kan4":   DrawKan4,
	"kaze4":  DrawWind4,
	"nm":     DrawNagashi,
}

var ReverseDrawMap = func() map[DrawType]string {
	ret := make(map[DrawType]string, len(DrawMap))
	for k, v := range DrawMap {
		ret[v] = k
	}
	return ret
}()
