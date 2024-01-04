package log

import (
	"context"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	hertzzap "github.com/hertz-contrib/logger/zap"
	"gopkg.in/natefinch/lumberjack.v2"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func Init() {
	// 可定制的输出目录。
	var logFilePath string
	logFilePath = "./logs/"
	if err := os.MkdirAll(logFilePath, 0o777); err != nil {
		log.Printf("err: %+v", err)
		return
	}

	// 将文件名设置为日期
	logFileName := time.Now().Format("2006-01-02") + ".log"
	fileName := path.Join(logFilePath, logFileName)
	if _, err := os.Stat(fileName); err != nil {
		if _, err := os.Create(fileName); err != nil {
			log.Printf("err: %+v", err)
			return
		}
	}

	logger := hertzzap.NewLogger()
	// 提供压缩和删除
	lumberjackLogger := &lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    1024, // 一个文件最大可达20M。
		MaxBackups: 5,    // 最多同时保存 5 个文件。
		MaxAge:     10,   // 一个文件最多可以保存 10 天。
		Compress:   true, // 用 gzip 压缩。
	}

	logger.SetOutput(lumberjackLogger)
	logger.SetLevel(hlog.LevelDebug)

	hlog.SetLogger(logger)
}

//func KeepNewDateLogFile() {
//	log.SetFlags(log.Llongfile | log.LstdFlags)
//
//	go func() {
//		defer func() {
//			if err := recover(); err != nil {
//				log.Printf("err: %+v", err)
//			}
//		}()
//
//		ticker := time.NewTicker(time.Minute)
//		select {
//		case <-ticker.C:
//			now := time.Now()
//			if now.Hour() == 0 && now.Minute() == 0 {
//				log.SetOutput(GetMultiWriter())
//
//				db.Db.Config.Logger = logger.New(log.New(utils.GetMultiWriter(), "", log.LstdFlags),
//					logger.Config{
//						SlowThreshold:             time.Second,
//						Colorful:                  true,
//						IgnoreRecordNotFoundError: true,
//						ParameterizedQueries:      false,
//						LogLevel:                  logger.Info,
//					})
//			}
//		}
//	}()
//}

func CtxTrace(ctx context.Context, str string) {
	//hlog.CtxTracef()
}

func Errorf(ctx context.Context, format string, v ...interface{}) {
	var traceId string
	if ctx != nil {
		if ctx.Value("trace_id") != nil {
			traceId = ctx.Value("trace_id").(string)
		}
	}

	log.Printf(FileWithLineNum()+" [ERROR] <"+traceId+"> "+logger.Red+format+logger.Reset, v...)
}

func Infof(ctx context.Context, format string, v ...interface{}) {
	var traceId string
	if ctx != nil {
		if ctx.Value("trace_id") != nil {
			traceId = ctx.Value("trace_id").(string)
		}
	}

	log.Printf(FileWithLineNum()+" [INFO] <"+traceId+"> "+logger.Green+format+logger.Reset, v...)
}

// FileWithLineNum return the file name and line number of the current file
func FileWithLineNum() string {
	// the second caller usually from gorm internal, so set i start from 2
	for i := 2; i < 15; i++ {
		_, file, line, ok := runtime.Caller(i)
		if ok && !(strings.Contains(file, "common/model") || strings.Contains(file, "db/model") || strings.Contains(file, "gorm")) {
			return file + ":" + strconv.FormatInt(int64(line), 10)
		}
	}

	return ""
}
