package render

import (
	"github.com/vistormu/go-dsa/errors"
)

const (
	ImagePathNotFound   errors.ErrorType = "image path not found"
	ImageLoadError      errors.ErrorType = "image load error"
	ImageHandleNotFound errors.ErrorType = "image handle not found"
	FilesystemNotSet    errors.ErrorType = "image filesystem not set"
)
