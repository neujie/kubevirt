package rest

import (
	"encoding/json"
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/onsi/gomega/types"
	"k8s.io/client-go/pkg/api/errors"
	"k8s.io/client-go/rest"
	rest2 "kubevirt.io/kubevirt/pkg/rest"
	"reflect"
)

func RepresentMimeType(expected interface{}) types.GomegaMatcher {
	return &representMimeTypeMatcher{
		expected: expected,
	}
}

type representMimeTypeMatcher struct {
	expected interface{}
	body     []byte
}

func (matcher *representMimeTypeMatcher) Match(actual interface{}) (success bool, err error) {
	result, ok := actual.(rest.Result)
	if !ok {
		return false, fmt.Errorf("RepresentMimeType matcher expects a kubernetes rest client Result")
	}

	//Ignore the error here, since receiving data when the return code is not 200 is valid
	matcher.body, _ = result.Raw()

	mimeType, ok := matcher.expected.(string)
	if !ok {
		return false, fmt.Errorf("Expected mime type needs to be a string")
	}

	switch mimeType {
	case rest2.MIME_JSON:
		var obj interface{}
		if err := json.Unmarshal(matcher.body, &obj); err != nil {
			return false, nil
		}
	case rest2.MIME_YAML:
		var obj interface{}
		// yaml.Unmarshal also accepts JSON, so let's check if it is JSON first
		if err := json.Unmarshal(matcher.body, &obj); err == nil {
			return false, nil
		}
		if err := yaml.Unmarshal(matcher.body, &obj); err != nil {
			return false, nil
		}
	default:
		return false, fmt.Errorf("Provided MIME-Type is not supported")
	}

	return true, nil
}

func (matcher *representMimeTypeMatcher) FailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected\n\t%#v\nto to be of type\n\t%#v", string(matcher.body), matcher.expected)
}

func (matcher *representMimeTypeMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected\n\t%#v\nnot to be of type\n\t%#v", string(matcher.body), matcher.expected)
}

func HaveBodyEqualTo(expected interface{}) types.GomegaMatcher {
	return &haveBodyEqualToMatcher{
		expected: expected,
	}
}

type haveBodyEqualToMatcher struct {
	expected interface{}
	obj      interface{}
}

func (matcher *haveBodyEqualToMatcher) Match(actual interface{}) (success bool, err error) {
	result, ok := actual.(rest.Result)
	if !ok {
		return false, fmt.Errorf("RepresentMimeType matcher expects a kubernetes rest client Result")
	}

	// Ignoring error here since failed requests can still contain a body
	matcher.obj, _ = result.Get()

	if reflect.TypeOf(matcher.expected).Kind() != reflect.Ptr {
		matcher.expected = &matcher.expected
	}
	if reflect.TypeOf(matcher.obj).Kind() != reflect.Ptr {
		matcher.obj = &matcher.obj
	}

	return reflect.DeepEqual(matcher.expected, matcher.obj), nil
}

func (matcher *haveBodyEqualToMatcher) FailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected\n\t%#v\nto to be equal to\n\t%#v", matcher.obj, matcher.expected)
}

func (matcher *haveBodyEqualToMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected\n\t%#v\nnot to be equal to\n\t%#v", matcher.obj, matcher.expected)
}

func HaveStatusCode(expected interface{}) types.GomegaMatcher {
	return &haveStatusCodeMatcher{
		expected: expected,
	}
}

type haveStatusCodeMatcher struct {
	expected   interface{}
	statusCode int
}

func (matcher *haveStatusCodeMatcher) Match(actual interface{}) (success bool, err error) {
	result, ok := actual.(rest.Result)
	if !ok {
		return false, fmt.Errorf("HaveStatusCode matcher expects a kubernetes rest client Result")
	}

	expectedStatusCode, ok := matcher.expected.(int)
	if !ok {
		return false, fmt.Errorf("Expected status code to be of type int")
	}

	result.StatusCode(&matcher.statusCode)

	if result.Error() != nil {
		matcher.statusCode = int(result.Error().(*errors.StatusError).Status().Code)
	}

	return matcher.statusCode == expectedStatusCode, nil
}

func (matcher *haveStatusCodeMatcher) FailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected status code \n\t%#v\nto to be\n\t%#v", matcher.statusCode, matcher.expected)
}

func (matcher *haveStatusCodeMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected status code \n\t%#v\nnot to be\n\t%#v", matcher.statusCode, matcher.expected)
}
