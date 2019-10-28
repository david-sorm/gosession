package gosession

import (
	"testing"
)

type mockGdb struct{}

func (mockGdb) Close() {
	return
}

func (mockGdb) CreateSession(SessionID string) {
	return
}

func (mockGdb) DestroySession(SessionID string) {
	return
}

func (mockGdb) SessionExists(SessionID string) bool {
	return false
}

func (mockGdb) ReadKey(SessionID string, Key string) interface{} {
	return ""
}

func (mockGdb) WriteKey(SessionID string, Key string, Value interface{}) {
	return
}

func (mockGdb) DeleteKey(SessionID string, Key string) {
	return
}

func Bench(b *testing.B) {
	for i := 0; i < b.N; i++{
		generateHashForSession(mockGdb{})
	}
}

func TestGenerateHashForSession(t *testing.T) {

}
