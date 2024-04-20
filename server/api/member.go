package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"server/config"
	"server/utils"
	"sync"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/yaml"
)

const DEFAULT_UUID = "00000000000000000000000000000000"

type WhitelistStruct struct {
	Whitelist map[string][]string `json:"whitelist"`
}

type APIMemberStruct struct {
	UUID         string `json:"uuid"`         // The uuid of the member.
	Name         string `json:"name"`         // The name of the member.
	Introduction string `json:"introduction"` // The introduction of the member.
}

// map[groups]APIMemberStruct
type APIWhitelistStruct = map[string][]APIMemberStruct

func HandlerGetWhitelistMembers(ctx *gin.Context) {
	apiWhitelist, err := parseWhiteListToAPIStruct(config.Get().API.MemberFile)
	if err != nil {
		log.WithError(err).Error("Failed to parse whitelist file")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to parse whitelist file",
		})
		return
	}
	ctx.JSON(http.StatusOK, apiWhitelist)
}

func parseWhiteListToAPIStruct(path string) (*APIWhitelistStruct, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	d := yaml.NewYAMLOrJSONDecoder(f, 4096)

	whitelist := WhitelistStruct{}
	if err := d.Decode(&whitelist); err != nil {
		return nil, err
	}

	uuidCache := loadUUIDCache()

	fetchUUID := getUUIDsToFetch(whitelist.Whitelist, uuidCache)

	onlineUUIDs := fetchOnlineUUIDs(fetchUUID)

	for name, uuid := range onlineUUIDs {
		uuidCache[name] = uuid
	}

	// Update the cache file
	cacheData, err := json.Marshal(uuidCache)
	if err == nil {
		if err := utils.AutoWriteFile(config.Get().GetCacheUUIDPath(), cacheData, 0644); err != nil {
			log.WithError(err).Error("Failed to update the cache file")
		}
	}

	apiWhitelist := convertToAPIStruct(whitelist, uuidCache)

	return &apiWhitelist, nil
}

func loadUUIDCache() map[string]string {
	uuidCache := map[string]string{}
	uuidCacheValue, err := os.ReadFile(config.Get().GetCacheUUIDPath())
	if err == nil {
		json.Unmarshal(uuidCacheValue, &uuidCache)
	} else if !errors.Is(err, os.ErrNotExist) {
		log.WithError(err).Error("Failed to read UUID cache file")
	}
	return uuidCache
}

func getUUIDsToFetch(whitelist map[string][]string, uuidCache map[string]string) []string {
	fetchUUID := map[string]struct{}{}
	for _, members := range whitelist {
		for _, member := range members {
			if _, ok := uuidCache[member]; !ok {
				fetchUUID[member] = struct{}{}
			}
		}
	}

	keys := make([]string, 0, len(fetchUUID))
	for k := range fetchUUID {
		keys = append(keys, k)
	}

	return keys
}

type OnlineUUIDStruct struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func fetchOnlineUUIDs(names []string) map[string]string {
	result := map[string]string{}
	var wg sync.WaitGroup
	var mu sync.Mutex

	max := len(names)
	for i := 0; i < max; i += 10 {
		end := i + 10
		if end > max {
			end = max
		}

		wg.Add(1)
		getOne := func(name string) error {
			resp, err := http.Get("https://api.mojang.com/users/profiles/minecraft/" + name)
			if err != nil {
				return err
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			onlineUUID := OnlineUUIDStruct{}
			if err := json.Unmarshal(body, &onlineUUID); err != nil {
				return err
			}

			mu.Lock()
			defer mu.Unlock()
			if onlineUUID.ID != name {
				result[name] = DEFAULT_UUID
			}
			result[onlineUUID.Name] = onlineUUID.ID

			return nil
		}
		getOnes := func(name ...string) {
			for _, n := range name {
				if err := getOne(n); err != nil {
					result[n] = DEFAULT_UUID
					log.WithError(err).Error(fmt.Sprintf("Failed to get the UUID of %s", n))
				}
			}
		}

		go func(start, end int) {
			defer wg.Done()

			sliceNames := names[start:end]
			nameString, err := json.Marshal(sliceNames)
			if err != nil {
				getOnes(sliceNames...)
				return
			}

			resp, err := http.Post("https://api.mojang.com/profiles/minecraft", "application/json", bytes.NewBuffer(nameString))
			if err != nil {
				getOnes(sliceNames...)
				return
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				getOnes(sliceNames...)
				return
			}

			onlineUUID := []OnlineUUIDStruct{}
			if err := json.Unmarshal(body, &onlineUUID); err != nil {
				getOnes(sliceNames...)
				return
			}

			// Check if the length is equal
			if len(onlineUUID) != len(sliceNames) {
				getOnes(sliceNames...)
				return
			}

			mu.Lock()
			defer mu.Unlock()
			for _, u := range onlineUUID {
				result[u.Name] = u.ID
			}
		}(i, end)
	}

	wg.Wait()
	return result
}

func convertToAPIStruct(whitelist WhitelistStruct, uuidCache map[string]string) APIWhitelistStruct {
	apiWhitelist := APIWhitelistStruct{}
	for group, members := range whitelist.Whitelist {
		for _, member := range members {
			apiMember := APIMemberStruct{Name: member}
			if uuid, ok := uuidCache[member]; ok && uuid != "" && uuid != DEFAULT_UUID {
				apiMember.UUID = uuid
			}
			apiWhitelist[group] = append(apiWhitelist[group], apiMember)
		}
	}
	return apiWhitelist
}
