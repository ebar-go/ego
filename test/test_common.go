package test

// 自定义测试，Assert风格
type TestCommon interface {
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fail()
	FailNow()
	Failed() bool
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Helper()
	Log(args ...interface{})
	Logf(format string, args ...interface{})
	Name() string
	Skip(args ...interface{})
	SkipNow()
	Skipf(format string, args ...interface{})
	Skipped() bool
}

// AssertEqual 判断两个变量是否相等
func AssertEqual(t TestCommon, a, b interface{}) {
	t.Helper()
	if a != b {
		t.Errorf("Not Equal. %d %d", a, b)
	}
}

// AssertTrue 判断变量是否为true
func AssertTrue(t TestCommon, a bool) {
	t.Helper()
	if !a {
		t.Errorf("Not True %t", a)
	}
}

// AssertFalse 判断变量是否为false
func AssertFalse(t TestCommon, a bool) {
	t.Helper()
	if a {
		t.Errorf("Not True %t", a)
	}
}

// AssertNil 判断变量是否为nil
func AssertNil(t TestCommon, a interface{}) {
	t.Helper()
	if a != nil {
		t.Error("Not Nil")
	}
}

// AssertNotNil 判断变量是否不为nil
func AssertNotNil(t TestCommon, a interface{}) {
	t.Helper()
	if a == nil {
		t.Error("Is Nil")
	}
}

// AssertEmpty 判断分片(数组)是否为空
func AssertEmpty(t TestCommon, a []interface{})  {
	t.Helper()
	if len(a) == 0{
		t.Error("Not Empty")
	}
}

// AssertNotEmpty 判断分片(数组)是否不为空
func AssertNotEmpty(t TestCommon, a []interface{})  {
	t.Helper()
	if len(a) == 0{
		t.Error("Is Empty")
	}
}


