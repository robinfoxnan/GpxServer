package main

import (
	"fmt"
	json "github.com/json-iterator/go"
	"github.com/mitchellh/mapstructure"
	"reflect"
	"strings"
)

type UserInfo struct {
	Id      int64  `json:"id"`
	Phone   string `json:"phone"`
	Name    string `json:"name"`
	Nick    string `json:"nick"`
	Email   string `json:"email"`
	Pwd     string `json:"pwd"`
	TempPwd string `json:"temp_pwd"`
	LastPt  string `json:"last_pt"`
	Region  string `json:"region"`
	Ipv4    string `json:"ipv4"`
	WxId    string `json:"wxid"`
	Age     int8   `json:"age"`
	Gender  int8   `json:"gender"`
}

type GroupMember struct {
	Id     int64  `json:"id"`
	Phone  string `json:"phone"`
	Name   string `json:"name"`
	Nick   string `json:"nick"`
	GNick  string `json:"gnick"`
	Region string `json:"region"`
	Ipv4   string `json:"ipv_4"`
	Age    int8   `json:"age"`
	Gender int8   `json:"gender"`
}

type FriendShip struct {
	Id       int64
	Phone    string
	Friend   string
	CreateTm int64
	Mask     string
}

type Friends struct {
	Phone   string
	Members []FriendShip
}

// 创建一个组，如果是id==0; 否则添加成员到组
type GroupUsers struct {
	GroupId int64    `json:"group_id"`
	Owner   string   `json:"phone"`
	List    []string `json:"list"`
}

type ShortMsg struct {
	MsgType string
	GId     string
	FromId  string
	ToId    string
	Tm      int64
	SendSeq int64
	SubType string
	Mesage  string
}

////////////////////////////////////////////////////////
func (member *GroupMember) GroupMemberToString() (string, error) {
	data, err := json.Marshal(member)
	if err != nil {
		fmt.Printf("序列化错误 err=%v\n", err)
		return "", err
	}
	str := string(data)
	//fmt.Println("序列化后: ", str)
	return str, err
}
func (user *UserInfo) UserInfoToString() (string, error) {
	data, err := json.Marshal(user)
	if err != nil {
		fmt.Printf("序列化错误 err=%v\n", err)
		return "", err
	}
	str := string(data)
	//fmt.Println("序列化后: ", str)
	return str, err
}
func UserFromString(strJson string) (*UserInfo, error) {
	user := UserInfo{}
	err := json.Unmarshal([]byte(strJson), &user)
	if err != nil {
		return &user, err
	}
	return &user, nil
}

// 这里与json的TAG不一样，直接使用名字小写作为键值
func UserFromMap(data map[string]string) (*UserInfo, error) {
	user := UserInfo{}
	err := mapstructure.Decode(data, &user)
	if err != nil {
		fmt.Println(err.Error())
	}

	return &user, nil
}

// 使用反射序列化
func (user *UserInfo) UserInfoToMap() map[string]interface{} {
	t := reflect.TypeOf(*user)
	v := reflect.ValueOf(*user)

	var data = make(map[string]interface{})

	for i := 0; i < t.NumField(); i++ {
		key := strings.ToLower(t.Field(i).Name)
		//value := InterfaceToString(v.Field(i).Interface())
		data[key] = v.Field(i).Interface()
	}
	return data

}

func GroupMemberFromString(strJson string) (*GroupMember, error) {
	mem := GroupMember{}
	err := json.Unmarshal([]byte(strJson), &mem)
	if err != nil {
		return &mem, err
	}
	return &mem, nil
}

// 由基本信息转为快速信息
func (user *UserInfo) UserInfoToMemberInfo() (*GroupMember, error) {
	if user == nil {
		return nil, nil
	}
	member := GroupMember{}
	member.Id = user.Id
	member.Phone = user.Phone
	member.Nick = user.Nick
	member.GNick = user.Nick
	member.Region = user.Region
	member.Ipv4 = user.Ipv4
	member.Gender = user.Gender
	member.Age = user.Age

	return &member, nil
}
