package controller

/*------------------------------------ 一层 -------------------------------------*/

type StatusParam struct {
	Status int `json:"status" query:"status"`
}

type UserIDParam struct {
	UserID string `json:"user_id" query:"user_id"`
}

type TypeParam struct {
	Type int `json:"type" query:"type"`
}

type SuffixParam struct {
	Suffix int `json:"suffix" query:"suffix"`
}

/*------------------------------------ 二层 -------------------------------------*/

type TypeSuffix struct {
	TypeParam
	SuffixParam
}
