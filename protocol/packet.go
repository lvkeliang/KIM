package protocol

import (
	"KIM/protocol/protoImpl"
	"KIM/util"
	"bytes"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"io"
	"reflect"
	"strconv"
)

type Magic [4]byte

var (
	MagicLogicPkt = Magic{0xc3, 0x11, 0xa3, 0x65} //逻辑协议
	MagicBasicPkt = Magic{0xc3, 0x15, 0xa7, 0x65} //基础协议
)

// Packet 一个包需要实现的接口
type Packet interface {
	Decode(r io.Reader) error
	Encode(w io.Writer) error
}

type Flag int32
type Status int32

// MarshalPacket 封包，包的统一封包方法，通过反射判断包类型并写入对应的魔数，之后调用其Encode方法
func MarshalPacket(p Packet) []byte {
	buf := new(bytes.Buffer)
	kind := reflect.TypeOf(p).Elem()

	if kind.AssignableTo(reflect.TypeOf(LogicPkt{})) {
		_, _ = buf.Write(MagicLogicPkt[:])
	} else if kind.AssignableTo(reflect.TypeOf(BasicPkt{})) {
		_, _ = buf.Write(MagicBasicPkt[:])
	}
	_ = p.Encode(buf)
	return buf.Bytes()
}

// UnMarshalPacket 拆包，包的统一读取接口，根据包的魔数区分包并调用对应的Decode方法，在读取后再转换为对应的包结构
func UnMarshalPacket(r io.Reader) (interface{}, error) {
	magic := Magic{}
	_, err := io.ReadFull(r, magic[:])
	if err != nil {
		return nil, err
	}
	switch magic {
	case MagicLogicPkt:
		p := new(LogicPkt)
		if err := p.Decode(r); err != nil {
			return nil, err
		}
		return p, nil
	case MagicBasicPkt:
		p := new(BasicPkt)
		if err := p.Decode(r); err != nil {
			return nil, err
		}
		return p, nil
	default:
		return nil, errors.New("magic code is incorrect")
	}
}

func MustUnMarshalLogicPkt(r io.Reader) (*LogicPkt, error) {
	val, err := UnMarshalPacket(r)
	if err != nil {
		return nil, err
	}
	if lp, ok := val.(*LogicPkt); ok {
		return lp, nil
	}
	return nil, fmt.Errorf("packet is not a logic packet")
}

func MustUnMarshalBasicPkt(r io.Reader) (*BasicPkt, error) {
	val, err := UnMarshalPacket(r)
	if err != nil {
		return nil, err
	}
	if bp, ok := val.(*BasicPkt); ok {
		return bp, nil
	}
	return nil, fmt.Errorf("packet is not a basic packet")
}

// LogicPkt 逻辑协议,网关对外的消息结构包,逻辑协议中的Header和Body使用protobuf序列化框架
type LogicPkt struct {
	protoImpl.Header
	Body []byte
}

func (p *LogicPkt) Encode(w io.Writer) error {
	//protobuf编码header
	headerBytes, err := proto.Marshal(&p.Header)
	if err != nil {
		return err
	}
	if err := util.WriteBytes(w, headerBytes); err != nil {
		return err
	}
	if err := util.WriteBytes(w, p.Body); err != nil {
		return err
	}
	return nil
}

func (p *LogicPkt) Decode(r io.Reader) error {
	headerBytes, err := util.ReadBytes(r)
	if err != nil {
		return err
	}

	// protobuf解码header
	if err := proto.Unmarshal(headerBytes, &p.Header); err != nil {
		return err
	}
	p.Body, err = util.ReadBytes(r)
	if err != nil {
		return err
	}
	return nil
}

func (p *LogicPkt) AddMeta(m ...*protoImpl.Meta) {
	p.Meta = append(p.Meta, m...)
}

func (p *LogicPkt) AddStringMeta(key, value string) {
	p.AddMeta(&protoImpl.Meta{
		Key:   key,
		Value: value,
		Type:  protoImpl.MetaType_string,
	})
}

func (p *LogicPkt) AddIntMeta(key string, value int) {
	p.AddMeta(&protoImpl.Meta{
		Key:   key,
		Value: strconv.Itoa(value),
		Type:  protoImpl.MetaType_int,
	})
}

func (p *LogicPkt) AddFloatMeta(key string, value float64) {
	p.AddMeta(&protoImpl.Meta{
		Key:   key,
		Value: strconv.FormatFloat(value, 'E', -1, 64),
		Type:  protoImpl.MetaType_float,
	})
}

func (p *LogicPkt) DelMeta(key string) {
	for i, m := range p.Meta {
		if m.Key == key {
			length := len(p.Meta)
			if i < length-1 {
				copy(p.Meta[i:], p.Meta[i+1:])
			}
			p.Meta = p.Meta[:length-1]
		}
	}
}

// GetSpecificMeta 获取指定的Meta
func (p *LogicPkt) GetSpecificMeta(key string) (interface{}, bool) {
	return FindSpecificMeta(p.Meta, key)
}

func FindSpecificMeta(meta []*protoImpl.Meta, key string) (interface{}, bool) {
	for _, m := range meta {
		if m.Key == key {
			switch m.Type {
			case protoImpl.MetaType_int:
				v, _ := strconv.Atoi(m.Value)
				return v, true
			case protoImpl.MetaType_float:
				v, _ := strconv.ParseFloat(m.Value, 64)
				return v, true
			case protoImpl.MetaType_string:
				return m.Value, true
			}
		}
	}
	return nil, false
}

// BasicPkt 基础协议，更轻量，用于各服务间的心跳等包
type BasicPkt struct {
	Code   uint16
	Length uint16
	Body   []byte
}

func (p *BasicPkt) Decode(r io.Reader) error {
	var err error
	if p.Code, err = util.ReadUint16(r); err != nil {
		return err
	}
	if p.Length, err = util.ReadUint16(r); err != nil {
		return err
	}
	if p.Length > 0 {
		if p.Body, err = util.ReadFixedBytes(int(p.Length), r); err != nil {
			return err
		}
	}
	return nil
}

func (p *BasicPkt) Encode(w io.Writer) error {
	if err := util.WriteUint16(w, p.Code); err != nil {
		return err
	}
	if err := util.WriteUint16(w, p.Length); err != nil {
		return err
	}
	if p.Length > 0 {
		if _, err := w.Write(p.Body); err != nil {
			return err
		}
	}
	return nil
}
