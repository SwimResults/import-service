package model

import "fmt"

type ImportFileStats struct {
	Created  ImportFileStatsEntries `json:"created"`
	Imported ImportFileStatsEntries `json:"imported"`
}

type ImportFileStatsEntries struct {
	Events            int `json:"events"`
	Athletes          int `json:"athletes"`
	Teams             int `json:"teams"`
	Starts            int `json:"starts"`
	Results           int `json:"results"`
	Disqualifications int `json:"disqualifications"`
}

func (stats *ImportFileStats) PrintReport() {
	fmt.Printf("\n\n--===[ IMPORT REPORT ]===--\n")
	fmt.Printf("  -> events: (%d, %d)\n", stats.Created.Events, stats.Imported.Events)
	fmt.Printf("  -> teams: (%d, %d)\n", stats.Created.Teams, stats.Imported.Teams)
	fmt.Printf("  -> athletes: (%d, %d)\n", stats.Created.Athletes, stats.Imported.Athletes)
	fmt.Printf("  -> starts: (%d, %d)\n", stats.Created.Starts, stats.Imported.Starts)
	fmt.Printf("  -> results: (%d, %d)\n", stats.Created.Results, stats.Imported.Results)
	fmt.Printf("  -> disqualifications: (%d, %d)\n", stats.Created.Disqualifications, stats.Imported.Disqualifications)
	fmt.Printf("\n---------------------------\n\n")
}
