package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

// Enums (as constants)
const (
	Caucasian = "Caucasian"
	Mongoloid = "Mongoloid"
	Negroid   = "Negroid"
	Android   = "Android"

	Neutral   = "Neutral"
	Happiness = "Happiness"
	Anger     = "Anger"
	Surprise  = "Surprise"
	Fear      = "Fear"
	Sadness   = "Sadness"
	Disgust   = "Disgust"

	Baby         = "Baby"
	Kid          = "Kid"
	Teenager     = "Teenager"
	YoungAdult   = "YoungAdult"
	MaturedAdult = "MaturedAdult"
	Senior       = "Senior"

	Male   = "Male"
	Female = "Female"
)

// Superhero struct definition
type Superhero struct {
	FileName string     `json:"file_name"`
	BBox     [4]float64 `json:"bbox"`
	Gender   string     `json:"gender"`
	Emotion  string     `json:"emotion"`
	Age      string     `json:"age"`
	Race     string     `json:"race"`
}

// GroupResult struct definition
type GroupResult struct {
	Type  string      `json:"type"`
	Group []Superhero `json:"group"`
}

// Helper functions to check the criteria
func isBalanceGuardians(group []Superhero) bool {
	if len(group) != 4 {
		return false
	}
	maleCount, femaleCount := 0, 0
	races := make(map[string]struct{})
	for _, hero := range group {
		if hero.Gender == Male {
			maleCount++
		} else if hero.Gender == Female {
			femaleCount++
		}
		races[hero.Race] = struct{}{}
	}
	return maleCount == 2 && femaleCount == 2 && len(races) >= 3
}

func isInsideOut(group []Superhero) bool {
	if len(group) != 4 {
		return false
	}
	emotions := make(map[string]struct{})
	races := make(map[string]struct{})
	for _, hero := range group {
		emotions[hero.Emotion] = struct{}{}
		races[hero.Race] = struct{}{}
	}
	return len(emotions) == 4 && len(races) >= 3
}

func isTheIncredibles(group []Superhero) bool {
	if len(group) != 4 {
		return false
	}
	ages := make(map[string]struct{})
	emotions := make(map[string]struct{})
	for _, hero := range group {
		ages[hero.Age] = struct{}{}
		emotions[hero.Emotion] = struct{}{}
	}
	return len(ages) == 4 && len(emotions) == 4
}

func isPotentiallyBalanceGuardians(group []Superhero, usedFileName map[string]struct{}) bool {
	for _, hero := range group {
		if _, used := usedFileName[hero.FileName]; used {
			return false
		}
	}
	if len(group) > 4 {
		return false
	}
	maleCount, femaleCount := 0, 0
	races := make(map[string]struct{})
	for _, hero := range group {
		if hero.Gender == Male {
			maleCount++
		} else if hero.Gender == Female {
			femaleCount++
		}
		races[hero.Race] = struct{}{}
	}

	// There should be no more than 2 males and 2 females in the current list
	if maleCount > 2 || femaleCount > 2 {
		return false
	}

	// If the list is not full yet, we only need to check the current constraints
	if len(group) < 4 {
		return true
	}

	// When the list is full, check if it meets all the criteria for Balance Guardians
	return maleCount == 2 && femaleCount == 2 && len(races) >= 3
}

// isPotentiallyInsideOut checks if the current array of heroes can still potentially form a valid Inside Out group
func isPotentiallyInsideOut(group []Superhero, usedFileName map[string]struct{}) bool {
	for _, hero := range group {
		if _, used := usedFileName[hero.FileName]; used {
			return false
		}
	}
	if len(group) > 4 {
		return false
	}
	emotions := make(map[string]struct{})
	races := make(map[string]struct{})
	for _, hero := range group {
		emotions[hero.Emotion] = struct{}{}
		races[hero.Race] = struct{}{}
	}

	// No duplicate emotions are allowed
	if len(emotions) != len(group) {
		return false
	}

	// If the list is not full yet, we only need to check the current constraints
	if len(group) < 4 {
		return true
	}

	// When the list is full, check if it meets all the criteria for Inside Out
	return len(emotions) == 4 && len(races) >= 3
}

// isPotentiallyTheIncredibles checks if the current array of heroes can still potentially form a valid The Incredibles group
func isPotentiallyTheIncredibles(group []Superhero, usedFileName map[string]struct{}) bool {
	for _, hero := range group {
		if _, used := usedFileName[hero.FileName]; used {
			return false
		}
	}
	if len(group) > 4 {
		return false
	}
	ages := make(map[string]struct{})
	emotions := make(map[string]struct{})
	for _, hero := range group {
		ages[hero.Age] = struct{}{}
		emotions[hero.Emotion] = struct{}{}
	}

	// No duplicate ages or emotions are allowed
	if len(ages) != len(group) || len(emotions) != len(group) {
		return false
	}

	// If the list is not full yet, we only need to check the current constraints
	if len(group) < 4 {
		return true
	}

	// When the list is full, check if it meets all the criteria for The Incredibles
	return len(ages) == 4 && len(emotions) == 4
}

// Function to perform greedy grouping
func formGreedyGroups(heroes []Superhero, isPotentialValidGroup func([]Superhero, map[string]struct{}) bool, groupName string, ctx context.Context, usedFileNames map[string]struct{}) [][]Superhero {
	var groups [][]Superhero

	var dfs func(currentGroup []Superhero, index int)
	dfs = func(currentGroup []Superhero, index int) {
		select {
		case <-ctx.Done():
			return
		default:
		}

		if len(currentGroup) == 4 && isPotentialValidGroup(currentGroup, usedFileNames) {
			groups = append(groups, append([]Superhero{}, currentGroup...))
			file, _ := json.MarshalIndent(groups, "", " ")
			_ = os.WriteFile(groupName, file, 0644)
			for _, hero := range currentGroup {
				usedFileNames[hero.FileName] = struct{}{}
			}
			return
		}

		if !isPotentialValidGroup(currentGroup, usedFileNames) {
			return
		}

		if index >= len(heroes) || len(currentGroup) >= 4 {
			return
		}

		for i := index; i < len(heroes); i++ {
			if _, used := usedFileNames[heroes[i].FileName]; !used {
				dfs(append(currentGroup, heroes[i]), i+1)
			}
		}
	}

	dfs([]Superhero{}, 0)

	return groups
}

// Function to group heroes into valid groups and format the result
func groupHeroes(heroes []Superhero) map[string][][]Superhero {
	points := 0
	result := map[string][][]Superhero{
		"type_1": {},
		"type_2": {},
		"type_3": {},
	}
	usedFileNames := make(map[string]struct{})

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	// Form groups for The Incredibles
	theIncrediblesGroups := formGreedyGroups(heroes, isPotentiallyTheIncredibles, "run/typeu_3.json", ctx, usedFileNames)
	result["type_3"] = theIncrediblesGroups
	points += len(theIncrediblesGroups)

	ctx2, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	// Form groups for Balance Guardians
	balanceGuardiansGroups := formGreedyGroups(heroes, isPotentiallyBalanceGuardians, "run/typeu_1.json", ctx2, usedFileNames)
	result["type_1"] = balanceGuardiansGroups
	points += len(balanceGuardiansGroups)

	ctx3, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	// Form groups for Inside Out
	insideOutGroups := formGreedyGroups(heroes, isPotentiallyInsideOut, "run/typeu_2.json", ctx3, usedFileNames)
	result["type_2"] = insideOutGroups
	points += len(insideOutGroups)

	// // write the rest of hero to json
	// rest_of_heroes_file := "after_run_input.json"
	// rest_of_heroes := make([]Superhero, 0)
	// for _, hero := range heroes {
	// 	if _, used := usedFileNames[hero.FileName]; !used {
	// 		rest_of_heroes = append(rest_of_heroes, hero)
	// 	}
	// }
	// rest_of_hero_data, _ := json.MarshalIndent(rest_of_heroes, "", "  ")
	// os.WriteFile(rest_of_heroes_file, rest_of_hero_data, 0644)

	fmt.Printf("Total points: %d\n", points)
	return result
}

func main() {
	// Read the JSON file
	data, err := os.ReadFile("./final_untrusted.json")
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	// Parse JSON data
	var heroes []Superhero
	if err := json.Unmarshal(data, &heroes); err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	// Group heroes
	result := groupHeroes(heroes)

	// Output the result
	resultJSON, _ := json.MarshalIndent(result, "", "  ")
	os.WriteFile("run/final_untrusted.json", resultJSON, 0644)
}
