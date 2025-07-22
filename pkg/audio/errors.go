package audio

import (
	"github.com/vistormu/go-dsa/errors"
)

const (
	AudioPathNotFound      errors.ErrorType = "audio path not found"
	UnsupportedAudioFormat errors.ErrorType = "unsupported audio format"
	LoadAudioError         errors.ErrorType = "load audio error"
	AudioHandleNotFound    errors.ErrorType = "audio handle not found"
	FilesystemNotSet       errors.ErrorType = "audio filesystem not set"
)
