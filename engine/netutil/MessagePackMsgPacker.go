package netutil

import (
	"bytes"

	"github.com/vmihailenco/msgpack"
)

// MessagePackMsgPacker packs and unpacks message in MessagePack format
type MessagePackMsgPacker struct{}

// PackMsg packs message to bytes in MessagePack format
func (mp MessagePackMsgPacker) PackMsg(msg interface{}, buf []byte) ([]byte, error) {
	buffer := bytes.NewBuffer(buf)

	encoder := msgpack.NewEncoder(buffer)
	err := encoder.Encode(msg)
	if err != nil {
		return buf, err
	}
	buf = buffer.Bytes()
	return buf, nil
}

// UnpackMsg unpacksbytes in MessagePack format to message
func (mp MessagePackMsgPacker) UnpackMsg(data []byte, msg interface{}) error {
	err := msgpack.Unmarshal(data, msg)
	return err
}

//func (mp MessagePackMsgPacker) convertToStringKeys(v interface{}) interface{} {
//	defer func() {
//		err := recover()
//		if err != nil {
//			gwlog.Errorf("MessagePackMsgPacker.convertToStringKeys failed")
//			panic(err)
//		}
//	}()
//
//	if rv, ok := v.(map[interface{}]interface{}); ok {
//		rrv := make(map[string]interface{}, len(rv))
//		for k, _v := range rv {
//			ks, ok := k.(string)
//			if !ok {
//				gwlog.Panicf("%v is not string, but %T", k, k)
//			}
//
//			rrv[ks] = mp.convertToStringKeys(_v)
//		}
//		return rrv
//	}
//
//	if rv, ok := v.(map[string]interface{}); ok {
//		for k, _v := range rv {
//			rv[k] = mp.convertToStringKeys(_v)
//		}
//		return rv
//	}
//
//	if rv, ok := v.([]interface{}); ok {
//		for i, _v := range rv {
//			rv[i] = mp.convertToStringKeys(_v)
//		}
//		return rv
//	}
//
//	return v
//}
