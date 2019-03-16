package log

import (
	"fmt"
	"os"
	"path"
	"sort"
	"strings"
	"time"
)

const (
	kChannelSize = 1024
	//	kMaxInt64			= int64(^uint64(0) >> 1)
	kLogExtensionLen   = 4
	kLogCreatedTimeLen = 15 + kLogExtensionLen
	kLogFilenameMinLen = 5 + kLogCreatedTimeLen
)

var gProgname = path.Base(os.Args[0])

// init is called after all the variable declarations in the package have evaluated their initializers,
// and those are evaluated only after all the imported packages have been initialized.
// Besides initializations that cannot be expressed as declarations, a common use of init functions is to verify
// or repair correctness of the program state before real execution begins.
func init() {
	tmpProgname := strings.Split(gProgname, "\\") // for compatible with `go run` under Windows
	gProgname = tmpProgname[len(tmpProgname)-1]
}

// logger
type fileLogger struct {
	config *LogConfig
	file   *os.File
	day    int
	size   int64
	read   chan []byte
	files  int // number of files under `logPath` currently
}

var gLogger fileLogger

func StartFileLogger(cfg *LogConfig) {
	gLogger.config = cfg
	if gLogger.config != nil && gLogger.config.LogLevel > LogLevelEmpty {
		gLogger.read = make(chan []byte, kChannelSize)
		go gLogger.work()
	}
}

func StopFileLogger() {
	if gLogger.config != nil && gLogger.config.LogLevel > LogLevelEmpty {
		gLogger.read <- []byte{}
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
				break
			}
		}
	}
	close(this.read)
	if this.file != nil {
		this.file.Close()
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
	files, err := getLogfilenames(this.config.LogPath)
	if err != nil {
		fmt.Print(err)
		return
	}
	this.files = len(files)
	if this.files >= this.config.MaxFiles {
		sort.Sort(byCreatedTime(files))
		//			fmt.Println(files)
		nfiles := this.files - this.config.MaxFiles + this.config.DelFiles
		if nfiles > this.files {
			nfiles = this.files
		}
		for i := 0; i < nfiles; i++ {
			fmt.Println("Remove file", files[i])
			err := os.RemoveAll(this.config.LogPath + files[i])
			if err == nil {
				this.files--
			} else {
				fmt.Print(err)
			}
		}
	}
}

func (this *fileLogger) reopenLogFile(t time.Time) {
	y, m, d := t.Date()
	hour, min, sec := t.Clock()
	filename := fmt.Sprintf("%s%s.%d%02d%02d_%02d%02d%02d.log",
		this.config.LogPath, this.config.LogFile, y, m, d, hour, min, sec)
	//		fmt.Println(filename)
	newfile, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Print(err)
		return
	}
	this.files++
	if this.file != nil {
		this.file.Close()
	}
	this.file = newfile
	this.day = d
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
func getLogfilenames(dir string) ([]string, error) {
	var filenames []string
	f, err := os.Open(dir)
	if err == nil {
		filenames, err = f.Readdirnames(0)
		f.Close()
		if err == nil {
		}
	}
	return filenames, err
}
