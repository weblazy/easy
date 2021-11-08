package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/tal-tech/go-zero/core/collection"
	"github.com/weblazy/easy/utils/list"
	"github.com/weblazy/easy/utils/threading"
	"github.com/weblazy/easy/utils/timex"
)

const (
	drainWorkers           = 8
	defaultMaxGoroutineNum = 16
)

type (
	// Execute defines the method to execute the task.
	Execute func(key, value interface{})

	// A TimingWheel is a timing wheel object to schedule tasks.
	TimingWheel struct {
		interval      time.Duration
		ticker        timex.Ticker
		slots         []map[int]*list.List
		timers        *collection.SafeMap
		tickedPos     int
		numSlots      int
		execute       Execute
		setChannel    chan timingEntry
		moveChannel   chan baseEntry
		removeChannel chan interface{}
		drainChannel  chan func(key, value interface{})
		stopChannel   chan struct{}

		cacheList       *list.List
		taskList        *list.List
		locked          bool
		lockCh          chan bool
		lock            sync.Mutex
		maxGoroutineNum int
		goroutineNum    int
		taskCh          chan *list.Element
	}

	timingEntry struct {
		baseEntry
		value  interface{}
		circle int
	}

	baseEntry struct {
		delay time.Duration
		key   interface{}
	}

	positionEntry struct {
		pos  int
		item *list.Element // item.Value ->  *timingEntry
	}

	timingTask struct {
		key   interface{}
		value interface{}
	}
)

// NewTimingWheel returns a TimingWheel.
func NewTimingWheel(interval time.Duration, numSlots int, execute Execute) (*TimingWheel, error) {
	if interval <= 0 || numSlots <= 0 || execute == nil {
		return nil, fmt.Errorf("interval: %v, slots: %d, execute: %p", interval, numSlots, execute)
	}

	return newTimingWheelWithClock(interval, numSlots, execute, timex.NewRealTicker(interval))
}

func newTimingWheelWithClock(interval time.Duration, numSlots int, execute Execute, ticker timex.Ticker) (
	*TimingWheel, error) {
	tw := &TimingWheel{
		interval:      interval,
		ticker:        ticker,
		slots:         make([]map[int]*list.List, numSlots),
		timers:        collection.NewSafeMap(),
		tickedPos:     numSlots - 1, // at previous virtual circle
		execute:       execute,
		numSlots:      numSlots,
		setChannel:    make(chan timingEntry),
		moveChannel:   make(chan baseEntry),
		removeChannel: make(chan interface{}),
		drainChannel:  make(chan func(key, value interface{})),
		stopChannel:   make(chan struct{}),

		cacheList:       list.New(),
		taskList:        list.New(),
		locked:          true,
		lockCh:          make(chan bool),
		taskCh:          make(chan *list.Element),
		maxGoroutineNum: defaultMaxGoroutineNum,
	}

	tw.initSlots()
	go tw.runTasks()
	go tw.run()

	return tw, nil
}

// Drain drains all items and executes them.
func (tw *TimingWheel) Drain(fn func(key, value interface{})) {
	tw.drainChannel <- fn
}

// MoveTimer moves the task with the given key to the given delay.
func (tw *TimingWheel) MoveTimer(key interface{}, delay time.Duration) {
	if delay <= 0 || key == nil {
		return
	}

	tw.moveChannel <- baseEntry{
		delay: delay,
		key:   key,
	}
}

// RemoveTimer removes the task with the given key.
func (tw *TimingWheel) RemoveTimer(key interface{}) {
	if key == nil {
		return
	}

	tw.removeChannel <- key
}

// SetTimer sets the task value with the given key to the delay.
func (tw *TimingWheel) SetTimer(key, value interface{}, delay time.Duration) {
	if delay <= 0 || key == nil {
		return
	}

	tw.setChannel <- timingEntry{
		baseEntry: baseEntry{
			delay: delay,
			key:   key,
		},
		value: value,
	}
}

// Stop stops tw.
func (tw *TimingWheel) Stop() {
	close(tw.stopChannel)
}

func (tw *TimingWheel) drainAll(fn func(key, value interface{})) {
	runner := threading.NewTaskRunner(drainWorkers)
	for k := range tw.slots {
		for k1 := range tw.slots[k] {
			slot := tw.slots[k][k1]
			for e := slot.Front(); e != nil; {
				task := e.Value.(*timingEntry)
				slot.Remove(e)
				e = slot.Front()
				runner.Schedule(func() {
					fn(task.key, task.value)
				})
			}
		}
		tw.slots[k] = make(map[int]*list.List)
	}
}

func (tw *TimingWheel) getPositionAndCircle(d time.Duration) (pos, circle int) {
	steps := int(d / tw.interval)
	pos = (tw.tickedPos + steps) % tw.numSlots
	circle = (steps - 1) / tw.numSlots
	return
}

func (tw *TimingWheel) initSlots() {
	for i := 0; i < tw.numSlots; i++ {
		tw.slots[i] = make(map[int]*list.List)
	}
}

func (tw *TimingWheel) moveTask(task baseEntry) {
	val, ok := tw.timers.Get(task.key)
	if !ok {
		return
	}

	timer := val.(*positionEntry)
	item := timer.item.Value.(*timingEntry)
	// immediate execution
	if task.delay < tw.interval {
		threading.GoSafe(func() {
			tw.execute(item.key, item.value)
		})
		return
	}

	pos, circle := tw.getPositionAndCircle(task.delay)
	// unchanged
	if pos == timer.pos && circle == item.circle {
		return
	}

	if _, ok := tw.slots[pos][circle]; !ok {
		tw.slots[pos][circle] = list.New()
	}
	// remove old element
	tw.slots[timer.pos][item.circle].Remove(timer.item)

	item.circle = circle
	timer.pos = pos

	// push new element
	tw.slots[pos][circle].PushBack(item)
	timer.item = tw.slots[pos][circle].Back()
}

func (tw *TimingWheel) onTick() {
	tw.tickedPos = (tw.tickedPos + 1) % tw.numSlots
	m := tw.slots[tw.tickedPos]
	if len(m) == 0 {
		return
	}
	newM := make(map[int]*list.List)
	for k1 := range m {
		if k1 > 0 {
			newM[k1-1] = m[k1]
		} else {
			for e := m[k1].Front(); e != nil; e = m[k1].Next(e) {
				timingEntry := e.Value.(*timingEntry)
				tw.timers.Del(timingEntry.key)
			}
			tw.lock.Lock()
			tw.cacheList.MergeBack(m[k1])
			if tw.locked {
				tw.locked = false
				tw.lockCh <- tw.locked
			}
			tw.lock.Unlock()
			m[k1].Init()
		}
	}
	tw.slots[tw.tickedPos] = newM
}

func (tw *TimingWheel) removeTask(key interface{}) {
	val, ok := tw.timers.Get(key)
	if !ok {
		return
	}

	timer := val.(*positionEntry)
	item := timer.item.Value.(*timingEntry)
	tw.timers.Del(key)
	tw.slots[timer.pos][item.circle].Remove(timer.item)
}

func (tw *TimingWheel) run() {
	for {
		select {
		case <-tw.ticker.Chan():
			tw.onTick()
		case task := <-tw.setChannel:
			tw.setTask(&task)
		case key := <-tw.removeChannel:
			tw.removeTask(key)
		case task := <-tw.moveChannel:
			tw.moveTask(task)
		case fn := <-tw.drainChannel:
			tw.drainAll(fn)
		case <-tw.stopChannel:
			tw.ticker.Stop()
			return
		}
	}
}

func (tw *TimingWheel) runTasks() {
	for {
		select {
		case <-tw.lockCh:
			for {
				e := tw.taskList.Back()
				if e == nil {
					tw.lock.Lock()
					tw.taskList.MergeBack(tw.cacheList)
					tw.cacheList.Init()
					e = tw.taskList.Back()
					if e == nil {
						for tw.goroutineNum > 0 {
							tw.goroutineNum--
							tw.taskCh <- nil
						}
						tw.locked = true
						tw.lock.Unlock()
						break
					}
					tw.lock.Unlock()
				}
				tw.taskList.Remove(e)
				if tw.goroutineNum < tw.maxGoroutineNum {
					tw.goroutineNum++
					go tw.worker()
				}
				tw.taskCh <- e
			}
		}
	}
}

func (tw *TimingWheel) worker() {
	for e := range tw.taskCh {
		if e == nil {
			break
		}
		tw.execute(e.Value.(*timingEntry).key, e.Value.(*timingEntry).value)
	}
}

func (tw *TimingWheel) setTask(task *timingEntry) {
	if task.delay < tw.interval {
		task.delay = tw.interval
	}

	if val, ok := tw.timers.Get(task.key); ok {
		entry := val.(*positionEntry)
		entry.item.Value.(*timingEntry).value = task.value
		tw.moveTask(task.baseEntry)
	} else {
		pos, circle := tw.getPositionAndCircle(task.delay)
		task.circle = circle
		if _, ok := tw.slots[pos][circle]; !ok {
			tw.slots[pos][circle] = list.New()
		}
		tw.slots[pos][circle].PushBack(task)
		tw.timers.Set(task.key, &positionEntry{
			pos:  pos,
			item: tw.slots[pos][circle].Back(),
		})
	}
}

func (tw *TimingWheel) GetValue(key string) (interface{}, bool) {
	val, ok := tw.timers.Get(key)
	if !ok {
		return nil, ok
	}
	return val.(*positionEntry).item.Value, ok
}

func (tw *TimingWheel) Len() int {
	return tw.timers.Size()
}
