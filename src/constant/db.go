package constant

const (

	/****************************************** mongo ****************************************/

	/****************************************** feedback ****************************************/
	FeedbackUnReadStatus = 0
	FeedbackReadedStatus = 1

	/****************************************** redis ****************************************/
	RedisUserLoc      = "user:loc:%s"    // format: user:loc:<id>
	RedisUserLocValue = "%f,%f"          // value: 使用逗号隔开 type: string
	RedisUserEQ       = "user:eq:%s"     // format: user:eq:<id>
	RedisUserStatus   = "user:status:%s" // format: user:status:<id>, 0表示下线，1表示在线

	RedisUserMsgQuene = "user:msg:quene:%s" // format: user:msg:quene:<id>，消息队列

	/****************************************** user ****************************************/
	UserStatusUnActive = 0
	UserStatusActive   = 1
)
