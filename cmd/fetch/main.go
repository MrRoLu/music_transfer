package main

import (
	"Test_api/internal/storage"
	"Test_api/internal/yandex"
	"fmt"
	"log"
)

func main() {
	ya := yandex.NewService()

	query := map[string]string{
		"resumeStream": "false",
		"richTracks":   "false",
	}
	resGet, err := ya.GetPlaylist(query)

	if err != nil {
		fmt.Printf("Err in main GET: %v", err)
		return
	}

	fmt.Printf("Len res: %v", len(resGet["trackIds"]))

	for k, v := range resGet {
		fmt.Printf("Key: %v, Value: %v", k, v)
	}

	res, err := ya.GetFinallist(resGet)

	if err != nil {
		fmt.Printf("Err in main POST: %v", err)
	}

	if err := storage.SaveSongs(res, storage.DefaultPath); err != nil {
		log.Fatal(err)
	}

	// for k, v := range res.Result {
	// 	fmt.Printf("Number: %v, ID: %v, RealID: %v, Title: %v, ArtistID: %v, Name: %v\n", k, *v.ID, *v.RealId, *v.Title, *v.Artists[0].ID, *v.Artists[0].Name)
	// }
}
