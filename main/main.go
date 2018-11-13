package main

import (
	"log"
	"os"
	"time"

	"github.com/TerrexTech/agg-inventory-cmd/inventory"
	"github.com/TerrexTech/go-commonutils/commonutil"
	"github.com/TerrexTech/go-eventspoll/poll"
	"github.com/coreos/etcd/clientv3"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

func validateEnv() error {
	missingVar, err := commonutil.ValidateEnv(
		"ETCD_HOSTS",
		"KAFKA_BROKERS",

		"KAFKA_CONSUMER_EVENT_GROUP",
		"KAFKA_CONSUMER_EVENT_QUERY_GROUP",

		"KAFKA_CONSUMER_EVENT_TOPIC",
		"KAFKA_CONSUMER_EVENT_QUERY_TOPIC",
		"KAFKA_PRODUCER_EVENT_QUERY_TOPIC",
		"KAFKA_PRODUCER_RESPONSE_TOPIC",

		"MONGO_HOSTS",
		"MONGO_DATABASE",
		"MONGO_AGG_COLLECTION",
		"MONGO_META_COLLECTION",

		"MONGO_CONNECTION_TIMEOUT_MS",
		"MONGO_RESOURCE_TIMEOUT_MS",
	)

	if err != nil {
		err = errors.Wrapf(err, "Env-var %s is required for testing, but is not set", missingVar)
		return err
	}
	return nil
}

func main() {
	log.Println("Reading environment file")
	err := godotenv.Load("./.env")
	if err != nil {
		err = errors.Wrap(err,
			".env file not found, env-vars will be read as set in environment",
		)
		log.Println(err)
	}

	err = validateEnv()
	if err != nil {
		log.Fatalln(err)
	}

	kc, err := loadKafkaConfig()
	if err != nil {
		err = errors.Wrap(err, "Error in KafkaConfig")
		log.Fatalln(err)
	}
	mc, err := loadMongoConfig()
	if err != nil {
		err = errors.Wrap(err, "Error in MongoConfig")
		log.Fatalln(err)
	}
	ioConfig := poll.IOConfig{
		ReadConfig: poll.ReadConfig{
			EnableInsert: true,
			EnableUpdate: true,
			EnableDelete: true,
		},
		KafkaConfig: *kc,
		MongoConfig: *mc,
	}

	eventPoll, err := poll.Init(ioConfig)
	if err != nil {
		err = errors.Wrap(err, "Error creating EventPoll service")
		log.Fatalln(err)
	}

	etcdHostsStr := os.Getenv("ETCD_HOSTS")
	etcdConfig := clientv3.Config{
		DialTimeout: 5 * time.Second,
		Endpoints:   *commonutil.ParseHosts(etcdHostsStr),
	}
	etcdUsername := os.Getenv("ETCD_USERNAME")
	etcdPassword := os.Getenv("ETCD_PASSWORD")
	if etcdUsername != "" {
		etcdConfig.Username = etcdUsername
	}
	if etcdPassword != "" {
		etcdConfig.Password = etcdPassword
	}
	etcd, err := clientv3.New(etcdConfig)
	if err != nil {
		err = errors.Wrap(err, "Failed to connect to ETCD")
		log.Fatalln(err)
	}
	log.Println("ETCD Ready")

	for {
		select {
		case <-eventPoll.RoutinesCtx().Done():
			err = errors.New("service-context closed")
			log.Fatalln(err)

		case eventResp := <-eventPoll.Delete():
			go func(eventResp *poll.EventResponse) {
				if eventResp == nil {
					return
				}
				err := eventResp.Error
				if err != nil {
					err = errors.Wrap(err, "Error in Delete-EventResponse")
					log.Println(err)
					return
				}
				kafkaResp := inventory.Delete(mc.AggCollection, &eventResp.Event)
				if kafkaResp != nil {
					eventPoll.ProduceResult() <- kafkaResp
				}
			}(eventResp)

		case eventResp := <-eventPoll.Insert():
			go func(eventResp *poll.EventResponse) {
				if eventResp == nil {
					return
				}
				err := eventResp.Error
				if err != nil {
					err = errors.Wrap(eventResp.Error, "Error in Insert-EventResponse")
					log.Println(err)
					return
				}
				kafkaResp := inventory.Insert(mc.AggCollection, &eventResp.Event)
				if kafkaResp != nil {
					eventPoll.ProduceResult() <- kafkaResp
				}
			}(eventResp)

		case eventResp := <-eventPoll.Update():
			go func(eventResp *poll.EventResponse) {
				if eventResp == nil {
					return
				}
				err := eventResp.Error
				if err != nil {
					err = errors.Wrap(err, "Error in Update-EventResponse")
					log.Println(err)
					return
				}
				kafkaResp := inventory.Update(etcd, mc.AggCollection, &eventResp.Event)
				if kafkaResp != nil {
					eventPoll.ProduceResult() <- kafkaResp
				}
			}(eventResp)
		}
	}
}
