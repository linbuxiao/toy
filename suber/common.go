package suber

import "fmt"

/**
What is mutex mean?
If there has a mutex and got locked by the Lock invoke. if you need
to invoke Lock at the same time in other goroutine, the goroutine will be blocked.
*/

type actionType int

const (
	sub actionType = iota + 1
	pub
	shutdown
)

type cmd struct {
	op     actionType
	topics []string
	msg    interface{}
	c      chan interface{}
}

type Suber struct {
	cmdChan chan cmd
}

type register struct {
	x map[string]chan interface{}
}

func (s *Suber) do(op actionType, msg interface{}, topics ...string) chan interface{} {
	// every time we sub, we will return a channel
	c := make(chan interface{}, 1)
	s.cmdChan <- cmd{
		op:     op,
		topics: topics,
		msg:    msg,
		c:      c,
	}
	return c
}

func (s *Suber) Sub(key string) chan interface{} {
	return s.do(sub, nil, key)
}

func (s *Suber) Pub(msg interface{}, topics ...string) {
	s.do(pub, msg, topics...)
}

func (s *Suber) Shutdown() {
	s.do(shutdown, nil)
}

func (s *Suber) start() {
	reg := &register{
		x: make(map[string]chan interface{}),
	}

loop:
	for {
		select {
		case c := <-s.cmdChan:
			switch c.op {
			case pub:
				for _, t := range c.topics {
					reg.x[t] <- c.msg
				}
			case sub:
				for _, t := range c.topics {
					reg.x[t] = c.c
				}
			case shutdown:
				break loop
			default:
				fmt.Printf("cannot parse op: %d", c.op)
				break loop
			}
		}
	}

	for _, c := range reg.x {
		close(c)
	}
}

func New() *Suber {
	s := &Suber{
		cmdChan: make(chan cmd),
	}
	// listen action case
	go s.start()
	return s
}
