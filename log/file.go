package log

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

const (
	kChannelSize       = 1024
	kLogExtensionLen   = 4
	kLogCreatedTimeLen = 15 + kLogExtensionLen
	kLogFilenameMinLen = 5 + kLogCreatedTimeLen
)

// file logger
type fileLogger struct {
	config *LogConfig
	file   *os.File
	day    int
	size   int64
	read   chan []byte
	files  int // number of files under `LogPath` currently
}

var gLogger fileLogger

func StartFileLogger(cfg *LogConfig) {
	if err := checkLogConfig(cfg); err != nil {
		fmt.Println(err)
		return
	}
	err := os.MkdirAll(cfg.LogPath, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	gLogger.config = cfg
	if gLogger.config.LogLevel > LogLevelEmpty {
		gLogger.read = make(chan []byte, kChannelSize)
		go gLogger.work()
	}
}

func StopFileLogger() {
	if gLogger.config != nil && gLogger.file != nil {
		gLogger.read <- []byte{}
		<-gLogger.read
	}
}

func LogToFile(level int, mesg string) {
	if len(mesg) > 0 && gLogger.config != nil && level > LogLevelEmpty {
		if level <= gLogger.config.LogLevel {
			gLogger.read <- []byte(mesg)
		}
		if level <= gLogger.config.ConsLevel {
			fmt.Printf(mesg)
		}
	}
}

func (this *fileLogger) work() {
	this.delOldFiles()
	this.reopenLogFile(time.Now())
	for {
		select {
		case mesg := <-this.read:
			if len(mesg) > 0 {
				this.logMsg(mesg)
			} else {
				this.logMsg([]byte("Close log file"))
				close(this.read)
				this.read = nil
				if this.file != nil {
					_ = this.file.Close()
					this.file = nil
				}
				return
			}
		}
	}
}

func (this *fileLogger) logMsg(data []byte) {
	if gLogger.config == nil {
		return
	}
	t := time.Now()
	_, _, d := t.Date()

	if this.size/1024 >= this.config.MaxSize || this.day != d || this.file == nil {
		this.delOldFiles()
		this.reopenLogFile(t)
	}
	if this.file != nil {
		n, _ := this.file.Write(data)
		this.size += int64(n)
	}
}

func (this *fileLogger) delOldFiles() {
	files, err := getLogfilenames(this.config.LogPath, this.config.LogFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	this.files = len(files)
	if this.files >= this.config.MaxFiles {
		nfiles := this.files - this.config.MaxFiles + this.config.DelFiles
		if nfiles > this.files {
			nfiles = this.files
		}
		for i := 0; i < nfiles; i++ {
			filename := this.config.LogPath + string(os.PathSeparator) + files[i]
			fmt.Println("Remove file", filename)
			err := os.RemoveAll(filename)
			if err == nil {
				this.files--
			} else {
				fmt.Print(err)
			}
		}
	}
}

func (this *fileLogger) reopenLogFile(t time.Time) {
	year, mon, day := t.Date()
	hour, min, sec := t.Clock()
	filename := fmt.Sprintf("%s%s%s.%d%02d%02d_%02d%02d%02d.log",
		this.config.LogPath, string(os.PathSeparator),
		this.config.LogFile, year, mon, day, hour, min, sec)
	//		fmt.Println(filename)
	newfile, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	this.files++
	if this.file != nil {
		err = this.file.Close()
	}
	this.file = newfile
	this.day = day
	this.size = 0
}

// sort files by created time embedded in the filename
type byCreatedTime []string

func (a byCreatedTime) Len() int {
	return len(a)
}

func (a byCreatedTime) Less(i, j int) bool {
	s1, s2 := a[i], a[j]
	if len(s1) < kLogFilenameMinLen {
		return true
	} else if len(s2) < kLogFilenameMinLen {
		return false
	} else {
		sa := s1[len(s1)-kLogCreatedTimeLen : len(s1)-kLogExtensionLen]
		sb := s2[len(s2)-kLogCreatedTimeLen : len(s2)-kLogExtensionLen]
		//		fmt.Println(sa, sb)
		return sa < sb
	}
}

func (a byCreatedTime) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

// helpers
func getLogfilenames(dir, name string) ([]string, error) {
	var filenames []string
	var goodnames byCreatedTime
	f, err := os.Open(dir)
	if err == nil {
		filenames, err = f.Readdirnames(0)
		err = f.Close()
		if err == nil {
		}
	}
	for _, file := range filenames {
		if strings.Contains(file, name) && strings.Contains(file, ".log") {
			goodnames = append(goodnames, file)
		}
	}
	sort.Sort(goodnames)
	return goodnames, err
}
