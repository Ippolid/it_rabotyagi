package logger

import (
	"os"
	"strings"

	"gopkg.in/natefinch/lumberjack.v2"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func getCore(level zap.AtomicLevel) zapcore.Core {
	// Настройка вывода в консоль
	stdout := zapcore.AddSync(os.Stdout)
	developmentCfg := zap.NewDevelopmentEncoderConfig()
	developmentCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	developmentCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder // Цветной вывод уровня
	consoleEncoder := zapcore.NewConsoleEncoder(developmentCfg)

	//Настройка вывода в файл с ротацией
	file := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "logs/app.log", // Путь к файлу логов
		MaxSize:    10,             // Размер файла в мегабайтах
		MaxBackups: 3,              // Количество старых файлов для хранения
		MaxAge:     7,              // Количество дней для хранения файлов
		Compress:   true,           // Сжимать старые файлы
	})

	productionCfg := zap.NewProductionEncoderConfig()
	productionCfg.TimeKey = "timestamp"
	productionCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	fileEncoder := zapcore.NewJSONEncoder(productionCfg)

	// Объединение выводов в консоль и файл
	return zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, stdout, level),
		zapcore.NewCore(fileEncoder, file, level),
	)
}

func getAtomicLevel(level string) (zap.AtomicLevel, error) {
	var zapLevel zapcore.Level
	// Преобразуем строку в уровень zap
	if err := zapLevel.Set(strings.ToLower(level)); err != nil {
		return zap.NewAtomicLevel(), err
	}
	return zap.NewAtomicLevelAt(zapLevel), nil
}

// InitLocalLogger инициализирует локальный логгер с заданным уровнем
func InitLocalLogger(level string) {
	zaplevel, err := getAtomicLevel(level)
	if err != nil {
		panic("Invalid log level: " + level)
	}

	globalLogger = zap.New(getCore(zaplevel))

}

var globalLogger *zap.Logger

// Init инициализирует глобальный логгер с заданным core и опциями
func Init(core zapcore.Core, options ...zap.Option) {
	globalLogger = zap.New(core, options...)
}

// Debug функция для логирования отладочных сообщений
func Debug(msg string, fields ...zap.Field) {
	globalLogger.Debug(msg, fields...)
}

// Info функция для логирования информационных сообщений
func Info(msg string, fields ...zap.Field) {
	globalLogger.Info(msg, fields...)
}

// Warn функция для логирования предупреждений
func Warn(msg string, fields ...zap.Field) {
	globalLogger.Warn(msg, fields...)
}

// Error функция для логирования ошибок
func Error(msg string, fields ...zap.Field) {
	globalLogger.Error(msg, fields...)
}

// Fatal функция для логирования фатальных ошибок
func Fatal(msg string, fields ...zap.Field) {
	globalLogger.Fatal(msg, fields...)
}

// Logger возвращает глобальный логгер
func Logger() *zap.Logger {
	return globalLogger
}

// WithOptions возвращает новый логгер с заданными опциями
func WithOptions(opts ...zap.Option) *zap.Logger {
	return globalLogger.WithOptions(opts...)
}
