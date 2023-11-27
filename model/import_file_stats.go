package model

import "fmt"

type ImportFileStats struct {
	Created  ImportFileStatsEntries `json:"created"`
	Imported ImportFileStatsEntries `json:"imported"`
}

type ImportFileStatsEntries struct {
	Events            int `json:"events"`
	AgeGroups         int `json:"age_groups"`
	Athletes          int `json:"athletes"`
	Teams             int `json:"teams"`
	Heats             int `json:"heats"`
	Starts            int `json:"starts"`
	Results           int `json:"results"`
	Disqualifications int `json:"disqualifications"`
}

func (stats *ImportFileStats) PrintReport() {
	if stats == nil {
		println("trying to print nil ImportFileStats")
		return
	}
	fmt.Printf("\n\n--===[ IMPORT REPORT ]===--\n")
	fmt.Printf("  -> events: (%d, %d)\n", stats.Created.Events, stats.Imported.Events)
	fmt.Printf("  -> age_groups: (%d, %d)\n", stats.Created.AgeGroups, stats.Imported.AgeGroups)
	fmt.Printf("  -> teams: (%d, %d)\n", stats.Created.Teams, stats.Imported.Teams)
	fmt.Printf("  -> athletes: (%d, %d)\n", stats.Created.Athletes, stats.Imported.Athletes)
	fmt.Printf("  -> heats: (%d, %d)\n", stats.Created.Heats, stats.Imported.Heats)
	fmt.Printf("  -> starts: (%d, %d)\n", stats.Created.Starts, stats.Imported.Starts)
	fmt.Printf("  -> results: (%d, %d)\n", stats.Created.Results, stats.Imported.Results)
	fmt.Printf("  -> disqualifications: (%d, %d)\n", stats.Created.Disqualifications, stats.Imported.Disqualifications)
	fmt.Printf("\n---------------------------\n\n")
}
