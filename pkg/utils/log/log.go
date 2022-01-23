/*
 * @Date: 2022.01.22 17:36
 * @Description: Omit
 * @LastEditors: Rustle Karl
 * @LastEditTime: 2022.01.22 17:36
 */

package log

import (
	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func init() {
	Log = logrus.New()
	Log.Out = colorable.NewColorableStdout()
	Log.Level = logrus.DebugLevel
	Log.Formatter = defaultFormatter
}

func Errorf(format string, v ...interface{}) {
	Log.Errorf(format, v...)
}

func Warnf(format string, v ...interface{}) {
	Log.Warnf(format, v...)
}

func Infof(format string, v ...interface{}) {
	Log.Infof(format, v...)
}

func Debugf(format string, v ...interface{}) {
	Log.Debugf(format, v...)
}
