package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
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

func getMinecraftInfosFrom10Names(names []string) (*[]OnlineUUIDStruct, error) {
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

func GetMinecraftInfosFromNames(names ...string) ([]*OnlineUUIDStruct, error) {
	results := []OnlineUUIDStruct{}
	max := len(names)
	for i := 0; i < max; i += 10 {
		end := i + 10
		if end > max {
			end = max
		}

		sliceNames := names[i:end]
		result, err := getMinecraftInfosFrom10Names(sliceNames)
		if err != nil {
			continue
		}

		results = append(results, *result...)
	}
	// if err != nil {
	// 	getOnes(sliceNames...)
	// 	return
	// }
	// defer resp.Body.Close()

	// body, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	getOnes(sliceNames...)
	// 	return
	// }

	// onlineUUID := []OnlineUUIDStruct{}
	// if err := json.Unmarshal(body, &onlineUUID); err != nil {
	// 	getOnes(sliceNames...)
	// 	return
	// }
	return nil, nil

	// go func(start, end int) {
	// 	defer wg.Done()

	// 	sliceNames := names[start:end]
	// 	nameString, err := json.Marshal(sliceNames)
	// 	if err != nil {
	// 		getOnes(sliceNames...)
	// 		return
	// 	}

	// 	resp, err := http.Post("https://api.mojang.com/profiles/minecraft", "application/json", bytes.NewBuffer(nameString))
	// 	if err != nil {
	// 		getOnes(sliceNames...)
	// 		return
	// 	}
	// 	defer resp.Body.Close()

	// 	body, err := io.ReadAll(resp.Body)
	// 	if err != nil {
	// 		getOnes(sliceNames...)
	// 		return
	// 	}

	// 	onlineUUID := []OnlineUUIDStruct{}
	// 	if err := json.Unmarshal(body, &onlineUUID); err != nil {
	// 		getOnes(sliceNames...)
	// 		return
	// 	}

	// 	// Check if the length is equal
	// 	if len(onlineUUID) != len(sliceNames) {
	// 		getOnes(sliceNames...)
	// 		return
	// 	}

	// 	mu.Lock()
	// 	defer mu.Unlock()
	// 	for _, u := range onlineUUID {
	// 		result[u.Name] = u.ID
	// 	}
	// }(i, end)

}
