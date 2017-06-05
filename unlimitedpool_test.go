package gopool

import "testing"

func TestBasicUnlimitedPool_NoQueue(t *testing.T) {
	pool := &unlimitedPool{}

	pool.new = func() interface{} {
		return 1
	}

	v := pool.Get()
	if v != nil {
		t.Error("Expected nil")
		t.FailNow()
	}

	pool.Put(1)
	v = pool.Get()
	if v != nil {
		t.Error("Expected nil")
		t.FailNow()
	}
}

func TestBasicUnlimitedPool_OneItemQueue(t *testing.T) {
	pool := &unlimitedPool{}

	pool.new = func() interface{} {
		return 1
	}

	pool.queue = make(chan interface{}, 1)

	if v := pool.Get(); v != 1 {
		t.Error("Expected 1")
		t.FailNow()
	}

	pool.new = func() interface{} {
		return 2
	}
	pool.Put(1)

	if v := pool.Get(); v != 1 {
		t.Error("Expected 1")
		t.FailNow()
	}
}

func TestBasicUnlimitedPool_Calbacks(t *testing.T) {
	pool := &unlimitedPool{}

	var (
		flagnewcalled     bool
		flagreleasecalled bool
		flagcheckcalled   bool
	)

	pool.new = func() interface{} {
		flagnewcalled = true
		return 1
	}

	pool.check = func(v interface{}) bool {
		flagcheckcalled = true
		if i, ok := v.(int); !ok || i != 1 {
			t.Error("Expected 1")
			t.FailNow()
		}
		return true
	}

	pool.release = func(v interface{}) {
		flagreleasecalled = true
		if i, ok := v.(int); !ok || i != 1 {
			t.Error("Expected 1")
			t.FailNow()
		}
	}

	pool.queue = make(chan interface{}, 1)

	if v := pool.Get(); v != 1 {
		t.Error("Expected 1")
		t.FailNow()
	}

	if !flagnewcalled {
		t.Error("New callback was not called as expected")
		t.FailNow()
	}

	if flagcheckcalled {
		t.Error("Check callback was called as NOT expected")
		t.FailNow()
	}

	if flagreleasecalled {
		t.Error("Release callback was called as NOT expected")
		t.FailNow()
	}

	flagnewcalled = false

	pool.new = func() interface{} {
		flagnewcalled = true
		return 2
	}
	pool.Put(1)

	if flagnewcalled {
		t.Error("New callback was called as NOT expected")
		t.FailNow()
	}

	if flagcheckcalled {
		t.Error("Check callback was called as NOT expected")
		t.FailNow()
	}

	if flagreleasecalled {
		t.Error("Release callback was called as NOT expected")
		t.FailNow()
	}

	if v := pool.Get(); v != 1 {
		t.Error("Expected 1")
		t.FailNow()
	}

	if flagnewcalled {
		t.Error("New callback was called as NOT expected")
		t.FailNow()
	}

	if !flagcheckcalled {
		t.Error("Check callback was NOT called as expected")
		t.FailNow()
	}

	if flagreleasecalled {
		t.Error("Release callback was called as NOT expected")
		t.FailNow()
	}

	flagcheckcalled = false
	pool.Put(2) // Should be kept
	pool.Put(1) // Should be released

	if flagnewcalled {
		t.Error("New callback was called as NOT expected")
		t.FailNow()
	}

	if flagcheckcalled {
		t.Error("Check callback was called as NOT expected")
		t.FailNow()
	}

	if !flagreleasecalled {
		t.Error("Release callback was NOT called as expected")
		t.FailNow()
	}
}

func TestBasicUnlimitedPool_API(t *testing.T) {
	var (
		flagnewcalled     bool
		flagreleasecalled bool
		flagcheckcalled   bool
	)

	fnnew := func() interface{} {
		flagnewcalled = true
		return 1
	}

	fncheck := func(v interface{}) bool {
		flagcheckcalled = true
		if i, ok := v.(int); !ok || i != 1 {
			t.Error("Expected 1")
			t.FailNow()
		}
		return true
	}

	fnrelease := func(v interface{}) {
		flagreleasecalled = true
		if i, ok := v.(int); !ok || i != 1 {
			t.Error("Expected 1")
			t.FailNow()
		}
	}

	pool, err := NewPool(1, 1, fnnew, fnrelease, fncheck)

	if err != nil {
		t.Error("Errror is not expected")
		t.FailNow()
	}

	if _, ok := pool.(*unlimitedPool); !ok {
		t.Error("Expected pointer to unlimitedPool structure")
		t.FailNow()
	}

	if !flagnewcalled {
		t.Error("New callback was not called as expected")
		t.FailNow()
	}

	flagnewcalled = false
	if v := pool.Get(); v != 1 {
		t.Error("Expected 1")
		t.FailNow()
	}

	if flagnewcalled {
		t.Error("New callback was called as NOT expected")
		t.FailNow()
	}

	if !flagcheckcalled {
		t.Error("Check callback was NOT called as expected")
		t.FailNow()
	}

	if flagreleasecalled {
		t.Error("Release callback was called as NOT expected")
		t.FailNow()
	}

	flagcheckcalled = false
	pool.Put(2) // Should be kept
	pool.Put(1) // Should be released

	if flagnewcalled {
		t.Error("New callback was called as NOT expected")
		t.FailNow()
	}

	if flagcheckcalled {
		t.Error("Check callback was called as NOT expected")
		t.FailNow()
	}

	if !flagreleasecalled {
		t.Error("Release callback was NOT called as expected")
		t.FailNow()
	}
}

func TestBasicUnlimitedPool_APIError(t *testing.T) {
	var (
		flagnewcalled     bool
		flagreleasecalled bool
		flagcheckcalled   bool
	)

	fnnew := func() interface{} {
		flagnewcalled = true
		return 1
	}

	fncheck := func(v interface{}) bool {
		flagcheckcalled = true
		if i, ok := v.(int); !ok || i != 1 {
			t.Error("Expected 1")
			t.FailNow()
		}
		return true
	}

	fnrelease := func(v interface{}) {
		flagreleasecalled = true
		if i, ok := v.(int); !ok || i != 1 {
			t.Error("Expected 1")
			t.FailNow()
		}
	}

	_, err := NewPool(2, 1, fnnew, fnrelease, fncheck)

	if err == nil {
		t.Error("Errror is expected")
		t.FailNow()
	}

	_, err = NewPool(0, 1, nil, fnrelease, fncheck)

	if err == nil {
		t.Error("Errror is expected")
		t.FailNow()
	}

	if flagnewcalled {
		t.Error("New callback was called as NOT expected")
		t.FailNow()
	}

	if flagcheckcalled {
		t.Error("Check callback was called as NOT expected")
		t.FailNow()
	}

	if flagreleasecalled {
		t.Error("Release callback was called as NOT expected")
		t.FailNow()
	}
}
