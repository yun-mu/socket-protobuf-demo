package model

import (
	"constant"
	"util/newtime"

	"gopkg.in/mgo.v2/bson"
)

type Feedback struct {
	ID         bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	Type       int           `bson:"type" json:"type"`     // 反馈类别
	Status     int           `bson:"status" json:"status"` // 类型：0 未查看，1 已经查看
	UserID     string        `bson:"user_id" json:"user_id"`
	ContactWay string        `bson:"contact_way" json:"contact_way"` // 联系方式
	Content    string        `bson:"content" json:"content"`
	Imgs       []Img         `bson:"imgs" json:"imgs"` // 图片

	IssueTime int64 `bson:"issue_time" json:"issue_time"` // 发布时间
}

func FindFeedbacksByStatus(status, page, perPage int) ([]Feedback, error) {
	query := bson.M{
		"status": status,
	}
	fields := []string{"-issue_time"}
	return findFeedbacks(query, page, perPage, fields...)
}

func CreateFeedback(feedback Feedback) error {
	feedback.ID = bson.NewObjectId()
	feedback.Status = constant.FeedbackUnReadStatus
	feedback.IssueTime = newtime.GetMilliTimestamp()
	return insertFeedback(feedback)
}

func UpdateFeedbackStatusToRead(ids []string) error {
	bsonIDs := make([]bson.ObjectId, len(ids))
	for i, id := range ids {
		bsonIDs[i] = bson.ObjectIdHex(id)
	}
	query := bson.M{
		"_id": bson.M{
			"$in": bsonIDs,
		},
		"status": constant.FeedbackUnReadStatus,
	}
	update := bson.M{
		"$set": bson.M{
			"status": constant.FeedbackReadedStatus,
		},
	}
	_, err := updateAllFeedback(query, update)
	return err
}

/****************************************** message db action ****************************************/

func findFeedbacks(query interface{}, page, perPage int, fields ...string) ([]Feedback, error) {
	sess := globalSess.Copy()
	defer sess.Close()
	feedbackTable := sess.DB(DBNAME).C(FEEDBACKTABLE)

	msgs := []Feedback{}
	err := feedbackTable.Find(query).Sort(fields...).Skip((page - 1) * perPage).Limit(perPage).All(&msgs)
	return msgs, err
}

func insertFeedback(docs ...interface{}) error {
	sess := globalSess.Clone()
	defer sess.Close()
	feedbackTable := sess.DB(DBNAME).C(FEEDBACKTABLE)

	return feedbackTable.Insert(docs...)
}

func updateFeedback(query, update interface{}) error {
	sess := globalSess.Clone()
	defer sess.Close()
	feedbackTable := sess.DB(DBNAME).C(FEEDBACKTABLE)

	return feedbackTable.Update(query, update)
}

func updateAllFeedback(query, update interface{}) (interface{}, error) {
	sess := globalSess.Clone()
	defer sess.Close()
	feedbackTable := sess.DB(DBNAME).C(FEEDBACKTABLE)

	return feedbackTable.UpdateAll(query, update)
}
