package ws

import (
	"github.com/hertz-contrib/websocket"
	"log"
	"sync"
)

type PoolManager struct {
	m map[uint64][]*websocket.Conn
	sync.Mutex
}

var Pool *PoolManager

func Init() {
	Pool = &PoolManager{
		m: make(map[uint64][]*websocket.Conn),
	}
}

func (p *PoolManager) AddConn(csId uint64, conn *websocket.Conn) {
	p.Lock()
	defer p.Unlock()
	log.Printf("客服 %d 已连接 %p", csId, conn)
	p.m[csId] = append(p.m[csId], conn)
}

func (p *PoolManager) RemoveConn(csId uint64, conn *websocket.Conn) {
	p.Lock()
	defer p.Unlock()
	log.Printf("客服 %d 断开连接 %p", csId, conn)
	for i, c := range p.m[csId] {
		if c == conn {
			p.m[csId] = append(p.m[csId][:i], p.m[csId][i+1:]...)
			return
		}
	}
}

type Response struct {
	Type int         `json:"type"`
	Data interface{} `json:"data"`
}

func (p *PoolManager) Send(csId uint64, v interface{}) {
	p.Lock()
	defer p.Unlock()

	log.Println("conn count: ", len(p.m[csId]))

	for _, conn := range p.m[csId] {
		log.Printf("send ws message %d : %+v", csId, v)

		err := conn.WriteJSON(v)
		if err != nil {
			log.Printf("err: %+v", err)
			continue
		}
	}
}
