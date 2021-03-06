// Code generated by MockGen. DO NOT EDIT.
// Source: ./db/trie/trie.go

// Package mock_trie is a generated GoMock package.
package mock_trie

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	trie "github.com/iotexproject/iotex-core/db/trie"
	reflect "reflect"
)

// MockTrie is a mock of Trie interface
type MockTrie struct {
	ctrl     *gomock.Controller
	recorder *MockTrieMockRecorder
}

// MockTrieMockRecorder is the mock recorder for MockTrie
type MockTrieMockRecorder struct {
	mock *MockTrie
}

// NewMockTrie creates a new mock instance
func NewMockTrie(ctrl *gomock.Controller) *MockTrie {
	mock := &MockTrie{ctrl: ctrl}
	mock.recorder = &MockTrieMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockTrie) EXPECT() *MockTrieMockRecorder {
	return m.recorder
}

// Start mocks base method
func (m *MockTrie) Start(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Start", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Start indicates an expected call of Start
func (mr *MockTrieMockRecorder) Start(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockTrie)(nil).Start), arg0)
}

// Stop mocks base method
func (m *MockTrie) Stop(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Stop", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Stop indicates an expected call of Stop
func (mr *MockTrieMockRecorder) Stop(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockTrie)(nil).Stop), arg0)
}

// Upsert mocks base method
func (m *MockTrie) Upsert(arg0, arg1 []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Upsert", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Upsert indicates an expected call of Upsert
func (mr *MockTrieMockRecorder) Upsert(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Upsert", reflect.TypeOf((*MockTrie)(nil).Upsert), arg0, arg1)
}

// Get mocks base method
func (m *MockTrie) Get(arg0 []byte) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get
func (mr *MockTrieMockRecorder) Get(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockTrie)(nil).Get), arg0)
}

// Delete mocks base method
func (m *MockTrie) Delete(arg0 []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockTrieMockRecorder) Delete(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockTrie)(nil).Delete), arg0)
}

// RootHash mocks base method
func (m *MockTrie) RootHash() []byte {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RootHash")
	ret0, _ := ret[0].([]byte)
	return ret0
}

// RootHash indicates an expected call of RootHash
func (mr *MockTrieMockRecorder) RootHash() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RootHash", reflect.TypeOf((*MockTrie)(nil).RootHash))
}

// SetRootHash mocks base method
func (m *MockTrie) SetRootHash(arg0 []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetRootHash", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetRootHash indicates an expected call of SetRootHash
func (mr *MockTrieMockRecorder) SetRootHash(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetRootHash", reflect.TypeOf((*MockTrie)(nil).SetRootHash), arg0)
}

// IsEmpty mocks base method
func (m *MockTrie) IsEmpty() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsEmpty")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsEmpty indicates an expected call of IsEmpty
func (mr *MockTrieMockRecorder) IsEmpty() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsEmpty", reflect.TypeOf((*MockTrie)(nil).IsEmpty))
}

// DB mocks base method
func (m *MockTrie) DB() trie.KVStore {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DB")
	ret0, _ := ret[0].(trie.KVStore)
	return ret0
}

// DB indicates an expected call of DB
func (mr *MockTrieMockRecorder) DB() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DB", reflect.TypeOf((*MockTrie)(nil).DB))
}

// deleteNodeFromDB mocks base method
func (m *MockTrie) deleteNodeFromDB(tn trie.Node) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "deleteNodeFromDB", tn)
	ret0, _ := ret[0].(error)
	return ret0
}

// deleteNodeFromDB indicates an expected call of deleteNodeFromDB
func (mr *MockTrieMockRecorder) deleteNodeFromDB(tn interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "deleteNodeFromDB", reflect.TypeOf((*MockTrie)(nil).deleteNodeFromDB), tn)
}

// putNodeIntoDB mocks base method
func (m *MockTrie) putNodeIntoDB(tn trie.Node) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "putNodeIntoDB", tn)
	ret0, _ := ret[0].(error)
	return ret0
}

// putNodeIntoDB indicates an expected call of putNodeIntoDB
func (mr *MockTrieMockRecorder) putNodeIntoDB(tn interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "putNodeIntoDB", reflect.TypeOf((*MockTrie)(nil).putNodeIntoDB), tn)
}

// loadNodeFromDB mocks base method
func (m *MockTrie) loadNodeFromDB(arg0 []byte) (trie.Node, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "loadNodeFromDB", arg0)
	ret0, _ := ret[0].(trie.Node)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// loadNodeFromDB indicates an expected call of loadNodeFromDB
func (mr *MockTrieMockRecorder) loadNodeFromDB(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "loadNodeFromDB", reflect.TypeOf((*MockTrie)(nil).loadNodeFromDB), arg0)
}

// isEmptyRootHash mocks base method
func (m *MockTrie) isEmptyRootHash(arg0 []byte) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "isEmptyRootHash", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// isEmptyRootHash indicates an expected call of isEmptyRootHash
func (mr *MockTrieMockRecorder) isEmptyRootHash(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "isEmptyRootHash", reflect.TypeOf((*MockTrie)(nil).isEmptyRootHash), arg0)
}

// emptyRootHash mocks base method
func (m *MockTrie) emptyRootHash() []byte {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "emptyRootHash")
	ret0, _ := ret[0].([]byte)
	return ret0
}

// emptyRootHash indicates an expected call of emptyRootHash
func (mr *MockTrieMockRecorder) emptyRootHash() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "emptyRootHash", reflect.TypeOf((*MockTrie)(nil).emptyRootHash))
}

// nodeHash mocks base method
func (m *MockTrie) nodeHash(tn trie.Node) []byte {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "nodeHash", tn)
	ret0, _ := ret[0].([]byte)
	return ret0
}

// nodeHash indicates an expected call of nodeHash
func (mr *MockTrieMockRecorder) nodeHash(tn interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "nodeHash", reflect.TypeOf((*MockTrie)(nil).nodeHash), tn)
}
