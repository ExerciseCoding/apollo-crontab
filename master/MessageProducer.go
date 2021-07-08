package master

import (
	"github.com/Shopify/sarama"
)
type Producer struct {
	client *sarama.Client
}

//func InitProducer(){
//
//}