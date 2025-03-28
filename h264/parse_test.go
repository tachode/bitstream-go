package h264_test

import (
	"encoding/json"
	"testing"

	"github.com/tachode/bitstream-go/h264"
)

func TestParse_ValidSPS(t *testing.T) {
	// Example SPS NAL unit (not real data, replace with actual SPS NAL unit bytes)
	spsNalUnit := spsBytes[:]

	nal, err := h264.Parse(spsNalUnit)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if nal.NalUnitType != h264.NalUnitTypeSPS {
		t.Errorf("Expected NalUnitTypeSPS, got %v", nal.NalUnitType)
	}

	if _, ok := nal.Payload.(*h264.SeqParameterSet); !ok {
		t.Errorf("Expected payload to be of type SeqParameterSet")
	}

	nalJSON, err := json.MarshalIndent(nal, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal nal to JSON: %v", err)
	}

	t.Logf("NAL JSON:\n%s", nalJSON)
}

func TestParse_InvalidData(t *testing.T) {
	// Invalid NAL unit data
	invalidNalUnit := []byte{0x00, 0x00, 0x00}

	_, err := h264.Parse(invalidNalUnit)
	if err == nil {
		t.Fatalf("Expected error for invalid NAL unit, got nil")
	}
}

func TestParse_EmptyBuffer(t *testing.T) {
	// Empty buffer
	emptyBuffer := []byte{}

	_, err := h264.Parse(emptyBuffer)
	if err == nil {
		t.Fatalf("Expected error for empty buffer, got nil")
	}
}
