package syntax

import "testing"

func TestSyntaxQueryScratchResetClearsValues(t *testing.T) {
	var scratch syntaxQueryScratch
	buffer := scratch.getSlice(2)
	buffer[0] = &struct{}{}
	buffer[1] = &struct{}{}

	scratch.reset()

	for index, value := range scratch.buffers[0][:2] {
		if value != nil {
			t.Fatalf("scratch value at index %d was retained", index)
		}
	}
}

func TestSyntaxQueryScratchResetDiscardsOversizedBuffers(t *testing.T) {
	var scratch syntaxQueryScratch
	scratch.getSlice(maxQueryScratchCapacity + 1)

	scratch.reset()

	if scratch.buffers[0] != nil {
		t.Fatal("oversized buffer was retained")
	}
}

func TestSyntaxQueryScratchResetLimitsBufferCount(t *testing.T) {
	var scratch syntaxQueryScratch
	for range maxQueryScratchBuffers + 1 {
		scratch.getSlice(1)
	}

	scratch.reset()

	if len(scratch.buffers) != maxQueryScratchBuffers {
		t.Fatalf("got %d buffers, want %d", len(scratch.buffers), maxQueryScratchBuffers)
	}
	if scratch.next != 0 {
		t.Fatalf("got next index %d, want 0", scratch.next)
	}
	if cap(scratch.buffers) > maxQueryScratchBuffers &&
		scratch.buffers[:maxQueryScratchBuffers+1][maxQueryScratchBuffers] != nil {
		t.Fatal("discarded buffer slot was retained")
	}
}
