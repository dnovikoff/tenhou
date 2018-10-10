// Code generated by protoc-gen-go. DO NOT EDIT.
// source: stats/stats.proto

package stats // import "github.com/dnovikoff/tenhou/protogen/stats"

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import duration "github.com/golang/protobuf/ptypes/duration"
import timestamp "github.com/golang/protobuf/ptypes/timestamp"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type LobbyType int32

const (
	LobbyType_LOBBY_TYPE_KU       LobbyType = 0
	LobbyType_LOBBY_TYPE_DAN      LobbyType = 1
	LobbyType_LOBBY_TYPE_UPPERDAN LobbyType = 2
	LobbyType_LOBBY_TYPE_PHOENIX  LobbyType = 3
	LobbyType_LOBBY_TYPE_DZ       LobbyType = 4
	LobbyType_LOBBY_TYPE_X1       LobbyType = 5
	LobbyType_LOBBY_TYPE_X2       LobbyType = 6
)

var LobbyType_name = map[int32]string{
	0: "LOBBY_TYPE_KU",
	1: "LOBBY_TYPE_DAN",
	2: "LOBBY_TYPE_UPPERDAN",
	3: "LOBBY_TYPE_PHOENIX",
	4: "LOBBY_TYPE_DZ",
	5: "LOBBY_TYPE_X1",
	6: "LOBBY_TYPE_X2",
}
var LobbyType_value = map[string]int32{
	"LOBBY_TYPE_KU":       0,
	"LOBBY_TYPE_DAN":      1,
	"LOBBY_TYPE_UPPERDAN": 2,
	"LOBBY_TYPE_PHOENIX":  3,
	"LOBBY_TYPE_DZ":       4,
	"LOBBY_TYPE_X1":       5,
	"LOBBY_TYPE_X2":       6,
}

func (x LobbyType) String() string {
	return proto.EnumName(LobbyType_name, int32(x))
}
func (LobbyType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_stats_eda27018f3e6cd88, []int{0}
}

type GameLength int32

const (
	GameLength_GAME_LENGTH_SOUTH GameLength = 0
	GameLength_GAME_LENGTH_EAST  GameLength = 1
	GameLength_GAME_LENGTH_ONE   GameLength = 2
)

var GameLength_name = map[int32]string{
	0: "GAME_LENGTH_SOUTH",
	1: "GAME_LENGTH_EAST",
	2: "GAME_LENGTH_ONE",
}
var GameLength_value = map[string]int32{
	"GAME_LENGTH_SOUTH": 0,
	"GAME_LENGTH_EAST":  1,
	"GAME_LENGTH_ONE":   2,
}

func (x GameLength) String() string {
	return proto.EnumName(GameLength_name, int32(x))
}
func (GameLength) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_stats_eda27018f3e6cd88, []int{1}
}

type GameType int32

const (
	GameType_GAME_TYPE_4 GameType = 0
	GameType_GAME_TYPE_3 GameType = 1
)

var GameType_name = map[int32]string{
	0: "GAME_TYPE_4",
	1: "GAME_TYPE_3",
}
var GameType_value = map[string]int32{
	"GAME_TYPE_4": 0,
	"GAME_TYPE_3": 1,
}

func (x GameType) String() string {
	return proto.EnumName(GameType_name, int32(x))
}
func (GameType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_stats_eda27018f3e6cd88, []int{2}
}

type GameSpeed int32

const (
	GameSpeed_GAME_SPEED_NORMAL GameSpeed = 0
	GameSpeed_GAME_SPEED_FAST   GameSpeed = 1
)

var GameSpeed_name = map[int32]string{
	0: "GAME_SPEED_NORMAL",
	1: "GAME_SPEED_FAST",
}
var GameSpeed_value = map[string]int32{
	"GAME_SPEED_NORMAL": 0,
	"GAME_SPEED_FAST":   1,
}

func (x GameSpeed) String() string {
	return proto.EnumName(GameSpeed_name, int32(x))
}
func (GameSpeed) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_stats_eda27018f3e6cd88, []int{3}
}

type Akkas int32

const (
	Akkas_AKKAS_YES Akkas = 0
	Akkas_AKKAS_NO  Akkas = 1
)

var Akkas_name = map[int32]string{
	0: "AKKAS_YES",
	1: "AKKAS_NO",
}
var Akkas_value = map[string]int32{
	"AKKAS_YES": 0,
	"AKKAS_NO":  1,
}

func (x Akkas) String() string {
	return proto.EnumName(Akkas_name, int32(x))
}
func (Akkas) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_stats_eda27018f3e6cd88, []int{4}
}

type Tanyao int32

const (
	Tanyao_TANYAO_YES Tanyao = 0
	Tanyao_TANYAO_NO  Tanyao = 1
)

var Tanyao_name = map[int32]string{
	0: "TANYAO_YES",
	1: "TANYAO_NO",
}
var Tanyao_value = map[string]int32{
	"TANYAO_YES": 0,
	"TANYAO_NO":  1,
}

func (x Tanyao) String() string {
	return proto.EnumName(Tanyao_name, int32(x))
}
func (Tanyao) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_stats_eda27018f3e6cd88, []int{5}
}

type Player struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Score                int64    `protobuf:"varint,2,opt,name=score,proto3" json:"score,omitempty"`
	Coins                int64    `protobuf:"varint,3,opt,name=coins,proto3" json:"coins,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Player) Reset()         { *m = Player{} }
func (m *Player) String() string { return proto.CompactTextString(m) }
func (*Player) ProtoMessage()    {}
func (*Player) Descriptor() ([]byte, []int) {
	return fileDescriptor_stats_eda27018f3e6cd88, []int{0}
}
func (m *Player) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Player.Unmarshal(m, b)
}
func (m *Player) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Player.Marshal(b, m, deterministic)
}
func (dst *Player) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Player.Merge(dst, src)
}
func (m *Player) XXX_Size() int {
	return xxx_messageInfo_Player.Size(m)
}
func (m *Player) XXX_DiscardUnknown() {
	xxx_messageInfo_Player.DiscardUnknown(m)
}

var xxx_messageInfo_Player proto.InternalMessageInfo

func (m *Player) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Player) GetScore() int64 {
	if m != nil {
		return m.Score
	}
	return 0
}

func (m *Player) GetCoins() int64 {
	if m != nil {
		return m.Coins
	}
	return 0
}

type Record struct {
	Time                 *timestamp.Timestamp `protobuf:"bytes,1,opt,name=time,proto3" json:"time,omitempty"`
	Duration             *duration.Duration   `protobuf:"bytes,2,opt,name=duration,proto3" json:"duration,omitempty"`
	Number               int64                `protobuf:"varint,3,opt,name=number,proto3" json:"number,omitempty"`
	Lobby                LobbyType            `protobuf:"varint,4,opt,name=lobby,proto3,enum=stats.LobbyType" json:"lobby,omitempty"`
	Length               GameLength           `protobuf:"varint,5,opt,name=length,proto3,enum=stats.GameLength" json:"length,omitempty"`
	Type                 GameType             `protobuf:"varint,6,opt,name=type,proto3,enum=stats.GameType" json:"type,omitempty"`
	Akkas                Akkas                `protobuf:"varint,7,opt,name=akkas,proto3,enum=stats.Akkas" json:"akkas,omitempty"`
	Tanyao               Tanyao               `protobuf:"varint,8,opt,name=tanyao,proto3,enum=stats.Tanyao" json:"tanyao,omitempty"`
	Players              []*Player            `protobuf:"bytes,9,rep,name=players,proto3" json:"players,omitempty"`
	Id                   string               `protobuf:"bytes,10,opt,name=id,proto3" json:"id,omitempty"`
	IsDz                 bool                 `protobuf:"varint,11,opt,name=is_dz,json=isDz,proto3" json:"is_dz,omitempty"`
	IsFive               bool                 `protobuf:"varint,12,opt,name=is_five,json=isFive,proto3" json:"is_five,omitempty"`
	IsChampionLobby      bool                 `protobuf:"varint,13,opt,name=is_champion_lobby,json=isChampionLobby,proto3" json:"is_champion_lobby,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *Record) Reset()         { *m = Record{} }
func (m *Record) String() string { return proto.CompactTextString(m) }
func (*Record) ProtoMessage()    {}
func (*Record) Descriptor() ([]byte, []int) {
	return fileDescriptor_stats_eda27018f3e6cd88, []int{1}
}
func (m *Record) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Record.Unmarshal(m, b)
}
func (m *Record) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Record.Marshal(b, m, deterministic)
}
func (dst *Record) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Record.Merge(dst, src)
}
func (m *Record) XXX_Size() int {
	return xxx_messageInfo_Record.Size(m)
}
func (m *Record) XXX_DiscardUnknown() {
	xxx_messageInfo_Record.DiscardUnknown(m)
}

var xxx_messageInfo_Record proto.InternalMessageInfo

func (m *Record) GetTime() *timestamp.Timestamp {
	if m != nil {
		return m.Time
	}
	return nil
}

func (m *Record) GetDuration() *duration.Duration {
	if m != nil {
		return m.Duration
	}
	return nil
}

func (m *Record) GetNumber() int64 {
	if m != nil {
		return m.Number
	}
	return 0
}

func (m *Record) GetLobby() LobbyType {
	if m != nil {
		return m.Lobby
	}
	return LobbyType_LOBBY_TYPE_KU
}

func (m *Record) GetLength() GameLength {
	if m != nil {
		return m.Length
	}
	return GameLength_GAME_LENGTH_SOUTH
}

func (m *Record) GetType() GameType {
	if m != nil {
		return m.Type
	}
	return GameType_GAME_TYPE_4
}

func (m *Record) GetAkkas() Akkas {
	if m != nil {
		return m.Akkas
	}
	return Akkas_AKKAS_YES
}

func (m *Record) GetTanyao() Tanyao {
	if m != nil {
		return m.Tanyao
	}
	return Tanyao_TANYAO_YES
}

func (m *Record) GetPlayers() []*Player {
	if m != nil {
		return m.Players
	}
	return nil
}

func (m *Record) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Record) GetIsDz() bool {
	if m != nil {
		return m.IsDz
	}
	return false
}

func (m *Record) GetIsFive() bool {
	if m != nil {
		return m.IsFive
	}
	return false
}

func (m *Record) GetIsChampionLobby() bool {
	if m != nil {
		return m.IsChampionLobby
	}
	return false
}

func init() {
	proto.RegisterType((*Player)(nil), "stats.Player")
	proto.RegisterType((*Record)(nil), "stats.Record")
	proto.RegisterEnum("stats.LobbyType", LobbyType_name, LobbyType_value)
	proto.RegisterEnum("stats.GameLength", GameLength_name, GameLength_value)
	proto.RegisterEnum("stats.GameType", GameType_name, GameType_value)
	proto.RegisterEnum("stats.GameSpeed", GameSpeed_name, GameSpeed_value)
	proto.RegisterEnum("stats.Akkas", Akkas_name, Akkas_value)
	proto.RegisterEnum("stats.Tanyao", Tanyao_name, Tanyao_value)
}

func init() { proto.RegisterFile("stats/stats.proto", fileDescriptor_stats_eda27018f3e6cd88) }

var fileDescriptor_stats_eda27018f3e6cd88 = []byte{
	// 653 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x64, 0x94, 0xdf, 0x6e, 0xda, 0x4a,
	0x10, 0x87, 0x31, 0x60, 0x07, 0x86, 0x00, 0xcb, 0x24, 0x27, 0xf1, 0xc9, 0x45, 0x8b, 0xd2, 0x3f,
	0xa1, 0x56, 0x04, 0x2a, 0x69, 0xd5, 0x6b, 0xa7, 0x38, 0xa1, 0x0a, 0x31, 0xc8, 0x38, 0x52, 0xc8,
	0x8d, 0x65, 0x60, 0x81, 0x55, 0xc0, 0x8b, 0xb0, 0x89, 0x44, 0x9e, 0xa5, 0xaf, 0xd6, 0x77, 0xa9,
	0xd8, 0x35, 0x94, 0xd2, 0x1b, 0xc4, 0x7c, 0xbf, 0xcf, 0x3b, 0xa3, 0x1d, 0xcb, 0x50, 0x0a, 0x23,
	0x3f, 0x0a, 0x6b, 0xe2, 0xb7, 0x3a, 0x5f, 0xf0, 0x88, 0xa3, 0x2a, 0x8a, 0xb3, 0xb7, 0x63, 0xce,
	0xc7, 0x53, 0x5a, 0x13, 0xb0, 0xbf, 0x1c, 0xd5, 0x22, 0x36, 0xa3, 0x61, 0xe4, 0xcf, 0xe6, 0xd2,
	0x3b, 0x7b, 0xb3, 0x2f, 0x0c, 0x97, 0x0b, 0x3f, 0x62, 0x3c, 0x90, 0xf9, 0x79, 0x13, 0xb4, 0xce,
	0xd4, 0x5f, 0xd1, 0x05, 0x22, 0xa4, 0x03, 0x7f, 0x46, 0x75, 0xa5, 0xac, 0x54, 0xb2, 0x8e, 0xf8,
	0x8f, 0xc7, 0xa0, 0x86, 0x03, 0xbe, 0xa0, 0x7a, 0xb2, 0xac, 0x54, 0x52, 0x8e, 0x2c, 0xd6, 0x74,
	0xc0, 0x59, 0x10, 0xea, 0x29, 0x49, 0x45, 0x71, 0xfe, 0x2b, 0x05, 0x9a, 0x43, 0x07, 0x7c, 0x31,
	0xc4, 0x2a, 0xa4, 0xd7, 0x73, 0x88, 0xa3, 0x72, 0xf5, 0xb3, 0xaa, 0x9c, 0xa1, 0xba, 0x99, 0xa1,
	0xea, 0x6e, 0x86, 0x74, 0x84, 0x87, 0x5f, 0x21, 0xb3, 0x19, 0x4b, 0x74, 0xca, 0xd5, 0xff, 0xff,
	0xe7, 0x99, 0x46, 0x2c, 0x38, 0x5b, 0x15, 0x4f, 0x40, 0x0b, 0x96, 0xb3, 0x3e, 0x5d, 0xc4, 0x83,
	0xc4, 0x15, 0x7e, 0x04, 0x75, 0xca, 0xfb, 0xfd, 0x95, 0x9e, 0x2e, 0x2b, 0x95, 0x42, 0x9d, 0x54,
	0xe5, 0xc5, 0xb5, 0xd6, 0xcc, 0x5d, 0xcd, 0xa9, 0x23, 0x63, 0xfc, 0x04, 0xda, 0x94, 0x06, 0xe3,
	0x68, 0xa2, 0xab, 0x42, 0x2c, 0xc5, 0xe2, 0xad, 0x3f, 0xa3, 0x2d, 0x11, 0x38, 0xb1, 0x80, 0xef,
	0x20, 0x1d, 0xad, 0xe6, 0x54, 0xd7, 0x84, 0x58, 0xdc, 0x11, 0xc5, 0x81, 0x22, 0xc4, 0x73, 0x50,
	0xfd, 0xe7, 0x67, 0x3f, 0xd4, 0x0f, 0x84, 0x75, 0x18, 0x5b, 0xe6, 0x9a, 0x39, 0x32, 0xc2, 0x0f,
	0xa0, 0x45, 0x7e, 0xb0, 0xf2, 0xb9, 0x9e, 0x11, 0x52, 0x3e, 0x96, 0x5c, 0x01, 0x9d, 0x38, 0xc4,
	0x0b, 0x38, 0x98, 0x8b, 0xb5, 0x84, 0x7a, 0xb6, 0x9c, 0xaa, 0xe4, 0xb6, 0x9e, 0x5c, 0x96, 0xb3,
	0x49, 0xb1, 0x00, 0x49, 0x36, 0xd4, 0x41, 0xec, 0x2c, 0xc9, 0x86, 0x78, 0x04, 0x2a, 0x0b, 0xbd,
	0xe1, 0xab, 0x9e, 0x2b, 0x2b, 0x95, 0x8c, 0x93, 0x66, 0x61, 0xe3, 0x15, 0x4f, 0xe1, 0x80, 0x85,
	0xde, 0x88, 0xbd, 0x50, 0xfd, 0x50, 0x60, 0x8d, 0x85, 0x37, 0xec, 0x85, 0xa2, 0x01, 0x25, 0x16,
	0x7a, 0x83, 0x89, 0x3f, 0x9b, 0x33, 0x1e, 0x78, 0xf2, 0xd6, 0xf2, 0x42, 0x29, 0xb2, 0xf0, 0x7b,
	0xcc, 0xc5, 0xc5, 0x19, 0x3f, 0x15, 0xc8, 0x6e, 0xaf, 0x10, 0x4b, 0x90, 0x6f, 0xb5, 0xaf, 0xaf,
	0x7b, 0x9e, 0xdb, 0xeb, 0x58, 0xde, 0xdd, 0x03, 0x49, 0x20, 0x42, 0x61, 0x07, 0x35, 0x4c, 0x9b,
	0x28, 0x78, 0x0a, 0x47, 0x3b, 0xec, 0xa1, 0xd3, 0xb1, 0x9c, 0x75, 0x90, 0xc4, 0x13, 0xc0, 0x9d,
	0xa0, 0xd3, 0x6c, 0x5b, 0xf6, 0x8f, 0x47, 0x92, 0xda, 0x3b, 0xb7, 0xf1, 0x44, 0xd2, 0x7b, 0xe8,
	0xf1, 0x33, 0x51, 0xf7, 0x51, 0x9d, 0x68, 0x86, 0x0d, 0xf0, 0x67, 0x6f, 0xf8, 0x1f, 0x94, 0x6e,
	0xcd, 0x7b, 0xcb, 0x6b, 0x59, 0xf6, 0xad, 0xdb, 0xf4, 0xba, 0xed, 0x07, 0xb7, 0x49, 0x12, 0x78,
	0x0c, 0x64, 0x17, 0x5b, 0x66, 0xd7, 0x25, 0x0a, 0x1e, 0x41, 0x71, 0x97, 0xb6, 0x6d, 0x8b, 0x24,
	0x8d, 0x4b, 0xc8, 0x6c, 0xd6, 0x8b, 0x45, 0xc8, 0x09, 0x41, 0x74, 0xfb, 0x42, 0x12, 0x7f, 0x83,
	0x2b, 0xa2, 0x18, 0xdf, 0x20, 0xbb, 0xb6, 0xbb, 0x73, 0x4a, 0x87, 0xdb, 0xe6, 0xdd, 0x8e, 0x65,
	0x35, 0x3c, 0xbb, 0xed, 0xdc, 0x9b, 0x2d, 0x92, 0xd8, 0xb6, 0x91, 0xf8, 0x46, 0xf4, 0x36, 0xde,
	0x83, 0x2a, 0xde, 0x0f, 0xcc, 0x43, 0xd6, 0xbc, 0xbb, 0x33, 0xbb, 0x5e, 0xcf, 0xea, 0x92, 0x04,
	0x1e, 0x42, 0x46, 0x96, 0x76, 0x9b, 0x28, 0xc6, 0x05, 0x68, 0xf2, 0x05, 0xc1, 0x02, 0x80, 0x6b,
	0xda, 0x3d, 0xb3, 0x1d, 0x7b, 0x79, 0xc8, 0xc6, 0xf5, 0x5a, 0xbc, 0xbe, 0x7c, 0x32, 0xc6, 0x2c,
	0x9a, 0x2c, 0xfb, 0xd5, 0x01, 0x9f, 0xd5, 0x86, 0x01, 0x7f, 0x61, 0xcf, 0x7c, 0x34, 0xaa, 0x45,
	0x34, 0x98, 0xf0, 0xa5, 0xfc, 0x0a, 0x8c, 0x69, 0x20, 0x3f, 0x25, 0x7d, 0x4d, 0xd4, 0x57, 0xbf,
	0x03, 0x00, 0x00, 0xff, 0xff, 0x66, 0x32, 0x13, 0x61, 0x60, 0x04, 0x00, 0x00,
}