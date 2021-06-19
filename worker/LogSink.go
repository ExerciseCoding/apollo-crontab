package worker
import(
	"context"
	"crontab/common"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	_"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)
//mongodb存储日志
type LogSink struct {
	client *mongo.Client
	logCollection *mongo.Collection
	logChan chan *common.JobLog
	autoCommitChan chan *common.LogBatch
}

var(
	//单例
	G_logSink *LogSink
)

//批量写入日志
func (LogSink *LogSink) saveLogs(batch *common.LogBatch){
	LogSink.logCollection.InsertMany(context.TODO(),batch.Logs)
}
func (logSink *LogSink) writeLoop(){
	var(
		log *common.JobLog
		logBatch *common.LogBatch //当前批次
		commitTimer *time.Timer
		timeoutBatch *common.LogBatch //超时批次
	)
	for{
		select {
		case log = <- logSink.logChan:
			if logBatch == nil{
				logBatch = &common.LogBatch{}
				//让这个批次超时自动提交(给1s的时间)
				commitTimer = time.AfterFunc(time.Duration(G_config.JobLogCommitTimeout)*time.Millisecond,
					//func(){
					//	//发出超时通知，不要直接提交batch
					//	logSink.autoCommitChan <- logBatch
					//},
					func(batch *common.LogBatch)func(){
						return func() {
							logSink.autoCommitChan <- batch
						}
					}(logBatch),
					)

			}
			//把新的日志追加到批次
			logBatch.Logs = append(logBatch.Logs,log)

			//如果批次满了立即发送
			if len(logBatch.Logs) >= G_config.JobLogBatchSize{
				//发送日志
				logSink.saveLogs(logBatch)
				//清空logBatch
				logBatch = nil

				//取消定时器
				commitTimer.Stop()
			}
		case timeoutBatch = <- logSink.autoCommitChan: //过期的批次
			//判断过期批次是否是当前的批次
			if timeoutBatch != logBatch{
				continue //跳过已经被提交的批次
			}
			//把批次写入到mongodb
			logSink.saveLogs(timeoutBatch)
			//清空logBatch
			logBatch = nil
		}
	}
}
func InitLogSink()(err error){
	var(
		client *mongo.Client
		clientOp *options.ClientOptions
		duration time.Duration
	)
	clientOp = options.Client().ApplyURI(G_config.MongodbUri)

	//超时时间
	duration = time.Duration(G_config.MongodbConnectTimeout) * time.Millisecond
	clientOp.ConnectTimeout = &duration

	//建立mongodb连接
	if client ,err = mongo.Connect(context.TODO(),clientOp); err != nil{
		return
	}

	//选择db和collection
	G_logSink = &LogSink{
		client:        client,
		logCollection: client.Database("cron").Collection("log"),
		logChan:       make(chan *common.JobLog,G_config.LogSinkChanLen),
		autoCommitChan:make(chan *common.LogBatch,1000),
	}

	//启动一个mongodb处理协程
	go G_logSink.writeLoop()
	return
}

//发送日志
func (logSink *LogSink) Append(jobLog *common.JobLog){
	select {
	case logSink.logChan <- jobLog:
	default:
		//队列满了丢弃
	}

}