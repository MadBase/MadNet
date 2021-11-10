package utils

import (
	"crypto/rand"

	"github.com/MadBase/MadNet/constants"
	"github.com/MadBase/MadNet/errorz"
	"github.com/sirupsen/logrus"
)

// ForceSliceToLength will return a byte slice of size length.
// It will left pad a byte slice to the specified number of zeros if the
// slice is not long enough. If the slice is too long, it will return the
// right-most bytes of the slice.
func ForceSliceToLength(inSlice []byte, length int) []byte {
	if len(inSlice) > length {
		return CopySlice(inSlice[len(inSlice)-length:])
	}
	outSlice := make([]byte, length-len(inSlice))
	outSlice = append(outSlice, CopySlice(inSlice)...)
	return outSlice
}

// CopySlice returns a copy of a passed byte slice.
func CopySlice(v []byte) []byte {
	out := make([]byte, len(v))
	copy(out, v)
	return out
}

// Epoch returns the epoch for the corresponding height.
func Epoch(height uint32) uint32 {
	if height <= constants.EpochLength {
		return 1
	}
	if height%constants.EpochLength == 0 {
		return height / constants.EpochLength
	}
	return (height / constants.EpochLength) + 1
}

// ValidateHash checks whether or not hsh has the correct length
func ValidateHash(hsh []byte) error {
	if len(hsh) != constants.HashLen {
		return errorz.ErrInvalid{}.New("the length of the hash is incorrect")
	}
	return nil
}

// RandomBytes will return a byte slice of num random bytes using crypto rand
func RandomBytes(num int) ([]byte, error) {
	b := make([]byte, num)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// DebugTrace allows a traceback to be generated that includes a file name,
// a line number, the error message, and an optional string. Calling this
// function using a logger that is set to anthything other than trace or debug
// level is a no-op. This filtering helps to minimize overhead during normal
// use but still allows error tracebacks to be created easily. The returned
// file and line number will point to where  this function was called.
// Although more than one string may be passed, only the first string will
// be displayed. The varadic property was only used to shorten calling syntax.
func DebugTrace(logger *logrus.Logger, err error, s ...string) {
	// TODO: make more generic, e.g. DebugTrace(l,err); DebugTrace(l,str); DebugTrace(l,pattern,v,v)
	if logger.GetLevel() > logrus.DebugLevel {
		return
	}

	trace := errorz.MakeTrace(1)

	if err != nil {
		if len(s) > 0 {
			logger.WithField("l", trace).Debugf("%v ::: %v", err.Error(), s[0])
			return
		}
		logger.WithField("l", trace).Debugf("%v", err.Error())
		return
	}
	if len(s) > 0 {
		logger.WithField("l", trace).Debugf("%v", s[0])
		return
	}
	logger.WithField("l", trace).Debug("")

}
