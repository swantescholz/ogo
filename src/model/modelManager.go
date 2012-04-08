package model

import (
	//"github.com/banthar/gl"
	//"fmt"
	//"math"
	//. "glmath/color"
	//"glmath/mat4"
	//. "glmath/util"
)

var GModelPath1, GModelPath2 = "res/models/", ".obj"
func SetDefaultModelPaths(start, end string) {
	GModelPath1, GModelPath2 = start, end
}
func modelStringToPath(s string) string {return GModelPath1 + s + GModelPath2}

var GModelManager = &ModelManager{make(map[string]*ModelData)}
type ModelManager struct {
	mtex map[string]*ModelData
}
func (m *ModelManager) ClearAllData() {
	for k,v := range m.mtex {
		v.Destroy()
		delete(m.mtex, k)
	}
}
func (m *ModelManager) Delete(path string) {
	path = modelStringToPath(path)
	v, ok := m.mtex[path]
	if !ok {return}
	v.Destroy() //cleaning gl-memory
	delete(m.mtex, path)
}
func (m *ModelManager) Load(path string) *Model {
	path = modelStringToPath(path)
	v,ok := m.mtex[path]
	if ok {return NewModel(v)}
	m.mtex[path] = NewModelData(path)
	return NewModel(m.mtex[path])
}
func (m *ModelManager) Get(path string) *Model {
	path = modelStringToPath(path)
	v,ok := m.mtex[path]
	if !ok {return nil}
	return NewModel(v)
}

func ClearAllData()      {GModelManager.ClearAllData()}
func Delete(path string) {GModelManager.Delete(path)}
func Load  (path string) *Model {return GModelManager.Load(path)}
func Get   (path string) *Model {return GModelManager.Get(path)}


