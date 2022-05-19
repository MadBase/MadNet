package logging

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/MadBase/MadNet/config"
)

func trimLogs(logFile string, maxAge time.Duration) error {
	f, err := os.OpenFile(logFile, os.O_RDWR, 0755)
	if err != nil {
		return err
	}
	defer f.Close()

	oldBytes, err := findTrimOffset(f, maxAge)
	if err != nil {
		return err
	}

	if oldBytes <= 0 {
		return nil
	}

	f2, err := os.OpenFile(makeTmpFileName(logFile), os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer f2.Close()

	_, err = f.Seek(oldBytes+1, 0)
	if err != nil {
		return err
	}

	_, err = io.Copy(f2, f)
	if err != nil {
		return err
	}

	f.Close()
	f2.Close()
	err = os.Remove(logFile)
	if err != nil {
		return err
	}

	return os.Rename(makeTmpFileName(logFile), logFile)
}

func getFileHandle(fileName string) (*os.File, error) {
	_, tmpFileErr := os.Stat(makeTmpFileName(fileName))
	// if tmp file exists, rename to normal filename if no normal file exists, else delete
	if tmpFileErr == nil {
		_, statErr := os.Stat(fileName)
		if statErr != nil {
			_ = os.Rename(makeTmpFileName(fileName), fileName)
		} else {
			os.Remove(makeTmpFileName(fileName))
		}
	}

	return os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
}

func makeTmpFileName(fileName string) string {
	return fileName + ".tmp"
}

func openLogFile(config config.LogfileConfig) (*os.File, error) {
	if config.FileName == "" {
		return nil, nil
	}

	if config.MaxAgeDays > 0.000001 {
		trimLogs(config.FileName, time.Duration(config.MaxAgeDays*24*float64(time.Hour)))
	}

	f, err := getFileHandle(config.FileName)
	if err != nil {
		return nil, fmt.Errorf("could not open log file %v: %v", config.FileName, err)
	}
	return f, nil
}

func findTrimOffset(f *os.File, maxAge time.Duration) (int64, error) {
	limit := time.Now().Add(-maxAge)
	dec := json.NewDecoder(f)
	var doc struct {
		Time string `json:"time"`
	}
	var oldBytes int64
	for {
		err := dec.Decode(&doc)
		if err != nil {
			if err == io.EOF {
				break
			}
			return 0, err
		}

		time, err := time.Parse(jsonTimeFormat, doc.Time)
		if err != nil {
			return 0, err
		}
		if time.After(limit) {
			break
		} else {
			oldBytes = dec.InputOffset()
		}
	}

	return oldBytes, nil
}
