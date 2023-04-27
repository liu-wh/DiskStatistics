package main

import (
	"context"
	"fmt"
	"github.com/panjf2000/ants/v2"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/disk"
)

var (
	RootPath string
)

type DiskTree struct {
	Name     string `json:"name"`
	Children []any  `json:"children"`
}

type DirNode struct {
	Children []any  `json:"children"`
	Name     string `json:"name"`
}

type Node struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

func scanDirTwo(path string, wait *sync.WaitGroup, node *DirNode, pool *ants.Pool) {
	//读取目录下的文件信息
	dirAry, err := os.ReadDir(path)
	if err != nil {
		return
	}
	for _, e := range dirAry {
		//是目录 递归调用,并且创建一个硬盘树节点
		if e.IsDir() {
			wait.Add(1)
			childNode := &DirNode{Name: e.Name(), Children: make([]any, 0, 100)}
			node.Children = append(node.Children, childNode)
			_ = pool.Submit(
				func() {
					scanDirTwo(filepath.Join(path, e.Name()), wait, childNode, pool)
					wait.Done()
				})
		} else {
			//是文件, 存入到硬盘树节点
			info, err := e.Info()
			if err != nil {
				continue
			}
			childNode := &Node{Name: e.Name(), Value: int(info.Size())}
			node.Children = append(node.Children, childNode)
		}

	}
}

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) DiskTreeMapStatistics() DiskTree {

	//初始化线程池
	pool, _ := ants.NewPool(1000)
	defer pool.Release()
	start := time.Now()

	//初始化一颗硬盘树
	diskTree := DiskTree{Name: RootPath}
	paths, _ := os.ReadDir(RootPath)
	diskTree.Children = make([]any, 0, len(paths))

	//获取根目录的剩余空间
	d, _ := disk.Usage(RootPath)
	diskTree.Children = append(diskTree.Children, Node{Name: "可用空间", Value: int(d.Free)})
	var (
		loopWg sync.WaitGroup
		info   os.FileInfo
		err    error
	)

	//开始扫描根目录下的目录
	for _, p := range paths {
		var wait sync.WaitGroup
		loopWg.Add(1)
		wait.Add(1)
		if info, err = p.Info(); err != nil {
			continue
		}
		//如果是目录, 使用另外一个函数开始扫描
		if p.IsDir() {
			node := &DirNode{Name: p.Name(), Children: make([]any, 0, 50)}
			diskTree.Children = append(diskTree.Children, node)
			_ = pool.Submit(func() {
				scanDirTwo(filepath.Join(RootPath, p.Name()), &wait, node, pool)
				wait.Done()
			})
		} else {
			//不是目录的话, 把文件信息存到硬盘树
			wait.Done()
			node := &Node{Name: p.Name(), Value: int(info.Size())}
			diskTree.Children = append(diskTree.Children, node)
		}
		wait.Wait()
		loopWg.Done()
	}
	loopWg.Wait()
	log.Printf("扫描结束, 用时%s秒\n", fmt.Sprintf("%0.2f", time.Since(start).Seconds()))
	return diskTree
}

// StartSan 前端点击扫描, 就会调用这个函数
func (a *App) StartSan() DiskTree {
	return a.DiskTreeMapStatistics()

}
