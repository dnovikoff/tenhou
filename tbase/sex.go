package tbase

type Sex int

const (
	SexUnknown Sex = iota
	SexMale
	SexFemale
	SexComputer
)

func (this Sex) Letter() string {
	switch this {
	case SexMale:
		return "M"
	case SexFemale:
		return "F"
	case SexComputer:
		return "C"
	}
	return "?"
}

func ParseSexLetter(in string) Sex {
	switch in {
	case "M":
		return SexMale
	case "F":
		return SexFemale
	case "C":
		return SexComputer
	}
	return SexUnknown
}
