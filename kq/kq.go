// Copyright 2020 Longxiao Zhang <zhanglongx@gmail.com>.
// All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package kq

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// URLT is the template of KQ
const URLT = "http://kq.oa.sumavision.com/Users.aspx?act=&Company=技术公司&bm=230&WorkerNO=%d&sttim=%s&entim=%s&orderby=0&pagecount=31&intoExcel=true"

// KQ *MUST* be initialized with all members
type KQ struct {
	StartDate time.Time
	EndDate   time.Time

	Workers []int
}

// Info contains a worker's all info
type Info struct {
	Days1 int
	Days2 int
}

// query should be initialized with startDate, endDate
type query struct {
	startDate time.Time
	endDate   time.Time
}

type qData struct {
	// TODO: add more, 上班时间, 下班时间, 出勤状态, 出勤时长
	duration float64
}

// query reads online URL, and return a map as key is row-date
func (q *query) query(w int) (map[time.Time]qData, error) {

	y, m, d := q.startDate.Date()
	sttim := fmt.Sprintf("%d-%02d-%02d", y, m, d)

	y, m, d = q.endDate.Date()
	entim := fmt.Sprintf("%d-%02d-%02d", y, m, d)

	url := fmt.Sprintf(URLT, w, sttim, entim)
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	r := csv.NewReader(res.Body)
	r.Comma = ','

	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	result := make(map[time.Time]qData)
	for _, row := range records {
		date, err := time.Parse("2006-1-2", row[1])
		if err != nil {
			continue
		}

		dur, err := strconv.ParseFloat(row[8], 64)
		if err != nil {
			continue
		}

		result[date] = qData{
			duration: dur,
		}
	}

	return result, nil
}

// Run runs a quering, and return a map of Info as Key
// is worker's ID
func (k *KQ) Run() (map[int]Info, error) {

	q := query{
		startDate: k.StartDate,
		endDate:   k.EndDate,
	}

	result := make(map[int]Info)
	for _, w := range k.Workers {
		data, err := q.query(w)
		if err != nil {
			continue
		}

		info := Info{}
		for _, row := range data {
			if row.duration >= 12.5 {
				info.Days2++
			} else if row.duration >= 10.5 {
				info.Days1++
			}
		}

		result[w] = info
	}

	return result, nil
}
