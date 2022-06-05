package mr

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"time"
)

type TaskStatus struct {
    start_time time.Time
    
    // Status code;
    // 0 not assigned
    // 1 in progress
    // 2 succeed
    // 3 failed (time up)
    status int 
}

type Master struct {
	// Your definitions here.
    files []string,
    mapper_inter_ids []int,
    nReduce int,

    global_mapper_id int,
    global_reducer_id int,

    mapper_status map[string]TaskStatus
    reducer_status map[int]TaskStatus


}

// Your code here -- RPC handlers for the worker to call.
func (m *Master) GetTask(args *GetTaskArgs, reply *GetTaskReply) error {
    reply := GetTaskReply{}
    for file, status := range mapper_status {
        if status.status == 0 | status.status == 3 {
            global_mapper_id++
            reply.worker_id = global_mapper_id
            reply.task_type = 1
            reply.file_name = file
            reply.reducer = m.Reduce

            mapper_status[file].start_time = time.Now()
            mapper_status[file].status = 1
            return nil
        }
    }

    for reducer_id, status := range reducer_status {
        if status.status == 0 
          | status.status == 3 
          | time.Now() - status.start_time > 10 {
            reply.worker_id = reducer_id
            reply.task_type = 2
            reply.intermediate_id = m.mapper_inter_ids
            
            reducer_status[reducer_id].start_time = time.Now()
            reducer_status[reducer_id].status = 1
            return nil
        }
    }

    reply.worker_id = 0
    reply.task_type = 0
    return nil
}

//
// an example RPC handler.
//
// the RPC argument and reply types are defined in rpc.go.
//
func (m *Master) Example(args *ExampleArgs, reply *ExampleReply) error {
	reply.Y = args.X + 1
	return nil
}


//
// start a thread that listens for RPCs from worker.go
//
func (m *Master) server() {
	rpc.Register(m)
	rpc.HandleHTTP()
	//l, e := net.Listen("tcp", ":1234")
	sockname := masterSock()
	os.Remove(sockname)
	l, e := net.Listen("unix", sockname)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}

//
// main/mrmaster.go calls Done() periodically to find out
// if the entire job has finished.
//
func (m *Master) Done() bool {
	ret := false

	// Your code here.


	return ret
}

//
// create a Master.
// main/mrmaster.go calls this function.
// nReduce is the number of reduce tasks to use.
//
func MakeMaster(files []string, nReduce int) *Master {
	m := Master{}

	// Your code here.


	m.server()
	return &m
}
