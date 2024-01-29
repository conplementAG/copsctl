package common

import "github.com/sirupsen/logrus"

func FatalOnError(err error) {
	if err != nil {
		logrus.Fatal(err)
	}
}

func ToPtr[T any](source T) *T {
	return &source
}
