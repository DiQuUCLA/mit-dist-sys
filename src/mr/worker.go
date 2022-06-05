package mr

import (
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"log"
	"net/rpc"
	"os"
	"strings"
	"unicode"
)

//
// Map functions return a slice of KeyValue.
//
type KeyValue struct {
	Key   string
	Value string
}

//
// use ihash(key) % NReduce to choose the reduce
// task number for each KeyValue emitted by Map.
//
func ihash(key string) int {
	h := fnv.New32a()
	h.Write([]byte(key))
	return int(h.Sum32() & 0x7fffffff)
}

//
// Mapper worker
//
func Mapping(file_name string, worker_id int, reducer int) {
    f, _ := os.Open(file_name)
    content, _ := ioutil.ReadAll(f)
    kvs := mapf(file_name, content)
    
    out_fs := make([]File*, 10)
    for i := 0; i < reducer; i++ {
        out_file_name := fmt.Sprintf("mr_%d_%d", worker_id, i)
        out_fs[i], err = os.Create(out_file_name)
    }

    for i < len(kvs) {
        out_index = ihash(kvs[i].Key) % reducer
        fmt.Fprintf(out_fs[out_index], "%v %v\n", kvs[i].Key, kvs[i].Value)
    }

    for _, out_f := range out_fs {
        out_f.Close()
    }
}

//
// Reducer worker
//
func Reducing(intermidate_ids []int, worker_id int) {
    intermiate := []mr.KeyValue{}
    for i, id := range intermidate_ids {
        file_name = fmt.Fprintf("mr_%d_%d", i, worker_id)
        if err != nil {
            log.Fatalf("cannot open %v", file_name)
        }
        content, err := ioutil.ReadAll(file_name)
        if err != nil {
            log.Fatalf("cannot read %v", filename)
        }
        file.Close()
        s := string(content)
        s_list = strings.Split(s, "\n")
        for _, w := range words {
            w_list = strings.Split(w, " ")
            kv := mr.KeyValue{w_list[0], "1"}
            intermiate = append(intermiate, kv)
        }
    }
    sort.Sort(ByKey(intermiate))

    oname := fmt.Fprintf("mr-out-%d", worker_id)
    ofile, _ := os.Create(oname)

    i := 0
    for i < len(intermiate) {
        j := i + 1
        for j < len(intermiate) && intermiate[j].Key == intermiate[i].Key {
            j++
        }
        values := []string{}
        for k := i; k < j; k++ {
            values = append(values, intermiate[k].Value)
        }
        output := reducef(intermiate[i].Key, values)

        fmt.Fprintf(ofile, "%v %v\n", intermiate[i].Key, output)

        i = j
    }

    ofile.Close()
}

//
// main/mrworker.go calls this function.
//
func Worker(mapf func(string, string) []KeyValue,
	reducef func(string, []string) string) {

	// Your worker implementation here.
    reply := CallGetTask(0, 0)
	// uncomment to send the Example RPC to the master.
	// CallExample()
    worker_id := reply.worker_id
    task_type := reply.task_type
    
    for task_type != 2 {
        if(task_type == 0) {
            //Mapper
            Mapping(reply.file_name, reply.worker_id, reply.reducers)                   
        }else if(task_type == 1) {
            //Reducer
            Reducing(reply.intermidate_ids, reply.worker_id)
        }
    }

    reply = CallGetTask(worker_id, 1)
    
}

//
// example function to show how to make an RPC call to the master.
//
// the RPC argument and reply types are defined in rpc.go.
//
func CallExample() {

	// declare an argument structure.
	args := ExampleArgs{}

	// fill in the argument(s).
	args.X = 99

	// declare a reply structure.
	reply := ExampleReply{}

	// send the RPC request, wait for the reply.
	call("Master.Example", &args, &reply)

	// reply.Y should be 100.
	fmt.Printf("reply.Y %v\n", reply.Y)
}

func CallGetTask(worker_id int, worker_status int) GetTaskReply {
    args := GetTaskArgs{}
    args.worker_id = worker_id
    args.worker_status = worker_status

    reply := GetTaskReply{}

    call("Master.GetTask", &args, &reply)

    return reply
}

func CallDone(worker_id int, worker_status int) GetTaskReply {
    args := GetTaskArgs{}
    args.worker_id = worker_id
    args.worker_status = worker_status

    reply := GetTaskReply{}

    call("Master.Done", &args, &reply)

    return reply
}
//
// send an RPC request to the master, wait for the response.
// usually returns true.
// returns false if something goes wrong.
//
func call(rpcname string, args interface{}, reply interface{}) bool {
	// c, err := rpc.DialHTTP("tcp", "127.0.0.1"+":1234")
	sockname := masterSock()
	c, err := rpc.DialHTTP("unix", sockname)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	defer c.Close()

	err = c.Call(rpcname, args, reply)
	if err == nil {
		return true
	}

	fmt.Println(err)
	return false
}
