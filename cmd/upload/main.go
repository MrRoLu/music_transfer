package main

import (
	"Test_api/internal/apple"
	"Test_api/internal/config"
	"Test_api/internal/storage"
	"fmt"
	"log"
)

func main() {
	cfg := config.LoadConfig()
	apple := apple.NewService(cfg.DevToken, cfg.MediaUserToken)

	storefront, err := apple.GetStorefront()

	if err != nil {
		fmt.Printf("Error in GetStorefront: %v\n", *storefront)
		log.Print(err)
	}

	// fmt.Printf("Storefront: %v\n", *storefront)

	songs, err := storage.LoadSongs("./data/songs.json")

	if err != nil {
		fmt.Printf("Error in main Apple LoadSongs: %v\n", err)
		log.Print(err)
	}

	// fmt.Printf("Songs: AName: %v Title: %v\n", *(*songs)[0].ArtistName, *(*songs)[0].Title)

	rightIds, notFoundSongs := apple.FindSongs(songs, *storefront, 5)

	fmt.Printf("All songs: %v\nFound songs: %v\nNot found songs: %v\n", len(*songs), len(*rightIds), len(*notFoundSongs))

	if err := apple.AddSongs(rightIds); err != nil {
		fmt.Printf("Error in main Apple AddSongs: %v", err)
		log.Print(err)
	}
}
