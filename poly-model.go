package main

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/peterhellberg/wavefront"
	"github.com/pkg/errors"
)

type face struct {
	a, b, c    vertex
	nA, nB, nC vertex
}

type polyModel struct {
	objs        map[string]*wavefront.Object
	displayList uint32
}

func loadObjFile(filename string) (polyModel, error) {
	var err error
	var pm polyModel

	pm.objs, err = wavefront.Read(filename)
	if err != nil {
		return polyModel{}, errors.Wrap(err, "loading OBJ file "+filename)
	}

	pm.displayList = objToDisplayList(pm.objs)

	return pm, nil
}

func (pm *polyModel) draw() {
	gl.PushMatrix()

	gl.CallList(pm.displayList)

	gl.PopMatrix()
}

func (pm *polyModel) destroy() {
	gl.DeleteLists(pm.displayList, 1)
}

func objToDisplayList(objs map[string]*wavefront.Object) uint32 {
	listID := gl.GenLists(1)

	gl.NewList(listID, gl.COMPILE)

	for _, obj := range objs {
		for _, group := range obj.Groups {
			gl.Materialfv(gl.FRONT, gl.AMBIENT, fPtr(group.Material.Ambient))
			gl.Materialfv(gl.FRONT, gl.DIFFUSE, fPtr(group.Material.Diffuse))
			gl.Materialfv(gl.FRONT, gl.SPECULAR, fPtr(group.Material.Specular))
			gl.Materialf(gl.FRONT, gl.SHININESS, group.Material.Shininess)

			gl.Begin(gl.TRIANGLES)

			for i := 0; i < len(group.Vertexes); i += 3 {
				gl.Vertex3f(group.Vertexes[i], group.Vertexes[i+1], group.Vertexes[i+2])
				gl.Normal3f(group.Normals[i], group.Normals[i+1], group.Normals[i+2])
			}

			gl.End()
		}
	}

	gl.EndList()

	return listID
}
