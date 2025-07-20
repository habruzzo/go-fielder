package fielder

import (
	"errors"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"reflect"
)

func (s *FieldWDefaultImpl) MarshalDynamoDBAttributeValue() (types.AttributeValue, error) {
	if s.IsDefault() {
		return nil, nil
	}

}

func (s *FieldWDefaultImpl) UnmarshalDynamoDBAttributeValue(av types.AttributeValue) error {
	convertedValS, ok := av.(*types.AttributeValueMemberS) // all types should be s
	if !ok {
		return &attributevalue.UnmarshalTypeError{
			Value: "string field",
			Type:  reflect.TypeOf(av),
			Err:   errors.New("attribute value is not string type"),
		}
	}

}
