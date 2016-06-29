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
