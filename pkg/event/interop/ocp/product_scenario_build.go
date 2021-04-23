package ocp

import (
	"sync"

	"github.com/adrianriobo/qe-eventmanager/pkg/logging"
	"github.com/adrianriobo/qe-eventmanager/pkg/messaging"
	stomp "github.com/go-stomp/stomp/v3"
)

const (
	topic string = "VirtualTopic.qe.ci.product-scenario.crcqe.test4"
)

type ProductScenarioBuild struct {
	connection   messaging.UMBConnection
	subscription *stomp.Subscription
	consumers    *sync.WaitGroup
	handlers     *sync.WaitGroup
}

func New(connection *messaging.UMBConnection) *ProductScenarioBuild {
	return &ProductScenarioBuild{connection: *connection}
}

func (p *ProductScenarioBuild) Init() {
	// this should be change to ack msg to ack when handle is done to avoid data loss
	p.subscription, _ = p.connection.FailoverSubscribe("Consumer.psi-crcqe-openstack.test5."+topic, stomp.AckClientIndividual)
	// group of consumers
	p.consumers = &sync.WaitGroup{}
	// group of handlers
	p.handlers = &sync.WaitGroup{}
	// async consume
	p.consumers.Add(1)
	go p.consume()
}

func (p *ProductScenarioBuild) Finish() {
	if err := p.subscription.Unsubscribe(); err != nil {
		logging.Error(err)
		// Force consume as finished ?
		p.consumers.Done()
	}
	p.consumers.Wait()
	p.handlers.Wait()
}

// TODO add selector based on regex??
func (p *ProductScenarioBuild) consume() {
	defer p.consumers.Done()
	for p.subscription.Active() {
		msg, err := p.subscription.Read()
		if err != nil {
			if !p.subscription.Active() {
				break
			}
			logging.Errorf("Error reading from topic: %s. %s", topic, err)
			break
		}
		p.handlers.Add(1)
		go p.handle(msg)
	}
}

func (p *ProductScenarioBuild) handle(msg *stomp.Message) {
	// when finish remove from group
	defer p.handlers.Done()
	// heavy consuming may regex over string
	interopMsg, err := Unmarshal(msg.Body)
	if err != nil {
		logging.Error("Error unmarshalling")
	}
	logging.Debugf("Start PipelineRun based on %s", (*interopMsg)[0].Topic)
	// Create N pipelineruns -> Evolve into CRD from PSI-Operator
	// On each pipelirun launch a pulling for results -> Evolve on listen back to some async mechanism generated by PSI-Operator
	// Send back message to interop
	// this shows with bad format
	if err := msg.Conn.Ack(msg); err != nil {
		logging.Error("Error ack message")
		// as we can not ack we should not process so
		// we need to run some compensation
		// and retry to process it
	}
}
