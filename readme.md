# Project 2 - Salary Analysis
Shabaan Saleem

## How to run
Make sure salaries.csv is in the same folder as main.go.
Open the terminal in that folder and type:

go run .

The program uses the fatih/color library for the colors.
If the dependencies are missing, run go mod download first.

## What the program shows
The program reads salaries.csv and compares salaries in US dollars.

It shows:
1. Average entry-level salaries for 2024 and 2023.
2. Median entry-level salaries for both years, their difference,
   and how each year's average compares to its median.
3. The difference between entry-level and mid-level median
   salaries in 2024.
4. Remote and onsite median salaries in 2024, for entry-level
   jobs and for all experience levels.
5. The five highest paying entry-level job titles in 2024,
   based on median salary. Research titles are left out.

The average is the total salary divided by the number of records.
The median is the middle salary after sorting. If there are two
middle salaries, the program takes their average.

## What the colors mean
Cyan is used for headings.

For the yearly averages, red means the 2024 average is lower.
Otherwise, the averages are green.

For differences, red means lower, green means higher,
and yellow means equal.

For remote and onsite salaries, both values are red when
onsite pays more. Otherwise, they are green.

The top job title is green, and the next four are yellow.

## Notes
The results only describe the data in the CSV file.
Moving to a mid-level job does not guarantee the raise shown.

## Unfinished parts
None.