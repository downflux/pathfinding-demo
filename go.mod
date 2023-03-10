module github.com/downflux/pathfinding-demo

go 1.19

require (
	github.com/downflux/go-boids v0.3.5
	github.com/downflux/go-collider v0.2.16
	github.com/downflux/go-database v0.4.1
	github.com/downflux/go-geometry v0.16.0
	golang.org/x/image v0.3.0
)

require (
	github.com/downflux/go-bvh v1.0.0 // indirect
	github.com/downflux/go-pq v0.3.0 // indirect
)

replace github.com/downflux/go-boids => ../go-boids
