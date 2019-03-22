// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"encoding/json"
	"fmt"
	"gitlab.com/artilligence/http-db-saver/binary_coder"
	"gitlab.com/artilligence/http-db-saver/db"
	"gitlab.com/artilligence/http-db-saver/domain"
	"gitlab.com/artilligence/http-db-saver/keeper"
	"gitlab.com/artilligence/http-db-saver/lru/in_memory"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

// httpCmd represents the http command
var httpCmd = &cobra.Command{
	Use:   "http",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("http called")
	},
}

func init() {
	rootCmd.AddCommand(httpCmd)
	cache := in_memory.NewInMemoryLRU()

	repos := domain.Repositories{
		Entity: db.NewEntityRepo(),
	}

	http.HandleFunc("/receive", func(writer http.ResponseWriter, request *http.Request) {
		if err := request.ParseForm(); err != nil {
			log.Panicf("failed to parse form: %s", err)
		}

		id := request.Form.Get("id")
		status, err := strconv.ParseBool(request.Form.Get("status"))
		if err != nil {
			log.Panicf("failed to convert status to bool: %s", err)
		}

		intId, err := strconv.Atoi(request.Form.Get("id"))
		if err != nil {
			log.Panicf("faield to convert id to int: %s", err)
		}

		mess, ok := cache.Get(id)

		if !ok {
			mess = repos.Entity.Get(int64(intId))
			cache.Put(id, status)
		} else {
			mess = mess.(*domain.Entity)
		}

		b, err := json.Marshal(mess)
		if err != nil {
			log.Panicf("failed to marshal message: %s", err)
		}

		writer.WriteHeader(http.StatusOK)
		if _, err := writer.Write(b); err != nil {
			log.Panicf("failed to write response: %s", err)
		}
	})

	http.HandleFunc("/send", func(writer http.ResponseWriter, request *http.Request) {
		if err := request.ParseForm(); err != nil {
			log.Panicf("failed to parse form: %s", err)
		}

		keep := keeper.NewKeeper(binary_coder.NewGob64Coder(), repos)
		id, err := strconv.Atoi(request.Form.Get("id"))
		if err != nil {
			log.Panicf("faield to convert id to int: %s", err)
		}

		status, err := strconv.ParseBool(request.Form.Get("status"))
		if err != nil {
			log.Panicf("failed to convert status to bool: %s", err)
		}

		keep.Send(&domain.Message{Name: "some-name", Type: "type", CreatedAt: time.Now(), Data: &domain.Entity{
			ID: int64(id),
			Status: status,
		}})

		writer.WriteHeader(http.StatusOK)
	})

	fmt.Println("[*] server is stated on port :10000")
	if err := http.ListenAndServe(":10000", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
