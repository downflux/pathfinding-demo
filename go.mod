module github.com/downflux/pathfinding-demo

go 1.19

require (
	github.com/downflux/go-boids/x v0.0.0-20230107091503-f8991804459f
	github.com/downflux/go-collider v0.2.11
	github.com/downflux/go-database v0.3.6
	github.com/downflux/go-geometry v0.15.4
	golang.org/x/image v0.3.0
)

require (
	github.com/downflux/go-bvh v1.0.0 // indirect
	github.com/downflux/go-pq v0.3.0 // indirect
)

replace github.com/downflux/go-boids/x => ../go-boids/x
