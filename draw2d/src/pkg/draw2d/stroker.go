// Copyright 2010 The draw2d Authors. All rights reserved.
// created: 13/12/2010 by Laurent Le Goff

package draw2d

type Cap int

const (
	RoundCap Cap = iota
	ButtCap
	SquareCap
)

type Join int

const (
	BevelJoin Join = iota
	RoundJoin
	MiterJoin
)


type LineStroker struct {
	Next          VertexConverter
	HalfLineWidth float
	Cap           Cap
	Join          Join
	vertices      []float
	rewind        []float
	x, y, nx, ny  float
	command       VertexCommand
}

func NewLineStroker(converter VertexConverter) *LineStroker {
	l := new(LineStroker)
	l.Next = converter
	l.HalfLineWidth = 0.5
	l.vertices = make([]float, 0)
	l.rewind = make([]float, 0)
	l.Cap = ButtCap
	l.Join = MiterJoin
	l.command = VertexNoCommand
	return l
}


func (l *LineStroker) NextCommand(command VertexCommand) {
	l.command = command
	if command == VertexStopCommand {
		l.Next.NextCommand(VertexStartCommand)
		for i, j := 0, 1; j < len(l.vertices); i, j = i+2, j+2 {
			l.Next.Vertex(l.vertices[i], l.vertices[j])
			l.Next.NextCommand(VertexNoCommand)
		}
		for i, j := len(l.rewind)-2, len(l.rewind)-1; j > 0; i, j = i-2, j-2 {
			l.Next.NextCommand(VertexNoCommand)
			l.Next.Vertex(l.rewind[i], l.rewind[j])
		}
		if len(l.vertices) > 1 {
			l.Next.NextCommand(VertexNoCommand)
			l.Next.Vertex(l.vertices[0], l.vertices[1])
		}
		l.Next.NextCommand(VertexStopCommand)
		// reinit vertices	
		l.vertices = make([]float, 0)
		l.rewind = make([]float, 0)
		l.x, l.y, l.nx, l.ny = 0, 0, 0, 0
	}
}

func (l *LineStroker) Vertex(x, y float) {
	switch l.command {
	case VertexNoCommand:
		l.line(l.x, l.y, x, y)
	case VertexStartCommand:
		l.x, l.y = x, y
	case VertexJoinCommand:
		l.joinLine(l.x, l.y, l.nx, l.ny, x, y)
	case VertexCloseCommand:
		l.line(l.x, l.y, x, y)
		l.joinLine(l.x, l.y, l.nx, l.ny, x, y)
		l.closePolygon()
	}
	l.command = VertexNoCommand
}

func (l *LineStroker) closePolygon() {
	if len(l.vertices) > 1 {
		l.vertices = append(l.vertices, l.vertices[0], l.vertices[1])
		l.rewind = append(l.rewind, l.rewind[0], l.rewind[1])
	}
}


func (l *LineStroker) line(x1, y1, x2, y2 float) {
	dx := (x2 - x1)
	dy := (y2 - y1)
	d := vectorDistance(dx, dy)
	if d != 0 {
		nx := dy * l.HalfLineWidth / d
		ny := -(dx * l.HalfLineWidth / d)
		l.vertices = append(l.vertices, x1+nx, y1+ny, x2+nx, y2+ny)
		l.rewind = append(l.rewind, x1-nx, y1-ny, x2-nx, y2-ny)
		l.x, l.y, l.nx, l.ny = x2, y2, nx, ny
	}
}

func (l *LineStroker) joinLine(x1, y1, nx1, ny1, x2, y2 float) {
	dx := (x2 - x1)
	dy := (y2 - y1)
	d := vectorDistance(dx, dy)

	if d != 0 {
		nx := dy * l.HalfLineWidth / d
		ny := -(dx * l.HalfLineWidth / d)
		/*	l.join(x1, y1, x1 + nx, y1 - ny, nx, ny, x1 + ny2, y1 + nx2, nx2, ny2)
			l.join(x1, y1, x1 - ny1, y1 - nx1, nx1, ny1, x1 - ny2, y1 - nx2, nx2, ny2)*/

		l.vertices = append(l.vertices, x1+nx, y1+ny, x2+nx, y2+ny)
		l.rewind = append(l.rewind, x1-nx, y1-ny, x2-nx, y2-ny)
		l.x, l.y, l.nx, l.ny = x2, y2, nx, ny
	}
}
