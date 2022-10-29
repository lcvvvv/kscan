package pool

import (
	"log"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

var logger = Logger(log.New(os.Stderr, "[pool]", log.Ldate|log.Ltime))

type Logger interface {
	Println(...interface{})
	Printf(string, ...interface{})
}

func SetLogger(log Logger) {
	logger = log
}

//创建worker，每一个worker抽象成一个可以执行任务的函数
type Worker struct {
	f func(interface{})
}

//通过NewTask来创建一个worker
func generateWorker(f func(interface{})) *Worker {
	return &Worker{
		f: func(in interface{}) {
			f(in)
		},
	}
}

//执行worker
func (t *Worker) run(in interface{}) {
	t.f(in)
}

//池
type Pool struct {
	//母版函数
	Function func(interface{})
	//Pool输入队列
	in chan interface{}
	//size用来表明池的大小，不能超发。
	threads int
	//启动协程等待时间
	Interval time.Duration
	//正在执行的任务清单
	JobsList *sync.Map
	//正在工作的协程数量
	length int32
	//用于阻塞
	wg *sync.WaitGroup
	//提前结束标识符
	Done bool
}

//实例化工作池使用
func New(threads int) *Pool {
	return &Pool{
		threads:  threads,
		JobsList: &sync.Map{},
		wg:       &sync.WaitGroup{},
		in:       make(chan interface{}),
		Function: nil,
		Done:     true,
		Interval: time.Duration(0),
	}
}

//结束整个工作
func (p *Pool) Push(i interface{}) {
	if p.Done {
		return
	}
	p.in <- i
}

//结束整个工作
func (p *Pool) Stop() {
	if p.Done != true {
		close(p.in)
	}
	p.Done = true
}

//执行工作池当中的任务
func (p *Pool) Run() {
	p.Done = false
	//只启动有限大小的协程，协程的数量不可以超过工作池设定的数量，防止计算资源崩溃
	for i := 0; i < p.threads; i++ {
		p.wg.Add(1)
		time.Sleep(p.Interval)
		go p.work()
		if p.Done == true {
			break
		}
	}
	p.wg.Wait()
}

//从jobs当中取出任务并执行。
func (p *Pool) work() {
	var Tick string
	var param interface{}
	//减少waitGroup计数器的值
	defer func() {
		defer func() {
			if e := recover(); e != nil {
				logger.Println(Tick, param, e)
			}
		}()
		p.wg.Done()
	}()
	for param = range p.in {
		if p.Done {
			return
		}
		atomic.AddInt32(&p.length, 1)
		//获取任务唯一票据
		Tick = p.generateTick()
		//压入工作任务到工作清单
		p.JobsList.Store(Tick, param)
		//设置工作内容
		f := generateWorker(p.Function)
		//开始工作，输出工作结果
		f.run(param)
		//工作结束，删除工作清单
		p.JobsList.Delete(Tick)
		atomic.AddInt32(&p.length, -1)
	}
}

//生成工作票据
func (p *Pool) generateTick() string {
	rand.Seed(time.Now().UnixNano())
	return strconv.FormatInt(rand.Int63(), 10)
}

//获取线程数
func (p *Pool) Threads() int {
	return p.threads
}

func (p *Pool) RunningThreads() int {
	return int(p.length)
}
