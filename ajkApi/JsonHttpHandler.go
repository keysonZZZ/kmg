package ajkApi

import (
	"encoding/json"
	//"errors"
	"fmt"
	//"github.com/bronze1man/kmg/kmgReflect"
	"github.com/bronze1man/kmg/kmgLog"
	"github.com/bronze1man/kmg/sessionStore"
	"net/http"
	"reflect"
	"time"
)

type JsonHttpInput struct {
	Name string
	Guid string //
	Data interface{}
}
type httpInput struct {
	Name string
	Guid string //
	Data json.RawMessage
}
type JsonHttpOutput struct {
	Err  string
	Guid string // "" as not set guid to peer
	Data interface{}
}
type JsonHttpHandler struct {
	ApiManager          ApiManagerInterface
	SessionStoreManager *sessionStore.Manager
	//	ReflectDecl         *kmgReflect.ContextDecl
}

func (handler *JsonHttpHandler) Filter(c *HttpApiContext, _ []HttpApiFilter) {
	handler.ServeHTTP(c.ResponseWriter, c.Request)
}
func (handler *JsonHttpHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	startTime := time.Now()
	var err error
	rawInput := &httpInput{}
	defer func() {
		go httpLog(httpLogRequest{
			Name:    rawInput.Name,
			Dur:     time.Since(startTime).String(),
			Err:     err,
			SessId:  rawInput.Guid,
			Ip:      req.RemoteAddr,
		})
	}()
	defer req.Body.Close()
	err = json.NewDecoder(req.Body).Decode(rawInput)
	if err != nil {
		handler.returnOutput(w, &JsonHttpOutput{Err: err.Error()})
		return
	}
	var apiOutput interface{}
	session, err := handler.SessionStoreManager.Load(rawInput.Guid)
	if err != nil {
		err = fmt.Errorf("[session.load] %s", err.Error())
		return
	}
	err = handler.ApiManager.RpcCall(session, rawInput.Name, func(meta *ApiFuncMeta) error {
		apiOutput, err = structRpcCall(meta, rawInput)
		return err
	})

	if err != nil {
		handler.returnOutput(w, &JsonHttpOutput{Err: err.Error(), Guid: session.Id})
		return
	}
	err = handler.SessionStoreManager.Save(session)
	if err != nil {
		err = fmt.Errorf("[session.Save] %s", err.Error())
		return
	}
	handler.returnOutput(w, &JsonHttpOutput{Data: apiOutput, Guid: session.Id})
}

type httpLogRequest struct {
	Name    string
	Dur     string
	Err     error
	SessId  string
	Ip      string
}

func httpLog(req httpLogRequest) {
	errStr := ""
	if req.Err != nil {
		errStr = req.Err.Error()
	}
	kmgLog.Log("apiAccess", errStr, req)
	if errStr != "" {
		kmgLog.Log("apiError", errStr, req)
	}
}

//TODO finish rpcCall by function param name
/*
func (handler *JsonHttpHandler) rpcCall(funcMeta *ApiFuncMeta, rawInput *httpInput) (interface{}, error) {
	if handler.ReflectDecl == nil {
		return structRpcCall(funcMeta, rawInput)
	}
	objectReflectType := funcMeta.AttachObject.Type()
	f, ok := handler.ReflectDecl.GetMethodDeclByReflectType(objectReflectType, funcMeta.MethodName)
	if !ok {
		return nil, fmt.Errorf("not found method in ReflectDecl %s.%s", objectReflectType.Name(), funcMeta.MethodName)
	}
	reflectFuncDecl, err := f.GetReflectFuncDecl(funcMeta.Func.Type(), funcMeta.IsMethod)
	if err != nil {
		return nil, fmt.Errorf("func %s.%s FuncDecl not match reflect err:%s", objectReflectType.Name(), funcMeta.MethodName, err.Error())
	}
	if len(reflectFuncDecl.Results) > 0 && !reflectFuncDecl.ResultHasNames {
		return nil, fmt.Errorf("func %s.%s need have result name to become a api func", objectReflectType.Name(), funcMeta.MethodName)
	}
	inValues := make([]reflect.Value, funcMeta.Func.Type().NumIn())
	if funcMeta.IsMethod {
		inValues[0] = funcMeta.AttachObject
	}
	if len(reflectFuncDecl.Params) > 0 {
		inRaw := map[string]json.RawMessage{}
		err := json.Unmarshal([]byte(rawInput.Data), inRaw)
		if err != nil {
			return nil, fmt.Errorf("api input shuold be a map :%s", err.Error())
		}
		for key, rawData := range inRaw {
			field, ok := reflectFuncDecl.ParamMap[key]
			if !ok {
				continue
			}
			thisValuePtr := reflect.New(field.Type)
			err := json.Unmarshal([]byte(rawData), thisValuePtr.Interface())
			if err != nil {
				return nil, fmt.Errorf("api input key: %s, type not match: %s", key, err.Error())
			}
			inValues[field.Index] = thisValuePtr.Elem()
		}
		//zero value input for key not in ParamMap
		for i, value := range inValues {
			if value.IsValid() {
				continue
			}
			inValues[i] = reflect.Zero(funcMeta.Func.Type().In(i))
		}
	}
	return nil, errors.New("not implement rpcCall by function param name")
}
*/
func structRpcCall(funcMeta *ApiFuncMeta, rawInput *httpInput) (interface{}, error) {
	funcType := funcMeta.Func.Type()
	var inValues []reflect.Value
	var apiOutputValue reflect.Value
	serviceValue := funcMeta.AttachObject
	switch funcType.NumIn() {
	case 1:
		inValues = []reflect.Value{serviceValue}
	case 2:
		apiInputValue, err := jsonUnmarshalFromPtrReflectType(funcType.In(1), []byte(rawInput.Data))
		if err != nil {
			return nil, err
		}
		inValues = []reflect.Value{serviceValue, apiInputValue}
	case 3:
		apiInputValue, err := jsonUnmarshalFromPtrReflectType(funcType.In(1), []byte(rawInput.Data))
		if err != nil {
			return nil, err
		}
		apiOutputValue = reflect.New(funcType.In(2).Elem())
		inValues = []reflect.Value{serviceValue, apiInputValue, apiOutputValue}
	default:
		return nil, &ApiFuncArgumentError{Reason: "only accept function input argument num 0,1,2", ApiName: rawInput.Name}
	}
	switch funcType.NumOut() {
	case 0:
	case 1:
		if funcType.Out(0).Kind() != reflect.Interface {
			return nil, &ApiFuncArgumentError{
				Reason:  "only accept function output one argument with error",
				ApiName: rawInput.Name,
			}
		}
	default:
		return nil, &ApiFuncArgumentError{Reason: "only accept function output argument num 0,1", ApiName: rawInput.Name}
	}
	outValues := funcMeta.Func.Call(inValues)
	var output interface{}
	if apiOutputValue.IsValid() && apiOutputValue.CanInterface() {
		output = apiOutputValue.Interface()
	}
	if len(outValues) == 1 {
		if outValues[0].IsNil() {
			return output, nil
		}
		err, ok := outValues[0].Interface().(error)
		if ok == false {
			return nil, &ApiFuncArgumentError{
				Reason:  "only accept function output one argument with error",
				ApiName: rawInput.Name,
			}
		}
		return nil, err
	}
	return output, nil
}

func jsonUnmarshalFromPtrReflectType(inputType reflect.Type, data []byte) (reflect.Value, error) {
	var apiInputValue = reflect.New(inputType.Elem())
	apiInput := apiInputValue.Interface()
	err := json.Unmarshal(data, apiInput)
	if err != nil {
		return reflect.Value{}, err
	}
	return apiInputValue, nil
}
func (handler *JsonHttpHandler) returnOutput(w http.ResponseWriter, output *JsonHttpOutput) {
	w.Header().Set("Content-Type", "text/json; charset=utf-8")
	err := json.NewEncoder(w).Encode(output)
	if err != nil {
		kmgLog.Log("apiError", "[JsonHttpHandler.returnOutput] json.NewEncoder(w).Encode(output)"+err.Error(), nil)
	}
}
