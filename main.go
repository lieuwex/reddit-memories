/*
Shows an older reddit post or comment on this date
Copyright © 2017 Lieuwe Rooijakkers

This library is free software; you can redistribute it and/or modify
it under the terms of the GNU Lesser General Public License as published
by the Free Software Foundation; either version 3 of the License, or
(at your option) any later version.

This library is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU Lesser General Public License for more details.

You should have received a copy of the GNU Lesser General Public License
along with this library; if not, see <http://www.gnu.org/licenses/>.
*/

package main

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"reddit-memories/geddit"
)

func trunc(t time.Time) time.Time {
	return t.Truncate(24 * time.Hour)
}

func dateMonthMatch(a time.Time, b time.Time) bool {
	return a.Month() <= b.Month() && a.Day() <= b.Day()
}

func fetchAllSubmissions(session *geddit.LoginSession) ([]*geddit.Submission, error) {
	res := make([]*geddit.Submission, 0)
	after := ""

	for {
		fetch, err := session.MyOverview(geddit.NewSubmissions, after)
		if err != nil {
			return nil, err
		}

		if len(fetch) == 0 {
			break
		}

		after = fetch[len(fetch)-1].FullID
		res = append(res, fetch...)
	}

	return res, nil
}

// Please don't handle errors this way.
func main() {
	// Login to reddit
	session, _ := geddit.NewLoginSession(
		"xxxx",
		"xxxx",
		"reddit-memories v1",
	)

	submissions, _ := fetchAllSubmissions(session)
	count := len(submissions)

	today := trunc(time.Now().UTC())
	index := sort.Search(count, func(i int) bool {
		submission := submissions[i]
		time := time.Unix(int64(submission.DateCreated), 0)
		date := trunc(time)
		return date.Year() != today.Year() && dateMonthMatch(date, today)
	})
	if index == count {
		return
	}

	submission := submissions[index]
	submissionTime := time.Unix(int64(submission.DateCreated), 0)
	yearDiff := today.Year() - submissionTime.Year()

	plural := "s"
	if yearDiff == 1 || yearDiff == -1 {
		plural = ""
	}

	str := fmt.Sprintf("%d year%s today:\n", yearDiff, plural)
	if submission.Title != "" {
		str += submission.Title + "\n"
	}
	if submission.Selftext != "" {
		str += submission.Selftext
	} else if submission.URL != "" {
		str += submission.URL
	} else if submission.Body != "" {
		str += submission.Body
	}
	fmt.Printf("%s\n", strings.Trim(str, " \n"))
}
