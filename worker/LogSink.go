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
}

var(
	//单例
	G_logSink *LogSink
)
func (logSink *LogSink) writeLoop(){
	var(
		log *common.JobLog
	)
	for{
		select {
		case log = <- logSink.logChan:
			//把log写入mongodb中
			//logSink.logCollection.InsertOne
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
	}

	//启动一个mongodb处理协程
	go G_logSink.writeLoop()
	return
}