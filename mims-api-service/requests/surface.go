package requests

type RefSurface struct {
	Type             string   `json:"type" validate:"nonzero"`
	Name             string   `json:"name" validate:"nonzero"`
	Drainage         float64  `json:"drainage" `
	LayerCoefficient float64  `json:"layer_coefficient" `
	SurfaceGroup     string   `json:"surface_group" validate:"nonzero"`
	A                float64  `json:"a" `
	B                float64  `json:"b" `
	C1               int      `json:"c1"`
	C2               int      `json:"c2" `
	CRT              *float64 `json:"crt"`
	RRF              *float64 `json:"rrf"`
}

type RefSurfacePointer struct {
	Type             string   `json:"type" validate:"nonzero"`
	Name             string   `json:"name" validate:"nonzero"`
	Drainage         *float64 `json:"drainage" `
	LayerCoefficient *float64 `json:"layer_coefficient" `
	SurfaceGroup     string   `json:"surface_group" validate:"nonzero"`
	A                *float64 `json:"a" `
	B                *float64 `json:"b" `
	C1               *int     `json:"c1" `
	C2               *int     `json:"c2" `
	CRT              *float64 `json:"crt"`
	RRF              *float64 `json:"rrf"`
}
