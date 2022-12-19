package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/fernandezvara/router"
	clientv3 "go.etcd.io/etcd/client/v3"
)

//
// run a etcd instance. example:
//
// docker run --name etcd -p 2379:2379 -e ALLOW_NONE_AUTHENTICATION=yes -d bitnami/etcd:latest
//

func main() {

	var (
		r       *router.Router
		client  *clientv3.Client
		watcher clientv3.WatchChan
		ctx     context.Context = context.Background()
		err     error
	)

	rand.Seed(time.Now().UnixNano())

	r = router.New()

	client, err = clientv3.New(clientv3.Config{
		Endpoints: []string{"localhost:2379"},
	})

	if err != nil {
		panic(err)
	}

	r.Method("PUT").Insert("example/:id", examplePutFunc)
	r.Method("DELETE").Insert("example/:id", exampleDeleteFunc)

	go func() {
		watcher = client.Watch(ctx, "", clientv3.WithPrefix())
		for resp := range watcher {
			for _, event := range resp.Events {
				fmt.Println("key:", string(event.Kv.Key))
				fmt.Println("err:", r.Method(event.Type.String()).Execute(string(event.Kv.Key)))
			}
		}
	}()

	// insert and delete random items
	for a := 0; a < 3; a++ {
		id := randomString()
		client.KV.Put(ctx, fmt.Sprintf("example/%s", id), "")
		time.Sleep(1 * time.Second)
		client.KV.Put(ctx, fmt.Sprintf("example/%s", id), "")
		time.Sleep(1 * time.Second)
		client.KV.Put(ctx, fmt.Sprintf("non/exists/%s", id), "")
		time.Sleep(1 * time.Second)
		client.KV.Delete(ctx, fmt.Sprintf("example/%s", id))
		time.Sleep(1 * time.Second)
	}

}

func examplePutFunc(p *router.Params) error {
	fmt.Println("func: PUT:", p.Param("id"))
	return nil
}

func exampleDeleteFunc(p *router.Params) error {
	fmt.Println("func: DEL:", p.Param("id"))
	return nil
}

func randomString() string {
	b := make([]byte, 10)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
