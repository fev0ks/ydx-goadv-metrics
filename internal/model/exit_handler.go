package model

import "context"

type ExitHandler struct {
	Cancel []context.CancelFunc
	Stop   []chan bool
}
