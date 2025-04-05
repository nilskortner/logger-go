package file

import (
	"fmt"
	"gurms/internal/supportpkgs/datastructures/dequeue"
	"gurms/internal/supportpkgs/mathsupport"
	"io/fs"
	"loggergo/infra/lang"
	"loggergo/logging/core/appender"
	"loggergo/logging/core/appender/file/logfile"
	"loggergo/logging/core/compression"
	"loggergo/logging/core/model/loglevel"
	"loggergo/logging/core/model/logrecord"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
	"timezone"
)

const FIELD_DELIMITER string = "_"
const FILE_MIDDLE = "20060102"

const ARCHIVE_FILE_SUFFIX = ".gz"

const COMPRESSION_LEVEL = 3

var READ_OPTIONS = []int{os.O_RDONLY}
var CREATE_NEW_OPTIONS = []int{os.O_CREATE | os.O_WRONLY | os.O_EXCL}

var APPEND_OPTIONS = os.O_CREATE | os.O_WRONLY | os.O_APPEND

type RollingFileAppender struct {
	channelAppender     *appender.ChannelAppender
	filePrefix          string
	fileSuffix          string
	fileDirectory       string
	fileDirectoryFile   fs.FileInfo
	maxFiles            int
	maxFilesBytes       int64
	minUsableSpaceBytes int64
	files               *dequeue.Dequeue
	currentFile         logfile.LogFile
	nextFileBytes       int64
	nextIndex           int64
	nextDay             int64
	enableCompression   bool
	gzipOutputStream    *compression.FastGzipOutputStream
}

func NewRollingFileAppender(
	level loglevel.LogLevel,
	file string,
	maxFiles int,
	maxFileMb int64,
	enableCompression bool,
) *RollingFileAppender {
	filePath, err := filepath.Abs(file)
	if err != nil {
		fmt.Println(err)
	}
	fileName := filepath.Base(filePath)
	index := strings.LastIndex(fileName, ".")

	var filePrefix string
	var fileSuffix string
	if index == -1 {
		filePrefix = fileName
		fileSuffix = ""
	} else {
		filePrefix = fileName[:index]
		fileSuffix = fileName[index:]
	}
	fileDirectory := filepath.Dir(filePath)
	fileDirectoryFile, err := os.Stat(fileDirectory)
	if err != nil {
		fmt.Println(err)
	}
	maxFiles = mathsupport.Max(maxFiles, 0)
	var maxFileBytes int64
	if maxFileMb > 0 {
		maxFileBytes = maxFileMb * 1024
	} else {
		maxFileBytes = 1024 * 1024 * 1024
	}
	minUsableSpaceBytes := int64(float64(maxFileBytes) * 2.5)

	var gzipOutputStream *compression.FastGzipOutputStream
	if enableCompression {
		gzipOutputStream, err = compression.NewFastGzipOutputStream("a", COMPRESSION_LEVEL, int(maxFileBytes/10))
		if err != nil {
			fmt.Println("internal Logger: ", err)
		}
	}

	err = os.MkdirAll(fileDirectory, 0755)
	if err != nil {
		fmt.Println("internal Logger: Failed to create the directory ("+fileDirectory+")for log files", err)
	}
	var files *dequeue.Dequeue
	files, err = Visit(fileDirectory, filePrefix, fileSuffix, FILE_MIDDLE, maxFiles)
	if err != nil {
		fmt.Println("internal Logger: ", err)
	}

	var nextIndex int64

	rfa := &RollingFileAppender{
		channelAppender:     appender.NewChannelAppender(level),
		filePrefix:          filePrefix,
		fileSuffix:          fileSuffix,
		fileDirectory:       fileDirectory,
		fileDirectoryFile:   fileDirectoryFile,
		maxFiles:            maxFiles,
		maxFilesBytes:       maxFileBytes,
		minUsableSpaceBytes: minUsableSpaceBytes,
		files:               files,
		nextIndex:           nextIndex,
		nextDay:             math.MinInt64,
		enableCompression:   enableCompression,
		gzipOutputStream:    gzipOutputStream,
	}

	logFile, ok := files.PeekLast()
	if !ok {
		rfa.openNewFile(false)
	} else {
		rfa.openExistingFile(logFile)
	}

	return rfa
}

func (r *RollingFileAppender) Append(record logrecord.LogRecord) int {
	if r.needRoll(record) {
		r.roll()
	}
	bytes := r.channelAppender.Append(record)
	r.nextFileBytes += int64(bytes)
	return bytes
}

func (r *RollingFileAppender) GetLevel() loglevel.LogLevel {
	return r.channelAppender.GetLevel()
}

func (r *RollingFileAppender) Close() {
	r.channelAppender.Close()
}

func (r *RollingFileAppender) needRoll(record logrecord.LogRecord) bool {
	return record.Timestamp() >= r.nextDay || r.nextFileBytes >= r.maxFilesBytes
}

func (r *RollingFileAppender) roll() {
	r.cleanFilesUntilSpaceEnough()
	r.closeCurrentFile()
	//if r.enableCompression {
	//compress()
	//}
	r.openNewFile(true)
	r.cleanExceededFiles()
}

func (r *RollingFileAppender) openNewFile(recoverFromError bool) {
	now := time.Now().In(timezone.ZONE_ID)
	next := now.Add(24 * time.Hour)
	dir := r.fileDirectory
	if recoverFromError {
		err := os.MkdirAll(r.fileDirectory, 0755)
		if err != nil {
			if os.IsExist(err) {
				if info, err := os.Stat(dir); err == nil && !info.IsDir() {
					if err := os.Remove(dir); err != nil {
						fmt.Println("internal Logger: Failed to delete file:", err)
					}
					err = os.MkdirAll(dir, os.ModePerm)
					if err != nil {
						fmt.Println("internal Logger: Failed to create directory after deleting file:", err)
					}
				}
			} else {
				fmt.Println("internal Logger: Failed to create directory:", err)
			}
		}
	} else {
		err := os.MkdirAll(r.fileDirectory, 0755)
		if err != nil {
			fmt.Println("internal Logger: Error while creating the directory: "+dir, err)
		}
	}
	var fileName string
	var filePath string
	for {
		nextIndexString := strconv.FormatInt(r.nextIndex, 10)
		if lang.IsBlank(r.filePrefix) {
			fileName = now.Format(FILE_MIDDLE) + FIELD_DELIMITER + nextIndexString + r.fileSuffix
		} else {
			fileName = r.filePrefix + FIELD_DELIMITER + now.Format(FILE_MIDDLE) + FIELD_DELIMITER + nextIndexString + r.fileSuffix
		}
		filePath = filepath.Join(r.fileDirectory, fileName)
		var err error
		r.channelAppender.File, err = os.OpenFile(filePath, APPEND_OPTIONS, 0666)
		if err != nil {
			message := "failed to open the file: " + filePath
			if recoverFromError {
				println("internal logger: " + message)
				r.nextIndex++
				time.Sleep(1 * time.Second)
			}
		} else {
			break
		}
	}
	var archivePath string
	if r.enableCompression {
		filepath.Join(filepath.Dir(filePath), fileName+ARCHIVE_FILE_SUFFIX)
	} else {
		archivePath = ""
	}
	r.currentFile = logfile.NewLogFile(filePath, archivePath, now, r.nextIndex)
	r.files.Add(r.currentFile)

	fileInfo, err := r.channelAppender.File.Stat()
	if err != nil {
		fmt.Println("internal Logger: failed to get file stats:"+filePath, err)
		if recoverFromError {
			fmt.Println(". Fallback to 0")
			r.nextFileBytes = 0
		}
	} else {
		r.nextFileBytes = fileInfo.Size()
	}
	r.nextDay = next.UnixNano()
	r.nextIndex++
}

func (r *RollingFileAppender) openExistingFile(existingFile logfile.LogFile) {
	now := existingFile.GetTime()
	next := now.AddDate(0, 0, 1)
	next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())

	filePath := existingFile.GetPath()

	var err error
	r.channelAppender.File, err = openFile(filePath)
	if err != nil {
		fmt.Println("internal Logger: failed to open file:", err)
	}
	r.currentFile = existingFile

	fileInfo, err := r.channelAppender.File.Stat()
	if err != nil {
		fmt.Println("internal Logger: failed to get file stats:"+filePath, err)
	} else {
		r.nextFileBytes = fileInfo.Size()
	}
	r.nextDay = next.UnixNano()
	r.nextIndex = existingFile.GetIndex() + 1
}

func openFile(filePath string) (*os.File, error) {
	directory := filepath.Dir(filePath)
	if directory != "." {
		err := os.MkdirAll(directory, 0755)
		if err != nil {
			return nil, fmt.Errorf(
				"internal logger: failed to create the directory (%s) for log files: %w", directory, err)
		}
	}
	file, err := os.OpenFile(filePath, APPEND_OPTIONS, 0666)
	if err != nil {
		return nil, fmt.Errorf("internal logger: failed to open the file: %s: %w", filePath, err)
	}
	return file, nil
}

func (r *RollingFileAppender) closeCurrentFile() {
	err := r.channelAppender.File.Close()
	if err != nil {
		fmt.Println("internal logger: couldnt close File: %w", err)
	}
}

func (r *RollingFileAppender) cleanExceededFiles() {
	for r.files.Size() > r.maxFiles {
		file, err := r.files.Remove()
		if err != nil {
			fmt.Println("internal logger: %w", err)
		}
		if r.maxFiles > 0 {
			deleteLogFile(file)
		}
	}
}

func (r *RollingFileAppender) cleanFilesUntilSpaceEnough() {
	file := r.fileDirectory
	if getUsableSpace(file) > r.minUsableSpaceBytes {
		return
	}
	lastElementIndex := r.files.Size() - 1
	i := 0
	for _, logFile := range r.files.Data() {
		if i == lastElementIndex {
			break
		}
		deleteLogFile(logFile)
		if getUsableSpace(file) > r.minUsableSpaceBytes {
			return
		}
		i++
	}
	fmt.Println("internal logger: freeze to wait for the space to be available")
	for {
		time.Sleep(1000 * time.Millisecond)
		if getUsableSpace(file) > r.minUsableSpaceBytes {
			break
		}
	}
	fmt.Println("internal logger: the space has become available")
}

func deleteLogFile(logFile logfile.LogFile) {
	path := logFile.GetPath()
	archivePath := logFile.GetArchivePath()

	err := os.Remove(path)
	if err != nil && !os.IsNotExist(err) {
		fmt.Println("internal logger: caught an error while deleting the log file: ", err)
	}
	if archivePath != "" {
		err = os.Remove(archivePath)
		if err != nil && !os.IsNotExist(err) {
			fmt.Println("internal logger: caught an error while deleting the log file: ", err)
		}
	}
}

func getUsableSpace(path string) int64 {
	var stat syscall.Statfs_t

	// Get file system statistics using Statfs
	err := syscall.Statfs(path, &stat)
	if err != nil {
		return 0
	}

	// Available blocks * size per block = available space in bytes
	usableSpace := int64(stat.Bavail) * stat.Bsize

	return usableSpace
}
