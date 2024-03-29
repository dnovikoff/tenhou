// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.21.2
// source: stats/stats.proto

package stats

import (
	proto "github.com/golang/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type LobbyType int32

const (
	LobbyType_LOBBY_TYPE_KU       LobbyType = 0 // 般
	LobbyType_LOBBY_TYPE_DAN      LobbyType = 1 // 上
	LobbyType_LOBBY_TYPE_UPPERDAN LobbyType = 2 // 特
	LobbyType_LOBBY_TYPE_PHOENIX  LobbyType = 3 // 鳳
	LobbyType_LOBBY_TYPE_DZ       LobbyType = 4 // 技
	LobbyType_LOBBY_TYPE_X1       LobbyType = 5 // 若
	LobbyType_LOBBY_TYPE_X2       LobbyType = 6 // 銀
	LobbyType_LOBBY_TYPE_X3       LobbyType = 7 // 琥
	LobbyType_LOBBY_TYPE_EXTERNAL LobbyType = 8 // －
)

// Enum value maps for LobbyType.
var (
	LobbyType_name = map[int32]string{
		0: "LOBBY_TYPE_KU",
		1: "LOBBY_TYPE_DAN",
		2: "LOBBY_TYPE_UPPERDAN",
		3: "LOBBY_TYPE_PHOENIX",
		4: "LOBBY_TYPE_DZ",
		5: "LOBBY_TYPE_X1",
		6: "LOBBY_TYPE_X2",
		7: "LOBBY_TYPE_X3",
		8: "LOBBY_TYPE_EXTERNAL",
	}
	LobbyType_value = map[string]int32{
		"LOBBY_TYPE_KU":       0,
		"LOBBY_TYPE_DAN":      1,
		"LOBBY_TYPE_UPPERDAN": 2,
		"LOBBY_TYPE_PHOENIX":  3,
		"LOBBY_TYPE_DZ":       4,
		"LOBBY_TYPE_X1":       5,
		"LOBBY_TYPE_X2":       6,
		"LOBBY_TYPE_X3":       7,
		"LOBBY_TYPE_EXTERNAL": 8,
	}
)

func (x LobbyType) Enum() *LobbyType {
	p := new(LobbyType)
	*p = x
	return p
}

func (x LobbyType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (LobbyType) Descriptor() protoreflect.EnumDescriptor {
	return file_stats_stats_proto_enumTypes[0].Descriptor()
}

func (LobbyType) Type() protoreflect.EnumType {
	return &file_stats_stats_proto_enumTypes[0]
}

func (x LobbyType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use LobbyType.Descriptor instead.
func (LobbyType) EnumDescriptor() ([]byte, []int) {
	return file_stats_stats_proto_rawDescGZIP(), []int{0}
}

type GameLength int32

const (
	GameLength_GAME_LENGTH_SOUTH GameLength = 0
	GameLength_GAME_LENGTH_EAST  GameLength = 1
	GameLength_GAME_LENGTH_ONE   GameLength = 2
)

// Enum value maps for GameLength.
var (
	GameLength_name = map[int32]string{
		0: "GAME_LENGTH_SOUTH",
		1: "GAME_LENGTH_EAST",
		2: "GAME_LENGTH_ONE",
	}
	GameLength_value = map[string]int32{
		"GAME_LENGTH_SOUTH": 0,
		"GAME_LENGTH_EAST":  1,
		"GAME_LENGTH_ONE":   2,
	}
)

func (x GameLength) Enum() *GameLength {
	p := new(GameLength)
	*p = x
	return p
}

func (x GameLength) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (GameLength) Descriptor() protoreflect.EnumDescriptor {
	return file_stats_stats_proto_enumTypes[1].Descriptor()
}

func (GameLength) Type() protoreflect.EnumType {
	return &file_stats_stats_proto_enumTypes[1]
}

func (x GameLength) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use GameLength.Descriptor instead.
func (GameLength) EnumDescriptor() ([]byte, []int) {
	return file_stats_stats_proto_rawDescGZIP(), []int{1}
}

type GameType int32

const (
	GameType_GAME_TYPE_4 GameType = 0
	GameType_GAME_TYPE_3 GameType = 1
)

// Enum value maps for GameType.
var (
	GameType_name = map[int32]string{
		0: "GAME_TYPE_4",
		1: "GAME_TYPE_3",
	}
	GameType_value = map[string]int32{
		"GAME_TYPE_4": 0,
		"GAME_TYPE_3": 1,
	}
)

func (x GameType) Enum() *GameType {
	p := new(GameType)
	*p = x
	return p
}

func (x GameType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (GameType) Descriptor() protoreflect.EnumDescriptor {
	return file_stats_stats_proto_enumTypes[2].Descriptor()
}

func (GameType) Type() protoreflect.EnumType {
	return &file_stats_stats_proto_enumTypes[2]
}

func (x GameType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use GameType.Descriptor instead.
func (GameType) EnumDescriptor() ([]byte, []int) {
	return file_stats_stats_proto_rawDescGZIP(), []int{2}
}

type GameSpeed int32

const (
	GameSpeed_GAME_SPEED_NORMAL GameSpeed = 0
	GameSpeed_GAME_SPEED_FAST   GameSpeed = 1
)

// Enum value maps for GameSpeed.
var (
	GameSpeed_name = map[int32]string{
		0: "GAME_SPEED_NORMAL",
		1: "GAME_SPEED_FAST",
	}
	GameSpeed_value = map[string]int32{
		"GAME_SPEED_NORMAL": 0,
		"GAME_SPEED_FAST":   1,
	}
)

func (x GameSpeed) Enum() *GameSpeed {
	p := new(GameSpeed)
	*p = x
	return p
}

func (x GameSpeed) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (GameSpeed) Descriptor() protoreflect.EnumDescriptor {
	return file_stats_stats_proto_enumTypes[3].Descriptor()
}

func (GameSpeed) Type() protoreflect.EnumType {
	return &file_stats_stats_proto_enumTypes[3]
}

func (x GameSpeed) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use GameSpeed.Descriptor instead.
func (GameSpeed) EnumDescriptor() ([]byte, []int) {
	return file_stats_stats_proto_rawDescGZIP(), []int{3}
}

type Akkas int32

const (
	Akkas_AKKAS_YES Akkas = 0
	Akkas_AKKAS_NO  Akkas = 1
)

// Enum value maps for Akkas.
var (
	Akkas_name = map[int32]string{
		0: "AKKAS_YES",
		1: "AKKAS_NO",
	}
	Akkas_value = map[string]int32{
		"AKKAS_YES": 0,
		"AKKAS_NO":  1,
	}
)

func (x Akkas) Enum() *Akkas {
	p := new(Akkas)
	*p = x
	return p
}

func (x Akkas) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Akkas) Descriptor() protoreflect.EnumDescriptor {
	return file_stats_stats_proto_enumTypes[4].Descriptor()
}

func (Akkas) Type() protoreflect.EnumType {
	return &file_stats_stats_proto_enumTypes[4]
}

func (x Akkas) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Akkas.Descriptor instead.
func (Akkas) EnumDescriptor() ([]byte, []int) {
	return file_stats_stats_proto_rawDescGZIP(), []int{4}
}

type Tanyao int32

const (
	Tanyao_TANYAO_YES Tanyao = 0
	Tanyao_TANYAO_NO  Tanyao = 1
)

// Enum value maps for Tanyao.
var (
	Tanyao_name = map[int32]string{
		0: "TANYAO_YES",
		1: "TANYAO_NO",
	}
	Tanyao_value = map[string]int32{
		"TANYAO_YES": 0,
		"TANYAO_NO":  1,
	}
)

func (x Tanyao) Enum() *Tanyao {
	p := new(Tanyao)
	*p = x
	return p
}

func (x Tanyao) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Tanyao) Descriptor() protoreflect.EnumDescriptor {
	return file_stats_stats_proto_enumTypes[5].Descriptor()
}

func (Tanyao) Type() protoreflect.EnumType {
	return &file_stats_stats_proto_enumTypes[5]
}

func (x Tanyao) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Tanyao.Descriptor instead.
func (Tanyao) EnumDescriptor() ([]byte, []int) {
	return file_stats_stats_proto_rawDescGZIP(), []int{5}
}

type NumberType int32

const (
	NumberType_NO_NUMBER NumberType = 0
	NumberType_NUMBER_2  NumberType = 2
	NumberType_NUMBER_5  NumberType = 5
	NumberType_NUMBER_0  NumberType = 10
)

// Enum value maps for NumberType.
var (
	NumberType_name = map[int32]string{
		0:  "NO_NUMBER",
		2:  "NUMBER_2",
		5:  "NUMBER_5",
		10: "NUMBER_0",
	}
	NumberType_value = map[string]int32{
		"NO_NUMBER": 0,
		"NUMBER_2":  2,
		"NUMBER_5":  5,
		"NUMBER_0":  10,
	}
)

func (x NumberType) Enum() *NumberType {
	p := new(NumberType)
	*p = x
	return p
}

func (x NumberType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (NumberType) Descriptor() protoreflect.EnumDescriptor {
	return file_stats_stats_proto_enumTypes[6].Descriptor()
}

func (NumberType) Type() protoreflect.EnumType {
	return &file_stats_stats_proto_enumTypes[6]
}

func (x NumberType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use NumberType.Descriptor instead.
func (NumberType) EnumDescriptor() ([]byte, []int) {
	return file_stats_stats_proto_rawDescGZIP(), []int{6}
}

type Player struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name  string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Score int64  `protobuf:"varint,2,opt,name=score,proto3" json:"score,omitempty"`
	Coins int64  `protobuf:"varint,3,opt,name=coins,proto3" json:"coins,omitempty"`
}

func (x *Player) Reset() {
	*x = Player{}
	if protoimpl.UnsafeEnabled {
		mi := &file_stats_stats_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Player) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Player) ProtoMessage() {}

func (x *Player) ProtoReflect() protoreflect.Message {
	mi := &file_stats_stats_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Player.ProtoReflect.Descriptor instead.
func (*Player) Descriptor() ([]byte, []int) {
	return file_stats_stats_proto_rawDescGZIP(), []int{0}
}

func (x *Player) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Player) GetScore() int64 {
	if x != nil {
		return x.Score
	}
	return 0
}

func (x *Player) GetCoins() int64 {
	if x != nil {
		return x.Coins
	}
	return 0
}

type Record struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Time     *timestamppb.Timestamp `protobuf:"bytes,1,opt,name=time,proto3" json:"time,omitempty"`
	Duration *durationpb.Duration   `protobuf:"bytes,2,opt,name=duration,proto3" json:"duration,omitempty"`
	Number   int64                  `protobuf:"varint,3,opt,name=number,proto3" json:"number,omitempty"`
	Lobby    LobbyType              `protobuf:"varint,4,opt,name=lobby,proto3,enum=stats.LobbyType" json:"lobby,omitempty"`
	Length   GameLength             `protobuf:"varint,5,opt,name=length,proto3,enum=stats.GameLength" json:"length,omitempty"`
	Type     GameType               `protobuf:"varint,6,opt,name=type,proto3,enum=stats.GameType" json:"type,omitempty"`
	Akkas    Akkas                  `protobuf:"varint,7,opt,name=akkas,proto3,enum=stats.Akkas" json:"akkas,omitempty"`
	Tanyao   Tanyao                 `protobuf:"varint,8,opt,name=tanyao,proto3,enum=stats.Tanyao" json:"tanyao,omitempty"`
	// Some strange numbers, that are not ascii numbers like '５' or '２'
	NumberType      NumberType `protobuf:"varint,9,opt,name=number_type,json=numberType,proto3,enum=stats.NumberType" json:"number_type,omitempty"`
	Players         []*Player  `protobuf:"bytes,10,rep,name=players,proto3" json:"players,omitempty"`
	Id              string     `protobuf:"bytes,11,opt,name=id,proto3" json:"id,omitempty"`
	IsDz            bool       `protobuf:"varint,12,opt,name=is_dz,json=isDz,proto3" json:"is_dz,omitempty"`
	IsChampionLobby bool       `protobuf:"varint,13,opt,name=is_champion_lobby,json=isChampionLobby,proto3" json:"is_champion_lobby,omitempty"`
}

func (x *Record) Reset() {
	*x = Record{}
	if protoimpl.UnsafeEnabled {
		mi := &file_stats_stats_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Record) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Record) ProtoMessage() {}

func (x *Record) ProtoReflect() protoreflect.Message {
	mi := &file_stats_stats_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Record.ProtoReflect.Descriptor instead.
func (*Record) Descriptor() ([]byte, []int) {
	return file_stats_stats_proto_rawDescGZIP(), []int{1}
}

func (x *Record) GetTime() *timestamppb.Timestamp {
	if x != nil {
		return x.Time
	}
	return nil
}

func (x *Record) GetDuration() *durationpb.Duration {
	if x != nil {
		return x.Duration
	}
	return nil
}

func (x *Record) GetNumber() int64 {
	if x != nil {
		return x.Number
	}
	return 0
}

func (x *Record) GetLobby() LobbyType {
	if x != nil {
		return x.Lobby
	}
	return LobbyType_LOBBY_TYPE_KU
}

func (x *Record) GetLength() GameLength {
	if x != nil {
		return x.Length
	}
	return GameLength_GAME_LENGTH_SOUTH
}

func (x *Record) GetType() GameType {
	if x != nil {
		return x.Type
	}
	return GameType_GAME_TYPE_4
}

func (x *Record) GetAkkas() Akkas {
	if x != nil {
		return x.Akkas
	}
	return Akkas_AKKAS_YES
}

func (x *Record) GetTanyao() Tanyao {
	if x != nil {
		return x.Tanyao
	}
	return Tanyao_TANYAO_YES
}

func (x *Record) GetNumberType() NumberType {
	if x != nil {
		return x.NumberType
	}
	return NumberType_NO_NUMBER
}

func (x *Record) GetPlayers() []*Player {
	if x != nil {
		return x.Players
	}
	return nil
}

func (x *Record) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Record) GetIsDz() bool {
	if x != nil {
		return x.IsDz
	}
	return false
}

func (x *Record) GetIsChampionLobby() bool {
	if x != nil {
		return x.IsChampionLobby
	}
	return false
}

var File_stats_stats_proto protoreflect.FileDescriptor

var file_stats_stats_proto_rawDesc = []byte{
	0x0a, 0x11, 0x73, 0x74, 0x61, 0x74, 0x73, 0x2f, 0x73, 0x74, 0x61, 0x74, 0x73, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x05, 0x73, 0x74, 0x61, 0x74, 0x73, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x75, 0x72,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x48, 0x0a, 0x06, 0x50,
	0x6c, 0x61, 0x79, 0x65, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x63, 0x6f,
	0x72, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x73, 0x63, 0x6f, 0x72, 0x65, 0x12,
	0x14, 0x0a, 0x05, 0x63, 0x6f, 0x69, 0x6e, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05,
	0x63, 0x6f, 0x69, 0x6e, 0x73, 0x22, 0xf8, 0x03, 0x0a, 0x06, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64,
	0x12, 0x2e, 0x0a, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x04, 0x74, 0x69, 0x6d, 0x65,
	0x12, 0x35, 0x0a, 0x08, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x08, 0x64,
	0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x16, 0x0a, 0x06, 0x6e, 0x75, 0x6d, 0x62, 0x65,
	0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12,
	0x26, 0x0a, 0x05, 0x6c, 0x6f, 0x62, 0x62, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x10,
	0x2e, 0x73, 0x74, 0x61, 0x74, 0x73, 0x2e, 0x4c, 0x6f, 0x62, 0x62, 0x79, 0x54, 0x79, 0x70, 0x65,
	0x52, 0x05, 0x6c, 0x6f, 0x62, 0x62, 0x79, 0x12, 0x29, 0x0a, 0x06, 0x6c, 0x65, 0x6e, 0x67, 0x74,
	0x68, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x11, 0x2e, 0x73, 0x74, 0x61, 0x74, 0x73, 0x2e,
	0x47, 0x61, 0x6d, 0x65, 0x4c, 0x65, 0x6e, 0x67, 0x74, 0x68, 0x52, 0x06, 0x6c, 0x65, 0x6e, 0x67,
	0x74, 0x68, 0x12, 0x23, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0e,
	0x32, 0x0f, 0x2e, 0x73, 0x74, 0x61, 0x74, 0x73, 0x2e, 0x47, 0x61, 0x6d, 0x65, 0x54, 0x79, 0x70,
	0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x22, 0x0a, 0x05, 0x61, 0x6b, 0x6b, 0x61, 0x73,
	0x18, 0x07, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0c, 0x2e, 0x73, 0x74, 0x61, 0x74, 0x73, 0x2e, 0x41,
	0x6b, 0x6b, 0x61, 0x73, 0x52, 0x05, 0x61, 0x6b, 0x6b, 0x61, 0x73, 0x12, 0x25, 0x0a, 0x06, 0x74,
	0x61, 0x6e, 0x79, 0x61, 0x6f, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0d, 0x2e, 0x73, 0x74,
	0x61, 0x74, 0x73, 0x2e, 0x54, 0x61, 0x6e, 0x79, 0x61, 0x6f, 0x52, 0x06, 0x74, 0x61, 0x6e, 0x79,
	0x61, 0x6f, 0x12, 0x32, 0x0a, 0x0b, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x5f, 0x74, 0x79, 0x70,
	0x65, 0x18, 0x09, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x11, 0x2e, 0x73, 0x74, 0x61, 0x74, 0x73, 0x2e,
	0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x54, 0x79, 0x70, 0x65, 0x52, 0x0a, 0x6e, 0x75, 0x6d, 0x62,
	0x65, 0x72, 0x54, 0x79, 0x70, 0x65, 0x12, 0x27, 0x0a, 0x07, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72,
	0x73, 0x18, 0x0a, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x73, 0x74, 0x61, 0x74, 0x73, 0x2e,
	0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x52, 0x07, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x73, 0x12,
	0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12,
	0x13, 0x0a, 0x05, 0x69, 0x73, 0x5f, 0x64, 0x7a, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x08, 0x52, 0x04,
	0x69, 0x73, 0x44, 0x7a, 0x12, 0x2a, 0x0a, 0x11, 0x69, 0x73, 0x5f, 0x63, 0x68, 0x61, 0x6d, 0x70,
	0x69, 0x6f, 0x6e, 0x5f, 0x6c, 0x6f, 0x62, 0x62, 0x79, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x0f, 0x69, 0x73, 0x43, 0x68, 0x61, 0x6d, 0x70, 0x69, 0x6f, 0x6e, 0x4c, 0x6f, 0x62, 0x62, 0x79,
	0x2a, 0xc8, 0x01, 0x0a, 0x09, 0x4c, 0x6f, 0x62, 0x62, 0x79, 0x54, 0x79, 0x70, 0x65, 0x12, 0x11,
	0x0a, 0x0d, 0x4c, 0x4f, 0x42, 0x42, 0x59, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x4b, 0x55, 0x10,
	0x00, 0x12, 0x12, 0x0a, 0x0e, 0x4c, 0x4f, 0x42, 0x42, 0x59, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f,
	0x44, 0x41, 0x4e, 0x10, 0x01, 0x12, 0x17, 0x0a, 0x13, 0x4c, 0x4f, 0x42, 0x42, 0x59, 0x5f, 0x54,
	0x59, 0x50, 0x45, 0x5f, 0x55, 0x50, 0x50, 0x45, 0x52, 0x44, 0x41, 0x4e, 0x10, 0x02, 0x12, 0x16,
	0x0a, 0x12, 0x4c, 0x4f, 0x42, 0x42, 0x59, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x50, 0x48, 0x4f,
	0x45, 0x4e, 0x49, 0x58, 0x10, 0x03, 0x12, 0x11, 0x0a, 0x0d, 0x4c, 0x4f, 0x42, 0x42, 0x59, 0x5f,
	0x54, 0x59, 0x50, 0x45, 0x5f, 0x44, 0x5a, 0x10, 0x04, 0x12, 0x11, 0x0a, 0x0d, 0x4c, 0x4f, 0x42,
	0x42, 0x59, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x58, 0x31, 0x10, 0x05, 0x12, 0x11, 0x0a, 0x0d,
	0x4c, 0x4f, 0x42, 0x42, 0x59, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x58, 0x32, 0x10, 0x06, 0x12,
	0x11, 0x0a, 0x0d, 0x4c, 0x4f, 0x42, 0x42, 0x59, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x58, 0x33,
	0x10, 0x07, 0x12, 0x17, 0x0a, 0x13, 0x4c, 0x4f, 0x42, 0x42, 0x59, 0x5f, 0x54, 0x59, 0x50, 0x45,
	0x5f, 0x45, 0x58, 0x54, 0x45, 0x52, 0x4e, 0x41, 0x4c, 0x10, 0x08, 0x2a, 0x4e, 0x0a, 0x0a, 0x47,
	0x61, 0x6d, 0x65, 0x4c, 0x65, 0x6e, 0x67, 0x74, 0x68, 0x12, 0x15, 0x0a, 0x11, 0x47, 0x41, 0x4d,
	0x45, 0x5f, 0x4c, 0x45, 0x4e, 0x47, 0x54, 0x48, 0x5f, 0x53, 0x4f, 0x55, 0x54, 0x48, 0x10, 0x00,
	0x12, 0x14, 0x0a, 0x10, 0x47, 0x41, 0x4d, 0x45, 0x5f, 0x4c, 0x45, 0x4e, 0x47, 0x54, 0x48, 0x5f,
	0x45, 0x41, 0x53, 0x54, 0x10, 0x01, 0x12, 0x13, 0x0a, 0x0f, 0x47, 0x41, 0x4d, 0x45, 0x5f, 0x4c,
	0x45, 0x4e, 0x47, 0x54, 0x48, 0x5f, 0x4f, 0x4e, 0x45, 0x10, 0x02, 0x2a, 0x2c, 0x0a, 0x08, 0x47,
	0x61, 0x6d, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0f, 0x0a, 0x0b, 0x47, 0x41, 0x4d, 0x45, 0x5f,
	0x54, 0x59, 0x50, 0x45, 0x5f, 0x34, 0x10, 0x00, 0x12, 0x0f, 0x0a, 0x0b, 0x47, 0x41, 0x4d, 0x45,
	0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x33, 0x10, 0x01, 0x2a, 0x37, 0x0a, 0x09, 0x47, 0x61, 0x6d,
	0x65, 0x53, 0x70, 0x65, 0x65, 0x64, 0x12, 0x15, 0x0a, 0x11, 0x47, 0x41, 0x4d, 0x45, 0x5f, 0x53,
	0x50, 0x45, 0x45, 0x44, 0x5f, 0x4e, 0x4f, 0x52, 0x4d, 0x41, 0x4c, 0x10, 0x00, 0x12, 0x13, 0x0a,
	0x0f, 0x47, 0x41, 0x4d, 0x45, 0x5f, 0x53, 0x50, 0x45, 0x45, 0x44, 0x5f, 0x46, 0x41, 0x53, 0x54,
	0x10, 0x01, 0x2a, 0x24, 0x0a, 0x05, 0x41, 0x6b, 0x6b, 0x61, 0x73, 0x12, 0x0d, 0x0a, 0x09, 0x41,
	0x4b, 0x4b, 0x41, 0x53, 0x5f, 0x59, 0x45, 0x53, 0x10, 0x00, 0x12, 0x0c, 0x0a, 0x08, 0x41, 0x4b,
	0x4b, 0x41, 0x53, 0x5f, 0x4e, 0x4f, 0x10, 0x01, 0x2a, 0x27, 0x0a, 0x06, 0x54, 0x61, 0x6e, 0x79,
	0x61, 0x6f, 0x12, 0x0e, 0x0a, 0x0a, 0x54, 0x41, 0x4e, 0x59, 0x41, 0x4f, 0x5f, 0x59, 0x45, 0x53,
	0x10, 0x00, 0x12, 0x0d, 0x0a, 0x09, 0x54, 0x41, 0x4e, 0x59, 0x41, 0x4f, 0x5f, 0x4e, 0x4f, 0x10,
	0x01, 0x2a, 0x45, 0x0a, 0x0a, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x54, 0x79, 0x70, 0x65, 0x12,
	0x0d, 0x0a, 0x09, 0x4e, 0x4f, 0x5f, 0x4e, 0x55, 0x4d, 0x42, 0x45, 0x52, 0x10, 0x00, 0x12, 0x0c,
	0x0a, 0x08, 0x4e, 0x55, 0x4d, 0x42, 0x45, 0x52, 0x5f, 0x32, 0x10, 0x02, 0x12, 0x0c, 0x0a, 0x08,
	0x4e, 0x55, 0x4d, 0x42, 0x45, 0x52, 0x5f, 0x35, 0x10, 0x05, 0x12, 0x0c, 0x0a, 0x08, 0x4e, 0x55,
	0x4d, 0x42, 0x45, 0x52, 0x5f, 0x30, 0x10, 0x0a, 0x42, 0x2c, 0x5a, 0x2a, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x64, 0x6e, 0x6f, 0x76, 0x69, 0x6b, 0x6f, 0x66, 0x66,
	0x2f, 0x74, 0x65, 0x6e, 0x68, 0x6f, 0x75, 0x2f, 0x67, 0x65, 0x6e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2f, 0x73, 0x74, 0x61, 0x74, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_stats_stats_proto_rawDescOnce sync.Once
	file_stats_stats_proto_rawDescData = file_stats_stats_proto_rawDesc
)

func file_stats_stats_proto_rawDescGZIP() []byte {
	file_stats_stats_proto_rawDescOnce.Do(func() {
		file_stats_stats_proto_rawDescData = protoimpl.X.CompressGZIP(file_stats_stats_proto_rawDescData)
	})
	return file_stats_stats_proto_rawDescData
}

var file_stats_stats_proto_enumTypes = make([]protoimpl.EnumInfo, 7)
var file_stats_stats_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_stats_stats_proto_goTypes = []interface{}{
	(LobbyType)(0),                // 0: stats.LobbyType
	(GameLength)(0),               // 1: stats.GameLength
	(GameType)(0),                 // 2: stats.GameType
	(GameSpeed)(0),                // 3: stats.GameSpeed
	(Akkas)(0),                    // 4: stats.Akkas
	(Tanyao)(0),                   // 5: stats.Tanyao
	(NumberType)(0),               // 6: stats.NumberType
	(*Player)(nil),                // 7: stats.Player
	(*Record)(nil),                // 8: stats.Record
	(*timestamppb.Timestamp)(nil), // 9: google.protobuf.Timestamp
	(*durationpb.Duration)(nil),   // 10: google.protobuf.Duration
}
var file_stats_stats_proto_depIdxs = []int32{
	9,  // 0: stats.Record.time:type_name -> google.protobuf.Timestamp
	10, // 1: stats.Record.duration:type_name -> google.protobuf.Duration
	0,  // 2: stats.Record.lobby:type_name -> stats.LobbyType
	1,  // 3: stats.Record.length:type_name -> stats.GameLength
	2,  // 4: stats.Record.type:type_name -> stats.GameType
	4,  // 5: stats.Record.akkas:type_name -> stats.Akkas
	5,  // 6: stats.Record.tanyao:type_name -> stats.Tanyao
	6,  // 7: stats.Record.number_type:type_name -> stats.NumberType
	7,  // 8: stats.Record.players:type_name -> stats.Player
	9,  // [9:9] is the sub-list for method output_type
	9,  // [9:9] is the sub-list for method input_type
	9,  // [9:9] is the sub-list for extension type_name
	9,  // [9:9] is the sub-list for extension extendee
	0,  // [0:9] is the sub-list for field type_name
}

func init() { file_stats_stats_proto_init() }
func file_stats_stats_proto_init() {
	if File_stats_stats_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_stats_stats_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Player); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_stats_stats_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Record); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_stats_stats_proto_rawDesc,
			NumEnums:      7,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_stats_stats_proto_goTypes,
		DependencyIndexes: file_stats_stats_proto_depIdxs,
		EnumInfos:         file_stats_stats_proto_enumTypes,
		MessageInfos:      file_stats_stats_proto_msgTypes,
	}.Build()
	File_stats_stats_proto = out.File
	file_stats_stats_proto_rawDesc = nil
	file_stats_stats_proto_goTypes = nil
	file_stats_stats_proto_depIdxs = nil
}
