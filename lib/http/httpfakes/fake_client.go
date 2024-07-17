// Code generated by counterfeiter. DO NOT EDIT.
package httpfakes

import (
	"context"
	httpa "net/http"
	"sync"

	retryablehttp "github.com/hashicorp/go-retryablehttp"
	"github.com/kava-forge/eve-alts/lib/http"
)

type FakeClient struct {
	ConfigureRetriesStub        func(http.RetryOptions)
	configureRetriesMutex       sync.RWMutex
	configureRetriesArgsForCall []struct {
		arg1 http.RetryOptions
	}
	DeleteBodyStub        func(context.Context, string, ...func(req *retryablehttp.Request) error) ([]byte, *httpa.Response, error)
	deleteBodyMutex       sync.RWMutex
	deleteBodyArgsForCall []struct {
		arg1 context.Context
		arg2 string
		arg3 []func(req *retryablehttp.Request) error
	}
	deleteBodyReturns struct {
		result1 []byte
		result2 *httpa.Response
		result3 error
	}
	deleteBodyReturnsOnCall map[int]struct {
		result1 []byte
		result2 *httpa.Response
		result3 error
	}
	DeleteJSONStub        func(context.Context, interface{}, string, ...func(req *retryablehttp.Request) error) (*httpa.Response, error)
	deleteJSONMutex       sync.RWMutex
	deleteJSONArgsForCall []struct {
		arg1 context.Context
		arg2 interface{}
		arg3 string
		arg4 []func(req *retryablehttp.Request) error
	}
	deleteJSONReturns struct {
		result1 *httpa.Response
		result2 error
	}
	deleteJSONReturnsOnCall map[int]struct {
		result1 *httpa.Response
		result2 error
	}
	GetBodyStub        func(context.Context, string, ...func(req *retryablehttp.Request) error) ([]byte, *httpa.Response, error)
	getBodyMutex       sync.RWMutex
	getBodyArgsForCall []struct {
		arg1 context.Context
		arg2 string
		arg3 []func(req *retryablehttp.Request) error
	}
	getBodyReturns struct {
		result1 []byte
		result2 *httpa.Response
		result3 error
	}
	getBodyReturnsOnCall map[int]struct {
		result1 []byte
		result2 *httpa.Response
		result3 error
	}
	GetJSONStub        func(context.Context, interface{}, string, ...func(req *retryablehttp.Request) error) (*httpa.Response, error)
	getJSONMutex       sync.RWMutex
	getJSONArgsForCall []struct {
		arg1 context.Context
		arg2 interface{}
		arg3 string
		arg4 []func(req *retryablehttp.Request) error
	}
	getJSONReturns struct {
		result1 *httpa.Response
		result2 error
	}
	getJSONReturnsOnCall map[int]struct {
		result1 *httpa.Response
		result2 error
	}
	PatchBodyStub        func(context.Context, string, ...func(req *retryablehttp.Request) error) ([]byte, *httpa.Response, error)
	patchBodyMutex       sync.RWMutex
	patchBodyArgsForCall []struct {
		arg1 context.Context
		arg2 string
		arg3 []func(req *retryablehttp.Request) error
	}
	patchBodyReturns struct {
		result1 []byte
		result2 *httpa.Response
		result3 error
	}
	patchBodyReturnsOnCall map[int]struct {
		result1 []byte
		result2 *httpa.Response
		result3 error
	}
	PatchJSONStub        func(context.Context, interface{}, string, ...func(req *retryablehttp.Request) error) (*httpa.Response, error)
	patchJSONMutex       sync.RWMutex
	patchJSONArgsForCall []struct {
		arg1 context.Context
		arg2 interface{}
		arg3 string
		arg4 []func(req *retryablehttp.Request) error
	}
	patchJSONReturns struct {
		result1 *httpa.Response
		result2 error
	}
	patchJSONReturnsOnCall map[int]struct {
		result1 *httpa.Response
		result2 error
	}
	PostBodyStub        func(context.Context, string, ...func(req *retryablehttp.Request) error) ([]byte, *httpa.Response, error)
	postBodyMutex       sync.RWMutex
	postBodyArgsForCall []struct {
		arg1 context.Context
		arg2 string
		arg3 []func(req *retryablehttp.Request) error
	}
	postBodyReturns struct {
		result1 []byte
		result2 *httpa.Response
		result3 error
	}
	postBodyReturnsOnCall map[int]struct {
		result1 []byte
		result2 *httpa.Response
		result3 error
	}
	PostJSONStub        func(context.Context, interface{}, string, ...func(req *retryablehttp.Request) error) (*httpa.Response, error)
	postJSONMutex       sync.RWMutex
	postJSONArgsForCall []struct {
		arg1 context.Context
		arg2 interface{}
		arg3 string
		arg4 []func(req *retryablehttp.Request) error
	}
	postJSONReturns struct {
		result1 *httpa.Response
		result2 error
	}
	postJSONReturnsOnCall map[int]struct {
		result1 *httpa.Response
		result2 error
	}
	PutBodyStub        func(context.Context, string, ...func(req *retryablehttp.Request) error) ([]byte, *httpa.Response, error)
	putBodyMutex       sync.RWMutex
	putBodyArgsForCall []struct {
		arg1 context.Context
		arg2 string
		arg3 []func(req *retryablehttp.Request) error
	}
	putBodyReturns struct {
		result1 []byte
		result2 *httpa.Response
		result3 error
	}
	putBodyReturnsOnCall map[int]struct {
		result1 []byte
		result2 *httpa.Response
		result3 error
	}
	PutJSONStub        func(context.Context, interface{}, string, ...func(req *retryablehttp.Request) error) (*httpa.Response, error)
	putJSONMutex       sync.RWMutex
	putJSONArgsForCall []struct {
		arg1 context.Context
		arg2 interface{}
		arg3 string
		arg4 []func(req *retryablehttp.Request) error
	}
	putJSONReturns struct {
		result1 *httpa.Response
		result2 error
	}
	putJSONReturnsOnCall map[int]struct {
		result1 *httpa.Response
		result2 error
	}
	RequestBodyStub        func(context.Context, string, string, ...func(req *retryablehttp.Request) error) ([]byte, *httpa.Response, error)
	requestBodyMutex       sync.RWMutex
	requestBodyArgsForCall []struct {
		arg1 context.Context
		arg2 string
		arg3 string
		arg4 []func(req *retryablehttp.Request) error
	}
	requestBodyReturns struct {
		result1 []byte
		result2 *httpa.Response
		result3 error
	}
	requestBodyReturnsOnCall map[int]struct {
		result1 []byte
		result2 *httpa.Response
		result3 error
	}
	RequestJSONStub        func(context.Context, interface{}, string, string, ...func(req *retryablehttp.Request) error) (*httpa.Response, error)
	requestJSONMutex       sync.RWMutex
	requestJSONArgsForCall []struct {
		arg1 context.Context
		arg2 interface{}
		arg3 string
		arg4 string
		arg5 []func(req *retryablehttp.Request) error
	}
	requestJSONReturns struct {
		result1 *httpa.Response
		result2 error
	}
	requestJSONReturnsOnCall map[int]struct {
		result1 *httpa.Response
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeClient) ConfigureRetries(arg1 http.RetryOptions) {
	fake.configureRetriesMutex.Lock()
	fake.configureRetriesArgsForCall = append(fake.configureRetriesArgsForCall, struct {
		arg1 http.RetryOptions
	}{arg1})
	stub := fake.ConfigureRetriesStub
	fake.recordInvocation("ConfigureRetries", []interface{}{arg1})
	fake.configureRetriesMutex.Unlock()
	if stub != nil {
		fake.ConfigureRetriesStub(arg1)
	}
}

func (fake *FakeClient) ConfigureRetriesCallCount() int {
	fake.configureRetriesMutex.RLock()
	defer fake.configureRetriesMutex.RUnlock()
	return len(fake.configureRetriesArgsForCall)
}

func (fake *FakeClient) ConfigureRetriesCalls(stub func(http.RetryOptions)) {
	fake.configureRetriesMutex.Lock()
	defer fake.configureRetriesMutex.Unlock()
	fake.ConfigureRetriesStub = stub
}

func (fake *FakeClient) ConfigureRetriesArgsForCall(i int) http.RetryOptions {
	fake.configureRetriesMutex.RLock()
	defer fake.configureRetriesMutex.RUnlock()
	argsForCall := fake.configureRetriesArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeClient) DeleteBody(arg1 context.Context, arg2 string, arg3 ...func(req *retryablehttp.Request) error) ([]byte, *httpa.Response, error) {
	fake.deleteBodyMutex.Lock()
	ret, specificReturn := fake.deleteBodyReturnsOnCall[len(fake.deleteBodyArgsForCall)]
	fake.deleteBodyArgsForCall = append(fake.deleteBodyArgsForCall, struct {
		arg1 context.Context
		arg2 string
		arg3 []func(req *retryablehttp.Request) error
	}{arg1, arg2, arg3})
	stub := fake.DeleteBodyStub
	fakeReturns := fake.deleteBodyReturns
	fake.recordInvocation("DeleteBody", []interface{}{arg1, arg2, arg3})
	fake.deleteBodyMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3...)
	}
	if specificReturn {
		return ret.result1, ret.result2, ret.result3
	}
	return fakeReturns.result1, fakeReturns.result2, fakeReturns.result3
}

func (fake *FakeClient) DeleteBodyCallCount() int {
	fake.deleteBodyMutex.RLock()
	defer fake.deleteBodyMutex.RUnlock()
	return len(fake.deleteBodyArgsForCall)
}

func (fake *FakeClient) DeleteBodyCalls(stub func(context.Context, string, ...func(req *retryablehttp.Request) error) ([]byte, *httpa.Response, error)) {
	fake.deleteBodyMutex.Lock()
	defer fake.deleteBodyMutex.Unlock()
	fake.DeleteBodyStub = stub
}

func (fake *FakeClient) DeleteBodyArgsForCall(i int) (context.Context, string, []func(req *retryablehttp.Request) error) {
	fake.deleteBodyMutex.RLock()
	defer fake.deleteBodyMutex.RUnlock()
	argsForCall := fake.deleteBodyArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeClient) DeleteBodyReturns(result1 []byte, result2 *httpa.Response, result3 error) {
	fake.deleteBodyMutex.Lock()
	defer fake.deleteBodyMutex.Unlock()
	fake.DeleteBodyStub = nil
	fake.deleteBodyReturns = struct {
		result1 []byte
		result2 *httpa.Response
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeClient) DeleteBodyReturnsOnCall(i int, result1 []byte, result2 *httpa.Response, result3 error) {
	fake.deleteBodyMutex.Lock()
	defer fake.deleteBodyMutex.Unlock()
	fake.DeleteBodyStub = nil
	if fake.deleteBodyReturnsOnCall == nil {
		fake.deleteBodyReturnsOnCall = make(map[int]struct {
			result1 []byte
			result2 *httpa.Response
			result3 error
		})
	}
	fake.deleteBodyReturnsOnCall[i] = struct {
		result1 []byte
		result2 *httpa.Response
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeClient) DeleteJSON(arg1 context.Context, arg2 interface{}, arg3 string, arg4 ...func(req *retryablehttp.Request) error) (*httpa.Response, error) {
	fake.deleteJSONMutex.Lock()
	ret, specificReturn := fake.deleteJSONReturnsOnCall[len(fake.deleteJSONArgsForCall)]
	fake.deleteJSONArgsForCall = append(fake.deleteJSONArgsForCall, struct {
		arg1 context.Context
		arg2 interface{}
		arg3 string
		arg4 []func(req *retryablehttp.Request) error
	}{arg1, arg2, arg3, arg4})
	stub := fake.DeleteJSONStub
	fakeReturns := fake.deleteJSONReturns
	fake.recordInvocation("DeleteJSON", []interface{}{arg1, arg2, arg3, arg4})
	fake.deleteJSONMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3, arg4...)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeClient) DeleteJSONCallCount() int {
	fake.deleteJSONMutex.RLock()
	defer fake.deleteJSONMutex.RUnlock()
	return len(fake.deleteJSONArgsForCall)
}

func (fake *FakeClient) DeleteJSONCalls(stub func(context.Context, interface{}, string, ...func(req *retryablehttp.Request) error) (*httpa.Response, error)) {
	fake.deleteJSONMutex.Lock()
	defer fake.deleteJSONMutex.Unlock()
	fake.DeleteJSONStub = stub
}

func (fake *FakeClient) DeleteJSONArgsForCall(i int) (context.Context, interface{}, string, []func(req *retryablehttp.Request) error) {
	fake.deleteJSONMutex.RLock()
	defer fake.deleteJSONMutex.RUnlock()
	argsForCall := fake.deleteJSONArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4
}

func (fake *FakeClient) DeleteJSONReturns(result1 *httpa.Response, result2 error) {
	fake.deleteJSONMutex.Lock()
	defer fake.deleteJSONMutex.Unlock()
	fake.DeleteJSONStub = nil
	fake.deleteJSONReturns = struct {
		result1 *httpa.Response
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) DeleteJSONReturnsOnCall(i int, result1 *httpa.Response, result2 error) {
	fake.deleteJSONMutex.Lock()
	defer fake.deleteJSONMutex.Unlock()
	fake.DeleteJSONStub = nil
	if fake.deleteJSONReturnsOnCall == nil {
		fake.deleteJSONReturnsOnCall = make(map[int]struct {
			result1 *httpa.Response
			result2 error
		})
	}
	fake.deleteJSONReturnsOnCall[i] = struct {
		result1 *httpa.Response
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) GetBody(arg1 context.Context, arg2 string, arg3 ...func(req *retryablehttp.Request) error) ([]byte, *httpa.Response, error) {
	fake.getBodyMutex.Lock()
	ret, specificReturn := fake.getBodyReturnsOnCall[len(fake.getBodyArgsForCall)]
	fake.getBodyArgsForCall = append(fake.getBodyArgsForCall, struct {
		arg1 context.Context
		arg2 string
		arg3 []func(req *retryablehttp.Request) error
	}{arg1, arg2, arg3})
	stub := fake.GetBodyStub
	fakeReturns := fake.getBodyReturns
	fake.recordInvocation("GetBody", []interface{}{arg1, arg2, arg3})
	fake.getBodyMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3...)
	}
	if specificReturn {
		return ret.result1, ret.result2, ret.result3
	}
	return fakeReturns.result1, fakeReturns.result2, fakeReturns.result3
}

func (fake *FakeClient) GetBodyCallCount() int {
	fake.getBodyMutex.RLock()
	defer fake.getBodyMutex.RUnlock()
	return len(fake.getBodyArgsForCall)
}

func (fake *FakeClient) GetBodyCalls(stub func(context.Context, string, ...func(req *retryablehttp.Request) error) ([]byte, *httpa.Response, error)) {
	fake.getBodyMutex.Lock()
	defer fake.getBodyMutex.Unlock()
	fake.GetBodyStub = stub
}

func (fake *FakeClient) GetBodyArgsForCall(i int) (context.Context, string, []func(req *retryablehttp.Request) error) {
	fake.getBodyMutex.RLock()
	defer fake.getBodyMutex.RUnlock()
	argsForCall := fake.getBodyArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeClient) GetBodyReturns(result1 []byte, result2 *httpa.Response, result3 error) {
	fake.getBodyMutex.Lock()
	defer fake.getBodyMutex.Unlock()
	fake.GetBodyStub = nil
	fake.getBodyReturns = struct {
		result1 []byte
		result2 *httpa.Response
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeClient) GetBodyReturnsOnCall(i int, result1 []byte, result2 *httpa.Response, result3 error) {
	fake.getBodyMutex.Lock()
	defer fake.getBodyMutex.Unlock()
	fake.GetBodyStub = nil
	if fake.getBodyReturnsOnCall == nil {
		fake.getBodyReturnsOnCall = make(map[int]struct {
			result1 []byte
			result2 *httpa.Response
			result3 error
		})
	}
	fake.getBodyReturnsOnCall[i] = struct {
		result1 []byte
		result2 *httpa.Response
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeClient) GetJSON(arg1 context.Context, arg2 interface{}, arg3 string, arg4 ...func(req *retryablehttp.Request) error) (*httpa.Response, error) {
	fake.getJSONMutex.Lock()
	ret, specificReturn := fake.getJSONReturnsOnCall[len(fake.getJSONArgsForCall)]
	fake.getJSONArgsForCall = append(fake.getJSONArgsForCall, struct {
		arg1 context.Context
		arg2 interface{}
		arg3 string
		arg4 []func(req *retryablehttp.Request) error
	}{arg1, arg2, arg3, arg4})
	stub := fake.GetJSONStub
	fakeReturns := fake.getJSONReturns
	fake.recordInvocation("GetJSON", []interface{}{arg1, arg2, arg3, arg4})
	fake.getJSONMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3, arg4...)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeClient) GetJSONCallCount() int {
	fake.getJSONMutex.RLock()
	defer fake.getJSONMutex.RUnlock()
	return len(fake.getJSONArgsForCall)
}

func (fake *FakeClient) GetJSONCalls(stub func(context.Context, interface{}, string, ...func(req *retryablehttp.Request) error) (*httpa.Response, error)) {
	fake.getJSONMutex.Lock()
	defer fake.getJSONMutex.Unlock()
	fake.GetJSONStub = stub
}

func (fake *FakeClient) GetJSONArgsForCall(i int) (context.Context, interface{}, string, []func(req *retryablehttp.Request) error) {
	fake.getJSONMutex.RLock()
	defer fake.getJSONMutex.RUnlock()
	argsForCall := fake.getJSONArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4
}

func (fake *FakeClient) GetJSONReturns(result1 *httpa.Response, result2 error) {
	fake.getJSONMutex.Lock()
	defer fake.getJSONMutex.Unlock()
	fake.GetJSONStub = nil
	fake.getJSONReturns = struct {
		result1 *httpa.Response
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) GetJSONReturnsOnCall(i int, result1 *httpa.Response, result2 error) {
	fake.getJSONMutex.Lock()
	defer fake.getJSONMutex.Unlock()
	fake.GetJSONStub = nil
	if fake.getJSONReturnsOnCall == nil {
		fake.getJSONReturnsOnCall = make(map[int]struct {
			result1 *httpa.Response
			result2 error
		})
	}
	fake.getJSONReturnsOnCall[i] = struct {
		result1 *httpa.Response
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) PatchBody(arg1 context.Context, arg2 string, arg3 ...func(req *retryablehttp.Request) error) ([]byte, *httpa.Response, error) {
	fake.patchBodyMutex.Lock()
	ret, specificReturn := fake.patchBodyReturnsOnCall[len(fake.patchBodyArgsForCall)]
	fake.patchBodyArgsForCall = append(fake.patchBodyArgsForCall, struct {
		arg1 context.Context
		arg2 string
		arg3 []func(req *retryablehttp.Request) error
	}{arg1, arg2, arg3})
	stub := fake.PatchBodyStub
	fakeReturns := fake.patchBodyReturns
	fake.recordInvocation("PatchBody", []interface{}{arg1, arg2, arg3})
	fake.patchBodyMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3...)
	}
	if specificReturn {
		return ret.result1, ret.result2, ret.result3
	}
	return fakeReturns.result1, fakeReturns.result2, fakeReturns.result3
}

func (fake *FakeClient) PatchBodyCallCount() int {
	fake.patchBodyMutex.RLock()
	defer fake.patchBodyMutex.RUnlock()
	return len(fake.patchBodyArgsForCall)
}

func (fake *FakeClient) PatchBodyCalls(stub func(context.Context, string, ...func(req *retryablehttp.Request) error) ([]byte, *httpa.Response, error)) {
	fake.patchBodyMutex.Lock()
	defer fake.patchBodyMutex.Unlock()
	fake.PatchBodyStub = stub
}

func (fake *FakeClient) PatchBodyArgsForCall(i int) (context.Context, string, []func(req *retryablehttp.Request) error) {
	fake.patchBodyMutex.RLock()
	defer fake.patchBodyMutex.RUnlock()
	argsForCall := fake.patchBodyArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeClient) PatchBodyReturns(result1 []byte, result2 *httpa.Response, result3 error) {
	fake.patchBodyMutex.Lock()
	defer fake.patchBodyMutex.Unlock()
	fake.PatchBodyStub = nil
	fake.patchBodyReturns = struct {
		result1 []byte
		result2 *httpa.Response
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeClient) PatchBodyReturnsOnCall(i int, result1 []byte, result2 *httpa.Response, result3 error) {
	fake.patchBodyMutex.Lock()
	defer fake.patchBodyMutex.Unlock()
	fake.PatchBodyStub = nil
	if fake.patchBodyReturnsOnCall == nil {
		fake.patchBodyReturnsOnCall = make(map[int]struct {
			result1 []byte
			result2 *httpa.Response
			result3 error
		})
	}
	fake.patchBodyReturnsOnCall[i] = struct {
		result1 []byte
		result2 *httpa.Response
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeClient) PatchJSON(arg1 context.Context, arg2 interface{}, arg3 string, arg4 ...func(req *retryablehttp.Request) error) (*httpa.Response, error) {
	fake.patchJSONMutex.Lock()
	ret, specificReturn := fake.patchJSONReturnsOnCall[len(fake.patchJSONArgsForCall)]
	fake.patchJSONArgsForCall = append(fake.patchJSONArgsForCall, struct {
		arg1 context.Context
		arg2 interface{}
		arg3 string
		arg4 []func(req *retryablehttp.Request) error
	}{arg1, arg2, arg3, arg4})
	stub := fake.PatchJSONStub
	fakeReturns := fake.patchJSONReturns
	fake.recordInvocation("PatchJSON", []interface{}{arg1, arg2, arg3, arg4})
	fake.patchJSONMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3, arg4...)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeClient) PatchJSONCallCount() int {
	fake.patchJSONMutex.RLock()
	defer fake.patchJSONMutex.RUnlock()
	return len(fake.patchJSONArgsForCall)
}

func (fake *FakeClient) PatchJSONCalls(stub func(context.Context, interface{}, string, ...func(req *retryablehttp.Request) error) (*httpa.Response, error)) {
	fake.patchJSONMutex.Lock()
	defer fake.patchJSONMutex.Unlock()
	fake.PatchJSONStub = stub
}

func (fake *FakeClient) PatchJSONArgsForCall(i int) (context.Context, interface{}, string, []func(req *retryablehttp.Request) error) {
	fake.patchJSONMutex.RLock()
	defer fake.patchJSONMutex.RUnlock()
	argsForCall := fake.patchJSONArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4
}

func (fake *FakeClient) PatchJSONReturns(result1 *httpa.Response, result2 error) {
	fake.patchJSONMutex.Lock()
	defer fake.patchJSONMutex.Unlock()
	fake.PatchJSONStub = nil
	fake.patchJSONReturns = struct {
		result1 *httpa.Response
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) PatchJSONReturnsOnCall(i int, result1 *httpa.Response, result2 error) {
	fake.patchJSONMutex.Lock()
	defer fake.patchJSONMutex.Unlock()
	fake.PatchJSONStub = nil
	if fake.patchJSONReturnsOnCall == nil {
		fake.patchJSONReturnsOnCall = make(map[int]struct {
			result1 *httpa.Response
			result2 error
		})
	}
	fake.patchJSONReturnsOnCall[i] = struct {
		result1 *httpa.Response
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) PostBody(arg1 context.Context, arg2 string, arg3 ...func(req *retryablehttp.Request) error) ([]byte, *httpa.Response, error) {
	fake.postBodyMutex.Lock()
	ret, specificReturn := fake.postBodyReturnsOnCall[len(fake.postBodyArgsForCall)]
	fake.postBodyArgsForCall = append(fake.postBodyArgsForCall, struct {
		arg1 context.Context
		arg2 string
		arg3 []func(req *retryablehttp.Request) error
	}{arg1, arg2, arg3})
	stub := fake.PostBodyStub
	fakeReturns := fake.postBodyReturns
	fake.recordInvocation("PostBody", []interface{}{arg1, arg2, arg3})
	fake.postBodyMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3...)
	}
	if specificReturn {
		return ret.result1, ret.result2, ret.result3
	}
	return fakeReturns.result1, fakeReturns.result2, fakeReturns.result3
}

func (fake *FakeClient) PostBodyCallCount() int {
	fake.postBodyMutex.RLock()
	defer fake.postBodyMutex.RUnlock()
	return len(fake.postBodyArgsForCall)
}

func (fake *FakeClient) PostBodyCalls(stub func(context.Context, string, ...func(req *retryablehttp.Request) error) ([]byte, *httpa.Response, error)) {
	fake.postBodyMutex.Lock()
	defer fake.postBodyMutex.Unlock()
	fake.PostBodyStub = stub
}

func (fake *FakeClient) PostBodyArgsForCall(i int) (context.Context, string, []func(req *retryablehttp.Request) error) {
	fake.postBodyMutex.RLock()
	defer fake.postBodyMutex.RUnlock()
	argsForCall := fake.postBodyArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeClient) PostBodyReturns(result1 []byte, result2 *httpa.Response, result3 error) {
	fake.postBodyMutex.Lock()
	defer fake.postBodyMutex.Unlock()
	fake.PostBodyStub = nil
	fake.postBodyReturns = struct {
		result1 []byte
		result2 *httpa.Response
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeClient) PostBodyReturnsOnCall(i int, result1 []byte, result2 *httpa.Response, result3 error) {
	fake.postBodyMutex.Lock()
	defer fake.postBodyMutex.Unlock()
	fake.PostBodyStub = nil
	if fake.postBodyReturnsOnCall == nil {
		fake.postBodyReturnsOnCall = make(map[int]struct {
			result1 []byte
			result2 *httpa.Response
			result3 error
		})
	}
	fake.postBodyReturnsOnCall[i] = struct {
		result1 []byte
		result2 *httpa.Response
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeClient) PostJSON(arg1 context.Context, arg2 interface{}, arg3 string, arg4 ...func(req *retryablehttp.Request) error) (*httpa.Response, error) {
	fake.postJSONMutex.Lock()
	ret, specificReturn := fake.postJSONReturnsOnCall[len(fake.postJSONArgsForCall)]
	fake.postJSONArgsForCall = append(fake.postJSONArgsForCall, struct {
		arg1 context.Context
		arg2 interface{}
		arg3 string
		arg4 []func(req *retryablehttp.Request) error
	}{arg1, arg2, arg3, arg4})
	stub := fake.PostJSONStub
	fakeReturns := fake.postJSONReturns
	fake.recordInvocation("PostJSON", []interface{}{arg1, arg2, arg3, arg4})
	fake.postJSONMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3, arg4...)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeClient) PostJSONCallCount() int {
	fake.postJSONMutex.RLock()
	defer fake.postJSONMutex.RUnlock()
	return len(fake.postJSONArgsForCall)
}

func (fake *FakeClient) PostJSONCalls(stub func(context.Context, interface{}, string, ...func(req *retryablehttp.Request) error) (*httpa.Response, error)) {
	fake.postJSONMutex.Lock()
	defer fake.postJSONMutex.Unlock()
	fake.PostJSONStub = stub
}

func (fake *FakeClient) PostJSONArgsForCall(i int) (context.Context, interface{}, string, []func(req *retryablehttp.Request) error) {
	fake.postJSONMutex.RLock()
	defer fake.postJSONMutex.RUnlock()
	argsForCall := fake.postJSONArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4
}

func (fake *FakeClient) PostJSONReturns(result1 *httpa.Response, result2 error) {
	fake.postJSONMutex.Lock()
	defer fake.postJSONMutex.Unlock()
	fake.PostJSONStub = nil
	fake.postJSONReturns = struct {
		result1 *httpa.Response
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) PostJSONReturnsOnCall(i int, result1 *httpa.Response, result2 error) {
	fake.postJSONMutex.Lock()
	defer fake.postJSONMutex.Unlock()
	fake.PostJSONStub = nil
	if fake.postJSONReturnsOnCall == nil {
		fake.postJSONReturnsOnCall = make(map[int]struct {
			result1 *httpa.Response
			result2 error
		})
	}
	fake.postJSONReturnsOnCall[i] = struct {
		result1 *httpa.Response
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) PutBody(arg1 context.Context, arg2 string, arg3 ...func(req *retryablehttp.Request) error) ([]byte, *httpa.Response, error) {
	fake.putBodyMutex.Lock()
	ret, specificReturn := fake.putBodyReturnsOnCall[len(fake.putBodyArgsForCall)]
	fake.putBodyArgsForCall = append(fake.putBodyArgsForCall, struct {
		arg1 context.Context
		arg2 string
		arg3 []func(req *retryablehttp.Request) error
	}{arg1, arg2, arg3})
	stub := fake.PutBodyStub
	fakeReturns := fake.putBodyReturns
	fake.recordInvocation("PutBody", []interface{}{arg1, arg2, arg3})
	fake.putBodyMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3...)
	}
	if specificReturn {
		return ret.result1, ret.result2, ret.result3
	}
	return fakeReturns.result1, fakeReturns.result2, fakeReturns.result3
}

func (fake *FakeClient) PutBodyCallCount() int {
	fake.putBodyMutex.RLock()
	defer fake.putBodyMutex.RUnlock()
	return len(fake.putBodyArgsForCall)
}

func (fake *FakeClient) PutBodyCalls(stub func(context.Context, string, ...func(req *retryablehttp.Request) error) ([]byte, *httpa.Response, error)) {
	fake.putBodyMutex.Lock()
	defer fake.putBodyMutex.Unlock()
	fake.PutBodyStub = stub
}

func (fake *FakeClient) PutBodyArgsForCall(i int) (context.Context, string, []func(req *retryablehttp.Request) error) {
	fake.putBodyMutex.RLock()
	defer fake.putBodyMutex.RUnlock()
	argsForCall := fake.putBodyArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeClient) PutBodyReturns(result1 []byte, result2 *httpa.Response, result3 error) {
	fake.putBodyMutex.Lock()
	defer fake.putBodyMutex.Unlock()
	fake.PutBodyStub = nil
	fake.putBodyReturns = struct {
		result1 []byte
		result2 *httpa.Response
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeClient) PutBodyReturnsOnCall(i int, result1 []byte, result2 *httpa.Response, result3 error) {
	fake.putBodyMutex.Lock()
	defer fake.putBodyMutex.Unlock()
	fake.PutBodyStub = nil
	if fake.putBodyReturnsOnCall == nil {
		fake.putBodyReturnsOnCall = make(map[int]struct {
			result1 []byte
			result2 *httpa.Response
			result3 error
		})
	}
	fake.putBodyReturnsOnCall[i] = struct {
		result1 []byte
		result2 *httpa.Response
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeClient) PutJSON(arg1 context.Context, arg2 interface{}, arg3 string, arg4 ...func(req *retryablehttp.Request) error) (*httpa.Response, error) {
	fake.putJSONMutex.Lock()
	ret, specificReturn := fake.putJSONReturnsOnCall[len(fake.putJSONArgsForCall)]
	fake.putJSONArgsForCall = append(fake.putJSONArgsForCall, struct {
		arg1 context.Context
		arg2 interface{}
		arg3 string
		arg4 []func(req *retryablehttp.Request) error
	}{arg1, arg2, arg3, arg4})
	stub := fake.PutJSONStub
	fakeReturns := fake.putJSONReturns
	fake.recordInvocation("PutJSON", []interface{}{arg1, arg2, arg3, arg4})
	fake.putJSONMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3, arg4...)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeClient) PutJSONCallCount() int {
	fake.putJSONMutex.RLock()
	defer fake.putJSONMutex.RUnlock()
	return len(fake.putJSONArgsForCall)
}

func (fake *FakeClient) PutJSONCalls(stub func(context.Context, interface{}, string, ...func(req *retryablehttp.Request) error) (*httpa.Response, error)) {
	fake.putJSONMutex.Lock()
	defer fake.putJSONMutex.Unlock()
	fake.PutJSONStub = stub
}

func (fake *FakeClient) PutJSONArgsForCall(i int) (context.Context, interface{}, string, []func(req *retryablehttp.Request) error) {
	fake.putJSONMutex.RLock()
	defer fake.putJSONMutex.RUnlock()
	argsForCall := fake.putJSONArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4
}

func (fake *FakeClient) PutJSONReturns(result1 *httpa.Response, result2 error) {
	fake.putJSONMutex.Lock()
	defer fake.putJSONMutex.Unlock()
	fake.PutJSONStub = nil
	fake.putJSONReturns = struct {
		result1 *httpa.Response
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) PutJSONReturnsOnCall(i int, result1 *httpa.Response, result2 error) {
	fake.putJSONMutex.Lock()
	defer fake.putJSONMutex.Unlock()
	fake.PutJSONStub = nil
	if fake.putJSONReturnsOnCall == nil {
		fake.putJSONReturnsOnCall = make(map[int]struct {
			result1 *httpa.Response
			result2 error
		})
	}
	fake.putJSONReturnsOnCall[i] = struct {
		result1 *httpa.Response
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) RequestBody(arg1 context.Context, arg2 string, arg3 string, arg4 ...func(req *retryablehttp.Request) error) ([]byte, *httpa.Response, error) {
	fake.requestBodyMutex.Lock()
	ret, specificReturn := fake.requestBodyReturnsOnCall[len(fake.requestBodyArgsForCall)]
	fake.requestBodyArgsForCall = append(fake.requestBodyArgsForCall, struct {
		arg1 context.Context
		arg2 string
		arg3 string
		arg4 []func(req *retryablehttp.Request) error
	}{arg1, arg2, arg3, arg4})
	stub := fake.RequestBodyStub
	fakeReturns := fake.requestBodyReturns
	fake.recordInvocation("RequestBody", []interface{}{arg1, arg2, arg3, arg4})
	fake.requestBodyMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3, arg4...)
	}
	if specificReturn {
		return ret.result1, ret.result2, ret.result3
	}
	return fakeReturns.result1, fakeReturns.result2, fakeReturns.result3
}

func (fake *FakeClient) RequestBodyCallCount() int {
	fake.requestBodyMutex.RLock()
	defer fake.requestBodyMutex.RUnlock()
	return len(fake.requestBodyArgsForCall)
}

func (fake *FakeClient) RequestBodyCalls(stub func(context.Context, string, string, ...func(req *retryablehttp.Request) error) ([]byte, *httpa.Response, error)) {
	fake.requestBodyMutex.Lock()
	defer fake.requestBodyMutex.Unlock()
	fake.RequestBodyStub = stub
}

func (fake *FakeClient) RequestBodyArgsForCall(i int) (context.Context, string, string, []func(req *retryablehttp.Request) error) {
	fake.requestBodyMutex.RLock()
	defer fake.requestBodyMutex.RUnlock()
	argsForCall := fake.requestBodyArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4
}

func (fake *FakeClient) RequestBodyReturns(result1 []byte, result2 *httpa.Response, result3 error) {
	fake.requestBodyMutex.Lock()
	defer fake.requestBodyMutex.Unlock()
	fake.RequestBodyStub = nil
	fake.requestBodyReturns = struct {
		result1 []byte
		result2 *httpa.Response
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeClient) RequestBodyReturnsOnCall(i int, result1 []byte, result2 *httpa.Response, result3 error) {
	fake.requestBodyMutex.Lock()
	defer fake.requestBodyMutex.Unlock()
	fake.RequestBodyStub = nil
	if fake.requestBodyReturnsOnCall == nil {
		fake.requestBodyReturnsOnCall = make(map[int]struct {
			result1 []byte
			result2 *httpa.Response
			result3 error
		})
	}
	fake.requestBodyReturnsOnCall[i] = struct {
		result1 []byte
		result2 *httpa.Response
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeClient) RequestJSON(arg1 context.Context, arg2 interface{}, arg3 string, arg4 string, arg5 ...func(req *retryablehttp.Request) error) (*httpa.Response, error) {
	fake.requestJSONMutex.Lock()
	ret, specificReturn := fake.requestJSONReturnsOnCall[len(fake.requestJSONArgsForCall)]
	fake.requestJSONArgsForCall = append(fake.requestJSONArgsForCall, struct {
		arg1 context.Context
		arg2 interface{}
		arg3 string
		arg4 string
		arg5 []func(req *retryablehttp.Request) error
	}{arg1, arg2, arg3, arg4, arg5})
	stub := fake.RequestJSONStub
	fakeReturns := fake.requestJSONReturns
	fake.recordInvocation("RequestJSON", []interface{}{arg1, arg2, arg3, arg4, arg5})
	fake.requestJSONMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3, arg4, arg5...)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeClient) RequestJSONCallCount() int {
	fake.requestJSONMutex.RLock()
	defer fake.requestJSONMutex.RUnlock()
	return len(fake.requestJSONArgsForCall)
}

func (fake *FakeClient) RequestJSONCalls(stub func(context.Context, interface{}, string, string, ...func(req *retryablehttp.Request) error) (*httpa.Response, error)) {
	fake.requestJSONMutex.Lock()
	defer fake.requestJSONMutex.Unlock()
	fake.RequestJSONStub = stub
}

func (fake *FakeClient) RequestJSONArgsForCall(i int) (context.Context, interface{}, string, string, []func(req *retryablehttp.Request) error) {
	fake.requestJSONMutex.RLock()
	defer fake.requestJSONMutex.RUnlock()
	argsForCall := fake.requestJSONArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4, argsForCall.arg5
}

func (fake *FakeClient) RequestJSONReturns(result1 *httpa.Response, result2 error) {
	fake.requestJSONMutex.Lock()
	defer fake.requestJSONMutex.Unlock()
	fake.RequestJSONStub = nil
	fake.requestJSONReturns = struct {
		result1 *httpa.Response
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) RequestJSONReturnsOnCall(i int, result1 *httpa.Response, result2 error) {
	fake.requestJSONMutex.Lock()
	defer fake.requestJSONMutex.Unlock()
	fake.RequestJSONStub = nil
	if fake.requestJSONReturnsOnCall == nil {
		fake.requestJSONReturnsOnCall = make(map[int]struct {
			result1 *httpa.Response
			result2 error
		})
	}
	fake.requestJSONReturnsOnCall[i] = struct {
		result1 *httpa.Response
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.configureRetriesMutex.RLock()
	defer fake.configureRetriesMutex.RUnlock()
	fake.deleteBodyMutex.RLock()
	defer fake.deleteBodyMutex.RUnlock()
	fake.deleteJSONMutex.RLock()
	defer fake.deleteJSONMutex.RUnlock()
	fake.getBodyMutex.RLock()
	defer fake.getBodyMutex.RUnlock()
	fake.getJSONMutex.RLock()
	defer fake.getJSONMutex.RUnlock()
	fake.patchBodyMutex.RLock()
	defer fake.patchBodyMutex.RUnlock()
	fake.patchJSONMutex.RLock()
	defer fake.patchJSONMutex.RUnlock()
	fake.postBodyMutex.RLock()
	defer fake.postBodyMutex.RUnlock()
	fake.postJSONMutex.RLock()
	defer fake.postJSONMutex.RUnlock()
	fake.putBodyMutex.RLock()
	defer fake.putBodyMutex.RUnlock()
	fake.putJSONMutex.RLock()
	defer fake.putJSONMutex.RUnlock()
	fake.requestBodyMutex.RLock()
	defer fake.requestBodyMutex.RUnlock()
	fake.requestJSONMutex.RLock()
	defer fake.requestJSONMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeClient) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ http.Client = new(FakeClient)
