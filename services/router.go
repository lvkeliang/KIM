package services

import (
	"KIM/protocol"
	"KIM/protocol/protoImpl"
	"errors"
	"sync"
)

type Router struct {
	middlewares []HandlerFunc // 中间件
	handlers    *FuncTree     // 注册的监听器列表
	pool        sync.Pool     // RouterContext的对象池
}

func NewRouter() *Router {
	r := &Router{
		handlers:    NewTree(),
		middlewares: make([]HandlerFunc, 0),
	}
	r.pool.New = func() interface{} {
		return NewRouterContext()
	}
	return r
}

func (r *Router) Serve(packet *protocol.LogicPkt, dispatcher Dispatcher, cache SessionStorage, session Session) error {
	if dispatcher == nil {
		return errors.New("dispatcher is nil")
	}
	if cache == nil {
		return errors.New("cache is nil")
	}

	// 从pool中获取一个RouterContextImpl
	ctx := r.pool.Get().(*RouterContextImpl)
	ctx.reset()
	// 注入
	ctx.packet = packet
	ctx.Dispatcher = dispatcher
	ctx.SessionStorage = cache
	ctx.session = session

	// 处理该RouterContextImpl
	r.serveContent(ctx)

	// 使用完后放回池中
	r.pool.Put(ctx)
	return nil
}

func (r *Router) serveContent(ctx *RouterContextImpl) {
	chain, ok := r.handlers.Get(ctx.Header().Command)
	if !ok {
		ctx.handlers = []HandlerFunc{handleNoFound}
		ctx.Next()
		return
	}
	ctx.handlers = chain
	ctx.Next()
}

// FuncTree is a tree structure
type FuncTree struct {
	nodes map[string]HandlersChain
}

// NewTree NewTree
func NewTree() *FuncTree {
	return &FuncTree{nodes: make(map[string]HandlersChain, 10)}
}

func (t *FuncTree) Add(path string, handlers ...HandlerFunc) {
	if t.nodes[path] == nil {
		t.nodes[path] = HandlersChain{}
	}

	t.nodes[path] = append(t.nodes[path], handlers...)
}

// Get a handler from tree
func (t *FuncTree) Get(path string) (HandlersChain, bool) {
	f, ok := t.nodes[path]
	return f, ok
}

// Handle register a command handler
func (r *Router) Handle(command string, handlers ...HandlerFunc) {
	r.handlers.Add(command, r.middlewares...)
	r.handlers.Add(command, handlers...)
}

func handleNoFound(ctx RouterContext) {
	_ = ctx.Resp(protoImpl.Status_NotImplemented, &protoImpl.ErrorResp{Message: "NotImplemented"})
}
