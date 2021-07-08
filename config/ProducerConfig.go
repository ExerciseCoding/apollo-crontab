package config

import "github.com/Shopify/sarama"

var (
	G_producerConf *sarama.Config
)

func InitProducerConfig()(producerConfig *sarama.Config){
	config := sarama.NewConfig()

	config.Producer.RequiredAcks = sarama.WaitForAll //发送完数据需要leader和follow确定
	config.Producer.Partitioner = sarama.NewRandomPartitioner //新选出一个partition
	config.Producer.Return.Successes = true //成功交付的消息将在success channel返回

	G_producerConf = config
	return G_producerConf

}

