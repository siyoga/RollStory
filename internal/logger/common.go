package logger

import (
	"runtime"
	"strings"
	"unicode"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func (l *logger) getFunctionsCallTrace(skip int) []runtime.Frame {
	var (
		stackSize = 64
		callers   = make([]uintptr, stackSize)
		n         = runtime.Callers(0, callers)
		frames    = runtime.CallersFrames(callers[:n])
		toSkip    = 2 + skip // 2 т.к. нужно скипнуть вызов runtime.Callers и самой функции getFunctionsCallTrace
	)

	var resultFrames []runtime.Frame
	for frame, hasNext := frames.Next(); hasNext; frame, hasNext = frames.Next() {
		if strings.HasPrefix(frame.Function, l.moduleName) {
			if strings.Contains(frame.File, "/handler/context.go") {
				break
			}
			resultFrames = append(resultFrames, frame)
		}
	}

	return resultFrames[toSkip:]
}

func (l *logger) getMethod(frames []runtime.Frame) string {
	for _, frame := range frames {
		functionName := l.getFunctionName(frame)
		if unicode.IsUpper(rune(functionName[0])) {
			return functionName
		}
	}

	return l.getFunctionName(frames[0])
}

func (l *logger) getFunctionName(frame runtime.Frame) string {
	return frame.Function[strings.LastIndex(frame.Function, ".")+1:]
}

func (l *logger) setZapMethod(frames []runtime.Frame) zap.Field {
	method := l.getMethod(frames)
	return zap.String(logMethodField, method)
}

func (l *logger) zapReason(reason string) zap.Field {
	return zap.String(logReasonField, reason)
}

func (l *logger) zapStacktrace(frames []runtime.Frame) zap.Field {
	return zap.Uintptr(logStacktraceField, frames[0].PC)
}

func (l *logger) zapFullStacktrace(frames []runtime.Frame) zap.Field {
	return zap.Reflect(logFullStacktraceField, frames)
}

func (l *logger) zapCaller(frames []runtime.Frame) zap.Field {
	f := frames[0]
	return zap.Reflect(logCallerField, zapcore.EntryCaller{
		Defined:  true,
		PC:       f.PC,
		File:     f.File,
		Line:     f.Line,
		Function: f.Function,
	})
}

func (l *logger) genError(
	err error, reason string,
) (string, zap.Field, zap.Field, zap.Field, zap.Field, zap.Field) {
	frames := l.getFunctionsCallTrace(1)
	return err.Error(), l.setZapMethod(frames), l.zapReason(reason), l.zapStacktrace(frames), l.zapCaller(frames), l.zapFullStacktrace(frames)
}
