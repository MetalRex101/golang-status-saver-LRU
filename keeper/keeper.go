package keeper

import (
	"fmt"
	"gitlab.com/artilligence/http-db-saver/domain"
	"io/ioutil"
	"log"
	"os"
	"time"
)

const cachePath = "storage/cache"
const currentCacheDir = "current"
const oldCacheDir = "old"
const flushTimeMinutes = 10

type keeper struct {
	receive  <-chan *domain.Message
	stop     <-chan bool
	storage  []*domain.Message
	gobCoder domain.BinaryCoder
	repos    domain.Repositories
}

func NewKeeper(
	gobCoder domain.BinaryCoder,
	repos domain.Repositories,
) domain.KeeperClient {
	var queue = make(chan *domain.Message)
	var stop = make(chan bool)

	k := &keeper{
		receive:  queue,
		stop:     stop,
		storage:  []*domain.Message{},
		repos:    repos,
		gobCoder: gobCoder,
	}

	c := &client{
		send: queue,
	}

	go k.start()

	return c
}

func (k *keeper) start() {
	flushTimer := time.NewTimer(flushTimeMinutes * time.Minute)

	for {
		select {
		case mess, ok := <-k.receive:
			k.store(mess)

			if !ok {
				if err := k.performStop(); err != nil {
					fmt.Println(fmt.Printf("failed to perform stop: %s", err))
				}
				return
			}
		case <-flushTimer.C:
			if err := k.flushToDB(); err != nil {
				log.Println(fmt.Printf("failed to flush storage to db: %s", err))
			}

			flushTimer.Reset(flushTimeMinutes * time.Minute)
		}
	}
}

func (k *keeper) store(message *domain.Message) {
	k.storage = append(k.storage, message)

	go k.writeOnDisk(message, time.Now().Unix())
}

func (k *keeper) performStop() error {
	if err := k.flushToDB(); err != nil {
		return fmt.Errorf("failed to flush storage to db: %s", err)
	}

	return nil
}

func (k *keeper) flushToDB() error {
	for _, m := range k.storage {
		switch m.Type {
		case domain.TypeEntity:
			en, ok := m.Data.(*domain.Entity)
			if !ok {
				return fmt.Errorf("failed to cast %+v to Entity", m.Data)
			}

			err := k.repos.Entity.Insert(en)
			if err != nil {
				return err
			}
		}
	}

	swapOldAndNew()

	return nil
}

func (k *keeper) writeOnDisk(mess *domain.Message, timestamp int64) {
	b, err := k.gobCoder.Encode(mess)
	if err != nil {
		fmt.Println(fmt.Printf("failed to encode value to GOB64: %s", err))
	}

	if err := os.MkdirAll(getCurrentDir(), 0700); err != nil {
		log.Panicf("failed to create new current cache dir: %s", err)
	}

	if err := ioutil.WriteFile(getFileName(mess.Name, timestamp), b, 0700); err != nil {
		fmt.Println(fmt.Printf("failed to write file on disk: %s", err))
	}
}

func (k *keeper) loadStorage() error {
	current := getCurrentDir()

	files, err := ioutil.ReadDir(current)
	if err != nil {
		log.Panicf("failed to read current dir: %s", err)
	}

	for i := len(files) - 1; i >= 0; i-- {
		var f = files[i]

		fileName := fmt.Sprintf("%s/%s", current, f.Name())

		content, err := ioutil.ReadFile(fileName)
		if err != nil {
			return fmt.Errorf("failed to read file: %s", err)
		}

		decoded, err := k.gobCoder.Decode(string(content))
		if err != nil {
			return fmt.Errorf("failed to decode file content: %s", err)
		}

		mess, ok := decoded.(*domain.Message)
		if !ok {
			return fmt.Errorf("failed to cast file content to message %+v", decoded)
		}

		k.storage = append(k.storage, mess)
	}

	return nil
}

func getFileName(key string, timestamp int64) string {
	return fmt.Sprintf("%s/%d-%s", getCurrentDir(), timestamp, key)
}

func getCurrentDir() string {
	return fmt.Sprintf("./%s/%s", cachePath, currentCacheDir)
}

func getOldDir() string {
	return fmt.Sprintf("./%s/%s", cachePath, oldCacheDir)
}

func swapOldAndNew() {
	if err := os.Remove(getOldDir()); err != nil {
		log.Panicf("failed to remove old cache dir: %s", err)
	}

	if err := os.Rename(getCurrentDir(), getOldDir()); err != nil {
		log.Panicf("failed to rename current cache dir: %s", err)
	}
}
