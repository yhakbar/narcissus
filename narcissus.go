// Package narcissus updates a struct with fields that have been tagged with `ssm:"Parameter"` according to the corresponding value in SSM Parameter Store using reflection.
package narcissus

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
)

// Wrapper wraps the SSM Client so that it can be mocked. Also allows for client reuse.
type Wrapper struct {
	Client ssmiface.SSMAPI
}

func (wrapper *Wrapper) getSSMParameter(name *string) (*string, error) {
	withDecryption := true
	input := ssm.GetParameterInput{
		Name:           name,
		WithDecryption: &withDecryption,
	}
	value, err := wrapper.Client.GetParameter(&input)
	if err != nil {
		return nil, err
	}
	return value.Parameter.Value, nil
}

func (wrapper *Wrapper) handleSSMUpdate(field *reflect.Value, fieldType *reflect.StructField, ssmPath *string, ssm *string) error {
	parameterName := fmt.Sprintf("%s%s", *ssmPath, *ssm)
	updatedValue, err := wrapper.getSSMParameter(&parameterName)
	if err != nil {
		return err
	}
	switch field.Kind() {
	case reflect.String:
		field.SetString(*updatedValue)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		ssmInt, err := strconv.ParseInt(*updatedValue, 0, 64)
		if err != nil {
			return err
		}
		field.SetInt(ssmInt)
	case reflect.Float32, reflect.Float64:
		ssmFloat, err := strconv.ParseFloat(*updatedValue, 64)
		if err != nil {
			return err
		}
		field.SetFloat(ssmFloat)
	case reflect.Bool:
		ssmBool, err := strconv.ParseBool(*updatedValue)
		if err != nil {
			return err
		}
		field.SetBool(ssmBool)
	default:
		return fmt.Errorf("Field %s is of a type that cannot be set", fieldType.Name)
	}
	return nil
}

func getSSMClient() *ssm.SSM {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	client := ssm.New(sess)
	return client
}

// UpdateBySSM updates a struct by fields tagged `ssm` using a given wrapped SSM client.
func (wrapper *Wrapper) UpdateBySSM(generic interface{}, ssmPath *string) error {
	valueOfGeneric := reflect.ValueOf(generic).Elem()
	typeOfGeneric := valueOfGeneric.Type()

	for i := 0; i < valueOfGeneric.NumField(); i++ {
		field := valueOfGeneric.Field(i)
		fieldType := typeOfGeneric.Field(i)
		if field.Kind() == reflect.Struct {
			err := wrapper.UpdateBySSM(field.Addr().Interface(), ssmPath)
			if err != nil {
				return err
			}
		} else if ssm, ok := fieldType.Tag.Lookup("ssm"); ok {
			err := wrapper.handleSSMUpdate(&field, &fieldType, ssmPath, &ssm)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// UpdateBySSM updates a struct by fields tagged `ssm`.
func UpdateBySSM(generic interface{}, ssmPath *string) error {
	wrapper := &Wrapper{
		Client: getSSMClient(),
	}
	return wrapper.UpdateBySSM(generic, ssmPath)
}
