package gwlog

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
	"time"
)

// https://studygolang.com/articles/25044?fr=sidebar
// https://blog.csdn.net/wdy_yx/article/details/79479484
var (
	encoder  zapcore.Encoder
	infoHook io.Writer
	warnHook io.Writer
	core     zapcore.Core
	configex zapcore.EncoderConfig
)

//var sugar *zap.SugaredLogger

func init() {
	// 设置一些基本日志格式 具体含义还比较好理解，直接看zap源码也不难懂
	configex = zapcore.EncoderConfig{
		MessageKey:  "msg",
		LevelKey:    "level",
		EncodeLevel: zapcore.CapitalLevelEncoder,
		TimeKey:     "ts",
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
		CallerKey:    "file",
		EncodeCaller: zapcore.ShortCallerEncoder,
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		},
	}

	currentLevel = DebugLevel

	// 输出默认到控制台
	core = zapcore.NewTee(zapcore.NewCore(zapcore.NewConsoleEncoder(configex),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)), currentLevel))

	rebuildLoggerFromCfgEx()
}

func getWriter(filename string) io.Writer {
	// 生成rotatelogs的Logger 实际生成的文件名 demo.log.YYmmddHH
	// demo.log是指向最新日志的链接
	// 保存7天内的日志，每1小时(整点)分割一次日志
	hook, err := rotatelogs.New(
		filename+".%Y%m%d%H", // 没有使用go风格反人类的format格式
		rotatelogs.WithLinkName(filename),
		rotatelogs.WithMaxAge(time.Hour*24*7),
		rotatelogs.WithRotationTime(time.Hour),
	)

	if err != nil {
		panic(err)
	}
	return hook
}

func rebuildLoggerFromCfgEx() {
	// 最后创建具体的Logger
	logger := zap.New(core /*zap.AddCaller()*/) // 需要传入 zap.AddCaller() 才会显示打日志点的文件名和行数, 有点小坑

	logger.With(zap.String("source", source))
	sugar = logger.Sugar()
}

// SetOutput sets the output writer
func SetOutputEx(outputs []string) {
	// 实现两个判断日志等级的interface (其实 zapcore.*Level 自身就是 interface)
	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.WarnLevel
	})

	warnLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.WarnLevel
	})

	// 获取 info、warn日志文件的io.Writer 抽象 getWriter() 在下方实现
	infoHook = getWriter("game.log")
	warnHook = getWriter("game_error.log")

	// 最后创建具体的Logger
	core = zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(infoHook), infoLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(warnHook), warnLevel),
		zapcore.NewCore(zapcore.NewConsoleEncoder(configex),
			zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stderr)), currentLevel),
	)

	rebuildLoggerFromCfgEx()
}

// SetSource sets the component name (dispatcher/gate/game) of gwlog module
func SetSourceEx(source_ string) {
	source = source_
	rebuildLoggerFromCfgEx()
}

// SetLevel sets the log level
func SetLevelEx(lv Level) {
	currentLevel = lv

	cfg.Level.SetLevel(lv)
}
