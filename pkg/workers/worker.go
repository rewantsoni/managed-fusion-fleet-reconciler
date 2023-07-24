package workers

import (
	"github.com/red-hat-storage/managed-fusion-fleet-reconciler/pkg/utils"
)

type WorkHandler[Req, Res any] func(Req) Res

func GoWorker[Req, Res any](
	reqChan <-chan Req,
	resChan chan<- utils.Pair[Req, Res],
	handler WorkHandler[Req, Res],
) {
	go func() {
		for req := range reqChan {
			resChan <- utils.NewPair(req, processRequest(req, handler))
		}
	}()
}

func processRequest[Req, Res any](req Req, handler WorkHandler[Req, Res]) Res {
	defer recoverWorker(req, handler)
	return handler(req)
}

func recoverWorker[Req, Res any](req Req, handler WorkHandler[Req, Res]) {
	if r := recover(); r != nil {
		// Re-process the request
		processRequest(req, handler)
	}
}
