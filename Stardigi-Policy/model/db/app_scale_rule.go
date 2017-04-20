package db

type AppScaleRule struct {
	MarathonName   string `staro:"column(marathon_name);union(pk)"`
	AppId          string `staro:"column(app_id);union(pk)"`
	ScaleType      string `staro:"column(scale_type);union(pk)"`
	ScaleThreshold int    `staro:"column(scale_threshold)"`
	PerAutoScale   int    `staro:"column(per_auto_scale)"`
	Memory         int    `staro:"column(memory)"`
	Cpu            int    `staro:"column(cpu)"`
	Thread         int    `staro:"column(thread)"`
	HaQueue        int    `staro:"column(ha_queue)"`
	Switch         int    `staro:"column(switch)"`
	ColdTime       int    `staro:"column(cold_time)"`
	CollectPeriod  int    `staro:"column(collect_period)"`
	ContinuePeriod int    `staro:"column(continue_period)"`
}
