package cache

import (
	"os"

	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
)

var redisClient *redis.Client

func Get() (*redis.Client, error) {
	if redisClient != nil {
		return redisClient, nil
	}

	address := os.Getenv("REDIS_URL")
	password := os.Getenv("REDIS_PASSWORD")

	rdb := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       0,
	})

	// check if the redis is available
	_, err := rdb.Ping().Result()
	if err != nil {
		log.WithFields(log.Fields{"err": err.Error()}).Info("could not connect to cache")
		return nil, err
	}
	redisClient = rdb

	return rdb, nil
}

func Write(key, val string) error {
	client, err := Get()
	if err != nil {
		log.WithFields(log.Fields{"err": err.Error()}).Info("could not connect to cache")
		return err
	}
	if _, err := client.Set(key, val, 0).Result(); err != nil {
		log.WithFields(log.Fields{"err": err.Error()}).Info("could not write to cache")

		return err
	}

	return nil
}

func Read(key string) (string, error) {
	client, err := Get()
	if err != nil {
		log.WithFields(log.Fields{"err": err.Error()}).Info("could not connect to cache")
		return "", err
	}
	return client.Get(key).Result()

}

func ReadHash(key, field string) (string, error) {
	client, err := Get()
	if err != nil {
		log.WithFields(log.Fields{"err": err.Error()}).Info("could not connect to cache")
		return "", err
	}
	return client.HGet(key, field).Result()

}

func WriteHash(key, field, val string) error {
	client, err := Get()
	if err != nil {
		log.WithFields(log.Fields{"err": err.Error()}).Info("could not connect to cache")
		return err
	}
	if _, err := client.HSet(key, field, val).Result(); err != nil {
		log.WithFields(log.Fields{"err": err.Error()}).Info("could not connect to cache")
		return err
	}

	return nil
}
