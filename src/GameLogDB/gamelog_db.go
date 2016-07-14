package GameLogDB

import (
	"MyUtility"
	"database/sql"
	// "encoding/json"
	"fmt"
	"log"
	// "reflect"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	// "github.com/novalagung/golpal"
)

type MySqlDb struct {
	// dbInstance *sql.DB
	log_fields []string
}

var mysqldb *MySqlDb

func init() {
	mysqldb = new(MySqlDb)
	mysqldb.log_fields = []string{
		"id", "user_id", "channel", "action", "time", "diamond",
		"coin", "attr", "achievepoint", "hp", "exp", "vip_point",
		"honor", "prop", "cloth", "other", "player"}

	//creat table
	mysqldb.ExecDB(GetSqlCommand("create_log_login_detail"))
	mysqldb.ExecDB(GetSqlCommand("create_login_log"))
	mysqldb.ExecDB(GetSqlCommand("create_retention_log"))
}

func GetGameLogDB() *MySqlDb {
	return mysqldb
}

func recharge_msg(records []string) {

	dmsg, _ := MyUtility.DealJsonFormat(records[6])
	day_time := strings.Split(strings.TrimSpace(records[1]), ",")[0]
	username := MyUtility.GetMapRetString(dmsg, "username", "")
	channel := strings.TrimSpace(records[7])
	total := MyUtility.GetMapRetString(dmsg, "total", "0")
	level := MyUtility.GetMapRetString(dmsg, "level", "0")
	viplevel := MyUtility.GetMapRetString(dmsg, "vip_level", "0")
	num := MyUtility.GetMapRetString(dmsg, "num", "0")

	//把新的总充值额和新的单日充值额存入数据库
	detail_sql := fmt.Sprintf("INSERT INTO `recharge_details` (`username`,`channel`,`level`,`viplevel`,`num`,`day_time`)"+
		"VALUE ('%s', '%s', %s, %s, %s, '%s')", username, channel, level, viplevel, num, day_time)

	total_sql := fmt.Sprintf("REPLACE INTO `recharge_total` (`username`,`channel`,`total_num`) VALUE ('%s', '%s', '%s')", username, channel, total)

	mysqldb.ExecDB(detail_sql)
	mysqldb.ExecDB(total_sql)

}

func login_msg(records []string) {

	user_id := strings.Split(strings.TrimSpace(records[4]), ":")[1]
	dmsg, _ := MyUtility.DealJsonFormat(records[6])

	log.Println("user_id:", user_id)
	log.Println("dmsg:", dmsg)

	last_time := time.Unix(MyUtility.GetMapRetInt(dmsg, "last", 0), 0)
	first_time := time.Unix(MyUtility.GetMapRetInt(dmsg, "first", 0), 0)

	last := last_time.Format("2006-01-02 15:04:05")
	first := first_time.Format("2006-01-02 15:04:05")

	last_day := last_time.Format("2006-01-02")
	first_day := first_time.Format("2006-01-02")

	level := MyUtility.GetMapRetString(dmsg, "level", "0")
	viplevel := MyUtility.GetMapRetString(dmsg, "vip_level", "0")
	channel := strings.TrimSpace(records[7])

	sql := fmt.Sprintf("INSERT INTO log_login_detail VALUE ('%s', '%s', '%s', %s, %s, '%s', '%s')",
		user_id, channel, MyUtility.GetMapRetString(dmsg, "dbid", ""), level, viplevel, last, first)

	mysqldb.ExecDB(sql)

	sql = "SELECT lastdate, firstdate FROM log_users_login where username ='" +
		MyUtility.GetMapRetString(dmsg, "username", "") + "'"

	row := mysqldb.QueryDB(sql)
	rowFlag := false
	var rowValue_lastdate string = ""
	var rowValue_firstdate string = ""
	if row != nil {
		for row.Next() {
			row.Scan(&rowValue_lastdate, &rowValue_firstdate)
			log.Println("rowValue_lastdate:", rowValue_lastdate)
			log.Println("rowValue_firstdate:", rowValue_firstdate)
			rowFlag = strings.EqualFold(rowValue_lastdate, last_day)
		}
	}

	if row == nil || rowFlag == false {
		sql = fmt.Sprintf("INSERT INTO log_users_login(username , cnt, lastdate, firstdate) SELECT '"+
			MyUtility.GetMapRetString(dmsg, "username", "")+"',1, '%s', '%s'", last_day, first_day)
		sql += "ON DUPLICATE KEY UPDATE cnt=1, lastdate=date(now())"

		mysqldb.ExecDB(sql)
		daycnt := last_time.Day() - first_time.Day()
		log.Println("daycnt:", daycnt)
		sql = "INSERT INTO log_retention_users(daytime,dbid,cid,first,second," +
			"third,fourth, fifth, sixth,week,halfmonth," +
			"month,interval_second_login)"

		if daycnt == 0 {
			sql += fmt.Sprintf("SELECT '%s','%s','%s',1,0,0,0,0,0,0,0,0,0 ON DUPLICATE KEY UPDATE first=first+1", last_day, MyUtility.GetMapRetString(dmsg, "dbid", ""), channel)
		} else if daycnt == 1 {
			sql += fmt.Sprintf("SELECT '%s','%s','%s',0,1,0,0,0,0,0,0,0,0 ON DUPLICATE KEY UPDATE second=second+1", last_day, MyUtility.GetMapRetString(dmsg, "dbid", ""), channel)
		} else if daycnt == 2 {
			sql += fmt.Sprintf("SELECT '%s','%s','%s',0,0,1,0,0,0,0,0,0,0 ON DUPLICATE KEY UPDATE third=third+1", last_day, MyUtility.GetMapRetString(dmsg, "dbid", ""), channel)
		} else if daycnt == 3 {
			sql += fmt.Sprintf("SELECT '%s','%s','%s',0,0,0,1,0,0,0,0,0,0 ON DUPLICATE KEY UPDATE fourth=fourth+1", last_day, MyUtility.GetMapRetString(dmsg, "dbid", ""), channel)
		} else if daycnt == 4 {
			sql += fmt.Sprintf("SELECT '%s','%s','%s',0,0,0,0,1,0,0,0,0,0 ON DUPLICATE KEY UPDATE fifth=fifth+1", last_day, MyUtility.GetMapRetString(dmsg, "dbid", ""), channel)
		} else if daycnt == 5 {
			sql += fmt.Sprintf("SELECT '%s','%s','%s',0,0,0,0,0,1,0,0,0,0 ON DUPLICATE KEY UPDATE sixth=sixth+1", last_day, MyUtility.GetMapRetString(dmsg, "dbid", ""), channel)
		} else if daycnt == 6 {
			sql += fmt.Sprintf("SELECT '%s','%s','%s',0,0,0,0,0,0,1,0,0,0 ON DUPLICATE KEY UPDATE week=week+1", last_day, MyUtility.GetMapRetString(dmsg, "dbid", ""), channel)
		} else if daycnt < 15 {
			sql += fmt.Sprintf("SELECT '%s','%s','%s',0,0,0,0,0,0,0,1,0,0 ON DUPLICATE KEY UPDATE halfmonth=halfmonth+1", last_day, MyUtility.GetMapRetString(dmsg, "dbid", ""), channel)
		} else if daycnt < 30 {
			sql += fmt.Sprintf("SELECT '%s','%s','%s',0,0,0,0,0,0,0,0,1,0 ON DUPLICATE KEY UPDATE month=month+1", last_day, MyUtility.GetMapRetString(dmsg, "dbid", ""), channel)
		} else {
			sql = ""
		}

		if len(sql) > 0 {
			mysqldb.QueryDB(sql)
		}

		if len(rowValue_firstdate) > 0 && strings.EqualFold(rowValue_firstdate, rowValue_lastdate) {
			sql = fmt.Sprintf("INSERT INTO log_retention_users(daytime,dbid,cid,"+
				"first,second,third,fourth, fifth, sixth, week"+
				",halfmonth,month,interval_second_login)"+
				"SELECT '%s','%s','%s',0,0,0,0,0,0,0,0,0,1 "+
				"ON DUPLICATE KEY UPDATE interval_second_login=interval_second_login+1", last_day,
				MyUtility.GetMapRetString(dmsg, "dbid", ""), channel)
			mysqldb.ExecDB(sql)
		}

	} else {
		sql = "UPDATE log_retention_users SET cnt = cnt + 1 WHERE username='" + MyUtility.GetMapRetString(dmsg, "username", "") + "'"
		mysqldb.ExecDB(sql)
	}

	// log.Println("lastValue:", lastValue, ",last:", last, ",last_day:", last_day)
	// log.Println("firstValue:", firstValue, ",first:", first, ",first_day:", first_day)

}

func log_save(records []string) {
	action := strings.TrimSpace(records[5])
	if strings.EqualFold(action, "cust_event") {
		return
	}

	data, _ := MyUtility.DealJsonFormat(records[6])

	day_time := strings.Split(strings.TrimSpace(records[1]), ",")[0]
	user_id := strings.Split(strings.TrimSpace(records[4]), ":")[1]
	var channel string = ""
	if len(records) > 8 {
		channel = strings.TrimSpace(records[7])
	} else {
		channel = MyUtility.GetMapRetString(data, "channel", "unknow")
	}

	if len(user_id) <= 0 {
		user_id = "0"
	}
	// date := strings.Replace(day_time, "-", "", -1)

	//开始数据分析啦

	//属性
	attr_num := map[string]int{"health": 5, "virtue": 6, "vigour": 7}
	attrs := make([]string, 0)
	attrs_len := 0
	tempattr, attrErr := MyUtility.GetMapWithDefault(data, "attr", make(map[string]interface{}, 0)).(map[string]interface{})
	if attrErr == false {
		delete(data, "attr")
	} else {
		maptempattr_len := len(tempattr)
		if maptempattr_len > 0 {
			for k, v := range tempattr {
				// log.Println("k:", k)
				// log.Println("v:", v)
				if len(strings.TrimSpace(k)) > 0 {
					target, attr_finded := attr_num[k]
					if attr_finded {
						attrs = append(attrs, strconv.Itoa(MyUtility.GetSum(target*1000, v)))
					}
				}
			}

			attrs_len = len(attrs)
			if attrs_len > 0 {
				data["attr"] = attrs[attrs_len-1]
				attrs = attrs[:attrs_len-1]
			}
		}
	}

	//服装
	cloth := make([]string, 0)
	cloth_len := 0
	tempCloth, clothErr := MyUtility.GetMapWithDefault(data, "cloth", make(map[string]interface{}, 0)).(map[string]interface{})
	if clothErr == false {
		delete(data, "cloth")
	} else {
		mapcloth_len := len(tempCloth)
		if mapcloth_len > 0 {
			// log.Println("cloth_len:", cloth_len)
			for k, v := range tempCloth {
				// log.Println("k:", k)
				// log.Println("v:", v)
				num := int(MyUtility.GetMapRetInt(v.(map[string]interface{}), "num", 0))
				// log.Println("num:", num)
				cloth_id, err := strconv.Atoi(k)
				if err == nil {
					cloth = append(cloth, strconv.Itoa(cloth_id*1000+num))
				} else {
					log.Println("strconv.Atoi err:", err)
				}
			}

			cloth_len = len(cloth)
			if cloth_len > 0 {
				data["cloth"] = cloth[cloth_len-1]
				cloth = cloth[:cloth_len-1]
			}
		}
	}

	//道具
	prop := make([]string, 0)
	prop_len := 0
	tempProp, propErr := MyUtility.GetMapWithDefault(data, "prop", make(map[string]interface{}, 0)).(map[string]interface{})
	if propErr == false {
		delete(data, "prop")
	} else {
		maptempProp_len := len(tempProp)
		if maptempProp_len > 0 {
			for k, v := range tempProp {
				// log.Println("k:", k)
				// log.Println("v:", v)
				if len(strings.TrimSpace(k)) > 0 {
					prop = append(prop, strconv.Itoa(MyUtility.GetSum(MyUtility.GetString2Int(k)*1000, v)))
					prop_len += 1
				}
			}

			if prop_len > 0 {
				data["prop"] = prop[prop_len-1]
				prop = prop[:prop_len-1]
			}
		}
	}

	// log.Println("before filter data:", data)
	//删除非log表字段数据
	unfields := make([]string, 0)
	for k, _ := range data {
		if MyUtility.ContainInTarget(k, mysqldb.log_fields) {
			// log.Println(k, "is in log_fields")
		} else {
			unfields = append(unfields, k)
		}
	}
	for _, v := range unfields {
		// log.Println("delete ", v, " in data")
		delete(data, v)
	}

	//sql语句拼接
	pre_values := []string{strings.TrimSpace(user_id), strings.TrimSpace(action), strings.TrimSpace(day_time), strings.TrimSpace(channel)}
	pre_columns := []string{"user_id", "action", "time", "channel"}

	mapKeyList := MyUtility.GetAllMapKey(data)
	colums := strings.Join(append(pre_columns, mapKeyList...), ",")
	values := strings.Join(append(pre_values, MyUtility.GetAllMapValueOrderByKey(data, mapKeyList)...), "','")

	if len(colums) > 0 {
		prepare := fmt.Sprintf("INSERT INTO log_%s (%s)VALUES ('%s')", "dialy", colums, values)

		mysqldb.ExecDB(prepare)
	}

	// log.Println("after filter data:", data)
	// log.Println("day_time:", day_time)
	// log.Println("date:", date)
	// log.Println("user_id:", user_id)
	// log.Println("channel:", channel)
	// log.Println("cloth:", cloth)
	// log.Println("prop:", prop)
	// log.Println("colums:", colums)
	// log.Println("values:", values)

	max_Len := MyUtility.GetIntMax(attrs_len, cloth_len, prop_len)

	if max_Len > 1 {
		// log.Println("|||||||||||||---------------------||||||||||||")
		rest_prop := make([]string, max_Len-prop_len)
		rest_cloth := make([]string, max_Len-cloth_len)
		rest_attrs := make([]string, max_Len-attrs_len)

		prop = append(prop, rest_prop...)
		cloth = append(cloth, rest_cloth...)
		attrs = append(attrs, rest_attrs...)

		// log.Println("prop:", prop)
		// log.Println("cloth:", cloth)
		// log.Println("attrs:", attrs)

		rest_values := make([]string, max_Len-1)
		for i := 0; i < max_Len-1; i++ {
			finalAttrs := attrs[i]
			finalProp := prop[i]
			finalCloth := cloth[i]
			if len(finalAttrs) <= 0 {
				finalAttrs = "0"
			}
			if len(finalProp) <= 0 {
				finalProp = "0"
			}
			if len(finalCloth) <= 0 {
				finalCloth = "0"
			}
			tmp := fmt.Sprintf("%s','%s','%s", finalAttrs, finalProp, finalCloth)
			// log.Println("tmp[a]:", i, "=", tmp)
			rest_values[i] = tmp
		}
		// log.Println("rest_values:", rest_values)
		rest_columns := append(pre_columns, "attr", "prop", "cloth", "other")
		values_sql := "'"

		for k, v := range rest_values {
			tmp := strings.Join(append(pre_values, v), "','") + "','" + MyUtility.GetMapRetString(data, "other", "")
			// log.Println("k:", k, ",v:", v)
			// log.Println("tmp[b]:", tmp)
			if k == len(rest_values)-1 {
				values_sql += tmp + "'"
			} else {
				values_sql += tmp + "'), ('"
			}
		}
		rest_sql := fmt.Sprintf("INSERT INTO log_%s (%s)VALUES (%s)", "dialy", strings.Join(rest_columns, ","), values_sql)

		mysqldb.ExecDB(rest_sql)
	}

	log.Println("----------------------------------------------------")

}

func (self *MySqlDb) connectDB() *sql.DB {
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/gamelog_t2?charset=utf8")
	checkError(err)
	// self.dbInstance = db
	return db
}

func (self *MySqlDb) QueryDB(prepare string) *sql.Rows {

	dbInstance := self.connectDB()

	if dbInstance == nil {
		return nil
	}

	defer func() {
		dbInstance.Close()
		dbInstance = nil
	}()

	log.Println("prepare:", prepare)
	res, err := dbInstance.Query(prepare)
	checkError(err)
	return res
	// id, err := res.LastInsertId()
	// log.Println("QueryDB LastInsertId:", id)
	// checkError(err)

}

func (self *MySqlDb) ExecDB(prepare string) sql.Result {

	dbInstance := self.connectDB()

	if dbInstance == nil {
		return nil
	}

	defer func() {
		dbInstance.Close()
		dbInstance = nil
	}()

	log.Println("prepare:", prepare)
	res, err := dbInstance.Exec(prepare)
	checkError(err)
	return res
	// id, err := res.LastInsertId()
	// log.Println("QueryDB LastInsertId:", id)
	// checkError(err)

}

func checkError(err error) {
	if err != nil {
		log.Println("err:", err)
		panic(err)
	}
}
