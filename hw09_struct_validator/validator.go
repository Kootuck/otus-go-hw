package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type (
	ValidationRuleStr string
	ValidationRuleInt string
)

const (
	StringStrictLength  ValidationRuleStr = "len"
	StringRegExp        ValidationRuleStr = "regexp"
	StringAllowedValues ValidationRuleStr = "in"
)

const (
	IntMinBoundary   ValidationRuleInt = "min"
	IntMaxBoundary   ValidationRuleInt = "max"
	IntAllowedValues ValidationRuleInt = "in"
)

type Validatable interface {
	Validate() ValidationError
}
type ValidationControl struct {
	FieldName string
	Result    bool
}
type StringValidationRule struct {
	Rule      ValidationRuleStr
	Condition interface{}
}
type IntValidationRule struct {
	Rule      ValidationRuleInt
	Condition interface{}
}
type StringValidation struct {
	StringValidationRule
	Value string
	ValidationControl
}
type IntValidation struct {
	IntValidationRule
	Value int
	ValidationControl
}

// Validation is (1) Validation Rule (how to validate) + (2) ValidatableValue - (what to validate).
func (v StringValidation) Validate() ValidationError {
	condition := v.Condition.(string)

	ret := ValidationError{
		Field: v.FieldName,
	}

	switch v.Rule {
	case StringStrictLength:
		condValue, _ := strconv.Atoi(condition)
		if len(v.Value) != condValue {
			ret.Error = NewStrictStringLengthError(len(v.Value), condValue)
		}

	case StringRegExp:
		r, _ := regexp.Compile(condition)
		ok := r.MatchString(v.Value)
		if !ok {
			ret.Error = NewStringRegExpError(condition)
		}

	case StringAllowedValues:
		splitted := strings.Split(condition, ",")
		found := false
		for _, allowed := range splitted {
			if strings.TrimSpace(allowed) == v.Value {
				found = true
			}
		}
		if !found {
			ret.Error = NewStringNotAllowedError(v.Value, condition)
		}
	}
	return ret
}

func (v IntValidation) Validate() ValidationError {
	ret := ValidationError{
		Field: v.FieldName,
	}

	switch v.Rule {
	case IntMinBoundary:
		c, _ := strconv.Atoi(v.Condition.(string))
		if v.Value < c {
			ret.Error = NewIntMustBeLargerThanError(c, v.Value)
		}

	case IntMaxBoundary:
		c, _ := strconv.Atoi(v.Condition.(string))
		if v.Value > c {
			ret.Error = NewIntMustBeLowerThanError(c, v.Value)
		}

	case IntAllowedValues:
		splitted := strings.Split(v.Condition.(string), ",")
		found := false
		for _, allowed := range splitted {
			if i, _ := strconv.Atoi(strings.TrimSpace(allowed)); i == v.Value {
				found = true
			}
		}
		if !found {
			ret.Error = NewIntNotAllowedError(v.Value, splitted)
		}
	}

	return ret
}

func Validate(v interface{}) error {
	// 0. Pre-execution check -> input argument has to be a struct
	if !isStruct(v) {
		return errors.New("only struct type is allowed as argument")
	}
	// 1. Parse validation tags into a slice of Validations to be executed
	validations, err := createValidations(reflect.TypeOf(v), reflect.ValueOf(v))
	if err != nil {
		return ProgramError{Msg: err.Error()}
	}
	// 2. Run validations, accumulate errors
	vErrors := ExecuteValidations(validations)
	if len(vErrors) == 0 {
		return nil
	}
	return vErrors
}

func createValidations(typeOf reflect.Type, valueOf reflect.Value) (validations []Validatable, err error) {
	// loop through struct's fields
	for i := 0; i < typeOf.NumField(); i++ {
		f := typeOf.Field(i)
		vTag := f.Tag.Get("validate")
		if vTag == "" {
			continue
		}

		kind := f.Type.Kind()
		value := valueOf.Field(i)

		var rulesStr []StringValidationRule
		var rulesInt []IntValidationRule
		var newValidations []Validatable
		values := make([]interface{}, 0)

		if kind == reflect.Slice {
			for j := 0; j < value.Len(); j++ {
				values = append(values, value.Index(j).Interface())
				kind = f.Type.Elem().Kind() // derive the kind from slice's element
			}
		} else {
			values = append(values, value.Interface())
		}

		if kind == reflect.String {
			rulesStr, err = parseTagIntoRulesStr(vTag)
		}
		if kind == reflect.Int {
			rulesInt, err = parseTagIntoRulesInt(vTag)
		}
		if err != nil {
			return nil, fmt.Errorf("field \"%v\": %w", f.Name, err)
		}

		for _, val := range values {
			if kind == reflect.String {
				newValidations, err = createValidationsStr(rulesStr, val, f.Name)
			}
			if kind == reflect.Int {
				newValidations, err = createValidationsInt(rulesInt, val, f.Name)
			}
			if err != nil {
				return nil, fmt.Errorf("field \"%v\": %w", f.Name, err)
			}
			validations = append(validations, newValidations...)
		}
	}
	return validations, nil
}

func ExecuteValidations(validations []Validatable) (vErrs ValidationErrors) {
	for _, v := range validations {
		err := v.Validate()
		if err.IsError() {
			vErrs = append(vErrs, err)
		}
	}
	return vErrs
}

func createValidationsInt(fieldRules []IntValidationRule,
	fieldValue interface{}, fieldName string,
) (validations []Validatable, err error) {
	valInt, ok := fieldValue.(int)
	if !ok {
		return nil, errors.New("value must be integer")
	}
	for _, fr := range fieldRules {
		validations = append(validations, IntValidation{
			IntValidationRule: IntValidationRule{
				Rule:      fr.Rule,
				Condition: fr.Condition,
			},
			Value: valInt,
			ValidationControl: ValidationControl{
				FieldName: fieldName,
				Result:    false,
			},
		})
	}
	return validations, nil
}

// validation = rule(s) + value + fieldName.
func createValidationsStr(fieldRules []StringValidationRule,
	fieldValue interface{}, fieldName string,
) (validations []Validatable, err error) {
	value, err := getFieldValueStr(fieldValue)
	if err != nil {
		return nil, err
	}

	for _, fr := range fieldRules {
		validations = append(validations, StringValidation{
			StringValidationRule: StringValidationRule{
				Rule:      fr.Rule,
				Condition: fr.Condition,
			},
			ValidationControl: ValidationControl{
				FieldName: fieldName,
				Result:    false,
			},
			Value: value,
		})
	}

	return validations, nil
}

// cast to string, custom types have to explicitly implement String() method or panic will happen.
func getFieldValueStr(val interface{}) (string, error) {
	switch v := val.(type) {
	case fmt.Stringer: // any type that has a String() string method
		return v.String(), nil // use the String() method to get a string
	default:
		if str, ok := v.(string); ok {
			return str, nil
		}
		return "", errors.New("not able to interpret as a string")
	}
}

func isStruct(i interface{}) bool {
	return reflect.TypeOf(i).Kind() == reflect.Struct
}

func checkRuleSupportedString(rule ValidationRuleStr) bool {
	switch rule {
	case StringStrictLength, StringRegExp, StringAllowedValues:
		return true
	default:
		return false
	}
}

func checkRuleSupportedInt(rule ValidationRuleInt) bool {
	switch rule {
	case IntMinBoundary, IntMaxBoundary, IntAllowedValues:
		return true
	default:
		return false
	}
}

func checkRuleConditionString(rule ValidationRuleStr, cond string) bool {
	switch rule {
	case StringStrictLength: // has to be int
		_, err := strconv.Atoi(cond)
		if err != nil {
			return false
		}
	case StringRegExp:
		_, err := regexp.Compile(cond)
		if err != nil {
			return false
		}
	case StringAllowedValues:
		return true
	default:
		return false
	}
	return true
}

func checkRuleConditionInt(rule ValidationRuleInt, cond string) bool {
	switch rule {
	case IntMinBoundary, IntMaxBoundary: // has to be int
		_, err := strconv.Atoi(cond)
		if err != nil {
			return false
		}
	case IntAllowedValues:
		allowedValues := strings.Split(cond, ",")
		for _, v := range allowedValues {
			_, err := strconv.Atoi(v)
			if err != nil {
				return false
			}
		}
	default:
		return false
	}
	return true
}

func parseTagIntoRulesInt(vTag string) (ret []IntValidationRule, err error) {
	// several validations combined possible for a single field
	validations := strings.Split(vTag, "|")
	// each boundary in combined validation is treated like a separate validation
	for _, v := range validations {
		vv := strings.Split(v, ":")
		rule := ValidationRuleInt(vv[0])
		condition := vv[1]
		if !checkRuleSupportedInt(rule) {
			return nil, fmt.Errorf("validation rule %v is not supported", rule)
		}
		if !checkRuleConditionInt(rule, condition) {
			return nil, fmt.Errorf("validation condition %v is incorrect", condition)
		}
		ret = append(ret, IntValidationRule{Rule: rule, Condition: condition})
	}
	return ret, nil
}

func parseTagIntoRulesStr(vTag string) (ret []StringValidationRule, err error) {
	// several validations combined possible for a single field
	validations := strings.Split(vTag, "|")
	// each boundary in combined validation is treated like a separate validation
	for _, v := range validations {
		vv := strings.Split(v, ":")
		rule := ValidationRuleStr(vv[0])
		condition := vv[1]
		if !checkRuleSupportedString(rule) {
			return nil, fmt.Errorf("validation rule %v is not supported", rule)
		}
		if !checkRuleConditionString(rule, condition) {
			return nil, fmt.Errorf("validation condition %v is incorrect", condition)
		}
		ret = append(ret, StringValidationRule{Rule: rule, Condition: condition})
	}
	return ret, nil
}
