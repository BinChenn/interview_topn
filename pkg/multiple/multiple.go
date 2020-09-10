package multiple

// MasterInit initialize service
func MasterInit() {

}

// WorkerInit initialize service
func WorkerInit(workerAddr string) {
	StartWorker(workerAddr)
}
