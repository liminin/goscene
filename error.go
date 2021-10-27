package goscene

import "errors"

var (
	errSceneNotFount        = errors.New("scene not found")
	errInvalidCurrentIndex  = errors.New("invalid current index")
	errUserHasNotActivePlay = errors.New("user hasn't active play")
)
