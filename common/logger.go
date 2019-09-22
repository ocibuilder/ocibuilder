/*
Copyright 2019 BlackRock, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package common

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Create a new instance of the logger. You can have any number of instances.
var log = GetLogger(false)

// GetLogger returns a new logger
func GetLogger(debug bool) *logrus.Logger {
	log := &logrus.Logger{
		Out:   os.Stdout,
		Level: logrus.InfoLevel,
		Formatter: &logrus.TextFormatter{
			TimestampFormat:  "2006-01-02 15:04:05",
			FullTimestamp:    true,
			ForceColors:      true,
			QuoteEmptyFields: true,
		},
	}
	if debug {
		log.SetLevel(logrus.DebugLevel)
	}
	return log
}
