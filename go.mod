module github.com/downflux/pathfinding-demo

go 1.19

require (
	github.com/downflux/go-collider v0.1.9-0.20221229204309-87ed56cfa83e
	github.com/downflux/go-geometry v0.15.4
	golang.org/x/image v0.2.0
)

require (
	github.com/downflux/go-bvh v1.0.0 // indirect
	github.com/downflux/go-pq v0.3.0 // indirect
)

replace github.com/downflux/go-collider => ../go-collider
