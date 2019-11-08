package common

import "github.com/sirupsen/logrus"

// Reader is the struct for the common reader library
type Reader struct {
	// Logger is the logrus logger pass to the reader
	Logger *logrus.Logger
}
