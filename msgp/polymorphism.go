package msgp

type MsgPackDeserializer interface {
	Unmarshaler
	Decodable
}

// PolymorphicResolver interface can be added to any model with polymorphic fields (interface{}s)
// allowing them to unmarshal themselves for MessagePack using a runtime hint (via a discriminator field).
// This greatly reduces amount of custom unmarshaling code needed by users of msgpack.
type PolymorphicResolver interface {
	ChooseType(field string) (MsgPackDeserializer, error)
}

func ResolveAndUnmarshalMsg(resolvable PolymorphicResolver, field string, bts []byte) (interface{}, []byte, error) {
	obj, err := resolvable.ChooseType(field)
	if err != nil {
		return nil, bts, err
	}

	// if polymorphic object is nil, drain nil from serialized stream first
	if IsNil(bts) || obj == nil {
		bts, err = ReadNilBytes(bts)
		if err != nil {
			return nil, bts, WrapError(err, "unable to read nil object from byte stream")
		}

		return nil, bts, nil
	}

	bts, err = obj.UnmarshalMsg(bts)
	if err != nil {
		return obj, bts, WrapError(err, "unable to unmarshal "+field)
	}

	return obj, bts, nil
}

func ResolveAndDecodeMsg(resolvable PolymorphicResolver, field string, reader *Reader) (interface{}, error) {
	obj, err := resolvable.ChooseType(field)
	if err != nil {
		return nil, err
	}

	// if polymorphic object is nil, drain nil from reader first
	if obj == nil || reader.IsNil() {
		err = reader.ReadNil()
		if err != nil {
			return nil, WrapError(err, "unable to read nil object from reader")
		}
		return nil, nil
	}

	err = obj.DecodeMsg(reader)
	if err != nil {
		return obj, WrapError(err, "unable to decode "+field)
	}

	return obj, nil
}
