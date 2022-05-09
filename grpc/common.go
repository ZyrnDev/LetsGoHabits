package grpc

import (
	"context"
	"reflect"
	"strings"

	"github.com/ZyrnDev/letsgohabits/database"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func isRequired(field reflect.StructField, operationType string) bool {
	requiredTag := field.Tag.Get("required")
	requiredOperations := strings.Split(requiredTag, ",")

	for _, operation := range requiredOperations {
		if operation == operationType {
			return true
		}
	}
	return false
}

func requireFields(data interface{}, operationType string) error {
	typ := reflect.TypeOf(data)
	val := reflect.ValueOf(data)
	for i := 0; i < typ.NumField(); i++ {
		if isRequired(typ.Field(i), operationType) && isEmptyValue(val.Field(i)) {
			return status.Errorf(codes.InvalidArgument, "field %s is required", typ.Field(i).Name)
		}
		if typ.Field(i).Type.Kind() == reflect.Struct {
			err := requireFields(val.Field(i).Interface(), operationType)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	default:
		return false
	}
}

type ProtobufSerialDatabaseObject[ProtobufType any] interface {
	database.HasID
	database.ProtobufSerial[ProtobufType]
}

func New[ProtobufType any, DatabaseType database.HasID, DatabaseTypePtr ProtobufSerialDatabaseObject[ProtobufType]](ctx context.Context, db database.Database, in *ProtobufType) (*ProtobufType, error) {
	var input DatabaseType
	var inputPtrAny interface{}
	inputPtrAny = &input
	inputPtr := (inputPtrAny).(DatabaseTypePtr)

	inputPtr.FromProtobuf(in)

	requireFields(input, "create")

	if res := db.Create(&input); res.Error != nil {
		return nil, status.Errorf(codes.Internal, "failed to create oject: %s", res.Error)
	}

	return inputPtr.ToProtobuf(), nil
}

func Get[ProtobufType any, DatabaseType database.HasID, DatabaseTypePtr ProtobufSerialDatabaseObject[ProtobufType]](ctx context.Context, db database.Database, in *ProtobufType) ([]*ProtobufType, error) {
	var input DatabaseType
	var inputPtrAny interface{}
	inputPtrAny = &input
	inputPtr := (inputPtrAny).(DatabaseTypePtr)

	inputPtr.FromProtobuf(in)

	requireFields(input, "get")

	var results []DatabaseType

	res := db.Find(&results, input)
	if res.Error != nil {
		return nil, status.Errorf(codes.Internal, "failed to get object: %s", res.Error)
	}

	resultsProto := make([]*ProtobufType, len(results))
	for i, object := range results {
		var objectAny interface{}
		objectAny = &object
		resultsProto[i] = (objectAny).(DatabaseTypePtr).ToProtobuf()
	}

	return resultsProto, nil
}

func Update[ProtobufType any, DatabaseType database.HasID, DatabaseTypePtr ProtobufSerialDatabaseObject[ProtobufType]](ctx context.Context, db database.Database, in *ProtobufType) (*ProtobufType, error) {
	var input DatabaseType
	var inputPtrAny interface{}
	inputPtrAny = &input
	inputPtr := (inputPtrAny).(DatabaseTypePtr)

	inputPtr.FromProtobuf(in)

	requireFields(input, "update")

	var updatedInput DatabaseType
	var updatedInputPtrAny interface{}
	updatedInputPtrAny = &updatedInput
	updatedInputPtr := (updatedInputPtrAny).(DatabaseTypePtr)

	if res := db.First(&updatedInput, input.GetID()); res.Error != nil {
		return nil, status.Errorf(codes.Internal, "failed to get object: %s", res.Error)
	}

	if res := db.Model(&updatedInput).Updates(input); res.Error != nil {
		return nil, status.Errorf(codes.Internal, "failed to update object: %s", res.Error)
	}

	if res := db.First(&updatedInput, input.GetID()); res.Error != nil {
		return nil, status.Errorf(codes.Internal, "failed to find object after update: %s", res.Error)
	}

	return updatedInputPtr.ToProtobuf(), nil
}

func Delete[ProtobufType any, DatabaseType database.HasID, DatabaseTypePtr ProtobufSerialDatabaseObject[ProtobufType]](ctx context.Context, db database.Database, in *ProtobufType) (*empty.Empty, error) {
	var input DatabaseType
	var inputPtrAny interface{}
	inputPtrAny = &input
	inputPtr := (inputPtrAny).(DatabaseTypePtr)

	inputPtr.FromProtobuf(in)
	requireFields(input, "delete")

	result := db.Unscoped().Delete(&input)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete object: %s", result.Error)
	}

	return &empty.Empty{}, nil
}
