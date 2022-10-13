package runtime

type errCaller func() error

// Call 同步依次执行函数，只要有一个返回err,直接返回
// 通过此方法可以减少if err != nil { return }的判断
func Call(callers ...errCaller) error {
	for _, caller := range callers {
		if caller == nil {
			continue
		}
		if err := caller(); err != nil {
			return err
		}
	}

	return nil
}

// IgnoreErrorCaller return new function that will ignore error
func IgnoreErrorCaller(caller errCaller) func() {
	return func() {
		_ = caller()
	}
}
