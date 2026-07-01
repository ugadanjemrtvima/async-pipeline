package main

import (
	"slices"
	"strconv"
	"strings"
	"sync"
)

var md5Mutex = &sync.Mutex{}

func ExecutePipeline(jobs ...job) {
	var wg sync.WaitGroup
	in := make(chan interface{})
	for _, worker := range jobs {
		out := make(chan interface{})
		wg.Add(1)
		go func(job job, in chan interface{}, out chan interface{}) {
			defer wg.Done()
			defer close(out)
			job(in, out)
		}(worker, in, out)
		in = out
	}
	wg.Wait()
}

func SingleHash(in chan interface{}, out chan interface{}) {
	var wg sync.WaitGroup
	for data := range in {
		wg.Add(1)
		dataInt := data.(int)
		dataStr := strconv.Itoa(dataInt)
		go func(data string, out chan interface{}, mu *sync.Mutex, wg *sync.WaitGroup) {
			defer wg.Done()
			crc32Chan := make(chan string)
			go func() { //
				crc32Chan <- DataSignerCrc32(data)//
			}()//
			mu.Lock()
			md5base := DataSignerMd5(data)
			mu.Unlock()
			crc32Md5 := DataSignerCrc32(md5base)
			crc32Data := <-crc32Chan
			all := crc32Data + "~" + crc32Md5
			out <- all
		}(dataStr, out, md5Mutex, &wg)
	}
	wg.Wait()
}

func MultiHash(in chan interface{}, out chan interface{}) {
	var wg sync.WaitGroup
	for data := range in {
		wg.Add(1)
		dataStr := data.(string)
		go func(data string, out chan interface{}, wg *sync.WaitGroup) {
			defer wg.Done()
			var wg1 sync.WaitGroup
			resSlice := make([]string, 6)
			for i := 0; i < 6; i++ {
				wg1.Add(1)
				go func(data string, wg *sync.WaitGroup, num int, slice []string) {
					defer wg.Done()
					trueNum := strconv.Itoa(num)
					crc32 := DataSignerCrc32(trueNum + data)
					slice[num] = crc32
				}(data, &wg1, i, resSlice)
			}
			wg1.Wait()
			resStr := strings.Join(resSlice, "")
			out <- resStr
		}(dataStr, out, &wg)
	}
	wg.Wait()
}

func CombineResults(in chan interface{}, out chan interface{}) {
	var data []string
	for rawData := range in {
		dataStr := rawData.(string)
		data = append(data, dataStr)
	}
	slices.Sort(data)
	out <- strings.Join(data, "_")
}
