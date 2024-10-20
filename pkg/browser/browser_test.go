package browser

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPrintableOpenURL(t *testing.T) {
	input := "http://proxmox.example.com:8006"
	expected := "# Open this URL in the browser: http://proxmox.example.com:8006"
	output := PrintableOpenURL(input)
	assert.Equal(t, expected, output)
}
