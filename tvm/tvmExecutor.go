package tvm

type TvmExecutionResult struct {
	ExitCode       int             `json:"exit_code"`
	GasConsumed    int             `json:"gas_consumed"`
	DataCell       string          `json:"data_cell"`
	ActionListCell string          `json:"action_list_cell"`
	Logs           string          `json:"logs"`
	Stack          []TvmStackEntry `json:"stack"`
}

func getVmFunctionSelector(name string) int {
	if name == "main" {
		return 0
	} else if name == "recv_internal" {
		return 0
	} else if name == "recv_external" {
		return -1
	} else {
		return int(Crc16String(name))
	}
}

//func

//type TvmExecutionConfig struct {
//}
