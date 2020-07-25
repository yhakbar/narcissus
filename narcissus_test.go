package narcissus

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
)

type mockedGetParameter struct {
	ssmiface.SSMAPI
	MockData map[string]string
}

func (m mockedGetParameter) GetParameter(input *ssm.GetParameterInput) (*ssm.GetParameterOutput, error) {
	mockValue := m.MockData[*input.Name]
	parameterOutput := ssm.GetParameterOutput{
		Parameter: &ssm.Parameter{
			Value: &mockValue,
		},
	}
	return &parameterOutput, nil
}

var wrapper = Wrapper{
	Client: mockedGetParameter{
		MockData: map[string]string{
			"/path/to/parameters/Name/FirstName":             "Jane",
			"/path/to/parameters/Name/LastName":              "Doe",
			"/path/to/parameters/Contact/Email":              "jane.doe@email.com",
			"/path/to/parameters/Contact/Number":             "(123)456-7890",
			"/path/to/parameters/FavoriteNumber":             "26",
			"/path/to/parameters/FavoriteInconvenientNumber": "26.24",
		},
	},
}

func TestUpdateBySSM(t *testing.T) {
	type Name struct {
		FirstName string `ssm:"Name/FirstName"`
		LastName  string `ssm:"Name/LastName"`
	}

	type Contact struct {
		Email  string `ssm:"Contact/Email"`
		Number string `ssm:"Contact/Number"`
	}

	type Person struct {
		Name                       Name
		Contact                    Contact
		FavoriteNumber             int     `ssm:"FavoriteNumber"`
		FavoriteInconvenientNumber float64 `ssm:"FavoriteInconvenientNumber"`
	}

	person := Person{}
	ssmPath := "/path/to/parameters/"
	err := wrapper.UpdateBySSM(&person, &ssmPath)
	if err != nil {
		t.Log(err)
	}
	expectedPerson := Person{
		Name: Name{
			FirstName: "Jane",
			LastName:  "Doe",
		},
		Contact: Contact{
			Email:  "jane.doe@email.com",
			Number: "(123)456-7890",
		},
		FavoriteNumber:             26,
		FavoriteInconvenientNumber: 26.24,
	}
	assert.Equal(t, person, expectedPerson, "Updated incorrectly by SSM")
}

func ExampleUpdateBySSM() {
	type Name struct {
		FirstName string `ssm:"Name/FirstName"`
		LastName  string `ssm:"Name/LastName"`
	}

	type Contact struct {
		Email  string `ssm:"Contact/Email"`
		Number string `ssm:"Contact/Number"`
	}

	type Person struct {
		Name                       Name
		Contact                    Contact
		FavoriteNumber             int     `ssm:"FavoriteNumber"`
		FavoriteInconvenientNumber float64 `ssm:"FavoriteInconvenientNumber"`
	}

	person := Person{}
	ssmPath := "/path/to/parameters/"
	// You can get this wrapper like so: wrapper := narcissus.Wrapper{Client: client}
	_ = wrapper.UpdateBySSM(&person, &ssmPath)
	// If you don't want to reuse a client, simply use narcissus.UpdateBySSM(&person, &ssmPath)
	fmt.Println(person)
	// Output: {{Jane Doe} {jane.doe@email.com (123)456-7890} 26 26.24}
}
