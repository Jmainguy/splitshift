package main

import (
	"encoding/csv"
	"fmt"
	"mime/multipart"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func percentOfHour(minutes string) string {
	if minutes == "15" {
		return "25"
	} else if minutes == "30" {
		return "50"
	} else if minutes == "45" {
		return "75"
	} else if minutes == "00" {
		return "00"
	}
	return "00"
}

func timeToFloat(shiftTime string) (shiftTimeFloat float64, err error) {
	shiftString := strings.Split(shiftTime, " ")[0]
	shiftArray := strings.Split(shiftString, ":")
	shiftHour := shiftArray[0]
	shiftMinute := shiftArray[1]
	shiftMinutePercentage := percentOfHour(shiftMinute)
	shiftTimeString := shiftHour + "." + shiftMinutePercentage
	shiftTimeFloat, err = strconv.ParseFloat(shiftTimeString, 64)
	if err != nil {
		return shiftTimeFloat, err
	}
	return shiftTimeFloat, err

}

func processFile(f multipart.File) (resultLines []string, err error) {

	// Read CSV File into a Variable
	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return []string{}, err
	}

	var monthlyTotal float64
	var monthlyHours float64
	//timeRegex, err := regexp.Compile(`^([0-1]?[0-9]|2[0-3]):[0-5][0-9]$`)
	timeRegex, err := regexp.Compile(`((1[0-2]|0?[1-9]):([0-5][0-9]) ?([AaPp][Mm]))`)
	if err != nil {
		return []string{}, err
	}
	// Loop through CSV lines & turn into object
	for _, line := range lines {
		shift := &Shift{}
		match := timeRegex.MatchString(line[4])
		if match {
			match = timeRegex.MatchString(line[5])
			if match {
				split := false
				if strings.Contains(line[4], "PM") {
					if strings.Contains(line[5], "AM") {
						split = true
					}
				}

				shift.AideName = line[0]
				rateString := line[7]
				shift.StartTime = line[4]
				shift.EndTime = line[5]
				shift.Date = line[3]
				//fmt.Println(line)
				shift.Hours, err = strconv.ParseFloat(line[6], 32)
				if err != nil {
					return []string{}, err
				}
				rateString = strings.Trim(rateString, "$")
				rateString = strings.Replace(rateString, " ", "", -1)
				shift.Rate, err = strconv.ParseFloat(rateString, 32)
				if err != nil {
					return []string{}, err
				}
				shift.Multiplier, err = strconv.ParseFloat(line[10], 32)
				if err != nil {
					return []string{}, err
				}
				shift.Amount = shift.Rate * shift.Hours * shift.Multiplier
				// Add to monthly totals
				monthlyTotal = monthlyTotal + shift.Amount
				monthlyHours = monthlyHours + shift.Hours
				// Split or not
				if !split {
					resultLine := fmt.Sprintf("AideName: %v, Date: %v, Start: %v, End: %v, Hours: %.2f, Rate: $%.2f, Total: $%.2f\n", shift.AideName, shift.Date, shift.StartTime, shift.EndTime, shift.Hours, shift.Rate, shift.Amount)
					resultLines = append(resultLines, resultLine)
				} else if shift.EndTime == "12:00 AM" {
					resultLine := fmt.Sprintf("AideName: %v, Date: %v, Start: %v, End: %v, Hours: %.2f, Rate: $%.2f, Total: $%.2f\n", shift.AideName, shift.Date, shift.StartTime, shift.EndTime, shift.Hours, shift.Rate, shift.Amount)
					resultLines = append(resultLines, resultLine)
				} else {
					// Split Shift
					shift.Hours, err = strconv.ParseFloat(line[6], 32)
					if err != nil {
						return []string{}, err
					}
					// Check time
					endTime, err := timeToFloat(shift.EndTime)
					if err != nil {
						return []string{}, err
					}

					shiftOneHours := shift.Hours - endTime
					shiftOneTotal := shiftOneHours * shift.Rate * shift.Multiplier
					shiftTwoTotal := endTime * shift.Rate * shift.Multiplier
					// Sanity Check, hours and total should match
					bothShiftHours := shiftOneHours + endTime

					bothShiftHoursString := fmt.Sprintf("%.2f", bothShiftHours)
					generationHours, err := strconv.ParseFloat(line[6], 32)
					if err != nil {
						return []string{}, err
					}

					if bothShiftHoursString != fmt.Sprintf("%.2f", generationHours) {
						return []string{}, fmt.Errorf("Hours did not match, %v and %v, Line in question is %v", bothShiftHoursString, fmt.Sprintf("%.2f", generationHours), line)
					}
					shiftOneDate, err := time.Parse("01/02/2006", shift.Date)
					if err != nil {
						shiftOneDate, err = time.Parse("01/2/2006", shift.Date)
						if err != nil {
							return []string{}, err
						}
					}
					shiftTwoDate := shiftOneDate.AddDate(0, 0, 1)

					// Payoff
					resultLine := fmt.Sprintf("AideName: %v, Date: %v, Start: %v, End: %v, Hours: %.2f, Rate: $%.2f, Total: $%.2f\n", shift.AideName, shiftOneDate.Format("01/02/2006"), shift.StartTime, "11:59 PM", shiftOneHours, shift.Rate, shiftOneTotal)
					resultLineTwo := fmt.Sprintf("AideName: %v, Date: %v, Start: %v, End: %v, Hours: %.2f, Rate: $%.2f, Total: $%.2f\n", shift.AideName, shiftTwoDate.Format("01/02/2006"), "12:00 AM", shift.EndTime, endTime, shift.Rate, shiftTwoTotal)
					resultLines = append(resultLines, resultLine)
					resultLines = append(resultLines, resultLineTwo)
				}
			}
		}
	}
	totals := fmt.Sprintf("Total Hours: %.2f, Total Amount $%.2f", monthlyHours, monthlyTotal)
	fmt.Println(totals)
	resultLines = append(resultLines, totals)

	return resultLines, err
}

func main() {
	http.HandleFunc("/", upload)

	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		fmt.Println(err)
	}

}
