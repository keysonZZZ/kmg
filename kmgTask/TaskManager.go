package kmgTask

/*
任务管理
有很多任务,可以开线程同时处理这些任务.任务管理器可以等待所有任务完成任务完成之后
*/
type Task interface {
	//运行任务(同步),在任务运行完成后或被终止时返回.返回运行结果
	Run()
	//在运行中途,结束任务(同步),在成功终止任务后返回.
	Stop()
}

//任务管理器,可以添加任务,可以等待所有添加的任务结束,
//这个类是并发安全的.
type TaskManager interface {
	//添加任务(异步,或同步),立刻返回,任务会在这个调用开始之后的某个时间开始运行.
	//在某些情况下,这个函数会被阻塞(排队任务过多..)
	AddTask(t Task)

	//等待所有已经添加的任务都完成
	Wait()

	//1.等待所有的任务都完成
	//2.释放某些资源
	//关闭管理器(同步),在所有已添加任务完成后返回,
	//关闭后添加任务,会panic
	//二次关闭会panic
	Close()
}

//用一个函数定义的任务(不能中途结束)
type TaskFunc func()

func (f TaskFunc) Run() {
	f()
}
func (f TaskFunc) Stop() {
}