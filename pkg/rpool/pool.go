package rpool

import (
	"errors"
	"github.com/gookit/slog"
	"io"
	"sync"
)

//https://www.jianshu.com/p/397b97e0943b

// Pool 管理一组可以安全地在多个goroutine间共享得到资源。被管理的资源必须实现io.Closer接口
type Pool struct {
	m         sync.Mutex
	resources chan io.Closer
	factory   func() (io.Closer, error)
	closed    bool
}

var ErrPoolClosed = errors.New("pool has been closed.")

//New创建一个用来管理资源的池，这个池需要一个可以分配新资源的函数，并规定池的大小
func New(fn func() (io.Closer, error), size uint) (*Pool, error) {
	if size < 0 {
		return nil, errors.New("size value too small")
	}
	return &Pool{
		factory:   fn,
		resources: make(chan io.Closer, size),
	}, nil
}

//Acquire从池中获取一个资源
func (p *Pool) Acquire() (io.Closer, error) {
	select {
	case r, ok := <-p.resources: //检查是否有空闲的资源
		slog.Println("ACquire:", "shared resource")
		if !ok {
			return nil, ErrPoolClosed
		}
		return r, nil
	default:
		slog.Println("Acquire:", "New resource")
		return p.factory()
	}
}

//Release将一个使用后的资源池释放回池里
func (p *Pool) Release(r io.Closer) {
	p.m.Lock() //确保本操作和Close操作安全
	defer p.m.Unlock()
	if p.closed { //如果资源池已经关闭，销毁这个资源
		r.Close()
		return
	}
	select {
	case p.resources <- r: //试图将这个资源放入队列
		slog.Println("Release:", "In Queue")
	default: //如果队列已满，则关闭这个资源
		slog.Println("Release", "Closing")
		r.Close()
	}
}

//close会让资源池停止工作，并关闭所有的资源
func (p *Pool) Close() {
	p.m.Lock() //确保本操作与release操作安全
	defer p.m.Unlock()
	if p.closed {
		return
	}
	p.closed = true    //将池关闭
	close(p.resources) //在清空通道里的资源前，将通道关闭，如果不这样做，会发生死锁
	//关闭资源
	for r := range p.resources {
		r.Close()
	}
}
