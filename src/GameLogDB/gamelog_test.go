package GameLogDB

import (
	"testing"
)

//测试工程
func TestOps(t *testing.T) {

	str := "INFO2015-11-30 10:49:15,677  10.224.32.84  /getmailaward uid: 93  mail_award" +
		" {'other': {'mailid': 11111301}, u'prop': {u'3020006': 1, u'3020011': 1}}"
	message := Message{logData: str}
	dispatchLog(message)
}
