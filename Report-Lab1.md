##前言
之前做了MIT6.828操作系统课程，过程很痛苦，也学到了很多。不过最大的收货还是发现了非常好的学习方式——刷顶级名校的闻名课程。自然同样是麻省理工久负盛名的分布式系统课程MIT6.824就顺理成章的加入了todolist。值得议一提的是，这两门课程都是同一个实验室出品，实在是良心。

####[我的项目地址在:github.com/T0mmyliu/distributed-system-course-MIT6.824](https://github.com/T0mmyliu/distributed-system-course-MIT6.824)
####[课程网址在:https://pdos.csail.mit.edu/6.824/](https://pdos.csail.mit.edu/6.824/)

---

##MapReduce

MapReduce是谷歌的并行计算框架，据说是谷歌的工程师觉得每天处理分布式的代码，开锁解锁实在是又难又烦，而且有很多重复工作。就对并行的过程进行了抽象，就搞了这样一套简单易用的框架。

要想理解其中奥妙，当然是要阅读谷歌的那一篇顶顶有名的论文["MapReduce: Simplified Data Processing on Large Clusters"](https://pdos.csail.mit.edu/6.824/papers/mapreduce.pdf)，这篇论文也在课程的schedule里面，值得一提的是，6.824的schedule里有很多assignment，其中不少是论文，能入选的都是经典中的经典，一定不能跳过，况且很多都和lab有很大的关系。

###什么是MapReduce
虽然MapReduce是一个威力强大的工具，但是一般来说威力强大的工具最好也要简单好用，而把复杂封装在背后，才能称之为好工具。具体背后有多复杂，比如如何处理各类异常，我暂时不能领会。这里只介绍下这个工具如何简单好用！

lab1需要我们实现一个toy版本的MapReduce框架出来，虽然代码量比较少，但是难度也还是不小，下面先解释下这个mapreduce的逻辑。

###直观解释
MapReduce的逻辑实在是非常的直观。简单的说就分为两步:Map 和 Reduce 。

####用一个生活中的例子来举例说明：炒菜

假如你是一个酒店的行政总厨(Master)，你不负责做菜，只负责指挥下面的厨师Worker。每天都进行着如下的步骤：

####Register
每天都会有人去采购很多食材回来，有牛肉，猪肉，藕，土豆........买回来这些食物后，需要将他们注册在一个本子上，这就是今天要处理的任务(task)了！

####Partition and Ｍap
做菜的第一个过程，当然是切菜，你是主厨(Master)，你不负责做菜，只负责通知。”喂，过来把这个藕切成丁“，于是就有一个空闲的厨师过来切藕，你就把藕这个任务分发给了厨师(Worker)。切东西的过程，就是partition。之后可能还需要对食物做一些预先处理，比如抹上盐和料酒，这个过程就是mapfunc。

厨师切好藕后，根据一个规则”噢，藕是要炒肉的，放到第一个盆子里面去“，规则对于每一个厨师都是一样的，这样保证了，同一个input，处理(map)完后肯定都在一起。比如藕片肯定就都放在第一个盆里。一般来说如果有R个输出对象(盆子)，一个简单的分配规则就是hash后MOD R。

####Reduce
在介绍这个环节之前有一点要注意：和做菜一样，炒菜之前务必要保证菜已经切好了，这样才能不至于要炒回锅肉，发现只有肉片没有蒜苗。当确认前一个环节(map)确实结束后，开始炒菜(Reduce)。在这个环节里，你还是负责通知空闲的厨师，”喂你去炒一份回锅肉“。就有一个厨师接收到这个任务，开始准备炒肉。

”我需要煮熟的五花肉“说完从第一个盆子里拿了肉，"我还需要辣椒"，说完从第二个盆子里拿了辣椒。。。。“最后来一点蒜”从第R个盆子里拿了蒜。原料准备好了，就开始炒菜，这个过程就是reduce function，炒完了过后，统一输出到一个盘子里。这个环节一共要炒M盘菜(M个任务)，那么就有了M盘菜了。

####Merge
这个环节最好理解，在所有菜都做好了之后，统一端出去，给客人(输入)。

解释的比较粗略，具体还是要去读["MapReduce: Simplified Data Processing on Large Clusters"](https://pdos.csail.mit.edu/6.824/papers/mapreduce.pdf)。也可以参考这一副图理解一下：

![这里写图片描述](http://i.stack.imgur.com/EH8tY.jpg)

##LAB1

前两个部分是完成非分布式的版本

###Part I: Map/Reduce input and output

domap这个函数没有什么好说的，就是读入一个指定文件，跑一遍mapfunc，得到一个key-value的slice。然后将slice分发到R个文件里去。分发部分这样写：
```
		for _, kv := range mapRes {
			index := int(ihash(kv.Key)) % nReduce
			encoders[index].Encode(&kv)
		}
```

###Part II: Single-worker word count

单词计数，作为用户写MapFunction和ReduceFunction，体会一下MapReduce框架使用起来有多简洁。

###Part III: Distributing MapReduce tasks

终于到分布式的版本了，厉害了我的哥！

虽然需要完成的代码量很少，但是我确实花了很长一段时间才搞懂要干什么，然后又花了很长一段时间debug。

首先得通读6.824/src/mapreduce模块下的每一个.go文件，尤其是不能忘记了test_test.go文件，多次经验证明，测试文件是最好的说明文档！！！

读说明我们知道，我们需要通过的单元测试是TestBasic
```
    func TestBasic(t *testing.T) {
        mr := setup()
        for i := 0; i < 2; i++ {
            go RunWorker(mr.address, port("worker"+strconv.Itoa(i)),
                MapFunc, ReduceFunc, -1)
        }
        mr.Wait()
        check(t, mr.files)
        checkWorker(t, mr.stats)
        cleanup(mr)
    }
```
其中setup是
```
    func setup() *Master {
        files := makeInputs(nMap)
        master := port("master")
        mr := Distributed("test", files, nReduce, master)
        return mr
    }
```
第一个主要的函数调用是`Distributed`,位于master.go
```
    func Distributed(jobName string, files []string, nreduce int, master string) (mr *Master) {
        mr = newMaster(master)　　
        mr.startRPCServer()   //初始化RPC　server
        go mr.run(jobName, files, nreduce, mr.schedule, func() {
            mr.stats = mr.killWorkers()
            mr.stopRPCServer()
        })
        return
    }
```
然后就是需要去master_rpc.go里去看`mr.startRPCServer()`。

...省略dsp过程若干，给的代码中间的注释挺多，应该只需要时间问题。于是弄懂了要在schedule.go里写的东西是什么。

开始干活...(这一步很快...因为实验一的代码量很小)

提几个要注意地方，有些是语言的注意事项，因为我是刚刚学golang，所以遇到了一些问题。

####Q:channel何时会阻塞？
Ａ:当chan为空，又要取出(<-chan)的时候，以及如果chan满了，要push进去的时候。尤其是第二点。给的代码里的chan都是没有buffered的，所以只要chan里有正在等待的element，再push就会阻塞
```
    type Master struct {
	sync.Mutex

	address         string
	registerChannel chan string
	doneChannel     chan bool
	workers         []string // protected by the mutex

	// Per-task information
	jobName string   // Name of currently executing job
	files   []string // Input files
	nReduce int      // Number of reduce partitions

	shutdown chan struct{}
	l        net.Listener
	stats    []int
}
```

处理方法很简单，就是把普通的chan<-element放到goroutine里
```
    go func(){
        chan<-element
    }()
```

####Q:为什么我的代码能通过map,但是一到reduce，运行前两个文件就失败了？
A:我被这个bug卡了一晚上...就如前面说的，就如同做菜一样，你要保证所有原材料都已经就绪了，才能开始炒菜。同理，你要保证map的所有go　thread已经结束了，所有的中间文件已经生产了，才能开始下一步！！！
可以采用计数的方法，具体是用了`sync.WaitGroup`,具体如下:
```
	var wg sync.WaitGroup 

	for taskIndex := 0; taskIndex < ntasks; taskIndex++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
            dotask...........
		}(taskIndex)
	}
	wg.Wait()
	fmt.Printf("Schedule: %v phase done\n", phase)
```
###Part IV: Handling worker failures

任务失败情况的处理在这个lab里很简单...就是把没完成的任务重新发出去就可以了。具体用循环来做：分发任务－>任务失败，继续到开头／任务成功，break退出。

###最后写出来的schedule代码如下：

```
    package mapreduce

    import (
        "fmt"
        "sync"
    )

    // schedule starts and waits for all tasks in the given phase (Map or Reduce).
    func (mr *Master) schedule(phase jobPhase) {
        var ntasks int
        var nios int // number of inputs (for reduce) or outputs (for map)
        switch phase {
        case mapPhase:
            ntasks = len(mr.files)
            nios = mr.nReduce
        case reducePhase:
            ntasks = mr.nReduce
            nios = len(mr.files)
        }

        fmt.Printf("Schedule: %v %v tasks (%d I/Os)\n", ntasks, phase, nios)

        // All ntasks tasks have to be scheduled on workers, and only once all of
        // them have been completed successfully should the function return.
        // Remember that workers may fail, and that any given worker may finish
        // multiple tasks.
        //
        // TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO
        //
        taskArgsList := make([]DoTaskArgs, ntasks)
        for i := 0; i < ntasks; i++ {
            taskArgsList[i].JobName = mr.jobName
            taskArgsList[i].Phase = phase
            taskArgsList[i].TaskNumber = i
            taskArgsList[i].NumOtherPhase = nios
            taskArgsList[i].File = mr.files[i]
        }

        fmt.Printf("Schedule: taskArgsList prepared.\n")
        var wg sync.WaitGroup //

        for taskIndex := 0; taskIndex < ntasks; taskIndex++ {
            wg.Add(1)
            var idleWorker string

            go func(index int) {
                defer wg.Done()
                for {
                    idleWorker = <-mr.registerChannel
                    ok := call(idleWorker, "Worker.DoTask", &taskArgsList[index], new(struct{}))
                    if ok == false {
                        fmt.Printf("Master: RPC %s DoTask error\n", idleWorker)
                        continue
                    } else {
                        go func() {
                            mr.registerChannel <- idleWorker
                        }()
                        break
                    }
                }
            }(taskIndex)
        }
        wg.Wait()
        fmt.Printf("Schedule: %v phase done\n", phase)
    }
```


