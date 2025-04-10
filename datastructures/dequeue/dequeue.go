package dequeue

import (
	"fmt"
	"loggergo/logging/core/appender/file/logfile"
)

type Dequeue struct {
	data []logfile.LogFile
}

func NewDequeue(files []logfile.LogFile) *Dequeue {
	return &Dequeue{data: files}
}

func (d *Dequeue) PeekLast() (logfile.LogFile, bool) {
	if len(d.data) > 0 {
		return d.data[len(d.data)-1], true
	}
	return logfile.LogFile{}, false
}

func (d *Dequeue) Add(file logfile.LogFile) {
	d.data = append(d.data, file)
}

func (d *Dequeue) Remove() (logfile.LogFile, error) {
	if len(d.data) == 0 {
		return logfile.LogFile{}, fmt.Errorf("can not remove. dequeue is empty")
	}
	file := d.data[0]
	d.data = d.data[1:]
	return file, nil
}

func (d *Dequeue) Size() int {
	return len(d.data)
}

func (d *Dequeue) Data() []logfile.LogFile {
	return d.data
}
