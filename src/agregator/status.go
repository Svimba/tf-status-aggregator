package agregator

// TFStatus struct
type TFStatus struct {
	PodName   string     `json:"podName"`
	Groups    []*TFGroup `json:"groups"`
	PlainText string     `json:"-"`
}

// TFGroup stuct to handling TF service groups e.g. Control, Config, Analytics
type TFGroup struct {
	Name     string     `json:"name"`
	Services []*Service `json:"services"`
}

// Service struct
type Service struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}
