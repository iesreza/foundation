package system

/*func SetAlarm()  {
	c := Config()
	if c.App.OS == "linux" {
		golarm.AddAlarm(golarm.SystemLoad(golarm.OneMinPeriod).AboveEqual(c.Alarm.Processor).Run(func() {

		}))

		golarm.AddAlarm(golarm.SystemMemory().Used().Above(c.Alarm.Memory).Percent().Run(func() {
			fmt.Println("Used system memory > 90% !!")
		}))

		if c.Alarm.Reset {
			// checks if the system has been running for less than 1 minute
			golarm.AddAlarm(golarm.SystemUptime().Below(1).Run(func() {
				fmt.Println("System just started !!")
			}))
		}
	}

}*/