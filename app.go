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

type DiskTree struct {
	Name     string `json:"name"`
	Children []any  `json:"children"`
}

type DiskTreeResult struct {
	Tree   *DiskTree
	Status string
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
	dirAry, err := os.ReadDir(path)
	if err != nil {
		return
	}
	for _, e := range dirAry {
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

func (a *App) DiskTreeMapStatistics(path string) DiskTree {
	pool, _ := ants.NewPool(10000)
	defer pool.Release()
	start := time.Now()
	var rootDir string
	if path == "root" {
		rootDir = "C:/"
	} else {
		rootDir = path
	}
	d, _ := disk.Usage(rootDir)
	diskTree := DiskTree{Name: rootDir}
	diskTree.Children = make([]any, 0, 30)
	diskTree.Children = append(diskTree.Children, Node{Name: "Free", Value: int(d.Free)})
	paths, _ := os.ReadDir(rootDir)
	var (
		loopWg sync.WaitGroup
		info   os.FileInfo
		err    error
	)
	for _, p := range paths {
		var wait sync.WaitGroup
		loopWg.Add(1)
		wait.Add(1)
		if info, err = p.Info(); err != nil {
			continue
		}
		if p.IsDir() {
			node := &DirNode{Name: p.Name(), Children: make([]any, 0, 50)}
			diskTree.Children = append(diskTree.Children, node)
			_ = pool.Submit(func() {
				scanDirTwo(filepath.Join(rootDir, p.Name()), &wait, node, pool)
				wait.Done()
			})
		} else {
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

func (a *App) StartSan(path string) DiskTree {
	return a.DiskTreeMapStatistics(path)

}
