package plugin

func (p *Plugin) Close() {
	p.DbRead.Close()
}
