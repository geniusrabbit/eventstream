package zlogger

import (
	"log"
	"net"
	"net/url"
	"os"

	"go.elastic.co/ecszap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New creates new zap logger object
func New(serviceName, logEncoder, logLevel, addr string, options ...zap.Option) (*zap.Logger, error) {
	var level zapcore.Level
	if err := level.UnmarshalText([]byte(logLevel)); err != nil {
		log.Println("parse log level error", err)
	}
	writer, err := zapWriter(addr)
	if err != nil {
		log.Println("connect to logger address", err)
		return nil, err
	}
	core := zapCore(logEncoder, writer, level)
	switch logEncoder {
	case "logstash", "es", "elastic", "elasticsearch":
		options = append(options, zap.AddCaller())
	default:
	}
	logger := zap.New(core, options...).Named(serviceName)
	return logger, nil
}

func zapNetConnect(address string) (zapcore.WriteSyncer, error) {
	udpAddr, err := url.Parse(address)
	if err != nil {
		return nil, err
	}
	conn, err := net.Dial(udpAddr.Scheme, udpAddr.Host)
	if err != nil {
		return nil, err
	}
	return zapcore.Lock(zapcore.AddSync(conn)), nil
}

func zapWriter(addr string) (zapcore.WriteSyncer, error) {
	writers := []zapcore.WriteSyncer{os.Stderr}
	if addr != "" {
		netConn, err := zapNetConnect(addr)
		if err != nil {
			return nil, err
		}
		writers = append(writers, netConn)
	}
	if len(writers) == 1 {
		return writers[0], nil
	}
	return zapcore.NewMultiWriteSyncer(writers...), nil
}

func zapCore(logEncoder string, writer zapcore.WriteSyncer, level zapcore.Level) zapcore.Core {
	var (
		encoder     zapcore.Encoder
		encoderConf = zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
		}
	)
	switch logEncoder {
	case "logstash", "es", "elastic", "elasticsearch":
		ecsconf := ecszap.NewDefaultEncoderConfig()
		return ecszap.NewCore(ecsconf, writer, level)
	case "json":
		encoder = zapcore.NewJSONEncoder(encoderConf)
	default:
		encoderConf.EncodeLevel = zapcore.LowercaseColorLevelEncoder
		encoder = zapcore.NewConsoleEncoder(encoderConf)
	}
	return zapcore.NewCore(encoder, writer, level)
}
