package utils

import (
	"fmt"
	"testing"

	mapset "github.com/deckarep/golang-set/v2"
)

func TestGetMinecraftInfoFromName(t *testing.T) {
	for _, name := range []string{"Steve", "StevE"} {
		info, err := GetMinecraftInfoFromName(name)
		if err != nil {
			t.Error(err)
			continue
		}

		if info == nil {
			t.Errorf("Expected info, got nil")
			continue
		}

		if info.Name != "Steve" {
			t.Errorf("Expected name Steve, got %s", info.Name)
		} else if info.ID != "8667ba71b85a4004af54457a9734eed7" {
			t.Errorf("Expected 8667ba71b85a4004af54457a9734eed7, got %s", info.ID)
		}
	}
}

func TestGetMinecraftInfosFrom10Names(t *testing.T) {
	names := []string{
		"Steve", "Alex", "Noor", "Sunny", "Ari",
		"Zuri", "Makena", "Kai", "Efe",
	}

	infos, err := getMinecraftInfosFrom10Names(names)
	if err != nil {
		t.Error(err)
		return
	}

	players := mapset.NewSet(
		OnlineUUIDStruct{ID: "ec561538f3fd461daff5086b22154bce", Name: "Alex"},
		OnlineUUIDStruct{ID: "938e960d50ab489b9b2aaf3751942989", Name: "Ari"},
		OnlineUUIDStruct{ID: "20bf454f34e34010a378613546e3d0f9", Name: "efe"},
		OnlineUUIDStruct{ID: "cf9858b6ed4946538e47f0e4214539f7", Name: "Kai"},
		OnlineUUIDStruct{ID: "6c4bc87ce82944efa1ad63d45e2b9545", Name: "Makena"},
		OnlineUUIDStruct{ID: "2d9f2227592b481d8433d13b69473ccc", Name: "noor"},
		OnlineUUIDStruct{ID: "8667ba71b85a4004af54457a9734eed7", Name: "Steve"},
		OnlineUUIDStruct{ID: "bafbe1cb77b348099fa3c89604bda644", Name: "Sunny"},
		OnlineUUIDStruct{ID: "f5e039b93b8a45109ee8e7552e098c55", Name: "Zuri"},
	)
	fmt.Println(players.Difference(mapset.NewSet(*infos...)))

	// waitMatchData := mapset.NewSet(*infos...)
	// waitMatchData.Difference()
	// for _, player := range players {
	// 	if waitMatchData.Contains(player) {
	// 		waitMatchData.Remove(player)
	// 	} else {
	// 		t.Errorf("Expected %v, got nil", player)
	// 	}
	// }
	// if waitMatchData.Cardinality() != 0 {
	// 	t.Errorf("Expected 0, got %d", waitMatchData.Cardinality())
	// }
}
