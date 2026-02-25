package service

import "fmt"

const (
	GitlabUrlSetting   = "GITLAB_URL"
	GitlabTokenSetting = "GITLAB_TOKEN"
	RecentItemsSetting = "RECENT_ITEMS"
	RecentCountInRowSetting = "RECENT_COUNT_IN_ROW"
	OtherIssueSetting = "OTHER_ISSUE"
)

var initSqls = []string{
	`
		create table if not exists actions
		(
			start TEXT PRIMARY KEY,
			finish TEXT,
			name TEXT NOT NULL,
			descr TEXT
		)
	`,
	`
		create table if not exists report_data
		(
			dt DATETIME NOT NULL,
			name TEXT NOT NULL,
			minutes int,
			registered boolean,
			PRIMARY KEY (dt, name)
		)
	`,
	`
		create table if not exists settings
		(
			ident TEXT PRIMARY KEY,
			val TEXT
		)
	`,
	fmt.Sprintf(`insert into settings (ident) values ('%s') ON CONFLICT(ident) DO NOTHING`, GitlabUrlSetting),
	fmt.Sprintf(`insert into settings (ident) values ('%s') ON CONFLICT(ident) DO NOTHING`, GitlabTokenSetting),
	fmt.Sprintf(`insert into settings (ident, val) values ('%s', 'Дейлик,Статус,Демо,Ретро,Предпланирование') ON CONFLICT(ident) DO NOTHING`, RecentItemsSetting),
	fmt.Sprintf(`insert into settings (ident, val) values ('%s', '3') ON CONFLICT(ident) DO NOTHING`, RecentCountInRowSetting),
	fmt.Sprintf(`insert into settings (ident) values ('%s') ON CONFLICT(ident) DO NOTHING`, OtherIssueSetting),
}

var insertActionSql = `
	insert into actions (start, name, descr) values (?, ?, ?)
`

var updateActionFinishSql = `
	update actions set finish = ? where finish is null
`

var getActiveActionSql = `
	select name, coalesce((unixepoch(CURRENT_TIMESTAMP) - unixepoch(start))/60, 0) t
	from actions where finish is null
`

var getMinutesTodaySql = `
	select coalesce(sum((unixepoch(coalesce(a.finish, CURRENT_TIMESTAMP)) - unixepoch(a.start)))/60, 0) t
	from actions a
	where date(a.start, 'localtime') = date(CURRENT_TIMESTAMP, 'localtime')
`

var getReportDataSql = `
		select dt, name, minutes, coalesce(registered, false)
		from report_data
		where date(dt, 'start of month') = date(?, 'start of month')
		order by dt, name

`

var updateReportDataRegisteredSql = `
	update report_data 
	set registered = true 
	where name = ?
	and date(dt, 'localtime') = date(?, 'localtime')
`

var recalcReportSql = `
	insert into report_data (dt, name, minutes)
	select date(a.start, 'localtime') dt, name, sum(
		(unixepoch(coalesce(datetime(a.finish, 'localtime'), date(a.start, '+1 day', 'localtime'))) -
		unixepoch(datetime(a.start, 'localtime'))))/60 t
	from actions a
	where date(a.start, 'localtime') < date(CURRENT_TIMESTAMP, 'localtime')
	and  (date(a.start, 'localtime'), name) not in (select dt, name from report_data)
	group by date(a.start, 'localtime'), name
`

var getSettingsSql = `
	select ident, coalesce(val, '') from settings
`

var updateSettingSql = `
	update settings
	set val = ?
	where ident = ?
`
