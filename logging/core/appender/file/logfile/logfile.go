package logfile

import "time"

type LogFile struct {
	path        string
	archivePath string
	dateTime    time.Time
	index       int64
}

func NewLogFile(path string, archivePath string, dateTime time.Time, index int64) LogFile {
	return LogFile{
		path:        path,
		archivePath: archivePath,
		dateTime:    dateTime,
		index:       index,
	}
}

func (l LogFile) GetIndex() int64 {
	return l.index
}

func (l LogFile) GetPath() string {
	return l.path
}

func (l LogFile) GetArchivePath() string {
	return l.archivePath
}

func (l LogFile) GetTime() time.Time {
	return l.dateTime
}
