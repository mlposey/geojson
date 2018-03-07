package main

// Point models a point in a trajectory.
type Point struct {
	// Latitude
	Lat float64 `json:"lat"`
	// Longitude
	Lng float64 `json:"lng"`
	// Altitude
	Alt float64 `json:"alt"`
}

// Trajectory models a sequence of connected points.
type Trajectory struct {
	// A unique id for the trajectory
	ID string `json:"id"`
	// The sequence of points which make the trajectory
	Path []Point `json:"path"`
}
