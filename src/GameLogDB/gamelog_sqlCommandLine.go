package GameLogDB

import (
	"sync"
)

var sql_col = make(map[string]string)
var onceDo sync.Once

func GetSqlCommand(name string) string {
	initCreatDialy()
	onceDo.Do(initOther)
	return sql_col[name]
}

func initOther() {
	sql_col["create_log_login_detail"] = "CREATE TABLE IF NOT EXISTS `log_login_detail` (" +
		" `user_id` INT NOT NULL ," +
		" `channel` VARCHAR(20) NOT NULL DEFAULT 'local'," +
		" `db_id` INT NOT NULL," +
		" `level` INT NOT NULL," +
		" `viplevel` INT NOT NULL," +
		" `time` DATETIME NOT NULL," +
		" `reg` DATETIME NOT NULL," +
		" KEY `uid_time`(`user_id`, `time`)" +
		" ) ENGINE=MyISAM DEFAULT CHARSET=utf8;"

	sql_col["create_login_log"] = "CREATE TABLE IF NOT EXISTS `log_users_login` (" +
		"`username` varchar(50) NOT NULL," +
		"`cnt` int(11) DEFAULT '1' COMMENT '登录次数'," +
		"`lastdate` date DEFAULT '0000-00-00' COMMENT '最后登录日期'," +
		"`firstdate` date NOT NULL," +
		"PRIMARY KEY (`username`)" +
		") ENGINE=MyISAM;"

	sql_col["create_retention_log"] = "CREATE TABLE IF NOT EXISTS `log_retention_users` (" +
		"`daytime` date NOT NULL DEFAULT '0000-00-00' COMMENT '统计日期'," +
		"`dbid` varchar(50) NOT NULL DEFAULT '0' COMMENT '数据库编号'," +
		"`cid` varchar(50) NOT NULL DEFAULT 'local' COMMENT '渠道编号'," +
		"`first` int(11) DEFAULT '0' COMMENT '新玩家数'," +
		"`second` int(11) DEFAULT '0' COMMENT '次日留存数'," +
		"`third` int(11) DEFAULT '0' COMMENT '3日留存数'," +
		"`fourth`int(11) DEFAULT '0'," +
		"`fifth`int(11) DEFAULT '0'," +
		"`sixth`int(11) DEFAULT '0'," +
		"`week` int(11) DEFAULT '0' COMMENT '7日留存数'," +
		"`halfmonth` int(11) DEFAULT '0' COMMENT '15日留存数'," +
		"`month` int(11) DEFAULT '0' COMMENT '30日留存数'," +
		"`interval_second_login` int(11) DEFAULT '0' COMMENT '第二次登陆为第二天的玩家人数统计'," +
		"PRIMARY KEY (`daytime`,`dbid`,`cid`)," +
		"KEY `idx_dbid` (`dbid`)," +
		"KEY `idx_cid` (`cid`)" +
		") ENGINE Myisam DEFAULT CHARSET=utf8;"
}

func initCreatDialy() {
	sql_col["create_daily_log"] = "CREATE TABLE IF NOT EXISTS `log_dialy` (" +
		" `id` INT AUTO_INCREMENT," +
		" `user_id` INT NOT NULL ," +
		" `channel` VARCHAR(20) NOT NULL DEFAULT 'local'," +
		" `action` VARCHAR(100) NOT NULL," +
		" `time` TIMESTAMP NOT NULL ," +
		" `diamond` INT DEFAULT '0'," +
		" `coin` INT DEFAULT '0'," +
		" `attr` INT DEFAULT '0'," +
		" `achievepoint` INT DEFAULT '0'," +
		" `hp` INT DEFAULT '0'," +
		" `exp` INT DEFAULT '0'," +
		" `vip_point` INT DEFAULT '0'," +
		" `honor` INT DEFAULT '0'," +
		" `prop` BIGINT DEFAULT '0'," +
		" `cloth` BIGINT DEFAULT '0'," +
		" `other` TEXT," +
		" `player` TEXT," +
		" PRIMARY KEY (`id`)," +
		" KEY ana1 (`action`, `channel`, `time`)," +
		" KEY (`action`)," +
		" KEY (`user_id`)" +
		" ) ENGINE Myisam;"
}
