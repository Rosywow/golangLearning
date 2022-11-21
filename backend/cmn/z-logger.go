package cmn

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v4"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"time"

	"github.com/spf13/viper"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var isConsoleEnabled bool
var isPostgresqlEnabled bool
var isBBoltEnabled bool

//dstLogFile log file
var dstLogFile string

//L global zap.Logger
var z *zap.Logger

const zLoggerTimeLayout = "2006-01-02 15:04:05.000-0700"
const consoleLoggerTimeLayout = "15:04:05.000"

var rLogLevel = zapcore.InfoLevel
var stacktraceLogLevel = zapcore.ErrorLevel

/*
DebugLevel Level = iota - 1
DebugLevel: logs are typically voluminous, and are usually disabled
	in production.

InfoLevel: is the default logging priority.

WarnLevel: logs are more important than Info, but don't need
	individual human review.

ErrorLevel: logs are high-priority. If an application is running
	smoothly, it shouldn't generate any error-level logs.

DPanicLevel: logs are particularly significant errors. In
	development the logger panics after writing the message.

PanicLevel logs a message, then panics.

FatalLevel logs a message, then calls cmn.Terminate(1).

_minLevel = DebugLevel
_maxLevel = FatalLevel
*/

//SetLogLevel set runtime log level
func SetLogLevel(runtimeLevel int8) {
	rLogLevel = zapcore.Level(runtimeLevel)
}

func GetLogger() *zap.Logger {
	if z == nil {
		InitLogger()
	}
	return z
}

//InitLogger initialize zap logger
func InitLogger() {
	if z != nil {
		return
	}

	if viper.IsSet("zLogger.console.enable") {
		isConsoleEnabled = viper.GetBool("zLogger.console.enable")
	}

	if viper.IsSet("zLogger.postgresql.enable") {
		isPostgresqlEnabled = viper.GetBool("zLogger.postgresql.enable")
	}

	if viper.IsSet("zLogger.bblot.enable") {
		isBBoltEnabled = viper.GetBool("zLogger.bblot.enable")
	}

	consoleTimeLayout := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(consoleLoggerTimeLayout))
	}

	logLevel := zap.LevelEnablerFunc(func(v zapcore.Level) bool {
		return v >= rLogLevel
	})

	//------------------------
	//console logger output
	consoleLoggerCfg := zapcore.EncoderConfig{
		// Keys can be anything except the empty string.
		TimeKey:  "T",
		LevelKey: "L",
		NameKey:  "N",

		CallerKey:  "C",
		MessageKey: "M",

		StacktraceKey: "S",

		LineEnding:  zapcore.DefaultLineEnding,
		EncodeLevel: zapcore.CapitalColorLevelEncoder,
		EncodeTime:  consoleTimeLayout,

		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,

		ConsoleSeparator: " ",
	}

	if runtime.GOOS == "windows" {
		consoleLoggerCfg.EncodeLevel = zapcore.CapitalLevelEncoder
	}

	consoleEncoder := zapcore.NewConsoleEncoder(consoleLoggerCfg)

	consoleSink := zapcore.Lock(os.Stdout)

	//file logger
	//-------------------------

	basePath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		z.Error(err.Error())
		Terminate(1003)
	}

	dstLogPath := basePath + "/logs"
	if _, err := os.Stat(dstLogPath); os.IsNotExist(err) {
		_ = os.Mkdir(dstLogPath, os.ModePerm)
	}
	dstLogFile = dstLogPath + "/log.txt.json"
	fd, err := os.OpenFile(dstLogFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("open " + dstLogFile + " failed by " + err.Error())
	}

	zLoggerTimeLayout := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(zLoggerTimeLayout))
	}

	fileSink := zapcore.Lock(fd)

	jsonEncoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		TimeKey:        "T",
		LevelKey:       "L",
		NameKey:        "N",
		CallerKey:      "C",
		MessageKey:     "M",
		StacktraceKey:  "S",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zLoggerTimeLayout,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})

	//----------------------------
	//dbms logger
	dbLogger := DbLoggerAdaptor{
		stack: newLogItemStack(),
	}
	dbSink := zapcore.Lock(&dbLogger)

	var dst []zapcore.Core

	dst = append(dst, zapcore.NewCore(jsonEncoder, fileSink, logLevel))

	if isConsoleEnabled {
		dst = append(dst, zapcore.NewCore(consoleEncoder, consoleSink, logLevel))
	}

	if isBBoltEnabled || isPostgresqlEnabled {
		dst = append(dst, zapcore.NewCore(jsonEncoder, dbSink, logLevel))
	}

	//-------------------------
	//
	core := zapcore.NewTee(dst...)

	z = zap.New(core,
		zap.AddCaller(),
		zap.AddStacktrace(stacktraceLogLevel),
		zap.WrapCore(func(core zapcore.Core) zapcore.Core {
			return zapcore.NewSamplerWithOptions(core,
				time.Second, 100, 100)
		}),
	)

	defer func() {
		err := z.Sync()
		if err != nil {
			log.Println(err.Error())
		}
	}()

	if buildVer != "" {
		s := `
***************************************************************
******** app version: %s
***************************************************************`
		s = fmt.Sprintf(s, buildVer)
		z.Info(s)

		if viper.IsSet("zLogger.zLogLevel") {
			logLevel := zapcore.Level(viper.GetInt("zLogger.zLogLevel"))
			if logLevel > zapcore.FatalLevel || logLevel < zapcore.DebugLevel {
				z.Error("zLogLevel should between -1 and 5")
				Terminate(-1)
			}
			SetLogLevel(int8(logLevel))
		}
	}

}

//DbLoggerAdaptor adaptor for zap logger
type DbLoggerAdaptor struct {
	stack LogItemStack
}

//Sync flush data to persist
func (t *DbLoggerAdaptor) Sync() error {
	t.stack.sync(true)
	return nil
}

/*

DROP TABLE IF EXISTS t_log;

CREATE TABLE IF NOT EXISTS t_log
(
    grade character varying(8),
    msg character varying(4096),
    caller character varying(256),
    stacktrace character varying(4096),
    namespace character varying(64),
    "createTime" time(6) with time zone DEFAULT ('now'::text)::time with time zone,
    "loginUserName" character varying(128),
    "loginUserID" character varying(128),
		original jsonb
)
WITH (
    OIDS = TRUE
);
*/

//LogItem for logging unit
type LogItem struct {
	CreateTime string `json:"T,omitempty"`
	Level      string `json:"L,omitempty"`
	NameSpace  string `json:"N,omitempty"`
	Caller     string `json:"C,omitempty"`
	Message    string `json:"M,omitempty"`
	Stacktrace string `json:"S,omitempty"`
	Original   string `json:"O,omitempty"`
}

//Write log data to dbms
func (t *DbLoggerAdaptor) Write(d []byte) (n int, err error) {
	n = len(d)
	err = nil
	t.stack.store(d)
	t.stack.sync(false)
	//fmt.Printf("total stacked %d logItems", len(t.stack.buf))
	return
}

//Close sql.DB
func (t *DbLoggerAdaptor) Close() error {
	t.stack.sync(true)
	return nil
}

//LogItemStack pool the logging json data
type LogItemStack struct {
	recentSyncTime time.Time
	buf            [][]byte
	syncInterval   time.Duration

	idx  int
	data []interface{}
	err  error
}

//Next goto next row
func (v *LogItemStack) Next() bool {
	d := v.buf[v.idx]
	i := LogItem{}
	v.err = json.Unmarshal(d, &i)
	if v.err != nil {
		log.Print(v.err.Error())
		return false
	}

	z := reflect.ValueOf(i)
	v.data = make([]interface{}, z.NumField())
	for j := 0; j < z.NumField(); j++ {
		v.data[j] = z.Field(j).Interface()
	}

	v.idx++
	return v.idx < len(v.buf)
}

//Values return current row
func (v *LogItemStack) Values() ([]interface{}, error) {
	if v.err != nil {
		return nil, v.err
	}
	return v.data, nil
}

//Err return error while Next()
func (v *LogItemStack) Err() error {
	return v.err
}

func (v *LogItemStack) store(jsonData []byte) {
	if v.buf == nil {
		log.Print("LogItemStack.buf is nil")
		return
	}
	if jsonData == nil {
		log.Print("jsonData is nil")
		return
	}
	data := make([]byte, len(jsonData), len(jsonData))
	copy(data, jsonData)
	v.buf = append(v.buf, data)
}

func flashToPostgreSQL(v *LogItemStack) {
	if pgxConn == nil {
		return
	}
	var err error
	var rows [][]interface{}

	for k := 0; k < len(v.buf); k++ {
		d := v.buf[k]
		i := LogItem{}
		err = json.Unmarshal(d, &i)
		if err != nil {
			log.Print(err.Error())
			return
		}

		var t time.Time
		t, err = time.Parse(zLoggerTimeLayout, i.CreateTime)
		if err != nil {
			log.Print(err.Error())
			return
		}
		rows = append(rows, []interface{}{i.Level, i.Message, i.Caller, i.Stacktrace, t.UnixNano() / 1e6})
	}

	columns := []string{"grade", "msg", "caller", "stacktrace", "create_time"}
	_, err = pgxConn.CopyFrom(context.Background(), pgx.Identifier{"t_log"}, columns, pgx.CopyFromRows(rows))
	if err != nil {
		log.Print(err.Error())
		return
	}
}

func flashToBBolt(v *LogItemStack) {

}

func (v *LogItemStack) sync(force bool) {
	if v.buf == nil || len(v.buf) <= 0 {
		//log.Print("LogItemStack.buf is nil")
		return
	}

	if len(v.buf) <= 0 {
		//log.Print("LogItemStack.buf is empty")
		return
	}

	if !force && time.Now().Sub(v.recentSyncTime) < v.syncInterval {
		//fmt.Println("时间未到")
		return
	}

	//fmt.Println("准备开始")

	//begin sync
	if isPostgresqlEnabled {
		flashToPostgreSQL(v)
	}

	if isBBoltEnabled {
		flashToBBolt(v)
	}
	//end syncc

	//clear the stack
	for k := 0; k < len(v.buf); k++ {
		v.buf[k] = nil
	}
	v.buf = v.buf[:0]

	v.recentSyncTime = time.Now()
}

func newLogItemStack() LogItemStack {
	return LogItemStack{
		buf:            make([][]byte, 0, 20000),
		recentSyncTime: time.Now(),
		syncInterval:   time.Second * 15,
	}
}
