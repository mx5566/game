package gwlog

import "testing"

func TestGwExLog(t *testing.T) {
	SetSourceEx("gwexlog_test")
	SetOutputEx(map[string]string{"errFile": "game.log", "logFile": "game_err.log"})
	SetLevelEx(DebugLevel)

	DebugfE("this is a debufdE %d", 1)
	/*SetLevelEx(InfoLevel)
	DebugfE("SHOULD NOT SEE THIS!")
	InfofE("this is an infoE %d", 2)
	WarnfE("this is a waringE %d", 3)
	TraceErrorEx("this is a traceEx error %d", 4)

	func() {
		defer func() {
			_ = recover()
		}()
		PanicfE("this is a panicEx %d", 5)
	}()

	func() {
		defer func() {
			_ = recover()
		}()
	}()*/
}
