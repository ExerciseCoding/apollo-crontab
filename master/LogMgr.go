package master

import (
	"context"
	"crontab/common"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type LogMgr struct {
	client *mongo.Client
	logCollection *mongo.Collection

}
var(
	G_logMgr *LogMgr
)
func InitLogMgr()(err error) {
	var (
		client   *mongo.Client
		clientOp *options.ClientOptions
		duration time.Duration
	)
	clientOp = options.Client().ApplyURI(G_config.MongodbUri)

	//超时时间
	duration = time.Duration(G_config.MongodbConnectTimeout) * time.Millisecond
	clientOp.ConnectTimeout = &duration

	//建立mongodb连接
	if client, err = mongo.Connect(context.TODO(), clientOp); err != nil {
		return
	}

	G_logMgr = &LogMgr{
		client:        client,
		logCollection: client.Database("cron").Collection("log"),
	}
	return
}

//查看任务日志
func (logMgr *LogMgr) ListLog(name string,skip int,limit int)(logArr []*common.JobLog,err error){
	var(
		filter *common.JobLogFilter
		logSort *common.SortLogByStartTime
		findOptions *options.FindOptions
		cursor *mongo.Cursor
		skipNew int64
		limitNew int64
		jobLog *common.JobLog

	)
	logArr = make([]*common.JobLog,0)
	skipNew , limitNew= int64(skip), int64(limit)
	//过滤条件
	filter = &common.JobLogFilter{JobName:name}
	//按照任务开始时间倒排
	logSort = &common.SortLogByStartTime{SortOrder:-1}
	findOptions = options.Find()
	findOptions.Sort = logSort
	findOptions.Skip = &skipNew
	findOptions.Limit = &limitNew
	//查询
	if cursor,err = logMgr.logCollection.Find(context.TODO(),filter,findOptions); err != nil{
		return
	}
	//延迟释放游标
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()){
		jobLog = &common.JobLog{}

		//反序列化BSON
		if err = cursor.Decode(jobLog); err != nil{
			continue //有日志不合法
		}

		logArr = append(logArr,jobLog)
	}
	return
}