package file

import (
	"fmt"
	"loggergo/datastructures/dequeue"
	"loggergo/datastructures/treeset"
	"loggergo/filesupport"
	"loggergo/infra/timezone"
	"loggergo/logging/core/appender/file/logfile"
	"loggergo/mathsupport"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type LogDirectoryVisitor struct {
	filePrefix        string
	fileSuffix        string
	fileMiddle        string
	maxFilesToKeep    int
	deleteExceedFiles bool
	files             *treeset.Tree[logfile.LogFile]
}

func NewLogDirectoryVisitor(
	filePrefix string,
	fileSuffix string,
	fileMiddle string,
	maxFiles int) *LogDirectoryVisitor {
	return &LogDirectoryVisitor{
		filePrefix:        filePrefix,
		fileSuffix:        fileSuffix,
		fileMiddle:        fileMiddle,
		maxFilesToKeep:    mathsupport.Max(1, maxFiles),
		deleteExceedFiles: maxFiles > 0,
		files:             treeset.New(treeset.LogComparator),
	}
}

func (l *LogDirectoryVisitor) visitFile(path string) {
	name := filepath.Base(path)
	if !l.isLogFile(name) {
		return
	}
	indexEnd := len(name) - len(l.fileSuffix)
	indexStart := strings.LastIndex(name[:indexEnd], FIELD_DELIMITER)
	if indexStart == len(l.filePrefix)+len(l.fileMiddle)+1 {
		index, err := strconv.ParseInt(name[indexStart+1:indexEnd], 10, 64)
		if err != nil {
			fmt.Println("Error parsing string to int64:", err)
		}
		timestamp, err := time.Parse(l.fileMiddle, name[len(l.filePrefix)+1:indexStart])
		if err != nil {
			fmt.Println("Error parsing string to timestamp:", err)
			return
		}
		timestamp = timestamp.In(timezone.ZONE_ID)
		l.handleNewLogFile(path, timestamp, index)
	}
}

func (l *LogDirectoryVisitor) isLogFile(name string) bool {
	return (len(name) > len(l.filePrefix)+len(l.fileSuffix)+len(l.fileMiddle)+1) &&
		strings.HasPrefix(name, l.filePrefix) &&
		strings.HasSuffix(name, l.fileSuffix)
}

func (l *LogDirectoryVisitor) handleNewLogFile(path string, timestamp time.Time, index int64) {
	fileName := filepath.Base(path)
	isArchive := strings.HasSuffix(fileName, ARCHIVE_FILE_SUFFIX)
	var filePath string
	var archivePath string
	if isArchive {
		baseFileName := strings.TrimSuffix(fileName, ARCHIVE_FILE_SUFFIX)
		filePath = resolveSibling(path, baseFileName)
		archivePath = path
	} else {
		filePath = path
		archivePath = resolveSibling(path, fileName+ARCHIVE_FILE_SUFFIX)
	}
	file := logfile.NewLogFile(filePath, archivePath, timestamp, index)
	l.files.Put(file)
	if l.files.GetSize() > l.maxFilesToKeep {
		firstLogFile := l.files.GetRootKey()
		l.files.Remove(firstLogFile)
		if l.deleteExceedFiles {
			filesupport.DeleteIfExists(firstLogFile.GetPath())
			firstArchivePath := firstLogFile.GetArchivePath()
			if firstArchivePath != "" {
				filesupport.DeleteIfExists(firstArchivePath)
			}
		}
	}
}

func resolveSibling(basePath, sibling string) string {
	dir := filepath.Dir(basePath)
	siblingPath := filepath.Join(dir, sibling)
	return siblingPath
}

func Visit(
	directory string,
	prefix string,
	suffix string,
	middle string,
	maxFiles int) (*dequeue.Dequeue, error) {
	visitor := NewLogDirectoryVisitor(prefix, suffix, middle, maxFiles)
	maxDepth := 1
	err := filepath.WalkDir(directory, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("failed to visit the directory: %v, %v", path, err)
		}

		depth := strings.Count(path, string(os.PathSeparator)) - strings.Count(directory, string(os.PathSeparator))
		if depth > maxDepth {
			return filepath.SkipDir
		}
		visitor.visitFile(path)
		return nil

	})
	if err != nil {
		return nil, err
	}
	return dequeue.NewDequeue(visitor.files.Keys()), nil
}
