package main

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestLabel(t *testing.T) {
	// Skip tests that actually hit the API in short mode. We run with the
	// -test.short flag in Travis from external pull requests because the
	// credentials are not available.
	if testing.Short() {
		return
	}
	assert := assert.New(t)
	creds := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	//const filename = "../data/label/1038471-1.jpg"
	const filename = "../data/land/photo-1449.jpeg"

	os.Clearenv()
	assert.Error(run(filename))

	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", creds)
	assert.Error(run("no_exists.jpg"))
	assert.NoError(run(filename))
}
