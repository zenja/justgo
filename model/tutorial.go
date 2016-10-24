package model

type Tutorial struct {
	Key            string
	Title          string
	Description    string
	Code           string
	ExpectedStdout string
}

type ExtendedTutorial struct {
	Tutorial
	PrevKey string
	NextKey string
}
