package jwtex

import (
	"fmt"
	"testing"
	"time"

	"github.com/go2s/o2s/util/timeutil"
	"github.com/stretchr/testify/assert"
)

func TestTime(t *testing.T) {
	w := timeutil.Now()
	fmt.Println(w)
	w = w.UTC()
	fmt.Println(w)

	w = w.Round(time.Second)
	fmt.Println(w)

	now := w
	fmt.Println(now.Unix())

	du := time.Hour * 2
	fmt.Println(du)

	expiresAt := now.Add(du).Unix()
	fmt.Println(expiresAt)

	du1 := time.Unix(expiresAt, 0).Sub(now)
	fmt.Println(du1)

	assert.Equal(t, du1, du)
}
