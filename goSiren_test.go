package goSiren

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestGoSiren_IsSending_환경_에러코드_호스트_모두일치(t *testing.T) {
	t.Run("", func(t *testing.T) {
		siren := &GoSiren{
			TargetCodes:        []int{http.StatusBadRequest, http.StatusInternalServerError},
			TargetEnvironments: []string{EnvironmentStaging},
			TargetHost:         []string{"www.test.com"},
		}
		environment := EnvironmentStaging
		code := http.StatusBadRequest
		host := "www.test.com"
		// when
		got := siren.IsSending(code, environment, host)
		// then
		expected := true
		assert.Equal(t, expected, got)
	})
}
func TestGoSiren_IsSending_환경이다른경우(t *testing.T) {
	t.Run("", func(t *testing.T) {
		//given
		siren := &GoSiren{
			TargetCodes:        []int{http.StatusBadRequest, http.StatusInternalServerError},
			TargetEnvironments: []string{EnvironmentProduction},
			TargetHost:         []string{"www.test.com"},
		}
		environment := EnvironmentStaging
		code := http.StatusBadRequest
		host := "www.test.com"
		// when
		got := siren.IsSending(code, environment, host)
		// then
		expected := false
		assert.Equal(t, expected, got)
	})
}

func TestGoSiren_Setting_기본_Alarm_으로_내용생성(t *testing.T) {
	//given
	siren := &GoSiren{
		TargetCodes:        []int{http.StatusBadRequest, http.StatusInternalServerError},
		TargetEnvironments: []string{EnvironmentProduction},
		TargetHost:         []string{"www.test.com"},
	}

	alarm := Alarm{
		Request: &http.Request{
			Host:       "www.test.com",
			RequestURI: "/",
			Method:     http.MethodPost,
		},
		Code:       http.StatusInternalServerError,
		StackTrace: "error  test.go line 230",
		UserID:     1,
		UserName:   "test",
	}
	// when
	title, _, _ := siren.Setting(&alarm)
	assert.Equal(t, "error  test.go line 230", title)
}

/*
직접 struct 를 만들어서 원하는 탬플릿을 만들 수 있습니다.
*/
type CustomAlarm struct {
}

func (g *CustomAlarm) make() (string, string, error) {
	return "title", "content", nil
}

func TestGoSiren_Setting_커스텀_Alarm_으로_내용생성(t *testing.T) {
	//given
	siren := &GoSiren{
		TargetCodes:        []int{http.StatusBadRequest, http.StatusInternalServerError},
		TargetEnvironments: []string{EnvironmentProduction},
		TargetHost:         []string{"www.test.com"},
	}

	// when
	title, content, _ := siren.Setting(&CustomAlarm{})
	assert.Equal(t, "title", title)
	assert.Equal(t, "content", content)
}
