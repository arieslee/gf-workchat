package internal

import "gf_workchat/core"

type UserInfo struct {
	core.ResponseError
	Userid string `json:"userid"`
	Name string `json:"name"`
	Department []int `json:"department"`
	Order []int `json:"order"`
	Position string `json:"position"`
	Mobile string `json:"mobile"`
	Gender string `json:"gender"`
	Email string `json:"email"`
	IsLeaderInDept []int `json:"is_leader_in_dept"`
	Avatar string `json:"avatar"`
	ThumbAvatar string `json:"thumb_avatar"`
	Telephone string `json:"telephone"`
	Alias string `json:"alias"`
	Address string `json:"address"`
	OpenUserid string `json:"open_userid"`
	MainDepartment int `json:"main_department"`
	ExtAttr ExtAttr `json:"ext_attr"`
	Status int `json:"status"`
	QrCode string `json:"qr_code"`
	ExternalPosition string `json:"external_position"`
	ExternalProfile ExternalProfile `json:"external_profile"`
}
type Text struct {
	Value string `json:"value"`
}
type Web struct {
	URL string `json:"url"`
	Title string `json:"title"`
}
type Attrs struct {
	Type int `json:"type"`
	Name string `json:"name"`
	Text Text `json:"text,omitempty"`
	Web Web `json:"web,omitempty"`
}
type ExtAttr struct {
	Attrs []Attrs `json:"attrs"`
}
type MiniProgram struct {
	AppId string `json:"app_id"`
	PagePath string `json:"page_path"`
	Title string `json:"title"`
}
type ExternalAttr struct {
	Type int `json:"type"`
	Name string `json:"name"`
	Text Text `json:"text,omitempty"`
	Web Web `json:"web,omitempty"`
	MiniProgram MiniProgram `json:"mini_program,omitempty"`
}
type ExternalProfile struct {
	ExternalCorpName string `json:"external_corp_name"`
	ExternalAttr []ExternalAttr `json:"external_attr"`
}
