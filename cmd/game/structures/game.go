package structures

import "path/filepath"

const (
	// Base structure and keywords for handbook
	Av = "AvailableResearch"
	Co = "CompletedResearch"
	Cr = "Creative"
	Ma = "Materials"
	Su = "Survival"
)

var (
	BaseDir  = filepath.Join("64", "Default")
	Handbook = filepath.Join(BaseDir, "Handbook")
	Data     = filepath.Join(BaseDir, "Data")
	Lang     = filepath.Join(BaseDir, "Lang")
)
