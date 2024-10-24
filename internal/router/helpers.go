package router

func (r *router) getExecRoute(userId int64) *Route {
	r.mu.Lock()
	defer r.mu.Unlock()

	route, ok := r.routesExec[userId]
	if !ok {
		return nil
	}

	return route
}

func (r *router) lockRoute(userId int64, route *Route) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.routesExec[userId]; !ok {
		r.routesExec[userId] = route
	}
}

func (r *router) unlockRoute(userId int64) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.routesExec, userId)
}

func (r *router) initRoutes() {
	r.ctxHandler.FillHandlers(r)
	r.gameHandler.FillHandlers(r)
}

func createButton(text string) button {
	return button{
		text: text,
	}
}
