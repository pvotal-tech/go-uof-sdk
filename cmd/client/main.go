package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/pvotal-tech/go-uof-sdk"
	"github.com/pvotal-tech/go-uof-sdk/api"
	"github.com/pvotal-tech/go-uof-sdk/sdk"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

const (
	EnvBookmakerID = "UOF_BOOKMAKER_ID"
	EnvToken       = "UOF_TOKEN"
)

func env(name string) string {
	e, ok := os.LookupEnv(name)
	if !ok {
		log.Printf("env %s not found", name)
	}
	return e
}

func debugHTTP() {
	if err := http.ListenAndServe("localhost:8124", nil); err != nil {
		log.Fatal(err)
	}
}

func exitSignal() context.Context {
	ctx, stop := context.WithCancel(context.Background())
	go func() {
		c := make(chan os.Signal, 1)
		//SIGINT je ctrl-C u shell-u, SIGTERM salje upstart kada se napravi sudo stop ...
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		<-c
		stop()
	}()
	return ctx
}

func main() {
	go func() {
		_ = http.ListenAndServe(fmt.Sprintf(":%d", 6060), nil)
	}()
	go debugHTTP()

	//preloadTo := time.Now().Add(1 * time.Hour)

	timestamp := uof.CurrentTimestamp() - 1000 // -5 minutes
	var ppp uof.ProducersChange
	ppp.Add(uof.ProducerPrematch, timestamp)
	//pc.Add(uof.ProducerLiveOdds, timestamp)
	//pc.Add(uof.ProducerBetPal, timestamp)
	//pc.Add(uof.ProducerPremiumCricket, timestamp)

	err := sdk.Run(exitSignal(),
		sdk.Credentials(38616, "j8l0CpSoytxhp7xb7r", 123456),
		sdk.Staging(),
		//sdk.Recovery(ppp),
		sdk.DisablePipeline(),
		sdk.ConfigThrottle(false, 10),
		//sdk.ConfigConcurrentAPIFetch(true),
		sdk.DisableAutoAck(),
		//sdk.Fixtures(preloadTo),
		sdk.Languages(uof.Languages("en")),
		//sdk.BufferedConsumer(pipe.FileStore("./tmp"), 1024),
		sdk.Consumer(logMessages),
		sdk.ListenErrors(listenSDKErrors),
	)
	if err != nil {
		fmt.Println("SDK FAILED!!", err.Error())
		log.Fatal(err)
	}
}

func startReplay(rpl *api.ReplayAPI) error {
	return nil
	//return rpl.StartEvent(eventURN, speed, maxDelay)
	//return rpl.StartScenario(scenarioID, speed, maxDelay)
}

// consumer of incoming messages
func logMessages(in <-chan *uof.Message) error {
	for m := range in {
		err := processMessage(m)
		if err != nil {
			fmt.Println("error during message process, nacking and requesting requeue", err.Error())

			if !m.EnabledAutoAck {
				if err := m.NackRequeue(); err != nil {
					fmt.Println("error during NackRequeue", err.Error())
				}
			}
			continue
		}
		if !m.EnabledAutoAck {
			if m.Delivery != nil {
				fmt.Println("Acking", m.Delivery.DeliveryTag)
			}
			if err := m.Ack(); err != nil {
				fmt.Println("error during Ack", err.Error())
				continue
			}
		}
	}
	return nil
}

var sportMap = make(map[string]string)

var errorCount int

var pc uof.ProducersChange

func processMessage(m *uof.Message) error {
	time.Sleep(time.Second * 2000)
	reqID := 0
	fmt.Printf("%-25s | %-25v | %-25s | %-25d\n", m.Type, m.Delivery != nil, m.EventURN.String(), m.SportID)
	return nil
	switch m.Type {
	case uof.MessageTypeConnection:
		fmt.Printf("%-25s status: %s, server: %s, local: %s, network: %s, tls: %s\n", m.Type, m.Connection.Status, m.Connection.ServerName, m.Connection.LocalAddr, m.Connection.Network, m.Connection.TLSVersionToString())
	case uof.MessageTypeFixture:
		//fmt.Printf("%-25s lang: %s, urn: %s raw: %d, timestamps %d\n", m.Type, m.Lang, m.FixtureBytes.URN, len(m.Raw), m.Timestamp)
	case uof.MessageTypeMarkets:
		//fmt.Printf("%-25s lang: %s, count: %d\n", m.Type, m.Lang, len(m.Markets))
	case uof.MessageTypePlayer:
		//fmt.Printf("%-25s lang: %s, count: %d\n", m.Type, m.Lang, len(m.Markets))
	case uof.MessageTypeAlive:
		fmt.Printf("%-25s producer: %s, timestamp: %d, subscribed: %d, nodeID: %d\n", m.Type, m.Alive.Producer, m.Alive.Timestamp, m.Alive.Subscribed, m.NodeID)
	case uof.MessageTypeSnapshotComplete:
		fmt.Printf("%-25s prod: %s, reqID: %d, nodeID: %d\n", m.Type, m.SnapshotComplete.Producer, m.SnapshotComplete.RequestID, m.NodeID)
	case uof.MessageTypeProducersChange:
		statuses := make([]string, len(m.Producers))
		for i, v := range m.Producers {
			status := ""
			switch v.Status {
			case -1:
				status = "down"
			case 1:
				status = "active"
			case 2:
				status = "in-recovery"
			}
			statuses[i] = v.Producer.Code() + "=" + status
		}
		fmt.Printf("%-25s producers: %s \n", m.Type, strings.Join(statuses, "|"))
	case uof.MessageTypeBetSettlement:
		//if !validEventURN(m.BetSettlement.EventURN) { return nil }
		if validRequestID(m.BetSettlement.RequestID) {
			reqID = *m.BetSettlement.RequestID
		} //return nil }
		fmt.Println(m.Type, m.Delivery.DeliveryTag, m.NodeID, reqID)
		//logObject(*m.BetSettlement)
	case uof.MessageTypeRollbackBetSettlement:
		//if !validEventURN(m.RollbackBetSettlement.EventURN) { return nil }
		if validRequestID(m.RollbackBetSettlement.RequestID) {
			reqID = *m.RollbackBetSettlement.RequestID
		} //return nil }
		fmt.Println(m.Type, m.Delivery.DeliveryTag, m.NodeID, reqID)
		//logObject(*m.RollbackBetSettlement)
	case uof.MessageTypeBetStop:
		//if !validEventURN(m.BetStop.EventURN) { return nil }
		if validRequestID(m.BetStop.RequestID) {
			reqID = *m.BetStop.RequestID
		} //return nil }
		fmt.Println(m.Type, m.Delivery.DeliveryTag, m.NodeID, reqID)
		//logObject(*m.BetStop)
	case uof.MessageTypeBetCancel:
		//if !validEventURN(m.BetCancel.EventURN) { return nil }
		if validRequestID(m.BetCancel.RequestID) {
			reqID = *m.BetCancel.RequestID
		} //return nil }
		fmt.Println(m.Type, m.Delivery.DeliveryTag, m.NodeID, reqID)
		//logObject(*m.BetCancel)
	case uof.MessageTypeRollbackBetCancel:
		//if !validEventURN(m.RollbackBetCancel.EventURN) { return nil }
		if validRequestID(m.RollbackBetCancel.RequestID) {
			reqID = *m.RollbackBetCancel.RequestID
		} //return nil }
		fmt.Println(m.Type, m.Delivery.DeliveryTag, m.NodeID, reqID)
		//logObject(*m.RollbackBetCancel)
	case uof.MessageTypeOddsChange:
		//if !validEventURN(m.OddsChange.EventURN) { return nil }
		if validRequestID(m.OddsChange.RequestID) {
			reqID = *m.OddsChange.RequestID
		} //return nil }
		fmt.Println(m.Type, m.Delivery.DeliveryTag, m.NodeID, reqID)
		//logObject(*m.OddsChange)
	}
	return nil
}

func validRequestID(requestID *int) bool {
	if requestID == nil {
		return false
	}
	return true //*requestID == 999999
}

func validEventURN(urn uof.URN) bool {
	return urn.ID() == 41763059
}

// listenSDKErrors listens all SDK errors for logging or any other pourpose
func listenSDKErrors(err error) {
	// example handling SDK typed errors
	var eu uof.Error
	if errors.As(err, &eu) {
		// use uof.Error attributes to build custom logging
		var logLine string
		if eu.Severity == uof.NoticeSeverity {
			logLine = fmt.Sprintf("NOTICE Operation:%s Details:", eu.Op)
		} else {
			logLine = fmt.Sprintf("ERROR Operation:%s Details:", eu.Op)
		}

		if eu.Inner != nil {
			var ea uof.APIError
			if errors.As(eu.Inner, &ea) {
				// use uof.APIError attributes for custom logging
				logLine = fmt.Sprintf("%s URL:%s", logLine, ea.URL)
				logLine = fmt.Sprintf("%s StatusCode:%d", logLine, ea.StatusCode)
				logLine = fmt.Sprintf("%s Response:%s", logLine, ea.Response)
				if ea.Inner != nil {
					logLine = fmt.Sprintf("%s Inner:%s", logLine, ea.Inner)
				}

				// or just log error as is...
				//log.Print(ea.Error())
			} else {
				// not an uof.APIError
				logLine = fmt.Sprintf("%s %s", logLine, eu.Inner)
			}
		}
		log.Println(logLine)

		// or just log error as is...
		//log.Println(eu.Error())
	} else {
		// any other error not uof.Error
		log.Println(err)
	}
}

func logObject(obj interface{}) {
	bytes, err := json.MarshalIndent(obj, "", " ")
	if err != nil {
		fmt.Println("ERROR WHILE UNMARSHALLING")
	}
	fmt.Println(string(bytes))
}

func getProducerChange(producers uof.ProducersChange, producer uof.Producer) (*uof.ProducerChange, error) {
	for _, p := range producers {
		if p.Producer.Code() == producer.Code() {
			return &p, nil
		}
	}
	return nil, fmt.Errorf("invalid producer: %s (forgot to do recovery?)", producer.Code())
}
