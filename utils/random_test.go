package utils

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRandomInt(t *testing.T) {
	result := RandomInt(1, 15)
	require.Equal(t, "int64", fmt.Sprintf("%T", result))
	require.True(t, result > 1 && result < 15)
}

func TestRandomString(t *testing.T) {
	result := RandomString(15)
	require.Equal(t, "string", fmt.Sprintf("%T", result))
	require.Len(t, result, 15)
}

func TestGenerateRandomOwner(t *testing.T) {
	result := GenerateRandomOwner()
	require.Equal(t, "string", fmt.Sprintf("%T", result))
	require.Len(t, result, 6)
}

func TestRandomCurrency(t *testing.T) {
	expectedCurrencies := []string{"EUR", "BRL", "USD"}

	result := RandomCurrency()

	require.Equal(t, "string", fmt.Sprintf("%T", result))
	require.Len(t, result, 3)

	require.Contains(t, expectedCurrencies, result)
}

func TestRandomMoney(t *testing.T) {
	result := RandomMoney()

	require.Equal(t, "int64", fmt.Sprintf("%T", result))
	require.True(t, result > 0 && result < 1000)
}
