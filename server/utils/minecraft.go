package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"sync"
)

type OnlineUUIDStruct struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func GetMinecraftInfoFromName(name string) (*OnlineUUIDStruct, error) {
	resp, err := http.Get("https://api.mojang.com/users/profiles/minecraft/" + name)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	onlineUUID := OnlineUUIDStruct{}
	if err = json.Unmarshal(body, &onlineUUID); err != nil {
		return nil, err
	}

	return &onlineUUID, nil
}

func GetMinecraftInfosFrom10Names(names []string) (*[]OnlineUUIDStruct, error) {
	size := len(names)
	if size > 10 {
		return nil, errors.New("Too many names")
	}

	if size == 0 {
		return &[]OnlineUUIDStruct{}, nil
	}

	nameString, err := json.Marshal(names)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post("https://api.mojang.com/profiles/minecraft", "application/json", bytes.NewBuffer(nameString))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	result := []OnlineUUIDStruct{}
	if err = json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func GetMinecraftInfosFromNames(names ...string) *[]OnlineUUIDStruct {
	results := []OnlineUUIDStruct{}
	max := len(names)
	mu, wg := sync.Mutex{}, &sync.WaitGroup{}
	for i := 0; i <= max; i += 10 {
		end := min(i+10, max)
		sliceNames := names[i:end]
		go func(names []string) {
			defer wg.Done()

			result, err := GetMinecraftInfosFrom10Names(names)
			if err == nil {
				mu.Lock()
				results = append(results, *result...)
				mu.Unlock()
				return
			}

			for _, name := range names {
				go func(name string) {
					defer wg.Done()
					info, _ := GetMinecraftInfoFromName(name)
					if info != nil {
						mu.Lock()
						results = append(results, *info)
						mu.Unlock()
					}
				}(name)
				wg.Add(1)
			}
		}(sliceNames)
		wg.Add(1)
	}
	wg.Wait()
	return &results
}
