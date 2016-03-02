package gock

import (
	"github.com/nbio/st"
	"testing"
)

func TestStoreRegister(t *testing.T) {
	defer after()
	st.Expect(t, len(mocks), 0)
	mock := New("foo").Mock
	Register(mock)
	st.Expect(t, len(mocks), 1)
	st.Expect(t, mock.Request().Mock, mock)
	st.Expect(t, mock.Response().Mock, mock)
}

func TestStoreGetAll(t *testing.T) {
	defer after()
	st.Expect(t, len(mocks), 0)
	mock := New("foo").Mock
	store := GetAll()
	st.Expect(t, len(mocks), 1)
	st.Expect(t, len(store), 1)
	st.Expect(t, store[0], mock)
}

func TestStoreExists(t *testing.T) {
	defer after()
	st.Expect(t, len(mocks), 0)
	mock := New("foo").Mock
	st.Expect(t, len(mocks), 1)
	st.Expect(t, Exists(mock), true)
}

func TestStoreRemove(t *testing.T) {
	defer after()
	st.Expect(t, len(mocks), 0)
	mock := New("foo").Mock
	st.Expect(t, len(mocks), 1)
	st.Expect(t, Exists(mock), true)

	Remove(mock)
	st.Expect(t, Exists(mock), false)

	Remove(mock)
	st.Expect(t, Exists(mock), false)
}

func TestStoreFlush(t *testing.T) {
	defer after()
	st.Expect(t, len(mocks), 0)

	mock1 := New("foo").Mock
	mock2 := New("foo").Mock
	st.Expect(t, len(mocks), 2)
	st.Expect(t, Exists(mock1), true)
	st.Expect(t, Exists(mock2), true)

	Flush()
	st.Expect(t, len(mocks), 0)
	st.Expect(t, Exists(mock1), false)
	st.Expect(t, Exists(mock2), false)
}
