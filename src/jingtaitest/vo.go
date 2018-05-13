package main
type password struct {

}
type project struct {
	Name string `json:"name,omitempty"`
}
type scope struct {
	Project project `json:"project,omitempty"`
}

type identity struct {
	Methods string

}

type auth struct {
	Name           string `json:"name,omitempty"`
	Time           string `json:"time,omitempty"`
	FileWriteSpeed int    `json:"writespeed,omitempty"`
	FileReadSpeed  int    `json:"readspeed,omitempty"`
	Success        int    `json:"success,omitempty"`
	Total          int    `json:"total,omitempty"`
}