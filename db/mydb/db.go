// Copyright 2018 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

/**
 * Copyright (c) 2010-2016 Yahoo! Inc., 2017 YCSB contributors All rights reserved.
 * <p>
 * Licensed under the Apache License, Version 2.0 (the "License"); you
 * may not use this file except in compliance with the License. You
 * may obtain a copy of the License at
 * <p>
 * http://www.apache.org/licenses/LICENSE-2.0
 * <p>
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
 * implied. See the License for the specific language governing
 * permissions and limitations under the License. See accompanying
 * LICENSE file.
 */

package mydb

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"time"
	"log"
	"os"
	"strings"

	pb "github.com/scalog/scalog/zookeeper/zookeeperpb"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"github.com/magiconair/properties"
	"github.com/pingcap/go-ycsb/pkg/ycsb"
)

const (

)

type contextKey string

const stateKey = contextKey("mydb")

type myState struct {
	r *rand.Rand
	buf *bytes.Buffer
	clients []pb.ZooKeeperClient
}

// myDB
type mydb struct {
	verbose        bool
	randomizeDelay bool
	toDelay        int64
}

func (db *mydb) InitThread(ctx context.Context, _ int, _ int) context.Context {
	state := new(myState)
	state.r = rand.New(rand.NewSource(time.Now().UnixNano()))
	state.buf = new(bytes.Buffer)
	state.clients := []pb.ZooKeeperClient{}
	
	zkNodesIp := viper.GetStringSlice("zk-servers")
	zkPort := int32(viper.GetInt("zk-port"))

	zids_ := viper.Get("zids").([]interface{})
	var zids []int

	for _, zid_ := range zids_ {
		zids = append(zids, zid_.(int))
	}

	for idx, ip := range zkNodesIp {
		fmt.Printf("initiating client to: %v:%v\n", ip, zkPort+int32(zids[idx]))
		serverAddress := fmt.Sprintf("%v:%v", ip, zkPort+int32(zids[idx]))

		conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Failed to connect: %v", err)
		}
		defer conn.Close()

		client := pb.NewZooKeeperClient(conn)

		clients = append(clients, client)
	}

	return context.WithValue(ctx, stateKey, state)
}

func (db *mydb) CleanupThread(_ context.Context) {

}

func (db *mydb) Close() error {
	return nil
}

func (db *mydb) Read(ctx context.Context, table string, key string, fields []string) (map[string][]byte, error) {
	state := ctx.Value(stateKey).(*myState)
	client := state.clients[state.r.Intn(len(state.clients))]

	pathObj := &pb.Path{Path: key}
	getResponse, err := client.GetZNode(context.Background(), pathObj)
	if err != nil {
		log.Printf("GetZNode failed: %v", err)
	} else {
		fmt.Printf("ZNode Path: %v\nZNode Data: %v\n", getResponse.Path, string(getResponse.Data))
	}

	return nil, nil
}

func (db *mydb) BatchRead(ctx context.Context, table string, keys []string, fields []string) ([]map[string][]byte, error) {
	panic("The mydb has not implemented the batch operation")
}

func (db *mydb) Scan(ctx context.Context, table string, startKey string, count int, fields []string) ([]map[string][]byte, error) {
	panic("The mydb has not implemented the scan operation")
}

func (db *mydb) Update(ctx context.Context, table string, key string, values map[string][]byte) error {
	state := ctx.Value(stateKey).(*myState)
	client := state.clients[state.r.Intn(len(state.clients))]

	var x int := 0
	var data []byte
	for _, val := range values {
		if x == 0 {
			data := val
			x := x + 1
		}
	}

	znode := &pb.ZNode{
		Path: key,
		Data: data,
	}

	createResponse, err := client.CreateZNode(ctx.Background(), znode)
	if err != nil {
		log.Printf("CreateZNode failed: %v", err)
	} else {
		fmt.Printf("Created ZNode: %v\n", createResponse.Path)
	}

	return nil
}

func (db *mydb) BatchUpdate(ctx context.Context, table string, keys []string, values []map[string][]byte) error {
	panic("The mydb has not implemented the batch update operation")
}

func (db *mydb) Insert(ctx context.Context, table string, key string, values map[string][]byte) error {
	state := ctx.Value(stateKey).(*myState)

	return nil
}

func (db *mydb) BatchInsert(ctx context.Context, table string, keys []string, values []map[string][]byte) error {
	panic("The mydb has not implemented the batch insert operation")
}

func (db *mydb) Delete(ctx context.Context, table string, key string) error {
	panic("The mydb has not implemented the delete operation")
}

func (db *mydb) BatchDelete(ctx context.Context, table string, keys []string) error {
	panic("The mydb has not implemented the batch delete operation")
}

type mydbCreator struct{}

func (mydbCreator) Create(p *properties.Properties) (ycsb.DB, error) {
	db := new(mydb)

	return db, nil
}

func init() {
	ycsb.RegisterDBCreator("mydb", mydbCreator{})
}
