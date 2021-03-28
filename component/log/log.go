package log

import (
	"github.com/ebar-go/ego/component/trace"
	"github.com/ebar-go/egu"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"log"
	"os"
	"time"
)

type Context map[string]interface{}

// format
func format(ctx Context) zap.Field {
	return zap.Any("context", ctx)
}

func traceId() zap.Field {
	return zap.String("trace_id", trace.Get())
}

// Info
func (logger *Logger) Info(message string, ctx Context) {
	logger.instance.Info(message, format(ctx), traceId())
}

// Debug
func (logger *Logger) Debug(message string, ctx Context) {
	logger.instance.Debug(message, format(ctx), traceId())
}

// Error
func (logger *Logger) Error(message string, ctx Context) {
	logger.instance.Error(message, format(ctx), traceId())
}

// Logger
type Logger struct {
	path     string
	debug    bool
	fields   map[string]interface{}
	instance *zap.Logger
}



// getInstance init logger instance
func New(debug bool, fields map[string]interface{}) *Logger {
	logger := new(Logger)
	logger.debug = debug
	logger.fields = fields
	logger.lazyInit()

	return logger

}

func (logger *Logger) lazyInit() {
	level := zap.InfoLevel
	if logger.debug {
		level = zap.DebugLevel
	}

	var fields []zap.Field
	for idx, val := range logger.fields {
		fields = append(fields, zap.Any(idx, val))
	}
	logger.instance = newZap(logger.path, level, fields...)
}


// newZap return zap logger instance
func newZap(filename string, enableLevel zapcore.LevelEnabler, initFields ...zap.Field) *zap.Logger {
	conf := zapcore.EncoderConfig{
		MessageKey:  "message",
		LevelKey:    "level",
		EncodeLevel: zapcore.LowercaseLevelEncoder, //将级别转换成大写
		TimeKey:     "time",

		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format(egu.TimeFormat))
		},
		CallerKey:    "file",
		EncodeCaller: zapcore.ShortCallerEncoder,
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		},
	}

	encoder := zapcore.NewJSONEncoder(conf)
	core := zapcore.NewTee(
		//zapcore.NewCore(encoder, zapcore.AddSync(getRotateWriter(filename)), enableLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), enableLevel),
	)

	// 需要传入 zap.AddCaller() 才会显示打日志点的文件名和行数
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(enableLevel), zap.Fields(initFields...))

	return logger
}

// getRotateWriter
func getRotateWriter(filename string) io.Writer {
	// 生成rotate logs的Logger 实际生成的文件名 demo.log.YYmmddHH
	// 分割文件
	//prefixName, ext := egu.SplitPathExt(filename)
	// demo.log是指向最新日志的链接
	hook, err := rotatelogs.New(
		filename, // 没有使用go风格反人类的format格式
		//rotatelogs.WithLinkName(filename),
		rotatelogs.WithMaxAge(time.Hour*24*30),    // 保存30天
		rotatelogs.WithRotationTime(time.Hour*24), //切割频率 24小时
	)

	if err != nil {
		log.Printf("getWriter:%v", err)
		panic(err)
	}

	return hook
}
