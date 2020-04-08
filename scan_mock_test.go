package libsql

// Code generated by http://github.com/gojuno/minimock (dev). DO NOT EDIT.

import (
	"sync"
	mm_atomic "sync/atomic"
	mm_time "time"

	"github.com/gojuno/minimock/v3"
)

// ScanDoerMock implements scanDoer
type ScanDoerMock struct {
	t minimock.Tester

	funcDo          func(rowScanner RowScanner, oneRow bool, query func() (sqlRows, error)) (err error)
	inspectFuncDo   func(rowScanner RowScanner, oneRow bool, query func() (sqlRows, error))
	afterDoCounter  uint64
	beforeDoCounter uint64
	DoMock          mScanDoerMockDo
}

// NewScanDoerMock returns a mock for scanDoer
func NewScanDoerMock(t minimock.Tester) *ScanDoerMock {
	m := &ScanDoerMock{t: t}
	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.DoMock = mScanDoerMockDo{mock: m}
	m.DoMock.callArgs = []*ScanDoerMockDoParams{}

	return m
}

type mScanDoerMockDo struct {
	mock               *ScanDoerMock
	defaultExpectation *ScanDoerMockDoExpectation
	expectations       []*ScanDoerMockDoExpectation

	callArgs []*ScanDoerMockDoParams
	mutex    sync.RWMutex
}

// ScanDoerMockDoExpectation specifies expectation struct of the scanDoer.Do
type ScanDoerMockDoExpectation struct {
	mock    *ScanDoerMock
	params  *ScanDoerMockDoParams
	results *ScanDoerMockDoResults
	Counter uint64
}

// ScanDoerMockDoParams contains parameters of the scanDoer.Do
type ScanDoerMockDoParams struct {
	rowScanner RowScanner
	oneRow     bool
	query      func() (sqlRows, error)
}

// ScanDoerMockDoResults contains results of the scanDoer.Do
type ScanDoerMockDoResults struct {
	err error
}

// Expect sets up expected params for scanDoer.Do
func (mmDo *mScanDoerMockDo) Expect(rowScanner RowScanner, oneRow bool, query func() (sqlRows, error)) *mScanDoerMockDo {
	if mmDo.mock.funcDo != nil {
		mmDo.mock.t.Fatalf("ScanDoerMock.Do mock is already set by Set")
	}

	if mmDo.defaultExpectation == nil {
		mmDo.defaultExpectation = &ScanDoerMockDoExpectation{}
	}

	mmDo.defaultExpectation.params = &ScanDoerMockDoParams{rowScanner, oneRow, query}
	for _, e := range mmDo.expectations {
		if minimock.Equal(e.params, mmDo.defaultExpectation.params) {
			mmDo.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmDo.defaultExpectation.params)
		}
	}

	return mmDo
}

// Inspect accepts an inspector function that has same arguments as the scanDoer.Do
func (mmDo *mScanDoerMockDo) Inspect(f func(rowScanner RowScanner, oneRow bool, query func() (sqlRows, error))) *mScanDoerMockDo {
	if mmDo.mock.inspectFuncDo != nil {
		mmDo.mock.t.Fatalf("Inspect function is already set for ScanDoerMock.Do")
	}

	mmDo.mock.inspectFuncDo = f

	return mmDo
}

// Return sets up results that will be returned by scanDoer.Do
func (mmDo *mScanDoerMockDo) Return(err error) *ScanDoerMock {
	if mmDo.mock.funcDo != nil {
		mmDo.mock.t.Fatalf("ScanDoerMock.Do mock is already set by Set")
	}

	if mmDo.defaultExpectation == nil {
		mmDo.defaultExpectation = &ScanDoerMockDoExpectation{mock: mmDo.mock}
	}
	mmDo.defaultExpectation.results = &ScanDoerMockDoResults{err}
	return mmDo.mock
}

//Set uses given function f to mock the scanDoer.Do method
func (mmDo *mScanDoerMockDo) Set(f func(rowScanner RowScanner, oneRow bool, query func() (sqlRows, error)) (err error)) *ScanDoerMock {
	if mmDo.defaultExpectation != nil {
		mmDo.mock.t.Fatalf("Default expectation is already set for the scanDoer.Do method")
	}

	if len(mmDo.expectations) > 0 {
		mmDo.mock.t.Fatalf("Some expectations are already set for the scanDoer.Do method")
	}

	mmDo.mock.funcDo = f
	return mmDo.mock
}

// When sets expectation for the scanDoer.Do which will trigger the result defined by the following
// Then helper
func (mmDo *mScanDoerMockDo) When(rowScanner RowScanner, oneRow bool, query func() (sqlRows, error)) *ScanDoerMockDoExpectation {
	if mmDo.mock.funcDo != nil {
		mmDo.mock.t.Fatalf("ScanDoerMock.Do mock is already set by Set")
	}

	expectation := &ScanDoerMockDoExpectation{
		mock:   mmDo.mock,
		params: &ScanDoerMockDoParams{rowScanner, oneRow, query},
	}
	mmDo.expectations = append(mmDo.expectations, expectation)
	return expectation
}

// Then sets up scanDoer.Do return parameters for the expectation previously defined by the When method
func (e *ScanDoerMockDoExpectation) Then(err error) *ScanDoerMock {
	e.results = &ScanDoerMockDoResults{err}
	return e.mock
}

// Do implements scanDoer
func (mmDo *ScanDoerMock) Do(rowScanner RowScanner, oneRow bool, query func() (sqlRows, error)) (err error) {
	mm_atomic.AddUint64(&mmDo.beforeDoCounter, 1)
	defer mm_atomic.AddUint64(&mmDo.afterDoCounter, 1)

	if mmDo.inspectFuncDo != nil {
		mmDo.inspectFuncDo(rowScanner, oneRow, query)
	}

	mm_params := &ScanDoerMockDoParams{rowScanner, oneRow, query}

	// Record call args
	mmDo.DoMock.mutex.Lock()
	mmDo.DoMock.callArgs = append(mmDo.DoMock.callArgs, mm_params)
	mmDo.DoMock.mutex.Unlock()

	for _, e := range mmDo.DoMock.expectations {
		if minimock.Equal(e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.err
		}
	}

	if mmDo.DoMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmDo.DoMock.defaultExpectation.Counter, 1)
		mm_want := mmDo.DoMock.defaultExpectation.params
		mm_got := ScanDoerMockDoParams{rowScanner, oneRow, query}
		if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmDo.t.Errorf("ScanDoerMock.Do got unexpected parameters, want: %#v, got: %#v%s\n", *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmDo.DoMock.defaultExpectation.results
		if mm_results == nil {
			mmDo.t.Fatal("No results are set for the ScanDoerMock.Do")
		}
		return (*mm_results).err
	}
	if mmDo.funcDo != nil {
		return mmDo.funcDo(rowScanner, oneRow, query)
	}
	mmDo.t.Fatalf("Unexpected call to ScanDoerMock.Do. %v %v %v", rowScanner, oneRow, query)
	return
}

// DoAfterCounter returns a count of finished ScanDoerMock.Do invocations
func (mmDo *ScanDoerMock) DoAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmDo.afterDoCounter)
}

// DoBeforeCounter returns a count of ScanDoerMock.Do invocations
func (mmDo *ScanDoerMock) DoBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmDo.beforeDoCounter)
}

// Calls returns a list of arguments used in each call to ScanDoerMock.Do.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmDo *mScanDoerMockDo) Calls() []*ScanDoerMockDoParams {
	mmDo.mutex.RLock()

	argCopy := make([]*ScanDoerMockDoParams, len(mmDo.callArgs))
	copy(argCopy, mmDo.callArgs)

	mmDo.mutex.RUnlock()

	return argCopy
}

// MinimockDoDone returns true if the count of the Do invocations corresponds
// the number of defined expectations
func (m *ScanDoerMock) MinimockDoDone() bool {
	for _, e := range m.DoMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.DoMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterDoCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcDo != nil && mm_atomic.LoadUint64(&m.afterDoCounter) < 1 {
		return false
	}
	return true
}

// MinimockDoInspect logs each unmet expectation
func (m *ScanDoerMock) MinimockDoInspect() {
	for _, e := range m.DoMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to ScanDoerMock.Do with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.DoMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterDoCounter) < 1 {
		if m.DoMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to ScanDoerMock.Do")
		} else {
			m.t.Errorf("Expected call to ScanDoerMock.Do with params: %#v", *m.DoMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcDo != nil && mm_atomic.LoadUint64(&m.afterDoCounter) < 1 {
		m.t.Error("Expected call to ScanDoerMock.Do")
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *ScanDoerMock) MinimockFinish() {
	if !m.minimockDone() {
		m.MinimockDoInspect()
		m.t.FailNow()
	}
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *ScanDoerMock) MinimockWait(timeout mm_time.Duration) {
	timeoutCh := mm_time.After(timeout)
	for {
		if m.minimockDone() {
			return
		}
		select {
		case <-timeoutCh:
			m.MinimockFinish()
			return
		case <-mm_time.After(10 * mm_time.Millisecond):
		}
	}
}

func (m *ScanDoerMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockDoDone()
}