package conn

import (
	"github.com/ipfs/go-log/v2"
	"go.uber.org/zap/zapcore"
	gormlogger "gorm.io/gorm/logger"
	"moul.io/zapgorm2"
)

func getLogger(lgr *log.ZapEventLogger) (gormlogger.Interface, gormlogger.LogLevel) {
	glogger := zapgorm2.New(lgr.Desugar())
	glogger.SetAsDefault()
	lvl := logLevel(lgr.Level())
	return glogger.LogMode(lvl), lvl
}

var levelMap = map[zapcore.Level]gormlogger.LogLevel{
	zapcore.ErrorLevel: gormlogger.Error,
	zapcore.WarnLevel:  gormlogger.Warn,
	zapcore.InfoLevel:  gormlogger.Info,
	zapcore.DebugLevel: gormlogger.Info,
}

func logLevel(lvl zapcore.Level) gormlogger.LogLevel {
	if l, ok := levelMap[lvl]; ok {
		return l
	}
	return gormlogger.Silent
}

/* Copyright 2023 Spiegel
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * 	http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
