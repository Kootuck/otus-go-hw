package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

// we are not providing done channel for Stages by task conditions
// to give our pipeline the ability to be shutdown via done/quit channel
// we're wrapping 1st and last stages into additional, closeable channels (wrappers).
func ExecutePipeline(in In, done In, stages ...Stage) Out {
	pipelineEntrance := channelWrapper(in, done)

	// assemble pipeline
	var stageIn In
	var stageOut Out

	stageIn = pipelineEntrance

	for _, s := range stages {
		stageOut = s(stageIn)
		stageIn = stageOut
	}

	return channelWrapper(stageOut, done)
}

func channelWrapper(wrappedChannel In, done In) (wrapper Bi) {
	var msg interface{}
	ok := true
	wrapper = make(Bi)

	go func() {
		defer func() {
			close(wrapper)
		}()

		for ok {
			select {
			case <-done:
				return
			case msg, ok = <-wrappedChannel:
				if !ok {
					return
				}
				wrapper <- msg
			}
		}
	}()

	return wrapper
}
