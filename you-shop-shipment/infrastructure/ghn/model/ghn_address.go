package model

type GhnProvinceAddress struct {
	ProvinceID    int      `json:"ProvinceID"`
	ProvinceName  string   `json:"ProvinceName"`
	CountryID     int      `json:"CountryID"`
	Code          string   `json:"Code"`
	NameExtension []string `json:"NameExtension"`
	IsEnable      int      `json:"IsEnable"`
	RegionID      int      `json:"RegionID"`
	CanUpdateCOD  bool     `json:"CanUpdateCOD"`
	Status        int      `json:"Status"`
}
type GhnDistrictAddress struct {
	DistrictID    int      `json:"DistrictID"`
	ProvinceID    int      `json:"ProvinceID"`
	DistrictName  string   `json:"DistrictName"`
	Code          string   `json:"Code"`
	Type          int      `json:"Type"`
	SupportType   int      `json:"SupportType"`
	NameExtension []string `json:"NameExtension"`
	IsEnable      int      `json:"IsEnable"`
	CanUpdateCOD  bool     `json:"CanUpdateCOD"`
	Status        int      `json:"Status"`
}

type GhnWardAddress struct {
	WardCode      string   `json:"WardCode"`
	DistrictID    int      `json:"DistrictID"`
	WardName      string   `json:"WardName"`
	NameExtension []string `json:"NameExtension"`
	CanUpdateCOD  bool     `json:"CanUpdateCOD"`
	SupportType   int      `json:"SupportType"`
	Status        int      `json:"Status"`
}
