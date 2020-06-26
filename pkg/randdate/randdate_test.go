package randdate

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"makeit.imfr.cgi.com/makeit2/scm/lino/pimo/pkg/model"
)

func TestMaskingShouldReplaceDateByRandom(t *testing.T) {
	dateMin := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	dateMax := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	dateRange := NewMask(dateMin, dateMax, 0)

	config := model.NewMaskConfiguration().
		WithEntry("date", dateRange)

	maskingEngine := model.MaskingEngineFactory(config)
	data := model.Dictionary{"date": time.Date(2019, 3, 2, 0, 0, 0, 0, time.UTC)}

	result, err := maskingEngine.Mask(data)
	assert.Equal(t, nil, err, "error should be nil")
	resultmap := result.(map[string]model.Entry)

	assert.Greater(t, dateMax.Sub(resultmap["date"].(time.Time)).Microseconds(), int64(0),
		"%v should be before max date %v", resultmap["date"].(time.Time), dateMax)

	assert.Less(t, dateMin.Sub(resultmap["date"].(time.Time)).Microseconds(), int64(0),
		"%v should be after min date %v", resultmap["date"].(time.Time), dateMin)

	equal := 0
	for i := 0; i < 1000; i++ {
		result, err := maskingEngine.Mask(data)
		assert.Equal(t, nil, err, "error should be nil")
		resultmap := result.(map[string]model.Entry)
		if resultmap["date"] == data["date"] {
			equal++
		}
	}

	assert.True(t, equal <= 1, "Shouldn't be the same date less than 0.1% of time")
}

func TestRegistryMaskToConfigurationShouldCreateAMask(t *testing.T) {
	min := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)
	max := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	maskingConfig := model.Masking{Mask: model.MaskType{RandDate: model.RandDateType{DateMin: min, DateMax: max}}}
	config, present, err := RegistryMaskToConfiguration(maskingConfig, model.NewMaskConfiguration(), 0)
	assert.NotNil(t, config, "config shouldn't be nil")
	assert.True(t, present, "should be true")
	assert.Nil(t, err, "error should be nil")
}

func TestRegistryMaskToConfigurationShouldNotCreateAMaskFromAnEmptyConfig(t *testing.T) {
	maskingConfig := model.Masking{Mask: model.MaskType{}}
	mask, present, err := RegistryMaskToConfiguration(maskingConfig, model.NewMaskConfiguration(), 0)
	assert.Nil(t, mask, "should be nil")
	assert.False(t, present, "should be false")
	assert.Nil(t, err, "error should be nil")
}