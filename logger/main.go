package logger

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

//prefixes
const (
	debugPrefix   = "Debug"
	infoPrefix    = "Info"
	warningPrefix = "Warning"
	errorPrefix   = "Error"
	fatalPrefix   = "Fatal"
)

//trace levels
const (
	DebugTraceLevel = iota + 1
	InfoTraceLevel
	WarningTraceLevel
	ErrorTraceLevel
	FatalTraceLevel
)

//time format
const (
	timeFormat = time.RFC3339
)

//file mode
const (
	fileMode = 0644
)

//default parameters
const (
	defaultFileSize        = 1000
	defaultFilename        = "./logs/error.log"
	defaultReverseFilename = "./logs/view_error.log"
	defaultTraceLevel      = DebugTraceLevel
)

type Logger struct {
	Filename   string
	TraceLevel int
	mutex      *sync.Mutex
}

func DefaultLogger() *Logger {
	return &Logger{
		Filename:   defaultFilename,
		TraceLevel: defaultTraceLevel,
		mutex:      &sync.Mutex{},
	}
}

func NewLogger(filename string, traceLevel int) *Logger {
	l := &Logger{
		Filename:   filename,
		TraceLevel: traceLevel,
		mutex:      &sync.Mutex{},
	}
	return l
}

func (l *Logger) Debug(values ...interface{}) {
	messageFormat := createMessageFormat(values...)
	l.log(DebugTraceLevel, debugPrefix, messageFormat, values...)
}

func (l *Logger) Debugf(messageFormat string, values ...interface{}) {
	l.log(DebugTraceLevel, debugPrefix, messageFormat, values...)
}

func (l *Logger) Info(values ...interface{}) {
	messageFormat := createMessageFormat(values...)
	l.log(InfoTraceLevel, infoPrefix, messageFormat, values...)
}

func (l *Logger) Infof(messageFormat string, values ...interface{}) {
	l.log(InfoTraceLevel, infoPrefix, messageFormat, values...)
}

func (l *Logger) Warning(values ...interface{}) {
	messageFormat := createMessageFormat(values...)
	l.log(WarningTraceLevel, warningPrefix, messageFormat, values...)
}

func (l *Logger) Warningf(messageFormat string, values ...interface{}) {
	l.log(WarningTraceLevel, warningPrefix, messageFormat, values...)
}

func (l *Logger) Error(values ...interface{}) {
	messageFormat := createMessageFormat(values...)
	l.log(ErrorTraceLevel, errorPrefix, messageFormat, values...)
}

func (l *Logger) Errorf(messageFormat string, values ...interface{}) {
	l.log(ErrorTraceLevel, errorPrefix, messageFormat, values...)
}

func (l *Logger) Fatal(values ...interface{}) {
	messageFormat := createMessageFormat(values...)
	l.log(FatalTraceLevel, fatalPrefix, messageFormat, values...)
	os.Exit(1)
}

func (l *Logger) Fatalf(messageFormat string, values ...interface{}) {
	l.log(FatalTraceLevel, fatalPrefix, messageFormat, values...)
	os.Exit(1)
}

func createMessageFormat(values ...interface{}) string {
	messageFormat := strings.Repeat("%v, ", len(values))
	messageFormat = strings.Trim(messageFormat, ", ")
	return messageFormat
}

func (l *Logger) log(traceLevel int, prefix string, messageFormat string, values ...interface{}) {
	//check trace level
	if l.TraceLevel > traceLevel {
		return
	}
	//synchronization
	l.mutex.Lock()
	defer l.mutex.Unlock()
	//create message
	message := fmt.Sprintf(messageFormat, values...)
	//replace new line characters with white spaces
	message = strings.Replace(message, "\n", " ", -1)
	//create formatted message
	message = fmt.Sprintf("[%s][%s][%s]\n", time.Now().Format(timeFormat), prefix, message)
	//open log file
	logFile, err := os.OpenFile(l.Filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, fileMode)
	if err != nil {
		fmt.Println("opening log file failed")
		return
	}
	defer func() {
		err := logFile.Close()
		if err != nil {
			fmt.Println("closing log file failed")
		}
		ReverseWriter()
	}()

	//write to log file
	_, err = logFile.WriteString(message)
	if err != nil {
		fmt.Println("writing log failed")
	}
}

func ReverseWriter() {
	readFile, err := os.Open(defaultFilename)
	if err != nil {
		fmt.Printf("failed to open file: %s", err)
	}

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	var fileTextLines []string
	for fileScanner.Scan() {
		fileTextLines = append(fileTextLines, fileScanner.Text())
	}

	logFile, err := os.OpenFile(defaultReverseFilename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, fileMode)
	if err != nil {
		fmt.Println("opening log file failed")
		return
	}
	defer func() {

		err := logFile.Close()
		if err != nil {
			fmt.Println("closing log file failed")
		}

	}()
	err = os.Truncate(defaultReverseFilename, 0)
	if err != nil {
		fmt.Println(err)

	}
	for i := len(fileTextLines) - 1; i >= 0; i-- {
		//write to log file
		_, err = logFile.WriteString(fileTextLines[i] + "\n")
		if err != nil {
			fmt.Println("writing log failed")
		}
	}

	fi, err := readFile.Stat()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(fi.Size())
	if fi.Size() > defaultFileSize {
		os.Remove(defaultFilename)
		os.Rename(defaultReverseFilename, defaultReverseFilename+"."+time.Now().String())
	}
	readFile.Close()

}

func ServeLogFiles() {
	port := flag.String("p", "8100", "port to serve on")
	directory := flag.String("d", "./logs", "the directory of static file to host")
	flag.Parse()
	http.Handle("/", http.FileServer(http.Dir(*directory)))

	panic(http.ListenAndServe(":"+*port, nil))
}
