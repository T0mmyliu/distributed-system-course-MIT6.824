package mapreduce

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
)

// doReduce does the job of a reduce worker: it reads the intermediate
// key/value pairs (produced by the map phase) for this task, sorts the
// intermediate key/value pairs by key, calls the user-defined reduce function
// (reduceF) for each key, and writes the output to disk.
func doReduce(
	jobName string, // the name of the whole MapReduce job
	reduceTaskNumber int, // which reduce task this is
	nMap int, // the number of map tasks that were run ("M" in the paper)
	reduceF func(key string, values []string) string,
) {
	// TODO:
	// You will need to write this function.
	// You can find the intermediate file for this reduce task from map task number
	// m using reduceName(jobName, m, reduceTaskNumber).
	// Remember that you've encoded the values in the intermediate files, so you
	// will need to decode them. If you chose to use JSON, you can read out
	// multiple decoded values by creating a decoder, and then repeatedly calling
	// .Decode() on it until Decode() returns an error.
	//
	// You should write the reduced output in as JSON encoded KeyValue
	// objects to a file named mergeName(jobName, reduceTaskNumber). We require
	// you to use JSON here because that is what the merger than combines the
	// output from all the reduce tasks expects. There is nothing "special" about
	// JSON -- it is just the marshalling format we chose to use. It will look
	// something like this:
	//
	// enc := json.NewEncoder(mergeFile)
	// for key in ... {
	// 	enc.Encode(KeyValue{key, reduceF(...)})
	// }
	// file.Close()

	keyValues := make(map[string][]string)

	i := 0
	for i < nMap {
		fileName := reduceName(jobName, i, reduceTaskNumber)

		file, err := os.Open(fileName)
		if err != nil {
			log.Fatal(err)
		}

		enc := json.NewDecoder(file)
		for {
			var kv KeyValue
			err := enc.Decode(&kv)
			if err != nil {
				break
			}
			_, ok := keyValues[kv.Key]
			if !ok {
				keyValues[kv.Key] = make([]string, 0)
			}
			keyValues[kv.Key] = append(keyValues[kv.Key], kv.Value)
		}
		i++
	}

	var keys []string
	for k, err := range keyValues {
		keys = append(keys, k)
		if err != nil {
			log.Fatal(err)
		}
	}

	sort.Strings(keys)

	file, err := os.Create(mergeName(jobName, reduceTaskNumber))
	if err != nil {
		fmt.Printf("reduce merge file:%s can't open\n", mergeName(jobName, reduceTaskNumber))
		return
	}
	enc := json.NewEncoder(file)

	for _, k := range keys {
		enc.Encode(KeyValue{k, reduceF(k, keyValues[k])})
	}
	file.Close()

}
