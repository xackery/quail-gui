package form

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/xackery/wlk/walk"
)

var (
	validateTexture  = &textureValidator{}
	validateMaterial = &layerMaterialValidator{}
)

func strValidate(in string) error {
	if in == "" {
		return fmt.Errorf("value is empty")
	}
	return nil
}

func intValidate(in string, min int, max int) error {
	val, err := strconv.Atoi(in)
	if err != nil {
		return fmt.Errorf("is not a number")
	}

	if val < min {
		return fmt.Errorf("is less than minimum (%d)", min)
	}
	if val > max {
		return fmt.Errorf("is greater than maximum (%d)", max)
	}
	return nil
}

func floatValidate(in string) error {
	_, err := strconv.ParseFloat(in, 32)
	if err != nil {
		return fmt.Errorf("is not a number")
	}
	return nil
}

type layerMaterialValidator struct {
}

func (v *layerMaterialValidator) Create() (walk.Validator, error) {
	return v, nil
}

func (v *layerMaterialValidator) Validate(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("value is not a string")
	}
	err := strValidate(str)
	if err != nil {
		return err
	}
	// regex must follow pattern C_RKP_S00_M01, C_RKP_S00_M00, C_RKP_S00_M02

	reg := regexp.MustCompile(`_S\d\d_M\d\d$`)
	if !reg.MatchString(str) {
		return fmt.Errorf("value is not a valid material, needs _S##_M## pattern suffix")
	}

	return nil
}

type textureValidator struct {
}

func (v *textureValidator) Create() (walk.Validator, error) {
	return v, nil
}

func (v *textureValidator) Validate(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("value is not a string")
	}
	err := strValidate(str)
	if err != nil {
		return err
	}

	return nil
}
