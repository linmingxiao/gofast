package logx

type LogConfig struct {
	ServiceName         string `json:",optional"`
	Mode                string `json:",default=console,options=console|file|volume"`
	Path                string `json:",default=logs"`
	Level               string `json:",default=info,options=info|error|severe"`
	Compress            bool   `json:",optional"`
	KeepDays            int    `json:",optional"`
	StackCooldownMillis int    `json:",default=100"`
	NeedCpuMem          bool   `json:",default=true"`
	Style               string `json:",default=json"`
}

var theConfig *LogConfig
