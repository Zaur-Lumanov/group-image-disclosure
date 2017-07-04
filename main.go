package main

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/urShadow/go-vk-api"
)

func main() {
	if len(os.Args) < 4 {
		return
	}

	accessToken := os.Args[1]

	api := vk.New("ru")

	err := api.Init(accessToken)

	if err != nil {
		log.Fatalln(err)
	}

	var members []int
	var offset int
	groupID := os.Args[2]
	photoURL := os.Args[3]

	getMembers(groupID, api, &members, &offset)

	uri := strings.Split(photoURL, "/")

	if len(uri) < 7 {
		log.Fatalln("Invalid photo url")
	}

	datapart := []byte(uri[4])[len(uri[4])-3:]

	for _, v := range members {
		strint := strconv.Itoa(v)
		if string(datapart) == string([]byte(strint)[len(strint)-3:]) {
			log.Println("https://vk.com/id" + strconv.Itoa(v))
		}
	}
}

func getMembers(groupID string, api *vk.VK, members *[]int, offset *int) {
	data, err := api.CallMethod("groups.getMembers", vk.RequestParams{
		"group_id": groupID,
		"offset":   strconv.Itoa(*offset * 1000),
	})

	time.Sleep(200 * time.Millisecond)

	if err != nil {
		log.Fatalln(err)
	}

	type JSONBody struct {
		Response struct {
			Count int   `json:"count"`
			Items []int `json:"items"`
		} `json:"response"`
	}

	var body JSONBody

	if err := json.Unmarshal(data, &body); err != nil {
		log.Fatalln(err)
	}

	for _, v := range body.Response.Items {
		*members = append(*members, v)
	}

	log.Println(body.Response.Count, len(*members))

	if body.Response.Count != len(*members) {
		*offset++
		getMembers(groupID, api, members, offset)
	}
}
