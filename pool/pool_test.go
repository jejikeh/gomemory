package pool

import (
	"runtime"
	"testing"
)

func TestNewObjectInPool(t *testing.T) { // @Incomplete.
	t.Parallel()

	pool := NewPool[int, int](1024)
	x := pool.StoreAt(0)

	if x == nil {
		t.Errorf("failed to allocate new object in pool")
	}

	if pool.Length() != 1 {
		t.Errorf("pool length expected to be %d, but got %d", 1, pool.Length())
	}
}

// @Cleanup: Cleanup this test.
func TestGetObjectFromPool(t *testing.T) {
	t.Parallel()

	pool := NewPool[int, int](1024)
	x := pool.StoreAt(0)
	*x = 123

	xFromPool, _ := pool.LoadAt(0)
	if xFromPool == nil {
		t.Error("failed to get object from pool by index")
	} else if *xFromPool != 123 {
		t.Errorf("expected %d but got %d", 123, *xFromPool)
	}

	*x = 234

	xFromPool, _ = pool.LoadAt(0)
	if xFromPool == nil {
		t.Error("failed to get object from pool by index")
	} else if *xFromPool != 234 {
		t.Errorf("expected %d but got %d", 234, *xFromPool)
	}

	y := pool.StoreAt(1)
	*y = 102

	yFromPool, _ := pool.LoadAt(1)
	if yFromPool == nil {
		t.Error("failed to get object from pool by index")
	} else if *yFromPool != 102 {
		t.Errorf("expected %d but got %d", 102, *yFromPool)
	}

	xFromPool, _ = pool.LoadAt(0)
	if xFromPool == nil {
		t.Error("failed to get object from pool by index")
	} else if *xFromPool != 234 {
		t.Errorf("expected %d but got %d", 234, *xFromPool)
	}

	if pool.Length() != 2 {
		t.Errorf("pool length expected to be %d, but got %d", 2, pool.Length())
	}
}

func TestPool(t *testing.T) {
	t.Parallel()

	untypedPool := NewPool[int, int](16)
	x := (*int)(untypedPool.StoreAt(1))
	*x = 123

	ptr, _ := untypedPool.LoadAt(1)

	xFromPool := (*int)(ptr)
	if xFromPool == nil {
		t.Error("failed to get object from pool by index")
	} else if *xFromPool != 123 {
		t.Errorf("expected %d but got %d", 123, *xFromPool)
	}
}

type Sprite struct {
	Path string
}

type Transform struct {
	X, Y   float32
	Sprite *Sprite
}

func TestPoolGC(t *testing.T) {
	t.Parallel()

	pool := NewPool[int, Transform](2)
	allocateObject(pool)

	runtime.GC()

	testGetAllocatedObject(t, pool)
}

func allocateObject(pool *Pool[int, Transform]) {
	t1 := pool.StoreAt(1)
	t1.Sprite = &Sprite{
		Path: "sprite",
	}
}

func testGetAllocatedObject(t *testing.T, pool *Pool[int, Transform]) {
	t.Helper()

	t1, _ := pool.LoadAt(1)
	if t1.Sprite == nil {
		t.Fail()
	}

	if t1.Sprite.Path != "sprite" {
		t.Fail()
	}
}
