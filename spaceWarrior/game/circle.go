package game

type Circle struct {
	X      float64
	Y      float64
	radius float64
}

func NewCircle(x, y, r float64) Circle {
	return Circle{
		X:      x,
		Y:      y,
		radius: r,
	}
}

func (c Circle) MaxX() float64 {
	return c.X + c.radius
}

func (c Circle) MaxY() float64 {
	return c.Y + c.radius
}

func (c Circle) IntersectsCircle(other Circle) bool {
	return c.X <= other.MaxX() &&
		other.X <= c.MaxX() &&
		c.Y <= other.MaxY() &&
		other.Y <= c.MaxY()
}
