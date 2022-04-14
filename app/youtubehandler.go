package main

//https://youtu.be/EUpdJjfhWJA?list=PLSvCAHoiHC_rqKbcu1ummWVpLTDBNZHH7&t=588
import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type YoutubeStats struct {
	Subscribers 	int `json:"subscribers"`
	ChannelName 	string `json:"channelName"`
	MinutesWatched 	int `json:"minutesWatched"`
	Views 			int `json:"views"`
}

func getChannelStats(apiKey string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// w.Write([]byte("response!"))
		yt := YoutubeStats{
			Subscribers: 	5,
			ChannelName: 	"my channel",
			MinutesWatched: 50,
			Views: 			100,
		}

		ctx := context.Background()
		yts, err := youtube.NewService(ctx, option.WithAPIKey(apiKey))
		if err != nil {
			fmt.Println("failed to create service")
			panic(err)
		}

		call := yts.Channels.List([]string{"snippet", "contentDetails", "statistics"})
		response, err := call.Id("UCt7T2EvYBqvlxNU3fbE4Y7g").Do()
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(response.Items[0].Snippet)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(yt); err != nil {
			panic(err)
		}
	}
}