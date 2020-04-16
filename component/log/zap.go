package log

import (
	"github.com/ebar-go/ego/utils"
	"github.com/ebar-go/ego/utils/date"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"log"
	"os"
	"time"
)

// NewZap return zap logger instance
func NewZap(filename string, enableLevel zapcore.LevelEnabler, initFields ...zap.Field) (*zap.Logger) {
	conf := zapcore.EncoderConfig{
		MessageKey:  "message",
		LevelKey:    "level_name",
		EncodeLevel: zapcore.LowercaseLevelEncoder, //将级别转换成大写
		TimeKey:     "datetime",

		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format(date.TimeFormat))
		},
		CallerKey:    "file",
		EncodeCaller: zapcore.ShortCallerEncoder,
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		},
	}

	encoder := zapcore.NewJSONEncoder(conf)
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(getRotateWriter(filename)), enableLevel),
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
	prefixName, ext := utils.SplitPathExt(filename)
	// demo.log是指向最新日志的链接
	hook, err := rotatelogs.New(
		prefixName+"-%Y%m%d%H%M" + ext, // 没有使用go风格反人类的format格式
		//rotatelogs.WithLinkName(filename),
		rotatelogs.WithMaxAge(time.Hour*24*30),      // 保存30天
		rotatelogs.WithRotationTime(time.Hour * 24), //切割频率 24小时
	)

	if err != nil {
		log.Printf("getWriter:%v", err)
		panic(err)
	}

	return hook
}
