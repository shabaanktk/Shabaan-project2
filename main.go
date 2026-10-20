package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

func main() {
	file, err := os.Open("salaries.csv")
	if err != nil {
		fmt.Println("Could not open salaries.csv:", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Could not read the CSV:", err)
		return
	}

	if len(rows) < 2 {
		fmt.Println("The file has no salary data.")
		return
	}

	var entry2024 []float64
	var entry2023 []float64
	var mid2024 []float64

	var remoteEntry []float64
	var onsiteEntry []float64
	var remoteAll []float64
	var onsiteAll []float64

	titleSalaries := make(map[string][]float64)

	// Start at 1 because row 0 contains the column names.
	for i := 1; i < len(rows); i++ {
		row := rows[i]

		if len(row) < 11 {
			fmt.Println("Missing columns on row", i+1)
			return
		}

		title := row[0]
		level := row[1]
		model := row[3]

		year, err := strconv.Atoi(row[4])
		if err != nil {
			fmt.Println("Invalid year on row", i+1)
			return
		}

		// Column 9 contains the salary converted to US dollars.
		// Read the salary in US dollars.
		salary, err := strconv.ParseFloat(row[8], 64)
		if err != nil {
			fmt.Println("Invalid salary on row", i+1)
			return
		}

		if year == 2023 && level == "Entry-level" {
			entry2023 = append(entry2023, salary)
		}

		if year == 2024 {
			if level == "Entry-level" {
				entry2024 = append(entry2024, salary)

				if model == "Remote" {
					remoteEntry = append(remoteEntry, salary)
				}

				if model == "On-site" {
					onsiteEntry = append(onsiteEntry, salary)
				}

				// Leave research jobs out of the job title ranking.
				if !strings.Contains(strings.ToLower(title), "research") {
					titleSalaries[title] = append(
						titleSalaries[title], salary)
				}
			}

			if level == "Mid-level" {
				mid2024 = append(mid2024, salary)
			}

			if model == "Remote" {
				remoteAll = append(remoteAll, salary)
			}

			if model == "On-site" {
				onsiteAll = append(onsiteAll, salary)
			}
		}
	}

	if len(entry2024) == 0 || len(entry2023) == 0 ||
		len(mid2024) == 0 {
		fmt.Println("Missing data needed for the yearly comparisons.")
		return
	}

	color.Cyan("SALARY REPORT")
	fmt.Println("All salaries are annual amounts in US dollars.")
	fmt.Printf("2024 entry-level records: %d\n", len(entry2024))
	fmt.Printf("2023 entry-level records: %d\n", len(entry2023))

	// Part 1
	average2024 := average(entry2024)
	average2023 := average(entry2023)

	color.Cyan("\n1. Average Entry-Level Salaries")

	if average2024 < average2023 {
		color.Red("2024 average: $%.2f", average2024)
		color.Red("2023 average: $%.2f", average2023)
	} else {
		color.Green("2024 average: $%.2f", average2024)
		color.Green("2023 average: $%.2f", average2023)
	}

	// Part 2
	median2024 := median(entry2024)
	median2023 := median(entry2023)

	color.Cyan("\n2. Median Entry-Level Salaries")
	fmt.Printf("2024 median: $%.2f\n", median2024)
	fmt.Printf("2023 median: $%.2f\n", median2023)

	showDifference("2024 median compared to 2023",
		median2024, median2023)

	fmt.Println("\nAverage compared to median:")
	showDifference("2024 average compared to its median",
		average2024, median2024)
	showDifference("2023 average compared to its median",
		average2023, median2023)

	// Part 3
	midMedian := median(mid2024)

	color.Cyan("\n3. Entry-Level and Mid-Level Salaries in 2024")
	fmt.Printf("Entry-level median: $%.2f\n", median2024)
	fmt.Printf("Mid-level median: $%.2f\n", midMedian)

	showDifference("Mid-level compared to entry-level",
		midMedian, median2024)

	fmt.Println("This difference shows a possible salary change.")
	fmt.Println("An individual raise is not guaranteed.")

	// Part 4
	color.Cyan("\n4. Remote and Onsite Salaries in 2024")

	fmt.Println("\nEntry-level jobs:")
	showRemote(remoteEntry, onsiteEntry)

	fmt.Println("\nAll experience levels:")
	showRemote(remoteAll, onsiteAll)

	// Part 5
	color.Cyan("\n5. Top Five Entry-Level Job Titles in 2024")
	fmt.Println("Research titles are excluded. Ranked by median salary.")

	// Find the highest remaining title during each loop.
	for rank := 1; rank <= 5; rank++ {
		bestTitle := ""
		bestSalary := 0.0

		for title, salaries := range titleSalaries {
			titleMedian := median(salaries)

			if bestTitle == "" || titleMedian > bestSalary ||
				(titleMedian == bestSalary && title < bestTitle) {
				bestTitle = title
				bestSalary = titleMedian
			}
		}

		if bestTitle == "" {
			break
		}

		if rank == 1 {
			color.Green("%d. %s: $%.2f", rank, bestTitle, bestSalary)
		} else {
			color.Yellow("%d. %s: $%.2f", rank, bestTitle, bestSalary)
		}

		// Remove this title so the next loop finds the next highest.
		delete(titleSalaries, bestTitle)
	}
}

func average(salaries []float64) float64 {
	total := 0.0

	for _, salary := range salaries {
		total += salary
	}

	return total / float64(len(salaries))
}

func median(salaries []float64) float64 {
	slices.Sort(salaries)
	middle := len(salaries) / 2

	if len(salaries)%2 == 0 {
		return (salaries[middle-1] + salaries[middle]) / 2
	}

	return salaries[middle]
}

func showDifference(label string, first, second float64) {
	difference := first - second

	if difference > 0 {
		color.Green("%s: $%.2f higher.", label, difference)
	} else if difference < 0 {
		color.Red("%s: $%.2f lower.", label, -difference)
	} else {
		color.Yellow("%s: equal, with a $0.00 difference.", label)
	}
}

func showRemote(remote, onsite []float64) {
	if len(remote) == 0 || len(onsite) == 0 {
		fmt.Println("Not enough data for this comparison.")
		return
	}

	remoteMedian := median(remote)
	onsiteMedian := median(onsite)

	if onsiteMedian > remoteMedian {
		color.Red("Fully remote median: $%.2f", remoteMedian)
		color.Red("Onsite median: $%.2f", onsiteMedian)
	} else {
		color.Green("Fully remote median: $%.2f", remoteMedian)
		color.Green("Onsite median: $%.2f", onsiteMedian)
	}
}
