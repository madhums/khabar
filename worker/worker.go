package worker

import (
	"bytes"
	"encoding/json"
	"github.com/changer/sc-notifications/config"
	"github.com/changer/sc-notifications/db"
	"github.com/changer/sc-notifications/dbapi/notification_instance"
	"gopkg.in/simversity/gotracer.v1"
	"gopkg.in/simversity/gottp.v2"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

func NotificationWorker(errChan chan int, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	gottp.OnSysExit(func() {
		errChan <- 1
		errChan <- 1
		wg.Wait()
	})

	go notificationWorkerInternal(errChan, wg)

	for {
		s := <-errChan
		if s == 0 {
			go notificationWorkerInternal(errChan, wg)
		} else {
			return
		}
	}
}

func errorSignaler(errChan chan<- int) {
	errChan <- 0
}

func notificationWorkerInternal(errChan chan int, wg *sync.WaitGroup) {
	wg.Add(1)

	defer wg.Done()
	defer errorSignaler(errChan)
	defer gotracer.Tracer{Dummy: true}.Notify(func() string {
		return "Error in worker"
	})

	for {

		select {
		case <-errChan:
			return
		default:
		}

		time.Sleep(5 * time.Second)
		ntfInsts := *fetchNotifications(*new(time.Time))
		for i, ntf := range ntfInsts {
			body, err := json.Marshal(ntf)
			if err != nil {
				panic(err)
			}
			randomPanic()
			http.DefaultClient.Post("http://"+config.Settings.Gottp.Listen+"/notifications/"+ntfInsts[i].DestinationUri, "application/json", bytes.NewReader(body))
		}

	}
}

func randomPanic() {
	if rand.Int()%4 == 0 {
		panic("Panicing...")
	}
}

func fetchNotifications(timeStamp time.Time) *[]notification_instance.NotificationInstance {
	return &[]notification_instance.NotificationInstance{
		notification_instance.NotificationInstance{
			Organization:     "org1",
			AppName:          "app1",
			User:             "user1",
			NotificationType: "type1",
			DestinationUri:   "destUri1",
			IsPending:        true,
			Context: db.M{
				"a": "b",
			},
		},
		notification_instance.NotificationInstance{
			Organization:     "org2",
			AppName:          "app2",
			User:             "user2",
			NotificationType: "type2",
			DestinationUri:   "destUri2",
			IsPending:        true,
			Context: db.M{
				"a": "b",
			},
		},
	}

}
