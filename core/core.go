package core

type Bundle struct {
	Status   int
	Code     int
	Plural   map[string]string
	Singular map[string]string
	Msg      map[string]string
}
