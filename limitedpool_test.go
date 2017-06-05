package gopool

import "testing"

func TestBasicLimitedPool_NoQueue(t *testing.T) {
	pool := &limitedPool{}

	pool.new = func() interface{} {
		return 1
	}

	if v := pool.Get(); v != nil {
		t.Error("Expected nil")
		t.FailNow()
	}

	pool.Put(2)
	if v := pool.Get(); v != nil {
		t.Error("Expected nil")
		t.FailNow()
	}
}

func TestBasicLimitedPool_OneItemQueue(t *testing.T) {
	pool := &limitedPool{}

	pool.new = func() interface{} {
		return 1
	}

	pool.queue = make(chan interface{}, 1)
	pool.max = 1

	if v := pool.Get(); v != 1 {
		t.Error("Expected nil")
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

func TestBasicLimitedPool_Calbacks(t *testing.T) {
	pool := &limitedPool{}

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
	pool.max = 1

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

func TestBasicLimitedPool_API(t *testing.T) {
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

	pool, err := NewLPool(1, 1, fnnew, fnrelease, fncheck)

	if err != nil {
		t.Error("Errror is not expected")
		t.FailNow()
	}

	if _, ok := pool.(*limitedPool); !ok {
		t.Error("Expected pointer to limitedPool structure")
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

func TestBasicLimitedPool_APIError(t *testing.T) {
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

	_, err := NewLPool(2, 1, fnnew, fnrelease, fncheck)

	if err == nil {
		t.Error("Errror is expected")
		t.FailNow()
	}

	_, err = NewLPool(0, 1, nil, fnrelease, fncheck)

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
