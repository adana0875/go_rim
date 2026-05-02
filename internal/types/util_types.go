package types

type LoadedResult int

const (
	NoneLoaded LoadedResult = iota
	LoadSuccess
	LoadFailure
)
