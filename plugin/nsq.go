package plugin

import (
	"github.com/nsqio/go-nsq"
    "log"
)

type NsqProducer struct {
	counter   uint64
	producers map[string]*nsq.Producer
	nsqdAddrs []string
}

func NewNsqProducer(NsqWriters []string) (*NsqProducer, error) {
	p := &NsqProducer{nsqdAddrs: NsqWriters}

	cfg := nsq.NewConfig()
	p.producers = make(map[string]*nsq.Producer)
	for _, addr := range p.nsqdAddrs {
		producer, err := nsq.NewProducer(addr, cfg)
		if err != nil {
			log.Fatalf("Connect to nsq host err")
			return nil, err
		}
		p.producers[addr] = producer
	}

	return p, nil

}

func (p *NsqProducer) SendNsqMsg(topic string, body []byte) error {
	idx := p.counter % uint64(len(p.nsqdAddrs))
	producer := p.producers[p.nsqdAddrs[idx]]
	err := producer.Publish(topic, body)
	p.counter++

	return err
}

func (p *NsqProducer) MultiSendNsqMsg(topic string, body [][]byte) error {
	idx := p.counter % uint64(len(p.nsqdAddrs))
	producer := p.producers[p.nsqdAddrs[idx]]
	err := producer.MultiPublish(topic, body)
	p.counter++

	return err
}

//func (p *NsqProducer) ProducersClose() {
//	for _, producer := range p.producers {
//		producer.Stop()
//	}
//}

