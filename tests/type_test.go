package tests

import (
	"fmt"
	"sober-api/internal/helper"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testingCreateAccountReq(t *testing.T) {

	acc, err := helper.CreateUserAccount("testing", "testing@testing-email.com", "testing")

	assert.Nil(t, err)

	fmt.Println(acc)

}

func testingOnBoardingReq(t *testing.T) {

	boarding, err := helper.AddOnBoardingFlow(1, "I want to be sober", "2021-09-01")

	assert.Nil(t, err)

	fmt.Println(boarding)

}
