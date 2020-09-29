package gwlog

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
	"runtime"
	"runtime/debug"
	"time"
)

// https://studygolang.com/articles/25044?fr=sidebar
// https://blog.csdn.net/wdy_yx/article/details/79479484
var (
	infoHook io.Writer
	warnHook io.Writer
	core     zapcore.Core
	configex zapcore.EncoderConfig
	atom     zap.AtomicLevel
	loggerEx *zap.Logger
	sourceEx string
)

func init() {
	// 设置一些基本日志格式 具体含义还比较好理解，直接看zap源码也不难懂
	configex = zapcore.EncoderConfig{
		MessageKey:  "msgkey",
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

	atom = zap.NewAtomicLevel()
	atom.SetLevel(zap.DebugLevel)

	// 输出默认到控制台
	core = zapcore.NewTee(zapcore.NewCore(zapcore.NewConsoleEncoder(configex),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)), atom))

	rebuildLoggerFromCfgEx()
}

func getWriter(filename string) io.Writer {
	// 生成rotatelogs的Logger 实际生成的文件名 demo.log.YYmmddHH
	// demo.log是指向最新日志的链接
	// 保存7天内的日志，每1小时(整点)分割一次日志
	str, _ := os.Getwd()

	ostype := runtime.GOOS // 获取系统类型
	//fmt.Println("dir cwd ----------------- " + str + "  type  " + ostype)

	if ostype == "windows" {
		str += "\\log\\"

	} else if ostype == "linux" {
		str += "/log/"
	}

	hook, erre := rotatelogs.New(
		str+filename+".%Y%m%d%H", // 没有使用go风格反人类的format格式
		//rotatelogs.WithLinkName(filename),
		rotatelogs.WithMaxAge(time.Hour*24*7),
		rotatelogs.WithRotationTime(time.Hour),
	)

	if erre != nil {
		panic("getWriter panic " + erre.Error())
	}

	return hook
}

func rebuildLoggerFromCfgEx() {
	if loggerEx != nil {
		_ = loggerEx.Sync()
	}
	// 最后创建具体的Logger
	loggerEx = zap.New(core /*zap.AddCaller()*/) // 需要传入 zap.AddCaller() 才会显示打日志点的文件名和行数, 有点小坑

	loggerEx.With(zap.String("source", sourceEx))
}

// SetOutput sets the output writer
func SetOutputEx(outputs map[string]string) {
	// 实现两个判断日志等级的interface (其实 zapcore.*Level 自身就是 interface)
	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.WarnLevel
	})

	warnLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.WarnLevel
	})

	// 获取 info、warn日志文件的io.Writer 抽象 getWriter() 在下方实现
	if _, ok := outputs["errFile"]; !ok {
		outputs["errFile"] = sourceEx + ".log"
	}

	if _, ok := outputs["logFile"]; !ok {
		outputs["logFile"] = sourceEx + "_error.log"
	}
	infoHook = getWriter(outputs["logFile"])
	warnHook = getWriter(outputs["errFile"])

	//strName := infoHook.(*rotatelogs.RotateLogs).CurrentFileName()

	//fmt.Printf("---------------------------- SetOutputEx %v", outputs)

	// 最后创建具体的Logger
	core = zapcore.NewTee(
		zapcore.NewCore(zapcore.NewConsoleEncoder(configex), zapcore.AddSync(infoHook), infoLevel),
		zapcore.NewCore(zapcore.NewConsoleEncoder(configex), zapcore.AddSync(warnHook), warnLevel),
		zapcore.NewCore(zapcore.NewConsoleEncoder(configex),
			zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)), atom),
	)

	rebuildLoggerFromCfgEx()
}

// SetSource sets the component name (dispatcher/gate/game) of gwlog module
func SetSourceEx(source_ string) {
	sourceEx = source_
	rebuildLoggerFromCfgEx()
}

// TraceError prints the stack and error
func TraceErrorEx(format string, args ...interface{}) {
	ErrorE(string(debug.Stack()))
	ErrorfE(format, args...)
}

// SetLevel sets the log level
func SetLevelEx(lv zapcore.Level) {
	atom.SetLevel(lv)
}

func DebugfE(format string, args ...interface{}) {
	loggerEx.Sugar().With(zap.Time("ts", time.Now())).Debugf(format, args...)
}

func InfofE(format string, args ...interface{}) {
	loggerEx.Sugar().With(zap.Time("ts", time.Now())).Infof(format, args...)
}

func WarnfE(format string, args ...interface{}) {
	loggerEx.Sugar().With(zap.Time("ts", time.Now())).Warnf(format, args...)
}

func ErrorfE(format string, args ...interface{}) {
	loggerEx.Sugar().With(zap.Time("ts", time.Now())).Errorf(format, args...)
}

func PanicfE(format string, args ...interface{}) {
	loggerEx.Sugar().With(zap.Time("ts", time.Now())).Panicf(format, args...)
}

func FatalfE(format string, args ...interface{}) {
	debug.PrintStack()
	loggerEx.Sugar().With(zap.Time("ts", time.Now())).Fatalf(format, args...)
}

func ErrorE(args ...interface{}) {
	loggerEx.Sugar().With(zap.Time("ts", time.Now())).Error(args...)
}

func PanicE(args ...interface{}) {
	loggerEx.Sugar().With(zap.Time("ts", time.Now())).Panic(args...)
}

func FatalE(args ...interface{}) {
	loggerEx.Sugar().With(zap.Time("ts", time.Now())).Fatal(args...)
}
